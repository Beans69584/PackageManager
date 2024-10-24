package cmd

import (
    "fmt"
    "os"
    "path/filepath"
    "runtime"
    "strings"

    "github.com/google/uuid"
    "github.com/spf13/cobra"
    "github.com/Beans69584/PackageManager/pkg"
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
        packageName := strings.TrimSuffix(strings.TrimSuffix(filepath.Base(archivePath), ".tar.gz"), ".tgz")
        installPath := filepath.Join(packagesDir, fmt.Sprintf("%s-%s", installUUID, packageName))

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
            for i, exec := range executables {
                fmt.Printf("  %d) %s\n", i+1, exec)
            }

            var choice int
            fmt.Printf("Select an executable to symlink (1-%d): ", len(executables))
            _, err := fmt.Scanf("%d", &choice)
            if err != nil || choice < 1 || choice > len(executables) {
                fmt.Println("Invalid selection.")
                os.Exit(1)
            }

            selectedExecutable = executables[choice-1]
            fmt.Printf("Selected executable: %s\n", filepath.Base(selectedExecutable))
        }

        // Create symlink in /usr/local/bin
        symlinkPath := filepath.Join("/usr/local/bin", filepath.Base(selectedExecutable))

        // Check if symlinkPath already exists
        if _, err := os.Lstat(symlinkPath); err == nil {
            fmt.Printf("Symlink %s already exists. Overwrite? (y/n): ", symlinkPath)
            var overwrite string
            fmt.Scanf("%s", &overwrite)
            if strings.ToLower(overwrite) != "y" {
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
        err = pkg.CreateDesktopFile(selectedExecutable, packageName)
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

        fmt.Printf("Package %s installed successfully.\n", packageName)
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

