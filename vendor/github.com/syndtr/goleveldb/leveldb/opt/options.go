






package opt

import (
	"math"

	"github.com/syndtr/goleveldb/leveldb/cache"
	"github.com/syndtr/goleveldb/leveldb/comparer"
	"github.com/syndtr/goleveldb/leveldb/filter"
)

const (
	KiB = 1024
	MiB = KiB * 1024
	GiB = MiB * 1024
)

var (
	DefaultBlockCacher                   = LRUCacher
	DefaultBlockCacheCapacity            = 8 * MiB
	DefaultBlockRestartInterval          = 16
	DefaultBlockSize                     = 4 * KiB
	DefaultCompactionExpandLimitFactor   = 25
	DefaultCompactionGPOverlapsFactor    = 10
	DefaultCompactionL0Trigger           = 4
	DefaultCompactionSourceLimitFactor   = 1
	DefaultCompactionTableSize           = 2 * MiB
	DefaultCompactionTableSizeMultiplier = 1.0
	DefaultCompactionTotalSize           = 10 * MiB
	DefaultCompactionTotalSizeMultiplier = 10.0
	DefaultCompressionType               = SnappyCompression
	DefaultIteratorSamplingRate          = 1 * MiB
	DefaultOpenFilesCacher               = LRUCacher
	DefaultOpenFilesCacheCapacity        = 500
	DefaultWriteBuffer                   = 4 * MiB
	DefaultWriteL0PauseTrigger           = 12
	DefaultWriteL0SlowdownTrigger        = 8
)


type Cacher interface {
	New(capacity int) cache.Cacher
}

type CacherFunc struct {
	NewFunc func(capacity int) cache.Cacher
}

func (f *CacherFunc) New(capacity int) cache.Cacher {
	if f.NewFunc != nil {
		return f.NewFunc(capacity)
	}
	return nil
}

func noCacher(int) cache.Cacher { return nil }

var (
	
	LRUCacher = &CacherFunc{cache.NewLRU}

	
	NoCacher = &CacherFunc{}
)


type Compression uint

func (c Compression) String() string {
	switch c {
	case DefaultCompression:
		return "default"
	case NoCompression:
		return "none"
	case SnappyCompression:
		return "snappy"
	}
	return "invalid"
}

const (
	DefaultCompression Compression = iota
	NoCompression
	SnappyCompression
	nCompression
)


type Strict uint

const (
	
	
	
	StrictManifest Strict = 1 << iota

	
	StrictJournalChecksum

	
	
	
	StrictJournal

	
	
	StrictBlockChecksum

	
	
	StrictCompaction

	
	StrictReader

	
	StrictRecovery

	
	
	StrictOverride

	
	StrictAll = StrictManifest | StrictJournalChecksum | StrictJournal | StrictBlockChecksum | StrictCompaction | StrictReader | StrictRecovery

	
	
	DefaultStrict = StrictJournalChecksum | StrictBlockChecksum | StrictCompaction | StrictReader

	
	NoStrict = ^StrictAll
)


type Options struct {
	
	
	
	
	
	AltFilters []filter.Filter

	
	
	
	
	BlockCacher Cacher

	
	
	
	
	BlockCacheCapacity int

	
	
	
	
	BlockRestartInterval int

	
	
	
	
	BlockSize int

	
	
	
	
	CompactionExpandLimitFactor int

	
	
	
	
	
	CompactionGPOverlapsFactor int

	
	
	
	
	CompactionL0Trigger int

	
	
	
	
	
	CompactionSourceLimitFactor int

	
	
	
	
	
	
	CompactionTableSize int

	
	
	
	CompactionTableSizeMultiplier float64

	
	
	
	
	
	CompactionTableSizeMultiplierPerLevel []float64

	
	
	
	
	
	
	
	CompactionTotalSize int

	
	
	
	CompactionTotalSizeMultiplier float64

	
	
	
	
	
	CompactionTotalSizeMultiplierPerLevel []float64

	
	
	
	
	
	Comparer comparer.Comparer

	
	
	
	Compression Compression

	
	
	
	DisableBufferPool bool

	
	
	
	
	DisableBlockCache bool

	
	
	
	DisableCompactionBackoff bool

	
	
	
	
	
	DisableLargeBatchTransaction bool

	
	
	
	
	ErrorIfExist bool

	
	
	
	
	
	ErrorIfMissing bool

	
	
	
	
	
	
	
	
	
	
	
	
	
	Filter filter.Filter

	
	
	
	
	
	IteratorSamplingRate int

	
	
	
	NoSync bool

	
	
	
	NoWriteMerge bool

	
	
	
	
	OpenFilesCacher Cacher

	
	
	
	
	OpenFilesCacheCapacity int

	
	
	
	ReadOnly bool

	
	Strict Strict

	
	
	
	
	
	
	
	WriteBuffer int

	
	
	
	
	WriteL0PauseTrigger int

	
	
	
	
	WriteL0SlowdownTrigger int
}

