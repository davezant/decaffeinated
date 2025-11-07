/*
Copyright © 2025 davezant
*/
package main

import (
	"github.com/davezant/decafein/src/server/database"
	"github.com/davezant/decafein/src/server/processes"
	"github.com/davezant/decafein/src/server/webserver"
)

func main() {
	database.CreateApp("notepad", "notepad.exe", "", "", "", true)
	webserver.OpenServer("", true)
	usuario := database.NewUser("deivid", "amo dogs")
	print(usuario.Name)
	processes.CurrentSession, _ = usuario.Login("amo dogs")
	processes.LocalWatcher.Login(processes.CurrentSession)
}
