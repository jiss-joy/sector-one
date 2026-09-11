package stream

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Handler(broadcaster *Broadcaster, source string) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		healthHandler(w, r, broadcaster, source)
	})
	mux.HandleFunc("GET /api/telemetry", func(w http.ResponseWriter, r *http.Request) {
		telemetryHandler(w, r, broadcaster)
	})
	mux.HandleFunc("OPTIONS /api/telemetry", getOptions)
	mux.HandleFunc("OPTIONS /api/health", getOptions)
	return mux
}

func telemetryHandler(w http.ResponseWriter, r *http.Request, broadcaster *Broadcaster) {
	withCORS(w, r)
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Subscribe to events
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
}

func getOptions(w http.ResponseWriter, r *http.Request) {
	withCORS(w, r)
	w.WriteHeader(http.StatusNoContent)
}

func healthHandler(w http.ResponseWriter, r *http.Request, broadcaster *Broadcaster, source string) {
	withCORS(w, r)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"ok":          true,
		"source":      source,
		"subscribers": broadcaster.Subscribers(),
	})
}

func withCORS(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	if origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	} else {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
}
