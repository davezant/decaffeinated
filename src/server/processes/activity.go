package processes

import (
	"log"
	"time"

	"github.com/davezant/decafein/src/server/policies"
	"github.com/davezant/decafein/src/server/tempo"
	"github.com/davezant/decafein/src/utils"
)

var Timeout = 5 * time.Second

func NewActivity(name, processBinary string) *Activity {
	log.Println("[INFO] activity: Registering new activity - '" + name + "'")
	return &Activity{
		Name:                 name,
		ExecutionBinary:      processBinary,
		IsUp:                 false,
		executionTime:        time.Duration(time.Duration.Seconds(0)),
		IsCounting:           true,
		DisplayExecutionTime: "",
		limitPassed:          false,
		timer:                tempo.NewSimpleTimer(),
		onPassedLimit:        func() { policies.NotifyPassedTime(name) },
		onAlmostEnding:       func() { policies.NotifyAlmostEndingTime(name) },
		onHalf:               func() { policies.NotifyHalfTime(name) },
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
			if !a.IsUp {
				return
			}

			a.executionTime += 1 * time.Second
			a.DisplayExecutionTime = a.executionTime.String()

			if a.Limit > 0 {

				half := a.Limit / 2
				almost := a.Limit - (a.Limit / 10) // 90% do tempo

				// --- METADE ---
				if a.onHalf != nil && a.executionTime == half {
					a.onHalf()
				}

				// --- QUASE ACABANDO (faltando 10%) ---
				if a.onAlmostEnding != nil && a.executionTime == almost {
					a.onAlmostEnding()
				}

				// --- PASSOU DO LIMITE ---
				if !a.limitPassed && a.executionTime >= a.Limit {
					a.limitPassed = true
					if a.onPassedLimit != nil {
						a.onPassedLimit()
					}
				}
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

func (a *Activity) detectPassedTime() {

}

func (a *Activity) detectAgeDiff() {

}
