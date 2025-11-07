package database

type Group struct {
	GroupName string
	Apps      []App
}

func newGroup(name string) *Group {
	return &Group{GroupName: name}
}

var Unlisted = newGroup("unlisted_apps")

func (g *Group) AddToGroup(app *App) {
	g.Apps = append(g.Apps, *app)
}

type userCreatedAppGroups struct {
	Groups []*Group
}

var UserCreatedGroups = &userCreatedAppGroups{}

func CreateGroup(name string) *Group {
	g := newGroup(name)
	UserCreatedGroups.Groups = append(UserCreatedGroups.Groups, g)
	return g
}
