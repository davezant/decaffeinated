package processes

import (
	"log"
	"time"

	"github.com/shirou/gopsutil/v4/process"
)

var GlobalSnapshot = NewProcessesSnapshot()

func NewProcessesSnapshot() *ProcessesSnapshot {
	return &ProcessesSnapshot{}
}

func StartGlobalSnapshotUpdater() {
	log.Printf("[INFO] bucket: Starting global processes bucket updater")
	go func() {
		for {
			GlobalSnapshot.UpdateSnapshot()
			time.Sleep(Timeout)
		}
	}()
}

func (p *ProcessesSnapshot) UpdateSnapshot() {
	log.Printf("[DEBUG] bucket: Updating processes bucket")
	p.Processes = nil

	rawBucket, err := process.Processes()
	if err == process.ErrorNotPermitted {
		log.Printf("[ERROR] bucket: Permission denied: try run Decafein as administrator or with sudo")
		return
	}

	for _, proc := range rawBucket {
		name, err := proc.Name()
		if err == nil {
			p.Processes = append(p.Processes, name)
		}
	}

	if len(p.Processes) == 0 {
		log.Fatalf("[FATAL] bucket: No processes captured: insufficient permissions")
	}
}
