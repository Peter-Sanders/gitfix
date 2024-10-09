package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func PrintHelp() {
	fmt.Println("Usage: gitfix [-t value] [-s value] [-h]")
	fmt.Println("Options:")
	fmt.Println("  -t value\tThe target branch. EX: dev ")
	fmt.Println("  -s value\tThe source feature branch. EX: DE-1234")
	fmt.Println("  -h\t\tPrint this help message")
}

func PickFiles(files []string, source string) []string {
	var newFiles []string

	fmt.Println("Pick only the files you want to include:")

	// Prompt for each file
	for _, file := range files {
		fmt.Printf("Include %s? (y/n): ", file)

		// Read user input
		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input: %v\n", err)
			continue
		}

		response = strings.TrimSpace(response)
		response = strings.ToLower(response)

		switch response {
		case "y":
			newFiles = append(newFiles, file)
		case "n":
			// Do nothing, file is not included
		default:
			fmt.Println("Invalid response. Skipping file.")
		}
	}

	finalFiles := CheckFiles(newFiles, source)

	return finalFiles
}

func CheckFiles(files []string, source string) []string {
	var newFiles []string

	fmt.Println("Files changed:")
	for _, fileName := range files {
		fmt.Println(fileName)
	}

	out := "CHECK THIS LIST OF FILES CAREFULLY!!!!!!\nType 'y' to proceed, 'q' to quit, 'p' to pick only the files you want included, 'v' to "
	out += "choose files you want to keep using vi, or 'r' to reset the files back to the original diff (y/q/p/v/r):"
	fmt.Println(out)
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Error reading input: %v\n", err)
	}

	response = strings.TrimSpace(response)
	response = strings.ToLower(response)

	switch response {
	case "y":
		newFiles = files
	case "r":
		fmt.Println("Resetting files to original diff")
		originalFiles, _, err := GitDiff(source, "")
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return newFiles
		}
		newFiles = CheckFiles(originalFiles, source)
	case "q":
		// Leave this empty so we can check for an empty response later
	case "p":
		newFiles = PickFiles(files, source)
	case "v":
		newFiles, err = SearchFiles(files, source)
		if err != nil {
			fmt.Printf("Error Searching for Files: %v\n", err)
		}
	default:
		fmt.Println("Invalid response. Type only 'y', 'r', 'q', 'p', or 'v'")
		newFiles = CheckFiles(files, source)
	}

	return newFiles
}

func WriteToFile(filename string, lines []string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		_, err := fmt.Fprintf(writer, "%s\n", line)
		if err != nil {
			return fmt.Errorf("error writing to file: %v", err)
		}
	}
	writer.Flush()

	return nil
}

func WriteHeader(filename string) error {
	header := []string{
		"\" GITFIX V1.0\"",
		"\"",
		"\" To mark a file to keep, append 'keep ' before its name.",
		"\" For example:",
		"\"",
		"\" keep filename.txt",
		"\"",
		"\" Make your selections below:",
		"\"",
		"\" Continue with the rest of the process by typing ':wq' to write and then quit the vi editor",
	}

	content, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	// Prepare the new content with header
	newContent := strings.Join(header, "\n") + "\n\n" + string(content)

	// Write the new content back to the file
	err = os.WriteFile(filename, []byte(newContent), 0644)
	if err != nil {
		return fmt.Errorf("error writing header to file: %v", err)
	}

	return nil
}

func OpenInVim(outputFile string) error {
	// Open file in Vi for editing or confirmation
	cmd := exec.Command("vi", outputFile)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error running Vi: %v", err)
	}

	return nil
}

func ReadAndFilterFile(filename string) ([]string, error) {
	var keptFiles []string

	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "keep ") {
			cleanedLine := strings.TrimPrefix(line, "keep ")
			keptFiles = append(keptFiles, cleanedLine)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning file: %v", err)
	}

	return keptFiles, nil
}

func SearchFiles(files []string, source string) ([]string, error) {
	var newFiles []string

	tempFile, err := os.CreateTemp("", "temp.txt")
	if err != nil {
		return nil, fmt.Errorf("error creating temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	err = WriteHeader(tempFile.Name())
	if err != nil {
		return nil, fmt.Errorf("error writing header to temp file: %v", err)
	}

	err = WriteToFile(tempFile.Name(), files)
	if err != nil {
		return nil, err
	}

	err = OpenInVim(tempFile.Name())
	if err != nil {
		return nil, err
	}

	newFiles, err = ReadAndFilterFile(tempFile.Name())
	if err != nil {
		return nil, err
	}

	finalFiles := CheckFiles(newFiles, source)

	return finalFiles, nil
}

func DeleteFiles(files []string) error {
	for file := range files {
		// Delete the file
		fmt.Printf("Deleting file: %s\n", files[file])
		cmd := exec.Command("rm", files[file])
		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("Error deleting file %s: %v\n", files[file], err)
		}
	}
	return nil
}

func ConfirmFiles(modifiedFiles []string, deletedFiles []string) bool {

	// Display all the modified files
	fmt.Printf("Files to be carried over:\n")
	for _, file := range modifiedFiles {
		fmt.Printf("%s\n ", file)
	}

	// Display all the deleted files
	fmt.Printf("Files to be deleted:\n")
	for _, file := range deletedFiles {
		fmt.Printf("%s\n ", file)
	}

	out := "CHECK THIS LIST OF FILES CAREFULLY!!!!!!\nType 'y' to proceed or 'n' to quit (y/n):"
	fmt.Println(out)
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Error reading input: %v\n", err)
	}

	response = strings.TrimSpace(response)
	response = strings.ToLower(response)

	switch response {
	case "y":
		return true
	case "n":
		return false
	default:
		fmt.Println("Invalid response. Type only 'y'or 'n'")
		ConfirmFiles(modifiedFiles, deletedFiles)
	}

	fmt.Println("Something weird happened you should never see this ")
	return false
}
