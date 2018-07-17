












package prometheus

import (
	"fmt"
	"math"
	"sort"
	"sync/atomic"

	"github.com/golang/protobuf/proto"

	dto "github.com/prometheus/client_model/go"
)
















type Histogram interface {
	Metric
	Collector

	
	Observe(float64)
}



const bucketLabel = "le"





var (
	DefBuckets = []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10}

	errBucketLabelNotAllowed = fmt.Errorf(
		"%q is not allowed as label name in histograms", bucketLabel,
	)
)







func LinearBuckets(start, width float64, count int) []float64 {
	if count < 1 {
		panic("LinearBuckets needs a positive count")
	}
	buckets := make([]float64, count)
	for i := range buckets {
		buckets[i] = start
		start += width
	}
	return buckets
}









func ExponentialBuckets(start, factor float64, count int) []float64 {
	if count < 1 {
		panic("ExponentialBuckets needs a positive count")
	}
	if start <= 0 {
		panic("ExponentialBuckets needs a positive start value")
	}
	if factor <= 1 {
		panic("ExponentialBuckets needs a factor greater than 1")
	}
	buckets := make([]float64, count)
	for i := range buckets {
		buckets[i] = start
		start *= factor
	}
	return buckets
}




type HistogramOpts struct {
	
	
	
	
	
	Namespace string
	Subsystem string
	Name      string

	
	
	
	
	Help string

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	ConstLabels Labels

	
	
	
	
	
	Buckets []float64
}



func NewHistogram(opts HistogramOpts) Histogram {
	return newHistogram(
		NewDesc(
			BuildFQName(opts.Namespace, opts.Subsystem, opts.Name),
			opts.Help,
			nil,
			opts.ConstLabels,
		),
		opts,
	)
}

func newHistogram(desc *Desc, opts HistogramOpts, labelValues ...string) Histogram {
	if len(desc.variableLabels) != len(labelValues) {
		panic(errInconsistentCardinality)
	}

	for _, n := range desc.variableLabels {
		if n == bucketLabel {
			panic(errBucketLabelNotAllowed)
		}
	}
	for _, lp := range desc.constLabelPairs {
		if lp.GetName() == bucketLabel {
			panic(errBucketLabelNotAllowed)
		}
	}

	if len(opts.Buckets) == 0 {
		opts.Buckets = DefBuckets
	}

	h := &histogram{
		desc:        desc,
		upperBounds: opts.Buckets,
		labelPairs:  makeLabelPairs(desc, labelValues),
	}
	for i, upperBound := range h.upperBounds {
		if i < len(h.upperBounds)-1 {
			if upperBound >= h.upperBounds[i+1] {
				panic(fmt.Errorf(
					"histogram buckets must be in increasing order: %f >= %f",
					upperBound, h.upperBounds[i+1],
				))
			}
		} else {
			if math.IsInf(upperBound, +1) {
				
				h.upperBounds = h.upperBounds[:i]
			}
		}
	}
	
	h.counts = make([]uint64, len(h.upperBounds))

	h.init(h) 
	return h
}

type histogram struct {
	
	
	
	
	sumBits uint64
	count   uint64

	selfCollector
	

	desc *Desc

	upperBounds []float64
	counts      []uint64

	labelPairs []*dto.LabelPair
}

func (h *histogram) Desc() *Desc {
	return h.desc
}

func (h *histogram) Observe(v float64) {
	
	
	
	
	
	
	
	
	
	i := sort.SearchFloat64s(h.upperBounds, v)
	if i < len(h.counts) {
		atomic.AddUint64(&h.counts[i], 1)
	}
	atomic.AddUint64(&h.count, 1)
	for {
		oldBits := atomic.LoadUint64(&h.sumBits)
		newBits := math.Float64bits(math.Float64frombits(oldBits) + v)
		if atomic.CompareAndSwapUint64(&h.sumBits, oldBits, newBits) {
			break
		}
	}
}

