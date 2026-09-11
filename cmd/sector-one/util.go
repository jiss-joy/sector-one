package main

import (
	"context"
	"fmt"
	"log"
	"sector-one/internal/acudp"
	"sector-one/internal/record"
	"time"
)

func logCar(car acudp.CarInfo, lastLog *time.Time, packets *int) {
	if lastLog.IsZero() {
		*lastLog = time.Now()
	}
	*packets++
	elapsed := time.Since(*lastLog)
	if elapsed < LogEvery {
		return
	}
	hz := float64(*packets) / elapsed.Seconds()
	fmt.Printf("speed_kmh=%.1f gear=%d rpm=%.0f  (%.0f Hz ingest)\n",
		car.SpeedKmh, car.Gear, car.EngineRPM, hz)
	*packets = 0
	*lastLog = time.Now()
}

func sleep(ctx context.Context, duration time.Duration) error {
	timer := time.NewTimer(duration)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func writeRecord(record *record.Writer, kind uint8, payload []byte) {
	if record == nil {
		return
	}
	err := record.Write(kind, payload)
	if err != nil {
		log.Println("Record write error: ", err)
	}
}
