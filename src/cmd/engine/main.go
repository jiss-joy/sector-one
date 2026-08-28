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
	startHTTP(*httpAddr, broadcaster, source)

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
