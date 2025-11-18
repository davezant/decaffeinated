package policies

import (
	"fmt"

	"github.com/gen2brain/beeep"
)

func init() {
	beeep.AppName = "Decaffeinated"
}
func PrintPassedTime(aplicationName string) {
	fmt.Println("[POLICIES]	the application", aplicationName, "has exceeded the allowed execution time limit.")
}
func PrintHalfTime(applicationName string) {
	fmt.Println("[POLICIES]\tthe application", applicationName, "has reached 50% of the allowed execution time.")
}
func PrintAlmostEndingTime(applicationName string) {
	fmt.Println("[POLICIES]\tthe application", applicationName, "is close to exceeding the execution time limit (90% reached).")
}

func NotifyPassedTime(applicationName string) {
	beeep.Notify("Time Limit Exceeded",
		fmt.Sprintf("The application %s has exceeded the allowed execution time.", applicationName),
		"")
}

func NotifyHalfTime(applicationName string) {
	beeep.Notify("50% Reached",
		fmt.Sprintf("The application %s has reached 50%% of the allowed execution time.", applicationName),
		"")
}

func NotifyAlmostEndingTime(applicationName string) {
	beeep.Notify("Almost Ending",
		fmt.Sprintf("The application %s is close to exceeding the execution time limit (90%% reached).", applicationName),
		"")
}
