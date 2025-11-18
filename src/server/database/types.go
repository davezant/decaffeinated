package database

import (
	"time"

	"github.com/davezant/decafein/src/server/policies"
	"github.com/davezant/decafein/src/server/processes"
)

type App struct {
	Name              string
	RootBinaryPath    string
	BinaryName        string
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
	GroupName string `json:"groupName"`
	Apps      []App  `json:"apps"`
	Policy    Policy `json:"policies"`
}

type Policy struct {
	Limit time.Duration

	halfAction  policies.Action
	closeAction policies.Action
	onEndAction policies.Action
}

type AppConfig struct {
	Name            string
	BinaryInProcess string
	BinaryInPath    string
	Path            string
	PrefixCommand   string
	SuffixComand    string

	AllowMinors     bool
	MainGroupPolicy Policy
	CustomPolicy    Policy
}

func NewTemplateAppConfig(name, binary string) *AppConfig {
	return &AppConfig{
		Name:            name,
		BinaryInProcess: binary,
		BinaryInPath:    "",
		Path:            "",
		PrefixCommand:   "",
		SuffixComand:    "",
		AllowMinors:     false,
	}
}

func NewCompleteAppConfig(name, binaryProcess, binaryPath, path, prefix, suffix string, minors bool) *AppConfig {
	return &AppConfig{
		Name:            name,
		BinaryInProcess: binaryProcess,
		BinaryInPath:    binaryPath,
		Path:            path,
		PrefixCommand:   prefix,
		SuffixComand:    suffix,
		AllowMinors:     minors,
	}
}
