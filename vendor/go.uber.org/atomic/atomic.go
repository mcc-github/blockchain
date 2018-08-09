





















package atomic

import (
	"math"
	"sync/atomic"
	"time"
)


type Int32 struct{ v int32 }


func NewInt32(i int32) *Int32 {
	return &Int32{i}
}


func (i *Int32) Load() int32 {
	return atomic.LoadInt32(&i.v)
}


func (i *Int32) Add(n int32) int32 {
	return atomic.AddInt32(&i.v, n)
}


func (i *Int32) Sub(n int32) int32 {
	return atomic.AddInt32(&i.v, -n)
}


func (i *Int32) Inc() int32 {
	return i.Add(1)
}


func (i *Int32) Dec() int32 {
	return i.Sub(1)
}


func (i *Int32) CAS(old, new int32) bool {
	return atomic.CompareAndSwapInt32(&i.v, old, new)
}


func (i *Int32) Store(n int32) {
	atomic.StoreInt32(&i.v, n)
}


func (i *Int32) Swap(n int32) int32 {
	return atomic.SwapInt32(&i.v, n)
}


type Int64 struct{ v int64 }


func NewInt64(i int64) *Int64 {
	return &Int64{i}
}


func (i *Int64) Load() int64 {
	return atomic.LoadInt64(&i.v)
}


func (i *Int64) Add(n int64) int64 {
	return atomic.AddInt64(&i.v, n)
}


func (i *Int64) Sub(n int64) int64 {
	return atomic.AddInt64(&i.v, -n)
}


func (i *Int64) Inc() int64 {
	return i.Add(1)
}


func (i *Int64) Dec() int64 {
	return i.Sub(1)
}


func (i *Int64) CAS(old, new int64) bool {
	return atomic.CompareAndSwapInt64(&i.v, old, new)
}


func (i *Int64) Store(n int64) {
	atomic.StoreInt64(&i.v, n)
}


func (i *Int64) Swap(n int64) int64 {
	return atomic.SwapInt64(&i.v, n)
}


type Uint32 struct{ v uint32 }


func NewUint32(i uint32) *Uint32 {
	return &Uint32{i}
}


func (i *Uint32) Load() uint32 {
	return atomic.LoadUint32(&i.v)
}


func (i *Uint32) Add(n uint32) uint32 {
	return atomic.AddUint32(&i.v, n)
}


func (i *Uint32) Sub(n uint32) uint32 {
	return atomic.AddUint32(&i.v, ^(n - 1))
}


func (i *Uint32) Inc() uint32 {
	return i.Add(1)
}


func (i *Uint32) Dec() uint32 {
	return i.Sub(1)
}


func (i *Uint32) CAS(old, new uint32) bool {
	return atomic.CompareAndSwapUint32(&i.v, old, new)
}


func (i *Uint32) Store(n uint32) {
	atomic.StoreUint32(&i.v, n)
}


func (i *Uint32) Swap(n uint32) uint32 {
	return atomic.SwapUint32(&i.v, n)
}


type Uint64 struct{ v uint64 }


func NewUint64(i uint64) *Uint64 {
	return &Uint64{i}
}


func (i *Uint64) Load() uint64 {
	return atomic.LoadUint64(&i.v)
}


func (i *Uint64) Add(n uint64) uint64 {
	return atomic.AddUint64(&i.v, n)
}


func (i *Uint64) Sub(n uint64) uint64 {
	return atomic.AddUint64(&i.v, ^(n - 1))
}


func (i *Uint64) Inc() uint64 {
	return i.Add(1)
}


func (i *Uint64) Dec() uint64 {
	return i.Sub(1)
}


func (i *Uint64) CAS(old, new uint64) bool {
	return atomic.CompareAndSwapUint64(&i.v, old, new)
}


func (i *Uint64) Store(n uint64) {
	atomic.StoreUint64(&i.v, n)
}


func (i *Uint64) Swap(n uint64) uint64 {
	return atomic.SwapUint64(&i.v, n)
}


type Bool struct{ v uint32 }


func NewBool(initial bool) *Bool {
	return &Bool{boolToInt(initial)}
}


func (b *Bool) Load() bool {
	return truthy(atomic.LoadUint32(&b.v))
}


func (b *Bool) CAS(old, new bool) bool {
	return atomic.CompareAndSwapUint32(&b.v, boolToInt(old), boolToInt(new))
}


func (b *Bool) Store(new bool) {
	atomic.StoreUint32(&b.v, boolToInt(new))
}


func (b *Bool) Swap(new bool) bool {
	return truthy(atomic.SwapUint32(&b.v, boolToInt(new)))
}


func (b *Bool) Toggle() bool {
	return truthy(atomic.AddUint32(&b.v, 1) - 1)
}

func truthy(n uint32) bool {
	return n&1 == 1
}

func boolToInt(b bool) uint32 {
	if b {
		return 1
	}
	return 0
}


type Float64 struct {
	v uint64
}


func NewFloat64(f float64) *Float64 {
	return &Float64{math.Float64bits(f)}
}


func (f *Float64) Load() float64 {
	return math.Float64frombits(atomic.LoadUint64(&f.v))
}


func (f *Float64) Store(s float64) {
	atomic.StoreUint64(&f.v, math.Float64bits(s))
}


func (f *Float64) Add(s float64) float64 {
	for {
		old := f.Load()
		new := old + s
		if f.CAS(old, new) {
			return new
		}
	}
}


func (f *Float64) Sub(s float64) float64 {
	return f.Add(-s)
}


func (f *Float64) CAS(old, new float64) bool {
	return atomic.CompareAndSwapUint64(&f.v, math.Float64bits(old), math.Float64bits(new))
}



type Duration struct {
	v Int64
}


func NewDuration(d time.Duration) *Duration {
	return &Duration{v: *NewInt64(int64(d))}
}


func (d *Duration) Load() time.Duration {
	return time.Duration(d.v.Load())
}


func (d *Duration) Store(n time.Duration) {
	d.v.Store(int64(n))
}


func (d *Duration) Add(n time.Duration) time.Duration {
	return time.Duration(d.v.Add(int64(n)))
}


func (d *Duration) Sub(n time.Duration) time.Duration {
	return time.Duration(d.v.Sub(int64(n)))
}


func (d *Duration) Swap(n time.Duration) time.Duration {
	return time.Duration(d.v.Swap(int64(n)))
}


func (d *Duration) CAS(old, new time.Duration) bool {
	return d.v.CAS(int64(old), int64(new))
}



type Value struct{ atomic.Value }
