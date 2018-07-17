

package grpc

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	"golang.org/x/net/trace"
)



var EnableTracing bool



func methodFamily(m string) string {
	m = strings.TrimPrefix(m, "/") 
	if i := strings.Index(m, "/"); i >= 0 {
		m = m[:i] 
	}
	if i := strings.LastIndex(m, "."); i >= 0 {
		m = m[i+1:] 
	}
	return m
}


type traceInfo struct {
	tr        trace.Trace
	firstLine firstLine
}


type firstLine struct {
	client     bool 
	remoteAddr net.Addr
	deadline   time.Duration 
}

func (f *firstLine) String() string {
	var line bytes.Buffer
	io.WriteString(&line, "RPC: ")
	if f.client {
		io.WriteString(&line, "to")
	} else {
		io.WriteString(&line, "from")
	}
	fmt.Fprintf(&line, " %v deadline:", f.remoteAddr)
	if f.deadline != 0 {
		fmt.Fprint(&line, f.deadline)
	} else {
		io.WriteString(&line, "none")
	}
	return line.String()
}

const truncateSize = 100

func truncate(x string, l int) string {
	if l > len(x) {
		return x
	}
	return x[:l]
}


type payload struct {
	sent bool        
	msg  interface{} 
	
}

func (p payload) String() string {
	if p.sent {
		return truncate(fmt.Sprintf("sent: %v", p.msg), truncateSize)
	}
	return truncate(fmt.Sprintf("recv: %v", p.msg), truncateSize)
}

type fmtStringer struct {
	format string
	a      []interface{}
}

func (f *fmtStringer) String() string {
	return fmt.Sprintf(f.format, f.a...)
}

type stringer string

func (s stringer) String() string { return string(s) }
