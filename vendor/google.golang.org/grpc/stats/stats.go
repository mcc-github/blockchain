






package stats 

import (
	"context"
	"net"
	"time"

	"google.golang.org/grpc/metadata"
)


type RPCStats interface {
	isRPCStats()
	
	IsClient() bool
}



type Begin struct {
	
	Client bool
	
	BeginTime time.Time
	
	FailFast bool
}


func (s *Begin) IsClient() bool { return s.Client }

func (s *Begin) isRPCStats() {}


type InPayload struct {
	
	Client bool
	
	Payload interface{}
	
	Data []byte
	
	Length int
	
	WireLength int
	
	RecvTime time.Time
}


func (s *InPayload) IsClient() bool { return s.Client }

func (s *InPayload) isRPCStats() {}


type InHeader struct {
	
	Client bool
	
	WireLength int

	
	
	FullMethod string
	
	RemoteAddr net.Addr
	
	LocalAddr net.Addr
	
	Compression string
}


func (s *InHeader) IsClient() bool { return s.Client }

func (s *InHeader) isRPCStats() {}


type InTrailer struct {
	
	Client bool
	
	WireLength int
}


func (s *InTrailer) IsClient() bool { return s.Client }

func (s *InTrailer) isRPCStats() {}


type OutPayload struct {
	
	Client bool
	
	Payload interface{}
	
	Data []byte
	
	Length int
	
	WireLength int
	
	SentTime time.Time
}


func (s *OutPayload) IsClient() bool { return s.Client }

func (s *OutPayload) isRPCStats() {}


type OutHeader struct {
	
	Client bool

	
	
	FullMethod string
	
	RemoteAddr net.Addr
	
	LocalAddr net.Addr
	
	Compression string
}


func (s *OutHeader) IsClient() bool { return s.Client }

func (s *OutHeader) isRPCStats() {}


type OutTrailer struct {
	
	Client bool
	
	WireLength int
}


func (s *OutTrailer) IsClient() bool { return s.Client }

func (s *OutTrailer) isRPCStats() {}


type End struct {
	
	Client bool
	
	BeginTime time.Time
	
	EndTime time.Time
	
	
	Trailer metadata.MD
	
	
	
	Error error
}


func (s *End) IsClient() bool { return s.Client }

func (s *End) isRPCStats() {}


type ConnStats interface {
	isConnStats()
	
	IsClient() bool
}


type ConnBegin struct {
	
	Client bool
}


func (s *ConnBegin) IsClient() bool { return s.Client }

func (s *ConnBegin) isConnStats() {}


type ConnEnd struct {
	
	Client bool
}


func (s *ConnEnd) IsClient() bool { return s.Client }

func (s *ConnEnd) isConnStats() {}

type incomingTagsKey struct{}
type outgoingTagsKey struct{}









func SetTags(ctx context.Context, b []byte) context.Context {
	return context.WithValue(ctx, outgoingTagsKey{}, b)
}







func Tags(ctx context.Context) []byte {
	b, _ := ctx.Value(incomingTagsKey{}).([]byte)
	return b
}





func SetIncomingTags(ctx context.Context, b []byte) context.Context {
	return context.WithValue(ctx, incomingTagsKey{}, b)
}




func OutgoingTags(ctx context.Context) []byte {
	b, _ := ctx.Value(outgoingTagsKey{}).([]byte)
	return b
}

type incomingTraceKey struct{}
type outgoingTraceKey struct{}









func SetTrace(ctx context.Context, b []byte) context.Context {
	return context.WithValue(ctx, outgoingTraceKey{}, b)
}







func Trace(ctx context.Context) []byte {
	b, _ := ctx.Value(incomingTraceKey{}).([]byte)
	return b
}




func SetIncomingTrace(ctx context.Context, b []byte) context.Context {
	return context.WithValue(ctx, incomingTraceKey{}, b)
}



func OutgoingTrace(ctx context.Context) []byte {
	b, _ := ctx.Value(outgoingTraceKey{}).([]byte)
	return b
}
