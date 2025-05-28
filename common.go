package main

import (
	"io/fs"
	"os"
	"path/filepath"
)

// readFileContent reads the entire content of a file.
// It returns the file content as a string and a nil error on success.
// If an error occurs, it returns an empty string and the error.
// This function does not use log.Fatalf or any other logging that terminates the program.
func readFileContent(filePath string) (string, error) {
	contentBytes, err := os.ReadFile(filePath)
	if err != nil {
		return "", err // Propagate error
	}
	return string(contentBytes), nil
}

// readFiles walks through the rootDirectory recursively.
// For each directory found that contains a go.mod file, it calls the callback function
// with the path to that directory.
// If the callback returns an error, readFiles stops iteration and returns the error from the callback.
// If filepath.WalkDir itself encounters an error, readFiles returns that error.
// If the directory traversal completes without the callback returning any errors, readFiles returns nil.
func readFiles(rootDirectory string, callback func(directoryPath string) error) error {
	return filepath.WalkDir(rootDirectory, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			// This error is passed by filepath.WalkDir if it cannot read a directory or
			// if there was an error from a previous step in the walk.
			// Propagate this error to stop the walk.
			return err
		}

		// Check if the entry is a file named as per goModFile constant (e.g., "go.mod").
		// The constant goModFile is defined in consts.go and should be accessible.
		if !d.IsDir() && d.Name() == goModFile {
			// Get the parent directory of the go.mod file.
			dir := filepath.Dir(path)

			// Call the callback with the directory containing the go.mod file.
			cbErr := callback(dir)
			if cbErr != nil {
				// Callback signaled an error. Stop walking by returning this error.
				// filepath.WalkDir will terminate and return this cbErr.
				return cbErr
			}
		}
		return nil // Continue walking
	})
}
