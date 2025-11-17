package database

import (
	"time"

	"github.com/davezant/decafein/src/server/processes"
)

type App struct {
	Name string

	RootBinaryPath string
	BinaryName     string

	commandLinePrefix string
	commandLineSuffix string

	Activity      *processes.Activity
	CanMinorsPlay bool
}

type User struct {
	Name           string
	Password       string
	LastLogged     time.Time
	TimeWasted     time.Duration
	TimeUntilReset time.Duration

	isLogged bool
	Session  *processes.Session
}

type Group struct {
	GroupName string        `json:"groupName"`
	Apps      []App         `json:"apps"`
	TimeLimit time.Duration `json:"timeLimit"`
}
