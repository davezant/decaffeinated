package tempo

import "time"

func TickerTimer(refreshRate time.Duration, checkCallback func() bool, responseCallback func()) {
	ticker := time.NewTicker(refreshRate)
	defer ticker.Stop()

	for range ticker.C {
		if checkCallback() {
			responseCallback()
			return
		}
	}
}

type SimpleTimer struct {
	isRunning bool
	ticker    *time.Ticker
	stopChan  chan bool
}

func NewSimpleTimer() *SimpleTimer {
	return &SimpleTimer{
		isRunning: false,
		stopChan:  make(chan bool),
	}
}

func (t *SimpleTimer) Start(interval time.Duration, onTick func()) {
	// já está rodando? não inicia de novo
	if t.isRunning {
		return
	}

	t.isRunning = true
	t.ticker = time.NewTicker(interval)

	go func() {
		for {
			select {
			case <-t.ticker.C:
				onTick()
			case <-t.stopChan:
				t.ticker.Stop()
				return
			}
		}
	}()
}

func (t *SimpleTimer) Stop() {
	if !t.isRunning {
		return
	}
	t.isRunning = false
	t.stopChan <- true
}
