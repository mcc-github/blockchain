





package leveldb

import (
	"container/list"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/journal"
	"github.com/syndtr/goleveldb/leveldb/memdb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/storage"
	"github.com/syndtr/goleveldb/leveldb/table"
	"github.com/syndtr/goleveldb/leveldb/util"
)


type DB struct {
	
	seq uint64

	
	cWriteDelay            int64 
	cWriteDelayN           int32 
	inWritePaused          int32 
	aliveSnaps, aliveIters int32

	
	s *session

	
	memMu           sync.RWMutex
	memPool         chan *memdb.DB
	mem, frozenMem  *memDB
	journal         *journal.Writer
	journalWriter   storage.Writer
	journalFd       storage.FileDesc
	frozenJournalFd storage.FileDesc
	frozenSeq       uint64

	
	snapsMu   sync.Mutex
	snapsList *list.List

	
	batchPool    sync.Pool
	writeMergeC  chan writeMerge
	writeMergedC chan bool
	writeLockC   chan struct{}
	writeAckC    chan error
	writeDelay   time.Duration
	writeDelayN  int
	tr           *Transaction

	
	compCommitLk     sync.Mutex
	tcompCmdC        chan cCmd
	tcompPauseC      chan chan<- struct{}
	mcompCmdC        chan cCmd
	compErrC         chan error
	compPerErrC      chan error
	compErrSetC      chan error
	compWriteLocking bool
	compStats        cStats
	memdbMaxLevel    int 

	
	closeW sync.WaitGroup
	closeC chan struct{}
	closed uint32
	closer io.Closer
}

func openDB(s *session) (*DB, error) {
	s.log("db@open opening")
	start := time.Now()
	db := &DB{
		s: s,
		
		seq: s.stSeqNum,
		
		memPool: make(chan *memdb.DB, 1),
		
		snapsList: list.New(),
		
		batchPool:    sync.Pool{New: newBatch},
		writeMergeC:  make(chan writeMerge),
		writeMergedC: make(chan bool),
		writeLockC:   make(chan struct{}, 1),
		writeAckC:    make(chan error),
		
		tcompCmdC:   make(chan cCmd),
		tcompPauseC: make(chan chan<- struct{}),
		mcompCmdC:   make(chan cCmd),
		compErrC:    make(chan error),
		compPerErrC: make(chan error),
		compErrSetC: make(chan error),
		
		closeC: make(chan struct{}),
	}

	
	readOnly := s.o.GetReadOnly()

	if readOnly {
		
		if err := db.recoverJournalRO(); err != nil {
			return nil, err
		}
	} else {
		
		if err := db.recoverJournal(); err != nil {
			return nil, err
		}

		
		if err := db.checkAndCleanFiles(); err != nil {
			
			if db.journal != nil {
				db.journal.Close()
				db.journalWriter.Close()
			}
			return nil, err
		}

	}

	
	go db.compactionError()
	go db.mpoolDrain()

	if readOnly {
		db.SetReadOnly()
	} else {
		db.closeW.Add(2)
		go db.tCompaction()
		go db.mCompaction()
		
	}

	s.logf("db@open done T·%v", time.Since(start))

	runtime.SetFinalizer(db, (*DB).Close)
	return db, nil
}












func Open(stor storage.Storage, o *opt.Options) (db *DB, err error) {
	s, err := newSession(stor, o)
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			s.close()
			s.release()
		}
	}()

	err = s.recover()
	if err != nil {
		if !os.IsNotExist(err) || s.o.GetErrorIfMissing() {
			return
		}
		err = s.create()
		if err != nil {
			return
		}
	} else if s.o.GetErrorIfExist() {
		err = os.ErrExist
		return
	}

	return openDB(s)
}















func OpenFile(path string, o *opt.Options) (db *DB, err error) {
	stor, err := storage.OpenFile(path, o.GetReadOnly())
	if err != nil {
		return
	}
	db, err = Open(stor, o)
	if err != nil {
		stor.Close()
	} else {
		db.closer = stor
	}
	return
}








func Recover(stor storage.Storage, o *opt.Options) (db *DB, err error) {
	s, err := newSession(stor, o)
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			s.close()
			s.release()
		}
	}()

	err = recoverTable(s, o)
	if err != nil {
		return
	}
	return openDB(s)
}











func RecoverFile(path string, o *opt.Options) (db *DB, err error) {
	stor, err := storage.OpenFile(path, false)
	if err != nil {
		return
	}
	db, err = Recover(stor, o)
	if err != nil {
		stor.Close()
	} else {
		db.closer = stor
	}
	return
}

