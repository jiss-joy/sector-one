package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"time"

	"sector-one/internal/acudp"
	"sector-one/internal/record"
	"sector-one/internal/stream"
)

// runReplay reads a .bin and feeds the same parsers as live. No UDP, no AC.
func runReplay(ctx context.Context, path string, rate float64, broadcaster *stream.Broadcaster, source string) error {
	rd, err := record.NewReader(path)
	if err != nil {
		return err
	}
	defer rd.Close()
	log.Println("replaying", path)

	var prev int64
	capped := false
	lastLog := time.Time{}
	packets := 0

	for {
		if ctx.Err() != nil {
			return nil
		}
		rec, err := rd.Next()
		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Println("replay finished")
				return nil
			}
			return err
		}
		if err := sleepDelta(ctx, prev, rec.Timestamp, rate, &capped); err != nil {
			return nil
		}
		prev = rec.Timestamp

		switch rec.Kind {
		case record.KindHello:
			session, err := acudp.ParseHandshakeResponse(rec.Payload)
			if err != nil {
				log.Println("replay handshake:", err)
				continue
			}
			fmt.Printf("car=%q driver=%q track=%q config=%q\n",
				session.CarName, session.DriverName, session.TrackName, session.TrackConfig)
		case record.KindCar:
			car, err := acudp.ParseCarInfo(rec.Payload)
			if err != nil {
				log.Println("skipping packet: ", err)
				continue
			}
			publishCar(broadcaster, car, source, rec.Timestamp)
			logCar(car, &lastLog, &packets)
		default:
			log.Printf("replay: unknown kind %d, skipping", rec.Kind)
		}
	}
}

func sleepDelta(ctx context.Context, prev, now int64, rate float64, loggedCap *bool) error {
	if prev == 0 || rate <= 0 {
		return nil
	}
	d := time.Duration(float64(now-prev) / rate)
	if d <= 0 {
		return nil
	}
	if d > maxReplayGap {
		if loggedCap != nil && !*loggedCap {
			log.Printf("replay: capping gap %s to %s (paused AC / CM)", d.Round(time.Second), maxReplayGap)
			*loggedCap = true
		}
		d = maxReplayGap
	}
	return sleep(ctx, d)
}
