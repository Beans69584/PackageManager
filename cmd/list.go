package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/Beans69584/PackageManager/pkg"
	"github.com/spf13/cobra"
)

// ListCmd represents the 'list' command for the PackageManager.
// It allows users to view all currently installed packages.
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all installed packages",
	Run: func(cmd *cobra.Command, args []string) {
		// Define the base directory where packages are installed.
		packagesDir := "/usr/local/share/packagemanager"

		// Initialise the PackageManager, which manages the tracking of installed packages.
		pm, err := pkg.NewPackageManager(filepath.Join(packagesDir, "packages.json"))
		if err != nil {
			// If there is an error initializing the PackageManager, inform the user and exit.
			fmt.Printf("Error initializing PackageManager: %v\n", err)
			os.Exit(1)
		}

		// Check if there are any packages installed.
		if len(pm.Packages) == 0 {
			fmt.Println("No packages installed.")
			return
		}

		// Set up a tab writer for formatted, aligned output in the terminal.
		// The tabwriter.Writer ensures that the columns are properly aligned.
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

		// Print the header row with column titles.
		fmt.Fprintln(w, "NAME\tINSTALL PATH\tEXECUTABLE")
		fmt.Fprintln(w, "----\t------------\t-----------")

		// Iterate over each installed package and print its details.
		for _, p := range pm.Packages {
			// Format each package's name, installation path, and executable path into the tabbed format.
			fmt.Fprintf(w, "%s\t%s\t%s\n", p.Name, p.InstallPath, p.Executable)
		}

		// Flush the writer to ensure all output is written to the terminal.
		w.Flush()
	},
}
