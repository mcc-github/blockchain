





package leveldb

import (
	"errors"
	"sync/atomic"
	"time"

	"github.com/syndtr/goleveldb/leveldb/journal"
	"github.com/syndtr/goleveldb/leveldb/memdb"
	"github.com/syndtr/goleveldb/leveldb/storage"
)

var (
	errHasFrozenMem = errors.New("has frozen mem")
)

type memDB struct {
	db *DB
	*memdb.DB
	ref int32
}

func (m *memDB) getref() int32 {
	return atomic.LoadInt32(&m.ref)
}

func (m *memDB) incref() {
	atomic.AddInt32(&m.ref, 1)
}

func (m *memDB) decref() {
	if ref := atomic.AddInt32(&m.ref, -1); ref == 0 {
		
		if m.Capacity() == m.db.s.o.GetWriteBuffer() {
			m.Reset()
			m.db.mpoolPut(m.DB)
		}
		m.db = nil
		m.DB = nil
	} else if ref < 0 {
		panic("negative memdb ref")
	}
}


func (db *DB) getSeq() uint64 {
	return atomic.LoadUint64(&db.seq)
}


func (db *DB) addSeq(delta uint64) {
	atomic.AddUint64(&db.seq, delta)
}

func (db *DB) setSeq(seq uint64) {
	atomic.StoreUint64(&db.seq, seq)
}

func (db *DB) sampleSeek(ikey internalKey) {
	v := db.s.version()
	if v.sampleSeek(ikey) {
		
		db.compTrigger(db.tcompCmdC)
	}
	v.release()
}

func (db *DB) mpoolPut(mem *memdb.DB) {
	if !db.isClosed() {
		select {
		case db.memPool <- mem:
		default:
		}
	}
}

func (db *DB) mpoolGet(n int) *memDB {
	var mdb *memdb.DB
	select {
	case mdb = <-db.memPool:
	default:
	}
	if mdb == nil || mdb.Capacity() < n {
		mdb = memdb.New(db.s.icmp, maxInt(db.s.o.GetWriteBuffer(), n))
	}
	return &memDB{
		db: db,
		DB: mdb,
	}
}

func (db *DB) mpoolDrain() {
	ticker := time.NewTicker(30 * time.Second)
	for {
		select {
		case <-ticker.C:
			select {
			case <-db.memPool:
			default:
			}
		case <-db.closeC:
			ticker.Stop()
			
			select {
			case <-db.memPool:
			case <-time.After(time.Second):
			}
			close(db.memPool)
			return
		}
	}
}



func (db *DB) newMem(n int) (mem *memDB, err error) {
	fd := storage.FileDesc{Type: storage.TypeJournal, Num: db.s.allocFileNum()}
	w, err := db.s.stor.Create(fd)
	if err != nil {
		db.s.reuseFileNum(fd.Num)
		return
	}

	db.memMu.Lock()
	defer db.memMu.Unlock()

	if db.frozenMem != nil {
		return nil, errHasFrozenMem
	}

	if db.journal == nil {
		db.journal = journal.NewWriter(w)
	} else {
		db.journal.Reset(w)
		db.journalWriter.Close()
		db.frozenJournalFd = db.journalFd
	}
	db.journalWriter = w
	db.journalFd = fd
	db.frozenMem = db.mem
	mem = db.mpoolGet(n)
	mem.incref() 
	mem.incref() 
	db.mem = mem
	
	
	db.frozenSeq = db.seq
	return
}


func (db *DB) getMems() (e, f *memDB) {
	db.memMu.RLock()
	defer db.memMu.RUnlock()
	if db.mem != nil {
		db.mem.incref()
	} else if !db.isClosed() {
		panic("nil effective mem")
	}
	if db.frozenMem != nil {
		db.frozenMem.incref()
	}
	return db.mem, db.frozenMem
}


func (db *DB) getEffectiveMem() *memDB {
	db.memMu.RLock()
	defer db.memMu.RUnlock()
	if db.mem != nil {
		db.mem.incref()
	} else if !db.isClosed() {
		panic("nil effective mem")
	}
	return db.mem
}


func (db *DB) hasFrozenMem() bool {
	db.memMu.RLock()
	defer db.memMu.RUnlock()
	return db.frozenMem != nil
}


func (db *DB) getFrozenMem() *memDB {
	db.memMu.RLock()
	defer db.memMu.RUnlock()
	if db.frozenMem != nil {
		db.frozenMem.incref()
	}
	return db.frozenMem
}


func (db *DB) dropFrozenMem() {
	db.memMu.Lock()
	if err := db.s.stor.Remove(db.frozenJournalFd); err != nil {
		db.logf("journal@remove removing @%d %q", db.frozenJournalFd.Num, err)
	} else {
		db.logf("journal@remove removed @%d", db.frozenJournalFd.Num)
	}
	db.frozenJournalFd = storage.FileDesc{}
	db.frozenMem.decref()
	db.frozenMem = nil
	db.memMu.Unlock()
}


func (db *DB) clearMems() {
	db.memMu.Lock()
	db.mem = nil
	db.frozenMem = nil
	db.memMu.Unlock()
}


func (db *DB) setClosed() bool {
	return atomic.CompareAndSwapUint32(&db.closed, 0, 1)
}


func (db *DB) isClosed() bool {
	return atomic.LoadUint32(&db.closed) != 0
}


func (db *DB) ok() error {
	if db.isClosed() {
		return ErrClosed
	}
	return nil
}