func recoverTable(s *session, o *opt.Options) error {
	o = dupOptions(o)
	
	o.Strict &= ^opt.StrictReader

	
	fds, err := s.stor.List(storage.TypeTable)
	if err != nil {
		return err
	}
	sortFds(fds)

	var (
		maxSeq                                                            uint64
		recoveredKey, goodKey, corruptedKey, corruptedBlock, droppedTable int

		
		strict = o.GetStrict(opt.StrictRecovery)
		noSync = o.GetNoSync()

		rec   = &sessionRecord{}
		bpool = util.NewBufferPool(o.GetBlockSize() + 5)
	)
	buildTable := func(iter iterator.Iterator) (tmpFd storage.FileDesc, size int64, err error) {
		tmpFd = s.newTemp()
		writer, err := s.stor.Create(tmpFd)
		if err != nil {
			return
		}
		defer func() {
			writer.Close()
			if err != nil {
				s.stor.Remove(tmpFd)
				tmpFd = storage.FileDesc{}
			}
		}()

		
		tw := table.NewWriter(writer, o)
		for iter.Next() {
			key := iter.Key()
			if validInternalKey(key) {
				err = tw.Append(key, iter.Value())
				if err != nil {
					return
				}
			}
		}
		err = iter.Error()
		if err != nil && !errors.IsCorrupted(err) {
			return
		}
		err = tw.Close()
		if err != nil {
			return
		}
		if !noSync {
			err = writer.Sync()
			if err != nil {
				return
			}
		}
		size = int64(tw.BytesLen())
		return
	}
	recoverTable := func(fd storage.FileDesc) error {
		s.logf("table@recovery recovering @%d", fd.Num)
		reader, err := s.stor.Open(fd)
		if err != nil {
			return err
		}
		var closed bool
		defer func() {
			if !closed {
				reader.Close()
			}
		}()

		
		size, err := reader.Seek(0, 2)
		if err != nil {
			return err
		}

		var (
			tSeq                                     uint64
			tgoodKey, tcorruptedKey, tcorruptedBlock int
			imin, imax                               []byte
		)
		tr, err := table.NewReader(reader, size, fd, nil, bpool, o)
		if err != nil {
			return err
		}
		iter := tr.NewIterator(nil, nil)
		if itererr, ok := iter.(iterator.ErrorCallbackSetter); ok {
			itererr.SetErrorCallback(func(err error) {
				if errors.IsCorrupted(err) {
					s.logf("table@recovery block corruption @%d %q", fd.Num, err)
					tcorruptedBlock++
				}
			})
		}

		
		for iter.Next() {
			key := iter.Key()
			_, seq, _, kerr := parseInternalKey(key)
			if kerr != nil {
				tcorruptedKey++
				continue
			}
			tgoodKey++
			if seq > tSeq {
				tSeq = seq
			}
			if imin == nil {
				imin = append([]byte{}, key...)
			}
			imax = append(imax[:0], key...)
		}
		if err := iter.Error(); err != nil && !errors.IsCorrupted(err) {
			iter.Release()
			return err
		}
		iter.Release()

		goodKey += tgoodKey
		corruptedKey += tcorruptedKey
		corruptedBlock += tcorruptedBlock

		if strict && (tcorruptedKey > 0 || tcorruptedBlock > 0) {
			droppedTable++
			s.logf("table@recovery dropped @%d Gk·%d Ck·%d Cb·%d S·%d Q·%d", fd.Num, tgoodKey, tcorruptedKey, tcorruptedBlock, size, tSeq)
			return nil
		}

		if tgoodKey > 0 {
			if tcorruptedKey > 0 || tcorruptedBlock > 0 {
				
				s.logf("table@recovery rebuilding @%d", fd.Num)
				iter := tr.NewIterator(nil, nil)
				tmpFd, newSize, err := buildTable(iter)
				iter.Release()
				if err != nil {
					return err
				}
				closed = true
				reader.Close()
				if err := s.stor.Rename(tmpFd, fd); err != nil {
					return err
				}
				size = newSize
			}
			if tSeq > maxSeq {
				maxSeq = tSeq
			}
			recoveredKey += tgoodKey
			
			rec.addTable(0, fd.Num, size, imin, imax)
			s.logf("table@recovery recovered @%d Gk·%d Ck·%d Cb·%d S·%d Q·%d", fd.Num, tgoodKey, tcorruptedKey, tcorruptedBlock, size, tSeq)
		} else {
			droppedTable++
			s.logf("table@recovery unrecoverable @%d Ck·%d Cb·%d S·%d", fd.Num, tcorruptedKey, tcorruptedBlock, size)
		}

		return nil
	}

	
	if len(fds) > 0 {
		s.logf("table@recovery F·%d", len(fds))

		
		s.markFileNum(fds[len(fds)-1].Num)

		for _, fd := range fds {
			if err := recoverTable(fd); err != nil {
				return err
			}
		}

		s.logf("table@recovery recovered F·%d N·%d Gk·%d Ck·%d Q·%d", len(fds), recoveredKey, goodKey, corruptedKey, maxSeq)
	}

	
	rec.setSeqNum(maxSeq)

	
	if err := s.create(); err != nil {
		return err
	}

	
	return s.commit(rec)
}

