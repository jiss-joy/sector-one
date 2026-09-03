package acudp

import (
	"encoding/binary"
	"fmt"
	"math"
)

const (
	CarInfoSize       = 328
	carInfoIdentifier = 'a'
)

// Gear: 0 = reverse, 1 = neutral, 2+ = forward.
// Wheel arrays are Kunos order: FL, FR, RL, RR ([Likely] — confirm on a known corner).
type CarInfo struct {
	SpeedKmh       float32
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
	Gas            float32
	Brake          float32
	Clutch         float32
	EngineRPM      float32
	Steer          float32
	Gear           int32
	WheelRadS      [4]float32
	SlipRatio      [4]float32
	LoadN          [4]float32
	SlipAngle      [4]float32
	NormalizedPos  float32
}

func ParseCarInfo(data []byte) (CarInfo, error) {
	var zero CarInfo
	if len(data) != CarInfoSize {
		return zero, fmt.Errorf("car info: want %d bytes, got %d", CarInfoSize, len(data))
	}
	if data[0] != carInfoIdentifier {
		return zero, fmt.Errorf("car info: identifier %q, want %q", data[0], carInfoIdentifier)
	}

	return CarInfo{
		SpeedKmh:       float32LE(data, 8),
		AbsEnabled:     data[20] != 0,
		AbsInAction:    data[21] != 0,
		TcEnabled:      data[22] != 0,
		TcInAction:     data[23] != 0,
		InPit:          data[24] != 0,
		EngineLimiter:  data[25] != 0,
		AccGVertical:   float32LE(data, 28), // after 6 bools + 2-byte pad
		AccGHorizontal: float32LE(data, 32),
		AccGFrontal:    float32LE(data, 36),
		LapTimeMs:      int32LE(data, 40),
		LastLapMs:      int32LE(data, 44),
		BestLapMs:      int32LE(data, 48),
		LapCount:       int32LE(data, 52),
		Gas:            float32LE(data, 56),
		Brake:          float32LE(data, 60),
		Clutch:         float32LE(data, 64),
		EngineRPM:      float32LE(data, 68),
		Steer:          float32LE(data, 72),
		Gear:           int32LE(data, 76),
		WheelRadS:      float32x4LE(data, 84),
		SlipRatio:      float32x4LE(data, 132),
		LoadN:          float32x4LE(data, 180),
		SlipAngle:      float32x4LE(data, 100),
		NormalizedPos:  float32LE(data, 308),
	}, nil
}

func float32x4LE(data []byte, offset int) [4]float32 {
	return [4]float32{
		float32LE(data, offset),
		float32LE(data, offset+4),
		float32LE(data, offset+8),
		float32LE(data, offset+12),
	}
}

func int32LE(data []byte, offset int) int32 {
	return int32(binary.LittleEndian.Uint32(data[offset : offset+4]))
}

func float32LE(data []byte, offset int) float32 {
	return math.Float32frombits(binary.LittleEndian.Uint32(data[offset : offset+4]))
}
