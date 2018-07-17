




package timeseries 

import (
	"fmt"
	"log"
	"time"
)

const (
	timeSeriesNumBuckets       = 64
	minuteHourSeriesNumBuckets = 60
)

var timeSeriesResolutions = []time.Duration{
	1 * time.Second,
	10 * time.Second,
	1 * time.Minute,
	10 * time.Minute,
	1 * time.Hour,
	6 * time.Hour,
	24 * time.Hour,          
	7 * 24 * time.Hour,      
	4 * 7 * 24 * time.Hour,  
	16 * 7 * 24 * time.Hour, 
}

var minuteHourSeriesResolutions = []time.Duration{
	1 * time.Second,
	1 * time.Minute,
}


type Observable interface {
	Multiply(ratio float64)    
	Add(other Observable)      
	Clear()                    
	CopyFrom(other Observable) 
}


type Float float64


func NewFloat() Observable {
	f := Float(0)
	return &f
}


func (f *Float) String() string { return fmt.Sprintf("%g", f.Value()) }


func (f *Float) Value() float64 { return float64(*f) }

func (f *Float) Multiply(ratio float64) { *f *= Float(ratio) }

func (f *Float) Add(other Observable) {
	o := other.(*Float)
	*f += *o
}

func (f *Float) Clear() { *f = 0 }

func (f *Float) CopyFrom(other Observable) {
	o := other.(*Float)
	*f = *o
}


type Clock interface {
	Time() time.Time
}

type defaultClock int

var defaultClockInstance defaultClock

func (defaultClock) Time() time.Time { return time.Now() }




type tsLevel struct {
	oldest   int               
	newest   int               
	end      time.Time         
	size     time.Duration     
	buckets  []Observable      
	provider func() Observable 
}

func (l *tsLevel) Clear() {
	l.oldest = 0
	l.newest = len(l.buckets) - 1
	l.end = time.Time{}
	for i := range l.buckets {
		if l.buckets[i] != nil {
			l.buckets[i].Clear()
			l.buckets[i] = nil
		}
	}
}

func (l *tsLevel) InitLevel(size time.Duration, numBuckets int, f func() Observable) {
	l.size = size
	l.provider = f
	l.buckets = make([]Observable, numBuckets)
}









type timeSeries struct {
	provider    func() Observable 
	numBuckets  int               
	levels      []*tsLevel        
	lastAdd     time.Time         
	total       Observable        
	clock       Clock             
	pending     Observable        
	pendingTime time.Time         
	dirty       bool              
}


func (ts *timeSeries) init(resolutions []time.Duration, f func() Observable, numBuckets int, clock Clock) {
	ts.provider = f
	ts.numBuckets = numBuckets
	ts.clock = clock
	ts.levels = make([]*tsLevel, len(resolutions))

	for i := range resolutions {
		if i > 0 && resolutions[i-1] >= resolutions[i] {
			log.Print("timeseries: resolutions must be monotonically increasing")
			break
		}
		newLevel := new(tsLevel)
		newLevel.InitLevel(resolutions[i], ts.numBuckets, ts.provider)
		ts.levels[i] = newLevel
	}

	ts.Clear()
}


func (ts *timeSeries) Clear() {
	ts.lastAdd = time.Time{}
	ts.total = ts.resetObservation(ts.total)
	ts.pending = ts.resetObservation(ts.pending)
	ts.pendingTime = time.Time{}
	ts.dirty = false

	for i := range ts.levels {
		ts.levels[i].Clear()
	}
}


func (ts *timeSeries) Add(observation Observable) {
	ts.AddWithTime(observation, ts.clock.Time())
}


func (ts *timeSeries) AddWithTime(observation Observable, t time.Time) {

	smallBucketDuration := ts.levels[0].size

	if t.After(ts.lastAdd) {
		ts.lastAdd = t
	}

	if t.After(ts.pendingTime) {
		ts.advance(t)
		ts.mergePendingUpdates()
		ts.pendingTime = ts.levels[0].end
		ts.pending.CopyFrom(observation)
		ts.dirty = true
	} else if t.After(ts.pendingTime.Add(-1 * smallBucketDuration)) {
		
		
		
		ts.pending.Add(observation)
		ts.dirty = true
	} else {
		ts.mergeValue(observation, t)
	}
}


