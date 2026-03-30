package dprocesses

import (
	"context"
	"fmt"
	"log"
	"time"
	"github.com/shirou/gopsutil/v4/process"
)

// Raw Process Management

type DProcess struct {
	Name string
	Filename string
	OpenTime time.Time

	context context.Context
} 

type Monitor struct {
	BootTime time.Time
	DProcesses map[*DProcess]int
	KnowProcessesByName map[string]bool
	RawLen int
}

func NewMonitor() Monitor {
	return Monitor{
		BootTime: time.Now(),
		KnowProcessesByName: make(map[string]bool),
		DProcesses: make(map[*DProcess]int),
		RawLen: 0,
	}
}

func NewDProcess(name string, filename string) DProcess{
	return DProcess{
		Name: name,
		Filename: filename,
	}
}

func (m *Monitor) RefreshCurrentProcesses() (bool, error) {
    timeStart := time.Now()
    
    currentPs, err := process.Processes()
    if err != nil {
        return false, err
    }

    namesThisRun := make(map[string]bool)
    hasChanged := false

    // 1. Mapeia todos os processos atuais sem travas de "len"
    for _, p := range currentPs {
        name, err := p.Name()
        if err != nil || name == "" {
            continue
        }

        namesThisRun[name] = true

        // Se for um processo novo que não conhecemos
        if !m.KnowProcessesByName[name] {
            file, _ := p.Exe()
            if file != "" {
                m.KnowProcessesByName[name] = true
                f := NewDProcess(name, file)
                m.DProcesses[&f] = int(p.Pid)
                
                fmt.Printf("🚀 Novo Software: %s\n", name)
                hasChanged = true
            }
        }

    }

    // 2. Compara o que conhecíamos com o que rodou AGORA
    for name := range m.KnowProcessesByName {
        if !namesThisRun[name] {
            delete(m.KnowProcessesByName, name)
            fmt.Printf("❌ Software Fechado: %s\n", name)
            hasChanged = true
        }
    }

    if hasChanged {
        m.RawLen = len(currentPs) // Atualiza o contador apenas para referência
    }
	fmt.Printf("Atualizado em: %s\n", time.Since(timeStart))
    
	return hasChanged, nil
}

func GetLiteralFromStruct(proc *DProcess) (*process.Process, error) {
	ps, err := process.Processes()
	for _, p := range ps {
		name, _ := p.Name()

		if name == proc.Name {
			return p, nil
		}
	}
	return nil, err
}

func GetState(proc *DProcess) (bool, error) {
	var err error
	literal, err := GetLiteralFromStruct(proc)
	if err != nil {
		return literal.IsRunning()
	}
	return false, err
}

func GetStateByName(proc string)(bool, error){
	ps, _ := process.Processes()
	for _, p := range ps {
		name, _ := p.Name()
		if name == proc {
			return true, nil
		} else {
			continue
		}
	}
	return false, nil
}

func MakeStateChannel(proc DProcess) (chan bool, error){
	isRunning := make(chan bool)
	return isRunning, nil
}

func KillProcessesByName(n []string) {
	for _, c := range n {
		err := KillProcessByName(c)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
}

func KillProcessByName(n string) error{
	ps, _ := process.Processes()
	for _, p := range ps {
		name , _ := p.Name()
		if name == n{
			log.Printf("killing %s process", n)
			p.Kill()
		}
	}
	return nil
}



