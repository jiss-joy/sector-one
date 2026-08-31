package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"sector-one/internal/acudp"
	"sector-one/internal/record"
	"sector-one/internal/stream"
)

func runLive(ctx context.Context, addr *net.UDPAddr, rec *record.Writer, broadcaster *stream.Broadcaster, source string) {
	for {
		if ctx.Err() != nil {
			return
		}
		conn, err := net.DialUDP("udp", nil, addr)
		if err != nil {
			printAndSleep(ctx, "error dialing UDP: ", err)
			continue
		}
		fmt.Println("local socket: ", conn.LocalAddr())

		packet := acudp.EncodeHandshake(acudp.OpHandshake)
		n, err := conn.Write(packet)
		if err != nil {
			_ = conn.Close()
			printAndSleep(ctx, "error sending handshake packet: ", err)
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
			printAndSleep(ctx, "waiting for a session (Content Manager home is fine)", err)
			continue
		}

		fmt.Println("got handshake reply: ", n, "bytes")
		session, err := acudp.ParseHandshakeResponse(buf[:n])
		if err != nil {
			_ = conn.Close()
			printAndSleep(ctx, "error parsing handshake response: ", err)
			continue
		}
		fmt.Printf("car=%q driver=%q track=%q config=%q\n",
			session.CarName, session.DriverName, session.TrackName, session.TrackConfig)

		writeRec(rec, record.KindHello, buf[:n])

		n, err = conn.Write(acudp.EncodeHandshake(acudp.OpSubscribeUpdate))
		if err != nil {
			_ = conn.Close()
			printAndSleep(ctx, "error sending subscribe update packet: ", err)
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
			publishCar(broadcaster, car, source, time.Now().UnixMilli())
			logCar(car, &lastLog, &packets)
		}
		dismiss(conn)
		_ = conn.Close()
		if quit {
			return
		}
	}
}

func printAndSleep(ctx context.Context, msg string, err error) {
	log.Println(msg, err)
	if err := sleep(ctx, retryEvery); err != nil {
		return
	}
}

func dismiss(conn *net.UDPConn) {
	_ = conn.SetWriteDeadline(time.Now().Add(time.Second))
	_, _ = conn.Write(acudp.EncodeHandshake(acudp.OpDismiss))
}
