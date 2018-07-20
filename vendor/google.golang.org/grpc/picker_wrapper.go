

package grpc

import (
	"io"
	"sync"
	"sync/atomic"

	"golang.org/x/net/context"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/internal/channelz"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/transport"
)



type pickerWrapper struct {
	mu         sync.Mutex
	done       bool
	blockingCh chan struct{}
	picker     balancer.Picker

	
	connErrMu sync.Mutex
	connErr   error

	stickinessMDKey atomic.Value
	stickiness      *stickyStore
}

func newPickerWrapper() *pickerWrapper {
	bp := &pickerWrapper{
		blockingCh: make(chan struct{}),
		stickiness: newStickyStore(),
	}
	return bp
}

func (bp *pickerWrapper) updateConnectionError(err error) {
	bp.connErrMu.Lock()
	bp.connErr = err
	bp.connErrMu.Unlock()
}

func (bp *pickerWrapper) connectionError() error {
	bp.connErrMu.Lock()
	err := bp.connErr
	bp.connErrMu.Unlock()
	return err
}

func (bp *pickerWrapper) updateStickinessMDKey(newKey string) {
	
	if oldKey, _ := bp.stickinessMDKey.Load().(string); oldKey != newKey {
		bp.stickinessMDKey.Store(newKey)
		bp.stickiness.reset(newKey)
	}
}

func (bp *pickerWrapper) getStickinessMDKey() string {
	
	mdKey, _ := bp.stickinessMDKey.Load().(string)
	return mdKey
}

func (bp *pickerWrapper) clearStickinessState() {
	if oldKey := bp.getStickinessMDKey(); oldKey != "" {
		
		bp.stickiness.reset(oldKey)
	}
}


func (bp *pickerWrapper) updatePicker(p balancer.Picker) {
	bp.mu.Lock()
	if bp.done {
		bp.mu.Unlock()
		return
	}
	bp.picker = p
	
	close(bp.blockingCh)
	bp.blockingCh = make(chan struct{})
	bp.mu.Unlock()
}

func doneChannelzWrapper(acw *acBalancerWrapper, done func(balancer.DoneInfo)) func(balancer.DoneInfo) {
	acw.mu.Lock()
	ac := acw.ac
	acw.mu.Unlock()
	ac.incrCallsStarted()
	return func(b balancer.DoneInfo) {
		if b.Err != nil && b.Err != io.EOF {
			ac.incrCallsFailed()
		} else {
			ac.incrCallsSucceeded()
		}
		if done != nil {
			done(b)
		}
	}
}








func (bp *pickerWrapper) pick(ctx context.Context, failfast bool, opts balancer.PickOptions) (transport.ClientTransport, func(balancer.DoneInfo), error) {

	mdKey := bp.getStickinessMDKey()
	stickyKey, isSticky := stickyKeyFromContext(ctx, mdKey)

	
	
	
	
	
	
	
	
	

	if isSticky {
		if t, ok := bp.stickiness.get(mdKey, stickyKey); ok {
			
			return t, nil, nil
		}
	}

	var (
		p  balancer.Picker
		ch chan struct{}
	)

	for {
		bp.mu.Lock()
		if bp.done {
			bp.mu.Unlock()
			return nil, nil, ErrClientConnClosing
		}

		if bp.picker == nil {
			ch = bp.blockingCh
		}
		if ch == bp.blockingCh {
			
			
			
			bp.mu.Unlock()
			select {
			case <-ctx.Done():
				return nil, nil, ctx.Err()
			case <-ch:
			}
			continue
		}

		ch = bp.blockingCh
		p = bp.picker
		bp.mu.Unlock()

		subConn, done, err := p.Pick(ctx, opts)

		if err != nil {
			switch err {
			case balancer.ErrNoSubConnAvailable:
				continue
			case balancer.ErrTransientFailure:
				if !failfast {
					continue
				}
				return nil, nil, status.Errorf(codes.Unavailable, "%v, latest connection error: %v", err, bp.connectionError())
			default:
				
				return nil, nil, toRPCErr(err)
			}
		}

		acw, ok := subConn.(*acBalancerWrapper)
		if !ok {
			grpclog.Infof("subconn returned from pick is not *acBalancerWrapper")
			continue
		}
		if t, ok := acw.getAddrConn().getReadyTransport(); ok {
			if isSticky {
				bp.stickiness.put(mdKey, stickyKey, acw)
			}
			if channelz.IsOn() {
				return t, doneChannelzWrapper(acw, done), nil
			}
			return t, done, nil
		}
		grpclog.Infof("blockingPicker: the picked transport is not ready, loop back to repick")
		
		
		
		
	}
}

func (bp *pickerWrapper) close() {
	bp.mu.Lock()
	defer bp.mu.Unlock()
	if bp.done {
		return
	}
	bp.done = true
	close(bp.blockingCh)
}

const stickinessKeyCountLimit = 1000

type stickyStoreEntry struct {
	acw  *acBalancerWrapper
	addr resolver.Address
}

type stickyStore struct {
	mu sync.Mutex
	
	
	curMDKey string
	store    *linkedMap
}

func newStickyStore() *stickyStore {
	return &stickyStore{
		store: newLinkedMap(),
	}
}


func (ss *stickyStore) reset(newMDKey string) {
	ss.mu.Lock()
	ss.curMDKey = newMDKey
	ss.store.clear()
	ss.mu.Unlock()
}



func (ss *stickyStore) put(mdKey, stickyKey string, acw *acBalancerWrapper) {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	if mdKey != ss.curMDKey {
		return
	}
	
	ss.store.put(stickyKey, &stickyStoreEntry{
		acw:  acw,
		addr: acw.getAddrConn().getCurAddr(),
	})
	if ss.store.len() > stickinessKeyCountLimit {
		ss.store.removeOldest()
	}
}



func (ss *stickyStore) get(mdKey, stickyKey string) (transport.ClientTransport, bool) {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	if mdKey != ss.curMDKey {
		return nil, false
	}
	entry, ok := ss.store.get(stickyKey)
	if !ok {
		return nil, false
	}
	ac := entry.acw.getAddrConn()
	if ac.getCurAddr() != entry.addr {
		ss.store.remove(stickyKey)
		return nil, false
	}
	t, ok := ac.getReadyTransport()
	if !ok {
		ss.store.remove(stickyKey)
		return nil, false
	}
	return t, true
}




func stickyKeyFromContext(ctx context.Context, stickinessMDKey string) (string, bool) {
	if stickinessMDKey == "" {
		return "", false
	}

	md, added, ok := metadata.FromOutgoingContextRaw(ctx)
	if !ok {
		return "", false
	}

	if vv, ok := md[stickinessMDKey]; ok {
		if len(vv) > 0 {
			return vv[0], true
		}
	}

	for _, ss := range added {
		for i := 0; i < len(ss)-1; i += 2 {
			if ss[i] == stickinessMDKey {
				return ss[i+1], true
			}
		}
	}

	return "", false
}
