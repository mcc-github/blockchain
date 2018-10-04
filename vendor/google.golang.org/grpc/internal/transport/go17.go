



package transport

import (
	"context"
	"net"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	netctx "golang.org/x/net/context"
)


func dialContext(ctx context.Context, network, address string) (net.Conn, error) {
	return (&net.Dialer{}).DialContext(ctx, network, address)
}


func ContextErr(err error) error {
	switch err {
	case context.DeadlineExceeded, netctx.DeadlineExceeded:
		return status.Error(codes.DeadlineExceeded, err.Error())
	case context.Canceled, netctx.Canceled:
		return status.Error(codes.Canceled, err.Error())
	}
	return status.Errorf(codes.Internal, "Unexpected error from context packet: %v", err)
}


func contextFromRequest(r *http.Request) context.Context {
	return r.Context()
}
