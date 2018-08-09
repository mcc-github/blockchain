



















package zap

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"sync"

	"go.uber.org/zap/zapcore"
)

const (
	_stdLogDefaultDepth      = 2
	_loggerWriterDepth       = 2
	_programmerErrorTemplate = "You've found a bug in zap! Please file a bug at " +
		"https://github.com/uber-go/zap/issues/new and reference this error: %v"
)

var (
	_globalMu sync.RWMutex
	_globalL  = NewNop()
	_globalS  = _globalL.Sugar()
)



func L() *Logger {
	_globalMu.RLock()
	l := _globalL
	_globalMu.RUnlock()
	return l
}



func S() *SugaredLogger {
	_globalMu.RLock()
	s := _globalS
	_globalMu.RUnlock()
	return s
}



func ReplaceGlobals(logger *Logger) func() {
	_globalMu.Lock()
	prev := _globalL
	_globalL = logger
	_globalS = logger.Sugar()
	_globalMu.Unlock()
	return func() { ReplaceGlobals(prev) }
}




func NewStdLog(l *Logger) *log.Logger {
	logger := l.WithOptions(AddCallerSkip(_stdLogDefaultDepth + _loggerWriterDepth))
	f := logger.Info
	return log.New(&loggerWriter{f}, "" , 0 )
}



func NewStdLogAt(l *Logger, level zapcore.Level) (*log.Logger, error) {
	logger := l.WithOptions(AddCallerSkip(_stdLogDefaultDepth + _loggerWriterDepth))
	logFunc, err := levelToFunc(logger, level)
	if err != nil {
		return nil, err
	}
	return log.New(&loggerWriter{logFunc}, "" , 0 ), nil
}








func RedirectStdLog(l *Logger) func() {
	f, err := redirectStdLogAt(l, InfoLevel)
	if err != nil {
		
		
		panic(fmt.Sprintf(_programmerErrorTemplate, err))
	}
	return f
}








func RedirectStdLogAt(l *Logger, level zapcore.Level) (func(), error) {
	return redirectStdLogAt(l, level)
}

func redirectStdLogAt(l *Logger, level zapcore.Level) (func(), error) {
	flags := log.Flags()
	prefix := log.Prefix()
	log.SetFlags(0)
	log.SetPrefix("")
	logger := l.WithOptions(AddCallerSkip(_stdLogDefaultDepth + _loggerWriterDepth))
	logFunc, err := levelToFunc(logger, level)
	if err != nil {
		return nil, err
	}
	log.SetOutput(&loggerWriter{logFunc})
	return func() {
		log.SetFlags(flags)
		log.SetPrefix(prefix)
		log.SetOutput(os.Stderr)
	}, nil
}

func levelToFunc(logger *Logger, lvl zapcore.Level) (func(string, ...Field), error) {
	switch lvl {
	case DebugLevel:
		return logger.Debug, nil
	case InfoLevel:
		return logger.Info, nil
	case WarnLevel:
		return logger.Warn, nil
	case ErrorLevel:
		return logger.Error, nil
	case DPanicLevel:
		return logger.DPanic, nil
	case PanicLevel:
		return logger.Panic, nil
	case FatalLevel:
		return logger.Fatal, nil
	}
	return nil, fmt.Errorf("unrecognized level: %q", lvl)
}

type loggerWriter struct {
	logFunc func(msg string, fields ...Field)
}

func (l *loggerWriter) Write(p []byte) (int, error) {
	p = bytes.TrimSpace(p)
	l.logFunc(string(p))
	return len(p), nil
}
