












package prometheus

import (
	"fmt"
	"math"
	"sort"
	"sync"
	"time"

	"github.com/beorn7/perks/quantile"
	"github.com/golang/protobuf/proto"

	dto "github.com/prometheus/client_model/go"
)



const quantileLabel = "quantile"
















type Summary interface {
	Metric
	Collector

	
	Observe(float64)
}


var (
	DefObjectives = map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001}

	errQuantileLabelNotAllowed = fmt.Errorf(
		"%q is not allowed as label name in summaries", quantileLabel,
	)
)


const (
	
	
	DefMaxAge time.Duration = 10 * time.Minute
	
	
	DefAgeBuckets = 5
	
	DefBufCap = 500
)




type SummaryOpts struct {
	
	
	
	
	
	Namespace string
	Subsystem string
	Name      string

	
	
	
	
	Help string

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	ConstLabels Labels

	
	
	
	
	Objectives map[float64]float64

	
	
	MaxAge time.Duration

	
	
	
	
	
	
	
	AgeBuckets uint32

	
	
	
	
	
	BufCap uint32
}

















func NewSummary(opts SummaryOpts) Summary {
	return newSummary(
		NewDesc(
			BuildFQName(opts.Namespace, opts.Subsystem, opts.Name),
			opts.Help,
			nil,
			opts.ConstLabels,
		),
		opts,
	)
}

func newSummary(desc *Desc, opts SummaryOpts, labelValues ...string) Summary {
	if len(desc.variableLabels) != len(labelValues) {
		panic(errInconsistentCardinality)
	}

	for _, n := range desc.variableLabels {
		if n == quantileLabel {
			panic(errQuantileLabelNotAllowed)
		}
	}
	for _, lp := range desc.constLabelPairs {
		if lp.GetName() == quantileLabel {
			panic(errQuantileLabelNotAllowed)
		}
	}

	if len(opts.Objectives) == 0 {
		opts.Objectives = DefObjectives
	}

	if opts.MaxAge < 0 {
		panic(fmt.Errorf("illegal max age MaxAge=%v", opts.MaxAge))
	}
	if opts.MaxAge == 0 {
		opts.MaxAge = DefMaxAge
	}

	if opts.AgeBuckets == 0 {
		opts.AgeBuckets = DefAgeBuckets
	}

	if opts.BufCap == 0 {
		opts.BufCap = DefBufCap
	}

	s := &summary{
		desc: desc,

		objectives:       opts.Objectives,
		sortedObjectives: make([]float64, 0, len(opts.Objectives)),

		labelPairs: makeLabelPairs(desc, labelValues),

		hotBuf:         make([]float64, 0, opts.BufCap),
		coldBuf:        make([]float64, 0, opts.BufCap),
		streamDuration: opts.MaxAge / time.Duration(opts.AgeBuckets),
	}
	s.headStreamExpTime = time.Now().Add(s.streamDuration)
	s.hotBufExpTime = s.headStreamExpTime

	for i := uint32(0); i < opts.AgeBuckets; i++ {
		s.streams = append(s.streams, s.newStream())
	}
	s.headStream = s.streams[0]

	for qu := range s.objectives {
		s.sortedObjectives = append(s.sortedObjectives, qu)
	}
	sort.Float64s(s.sortedObjectives)

	s.init(s) 
	return s
}

type summary struct {
	selfCollector

	bufMtx sync.Mutex 
	mtx    sync.Mutex 
	

	desc *Desc

	objectives       map[float64]float64
	sortedObjectives []float64

	labelPairs []*dto.LabelPair

	sum float64
	cnt uint64

	hotBuf, coldBuf []float64

	streams                          []*quantile.Stream
	streamDuration                   time.Duration
	headStream                       *quantile.Stream
	headStreamIdx                    int
	headStreamExpTime, hotBufExpTime time.Time
}

func (s *summary) Desc() *Desc {
	return s.desc
}

func (s *summary) Observe(v float64) {
	s.bufMtx.Lock()
	defer s.bufMtx.Unlock()

	now := time.Now()
	if now.After(s.hotBufExpTime) {
		s.asyncFlush(now)
	}
	s.hotBuf = append(s.hotBuf, v)
	if len(s.hotBuf) == cap(s.hotBuf) {
		s.asyncFlush(now)
	}
}

func (s *summary) Write(out *dto.Metric) error {
	sum := &dto.Summary{}
	qs := make([]*dto.Quantile, 0, len(s.objectives))

	s.bufMtx.Lock()
	s.mtx.Lock()
	
	s.swapBufs(time.Now())
	s.bufMtx.Unlock()

	s.flushColdBuf()
	sum.SampleCount = proto.Uint64(s.cnt)
	sum.SampleSum = proto.Float64(s.sum)

	for _, rank := range s.sortedObjectives {
		var q float64
		if s.headStream.Count() == 0 {
			q = math.NaN()
		} else {
			q = s.headStream.Query(rank)
		}
		qs = append(qs, &dto.Quantile{
			Quantile: proto.Float64(rank),
			Value:    proto.Float64(q),
		})
	}

	s.mtx.Unlock()

	if len(qs) > 0 {
		sort.Sort(quantSort(qs))
	}
	sum.Quantile = qs

	out.Summary = sum
	out.Label = s.labelPairs
	return nil
}

