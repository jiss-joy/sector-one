package physics

import (
	"sector-one/internal/acudp"
)

type Frame struct {
	TS            int64      `json:"ts"`
	SpeedKmh      float32    `json:"speed_kmh"`
	RPM           float32    `json:"rpm"`
	Gear          int32      `json:"gear"`
	Throttle      float32    `json:"throttle"`
	Brake         float32    `json:"brake"`
	Clutch        float32    `json:"clutch"`
	Steer         float32    `json:"steer"`
	GLat          float32    `json:"g_lat"`
	GLong         float32    `json:"g_long"`
	GVert         float32    `json:"g_vert"`
	LapTimeMs     int32      `json:"lap_time_ms"`
	LastLapMs     int32      `json:"last_lap_ms"`
	BestLapMs     int32      `json:"best_lap_ms"`
	LapCount      int32      `json:"lap_count"`
	AbsEnabled    bool       `json:"abs_enabled"`
	AbsInAction   bool       `json:"abs_in_action"`
	TcEnabled     bool       `json:"tc_enabled"`
	TcInAction    bool       `json:"tc_in_action"`
	InPit         bool       `json:"in_pit"`
	EngineLimiter bool       `json:"engine_limiter"`
	WheelRadS     [4]float32 `json:"wheel_rad_s"`
	SlipRatio     [4]float32 `json:"slip_ratio"`
	SlipAngle     [4]float32 `json:"slip_angle"`
	NormalizedPos float32    `json:"normalized_pos"`
	LoadN         [4]float32 `json:"load_n"`
	MaxRPM        float32    `json:"max_rpm"`
	MaxLoad       float32    `json:"max_load"`
	Source        string     `json:"source"`
	Car           string     `json:"car"`
}

func FromCar(c acudp.CarInfo, source string, ts int64, spec CarSpec, carName string) Frame {
	return Frame{
		TS:            ts,
		SpeedKmh:      c.SpeedKmh,
		RPM:           c.EngineRPM,
		Gear:          c.Gear,
		Throttle:      c.Gas,
		Brake:         c.Brake,
		Clutch:        c.Clutch,
		Steer:         c.Steer,
		GLat:          c.AccGHorizontal,
		GLong:         c.AccGFrontal,
		GVert:         c.AccGVertical,
		LapTimeMs:     c.LapTimeMs,
		LastLapMs:     c.LastLapMs,
		BestLapMs:     c.BestLapMs,
		LapCount:      c.LapCount,
		AbsEnabled:    c.AbsEnabled,
		AbsInAction:   c.AbsInAction,
		TcEnabled:     c.TcEnabled,
		TcInAction:    c.TcInAction,
		InPit:         c.InPit,
		EngineLimiter: c.EngineLimiter,
		WheelRadS:     c.WheelRadS,
		SlipRatio:     c.SlipRatio,
		SlipAngle:     c.SlipAngle,
		NormalizedPos: c.NormalizedPos,
		LoadN:         c.LoadN,
		MaxRPM:        spec.MaxRPM,
		MaxLoad:       spec.MaxLoad,
		Source:        source,
		Car:           carName,
	}
}
