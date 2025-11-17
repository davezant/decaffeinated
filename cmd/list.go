package cmd

import (
	"fmt"

	"github.com/davezant/decafein/src/client/webclient"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List resources from the server",
}

var listApps = &cobra.Command{
	Use:   "apps",
	Short: "List apps being watched",
	Run: func(cmd *cobra.Command, args []string) {
		client := webclient.New("http://localhost:1337")

		apps, err := client.GetApps()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		if len(apps) == 0 {
			fmt.Println("No apps found.")
			return
		}

		for _, app := range apps {
			fmt.Printf("- %s (%s)\n", app.Name, app.Binary)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.AddCommand(listApps)
}
