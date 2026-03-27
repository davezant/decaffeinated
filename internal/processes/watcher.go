package processes

import "time"

type Process struct {
	UUID string
	Name string
	Path string
	Filename string
}

type Monitor struct {
	BootTime string
	Processes []Process
}

type Rule struct {
	TimeLimit time.Time
	OpenLimit int
	
	OnHalfDo func()
	OnCloseDo func()
	OnEndDo func()
}

type AppRule struct {
	UUID string
	Rule
}

type CategoryRule struct {
	UUIDs []string
	Rule
}

type WatcherRules struct {
	WatchedAppsUUID []string
	AppRules []AppRule
	SectionRules []CategoryRule
	ComputerTimeLimit time.Time
}

type Manager interface {
	ManageApplications()
}

func NewMonitor() Monitor {
	/*

	*/
	return Monitor{
		BootTime: time.Now().String(),
	}
}

func (m Monitor) CurrentProcesses(){

}

func ManageApplications(m Monitor, r WatcherRules){
	
}
