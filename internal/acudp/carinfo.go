package acudp

import (
	"encoding/binary"
	"fmt"
	"math"
)

const (
	CarInfoSize       = 328
	CarInfoIdentifier = 'a'
)

type CarInfo struct {
	AbsEnabled     bool
	AbsInAction    bool
	TcEnabled      bool
	TcInAction     bool
	InPit          bool
	EngineLimiter  bool
	AccGVertical   float32
	AccGHorizontal float32
	AccGFrontal    float32
	LapTimeMs      int32
	LastLapMs      int32
	BestLapMs      int32
	LapCount       int32
	Gear           int32
	SpeedKmh       float32
	Gas            float32
	Brake          float32
	Clutch         float32
	EngineRPM      float32
	Steer          float32
	TrackProgress  float32    // between 0 and 1
	WheelSpeed     [4]float32 // radians per second (rad/s)
	SlipAngle      [4]float32 // Radians
	SlipRatio      [4]float32 // Ratio
	TyreLoad       [4]float32 // Newtons (N)
}

func ParseCarInfo(data []byte) (CarInfo, error) {
	var zero CarInfo
	dataLength := len(data)
	if dataLength != CarInfoSize {
		return zero, fmt.Errorf("Car info: Want %d bytes, got %d", CarInfoSize, dataLength)
	}
	if data[0] != CarInfoIdentifier {
		return zero, fmt.Errorf("Car info: Want %q as identifier, got %q", CarInfoIdentifier, data[0])
	}

	return CarInfo{
		AbsEnabled:     data[20] != 0,
		AbsInAction:    data[21] != 0,
		TcEnabled:      data[22] != 0,
		TcInAction:     data[23] != 0,
		InPit:          data[24] != 0,
		EngineLimiter:  data[25] != 0,
		AccGVertical:   toFloat32(data, 28),
		AccGHorizontal: toFloat32(data, 32),
		AccGFrontal:    toFloat32(data, 36),
		LapTimeMs:      toInt32(data, 40),
		LastLapMs:      toInt32(data, 44),
		BestLapMs:      toInt32(data, 48),
		LapCount:       toInt32(data, 52),
		Gear:           toInt32(data, 76),
		SpeedKmh:       toFloat32(data, 8),
		Gas:            toFloat32(data, 56),
		Brake:          toFloat32(data, 60),
		Clutch:         toFloat32(data, 64),
		EngineRPM:      toFloat32(data, 68),
		Steer:          toFloat32(data, 72),
		TrackProgress:  toFloat32(data, 308),
		WheelSpeed:     toFloat32x4(data, 84),
		SlipAngle:      toFloat32x4(data, 100),
		SlipRatio:      toFloat32x4(data, 132),
		TyreLoad:       toFloat32x4(data, 180),
	}, nil
}

func toFloat32x4(data []byte, offset int) [4]float32 {
	return [4]float32{
		toFloat32(data, offset),
		toFloat32(data, offset+4),
		toFloat32(data, offset+8),
		toFloat32(data, offset+12),
	}
}

func toInt32(data []byte, offset int) int32 {
	return int32(binary.LittleEndian.Uint32(data[offset : offset+4]))
}

func toFloat32(data []byte, offset int) float32 {
	return math.Float32frombits(binary.LittleEndian.Uint32(data[offset : offset+4]))
}
