package dwatchdog

import (
	"decaffeinated/internal/dtime"
	"decaffeinated/pkg/measures"
	"time"
)



var (
	NotifyOnlyTimestamps = []dtime.CallbackTimestamp{
		{
			Timestamp: dtime.InitTimestamp, 
			Callback: func() {measures.Notification("Using", time.Now())},
	},
		{
			Timestamp: dtime.HalfTimestamp, 
			Callback: func() {measures.Notification("Using", time.Now())},
	},
	}
	BlockOnEndTimestamps = []dtime.CallbackTimestamp{}
	OnlyLogTimestamps = []dtime.CallbackTimestamp{}
)

func EasyRestrictMode(){

}

func EasyPersonalMode(){

}


