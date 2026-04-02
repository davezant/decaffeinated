/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"decaffeinated/internal/dwatchdog"
	"decaffeinated/internal/hlnet"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"github.com/spf13/cobra"
)

var watchdogIPCPath string

var watchdogCmd = &cobra.Command{
	Use:   "watchdog",
	Short: "Start the WatchDog with IPC enabled",
	Long:  "Starts the WatchDog background monitor with IPC endpoint so clients can manage rules in real time.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if watchdogIPCPath == "" {
			return errors.New("--ipc-path is required")
		}

		wd := dwatchdog.NewWatchDog(nil)
		wd.IPCConfig = hlnet.IPCConfig{Path: watchdogIPCPath}

		if err := wd.StartIPC(); err != nil {
			return fmt.Errorf("failed to start IPC server: %w", err)
		}
		defer wd.StopIPC()

		wd.Start()

		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		fmt.Printf("WatchDog started (ipc=%s). Press Ctrl-C to stop.\n", watchdogIPCPath)
		<-ctx.Done()

		fmt.Println("Shutting down WatchDog...")
		return wd.StopIPC()
	},
}

var serviceInstallCmd = &cobra.Command{

}

func init() {
	rootCmd.AddCommand(watchdogCmd)
	watchdogCmd.Flags().StringVar(&watchdogIPCPath, "ipc-path", "", "IPC socket/pipe path")
	}

