package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
)

func CheckBranchExistence(branch string) (bool, error) {
	cmdGit := exec.Command("git", "branch", "--list", branch)
	outputGit, err := cmdGit.CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("error running git command: %v", err)
	}

	branchesCount := strings.Count(string(outputGit), "\n")
	return branchesCount == 1, nil
}

func ConfirmBranches(branches [2]string) (bool, error) {
	var wg sync.WaitGroup
	wg.Add(len(branches))
	results := make(chan bool, len(branches))

	for _, value := range branches {
		go func(branch string) {
			defer wg.Done()
			exists, err := CheckBranchExistence(branch)
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

func GitDiff(branch string) ([]string, error) {
	cmd := exec.Command("git", "diff", branch, "--name-only")

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error running git diff command: %v", err)
	}

	fileNames := strings.Split(strings.TrimSpace(string(output)), "\n")

	return fileNames, nil
}

func MakeNewFeatureBranch(envFeatureBranch string) error {
	exists, err := CheckBranchExistence(envFeatureBranch)
	if err != nil {
		return err
	}

	if exists {
		fmt.Printf("Branch '%s' already exists.\n", envFeatureBranch)
		cmdCheckout := exec.Command("git", "checkout", envFeatureBranch)
		cmdCheckout.Stdout = os.Stdout
		cmdCheckout.Stderr = os.Stderr
		if err := cmdCheckout.Run(); err != nil {
			return fmt.Errorf("error checking out branch: %v", err)
		}
	} else {
		cmdCheckout := exec.Command("git", "checkout", "-b", envFeatureBranch)
		cmdCheckout.Stdout = os.Stdout
		cmdCheckout.Stderr = os.Stderr
		if err := cmdCheckout.Run(); err != nil {
			return fmt.Errorf("error creating and checking out branch: %v", err)
		}
		fmt.Printf("Created and checked out branch '%s'.\n", envFeatureBranch)
	}

	return nil
}

func CheckoutFiles(files []string, source string) error {
	for _, fileName := range files {
		cmd := exec.Command("git", "checkout", source, fileName)

		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("error checking out files: %v", err)
		}
	}

	return nil
}
