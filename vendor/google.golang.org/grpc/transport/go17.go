



package transport

import (
	"context"
	"net"
	"net/http"

	"google.golang.org/grpc/codes"

	netctx "golang.org/x/net/context"
)


func dialContext(ctx context.Context, network, address string) (net.Conn, error) {
	return (&net.Dialer{}).DialContext(ctx, network, address)
}


func ContextErr(err error) StreamError {
	switch err {
	case context.DeadlineExceeded, netctx.DeadlineExceeded:
		return streamErrorf(codes.DeadlineExceeded, "%v", err)
	case context.Canceled, netctx.Canceled:
		return streamErrorf(codes.Canceled, "%v", err)
	}
	return streamErrorf(codes.Internal, "Unexpected error from context packet: %v", err)
}


func contextFromRequest(r *http.Request) context.Context {
	return r.Context()
}
