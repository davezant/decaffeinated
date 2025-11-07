package webserver

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/davezant/decafein/src/server/database"
	"github.com/davezant/decafein/src/server/processes"
)

var watch = &processes.LocalWatcher

var groups = map[string]interface{}{
	"unlisted":       &database.Unlisted,
	"created_groups": &database.UserCreatedGroups,
}

func WatcherJSONHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(watch); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func DatabaseJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	group := r.URL.Query().Get("group")

	if group != "" {
		if g, ok := groups[group]; ok {
			json.NewEncoder(w).Encode(g)
			return
		}
		http.Error(w, "group not found", http.StatusNotFound)
		return
	}

	http.Error(w, "argument missing", http.StatusNotFound)
}

func HandleInfo(w http.ResponseWriter, r *http.Request) {

}

func OpenServer(commonSecret string, proxyEnable bool) {
	http.HandleFunc("/info", HandleInfo)
	http.HandleFunc("/database-json", DatabaseJSON)
	http.HandleFunc("/watcher-json", WatcherJSONHandler)

	go func() {
		if proxyEnable {
			StartProxyServer()
		}
	}()
	log.Println("[INFO] webserver: Server Open on localhost:1337")
	log.Fatal(http.ListenAndServe("localhost:1337", nil))
}
