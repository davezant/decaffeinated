package dwatchdog

import (
	"decaffeinated/internal/dprocesses"
	"decaffeinated/internal/dtime"
	"decaffeinated/internal/hlnet"
	"log"
	"sync"
	"time"
	"fmt"
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

	AvailableCallbacks map[string]func(string)
}

func NewRule(name string, appNames []string, isBlocked bool, timeLimit time.Duration, timestamps []dtime.CallbackTimestamp) Rule {
	// Inicializa o controle de limite de tempo
	limit := dtime.NewLimit(name, int(timeLimit.Seconds()))
	limit.SetCallbackTimestamps(timestamps)
	
	// Se a regra já deve começar ativa/bloqueada, você decide se inicia o timer aqui
	limit.StartLimit() 

	return Rule{
		RuleName:     name,
		AppNames:     appNames,
		TimeLimit:    timeLimit,
		IsBlocked:    isBlocked,
		Timestamps:   timestamps,
		limitControl: limit,
		active:       false,
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
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			changed, err := w.Monitor.RefreshCurrentProcesses()
			if err != nil {
				log.Println("error updating processes:", err)
				continue
			}

			if changed {
					for _, rule := range w.Rules {
						fmt.Println("executing rule"+rule.RuleName)
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
				}
			}
		}
	}()
}

func (w *Watchdog) handleCallback(category string, callbackName string) {
	log.Printf("[CALLBACK] Rule: %s, Action: %s", category, callbackName)
	
	switch callbackName {
	case "notify":
		fmt.Printf("ALERTA: Você atingiu um marco na categoria %s!\n", category)
	case "kill_all":
		w.rulesMu.RLock()
		for _, r := range w.Rules {
			if r.RuleName == category {
				for _, app := range r.AppNames {
					w.KillProcess(app)
				}
			}
		}
		w.rulesMu.RUnlock()
	default:
		log.Printf("Callback %s não implementado", callbackName)
	}
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
