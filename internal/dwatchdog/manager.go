package dwatchdog

import (
	"decaffeinated/internal/dprocesses"
	"decaffeinated/internal/dtime"

	"log"
	"time"
)

type Rule struct {
	AppName      string
	TimeLimit    time.Duration
	IsBlocked    bool
	Timestamps   []dtime.CallbackTimestamp
	OnLimitReach func() 
	
	limitControl *dtime.DLimit
	active       bool
}

type WatchDog struct {
	Rules   map[string]*Rule
	NetConfig *NetConfig
	Monitor dprocesses.Monitor
	RefreshInterval time.Duration
	IPCPath string
}

func NewWatchDog(rules []Rule) *WatchDog {
	ruleMap := make(map[string]*Rule)
	
	for i := range rules {
		r := &rules[i]
		r.limitControl = dtime.NewLimit(
			r.AppName, 
			int(r.TimeLimit.Seconds()), 
			r.Timestamps,
		)
		r.limitControl.StartLimit()
		ruleMap[r.AppName] = r
	}

	return &WatchDog{
		Rules:           ruleMap,
		Monitor:         dprocesses.NewMonitor(),
		RefreshInterval: time.Second,
	}
}

func (w *WatchDog) Start() {
	log.Println("Watchdog Iniciado")

	if w.NetConfig != nil {
		log.Println("Proxy Iniciado")
		// Proxy Iniciar
		// TODO
		// defer fechar
	}
	go func() {
		ticker := time.NewTicker(5 * w.RefreshInterval)
		defer ticker.Stop()

		for range ticker.C {
			err := w.Monitor.RefreshCurrentProcesses()
			if err != nil {
				log.Println("Erro ao atualizar processos:", err)
				continue
			}
		for name, rule := range w.Rules {
		
			if rule.IsBlocked {
				w.KillProcess(name)
				continue
			}

			currentlyRunning := w.isAppRunning(name)
			if currentlyRunning && !rule.active {
				rule.limitControl.Toggle(true) 
				rule.active = true
			} else if !currentlyRunning && rule.active {
				rule.limitControl.Toggle(false)
				rule.active = false
			}
		}
	}
	}()
}

func (w *WatchDog) isAppRunning(name string) bool {
	state, _ := dprocesses.GetStateByName(name)
	return state
}

func (w *WatchDog) KillProcess(name string) {
	for c := range w.Monitor.DProcesses{
		if c.Name == name {
			dprocesses.KillProcessByName(name)
		}
	}
}

func (w *WatchDog) StartIPC(){
//TODO
}

func (w *WatchDog) StopIPC(){
//TODO
}
