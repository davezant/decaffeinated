package dtime

import "time"

type CallbackTimestamp struct {
	Timestamp float32
	Callback func()
}

func CreateTicker(intervals time.Duration, callback func()){

}

func CreateTime(){

}

func CreateLimit(time time.Duration){

}
