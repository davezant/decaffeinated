package processes

import (
	"log"
	"time"

	"github.com/shirou/gopsutil/v4/process"
)

type ProcessesBucket struct {
	ProcessesStrings []string
}

var CommonBucket = NewProcessesBucket()

func NewProcessesBucket() *ProcessesBucket {
	return &ProcessesBucket{}
}

func StartGlobalBucketUpdater() {
	log.Printf("[INFO] bucket: Starting global processes bucket updater")
	go func() {
		for {
			CommonBucket.UpdateBucket()
			time.Sleep(Timeout)
		}
	}()
}

func (p *ProcessesBucket) UpdateBucket() {
	log.Printf("[DEBUG] bucket: Updating processes bucket")
	p.ProcessesStrings = nil

	rawBucket, err := process.Processes()
	if err == process.ErrorNotPermitted {
		log.Printf("[ERROR] bucket: Permission denied: run Decafein as administrator or with sudo")
		return
	}

	for _, proc := range rawBucket {
		name, err := proc.Name()
		if err == nil {
			p.ProcessesStrings = append(p.ProcessesStrings, name)
		}
	}

	if len(p.ProcessesStrings) == 0 {
		log.Fatalf("[FATAL] bucket: No processes captured: insufficient permissions")
	}
}
