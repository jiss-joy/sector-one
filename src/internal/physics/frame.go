package physics

import (
	"time"

	"sector-one/internal/acudp"
)

type Frame struct {
	TS        int64   `json:"ts"`
	SpeedKmh  float32 `json:"speed_kmh"`
	RPM       float32 `json:"rpm"`
	Gear      int32   `json:"gear"`
	Throttle  float32 `json:"throttle"`
	Brake     float32 `json:"brake"`
	Clutch    float32 `json:"clutch"`
	Steer     float32 `json:"steer"`
	GLat      float32 `json:"g_lat"`
	GLong     float32 `json:"g_long"`
	GVert     float32 `json:"g_vert"`
	LapTimeMs int32   `json:"lap_time_ms"`
	LastLapMs int32   `json:"last_lap_ms"`
	BestLapMs int32   `json:"best_lap_ms"`
	LapCount  int32   `json:"lap_count"`
	Source    string  `json:"source"`
}

func FromCar(c acudp.CarInfo, source string) Frame {
	return Frame{
		TS:        time.Now().UnixMilli(),
		SpeedKmh:  c.SpeedKmh,
		RPM:       c.EngineRPM,
		Gear:      c.Gear,
		Throttle:  c.Gas,
		Brake:     c.Brake,
		Clutch:    c.Clutch,
		Steer:     c.Steer,
		GLat:      c.AccGHorizontal,
		GLong:     c.AccGFrontal,
		GVert:     c.AccGVertical,
		LapTimeMs: c.LapTimeMs,
		LastLapMs: c.LastLapMs,
		BestLapMs: c.BestLapMs,
		LapCount:  c.LapCount,
		Source:    source,
	}
}
