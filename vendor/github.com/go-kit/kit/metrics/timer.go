package metrics

import "time"



type Timer struct {
	h Histogram
	t time.Time
	u time.Duration
}


func NewTimer(h Histogram) *Timer {
	return &Timer{
		h: h,
		t: time.Now(),
		u: time.Second,
	}
}



func (t *Timer) ObserveDuration() {
	d := float64(time.Since(t.t).Nanoseconds()) / float64(t.u)
	if d < 0 {
		d = 0
	}
	t.h.Observe(d)
}



func (t *Timer) Unit(u time.Duration) {
	t.u = u
}
