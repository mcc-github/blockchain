












package prometheus








type Untyped interface {
	Metric
	Collector

	
	Set(float64)
	
	Inc()
	
	Dec()
	
	
	Add(float64)
	
	
	Sub(float64)
}


type UntypedOpts Opts


func NewUntyped(opts UntypedOpts) Untyped {
	return newValue(NewDesc(
		BuildFQName(opts.Namespace, opts.Subsystem, opts.Name),
		opts.Help,
		nil,
		opts.ConstLabels,
	), UntypedValue, 0)
}





type UntypedVec struct {
	*MetricVec
}




func NewUntypedVec(opts UntypedOpts, labelNames []string) *UntypedVec {
	desc := NewDesc(
		BuildFQName(opts.Namespace, opts.Subsystem, opts.Name),
		opts.Help,
		labelNames,
		opts.ConstLabels,
	)
	return &UntypedVec{
		MetricVec: newMetricVec(desc, func(lvs ...string) Metric {
			return newValue(desc, UntypedValue, 0, lvs...)
		}),
	}
}




func (m *UntypedVec) GetMetricWithLabelValues(lvs ...string) (Untyped, error) {
	metric, err := m.MetricVec.GetMetricWithLabelValues(lvs...)
	if metric != nil {
		return metric.(Untyped), err
	}
	return nil, err
}




func (m *UntypedVec) GetMetricWith(labels Labels) (Untyped, error) {
	metric, err := m.MetricVec.GetMetricWith(labels)
	if metric != nil {
		return metric.(Untyped), err
	}
	return nil, err
}





func (m *UntypedVec) WithLabelValues(lvs ...string) Untyped {
	return m.MetricVec.WithLabelValues(lvs...).(Untyped)
}




func (m *UntypedVec) With(labels Labels) Untyped {
	return m.MetricVec.With(labels).(Untyped)
}





type UntypedFunc interface {
	Metric
	Collector
}







func NewUntypedFunc(opts UntypedOpts, function func() float64) UntypedFunc {
	return newValueFunc(NewDesc(
		BuildFQName(opts.Namespace, opts.Subsystem, opts.Name),
		opts.Help,
		nil,
		opts.ConstLabels,
	), UntypedValue, function)
}
