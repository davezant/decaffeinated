package main

import (
	"github.com/davezant/decafein/internal/server/database"
	"github.com/davezant/decafein/internal/server/processes"
	"github.com/davezant/decafein/internal/server/webserver"
)

func main() {
	database.CreateApp(database.NewTemplateAppConfig("opera", "opera.exe"))
	user := database.NewUser("Deivid", "Santana")
	sess, _ := user.Login("Santana")
	processes.GlobalWatcher.Login(sess)
	webserver.OpenServer("", true)
}
