












package prometheus

import "time"



type Timer struct {
	begin    time.Time
	observer Observer
}









func NewTimer(o Observer) *Timer {
	return &Timer{
		begin:    time.Now(),
		observer: o,
	}
}








func (t *Timer) ObserveDuration() {
	if t.observer != nil {
		t.observer.Observe(time.Since(t.begin).Seconds())
	}
}
