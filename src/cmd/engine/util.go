package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"sector-one/internal/acudp"
	"sector-one/internal/record"
)

func logCar(car acudp.CarInfo, lastLog *time.Time, packets *int) {
	if lastLog.IsZero() {
		*lastLog = time.Now()
	}
	*packets++
	elapsed := time.Since(*lastLog)
	if elapsed < logEvery {
		return
	}
	hz := float64(*packets) / elapsed.Seconds()
	fmt.Printf("speed_kmh=%.1f gear=%d rpm=%.0f  (%.0f Hz ingest)\n",
		car.SpeedKmh, car.Gear, car.EngineRPM, hz)
	*packets = 0
	*lastLog = time.Now()
}

func sleep(ctx context.Context, duration time.Duration) error {
	t := time.NewTimer(duration)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}

func writeRec(rec *record.Writer, kind uint8, payload []byte) {
	if rec == nil {
		return
	}
	if err := rec.Write(kind, payload); err != nil {
		log.Println("record write:", err)
	}
}
