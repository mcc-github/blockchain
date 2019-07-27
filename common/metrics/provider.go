/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package metrics



type Provider interface {
	
	NewCounter(CounterOpts) Counter
	
	NewGauge(GaugeOpts) Gauge
	
	NewHistogram(HistogramOpts) Histogram
}


type Counter interface {
	
	
	With(labelValues ...string) Counter

	
	Add(delta float64)
}



type CounterOpts struct {
	
	
	
	
	Namespace string
	Subsystem string
	Name      string

	
	Help string

	
	
	
	LabelNames []string

	
	
	LabelHelp map[string]string

	
	
	
	
	
	
	
	
	
	
	
	
	StatsdFormat string
}


type Gauge interface {
	
	
	With(labelValues ...string) Gauge

	
	Add(delta float64) 

	
	Set(value float64)
}



type GaugeOpts struct {
	
	
	
	
	Namespace string
	Subsystem string
	Name      string

	
	Help string

	
	
	
	LabelNames []string

	
	
	LabelHelp map[string]string

	
	
	
	
	
	
	
	
	
	
	
	
	StatsdFormat string
}



type Histogram interface {
	
	
	
	With(labelValues ...string) Histogram
	Observe(value float64)
}



type HistogramOpts struct {
	
	
	
	
	Namespace string
	Subsystem string
	Name      string

	
	Help string

	
	
	Buckets []float64

	
	
	
	LabelNames []string

	
	
	LabelHelp map[string]string

	
	
	
	
	
	
	
	
	
	
	
	
	StatsdFormat string
}
