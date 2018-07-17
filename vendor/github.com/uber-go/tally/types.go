



















package tally

import (
	"fmt"
	"sort"
	"time"
)



type Scope interface {
	
	Counter(name string) Counter

	
	Gauge(name string) Gauge

	
	Timer(name string) Timer

	
	
	
	
	
	
	
	
	
	Histogram(name string, buckets Buckets) Histogram

	
	Tagged(tags map[string]string) Scope

	
	SubScope(name string) Scope

	
	Capabilities() Capabilities
}


type Counter interface {
	
	Inc(delta int64)
}


type Gauge interface {
	
	Update(value float64)
}


type Timer interface {
	
	Record(value time.Duration)

	
	Start() Stopwatch
}


type Histogram interface {
	
	
	RecordValue(value float64)

	
	
	RecordDuration(value time.Duration)

	
	
	Start() Stopwatch
}




type Stopwatch struct {
	start    time.Time
	recorder StopwatchRecorder
}



func NewStopwatch(start time.Time, r StopwatchRecorder) Stopwatch {
	return Stopwatch{start: start, recorder: r}
}


func (sw Stopwatch) Stop() {
	sw.recorder.RecordStopwatch(sw.start)
}



type StopwatchRecorder interface {
	RecordStopwatch(stopwatchStart time.Time)
}



type Buckets interface {
	fmt.Stringer
	sort.Interface

	
	AsValues() []float64

	
	AsDurations() []time.Duration
}



type BucketPair interface {
	LowerBoundValue() float64
	UpperBoundValue() float64
	LowerBoundDuration() time.Duration
	UpperBoundDuration() time.Duration
}


type Capabilities interface {
	
	Reporting() bool

	
	Tagging() bool
}
