package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	cli "github.com/jawher/mow.cli"
)

// --- Mocking infrastructure for action functions ---
var (
	mockUpdateDependenciesCalled bool
	mockUpdateDependenciesDir    string
	mockUpdateDependenciesModule string
	mockUpdateDependenciesCommit string
	realUpdateDependencies       = UpdateDependencies // Store real function
)

func mockUpdateDependenciesImpl(dir, module, commit string) {
	mockUpdateDependenciesCalled = true
	mockUpdateDependenciesDir = dir
	mockUpdateDependenciesModule = module
	mockUpdateDependenciesCommit = commit
}

var (
	mockUpdateGoVersionCalled bool
	mockUpdateGoVersionDir    string
	mockUpdateGoVersionOld    string
	mockUpdateGoVersionNew    string
	realUpdateGoVersion       = UpdateGoVersion // Store real function
)

func mockUpdateGoVersionImpl(dir, oldVersion, newVersion string) {
	mockUpdateGoVersionCalled = true
	mockUpdateGoVersionDir = dir
	mockUpdateGoVersionOld = oldVersion
	mockUpdateGoVersionNew = newVersion
}

var (
	mockActionVersionCalled bool
	realActionVersion       = actionVersion // Store real function
)

func mockActionVersionImpl() {
	mockActionVersionCalled = true
	// Simulate original actionVersion behavior for output testing if needed
	// For now, just record call. The actual actionVersion prints to stdout.
	// To test output, actionVersion would need to take a writer or use a global one we can swap.
	// Or, runTestCLI can capture it. For `version`, output is key.
	// Let's assume original actionVersion prints `appName + " version " + appVersion`
	fmt.Println(appName + " version " + appVersion)
}

var (
	mockCommandSayFinishedCalled bool
	realCommandSayFinished       = CommandSayFinished // Store real function
)

func mockCommandSayFinishedImpl() error {
	mockCommandSayFinishedCalled = true
	return nil
}

// Helper to set up all mocks
func setupActionMocks() {
	UpdateDependencies = mockUpdateDependenciesImpl
	UpdateGoVersion = mockUpdateGoVersionImpl
	actionVersion = mockActionVersionImpl
	CommandSayFinished = mockCommandSayFinishedImpl

	resetMockCallTrackers()
}

// Helper to reset mock trackers
func resetMockCallTrackers() {
	mockUpdateDependenciesCalled = false
	mockUpdateDependenciesDir = ""
	mockUpdateDependenciesModule = ""
	mockUpdateDependenciesCommit = ""

	mockUpdateGoVersionCalled = false
	mockUpdateGoVersionDir = ""
	mockUpdateGoVersionOld = ""
	mockUpdateGoVersionNew = ""

	mockActionVersionCalled = false
	mockCommandSayFinishedCalled = false
}

// Helper to tear down all mocks
func teardownActionMocks() {
	UpdateDependencies = realUpdateDependencies
	UpdateGoVersion = realUpdateGoVersion
	actionVersion = realActionVersion
	CommandSayFinished = realCommandSayFinished
}

// --- runTestCLI Helper ---
func runTestCLI(args []string, t *testing.T) (stdout, stderr string, exitCode int) {
	origArgs := os.Args
	origExitFunc := cli.ExitFunc
	origStdout := os.Stdout
	origStderr := os.Stderr

	// Defer restoration of original values
	defer func() {
		os.Args = origArgs
		cli.ExitFunc = origExitFunc
		os.Stdout = origStdout
		os.Stderr = origStderr
		GlobalSilentMode = false // Reset global flag
	}()

	// Set up args for this run
	os.Args = append([]string{"testapp"}, args...)

	// Mock cli.Exit
	var exitCodeStore int
	cli.ExitFunc = func(code int) {
		exitCodeStore = code
		panic(fmt.Sprintf("exit %d", code)) // Use panic to halt execution
	}

	// Capture stdout and stderr
	stdoutBuf := new(bytes.Buffer)
	stderrBuf := new(bytes.Buffer)
	os.Stdout = stdoutBuf
	os.Stderr = stderrBuf

	// Recover from panic caused by mock cli.ExitFunc
	defer func() {
		if r := recover(); r != nil {
			// Ensure it's our specific panic
			if !strings.HasPrefix(fmt.Sprint(r), "exit ") {
				panic(r) // Re-panic if it's not the one we expect
			}
		}
		stdout = stdoutBuf.String()
		stderr = stderrBuf.String()
		exitCode = exitCodeStore
	}()
	
	// Reset global flags that might be set by mow.cli parsing
	GlobalSilentMode = false 

	startCLI() // Execute the CLI logic

	// If startCLI completes without panic (i.e. no cli.Exit was called),
	// capture current buffer contents and set exitCode to 0.
	// This handles cases where commands succeed without explicit exit.
	if exitCodeStore == 0 && exitCode == 0 { // Check if panic didn't set exitCodeStore through cli.Exit
		stdout = stdoutBuf.String()
		stderr = stderrBuf.String()
		// exitCode remains 0
	}
	
	return // Named return values are assigned
}

