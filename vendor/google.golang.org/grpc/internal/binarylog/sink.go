

package binarylog

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	pb "google.golang.org/grpc/binarylog/grpc_binarylog_v1"
	"google.golang.org/grpc/grpclog"
)

var (
	defaultSink Sink = &noopSink{} 
)




func SetDefaultSink(s Sink) {
	if defaultSink != nil {
		defaultSink.Close()
	}
	defaultSink = s
}


type Sink interface {
	
	
	
	Write(*pb.GrpcLogEntry) error
	
	Close() error
}

type noopSink struct{}

func (ns *noopSink) Write(*pb.GrpcLogEntry) error { return nil }
func (ns *noopSink) Close() error                 { return nil }







func newWriterSink(w io.Writer) *writerSink {
	return &writerSink{out: w}
}

type writerSink struct {
	out io.Writer
}

func (ws *writerSink) Write(e *pb.GrpcLogEntry) error {
	b, err := proto.Marshal(e)
	if err != nil {
		grpclog.Infof("binary logging: failed to marshal proto message: %v", err)
	}
	hdr := make([]byte, 4)
	binary.BigEndian.PutUint32(hdr, uint32(len(b)))
	if _, err := ws.out.Write(hdr); err != nil {
		return err
	}
	if _, err := ws.out.Write(b); err != nil {
		return err
	}
	return nil
}

func (ws *writerSink) Close() error { return nil }

type bufWriteCloserSink struct {
	mu     sync.Mutex
	closer io.Closer
	out    *writerSink   
	buf    *bufio.Writer 

	writeStartOnce sync.Once
	writeTicker    *time.Ticker
}

func (fs *bufWriteCloserSink) Write(e *pb.GrpcLogEntry) error {
	
	fs.writeStartOnce.Do(fs.startFlushGoroutine)
	fs.mu.Lock()
	if err := fs.out.Write(e); err != nil {
		fs.mu.Unlock()
		return err
	}
	fs.mu.Unlock()
	return nil
}

const (
	bufFlushDuration = 60 * time.Second
)

func (fs *bufWriteCloserSink) startFlushGoroutine() {
	fs.writeTicker = time.NewTicker(bufFlushDuration)
	go func() {
		for range fs.writeTicker.C {
			fs.mu.Lock()
			fs.buf.Flush()
			fs.mu.Unlock()
		}
	}()
}

func (fs *bufWriteCloserSink) Close() error {
	if fs.writeTicker != nil {
		fs.writeTicker.Stop()
	}
	fs.mu.Lock()
	fs.buf.Flush()
	fs.closer.Close()
	fs.out.Close()
	fs.mu.Unlock()
	return nil
}

func newBufWriteCloserSink(o io.WriteCloser) Sink {
	bufW := bufio.NewWriter(o)
	return &bufWriteCloserSink{
		closer: o,
		out:    newWriterSink(bufW),
		buf:    bufW,
	}
}



func NewTempFileSink() (Sink, error) {
	tempFile, err := ioutil.TempFile("/tmp", "grpcgo_binarylog_*.txt")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %v", err)
	}
	return newBufWriteCloserSink(tempFile), nil
}
