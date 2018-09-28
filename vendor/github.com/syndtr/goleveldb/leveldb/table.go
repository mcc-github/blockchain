





package leveldb

import (
	"fmt"
	"sort"
	"sync/atomic"

	"github.com/syndtr/goleveldb/leveldb/cache"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/storage"
	"github.com/syndtr/goleveldb/leveldb/table"
	"github.com/syndtr/goleveldb/leveldb/util"
)


type tFile struct {
	fd         storage.FileDesc
	seekLeft   int32
	size       int64
	imin, imax internalKey
}


func (t *tFile) after(icmp *iComparer, ukey []byte) bool {
	return ukey != nil && icmp.uCompare(ukey, t.imax.ukey()) > 0
}


func (t *tFile) before(icmp *iComparer, ukey []byte) bool {
	return ukey != nil && icmp.uCompare(ukey, t.imin.ukey()) < 0
}


func (t *tFile) overlaps(icmp *iComparer, umin, umax []byte) bool {
	return !t.after(icmp, umin) && !t.before(icmp, umax)
}


func (t *tFile) consumeSeek() int32 {
	return atomic.AddInt32(&t.seekLeft, -1)
}


func newTableFile(fd storage.FileDesc, size int64, imin, imax internalKey) *tFile {
	f := &tFile{
		fd:   fd,
		size: size,
		imin: imin,
		imax: imax,
	}

	
	
	
	
	
	
	
	
	
	
	
	
	
	f.seekLeft = int32(size / 16384)
	if f.seekLeft < 100 {
		f.seekLeft = 100
	}

	return f
}

func tableFileFromRecord(r atRecord) *tFile {
	return newTableFile(storage.FileDesc{storage.TypeTable, r.num}, r.size, r.imin, r.imax)
}


type tFiles []*tFile

func (tf tFiles) Len() int      { return len(tf) }
func (tf tFiles) Swap(i, j int) { tf[i], tf[j] = tf[j], tf[i] }

func (tf tFiles) nums() string {
	x := "[ "
	for i, f := range tf {
		if i != 0 {
			x += ", "
		}
		x += fmt.Sprint(f.fd.Num)
	}
	x += " ]"
	return x
}



func (tf tFiles) lessByKey(icmp *iComparer, i, j int) bool {
	a, b := tf[i], tf[j]
	n := icmp.Compare(a.imin, b.imin)
	if n == 0 {
		return a.fd.Num < b.fd.Num
	}
	return n < 0
}



func (tf tFiles) lessByNum(i, j int) bool {
	return tf[i].fd.Num > tf[j].fd.Num
}


func (tf tFiles) sortByKey(icmp *iComparer) {
	sort.Sort(&tFilesSortByKey{tFiles: tf, icmp: icmp})
}


func (tf tFiles) sortByNum() {
	sort.Sort(&tFilesSortByNum{tFiles: tf})
}


func (tf tFiles) size() (sum int64) {
	for _, t := range tf {
		sum += t.size
	}
	return sum
}



func (tf tFiles) searchMin(icmp *iComparer, ikey internalKey) int {
	return sort.Search(len(tf), func(i int) bool {
		return icmp.Compare(tf[i].imin, ikey) >= 0
	})
}



func (tf tFiles) searchMax(icmp *iComparer, ikey internalKey) int {
	return sort.Search(len(tf), func(i int) bool {
		return icmp.Compare(tf[i].imax, ikey) >= 0
	})
}



func (tf tFiles) overlaps(icmp *iComparer, umin, umax []byte, unsorted bool) bool {
	if unsorted {
		
		for _, t := range tf {
			if t.overlaps(icmp, umin, umax) {
				return true
			}
		}
		return false
	}

	i := 0
	if len(umin) > 0 {
		
		i = tf.searchMax(icmp, makeInternalKey(nil, umin, keyMaxSeq, keyTypeSeek))
	}
	if i >= len(tf) {
		
		return false
	}
	return !tf[i].before(icmp, umax)
}






