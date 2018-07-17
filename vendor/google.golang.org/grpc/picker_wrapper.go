

package grpc

import (
	"sync"

	"golang.org/x/net/context"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
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
}

func newPickerWrapper() *pickerWrapper {
	bp := &pickerWrapper{blockingCh: make(chan struct{})}
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








func (bp *pickerWrapper) pick(ctx context.Context, failfast bool, opts balancer.PickOptions) (transport.ClientTransport, func(balancer.DoneInfo), error) {
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
