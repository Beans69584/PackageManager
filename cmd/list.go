package cmd

import (
    "fmt"
    "os"
    "path/filepath"
    "text/tabwriter"

    "github.com/spf13/cobra"
    "github.com/Beans69584/PackageManager/pkg"
)

var ListCmd = &cobra.Command{
    Use:   "list",
    Short: "List all installed packages",
    Run: func(cmd *cobra.Command, args []string) {
        // Define the installation directory
        packagesDir := "/usr/local/share/packagemanager"

        // Create PackageManager
        pm, err := pkg.NewPackageManager(filepath.Join(packagesDir, "packages.json"))
        if err != nil {
            fmt.Printf("Error initializing PackageManager: %v\n", err)
            os.Exit(1)
        }

        if len(pm.Packages) == 0 {
            fmt.Println("No packages installed.")
            return
        }

        // Setup tabwriter for formatted output
        w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
        fmt.Fprintln(w, "NAME\tINSTALL PATH\tEXECUTABLE")
        fmt.Fprintln(w, "----\t------------\t-----------")

        for _, p := range pm.Packages {
            fmt.Fprintf(w, "%s\t%s\t%s\n", p.Name, p.InstallPath, p.Executable)
        }

        w.Flush()
    },
}

