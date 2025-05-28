package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func UpdateGoVersion(defaultDirectory, oldVersion, newVersion string) {
	PrintMessage(fmt.Sprintf("Starting Go version update from '%s' to '%s'", oldVersion, newVersion))

	now := time.Now()
	var filesUpdated []string // Renamed for clarity and consistency

	// Assuming readFiles iterates and calls the callback.
	// Errors within the callback are handled by printing a message and returning from the callback,
	// allowing readFiles to proceed with the next directory.
	// readFileContent is assumed to be modified to return (string, error) and not call log.Fatalf.
	err := readFiles(defaultDirectory, func(directory string) {
		filePath := fmt.Sprintf("%s/%s", directory, goModFile) // Using goModFile constant
		PrintMessage(fmt.Sprintf("Checking for Go version update in: %s", filePath))

		// Assume readFileContent now returns (string, error)
		fileContent, errFile := readFileContent(filePath)
		if errFile != nil {
			PrintErrorMessage(fmt.Sprintf("Error reading %s: %v. Skipping this file.", filePath, errFile))
			return // Skip this file
		}

		lines := strings.Split(fileContent, "\n")
		foundAndReplaced := false
		goVersionLineOld := fmt.Sprintf("go %s", oldVersion)
		goVersionLineNew := fmt.Sprintf("go %s", newVersion)

		for i, line := range lines {
			// TrimSpace helps ensure direct comparison, e.g. "go 1.18" vs "go 1.18  "
			if strings.TrimSpace(line) == goVersionLineOld {
				lines[i] = goVersionLineNew
				foundAndReplaced = true
				// Assuming only one 'go' version directive per go.mod, so we can break.
				break
			}
		}

		if foundAndReplaced {
			output := strings.Join(lines, "\n")
			// Using 0644 as a more standard file permission
			if errWrite := os.WriteFile(filePath, []byte(output), 0644); errWrite != nil {
				PrintErrorMessage(fmt.Sprintf("Error writing updated content to %s: %v. File may not have been updated.", filePath, errWrite))
				return // Skip this file (or at least, don't add to filesUpdated)
			}
			filesUpdated = append(filesUpdated, filePath)
			PrintSuccessMessage(fmt.Sprintf("Successfully updated Go version in %s to %s", filePath, newVersion))
		} else {
			PrintMessage(fmt.Sprintf("Old Go version directive 'go %s' not found in %s. No changes made to this file.", oldVersion, filePath))
		}
	})

	// This error is from readFiles itself (e.g., initial directory scan failed)
	if err != nil {
		PrintErrorMessage(fmt.Sprintf("Error during directory scanning: %v. Go version update process may be incomplete.", err))
		return // Exit UpdateGoVersion if readFiles itself fails
	}

	PrintSuccessMessage(fmt.Sprintf("Go version update process from '%s' to '%s' finished in %s.", oldVersion, newVersion, time.Since(now)))
	if len(filesUpdated) > 0 {
		PrintSuccessMessage("Successfully updated Go version in the following files:")
		for _, file := range filesUpdated { // Corrected variable name here
			PrintMessage(file)
		}
	} else {
		PrintMessage(fmt.Sprintf("No files were updated from Go version '%s' to '%s'.", oldVersion, newVersion))
	}
}
