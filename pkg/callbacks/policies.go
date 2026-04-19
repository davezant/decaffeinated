package callbacks

import (
	"fmt"
	"time"

	"github.com/gen2brain/beeep"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

var (
	
)
func QuickPopup() {
	beeep.Alert("Decaffeinated", "", nil)
}

func Warn(){

}

func Notification(appName string, time time.Time){
	str := fmt.Sprintf("You are using %s for %s, Be careful!", appName, time.String())
	beeep.Notify("Decaffeinated", str, icons.AVGames)
}

func ScreenBlocker(){

}

func MakeItShutdown(){

}
