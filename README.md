# sync-dependencies-go

## Description
This package provides CLI tools to help synchronize and update dependencies for Go projects. It operates by modifying `go.mod` files based on the provided commands and flags.

## Commands Available

All commands now use flags for argument input instead of interactive prompts.

---

### `update-deps`
Updates a specific dependency module to a given commit ID within the `go.mod` files found in the target directory (and its subdirectories).

**Flags:**
*   `-d, --directory=<path>`: (Required) Path to the target root directory to scan for `go.mod` files.
*   `-m, --module=<module_name>`: (Required) The full name of the module to update (e.g., `example.com/your/module`).
*   `-c, --commit=<commit_id>`: (Required) The commit ID, tag, or branch to update the module to.

**Example:**
```bash
sync-dependencies-go update-deps -d ./myproject -m example.com/mymodule -c abc123xyz
```

---

### `update-go-version`
Updates the Go version directive in `go.mod` files found within the target directory (and its subdirectories).

**Flags:**
*   `-d, --directory=<path>`: (Required) Path to the target root directory to scan for `go.mod` files.
*   `-o, --old-version=<version>`: (Required) The old Go version string to replace (e.g., "1.18").
*   `-n, --new-version=<version>`: (Required) The new Go version string to set (e.g., "1.19").

**Example:**
```bash
sync-dependencies-go update-go-version -d ./myproject -o 1.18 -n 1.19
```

---

### `version`
Displays the application version.

**Example:**
```bash
sync-dependencies-go version
```

---

## Global Flags

These flags can be used with any command.

*   `-s, --silent`: Suppresses the "job finished" voice notification if the system supports it.

**Example:**
```bash
sync-dependencies-go update-deps -d ./myproject -m example.com/mymodule -c abc123xyz -s
```
```bash
sync-dependencies-go update-go-version -d ./myproject -o 1.18 -n 1.19 --silent
```
