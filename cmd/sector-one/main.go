package main

import (
	"context"
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
	RetryEvery = 2 * time.Second
	LogEvery   = 2 * time.Second
	ACHost     = "127.0.0.1"
	ACPort     = 9996
)

func main() {
	replayPath := flag.String("replay", "", "Load a .bin file instead of live data")
	recordPath := flag.String("record", "", "Write framed session to this .bin")
	replayRate := flag.Float64("replay-rate", 1.0, "Replay speed multiplier")
	httpAddr := flag.String("http", "127.0.0.1:8080", "HTTP listen address")
	flag.Parse()

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
			err := runReplay(ctx, *replayPath, *replayRate, sm, broadcaster)
			if err != nil {
				log.Println("Error finding file", err)
			}
			cancel()
			return
		}

		addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", ACHost, ACPort))
		if err != nil {
			log.Println("UDP Addr: ", err)
			cancel()
			return
		}
		fmt.Println("Assetto Corsa Address: ", addr)

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

	openBrowser("http://" + *httpAddr)
	hideConsole()

	go func() {
		<-ctx.Done()
		systrayQuit()
	}()

	runTray(*httpAddr, cancel)
	<-ctx.Done()
}
