/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const text = `

`
var activismCmd = &cobra.Command{
	Use:   "activism",
	Short: "",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(text)
	},
}

func init() {
	rootCmd.AddCommand(activismCmd)
}
