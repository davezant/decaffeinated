/*
Copyright © 2025 Davezant <dsndeividdsn1@gmail.com>
*/
package cmd

import (
	"fmt"

	"github.com/davezant/decafein/src/client/webclient"
	"github.com/spf13/cobra"
)

// removeCmd representa o comando principal "remove"
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove resources from the server",
}

// removeAppCmd remove um app específico
var removeAppCmd = &cobra.Command{
	Use:   "app [name]",
	Short: "Remove an app by name",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		name := args[0]
		client := webclient.New("http://localhost:1337")

		err := client.DeleteApp(name)
		if err != nil {
			fmt.Println("Error removing app:", err)
			return
		}

		fmt.Printf("App '%s' removed successfully.\n", name)
	},
}

func init() {
	if globalFlags.IsAdmin {
		rootCmd.AddCommand(removeCmd)
		removeCmd.AddCommand(removeAppCmd)
	}
}
