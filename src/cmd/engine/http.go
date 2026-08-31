package main

import (
	"encoding/json"
	"log"
	"net/http"

	"sector-one/internal/acudp"
	"sector-one/internal/physics"
	"sector-one/internal/stream"
	"sector-one/ui"
)

func startHTTP(addr string, broadcaster *stream.Broadcaster, source string) {
	apiHandler := stream.Handler(broadcaster, func() map[string]any {
		return map[string]any{
			"ok":          true,
			"source":      source,
			"subscribers": broadcaster.Subscribers(),
		}
	})

	staticHandler := ui.Handler()

	mux := http.NewServeMux()
	mux.Handle("/api/", apiHandler)
	mux.Handle("/", staticHandler)

	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	go func() {
		log.Println("engine listening on", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Println("http:", err)
		}
	}()
}

func publishCar(broadcaster *stream.Broadcaster, car acudp.CarInfo, source string, ts int64) {
	b, err := json.Marshal(physics.FromCar(car, source, ts))
	if err != nil {
		return
	}
	broadcaster.Publish(b)
}
