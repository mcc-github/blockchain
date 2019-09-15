
package metricsfakes

import (
	"sync"

	"github.com/mcc-github/blockchain/common/metrics"
)

type Gauge struct {
	AddStub        func(float64)
	addMutex       sync.RWMutex
	addArgsForCall []struct {
		arg1 float64
	}
	SetStub        func(float64)
	setMutex       sync.RWMutex
	setArgsForCall []struct {
		arg1 float64
	}
	WithStub        func(...string) metrics.Gauge
	withMutex       sync.RWMutex
	withArgsForCall []struct {
		arg1 []string
	}
	withReturns struct {
		result1 metrics.Gauge
	}
	withReturnsOnCall map[int]struct {
		result1 metrics.Gauge
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *Gauge) Add(arg1 float64) {
	fake.addMutex.Lock()
	fake.addArgsForCall = append(fake.addArgsForCall, struct {
		arg1 float64
	}{arg1})
	fake.recordInvocation("Add", []interface{}{arg1})
	fake.addMutex.Unlock()
	if fake.AddStub != nil {
		fake.AddStub(arg1)
	}
}

func (fake *Gauge) AddCallCount() int {
	fake.addMutex.RLock()
	defer fake.addMutex.RUnlock()
	return len(fake.addArgsForCall)
}

func (fake *Gauge) AddCalls(stub func(float64)) {
	fake.addMutex.Lock()
	defer fake.addMutex.Unlock()
	fake.AddStub = stub
}

func (fake *Gauge) AddArgsForCall(i int) float64 {
	fake.addMutex.RLock()
	defer fake.addMutex.RUnlock()
	argsForCall := fake.addArgsForCall[i]
	return argsForCall.arg1
}

func (fake *Gauge) Set(arg1 float64) {
	fake.setMutex.Lock()
	fake.setArgsForCall = append(fake.setArgsForCall, struct {
		arg1 float64
	}{arg1})
	fake.recordInvocation("Set", []interface{}{arg1})
	fake.setMutex.Unlock()
	if fake.SetStub != nil {
		fake.SetStub(arg1)
	}
}

func (fake *Gauge) SetCallCount() int {
	fake.setMutex.RLock()
	defer fake.setMutex.RUnlock()
	return len(fake.setArgsForCall)
}

func (fake *Gauge) SetCalls(stub func(float64)) {
	fake.setMutex.Lock()
	defer fake.setMutex.Unlock()
	fake.SetStub = stub
}

func (fake *Gauge) SetArgsForCall(i int) float64 {
	fake.setMutex.RLock()
	defer fake.setMutex.RUnlock()
	argsForCall := fake.setArgsForCall[i]
	return argsForCall.arg1
}

func (fake *Gauge) With(arg1 ...string) metrics.Gauge {
	fake.withMutex.Lock()
	ret, specificReturn := fake.withReturnsOnCall[len(fake.withArgsForCall)]
	fake.withArgsForCall = append(fake.withArgsForCall, struct {
		arg1 []string
	}{arg1})
	fake.recordInvocation("With", []interface{}{arg1})
	fake.withMutex.Unlock()
	if fake.WithStub != nil {
		return fake.WithStub(arg1...)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.withReturns
	return fakeReturns.result1
}

func (fake *Gauge) WithCallCount() int {
	fake.withMutex.RLock()
	defer fake.withMutex.RUnlock()
	return len(fake.withArgsForCall)
}

func (fake *Gauge) WithCalls(stub func(...string) metrics.Gauge) {
	fake.withMutex.Lock()
	defer fake.withMutex.Unlock()
	fake.WithStub = stub
}

func (fake *Gauge) WithArgsForCall(i int) []string {
	fake.withMutex.RLock()
	defer fake.withMutex.RUnlock()
	argsForCall := fake.withArgsForCall[i]
	return argsForCall.arg1
}

func (fake *Gauge) WithReturns(result1 metrics.Gauge) {
	fake.withMutex.Lock()
	defer fake.withMutex.Unlock()
	fake.WithStub = nil
	fake.withReturns = struct {
		result1 metrics.Gauge
	}{result1}
}

func (fake *Gauge) WithReturnsOnCall(i int, result1 metrics.Gauge) {
	fake.withMutex.Lock()
	defer fake.withMutex.Unlock()
	fake.WithStub = nil
	if fake.withReturnsOnCall == nil {
		fake.withReturnsOnCall = make(map[int]struct {
			result1 metrics.Gauge
		})
	}
	fake.withReturnsOnCall[i] = struct {
		result1 metrics.Gauge
	}{result1}
}

func (fake *Gauge) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.addMutex.RLock()
	defer fake.addMutex.RUnlock()
	fake.setMutex.RLock()
	defer fake.setMutex.RUnlock()
	fake.withMutex.RLock()
	defer fake.withMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *Gauge) recordInvocation(key string, args []interface{}) {
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

var _ metrics.Gauge = new(Gauge)
