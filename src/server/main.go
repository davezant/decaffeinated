package main

import (
	"time"

	"github.com/davezant/decafein/src/server/database"
	"github.com/davezant/decafein/src/server/processes"
	"github.com/davezant/decafein/src/server/webserver"
)

func main() {
	database.CreateApp("opera", "opera.exe", "", "", "", true, 2*time.Minute)
	user := database.NewUser("Deivid", "Santana")
	sess, _ := user.Login("Santana")
	processes.GlobalWatcher.Login(sess)
	webserver.OpenServer("", true)
}