func (db *DB) recoverJournal() error {
	
	rawFds, err := db.s.stor.List(storage.TypeJournal)
	if err != nil {
		return err
	}
	sortFds(rawFds)

	
	var fds []storage.FileDesc
	for _, fd := range rawFds {
		if fd.Num >= db.s.stJournalNum || fd.Num == db.s.stPrevJournalNum {
			fds = append(fds, fd)
		}
	}

	var (
		ofd storage.FileDesc 
		rec = &sessionRecord{}
	)

	
	if len(fds) > 0 {
		db.logf("journal@recovery F·%d", len(fds))

		
		db.s.markFileNum(fds[len(fds)-1].Num)

		var (
			
			strict      = db.s.o.GetStrict(opt.StrictJournal)
			checksum    = db.s.o.GetStrict(opt.StrictJournalChecksum)
			writeBuffer = db.s.o.GetWriteBuffer()

			jr       *journal.Reader
			mdb      = memdb.New(db.s.icmp, writeBuffer)
			buf      = &util.Buffer{}
			batchSeq uint64
			batchLen int
		)

		for _, fd := range fds {
			db.logf("journal@recovery recovering @%d", fd.Num)

			fr, err := db.s.stor.Open(fd)
			if err != nil {
				return err
			}

			
			if jr == nil {
				jr = journal.NewReader(fr, dropper{db.s, fd}, strict, checksum)
			} else {
				jr.Reset(fr, dropper{db.s, fd}, strict, checksum)
			}

			
			if !ofd.Zero() {
				if mdb.Len() > 0 {
					if _, err := db.s.flushMemdb(rec, mdb, 0); err != nil {
						fr.Close()
						return err
					}
				}

				rec.setJournalNum(fd.Num)
				rec.setSeqNum(db.seq)
				if err := db.s.commit(rec); err != nil {
					fr.Close()
					return err
				}
				rec.resetAddedTables()

				db.s.stor.Remove(ofd)
				ofd = storage.FileDesc{}
			}

			
			mdb.Reset()
			for {
				r, err := jr.Next()
				if err != nil {
					if err == io.EOF {
						break
					}

					fr.Close()
					return errors.SetFd(err, fd)
				}

				buf.Reset()
				if _, err := buf.ReadFrom(r); err != nil {
					if err == io.ErrUnexpectedEOF {
						
						continue
					}

					fr.Close()
					return errors.SetFd(err, fd)
				}
				batchSeq, batchLen, err = decodeBatchToMem(buf.Bytes(), db.seq, mdb)
				if err != nil {
					if !strict && errors.IsCorrupted(err) {
						db.s.logf("journal error: %v (skipped)", err)
						
						continue
					}

					fr.Close()
					return errors.SetFd(err, fd)
				}

				
				db.seq = batchSeq + uint64(batchLen)

				
				if mdb.Size() >= writeBuffer {
					if _, err := db.s.flushMemdb(rec, mdb, 0); err != nil {
						fr.Close()
						return err
					}

					mdb.Reset()
				}
			}

			fr.Close()
			ofd = fd
		}

		
		if mdb.Len() > 0 {
			if _, err := db.s.flushMemdb(rec, mdb, 0); err != nil {
				return err
			}
		}
	}

	
	if _, err := db.newMem(0); err != nil {
		return err
	}

	
	rec.setJournalNum(db.journalFd.Num)
	rec.setSeqNum(db.seq)
	if err := db.s.commit(rec); err != nil {
		
		if db.journal != nil {
			db.journal.Close()
			db.journalWriter.Close()
		}
		return err
	}

	
	if !ofd.Zero() {
		db.s.stor.Remove(ofd)
	}

	return nil
}

