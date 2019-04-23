
package mock

import (
	sync "sync"

	metrics "github.com/rcrowley/go-metrics"
)

type MetricsHistogram struct {
	ClearStub        func()
	clearMutex       sync.RWMutex
	clearArgsForCall []struct {
	}
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
	MaxStub        func() int64
	maxMutex       sync.RWMutex
	maxArgsForCall []struct {
	}
	maxReturns struct {
		result1 int64
	}
	maxReturnsOnCall map[int]struct {
		result1 int64
	}
	MeanStub        func() float64
	meanMutex       sync.RWMutex
	meanArgsForCall []struct {
	}
	meanReturns struct {
		result1 float64
	}
	meanReturnsOnCall map[int]struct {
		result1 float64
	}
	MinStub        func() int64
	minMutex       sync.RWMutex
	minArgsForCall []struct {
	}
	minReturns struct {
		result1 int64
	}
	minReturnsOnCall map[int]struct {
		result1 int64
	}
	PercentileStub        func(float64) float64
	percentileMutex       sync.RWMutex
	percentileArgsForCall []struct {
		arg1 float64
	}
	percentileReturns struct {
		result1 float64
	}
	percentileReturnsOnCall map[int]struct {
		result1 float64
	}
	PercentilesStub        func([]float64) []float64
	percentilesMutex       sync.RWMutex
	percentilesArgsForCall []struct {
		arg1 []float64
	}
	percentilesReturns struct {
		result1 []float64
	}
	percentilesReturnsOnCall map[int]struct {
		result1 []float64
	}
	SampleStub        func() metrics.Sample
	sampleMutex       sync.RWMutex
	sampleArgsForCall []struct {
	}
	sampleReturns struct {
		result1 metrics.Sample
	}
	sampleReturnsOnCall map[int]struct {
		result1 metrics.Sample
	}
	SnapshotStub        func() metrics.Histogram
	snapshotMutex       sync.RWMutex
	snapshotArgsForCall []struct {
	}
	snapshotReturns struct {
		result1 metrics.Histogram
	}
	snapshotReturnsOnCall map[int]struct {
		result1 metrics.Histogram
	}
	StdDevStub        func() float64
	stdDevMutex       sync.RWMutex
	stdDevArgsForCall []struct {
	}
	stdDevReturns struct {
		result1 float64
	}
	stdDevReturnsOnCall map[int]struct {
		result1 float64
	}
	SumStub        func() int64
	sumMutex       sync.RWMutex
	sumArgsForCall []struct {
	}
	sumReturns struct {
		result1 int64
	}
	sumReturnsOnCall map[int]struct {
		result1 int64
	}
	UpdateStub        func(int64)
	updateMutex       sync.RWMutex
	updateArgsForCall []struct {
		arg1 int64
	}
	VarianceStub        func() float64
	varianceMutex       sync.RWMutex
	varianceArgsForCall []struct {
	}
	varianceReturns struct {
		result1 float64
	}
	varianceReturnsOnCall map[int]struct {
		result1 float64
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *MetricsHistogram) Clear() {
	fake.clearMutex.Lock()
	fake.clearArgsForCall = append(fake.clearArgsForCall, struct {
	}{})
	fake.recordInvocation("Clear", []interface{}{})
	fake.clearMutex.Unlock()
	if fake.ClearStub != nil {
		fake.ClearStub()
	}
}

func (fake *MetricsHistogram) ClearCallCount() int {
	fake.clearMutex.RLock()
	defer fake.clearMutex.RUnlock()
	return len(fake.clearArgsForCall)
}

func (fake *MetricsHistogram) ClearCalls(stub func()) {
	fake.clearMutex.Lock()
	defer fake.clearMutex.Unlock()
	fake.ClearStub = stub
}

func (fake *MetricsHistogram) Count() int64 {
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

func (fake *MetricsHistogram) CountCallCount() int {
	fake.countMutex.RLock()
	defer fake.countMutex.RUnlock()
	return len(fake.countArgsForCall)
}

func (fake *MetricsHistogram) CountCalls(stub func() int64) {
	fake.countMutex.Lock()
	defer fake.countMutex.Unlock()
	fake.CountStub = stub
}

func (fake *MetricsHistogram) CountReturns(result1 int64) {
	fake.countMutex.Lock()
	defer fake.countMutex.Unlock()
	fake.CountStub = nil
	fake.countReturns = struct {
		result1 int64
	}{result1}
}

func (fake *MetricsHistogram) CountReturnsOnCall(i int, result1 int64) {
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

func (fake *MetricsHistogram) Max() int64 {
	fake.maxMutex.Lock()
	ret, specificReturn := fake.maxReturnsOnCall[len(fake.maxArgsForCall)]
	fake.maxArgsForCall = append(fake.maxArgsForCall, struct {
	}{})
	fake.recordInvocation("Max", []interface{}{})
	fake.maxMutex.Unlock()
	if fake.MaxStub != nil {
		return fake.MaxStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.maxReturns
	return fakeReturns.result1
}

func (fake *MetricsHistogram) MaxCallCount() int {
	fake.maxMutex.RLock()
	defer fake.maxMutex.RUnlock()
	return len(fake.maxArgsForCall)
}

func (fake *MetricsHistogram) MaxCalls(stub func() int64) {
	fake.maxMutex.Lock()
	defer fake.maxMutex.Unlock()
	fake.MaxStub = stub
}

func (fake *MetricsHistogram) MaxReturns(result1 int64) {
	fake.maxMutex.Lock()
	defer fake.maxMutex.Unlock()
	fake.MaxStub = nil
	fake.maxReturns = struct {
		result1 int64
	}{result1}
}

func (fake *MetricsHistogram) MaxReturnsOnCall(i int, result1 int64) {
	fake.maxMutex.Lock()
	defer fake.maxMutex.Unlock()
	fake.MaxStub = nil
	if fake.maxReturnsOnCall == nil {
		fake.maxReturnsOnCall = make(map[int]struct {
			result1 int64
		})
	}
	fake.maxReturnsOnCall[i] = struct {
		result1 int64
	}{result1}
}

func (fake *MetricsHistogram) Mean() float64 {
	fake.meanMutex.Lock()
	ret, specificReturn := fake.meanReturnsOnCall[len(fake.meanArgsForCall)]
	fake.meanArgsForCall = append(fake.meanArgsForCall, struct {
	}{})
	fake.recordInvocation("Mean", []interface{}{})
	fake.meanMutex.Unlock()
	if fake.MeanStub != nil {
		return fake.MeanStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.meanReturns
	return fakeReturns.result1
}

func (fake *MetricsHistogram) MeanCallCount() int {
	fake.meanMutex.RLock()
	defer fake.meanMutex.RUnlock()
	return len(fake.meanArgsForCall)
}

func (fake *MetricsHistogram) MeanCalls(stub func() float64) {
	fake.meanMutex.Lock()
	defer fake.meanMutex.Unlock()
	fake.MeanStub = stub
}

func (fake *MetricsHistogram) MeanReturns(result1 float64) {
	fake.meanMutex.Lock()
	defer fake.meanMutex.Unlock()
	fake.MeanStub = nil
	fake.meanReturns = struct {
		result1 float64
	}{result1}
}

func (fake *MetricsHistogram) MeanReturnsOnCall(i int, result1 float64) {
	fake.meanMutex.Lock()
	defer fake.meanMutex.Unlock()
	fake.MeanStub = nil
	if fake.meanReturnsOnCall == nil {
		fake.meanReturnsOnCall = make(map[int]struct {
			result1 float64
		})
	}
	fake.meanReturnsOnCall[i] = struct {
		result1 float64
	}{result1}
}

func (fake *MetricsHistogram) Min() int64 {
	fake.minMutex.Lock()
	ret, specificReturn := fake.minReturnsOnCall[len(fake.minArgsForCall)]
	fake.minArgsForCall = append(fake.minArgsForCall, struct {
	}{})
	fake.recordInvocation("Min", []interface{}{})
	fake.minMutex.Unlock()
	if fake.MinStub != nil {
		return fake.MinStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.minReturns
	return fakeReturns.result1
}

func (fake *MetricsHistogram) MinCallCount() int {
	fake.minMutex.RLock()
	defer fake.minMutex.RUnlock()
	return len(fake.minArgsForCall)
}

func (fake *MetricsHistogram) MinCalls(stub func() int64) {
	fake.minMutex.Lock()
	defer fake.minMutex.Unlock()
	fake.MinStub = stub
}

func (fake *MetricsHistogram) MinReturns(result1 int64) {
	fake.minMutex.Lock()
	defer fake.minMutex.Unlock()
	fake.MinStub = nil
	fake.minReturns = struct {
		result1 int64
	}{result1}
}

func (fake *MetricsHistogram) MinReturnsOnCall(i int, result1 int64) {
	fake.minMutex.Lock()
	defer fake.minMutex.Unlock()
	fake.MinStub = nil
	if fake.minReturnsOnCall == nil {
		fake.minReturnsOnCall = make(map[int]struct {
			result1 int64
		})
	}
	fake.minReturnsOnCall[i] = struct {
		result1 int64
	}{result1}
}

func (fake *MetricsHistogram) Percentile(arg1 float64) float64 {
	fake.percentileMutex.Lock()
	ret, specificReturn := fake.percentileReturnsOnCall[len(fake.percentileArgsForCall)]
	fake.percentileArgsForCall = append(fake.percentileArgsForCall, struct {
		arg1 float64
	}{arg1})
	fake.recordInvocation("Percentile", []interface{}{arg1})
	fake.percentileMutex.Unlock()
	if fake.PercentileStub != nil {
		return fake.PercentileStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.percentileReturns
	return fakeReturns.result1
}

func (fake *MetricsHistogram) PercentileCallCount() int {
	fake.percentileMutex.RLock()
	defer fake.percentileMutex.RUnlock()
	return len(fake.percentileArgsForCall)
}

func (fake *MetricsHistogram) PercentileCalls(stub func(float64) float64) {
	fake.percentileMutex.Lock()
	defer fake.percentileMutex.Unlock()
	fake.PercentileStub = stub
}

func (fake *MetricsHistogram) PercentileArgsForCall(i int) float64 {
	fake.percentileMutex.RLock()
	defer fake.percentileMutex.RUnlock()
	argsForCall := fake.percentileArgsForCall[i]
	return argsForCall.arg1
}

func (fake *MetricsHistogram) PercentileReturns(result1 float64) {
	fake.percentileMutex.Lock()
	defer fake.percentileMutex.Unlock()
	fake.PercentileStub = nil
	fake.percentileReturns = struct {
		result1 float64
	}{result1}
}

func (fake *MetricsHistogram) PercentileReturnsOnCall(i int, result1 float64) {
	fake.percentileMutex.Lock()
	defer fake.percentileMutex.Unlock()
	fake.PercentileStub = nil
	if fake.percentileReturnsOnCall == nil {
		fake.percentileReturnsOnCall = make(map[int]struct {
			result1 float64
		})
	}
	fake.percentileReturnsOnCall[i] = struct {
		result1 float64
	}{result1}
}

func (fake *MetricsHistogram) Percentiles(arg1 []float64) []float64 {
	var arg1Copy []float64
	if arg1 != nil {
		arg1Copy = make([]float64, len(arg1))
		copy(arg1Copy, arg1)
	}
	fake.percentilesMutex.Lock()
	ret, specificReturn := fake.percentilesReturnsOnCall[len(fake.percentilesArgsForCall)]
	fake.percentilesArgsForCall = append(fake.percentilesArgsForCall, struct {
		arg1 []float64
	}{arg1Copy})
	fake.recordInvocation("Percentiles", []interface{}{arg1Copy})
	fake.percentilesMutex.Unlock()
	if fake.PercentilesStub != nil {
		return fake.PercentilesStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.percentilesReturns
	return fakeReturns.result1
}

func (fake *MetricsHistogram) PercentilesCallCount() int {
	fake.percentilesMutex.RLock()
	defer fake.percentilesMutex.RUnlock()
	return len(fake.percentilesArgsForCall)
}

func (fake *MetricsHistogram) PercentilesCalls(stub func([]float64) []float64) {
	fake.percentilesMutex.Lock()
	defer fake.percentilesMutex.Unlock()
	fake.PercentilesStub = stub
}

func (fake *MetricsHistogram) PercentilesArgsForCall(i int) []float64 {
	fake.percentilesMutex.RLock()
	defer fake.percentilesMutex.RUnlock()
	argsForCall := fake.percentilesArgsForCall[i]
	return argsForCall.arg1
}

func (fake *MetricsHistogram) PercentilesReturns(result1 []float64) {
	fake.percentilesMutex.Lock()
	defer fake.percentilesMutex.Unlock()
	fake.PercentilesStub = nil
	fake.percentilesReturns = struct {
		result1 []float64
	}{result1}
}

func (fake *MetricsHistogram) PercentilesReturnsOnCall(i int, result1 []float64) {
	fake.percentilesMutex.Lock()
	defer fake.percentilesMutex.Unlock()
	fake.PercentilesStub = nil
	if fake.percentilesReturnsOnCall == nil {
		fake.percentilesReturnsOnCall = make(map[int]struct {
			result1 []float64
		})
	}
	fake.percentilesReturnsOnCall[i] = struct {
		result1 []float64
	}{result1}
}

func (fake *MetricsHistogram) Sample() metrics.Sample {
	fake.sampleMutex.Lock()
	ret, specificReturn := fake.sampleReturnsOnCall[len(fake.sampleArgsForCall)]
	fake.sampleArgsForCall = append(fake.sampleArgsForCall, struct {
	}{})
	fake.recordInvocation("Sample", []interface{}{})
	fake.sampleMutex.Unlock()
	if fake.SampleStub != nil {
		return fake.SampleStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.sampleReturns
	return fakeReturns.result1
}

func (fake *MetricsHistogram) SampleCallCount() int {
	fake.sampleMutex.RLock()
	defer fake.sampleMutex.RUnlock()
	return len(fake.sampleArgsForCall)
}

func (fake *MetricsHistogram) SampleCalls(stub func() metrics.Sample) {
	fake.sampleMutex.Lock()
	defer fake.sampleMutex.Unlock()
	fake.SampleStub = stub
}

func (fake *MetricsHistogram) SampleReturns(result1 metrics.Sample) {
	fake.sampleMutex.Lock()
	defer fake.sampleMutex.Unlock()
	fake.SampleStub = nil
	fake.sampleReturns = struct {
		result1 metrics.Sample
	}{result1}
}

func (fake *MetricsHistogram) SampleReturnsOnCall(i int, result1 metrics.Sample) {
	fake.sampleMutex.Lock()
	defer fake.sampleMutex.Unlock()
	fake.SampleStub = nil
	if fake.sampleReturnsOnCall == nil {
		fake.sampleReturnsOnCall = make(map[int]struct {
			result1 metrics.Sample
		})
	}
	fake.sampleReturnsOnCall[i] = struct {
		result1 metrics.Sample
	}{result1}
}

func (fake *MetricsHistogram) Snapshot() metrics.Histogram {
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

func (fake *MetricsHistogram) SnapshotCallCount() int {
	fake.snapshotMutex.RLock()
	defer fake.snapshotMutex.RUnlock()
	return len(fake.snapshotArgsForCall)
}

func (fake *MetricsHistogram) SnapshotCalls(stub func() metrics.Histogram) {
	fake.snapshotMutex.Lock()
	defer fake.snapshotMutex.Unlock()
	fake.SnapshotStub = stub
}

func (fake *MetricsHistogram) SnapshotReturns(result1 metrics.Histogram) {
	fake.snapshotMutex.Lock()
	defer fake.snapshotMutex.Unlock()
	fake.SnapshotStub = nil
	fake.snapshotReturns = struct {
		result1 metrics.Histogram
	}{result1}
}

func (fake *MetricsHistogram) SnapshotReturnsOnCall(i int, result1 metrics.Histogram) {
	fake.snapshotMutex.Lock()
	defer fake.snapshotMutex.Unlock()
	fake.SnapshotStub = nil
	if fake.snapshotReturnsOnCall == nil {
		fake.snapshotReturnsOnCall = make(map[int]struct {
			result1 metrics.Histogram
		})
	}
	fake.snapshotReturnsOnCall[i] = struct {
		result1 metrics.Histogram
	}{result1}
}

func (fake *MetricsHistogram) StdDev() float64 {
	fake.stdDevMutex.Lock()
	ret, specificReturn := fake.stdDevReturnsOnCall[len(fake.stdDevArgsForCall)]
	fake.stdDevArgsForCall = append(fake.stdDevArgsForCall, struct {
	}{})
	fake.recordInvocation("StdDev", []interface{}{})
	fake.stdDevMutex.Unlock()
	if fake.StdDevStub != nil {
		return fake.StdDevStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.stdDevReturns
	return fakeReturns.result1
}

func (fake *MetricsHistogram) StdDevCallCount() int {
	fake.stdDevMutex.RLock()
	defer fake.stdDevMutex.RUnlock()
	return len(fake.stdDevArgsForCall)
}

func (fake *MetricsHistogram) StdDevCalls(stub func() float64) {
	fake.stdDevMutex.Lock()
	defer fake.stdDevMutex.Unlock()
	fake.StdDevStub = stub
}

func (fake *MetricsHistogram) StdDevReturns(result1 float64) {
	fake.stdDevMutex.Lock()
	defer fake.stdDevMutex.Unlock()
	fake.StdDevStub = nil
	fake.stdDevReturns = struct {
		result1 float64
	}{result1}
}

func (fake *MetricsHistogram) StdDevReturnsOnCall(i int, result1 float64) {
	fake.stdDevMutex.Lock()
	defer fake.stdDevMutex.Unlock()
	fake.StdDevStub = nil
	if fake.stdDevReturnsOnCall == nil {
		fake.stdDevReturnsOnCall = make(map[int]struct {
			result1 float64
		})
	}
	fake.stdDevReturnsOnCall[i] = struct {
		result1 float64
	}{result1}
}

func (fake *MetricsHistogram) Sum() int64 {
	fake.sumMutex.Lock()
	ret, specificReturn := fake.sumReturnsOnCall[len(fake.sumArgsForCall)]
	fake.sumArgsForCall = append(fake.sumArgsForCall, struct {
	}{})
	fake.recordInvocation("Sum", []interface{}{})
	fake.sumMutex.Unlock()
	if fake.SumStub != nil {
		return fake.SumStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.sumReturns
	return fakeReturns.result1
}

func (fake *MetricsHistogram) SumCallCount() int {
	fake.sumMutex.RLock()
	defer fake.sumMutex.RUnlock()
	return len(fake.sumArgsForCall)
}

func (fake *MetricsHistogram) SumCalls(stub func() int64) {
	fake.sumMutex.Lock()
	defer fake.sumMutex.Unlock()
	fake.SumStub = stub
}

func (fake *MetricsHistogram) SumReturns(result1 int64) {
	fake.sumMutex.Lock()
	defer fake.sumMutex.Unlock()
	fake.SumStub = nil
	fake.sumReturns = struct {
		result1 int64
	}{result1}
}

func (fake *MetricsHistogram) SumReturnsOnCall(i int, result1 int64) {
	fake.sumMutex.Lock()
	defer fake.sumMutex.Unlock()
	fake.SumStub = nil
	if fake.sumReturnsOnCall == nil {
		fake.sumReturnsOnCall = make(map[int]struct {
			result1 int64
		})
	}
	fake.sumReturnsOnCall[i] = struct {
		result1 int64
	}{result1}
}

func (fake *MetricsHistogram) Update(arg1 int64) {
	fake.updateMutex.Lock()
	fake.updateArgsForCall = append(fake.updateArgsForCall, struct {
		arg1 int64
	}{arg1})
	fake.recordInvocation("Update", []interface{}{arg1})
	fake.updateMutex.Unlock()
	if fake.UpdateStub != nil {
		fake.UpdateStub(arg1)
	}
}

func (fake *MetricsHistogram) UpdateCallCount() int {
	fake.updateMutex.RLock()
	defer fake.updateMutex.RUnlock()
	return len(fake.updateArgsForCall)
}

func (fake *MetricsHistogram) UpdateCalls(stub func(int64)) {
	fake.updateMutex.Lock()
	defer fake.updateMutex.Unlock()
	fake.UpdateStub = stub
}

func (fake *MetricsHistogram) UpdateArgsForCall(i int) int64 {
	fake.updateMutex.RLock()
	defer fake.updateMutex.RUnlock()
	argsForCall := fake.updateArgsForCall[i]
	return argsForCall.arg1
}

func (fake *MetricsHistogram) Variance() float64 {
	fake.varianceMutex.Lock()
	ret, specificReturn := fake.varianceReturnsOnCall[len(fake.varianceArgsForCall)]
	fake.varianceArgsForCall = append(fake.varianceArgsForCall, struct {
	}{})
	fake.recordInvocation("Variance", []interface{}{})
	fake.varianceMutex.Unlock()
	if fake.VarianceStub != nil {
		return fake.VarianceStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.varianceReturns
	return fakeReturns.result1
}

func (fake *MetricsHistogram) VarianceCallCount() int {
	fake.varianceMutex.RLock()
	defer fake.varianceMutex.RUnlock()
	return len(fake.varianceArgsForCall)
}

func (fake *MetricsHistogram) VarianceCalls(stub func() float64) {
	fake.varianceMutex.Lock()
	defer fake.varianceMutex.Unlock()
	fake.VarianceStub = stub
}

func (fake *MetricsHistogram) VarianceReturns(result1 float64) {
	fake.varianceMutex.Lock()
	defer fake.varianceMutex.Unlock()
	fake.VarianceStub = nil
	fake.varianceReturns = struct {
		result1 float64
	}{result1}
}

func (fake *MetricsHistogram) VarianceReturnsOnCall(i int, result1 float64) {
	fake.varianceMutex.Lock()
	defer fake.varianceMutex.Unlock()
	fake.VarianceStub = nil
	if fake.varianceReturnsOnCall == nil {
		fake.varianceReturnsOnCall = make(map[int]struct {
			result1 float64
		})
	}
	fake.varianceReturnsOnCall[i] = struct {
		result1 float64
	}{result1}
}

func (fake *MetricsHistogram) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.clearMutex.RLock()
	defer fake.clearMutex.RUnlock()
	fake.countMutex.RLock()
	defer fake.countMutex.RUnlock()
	fake.maxMutex.RLock()
	defer fake.maxMutex.RUnlock()
	fake.meanMutex.RLock()
	defer fake.meanMutex.RUnlock()
	fake.minMutex.RLock()
	defer fake.minMutex.RUnlock()
	fake.percentileMutex.RLock()
	defer fake.percentileMutex.RUnlock()
	fake.percentilesMutex.RLock()
	defer fake.percentilesMutex.RUnlock()
	fake.sampleMutex.RLock()
	defer fake.sampleMutex.RUnlock()
	fake.snapshotMutex.RLock()
	defer fake.snapshotMutex.RUnlock()
	fake.stdDevMutex.RLock()
	defer fake.stdDevMutex.RUnlock()
	fake.sumMutex.RLock()
	defer fake.sumMutex.RUnlock()
	fake.updateMutex.RLock()
	defer fake.updateMutex.RUnlock()
	fake.varianceMutex.RLock()
	defer fake.varianceMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *MetricsHistogram) recordInvocation(key string, args []interface{}) {
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
