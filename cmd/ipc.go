package cmd

import (
	"decaffeinated/internal/hlnet"
	"fmt"
	"github.com/spf13/cobra"
)

var (
	ipcPath      string
	ipcAppName   string
	ipcTimeLimit int64
	ipcBlocked   bool
)

var appCmd = &cobra.Command{
	Use:   "app",
	Short: "Manage applications via IPC",
}

func requestIPC(action string, cmd *cobra.Command, args []string) error {
	if ipcPath == "" {
		return fmt.Errorf("--ipc-path is required")
	}
	if ipcAppName == "" {
		return fmt.Errorf("--app is required")
	}

	client, err := hlnet.NewClient(hlnet.IPCConfig{Path: ipcPath})
	if err != nil {
		return fmt.Errorf("could not connect to IPC: %w", err)
	}

	cmdPayload := hlnet.IPCCommand{
		Action:           action,
		AppName:          ipcAppName,
		TimeLimitSeconds: ipcTimeLimit,
		IsBlocked:        ipcBlocked,
	}

	if err := client.SendCommand(cmdPayload); err != nil {
		return fmt.Errorf("IPC command failed: %w", err)
	}

	fmt.Printf("Action '%s' successfully sent for app: %s\n", action, ipcAppName)
	return nil
}

var addIPC = &cobra.Command{
	Use:   "add",
	Short: "Add a new application",
	RunE: func(cmd *cobra.Command, args []string) error {
		return requestIPC("add", cmd, args)
	},
}

var blockIPC = &cobra.Command{
	Use:   "block",
	Short: "Block a managed application",
	RunE: func(cmd *cobra.Command, args []string) error {
		return requestIPC("block", cmd, args)
	},
}

var updateIPC = &cobra.Command{
	Use:   "update",
	Short: "Update an application",
	RunE: func(cmd *cobra.Command, args []string) error {
		return requestIPC("update", cmd, args)
	},
}

var removeIPC = &cobra.Command{
	Use:   "remove",
	Short: "Remove an application",
	RunE: func(cmd *cobra.Command, args []string) error {
		return requestIPC("remove", cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(appCmd)

	appCmd.AddCommand(addIPC)
	appCmd.AddCommand(blockIPC)
	appCmd.AddCommand(updateIPC)
	appCmd.AddCommand(removeIPC)

	appCmd.PersistentFlags().StringVar(&ipcPath, "ipc-path", "", "Path to the IPC Unix socket")
	appCmd.PersistentFlags().StringVar(&ipcAppName, "app", "", "Name of the target application")
	appCmd.PersistentFlags().Int64Var(&ipcTimeLimit, "limit", 0, "Time limit in seconds")
	appCmd.PersistentFlags().BoolVar(&ipcBlocked, "blocked", false, "Set application to blocked state")

	appCmd.MarkPersistentFlagRequired("ipc-path")
	appCmd.MarkPersistentFlagRequired("app")
}
