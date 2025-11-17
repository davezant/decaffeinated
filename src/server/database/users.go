// database/user.go
package database

import (
	"fmt"
	"time"

	"github.com/davezant/decafein/src/server/processes"
)

func NewUser(name, password string) User {
	return User{Name: name, Password: password}
}

func (u *User) Login(password string) (*processes.Session, error) {
	if u.Password != password {
		return nil, fmt.Errorf("senha incorreta")
	}

	u.Session = &processes.Session{
		UserID:    u.Name,
		LoginTime: time.Now(),
		Limit:     time.Hour,
		IsMinor:   false,
	}

	u.LastLogged = time.Now()
	u.isLogged = true

	return u.Session, nil
}

func (u *User) Logoff() {
	if u.Session != nil {
		u.TimeWasted += u.Session.Elapsed()
	}

	u.Session = nil
	u.isLogged = false
}

func LoadDatabaseUsers() {

}
