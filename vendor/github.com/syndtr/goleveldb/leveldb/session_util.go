





package leveldb

import (
	"fmt"
	"sync/atomic"

	"github.com/syndtr/goleveldb/leveldb/journal"
	"github.com/syndtr/goleveldb/leveldb/storage"
)



type dropper struct {
	s  *session
	fd storage.FileDesc
}

func (d dropper) Drop(err error) {
	if e, ok := err.(*journal.ErrCorrupted); ok {
		d.s.logf("journal@drop %s-%d S·%s %q", d.fd.Type, d.fd.Num, shortenb(e.Size), e.Reason)
	} else {
		d.s.logf("journal@drop %s-%d %q", d.fd.Type, d.fd.Num, err)
	}
}

func (s *session) log(v ...interface{})                 { s.stor.Log(fmt.Sprint(v...)) }
func (s *session) logf(format string, v ...interface{}) { s.stor.Log(fmt.Sprintf(format, v...)) }



func (s *session) newTemp() storage.FileDesc {
	num := atomic.AddInt64(&s.stTempFileNum, 1) - 1
	return storage.FileDesc{storage.TypeTemp, num}
}

func (s *session) addFileRef(fd storage.FileDesc, ref int) int {
	ref += s.fileRef[fd.Num]
	if ref > 0 {
		s.fileRef[fd.Num] = ref
	} else if ref == 0 {
		delete(s.fileRef, fd.Num)
	} else {
		panic(fmt.Sprintf("negative ref: %v", fd))
	}
	return ref
}





func (s *session) version() *version {
	s.vmu.Lock()
	defer s.vmu.Unlock()
	s.stVersion.incref()
	return s.stVersion
}

func (s *session) tLen(level int) int {
	s.vmu.Lock()
	defer s.vmu.Unlock()
	return s.stVersion.tLen(level)
}


func (s *session) setVersion(v *version) {
	s.vmu.Lock()
	defer s.vmu.Unlock()
	
	
	v.incref()
	if s.stVersion != nil {
		
		s.stVersion.releaseNB()
	}
	s.stVersion = v
}


func (s *session) nextFileNum() int64 {
	return atomic.LoadInt64(&s.stNextFileNum)
}


func (s *session) setNextFileNum(num int64) {
	atomic.StoreInt64(&s.stNextFileNum, num)
}


func (s *session) markFileNum(num int64) {
	nextFileNum := num + 1
	for {
		old, x := s.stNextFileNum, nextFileNum
		if old > x {
			x = old
		}
		if atomic.CompareAndSwapInt64(&s.stNextFileNum, old, x) {
			break
		}
	}
}


func (s *session) allocFileNum() int64 {
	return atomic.AddInt64(&s.stNextFileNum, 1) - 1
}


func (s *session) reuseFileNum(num int64) {
	for {
		old, x := s.stNextFileNum, num
		if old != x+1 {
			x = old
		}
		if atomic.CompareAndSwapInt64(&s.stNextFileNum, old, x) {
			break
		}
	}
}


func (s *session) setCompPtr(level int, ik internalKey) {
	if level >= len(s.stCompPtrs) {
		newCompPtrs := make([]internalKey, level+1)
		copy(newCompPtrs, s.stCompPtrs)
		s.stCompPtrs = newCompPtrs
	}
	s.stCompPtrs[level] = append(internalKey{}, ik...)
}


func (s *session) getCompPtr(level int) internalKey {
	if level >= len(s.stCompPtrs) {
		return nil
	}
	return s.stCompPtrs[level]
}





func (s *session) fillRecord(r *sessionRecord, snapshot bool) {
	r.setNextFileNum(s.nextFileNum())

	if snapshot {
		if !r.has(recJournalNum) {
			r.setJournalNum(s.stJournalNum)
		}

		if !r.has(recSeqNum) {
			r.setSeqNum(s.stSeqNum)
		}

		for level, ik := range s.stCompPtrs {
			if ik != nil {
				r.addCompPtr(level, ik)
			}
		}

		r.setComparer(s.icmp.uName())
	}
}



func (s *session) recordCommited(rec *sessionRecord) {
	if rec.has(recJournalNum) {
		s.stJournalNum = rec.journalNum
	}

	if rec.has(recPrevJournalNum) {
		s.stPrevJournalNum = rec.prevJournalNum
	}

	if rec.has(recSeqNum) {
		s.stSeqNum = rec.seqNum
	}

	for _, r := range rec.compPtrs {
		s.setCompPtr(r.level, internalKey(r.ikey))
	}
}


func (s *session) newManifest(rec *sessionRecord, v *version) (err error) {
	fd := storage.FileDesc{storage.TypeManifest, s.allocFileNum()}
	writer, err := s.stor.Create(fd)
	if err != nil {
		return
	}
	jw := journal.NewWriter(writer)

	if v == nil {
		v = s.version()
		defer v.release()
	}
	if rec == nil {
		rec = &sessionRecord{}
	}
	s.fillRecord(rec, true)
	v.fillRecord(rec)

	defer func() {
		if err == nil {
			s.recordCommited(rec)
			if s.manifest != nil {
				s.manifest.Close()
			}
			if s.manifestWriter != nil {
				s.manifestWriter.Close()
			}
			if !s.manifestFd.Zero() {
				s.stor.Remove(s.manifestFd)
			}
			s.manifestFd = fd
			s.manifestWriter = writer
			s.manifest = jw
		} else {
			writer.Close()
			s.stor.Remove(fd)
			s.reuseFileNum(fd.Num)
		}
	}()

	w, err := jw.Next()
	if err != nil {
		return
	}
	err = rec.encode(w)
	if err != nil {
		return
	}
	err = jw.Flush()
	if err != nil {
		return
	}
	err = s.stor.SetMeta(fd)
	return
}


func (s *session) flushManifest(rec *sessionRecord) (err error) {
	s.fillRecord(rec, false)
	w, err := s.manifest.Next()
	if err != nil {
		return
	}
	err = rec.encode(w)
	if err != nil {
		return
	}
	err = s.manifest.Flush()
	if err != nil {
		return
	}
	if !s.o.GetNoSync() {
		err = s.manifestWriter.Sync()
		if err != nil {
			return
		}
	}
	s.recordCommited(rec)
	return
}
