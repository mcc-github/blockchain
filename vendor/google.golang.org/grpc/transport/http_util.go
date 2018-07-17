

package transport

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/hpack"
	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	
	http2MaxFrameLen = 16384 
	
	http2InitHeaderTableSize = 4096
	
	defaultWriteBufSize = 32 * 1024
	defaultReadBufSize  = 32 * 1024
	
	
	
	
	
	baseContentType = "application/grpc"
)

var (
	clientPreface   = []byte(http2.ClientPreface)
	http2ErrConvTab = map[http2.ErrCode]codes.Code{
		http2.ErrCodeNo:                 codes.Internal,
		http2.ErrCodeProtocol:           codes.Internal,
		http2.ErrCodeInternal:           codes.Internal,
		http2.ErrCodeFlowControl:        codes.ResourceExhausted,
		http2.ErrCodeSettingsTimeout:    codes.Internal,
		http2.ErrCodeStreamClosed:       codes.Internal,
		http2.ErrCodeFrameSize:          codes.Internal,
		http2.ErrCodeRefusedStream:      codes.Unavailable,
		http2.ErrCodeCancel:             codes.Canceled,
		http2.ErrCodeCompression:        codes.Internal,
		http2.ErrCodeConnect:            codes.Internal,
		http2.ErrCodeEnhanceYourCalm:    codes.ResourceExhausted,
		http2.ErrCodeInadequateSecurity: codes.PermissionDenied,
		http2.ErrCodeHTTP11Required:     codes.Internal,
	}
	statusCodeConvTab = map[codes.Code]http2.ErrCode{
		codes.Internal:          http2.ErrCodeInternal,
		codes.Canceled:          http2.ErrCodeCancel,
		codes.Unavailable:       http2.ErrCodeRefusedStream,
		codes.ResourceExhausted: http2.ErrCodeEnhanceYourCalm,
		codes.PermissionDenied:  http2.ErrCodeInadequateSecurity,
	}
	httpStatusConvTab = map[int]codes.Code{
		
		http.StatusBadRequest: codes.Internal,
		
		http.StatusUnauthorized: codes.Unauthenticated,
		
		http.StatusForbidden: codes.PermissionDenied,
		
		http.StatusNotFound: codes.Unimplemented,
		
		http.StatusTooManyRequests: codes.Unavailable,
		
		http.StatusBadGateway: codes.Unavailable,
		
		http.StatusServiceUnavailable: codes.Unavailable,
		
		http.StatusGatewayTimeout: codes.Unavailable,
	}
)



type decodeState struct {
	encoding string
	
	
	
	statusGen *status.Status
	
	
	rawStatusCode *int
	rawStatusMsg  string
	httpStatus    *int
	
	timeoutSet bool
	timeout    time.Duration
	method     string
	
	mdata          map[string][]string
	statsTags      []byte
	statsTrace     []byte
	contentSubtype string
}




func isReservedHeader(hdr string) bool {
	if hdr != "" && hdr[0] == ':' {
		return true
	}
	switch hdr {
	case "content-type",
		"grpc-message-type",
		"grpc-encoding",
		"grpc-message",
		"grpc-status",
		"grpc-timeout",
		"grpc-status-details-bin",
		"te":
		return true
	default:
		return false
	}
}



func isWhitelistedPseudoHeader(hdr string) bool {
	switch hdr {
	case ":authority":
		return true
	default:
		return false
	}
}














func contentSubtype(contentType string) (string, bool) {
	if contentType == baseContentType {
		return "", true
	}
	if !strings.HasPrefix(contentType, baseContentType) {
		return "", false
	}
	
	switch contentType[len(baseContentType)] {
	case '+', ';':
		
		
		
		return contentType[len(baseContentType)+1:], true
	default:
		return "", false
	}
}


func contentType(contentSubtype string) string {
	if contentSubtype == "" {
		return baseContentType
	}
	return baseContentType + "+" + contentSubtype
}

func (d *decodeState) status() *status.Status {
	if d.statusGen == nil {
		
		d.statusGen = status.New(codes.Code(int32(*(d.rawStatusCode))), d.rawStatusMsg)
	}
	return d.statusGen
}

const binHdrSuffix = "-bin"

func encodeBinHeader(v []byte) string {
	return base64.RawStdEncoding.EncodeToString(v)
}

func decodeBinHeader(v string) ([]byte, error) {
	if len(v)%4 == 0 {
		
		return base64.StdEncoding.DecodeString(v)
	}
	return base64.RawStdEncoding.DecodeString(v)
}

