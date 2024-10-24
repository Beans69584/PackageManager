package main

import (
	"fmt"
	"os"

	"github.com/Beans69584/PackageManager/cmd"
	"github.com/spf13/cobra"
)

func main() {
	if os.Geteuid() != 0 {
		fmt.Println("You need to have root privileges to run this program.")
		os.Exit(1)
	}

	rootCmd := &cobra.Command{
		Use:   "packagemanager",
		Short: "A simple package manager",
		Long:  `PackageManager is a simple tool to install, uninstall, and manage software packages.`,
	}

	rootCmd.AddCommand(cmd.InstallCmd)
	rootCmd.AddCommand(cmd.UninstallCmd)
	rootCmd.AddCommand(cmd.ListCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
