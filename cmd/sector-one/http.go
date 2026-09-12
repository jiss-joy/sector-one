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
	apiHandler := stream.Handler(broadcaster, source)
	uiHandler := ui.Handler()

	mux := http.NewServeMux()
	mux.Handle("/api/", apiHandler)
	mux.Handle("/", uiHandler)

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		log.Println("sector-one engine listening on:", server.Addr)
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Println("http: ", err)
		}
	}()
}

func publishCar(sm *physics.SpecManager, broadcaster *stream.Broadcaster, car acudp.CarInfo, source string, ts int64, carName, driverName string) {
	spec := sm.Get(carName)
	frame := physics.FromCar(car, source, ts, spec, carName, driverName)
	bytes, err := json.Marshal(frame)
	if err != nil {
		return
	}
	broadcaster.Publish(bytes)
}
