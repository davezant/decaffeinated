/*
Copyright © 2025 davezant
*/
package main

import (
	"github.com/davezant/decafein/src/server/policies"
	"github.com/davezant/decafein/src/server/webserver"
)

func main() {
	policies.CreateBlockWall()
	webserver.OpenServer("", true)
	select {}
}
