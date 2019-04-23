

package mocks

import metrics "github.com/mcc-github/blockchain/common/metrics"
import mock "github.com/stretchr/testify/mock"


type MetricsProvider struct {
	mock.Mock
}


func (_m *MetricsProvider) NewCounter(opts metrics.CounterOpts) metrics.Counter {
	ret := _m.Called(opts)

	var r0 metrics.Counter
	if rf, ok := ret.Get(0).(func(metrics.CounterOpts) metrics.Counter); ok {
		r0 = rf(opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(metrics.Counter)
		}
	}

	return r0
}


func (_m *MetricsProvider) NewGauge(opts metrics.GaugeOpts) metrics.Gauge {
	ret := _m.Called(opts)

	var r0 metrics.Gauge
	if rf, ok := ret.Get(0).(func(metrics.GaugeOpts) metrics.Gauge); ok {
		r0 = rf(opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(metrics.Gauge)
		}
	}

	return r0
}


func (_m *MetricsProvider) NewHistogram(opts metrics.HistogramOpts) metrics.Histogram {
	ret := _m.Called(opts)

	var r0 metrics.Histogram
	if rf, ok := ret.Get(0).(func(metrics.HistogramOpts) metrics.Histogram); ok {
		r0 = rf(opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(metrics.Histogram)
		}
	}

	return r0
}
