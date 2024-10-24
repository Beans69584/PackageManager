package pkg

import (
    "archive/tar"
    "compress/gzip"
    "fmt"
    "io"
    "os"
    "path/filepath"
)

// ExtractTarGz extracts a .tar.gz archive to the specified destination directory
func ExtractTarGz(archivePath, destDir string) error {
    // Open the archive file
    file, err := os.Open(archivePath)
    if err != nil {
        return fmt.Errorf("error opening archive: %v", err)
    }
    defer file.Close()

    // Create a gzip reader
    gzReader, err := gzip.NewReader(file)
    if err != nil {
        return fmt.Errorf("error creating gzip reader: %v", err)
    }
    defer gzReader.Close()

    // Create a tar reader
    tarReader := tar.NewReader(gzReader)

    // Ensure destination directory exists
    err = os.MkdirAll(destDir, os.ModePerm)
    if err != nil {
        return fmt.Errorf("error creating destination directory: %v", err)
    }

    // Iterate through the files in the archive
    for {
        header, err := tarReader.Next()

        if err == io.EOF {
            // End of archive
            break
        }

        if err != nil {
            return fmt.Errorf("error reading tar archive: %v", err)
        }

        // Determine the proper file path
        targetPath := filepath.Join(destDir, header.Name)

        switch header.Typeflag {
        case tar.TypeDir:
            // Create Directory
            if err := os.MkdirAll(targetPath, os.FileMode(header.Mode)); err != nil {
                return fmt.Errorf("error creating directory %s: %v", targetPath, err)
            }
        case tar.TypeReg:
            // Create File
            if err := os.MkdirAll(filepath.Dir(targetPath), os.ModePerm); err != nil {
                return fmt.Errorf("error creating directory for file %s: %v", targetPath, err)
            }

            outFile, err := os.Create(targetPath)
            if err != nil {
                return fmt.Errorf("error creating file %s: %v", targetPath, err)
            }

            // Copy file contents
            if _, err := io.Copy(outFile, tarReader); err != nil {
                outFile.Close()
                return fmt.Errorf("error writing to file %s: %v", targetPath, err)
            }

            outFile.Close()

            // Set file permissions
            if err := os.Chmod(targetPath, os.FileMode(header.Mode)); err != nil {
                return fmt.Errorf("error setting permissions for file %s: %v", targetPath, err)
            }

        default:
            fmt.Printf("Skipping unknown type: %v in %s\n", header.Typeflag, header.Name)
        }
    }

    return nil
}

