package log

import (
	"io"
	"sync"
	"sync/atomic"
)







type SwapLogger struct {
	logger atomic.Value
}

type loggerStruct struct {
	Logger
}



func (l *SwapLogger) Log(keyvals ...interface{}) error {
	s, ok := l.logger.Load().(loggerStruct)
	if !ok || s.Logger == nil {
		return nil
	}
	return s.Log(keyvals...)
}



func (l *SwapLogger) Swap(logger Logger) {
	l.logger.Store(loggerStruct{logger})
}











func NewSyncWriter(w io.Writer) io.Writer {
	switch w := w.(type) {
	case fdWriter:
		return &fdSyncWriter{fdWriter: w}
	default:
		return &syncWriter{Writer: w}
	}
}


type syncWriter struct {
	sync.Mutex
	io.Writer
}



func (w *syncWriter) Write(p []byte) (n int, err error) {
	w.Lock()
	n, err = w.Writer.Write(p)
	w.Unlock()
	return n, err
}



type fdWriter interface {
	io.Writer
	Fd() uintptr
}


type fdSyncWriter struct {
	sync.Mutex
	fdWriter
}



func (w *fdSyncWriter) Write(p []byte) (n int, err error) {
	w.Lock()
	n, err = w.fdWriter.Write(p)
	w.Unlock()
	return n, err
}


type syncLogger struct {
	mu     sync.Mutex
	logger Logger
}





func NewSyncLogger(logger Logger) Logger {
	return &syncLogger{logger: logger}
}



func (l *syncLogger) Log(keyvals ...interface{}) error {
	l.mu.Lock()
	err := l.logger.Log(keyvals...)
	l.mu.Unlock()
	return err
}
