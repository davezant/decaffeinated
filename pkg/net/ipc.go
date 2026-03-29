package net

import "net"

type IPCManager interface {
	MakeIpcChannelLinux()
	MakeIpcChannelWindows()
}

func MakeIpcChannelLinux(){
		
}

func MakeIpcChannelWindows(){

}