func (tf tFiles) getOverlaps(dst tFiles, icmp *iComparer, umin, umax []byte, overlapped bool) tFiles {
	dst = dst[:0]
	for i := 0; i < len(tf); {
		t := tf[i]
		if t.overlaps(icmp, umin, umax) {
			if umin != nil && icmp.uCompare(t.imin.ukey(), umin) < 0 {
				umin = t.imin.ukey()
				dst = dst[:0]
				i = 0
				continue
			} else if umax != nil && icmp.uCompare(t.imax.ukey(), umax) > 0 {
				umax = t.imax.ukey()
				
				if overlapped {
					dst = dst[:0]
					i = 0
					continue
				}
			}

			dst = append(dst, t)
		}
		i++
	}

	return dst
}


func (tf tFiles) getRange(icmp *iComparer) (imin, imax internalKey) {
	for i, t := range tf {
		if i == 0 {
			imin, imax = t.imin, t.imax
			continue
		}
		if icmp.Compare(t.imin, imin) < 0 {
			imin = t.imin
		}
		if icmp.Compare(t.imax, imax) > 0 {
			imax = t.imax
		}
	}

	return
}


func (tf tFiles) newIndexIterator(tops *tOps, icmp *iComparer, slice *util.Range, ro *opt.ReadOptions) iterator.IteratorIndexer {
	if slice != nil {
		var start, limit int
		if slice.Start != nil {
			start = tf.searchMax(icmp, internalKey(slice.Start))
		}
		if slice.Limit != nil {
			limit = tf.searchMin(icmp, internalKey(slice.Limit))
		} else {
			limit = tf.Len()
		}
		tf = tf[start:limit]
	}
	return iterator.NewArrayIndexer(&tFilesArrayIndexer{
		tFiles: tf,
		tops:   tops,
		icmp:   icmp,
		slice:  slice,
		ro:     ro,
	})
}


type tFilesArrayIndexer struct {
	tFiles
	tops  *tOps
	icmp  *iComparer
	slice *util.Range
	ro    *opt.ReadOptions
}

func (a *tFilesArrayIndexer) Search(key []byte) int {
	return a.searchMax(a.icmp, internalKey(key))
}

func (a *tFilesArrayIndexer) Get(i int) iterator.Iterator {
	if i == 0 || i == a.Len()-1 {
		return a.tops.newIterator(a.tFiles[i], a.slice, a.ro)
	}
	return a.tops.newIterator(a.tFiles[i], nil, a.ro)
}


type tFilesSortByKey struct {
	tFiles
	icmp *iComparer
}

func (x *tFilesSortByKey) Less(i, j int) bool {
	return x.lessByKey(x.icmp, i, j)
}


type tFilesSortByNum struct {
	tFiles
}

func (x *tFilesSortByNum) Less(i, j int) bool {
	return x.lessByNum(i, j)
}


type tOps struct {
	s      *session
	noSync bool
	cache  *cache.Cache
	bcache *cache.Cache
	bpool  *util.BufferPool
}


func (t *tOps) create() (*tWriter, error) {
	fd := storage.FileDesc{storage.TypeTable, t.s.allocFileNum()}
	fw, err := t.s.stor.Create(fd)
	if err != nil {
		return nil, err
	}
	return &tWriter{
		t:  t,
		fd: fd,
		w:  fw,
		tw: table.NewWriter(fw, t.s.o.Options),
	}, nil
}


func (t *tOps) createFrom(src iterator.Iterator) (f *tFile, n int, err error) {
	w, err := t.create()
	if err != nil {
		return
	}

	defer func() {
		if err != nil {
			w.drop()
		}
	}()

	for src.Next() {
		err = w.append(src.Key(), src.Value())
		if err != nil {
			return
		}
	}
	err = src.Error()
	if err != nil {
		return
	}

	n = w.tw.EntriesLen()
	f, err = w.finish()
	return
}



func (t *tOps) open(f *tFile) (ch *cache.Handle, err error) {
	ch = t.cache.Get(0, uint64(f.fd.Num), func() (size int, value cache.Value) {
		var r storage.Reader
		r, err = t.s.stor.Open(f.fd)
		if err != nil {
			return 0, nil
		}

		var bcache *cache.NamespaceGetter
		if t.bcache != nil {
			bcache = &cache.NamespaceGetter{Cache: t.bcache, NS: uint64(f.fd.Num)}
		}

		var tr *table.Reader
		tr, err = table.NewReader(r, f.size, f.fd, bcache, t.bpool, t.s.o.Options)
		if err != nil {
			r.Close()
			return 0, nil
		}
		return 1, tr

	})
	if ch == nil && err == nil {
		err = ErrClosed
	}
	return
}



