



package status

import (
	"context"

	netctx "golang.org/x/net/context"
	"google.golang.org/grpc/codes"
)




func FromContextError(err error) *Status {
	switch err {
	case nil:
		return New(codes.OK, "")
	case context.DeadlineExceeded, netctx.DeadlineExceeded:
		return New(codes.DeadlineExceeded, err.Error())
	case context.Canceled, netctx.Canceled:
		return New(codes.Canceled, err.Error())
	default:
		return New(codes.Unknown, err.Error())
	}
}
