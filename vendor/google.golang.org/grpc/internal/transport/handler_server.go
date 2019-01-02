






package transport

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	"golang.org/x/net/http2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/stats"
	"google.golang.org/grpc/status"
)




func NewServerHandlerTransport(w http.ResponseWriter, r *http.Request, stats stats.Handler) (ServerTransport, error) {
	if r.ProtoMajor != 2 {
		return nil, errors.New("gRPC requires HTTP/2")
	}
	if r.Method != "POST" {
		return nil, errors.New("invalid gRPC request method")
	}
	contentType := r.Header.Get("Content-Type")
	
	contentSubtype, validContentType := contentSubtype(contentType)
	if !validContentType {
		return nil, errors.New("invalid gRPC request content-type")
	}
	if _, ok := w.(http.Flusher); !ok {
		return nil, errors.New("gRPC requires a ResponseWriter supporting http.Flusher")
	}
	if _, ok := w.(http.CloseNotifier); !ok {
		return nil, errors.New("gRPC requires a ResponseWriter supporting http.CloseNotifier")
	}

	st := &serverHandlerTransport{
		rw:             w,
		req:            r,
		closedCh:       make(chan struct{}),
		writes:         make(chan func()),
		contentType:    contentType,
		contentSubtype: contentSubtype,
		stats:          stats,
	}

	if v := r.Header.Get("grpc-timeout"); v != "" {
		to, err := decodeTimeout(v)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "malformed time-out: %v", err)
		}
		st.timeoutSet = true
		st.timeout = to
	}

	metakv := []string{"content-type", contentType}
	if r.Host != "" {
		metakv = append(metakv, ":authority", r.Host)
	}
	for k, vv := range r.Header {
		k = strings.ToLower(k)
		if isReservedHeader(k) && !isWhitelistedHeader(k) {
			continue
		}
		for _, v := range vv {
			v, err := decodeMetadataHeader(k, v)
			if err != nil {
				return nil, status.Errorf(codes.Internal, "malformed binary metadata: %v", err)
			}
			metakv = append(metakv, k, v)
		}
	}
	st.headerMD = metadata.Pairs(metakv...)

	return st, nil
}






type serverHandlerTransport struct {
	rw               http.ResponseWriter
	req              *http.Request
	timeoutSet       bool
	timeout          time.Duration
	didCommonHeaders bool

	headerMD metadata.MD

	closeOnce sync.Once
	closedCh  chan struct{} 

	
	
	
	writes chan func()

	
	
	writeStatusMu sync.Mutex

	
	contentType string
	
	
	contentSubtype string

	stats stats.Handler
}

func (ht *serverHandlerTransport) Close() error {
	ht.closeOnce.Do(ht.closeCloseChanOnce)
	return nil
}

func (ht *serverHandlerTransport) closeCloseChanOnce() { close(ht.closedCh) }

func (ht *serverHandlerTransport) RemoteAddr() net.Addr { return strAddr(ht.req.RemoteAddr) }



type strAddr string

func (a strAddr) Network() string {
	if a != "" {
		
		
		
		
		
		
		
		
		
		return "tcp"
	}
	return ""
}

func (a strAddr) String() string { return string(a) }


func (ht *serverHandlerTransport) do(fn func()) error {
	
	select {
	case <-ht.closedCh:
		return ErrConnClosing
	default:
		select {
		case ht.writes <- fn:
			return nil
		case <-ht.closedCh:
			return ErrConnClosing
		}
	}
}

func (ht *serverHandlerTransport) WriteStatus(s *Stream, st *status.Status) error {
	ht.writeStatusMu.Lock()
	defer ht.writeStatusMu.Unlock()

	err := ht.do(func() {
		ht.writeCommonHeaders(s)

		
		
		
		ht.rw.(http.Flusher).Flush()

		h := ht.rw.Header()
		h.Set("Grpc-Status", fmt.Sprintf("%d", st.Code()))
		if m := st.Message(); m != "" {
			h.Set("Grpc-Message", encodeGrpcMessage(m))
		}

		if p := st.Proto(); p != nil && len(p.Details) > 0 {
			stBytes, err := proto.Marshal(p)
			if err != nil {
				
				panic(err)
			}

			h.Set("Grpc-Status-Details-Bin", encodeBinHeader(stBytes))
		}

		if md := s.Trailer(); len(md) > 0 {
			for k, vv := range md {
				
				if isReservedHeader(k) {
					continue
				}
				for _, v := range vv {
					
					
					h.Add(http2.TrailerPrefix+k, encodeMetadataHeader(k, v))
				}
			}
		}
	})

	if err == nil { 
		if ht.stats != nil {
			ht.stats.HandleRPC(s.Context(), &stats.OutTrailer{})
		}
		close(ht.writes)
	}
	ht.Close()
	return err
}



