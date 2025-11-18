package processes

import (
	"time"

	"github.com/davezant/decafein/src/server/tempo"
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
	timer                *tempo.SimpleTimer
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
	overlayTimer         *tempo.SimpleTimer
}
