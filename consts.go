package main

const (
	appVersion = "0.1-30122023"
	appName    = "sync-dependencies"
	appDesc    = "Sync Dependencies CLI tool"
	appSPEC    = "[ -c=<config> ]"

	appArgDirectory        = "Source Directory"
	appArgDirectoryDefault = ""
	appArgDirectoryDesc    = "Source directory path"

	vendorFolder           = "vendor"
	goModFile              = "go.mod"
	commandGoGet           = "go get"
	commandGoModTidy       = "go mod tidy"
	commandGoMod           = "go mod vendor"
	commandGoCleanModCache = "go clean -modcache"
	commandGoGetU          = "go get -u"
	commandBash            = "bash"
	commandC               = "-c"

	appActionVersion     = "version"
	appActionVersionDesc = "Displays app version"

	appActionUpdateDeps     = "update-deps"
	appActionUpdateDepsDesc = "Finds all modules and updates dependencies based on commit id"
	appActionUpdateDepsSPEC = ""

	appActionUpdateDepsArgumentModule        = "MODULE"
	appActionUpdateDepsArgumentModuleDefault = ""
	appActionUpdateDepsArgumentModuleDesc    = "Module to update"

	appActionUpdateDepsArgumentCommit        = "COMMIT"
	appActionUpdateDepsArgumentCommitDefault = ""
	appActionUpdateDepsArgumentCommitDesc    = "Commit to update"
)
