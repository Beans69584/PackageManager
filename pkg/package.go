package pkg

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Package represents an installed package with its essential metadata.
// This struct is used to track and manage packages within the PackageManager.
type Package struct {
	UUID        string `json:"uuid"`         // A unique identifier for the package installation.
	Name        string `json:"name"`         // The user-friendly name of the package.
	InstallPath string `json:"install_path"` // The filesystem path where the package is installed.
	Executable  string `json:"executable"`   // The path to the package's main executable file.
}

// PackageManager manages the collection of installed packages.
// It handles loading from and saving to the packages database file.
type PackageManager struct {
	PackagesFile string    // The path to the JSON file that stores package metadata.
	Packages     []Package // A slice containing all the currently installed packages.
}

// NewPackageManager creates and initializes a new PackageManager.
// It loads existing packages from the specified packages file or creates a new one if it doesn't exist.
//
// Parameters:
//   - packagesFile (string): The path to the JSON file that stores package metadata.
//
// Returns:
//   - *PackageManager: A pointer to the initialized PackageManager.
//   - error: An error object if initialization fails, otherwise nil.
func NewPackageManager(packagesFile string) (*PackageManager, error) {
	pm := &PackageManager{
		PackagesFile: packagesFile,
		Packages:     []Package{},
	}

	// Check if the packages file exists.
	if _, err := os.Stat(packagesFile); os.IsNotExist(err) {
		// Create the packages directory if it doesn't exist.
		packagesDir := filepath.Dir(packagesFile)
		err := os.MkdirAll(packagesDir, 0755)
		if err != nil {
			return nil, fmt.Errorf("error creating packages directory: %v", err)
		}

		// Create an empty packages file.
		file, err := os.Create(packagesFile)
		if err != nil {
			return nil, fmt.Errorf("error creating packages file: %v", err)
		}
		defer file.Close()

		// Initialize the packages file with an empty JSON array.
		_, err = file.Write([]byte("[]"))
		if err != nil {
			return nil, fmt.Errorf("error initializing packages file: %v", err)
		}
	}

	// Load existing packages from the packages file.
	data, err := os.ReadFile(packagesFile)
	if err != nil {
		return nil, fmt.Errorf("error reading packages file: %v", err)
	}

	// Unmarshal the JSON data into the Packages slice.
	if err := json.Unmarshal(data, &pm.Packages); err != nil {
		return nil, fmt.Errorf("error unmarshalling packages file: %v", err)
	}

	return pm, nil
}

// Save persists the current state of installed packages to the packages file.
// It serializes the Packages slice into JSON format and writes it to the file.
//
// Returns:
//   - error: An error object if saving fails, otherwise nil.
func (pm *PackageManager) Save() error {
	// Marshal the Packages slice into indented JSON for readability.
	data, err := json.MarshalIndent(pm.Packages, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling packages: %v", err)
	}

	// Write the JSON data to the packages file with appropriate permissions.
	if err := os.WriteFile(pm.PackagesFile, data, 0644); err != nil {
		return fmt.Errorf("error writing packages file: %v", err)
	}

	return nil
}

// AddPackage adds a new package to the PackageManager's tracking system.
// It appends the package to the Packages slice and saves the updated list.
//
// Parameters:
//   - pkg (Package): The package to be added.
//
// Returns:
//   - error: An error object if adding or saving fails, otherwise nil.
func (pm *PackageManager) AddPackage(pkg Package) error {
	pm.Packages = append(pm.Packages, pkg)
	return pm.Save()
}

// RemovePackage removes a package from the PackageManager's tracking system based on its UUID.
// It searches for the package, removes it from the Packages slice, and saves the updated list.
//
// Parameters:
//   - uuid (string): The unique identifier of the package to be removed.
//
// Returns:
//   - error: An error object if the package is not found or saving fails, otherwise nil.
func (pm *PackageManager) RemovePackage(uuid string) error {
	index := -1
	// Iterate through the Packages slice to find the package with the matching UUID.
	for i, pkg := range pm.Packages {
		if pkg.UUID == uuid {
			index = i
			break
		}
	}

	// If the package is not found, return an error.
	if index == -1 {
		return fmt.Errorf("package with UUID %s not found", uuid)
	}

	// Remove the package from the slice.
	pm.Packages = append(pm.Packages[:index], pm.Packages[index+1:]...)
	return pm.Save()
}