func (o *Options) GetAltFilters() []filter.Filter {
	if o == nil {
		return nil
	}
	return o.AltFilters
}

func (o *Options) GetBlockCacher() Cacher {
	if o == nil || o.BlockCacher == nil {
		return DefaultBlockCacher
	} else if o.BlockCacher == NoCacher {
		return nil
	}
	return o.BlockCacher
}

func (o *Options) GetBlockCacheCapacity() int {
	if o == nil || o.BlockCacheCapacity == 0 {
		return DefaultBlockCacheCapacity
	} else if o.BlockCacheCapacity < 0 {
		return 0
	}
	return o.BlockCacheCapacity
}

func (o *Options) GetBlockRestartInterval() int {
	if o == nil || o.BlockRestartInterval <= 0 {
		return DefaultBlockRestartInterval
	}
	return o.BlockRestartInterval
}

func (o *Options) GetBlockSize() int {
	if o == nil || o.BlockSize <= 0 {
		return DefaultBlockSize
	}
	return o.BlockSize
}

func (o *Options) GetCompactionExpandLimit(level int) int {
	factor := DefaultCompactionExpandLimitFactor
	if o != nil && o.CompactionExpandLimitFactor > 0 {
		factor = o.CompactionExpandLimitFactor
	}
	return o.GetCompactionTableSize(level+1) * factor
}

func (o *Options) GetCompactionGPOverlaps(level int) int {
	factor := DefaultCompactionGPOverlapsFactor
	if o != nil && o.CompactionGPOverlapsFactor > 0 {
		factor = o.CompactionGPOverlapsFactor
	}
	return o.GetCompactionTableSize(level+2) * factor
}

func (o *Options) GetCompactionL0Trigger() int {
	if o == nil || o.CompactionL0Trigger == 0 {
		return DefaultCompactionL0Trigger
	}
	return o.CompactionL0Trigger
}

func (o *Options) GetCompactionSourceLimit(level int) int {
	factor := DefaultCompactionSourceLimitFactor
	if o != nil && o.CompactionSourceLimitFactor > 0 {
		factor = o.CompactionSourceLimitFactor
	}
	return o.GetCompactionTableSize(level+1) * factor
}

func (o *Options) GetCompactionTableSize(level int) int {
	var (
		base = DefaultCompactionTableSize
		mult float64
	)
	if o != nil {
		if o.CompactionTableSize > 0 {
			base = o.CompactionTableSize
		}
		if level < len(o.CompactionTableSizeMultiplierPerLevel) && o.CompactionTableSizeMultiplierPerLevel[level] > 0 {
			mult = o.CompactionTableSizeMultiplierPerLevel[level]
		} else if o.CompactionTableSizeMultiplier > 0 {
			mult = math.Pow(o.CompactionTableSizeMultiplier, float64(level))
		}
	}
	if mult == 0 {
		mult = math.Pow(DefaultCompactionTableSizeMultiplier, float64(level))
	}
	return int(float64(base) * mult)
}

func (o *Options) GetCompactionTotalSize(level int) int64 {
	var (
		base = DefaultCompactionTotalSize
		mult float64
	)
	if o != nil {
		if o.CompactionTotalSize > 0 {
			base = o.CompactionTotalSize
		}
		if level < len(o.CompactionTotalSizeMultiplierPerLevel) && o.CompactionTotalSizeMultiplierPerLevel[level] > 0 {
			mult = o.CompactionTotalSizeMultiplierPerLevel[level]
		} else if o.CompactionTotalSizeMultiplier > 0 {
			mult = math.Pow(o.CompactionTotalSizeMultiplier, float64(level))
		}
	}
	if mult == 0 {
		mult = math.Pow(DefaultCompactionTotalSizeMultiplier, float64(level))
	}
	return int64(float64(base) * mult)
}

