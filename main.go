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
	h := flag.Bool("h", false, "Print help message")
	flag.Parse()

	if *h {
		PrintHelp()
		os.Exit(0)
	}

	if *t == "" && *s == "" {
		fmt.Println("GITFIX v1.0\n\nMade by Peter Sanders")
		return
	}

	if *t == "" || *s == "" {
		fmt.Println("Please specify a target and source branch using -t and -s")
		return
	}

	source := *s
	target := *t
	feature := source + "-" + target

	fmt.Printf("Want to merge %s into %s via %s\n", source, target, feature)

	// Confirm that the source and target branches exist
	branchesExist, err := ConfirmBranches([2]string{source, target})
	if err != nil {
		fmt.Printf("Error confirming branches: %v\n", err)
		return
	}
	if !branchesExist {
		fmt.Println("Not all branches exist")
		return
	}

	// Swap to the target branch and update it
	fmt.Printf("Moving to %s\n", target)
	if err := exec.Command("git", "checkout", target).Run(); err != nil {
		fmt.Printf("Error running git pull command: %v\n", err)
		return
	}

	if err := exec.Command("git", "pull").Run(); err != nil {
		fmt.Printf("Error running git pull command: %v\n", err)
		return
	}

	// Get all the files in the diff
	files, err := GitDiff(source)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if len(files) == 0 {
		fmt.Printf("There are no differences between %s and %s, exiting\n", target, source)
		return
	}

	// Get the list of files we want to checkout for the new branch
	good_files := CheckFiles(files, source)

	if len(good_files) == 0 {
		return
	}

	// Create the new feature branch
	err = MakeNewFeatureBranch(feature)
	if err != nil {
		fmt.Printf("Couldn't make the new branch, exiting: %v\n", err)
		return
	}

	err = CheckoutFiles(good_files, source)
	if err != nil {
		fmt.Printf("Couldn't checkout files, exiting: %v\n", err)
		return
	}

	finalPrompt := "Gitfix completed. Please check each file you copied over for completeness and to confirm "
	finalPrompt += "that there are no changes brought in by other commits to the same file. This script is not "
	finalPrompt += "equipped to deal with such things and you'll need to deal with them manually"

	fmt.Println(finalPrompt)
}
