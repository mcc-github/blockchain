



















package tally

import "time"


type BaseStatsReporter interface {
	
	Capabilities() Capabilities

	
	Flush()
}


type StatsReporter interface {
	BaseStatsReporter

	
	ReportCounter(
		name string,
		tags map[string]string,
		value int64,
	)

	
	ReportGauge(
		name string,
		tags map[string]string,
		value float64,
	)

	
	ReportTimer(
		name string,
		tags map[string]string,
		interval time.Duration,
	)

	
	ReportHistogramValueSamples(
		name string,
		tags map[string]string,
		buckets Buckets,
		bucketLowerBound,
		bucketUpperBound float64,
		samples int64,
	)

	
	ReportHistogramDurationSamples(
		name string,
		tags map[string]string,
		buckets Buckets,
		bucketLowerBound,
		bucketUpperBound time.Duration,
		samples int64,
	)
}



type CachedStatsReporter interface {
	BaseStatsReporter

	
	AllocateCounter(
		name string,
		tags map[string]string,
	) CachedCount

	
	AllocateGauge(
		name string,
		tags map[string]string,
	) CachedGauge

	
	AllocateTimer(
		name string,
		tags map[string]string,
	) CachedTimer

	
	
	AllocateHistogram(
		name string,
		tags map[string]string,
		buckets Buckets,
	) CachedHistogram
}


type CachedCount interface {
	ReportCount(value int64)
}


type CachedGauge interface {
	ReportGauge(value float64)
}


type CachedTimer interface {
	ReportTimer(interval time.Duration)
}


type CachedHistogram interface {
	ValueBucket(
		bucketLowerBound, bucketUpperBound float64,
	) CachedHistogramBucket
	DurationBucket(
		bucketLowerBound, bucketUpperBound time.Duration,
	) CachedHistogramBucket
}


type CachedHistogramBucket interface {
	ReportSamples(value int64)
}
