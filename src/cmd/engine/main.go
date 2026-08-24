package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"sector-one/internal/acudp"
)

func main() {
	host := flag.String("ac-host", "127.0.0.1", "Assetto Corsa host")
	port := flag.Int("ac-port", 9996, "Assetto Corsa UDP port")
	flag.Parse()

	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", *host, *port))
	if err != nil {
		log.Fatal(err)
	}

	// UDP has no real connection. This just binds a local port and
	// remembers that AC lives at addr, so Write/Read have a default peer.
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	if err := handshakeAndSubscribe(conn); err != nil {
		log.Fatal(err)
	}

	fmt.Println("subscribed — AC will now push 328-byte packets to this socket")
	fmt.Println("next step (you write this): Read() in a loop and parse RTCarInfo")
}

func handshakeAndSubscribe(conn *net.UDPConn) error {
	// 1) Ask AC who is driving / which track.
	if err := conn.SetReadDeadline(time.Now().Add(5 * time.Second)); err != nil {
		return err
	}
	if _, err := conn.Write(acudp.EncodeHandshake(acudp.OpHandshake)); err != nil {
		return fmt.Errorf("send handshake: %w", err)
	}

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		if os.IsTimeout(err) {
			return fmt.Errorf("AC did not answer on %s — is the 2014 game in a session (not the menu)?", conn.RemoteAddr())
		}
		return fmt.Errorf("read handshake: %w", err)
	}

	session, err := acudp.ParseHandshakeResponse(buf[:n])
	if err != nil {
		return err
	}
	fmt.Printf("car=%q driver=%q track=%q config=%q\n",
		session.CarName, session.DriverName, session.TrackName, session.TrackConfig)

	// 2) Tell AC to start sending live car updates to this socket.
	if _, err := conn.Write(acudp.EncodeHandshake(acudp.OpSubscribeUpdate)); err != nil {
		return fmt.Errorf("send subscribe: %w", err)
	}

	// Handshake is done. Later reads should not die after 5s of driving.
	return conn.SetReadDeadline(time.Time{})
}
