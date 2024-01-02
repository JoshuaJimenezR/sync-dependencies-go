package main

import (
	"fmt"
	"os/exec"
)

func CommandGoGet(directory string) error {
	cmd := exec.Command(commandBash, commandC, commandGoGet)
	cmd.Dir = directory

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func CommandGoGetU(directory string) error {
	cmd := exec.Command(commandBash, commandC, commandGoGetU)
	cmd.Dir = directory

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func CommandGoModTidy(directory string) error {
	cmd := exec.Command(commandBash, commandC, commandGoModTidy)
	cmd.Dir = directory

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func CommandGoModVendor(directory string) error {
	cmd := exec.Command(commandBash, commandC, commandGoMod)
	cmd.Dir = directory

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func CommandGoCleanModCache(directory string) error {
	cmd := exec.Command(commandBash, commandC, commandGoCleanModCache)
	cmd.Dir = directory

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func CommandGoCommit(directory, module, commit string) error {
	command := fmt.Sprintf("%v %v@%v", commandGoGet, module, commit)
	cmd := exec.Command(commandBash, commandC, command)
	cmd.Dir = directory

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func CommandSayFinished() error {
	command := fmt.Sprintf("say %v", "job has finished")
	cmd := exec.Command(commandBash, commandC, command)

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
