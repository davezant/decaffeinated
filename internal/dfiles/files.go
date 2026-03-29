package dfiles

import (
	"os/exec"
)

// File Management
func RunProgram(execPath string){
	exec.Command(execPath)
}
