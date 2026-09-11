package physics

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

type CarSpec struct {
	MaxRPM  float32 `json:"max_rpm"`
	MaxLoad float32 `json:"max_load"`
}

type SpecManager struct {
	mu    sync.RWMutex
	path  string
	specs map[string]*CarSpec
}

func NewSpecManager(path string) *SpecManager {
	sm := &SpecManager{
		path:  path,
		specs: make(map[string]*CarSpec),
	}
	sm.load()
	return sm
}

func (sm *SpecManager) Get(carName string) CarSpec {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	if s, ok := sm.specs[carName]; ok {
		return *s
	}
	// Defaults if unknown
	return CarSpec{MaxRPM: 6000, MaxLoad: 5000}
}

func (sm *SpecManager) Update(carName string, rpm, load float32, limiter bool) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	s, ok := sm.specs[carName]
	if !ok {
		s = &CarSpec{MaxRPM: 6000, MaxLoad: 5000}
		sm.specs[carName] = s
	}

	changed := false
	if rpm > s.MaxRPM {
		s.MaxRPM = rpm
		changed = true
	}
	// If limiter is active, we treat this as the definitive max
	if limiter && rpm > 0 {
		s.MaxRPM = rpm
		changed = true
	}
	if load > s.MaxLoad {
		s.MaxLoad = load
		changed = true
	}

	if changed {
		// Periodically save or just save every change for now (small file)
		sm.save()
	}
}

func (sm *SpecManager) load() {
	b, err := os.ReadFile(sm.path)
	if err != nil {
		return
	}
	_ = json.Unmarshal(b, &sm.specs)
}

func (sm *SpecManager) save() {
	b, err := json.MarshalIndent(sm.specs, "", "  ")
	if err != nil {
		log.Println("spec save error:", err)
		return
	}
	_ = os.WriteFile(sm.path, b, 0644)
}
