package acudp

import (
	"encoding/binary"
	"fmt"
	"strings"
	"unicode/utf16"
)

const (
	OPHandshake       int32 = 0
	OPSubscribeUpdate int32 = 1
	OPDismiss         int32 = 2
	DefaultIdentifier int32 = 0
	DefaultVersion    int32 = 1
	HandshakeSize           = 12
)

type SessionInfo struct {
	CarName     string
	DriverName  string
	TrackName   string
	TrackConfig string
	Identifier  int32
	Version     int32
}

// ParseHandshakeResponse reads AC's 408- or 808-byte reply after op=0.
func ParseHandshakeResponse(data []byte) (SessionInfo, error) {
	var zero SessionInfo
	if len(data) != 408 && len(data) != 808 {
		return zero, fmt.Errorf("Handshake response: want 408 or 808 bytes, got %d", len(data))
	}

	// Four UTF-16LE strings + two int32s. n is bytes per name field.
	// byte 0                                                    byte 407
	// |---- car ----|---- driver ----|-- id --|-- ver --|---- track ----|---- config ----|
	//  100 bytes       100 bytes      4 B      4 B        100 bytes        100 bytes
	n := (len(data) - 8) / 4
	carEnd := n
	driverEnd := carEnd * 2
	trackStart := driverEnd + 8
	trackEnd := trackStart + n

	return SessionInfo{
		CarName:     decodeAssettoCorsaString(data[0:carEnd]),
		DriverName:  decodeAssettoCorsaString(data[carEnd:driverEnd]),
		Identifier:  toInt(data[driverEnd : driverEnd+4]),
		Version:     toInt(data[driverEnd+4 : driverEnd+8]),
		TrackName:   decodeAssettoCorsaString(data[trackStart:trackEnd]),
		TrackConfig: decodeAssettoCorsaString(data[trackEnd : trackEnd+n]),
	}, nil
}

// EncodeHandshake builds the 12-byte packet AC expects: three little-endian int32s.
func EncodeHandshake(operation int32) []byte {
	buf := make([]byte, HandshakeSize)
	binary.LittleEndian.PutUint32(buf[0:4], uint32(DefaultIdentifier))
	binary.LittleEndian.PutUint32(buf[4:8], uint32(DefaultVersion))
	binary.LittleEndian.PutUint32(buf[8:12], uint32(operation))

	return buf
}

func toInt(data []byte) int32 {
	return int32(binary.LittleEndian.Uint32(data))
}

func decodeAssettoCorsaString(buf []byte) string {
	// UTF-16 is pairs. An odd length cannot be decoded. Slice off the last byte.
	if len(buf)%2 != 0 {
		buf = buf[:len(buf)-1]
	}

	u16 := make([]uint16, len(buf)/2)
	for i := range u16 {
		u16[i] = binary.LittleEndian.Uint16(buf[i*2:])
	}
	str := string(utf16.Decode(u16))

	cutLength := len(str)
	// Calculate the cut length
	for i, r := range str {
		if r == '%' || r == 0 {
			cutLength = i
			break
		}
	}
	return strings.TrimSpace(str[:cutLength])
}
