

package transport

import (
	"fmt"
	"io"
	"math"
	"sync"
	"time"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/hpack"
)

const (
	
	defaultWindowSize = 65535
	
	initialWindowSize             = defaultWindowSize 
	infinity                      = time.Duration(math.MaxInt64)
	defaultClientKeepaliveTime    = infinity
	defaultClientKeepaliveTimeout = time.Duration(20 * time.Second)
	defaultMaxStreamsClient       = 100
	defaultMaxConnectionIdle      = infinity
	defaultMaxConnectionAge       = infinity
	defaultMaxConnectionAgeGrace  = infinity
	defaultServerKeepaliveTime    = time.Duration(2 * time.Hour)
	defaultServerKeepaliveTimeout = time.Duration(20 * time.Second)
	defaultKeepalivePolicyMinTime = time.Duration(5 * time.Minute)
	
	maxWindowSize = math.MaxInt32
	
	
	
	defaultLocalSendQuota = 128 * 1024
)





type headerFrame struct {
	streamID  uint32
	hf        []hpack.HeaderField
	endStream bool
}

func (*headerFrame) item() {}

type continuationFrame struct {
	streamID            uint32
	endHeaders          bool
	headerBlockFragment []byte
}

type dataFrame struct {
	streamID  uint32
	endStream bool
	d         []byte
	f         func()
}

func (*dataFrame) item() {}

func (*continuationFrame) item() {}

type windowUpdate struct {
	streamID  uint32
	increment uint32
}

func (*windowUpdate) item() {}

type settings struct {
	ss []http2.Setting
}

func (*settings) item() {}

type settingsAck struct {
}

func (*settingsAck) item() {}

type resetStream struct {
	streamID uint32
	code     http2.ErrCode
}

func (*resetStream) item() {}

type goAway struct {
	code      http2.ErrCode
	debugData []byte
	headsUp   bool
	closeConn bool
}

func (*goAway) item() {}

type flushIO struct {
	closeTr bool
}

func (*flushIO) item() {}

type ping struct {
	ack  bool
	data [8]byte
}

func (*ping) item() {}



type quotaPool struct {
	mu      sync.Mutex
	c       chan struct{}
	version uint32
	quota   int
}


func newQuotaPool(q int) *quotaPool {
	qb := &quotaPool{
		quota: q,
		c:     make(chan struct{}, 1),
	}
	return qb
}



func (qb *quotaPool) add(v int) {
	qb.mu.Lock()
	defer qb.mu.Unlock()
	qb.lockedAdd(v)
}

func (qb *quotaPool) lockedAdd(v int) {
	var wakeUp bool
	if qb.quota <= 0 {
		wakeUp = true 
	}
	qb.quota += v
	if wakeUp && qb.quota > 0 {
		select {
		case qb.c <- struct{}{}:
		default:
		}
	}
}

func (qb *quotaPool) addAndUpdate(v int) {
	qb.mu.Lock()
	qb.lockedAdd(v)
	qb.version++
	qb.mu.Unlock()
}

func (qb *quotaPool) get(v int, wc waiters) (int, uint32, error) {
	qb.mu.Lock()
	if qb.quota > 0 {
		if v > qb.quota {
			v = qb.quota
		}
		qb.quota -= v
		ver := qb.version
		qb.mu.Unlock()
		return v, ver, nil
	}
	qb.mu.Unlock()
	for {
		select {
		case <-wc.ctx.Done():
			return 0, 0, ContextErr(wc.ctx.Err())
		case <-wc.tctx.Done():
			return 0, 0, ErrConnClosing
		case <-wc.done:
			return 0, 0, io.EOF
		case <-wc.goAway:
			return 0, 0, errStreamDrain
		case <-qb.c:
			qb.mu.Lock()
			if qb.quota > 0 {
				if v > qb.quota {
					v = qb.quota
				}
				qb.quota -= v
				ver := qb.version
				if qb.quota > 0 {
					select {
					case qb.c <- struct{}{}:
					default:
					}
				}
				qb.mu.Unlock()
				return v, ver, nil

			}
			qb.mu.Unlock()
		}
	}
}

func (qb *quotaPool) compareAndExecute(version uint32, success, failure func()) bool {
	qb.mu.Lock()
	if version == qb.version {
		success()
		qb.mu.Unlock()
		return true
	}
	failure()
	qb.mu.Unlock()
	return false
}


type inFlow struct {
	mu sync.Mutex
	
	limit uint32
	
	
	pendingData uint32
	
	
	pendingUpdate uint32
	
	
	delta uint32
}



func (f *inFlow) newLimit(n uint32) uint32 {
	f.mu.Lock()
	defer f.mu.Unlock()
	d := n - f.limit
	f.limit = n
	return d
}

func (f *inFlow) maybeAdjust(n uint32) uint32 {
	if n > uint32(math.MaxInt32) {
		n = uint32(math.MaxInt32)
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	
	
	estSenderQuota := int32(f.limit - (f.pendingData + f.pendingUpdate))
	
	
	
	estUntransmittedData := int32(n - f.pendingData) 
	
	
	
	if estUntransmittedData > estSenderQuota {
		
		if f.limit+n > maxWindowSize {
			f.delta = maxWindowSize - f.limit
		} else {
			
			
			
			f.delta = n
		}
		return f.delta
	}
	return 0
}


func (f *inFlow) onData(n uint32) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.pendingData += n
	if f.pendingData+f.pendingUpdate > f.limit+f.delta {
		return fmt.Errorf("received %d-bytes data exceeding the limit %d bytes", f.pendingData+f.pendingUpdate, f.limit)
	}
	return nil
}



func (f *inFlow) onRead(n uint32) uint32 {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.pendingData == 0 {
		return 0
	}
	f.pendingData -= n
	if n > f.delta {
		n -= f.delta
		f.delta = 0
	} else {
		f.delta -= n
		n = 0
	}
	f.pendingUpdate += n
	if f.pendingUpdate >= f.limit/4 {
		wu := f.pendingUpdate
		f.pendingUpdate = 0
		return wu
	}
	return 0
}

func (f *inFlow) resetPendingUpdate() uint32 {
	f.mu.Lock()
	defer f.mu.Unlock()
	n := f.pendingUpdate
	f.pendingUpdate = 0
	return n
}
