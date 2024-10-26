package pkg

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// CreateDesktopFile generates a .desktop file for the given executable.
// The .desktop file is used to integrate the application with desktop environments,
// allowing it to appear in application menus and support desktop shortcuts.
func CreateDesktopFile(executablePath, packageName, installPath string) error {
	// Define the directory where .desktop files are stored.
	desktopDir := "/usr/share/applications"

	// Construct the full path to the .desktop file, using the package name in lowercase.
	desktopFilePath := filepath.Join(desktopDir, fmt.Sprintf("%s.desktop", strings.ToLower(packageName)))

	// Retrieve the path to the application's icon.
	iconPath := getDefaultIcon(installPath)

	// Define the content of the .desktop file following the Desktop Entry Specification.
	desktopContent := fmt.Sprintf(`[Desktop Entry]
Type=Application
Name=%s
Exec=%s
Icon=%s
Terminal=false
Categories=Utility;`, packageName, executablePath, iconPath)

	// Create or overwrite the .desktop file with the defined content.
	file, err := os.Create(desktopFilePath)
	if err != nil {
		return fmt.Errorf("error creating .desktop file: %v", err)
	}
	defer file.Close()

	// Write the desktop content to the file.
	_, err = file.WriteString(desktopContent)
	if err != nil {
		return fmt.Errorf("error writing to .desktop file: %v", err)
	}

	// Inform the user that the .desktop file has been created successfully.
	fmt.Printf("Created .desktop file at %s\n", desktopFilePath)
	return nil
}

// RemoveDesktopFile deletes the .desktop file associated with the specified package.
// This function ensures that the application is removed from desktop environment menus.
func RemoveDesktopFile(packageName string) error {
	// Define the directory where .desktop files are stored.
	desktopDir := "/usr/share/applications"

	// Construct the full path to the .desktop file, using the package name in lowercase.
	desktopFilePath := filepath.Join(desktopDir, fmt.Sprintf("%s.desktop", strings.ToLower(packageName)))

	// Check if the .desktop file exists.
	if _, err := os.Stat(desktopFilePath); os.IsNotExist(err) {
		// If the .desktop file does not exist, there's nothing to remove.
		return nil
	}

	// Attempt to remove the .desktop file.
	err := os.Remove(desktopFilePath)
	if err != nil {
		return fmt.Errorf("error removing .desktop file: %v", err)
	}

	// Inform the user that the .desktop file has been removed successfully.
	fmt.Printf("Removed .desktop file at %s\n", desktopFilePath)
	return nil
}

// getDefaultIcon searches for a default icon within the installation directory.
// If no icon is found, it returns a predefined default icon path.
func getDefaultIcon(installPath string) string {
	// Define possible icon file extensions to search for.
	iconExtensions := []string{".png", ".jpg", ".jpeg", ".ico", ".svg"}

	var iconPath string

	// Walk through the installation directory to find an icon file.
	err := filepath.WalkDir(installPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			// If there's an error accessing a path, log it and continue.
			fmt.Printf("Error accessing path %q: %v\n", path, err)
			return err
		}

		// Skip directories; we're only interested in files.
		if d.IsDir() {
			return nil
		}

		// Check if the current file has one of the specified icon extensions.
		for _, ext := range iconExtensions {
			if strings.HasSuffix(strings.ToLower(d.Name()), ext) {
				// If an icon file is found, set the iconPath and stop walking the directory.
				iconPath = path
				fmt.Printf("Found icon file: %s\n", iconPath)
				return fs.SkipDir
			}
		}
		return nil
	})

	// Handle any errors encountered during the directory walk, excluding fs.SkipDir.
	if err != nil && err != fs.SkipDir {
		fmt.Printf("Error walking the path %q: %v\n", installPath, err)
	}

	// If an icon was found, return its path.
	if iconPath != "" {
		return iconPath
	}

	// If no icon was found, return the path to the default icon.
	return "/usr/share/pixmaps/default-icon.png"
}
