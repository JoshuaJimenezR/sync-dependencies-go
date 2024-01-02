package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type baseFunc func(string)

func readFiles(defaultDirectory string, fn baseFunc) error {

	err := filepath.Walk(defaultDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip vendor folders
		if !strings.Contains(path, vendorFolder) {

			// Check if path contains go.mod
			if strings.Contains(path, goModFile) {
				// Get directory
				directory := fmt.Sprintf("..%s", strings.Trim(path, goModFile))

				fn(directory)
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
