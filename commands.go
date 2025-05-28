package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func CommandGoGet(directory string) error {
	cmd := exec.Command(commandBash, commandC, commandGoGet)
	cmd.Dir = directory

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error executing '%s': %w. Stderr: %s", strings.Join(cmd.Args, " "), err, stderr.String())
	}

	return nil
}

func CommandGoGetU(directory string) error {
	cmd := exec.Command(commandBash, commandC, commandGoGetU)
	cmd.Dir = directory

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error executing '%s': %w. Stderr: %s", strings.Join(cmd.Args, " "), err, stderr.String())
	}

	return nil
}

func CommandGoModTidy(directory string) error {
	cmd := exec.Command(commandBash, commandC, commandGoModTidy)
	cmd.Dir = directory

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error executing '%s': %w. Stderr: %s", strings.Join(cmd.Args, " "), err, stderr.String())
	}

	return nil
}

func CommandGoModVendor(directory string) error {
	cmd := exec.Command(commandBash, commandC, commandGoMod)
	cmd.Dir = directory

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error executing '%s': %w. Stderr: %s", strings.Join(cmd.Args, " "), err, stderr.String())
	}

	return nil
}

func CommandGoCleanModCache(directory string) error {
	cmd := exec.Command(commandBash, commandC, commandGoCleanModCache)
	cmd.Dir = directory

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error executing '%s': %w. Stderr: %s", strings.Join(cmd.Args, " "), err, stderr.String())
	}

	return nil
}

func CommandGoCommit(directory, module, commit string) error {
	commandString := fmt.Sprintf("%s %s@%s", commandGoGet, module, commit)
	cmd := exec.Command(commandBash, commandC, commandString)
	cmd.Dir = directory

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		// For this one, cmd.Args would be [bash -c "go get module@commit"], so using commandString is more direct.
		return fmt.Errorf("error executing '%s': %w. Stderr: %s", commandString, err, stderr.String())
	}

	return nil
}

func CommandSayFinished() error {
	commandString := fmt.Sprintf("say %s", "job has finished")
	cmd := exec.Command(commandBash, commandC, commandString)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		// Similar to CommandGoCommit, using commandString is better than cmd.Args
		return fmt.Errorf("error executing '%s': %w. Stderr: %s", commandString, err, stderr.String())
	}

	return nil
}
