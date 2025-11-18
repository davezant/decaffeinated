package processes

import (
	"log"
	"time"

	"github.com/davezant/decafein/src/server/tempo"
	"github.com/shirou/gopsutil/v4/process"
)

var GlobalWatcher = NewWatcher(GlobalSnapshot, GlobalRegistry, CurrentSession)

func NewWatcher(bucket *ProcessesSnapshot, soup *ActivitiesRegistry, session *Session) *Watcher {
	return &Watcher{
		ProcessesSnapshot:    bucket,
		ActivitiesUp:         soup,
		ActiveSession:        session,
		ServiceStartTime:     time.Now(),
		SessionExecutionTime: "",
		overlayTimer:         tempo.NewSimpleTimer(),
	}
}

func (w *Watcher) Start() {
	w.overlayTimer.Start(time.Second, func() {
		w.ProcessesSnapshot.UpdateSnapshot()
		if w.ActiveSession == nil {
			return
		}

		elapsed := time.Since(w.ServiceStartTime).Truncate(time.Second)
		w.SessionExecutionTime = elapsed.String()

		if w.ActiveSession.Expired() {
			log.Printf("[INFO] Sessão expirada para user %s, deslogando.", w.ActiveSession.UserID)
			w.Logout()
		}
	})
}

func (w *Watcher) Login(sess *Session) {
	w.ActiveSession = sess
	w.ServiceStartTime = time.Now()
	log.Printf("[INFO] Sessão iniciada para %s", sess.UserID)
}

func (w *Watcher) Logout() {
	if w.ActiveSession != nil {
		log.Printf("[INFO] Sessão finalizada para %s", w.ActiveSession.UserID)
		w.ActiveSession = nil
	}
}

func (w *Watcher) RegisterActivity(a *Activity) {
	log.Printf("[INFO] watcher: '%s' registered", a.Name)
	w.ActivitiesUp.Active = append(w.ActivitiesUp.Active, a)
}

func (w *Watcher) RemoveActivity(a *Activity) {
	filtered := []*Activity{}
	for _, act := range w.ActivitiesUp.Active {
		if act.Name != a.Name {
			filtered = append(filtered, act)
		} else {
			log.Printf("[INFO] watcher: popped activity '%s'", a.Name)
		}
	}
	w.ActivitiesUp.Active = filtered
}

func (w *Watcher) KillActivity(a *Activity) {
	if a == nil {
		return
	}

	proc, err := process.NewProcess(a.Pid)
	if err != nil {
		log.Printf("[ERROR] watcher: cannot find process %d for '%s': %v", a.Pid, a.Name, err)
		return
	}

	// Tenta matar com Terminate primeiro (mais amigável)
	if err := proc.Terminate(); err != nil {
		log.Printf("[WARN] watcher: terminate failed for '%s' (%d): %v", a.Name, a.Pid, err)

		// Se falhar, tenta Kill (SIGKILL)
		if err := proc.Kill(); err != nil {
			log.Printf("[ERROR] watcher: kill failed for '%s' (%d): %v", a.Name, a.Pid, err)
			return
		}
	}

	log.Printf("[INFO] watcher: activity '%s' (pid %d) terminated", a.Name, a.Pid)

	// Remove da lista
	w.RemoveActivity(a)
}
