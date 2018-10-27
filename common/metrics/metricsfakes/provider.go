
package metricsfakes

import (
	"sync"

	metricskit "github.com/go-kit/kit/metrics"
	"github.com/mcc-github/blockchain/common/metrics"
)

type Provider struct {
	NewCounterStub        func(name string) metricskit.Counter
	newCounterMutex       sync.RWMutex
	newCounterArgsForCall []struct {
		name string
	}
	newCounterReturns struct {
		result1 metricskit.Counter
	}
	newCounterReturnsOnCall map[int]struct {
		result1 metricskit.Counter
	}
	NewGaugeStub        func(name string) metricskit.Gauge
	newGaugeMutex       sync.RWMutex
	newGaugeArgsForCall []struct {
		name string
	}
	newGaugeReturns struct {
		result1 metricskit.Gauge
	}
	newGaugeReturnsOnCall map[int]struct {
		result1 metricskit.Gauge
	}
	NewHistogramStub        func(name string) metricskit.Histogram
	newHistogramMutex       sync.RWMutex
	newHistogramArgsForCall []struct {
		name string
	}
	newHistogramReturns struct {
		result1 metricskit.Histogram
	}
	newHistogramReturnsOnCall map[int]struct {
		result1 metricskit.Histogram
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *Provider) NewCounter(name string) metricskit.Counter {
	fake.newCounterMutex.Lock()
	ret, specificReturn := fake.newCounterReturnsOnCall[len(fake.newCounterArgsForCall)]
	fake.newCounterArgsForCall = append(fake.newCounterArgsForCall, struct {
		name string
	}{name})
	fake.recordInvocation("NewCounter", []interface{}{name})
	fake.newCounterMutex.Unlock()
	if fake.NewCounterStub != nil {
		return fake.NewCounterStub(name)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.newCounterReturns.result1
}

func (fake *Provider) NewCounterCallCount() int {
	fake.newCounterMutex.RLock()
	defer fake.newCounterMutex.RUnlock()
	return len(fake.newCounterArgsForCall)
}

func (fake *Provider) NewCounterArgsForCall(i int) string {
	fake.newCounterMutex.RLock()
	defer fake.newCounterMutex.RUnlock()
	return fake.newCounterArgsForCall[i].name
}

func (fake *Provider) NewCounterReturns(result1 metricskit.Counter) {
	fake.NewCounterStub = nil
	fake.newCounterReturns = struct {
		result1 metricskit.Counter
	}{result1}
}

func (fake *Provider) NewCounterReturnsOnCall(i int, result1 metricskit.Counter) {
	fake.NewCounterStub = nil
	if fake.newCounterReturnsOnCall == nil {
		fake.newCounterReturnsOnCall = make(map[int]struct {
			result1 metricskit.Counter
		})
	}
	fake.newCounterReturnsOnCall[i] = struct {
		result1 metricskit.Counter
	}{result1}
}

func (fake *Provider) NewGauge(name string) metricskit.Gauge {
	fake.newGaugeMutex.Lock()
	ret, specificReturn := fake.newGaugeReturnsOnCall[len(fake.newGaugeArgsForCall)]
	fake.newGaugeArgsForCall = append(fake.newGaugeArgsForCall, struct {
		name string
	}{name})
	fake.recordInvocation("NewGauge", []interface{}{name})
	fake.newGaugeMutex.Unlock()
	if fake.NewGaugeStub != nil {
		return fake.NewGaugeStub(name)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.newGaugeReturns.result1
}

func (fake *Provider) NewGaugeCallCount() int {
	fake.newGaugeMutex.RLock()
	defer fake.newGaugeMutex.RUnlock()
	return len(fake.newGaugeArgsForCall)
}

func (fake *Provider) NewGaugeArgsForCall(i int) string {
	fake.newGaugeMutex.RLock()
	defer fake.newGaugeMutex.RUnlock()
	return fake.newGaugeArgsForCall[i].name
}

func (fake *Provider) NewGaugeReturns(result1 metricskit.Gauge) {
	fake.NewGaugeStub = nil
	fake.newGaugeReturns = struct {
		result1 metricskit.Gauge
	}{result1}
}

func (fake *Provider) NewGaugeReturnsOnCall(i int, result1 metricskit.Gauge) {
	fake.NewGaugeStub = nil
	if fake.newGaugeReturnsOnCall == nil {
		fake.newGaugeReturnsOnCall = make(map[int]struct {
			result1 metricskit.Gauge
		})
	}
	fake.newGaugeReturnsOnCall[i] = struct {
		result1 metricskit.Gauge
	}{result1}
}

func (fake *Provider) NewHistogram(name string) metricskit.Histogram {
	fake.newHistogramMutex.Lock()
	ret, specificReturn := fake.newHistogramReturnsOnCall[len(fake.newHistogramArgsForCall)]
	fake.newHistogramArgsForCall = append(fake.newHistogramArgsForCall, struct {
		name string
	}{name})
	fake.recordInvocation("NewHistogram", []interface{}{name})
	fake.newHistogramMutex.Unlock()
	if fake.NewHistogramStub != nil {
		return fake.NewHistogramStub(name)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.newHistogramReturns.result1
}

func (fake *Provider) NewHistogramCallCount() int {
	fake.newHistogramMutex.RLock()
	defer fake.newHistogramMutex.RUnlock()
	return len(fake.newHistogramArgsForCall)
}

func (fake *Provider) NewHistogramArgsForCall(i int) string {
	fake.newHistogramMutex.RLock()
	defer fake.newHistogramMutex.RUnlock()
	return fake.newHistogramArgsForCall[i].name
}

func (fake *Provider) NewHistogramReturns(result1 metricskit.Histogram) {
	fake.NewHistogramStub = nil
	fake.newHistogramReturns = struct {
		result1 metricskit.Histogram
	}{result1}
}

func (fake *Provider) NewHistogramReturnsOnCall(i int, result1 metricskit.Histogram) {
	fake.NewHistogramStub = nil
	if fake.newHistogramReturnsOnCall == nil {
		fake.newHistogramReturnsOnCall = make(map[int]struct {
			result1 metricskit.Histogram
		})
	}
	fake.newHistogramReturnsOnCall[i] = struct {
		result1 metricskit.Histogram
	}{result1}
}

func (fake *Provider) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.newCounterMutex.RLock()
	defer fake.newCounterMutex.RUnlock()
	fake.newGaugeMutex.RLock()
	defer fake.newGaugeMutex.RUnlock()
	fake.newHistogramMutex.RLock()
	defer fake.newHistogramMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *Provider) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ metrics.Provider = new(Provider)
