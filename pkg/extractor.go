package pkg

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// ExtractTarGz extracts a .tar.gz archive to the specified destination directory.
// It handles the creation of directories and files, sets appropriate permissions,
// and ensures that the extraction process is secure and efficient.
//
// Parameters:
//   - archivePath (string): The file system path to the .tar.gz archive.
//   - destDir (string): The destination directory where the archive will be extracted.
//
// Returns:
//   - error: An error object if the extraction fails, otherwise nil.
func ExtractTarGz(archivePath, destDir string) error {
	// Open the archive file for reading.
	file, err := os.Open(archivePath)
	if err != nil {
		return fmt.Errorf("error opening archive: %v", err)
	}
	defer file.Close()

	// Create a gzip reader to decompress the .tar.gz archive.
	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("error creating gzip reader: %v", err)
	}
	defer gzReader.Close()

	// Create a tar reader to read the decompressed archive contents.
	tarReader := tar.NewReader(gzReader)

	// Ensure that the destination directory exists; create it if necessary.
	err = os.MkdirAll(destDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating destination directory: %v", err)
	}

	// Iterate through each entry in the tar archive.
	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			// Reached the end of the archive.
			break
		}

		if err != nil {
			return fmt.Errorf("error reading tar archive: %v", err)
		}

		// Determine the full path for the current entry.
		targetPath := filepath.Join(destDir, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			// Handle directory entries by creating the directory.
			if err := os.MkdirAll(targetPath, os.FileMode(header.Mode)); err != nil {
				return fmt.Errorf("error creating directory %s: %v", targetPath, err)
			}

		case tar.TypeReg:
			// Handle regular file entries.
			// Ensure that the parent directory exists.
			if err := os.MkdirAll(filepath.Dir(targetPath), os.ModePerm); err != nil {
				return fmt.Errorf("error creating directory for file %s: %v", targetPath, err)
			}

			// Create the file with appropriate permissions.
			outFile, err := os.Create(targetPath)
			if err != nil {
				return fmt.Errorf("error creating file %s: %v", targetPath, err)
			}

			// Copy the file contents from the archive to the newly created file.
			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return fmt.Errorf("error writing to file %s: %v", targetPath, err)
			}

			// Close the file to flush the write buffer.
			outFile.Close()

			// Set the file permissions as specified in the archive header.
			if err := os.Chmod(targetPath, os.FileMode(header.Mode)); err != nil {
				return fmt.Errorf("error setting permissions for file %s: %v", targetPath, err)
			}

		default:
			// Skip any unknown file types and inform the user.
			fmt.Printf("Skipping unknown type: %v in %s\n", header.Typeflag, header.Name)
		}
	}

	return nil
}