func (ht *serverHandlerTransport) writeCommonHeaders(s *Stream) {
	if ht.didCommonHeaders {
		return
	}
	ht.didCommonHeaders = true

	h := ht.rw.Header()
	h["Date"] = nil 
	h.Set("Content-Type", ht.contentType)

	
	
	
	
	
	h.Add("Trailer", "Grpc-Status")
	h.Add("Trailer", "Grpc-Message")
	h.Add("Trailer", "Grpc-Status-Details-Bin")

	if s.sendCompress != "" {
		h.Set("Grpc-Encoding", s.sendCompress)
	}
}

func (ht *serverHandlerTransport) Write(s *Stream, hdr []byte, data []byte, opts *Options) error {
	return ht.do(func() {
		ht.writeCommonHeaders(s)
		ht.rw.Write(hdr)
		ht.rw.Write(data)
		ht.rw.(http.Flusher).Flush()
	})
}

func (ht *serverHandlerTransport) WriteHeader(s *Stream, md metadata.MD) error {
	err := ht.do(func() {
		ht.writeCommonHeaders(s)
		h := ht.rw.Header()
		for k, vv := range md {
			
			if isReservedHeader(k) {
				continue
			}
			for _, v := range vv {
				v = encodeMetadataHeader(k, v)
				h.Add(k, v)
			}
		}
		ht.rw.WriteHeader(200)
		ht.rw.(http.Flusher).Flush()
	})

	if err == nil {
		if ht.stats != nil {
			ht.stats.HandleRPC(s.Context(), &stats.OutHeader{})
		}
	}
	return err
}

func (ht *serverHandlerTransport) HandleStreams(startStream func(*Stream), traceCtx func(context.Context, string) context.Context) {
	

	ctx := ht.req.Context()
	var cancel context.CancelFunc
	if ht.timeoutSet {
		ctx, cancel = context.WithTimeout(ctx, ht.timeout)
	} else {
		ctx, cancel = context.WithCancel(ctx)
	}

	
	
	requestOver := make(chan struct{})

	
	
	
	clientGone := ht.rw.(http.CloseNotifier).CloseNotify()
	go func() {
		select {
		case <-requestOver:
		case <-ht.closedCh:
		case <-clientGone:
		}
		cancel()
		ht.Close()
	}()

	req := ht.req

	s := &Stream{
		id:             0, 
		requestRead:    func(int) {},
		cancel:         cancel,
		buf:            newRecvBuffer(),
		st:             ht,
		method:         req.URL.Path,
		recvCompress:   req.Header.Get("grpc-encoding"),
		contentSubtype: ht.contentSubtype,
	}
	pr := &peer.Peer{
		Addr: ht.RemoteAddr(),
	}
	if req.TLS != nil {
		pr.AuthInfo = credentials.TLSInfo{State: *req.TLS}
	}
	ctx = metadata.NewIncomingContext(ctx, ht.headerMD)
	s.ctx = peer.NewContext(ctx, pr)
	if ht.stats != nil {
		s.ctx = ht.stats.TagRPC(s.ctx, &stats.RPCTagInfo{FullMethodName: s.method})
		inHeader := &stats.InHeader{
			FullMethod:  s.method,
			RemoteAddr:  ht.RemoteAddr(),
			Compression: s.recvCompress,
		}
		ht.stats.HandleRPC(s.ctx, inHeader)
	}
	s.trReader = &transportReader{
		reader:        &recvBufferReader{ctx: s.ctx, ctxDone: s.ctx.Done(), recv: s.buf},
		windowHandler: func(int) {},
	}

	
	readerDone := make(chan struct{})
	go func() {
		defer close(readerDone)

		
		const readSize = 8196
		for buf := make([]byte, readSize); ; {
			n, err := req.Body.Read(buf)
			if n > 0 {
				s.buf.put(recvMsg{data: buf[:n:n]})
				buf = buf[n:]
			}
			if err != nil {
				s.buf.put(recvMsg{err: mapRecvMsgError(err)})
				return
			}
			if len(buf) == 0 {
				buf = make([]byte, readSize)
			}
		}
	}()

	
	
	
	
	startStream(s)

	ht.runStream()
	close(requestOver)

	
	req.Body.Close()
	<-readerDone
}

func (ht *serverHandlerTransport) runStream() {
	for {
		select {
		case fn, ok := <-ht.writes:
			if !ok {
				return
			}
			fn()
		case <-ht.closedCh:
			return
		}
	}
}

func (ht *serverHandlerTransport) IncrMsgSent() {}

func (ht *serverHandlerTransport) IncrMsgRecv() {}

func (ht *serverHandlerTransport) Drain() {
	panic("Drain() is not implemented")
}








func mapRecvMsgError(err error) error {
	if err == io.EOF || err == io.ErrUnexpectedEOF {
		return err
	}
	if se, ok := err.(http2.StreamError); ok {
		if code, ok := http2ErrConvTab[se.Code]; ok {
			return status.Error(code, se.Error())
		}
	}
	if strings.Contains(err.Error(), "body closed by handler") {
		return status.Error(codes.Canceled, err.Error())
	}
	return connectionErrorf(true, err, err.Error())
}
