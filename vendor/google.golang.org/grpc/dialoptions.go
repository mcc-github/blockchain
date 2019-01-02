

package grpc

import (
	"context"
	"fmt"
	"net"
	"time"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/internal"
	"google.golang.org/grpc/internal/backoff"
	"google.golang.org/grpc/internal/envconfig"
	"google.golang.org/grpc/internal/transport"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/stats"
)



type dialOptions struct {
	unaryInt    UnaryClientInterceptor
	streamInt   StreamClientInterceptor
	cp          Compressor
	dc          Decompressor
	bs          backoff.Strategy
	block       bool
	insecure    bool
	timeout     time.Duration
	scChan      <-chan ServiceConfig
	authority   string
	copts       transport.ConnectOptions
	callOptions []CallOption
	
	
	balancerBuilder balancer.Builder
	
	resolverBuilder      resolver.Builder
	reqHandshake         envconfig.RequireHandshakeSetting
	channelzParentID     int64
	disableServiceConfig bool
	disableRetry         bool
	disableHealthCheck   bool
}


type DialOption interface {
	apply(*dialOptions)
}





type EmptyDialOption struct{}

func (EmptyDialOption) apply(*dialOptions) {}



type funcDialOption struct {
	f func(*dialOptions)
}

func (fdo *funcDialOption) apply(do *dialOptions) {
	fdo.f(do)
}

func newFuncDialOption(f func(*dialOptions)) *funcDialOption {
	return &funcDialOption{
		f: f,
	}
}








func WithWaitForHandshake() DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.reqHandshake = envconfig.RequireHandshakeOn
	})
}








func WithWriteBufferSize(s int) DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.copts.WriteBufferSize = s
	})
}






func WithReadBufferSize(s int) DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.copts.ReadBufferSize = s
	})
}




func WithInitialWindowSize(s int32) DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.copts.InitialWindowSize = s
	})
}




func WithInitialConnWindowSize(s int32) DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.copts.InitialConnWindowSize = s
	})
}





func WithMaxMsgSize(s int) DialOption {
	return WithDefaultCallOptions(MaxCallRecvMsgSize(s))
}



func WithDefaultCallOptions(cos ...CallOption) DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.callOptions = append(o.callOptions, cos...)
	})
}





func WithCodec(c Codec) DialOption {
	return WithDefaultCallOptions(CallCustomCodec(c))
}






func WithCompressor(cp Compressor) DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.cp = cp
	})
}










func WithDecompressor(dc Decompressor) DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.dc = dc
	})
}






func WithBalancer(b Balancer) DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.balancerBuilder = &balancerWrapperBuilder{
			b: b,
		}
	})
}









func WithBalancerName(balancerName string) DialOption {
	builder := balancer.Get(balancerName)
	if builder == nil {
		panic(fmt.Sprintf("grpc.WithBalancerName: no balancer is registered for name %v", balancerName))
	}
	return newFuncDialOption(func(o *dialOptions) {
		o.balancerBuilder = builder
	})
}


func withResolverBuilder(b resolver.Builder) DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.resolverBuilder = b
	})
}







func WithServiceConfig(c <-chan ServiceConfig) DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.scChan = c
	})
}



func WithBackoffMaxDelay(md time.Duration) DialOption {
	return WithBackoffConfig(BackoffConfig{MaxDelay: md})
}






func WithBackoffConfig(b BackoffConfig) DialOption {
	return withBackoff(backoff.Exponential{
		MaxDelay: b.MaxDelay,
	})
}





func withBackoff(bs backoff.Strategy) DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.bs = bs
	})
}




func WithBlock() DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.block = true
	})
}




func WithInsecure() DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.insecure = true
	})
}




func WithTransportCredentials(creds credentials.TransportCredentials) DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.copts.TransportCredentials = creds
	})
}



func WithPerRPCCredentials(creds credentials.PerRPCCredentials) DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.copts.PerRPCCredentials = append(o.copts.PerRPCCredentials, creds)
	})
}






func WithCredentialsBundle(b credentials.Bundle) DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.copts.CredsBundle = b
	})
}





func WithTimeout(d time.Duration) DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.timeout = d
	})
}

func withContextDialer(f func(context.Context, string) (net.Conn, error)) DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.copts.Dialer = f
	})
}

func init() {
	internal.WithContextDialer = withContextDialer
	internal.WithResolverBuilder = withResolverBuilder
}





func WithDialer(f func(string, time.Duration) (net.Conn, error)) DialOption {
	return withContextDialer(
		func(ctx context.Context, addr string) (net.Conn, error) {
			if deadline, ok := ctx.Deadline(); ok {
				return f(addr, deadline.Sub(time.Now()))
			}
			return f(addr, 0)
		})
}



func WithStatsHandler(h stats.Handler) DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.copts.StatsHandler = h
	})
}










func FailOnNonTempDialError(f bool) DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.copts.FailOnNonTempDialError = f
	})
}



func WithUserAgent(s string) DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.copts.UserAgent = s
	})
}



func WithKeepaliveParams(kp keepalive.ClientParameters) DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.copts.KeepaliveParams = kp
	})
}



func WithUnaryInterceptor(f UnaryClientInterceptor) DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.unaryInt = f
	})
}



func WithStreamInterceptor(f StreamClientInterceptor) DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.streamInt = f
	})
}




func WithAuthority(a string) DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.authority = a
	})
}




func WithChannelzParentID(id int64) DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.channelzParentID = id
	})
}




func WithDisableServiceConfig() DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.disableServiceConfig = true
	})
}











func WithDisableRetry() DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.disableRetry = true
	})
}



func WithMaxHeaderListSize(s uint32) DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.copts.MaxHeaderListSize = &s
	})
}




func WithDisableHealthCheck() DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.disableHealthCheck = true
	})
}
func defaultDialOptions() dialOptions {
	return dialOptions{
		disableRetry: !envconfig.Retry,
		reqHandshake: envconfig.RequireHandshake,
		copts: transport.ConnectOptions{
			WriteBufferSize: defaultWriteBufSize,
			ReadBufferSize:  defaultReadBufSize,
		},
	}
}
