package webserver

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/davezant/decafein/src/server/database"
	"github.com/davezant/decafein/src/server/processes"
)

var watch = processes.LocalWatcher

func respondJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	http.Error(w, message, status)
}

func WatcherHandler(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, watch)
}

func AppsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/apps"), "/")
		if len(parts) > 1 && parts[1] != "" {
			name := parts[1]
			for _, app := range database.Unlisted.Apps {
				if app.Name == name {
					respondJSON(w, app)
					return
				}
			}
			respondError(w, http.StatusNotFound, "app not found")
			return
		}
		respondJSON(w, database.Unlisted.Apps)
	case http.MethodPost:
		var data struct {
			Name          string `json:"name"`
			Binary        string `json:"binary"`
			Path          string `json:"path"`
			CommandPrefix string `json:"command_prefix"`
			CommandSuffix string `json:"command_suffix"`
			CanMinorsPlay bool   `json:"can_minors_play"`
		}
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		app := database.CreateApp(data.Name, data.Binary, data.Path, data.CommandPrefix, data.CommandSuffix, data.CanMinorsPlay)
		respondJSON(w, app)
	case http.MethodDelete:
		name := strings.TrimPrefix(r.URL.Path, "/apps/")
		for i, app := range database.Unlisted.Apps {
			if app.Name == name {
				app.Remove()
				database.Unlisted.Apps = append(database.Unlisted.Apps[:i], database.Unlisted.Apps[i+1:]...)
				respondJSON(w, map[string]string{"status": "deleted", "name": name})
				return
			}
		}
		respondError(w, http.StatusNotFound, "app not found")
	default:
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func GroupsHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/groups")
	path = strings.Trim(path, "/")

	switch r.Method {
	case http.MethodGet:
		if path != "" {
			if path == database.Unlisted.GroupName {
				respondJSON(w, database.Unlisted)
				return
			}
			for _, g := range database.UserCreatedGroups.Groups {
				if g.GroupName == path {
					respondJSON(w, g)
					return
				}
			}
			respondError(w, http.StatusNotFound, "group not found")
			return
		}

		resp := map[string]interface{}{
			"unlistedApps": database.Unlisted.Apps,
			"userCreated":  database.UserCreatedGroups.Groups,
		}
		respondJSON(w, resp)

	case http.MethodPut:
		if path != "" {
			respondError(w, http.StatusBadRequest, "invalid POST path")
			return
		}

		var payload struct {
			Name string `json:"name"`
		}

		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}

		if payload.Name == "" {
			respondError(w, http.StatusBadRequest, "missing name")
			return
		}

		group := database.CreateGroup(payload.Name)
		respondJSON(w, map[string]interface{}{
			"status": "created",
			"group":  group,
		})

	case http.MethodDelete:
		if path == "" {
			respondError(w, http.StatusBadRequest, "missing group name")
			return
		}

		for i, g := range database.UserCreatedGroups.Groups {
			if g.GroupName == path {
				database.UserCreatedGroups.Groups = append(
					database.UserCreatedGroups.Groups[:i],
					database.UserCreatedGroups.Groups[i+1:]...,
				)
				respondJSON(w, map[string]string{
					"status": "deleted",
					"group":  g.GroupName,
				})
				return
			}
		}
		respondError(w, http.StatusNotFound, "group not found")

	default:
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func SessionHandler(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, processes.CurrentSession)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	user := database.NewUser(creds.Name, creds.Password)
	session, err := user.Login(creds.Password)
	if err != nil {
		respondError(w, http.StatusUnauthorized, err.Error())
		return
	}
	processes.CurrentSession = session
	processes.LocalWatcher.Login(session)
	respondJSON(w, session)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if processes.CurrentSession != nil {
		processes.CurrentSession = nil
		watch.ActiveSession = nil
	}
	respondJSON(w, map[string]string{"status": "logged out"})
}

func OpenServer(commonSecret string, proxyEnable bool) {
	watch.Start()
	http.HandleFunc("/apps/", AppsHandler)
	http.HandleFunc("/groups", GroupsHandler)
	http.HandleFunc("/groups/", GroupsHandler)
	http.HandleFunc("/watcher/", WatcherHandler)
	http.HandleFunc("/session/", SessionHandler)
	http.HandleFunc("/login/", LoginHandler)
	http.HandleFunc("/logout/", LogoutHandler)

	if proxyEnable {
		go StartProxyServer()
	}

	addr := "localhost:1337"
	log.Printf("[INFO] webserver: Server open on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
