



package transport

import (
	"net"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"golang.org/x/net/context"
)


func dialContext(ctx context.Context, network, address string) (net.Conn, error) {
	return (&net.Dialer{Cancel: ctx.Done()}).Dial(network, address)
}


func ContextErr(err error) error {
	switch err {
	case context.DeadlineExceeded:
		return status.Error(codes.DeadlineExceeded, err.Error())
	case context.Canceled:
		return status.Error(codes.Canceled, err.Error())
	}
	return status.Errorf(codes.Internal, "Unexpected error from context packet: %v", err)
}


func contextFromRequest(r *http.Request) context.Context {
	return context.Background()
}
