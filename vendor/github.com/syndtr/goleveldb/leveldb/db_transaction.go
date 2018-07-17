





package leveldb

import (
	"errors"
	"sync"
	"time"

	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

var errTransactionDone = errors.New("leveldb: transaction already closed")


type Transaction struct {
	db        *DB
	lk        sync.RWMutex
	seq       uint64
	mem       *memDB
	tables    tFiles
	ikScratch []byte
	rec       sessionRecord
	stats     cStatStaging
	closed    bool
}







func (tr *Transaction) Get(key []byte, ro *opt.ReadOptions) ([]byte, error) {
	tr.lk.RLock()
	defer tr.lk.RUnlock()
	if tr.closed {
		return nil, errTransactionDone
	}
	return tr.db.get(tr.mem.DB, tr.tables, key, tr.seq, ro)
}




func (tr *Transaction) Has(key []byte, ro *opt.ReadOptions) (bool, error) {
	tr.lk.RLock()
	defer tr.lk.RUnlock()
	if tr.closed {
		return false, errTransactionDone
	}
	return tr.db.has(tr.mem.DB, tr.tables, key, tr.seq, ro)
}















func (tr *Transaction) NewIterator(slice *util.Range, ro *opt.ReadOptions) iterator.Iterator {
	tr.lk.RLock()
	defer tr.lk.RUnlock()
	if tr.closed {
		return iterator.NewEmptyIterator(errTransactionDone)
	}
	tr.mem.incref()
	return tr.db.newIterator(tr.mem, tr.tables, tr.seq, slice, ro)
}

func (tr *Transaction) flush() error {
	
	if tr.mem.Len() != 0 {
		tr.stats.startTimer()
		iter := tr.mem.NewIterator(nil)
		t, n, err := tr.db.s.tops.createFrom(iter)
		iter.Release()
		tr.stats.stopTimer()
		if err != nil {
			return err
		}
		if tr.mem.getref() == 1 {
			tr.mem.Reset()
		} else {
			tr.mem.decref()
			tr.mem = tr.db.mpoolGet(0)
			tr.mem.incref()
		}
		tr.tables = append(tr.tables, t)
		tr.rec.addTableFile(0, t)
		tr.stats.write += t.size
		tr.db.logf("transaction@flush created L0@%d N·%d S·%s %q:%q", t.fd.Num, n, shortenb(int(t.size)), t.imin, t.imax)
	}
	return nil
}

func (tr *Transaction) put(kt keyType, key, value []byte) error {
	tr.ikScratch = makeInternalKey(tr.ikScratch, key, tr.seq+1, kt)
	if tr.mem.Free() < len(tr.ikScratch)+len(value) {
		if err := tr.flush(); err != nil {
			return err
		}
	}
	if err := tr.mem.Put(tr.ikScratch, value); err != nil {
		return err
	}
	tr.seq++
	return nil
}







func (tr *Transaction) Put(key, value []byte, wo *opt.WriteOptions) error {
	tr.lk.Lock()
	defer tr.lk.Unlock()
	if tr.closed {
		return errTransactionDone
	}
	return tr.put(keyTypeVal, key, value)
}






func (tr *Transaction) Delete(key []byte, wo *opt.WriteOptions) error {
	tr.lk.Lock()
	defer tr.lk.Unlock()
	if tr.closed {
		return errTransactionDone
	}
	return tr.put(keyTypeDel, key, nil)
}







func (tr *Transaction) Write(b *Batch, wo *opt.WriteOptions) error {
	if b == nil || b.Len() == 0 {
		return nil
	}

	tr.lk.Lock()
	defer tr.lk.Unlock()
	if tr.closed {
		return errTransactionDone
	}
	return b.replayInternal(func(i int, kt keyType, k, v []byte) error {
		return tr.put(kt, k, v)
	})
}

func (tr *Transaction) setDone() {
	tr.closed = true
	tr.db.tr = nil
	tr.mem.decref()
	<-tr.db.writeLockC
}





func (tr *Transaction) Commit() error {
	if err := tr.db.ok(); err != nil {
		return err
	}

	tr.lk.Lock()
	defer tr.lk.Unlock()
	if tr.closed {
		return errTransactionDone
	}
	if err := tr.flush(); err != nil {
		
		
		return err
	}
	if len(tr.tables) != 0 {
		
		tr.rec.setSeqNum(tr.seq)
		tr.db.compCommitLk.Lock()
		tr.stats.startTimer()
		var cerr error
		for retry := 0; retry < 3; retry++ {
			cerr = tr.db.s.commit(&tr.rec)
			if cerr != nil {
				tr.db.logf("transaction@commit error R·%d %q", retry, cerr)
				select {
				case <-time.After(time.Second):
				case <-tr.db.closeC:
					tr.db.logf("transaction@commit exiting")
					tr.db.compCommitLk.Unlock()
					return cerr
				}
			} else {
				
				tr.db.setSeq(tr.seq)
				break
			}
		}
		tr.stats.stopTimer()
		if cerr != nil {
			
			
			return cerr
		}

		
		tr.db.compStats.addStat(0, &tr.stats)

		
		tr.db.compTrigger(tr.db.tcompCmdC)
		tr.db.compCommitLk.Unlock()

		
		
		tr.db.waitCompaction()
	}
	
	tr.setDone()
	return nil
}

func (tr *Transaction) discard() {
	
	for _, t := range tr.tables {
		tr.db.logf("transaction@discard @%d", t.fd.Num)
		if err1 := tr.db.s.stor.Remove(t.fd); err1 == nil {
			tr.db.s.reuseFileNum(t.fd.Num)
		}
	}
}




func (tr *Transaction) Discard() {
	tr.lk.Lock()
	if !tr.closed {
		tr.discard()
		tr.setDone()
	}
	tr.lk.Unlock()
}

func (db *DB) waitCompaction() error {
	if db.s.tLen(0) >= db.s.o.GetWriteL0PauseTrigger() {
		return db.compTriggerWait(db.tcompCmdC)
	}
	return nil
}












func (db *DB) OpenTransaction() (*Transaction, error) {
	if err := db.ok(); err != nil {
		return nil, err
	}

	
	select {
	case db.writeLockC <- struct{}{}:
	case err := <-db.compPerErrC:
		return nil, err
	case <-db.closeC:
		return nil, ErrClosed
	}

	if db.tr != nil {
		panic("leveldb: has open transaction")
	}

	
	if db.mem != nil && db.mem.Len() != 0 {
		if _, err := db.rotateMem(0, true); err != nil {
			return nil, err
		}
	}

	
	if err := db.waitCompaction(); err != nil {
		return nil, err
	}

	tr := &Transaction{
		db:  db,
		seq: db.seq,
		mem: db.mpoolGet(0),
	}
	tr.mem.incref()
	db.tr = tr
	return tr, nil
}
