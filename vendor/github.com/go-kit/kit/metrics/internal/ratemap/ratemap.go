



package ratemap

import "sync"


type RateMap struct {
	mtx sync.RWMutex
	m   map[string]float64
}


func New() *RateMap {
	return &RateMap{
		m: map[string]float64{},
	}
}



func (m *RateMap) Set(name string, rate float64) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.m[name] = rate
}



func (m *RateMap) Get(name string) float64 {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	f, ok := m.m[name]
	if !ok {
		f = 1.0
	}
	return f
}
