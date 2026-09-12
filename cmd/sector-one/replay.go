package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"sector-one/internal/acudp"
	"sector-one/internal/physics"
	"sector-one/internal/record"
	"sector-one/internal/stream"
	"time"
)

const (
	maxReplayGap = 200 * time.Millisecond
)

func runReplay(ctx context.Context, path string, replayRate float64, sm *physics.SpecManager, broadcaster *stream.Broadcaster) error {
	reader, err := record.NewReader(path)
	if err != nil {
		return err
	}
	defer reader.Close()
	log.Println("replaying", path)

	var prev int64
	carName := "Unknown Vehicle"

	for {
		if ctx.Err() != nil {
			return nil
		}
		rec, err := reader.Next()
		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Println("replay finished")
				return nil
			}
			return err
		}
		err = sleepDelta(ctx, prev, rec.Timestamp, replayRate)
		if err != nil {
			return nil
		}
		prev = rec.Timestamp

		switch rec.Kind {
		case record.KindHandshake:
			handleKindHandshake(rec.Payload, &carName)
		case record.KindCarInfo:
			handleKindCarInfo(rec.Payload, sm, broadcaster, carName, rec)
		default:
			log.Printf("replay: unknown kind %d, skipping", rec.Kind)
		}
	}
}

func handleKindHandshake(payload []byte, carName *string) {
	session, err := acudp.ParseHandshakeResponse(payload)
	if err != nil {
		log.Println("[replay] handshake error: ", err)
		return
	}
	*carName = session.CarName
	fmt.Printf("car=%q driver=%q track=%q config=%q\n",
		session.CarName, session.DriverName, session.TrackName, session.TrackConfig)
}

func handleKindCarInfo(payload []byte, sm *physics.SpecManager, broadcaster *stream.Broadcaster, carName string, rec record.Record) {
	car, err := acudp.ParseCarInfo(payload)
	if err != nil {
		log.Println("[replay] parsing error: ", err)
		return
	}

	// Auto calibrate
	maxLoad := float32(0)
	for _, load := range car.TyreLoad {
		if load > maxLoad {
			maxLoad = load
		}
	}
	source := "replay"
	sm.Update(carName, car.EngineRPM, maxLoad, car.EngineLimiter)
	// File timestamps are UnixNano (sleepDelta depends on that). Frame.ts is UnixMilli, same as live.
	publishCar(sm, broadcaster, car, source, rec.Timestamp/int64(time.Millisecond), carName)
}

func sleepDelta(ctx context.Context, prev, now int64, replayRate float64) error {
	if prev == 0 || replayRate <= 0 {
		return nil
	}
	duration := time.Duration(float64(now-prev) / replayRate)
	if duration <= 0 {
		return nil
	}
	if duration > maxReplayGap {
		duration = maxReplayGap
	}

	return sleep(ctx, duration)
}