func (s *summary) newStream() *quantile.Stream {
	return quantile.NewTargeted(s.objectives)
}


func (s *summary) asyncFlush(now time.Time) {
	s.mtx.Lock()
	s.swapBufs(now)

	
	
	
	go func() {
		s.flushColdBuf()
		s.mtx.Unlock()
	}()
}


func (s *summary) maybeRotateStreams() {
	for !s.hotBufExpTime.Equal(s.headStreamExpTime) {
		s.headStream.Reset()
		s.headStreamIdx++
		if s.headStreamIdx >= len(s.streams) {
			s.headStreamIdx = 0
		}
		s.headStream = s.streams[s.headStreamIdx]
		s.headStreamExpTime = s.headStreamExpTime.Add(s.streamDuration)
	}
}


func (s *summary) flushColdBuf() {
	for _, v := range s.coldBuf {
		for _, stream := range s.streams {
			stream.Insert(v)
		}
		s.cnt++
		s.sum += v
	}
	s.coldBuf = s.coldBuf[0:0]
	s.maybeRotateStreams()
}


func (s *summary) swapBufs(now time.Time) {
	if len(s.coldBuf) != 0 {
		panic("coldBuf is not empty")
	}
	s.hotBuf, s.coldBuf = s.coldBuf, s.hotBuf
	
	for now.After(s.hotBufExpTime) {
		s.hotBufExpTime = s.hotBufExpTime.Add(s.streamDuration)
	}
}

type quantSort []*dto.Quantile

func (s quantSort) Len() int {
	return len(s)
}

func (s quantSort) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s quantSort) Less(i, j int) bool {
	return s[i].GetQuantile() < s[j].GetQuantile()
}






type SummaryVec struct {
	*MetricVec
}




func NewSummaryVec(opts SummaryOpts, labelNames []string) *SummaryVec {
	desc := NewDesc(
		BuildFQName(opts.Namespace, opts.Subsystem, opts.Name),
		opts.Help,
		labelNames,
		opts.ConstLabels,
	)
	return &SummaryVec{
		MetricVec: newMetricVec(desc, func(lvs ...string) Metric {
			return newSummary(desc, opts, lvs...)
		}),
	}
}




func (m *SummaryVec) GetMetricWithLabelValues(lvs ...string) (Summary, error) {
	metric, err := m.MetricVec.GetMetricWithLabelValues(lvs...)
	if metric != nil {
		return metric.(Summary), err
	}
	return nil, err
}




func (m *SummaryVec) GetMetricWith(labels Labels) (Summary, error) {
	metric, err := m.MetricVec.GetMetricWith(labels)
	if metric != nil {
		return metric.(Summary), err
	}
	return nil, err
}





func (m *SummaryVec) WithLabelValues(lvs ...string) Summary {
	return m.MetricVec.WithLabelValues(lvs...).(Summary)
}




func (m *SummaryVec) With(labels Labels) Summary {
	return m.MetricVec.With(labels).(Summary)
}

type constSummary struct {
	desc       *Desc
	count      uint64
	sum        float64
	quantiles  map[float64]float64
	labelPairs []*dto.LabelPair
}

func (s *constSummary) Desc() *Desc {
	return s.desc
}

func (s *constSummary) Write(out *dto.Metric) error {
	sum := &dto.Summary{}
	qs := make([]*dto.Quantile, 0, len(s.quantiles))

	sum.SampleCount = proto.Uint64(s.count)
	sum.SampleSum = proto.Float64(s.sum)

	for rank, q := range s.quantiles {
		qs = append(qs, &dto.Quantile{
			Quantile: proto.Float64(rank),
			Value:    proto.Float64(q),
		})
	}

	if len(qs) > 0 {
		sort.Sort(quantSort(qs))
	}
	sum.Quantile = qs

	out.Summary = sum
	out.Label = s.labelPairs

	return nil
}















func NewConstSummary(
	desc *Desc,
	count uint64,
	sum float64,
	quantiles map[float64]float64,
	labelValues ...string,
) (Metric, error) {
	if len(desc.variableLabels) != len(labelValues) {
		return nil, errInconsistentCardinality
	}
	return &constSummary{
		desc:       desc,
		count:      count,
		sum:        sum,
		quantiles:  quantiles,
		labelPairs: makeLabelPairs(desc, labelValues),
	}, nil
}



func MustNewConstSummary(
	desc *Desc,
	count uint64,
	sum float64,
	quantiles map[float64]float64,
	labelValues ...string,
) Metric {
	m, err := NewConstSummary(desc, count, sum, quantiles, labelValues...)
	if err != nil {
		panic(err)
	}
	return m
}
