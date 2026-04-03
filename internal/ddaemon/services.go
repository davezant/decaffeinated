package ddaemon

import (
	"fmt"
	"github.com/kardianos/service"
)

// DaemonManager wraps the kardianos/service logic
type DaemonManager struct {
	Service service.Service
	Config  *service.Config
}

// program implements service.Interface
// This is the bridge between the OS service manager and your Go code
type program struct {
	exit    chan struct{}
	execute func() // The actual watchdog/server logic
}

func (p *program) Start(s service.Service) error {
	// Start the logic in a background goroutine so Start returns quickly
	go p.run()
	return nil
}

func (p *program) run() {
	// Execute the provided server/watchdog logic
	if p.execute != nil {
		p.execute()
	}
}

func (p *program) Stop(s service.Service) error {
	// Signal shutdown if necessary
	close(p.exit)
	return nil
}

// NewDaemonManager creates a manager. 
// 'work' is a function containing your infinite loop (the watchdog server).
func NewDaemonManager(name, displayName, description string, work func()) (*DaemonManager, error) {
	prg := &program{
		exit:    make(chan struct{}),
		execute: work,
	}

	svcConfig := &service.Config{
		Name:        name,
		DisplayName: displayName,
		Description: description,
		Option: service.KeyValue{
			"UserService" : true,
		},
	}

	s, err := service.New(prg, svcConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create service: %w", err)
	}

	return &DaemonManager{
		Service: s,
		Config:  svcConfig,
	}, nil
}

// --- Implementation Methods ---

func (dm *DaemonManager) InstallService() error {
	// Attempt to install. Note: May require sudo/admin depending on OS.
	return dm.Service.Install()
}

func (dm *DaemonManager) UninstallService() error {
	return dm.Service.Uninstall()
}

func (dm *DaemonManager) StartService() error {
	return dm.Service.Start()
}

func (dm *DaemonManager) StopService() error {
	return dm.Service.Stop()
}

func (dm *DaemonManager) Run() error {
	// This is used inside the actual service executable to block and wait
	return dm.Service.Run()
}
