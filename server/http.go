package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/ChaitanyaSai-Meka/Taskdispatcher/internal/dispatcher"
)

func StartHTTP(addr string, d *dispatcher.Dispatcher) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/status/", func(w http.ResponseWriter, r *http.Request) {
		idStr := strings.TrimPrefix(r.URL.Path, "/status/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "invalid task id", http.StatusBadRequest)
			return
		}

		task, ok := d.Status(id)
		if !ok {
			http.Error(w, "task not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(task)
	})

	return http.ListenAndServe(addr, mux)
}
