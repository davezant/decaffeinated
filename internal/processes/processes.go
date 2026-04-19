package processesmanager

import (

	"time"
	"github.com/shirou/gopsutil/v4/process"
)

type SmallProcess struct {
	Name string
	Filename string
} 

type Monitor struct {
	BootTime time.Time
	Processes map[int32]*SmallProcess
	LenOfProcesses int
}

func NewMonitor() Monitor {
	return Monitor{
		BootTime: time.Now(),
		Processes: make(map[int32]*SmallProcess),
		LenOfProcesses: 0,
	}
}

func NewProcess(name string, filename string) SmallProcess{
	return SmallProcess{
		Name: name,
		Filename: filename,
	}
}

func (m *Monitor) RefreshCurrentProcesses() (bool, error) {
    currentPs, err := process.Processes()
    if err != nil { return false, err }

	newMap := make(map[int32]*SmallProcess)
	
    changed := false

    for _, p := range currentPs {
        pid := p.Pid
        if existing, ok := m.Processes[pid]; ok {
            newMap[pid] = existing
        } else {
            name, _ := p.Name()
            exe, _ := p.Exe()
            newMap[pid] = &SmallProcess{Name: name, Filename: exe}
            changed = true
        }
    }

    if len(newMap) != len(m.Processes) { changed = true }
    m.Processes = newMap
    return changed, nil
}

func (m *Monitor) KillPID(pid int32) error {
    p, err := process.NewProcess(pid)
    if err != nil {
        return err
    }
    return p.Kill()
}

func (m *Monitor) GetPidsByName(targetName string) []int32 {
	var pids []int32
	for pid, proc := range m.Processes {
		if proc.Name == targetName {
			pids = append(pids, pid)
		}
	}
	return pids
}
