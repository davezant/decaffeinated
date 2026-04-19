package watchdog

import (
	"decaffeinated/internal/processes"
	"decaffeinated/internal/timers"
	"log"
	"sync"
	"time"
)

type Rule struct {
	GroupName     string
	AppsNames     []string
	IsBlocked    bool
	LimitControl timers.TLimit
	active       bool
}

type Watchdog struct {
	Rules           []Rule
	rulesMu         sync.RWMutex
	Monitor         processesmanager.Monitor
	RefreshInterval time.Duration
}

func NewWatchdog(intervalSeconds int) *Watchdog {
	return &Watchdog{
		Monitor:         processesmanager.NewMonitor(),
		RefreshInterval: time.Duration(intervalSeconds) * time.Second,
		Rules:           []Rule{},
	}
}

// Rules 

func (w *Watchdog) SetRules(rules []Rule){
	w.rulesMu.Lock()
	defer w.rulesMu.Unlock()
	w.Rules = rules
}

func (w *Watchdog) AddRule(rule Rule) {
	w.rulesMu.Lock()
	defer w.rulesMu.Unlock()
	w.Rules = append(w.Rules, rule)
}

func (w *Watchdog) RemoveRule(name string){
	w.rulesMu.Lock()
	defer w.rulesMu.Unlock()
	var newRules []Rule

	for _, c := range w.Rules {
		if c.GroupName != name {
			newRules = append(newRules, c)
		} else {
			continue
		}
	}
	w.Rules = newRules
}

func (w *Watchdog) BlockRule() {

}

func (w *Watchdog) UnblockRule() {

}

func (w *Watchdog) RenameRule() {

}

func (w *Watchdog) ActivateRule() {

}

func (w *Watchdog) DeactivateRule() {

}
// Apps 

func (w *Watchdog) AddAppToRule() {

}

func (w *Watchdog) RemoveAppFromRule() {

}

func (w *Watchdog) Start(verbose bool) {
	log.Println("[Watchdog] Started")
	ticker := time.NewTicker(w.RefreshInterval)
	defer ticker.Stop()

	for range ticker.C {
		w.Monitor.RefreshCurrentProcesses()
	
		w.rulesMu.RLock()
		for i := range w.Rules {
			rule := &w.Rules[i]
			
			var activePIDs []int32
			for _, name := range rule.AppsNames {
				activePIDs = append(activePIDs, w.Monitor.GetPidsByName(name)...)
			}

			isRunning := len(activePIDs) > 0

			if rule.IsBlocked && isRunning {
				for _, pid := range activePIDs {
					_ = w.Monitor.KillPID(pid)
				}
				continue
			}

			if isRunning && !rule.active {
				rule.LimitControl.Toggle(true)
				rule.active = true
				if verbose { log.Printf("[Watchdog] Rule %s activated", rule.GroupName) }
			} else if !isRunning && rule.active {
				rule.LimitControl.Toggle(false)
				rule.active = false
				if verbose { log.Printf("[Watchdog] Rule %s deactivated", rule.GroupName) }
			}
		}
		w.rulesMu.RUnlock()
	}
}
