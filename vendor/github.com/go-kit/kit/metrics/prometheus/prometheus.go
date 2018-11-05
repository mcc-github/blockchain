



package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/internal/lv"
)


type Counter struct {
	cv  *prometheus.CounterVec
	lvs lv.LabelValues
}



func NewCounterFrom(opts prometheus.CounterOpts, labelNames []string) *Counter {
	cv := prometheus.NewCounterVec(opts, labelNames)
	prometheus.MustRegister(cv)
	return NewCounter(cv)
}


func NewCounter(cv *prometheus.CounterVec) *Counter {
	return &Counter{
		cv: cv,
	}
}


func (c *Counter) With(labelValues ...string) metrics.Counter {
	return &Counter{
		cv:  c.cv,
		lvs: c.lvs.With(labelValues...),
	}
}


func (c *Counter) Add(delta float64) {
	c.cv.With(makeLabels(c.lvs...)).Add(delta)
}


type Gauge struct {
	gv  *prometheus.GaugeVec
	lvs lv.LabelValues
}



func NewGaugeFrom(opts prometheus.GaugeOpts, labelNames []string) *Gauge {
	gv := prometheus.NewGaugeVec(opts, labelNames)
	prometheus.MustRegister(gv)
	return NewGauge(gv)
}


func NewGauge(gv *prometheus.GaugeVec) *Gauge {
	return &Gauge{
		gv: gv,
	}
}


func (g *Gauge) With(labelValues ...string) metrics.Gauge {
	return &Gauge{
		gv:  g.gv,
		lvs: g.lvs.With(labelValues...),
	}
}


func (g *Gauge) Set(value float64) {
	g.gv.With(makeLabels(g.lvs...)).Set(value)
}


func (g *Gauge) Add(delta float64) {
	g.gv.With(makeLabels(g.lvs...)).Add(delta)
}




type Summary struct {
	sv  *prometheus.SummaryVec
	lvs lv.LabelValues
}



func NewSummaryFrom(opts prometheus.SummaryOpts, labelNames []string) *Summary {
	sv := prometheus.NewSummaryVec(opts, labelNames)
	prometheus.MustRegister(sv)
	return NewSummary(sv)
}


func NewSummary(sv *prometheus.SummaryVec) *Summary {
	return &Summary{
		sv: sv,
	}
}


func (s *Summary) With(labelValues ...string) metrics.Histogram {
	return &Summary{
		sv:  s.sv,
		lvs: s.lvs.With(labelValues...),
	}
}


func (s *Summary) Observe(value float64) {
	s.sv.With(makeLabels(s.lvs...)).Observe(value)
}




type Histogram struct {
	hv  *prometheus.HistogramVec
	lvs lv.LabelValues
}



func NewHistogramFrom(opts prometheus.HistogramOpts, labelNames []string) *Histogram {
	hv := prometheus.NewHistogramVec(opts, labelNames)
	prometheus.MustRegister(hv)
	return NewHistogram(hv)
}


func NewHistogram(hv *prometheus.HistogramVec) *Histogram {
	return &Histogram{
		hv: hv,
	}
}


func (h *Histogram) With(labelValues ...string) metrics.Histogram {
	return &Histogram{
		hv:  h.hv,
		lvs: h.lvs.With(labelValues...),
	}
}


func (h *Histogram) Observe(value float64) {
	h.hv.With(makeLabels(h.lvs...)).Observe(value)
}

func makeLabels(labelValues ...string) prometheus.Labels {
	labels := prometheus.Labels{}
	for i := 0; i < len(labelValues); i += 2 {
		labels[labelValues[i]] = labelValues[i+1]
	}
	return labels
}
