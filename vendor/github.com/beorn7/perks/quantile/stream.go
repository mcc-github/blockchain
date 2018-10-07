













package quantile

import (
	"math"
	"sort"
)



type Sample struct {
	Value float64 `json:",string"`
	Width float64 `json:",string"`
	Delta float64 `json:",string"`
}


type Samples []Sample

func (a Samples) Len() int           { return len(a) }
func (a Samples) Less(i, j int) bool { return a[i].Value < a[j].Value }
func (a Samples) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type invariant func(s *stream, r float64) float64











func NewLowBiased(epsilon float64) *Stream {
	ƒ := func(s *stream, r float64) float64 {
		return 2 * epsilon * r
	}
	return newStream(ƒ)
}











func NewHighBiased(epsilon float64) *Stream {
	ƒ := func(s *stream, r float64) float64 {
		return 2 * epsilon * (s.n - r)
	}
	return newStream(ƒ)
}








func NewTargeted(targetMap map[float64]float64) *Stream {
	
	
	
	targets := targetMapToSlice(targetMap)

	ƒ := func(s *stream, r float64) float64 {
		var m = math.MaxFloat64
		var f float64
		for _, t := range targets {
			if t.quantile*s.n <= r {
				f = (2 * t.epsilon * r) / t.quantile
			} else {
				f = (2 * t.epsilon * (s.n - r)) / (1 - t.quantile)
			}
			if f < m {
				m = f
			}
		}
		return m
	}
	return newStream(ƒ)
}

type target struct {
	quantile float64
	epsilon  float64
}

func targetMapToSlice(targetMap map[float64]float64) []target {
	targets := make([]target, 0, len(targetMap))

	for quantile, epsilon := range targetMap {
		t := target{
			quantile: quantile,
			epsilon:  epsilon,
		}
		targets = append(targets, t)
	}

	return targets
}



type Stream struct {
	*stream
	b      Samples
	sorted bool
}

func newStream(ƒ invariant) *Stream {
	x := &stream{ƒ: ƒ}
	return &Stream{x, make(Samples, 0, 500), true}
}


func (s *Stream) Insert(v float64) {
	s.insert(Sample{Value: v, Width: 1})
}

func (s *Stream) insert(sample Sample) {
	s.b = append(s.b, sample)
	s.sorted = false
	if len(s.b) == cap(s.b) {
		s.flush()
	}
}




func (s *Stream) Query(q float64) float64 {
	if !s.flushed() {
		
		
		l := len(s.b)
		if l == 0 {
			return 0
		}
		i := int(math.Ceil(float64(l) * q))
		if i > 0 {
			i -= 1
		}
		s.maybeSort()
		return s.b[i].Value
	}
	s.flush()
	return s.stream.query(q)
}






func (s *Stream) Merge(samples Samples) {
	sort.Sort(samples)
	s.stream.merge(samples)
}


func (s *Stream) Reset() {
	s.stream.reset()
	s.b = s.b[:0]
}


func (s *Stream) Samples() Samples {
	if !s.flushed() {
		return s.b
	}
	s.flush()
	return s.stream.samples()
}



func (s *Stream) Count() int {
	return len(s.b) + s.stream.count()
}

func (s *Stream) flush() {
	s.maybeSort()
	s.stream.merge(s.b)
	s.b = s.b[:0]
}

func (s *Stream) maybeSort() {
	if !s.sorted {
		s.sorted = true
		sort.Sort(s.b)
	}
}

func (s *Stream) flushed() bool {
	return len(s.stream.l) > 0
}

type stream struct {
	n float64
	l []Sample
	ƒ invariant
}

func (s *stream) reset() {
	s.l = s.l[:0]
	s.n = 0
}

func (s *stream) insert(v float64) {
	s.merge(Samples{{v, 1, 0}})
}

func (s *stream) merge(samples Samples) {
	
	
	
	
	var r float64
	i := 0
	for _, sample := range samples {
		for ; i < len(s.l); i++ {
			c := s.l[i]
			if c.Value > sample.Value {
				
				s.l = append(s.l, Sample{})
				copy(s.l[i+1:], s.l[i:])
				s.l[i] = Sample{
					sample.Value,
					sample.Width,
					math.Max(sample.Delta, math.Floor(s.ƒ(s, r))-1),
					
				}
				i++
				goto inserted
			}
			r += c.Width
		}
		s.l = append(s.l, Sample{sample.Value, sample.Width, 0})
		i++
	inserted:
		s.n += sample.Width
		r += sample.Width
	}
	s.compress()
}

func (s *stream) count() int {
	return int(s.n)
}

func (s *stream) query(q float64) float64 {
	t := math.Ceil(q * s.n)
	t += math.Ceil(s.ƒ(s, t) / 2)
	p := s.l[0]
	var r float64
	for _, c := range s.l[1:] {
		r += p.Width
		if r+c.Width+c.Delta > t {
			return p.Value
		}
		p = c
	}
	return p.Value
}

func (s *stream) compress() {
	if len(s.l) < 2 {
		return
	}
	x := s.l[len(s.l)-1]
	xi := len(s.l) - 1
	r := s.n - 1 - x.Width

	for i := len(s.l) - 2; i >= 0; i-- {
		c := s.l[i]
		if c.Width+x.Width+x.Delta <= s.ƒ(s, r) {
			x.Width += c.Width
			s.l[xi] = x
			
			copy(s.l[i:], s.l[i+1:])
			s.l = s.l[:len(s.l)-1]
			xi -= 1
		} else {
			x = c
			xi = i
		}
		r -= c.Width
	}
}

func (s *stream) samples() Samples {
	samples := make(Samples, len(s.l))
	copy(samples, s.l)
	return samples
}
