package database

import "time"

type Group struct {
	GroupName string        `json:"groupName"`
	Apps      []App         `json:"apps"`
	TimeLimit time.Duration `json:"timeLimit"`
}

func newGroup(name string) *Group {
	// Garante que Apps nunca é nil
	return &Group{
		GroupName: name,
		Apps:      []App{},
	}
}

var Unlisted = newGroup("unlisted_apps")

func (g *Group) AddToGroup(app *App) {
	if g.Apps == nil {
		g.Apps = []App{}
	}
	g.Apps = append(g.Apps, *app)
}

type userCreatedAppGroups struct {
	Groups []*Group `json:"groups"`
}

var UserCreatedGroups = &userCreatedAppGroups{
	Groups: []*Group{},
}

func CreateGroup(name string) *Group {
	g := newGroup(name)
	UserCreatedGroups.Groups = append(UserCreatedGroups.Groups, g)
	return g
}

func AddTimeLimitToGroup(group *Group, limit time.Duration) {
	group.TimeLimit = limit
}
