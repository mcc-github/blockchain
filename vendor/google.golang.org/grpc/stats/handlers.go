

package stats

import (
	"net"

	"golang.org/x/net/context"
)


type ConnTagInfo struct {
	
	RemoteAddr net.Addr
	
	LocalAddr net.Addr
}


type RPCTagInfo struct {
	
	FullMethodName string
	
	
	FailFast bool
}


type Handler interface {
	
	
	
	TagRPC(context.Context, *RPCTagInfo) context.Context
	
	HandleRPC(context.Context, RPCStats)

	
	
	
	
	
	
	
	
	TagConn(context.Context, *ConnTagInfo) context.Context
	
	HandleConn(context.Context, ConnStats)
}
