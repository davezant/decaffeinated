package processes

import (
	"time"
)

type Activity struct {
	Name                 string
	ExecutionBinary      string
	DisplayExecutionTime string
	IsUp                 bool
	limitPassed          bool
	IsCounting           bool
	Limit                time.Duration
	executionTime        time.Duration
	timer                *SimpleTimer
	onHalf               func()
	onAlmostEnding       func()
	onPassedLimit        func()
}

type ActivitiesRegistry struct {
	Active   []*Activity
	Inactive []*Activity
}

type Session struct {
	UserID    string
	IsMinor   bool
	LoginTime time.Time
	Limit     time.Duration
	OnEnding  func()
}

type ProcessesSnapshot struct {
	Processes []string
}

type Watcher struct {
	ProcessesSnapshot    *ProcessesSnapshot
	ActivitiesUp         *ActivitiesRegistry
	ActiveSession        *Session
	ServiceStartTime     time.Time
	SessionExecutionTime string
	overlayTimer         *SimpleTimer
}
