package watchdog

import (
	"decaffeinated/internal/processes"
	"decaffeinated/internal/timers"
	"log"
	"slices"
	"sync"
	"time"
)

type Rule struct {
	RuleName     string
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

func NewRule(groupName string) Rule {
	return Rule{
		RuleName: groupName,
	}
}

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
		if c.RuleName != name {
			newRules = append(newRules, c)
		} else {
			continue
		}
	}
	w.Rules = newRules
}
// --- Gestão de Regras ---

func (w *Watchdog) BlockRule(name string) {
	w.rulesMu.Lock()
	defer w.rulesMu.Unlock()
	for i := range w.Rules {
		if w.Rules[i].RuleName == name {
			w.Rules[i].IsBlocked = true
		}
	}
}

func (w *Watchdog) UnblockRule(name string) {
	w.rulesMu.Lock()
	defer w.rulesMu.Unlock()
	for i := range w.Rules {
		if w.Rules[i].RuleName == name {
			w.Rules[i].IsBlocked = false
		}
	}
}

func (w *Watchdog) RenameRule(oldName string, newName string) {
	w.rulesMu.Lock()
	defer w.rulesMu.Unlock()
	for i := range w.Rules {
		if w.Rules[i].RuleName == oldName {
			w.Rules[i].RuleName = newName
		}
	}
}

func (w *Watchdog) ActivateRule(name string) {
	w.rulesMu.Lock()
	defer w.rulesMu.Unlock()
	for i := range w.Rules {
		if w.Rules[i].RuleName == name {
			w.Rules[i].active = true
			w.Rules[i].LimitControl.Toggle(true)
		}
	}
}

func (w *Watchdog) DeactivateRule(name string) {
	w.rulesMu.Lock()
	defer w.rulesMu.Unlock()
	for i := range w.Rules {
		if w.Rules[i].RuleName == name {
			w.Rules[i].active = false
			w.Rules[i].LimitControl.Toggle(false)
		}
	}
}

func (w *Watchdog) AddAppToRule(ruleName string, appName string) {
	w.rulesMu.Lock()
	defer w.rulesMu.Unlock()
	for i := range w.Rules {
		if w.Rules[i].RuleName == ruleName {
			if slices.Contains(w.Rules[i].AppsNames, appName){
				return
			}
			w.Rules[i].AppsNames = append(w.Rules[i].AppsNames, appName)
		}
	}
}

func (w *Watchdog) RemoveAppFromRule(ruleName string, appName string) {
	w.rulesMu.Lock()
	defer w.rulesMu.Unlock()
	for i := range w.Rules {
		if w.Rules[i].RuleName == ruleName {
			var updatedApps []string
			for _, name := range w.Rules[i].AppsNames {
				if name != appName {
					updatedApps = append(updatedApps, name)
				}
			}
			w.Rules[i].AppsNames = updatedApps
		}
	}
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
				if verbose { log.Printf("[Watchdog] Rule %s activated", rule.RuleName) }
			} else if !isRunning && rule.active {
				rule.LimitControl.Toggle(false)
				rule.active = false
				if verbose { log.Printf("[Watchdog] Rule %s deactivated", rule.RuleName) }
			}
		}
		w.rulesMu.RUnlock()
	}
}
