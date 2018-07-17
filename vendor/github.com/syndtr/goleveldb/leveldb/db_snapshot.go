





package leveldb

import (
	"container/list"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type snapshotElement struct {
	seq uint64
	ref int
	e   *list.Element
}


func (db *DB) acquireSnapshot() *snapshotElement {
	db.snapsMu.Lock()
	defer db.snapsMu.Unlock()

	seq := db.getSeq()

	if e := db.snapsList.Back(); e != nil {
		se := e.Value.(*snapshotElement)
		if se.seq == seq {
			se.ref++
			return se
		} else if seq < se.seq {
			panic("leveldb: sequence number is not increasing")
		}
	}
	se := &snapshotElement{seq: seq, ref: 1}
	se.e = db.snapsList.PushBack(se)
	return se
}


func (db *DB) releaseSnapshot(se *snapshotElement) {
	db.snapsMu.Lock()
	defer db.snapsMu.Unlock()

	se.ref--
	if se.ref == 0 {
		db.snapsList.Remove(se.e)
		se.e = nil
	} else if se.ref < 0 {
		panic("leveldb: Snapshot: negative element reference")
	}
}


func (db *DB) minSeq() uint64 {
	db.snapsMu.Lock()
	defer db.snapsMu.Unlock()

	if e := db.snapsList.Front(); e != nil {
		return e.Value.(*snapshotElement).seq
	}

	return db.getSeq()
}


type Snapshot struct {
	db       *DB
	elem     *snapshotElement
	mu       sync.RWMutex
	released bool
}


func (db *DB) newSnapshot() *Snapshot {
	snap := &Snapshot{
		db:   db,
		elem: db.acquireSnapshot(),
	}
	atomic.AddInt32(&db.aliveSnaps, 1)
	runtime.SetFinalizer(snap, (*Snapshot).Release)
	return snap
}

func (snap *Snapshot) String() string {
	return fmt.Sprintf("leveldb.Snapshot{%d}", snap.elem.seq)
}






func (snap *Snapshot) Get(key []byte, ro *opt.ReadOptions) (value []byte, err error) {
	err = snap.db.ok()
	if err != nil {
		return
	}
	snap.mu.RLock()
	defer snap.mu.RUnlock()
	if snap.released {
		err = ErrSnapshotReleased
		return
	}
	return snap.db.get(nil, nil, key, snap.elem.seq, ro)
}




func (snap *Snapshot) Has(key []byte, ro *opt.ReadOptions) (ret bool, err error) {
	err = snap.db.ok()
	if err != nil {
		return
	}
	snap.mu.RLock()
	defer snap.mu.RUnlock()
	if snap.released {
		err = ErrSnapshotReleased
		return
	}
	return snap.db.has(nil, nil, key, snap.elem.seq, ro)
}


















func (snap *Snapshot) NewIterator(slice *util.Range, ro *opt.ReadOptions) iterator.Iterator {
	if err := snap.db.ok(); err != nil {
		return iterator.NewEmptyIterator(err)
	}
	snap.mu.Lock()
	defer snap.mu.Unlock()
	if snap.released {
		return iterator.NewEmptyIterator(ErrSnapshotReleased)
	}
	
	
	return snap.db.newIterator(nil, nil, snap.elem.seq, slice, ro)
}






func (snap *Snapshot) Release() {
	snap.mu.Lock()
	defer snap.mu.Unlock()

	if !snap.released {
		
		runtime.SetFinalizer(snap, nil)

		snap.released = true
		snap.db.releaseSnapshot(snap.elem)
		atomic.AddInt32(&snap.db.aliveSnaps, -1)
		snap.db = nil
		snap.elem = nil
	}
}
