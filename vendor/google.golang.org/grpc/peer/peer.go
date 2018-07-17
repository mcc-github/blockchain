



package peer

import (
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc/credentials"
)



type Peer struct {
	
	Addr net.Addr
	
	
	AuthInfo credentials.AuthInfo
}

type peerKey struct{}


func NewContext(ctx context.Context, p *Peer) context.Context {
	return context.WithValue(ctx, peerKey{}, p)
}


func FromContext(ctx context.Context) (p *Peer, ok bool) {
	p, ok = ctx.Value(peerKey{}).(*Peer)
	return
}
