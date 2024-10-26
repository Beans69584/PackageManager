package pkg

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// CreateDesktopFile generates a .desktop file for the given executable
func CreateDesktopFile(executablePath, packageName, installPath string) error {
	desktopDir := "/usr/share/applications"
	desktopFilePath := filepath.Join(desktopDir, fmt.Sprintf("%s.desktop", strings.ToLower(packageName)))

	// Get the icon path from the install directory
	iconPath := getDefaultIcon(installPath)

	// Define the content of the .desktop file
	desktopContent := fmt.Sprintf(`[Desktop Entry]
Type=Application
Name=%s
Exec=%s
Icon=%s
Terminal=false
Categories=Utility;`, packageName, executablePath, iconPath)

	// Create or overwrite the .desktop file
	file, err := os.Create(desktopFilePath)
	if err != nil {
		return fmt.Errorf("error creating .desktop file: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(desktopContent)
	if err != nil {
		return fmt.Errorf("error writing to .desktop file: %v", err)
	}

	fmt.Printf("Created .desktop file at %s\n", desktopFilePath)
	return nil
}

// RemoveDesktopFile deletes the .desktop file associated with the package
func RemoveDesktopFile(packageName string) error {
	desktopDir := "/usr/share/applications"
	desktopFilePath := filepath.Join(desktopDir, fmt.Sprintf("%s.desktop", strings.ToLower(packageName)))

	if _, err := os.Stat(desktopFilePath); os.IsNotExist(err) {
		// .desktop file does not exist; nothing to do
		return nil
	}

	err := os.Remove(desktopFilePath)
	if err != nil {
		return fmt.Errorf("error removing .desktop file: %v", err)
	}

	fmt.Printf("Removed .desktop file at %s\n", desktopFilePath)
	return nil
}

func getDefaultIcon(installPath string) string {
	// Possible icon file extensions
	iconExtensions := []string{".png", ".jpg", ".jpeg", ".ico", ".svg"}

	var iconPath string

	err := filepath.WalkDir(installPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("Error accessing path %q: %v\n", path, err)
			return err
		}

		if !d.IsDir() {
			for _, ext := range iconExtensions {
				if strings.HasSuffix(strings.ToLower(d.Name()), ext) {
					// Found an icon file
					iconPath = path
					fmt.Printf("Found icon file: %s\n", iconPath)
					return fs.SkipDir // Stop walking once an icon is found
				}
			}
		}
		return nil
	})

	if err != nil && err != fs.SkipDir {
		fmt.Printf("Error walking the path %q: %v\n", installPath, err)
	}

	if iconPath != "" {
		return iconPath
	}

	// If no icon found, return default icon
	return "/usr/share/pixmaps/default-icon.png"
}