func (db *DB) recoverJournalRO() error {
	
	rawFds, err := db.s.stor.List(storage.TypeJournal)
	if err != nil {
		return err
	}
	sortFds(rawFds)

	
	var fds []storage.FileDesc
	for _, fd := range rawFds {
		if fd.Num >= db.s.stJournalNum || fd.Num == db.s.stPrevJournalNum {
			fds = append(fds, fd)
		}
	}

	var (
		
		strict      = db.s.o.GetStrict(opt.StrictJournal)
		checksum    = db.s.o.GetStrict(opt.StrictJournalChecksum)
		writeBuffer = db.s.o.GetWriteBuffer()

		mdb = memdb.New(db.s.icmp, writeBuffer)
	)

	
	if len(fds) > 0 {
		db.logf("journal@recovery RO·Mode F·%d", len(fds))

		var (
			jr       *journal.Reader
			buf      = &util.Buffer{}
			batchSeq uint64
			batchLen int
		)

		for _, fd := range fds {
			db.logf("journal@recovery recovering @%d", fd.Num)

			fr, err := db.s.stor.Open(fd)
			if err != nil {
				return err
			}

			
			if jr == nil {
				jr = journal.NewReader(fr, dropper{db.s, fd}, strict, checksum)
			} else {
				jr.Reset(fr, dropper{db.s, fd}, strict, checksum)
			}

			
			for {
				r, err := jr.Next()
				if err != nil {
					if err == io.EOF {
						break
					}

					fr.Close()
					return errors.SetFd(err, fd)
				}

				buf.Reset()
				if _, err := buf.ReadFrom(r); err != nil {
					if err == io.ErrUnexpectedEOF {
						
						continue
					}

					fr.Close()
					return errors.SetFd(err, fd)
				}
				batchSeq, batchLen, err = decodeBatchToMem(buf.Bytes(), db.seq, mdb)
				if err != nil {
					if !strict && errors.IsCorrupted(err) {
						db.s.logf("journal error: %v (skipped)", err)
						
						continue
					}

					fr.Close()
					return errors.SetFd(err, fd)
				}

				
				db.seq = batchSeq + uint64(batchLen)
			}

			fr.Close()
		}
	}

	
	db.mem = &memDB{db: db, DB: mdb, ref: 1}

	return nil
}

func memGet(mdb *memdb.DB, ikey internalKey, icmp *iComparer) (ok bool, mv []byte, err error) {
	mk, mv, err := mdb.Find(ikey)
	if err == nil {
		ukey, _, kt, kerr := parseInternalKey(mk)
		if kerr != nil {
			
			panic(kerr)
		}
		if icmp.uCompare(ukey, ikey.ukey()) == 0 {
			if kt == keyTypeDel {
				return true, nil, ErrNotFound
			}
			return true, mv, nil

		}
	} else if err != ErrNotFound {
		return true, nil, err
	}
	return
}

func (db *DB) get(auxm *memdb.DB, auxt tFiles, key []byte, seq uint64, ro *opt.ReadOptions) (value []byte, err error) {
	ikey := makeInternalKey(nil, key, seq, keyTypeSeek)

	if auxm != nil {
		if ok, mv, me := memGet(auxm, ikey, db.s.icmp); ok {
			return append([]byte{}, mv...), me
		}
	}

	em, fm := db.getMems()
	for _, m := range [...]*memDB{em, fm} {
		if m == nil {
			continue
		}
		defer m.decref()

		if ok, mv, me := memGet(m.DB, ikey, db.s.icmp); ok {
			return append([]byte{}, mv...), me
		}
	}

	v := db.s.version()
	value, cSched, err := v.get(auxt, ikey, ro, false)
	v.release()
	if cSched {
		
		db.compTrigger(db.tcompCmdC)
	}
	return
}

func nilIfNotFound(err error) error {
	if err == ErrNotFound {
		return nil
	}
	return err
}

