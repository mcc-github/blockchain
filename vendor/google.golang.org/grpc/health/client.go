

package health

import (
	"context"
	"fmt"
	"io"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/internal"
	"google.golang.org/grpc/internal/backoff"
	"google.golang.org/grpc/status"
)

const maxDelay = 120 * time.Second

var backoffStrategy = backoff.Exponential{MaxDelay: maxDelay}
var backoffFunc = func(ctx context.Context, retries int) bool {
	d := backoffStrategy.Backoff(retries)
	timer := time.NewTimer(d)
	select {
	case <-timer.C:
		return true
	case <-ctx.Done():
		timer.Stop()
		return false
	}
}

func init() {
	internal.HealthCheckFunc = clientHealthCheck
}

func clientHealthCheck(ctx context.Context, newStream func() (interface{}, error), reportHealth func(bool), service string) error {
	tryCnt := 0

retryConnection:
	for {
		
		if tryCnt > 0 && !backoffFunc(ctx, tryCnt-1) {
			return nil
		}
		tryCnt++

		if ctx.Err() != nil {
			return nil
		}
		rawS, err := newStream()
		if err != nil {
			continue retryConnection
		}

		s, ok := rawS.(grpc.ClientStream)
		
		if !ok {
			reportHealth(true)
			return fmt.Errorf("newStream returned %v (type %T); want grpc.ClientStream", rawS, rawS)
		}

		if err = s.SendMsg(&healthpb.HealthCheckRequest{Service: service}); err != nil && err != io.EOF {
			
			continue retryConnection
		}
		s.CloseSend()

		resp := new(healthpb.HealthCheckResponse)
		for {
			err = s.RecvMsg(resp)

			
			if status.Code(err) == codes.Unimplemented {
				reportHealth(true)
				return err
			}

			
			if err != nil {
				reportHealth(false)
				continue retryConnection
			}

			
			tryCnt = 0
			reportHealth(resp.Status == healthpb.HealthCheckResponse_SERVING)
		}
	}
}
