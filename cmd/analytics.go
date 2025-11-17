/*
Copyright © 2025 Davezant
*/
package cmd

import (
	"fmt"

	"github.com/davezant/decafein/src/client/webclient"
	"github.com/spf13/cobra"
)

var analyticsCmd = &cobra.Command{
	Use:   "analytics",
	Short: "Analytics and detailed app information",
	Long:  `Provides detailed information and analytics about specific apps.`,
}

var analyticsShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show analytics details",
}

var analyticsShowAppCmd = &cobra.Command{
	Use:   "app [name]",
	Short: "Show analytics for a specific app",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		name := args[0]

		// Cria cliente HTTP
		client := webclient.New("http://localhost:1337")

		app, err := client.GetApp(name)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		// Exibição formatada
		fmt.Printf("\n=== Analytics for App: %s ===\n", app.Name)
		fmt.Printf("Binary:          %s\n", app.Binary)
		fmt.Printf("Path:            %s\n", app.Path)
		fmt.Printf("Prefix:          %s\n", app.CommandPrefix)
		fmt.Printf("Suffix:          %s\n", app.CommandSuffix)
		fmt.Printf("Limit:           %s\n", app.Limit)
		fmt.Printf("Minors Allowed:  %v\n", app.CanMinorsPlay)

		// Se existir histórico ou métricas, você coloca aqui depois
	},
}
var analyticsWatcherCmd = &cobra.Command{
	Use:   "watcher",
	Short: "Show full watcher status",
	Run: func(cmd *cobra.Command, args []string) {

		client := webclient.New("http://localhost:1337")

		watch, err := client.GetWatcher()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Println("\n=== Watcher Status ===")

		// ---- Processes Snapshot ----
		fmt.Println("\n--- Running Processes ---")
		for _, p := range watch.ProcessesSnapshot.Processes {
			fmt.Println("•", p)
		}

		// ---- Activities ----
		fmt.Println("\n--- Active Activities ---")
		for _, a := range watch.ActivitiesUp.Active {
			fmt.Printf("[ACTIVE] %s (%s) - Up: %v - Time: %s\n",
				a.Name, a.ExecutionBinary, a.IsUp, a.DisplayExecutionTime)
		}

		fmt.Println("\n--- Inactive Activities ---")
		for _, a := range watch.ActivitiesUp.Inactive {
			fmt.Printf("[INACTIVE] %s (%s)\n", a.Name, a.ExecutionBinary)
		}

		// ---- Active Session ----
		fmt.Println("\n--- Active Session ---")
		fmt.Printf("User ID: %s\n", watch.ActiveSession.UserID)
		fmt.Printf("Login:   %s\n", watch.ActiveSession.LoginTime)
		fmt.Printf("Limit:   %d\n", watch.ActiveSession.Limit)
		fmt.Printf("Minor:   %v\n", watch.ActiveSession.IsMinor)

		// ---- Service Time ----
		fmt.Println("\n--- Service Info ---")
		fmt.Printf("Service Start Time:  %s\n", watch.ServiceStartTime)
		fmt.Printf("Session Exec Time:   %s\n", watch.SessionExecutionTime)
	},
}

func init() {
	rootCmd.AddCommand(analyticsCmd)
	analyticsCmd.AddCommand(analyticsWatcherCmd)
	analyticsCmd.AddCommand(analyticsShowCmd)
	analyticsShowCmd.AddCommand(analyticsShowAppCmd)
}
