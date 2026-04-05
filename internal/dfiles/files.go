package dfiles

import (
	"os/exec"
)

// File Management

func SaveAppPath(appName string, execPath string, configPath string){

}

func LoadAppPaths(configPath string){

}

func RunProgram(execPath string){
	exec.Command(execPath)
}


