
package mock

import (
	sync "sync"

	metrics "github.com/rcrowley/go-metrics"
)

type MetricsMeter struct {
	CountStub        func() int64
	countMutex       sync.RWMutex
	countArgsForCall []struct {
	}
	countReturns struct {
		result1 int64
	}
	countReturnsOnCall map[int]struct {
		result1 int64
	}
	MarkStub        func(int64)
	markMutex       sync.RWMutex
	markArgsForCall []struct {
		arg1 int64
	}
	Rate1Stub        func() float64
	rate1Mutex       sync.RWMutex
	rate1ArgsForCall []struct {
	}
	rate1Returns struct {
		result1 float64
	}
	rate1ReturnsOnCall map[int]struct {
		result1 float64
	}
	Rate15Stub        func() float64
	rate15Mutex       sync.RWMutex
	rate15ArgsForCall []struct {
	}
	rate15Returns struct {
		result1 float64
	}
	rate15ReturnsOnCall map[int]struct {
		result1 float64
	}
	Rate5Stub        func() float64
	rate5Mutex       sync.RWMutex
	rate5ArgsForCall []struct {
	}
	rate5Returns struct {
		result1 float64
	}
	rate5ReturnsOnCall map[int]struct {
		result1 float64
	}
	RateMeanStub        func() float64
	rateMeanMutex       sync.RWMutex
	rateMeanArgsForCall []struct {
	}
	rateMeanReturns struct {
		result1 float64
	}
	rateMeanReturnsOnCall map[int]struct {
		result1 float64
	}
	SnapshotStub        func() metrics.Meter
	snapshotMutex       sync.RWMutex
	snapshotArgsForCall []struct {
	}
	snapshotReturns struct {
		result1 metrics.Meter
	}
	snapshotReturnsOnCall map[int]struct {
		result1 metrics.Meter
	}
	StopStub        func()
	stopMutex       sync.RWMutex
	stopArgsForCall []struct {
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *MetricsMeter) Count() int64 {
	fake.countMutex.Lock()
	ret, specificReturn := fake.countReturnsOnCall[len(fake.countArgsForCall)]
	fake.countArgsForCall = append(fake.countArgsForCall, struct {
	}{})
	fake.recordInvocation("Count", []interface{}{})
	fake.countMutex.Unlock()
	if fake.CountStub != nil {
		return fake.CountStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.countReturns
	return fakeReturns.result1
}

func (fake *MetricsMeter) CountCallCount() int {
	fake.countMutex.RLock()
	defer fake.countMutex.RUnlock()
	return len(fake.countArgsForCall)
}

func (fake *MetricsMeter) CountCalls(stub func() int64) {
	fake.countMutex.Lock()
	defer fake.countMutex.Unlock()
	fake.CountStub = stub
}

func (fake *MetricsMeter) CountReturns(result1 int64) {
	fake.countMutex.Lock()
	defer fake.countMutex.Unlock()
	fake.CountStub = nil
	fake.countReturns = struct {
		result1 int64
	}{result1}
}

func (fake *MetricsMeter) CountReturnsOnCall(i int, result1 int64) {
	fake.countMutex.Lock()
	defer fake.countMutex.Unlock()
	fake.CountStub = nil
	if fake.countReturnsOnCall == nil {
		fake.countReturnsOnCall = make(map[int]struct {
			result1 int64
		})
	}
	fake.countReturnsOnCall[i] = struct {
		result1 int64
	}{result1}
}

func (fake *MetricsMeter) Mark(arg1 int64) {
	fake.markMutex.Lock()
	fake.markArgsForCall = append(fake.markArgsForCall, struct {
		arg1 int64
	}{arg1})
	fake.recordInvocation("Mark", []interface{}{arg1})
	fake.markMutex.Unlock()
	if fake.MarkStub != nil {
		fake.MarkStub(arg1)
	}
}

func (fake *MetricsMeter) MarkCallCount() int {
	fake.markMutex.RLock()
	defer fake.markMutex.RUnlock()
	return len(fake.markArgsForCall)
}

func (fake *MetricsMeter) MarkCalls(stub func(int64)) {
	fake.markMutex.Lock()
	defer fake.markMutex.Unlock()
	fake.MarkStub = stub
}

func (fake *MetricsMeter) MarkArgsForCall(i int) int64 {
	fake.markMutex.RLock()
	defer fake.markMutex.RUnlock()
	argsForCall := fake.markArgsForCall[i]
	return argsForCall.arg1
}

func (fake *MetricsMeter) Rate1() float64 {
	fake.rate1Mutex.Lock()
	ret, specificReturn := fake.rate1ReturnsOnCall[len(fake.rate1ArgsForCall)]
	fake.rate1ArgsForCall = append(fake.rate1ArgsForCall, struct {
	}{})
	fake.recordInvocation("Rate1", []interface{}{})
	fake.rate1Mutex.Unlock()
	if fake.Rate1Stub != nil {
		return fake.Rate1Stub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.rate1Returns
	return fakeReturns.result1
}

func (fake *MetricsMeter) Rate1CallCount() int {
	fake.rate1Mutex.RLock()
	defer fake.rate1Mutex.RUnlock()
	return len(fake.rate1ArgsForCall)
}

func (fake *MetricsMeter) Rate1Calls(stub func() float64) {
	fake.rate1Mutex.Lock()
	defer fake.rate1Mutex.Unlock()
	fake.Rate1Stub = stub
}

func (fake *MetricsMeter) Rate1Returns(result1 float64) {
	fake.rate1Mutex.Lock()
	defer fake.rate1Mutex.Unlock()
	fake.Rate1Stub = nil
	fake.rate1Returns = struct {
		result1 float64
	}{result1}
}

func (fake *MetricsMeter) Rate1ReturnsOnCall(i int, result1 float64) {
	fake.rate1Mutex.Lock()
	defer fake.rate1Mutex.Unlock()
	fake.Rate1Stub = nil
	if fake.rate1ReturnsOnCall == nil {
		fake.rate1ReturnsOnCall = make(map[int]struct {
			result1 float64
		})
	}
	fake.rate1ReturnsOnCall[i] = struct {
		result1 float64
	}{result1}
}

func (fake *MetricsMeter) Rate15() float64 {
	fake.rate15Mutex.Lock()
	ret, specificReturn := fake.rate15ReturnsOnCall[len(fake.rate15ArgsForCall)]
	fake.rate15ArgsForCall = append(fake.rate15ArgsForCall, struct {
	}{})
	fake.recordInvocation("Rate15", []interface{}{})
	fake.rate15Mutex.Unlock()
	if fake.Rate15Stub != nil {
		return fake.Rate15Stub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.rate15Returns
	return fakeReturns.result1
}

func (fake *MetricsMeter) Rate15CallCount() int {
	fake.rate15Mutex.RLock()
	defer fake.rate15Mutex.RUnlock()
	return len(fake.rate15ArgsForCall)
}

func (fake *MetricsMeter) Rate15Calls(stub func() float64) {
	fake.rate15Mutex.Lock()
	defer fake.rate15Mutex.Unlock()
	fake.Rate15Stub = stub
}

func (fake *MetricsMeter) Rate15Returns(result1 float64) {
	fake.rate15Mutex.Lock()
	defer fake.rate15Mutex.Unlock()
	fake.Rate15Stub = nil
	fake.rate15Returns = struct {
		result1 float64
	}{result1}
}

func (fake *MetricsMeter) Rate15ReturnsOnCall(i int, result1 float64) {
	fake.rate15Mutex.Lock()
	defer fake.rate15Mutex.Unlock()
	fake.Rate15Stub = nil
	if fake.rate15ReturnsOnCall == nil {
		fake.rate15ReturnsOnCall = make(map[int]struct {
			result1 float64
		})
	}
	fake.rate15ReturnsOnCall[i] = struct {
		result1 float64
	}{result1}
}

func (fake *MetricsMeter) Rate5() float64 {
	fake.rate5Mutex.Lock()
	ret, specificReturn := fake.rate5ReturnsOnCall[len(fake.rate5ArgsForCall)]
	fake.rate5ArgsForCall = append(fake.rate5ArgsForCall, struct {
	}{})
	fake.recordInvocation("Rate5", []interface{}{})
	fake.rate5Mutex.Unlock()
	if fake.Rate5Stub != nil {
		return fake.Rate5Stub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.rate5Returns
	return fakeReturns.result1
}

func (fake *MetricsMeter) Rate5CallCount() int {
	fake.rate5Mutex.RLock()
	defer fake.rate5Mutex.RUnlock()
	return len(fake.rate5ArgsForCall)
}

func (fake *MetricsMeter) Rate5Calls(stub func() float64) {
	fake.rate5Mutex.Lock()
	defer fake.rate5Mutex.Unlock()
	fake.Rate5Stub = stub
}

func (fake *MetricsMeter) Rate5Returns(result1 float64) {
	fake.rate5Mutex.Lock()
	defer fake.rate5Mutex.Unlock()
	fake.Rate5Stub = nil
	fake.rate5Returns = struct {
		result1 float64
	}{result1}
}

func (fake *MetricsMeter) Rate5ReturnsOnCall(i int, result1 float64) {
	fake.rate5Mutex.Lock()
	defer fake.rate5Mutex.Unlock()
	fake.Rate5Stub = nil
	if fake.rate5ReturnsOnCall == nil {
		fake.rate5ReturnsOnCall = make(map[int]struct {
			result1 float64
		})
	}
	fake.rate5ReturnsOnCall[i] = struct {
		result1 float64
	}{result1}
}

func (fake *MetricsMeter) RateMean() float64 {
	fake.rateMeanMutex.Lock()
	ret, specificReturn := fake.rateMeanReturnsOnCall[len(fake.rateMeanArgsForCall)]
	fake.rateMeanArgsForCall = append(fake.rateMeanArgsForCall, struct {
	}{})
	fake.recordInvocation("RateMean", []interface{}{})
	fake.rateMeanMutex.Unlock()
	if fake.RateMeanStub != nil {
		return fake.RateMeanStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.rateMeanReturns
	return fakeReturns.result1
}

func (fake *MetricsMeter) RateMeanCallCount() int {
	fake.rateMeanMutex.RLock()
	defer fake.rateMeanMutex.RUnlock()
	return len(fake.rateMeanArgsForCall)
}

func (fake *MetricsMeter) RateMeanCalls(stub func() float64) {
	fake.rateMeanMutex.Lock()
	defer fake.rateMeanMutex.Unlock()
	fake.RateMeanStub = stub
}

func (fake *MetricsMeter) RateMeanReturns(result1 float64) {
	fake.rateMeanMutex.Lock()
	defer fake.rateMeanMutex.Unlock()
	fake.RateMeanStub = nil
	fake.rateMeanReturns = struct {
		result1 float64
	}{result1}
}

func (fake *MetricsMeter) RateMeanReturnsOnCall(i int, result1 float64) {
	fake.rateMeanMutex.Lock()
	defer fake.rateMeanMutex.Unlock()
	fake.RateMeanStub = nil
	if fake.rateMeanReturnsOnCall == nil {
		fake.rateMeanReturnsOnCall = make(map[int]struct {
			result1 float64
		})
	}
	fake.rateMeanReturnsOnCall[i] = struct {
		result1 float64
	}{result1}
}

func (fake *MetricsMeter) Snapshot() metrics.Meter {
	fake.snapshotMutex.Lock()
	ret, specificReturn := fake.snapshotReturnsOnCall[len(fake.snapshotArgsForCall)]
	fake.snapshotArgsForCall = append(fake.snapshotArgsForCall, struct {
	}{})
	fake.recordInvocation("Snapshot", []interface{}{})
	fake.snapshotMutex.Unlock()
	if fake.SnapshotStub != nil {
		return fake.SnapshotStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.snapshotReturns
	return fakeReturns.result1
}

func (fake *MetricsMeter) SnapshotCallCount() int {
	fake.snapshotMutex.RLock()
	defer fake.snapshotMutex.RUnlock()
	return len(fake.snapshotArgsForCall)
}

func (fake *MetricsMeter) SnapshotCalls(stub func() metrics.Meter) {
	fake.snapshotMutex.Lock()
	defer fake.snapshotMutex.Unlock()
	fake.SnapshotStub = stub
}

func (fake *MetricsMeter) SnapshotReturns(result1 metrics.Meter) {
	fake.snapshotMutex.Lock()
	defer fake.snapshotMutex.Unlock()
	fake.SnapshotStub = nil
	fake.snapshotReturns = struct {
		result1 metrics.Meter
	}{result1}
}

func (fake *MetricsMeter) SnapshotReturnsOnCall(i int, result1 metrics.Meter) {
	fake.snapshotMutex.Lock()
	defer fake.snapshotMutex.Unlock()
	fake.SnapshotStub = nil
	if fake.snapshotReturnsOnCall == nil {
		fake.snapshotReturnsOnCall = make(map[int]struct {
			result1 metrics.Meter
		})
	}
	fake.snapshotReturnsOnCall[i] = struct {
		result1 metrics.Meter
	}{result1}
}

func (fake *MetricsMeter) Stop() {
	fake.stopMutex.Lock()
	fake.stopArgsForCall = append(fake.stopArgsForCall, struct {
	}{})
	fake.recordInvocation("Stop", []interface{}{})
	fake.stopMutex.Unlock()
	if fake.StopStub != nil {
		fake.StopStub()
	}
}

func (fake *MetricsMeter) StopCallCount() int {
	fake.stopMutex.RLock()
	defer fake.stopMutex.RUnlock()
	return len(fake.stopArgsForCall)
}

func (fake *MetricsMeter) StopCalls(stub func()) {
	fake.stopMutex.Lock()
	defer fake.stopMutex.Unlock()
	fake.StopStub = stub
}

func (fake *MetricsMeter) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.countMutex.RLock()
	defer fake.countMutex.RUnlock()
	fake.markMutex.RLock()
	defer fake.markMutex.RUnlock()
	fake.rate1Mutex.RLock()
	defer fake.rate1Mutex.RUnlock()
	fake.rate15Mutex.RLock()
	defer fake.rate15Mutex.RUnlock()
	fake.rate5Mutex.RLock()
	defer fake.rate5Mutex.RUnlock()
	fake.rateMeanMutex.RLock()
	defer fake.rateMeanMutex.RUnlock()
	fake.snapshotMutex.RLock()
	defer fake.snapshotMutex.RUnlock()
	fake.stopMutex.RLock()
	defer fake.stopMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *MetricsMeter) recordInvocation(key string, args []interface{}) {
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
