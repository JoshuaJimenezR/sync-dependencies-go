package main

import (
	"fmt"
	"os"
	"regexp"

	cli "github.com/jawher/mow.cli"
)

var GlobalSilentMode bool

func startCLI() {
	var (
		app = cli.App(appName, appDesc)
	)

	// Register global flags
	app.BoolOptPtr(&GlobalSilentMode, "s silent", false, "Silent mode, suppress 'job finished' notification")

	// *********************************************************************************************
	// ** Before
	// *********************************************************************************************

	// execute before commands
	app.Before = func() {

	}

	// Version
	app.Command(appActionVersion, appActionVersionDesc, func(cmd *cli.Cmd) {
		cmd.Action = func() {
			actionVersion()
		}
	})

	// Update Dependencies
	app.Command(appActionUpdateDeps, appActionUpdateDepsDesc, func(cmd *cli.Cmd) {
		cmd.Spec = appActionUpdateDepsSPEC
		var (
			argDirectoryPath = cmd.StringOpt("d directory", "", "Directory path")
			argModule        = cmd.StringOpt("m module", "", "Module name")
			argCommit        = cmd.StringOpt("c commit", "", "Commit ID")
		)

		cmd.Action = func() {
			directoryPath := *argDirectoryPath
			module := *argModule
			commit := *argCommit

			// Validate Directory Path
			fileInfo, err := os.Stat(directoryPath)
			if err != nil {
				if os.IsNotExist(err) {
					PrintErrorMessage(fmt.Sprintf("Error: Directory '--directory %s' not found.", directoryPath))
				} else {
					PrintErrorMessage(fmt.Sprintf("Error: Could not stat directory '--directory %s': %v", directoryPath, err))
				}
				cli.Exit(1)
			}
			if !fileInfo.IsDir() {
				PrintErrorMessage(fmt.Sprintf("Error: Path '--directory %s' is not a directory.", directoryPath))
				cli.Exit(1)
			}

			// If validations pass, call the action function
			UpdateDependencies(directoryPath, module, commit)

			if !GlobalSilentMode {
				if errNotify := CommandSayFinished(); errNotify != nil {
					PrintErrorMessage(fmt.Sprintf("Finished notification error: %v", errNotify))
				}
			}
		}
	})

	app.Command(appActionUpdateGoVersion, appActionUpdateGoVersionDesc, func(cmd *cli.Cmd) {
		cmd.Spec = appActionUpdateGoVersionSPEC
		var (
			argDirectoryPath = cmd.StringOpt("d directory", "", "Directory path")
			argOldVersion    = cmd.StringOpt("o old-version", "", "Old Go version")
			argNewVersion    = cmd.StringOpt("n new-version", "", "New Go version")
		)

		cmd.Action = func() {
			directoryPath := *argDirectoryPath
			oldVersion := *argOldVersion
			newVersion := *argNewVersion

			// Validate Directory Path
			fileInfo, err := os.Stat(directoryPath)
			if err != nil {
				if os.IsNotExist(err) {
					PrintErrorMessage(fmt.Sprintf("Error: Directory '--directory %s' not found.", directoryPath))
				} else {
					PrintErrorMessage(fmt.Sprintf("Error: Could not stat directory '--directory %s': %v", directoryPath, err))
				}
				cli.Exit(1)
			}
			if !fileInfo.IsDir() {
				PrintErrorMessage(fmt.Sprintf("Error: Path '--directory %s' is not a directory.", directoryPath))
				cli.Exit(1)
			}

			// Validate Go Version Strings
			versionRegex := regexp.MustCompile(`^\d+\.\d+(\.\d+)?$`)

			if !versionRegex.MatchString(oldVersion) {
				PrintErrorMessage(fmt.Sprintf("Error: Invalid Go version format for '--old-version %s'. Expected format like '1.18' or '1.19.2'.", oldVersion))
				cli.Exit(1)
			}

			if !versionRegex.MatchString(newVersion) {
				PrintErrorMessage(fmt.Sprintf("Error: Invalid Go version format for '--new-version %s'. Expected format like '1.18' or '1.19.2'.", newVersion))
				cli.Exit(1)
			}

			// If validations pass, call the action function
			UpdateGoVersion(directoryPath, oldVersion, newVersion)

			if !GlobalSilentMode {
				if errNotify := CommandSayFinished(); errNotify != nil {
					PrintErrorMessage(fmt.Sprintf("Finished notification error: %v", errNotify))
				}
			}
		}
	})
	// *********************************************************************************************
	// ** Run CLI
	// *********************************************************************************************

	// With the app configured, execute it, passing in the os.Args array
	app.Run(os.Args)
}
