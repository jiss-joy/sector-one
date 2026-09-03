package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
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
	httpAddr := flag.String("http", "127.0.0.1:8080", "HTTP listen address")
	noTray := flag.Bool("no-tray", false, "run in the foreground without a system tray icon")
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
	sm := physics.NewSpecManager("cars.json")
	startHTTP(*httpAddr, broadcaster, source)

	go func() {
		if *replayPath != "" {
			if err := runReplay(ctx, *replayPath, *replayRate, sm, broadcaster, source); err != nil && !errors.Is(err, context.Canceled) {
				log.Println("replay:", err)
			}
			if *noTray {
				cancel()
			}
			return
		}

		addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", acHost, acPort))
		if err != nil {
			log.Println("udp addr:", err)
			cancel()
			return
		}
		fmt.Println("Assetto Corsa address:", addr)

		var rec *record.Writer
		if *recordPath != "" {
			rec, err = record.NewWriter(*recordPath)
			if err != nil {
				log.Println("record:", err)
				cancel()
				return
			}
			defer rec.Close()
			log.Println("recording to:", *recordPath)
		}

		runLive(ctx, addr, rec, sm, broadcaster, source)
		cancel()
	}()

	if *noTray {
		openBrowser("http://" + *httpAddr)
		<-ctx.Done()
		return
	}

	hideConsole()
	go func() {
		<-ctx.Done()
		systrayQuit()
	}()
	runTray(*httpAddr, cancel)
}
