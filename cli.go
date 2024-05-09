package main

import (
	cli "github.com/jawher/mow.cli"
	"os"
)

func startCLI() {
	var (
		app = cli.App(appName, appDesc)
	)

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

		cmd.Action = func() {
			// Define arguments
			var (
				argDirectoryPath string
				argModule        string
				argCommit        string
				err              error
			)

			// Read Directory Path
			PrintQuestion("Please provide the directory Path: ")
			argDirectoryPath, err = readCliInput()
			if err != nil {
				PrintErrorMessage("Error reading directory path", err)
				cli.Exit(1)
			}

			// Read Module
			PrintQuestion("Please provide the module: ")
			argModule, err = readCliInput()
			if err != nil {
				PrintErrorMessage("Error reading module", err)
				cli.Exit(1)
			}

			// Read Commit
			PrintQuestion("Please provide the commit ID: ")
			argCommit, err = readCliInput()
			if err != nil {
				PrintErrorMessage("Error reading commit", err)
				cli.Exit(1)
			}

			// Check if all arguments are provided
			if argDirectoryPath == "" || argModule == "" || argCommit == "" {
				PrintErrorMessage("Error reading arguments", "Directory Path, Module and commit are required ")
				cli.Exit(1)
			}

			UpdateDependencies(argDirectoryPath, argModule, argCommit)
		}
	})

	app.Command(appActionUpdateGoVersion, appActionUpdateGoVersionDesc, func(cmd *cli.Cmd) {
		cmd.Spec = appActionUpdateGoVersionSPEC

		cmd.Action = func() {
			// Define arguments
			var (
				argDirectoryPath string
				argOldVersion    string
				argVersion       string
				err              error
			)

			// Read Directoy Path
			PrintQuestion("Please provide the directory path: ")
			argDirectoryPath, err = readCliInput()
			if err != nil {
				PrintErrorMessage("Error reading directory path", err)
			}

			// Read old version
			PrintQuestion("Please provide the old go version: ")
			argDirectoryPath, err = readCliInput()
			if err != nil {
				PrintErrorMessage("Error reading old go version", err)
			}

			// Read new Version
			PrintQuestion("Please provide the new go version: ")
			argDirectoryPath, err = readCliInput()
			if err != nil {
				PrintErrorMessage("Error reading new go version", err)
			}

			// Check if all arguments are provided
			if argDirectoryPath == "" || argVersion == "" || argOldVersion == "" {
				PrintErrorMessage("Error Reading arguments", "Directory Path, old version and new version are required")
				cli.Exit(1)
			}

			UpdateGoVersion(argDirectoryPath, argOldVersion, argVersion)
		}
	})
	// *********************************************************************************************
	// ** Run CLI
	// *********************************************************************************************

	// With the app configured, execute it, passing in the os.Args array
	app.Run(os.Args)
}
