package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Beans69584/PackageManager/pkg"
	"github.com/spf13/cobra"
)

// UninstallCmd represents the 'uninstall' command for the PackageManager.
// It enables users to remove an installed package by specifying its name.
var UninstallCmd = &cobra.Command{
	Use:   "uninstall [package_name]",
	Short: "Uninstall a package by name",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Retrieve the package name from the command arguments.
		packageName := args[0]

		// Define the base directory where packages are installed.
		packagesDir := "/usr/local/share/packagemanager"

		// Initialise the PackageManager, which manages the tracking of installed packages.
		pm, err := pkg.NewPackageManager(filepath.Join(packagesDir, "packages.json"))
		if err != nil {
			// If there is an error initializing the PackageManager, inform the user and exit.
			fmt.Printf("Error initialising PackageManager: %v\n", err)
			os.Exit(1)
		}

		// Search for the target package by name within the list of installed packages.
		var targetPackage *pkg.Package
		for _, p := range pm.Packages {
			if p.Name == packageName {
				targetPackage = &p
				break
			}
		}

		// If the package is not found, inform the user and exit.
		if targetPackage == nil {
			fmt.Printf("Package %s not found.\n", packageName)
			os.Exit(1)
		}

		// Construct the path to the symbolic link in /usr/local/bin.
		symlinkPath := filepath.Join("/usr/local/bin", filepath.Base(targetPackage.Executable))

		// Attempt to remove the symbolic link.
		err = os.Remove(symlinkPath)
		if err != nil {
			// If removing the symlink fails, inform the user but proceed with uninstallation.
			fmt.Printf("Error removing symlink: %v\n", err)
		} else {
			// Inform the user that the symlink has been removed successfully.
			fmt.Printf("Removed symlink: %s\n", symlinkPath)
		}

		// Attempt to remove the associated .desktop file.
		err = pkg.RemoveDesktopFile(packageName)
		if err != nil {
			// If removing the .desktop file fails, inform the user but proceed.
			fmt.Printf("Error removing .desktop file: %v\n", err)
		} else {
			// Inform the user that the .desktop file has been removed successfully.
			fmt.Printf("Removed .desktop file for package: %s\n", packageName)
		}

		// Attempt to remove the installation directory and all its contents.
		err = os.RemoveAll(targetPackage.InstallPath)
		if err != nil {
			// If removing the installation directory fails, inform the user but proceed.
			fmt.Printf("Error removing installation directory: %v\n", err)
		} else {
			// Inform the user that the installation directory has been removed successfully.
			fmt.Printf("Removed installation directory: %s\n", targetPackage.InstallPath)
		}

		// Attempt to remove the package entry from the PackageManager's tracking system.
		err = pm.RemovePackage(targetPackage.UUID)
		if err != nil {
			// If removing the package from tracking fails, inform the user and exit with an error.
			fmt.Printf("Error removing package from PackageManager: %v\n", err)
			os.Exit(1)
		}

		// Inform the user that the package has been uninstalled successfully.
		fmt.Printf("Package '%s' uninstalled successfully.\n", packageName)

		// Attempt to terminate the AGS bus to refresh desktop entries.
		killCmd := exec.Command("ags", "quit")
		err = killCmd.Run()
		if err != nil {
			// If killing the AGS bus fails, inform the user but do not treat it as a critical error.
			fmt.Printf("Error killing AGS bus: %v\n", err)
		} else {
			// Inform the user that the AGS bus has been terminated successfully.
			fmt.Println("AGS bus killed successfully.")
		}
	},
}
