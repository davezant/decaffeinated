/*
Copyright © 2025 davezant
*/
package main

import (
	"time"

	"github.com/davezant/decafein/src/server/database"
	"github.com/davezant/decafein/src/server/policies"
	"github.com/davezant/decafein/src/server/webserver"
)

func main() {
	minecraft := database.CreateApp("Minecraft", "tlauncher.jar", "./", "", "", true, 2*time.Hour)
	policies.UpdateRunningBinaryByGui(minecraft.Activity)
	webserver.OpenServer("", true)
	select {}
}