func (o *Options) GetComparer() comparer.Comparer {
	if o == nil || o.Comparer == nil {
		return comparer.DefaultComparer
	}
	return o.Comparer
}

func (o *Options) GetCompression() Compression {
	if o == nil || o.Compression <= DefaultCompression || o.Compression >= nCompression {
		return DefaultCompressionType
	}
	return o.Compression
}

func (o *Options) GetDisableBufferPool() bool {
	if o == nil {
		return false
	}
	return o.DisableBufferPool
}

func (o *Options) GetDisableBlockCache() bool {
	if o == nil {
		return false
	}
	return o.DisableBlockCache
}

func (o *Options) GetDisableCompactionBackoff() bool {
	if o == nil {
		return false
	}
	return o.DisableCompactionBackoff
}

func (o *Options) GetDisableLargeBatchTransaction() bool {
	if o == nil {
		return false
	}
	return o.DisableLargeBatchTransaction
}

func (o *Options) GetErrorIfExist() bool {
	if o == nil {
		return false
	}
	return o.ErrorIfExist
}

func (o *Options) GetErrorIfMissing() bool {
	if o == nil {
		return false
	}
	return o.ErrorIfMissing
}

func (o *Options) GetFilter() filter.Filter {
	if o == nil {
		return nil
	}
	return o.Filter
}

func (o *Options) GetIteratorSamplingRate() int {
	if o == nil || o.IteratorSamplingRate <= 0 {
		return DefaultIteratorSamplingRate
	}
	return o.IteratorSamplingRate
}

func (o *Options) GetNoSync() bool {
	if o == nil {
		return false
	}
	return o.NoSync
}

func (o *Options) GetNoWriteMerge() bool {
	if o == nil {
		return false
	}
	return o.NoWriteMerge
}

func (o *Options) GetOpenFilesCacher() Cacher {
	if o == nil || o.OpenFilesCacher == nil {
		return DefaultOpenFilesCacher
	}
	if o.OpenFilesCacher == NoCacher {
		return nil
	}
	return o.OpenFilesCacher
}

func (o *Options) GetOpenFilesCacheCapacity() int {
	if o == nil || o.OpenFilesCacheCapacity == 0 {
		return DefaultOpenFilesCacheCapacity
	} else if o.OpenFilesCacheCapacity < 0 {
		return 0
	}
	return o.OpenFilesCacheCapacity
}

func (o *Options) GetReadOnly() bool {
	if o == nil {
		return false
	}
	return o.ReadOnly
}

func (o *Options) GetStrict(strict Strict) bool {
	if o == nil || o.Strict == 0 {
		return DefaultStrict&strict != 0
	}
	return o.Strict&strict != 0
}

func (o *Options) GetWriteBuffer() int {
	if o == nil || o.WriteBuffer <= 0 {
		return DefaultWriteBuffer
	}
	return o.WriteBuffer
}

func (o *Options) GetWriteL0PauseTrigger() int {
	if o == nil || o.WriteL0PauseTrigger == 0 {
		return DefaultWriteL0PauseTrigger
	}
	return o.WriteL0PauseTrigger
}

func (o *Options) GetWriteL0SlowdownTrigger() int {
	if o == nil || o.WriteL0SlowdownTrigger == 0 {
		return DefaultWriteL0SlowdownTrigger
	}
	return o.WriteL0SlowdownTrigger
}



type ReadOptions struct {
	
	
	
	
	
	DontFillCache bool

	
	
	Strict Strict
}

func (ro *ReadOptions) GetDontFillCache() bool {
	if ro == nil {
		return false
	}
	return ro.DontFillCache
}

func (ro *ReadOptions) GetStrict(strict Strict) bool {
	if ro == nil {
		return false
	}
	return ro.Strict&strict != 0
}



type WriteOptions struct {
	
	
	
	NoWriteMerge bool

	
	
	
	
	
	
	
	
	
	
	
	
	Sync bool
}

func (wo *WriteOptions) GetNoWriteMerge() bool {
	if wo == nil {
		return false
	}
	return wo.NoWriteMerge
}

func (wo *WriteOptions) GetSync() bool {
	if wo == nil {
		return false
	}
	return wo.Sync
}

func GetStrict(o *Options, ro *ReadOptions, strict Strict) bool {
	if ro.GetStrict(StrictOverride) {
		return ro.GetStrict(strict)
	} else {
		return o.GetStrict(strict) || ro.GetStrict(strict)
	}
}