func (ts *timeSeries) mergeValue(observation Observable, t time.Time) {
	for _, level := range ts.levels {
		index := (ts.numBuckets - 1) - int(level.end.Sub(t)/level.size)
		if 0 <= index && index < ts.numBuckets {
			bucketNumber := (level.oldest + index) % ts.numBuckets
			if level.buckets[bucketNumber] == nil {
				level.buckets[bucketNumber] = level.provider()
			}
			level.buckets[bucketNumber].Add(observation)
		}
	}
	ts.total.Add(observation)
}


func (ts *timeSeries) mergePendingUpdates() {
	if ts.dirty {
		ts.mergeValue(ts.pending, ts.pendingTime)
		ts.pending = ts.resetObservation(ts.pending)
		ts.dirty = false
	}
}



func (ts *timeSeries) advance(t time.Time) {
	if !t.After(ts.levels[0].end) {
		return
	}
	for i := 0; i < len(ts.levels); i++ {
		level := ts.levels[i]
		if !level.end.Before(t) {
			break
		}

		
		
		if !t.Before(level.end.Add(level.size * time.Duration(ts.numBuckets))) {
			for _, b := range level.buckets {
				ts.resetObservation(b)
			}
			level.end = time.Unix(0, (t.UnixNano()/level.size.Nanoseconds())*level.size.Nanoseconds())
		}

		for t.After(level.end) {
			level.end = level.end.Add(level.size)
			level.newest = level.oldest
			level.oldest = (level.oldest + 1) % ts.numBuckets
			ts.resetObservation(level.buckets[level.newest])
		}

		t = level.end
	}
}


func (ts *timeSeries) Latest(level, num int) Observable {
	now := ts.clock.Time()
	if ts.levels[0].end.Before(now) {
		ts.advance(now)
	}

	ts.mergePendingUpdates()

	result := ts.provider()
	l := ts.levels[level]
	index := l.newest

	for i := 0; i < num; i++ {
		if l.buckets[index] != nil {
			result.Add(l.buckets[index])
		}
		if index == 0 {
			index = ts.numBuckets
		}
		index--
	}

	return result
}


func (ts *timeSeries) LatestBuckets(level, num int) []Observable {
	if level < 0 || level > len(ts.levels) {
		log.Print("timeseries: bad level argument: ", level)
		return nil
	}
	if num < 0 || num >= ts.numBuckets {
		log.Print("timeseries: bad num argument: ", num)
		return nil
	}

	results := make([]Observable, num)
	now := ts.clock.Time()
	if ts.levels[0].end.Before(now) {
		ts.advance(now)
	}

	ts.mergePendingUpdates()

	l := ts.levels[level]
	index := l.newest

	for i := 0; i < num; i++ {
		result := ts.provider()
		results[i] = result
		if l.buckets[index] != nil {
			result.CopyFrom(l.buckets[index])
		}

		if index == 0 {
			index = ts.numBuckets
		}
		index -= 1
	}
	return results
}


func (ts *timeSeries) ScaleBy(factor float64) {
	for _, l := range ts.levels {
		for i := 0; i < ts.numBuckets; i++ {
			l.buckets[i].Multiply(factor)
		}
	}

	ts.total.Multiply(factor)
	ts.pending.Multiply(factor)
}




func (ts *timeSeries) Range(start, finish time.Time) Observable {
	return ts.ComputeRange(start, finish, 1)[0]
}


func (ts *timeSeries) Recent(delta time.Duration) Observable {
	now := ts.clock.Time()
	return ts.Range(now.Add(-delta), now)
}


func (ts *timeSeries) Total() Observable {
	ts.mergePendingUpdates()
	return ts.total
}






