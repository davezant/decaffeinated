package timers

import (
	"log"
	"time"
)

type CallbackTimestamp struct {
	Timestamp float32
	Callback  func()
}

type TLimit struct {
	LimitName      string
	TargetDuration time.Duration
	Timestamps     []CallbackTimestamp
	
	Done           chan bool
	Control        chan bool 
	
	activeDuration time.Duration
}

func NewLimit(name string, totalSeconds int) *TLimit {
	return &TLimit{
		LimitName:      name,
		TargetDuration: time.Duration(totalSeconds) * time.Second,
		Done:           make(chan bool, 1),
		Control:        make(chan bool, 1), 
	}
}

func (d *TLimit) SetCallbackTimestamps(timestamps []CallbackTimestamp) {
	log.Printf("[Timers] Defining %d timestamps for %s", len(timestamps), d.LimitName)
	d.Timestamps = timestamps
}

func (d *TLimit) CreateTimestamp(timestamp float32, callback func()) {
	d.Timestamps = append(d.Timestamps, CallbackTimestamp{Timestamp: timestamp, Callback: callback})
}

func (d *TLimit) StartLimit() {
	ticker := time.NewTicker(1 * time.Second)
	fired := make(map[float32]bool)
	isRunning := false 

	go func() {
		defer ticker.Stop()
		for {
			select {
			case state := <-d.Control:
				isRunning = state

			case <-ticker.C:
				if !isRunning {
					continue
				}

				d.activeDuration += time.Second
				
				totalSecs := d.TargetDuration.Seconds()
				if totalSecs <= 0 { continue }
				
				progress := float32(d.activeDuration.Seconds()) / float32(totalSecs)

				for _, ts := range d.Timestamps {
					if progress >= ts.Timestamp && !fired[ts.Timestamp] {
						log.Printf("[%s] Milestone %.2f reached [%.2f]", d.LimitName, ts.Timestamp, progress)
						if ts.Callback != nil {
							ts.Callback()
						}
						fired[ts.Timestamp] = true
					}
				}

				if progress >= 1.0 {
					log.Printf("[%s] Reached Limit. Ending monitoring...", d.LimitName)
					return
				}

			case <-d.Done:
				return
			}
		}
	}()
}

func (d *TLimit) Toggle(status bool) {
	select {
	case d.Control <- status:
	default:
	}
}

func (d *TLimit) StopLimit() {
	select {
	case d.Done <- true:
	default:
	}
}