// --- Test Cases ---

func TestCliVersionCommand(t *testing.T) {
	setupActionMocks()
	defer teardownActionMocks()

	stdout, stderr, exitCode := runTestCLI([]string{"version"}, t)

	if exitCode != 0 {
		t.Errorf("Expected exit code 0 for 'version' command, got %d. Stderr: %s", exitCode, stderr)
	}
	if !mockActionVersionCalled {
		t.Errorf("Expected actionVersion to be called, but it wasn't.")
	}
	expectedVersionString := appName + " version " + appVersion
	if !strings.Contains(stdout, expectedVersionString) {
		t.Errorf("Expected stdout to contain '%s', got '%s'", expectedVersionString, stdout)
	}
}

func TestCliUpdateDepsCommand(t *testing.T) {
	setupActionMocks()
	defer teardownActionMocks()

	t.Run("success", func(t *testing.T) {
		resetMockCallTrackers()
		tempDir := t.TempDir()
		// Create the directory to pass validation
		err := os.MkdirAll(tempDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}

		args := []string{"update-deps", "-d", tempDir, "-m", "example.com/mod", "-c", "abc123xyz"}
		_, stderr, exitCode := runTestCLI(args, t)

		if exitCode != 0 {
			t.Errorf("Expected exit code 0, got %d. Stderr: %s", exitCode, stderr)
		}
		if !mockUpdateDependenciesCalled {
			t.Errorf("Expected UpdateDependencies to be called, but it wasn't.")
		}
		if mockUpdateDependenciesDir != tempDir {
			t.Errorf("UpdateDependencies called with wrong dir: got %s, want %s", mockUpdateDependenciesDir, tempDir)
		}
		if mockUpdateDependenciesModule != "example.com/mod" {
			t.Errorf("UpdateDependencies called with wrong module: got %s, want example.com/mod", mockUpdateDependenciesModule)
		}
		if mockUpdateDependenciesCommit != "abc123xyz" {
			t.Errorf("UpdateDependencies called with wrong commit: got %s, want abc123xyz", mockUpdateDependenciesCommit)
		}
	})

	t.Run("missing required flags", func(t *testing.T) {
		resetMockCallTrackers()
		// Missing -m and -c
		args := []string{"update-deps", "-d", t.TempDir()} 
		// mow.cli will call cli.Exit(1) if spec implies required and not provided.
		// Our appActionUpdateDepsSPEC is "-d --directory=<directory> -m --module=<module> -c --commit=<commit>"
		// So all are required by mow.cli itself.
		_, stderr, exitCode := runTestCLI(args, t)

		if exitCode == 0 { // mow.cli should exit with non-zero
			t.Errorf("Expected non-zero exit code for missing flags, got %d", exitCode)
		}
		if mockUpdateDependenciesCalled {
			t.Errorf("UpdateDependencies should not be called when flags are missing.")
		}
		// Check stderr for message from mow.cli about missing flags
		// This message can be specific to mow.cli's internals, e.g. "Missing required option"
		// For robustness, just check that stderr is not empty.
		if stderr == "" {
			t.Errorf("Expected stderr to contain missing flags message, but it was empty.")
		}
		// A more specific check might be:
		// if !strings.Contains(stderr, "Missing required") && !strings.Contains(stderr, "expected") {
		//  t.Errorf("Expected stderr to contain message about missing flags, got: %s", stderr)
		// }
	})
	
	t.Run("invalid directory", func(t *testing.T) {
		resetMockCallTrackers()
		nonExistentDir := filepath.Join(t.TempDir(), "nonexistent")
		args := []string{"update-deps", "-d", nonExistentDir, "-m", "mod", "-c", "commit"}
		_, stderr, exitCode := runTestCLI(args, t)

		if exitCode != 1 { // Our validation exits with 1
			t.Errorf("Expected exit code 1 for invalid directory, got %d", exitCode)
		}
		if mockUpdateDependenciesCalled {
			t.Errorf("UpdateDependencies should not be called with invalid directory.")
		}
		if !strings.Contains(stderr, "not found") && !strings.Contains(stderr, "Could not stat directory") {
			t.Errorf("Expected stderr to contain 'not found' or 'Could not stat' message, got: %s", stderr)
		}
	})
}

