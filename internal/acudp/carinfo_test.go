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

func TestParseCarInfoExtendedOffsets(t *testing.T) {
	buf := make([]byte, CarInfoSize)
	buf[0] = 'a'
	putFloat32LE(buf, 32, 1.5)
	putFloat32LE(buf, 36, -0.75)
	putFloat32LE(buf, 72, 0.25)
	binary.LittleEndian.PutUint32(buf[52:], 4)

	got, err := ParseCarInfo(buf)
	if err != nil {
		t.Fatal(err)
	}
	if got.AccGHorizontal != 1.5 || got.AccGFrontal != -0.75 || got.Steer != 0.25 || got.LapCount != 4 {
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

func TestParseCarInfoWheelsAndFlags(t *testing.T) {
	buf := make([]byte, CarInfoSize)
	buf[0] = 'a'
	buf[21] = 1 // ABS in action
	buf[23] = 1 // TC in action
	putFloat32LE(buf, 84, 10)
	putFloat32LE(buf, 88, 20)
	putFloat32LE(buf, 92, 30)
	putFloat32LE(buf, 96, 40)
	putFloat32LE(buf, 132, 0.1)
	putFloat32LE(buf, 136, 0.2)
	putFloat32LE(buf, 140, 0.3)
	putFloat32LE(buf, 144, 0.4)
	putFloat32LE(buf, 180, 1000)
	putFloat32LE(buf, 184, 2000)
	putFloat32LE(buf, 188, 3000)
	putFloat32LE(buf, 192, 4000)

	got, err := ParseCarInfo(buf)
	if err != nil {
		t.Fatal(err)
	}
	if !got.AbsInAction || !got.TcInAction || got.AbsEnabled {
		t.Fatalf("flags %+v", got)
	}
	if got.WheelSpeed != [4]float32{10, 20, 30, 40} {
		t.Fatalf("Wheel Speed=%v", got.WheelSpeed)
	}
	if got.SlipRatio != [4]float32{0.1, 0.2, 0.3, 0.4} {
		t.Fatalf("SlipRatio=%v", got.SlipRatio)
	}
	if got.TyreLoad != [4]float32{1000, 2000, 3000, 4000} {
		t.Fatalf("TyreLoad=%v", got.TyreLoad)
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
