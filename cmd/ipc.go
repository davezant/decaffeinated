package cmd

import (
"decaffeinated/internal/hlnet"
"fmt"
"github.com/spf13/cobra"
)

var ipcPath string
var ipcAppName string
var ipcAction string
var ipcTimeLimit int64
var ipcBlocked bool

var ipcCmd = &cobra.Command{
	Use:   "ipc",
	Short: "Send IPC commands to the WatchDog daemon",
	Long:  "Manage WatchDog rules in real time via IPC (add/update/block/unblock/remove).",
	RunE: func(cmd *cobra.Command, args []string) error {
	if ipcPath == "" {
		return fmt.Errorf("--ipc-path is required")
	}
	if ipcAction == "" {
		return fmt.Errorf("--action is required")
	}
	if ipcAppName == "" {
		return fmt.Errorf("--app is required")
	}

	client, err := hlnet.NewClient(hlnet.IPCConfig{Path: ipcPath})
	if err != nil {
		return err
	}

	cmdPayload := hlnet.IPCCommand{
		Action:           ipcAction,
		AppName:          ipcAppName,
		TimeLimitSeconds: ipcTimeLimit,
		IsBlocked:        ipcBlocked,
	}

	resp, err := client.SendCommand(cmdPayload)
		if err != nil {
		return err
	}

	fmt.Printf("ipc response: %s - %s\n", resp.Status, resp.Message)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(ipcCmd)

	ipcCmd.Flags().StringVar(&ipcPath, "ipc-path", "", "IPC socket/pipe path")
	ipcCmd.Flags().StringVar(&ipcAppName, "app", "", "Application name")
	ipcCmd.Flags().StringVar(&ipcAction, "action", "", "Action: add, update, block, unblock, remove")
	ipcCmd.Flags().Int64Var(&ipcTimeLimit, "time-limit", 0, "Time limit in seconds")
	ipcCmd.Flags().BoolVar(&ipcBlocked, "blocked", false, "Set rule blocked state")
}