func (db *DB) has(auxm *memdb.DB, auxt tFiles, key []byte, seq uint64, ro *opt.ReadOptions) (ret bool, err error) {
	ikey := makeInternalKey(nil, key, seq, keyTypeSeek)

	if auxm != nil {
		if ok, _, me := memGet(auxm, ikey, db.s.icmp); ok {
			return me == nil, nilIfNotFound(me)
		}
	}

	em, fm := db.getMems()
	for _, m := range [...]*memDB{em, fm} {
		if m == nil {
			continue
		}
		defer m.decref()

		if ok, _, me := memGet(m.DB, ikey, db.s.icmp); ok {
			return me == nil, nilIfNotFound(me)
		}
	}

	v := db.s.version()
	_, cSched, err := v.get(auxt, ikey, ro, true)
	v.release()
	if cSched {
		
		db.compTrigger(db.tcompCmdC)
	}
	if err == nil {
		ret = true
	} else if err == ErrNotFound {
		err = nil
	}
	return
}







func (db *DB) Get(key []byte, ro *opt.ReadOptions) (value []byte, err error) {
	err = db.ok()
	if err != nil {
		return
	}

	se := db.acquireSnapshot()
	defer db.releaseSnapshot(se)
	return db.get(nil, nil, key, se.seq, ro)
}




func (db *DB) Has(key []byte, ro *opt.ReadOptions) (ret bool, err error) {
	err = db.ok()
	if err != nil {
		return
	}

	se := db.acquireSnapshot()
	defer db.releaseSnapshot(se)
	return db.has(nil, nil, key, se.seq, ro)
}

















func (db *DB) NewIterator(slice *util.Range, ro *opt.ReadOptions) iterator.Iterator {
	if err := db.ok(); err != nil {
		return iterator.NewEmptyIterator(err)
	}

	se := db.acquireSnapshot()
	defer db.releaseSnapshot(se)
	
	
	return db.newIterator(nil, nil, se.seq, slice, ro)
}






func (db *DB) GetSnapshot() (*Snapshot, error) {
	if err := db.ok(); err != nil {
		return nil, err
	}

	return db.newSnapshot(), nil
}
























func (db *DB) GetProperty(name string) (value string, err error) {
	err = db.ok()
	if err != nil {
		return
	}

	const prefix = "leveldb."
	if !strings.HasPrefix(name, prefix) {
		return "", ErrNotFound
	}
	p := name[len(prefix):]

	v := db.s.version()
	defer v.release()

	numFilesPrefix := "num-files-at-level"
	switch {
	case strings.HasPrefix(p, numFilesPrefix):
		var level uint
		var rest string
		n, _ := fmt.Sscanf(p[len(numFilesPrefix):], "%d%s", &level, &rest)
		if n != 1 {
			err = ErrNotFound
		} else {
			value = fmt.Sprint(v.tLen(int(level)))
		}
	case p == "stats":
		value = "Compactions\n" +
			" Level |   Tables   |    Size(MB)   |    Time(sec)  |    Read(MB)   |   Write(MB)\n" +
			"-------+------------+---------------+---------------+---------------+---------------\n"
		for level, tables := range v.levels {
			duration, read, write := db.compStats.getStat(level)
			if len(tables) == 0 && duration == 0 {
				continue
			}
			value += fmt.Sprintf(" %3d   | %10d | %13.5f | %13.5f | %13.5f | %13.5f\n",
				level, len(tables), float64(tables.size())/1048576.0, duration.Seconds(),
				float64(read)/1048576.0, float64(write)/1048576.0)
		}
	case p == "iostats":
		value = fmt.Sprintf("Read(MB):%.5f Write(MB):%.5f",
			float64(db.s.stor.reads())/1048576.0,
			float64(db.s.stor.writes())/1048576.0)
	case p == "writedelay":
		writeDelayN, writeDelay := atomic.LoadInt32(&db.cWriteDelayN), time.Duration(atomic.LoadInt64(&db.cWriteDelay))
		paused := atomic.LoadInt32(&db.inWritePaused) == 1
		value = fmt.Sprintf("DelayN:%d Delay:%s Paused:%t", writeDelayN, writeDelay, paused)
	case p == "sstables":
		for level, tables := range v.levels {
			value += fmt.Sprintf("--- level %d ---\n", level)
			for _, t := range tables {
				value += fmt.Sprintf("%d:%d[%q .. %q]\n", t.fd.Num, t.size, t.imin, t.imax)
			}
		}
	case p == "blockpool":
		value = fmt.Sprintf("%v", db.s.tops.bpool)
	case p == "cachedblock":
		if db.s.tops.bcache != nil {
			value = fmt.Sprintf("%d", db.s.tops.bcache.Size())
		} else {
			value = "<nil>"
		}
	case p == "openedtables":
		value = fmt.Sprintf("%d", db.s.tops.cache.Size())
	case p == "alivesnaps":
		value = fmt.Sprintf("%d", atomic.LoadInt32(&db.aliveSnaps))
	case p == "aliveiters":
		value = fmt.Sprintf("%d", atomic.LoadInt32(&db.aliveIters))
	default:
		err = ErrNotFound
	}

	return
}


