package dtime

import (
	"log"
	"time"
)

// TIMERS AND TIME REGISTRATION

const (
	InitTimestamp  = 0.05
	HalfTimestamp  = 0.5
	CloseTimestamp = 0.75
	EndTimestamp   = 0.9
)

type CallbackTimestamp struct {
	Timestamp float32
	Callback func()
}

type DTicker struct {
	TickerName string
	Intervals time.Duration
	Callback func()
	Done chan bool
}

type DTimer struct{

}

type DLimit struct {

}

func NewTicker(name string, secondsInterval int, callback func()) DTicker{
	return DTicker{TickerName: name, Intervals: time.Duration(secondsInterval) * time.Second, Callback: callback}
}

func (d DTicker) StartTicker(){
	t := time.NewTicker(d.Intervals)
	
	go func(){
		for {
			select{
				case <- t.C:
					log.Println("ticker - " + d.TickerName)
					d.Callback()
				case <- d.Done:
					return
				}
		}
	}()
}

func StopTicker(){

}

func NewTimer(){

}

func StartTimer(){

}

func StopTimer(){

}

func NewLimit(){

}

func GetLimitInfo(){

}

func StopLimit(){

}
