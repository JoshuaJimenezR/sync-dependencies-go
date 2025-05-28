package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// Backup of the original exec.Command
var origExecCommand = exec.Command

// mockExecCommand will be reassigned by tests to a function that returns a mock *exec.Cmd
var mockExecCommand func(command string, args ...string) *exec.Cmd

// TestMain sets up and tears down the mock for exec.Command
func TestMain(m *testing.M) {
	// Override exec.Command with our mock version
	exec.Command = func(name string, arg ...string) *exec.Cmd {
		// If mockExecCommand is not set, it means we are in TestHelperProcess itself or a real command is unexpected.
		// For simplicity in this setup, TestHelperProcess will use the original exec.Command.
		// Tests should always set mockExecCommand.
		if mockExecCommand == nil {
			// This path can be taken by TestHelperProcess when it re-executes os.Args[0]
			// Or if a test doesn't set up mockExecCommand (which would be a test setup error)
			return origExecCommand(name, arg...)
		}
		return mockExecCommand(name, arg...)
	}

	// Run all tests
	exitCode := m.Run()

	// Restore the original exec.Command
	exec.Command = origExecCommand
	os.Exit(exitCode)
}

// newMockCmd creates a mock *exec.Cmd that calls TestHelperProcess.
// TestHelperProcess will then simulate success or failure based on environment variables.
func newMockCmd(testName string, wantErr bool, expectedCmdParts ...string) *exec.Cmd {
	// When this command is run, it will re-execute the test binary itself.
	// We pass flags or env vars to tell TestHelperProcess what to do.
	cmd := origExecCommand(os.Args[0], "-test.run=TestHelperProcess") // Use origExecCommand to avoid recursion
	cmd.Env = append(os.Environ(), "GO_TEST_HELPER_PROCESS=1")
	if wantErr {
		cmd.Env = append(cmd.Env, "GO_WANT_HELPER_ERROR=1")
	}
	// Pass expected command parts for TestHelperProcess to potentially validate (optional)
	// For now, TestHelperProcess doesn't use this, but it's good for debugging.
	cmd.Env = append(cmd.Env, fmt.Sprintf("GO_TEST_EXPECTED_CMD=%s", strings.Join(expectedCmdParts, " ")))
	return cmd
}

// TestHelperProcess is not a real test but a helper sub-process.
// It's executed when newMockCmd runs os.Args[0].
func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_TEST_HELPER_PROCESS") != "1" {
		return // Not in helper mode
	}

	// Simulate command execution
	// expectedCmd := os.Getenv("GO_TEST_EXPECTED_CMD") // Can be used for more detailed validation
	// args := os.Args[3:] // Actual arguments passed to the "command" if needed

	if os.Getenv("GO_WANT_HELPER_ERROR") == "1" {
		fmt.Fprint(os.Stderr, "mock command error stderr") // Simulate stderr output
		os.Exit(1)                                 // Simulate failure
	}

	fmt.Fprint(os.Stdout, "mock command success stdout") // Simulate stdout output
	os.Exit(0)                                  // Simulate success
}