type DBStats struct {
	WriteDelayCount    int32
	WriteDelayDuration time.Duration
	WritePaused        bool

	AliveSnapshots int32
	AliveIterators int32

	IOWrite uint64
	IORead  uint64

	BlockCacheSize    int
	OpenedTablesCount int

	LevelSizes        []int64
	LevelTablesCounts []int
	LevelRead         []int64
	LevelWrite        []int64
	LevelDurations    []time.Duration
}


func (db *DB) Stats(s *DBStats) error {
	err := db.ok()
	if err != nil {
		return err
	}

	s.IORead = db.s.stor.reads()
	s.IOWrite = db.s.stor.writes()
	s.WriteDelayCount = atomic.LoadInt32(&db.cWriteDelayN)
	s.WriteDelayDuration = time.Duration(atomic.LoadInt64(&db.cWriteDelay))
	s.WritePaused = atomic.LoadInt32(&db.inWritePaused) == 1

	s.OpenedTablesCount = db.s.tops.cache.Size()
	if db.s.tops.bcache != nil {
		s.BlockCacheSize = db.s.tops.bcache.Size()
	} else {
		s.BlockCacheSize = 0
	}

	s.AliveIterators = atomic.LoadInt32(&db.aliveIters)
	s.AliveSnapshots = atomic.LoadInt32(&db.aliveSnaps)

	s.LevelDurations = s.LevelDurations[:0]
	s.LevelRead = s.LevelRead[:0]
	s.LevelWrite = s.LevelWrite[:0]
	s.LevelSizes = s.LevelSizes[:0]
	s.LevelTablesCounts = s.LevelTablesCounts[:0]

	v := db.s.version()
	defer v.release()

	for level, tables := range v.levels {
		duration, read, write := db.compStats.getStat(level)
		if len(tables) == 0 && duration == 0 {
			continue
		}
		s.LevelDurations = append(s.LevelDurations, duration)
		s.LevelRead = append(s.LevelRead, read)
		s.LevelWrite = append(s.LevelWrite, write)
		s.LevelSizes = append(s.LevelSizes, tables.size())
		s.LevelTablesCounts = append(s.LevelTablesCounts, len(tables))
	}

	return nil
}







func (db *DB) SizeOf(ranges []util.Range) (Sizes, error) {
	if err := db.ok(); err != nil {
		return nil, err
	}

	v := db.s.version()
	defer v.release()

	sizes := make(Sizes, 0, len(ranges))
	for _, r := range ranges {
		imin := makeInternalKey(nil, r.Start, keyMaxSeq, keyTypeSeek)
		imax := makeInternalKey(nil, r.Limit, keyMaxSeq, keyTypeSeek)
		start, err := v.offsetOf(imin)
		if err != nil {
			return nil, err
		}
		limit, err := v.offsetOf(imax)
		if err != nil {
			return nil, err
		}
		var size int64
		if limit >= start {
			size = limit - start
		}
		sizes = append(sizes, size)
	}

	return sizes, nil
}







func (db *DB) Close() error {
	if !db.setClosed() {
		return ErrClosed
	}

	start := time.Now()
	db.log("db@close closing")

	
	runtime.SetFinalizer(db, nil)

	
	var err error
	select {
	case err = <-db.compErrC:
		if err == ErrReadOnly {
			err = nil
		}
	default:
	}

	
	close(db.closeC)

	
	if db.tr != nil {
		db.tr.Discard()
	}

	
	db.writeLockC <- struct{}{}

	
	db.closeW.Wait()

	
	if db.journal != nil {
		db.journal.Close()
		db.journalWriter.Close()
		db.journal = nil
		db.journalWriter = nil
	}

	if db.writeDelayN > 0 {
		db.logf("db@write was delayed N·%d T·%v", db.writeDelayN, db.writeDelay)
	}

	
	db.s.close()
	db.logf("db@close done T·%v", time.Since(start))
	db.s.release()

	if db.closer != nil {
		if err1 := db.closer.Close(); err == nil {
			err = err1
		}
		db.closer = nil
	}

	
	db.clearMems()

	return err
}
