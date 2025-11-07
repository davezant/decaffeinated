/*
Copyright © 2025 Davezant <dsndeividdsn1@gmail.com>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// analyticsCmd represents the analytics command
var analyticsCmd = &cobra.Command{
	Use:   "analytics",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("analytics called")
	},
}

func init() {
	rootCmd.AddCommand(analyticsCmd)
}
