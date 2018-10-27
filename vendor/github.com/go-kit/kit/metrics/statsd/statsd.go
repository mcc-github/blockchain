







package statsd

import (
	"fmt"
	"io"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/internal/lv"
	"github.com/go-kit/kit/metrics/internal/ratemap"
	"github.com/go-kit/kit/util/conn"
)











type Statsd struct {
	prefix string
	rates  *ratemap.RateMap

	
	
	
	
	
	counters *lv.Space
	gauges   *lv.Space
	timings  *lv.Space

	logger log.Logger
}




func New(prefix string, logger log.Logger) *Statsd {
	return &Statsd{
		prefix:   prefix,
		rates:    ratemap.New(),
		counters: lv.NewSpace(),
		gauges:   lv.NewSpace(),
		timings:  lv.NewSpace(),
		logger:   logger,
	}
}


func (s *Statsd) NewCounter(name string, sampleRate float64) *Counter {
	s.rates.Set(s.prefix+name, sampleRate)
	return &Counter{
		name: s.prefix + name,
		obs:  s.counters.Observe,
	}
}


func (s *Statsd) NewGauge(name string) *Gauge {
	return &Gauge{
		name: s.prefix + name,
		obs:  s.gauges.Observe,
		add:  s.gauges.Add,
	}
}



func (s *Statsd) NewTiming(name string, sampleRate float64) *Timing {
	s.rates.Set(s.prefix+name, sampleRate)
	return &Timing{
		name: s.prefix + name,
		obs:  s.timings.Observe,
	}
}





func (s *Statsd) WriteLoop(c <-chan time.Time, w io.Writer) {
	for range c {
		if _, err := s.WriteTo(w); err != nil {
			s.logger.Log("during", "WriteTo", "err", err)
		}
	}
}






func (s *Statsd) SendLoop(c <-chan time.Time, network, address string) {
	s.WriteLoop(c, conn.NewDefaultManager(network, address, s.logger))
}





func (s *Statsd) WriteTo(w io.Writer) (count int64, err error) {
	var n int

	s.counters.Reset().Walk(func(name string, _ lv.LabelValues, values []float64) bool {
		n, err = fmt.Fprintf(w, "%s:%f|c%s\n", name, sum(values), sampling(s.rates.Get(name)))
		if err != nil {
			return false
		}
		count += int64(n)
		return true
	})
	if err != nil {
		return count, err
	}

	s.gauges.Reset().Walk(func(name string, _ lv.LabelValues, values []float64) bool {
		n, err = fmt.Fprintf(w, "%s:%f|g\n", name, last(values))
		if err != nil {
			return false
		}
		count += int64(n)
		return true
	})
	if err != nil {
		return count, err
	}

	s.timings.Reset().Walk(func(name string, _ lv.LabelValues, values []float64) bool {
		sampleRate := s.rates.Get(name)
		for _, value := range values {
			n, err = fmt.Fprintf(w, "%s:%f|ms%s\n", name, value, sampling(sampleRate))
			if err != nil {
				return false
			}
			count += int64(n)
		}
		return true
	})
	if err != nil {
		return count, err
	}

	return count, err
}

func sum(a []float64) float64 {
	var v float64
	for _, f := range a {
		v += f
	}
	return v
}

func last(a []float64) float64 {
	return a[len(a)-1]
}

func sampling(r float64) string {
	var sv string
	if r < 1.0 {
		sv = fmt.Sprintf("|@%f", r)
	}
	return sv
}

type observeFunc func(name string, lvs lv.LabelValues, value float64)



type Counter struct {
	name string
	obs  observeFunc
}


func (c *Counter) With(...string) metrics.Counter {
	return c
}


func (c *Counter) Add(delta float64) {
	c.obs(c.name, lv.LabelValues{}, delta)
}



type Gauge struct {
	name string
	obs  observeFunc
	add  observeFunc
}


func (g *Gauge) With(...string) metrics.Gauge {
	return g
}


func (g *Gauge) Set(value float64) {
	g.obs(g.name, lv.LabelValues{}, value)
}


func (g *Gauge) Add(delta float64) {
	g.add(g.name, lv.LabelValues{}, delta)
}




type Timing struct {
	name string
	obs  observeFunc
}


func (t *Timing) With(...string) metrics.Histogram {
	return t
}


func (t *Timing) Observe(value float64) {
	t.obs(t.name, lv.LabelValues{}, value)
}
