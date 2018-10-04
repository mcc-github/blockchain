



package grpc

import (
	"fmt"
	"io"
	"net"
	"net/http"

	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/internal/transport"
	"google.golang.org/grpc/status"
)


func dialContext(ctx context.Context, network, address string) (net.Conn, error) {
	return (&net.Dialer{Cancel: ctx.Done()}).Dial(network, address)
}

func sendHTTPRequest(ctx context.Context, req *http.Request, conn net.Conn) error {
	req.Cancel = ctx.Done()
	if err := req.Write(conn); err != nil {
		return fmt.Errorf("failed to write the HTTP request: %v", err)
	}
	return nil
}


func toRPCErr(err error) error {
	if err == nil || err == io.EOF {
		return err
	}
	if err == io.ErrUnexpectedEOF {
		return status.Error(codes.Internal, err.Error())
	}
	if _, ok := status.FromError(err); ok {
		return err
	}
	switch e := err.(type) {
	case transport.ConnectionError:
		return status.Error(codes.Unavailable, e.Desc)
	default:
		switch err {
		case context.DeadlineExceeded:
			return status.Error(codes.DeadlineExceeded, err.Error())
		case context.Canceled:
			return status.Error(codes.Canceled, err.Error())
		}
	}
	return status.Error(codes.Unknown, err.Error())
}