func encodeMetadataHeader(k, v string) string {
	if strings.HasSuffix(k, binHdrSuffix) {
		return encodeBinHeader(([]byte)(v))
	}
	return v
}

func decodeMetadataHeader(k, v string) (string, error) {
	if strings.HasSuffix(k, binHdrSuffix) {
		b, err := decodeBinHeader(v)
		return string(b), err
	}
	return v, nil
}

func (d *decodeState) decodeResponseHeader(frame *http2.MetaHeadersFrame) error {
	for _, hf := range frame.Fields {
		if err := d.processHeaderField(hf); err != nil {
			return err
		}
	}

	
	if d.rawStatusCode != nil || d.statusGen != nil {
		return nil
	}

	
	
	if d.httpStatus == nil {
		return streamErrorf(codes.Internal, "malformed header: doesn't contain status(gRPC or HTTP)")
	}

	if *(d.httpStatus) != http.StatusOK {
		code, ok := httpStatusConvTab[*(d.httpStatus)]
		if !ok {
			code = codes.Unknown
		}
		return streamErrorf(code, http.StatusText(*(d.httpStatus)))
	}

	
	
	
	
	
	
	code := int(codes.Unknown)
	d.rawStatusCode = &code
	return nil

}

func (d *decodeState) addMetadata(k, v string) {
	if d.mdata == nil {
		d.mdata = make(map[string][]string)
	}
	d.mdata[k] = append(d.mdata[k], v)
}

func (d *decodeState) processHeaderField(f hpack.HeaderField) error {
	switch f.Name {
	case "content-type":
		contentSubtype, validContentType := contentSubtype(f.Value)
		if !validContentType {
			return streamErrorf(codes.Internal, "transport: received the unexpected content-type %q", f.Value)
		}
		d.contentSubtype = contentSubtype
		
		
		
		
		d.addMetadata(f.Name, f.Value)
	case "grpc-encoding":
		d.encoding = f.Value
	case "grpc-status":
		code, err := strconv.Atoi(f.Value)
		if err != nil {
			return streamErrorf(codes.Internal, "transport: malformed grpc-status: %v", err)
		}
		d.rawStatusCode = &code
	case "grpc-message":
		d.rawStatusMsg = decodeGrpcMessage(f.Value)
	case "grpc-status-details-bin":
		v, err := decodeBinHeader(f.Value)
		if err != nil {
			return streamErrorf(codes.Internal, "transport: malformed grpc-status-details-bin: %v", err)
		}
		s := &spb.Status{}
		if err := proto.Unmarshal(v, s); err != nil {
			return streamErrorf(codes.Internal, "transport: malformed grpc-status-details-bin: %v", err)
		}
		d.statusGen = status.FromProto(s)
	case "grpc-timeout":
		d.timeoutSet = true
		var err error
		if d.timeout, err = decodeTimeout(f.Value); err != nil {
			return streamErrorf(codes.Internal, "transport: malformed time-out: %v", err)
		}
	case ":path":
		d.method = f.Value
	case ":status":
		code, err := strconv.Atoi(f.Value)
		if err != nil {
			return streamErrorf(codes.Internal, "transport: malformed http-status: %v", err)
		}
		d.httpStatus = &code
	case "grpc-tags-bin":
		v, err := decodeBinHeader(f.Value)
		if err != nil {
			return streamErrorf(codes.Internal, "transport: malformed grpc-tags-bin: %v", err)
		}
		d.statsTags = v
		d.addMetadata(f.Name, string(v))
	case "grpc-trace-bin":
		v, err := decodeBinHeader(f.Value)
		if err != nil {
			return streamErrorf(codes.Internal, "transport: malformed grpc-trace-bin: %v", err)
		}
		d.statsTrace = v
		d.addMetadata(f.Name, string(v))
	default:
		if isReservedHeader(f.Name) && !isWhitelistedPseudoHeader(f.Name) {
			break
		}
		v, err := decodeMetadataHeader(f.Name, f.Value)
		if err != nil {
			errorf("Failed to decode metadata header (%q, %q): %v", f.Name, f.Value, err)
			return nil
		}
		d.addMetadata(f.Name, string(v))
	}
	return nil
}

type timeoutUnit uint8

const (
	hour        timeoutUnit = 'H'
	minute      timeoutUnit = 'M'
	second      timeoutUnit = 'S'
	millisecond timeoutUnit = 'm'
	microsecond timeoutUnit = 'u'
	nanosecond  timeoutUnit = 'n'
)

