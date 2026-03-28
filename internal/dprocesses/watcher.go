package dprocesses

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/shirou/gopsutil/v4/process"
)

type RuleMonitoring interface {

}

type ProcessManagement interface {

}

type DProcess struct {
	Name string
	Filename string
	
	IsOn bool
	OpenTime time.Time

	context context.Context
}

type Monitor struct {
	BootTime time.Time
	DProcesses []DProcess
}

type Rule struct {
	TimeLimit time.Time
	OpenLimit int
	
	OnOpenDo func()
	OnHalfDo func()
	OnCloseDo func()
	OnEndDo func()
}

type AppRule struct {
	Name string
	Rule
}

type CategoryRule struct {
	Name []string
	Rule
}

type WatcherRules struct {
	AppRules []AppRule
	CategoryRules []CategoryRule
}

type Manager interface {
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

func (m Monitor) IncludeCurrentProcesses() error {
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

func SaveRunningState(proc DProcess) error{
	var err error
	literal, err := GetLiteralFromStruct(proc)
	if err != nil {
		proc.IsOn, err = literal.IsRunning()
	}
	return err
}

func StartProcess(proc DProcess) {

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
		log.Println("killed -", n.Name)
		p.Kill()
	}
	return nil
}

func GetTimeRunnedFromProcess(){

}

/*
func (m Monitor) 

func (m Monitor) StartWatcher(rules WatcherRules){
	err := m.IncludeCurrentProcesses()
	// TODO Create Ticker
	func (){
		err0 := m.IncludeCurrentProcesses() 
		for p := range m.DProcesses{
			
		}
	}()
}*/
