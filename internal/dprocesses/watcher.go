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

func (m *Monitor) RefreshCurrentProcesses() error {
    timeStart := time.Now()
    
    currentPs, err := process.Processes()
    if err != nil {
        return err
    }

    namesThisRun := make(map[string]bool)
    hasChanged := false

	if len(currentPs) > m.RawLen + 4 || len(currentPs) < m.RawLen - 4{
		m.RawLen = len(currentPs)
    	fmt.Printf("Novo valor de RawLen %d", m.RawLen)
		for _, p := range currentPs {
        	name, err := p.Name()
        	if err != nil || name == "" {
        	    continue
        	}

			namesThisRun[name] = true

        	// 2. Se o NOME não estiver no cache global, é um software NOVO
        	if !m.KnowProcessesByName[name] {
        	    file, _ := p.Exe() // Só chama o Exe() se for um nome novo (mais rápido)
            
        	    if file != "" {
        	        m.KnowProcessesByName[name] = true
        	        
        	        f := NewDProcess(name, file)
        	        m.DProcesses[&f] = int(p.Pid)
                
        	        fmt.Printf("🚀 Novo Software: %s\n", name)
        	        hasChanged = true
        	    }
        	}

    	for name := range m.KnowProcessesByName {
    	    if !namesThisRun[name] {
    	        delete(m.KnowProcessesByName, name)
    	        fmt.Printf("❌ Software Fechado: %s\n", name)
    	        hasChanged = true
    	    }
    	}
	}
	if hasChanged {
        fmt.Printf("Atualizado em: %s\n", time.Since(timeStart))
    	}
	}
    return nil
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



