

package grpc

import (
	"golang.org/x/net/context"
)





func (cc *ClientConn) Invoke(ctx context.Context, method string, args, reply interface{}, opts ...CallOption) error {
	
	
	opts = combine(cc.dopts.callOptions, opts)

	if cc.dopts.unaryInt != nil {
		return cc.dopts.unaryInt(ctx, method, args, reply, cc, invoke, opts...)
	}
	return invoke(ctx, method, args, reply, cc, opts...)
}

func combine(o1 []CallOption, o2 []CallOption) []CallOption {
	
	
	
	if len(o1) == 0 {
		return o2
	} else if len(o2) == 0 {
		return o1
	}
	ret := make([]CallOption, len(o1)+len(o2))
	copy(ret, o1)
	copy(ret[len(o1):], o2)
	return ret
}





func Invoke(ctx context.Context, method string, args, reply interface{}, cc *ClientConn, opts ...CallOption) error {
	return cc.Invoke(ctx, method, args, reply, opts...)
}

var unaryStreamDesc = &StreamDesc{ServerStreams: false, ClientStreams: false}

func invoke(ctx context.Context, method string, req, reply interface{}, cc *ClientConn, opts ...CallOption) error {
	
	
	firstAttempt := true
	for {
		csInt, err := newClientStream(ctx, unaryStreamDesc, cc, method, opts...)
		if err != nil {
			return err
		}
		cs := csInt.(*clientStream)
		if err := cs.SendMsg(req); err != nil {
			if !cs.c.failFast && cs.attempt.s.Unprocessed() && firstAttempt {
				
				firstAttempt = false
				continue
			}
			return err
		}
		if err := cs.RecvMsg(reply); err != nil {
			if !cs.c.failFast && cs.attempt.s.Unprocessed() && firstAttempt {
				
				firstAttempt = false
				continue
			}
			return err
		}
		return nil
	}
}