func (h *histogram) Write(out *dto.Metric) error {
	his := &dto.Histogram{}
	buckets := make([]*dto.Bucket, len(h.upperBounds))

	his.SampleSum = proto.Float64(math.Float64frombits(atomic.LoadUint64(&h.sumBits)))
	his.SampleCount = proto.Uint64(atomic.LoadUint64(&h.count))
	var count uint64
	for i, upperBound := range h.upperBounds {
		count += atomic.LoadUint64(&h.counts[i])
		buckets[i] = &dto.Bucket{
			CumulativeCount: proto.Uint64(count),
			UpperBound:      proto.Float64(upperBound),
		}
	}
	his.Bucket = buckets
	out.Histogram = his
	out.Label = h.labelPairs
	return nil
}






type HistogramVec struct {
	*MetricVec
}




func NewHistogramVec(opts HistogramOpts, labelNames []string) *HistogramVec {
	desc := NewDesc(
		BuildFQName(opts.Namespace, opts.Subsystem, opts.Name),
		opts.Help,
		labelNames,
		opts.ConstLabels,
	)
	return &HistogramVec{
		MetricVec: newMetricVec(desc, func(lvs ...string) Metric {
			return newHistogram(desc, opts, lvs...)
		}),
	}
}




func (m *HistogramVec) GetMetricWithLabelValues(lvs ...string) (Histogram, error) {
	metric, err := m.MetricVec.GetMetricWithLabelValues(lvs...)
	if metric != nil {
		return metric.(Histogram), err
	}
	return nil, err
}




func (m *HistogramVec) GetMetricWith(labels Labels) (Histogram, error) {
	metric, err := m.MetricVec.GetMetricWith(labels)
	if metric != nil {
		return metric.(Histogram), err
	}
	return nil, err
}





func (m *HistogramVec) WithLabelValues(lvs ...string) Histogram {
	return m.MetricVec.WithLabelValues(lvs...).(Histogram)
}




func (m *HistogramVec) With(labels Labels) Histogram {
	return m.MetricVec.With(labels).(Histogram)
}

type constHistogram struct {
	desc       *Desc
	count      uint64
	sum        float64
	buckets    map[float64]uint64
	labelPairs []*dto.LabelPair
}

func (h *constHistogram) Desc() *Desc {
	return h.desc
}

func (h *constHistogram) Write(out *dto.Metric) error {
	his := &dto.Histogram{}
	buckets := make([]*dto.Bucket, 0, len(h.buckets))

	his.SampleCount = proto.Uint64(h.count)
	his.SampleSum = proto.Float64(h.sum)

	for upperBound, count := range h.buckets {
		buckets = append(buckets, &dto.Bucket{
			CumulativeCount: proto.Uint64(count),
			UpperBound:      proto.Float64(upperBound),
		})
	}

	if len(buckets) > 0 {
		sort.Sort(buckSort(buckets))
	}
	his.Bucket = buckets

	out.Histogram = his
	out.Label = h.labelPairs

	return nil
}














func NewConstHistogram(
	desc *Desc,
	count uint64,
	sum float64,
	buckets map[float64]uint64,
	labelValues ...string,
) (Metric, error) {
	if len(desc.variableLabels) != len(labelValues) {
		return nil, errInconsistentCardinality
	}
	return &constHistogram{
		desc:       desc,
		count:      count,
		sum:        sum,
		buckets:    buckets,
		labelPairs: makeLabelPairs(desc, labelValues),
	}, nil
}



func MustNewConstHistogram(
	desc *Desc,
	count uint64,
	sum float64,
	buckets map[float64]uint64,
	labelValues ...string,
) Metric {
	m, err := NewConstHistogram(desc, count, sum, buckets, labelValues...)
	if err != nil {
		panic(err)
	}
	return m
}

type buckSort []*dto.Bucket

func (s buckSort) Len() int {
	return len(s)
}

func (s buckSort) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s buckSort) Less(i, j int) bool {
	return s[i].GetUpperBound() < s[j].GetUpperBound()
}
