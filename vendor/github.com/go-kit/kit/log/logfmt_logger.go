package log

import (
	"bytes"
	"io"
	"sync"

	"github.com/go-logfmt/logfmt"
)

type logfmtEncoder struct {
	*logfmt.Encoder
	buf bytes.Buffer
}

func (l *logfmtEncoder) Reset() {
	l.Encoder.Reset()
	l.buf.Reset()
}

var logfmtEncoderPool = sync.Pool{
	New: func() interface{} {
		var enc logfmtEncoder
		enc.Encoder = logfmt.NewEncoder(&enc.buf)
		return &enc
	},
}

type logfmtLogger struct {
	w io.Writer
}





func NewLogfmtLogger(w io.Writer) Logger {
	return &logfmtLogger{w}
}

func (l logfmtLogger) Log(keyvals ...interface{}) error {
	enc := logfmtEncoderPool.Get().(*logfmtEncoder)
	enc.Reset()
	defer logfmtEncoderPool.Put(enc)

	if err := enc.EncodeKeyvals(keyvals...); err != nil {
		return err
	}

	
	if err := enc.EndRecord(); err != nil {
		return err
	}

	
	
	
	if _, err := l.w.Write(enc.buf.Bytes()); err != nil {
		return err
	}
	return nil
}
