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
		return errors.New("you need to at least have a group name or application name")
	}

	w.rulesMu.Lock()
	defer w.rulesMu.Unlock()

	switch cmd.Action {
	case "add":
		var timestamps []dtime.CallbackTimestamp
		for _, ct := range cmd.CustomTimestamps {
			timestamps = append(timestamps, dtime.CallbackTimestamp{
				Timestamp: ct.Timestamp,
				Callback:  nil,
			})
		}

		newRule := NewRule(
			cmd.Category,
			cmd.AppNames,
			cmd.IsBlocked,
			time.Duration(cmd.TimeLimitSeconds),
			timestamps,
		)
		w.Rules = append(w.Rules, newRule)

	case "update":
		found := false
		for i, r := range w.Rules {
			if r.RuleName == cmd.Category {
				var timestamps []dtime.CallbackTimestamp
				for _, ct := range cmd.CustomTimestamps {
					timestamps = append(timestamps, dtime.CallbackTimestamp{
						Timestamp: ct.Timestamp,
						Callback:  nil,
					})
				}
				w.Rules[i] = NewRule(cmd.Category, cmd.AppNames, cmd.IsBlocked, time.Duration(cmd.TimeLimitSeconds), timestamps)
				found = true
				break
			}
		}
		if !found {
			return errors.New("rule not found for update")
		}

	case "delete":
		for i, r := range w.Rules {
			if r.RuleName == cmd.Category {
				w.Rules = append(w.Rules[:i], w.Rules[i+1:]...)
				break
			}
		}

	case "start":
		for i, r := range w.Rules {
			if r.RuleName == cmd.Category {
				w.Rules[i].IsBlocked = true
				break
			}
		}

	case "stop":
		for i, r := range w.Rules {
			if r.RuleName == cmd.Category {
				w.Rules[i].IsBlocked = false
				break
			}
		}

	case "setting":
		fmt.Println("Applying settings")

	default:
		return errors.New("invalid action mode")
	}

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
