package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
)

func main() {
	// Parse command-line flags
	t := flag.String("t", "", "Target branch to merge into")
	s := flag.String("s", "", "Source branch which has staged changes")
	d := flag.String("d", "main", "Default branch which the source branch was originally branched off of")
	h := flag.Bool("h", false, "Print help message")
	flag.Parse()

	if *h {
		PrintHelp()
		os.Exit(0)
	}

	if *t == "" && *s == "" {
		fmt.Println("GITFIX v1.1\n\nMade by Peter Sanders")
		return
	}

	if *t == "" || *s == "" {
		fmt.Println("Please specify a target and source branch using -t and -s")
		return
	}

	source := *s
	target := *t
	default_branch := *d
	feature := source + "-" + target

	fmt.Printf("Want to merge %s into %s via %s originally based on %s\n", source, target, feature, default_branch)

	// Fetch Origin
	fmt.Printf("Fetching Origin\n")
	if err := exec.Command("git", "fetch", "origin").Run(); err != nil {
		fmt.Printf("Error running git fetch origin command: %v\n", err)
		return
	}

	// Confirm that the source, default, and target branches exist
	branchesExist, err := ConfirmBranches([3]string{source, target, default_branch})
	if err != nil {
		fmt.Printf("Error confirming branches: %v\n", err)
		return
	}
	if !branchesExist {
		fmt.Println("Not all branches exist")
		return
	}

	// Swap to the default branch and update it
	if err := MoveAndPull(default_branch); err != nil {
		fmt.Printf("Error moving to branch %s: %v\n", default_branch, err)
		return
	}

	// Get all the files in the diff
	modifiedFiles, deletedFiles, err := GitDiff(source, default_branch)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	proceed := ConfirmFiles(modifiedFiles, deletedFiles)
	if !proceed {
		return
	}

	if len(modifiedFiles)+len(deletedFiles) == 0 {
		fmt.Printf("There are no differences between %s and %s, exiting\n", default_branch, source)
		return
	}

	// Swap to the target branch and update it
	if err := MoveAndPull(target); err != nil {
		fmt.Printf("Error moving to branch %s: %v\n", target, err)
		return
	}

	// Create the new feature branch
	err = MakeNewFeatureBranch(feature)
	if err != nil {
		fmt.Printf("Couldn't make the new branch, exiting: %v\n", err)
		return
	}

	err = CheckoutFiles(modifiedFiles, source)
	if err != nil {
		fmt.Printf("Couldn't checkout files, exiting: %v\n", err)
		return
	}

	err = DeleteFiles(deletedFiles)
	if err != nil {
		fmt.Printf("Couldn't delete files, exiting %v\n", err)
		return
	}

	finalPrompt := "gitfix completed.\nPlease check each file you copied over for completeness and to confirm "
	finalPrompt += "that there are no changes brought in by other commits to the same file. This script is not "
	finalPrompt += "equipped to deal with such things and you'll need to deal with them manually"

	fmt.Println(finalPrompt)
}