func TestCommandGoGet(t *testing.T) {
	tests := []struct {
		name                 string
		directory            string
		wantErr              bool
		expectedErrSubstring string
		expectedCommand      []string
	}{
		{
			name:            "success",
			directory:       "testdir_success_goget",
			wantErr:         false,
			expectedCommand: []string{"bash", "-c", commandGoGet},
		},
		{
			name:                 "failure",
			directory:            "testdir_failure_goget",
			wantErr:              true,
			expectedErrSubstring: "mock command error stderr",
			expectedCommand:      []string{"bash", "-c", commandGoGet},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock for this specific test run
			mockExecCommand = func(command string, args ...string) *exec.Cmd {
				fullCmd := append([]string{command}, args...)
				if len(fullCmd) != len(tt.expectedCommand) {
					t.Fatalf("mockExecCommand: expected %d command parts, got %d. Expected: %v, Got: %v", len(tt.expectedCommand), len(fullCmd), tt.expectedCommand, fullCmd)
				}
				for i, part := range tt.expectedCommand {
					if fullCmd[i] != part {
						t.Errorf("mockExecCommand: command part %d: expected '%s', got '%s'", i, part, fullCmd[i])
					}
				}
				return newMockCmd(t.Name(), tt.wantErr, tt.expectedCommand...)
			}

			err := CommandGoGet(tt.directory)

			if tt.wantErr {
				if err == nil {
					t.Errorf("CommandGoGet() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !strings.Contains(err.Error(), tt.expectedErrSubstring) {
					t.Errorf("CommandGoGet() error string = %q, expected to contain %q", err.Error(), tt.expectedErrSubstring)
				}
			} else {
				if err != nil {
					t.Errorf("CommandGoGet() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
			// Clean up mock after test if necessary, or rely on next test to set it.
			// For this structure, each test/subtest must set mockExecCommand.
		})
	}
}

func TestCommandGoGetU(t *testing.T) {
	tests := []struct {
		name                 string
		directory            string
		wantErr              bool
		expectedErrSubstring string
		expectedCommand      []string
	}{
		{
			name:            "success",
			directory:       "testdir_success_gogetu",
			wantErr:         false,
			expectedCommand: []string{"bash", "-c", commandGoGetU},
		},
		{
			name:                 "failure",
			directory:            "testdir_failure_gogetu",
			wantErr:              true,
			expectedErrSubstring: "mock command error stderr",
			expectedCommand:      []string{"bash", "-c", commandGoGetU},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockExecCommand = func(command string, args ...string) *exec.Cmd {
				fullCmd := append([]string{command}, args...)
				if len(fullCmd) != len(tt.expectedCommand) {
					t.Fatalf("mockExecCommand: expected %d command parts, got %d. Expected: %v, Got: %v", len(tt.expectedCommand), len(fullCmd), tt.expectedCommand, fullCmd)
				}
				for i, part := range tt.expectedCommand {
					if fullCmd[i] != part {
						t.Errorf("mockExecCommand: command part %d: expected '%s', got '%s'", i, part, fullCmd[i])
					}
				}
				return newMockCmd(t.Name(), tt.wantErr, tt.expectedCommand...)
			}

			err := CommandGoGetU(tt.directory)

			if tt.wantErr {
				if err == nil {
					t.Errorf("CommandGoGetU() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !strings.Contains(err.Error(), tt.expectedErrSubstring) {
					t.Errorf("CommandGoGetU() error string = %q, expected to contain %q", err.Error(), tt.expectedErrSubstring)
				}
			} else {
				if err != nil {
					t.Errorf("CommandGoGetU() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestCommandGoModTidy(t *testing.T) {
	tests := []struct {
		name                 string
		directory            string
		wantErr              bool
		expectedErrSubstring string
		expectedCommand      []string
	}{
		{
			name:            "success",
			directory:       "testdir_success_gomodtidy",
			wantErr:         false,
			expectedCommand: []string{"bash", "-c", commandGoModTidy},
		},
		{
			name:                 "failure",
			directory:            "testdir_failure_gomodtidy",
			wantErr:              true,
			expectedErrSubstring: "mock command error stderr",
			expectedCommand:      []string{"bash", "-c", commandGoModTidy},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockExecCommand = func(command string, args ...string) *exec.Cmd {
				fullCmd := append([]string{command}, args...)
				if len(fullCmd) != len(tt.expectedCommand) {
					t.Fatalf("mockExecCommand: expected %d command parts, got %d. Expected: %v, Got: %v", len(tt.expectedCommand), len(fullCmd), tt.expectedCommand, fullCmd)
				}
				for i, part := range tt.expectedCommand {
					if fullCmd[i] != part {
						t.Errorf("mockExecCommand: command part %d: expected '%s', got '%s'", i, part, fullCmd[i])
					}
				}
				return newMockCmd(t.Name(), tt.wantErr, tt.expectedCommand...)
			}

			err := CommandGoModTidy(tt.directory)

			if tt.wantErr {
				if err == nil {
					t.Errorf("CommandGoModTidy() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !strings.Contains(err.Error(), tt.expectedErrSubstring) {
					t.Errorf("CommandGoModTidy() error string = %q, expected to contain %q", err.Error(), tt.expectedErrSubstring)
				}
			} else {
				if err != nil {
					t.Errorf("CommandGoModTidy() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestCommandGoModVendor(t *testing.T) {
	tests := []struct {
		name                 string
		directory            string
		wantErr              bool
		expectedErrSubstring string
		expectedCommand      []string
	}{
		{
			name:            "success",
			directory:       "testdir_success_gomodvendor",
			wantErr:         false,
			expectedCommand: []string{"bash", "-c", commandGoMod}, // commandGoMod is used for vendor
		},
		{
			name:                 "failure",
			directory:            "testdir_failure_gomodvendor",
			wantErr:              true,
			expectedErrSubstring: "mock command error stderr",
			expectedCommand:      []string{"bash", "-c", commandGoMod}, // commandGoMod is used for vendor
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockExecCommand = func(command string, args ...string) *exec.Cmd {
				fullCmd := append([]string{command}, args...)
				if len(fullCmd) != len(tt.expectedCommand) {
					t.Fatalf("mockExecCommand: expected %d command parts, got %d. Expected: %v, Got: %v", len(tt.expectedCommand), len(fullCmd), tt.expectedCommand, fullCmd)
				}
				for i, part := range tt.expectedCommand {
					if fullCmd[i] != part {
						t.Errorf("mockExecCommand: command part %d: expected '%s', got '%s'", i, part, fullCmd[i])
					}
				}
				return newMockCmd(t.Name(), tt.wantErr, tt.expectedCommand...)
			}

			err := CommandGoModVendor(tt.directory)

			if tt.wantErr {
				if err == nil {
					t.Errorf("CommandGoModVendor() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !strings.Contains(err.Error(), tt.expectedErrSubstring) {
					t.Errorf("CommandGoModVendor() error string = %q, expected to contain %q", err.Error(), tt.expectedErrSubstring)
				}
			} else {
				if err != nil {
					t.Errorf("CommandGoModVendor() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestCommandGoCleanModCache(t *testing.T) {
	tests := []struct {
		name                 string
		directory            string
		wantErr              bool
		expectedErrSubstring string
		expectedCommand      []string
	}{
		{
			name:            "success",
			directory:       "testdir_success_gocleanmodcache",
			wantErr:         false,
			expectedCommand: []string{"bash", "-c", commandGoCleanModCache},
		},
		{
			name:                 "failure",
			directory:            "testdir_failure_gocleanmodcache",
			wantErr:              true,
			expectedErrSubstring: "mock command error stderr",
			expectedCommand:      []string{"bash", "-c", commandGoCleanModCache},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockExecCommand = func(command string, args ...string) *exec.Cmd {
				fullCmd := append([]string{command}, args...)
				if len(fullCmd) != len(tt.expectedCommand) {
					t.Fatalf("mockExecCommand: expected %d command parts, got %d. Expected: %v, Got: %v", len(tt.expectedCommand), len(fullCmd), tt.expectedCommand, fullCmd)
				}
				for i, part := range tt.expectedCommand {
					if fullCmd[i] != part {
						t.Errorf("mockExecCommand: command part %d: expected '%s', got '%s'", i, part, fullCmd[i])
					}
				}
				return newMockCmd(t.Name(), tt.wantErr, tt.expectedCommand...)
			}

			err := CommandGoCleanModCache(tt.directory)

			if tt.wantErr {
				if err == nil {
					t.Errorf("CommandGoCleanModCache() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !strings.Contains(err.Error(), tt.expectedErrSubstring) {
					t.Errorf("CommandGoCleanModCache() error string = %q, expected to contain %q", err.Error(), tt.expectedErrSubstring)
				}
			} else {
				if err != nil {
					t.Errorf("CommandGoCleanModCache() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestCommandGoCommit(t *testing.T) {
	tests := []struct {
		name                 string
		directory            string
		module               string
		commit               string
		wantErr              bool
		expectedErrSubstring string
		expectedCommandParts []string // Parts of the actual command string like "go get module@commit"
	}{
		{
			name:                 "success",
			directory:            "testdir_success_gocommit",
			module:               "example.com/mymodule",
			commit:               "v1.2.3",
			wantErr:              false,
			expectedCommandParts: []string{commandGoGet, "example.com/mymodule@v1.2.3"},
		},
		{
			name:                 "failure",
			directory:            "testdir_failure_gocommit",
			module:               "example.com/othermodule",
			commit:               "abcdef1",
			wantErr:              true,
			expectedErrSubstring: "mock command error stderr",
			expectedCommandParts: []string{commandGoGet, "example.com/othermodule@abcdef1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expectedFullCommand := []string{"bash", "-c", strings.Join(tt.expectedCommandParts, " ")}
			mockExecCommand = func(command string, args ...string) *exec.Cmd {
				fullCmd := append([]string{command}, args...)
				if len(fullCmd) != len(expectedFullCommand) {
					t.Fatalf("mockExecCommand: expected %d command parts, got %d. Expected: %v, Got: %v", len(expectedFullCommand), len(fullCmd), expectedFullCommand, fullCmd)
				}
				for i, part := range expectedFullCommand {
					if fullCmd[i] != part {
						// Special check for the combined command string part
						if i == 2 && strings.Contains(part, "@") { // This is the "go get module@commit" part
							if !strings.Contains(fullCmd[i], tt.module) || !strings.Contains(fullCmd[i], tt.commit) {
								t.Errorf("mockExecCommand: command part %d: expected to contain '%s' and '%s', got '%s'", i, tt.module, tt.commit, fullCmd[i])
							}
						} else if fullCmd[i] != part {
							t.Errorf("mockExecCommand: command part %d: expected '%s', got '%s'", i, part, fullCmd[i])
						}
					}
				}
				return newMockCmd(t.Name(), tt.wantErr, expectedFullCommand...)
			}

			err := CommandGoCommit(tt.directory, tt.module, tt.commit)

			if tt.wantErr {
				if err == nil {
					t.Errorf("CommandGoCommit() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !strings.Contains(err.Error(), tt.expectedErrSubstring) {
					t.Errorf("CommandGoCommit() error string = %q, expected to contain %q", err.Error(), tt.expectedErrSubstring)
				}
			} else {
				if err != nil {
					t.Errorf("CommandGoCommit() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestCommandSayFinished(t *testing.T) {
	tests := []struct {
		name                 string
		wantErr              bool
		expectedErrSubstring string
		expectedCommandParts []string
	}{
		{
			name:                 "success",
			wantErr:              false,
			expectedCommandParts: []string{"say", "job has finished"},
		},
		{
			name:                 "failure",
			wantErr:              true,
			expectedErrSubstring: "mock command error stderr",
			expectedCommandParts: []string{"say", "job has finished"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expectedFullCommand := []string{"bash", "-c", strings.Join(tt.expectedCommandParts, " ")}
			mockExecCommand = func(command string, args ...string) *exec.Cmd {
				fullCmd := append([]string{command}, args...)
				if len(fullCmd) != len(expectedFullCommand) {
					t.Fatalf("mockExecCommand: expected %d command parts, got %d. Expected: %v, Got: %v", len(expectedFullCommand), len(fullCmd), expectedFullCommand, fullCmd)
				}
				for i, part := range expectedFullCommand {
					if fullCmd[i] != part {
						t.Errorf("mockExecCommand: command part %d: expected '%s', got '%s'", i, part, fullCmd[i])
					}
				}
				return newMockCmd(t.Name(), tt.wantErr, expectedFullCommand...)
			}

			err := CommandSayFinished()

			if tt.wantErr {
				if err == nil {
					t.Errorf("CommandSayFinished() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !strings.Contains(err.Error(), tt.expectedErrSubstring) {
					t.Errorf("CommandSayFinished() error string = %q, expected to contain %q", err.Error(), tt.expectedErrSubstring)
				}
			} else {
				if err != nil {
					t.Errorf("CommandSayFinished() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}
