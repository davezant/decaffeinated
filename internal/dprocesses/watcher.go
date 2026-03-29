package dprocesses

import (
	"context"
	"fmt"
//	"log"
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
}

type Manager interface {
	IncludeCurrentProcesses()
	ManageApplications()
}

func NewMonitor() Monitor {
	return Monitor{
		BootTime: time.Now(),
		KnowProcessesByName: make(map[string]bool),
		DProcesses: make(map[*DProcess]int),
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

    // Usaremos um mapa temporário de nomes detectados NESTA rodada
    namesThisRun := make(map[string]bool)
    hasChanged := false

    for _, p := range currentPs {
        name, err := p.Name()
        if err != nil || name == "" {
            continue
        }

        // 1. Marca que o nome existe no sistema agora
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
    }

    // 3. Limpeza: Se o NOME não apareceu em nenhuma instância, ele foi fechado
    for name := range m.KnowProcessesByName {
        if !namesThisRun[name] {
            delete(m.KnowProcessesByName, name)
            fmt.Printf("❌ Software Fechado: %s\n", name)
            hasChanged = true
        }
    }

    if hasChanged {
        fmt.Printf("Atualizado em: %s\n", time.Since(timeStart))
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
			fmt.Println("Killing someting")
		}
	}
	return nil
}



