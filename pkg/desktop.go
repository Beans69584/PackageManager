package pkg

import (
    "fmt"
    "os"
    "path/filepath"
    "runtime"
    "strings"
)

// CreateDesktopFile generates a .desktop file for the given executable
func CreateDesktopFile(executablePath, packageName string) error {
    desktopDir := "/usr/share/applications"
    desktopFilePath := filepath.Join(desktopDir, fmt.Sprintf("%s.desktop", strings.ToLower(packageName)))

    // Define the content of the .desktop file
    desktopContent := fmt.Sprintf(`[Desktop Entry]
Type=Application
Name=%s
Exec=%s
Icon=%s
Terminal=false
Categories=Utility;`, packageName, executablePath, getDefaultIcon())

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

// getDefaultIcon returns a default icon path based on the operating system
func getDefaultIcon() string {
    if runtime.GOOS == "linux" {
        return "/usr/share/pixmaps/default-icon.png" // Update this path if you have a specific icon
    }
    // Add more OS-specific icon paths if needed
    return "/usr/share/pixmaps/default-icon.png"
}

