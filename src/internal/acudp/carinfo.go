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
type CarInfo struct {
	SpeedKmh     float32
	Gas          float32
	Brake        float32
	EngineRPM    float32
	Gear         int32
	AccGVertical float32
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
		SpeedKmh:     float32LE(data, 8),
		AccGVertical: float32LE(data, 28), // after 6 bools + 2-byte pad
		Gas:          float32LE(data, 56),
		Brake:        float32LE(data, 60),
		EngineRPM:    float32LE(data, 68),
		Gear:         int32(binary.LittleEndian.Uint32(data[76:80])),
	}, nil
}

func float32LE(data []byte, offset int) float32 {
	return math.Float32frombits(binary.LittleEndian.Uint32(data[offset : offset+4]))
}
