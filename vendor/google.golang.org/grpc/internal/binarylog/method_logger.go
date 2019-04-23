

package binarylog

import (
	"net"
	"strings"
	"sync/atomic"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	pb "google.golang.org/grpc/binarylog/grpc_binarylog_v1"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type callIDGenerator struct {
	id uint64
}

func (g *callIDGenerator) next() uint64 {
	id := atomic.AddUint64(&g.id, 1)
	return id
}


func (g *callIDGenerator) reset() {
	g.id = 0
}

var idGen callIDGenerator


type MethodLogger struct {
	headerMaxLen, messageMaxLen uint64

	callID          uint64
	idWithinCallGen *callIDGenerator

	sink Sink 
}

func newMethodLogger(h, m uint64) *MethodLogger {
	return &MethodLogger{
		headerMaxLen:  h,
		messageMaxLen: m,

		callID:          idGen.next(),
		idWithinCallGen: &callIDGenerator{},

		sink: defaultSink, 
	}
}


func (ml *MethodLogger) Log(c LogEntryConfig) {
	m := c.toProto()
	timestamp, _ := ptypes.TimestampProto(time.Now())
	m.Timestamp = timestamp
	m.CallId = ml.callID
	m.SequenceIdWithinCall = ml.idWithinCallGen.next()

	switch pay := m.Payload.(type) {
	case *pb.GrpcLogEntry_ClientHeader:
		m.PayloadTruncated = ml.truncateMetadata(pay.ClientHeader.GetMetadata())
	case *pb.GrpcLogEntry_ServerHeader:
		m.PayloadTruncated = ml.truncateMetadata(pay.ServerHeader.GetMetadata())
	case *pb.GrpcLogEntry_Message:
		m.PayloadTruncated = ml.truncateMessage(pay.Message)
	}

	ml.sink.Write(m)
}

func (ml *MethodLogger) truncateMetadata(mdPb *pb.Metadata) (truncated bool) {
	if ml.headerMaxLen == maxUInt {
		return false
	}
	var (
		bytesLimit = ml.headerMaxLen
		index      int
	)
	
	
	
	
	for ; index < len(mdPb.Entry); index++ {
		entry := mdPb.Entry[index]
		if entry.Key == "grpc-trace-bin" {
			
			
			continue
		}
		currentEntryLen := uint64(len(entry.Value))
		if currentEntryLen > bytesLimit {
			break
		}
		bytesLimit -= currentEntryLen
	}
	truncated = index < len(mdPb.Entry)
	mdPb.Entry = mdPb.Entry[:index]
	return truncated
}

func (ml *MethodLogger) truncateMessage(msgPb *pb.Message) (truncated bool) {
	if ml.messageMaxLen == maxUInt {
		return false
	}
	if ml.messageMaxLen >= uint64(len(msgPb.Data)) {
		return false
	}
	msgPb.Data = msgPb.Data[:ml.messageMaxLen]
	return true
}


type LogEntryConfig interface {
	toProto() *pb.GrpcLogEntry
}


type ClientHeader struct {
	OnClientSide bool
	Header       metadata.MD
	MethodName   string
	Authority    string
	Timeout      time.Duration
	
	PeerAddr net.Addr
}

func (c *ClientHeader) toProto() *pb.GrpcLogEntry {
	
	
	clientHeader := &pb.ClientHeader{
		Metadata:   mdToMetadataProto(c.Header),
		MethodName: c.MethodName,
		Authority:  c.Authority,
	}
	if c.Timeout > 0 {
		clientHeader.Timeout = ptypes.DurationProto(c.Timeout)
	}
	ret := &pb.GrpcLogEntry{
		Type: pb.GrpcLogEntry_EVENT_TYPE_CLIENT_HEADER,
		Payload: &pb.GrpcLogEntry_ClientHeader{
			ClientHeader: clientHeader,
		},
	}
	if c.OnClientSide {
		ret.Logger = pb.GrpcLogEntry_LOGGER_CLIENT
	} else {
		ret.Logger = pb.GrpcLogEntry_LOGGER_SERVER
	}
	if c.PeerAddr != nil {
		ret.Peer = addrToProto(c.PeerAddr)
	}
	return ret
}


type ServerHeader struct {
	OnClientSide bool
	Header       metadata.MD
	
	PeerAddr net.Addr
}

func (c *ServerHeader) toProto() *pb.GrpcLogEntry {
	ret := &pb.GrpcLogEntry{
		Type: pb.GrpcLogEntry_EVENT_TYPE_SERVER_HEADER,
		Payload: &pb.GrpcLogEntry_ServerHeader{
			ServerHeader: &pb.ServerHeader{
				Metadata: mdToMetadataProto(c.Header),
			},
		},
	}
	if c.OnClientSide {
		ret.Logger = pb.GrpcLogEntry_LOGGER_CLIENT
	} else {
		ret.Logger = pb.GrpcLogEntry_LOGGER_SERVER
	}
	if c.PeerAddr != nil {
		ret.Peer = addrToProto(c.PeerAddr)
	}
	return ret
}


type ClientMessage struct {
	OnClientSide bool
	
	
	Message interface{}
}

func (c *ClientMessage) toProto() *pb.GrpcLogEntry {
	var (
		data []byte
		err  error
	)
	if m, ok := c.Message.(proto.Message); ok {
		data, err = proto.Marshal(m)
		if err != nil {
			grpclog.Infof("binarylogging: failed to marshal proto message: %v", err)
		}
	} else if b, ok := c.Message.([]byte); ok {
		data = b
	} else {
		grpclog.Infof("binarylogging: message to log is neither proto.message nor []byte")
	}
	ret := &pb.GrpcLogEntry{
		Type: pb.GrpcLogEntry_EVENT_TYPE_CLIENT_MESSAGE,
		Payload: &pb.GrpcLogEntry_Message{
			Message: &pb.Message{
				Length: uint32(len(data)),
				Data:   data,
			},
		},
	}
	if c.OnClientSide {
		ret.Logger = pb.GrpcLogEntry_LOGGER_CLIENT
	} else {
		ret.Logger = pb.GrpcLogEntry_LOGGER_SERVER
	}
	return ret
}


type ServerMessage struct {
	OnClientSide bool
	
	
	Message interface{}
}

func (c *ServerMessage) toProto() *pb.GrpcLogEntry {
	var (
		data []byte
		err  error
	)
	if m, ok := c.Message.(proto.Message); ok {
		data, err = proto.Marshal(m)
		if err != nil {
			grpclog.Infof("binarylogging: failed to marshal proto message: %v", err)
		}
	} else if b, ok := c.Message.([]byte); ok {
		data = b
	} else {
		grpclog.Infof("binarylogging: message to log is neither proto.message nor []byte")
	}
	ret := &pb.GrpcLogEntry{
		Type: pb.GrpcLogEntry_EVENT_TYPE_SERVER_MESSAGE,
		Payload: &pb.GrpcLogEntry_Message{
			Message: &pb.Message{
				Length: uint32(len(data)),
				Data:   data,
			},
		},
	}
	if c.OnClientSide {
		ret.Logger = pb.GrpcLogEntry_LOGGER_CLIENT
	} else {
		ret.Logger = pb.GrpcLogEntry_LOGGER_SERVER
	}
	return ret
}


type ClientHalfClose struct {
	OnClientSide bool
}

func (c *ClientHalfClose) toProto() *pb.GrpcLogEntry {
	ret := &pb.GrpcLogEntry{
		Type:    pb.GrpcLogEntry_EVENT_TYPE_CLIENT_HALF_CLOSE,
		Payload: nil, 
	}
	if c.OnClientSide {
		ret.Logger = pb.GrpcLogEntry_LOGGER_CLIENT
	} else {
		ret.Logger = pb.GrpcLogEntry_LOGGER_SERVER
	}
	return ret
}


type ServerTrailer struct {
	OnClientSide bool
	Trailer      metadata.MD
	
	Err error
	
	
	PeerAddr net.Addr
}

func (c *ServerTrailer) toProto() *pb.GrpcLogEntry {
	st, ok := status.FromError(c.Err)
	if !ok {
		grpclog.Info("binarylogging: error in trailer is not a status error")
	}
	var (
		detailsBytes []byte
		err          error
	)
	stProto := st.Proto()
	if stProto != nil && len(stProto.Details) != 0 {
		detailsBytes, err = proto.Marshal(stProto)
		if err != nil {
			grpclog.Infof("binarylogging: failed to marshal status proto: %v", err)
		}
	}
	ret := &pb.GrpcLogEntry{
		Type: pb.GrpcLogEntry_EVENT_TYPE_SERVER_TRAILER,
		Payload: &pb.GrpcLogEntry_Trailer{
			Trailer: &pb.Trailer{
				Metadata:      mdToMetadataProto(c.Trailer),
				StatusCode:    uint32(st.Code()),
				StatusMessage: st.Message(),
				StatusDetails: detailsBytes,
			},
		},
	}
	if c.OnClientSide {
		ret.Logger = pb.GrpcLogEntry_LOGGER_CLIENT
	} else {
		ret.Logger = pb.GrpcLogEntry_LOGGER_SERVER
	}
	if c.PeerAddr != nil {
		ret.Peer = addrToProto(c.PeerAddr)
	}
	return ret
}


type Cancel struct {
	OnClientSide bool
}

func (c *Cancel) toProto() *pb.GrpcLogEntry {
	ret := &pb.GrpcLogEntry{
		Type:    pb.GrpcLogEntry_EVENT_TYPE_CANCEL,
		Payload: nil,
	}
	if c.OnClientSide {
		ret.Logger = pb.GrpcLogEntry_LOGGER_CLIENT
	} else {
		ret.Logger = pb.GrpcLogEntry_LOGGER_SERVER
	}
	return ret
}



func metadataKeyOmit(key string) bool {
	switch key {
	case "lb-token", ":path", ":authority", "content-encoding", "content-type", "user-agent", "te":
		return true
	case "grpc-trace-bin": 
		return false
	}
	return strings.HasPrefix(key, "grpc-")
}

func mdToMetadataProto(md metadata.MD) *pb.Metadata {
	ret := &pb.Metadata{}
	for k, vv := range md {
		if metadataKeyOmit(k) {
			continue
		}
		for _, v := range vv {
			ret.Entry = append(ret.Entry,
				&pb.MetadataEntry{
					Key:   k,
					Value: []byte(v),
				},
			)
		}
	}
	return ret
}

func addrToProto(addr net.Addr) *pb.Address {
	ret := &pb.Address{}
	switch a := addr.(type) {
	case *net.TCPAddr:
		if a.IP.To4() != nil {
			ret.Type = pb.Address_TYPE_IPV4
		} else if a.IP.To16() != nil {
			ret.Type = pb.Address_TYPE_IPV6
		} else {
			ret.Type = pb.Address_TYPE_UNKNOWN
			
			break
		}
		ret.Address = a.IP.String()
		ret.IpPort = uint32(a.Port)
	case *net.UnixAddr:
		ret.Type = pb.Address_TYPE_UNIX
		ret.Address = a.String()
	default:
		ret.Type = pb.Address_TYPE_UNKNOWN
	}
	return ret
}
