package stream

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Handler(broadcaster *Broadcaster, health func() map[string]any) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		withCORS(w, r)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(health())
	})
	mux.HandleFunc("GET /api/telemetry", func(w http.ResponseWriter, r *http.Request) {
		withCORS(w, r)
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "streaming unsupported", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		ch := broadcaster.Subscribe()
		defer broadcaster.Unsubscribe(ch)

		for {
			select {
			case <-r.Context().Done():
				return
			case payload, ok := <-ch:
				if !ok {
					return
				}
				fmt.Fprintf(w, "data: %s\n\n", payload)
				flusher.Flush()
			}
		}
	})
	mux.HandleFunc("OPTIONS /api/telemetry", options)
	mux.HandleFunc("OPTIONS /api/health", options)
	return mux
}

func options(w http.ResponseWriter, r *http.Request) {
	withCORS(w, r)
	w.WriteHeader(http.StatusNoContent)
}

func withCORS(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	if origin != "" {
		// Allow any localhost origin (any port) for dev
		w.Header().Set("Access-Control-Allow-Origin", origin)
	} else {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
}
