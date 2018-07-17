



















package tally

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"time"
)

var (
	
	DefaultBuckets Buckets

	errBucketsCountNeedsGreaterThanZero = errors.New("n needs to be > 0")
	errBucketsStartNeedsGreaterThanZero = errors.New("start needs to be > 0")
	errBucketsFactorNeedsGreaterThanOne = errors.New("factor needs to be > 1")
)


type ValueBuckets []float64


func (v ValueBuckets) Len() int {
	return len(v)
}


func (v ValueBuckets) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}


func (v ValueBuckets) Less(i, j int) bool {
	return v[i] < v[j]
}

func (v ValueBuckets) String() string {
	values := make([]string, len(v))
	for i := range values {
		values[i] = fmt.Sprintf("%f", v[i])
	}
	return fmt.Sprint(values)
}


func (v ValueBuckets) AsValues() []float64 {
	return []float64(v)
}



func (v ValueBuckets) AsDurations() []time.Duration {
	values := make([]time.Duration, len(v))
	for i := range values {
		values[i] = time.Duration(v[i] * float64(time.Second))
	}
	return values
}


type DurationBuckets []time.Duration


func (v DurationBuckets) Len() int {
	return len(v)
}


func (v DurationBuckets) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}


func (v DurationBuckets) Less(i, j int) bool {
	return v[i] < v[j]
}

func (v DurationBuckets) String() string {
	values := make([]string, len(v))
	for i := range values {
		values[i] = v[i].String()
	}
	return fmt.Sprintf("%v", values)
}



func (v DurationBuckets) AsValues() []float64 {
	values := make([]float64, len(v))
	for i := range values {
		values[i] = float64(v[i]) / float64(time.Second)
	}
	return values
}


func (v DurationBuckets) AsDurations() []time.Duration {
	return []time.Duration(v)
}




func BucketPairs(buckets Buckets) []BucketPair {
	if buckets == nil || buckets.Len() < 1 {
		return []BucketPair{

			bucketPair{
				lowerBoundValue:    -math.MaxFloat64,
				upperBoundValue:    math.MaxFloat64,
				lowerBoundDuration: time.Duration(math.MinInt64),
				upperBoundDuration: time.Duration(math.MaxInt64),
			},
		}
	}

	var (
		asValueBuckets    = copyAndSortValues(buckets.AsValues())
		asDurationBuckets = copyAndSortDurations(buckets.AsDurations())
		pairs             = make([]BucketPair, 0, buckets.Len()+2)
	)

	pairs = append(pairs, bucketPair{
		lowerBoundValue:    -math.MaxFloat64,
		upperBoundValue:    asValueBuckets[0],
		lowerBoundDuration: time.Duration(math.MinInt64),
		upperBoundDuration: asDurationBuckets[0],
	})

	prevValueBucket, prevDurationBucket :=
		asValueBuckets[0], asDurationBuckets[0]

	for i := 1; i < buckets.Len(); i++ {
		pairs = append(pairs, bucketPair{
			lowerBoundValue:    prevValueBucket,
			upperBoundValue:    asValueBuckets[i],
			lowerBoundDuration: prevDurationBucket,
			upperBoundDuration: asDurationBuckets[i],
		})
		prevValueBucket, prevDurationBucket =
			asValueBuckets[i], asDurationBuckets[i]
	}

	pairs = append(pairs, bucketPair{
		lowerBoundValue:    prevValueBucket,
		upperBoundValue:    math.MaxFloat64,
		lowerBoundDuration: prevDurationBucket,
		upperBoundDuration: time.Duration(math.MaxInt64),
	})

	return pairs
}

func copyAndSortValues(values []float64) []float64 {
	valuesCopy := make([]float64, len(values))
	for i := range values {
		valuesCopy[i] = values[i]
	}
	sort.Sort(ValueBuckets(valuesCopy))
	return valuesCopy
}

func copyAndSortDurations(durations []time.Duration) []time.Duration {
	durationsCopy := make([]time.Duration, len(durations))
	for i := range durations {
		durationsCopy[i] = durations[i]
	}
	sort.Sort(DurationBuckets(durationsCopy))
	return durationsCopy
}

type bucketPair struct {
	lowerBoundValue    float64
	upperBoundValue    float64
	lowerBoundDuration time.Duration
	upperBoundDuration time.Duration
}

func (p bucketPair) LowerBoundValue() float64 {
	return p.lowerBoundValue
}

func (p bucketPair) UpperBoundValue() float64 {
	return p.upperBoundValue
}

func (p bucketPair) LowerBoundDuration() time.Duration {
	return p.lowerBoundDuration
}

func (p bucketPair) UpperBoundDuration() time.Duration {
	return p.upperBoundDuration
}


func LinearValueBuckets(start, width float64, n int) (ValueBuckets, error) {
	if n <= 0 {
		return nil, errBucketsCountNeedsGreaterThanZero
	}
	buckets := make([]float64, n)
	for i := range buckets {
		buckets[i] = start + (float64(i) * width)
	}
	return ValueBuckets(buckets), nil
}



func MustMakeLinearValueBuckets(start, width float64, n int) ValueBuckets {
	buckets, err := LinearValueBuckets(start, width, n)
	if err != nil {
		panic(err)
	}
	return buckets
}


func LinearDurationBuckets(start, width time.Duration, n int) (DurationBuckets, error) {
	if n <= 0 {
		return nil, errBucketsCountNeedsGreaterThanZero
	}
	buckets := make([]time.Duration, n)
	for i := range buckets {
		buckets[i] = start + (time.Duration(i) * width)
	}
	return DurationBuckets(buckets), nil
}



func MustMakeLinearDurationBuckets(start, width time.Duration, n int) DurationBuckets {
	buckets, err := LinearDurationBuckets(start, width, n)
	if err != nil {
		panic(err)
	}
	return buckets
}


func ExponentialValueBuckets(start, factor float64, n int) (ValueBuckets, error) {
	if n <= 0 {
		return nil, errBucketsCountNeedsGreaterThanZero
	}
	if start <= 0 {
		return nil, errBucketsStartNeedsGreaterThanZero
	}
	if factor <= 1 {
		return nil, errBucketsFactorNeedsGreaterThanOne
	}
	buckets := make([]float64, n)
	curr := start
	for i := range buckets {
		buckets[i] = curr
		curr *= factor
	}
	return ValueBuckets(buckets), nil
}



func MustMakeExponentialValueBuckets(start, factor float64, n int) ValueBuckets {
	buckets, err := ExponentialValueBuckets(start, factor, n)
	if err != nil {
		panic(err)
	}
	return buckets
}


func ExponentialDurationBuckets(start time.Duration, factor float64, n int) (DurationBuckets, error) {
	if n <= 0 {
		return nil, errBucketsCountNeedsGreaterThanZero
	}
	if start <= 0 {
		return nil, errBucketsStartNeedsGreaterThanZero
	}
	if factor <= 1 {
		return nil, errBucketsFactorNeedsGreaterThanOne
	}
	buckets := make([]time.Duration, n)
	curr := start
	for i := range buckets {
		buckets[i] = curr
		curr = time.Duration(float64(curr) * factor)
	}
	return DurationBuckets(buckets), nil
}



func MustMakeExponentialDurationBuckets(start time.Duration, factor float64, n int) DurationBuckets {
	buckets, err := ExponentialDurationBuckets(start, factor, n)
	if err != nil {
		panic(err)
	}
	return buckets
}
