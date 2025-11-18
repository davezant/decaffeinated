package database

import (
	"log"

	"github.com/davezant/decafein/src/server/processes"
)

func CreateApp(a *AppConfig) App {
	activity := processes.NewActivity(a.Name, a.BinaryInProcess)
	activity.Limit = a.MainGroupPolicy.Limit

	application := App{
		Name:              a.Name,
		RootBinaryPath:    a.Path,
		BinaryName:        a.BinaryInPath,
		commandLinePrefix: a.PrefixCommand,
		commandLineSuffix: a.SuffixComand,
		Activity:          activity,
		CanMinorsPlay:     a.AllowMinors,
	}

	application.EnterInGroup(Unlisted)
	processes.GlobalWatcher.RegisterActivity(application.Activity)
	log.Println("[INFO] database: Creating process named '" + activity.Name + "' Executable : '" + activity.ExecutionBinary + "'")

	if activity.CheckIsRunning() {
		activity.Up()
	} else {
		activity.Down()
	}

	return application
}

func (application *App) EnterInGroup(group *Group) {
	application.Activity.Limit = group.Policy.Limit
	group.AddToGroup(application)
}

func (application *App) Remove() {
	processes.GlobalWatcher.RemoveActivity(application.Activity)
	log.Println("[INFO] database: Removing process named '" + application.Name + "'")
	application.Activity = nil
	application = nil
}
