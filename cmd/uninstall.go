package cmd

import (
    "fmt"
    "os"
    "path/filepath"

    "github.com/spf13/cobra"
    "github.com/Beans69584/PackageManager/pkg"
)

var UninstallCmd = &cobra.Command{
    Use:   "uninstall [package_name]",
    Short: "Uninstall a package by name",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        packageName := args[0]

        // Define the installation directory
        packagesDir := "/usr/local/share/packagemanager"

        // Create PackageManager
        pm, err := pkg.NewPackageManager(filepath.Join(packagesDir, "packages.json"))
        if err != nil {
            fmt.Printf("Error initializing PackageManager: %v\n", err)
            os.Exit(1)
        }

        // Find the package
        var targetPackage *pkg.Package
        for _, p := range pm.Packages {
            if p.Name == packageName {
                targetPackage = &p
                break
            }
        }

        if targetPackage == nil {
            fmt.Printf("Package %s not found.\n", packageName)
            os.Exit(1)
        }

        // Remove symlink
        symlinkPath := filepath.Join("/usr/local/bin", filepath.Base(targetPackage.Executable))
        err = os.Remove(symlinkPath)
        if err != nil {
            fmt.Printf("Error removing symlink: %v\n", err)
            // Proceeding even if symlink removal fails
        } else {
            fmt.Printf("Removed symlink: %s\n", symlinkPath)
        }

        // Remove .desktop file
        err = pkg.RemoveDesktopFile(packageName)
        if err != nil {
            fmt.Printf("Error removing .desktop file: %v\n", err)
            // Proceeding even if .desktop removal fails
        }

        // Remove installation directory
        err = os.RemoveAll(targetPackage.InstallPath)
        if err != nil {
            fmt.Printf("Error removing installation directory: %v\n", err)
            // Proceeding even if directory removal fails
        } else {
            fmt.Printf("Removed installation directory: %s\n", targetPackage.InstallPath)
        }

        // Remove package from PackageManager
        err = pm.RemovePackage(targetPackage.UUID)
        if err != nil {
            fmt.Printf("Error removing package from PackageManager: %v\n", err)
            os.Exit(1)
        }

        fmt.Printf("Package %s uninstalled successfully.\n", packageName)
    },
}

