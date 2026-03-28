package measures

import (
	"fmt"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/shirou/gopsutil/v4/process"
)

type BuiltinPunish interface{
	Popup()
	Warn()
	Notification()
	ScreenBlocker()
	MakeItShutdown()
}

func QuickPopup() {

}

func Warn(){

}

func Notification(appName string, time time.Time) {
	str := fmt.Sprintf("You are using %d for %t, Be careful!")
	beeep.Notify(str)
}

func Kill(pid int32) error{
	ps, err := process.Processes()
	if err != nil {
		return err
	}

	for _, p := range ps{
		if p.Pid == pid {
			e := p.Kill()
			if e != nil {
				return e	
			}
		}
	}
	return nil
}

func ScreenBlocker(){

}

func MakeItShutdown(){

}
