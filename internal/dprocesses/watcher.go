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
	DProcesses []DProcess
}

type Manager interface {
	IncludeCurrentProcesses()
	ManageApplications()
}

func NewMonitor() Monitor {
	return Monitor{
		BootTime: time.Now(),
	}
}

func NewDProcess(name string, filename string) DProcess{
	return DProcess{
		Name: name,
		Filename: filename,
	}
}

func (m Monitor) RefreshCurrentProcesses() error {
	ps, err := process.Processes()
	log.Println(m.BootTime)
	m.DProcesses = []DProcess{} 
	if err != nil {
		return err
	}

	for _, p := range ps {
		
		name, _ := p.Name()
		file, _ := p.Exe()
		if file != "" {
			log.Println("added - " + name + " to the pool > " + file)
			f := NewDProcess(name, file)
			m.DProcesses = append(m.DProcesses, f)
		}
	}
	return nil
}

func GetLiteralFromStruct(proc DProcess) (*process.Process, error) {
	ps, err := process.Processes()
	for _, p := range ps {
		name, _ := p.Name()
		if name == proc.Name {
			return p, nil
		}
	}
	return nil, err
}

func GetState(proc DProcess) (bool, error) {
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

func KillProcesses(n []DProcess) {
	for _, p := range n {
		err := KillProcess(p)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
}

func KillProcess(n DProcess) error{
	p, err := GetLiteralFromStruct(n)
	if err != nil{
		log.Println("killed - ", n.Name)
		p.Kill()
	}
	return nil
}


