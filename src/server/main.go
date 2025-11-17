/*
Copyright © 2025 davezant
*/
package main

import (
	"time"

	"github.com/davezant/decafein/src/server/database"
	"github.com/davezant/decafein/src/server/ui"
	"github.com/davezant/decafein/src/server/webserver"
)

func main() {
	minecraft := database.CreateApp("Minecraft", "tlauncher.jar", "./", "", "", true, 2*time.Hour)
	ui.UpdateRunningBinaryByGui(minecraft.Activity)
	webserver.OpenServer("", true)
	select {}
}
