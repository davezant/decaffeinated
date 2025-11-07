// processes/session.go
package processes

import "time"

type Session struct {
	UserID    string
	LoginTime time.Time
	Limit     time.Duration
	IsMinor   bool
}

var CurrentSession = newSession()

func newSession() *Session {
	return &Session{}
}
func (s *Session) Elapsed() time.Duration {
	return time.Since(s.LoginTime).Truncate(time.Second)
}

func (s *Session) Remaining() time.Duration {
	if s.Limit == 0 {
		return -1
	}
	return (s.Limit - s.Elapsed()).Truncate(time.Second)
}

func (s *Session) Expired() bool {
	if s.Limit == 0 {
		return false
	}
	return s.Elapsed() >= s.Limit
}
