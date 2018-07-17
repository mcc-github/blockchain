



















package statsd

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/uber-go/tally"

	"github.com/cactus/go-statsd-client/statsd"
)

const (
	
	
	
	DefaultHistogramBucketNamePrecision = uint(6)
)

type cactusStatsReporter struct {
	statter    statsd.Statter
	sampleRate float32
	bucketFmt  string
}


type Options struct {
	
	
	SampleRate float32

	
	
	
	HistogramBucketNamePrecision uint
}



func NewReporter(statsd statsd.Statter, opts Options) tally.StatsReporter {
	var nilSampleRate float32
	if opts.SampleRate == nilSampleRate {
		opts.SampleRate = 1.0
	}
	if opts.HistogramBucketNamePrecision == 0 {
		opts.HistogramBucketNamePrecision = DefaultHistogramBucketNamePrecision
	}
	return &cactusStatsReporter{
		statter:    statsd,
		sampleRate: opts.SampleRate,
		bucketFmt:  "%." + strconv.Itoa(int(opts.HistogramBucketNamePrecision)) + "f",
	}
}

func (r *cactusStatsReporter) ReportCounter(name string, tags map[string]string, value int64) {
	r.statter.Inc(name, value, r.sampleRate)
}

func (r *cactusStatsReporter) ReportGauge(name string, tags map[string]string, value float64) {
	r.statter.Gauge(name, int64(value), r.sampleRate)
}

func (r *cactusStatsReporter) ReportTimer(name string, tags map[string]string, interval time.Duration) {
	r.statter.TimingDuration(name, interval, r.sampleRate)
}

func (r *cactusStatsReporter) ReportHistogramValueSamples(
	name string,
	tags map[string]string,
	buckets tally.Buckets,
	bucketLowerBound,
	bucketUpperBound float64,
	samples int64,
) {
	r.statter.Inc(
		fmt.Sprintf("%s.%s-%s", name,
			r.valueBucketString(bucketLowerBound),
			r.valueBucketString(bucketUpperBound)),
		samples, r.sampleRate)
}

func (r *cactusStatsReporter) ReportHistogramDurationSamples(
	name string,
	tags map[string]string,
	buckets tally.Buckets,
	bucketLowerBound,
	bucketUpperBound time.Duration,
	samples int64,
) {
	r.statter.Inc(
		fmt.Sprintf("%s.%s-%s", name,
			r.durationBucketString(bucketLowerBound),
			r.durationBucketString(bucketUpperBound)),
		samples, r.sampleRate)
}

func (r *cactusStatsReporter) valueBucketString(
	upperBound float64,
) string {
	if upperBound == math.MaxFloat64 {
		return "infinity"
	}
	if upperBound == -math.MaxFloat64 {
		return "-infinity"
	}
	return fmt.Sprintf(r.bucketFmt, upperBound)
}

func (r *cactusStatsReporter) durationBucketString(
	upperBound time.Duration,
) string {
	if upperBound == time.Duration(math.MaxInt64) {
		return "infinity"
	}
	if upperBound == time.Duration(math.MinInt64) {
		return "-infinity"
	}
	return upperBound.String()
}

func (r *cactusStatsReporter) Capabilities() tally.Capabilities {
	return r
}

func (r *cactusStatsReporter) Reporting() bool {
	return true
}

func (r *cactusStatsReporter) Tagging() bool {
	return false
}

func (r *cactusStatsReporter) Flush() {
	
}
