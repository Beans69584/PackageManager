package main

import (
	"fmt"
	"os"

	"github.com/Beans69584/PackageManager/cmd"
	"github.com/spf13/cobra"
)

func main() {
	// Check if the program is running with root privileges.
	// Package installation and system modifications typically require elevated permissions.
	if os.Geteuid() != 0 {
		fmt.Println("You need to have root privileges to run this program.")
		os.Exit(1)
	}

	// Define the root command for the CLI application using Cobra.
	// This command acts as the base for all subcommands like install, uninstall, and list.
	rootCmd := &cobra.Command{
		Use:   "packagemanager",
		Short: "A simple package manager",
		Long:  `PackageManager is a simple tool to install, uninstall, and manage software packages.`,
	}

	// Add subcommands to the root command.
	// These subcommands are defined in the 'cmd' package and handle specific package management tasks.
	rootCmd.AddCommand(cmd.InstallCmd)
	rootCmd.AddCommand(cmd.UninstallCmd)
	rootCmd.AddCommand(cmd.ListCmd)

	// Execute the root command, which parses the CLI input and invokes the appropriate subcommand.
	if err := rootCmd.Execute(); err != nil {
		// If an error occurs during command execution, print the error and exit with a non-zero status code.
		fmt.Println(err)
		os.Exit(1)
	}
}
