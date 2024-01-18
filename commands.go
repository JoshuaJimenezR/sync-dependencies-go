package main

import (
	"bytes"
	"fmt"
	"os/exec"
)

func CommandGoGet(directory string) error {
	cmd := exec.Command(commandBash, commandC, commandGoGet)
	cmd.Dir = directory

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		fmt.Println(fmt.Sprintf("error: %v - %v", err, stderr.String()))
		return err
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
		fmt.Println(fmt.Sprintf("error: %v - %v", err, stderr.String()))
		return err
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
		fmt.Println(fmt.Sprintf("error: %v - %v", err, stderr.String()))
		return err
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
		fmt.Println(fmt.Sprintf("error: %v - %v", err, stderr.String()))
		return err
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
		fmt.Println(fmt.Sprintf("error: %v - %v", err, stderr.String()))
		return err
	}

	return nil
}

func CommandGoCommit(directory, module, commit string) error {
	command := fmt.Sprintf("%v %v@%v", commandGoGet, module, commit)
	cmd := exec.Command(commandBash, commandC, command)
	cmd.Dir = directory

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		fmt.Println(fmt.Sprintf("error: %v - %v", err, stderr.String()))
		return err
	}

	return nil
}

func CommandSayFinished() error {
	command := fmt.Sprintf("say %v", "job has finished")
	cmd := exec.Command(commandBash, commandC, command)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		fmt.Println(fmt.Sprintf("error: %v - %v", err, stderr.String()))
		return err
	}

	return nil
}
