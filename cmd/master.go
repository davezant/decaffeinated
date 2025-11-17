/*
Copyright © 2025 Davezant <dsndeividdsn1@gmail.com>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// masterCmd represents the master command
var masterCmd = &cobra.Command{
	Use:   "master",
	Short: "Only for master admin users - High ceil node management",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("master called")
	},
}

func init() {
	if globalFlags.IsAdmin {
		rootCmd.AddCommand(masterCmd)
	}
}
