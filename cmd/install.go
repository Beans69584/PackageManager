package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/Beans69584/PackageManager/pkg"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

// InstallCmd represents the 'install' command for the PackageManager.
// It enables users to install a package from a tar.gz archive.
var InstallCmd = &cobra.Command{
	Use:   "install [archive.tar.gz]",
	Short: "Install a package from a tar.gz archive",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Retrieve the path to the archive from the command arguments.
		archivePath := args[0]

		// Verify that the specified archive exists.
		if _, err := os.Stat(archivePath); os.IsNotExist(err) {
			fmt.Printf("Error: Archive %s does not exist.\n", archivePath)
			os.Exit(1)
		}

		// Define the base directory where packages will be installed.
		packagesDir := "/usr/local/share/packagemanager"

		// Generate a unique identifier for this installation instance.
		installUUID := uuid.New().String()

		// Determine the default package name by stripping extensions from the archive filename.
		defaultPackageName := strings.TrimSuffix(strings.TrimSuffix(filepath.Base(archivePath), ".tar.gz"), ".tgz")

		// Construct the full installation path using the base directory, UUID, and default package name.
		installPath := filepath.Join(packagesDir, fmt.Sprintf("%s-%s", installUUID, defaultPackageName))

		// Initialise the PackageManager, responsible for tracking installed packages.
		pm, err := pkg.NewPackageManager(filepath.Join(packagesDir, "packages.json"))
		if err != nil {
			fmt.Printf("Error initialising PackageManager: %v\n", err)
			os.Exit(1)
		}

		// Extract the contents of the archive to the designated installation path.
		fmt.Printf("Extracting %s to %s...\n", archivePath, installPath)
		err = pkg.ExtractTarGz(archivePath, installPath)
		if err != nil {
			fmt.Printf("Error extracting archive: %v\n", err)
			os.Exit(1)
		}

		// Prompt the user to input a friendly name for the package.
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("Enter a friendly name for the package [%s]: ", defaultPackageName)
		inputName, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input: %v\n", err)
			os.Exit(1)
		}
		inputName = strings.TrimSpace(inputName)
		if inputName == "" {
			inputName = defaultPackageName
		}
		packageName := inputName

		// Recursively search for executable files within the installation directory.
		executables, err := findExecutablesRecursively(installPath)
		if err != nil {
			fmt.Printf("Error searching for executables: %v\n", err)
			os.Exit(1)
		}

		// If no executables are found, inform the user and exit.
		if len(executables) == 0 {
			fmt.Println("No executables found in the package.")
			os.Exit(1)
		}

		var selectedExecutable string

		// If only one executable is found, select it automatically.
		if len(executables) == 1 {
			selectedExecutable = executables[0]
			fmt.Printf("Automatically selected executable: %s\n", filepath.Base(selectedExecutable))
		} else {
			// If multiple executables are found, list them and prompt the user to select one.
			fmt.Println("Multiple executables found:")
			for i, execPath := range executables {
				relPath, _ := filepath.Rel(installPath, execPath)
				fmt.Printf("  %d) %s\n", i+1, relPath)
			}

			// Improved executable selection process.
			for {
				fmt.Printf("Select an executable to symlink (1-%d): ", len(executables))
				input, err := reader.ReadString('\n')
				if err != nil {
					fmt.Printf("Error reading input: %v\n", err)
					continue
				}
				input = strings.TrimSpace(input)
				choice, err := strconv.Atoi(input)
				if err != nil || choice < 1 || choice > len(executables) {
					fmt.Println("Invalid selection. Please enter a valid number.")
					continue
				}
				selectedExecutable = executables[choice-1]
				fmt.Printf("Selected executable: %s\n", filepath.Base(selectedExecutable))
				break
			}
		}

		// Create a symbolic link in /usr/local/bin pointing to the selected executable.
		symlinkName := filepath.Base(selectedExecutable)
		symlinkPath := filepath.Join("/usr/local/bin", symlinkName)

		// Check if the symlink path already exists and handle accordingly.
		if _, err := os.Lstat(symlinkPath); err == nil {
			fmt.Printf("Symlink %s already exists. Overwrite? (y/n): ", symlinkPath)
			overwriteInput, err := reader.ReadString('\n')
			if err != nil {
				fmt.Printf("Error reading input: %v\n", err)
				os.Exit(1)
			}
			overwriteInput = strings.TrimSpace(strings.ToLower(overwriteInput))
			if overwriteInput != "y" && overwriteInput != "yes" {
				fmt.Println("Installation aborted by user.")
				os.Exit(0)
			}

			// Remove the existing symlink to make way for the new one.
			err = os.Remove(symlinkPath)
			if err != nil {
				fmt.Printf("Error removing existing symlink: %v\n", err)
				os.Exit(1)
			}
		}

		// Create the new symlink.
		err = os.Symlink(selectedExecutable, symlinkPath)
		if err != nil {
			fmt.Printf("Error creating symlink: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Created symlink: %s -> %s\n", symlinkPath, selectedExecutable)

		// Create a .desktop file to integrate the application with desktop environments.
		err = pkg.CreateDesktopFile(selectedExecutable, packageName, installPath)
		if err != nil {
			fmt.Printf("Error creating .desktop file: %v\n", err)
			// Optionally, remove the symlink if .desktop creation fails.
			os.Remove(symlinkPath)
			os.Exit(1)
		}

		// Add the package to the PackageManager's tracking system.
		newPackage := pkg.Package{
			UUID:        installUUID,
			Name:        packageName,
			InstallPath: installPath,
			Executable:  selectedExecutable,
		}

		err = pm.AddPackage(newPackage)
		if err != nil {
			fmt.Printf("Error adding package to PackageManager: %v\n", err)
			// Optionally, remove symlink and .desktop file if tracking fails.
			os.Remove(symlinkPath)
			pkg.RemoveDesktopFile(packageName)
			os.Exit(1)
		}

		fmt.Printf("Package '%s' installed successfully.\n", packageName)

		// Attempt to terminate the AGS bus to refresh desktop entries.
		killCmd := exec.Command("ags", "quit")
		err = killCmd.Run()
		if err != nil {
			fmt.Printf("Warning: Failed to kill AGS bus: %v\n", err)
		} else {
			fmt.Println("AGS bus killed successfully.")
		}
	},
}

// findExecutablesRecursively searches for executable files within the given directory and its subdirectories.
// It returns a slice of paths to executable files found.
func findExecutablesRecursively(root string) ([]string, error) {
	var executables []string

	// Walk through the directory tree rooted at 'root'.
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip non-regular files (e.g., directories, symlinks).
		if !info.Mode().IsRegular() {
			return nil
		}

		// Check if the file is executable.
		if isExecutable(path, info) {
			executables = append(executables, path)
		}

		return nil
	})

	return executables, err
}

// isExecutable determines if a file is executable based on its permissions.
// For Unix-like systems, it checks the executable bits in the file mode.
// For Windows, it checks for common executable file extensions.
func isExecutable(path string, info os.FileInfo) bool {
	mode := info.Mode()

	// On Unix-like systems, check if any of the executable bits are set.
	if mode&0111 != 0 {
		return true
	}

	// On Windows systems, check for known executable file extensions.
	if runtime.GOOS == "windows" {
		ext := strings.ToLower(filepath.Ext(path))
		executableExtensions := []string{".exe", ".bat", ".cmd", ".ps1"}
		for _, e := range executableExtensions {
			if ext == e {
				return true
			}
		}
	}

	return false
}
