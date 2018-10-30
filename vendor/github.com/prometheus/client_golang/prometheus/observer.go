












package prometheus



type Observer interface {
	Observe(float64)
}














type ObserverFunc func(float64)


func (f ObserverFunc) Observe(value float64) {
	f(value)
}


type ObserverVec interface {
	GetMetricWith(Labels) (Observer, error)
	GetMetricWithLabelValues(lvs ...string) (Observer, error)
	With(Labels) Observer
	WithLabelValues(...string) Observer
	CurryWith(Labels) (ObserverVec, error)
	MustCurryWith(Labels) ObserverVec

	Collector
}
