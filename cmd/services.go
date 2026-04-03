package cmd

import (
	"decaffeinated/internal/ddaemon"
	"decaffeinated/internal/dwatchdog"
	"fmt"

	"github.com/spf13/cobra"
)

// Helper to get a manager instance with your app's specific config
func getManager() (*ddaemon.DaemonManager, error) {
    // This is the function that actually runs when the service starts
    watchdogLogic := func() {
        // Start your IPC server here
        // hlnet.StartServer(...) 
    }

    return ddaemon.NewDaemonManager(
        "decaffeinated",
        "Decaffeinated Watchdog",
        "Monitors application usage and handles IPC requests.",
        watchdogLogic,
    )
}

var startCmd = &cobra.Command{
	Use: "start",
	RunE: func(cmd *cobra.Command, args []string) error {
		wd := dwatchdog.NewWatchDog(nil)
		wd.Start()
		wd.IPCConfig.Path = "/tmp/deca-sock"
		wd.StartIPC()
		select {}
		return nil
	},
}
var serviceInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install the watchdog as a system service or windows task",
	RunE: func(cmd *cobra.Command, args []string) error {
		dm, err := getManager()
		if err != nil {
			return err
		}
		err = dm.InstallService()
		if err == nil {
			fmt.Println("Service installed successfully.")
		}
		return err
	},
}

var serviceUninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall the watchdog service",
	RunE: func(cmd *cobra.Command, args []string) error {
		dm, err := getManager()
		if err != nil {
			return err
		}
		err = dm.UninstallService()
		if err == nil {
			fmt.Println("Service uninstalled successfully.")
		}
		return err
	},
}

var serviceStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the watchdog service",
	RunE: func(cmd *cobra.Command, args []string) error {
		dm, err := getManager()
		if err != nil {
			return err
		}
		return dm.StartService()
	},
}

var serviceStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the watchdog service",
	RunE: func(cmd *cobra.Command, args []string) error {
		dm, err := getManager()
		if err != nil {
			return err
		}
		return dm.StopService()
	},
}

func init() {
	// You can group these under a "service" parent command
	var serviceCmd = &cobra.Command{
		Use:   "service",
		Short: "Manage the background watchdog service",
	}

	rootCmd.AddCommand(serviceCmd)
	rootCmd.AddCommand(startCmd)
	serviceCmd.AddCommand(serviceInstallCmd)
	serviceCmd.AddCommand(serviceUninstallCmd)
	serviceCmd.AddCommand(serviceStartCmd)
	serviceCmd.AddCommand(serviceStopCmd)
}
