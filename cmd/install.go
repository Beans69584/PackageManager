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

var InstallCmd = &cobra.Command{
	Use:   "install [archive.tar.gz]",
	Short: "Install a package from a tar.gz archive",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		archivePath := args[0]

		// Verify the archive exists
		if _, err := os.Stat(archivePath); os.IsNotExist(err) {
			fmt.Printf("Error: Archive %s does not exist.\n", archivePath)
			os.Exit(1)
		}

		// Define the installation directory
		packagesDir := "/usr/local/share/packagemanager"
		installUUID := uuid.New().String()
		defaultPackageName := strings.TrimSuffix(strings.TrimSuffix(filepath.Base(archivePath), ".tar.gz"), ".tgz")
		installPath := filepath.Join(packagesDir, fmt.Sprintf("%s-%s", installUUID, defaultPackageName))

		// Create PackageManager
		pm, err := pkg.NewPackageManager(filepath.Join(packagesDir, "packages.json"))
		if err != nil {
			fmt.Printf("Error initializing PackageManager: %v\n", err)
			os.Exit(1)
		}

		// Extract the archive
		fmt.Printf("Extracting %s to %s...\n", archivePath, installPath)
		err = pkg.ExtractTarGz(archivePath, installPath)
		if err != nil {
			fmt.Printf("Error extracting archive: %v\n", err)
			os.Exit(1)
		}

		// Prompt the user for a friendly name
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

		// Recursively find executables in the installation directory
		executables, err := findExecutablesRecursively(installPath)
		if err != nil {
			fmt.Printf("Error searching for executables: %v\n", err)
			os.Exit(1)
		}

		if len(executables) == 0 {
			fmt.Println("No executables found in the package.")
			os.Exit(1)
		}

		var selectedExecutable string

		if len(executables) == 1 {
			selectedExecutable = executables[0]
			fmt.Printf("Automatically selected executable: %s\n", filepath.Base(selectedExecutable))
		} else {
			// List executables and allow user to select
			fmt.Println("Multiple executables found:")
			for i, execPath := range executables {
				relPath, _ := filepath.Rel(installPath, execPath)
				fmt.Printf("  %d) %s\n", i+1, relPath)
			}

			// Improved executable selection
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

		// Create symlink in /usr/local/bin
		symlinkName := filepath.Base(selectedExecutable)
		symlinkPath := filepath.Join("/usr/local/bin", symlinkName)

		// Check if symlinkPath already exists
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

			// Remove existing symlink
			err = os.Remove(symlinkPath)
			if err != nil {
				fmt.Printf("Error removing existing symlink: %v\n", err)
				os.Exit(1)
			}
		}

		// Create the symlink
		err = os.Symlink(selectedExecutable, symlinkPath)
		if err != nil {
			fmt.Printf("Error creating symlink: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Created symlink: %s -> %s\n", symlinkPath, selectedExecutable)

		// Create a .desktop file
		err = pkg.CreateDesktopFile(selectedExecutable, packageName, installPath)
		if err != nil {
			fmt.Printf("Error creating .desktop file: %v\n", err)
			// Optionally, remove the symlink if .desktop creation fails
			os.Remove(symlinkPath)
			os.Exit(1)
		}

		// Add package to PackageManager
		newPackage := pkg.Package{
			UUID:        installUUID,
			Name:        packageName,
			InstallPath: installPath,
			Executable:  selectedExecutable,
		}

		err = pm.AddPackage(newPackage)
		if err != nil {
			fmt.Printf("Error adding package to PackageManager: %v\n", err)
			// Optionally, remove symlink and .desktop file if tracking fails
			os.Remove(symlinkPath)
			pkg.RemoveDesktopFile(packageName)
			os.Exit(1)
		}

		fmt.Printf("Package '%s' installed successfully.\n", packageName)

		// Kill ags
		killCmd := exec.Command("ags", "quit")
		err = killCmd.Run()
		if err != nil {
			fmt.Printf("Warning: Failed to kill AGS bus: %v\n", err)
		} else {
			fmt.Println("AGS bus killed successfully.")
		}
	},
}

// findExecutablesRecursively searches for executable files within the given directory and its subdirectories
func findExecutablesRecursively(root string) ([]string, error) {
	var executables []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.Mode().IsRegular() {
			return nil
		}

		if isExecutable(path, info) {
			executables = append(executables, path)
		}

		return nil
	})

	return executables, err
}

// isExecutable determines if a file is executable based on its permissions
func isExecutable(path string, info os.FileInfo) bool {
	mode := info.Mode()

	// On Unix-like systems
	if mode&0111 != 0 {
		return true
	}

	// On Windows, check for executable extensions
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
