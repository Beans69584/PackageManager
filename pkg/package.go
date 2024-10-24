package pkg

import (
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
)

// Package represents an installed package
type Package struct {
    UUID        string `json:"uuid"`
    Name        string `json:"name"`
    InstallPath string `json:"install_path"`
    Executable  string `json:"executable"`
}

// PackageManager manages installed packages
type PackageManager struct {
    PackagesFile string
    Packages     []Package
}

// NewPackageManager creates a new PackageManager
func NewPackageManager(packagesFile string) (*PackageManager, error) {
    pm := &PackageManager{
        PackagesFile: packagesFile,
        Packages:     []Package{},
    }

    // Check if the packages file exists
    if _, err := os.Stat(packagesFile); os.IsNotExist(err) {
        // Create the packages directory
        packagesDir := filepath.Dir(packagesFile)
        err := os.MkdirAll(packagesDir, 0755)
        if err != nil {
            return nil, fmt.Errorf("error creating packages directory: %v", err)
        }

        // Create an empty packages file
        file, err := os.Create(packagesFile)
        if err != nil {
            return nil, fmt.Errorf("error creating packages file: %v", err)
        }
        defer file.Close()

        // Initialize with an empty JSON array
        _, err = file.Write([]byte("[]"))
        if err != nil {
            return nil, fmt.Errorf("error initializing packages file: %v", err)
        }
    }

    // Load existing packages
    data, err := os.ReadFile(packagesFile)
    if err != nil {
        return nil, fmt.Errorf("error reading packages file: %v", err)
    }

    if err := json.Unmarshal(data, &pm.Packages); err != nil {
        return nil, fmt.Errorf("error unmarshalling packages file: %v", err)
    }

    return pm, nil
}

// Save saves the current packages to the packages file
func (pm *PackageManager) Save() error {
    data, err := json.MarshalIndent(pm.Packages, "", "  ")
    if err != nil {
        return fmt.Errorf("error marshalling packages: %v", err)
    }

    if err := os.WriteFile(pm.PackagesFile, data, 0644); err != nil {
        return fmt.Errorf("error writing packages file: %v", err)
    }

    return nil
}

// AddPackage adds a new package to the manager
func (pm *PackageManager) AddPackage(pkg Package) error {
    pm.Packages = append(pm.Packages, pkg)
    return pm.Save()
}

// RemovePackage removes a package by UUID
func (pm *PackageManager) RemovePackage(uuid string) error {
    index := -1
    for i, pkg := range pm.Packages {
        if pkg.UUID == uuid {
            index = i
            break
        }
    }

    if index == -1 {
        return fmt.Errorf("package with UUID %s not found", uuid)
    }

    pm.Packages = append(pm.Packages[:index], pm.Packages[index+1:]...)
    return pm.Save()
}

