package dwatchdog

import (
	"decaffeinated/internal/dprocesses"
	"decaffeinated/internal/dtime"
	"decaffeinated/internal/hlnet"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net"
	"sync"
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

type IPCCommand struct {
	Action          string `json:"action"`
	AppName         string `json:"app_name"`
	TimeLimitSeconds int64  `json:"time_limit_seconds,omitempty"`
	IsBlocked       bool   `json:"is_blocked,omitempty"`
}

type WatchDog struct {
	Rules      map[string]*Rule
	rulesMu    sync.RWMutex
	NetConfig  *NetConfig
	Monitor    dprocesses.Monitor
	RefreshInterval time.Duration
	IPCConfig  hlnet.IPCConfig

	TimeOnStart time.Time
	TimeOnSession time.Time

	ipcServer  *hlnet.Server
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
	w.TimeOnSession = time.Now()
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
			changed, err := w.Monitor.RefreshCurrentProcesses()
			if err != nil {
				log.Println("Erro ao atualizar processos:", err)
				continue
			}
			if changed {
				w.rulesMu.RLock()
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
				w.rulesMu.RUnlock()
			}
	}
	}()
}

func (w *WatchDog) checkTimeSinceStart(){

}
func (w *WatchDog) sendIPCResponse(c net.Conn, status string, message string) {
	if c == nil {
		return
	}

	resp := hlnet.IPCResponse{Status: status, Message: message}
	payload, err := json.Marshal(resp)
	if err != nil {
		log.Println("sendIPCResponse marshal error:", err)
		return
	}

	_, err = c.Write(payload)
	if err != nil {
		log.Println("sendIPCResponse write error:", err)
	}
}

func (w *WatchDog) applyIPCCommand(cmd hlnet.IPCCommand) error {
	if cmd.AppName == "" {
		return errors.New("app_name is required")
	}

	w.rulesMu.Lock()
	defer w.rulesMu.Unlock()

	rule, exists := w.Rules[cmd.AppName]
	if cmd.Action == "remove" {
		if !exists {
			return errors.New("rule not found")
		}
		delete(w.Rules, cmd.AppName)
		return nil
	}

	limit := time.Duration(cmd.TimeLimitSeconds) * time.Second
	if limit == 0 {
		limit = time.Second
	}

	if !exists {
		rule = &Rule{AppName: cmd.AppName, TimeLimit: limit, IsBlocked: cmd.IsBlocked}
		rule.limitControl = dtime.NewLimit(rule.AppName, int(rule.TimeLimit.Seconds()), rule.Timestamps)
		rule.limitControl.StartLimit()
		w.Rules[cmd.AppName] = rule
	}

	// update existing rule
	rule.TimeLimit = limit
	rule.IsBlocked = cmd.IsBlocked

	if cmd.Action == "block" {
		rule.IsBlocked = true
	} else if cmd.Action == "unblock" {
		rule.IsBlocked = false
	}
	
	w.Monitor.RefreshCurrentProcesses()
	return nil
}

func (w *WatchDog) handleIPCConn(c net.Conn) {
	defer c.Close()
	body, err := io.ReadAll(c)
	if err != nil {
		w.sendIPCResponse(c, "error", "failed reading command")
		return
	}

	var cmd hlnet.IPCCommand
	if err := json.Unmarshal(body, &cmd); err != nil {
		w.sendIPCResponse(c, "error", "invalid json")
		return
	}

	if err := w.applyIPCCommand(cmd); err != nil {
		w.sendIPCResponse(c, "error", err.Error())
		return
	}

	w.sendIPCResponse(c, "ok", "command applied")
}

func (w *WatchDog) StartIPC() error {
	if w.IPCConfig.Path == "" {
		return errors.New("IPCPath required")
	}
	if w.ipcServer != nil {
		return errors.New("IPC already started")
	}

	handler := func(cmd hlnet.IPCCommand) (hlnet.IPCResponse, error) {
		if err := w.applyIPCCommand(cmd); err != nil {
			return hlnet.IPCResponse{Status: "error", Message: err.Error()}, nil
		}
		return hlnet.IPCResponse{Status: "ok", Message: "command applied"}, nil
	}

	server, err := hlnet.NewServer(w.IPCConfig, handler)
	if err != nil {
		return err
	}

	if err := server.Start(); err != nil {
		return err
	}

	w.ipcServer = server
	return nil
}

func (w *WatchDog) StopIPC() error {
	if w.ipcServer == nil {
		return nil
	}
	return w.ipcServer.Stop()
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

