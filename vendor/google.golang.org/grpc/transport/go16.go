



package transport

import (
	"net"
	"net/http"

	"google.golang.org/grpc/codes"

	"golang.org/x/net/context"
)


func dialContext(ctx context.Context, network, address string) (net.Conn, error) {
	return (&net.Dialer{Cancel: ctx.Done()}).Dial(network, address)
}


func ContextErr(err error) StreamError {
	switch err {
	case context.DeadlineExceeded:
		return streamErrorf(codes.DeadlineExceeded, "%v", err)
	case context.Canceled:
		return streamErrorf(codes.Canceled, "%v", err)
	}
	return streamErrorf(codes.Internal, "Unexpected error from context packet: %v", err)
}


func contextFromRequest(r *http.Request) context.Context {
	return context.Background()
}
