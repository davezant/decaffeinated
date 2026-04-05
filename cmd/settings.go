/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// settingsCmd represents the settings command
var settingsCmd = &cobra.Command{
	Use:   "settings",
	Short: "Create and manage settings",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("settings called")
	},
}

func init() {
	rootCmd.AddCommand(settingsCmd)
}
