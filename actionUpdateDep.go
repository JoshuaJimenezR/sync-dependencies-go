package main

import (
	"fmt"
	"log"
	"strings"
	"time"
)

func UpdateDependencies(defaultDirectory, module, commit string) {
	PrintMessage("Updating dependencies")

	now := time.Now()
	var directories []string

	err := readFiles(defaultDirectory, func(directory string) {
		fileContent, err := readFileContent(fmt.Sprintf("%s/%s", directory, goModFile))
		if err != nil {
			log.Fatalf("unable to read file: %v", err)
		}

		// remove title from file content
		title := strings.Replace(fileContent, fmt.Sprintf("module %v", module), "", -1)

		// Check if module exists in file
		if strings.Contains(title, module) {
			// Print directory
			fmt.Println(directory)

			// Run go get commit
			PrintMessage("Running go get commit")
			if err = CommandGoCommit(directory, module, commit); err != nil {
				PrintErrorMessage(err.Error())
			}

			// Run go mod tidy
			PrintMessage("Running go mod tidy")
			if err = CommandGoModTidy(directory); err != nil {
				PrintErrorMessage(err.Error())
			}

			// Add directory to list
			directories = append(directories, directory)
		}
	})
	if err != nil {
		PrintErrorMessage(err.Error())
	}

	PrintSuccessMessage("Updating dependencies done in", time.Since(now))
	PrintSuccessMessage("Directories Updated:")
	for _, dir := range directories {
		PrintMessage(dir)
	}
}
