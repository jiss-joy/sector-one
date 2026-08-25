package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sector-one/internal/acudp"
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
			if time.Since(lastLog) > 2*time.Second {
				fmt.Printf("speed_kmh=%.1f gear=%d rpm=%.0f\n", car.SpeedKmh, car.Gear, car.EngineRPM)
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
