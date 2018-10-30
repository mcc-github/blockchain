












package prometheus

import (
	"fmt"
	"math"
	"runtime"
	"sort"
	"sync"
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
		counts:      [2]*histogramCounts{&histogramCounts{}, &histogramCounts{}},
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
	
	
	h.counts[0].buckets = make([]uint64, len(h.upperBounds))
	h.counts[1].buckets = make([]uint64, len(h.upperBounds))

	h.init(h) 
	return h
}

type histogramCounts struct {
	
	
	
	
	sumBits uint64
	count   uint64
	buckets []uint64
}

type histogram struct {
	
	
	
	
	
	
	
	
	
	
	
	
	
	countAndHotIdx uint64

	selfCollector
	desc     *Desc
	writeMtx sync.Mutex 

	upperBounds []float64

	
	
	
	
	counts [2]*histogramCounts
	hotIdx int 

	labelPairs []*dto.LabelPair
}

func (h *histogram) Desc() *Desc {
	return h.desc
}

func (h *histogram) Observe(v float64) {
	
	
	
	
	
	
	
	
	
	i := sort.SearchFloat64s(h.upperBounds, v)

	
	
	
	n := atomic.AddUint64(&h.countAndHotIdx, 2)
	hotCounts := h.counts[n%2]

	if i < len(h.upperBounds) {
		atomic.AddUint64(&hotCounts.buckets[i], 1)
	}
	for {
		oldBits := atomic.LoadUint64(&hotCounts.sumBits)
		newBits := math.Float64bits(math.Float64frombits(oldBits) + v)
		if atomic.CompareAndSwapUint64(&hotCounts.sumBits, oldBits, newBits) {
			break
		}
	}
	
	
	atomic.AddUint64(&hotCounts.count, 1)
}

func (h *histogram) Write(out *dto.Metric) error {
	var (
		his                   = &dto.Histogram{}
		buckets               = make([]*dto.Bucket, len(h.upperBounds))
		hotCounts, coldCounts *histogramCounts
		count                 uint64
	)

	
	
	
	h.writeMtx.Lock()
	defer h.writeMtx.Unlock()

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	if h.hotIdx == 0 {
		count = atomic.AddUint64(&h.countAndHotIdx, 1) >> 1
		h.hotIdx = 1
		hotCounts = h.counts[1]
		coldCounts = h.counts[0]
	} else {
		count = atomic.AddUint64(&h.countAndHotIdx, ^uint64(0)) >> 1 
		h.hotIdx = 0
		hotCounts = h.counts[0]
		coldCounts = h.counts[1]
	}

	
	
	
	
	for {
		if count == atomic.LoadUint64(&coldCounts.count) {
			break
		}
		runtime.Gosched() 
	}

	his.SampleCount = proto.Uint64(count)
	his.SampleSum = proto.Float64(math.Float64frombits(atomic.LoadUint64(&coldCounts.sumBits)))
	var cumCount uint64
	for i, upperBound := range h.upperBounds {
		cumCount += atomic.LoadUint64(&coldCounts.buckets[i])
		buckets[i] = &dto.Bucket{
			CumulativeCount: proto.Uint64(cumCount),
			UpperBound:      proto.Float64(upperBound),
		}
	}

	his.Bucket = buckets
	out.Histogram = his
	out.Label = h.labelPairs

	
	atomic.AddUint64(&hotCounts.count, count)
	atomic.StoreUint64(&coldCounts.count, 0)
	for {
		oldBits := atomic.LoadUint64(&hotCounts.sumBits)
		newBits := math.Float64bits(math.Float64frombits(oldBits) + his.GetSampleSum())
		if atomic.CompareAndSwapUint64(&hotCounts.sumBits, oldBits, newBits) {
			atomic.StoreUint64(&coldCounts.sumBits, 0)
			break
		}
	}
	for i := range h.upperBounds {
		atomic.AddUint64(&hotCounts.buckets[i], atomic.LoadUint64(&coldCounts.buckets[i]))
		atomic.StoreUint64(&coldCounts.buckets[i], 0)
	}
	return nil
}






type HistogramVec struct {
	*metricVec
}



func NewHistogramVec(opts HistogramOpts, labelNames []string) *HistogramVec {
	desc := NewDesc(
		BuildFQName(opts.Namespace, opts.Subsystem, opts.Name),
		opts.Help,
		labelNames,
		opts.ConstLabels,
	)
	return &HistogramVec{
		metricVec: newMetricVec(desc, func(lvs ...string) Metric {
			return newHistogram(desc, opts, lvs...)
		}),
	}
}

























func (v *HistogramVec) GetMetricWithLabelValues(lvs ...string) (Observer, error) {
	metric, err := v.metricVec.getMetricWithLabelValues(lvs...)
	if metric != nil {
		return metric.(Observer), err
	}
	return nil, err
}













func (v *HistogramVec) GetMetricWith(labels Labels) (Observer, error) {
	metric, err := v.metricVec.getMetricWith(labels)
	if metric != nil {
		return metric.(Observer), err
	}
	return nil, err
}





func (v *HistogramVec) WithLabelValues(lvs ...string) Observer {
	h, err := v.GetMetricWithLabelValues(lvs...)
	if err != nil {
		panic(err)
	}
	return h
}




func (v *HistogramVec) With(labels Labels) Observer {
	h, err := v.GetMetricWith(labels)
	if err != nil {
		panic(err)
	}
	return h
}














func (v *HistogramVec) CurryWith(labels Labels) (ObserverVec, error) {
	vec, err := v.curryWith(labels)
	if vec != nil {
		return &HistogramVec{vec}, err
	}
	return nil, err
}



func (v *HistogramVec) MustCurryWith(labels Labels) ObserverVec {
	vec, err := v.CurryWith(labels)
	if err != nil {
		panic(err)
	}
	return vec
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
	if desc.err != nil {
		return nil, desc.err
	}
	if err := validateLabelValues(labelValues, len(desc.variableLabels)); err != nil {
		return nil, err
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
