/*
Copyright © 2025 Davezant <dsndeividdsn1@gmail.com>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/davezant/decafein/internal/client/api"
	"github.com/spf13/cobra"
)

type Flags struct {
	AdminToken string

	VerboseMode bool
	IsAdmin     bool
	AllowsX11   bool
}

var globalFlags = Flags{
	IsAdmin:     true,
	VerboseMode: false,
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "decafein",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

func Execute() {
	err := rootCmd.Execute()
	_, netErr := api.New("http://localhost:1337").GetWatcher()

	if netErr != nil {
		fmt.Println("SERVER OFFLINE ✕")
	} else {
		fmt.Println("SERVER ONLINE ✓")
	}
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
