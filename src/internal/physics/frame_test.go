package physics

import (
	"encoding/json"
	"testing"

	"sector-one/internal/acudp"
)

func TestFromCarJSONKeys(t *testing.T) {
	f := FromCar(acudp.CarInfo{SpeedKmh: 100, EngineRPM: 5000, Gear: 3}, "replay", 12345, CarSpec{MaxRPM: 8000, MaxLoad: 5000}, "ks_porsche")
	b, err := json.Marshal(f)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatal(err)
	}
	for _, k := range []string{"ts", "speed_kmh", "rpm", "gear", "throttle", "g_lat", "source", "car", "slip_ratio", "load_n", "abs_in_action"} {
		if _, ok := m[k]; !ok {
			t.Fatalf("missing key %q in %s", k, b)
		}
	}
	if m["source"] != "replay" {
		t.Fatalf("source=%v", m["source"])
	}
	if m["car"] != "ks_porsche" {
		t.Fatalf("car=%v", m["car"])
	}
}
