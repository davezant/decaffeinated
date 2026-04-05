/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const activismText = `


`
// activismCmd represents the activism command
var activismCmd = &cobra.Command{
	Use:   "activism",
	Short: "More information about donations for victims conflict",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(activismText)
	},
}

func init() {
	rootCmd.AddCommand(activismCmd)
}
