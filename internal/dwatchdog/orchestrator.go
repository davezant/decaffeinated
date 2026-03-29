package dwatchdog

import (
	"decaffeinated/internal/dprocesses"
	"decaffeinated/internal/dtime"
	"fmt"
	"log"
	"time"
)

type Rule struct {
	AppName      string
	TimeLimit    time.Duration
	IsBlocked    bool
	Timestamps   []dtime.CallbackTimestamp
	OnLimitReach func() // Callback específico para quando o tempo esgotar
	
	// Controle interno
	limitControl *dtime.DLimit
	active       bool
}

type WatchDog struct {
	Rules   map[string]*Rule
	Monitor dprocesses.Monitor
	RefreshInterval time.Duration
}

func NewWatchDog(rules []Rule) *WatchDog {
	ruleMap := make(map[string]*Rule)
	
	for i := range rules {
		r := &rules[i]
		// Inicializa o motor de tempo para cada regra
		r.limitControl = dtime.NewLimit(
			r.AppName, 
			int(r.TimeLimit.Seconds()), 
			r.Timestamps,
		)
		// Inicia a goroutine do DLimit (que começa pausada)
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
	go func() {
		ticker := time.NewTicker(5 * w.RefreshInterval)
		defer ticker.Stop()

		for range ticker.C {
			// 1. Atualiza a lista de processos do SO
			err := w.Monitor.RefreshCurrentProcesses()
			if err != nil {
				log.Println("Erro ao atualizar processos:", err)
				continue
			}

			// 2. Verifica cada regra configurada
			for name, rule := range w.Rules {
				// Se a regra for de bloqueio imediato
				//w.PrintProcesses()
				if rule.IsBlocked {
					//print("executing rule")
					//checkAndKill(w.Monitor, name)
					continue
				}

				// Verifica se o processo está rodando agora
				currentlyRunning := w.isAppRunning(name)
				// 3. Gerenciamento de Estado (Pause/Resume)
				if currentlyRunning && !rule.active {
					rule.limitControl.Toggle(true) // Resume
					rule.active = true
				} else if !currentlyRunning && rule.active {
					rule.limitControl.Toggle(false) // Pause
					rule.active = false
				}
			}
		}
	}()
}

// isAppRunning verifica se o nome do processo consta na última varredura do monitor
func (w *WatchDog) isAppRunning(name string) bool {
	return true
}
func (w *WatchDog) PrintProcesses() {
	for c := range w.Monitor.DProcesses{
		fmt.Println(c.Name)
		if c.Name == "firefox" {
			dprocesses.KillProcessByName("firefox")
		}
	}
}
/*
func checkAndKill(m dprocesses.Monitor, targetName string) bool {
	for _, p := range m.DProcesses{
		print(p.Name)
		if p.Name == targetName{
			print("KILLIT")
			return true
		}
	} 
	return false
}*/
