package processes

import (
	"time"

	"github.com/davezant/decafein/src/server/tempo"
)

type Activity struct {
	Name                 string
	ExecutionBinary      string
	IsUp                 bool
	Limit                time.Duration
	DisplayExecutionTime string
	IsCounting           bool
	executionTime        time.Duration
	timer                *tempo.SimpleTimer
}

type ActivitiesRegistry struct {
	Active   []*Activity
	Inactive []*Activity
}

type Session struct {
	UserID    string
	LoginTime time.Time
	Limit     time.Duration
	IsMinor   bool
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