func (ts *timeSeries) ComputeRange(start, finish time.Time, num int) []Observable {
	if start.After(finish) {
		log.Printf("timeseries: start > finish, %v>%v", start, finish)
		return nil
	}

	if num < 0 {
		log.Printf("timeseries: num < 0, %v", num)
		return nil
	}

	results := make([]Observable, num)

	for _, l := range ts.levels {
		if !start.Before(l.end.Add(-l.size * time.Duration(ts.numBuckets))) {
			ts.extract(l, start, finish, num, results)
			return results
		}
	}

	
	
	
	ts.extract(ts.levels[len(ts.levels)-1], start, finish, num, results)

	return results
}



func (ts *timeSeries) RecentList(delta time.Duration, num int) []Observable {
	if delta < 0 {
		return nil
	}
	now := ts.clock.Time()
	return ts.ComputeRange(now.Add(-delta), now, num)
}



func (ts *timeSeries) extract(l *tsLevel, start, finish time.Time, num int, results []Observable) {
	ts.mergePendingUpdates()

	srcInterval := l.size
	dstInterval := finish.Sub(start) / time.Duration(num)
	dstStart := start
	srcStart := l.end.Add(-srcInterval * time.Duration(ts.numBuckets))

	srcIndex := 0

	
	if dstStart.After(srcStart) {
		advance := dstStart.Sub(srcStart) / srcInterval
		srcIndex += int(advance)
		srcStart = srcStart.Add(advance * srcInterval)
	}

	
	
	
	
	
	for i := 0; i < num; i++ {
		results[i] = ts.resetObservation(results[i])
		dstEnd := dstStart.Add(dstInterval)
		for srcIndex < ts.numBuckets && srcStart.Before(dstEnd) {
			srcEnd := srcStart.Add(srcInterval)
			if srcEnd.After(ts.lastAdd) {
				srcEnd = ts.lastAdd
			}

			if !srcEnd.Before(dstStart) {
				srcValue := l.buckets[(srcIndex+l.oldest)%ts.numBuckets]
				if !srcStart.Before(dstStart) && !srcEnd.After(dstEnd) {
					
					if srcValue != nil {
						results[i].Add(srcValue)
					}
				} else {
					
					overlapStart := maxTime(srcStart, dstStart)
					overlapEnd := minTime(srcEnd, dstEnd)
					base := srcEnd.Sub(srcStart)
					fraction := overlapEnd.Sub(overlapStart).Seconds() / base.Seconds()

					used := ts.provider()
					if srcValue != nil {
						used.CopyFrom(srcValue)
					}
					used.Multiply(fraction)
					results[i].Add(used)
				}

				if srcEnd.After(dstEnd) {
					break
				}
			}
			srcIndex++
			srcStart = srcStart.Add(srcInterval)
		}
		dstStart = dstStart.Add(dstInterval)
	}
}


func (ts *timeSeries) resetObservation(observation Observable) Observable {
	if observation == nil {
		observation = ts.provider()
	} else {
		observation.Clear()
	}
	return observation
}


type TimeSeries struct {
	timeSeries
}


func NewTimeSeries(f func() Observable) *TimeSeries {
	return NewTimeSeriesWithClock(f, defaultClockInstance)
}



func NewTimeSeriesWithClock(f func() Observable, clock Clock) *TimeSeries {
	ts := new(TimeSeries)
	ts.timeSeries.init(timeSeriesResolutions, f, timeSeriesNumBuckets, clock)
	return ts
}


type MinuteHourSeries struct {
	timeSeries
}


func NewMinuteHourSeries(f func() Observable) *MinuteHourSeries {
	return NewMinuteHourSeriesWithClock(f, defaultClockInstance)
}



func NewMinuteHourSeriesWithClock(f func() Observable, clock Clock) *MinuteHourSeries {
	ts := new(MinuteHourSeries)
	ts.timeSeries.init(minuteHourSeriesResolutions, f,
		minuteHourSeriesNumBuckets, clock)
	return ts
}

func (ts *MinuteHourSeries) Minute() Observable {
	return ts.timeSeries.Latest(0, 60)
}

func (ts *MinuteHourSeries) Hour() Observable {
	return ts.timeSeries.Latest(1, 60)
}

func minTime(a, b time.Time) time.Time {
	if a.Before(b) {
		return a
	}
	return b
}

func maxTime(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}
