package acudp

import (
	"encoding/binary"
	"fmt"
	"strings"
	"unicode/utf16"
)

const (
	OpHandshake        int32 = 0
	OpSubscribeUpdate  int32 = 1
	OpSubscribeSpot    int32 = 2
	OpDismiss          int32 = 3
	HandshakeSize            = 12
	DefaultIdentifier  int32 = 0
	DefaultVersion     int32 = 1
)

type SessionInfo struct {
	CarName      string
	DriverName   string
	TrackName    string
	TrackConfig  string
	Identifier   int32
	Version      int32
}

// EncodeHandshake builds the 12-byte packet AC expects: three little-endian int32s.
func EncodeHandshake(operation int32) []byte {
	buf := make([]byte, HandshakeSize)
	binary.LittleEndian.PutUint32(buf[0:4], uint32(DefaultIdentifier))
	binary.LittleEndian.PutUint32(buf[4:8], uint32(DefaultVersion))
	binary.LittleEndian.PutUint32(buf[8:12], uint32(operation))
	return buf
}

// ParseHandshakeResponse reads AC's 408- or 808-byte reply after op=0.
func ParseHandshakeResponse(data []byte) (SessionInfo, error) {
	var zero SessionInfo
	if len(data) != 408 && len(data) != 808 {
		return zero, fmt.Errorf("handshake response: want 408 or 808 bytes, got %d", len(data))
	}

	// Four UTF-16LE strings + two int32s. n is bytes per name field.
	n := (len(data) - 8) / 4
	carEnd := n
	drvEnd := 2 * n
	trackStart := drvEnd + 8
	trackEnd := trackStart + n

	return SessionInfo{
		CarName:     decodeACString(data[0:carEnd]),
		DriverName:  decodeACString(data[carEnd:drvEnd]),
		Identifier:  int32(binary.LittleEndian.Uint32(data[drvEnd : drvEnd+4])),
		Version:     int32(binary.LittleEndian.Uint32(data[drvEnd+4 : drvEnd+8])),
		TrackName:   decodeACString(data[trackStart:trackEnd]),
		TrackConfig: decodeACString(data[trackEnd : trackEnd+n]),
	}, nil
}

// AC strings are UTF-16LE and often end at '%' (Kunos quirk) or a NUL.
func decodeACString(buf []byte) string {
	if len(buf)%2 != 0 {
		buf = buf[:len(buf)-1]
	}
	u16s := make([]uint16, len(buf)/2)
	for i := range u16s {
		u16s[i] = binary.LittleEndian.Uint16(buf[i*2:])
	}
	s := string(utf16.Decode(u16s))

	cut := len(s)
	for i, r := range s {
		if r == '%' || r == 0 {
			cut = i
			break
		}
	}
	return strings.TrimSpace(s[:cut])
}
