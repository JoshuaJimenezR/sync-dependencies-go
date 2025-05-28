package main

import (
	"fmt"
	"strings"
	"time"
)

func UpdateDependencies(defaultDirectory, module, commit string) {
	PrintMessage(fmt.Sprintf("Starting update process for module '%s' to commit '%s'", module, commit))

	now := time.Now()
	var updatedDirectories []string

	// Assuming readFiles iterates and calls the callback.
	// Errors within the callback are handled by printing a message and returning from the callback,
	// allowing readFiles to proceed with the next directory.
	// readFileContent is assumed to be modified to return (string, error) and not call log.Fatalf.
	err := readFiles(defaultDirectory, func(directory string) {
		filePath := fmt.Sprintf("%s/%s", directory, goModFile)
		PrintMessage(fmt.Sprintf("Checking directory: %s", directory))

		fileContent, errFile := readFileContent(filePath)
		if errFile != nil {
			PrintErrorMessage(fmt.Sprintf("Error reading %s: %v. Skipping this directory.", filePath, errFile))
			return // Skip this directory
		}

		// Original logic for checking module presence:
		// title := strings.Replace(fileContent, fmt.Sprintf("module %v", module), "", -1)
		// if strings.Contains(title, module)
		// This logic is preserved. If `module` is the module *declaring* itself in go.mod, this might be problematic.
		// If `module` is a dependency name, then `strings.Contains(fileContent, module)` would be more direct.
		// For now, adhering to the pattern of the original code.
		moduleDeclaration := fmt.Sprintf("module %s", module) // This is if `module` is the *name* of the go.mod itself.
		titleRemovedContent := strings.Replace(fileContent, moduleDeclaration, "", -1)

		if strings.Contains(titleRemovedContent, module) { // Check if the dependency 'module' is mentioned elsewhere.
			PrintMessage(fmt.Sprintf("Found module %s in %s. Proceeding with update.", module, filePath))

			// Run go get commit
			PrintMessage(fmt.Sprintf("Attempting 'go get %s@%s' in %s", module, commit, directory))
			if errCmd := CommandGoCommit(directory, module, commit); errCmd != nil {
				PrintErrorMessage(fmt.Sprintf("Failed to update module %s to commit %s in directory %s: %v. Skipping further processing for this directory.", module, commit, directory, errCmd))
				return // Skip this directory
			}
			PrintSuccessMessage(fmt.Sprintf("Successfully ran 'go get %s@%s' in %s", module, commit, directory))

			// Run go mod tidy
			PrintMessage(fmt.Sprintf("Attempting 'go mod tidy' in %s", directory))
			if errCmd := CommandGoModTidy(directory); errCmd != nil {
				PrintErrorMessage(fmt.Sprintf("Failed to run 'go mod tidy' in directory %s: %v. Dependency update may be incomplete here.", directory, errCmd))
				return // Skip this directory, as go mod tidy failed
			}
			PrintSuccessMessage(fmt.Sprintf("Successfully ran 'go mod tidy' in %s", directory))

			// Add directory to list if all operations succeeded for it
			updatedDirectories = append(updatedDirectories, directory)
		} else {
			PrintMessage(fmt.Sprintf("Module %s not found (as a dependency) in %s. Skipping.", module, filePath))
		}
	})

	// This error is from readFiles itself (e.g., initial directory scan failed or callback propagated error)
	if err != nil {
		PrintErrorMessage(fmt.Sprintf("Error during directory scanning or processing: %v. The update process may be incomplete.", err))
		// Depending on desired behavior, we might return here or let it report processed directories.
		// For now, let it proceed to print summary.
	}

	PrintSuccessMessage(fmt.Sprintf("Dependency update process for module '%s' finished in %s.", module, time.Since(now)))
	if len(updatedDirectories) > 0 {
		PrintSuccessMessage("Successfully updated dependencies in the following directories:")
		for _, dir := range updatedDirectories {
			PrintMessage(dir)
		}
	} else {
		PrintMessage("No directories were updated for this module and commit.")
	}
}
