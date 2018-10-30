












package prometheus

import (
	"math"
	"sync/atomic"
	"time"

	dto "github.com/prometheus/client_model/go"
)









type Gauge interface {
	Metric
	Collector

	
	Set(float64)
	
	
	Inc()
	
	
	Dec()
	
	
	Add(float64)
	
	
	Sub(float64)

	
	SetToCurrentTime()
}


type GaugeOpts Opts









func NewGauge(opts GaugeOpts) Gauge {
	desc := NewDesc(
		BuildFQName(opts.Namespace, opts.Subsystem, opts.Name),
		opts.Help,
		nil,
		opts.ConstLabels,
	)
	result := &gauge{desc: desc, labelPairs: desc.constLabelPairs}
	result.init(result) 
	return result
}

type gauge struct {
	
	
	
	valBits uint64

	selfCollector

	desc       *Desc
	labelPairs []*dto.LabelPair
}

func (g *gauge) Desc() *Desc {
	return g.desc
}

func (g *gauge) Set(val float64) {
	atomic.StoreUint64(&g.valBits, math.Float64bits(val))
}

func (g *gauge) SetToCurrentTime() {
	g.Set(float64(time.Now().UnixNano()) / 1e9)
}

func (g *gauge) Inc() {
	g.Add(1)
}

func (g *gauge) Dec() {
	g.Add(-1)
}

func (g *gauge) Add(val float64) {
	for {
		oldBits := atomic.LoadUint64(&g.valBits)
		newBits := math.Float64bits(math.Float64frombits(oldBits) + val)
		if atomic.CompareAndSwapUint64(&g.valBits, oldBits, newBits) {
			return
		}
	}
}

func (g *gauge) Sub(val float64) {
	g.Add(val * -1)
}

func (g *gauge) Write(out *dto.Metric) error {
	val := math.Float64frombits(atomic.LoadUint64(&g.valBits))
	return populateMetric(GaugeValue, val, g.labelPairs, out)
}






type GaugeVec struct {
	*metricVec
}



func NewGaugeVec(opts GaugeOpts, labelNames []string) *GaugeVec {
	desc := NewDesc(
		BuildFQName(opts.Namespace, opts.Subsystem, opts.Name),
		opts.Help,
		labelNames,
		opts.ConstLabels,
	)
	return &GaugeVec{
		metricVec: newMetricVec(desc, func(lvs ...string) Metric {
			if len(lvs) != len(desc.variableLabels) {
				panic(errInconsistentCardinality)
			}
			result := &gauge{desc: desc, labelPairs: makeLabelPairs(desc, lvs)}
			result.init(result) 
			return result
		}),
	}
}
























func (v *GaugeVec) GetMetricWithLabelValues(lvs ...string) (Gauge, error) {
	metric, err := v.metricVec.getMetricWithLabelValues(lvs...)
	if metric != nil {
		return metric.(Gauge), err
	}
	return nil, err
}













func (v *GaugeVec) GetMetricWith(labels Labels) (Gauge, error) {
	metric, err := v.metricVec.getMetricWith(labels)
	if metric != nil {
		return metric.(Gauge), err
	}
	return nil, err
}





func (v *GaugeVec) WithLabelValues(lvs ...string) Gauge {
	g, err := v.GetMetricWithLabelValues(lvs...)
	if err != nil {
		panic(err)
	}
	return g
}




func (v *GaugeVec) With(labels Labels) Gauge {
	g, err := v.GetMetricWith(labels)
	if err != nil {
		panic(err)
	}
	return g
}














func (v *GaugeVec) CurryWith(labels Labels) (*GaugeVec, error) {
	vec, err := v.curryWith(labels)
	if vec != nil {
		return &GaugeVec{vec}, err
	}
	return nil, err
}



func (v *GaugeVec) MustCurryWith(labels Labels) *GaugeVec {
	vec, err := v.CurryWith(labels)
	if err != nil {
		panic(err)
	}
	return vec
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
