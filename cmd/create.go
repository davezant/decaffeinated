/*
Copyright © 2025 Davezant <dsndeividdsn1@gmail.com>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

var (
	appName   string
	appPrefix string
	appPath   string
	appExec   string
	appArgs   string
	appGroup  string
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create system resources (applications, users, and application groups)",
	Long: `The 'create' command allows you to create different types of resources in the system:

  - Applications: register applications with path, executable, and arguments.
  - Users: create new users who can interact with applications.
  - Application Groups: group multiple applications for easier management.

Examples:

  # Create an application
  decafein create app --name MyApp --path /usr/local/bin/ --executable myapp

  # Create a user
  decafein create user --name Alice

  # Create an application group
  decafein create group --name DevTools

This command organizes resource creation in a consistent and centralized way.`,
}

var createApp = &cobra.Command{
	Use:   "app",
	Short: "Create a new app",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.AddCommand(createApp)

	createApp.Flags().StringVarP(&appName, "name", "n", "", "Name of the application")
	createApp.Flags().StringVarP(&appPrefix, "prefix", "p", "", "Prefix for the application")
	createApp.Flags().StringVarP(&appPath, "path", "d", "", "Path to the application")
	createApp.Flags().StringVarP(&appExec, "executable", "e", "", "Executable file of the application")
	createApp.Flags().StringVarP(&appGroup, "group", "g", "", "Add to a group automatically")
	createApp.Flags().StringVarP(&appArgs, "args", "a", "", "Arguments to pass to the application")

	createApp.MarkFlagRequired("name")
	createApp.MarkFlagRequired("executable")
}
