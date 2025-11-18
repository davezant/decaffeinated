package api

import "time"

type App struct {
	Name          string        `json:"name"`
	Binary        string        `json:"binary"`
	Path          string        `json:"path"`
	Limit         time.Duration `json:"limit"`
	CommandPrefix string        `json:"command_prefix"`
	CommandSuffix string        `json:"command_suffix"`
	CanMinorsPlay bool          `json:"can_minors_play"`
}

type Group struct {
	GroupName string `json:"groupName"`
	Apps      []App  `json:"apps"`
}

type Watcher struct {
	ActiveSession interface{} `json:"activeSession"`
}

type Session struct {
	User string `json:"user"`
}
