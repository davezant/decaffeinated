package control

import "decaffeinated/internal/watchdog"

type Manager struct {
	CurrentSessionHash string
	Watchdog *watchdog.Watchdog
}

func NewManager(intervalsInSeconds int) Manager{
	return Manager{
		Watchdog: watchdog.NewWatchdog(intervalsInSeconds),
	}
}

func (m *Manager) IsAdmin(userhash string){

}

func (m *Manager) LoginParameterstoHash(username string, password string) {

}

func (m *Manager) LoginUser(userhash string){
	// if userhash right
	m.CurrentSessionHash = userhash
}
