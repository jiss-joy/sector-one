package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sector-one/internal/acudp"
	"sector-one/internal/record"
	"syscall"
	"time"
)

const (
	retryEvery = 2 * time.Second
)

func main() {
	host := "127.0.0.1"
	port := 9996

	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Assetto Corsa address:", addr)

	recordPath := flag.String("record", "", "write framed session to this .bin")
	flag.Parse()

	var rec *record.Writer
	if *recordPath != "" {
		var err error
		rec, err = record.NewWriter(*recordPath)
		if err != nil {
			log.Fatal(err)
		}
		defer rec.Close()
		log.Println("recording to: ", *recordPath)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

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
		// Parse the handshake response
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

		// Subscribe to updates (AC will send 328 byte packets to conn)
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

			if lastLog.IsZero() {
				lastLog = time.Now()
			}
			packets++
			if elapsed := time.Since(lastLog); elapsed >= 2*time.Second {
				hz := float64(packets) / elapsed.Seconds()
				fmt.Printf("speed_kmh=%.1f gear=%d rpm=%.0f  (%.0f Hz ingest)\n",
					car.SpeedKmh, car.Gear, car.EngineRPM, hz)
				packets = 0
				lastLog = time.Now()
			}
		}
		dismiss(conn)
		_ = conn.Close()
		if quit {
			return
		}
	}
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
