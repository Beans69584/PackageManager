# PackageManager API Reference

**Temporary Note:**

- **This document is a placeholder and will be rewritten**

Welcome to the **PackageManager** API Reference. This document provides comprehensive details about the various packages, types, functions, and commands that constitute the PackageManager application. The PackageManager is a simple tool designed to install, uninstall, and manage software packages seamlessly.

---

## Table of Contents

- [PackageManager API Reference](#packagemanager-api-reference)
  - [Table of Contents](#table-of-contents)
  - [Overview](#overview)
  - [Packages](#packages)
    - [cmd Package](#cmd-package)
      - [Commands](#commands)
        - [`InstallCmd`](#installcmd)
        - [`ListCmd`](#listcmd)
        - [`UninstallCmd`](#uninstallcmd)
      - [Helper Functions](#helper-functions)
        - [`findExecutablesRecursively`](#findexecutablesrecursively)
        - [`isExecutable`](#isexecutable)
    - [pkg Package](#pkg-package)
      - [Types](#types)
        - [`Package`](#package)
        - [`PackageManager`](#packagemanager)
      - [Functions](#functions)
        - [`NewPackageManager`](#newpackagemanager)
        - [`Save`](#save)
        - [`AddPackage`](#addpackage)
        - [`RemovePackage`](#removepackage)
        - [`CreateDesktopFile`](#createdesktopfile)
        - [`RemoveDesktopFile`](#removedesktopfile)
        - [`ExtractTarGz`](#extracttargz)
        - [`getDefaultIcon`](#getdefaulticon)
      - [Detailed Function Descriptions](#detailed-function-descriptions)
        - [`ExtractTarGz`](#extracttargz-1)
    - [main Package](#main-package)
      - [`main`](#main)
  - [Directory Structure](#directory-structure)
    - [Breakdown](#breakdown)
  - [Additional Notes](#additional-notes)

---

## Overview

The **PackageManager** application is structured into several packages, each responsible for specific functionalities:

- **cmd**: Contains the command definitions (`install`, `list`, `uninstall`) using the Cobra library to create a CLI interface.
- **pkg**: Houses utility functions and structures for managing package metadata, handling desktop integrations, and extracting package archives.
- **main**: Serves as the entry point of the application, orchestrating the CLI and integrating all commands.

---

## Packages

### cmd Package

The `cmd` package defines the command-line interface commands for the PackageManager application. It utilises the [Cobra](https://github.com/spf13/cobra) library to create robust CLI commands.

#### Commands

##### `InstallCmd`

**Definition:**

```go
var InstallCmd = &cobra.Command{
    Use:   "install [archive.tar.gz]",
    Short: "Install a package from a tar.gz archive",
    Args:  cobra.ExactArgs(1),
    Run:   func(cmd *cobra.Command, args []string) { /* Implementation */ },
}
```

**Description:**

The `install` command allows users to install a software package from a `.tar.gz` archive. It performs the following actions:

1. **Archive Verification:** Ensures the specified archive exists.
2. **Extraction:** Extracts the archive to a designated installation directory.
3. **Executable Handling:** Identifies executables within the package and creates symbolic links for easy access.
4. **Desktop Integration:** Generates a `.desktop` file for the application.
5. **Package Tracking:** Adds the package to the PackageManager's tracking system.
6. **System Refresh:** Attempts to terminate the AGS bus to refresh desktop entries.

**Usage:**

```sh
packagemanager install /path/to/package.tar.gz
```

##### `ListCmd`

**Definition:**

```go
var ListCmd = &cobra.Command{
    Use:   "list",
    Short: "List all installed packages",
    Run:   func(cmd *cobra.Command, args []string) { /* Implementation */ },
}
```

**Description:**

The `list` command displays all currently installed packages managed by the PackageManager. It presents the package name, installation path, and executable path in a neatly formatted table.

**Usage:**

```sh
packagemanager list
```

##### `UninstallCmd`

**Definition:**

```go
var UninstallCmd = &cobra.Command{
    Use:   "uninstall [package_name]",
    Short: "Uninstall a package by name",
    Args:  cobra.ExactArgs(1),
    Run:   func(cmd *cobra.Command, args []string) { /* Implementation */ },
}
```

**Description:**

The `uninstall` command removes an installed package based on its name. The process includes:

1. **Package Identification:** Locates the package in the tracking system.
2. **Symlink Removal:** Deletes the symbolic link associated with the package's executable.
3. **Desktop Integration Cleanup:** Removes the `.desktop` file.
4. **Directory Cleanup:** Deletes the installation directory and its contents.
5. **Package Tracking Update:** Removes the package from the PackageManager's tracking system.
6. **System Refresh:** Attempts to terminate the AGS bus to refresh desktop entries.

**Usage:**

```sh
packagemanager uninstall SamplePackage
```

#### Helper Functions

##### `findExecutablesRecursively`

**Definition:**

```go
func findExecutablesRecursively(root string) ([]string, error)
```

**Description:**

Searches for executable files within a specified directory and its subdirectories. It returns a slice of paths to executable files found.

**Parameters:**

- `root` (`string`): The root directory from which the search begins.

**Returns:**

- `[]string`: A slice containing the paths of found executable files.
- `error`: An error object if the search fails, otherwise `nil`.

**Usage:**

```go
executables, err := findExecutablesRecursively("/path/to/install")
```

##### `isExecutable`

**Definition:**

```go
func isExecutable(path string, info os.FileInfo) bool
```

**Description:**

Determines if a given file is executable based on its permissions and, on Windows, its file extension.

**Parameters:**

- `path` (`string`): The file path to check.
- `info` (`os.FileInfo`): The file information object containing metadata about the file.

**Returns:**

- `bool`: `true` if the file is executable, `false` otherwise.

**Usage:**

```go
executable := isExecutable("/path/to/file", fileInfo)
```

---

### pkg Package

The `pkg` package encapsulates utility functions and structures essential for managing package metadata, desktop integrations, and archive extraction within the PackageManager application.

#### Types

##### `Package`

**Definition:**

```go
type Package struct {
    UUID        string `json:"uuid"`
    Name        string `json:"name"`
    InstallPath string `json:"install_path"`
    Executable  string `json:"executable"`
}
```

**Description:**

Represents an installed package with its essential metadata. This structure is used to track and manage packages within the PackageManager.

**Fields:**

- `UUID` (`string`): A unique identifier for the package installation.
- `Name` (`string`): The user-friendly name of the package.
- `InstallPath` (`string`): The filesystem path where the package is installed.
- `Executable` (`string`): The path to the package's main executable file.

##### `PackageManager`

**Definition:**

```go
type PackageManager struct {
    PackagesFile string
    Packages     []Package
}
```

**Description:**

Manages the collection of installed packages. It handles loading from and saving to the packages database file.

**Fields:**

- `PackagesFile` (`string`): The path to the JSON file that stores package metadata.
- `Packages` (`[]Package`): A slice containing all the currently installed packages.

#### Functions

##### `NewPackageManager`

**Definition:**

```go
func NewPackageManager(packagesFile string) (*PackageManager, error)
```

**Description:**

Creates and initialises a new `PackageManager`. It loads existing packages from the specified packages file or creates a new one if it doesn't exist.

**Parameters:**

- `packagesFile` (`string`): The path to the JSON file that stores package metadata.

**Returns:**

- `*PackageManager`: A pointer to the initialised `PackageManager`.
- `error`: An error object if initialisation fails, otherwise `nil`.

**Usage:**

```go
pm, err := NewPackageManager("/usr/local/share/packagemanager/packages.json")
if err != nil {
    // Handle error
}
```

##### `Save`

**Definition:**

```go
func (pm *PackageManager) Save() error
```

**Description:**

Persists the current state of installed packages to the packages file. It serialises the `Packages` slice into JSON format and writes it to the file.

**Parameters:** None.

**Returns:**

- `error`: An error object if saving fails, otherwise `nil`.

**Usage:**

```go
err := pm.Save()
if err != nil {
    // Handle error
}
```

##### `AddPackage`

**Definition:**

```go
func (pm *PackageManager) AddPackage(pkg Package) error
```

**Description:**

Adds a new package to the `PackageManager`'s tracking system. It appends the package to the `Packages` slice and saves the updated list.

**Parameters:**

- `pkg` (`Package`): The package to be added.

**Returns:**

- `error`: An error object if adding or saving fails, otherwise `nil`.

**Usage:**

```go
newPkg := Package{
    UUID:        "unique-uuid",
    Name:        "SamplePackage",
    InstallPath: "/usr/local/share/packagemanager/SamplePackage",
    Executable:  "/usr/local/share/packagemanager/SamplePackage/bin/sample-exec",
}
err := pm.AddPackage(newPkg)
if err != nil {
    // Handle error
}
```

##### `RemovePackage`

**Definition:**

```go
func (pm *PackageManager) RemovePackage(uuid string) error
```

**Description:**

Removes a package from the `PackageManager`'s tracking system based on its UUID. It searches for the package, removes it from the `Packages` slice, and saves the updated list.

**Parameters:**

- `uuid` (`string`): The unique identifier of the package to be removed.

**Returns:**

- `error`: An error object if the package is not found or saving fails, otherwise `nil`.

**Usage:**

```go
err := pm.RemovePackage("unique-uuid")
if err != nil {
    // Handle error
}
```

##### `CreateDesktopFile`

**Definition:**

```go
func CreateDesktopFile(executablePath, packageName, installPath string) error
```

**Description:**

Generates a `.desktop` file for the given executable. The `.desktop` file integrates the application with desktop environments, allowing it to appear in application menus and support desktop shortcuts.

**Parameters:**

- `executablePath` (`string`): The absolute path to the executable file of the application.
- `packageName` (`string`): The user-friendly name of the package.
- `installPath` (`string`): The directory path where the package is installed.

**Returns:**

- `error`: An error object if the creation fails, otherwise `nil`.

**Usage:**

```go
err := CreateDesktopFile("/usr/local/share/packagemanager/SamplePackage/bin/sample-exec", "SamplePackage", "/usr/local/share/packagemanager/SamplePackage")
if err != nil {
    // Handle error
}
```

##### `RemoveDesktopFile`

**Definition:**

```go
func RemoveDesktopFile(packageName string) error
```

**Description:**

Deletes the `.desktop` file associated with the specified package. This ensures that the application is removed from desktop environment menus.

**Parameters:**

- `packageName` (`string`): The name of the package whose `.desktop` file needs to be removed.

**Returns:**

- `error`: An error object if the removal fails, otherwise `nil`.

**Usage:**

```go
err := RemoveDesktopFile("SamplePackage")
if err != nil {
    // Handle error
}
```

##### `ExtractTarGz`

**Definition:**

```go
func ExtractTarGz(archivePath, destDir string) error
```

**Description:**

Extracts a `.tar.gz` archive to the specified destination directory. It handles the creation of directories and files, sets appropriate permissions, and ensures that the extraction process is secure and efficient.

**Parameters:**

- `archivePath` (`string`): The file system path to the `.tar.gz` archive.
- `destDir` (`string`): The destination directory where the archive will be extracted.

**Returns:**

- `error`: An error object if the extraction fails, otherwise `nil`.

**Usage:**

```go
err := ExtractTarGz("/path/to/package.tar.gz", "/usr/local/share/packagemanager/SamplePackage")
if err != nil {
    // Handle error
}
```

##### `getDefaultIcon`

**Definition:**

```go
func getDefaultIcon(installPath string) string
```

**Description:**

Searches for a default icon within the installation directory. If no icon is found, it returns a predefined default icon path.

**Parameters:**

- `installPath` (`string`): The directory path where the package is installed.

**Returns:**

- `string`: The path to the found icon or the default icon path.

**Usage:**

```go
iconPath := getDefaultIcon("/usr/local/share/packagemanager/SamplePackage")
```

#### Detailed Function Descriptions

##### `ExtractTarGz`

```go
func ExtractTarGz(archivePath, destDir string) error {
    // Implementation as provided above
}
```

**Description:**

Extracts a `.tar.gz` archive to the specified destination directory. It meticulously handles the extraction process by:

1. **Opening the Archive:**  
   Opens the `.tar.gz` file for reading.

2. **Decompressing the Archive:**  
   Utilises a gzip reader to decompress the archive.

3. **Reading the Tar Archive:**  
   Iterates through each entry in the tar archive.

4. **Handling Directories and Files:**  
   - **Directories:** Creates directories with the specified permissions.
   - **Files:** Creates files, copies their contents, and sets appropriate permissions.

5. **Skipping Unknown Types:**  
   Ignores any unknown file types within the archive, informing the user.

**Error Handling:**

Returns descriptive errors at each failure point, ensuring that issues are communicated clearly to the caller.

---

### main Package

The `main` package serves as the entry point for the PackageManager application. It initialises the command-line interface, ensures necessary privileges, and orchestrates the execution of various package management commands.

#### `main`

**Definition:**

```go
func main() {
    // Function implementation
}
```

**Description:**

The `main` function performs the following operations:

1. **Privilege Verification:**

    ```go
    if os.Geteuid() != 0 {
        fmt.Println("You need to have root privileges to run this program.")
        os.Exit(1)
    }
    ```

    - **Purpose:**  
      Ensures that the program is executed with root privileges, which are typically required for installing and managing system-wide packages.

    - **Behavior:**  
      - If the program is not run as root, it prints an error message and exits with a non-zero status code.
      - This prevents unauthorised users from performing sensitive package management operations.

2. **CLI Initialisation:**

    ```go
    rootCmd := &cobra.Command{
        Use:   "packagemanager",
        Short: "A simple package manager",
        Long:  `PackageManager is a simple tool to install, uninstall, and manage software packages.`,
    }
    ```

    - **Purpose:**  
      Defines the root command for the CLI application using the Cobra library.

    - **Fields:**
      - **Use:**  
        Specifies the command's invocation name (`packagemanager`).

      - **Short:**  
        Provides a brief description of the command, displayed in help messages.

      - **Long:**  
        Offers a more detailed description, enhancing user understanding.

3. **Adding Subcommands:**

    ```go
    rootCmd.AddCommand(cmd.InstallCmd)
    rootCmd.AddCommand(cmd.UninstallCmd)
    rootCmd.AddCommand(cmd.ListCmd)
    ```

    - **Purpose:**  
      Integrates the `install`, `uninstall`, and `list` subcommands into the root command.

    - **Implementation:**  
      Each subcommand is defined in the `cmd` package and encapsulates its own logic and argument parsing.

4. **Executing the Root Command:**

    ```go
    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    ```

    - **Purpose:**  
      Parses the command-line arguments and executes the appropriate subcommand based on user input.

    - **Behavior:**  
      - If an error occurs during command execution (e.g., invalid arguments, execution failure), it prints the error message and exits with a non-zero status code.
      - Successful execution completes without triggering this error block.

    - **Error Handling:**  
      Ensures that users are informed of any issues that arise during command execution, facilitating troubleshooting.

**Usage Example:**

```sh
# To install a package
sudo packagemanager install /path/to/package.tar.gz

# To list installed packages
sudo packagemanager list

# To uninstall a package
sudo packagemanager uninstall SamplePackage
```

**Sample Output:**

```sh
$ sudo packagemanager list
NAME            INSTALL PATH                                     EXECUTABLE
----            ------------                                     -----------
SamplePackage   /usr/local/share/packagemanager/123e4567-89ab-cdef-0123-456789abcdef-SamplePackage   /usr/local/share/packagemanager/123e4567-89ab-cdef-0123-456789abcdef-SamplePackage/bin/sample-exec
AnotherPackage  /usr/local/share/packagemanager/abcdef12-3456-7890-abcd-ef1234567890-AnotherPackage   /usr/local/bin/another-exec
```

---

## Directory Structure

The PackageManager application is organised into a clear and modular directory structure, promoting maintainability and scalability.

```text
PackageManager/
├── cmd/
│   ├── install.go
│   ├── list.go
│   └── uninstall.go
├── pkg/
│   ├── desktop.go
│   ├── extractor.go
│   └── package.go
├── main.go
└── README.md
```

### Breakdown

- **cmd/**:  
  Contains the command definitions for the CLI application. Each `.go` file corresponds to a specific subcommand.
  
  - `install.go`: Defines the `install` command.
  - `list.go`: Defines the `list` command.
  - `uninstall.go`: Defines the `uninstall` command.

- **pkg/**:  
  Houses utility functions and structures that support package management operations.
  
  - `desktop.go`: Manages the creation and removal of `.desktop` files for desktop environment integration.
  - `extractor.go`: Handles the extraction of `.tar.gz` archives.
  - `package.go`: Defines the `Package` and `PackageManager` types and their associated methods.

- **main.go**:  
  The entry point of the application. It initialises the CLI, ensures necessary privileges, and integrates all commands.

- **README.md**:  
  Provides an overview, installation instructions, usage examples, and other relevant information about the PackageManager application.

---

## Additional Notes

- **Root Privileges:**  
  The PackageManager requires root privileges to perform operations that modify system directories (e.g., `/usr/local/bin`, `/usr/share/applications`). Ensure that you run the commands with appropriate permissions (e.g., using `sudo`).

- **Error Handling:**  
  All functions and commands are designed to handle errors gracefully, providing descriptive messages to aid in troubleshooting.

- **Extensibility:**  
  The modular design allows for easy addition of new commands and functionalities. Developers can extend the `cmd` and `pkg` packages to introduce new features as needed.

- **Cross-Platform Considerations:**  
  While the current implementation targets Unix-like systems, further enhancements can be made to support other operating systems by handling platform-specific paths and permissions.
