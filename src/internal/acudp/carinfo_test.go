package acudp

import (
	"encoding/binary"
	"math"
	"testing"
)

func TestParseCarInfo(t *testing.T) {
	buf := make([]byte, CarInfoSize)
	buf[0] = 'a'
	binary.LittleEndian.PutUint32(buf[4:], CarInfoSize)
	putFloat32LE(buf, 8, 123.5)
	putFloat32LE(buf, 68, 5120)
	binary.LittleEndian.PutUint32(buf[76:], 3)

	got, err := ParseCarInfo(buf)
	if err != nil {
		t.Fatal(err)
	}
	if got.SpeedKmh != 123.5 || got.EngineRPM != 5120 || got.Gear != 3 {
		t.Fatalf("got %+v", got)
	}
}

func TestParseCarInfoSkipsBoolPad(t *testing.T) {
	buf := make([]byte, CarInfoSize)
	buf[0] = 'a'
	putFloat32LE(buf, 28, -1.25)

	got, err := ParseCarInfo(buf)
	if err != nil {
		t.Fatal(err)
	}
	if got.AccGVertical != -1.25 {
		t.Fatalf("AccGVertical=%v, pad was consumed", got.AccGVertical)
	}
}

func TestParseCarInfoRejects(t *testing.T) {
	if _, err := ParseCarInfo(make([]byte, 100)); err == nil {
		t.Fatal("expected size error")
	}
	buf := make([]byte, CarInfoSize)
	buf[0] = 'x'
	if _, err := ParseCarInfo(buf); err == nil {
		t.Fatal("expected identifier error")
	}
}

func putFloat32LE(buf []byte, offset int, v float32) {
	binary.LittleEndian.PutUint32(buf[offset:], math.Float32bits(v))
}
