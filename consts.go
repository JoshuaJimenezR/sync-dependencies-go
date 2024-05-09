package main

const (
	appVersion = "0.1-30122023"
	appName    = "sync-dependencies"
	appDesc    = "Sync Dependencies CLI tool"
	appSPEC    = "[ -c=<config> ]"

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

	appActionUpdateGoVersion     = "update-go-version"
	appActionUpdateGoVersionDesc = "Updates the go version in all go.mod files"
	appActionUpdateGoVersionSPEC = ""
)
