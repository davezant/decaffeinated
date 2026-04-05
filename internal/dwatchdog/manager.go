package dwatchdog

import (
	"decaffeinated/internal/dprocesses"
	"decaffeinated/internal/dtime"
	"decaffeinated/internal/hlnet"
	"log"
	"sync"
	"time"
)

// Structs 

// Rules are
type Rule struct {
	RuleName	 string
	AppNames     []string
	TimeLimit    time.Duration
	IsBlocked    bool
	Timestamps   []dtime.CallbackTimestamp

	limitControl dtime.DLimit
	active       bool
}

// Watchdog
type Watchdog struct {
	CurrentUser string
	
	Rules      []Rule
	rulesMu    sync.RWMutex
	
	NetConfig  *NetConfig
	Monitor    dprocesses.Monitor
	RefreshInterval time.Duration
	IPCConfig  hlnet.IPCConfig

	TimeOnStart time.Time
	TimeOnSession time.Time

	ipcServer  *hlnet.Server
}

func NewRule(name string, appName []string, isBlocked bool, timeLimit time.Duration, timestamps []dtime.CallbackTimestamp) Rule{
	return Rule{
		RuleName: name,
		AppNames: appName,
		TimeLimit: timeLimit * time.Second,
		Timestamps: timestamps,
		IsBlocked: isBlocked,
	}
}

// Watchdog

// Create a Watchdog
func NewWatchog(refreshIntervalInSeconds int) *Watchdog {

	return &Watchdog{
		Monitor:         dprocesses.NewMonitor(),
		RefreshInterval: time.Duration(refreshIntervalInSeconds) * time.Second,
	}
}

func (w *Watchdog) Start() {
	w.TimeOnSession = time.Now()
	log.Println("Watchdog: Started")

	if w.NetConfig != nil {
		log.Println("Proxy: Started")
		// TODO
	}
	go func(){
		ticker := time.NewTicker(w.RefreshInterval)
		defer ticker.Stop()
		for range ticker.C {
			changed, err := w.Monitor.RefreshCurrentProcesses()
			if err != nil {
				log.Println("error updating processes:", err)
				continue
			}

			if changed {
				w.rulesMu.RLock()
					for _, rule := range w.Rules {
						for _, apps := range rule.AppNames{
						if rule.IsBlocked {
							w.KillProcess(apps)
							continue
						}

						currentlyRunning := w.isAppRunning(apps)
						if currentlyRunning && !rule.active {
							rule.limitControl.Toggle(true)
							rule.active = true
						} else if !currentlyRunning && rule.active {
							rule.limitControl.Toggle(false)
							rule.active = false
						}
					}
					w.rulesMu.RUnlock()
				}
			}
		}
	}()
}

func (w *Watchdog) checkTimeSinceStart() time.Duration {
	return time.Since(w.TimeOnStart)
}

func (w *Watchdog) isAppRunning(name string) bool {
	state, _ := dprocesses.GetStateByName(name)
	return state
}

func (w *Watchdog) KillProcess(name string) {
	for c := range w.Monitor.DProcesses{
		if c.Name == name {
			dprocesses.KillProcessByName(name)
		}
	}
}
