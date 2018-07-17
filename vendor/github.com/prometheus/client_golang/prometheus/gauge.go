












package prometheus









type Gauge interface {
	Metric
	Collector

	
	Set(float64)
	
	Inc()
	
	Dec()
	
	
	Add(float64)
	
	
	Sub(float64)
}


type GaugeOpts Opts


func NewGauge(opts GaugeOpts) Gauge {
	return newValue(NewDesc(
		BuildFQName(opts.Namespace, opts.Subsystem, opts.Name),
		opts.Help,
		nil,
		opts.ConstLabels,
	), GaugeValue, 0)
}






type GaugeVec struct {
	*MetricVec
}




func NewGaugeVec(opts GaugeOpts, labelNames []string) *GaugeVec {
	desc := NewDesc(
		BuildFQName(opts.Namespace, opts.Subsystem, opts.Name),
		opts.Help,
		labelNames,
		opts.ConstLabels,
	)
	return &GaugeVec{
		MetricVec: newMetricVec(desc, func(lvs ...string) Metric {
			return newValue(desc, GaugeValue, 0, lvs...)
		}),
	}
}




func (m *GaugeVec) GetMetricWithLabelValues(lvs ...string) (Gauge, error) {
	metric, err := m.MetricVec.GetMetricWithLabelValues(lvs...)
	if metric != nil {
		return metric.(Gauge), err
	}
	return nil, err
}




func (m *GaugeVec) GetMetricWith(labels Labels) (Gauge, error) {
	metric, err := m.MetricVec.GetMetricWith(labels)
	if metric != nil {
		return metric.(Gauge), err
	}
	return nil, err
}





func (m *GaugeVec) WithLabelValues(lvs ...string) Gauge {
	return m.MetricVec.WithLabelValues(lvs...).(Gauge)
}




func (m *GaugeVec) With(labels Labels) Gauge {
	return m.MetricVec.With(labels).(Gauge)
}





type GaugeFunc interface {
	Metric
	Collector
}







func NewGaugeFunc(opts GaugeOpts, function func() float64) GaugeFunc {
	return newValueFunc(NewDesc(
		BuildFQName(opts.Namespace, opts.Subsystem, opts.Name),
		opts.Help,
		nil,
		opts.ConstLabels,
	), GaugeValue, function)
}
