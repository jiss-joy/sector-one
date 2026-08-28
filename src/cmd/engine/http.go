package main

import (
	"encoding/json"
	"log"
	"net/http"

	"sector-one/internal/acudp"
	"sector-one/internal/physics"
	"sector-one/internal/stream"
)

func startHTTP(broadcaster *stream.Broadcaster, source string) {
	srv := &http.Server{
		Addr: "127.0.0.1:8080",
		Handler: stream.Handler(broadcaster, func() map[string]any {
			return map[string]any{
				"ok":          true,
				"source":      source,
				"subscribers": broadcaster.Subscribers(),
			}
		}),
	}
	go func() {
		log.Println("sse listening on", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Println("http:", err)
		}
	}()
}

func publishCar(broadcaster *stream.Broadcaster, car acudp.CarInfo, source string) {
	b, err := json.Marshal(physics.FromCar(car, source))
	if err != nil {
		return
	}
	broadcaster.Publish(b)
}
