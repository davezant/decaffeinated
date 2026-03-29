package dwatchdog

import (
	"decaffeinated/internal/dprocesses"
	"log"
)

// High Level Time Process Management

type Rules struct {

}

func NewRules() {

}

func NewWatchDog(){

}

func StartWatchDog(){
	monitor := dprocesses.NewMonitor()
	err := monitor.RefreshCurrentProcesses()
	
	if err != nil {
		log.Println(err)
	}

	go func(){
		// Monitor Refresh Loop
		// Blocking Loop
		// Counting Apps
		// Measures Apply
	}()
}
