package main

import (
	"flag"
	"fmt"
	"os/exec"
	"strings"
	"sync"
)

func checkBranchExistence(branch string) (bool, error) {
	cmdGit := exec.Command("git", "branch", "--list", branch)
	outputGit, err := cmdGit.CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("error running git command: %v", err)
	}

	branchesCount := strings.Count(string(outputGit), "\n")
	return branchesCount == 1, nil
}

func confirmBranches(branches [2]string) (bool, error) {
	var wg sync.WaitGroup
	wg.Add(len(branches))
	results := make(chan bool, len(branches))

	for _, value := range branches {
		go func(branch string) {
			defer wg.Done()
			exists, err := checkBranchExistence(branch)
			if err != nil {
				fmt.Printf("Error checking branch existence for %s: %v\n", branch, err)
				return
			}
			results <- exists
		}(value)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for res := range results {
		if !res {
			return false, nil
		}
	}
	return true, nil
}

func gitDiff(branch string) ([]string, error) {
	cmd := exec.Command("git", "diff", branch, "--name-only")

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error running git diff command: %v", err)
	}

	fileNames := strings.Split(strings.TrimSpace(string(output)), "\n")

	return fileNames, nil
}

func main() {
	// Parse command-line flags
	t := flag.String("t", "", "Target branch to merge into")
	s := flag.String("s", "", "Source branch which has staged changes")
	flag.Parse()

	if *t == "" || *s == "" {
		fmt.Println("Both -t and -s flags are required")
		return
	}

	source := *s
	target := *t
	feature := source + "-" + target

	fmt.Printf("Want to merge %s into %s via %s\n", source, target, feature)

	// Confirm that the source and target branches exist
	branchesExist, err := confirmBranches([2]string{source, target})
	if err != nil {
		fmt.Printf("Error confirming branches: %v\n", err)
		return
	}
	if !branchesExist {
		fmt.Println("Not all branches exist")
		return
	}

	// Swap to the target branch and update it
	if err := exec.Command("git", "checkout", target).Run(); err != nil {
		fmt.Printf("Error running git pull command: %v\n", err)
		return
	}

	if err := exec.Command("git", "pull").Run(); err != nil {
		fmt.Printf("Error running git pull command: %v\n", err)
		return
	}

	fileNames, err := gitDiff(source)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Files changed in '%s' compared to '%s':\n", source, target)
	for _, fileName := range fileNames {
		fmt.Println(fileName)
	}

}
