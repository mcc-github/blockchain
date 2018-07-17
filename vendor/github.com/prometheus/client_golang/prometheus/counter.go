












package prometheus

import (
	"errors"
)










type Counter interface {
	Metric
	Collector

	
	
	
	
	
	
	
	
	Set(float64)
	
	Inc()
	
	
	Add(float64)
}


type CounterOpts Opts


func NewCounter(opts CounterOpts) Counter {
	desc := NewDesc(
		BuildFQName(opts.Namespace, opts.Subsystem, opts.Name),
		opts.Help,
		nil,
		opts.ConstLabels,
	)
	result := &counter{value: value{desc: desc, valType: CounterValue, labelPairs: desc.constLabelPairs}}
	result.init(result) 
	return result
}

type counter struct {
	value
}

func (c *counter) Add(v float64) {
	if v < 0 {
		panic(errors.New("counter cannot decrease in value"))
	}
	c.value.Add(v)
}









type CounterVec struct {
	*MetricVec
}




func NewCounterVec(opts CounterOpts, labelNames []string) *CounterVec {
	desc := NewDesc(
		BuildFQName(opts.Namespace, opts.Subsystem, opts.Name),
		opts.Help,
		labelNames,
		opts.ConstLabels,
	)
	return &CounterVec{
		MetricVec: newMetricVec(desc, func(lvs ...string) Metric {
			result := &counter{value: value{
				desc:       desc,
				valType:    CounterValue,
				labelPairs: makeLabelPairs(desc, lvs),
			}}
			result.init(result) 
			return result
		}),
	}
}




func (m *CounterVec) GetMetricWithLabelValues(lvs ...string) (Counter, error) {
	metric, err := m.MetricVec.GetMetricWithLabelValues(lvs...)
	if metric != nil {
		return metric.(Counter), err
	}
	return nil, err
}




func (m *CounterVec) GetMetricWith(labels Labels) (Counter, error) {
	metric, err := m.MetricVec.GetMetricWith(labels)
	if metric != nil {
		return metric.(Counter), err
	}
	return nil, err
}





func (m *CounterVec) WithLabelValues(lvs ...string) Counter {
	return m.MetricVec.WithLabelValues(lvs...).(Counter)
}




func (m *CounterVec) With(labels Labels) Counter {
	return m.MetricVec.With(labels).(Counter)
}





type CounterFunc interface {
	Metric
	Collector
}









func NewCounterFunc(opts CounterOpts, function func() float64) CounterFunc {
	return newValueFunc(NewDesc(
		BuildFQName(opts.Namespace, opts.Subsystem, opts.Name),
		opts.Help,
		nil,
		opts.ConstLabels,
	), CounterValue, function)
}
