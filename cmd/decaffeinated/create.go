/*
Copyright © 2025 Davezant <dsndeividdsn1@gmail.com>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/davezant/decafein/internal/client/api"
	"github.com/spf13/cobra"
)

var (
	createName   string
	createBinary string
	createPath   string
	createPrefix string
	createSuffix string
	createLimit  time.Duration
	createMinors bool
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create resources (apps, policies, groups...)",
	Long:  "Create resources. Example: decafein create app --name MyApp --binary mybin.exe ...",
}

var createAppCmd = &cobra.Command{
	Use:   "app",
	Short: "Create a new app watched by the server",
	Run: func(cmd *cobra.Command, args []string) {
		cli := api.New("http://localhost:1337")

		app, err := cli.CreateApp(createName, createBinary, createPath, createPrefix, createSuffix, createLimit, createMinors)
		if err != nil {
			log.Fatalf("[ERROR] create app: %v", err)
		}

		b, _ := json.MarshalIndent(app, "", "  ")
		fmt.Println(string(b))
	},
}

func init() {
	if globalFlags.IsAdmin {
		rootCmd.AddCommand(createCmd)
		createCmd.AddCommand(createAppCmd)

		createAppCmd.Flags().StringVar(&createName, "name", "", "App name (required)")
		createAppCmd.Flags().StringVar(&createBinary, "binary", "", "Binary filename (required)")
		createAppCmd.Flags().StringVar(&createPath, "path", "", "Binary root path")
		createAppCmd.Flags().StringVar(&createPrefix, "prefix", "", "Command line prefix")
		createAppCmd.Flags().StringVar(&createSuffix, "suffix", "", "Command line suffix")
		createAppCmd.Flags().DurationVar(&createLimit, "limit", 0, "Time limit for the app (e.g. 2h, 30m)")
		createAppCmd.Flags().BoolVar(&createMinors, "minors", false, "Allow minors to play this app")

		createAppCmd.MarkFlagRequired("name")
		createAppCmd.MarkFlagRequired("binary")
	}
}
