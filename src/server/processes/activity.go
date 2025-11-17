package processes

import (
	"log"
	"time"

	"github.com/davezant/decafein/src/server/tempo"
	"github.com/davezant/decafein/src/utils"
)

var Timeout = 5 * time.Second

type Activity struct {
	Name                 string
	ExecutionBinary      string
	IsUp                 bool
	Limit                time.Duration
	DisplayExecutionTime string
	IsCounting           bool
	executionTime        time.Duration
	timer                *tempo.SimpleTimer
}

func NewActivity(name, processBinary string) *Activity {
	log.Println("[INFO] activity: Registering new activity - '" + name + "'")
	return &Activity{
		Name:                 name,
		ExecutionBinary:      processBinary,
		IsUp:                 false,
		executionTime:        time.Duration(time.Duration.Seconds(0)),
		IsCounting:           true,
		DisplayExecutionTime: "",
		timer:                tempo.NewSimpleTimer(),
	}
}

func (a *Activity) CheckIsRunning() bool {
	GlobalSnapshot.UpdateSnapshot()
	for _, processName := range GlobalSnapshot.Processes {
		if utils.EqualIgnoreCase(processName, a.ExecutionBinary) {
			return true
		}
	}
	return false
}

func (a *Activity) Up() {
	if a.IsCounting {
		a.IsUp = true
		log.Println("[DEBUG] activity: '" + a.Name + "' UP")

		a.timer.Start(1*time.Second, func() {
			if a.IsUp {
				a.executionTime += 1 * time.Second
				a.DisplayExecutionTime = a.executionTime.String()
			}
		})

		go tempo.TickerTimer(
			Timeout,
			func() bool {
				return !a.CheckIsRunning()
			},
			func() {
				a.IsUp = false
				a.Down()
			},
		)
	}
}

func (a *Activity) Down() {
	if a.IsCounting {
		a.IsUp = false
		a.timer.Stop()

		log.Println("[DEBUG] activity: '" + a.Name + "' DOWN")

		go tempo.TickerTimer(
			Timeout,
			func() bool {
				return a.CheckIsRunning()
			},
			func() {
				a.IsUp = true
				a.Up()
			},
		)
	}
}
