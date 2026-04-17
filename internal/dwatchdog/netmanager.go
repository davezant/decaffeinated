package dwatchdog

import (
	"decaffeinated/internal/dtime"
	"decaffeinated/internal/hlnet"
	"errors"
	"fmt"
	"log"
	"time"
)

type NetConfig struct {
	BlockedIPS      map[string]bool
	BlockedHostnames map[string]bool
	Host             string
	Port             string
}

type NetManager interface {
	NewProxy(host, port string) error
	StartProxy() error
	StopProxy() error
	BlockIP(ip string) error
	BlockHostname(hostname string) error
}

type IPC interface {
	StartIPC(sockpath string) error
	StopIPC() error
}

func (w *Watchdog) applyIPCCommand(cmd hlnet.IPCCommand) error {
	if cmd.Category == "" {
		return errors.New("category (rule name) is required")
	}

	w.rulesMu.Lock()
	defer w.rulesMu.Unlock()

	// Helper para converter CustomTimestamps do IPC para CallbackTimestamps do DTime
	parseTimestamps := func(cts []hlnet.CustomTimestamp) []dtime.CallbackTimestamp {
		var ts []dtime.CallbackTimestamp
		for _, ct := range cts {
			callbackName := ct.Callback // Nome vindo do JSON (ex: "warn_user")
			ts = append(ts, dtime.CallbackTimestamp{
				Timestamp: ct.Timestamp,
				Callback: func() {
					w.handleCallback(cmd.Category, callbackName)
				},
			})
		}
		return ts
	}

	switch cmd.Action {
	case "add":
		// Verifica se já existe para evitar duplicatas na mesma categoria
		for _, r := range w.Rules {
			if r.RuleName == cmd.Category {
				return errors.New("rule already exists, use 'update'")
			}
		}
		newRule := NewRule(
			cmd.Category,
			cmd.AppNames,
			cmd.IsBlocked,
			time.Duration(cmd.TimeLimitSeconds)*time.Second,
			parseTimestamps(cmd.CustomTimestamps),
		)
		w.Rules = append(w.Rules, newRule)
	case "update":
		for i, r := range w.Rules {
			if r.RuleName == cmd.Category {
				w.Rules[i].limitControl.StopLimit()
				
				w.Rules[i] = NewRule(
					cmd.Category,
					cmd.AppNames,
					cmd.IsBlocked,
					time.Duration(cmd.TimeLimitSeconds)*time.Second,
					parseTimestamps(cmd.CustomTimestamps),
				)
				return nil
			}
		}
		return errors.New("rule not found")

	case "delete":
		for i, r := range w.Rules {
			if r.RuleName == cmd.Category {
				w.Rules[i].limitControl.StopLimit()
				w.Rules = append(w.Rules[:i], w.Rules[i+1:]...)
				return nil
			}
		}

	case "start", "stop":
		status := (cmd.Action == "start")
		for i, r := range w.Rules {
			if r.RuleName == cmd.Category {
				w.Rules[i].IsBlocked = status
				// Opcional: Se 'stop' deve resetar o timer, chame Toggle(false)
				w.Rules[i].limitControl.Toggle(status)
				return nil
			}
		}

	default:
		return fmt.Errorf("unknown action: %s", cmd.Action)
	}

	w.Monitor.RefreshCurrentProcesses()
	return nil
}

func (w *Watchdog) StartIPC(sockpath string) error {
	w.IPCConfig.Path = hlnet.DefaultLinuxSockPath
	
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

	w.ipcServer = server
	log.Println("IPC Server Starting at:", sockpath)
	return w.ipcServer.Start()
}

func (w *Watchdog) StopIPC() error {
	if w.ipcServer == nil {
		return nil
	}
	err := w.ipcServer.Stop()
	w.ipcServer = nil
	return err
}