func (t *tOps) find(f *tFile, key []byte, ro *opt.ReadOptions) (rkey, rvalue []byte, err error) {
	ch, err := t.open(f)
	if err != nil {
		return nil, nil, err
	}
	defer ch.Release()
	return ch.Value().(*table.Reader).Find(key, true, ro)
}


func (t *tOps) findKey(f *tFile, key []byte, ro *opt.ReadOptions) (rkey []byte, err error) {
	ch, err := t.open(f)
	if err != nil {
		return nil, err
	}
	defer ch.Release()
	return ch.Value().(*table.Reader).FindKey(key, true, ro)
}


func (t *tOps) offsetOf(f *tFile, key []byte) (offset int64, err error) {
	ch, err := t.open(f)
	if err != nil {
		return
	}
	defer ch.Release()
	return ch.Value().(*table.Reader).OffsetOf(key)
}


func (t *tOps) newIterator(f *tFile, slice *util.Range, ro *opt.ReadOptions) iterator.Iterator {
	ch, err := t.open(f)
	if err != nil {
		return iterator.NewEmptyIterator(err)
	}
	iter := ch.Value().(*table.Reader).NewIterator(slice, ro)
	iter.SetReleaser(ch)
	return iter
}



func (t *tOps) remove(f *tFile) {
	t.cache.Delete(0, uint64(f.fd.Num), func() {
		if err := t.s.stor.Remove(f.fd); err != nil {
			t.s.logf("table@remove removing @%d %q", f.fd.Num, err)
		} else {
			t.s.logf("table@remove removed @%d", f.fd.Num)
		}
		if t.bcache != nil {
			t.bcache.EvictNS(uint64(f.fd.Num))
		}
	})
}



func (t *tOps) close() {
	t.bpool.Close()
	t.cache.Close()
	if t.bcache != nil {
		t.bcache.CloseWeak()
	}
}


func newTableOps(s *session) *tOps {
	var (
		cacher cache.Cacher
		bcache *cache.Cache
		bpool  *util.BufferPool
	)
	if s.o.GetOpenFilesCacheCapacity() > 0 {
		cacher = cache.NewLRU(s.o.GetOpenFilesCacheCapacity())
	}
	if !s.o.GetDisableBlockCache() {
		var bcacher cache.Cacher
		if s.o.GetBlockCacheCapacity() > 0 {
			bcacher = s.o.GetBlockCacher().New(s.o.GetBlockCacheCapacity())
		}
		bcache = cache.NewCache(bcacher)
	}
	if !s.o.GetDisableBufferPool() {
		bpool = util.NewBufferPool(s.o.GetBlockSize() + 5)
	}
	return &tOps{
		s:      s,
		noSync: s.o.GetNoSync(),
		cache:  cache.NewCache(cacher),
		bcache: bcache,
		bpool:  bpool,
	}
}



type tWriter struct {
	t *tOps

	fd storage.FileDesc
	w  storage.Writer
	tw *table.Writer

	first, last []byte
}


func (w *tWriter) append(key, value []byte) error {
	if w.first == nil {
		w.first = append([]byte{}, key...)
	}
	w.last = append(w.last[:0], key...)
	return w.tw.Append(key, value)
}


func (w *tWriter) empty() bool {
	return w.first == nil
}


func (w *tWriter) close() {
	if w.w != nil {
		w.w.Close()
		w.w = nil
	}
}


func (w *tWriter) finish() (f *tFile, err error) {
	defer w.close()
	err = w.tw.Close()
	if err != nil {
		return
	}
	if !w.t.noSync {
		err = w.w.Sync()
		if err != nil {
			return
		}
	}
	f = newTableFile(w.fd, int64(w.tw.BytesLen()), internalKey(w.first), internalKey(w.last))
	return
}


func (w *tWriter) drop() {
	w.close()
	w.t.s.stor.Remove(w.fd)
	w.t.s.reuseFileNum(w.fd.Num)
	w.tw = nil
	w.first = nil
	w.last = nil
}
