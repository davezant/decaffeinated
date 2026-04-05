/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"decaffeinated/internal/hlnet"
	"github.com/spf13/cobra"
)

var (
	appNames         []string
	groupName        string
	timeLimitSeconds int64
	isBlocked        bool
)

var appsCmd = &cobra.Command{
	Use:   "apps",
	Short: "Manage apps and groups",
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new rule",
	Run: func(cmd *cobra.Command, args []string) {
		payload := hlnet.IPCCommand{
			Action:           "add",
			Category:         groupName,
			AppNames:         appNames,
			TimeLimitSeconds: timeLimitSeconds,
			IsBlocked:        isBlocked,
		}
		sendIPC(payload)
	},
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an existing rule",
	Run: func(cmd *cobra.Command, args []string) {
		payload := hlnet.IPCCommand{
			Action:           "update",
			Category:         groupName,
			AppNames:         appNames,
			TimeLimitSeconds: timeLimitSeconds,
			IsBlocked:        isBlocked,
		}
		sendIPC(payload)
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a rule",
	Run: func(cmd *cobra.Command, args []string) {
		payload := hlnet.IPCCommand{
			Action:   "delete",
			Category: groupName,
		}
		sendIPC(payload)
	},
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start blocking for a rule",
	Run: func(cmd *cobra.Command, args []string) {
		payload := hlnet.IPCCommand{
			Action:   "start",
			Category: groupName,
		}
		sendIPC(payload)
	},
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop blocking for a rule",
	Run: func(cmd *cobra.Command, args []string) {
		payload := hlnet.IPCCommand{
			Action:   "stop",
			Category: groupName,
		}
		sendIPC(payload)
	},
}

func init() {
	rootCmd.AddCommand(appsCmd)
	appsCmd.AddCommand(addCmd, updateCmd, deleteCmd, startCmd, stopCmd)

	appsCmd.PersistentFlags().StringVarP(&groupName, "group", "g", "", "Group or Category name")
	appsCmd.MarkPersistentFlagRequired("group")

	addCmd.Flags().StringSliceVarP(&appNames, "app-name", "a", []string{}, "Application names")
	addCmd.Flags().Int64VarP(&timeLimitSeconds, "limit", "l", 0, "Time limit in seconds")
	addCmd.Flags().BoolVarP(&isBlocked, "blocked", "b", false, "Block status")

	updateCmd.Flags().StringSliceVarP(&appNames, "app-name", "a", []string{}, "Application names")
	updateCmd.Flags().Int64VarP(&timeLimitSeconds, "limit", "l", 0, "Time limit in seconds")
	updateCmd.Flags().BoolVarP(&isBlocked, "blocked", "b", false, "Block status")
}

func sendIPC(cmd hlnet.IPCCommand) {
	client, _ := hlnet.NewClient(hlnet.IPCConfig{Path: hlnet.DefaultLinuxSockPath})
	client.SendCommand(cmd)
}
