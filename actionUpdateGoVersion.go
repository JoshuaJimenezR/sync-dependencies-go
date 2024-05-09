package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func UpdateGoVersion(defaultDirectory, oldVersion, newVersion string) {
	PrintMessage("Updating go version")

	now := time.Now()
	var files []string

	err := readFiles(defaultDirectory, func(directory string) {
		file := fmt.Sprintf("%s%s", directory, goModFile)

		fileContent, err := readFileContent(file)
		if err != nil {
			log.Fatalf("unable to read file %v", err)
		}

		lines := strings.Split(fileContent, "\n")

		for i, line := range lines {
			if strings.Contains(line, fmt.Sprintf("go %v", oldVersion)) {
				lines[i] = fmt.Sprintf("go %v", newVersion)
			}
		}

		output := strings.Join(lines, "\n")

		err = os.WriteFile(file, []byte(output), 064)
		if err != nil {
			log.Fatalln(err)
		}

		files = append(files, file)
	})

	if err != nil {
		PrintErrorMessage(err.Error())
	}

	PrintSuccessMessage("Updated Go version done in", time.Since(now))
	PrintSuccessMessage("Files Updated")

	for _, file := range files {
		PrintMessage(file)
	}
}
