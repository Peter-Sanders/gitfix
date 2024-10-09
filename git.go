package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
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

func ConfirmBranches(branches [3]string) (bool, error) {
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

func GitDiff(branch string, default_branch string) ([]string, []string, error) {
	if default_branch == "" {
		default_branch = "main"
	}

	var modifiedBranches []string
	var deletedBranches []string

	cmd := exec.Command("git", "diff", default_branch, branch, "--name-status")

	output, err := cmd.Output()
	if err != nil {
		return nil, nil, fmt.Errorf("error running git diff command: %v", err)
	}

	r1 := regexp.MustCompile(":")
	results := strings.Split(r1.ReplaceAllString(strings.TrimSpace(string(output)), ""), "\n")

	for file := range results {
		fileData := strings.Fields(results[file])
		action := strings.TrimSpace(fileData[0])
		name := strings.TrimSpace(fileData[1])

		if action == "D" {
			deletedBranches = append(deletedBranches, name)
		} else {
			modifiedBranches = append(modifiedBranches, name)
		}
	}

	return modifiedBranches, deletedBranches, nil
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

func MoveAndPull(branch string) error {
	// Swap to the branch and update it
	fmt.Printf("Moving to %s\n", branch)
	cmd1 := exec.Command("git", "checkout", branch)
	err1 := cmd1.Run()
	if err1 != nil {
		return fmt.Errorf("Error running git checkout command: %v\n", err1)
	}

	// Pull the latest
	cmd2 := exec.Command("git", "pull")
	err2 := cmd2.Run()
	if err2 != nil {
		return fmt.Errorf("Error running git pull command: %v\n", err2)
	}

	return nil
}
