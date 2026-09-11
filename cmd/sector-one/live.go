package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sector-one/internal/acudp"
	"sector-one/internal/physics"
	"sector-one/internal/record"
	"sector-one/internal/stream"
	"time"
)

func runLive(ctx context.Context, addr *net.UDPAddr, rec *record.Writer, sm *physics.SpecManager, broadcaster *stream.Broadcaster, source string) {
	for {
		if ctx.Err() != nil {
			return
		}
		connection, err := net.DialUDP("udp", nil, addr)
		if err != nil {
			printAndSleep(ctx, "Error dialing UDP", err)
			continue
		}
		fmt.Println("Local socket: ", connection.LocalAddr())

		packet := acudp.EncodeHandshake(acudp.OPHandshake)
		n, err := connection.Write(packet)
		if err != nil {
			connection.Close()
			printAndSleep(ctx, "Error sending handshake packet", err)
			continue
		}
		fmt.Println("Sent handshake packet: ", n, "bytes")
		err = connection.SetReadDeadline(time.Now().Add(5 * time.Second))
		if err != nil {
			connection.Close()
			log.Println("Error setting read deadline: ", err)
			if err := sleep(ctx, RetryEvery); err != nil {
				return
			}
			continue
		}
		buf := make([]byte, 1024)
		n, err = connection.Read(buf)
		if err != nil {
			connection.Close()
			printAndSleep(ctx, "Waiting for a session (Content Manager home is fine)", err)
			continue
		}
		session, err := acudp.ParseHandshakeResponse(buf[:n])
		if err != nil {
			connection.Close()
			printAndSleep(ctx, "Error parsing handshake response", err)
			continue
		}
		fmt.Printf("car=%q driver=%q track=%q config=%q\n",
			session.CarName, session.DriverName, session.TrackName, session.TrackConfig)

		writeRecord(rec, record.KindHandshake, buf[:n])

		packet = acudp.EncodeHandshake(acudp.OPSubscribeUpdate)
		n, err = connection.Write(packet)
		if err != nil {
			connection.Close()
			printAndSleep(ctx, "Error sending subscribe update packet: ", err)
			continue
		}
		fmt.Println("Sent subscribe update packet: ", n, "bytes")

		quit := false
		for {
			if ctx.Err() != nil {
				quit = true
				break
			}
			err := connection.SetReadDeadline(time.Now().Add(5 * time.Second))
			if err != nil {
				log.Println("Error setting read deadline", err)
				break
			}
			n, err = connection.Read(buf)
			if err != nil {
				if ctx.Err() != nil {
					quit = true
					break
				}
				log.Println("Session ended: ", err)
				break
			}
			car, err := acudp.ParseCarInfo(buf[:n])
			if err != nil {
				log.Println("Skipping packet: ", err)
				continue
			}
			writeRecord(rec, record.KindCarInfo, buf[:n])

			// Auto-calibrate
			maxLoad := float32(0)
			for _, load := range car.TyreLoad {
				if load > maxLoad {
					maxLoad = load
				}
			}
			sm.Update(session.CarName, car.EngineRPM, maxLoad, car.EngineLimiter)

			publishCar(sm, broadcaster, car, source, time.Now().UnixMilli(), session.CarName)
		}
		dismiss(connection)
		connection.Close()
		if quit {
			return
		}
	}
}

func printAndSleep(ctx context.Context, msg string, err error) {
	log.Println(msg, err)
	if err := sleep(ctx, RetryEvery); err != nil {
		return
	}
}

func dismiss(connection *net.UDPConn) {
	connection.SetWriteDeadline(time.Now().Add(time.Second))
	connection.Write(acudp.EncodeHandshake(acudp.OPDismiss))
}
