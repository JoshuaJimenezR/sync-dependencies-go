package main

import (
	"bufio"
	"log"
	"os"
)

func readFileContent(file string) (string, error) {
	body, err := os.ReadFile(file)
	if err != nil {
		log.Fatalf("unable to read file: %v", err)
		return "", err
	}

	return string(body), nil
}

func readCliInput() (string, error) {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	err := scanner.Err()
	if err != nil {
		return "", err
	}
	return scanner.Text(), nil
}
