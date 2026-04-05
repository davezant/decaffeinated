/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"decaffeinated/internal/dwatchdog"
	"decaffeinated/internal/hlnet"

	"github.com/spf13/cobra"
)

// servicesCmd represents the services command
var watchdogCmd = &cobra.Command{
	Use:   "watchdog",
	Short: "Start or install watchdog service",
	Long: `Watchdog`,
}

var runWatchdogCmd = &cobra.Command{
	Use: "run",
	Short:"Start Watchdog without service",
	RunE: func(cmd *cobra.Command, args []string) error {
		dw := dwatchdog.NewWatchog(2)
		dw.Start()
		dw.StartIPC(hlnet.DefaultLinuxSockPath)
		select {}
	},
}

var installWatchdogCmd = &cobra.Command{
	Use: "install",
	Short: "Install Watchdog service",
	RunE: func(cmd *cobra.Command, args []string) error {
		print()
		return nil
	},
}

var removeWatchdogCmd = &cobra.Command{
	Use: "delete",
	Short: "Delete Watchdog service",
	RunE: func(cmd *cobra.Command, args []string) error {
		print()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(watchdogCmd)
	watchdogCmd.AddCommand(runWatchdogCmd)
	watchdogCmd.AddCommand(installWatchdogCmd)
	watchdogCmd.AddCommand(removeWatchdogCmd)
}
