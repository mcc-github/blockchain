












package prometheus

import (
	"errors"
	"math"
	"sync/atomic"

	dto "github.com/prometheus/client_model/go"
)










type Counter interface {
	Metric
	Collector

	
	
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
	result := &counter{desc: desc, labelPairs: desc.constLabelPairs}
	result.init(result) 
	return result
}

type counter struct {
	
	
	
	
	valBits uint64
	valInt  uint64

	selfCollector
	desc *Desc

	labelPairs []*dto.LabelPair
}

func (c *counter) Desc() *Desc {
	return c.desc
}

func (c *counter) Add(v float64) {
	if v < 0 {
		panic(errors.New("counter cannot decrease in value"))
	}
	ival := uint64(v)
	if float64(ival) == v {
		atomic.AddUint64(&c.valInt, ival)
		return
	}

	for {
		oldBits := atomic.LoadUint64(&c.valBits)
		newBits := math.Float64bits(math.Float64frombits(oldBits) + v)
		if atomic.CompareAndSwapUint64(&c.valBits, oldBits, newBits) {
			return
		}
	}
}

func (c *counter) Inc() {
	atomic.AddUint64(&c.valInt, 1)
}

func (c *counter) Write(out *dto.Metric) error {
	fval := math.Float64frombits(atomic.LoadUint64(&c.valBits))
	ival := atomic.LoadUint64(&c.valInt)
	val := fval + float64(ival)

	return populateMetric(CounterValue, val, c.labelPairs, out)
}






type CounterVec struct {
	*metricVec
}



func NewCounterVec(opts CounterOpts, labelNames []string) *CounterVec {
	desc := NewDesc(
		BuildFQName(opts.Namespace, opts.Subsystem, opts.Name),
		opts.Help,
		labelNames,
		opts.ConstLabels,
	)
	return &CounterVec{
		metricVec: newMetricVec(desc, func(lvs ...string) Metric {
			if len(lvs) != len(desc.variableLabels) {
				panic(errInconsistentCardinality)
			}
			result := &counter{desc: desc, labelPairs: makeLabelPairs(desc, lvs)}
			result.init(result) 
			return result
		}),
	}
}
























func (v *CounterVec) GetMetricWithLabelValues(lvs ...string) (Counter, error) {
	metric, err := v.metricVec.getMetricWithLabelValues(lvs...)
	if metric != nil {
		return metric.(Counter), err
	}
	return nil, err
}













func (v *CounterVec) GetMetricWith(labels Labels) (Counter, error) {
	metric, err := v.metricVec.getMetricWith(labels)
	if metric != nil {
		return metric.(Counter), err
	}
	return nil, err
}





func (v *CounterVec) WithLabelValues(lvs ...string) Counter {
	c, err := v.GetMetricWithLabelValues(lvs...)
	if err != nil {
		panic(err)
	}
	return c
}




func (v *CounterVec) With(labels Labels) Counter {
	c, err := v.GetMetricWith(labels)
	if err != nil {
		panic(err)
	}
	return c
}














func (v *CounterVec) CurryWith(labels Labels) (*CounterVec, error) {
	vec, err := v.curryWith(labels)
	if vec != nil {
		return &CounterVec{vec}, err
	}
	return nil, err
}



func (v *CounterVec) MustCurryWith(labels Labels) *CounterVec {
	vec, err := v.CurryWith(labels)
	if err != nil {
		panic(err)
	}
	return vec
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