func timeoutUnitToDuration(u timeoutUnit) (d time.Duration, ok bool) {
	switch u {
	case hour:
		return time.Hour, true
	case minute:
		return time.Minute, true
	case second:
		return time.Second, true
	case millisecond:
		return time.Millisecond, true
	case microsecond:
		return time.Microsecond, true
	case nanosecond:
		return time.Nanosecond, true
	default:
	}
	return
}

const maxTimeoutValue int64 = 100000000 - 1



func div(d, r time.Duration) int64 {
	if m := d % r; m > 0 {
		return int64(d/r + 1)
	}
	return int64(d / r)
}


func encodeTimeout(t time.Duration) string {
	if t <= 0 {
		return "0n"
	}
	if d := div(t, time.Nanosecond); d <= maxTimeoutValue {
		return strconv.FormatInt(d, 10) + "n"
	}
	if d := div(t, time.Microsecond); d <= maxTimeoutValue {
		return strconv.FormatInt(d, 10) + "u"
	}
	if d := div(t, time.Millisecond); d <= maxTimeoutValue {
		return strconv.FormatInt(d, 10) + "m"
	}
	if d := div(t, time.Second); d <= maxTimeoutValue {
		return strconv.FormatInt(d, 10) + "S"
	}
	if d := div(t, time.Minute); d <= maxTimeoutValue {
		return strconv.FormatInt(d, 10) + "M"
	}
	
	return strconv.FormatInt(div(t, time.Hour), 10) + "H"
}

func decodeTimeout(s string) (time.Duration, error) {
	size := len(s)
	if size < 2 {
		return 0, fmt.Errorf("transport: timeout string is too short: %q", s)
	}
	unit := timeoutUnit(s[size-1])
	d, ok := timeoutUnitToDuration(unit)
	if !ok {
		return 0, fmt.Errorf("transport: timeout unit is not recognized: %q", s)
	}
	t, err := strconv.ParseInt(s[:size-1], 10, 64)
	if err != nil {
		return 0, err
	}
	return d * time.Duration(t), nil
}

const (
	spaceByte   = ' '
	tildaByte   = '~'
	percentByte = '%'
)







func encodeGrpcMessage(msg string) string {
	if msg == "" {
		return ""
	}
	lenMsg := len(msg)
	for i := 0; i < lenMsg; i++ {
		c := msg[i]
		if !(c >= spaceByte && c < tildaByte && c != percentByte) {
			return encodeGrpcMessageUnchecked(msg)
		}
	}
	return msg
}

func encodeGrpcMessageUnchecked(msg string) string {
	var buf bytes.Buffer
	lenMsg := len(msg)
	for i := 0; i < lenMsg; i++ {
		c := msg[i]
		if c >= spaceByte && c < tildaByte && c != percentByte {
			buf.WriteByte(c)
		} else {
			buf.WriteString(fmt.Sprintf("%%%02X", c))
		}
	}
	return buf.String()
}


func decodeGrpcMessage(msg string) string {
	if msg == "" {
		return ""
	}
	lenMsg := len(msg)
	for i := 0; i < lenMsg; i++ {
		if msg[i] == percentByte && i+2 < lenMsg {
			return decodeGrpcMessageUnchecked(msg)
		}
	}
	return msg
}

func decodeGrpcMessageUnchecked(msg string) string {
	var buf bytes.Buffer
	lenMsg := len(msg)
	for i := 0; i < lenMsg; i++ {
		c := msg[i]
		if c == percentByte && i+2 < lenMsg {
			parsed, err := strconv.ParseUint(msg[i+1:i+3], 16, 8)
			if err != nil {
				buf.WriteByte(c)
			} else {
				buf.WriteByte(byte(parsed))
				i += 2
			}
		} else {
			buf.WriteByte(c)
		}
	}
	return buf.String()
}

type framer struct {
	numWriters int32
	reader     io.Reader
	writer     *bufio.Writer
	fr         *http2.Framer
}

func newFramer(conn net.Conn, writeBufferSize, readBufferSize int) *framer {
	f := &framer{
		reader: bufio.NewReaderSize(conn, readBufferSize),
		writer: bufio.NewWriterSize(conn, writeBufferSize),
	}
	f.fr = http2.NewFramer(f.writer, f.reader)
	
	
	f.fr.SetReuseFrames()
	f.fr.ReadMetaHeaders = hpack.NewDecoder(http2InitHeaderTableSize, nil)
	return f
}
