package dtime

import (
	"log"
	"time"
)

const (
	InitTimestamp  = 0.05
	HalfTimestamp  = 0.5
	CloseTimestamp = 0.75
	EndTimestamp   = 0.9
)

type CallbackTimestamp struct {
	Timestamp float32
	Callback  func()
}

type DTicker struct {
	TickerName string
	Intervals  time.Duration
	Callback   func()
	Done       chan bool
}

type DTimer struct {
	TimerName string
	Duration  time.Duration
	Callback  func()
	Done       chan bool
}

type DLimit struct {
	LimitName      string
	TargetDuration time.Duration
	Timestamps     []CallbackTimestamp
	
	Done           chan bool
	Control        chan bool // true = resume, false = pause
	
	activeDuration time.Duration
}

func NewTicker(name string, tickInterval int, callback func()) *DTicker {
	return &DTicker{
		TickerName: name, 
		Intervals: time.Duration(tickInterval) * time.Second, 
		Callback: callback, 
		Done: make(chan bool, 1),
	}
}

func (d *DTicker) StartTicker() {
	t := time.NewTicker(d.Intervals)
	go func() {
		defer t.Stop()
		for {
			select {
			case <-t.C:
				log.Println("ticker - " + d.TickerName)
				d.Callback()
			case <-d.Done:
				return
			}
		}
	}()
}

func (d *DTicker) StopTicker() {
	d.Done <- true
}

func NewTimer(name string, seconds int, callback func()) *DTimer {
	return &DTimer{
		TimerName: name,
		Duration:  time.Duration(seconds) * time.Second,
		Callback:  callback,
		Done:      make(chan bool, 1),
	}
}

func (d *DTimer) StartTimer() {
	t := time.NewTimer(d.Duration)
	go func() {
		select {
		case <-t.C:
			log.Println("timer finished - " + d.TimerName)
			d.Callback()
		case <-d.Done:
			t.Stop()
			return
		}
	}()
}

func (d *DTimer) StopTimer() {
	d.Done <- true
}

func NewLimit(name string, totalSeconds int) DLimit {
	return DLimit{
		LimitName:      name,
		TargetDuration: time.Duration(totalSeconds) * time.Second,
		Done:           make(chan bool, 1),
		Control:        make(chan bool),
	}
}

func (d DLimit) SetCallbackTimestamps(timestamps []CallbackTimestamp){
	log.Println("setting new timestamps")
	d.Timestamps = timestamps
}

func (d DLimit) CreateTimestamp(timestamp float32, callback func()){
	d.Timestamps = append(d.Timestamps, CallbackTimestamp{Timestamp:timestamp, Callback: callback})
}

func (d DLimit) RemoveTimestamp(timestamp float32){
	var newCallbackTimestamp []CallbackTimestamp
	for _, dt := range d.Timestamps{
		if dt.Timestamp == timestamp{
			continue
		} else {
			newCallbackTimestamp = append(newCallbackTimestamp, dt)
		}
	}
	d.SetCallbackTimestamps(newCallbackTimestamp)
}


func (d *DLimit) StartLimit() {
	ticker := time.NewTicker(1 * time.Second)
	fired := make(map[float32]bool)
	isRunning := false // Começa pausado até receber o primeiro sinal

	go func() {
		defer ticker.Stop()
		for {
			select {
			case state := <-d.Control:
				isRunning = state
				status := "PAUSED"
				if isRunning { status = "RESUMED" }
				log.Printf("Limit %s is now %s", d.LimitName, status)

			case <-ticker.C:
				if !isRunning {
					continue
				}

				d.activeDuration += time.Second
				progress := float32(d.activeDuration.Seconds()) / float32(d.TargetDuration.Seconds())

				for _, ts := range d.Timestamps {
					if progress >= ts.Timestamp && !fired[ts.Timestamp] {
						log.Printf("[%s] Milestone %.2f reached", d.LimitName, ts.Timestamp)
						ts.Callback()
						fired[ts.Timestamp] = true
					}
				}

				if progress >= 1.0 {
					log.Printf("[%s] Limit reached, stopping...", d.LimitName)
					return
				}

			case <-d.Done:
				return
			}
		}
	}()
}

func (d *DLimit) Toggle(status bool) {
	d.Control <- status
}

func (d *DLimit) StopLimit() {
	d.Done <- true
}