func TestCliUpdateGoVersionCommand(t *testing.T) {
	setupActionMocks()
	defer teardownActionMocks()

	t.Run("success", func(t *testing.T) {
		resetMockCallTrackers()
		tempDir := t.TempDir()
		// Create the directory to pass validation
		err := os.MkdirAll(tempDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}

		args := []string{"update-go-version", "-d", tempDir, "-o", "1.18", "-n", "1.19"}
		_, stderr, exitCode := runTestCLI(args, t)

		if exitCode != 0 {
			t.Errorf("Expected exit code 0, got %d. Stderr: %s", exitCode, stderr)
		}
		if !mockUpdateGoVersionCalled {
			t.Errorf("Expected UpdateGoVersion to be called, but it wasn't.")
		}
		if mockUpdateGoVersionDir != tempDir {
			t.Errorf("UpdateGoVersion called with wrong dir: got %s, want %s", mockUpdateGoVersionDir, tempDir)
		}
		if mockUpdateGoVersionOld != "1.18" {
			t.Errorf("UpdateGoVersion called with wrong old version: got %s, want 1.18", mockUpdateGoVersionOld)
		}
		if mockUpdateGoVersionNew != "1.19" {
			t.Errorf("UpdateGoVersion called with wrong new version: got %s, want 1.19", mockUpdateGoVersionNew)
		}
	})

	t.Run("missing required flags", func(t *testing.T) {
		resetMockCallTrackers()
		// Missing -o and -n
		args := []string{"update-go-version", "-d", t.TempDir()}
		_, stderr, exitCode := runTestCLI(args, t)

		if exitCode == 0 { // mow.cli should exit
			t.Errorf("Expected non-zero exit code for missing flags, got %d", exitCode)
		}
		if mockUpdateGoVersionCalled {
			t.Errorf("UpdateGoVersion should not be called when flags are missing.")
		}
		if stderr == "" {
			t.Errorf("Expected stderr to contain missing flags message, but it was empty.")
		}
	})

	t.Run("invalid directory", func(t *testing.T) {
		resetMockCallTrackers()
		nonExistentDir := filepath.Join(t.TempDir(), "nonexistent")
		args := []string{"update-go-version", "-d", nonExistentDir, "-o", "1.18", "-n", "1.19"}
		_, stderr, exitCode := runTestCLI(args, t)

		if exitCode != 1 {
			t.Errorf("Expected exit code 1 for invalid directory, got %d", exitCode)
		}
		if mockUpdateGoVersionCalled {
			t.Errorf("UpdateGoVersion should not be called with invalid directory.")
		}
		if !strings.Contains(stderr, "not found") && !strings.Contains(stderr, "Could not stat directory") {
			t.Errorf("Expected stderr to contain 'not found' or 'Could not stat' message, got: %s", stderr)
		}
	})

	t.Run("invalid old Go version format", func(t *testing.T) {
		resetMockCallTrackers()
		tempDir := t.TempDir()
		err := os.MkdirAll(tempDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		args := []string{"update-go-version", "-d", tempDir, "-o", "1.bad", "-n", "1.19"}
		_, stderr, exitCode := runTestCLI(args, t)

		if exitCode != 1 {
			t.Errorf("Expected exit code 1 for invalid Go version format, got %d", exitCode)
		}
		if mockUpdateGoVersionCalled {
			t.Errorf("UpdateGoVersion should not be called with invalid Go version format.")
		}
		if !strings.Contains(stderr, "Invalid Go version format for '--old-version") {
			t.Errorf("Expected stderr to contain 'Invalid Go version format for --old-version' message, got: %s", stderr)
		}
	})

	t.Run("invalid new Go version format", func(t *testing.T) {
		resetMockCallTrackers()
		tempDir := t.TempDir()
		err := os.MkdirAll(tempDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		args := []string{"update-go-version", "-d", tempDir, "-o", "1.18", "-n", "1.bad.2"}
		_, stderr, exitCode := runTestCLI(args, t)

		if exitCode != 1 {
			t.Errorf("Expected exit code 1 for invalid Go version format, got %d", exitCode)
		}
		if mockUpdateGoVersionCalled {
			t.Errorf("UpdateGoVersion should not be called with invalid Go version format.")
		}
		if !strings.Contains(stderr, "Invalid Go version format for '--new-version") {
			t.Errorf("Expected stderr to contain 'Invalid Go version format for --new-version' message, got: %s", stderr)
		}
	})
}

func TestCliSilentFlag(t *testing.T) {
	setupActionMocks()
	defer teardownActionMocks()

	tempDir := t.TempDir() // Used by both sub-commands
	err := os.MkdirAll(tempDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create temp dir for silent flag tests: %v", err)
	}

	// --- Test with update-deps ---
	t.Run("update-deps without silent flag", func(t *testing.T) {
		resetMockCallTrackers()
		GlobalSilentMode = false // Explicitly set for clarity, though runTestCLI also resets
		args := []string{"update-deps", "-d", tempDir, "-m", "mod", "-c", "commit"}
		_, stderr, exitCode := runTestCLI(args, t)

		if exitCode != 0 {
			t.Errorf("Expected exit code 0, got %d. Stderr: %s", exitCode, stderr)
		}
		if !mockUpdateDependenciesCalled {
			t.Error("UpdateDependencies should have been called.")
		}
		if !mockCommandSayFinishedCalled {
			t.Error("CommandSayFinished should have been called when --silent is not used.")
		}
	})

	t.Run("update-deps with -s flag", func(t *testing.T) {
		resetMockCallTrackers()
		// GlobalSilentMode will be set by mow.cli parsing within startCLI via runTestCLI
		args := []string{"update-deps", "-s", "-d", tempDir, "-m", "mod", "-c", "commit"}
		_, stderr, exitCode := runTestCLI(args, t)

		if exitCode != 0 {
			t.Errorf("Expected exit code 0, got %d. Stderr: %s", exitCode, stderr)
		}
		if !mockUpdateDependenciesCalled {
			t.Error("UpdateDependencies should have been called.")
		}
		if mockCommandSayFinishedCalled {
			t.Error("CommandSayFinished should NOT have been called when -s is used.")
		}
	})
	
	t.Run("update-deps with --silent flag", func(t *testing.T) {
		resetMockCallTrackers()
		args := []string{"update-deps", "--silent", "-d", tempDir, "-m", "mod", "-c", "commit"}
		_, stderr, exitCode := runTestCLI(args, t)

		if exitCode != 0 {
			t.Errorf("Expected exit code 0, got %d. Stderr: %s", exitCode, stderr)
		}
		if !mockUpdateDependenciesCalled {
			t.Error("UpdateDependencies should have been called.")
		}
		if mockCommandSayFinishedCalled {
			t.Error("CommandSayFinished should NOT have been called when --silent is used.")
		}
	})

	// --- Test with update-go-version ---
	t.Run("update-go-version without silent flag", func(t *testing.T) {
		resetMockCallTrackers()
		GlobalSilentMode = false
		args := []string{"update-go-version", "-d", tempDir, "-o", "1.18", "-n", "1.19"}
		_, stderr, exitCode := runTestCLI(args, t)

		if exitCode != 0 {
			t.Errorf("Expected exit code 0, got %d. Stderr: %s", exitCode, stderr)
		}
		if !mockUpdateGoVersionCalled {
			t.Error("UpdateGoVersion should have been called.")
		}
		if !mockCommandSayFinishedCalled {
			t.Error("CommandSayFinished should have been called when --silent is not used.")
		}
	})

	t.Run("update-go-version with -s flag", func(t *testing.T) {
		resetMockCallTrackers()
		args := []string{"update-go-version", "-s", "-d", tempDir, "-o", "1.18", "-n", "1.19"}
		_, stderr, exitCode := runTestCLI(args, t)

		if exitCode != 0 {
			t.Errorf("Expected exit code 0, got %d. Stderr: %s", exitCode, stderr)
		}
		if !mockUpdateGoVersionCalled {
			t.Error("UpdateGoVersion should have been called.")
		}
		if mockCommandSayFinishedCalled {
			t.Error("CommandSayFinished should NOT have been called when -s is used.")
		}
	})

	t.Run("update-go-version with --silent flag", func(t *testing.T) {
		resetMockCallTrackers()
		args := []string{"update-go-version", "--silent", "-d", tempDir, "-o", "1.18", "-n", "1.19"}
		_, stderr, exitCode := runTestCLI(args, t)

		if exitCode != 0 {
			t.Errorf("Expected exit code 0, got %d. Stderr: %s", exitCode, stderr)
		}
		if !mockUpdateGoVersionCalled {
			t.Error("UpdateGoVersion should have been called.")
		}
		if mockCommandSayFinishedCalled {
			t.Error("CommandSayFinished should NOT have been called when --silent is used.")
		}
	})
}
```
