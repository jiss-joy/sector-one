package physics

import "sector-one/internal/acudp"

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
	WheelSpeed    [4]float32 `json:"wheel_speed"`
	SlipRatio     [4]float32 `json:"slip_ratio"`
	SlipAngle     [4]float32 `json:"slip_angle"`
	TrackProgress float32    `json:"track_progress"`
	TyreLoad      [4]float32 `json:"tyre_load"`
	MaxRPM        float32    `json:"max_rpm"`
	MaxLoad       float32    `json:"max_load"`
	Source        string     `json:"source"`
	Car           string     `json:"car"`
}

func FromCar(carInfo acudp.CarInfo, source string, timestamp int64, spec CarSpec, carName string) Frame {
	return Frame{
		TS:            timestamp,
		SpeedKmh:      carInfo.SpeedKmh,
		RPM:           carInfo.EngineRPM,
		Gear:          carInfo.Gear,
		Throttle:      carInfo.Gas,
		Brake:         carInfo.Brake,
		Clutch:        carInfo.Clutch,
		Steer:         carInfo.Steer,
		GLat:          carInfo.AccGHorizontal,
		GLong:         carInfo.AccGFrontal,
		GVert:         carInfo.AccGVertical,
		LapTimeMs:     carInfo.LapTimeMs,
		LastLapMs:     carInfo.LastLapMs,
		BestLapMs:     carInfo.BestLapMs,
		LapCount:      carInfo.LapCount,
		AbsEnabled:    carInfo.AbsEnabled,
		AbsInAction:   carInfo.AbsInAction,
		TcEnabled:     carInfo.TcEnabled,
		TcInAction:    carInfo.TcInAction,
		InPit:         carInfo.InPit,
		EngineLimiter: carInfo.EngineLimiter,
		WheelSpeed:    carInfo.WheelSpeed,
		SlipRatio:     carInfo.SlipRatio,
		SlipAngle:     carInfo.SlipAngle,
		TrackProgress: carInfo.TrackProgress,
		TyreLoad:      carInfo.TyreLoad,
		MaxRPM:        spec.MaxRPM,
		MaxLoad:       spec.MaxLoad,
		Source:        source,
		Car:           carName,
	}
}
