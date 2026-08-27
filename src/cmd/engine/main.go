package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sector-one/internal/acudp"
	"sector-one/internal/physics"
	"sector-one/internal/record"
	"sector-one/internal/stream"
	"syscall"
	"time"
)

const (
	retryEvery   = 2 * time.Second
	maxReplayGap = 200 * time.Millisecond
	logEvery     = 2 * time.Second
	acHost       = "127.0.0.1"
	acPort       = 9996
)

func main() {
	recordPath := flag.String("record", "", "write framed session to this .bin")
	replayPath := flag.String("replay", "", "play a .bin instead of talking to AC")
	replayRate := flag.Float64("replay-rate", 1.0, "replay speed multiplier")
	flag.Parse()

	if *recordPath != "" && *replayPath != "" {
		log.Fatal("use --record or --replay, not both")
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	source := "live"
	if *replayPath != "" {
		source = "replay"
	}
	broadcaster := stream.NewBroadcaster()
	startHTTP(broadcaster, source)

	if *replayPath != "" {
		if err := runReplay(ctx, *replayPath, *replayRate, broadcaster, source); err != nil && !errors.Is(err, context.Canceled) {
			log.Fatal(err)
		}
		return
	}

	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", acHost, acPort))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Assetto Corsa address:", addr)

	var rec *record.Writer
	if *recordPath != "" {
		rec, err = record.NewWriter(*recordPath)
		if err != nil {
			log.Fatal(err)
		}
		defer rec.Close()
		log.Println("recording to:", *recordPath)
	}

	runLive(ctx, addr, rec, broadcaster, source)
}

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

func runLive(ctx context.Context, addr *net.UDPAddr, rec *record.Writer, broadcaster *stream.Broadcaster, source string) {
	for {
		if ctx.Err() != nil {
			return
		}
		conn, err := net.DialUDP("udp", nil, addr)
		if err != nil {
			log.Println("error dialing UDP: ", err)
			if err := sleep(ctx, retryEvery); err != nil {
				return
			}
			continue
		}
		fmt.Println("local socket: ", conn.LocalAddr())

		packet := acudp.EncodeHandshake(acudp.OpHandshake)
		n, err := conn.Write(packet)
		if err != nil {
			_ = conn.Close()
			log.Println("error sending handshake packet: ", err)
			if err := sleep(ctx, retryEvery); err != nil {
				return
			}
			continue
		}
		fmt.Println("sent handshake packet: ", n, "bytes")
		if err := conn.SetReadDeadline(time.Now().Add(5 * time.Second)); err != nil {
			_ = conn.Close()
			log.Println("error setting read deadline: ", err)
			if err := sleep(ctx, retryEvery); err != nil {
				return
			}
			continue
		}
		buf := make([]byte, 1024)
		n, err = conn.Read(buf)
		if err != nil {
			_ = conn.Close()
			log.Println("waiting for a session (Content Manager home is fine)")
			if err := sleep(ctx, retryEvery); err != nil {
				return
			}
			continue
		}

		fmt.Println("got handshake reply: ", n, "bytes")
		session, err := acudp.ParseHandshakeResponse(buf[:n])
		if err != nil {
			_ = conn.Close()
			log.Println("error parsing handshake response: ", err)
			if err := sleep(ctx, retryEvery); err != nil {
				return
			}
			continue
		}
		fmt.Printf("car=%q driver=%q track=%q config=%q\n",
			session.CarName, session.DriverName, session.TrackName, session.TrackConfig)

		writeRec(rec, record.KindHello, buf[:n])

		n, err = conn.Write(acudp.EncodeHandshake(acudp.OpSubscribeUpdate))
		if err != nil {
			_ = conn.Close()
			log.Println("error sending subscribe update packet: ", err)
			if err := sleep(ctx, retryEvery); err != nil {
				return
			}
			continue
		}
		fmt.Println("sent subscribe update packet: ", n, "bytes")

		lastLog := time.Time{}
		packets := 0
		quit := false
		for {
			if ctx.Err() != nil {
				quit = true
				break
			}
			if err := conn.SetReadDeadline(time.Now().Add(5 * time.Second)); err != nil {
				log.Println("error setting read deadline: ", err)
				break
			}
			n, err = conn.Read(buf)
			if err != nil {
				if ctx.Err() != nil {
					quit = true
					break
				}
				log.Println("session ended: ", err)
				break
			}
			car, err := acudp.ParseCarInfo(buf[:n])
			if err != nil {
				log.Println("skipping packet: ", err)
				continue
			}
			writeRec(rec, record.KindCar, buf[:n])
			publishCar(broadcaster, car, source)
			logCar(car, &lastLog, &packets)
		}
		dismiss(conn)
		_ = conn.Close()
		if quit {
			return
		}
	}
}

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
			publishCar(broadcaster, car, source)
			logCar(car, &lastLog, &packets)
		default:
			log.Printf("replay: unknown kind %d, skipping", rec.Kind)
		}
	}
}

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

func dismiss(conn *net.UDPConn) {
	_ = conn.SetWriteDeadline(time.Now().Add(time.Second))
	_, _ = conn.Write(acudp.EncodeHandshake(acudp.OpDismiss))
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
