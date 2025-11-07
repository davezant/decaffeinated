package database

import (
	"log"

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

func CreateApp(name, binary, path, commandPrefix, commandSuffix string, canMinorsPlay bool) App {
	activity := processes.NewActivity(name, binary)

	application := App{
		Name:              name,
		RootBinaryPath:    path,
		BinaryName:        binary,
		commandLinePrefix: commandPrefix,
		commandLineSuffix: commandSuffix,
		Activity:          activity,
		CanMinorsPlay:     canMinorsPlay,
	}

	application.EnterInGroup(Unlisted)
	processes.LocalWatcher.RegisterActivity(application.Activity)
	log.Println("[INFO] database: Creating process named '" + activity.Name + "' Executable : '" + activity.ExecutionBinary + "'")

	if activity.CheckIsRunning() {
		activity.Up()
	} else {
		activity.Down()
	}

	return application
}

func (application *App) EnterInGroup(group *Group) {
	group.AddToGroup(application)
}

func (application *App) Remove() {
	processes.LocalWatcher.RemoveActivity(application.Activity)
	log.Println("[INFO] database: Removing process named '" + application.Name + "'")
	application.Activity = nil
	application = nil
}
