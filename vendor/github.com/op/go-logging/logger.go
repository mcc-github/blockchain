






package logging

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"sync/atomic"
	"time"
)




type Redactor interface {
	Redacted() interface{}
}


func Redact(s string) string {
	return strings.Repeat("*", len(s))
}

var (
	
	sequenceNo uint64

	
	timeNow = time.Now
)




type Record struct {
	ID     uint64
	Time   time.Time
	Module string
	Level  Level
	Args   []interface{}

	
	
	message   *string
	fmt       *string
	formatter Formatter
	formatted string
}


func (r *Record) Formatted(calldepth int) string {
	if r.formatted == "" {
		var buf bytes.Buffer
		r.formatter.Format(calldepth+1, r, &buf)
		r.formatted = buf.String()
	}
	return r.formatted
}


func (r *Record) Message() string {
	if r.message == nil {
		
		for i, arg := range r.Args {
			if redactor, ok := arg.(Redactor); ok == true {
				r.Args[i] = redactor.Redacted()
			}
		}
		var buf bytes.Buffer
		if r.fmt != nil {
			fmt.Fprintf(&buf, *r.fmt, r.Args...)
		} else {
			
			fmt.Fprintln(&buf, r.Args...)
			buf.Truncate(buf.Len() - 1) 
		}
		msg := buf.String()
		r.message = &msg
	}
	return *r.message
}



type Logger struct {
	Module      string
	backend     LeveledBackend
	haveBackend bool

	
	
	ExtraCalldepth int
}


func (l *Logger) SetBackend(backend LeveledBackend) {
	l.backend = backend
	l.haveBackend = true
}




func GetLogger(module string) (*Logger, error) {
	return &Logger{Module: module}, nil
}



func MustGetLogger(module string) *Logger {
	logger, err := GetLogger(module)
	if err != nil {
		panic("logger: " + module + ": " + err.Error())
	}
	return logger
}


func Reset() {
	
	
	
	sequenceNo = 0
	b := SetBackend(NewLogBackend(os.Stderr, "", log.LstdFlags))
	b.SetLevel(DEBUG, "")
	SetFormatter(DefaultFormatter)
	timeNow = time.Now
}


func (l *Logger) IsEnabledFor(level Level) bool {
	return defaultBackend.IsEnabledFor(level, l.Module)
}

func (l *Logger) log(lvl Level, format *string, args ...interface{}) {
	if !l.IsEnabledFor(lvl) {
		return
	}

	
	record := &Record{
		ID:     atomic.AddUint64(&sequenceNo, 1),
		Time:   timeNow(),
		Module: l.Module,
		Level:  lvl,
		fmt:    format,
		Args:   args,
	}

	
	

	
	
	
	
	if l.haveBackend {
		l.backend.Log(lvl, 2+l.ExtraCalldepth, record)
		return
	}

	defaultBackend.Log(lvl, 2+l.ExtraCalldepth, record)
}


func (l *Logger) Fatal(args ...interface{}) {
	l.log(CRITICAL, nil, args...)
	os.Exit(1)
}


func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.log(CRITICAL, &format, args...)
	os.Exit(1)
}


func (l *Logger) Panic(args ...interface{}) {
	l.log(CRITICAL, nil, args...)
	panic(fmt.Sprint(args...))
}


func (l *Logger) Panicf(format string, args ...interface{}) {
	l.log(CRITICAL, &format, args...)
	panic(fmt.Sprintf(format, args...))
}


func (l *Logger) Critical(args ...interface{}) {
	l.log(CRITICAL, nil, args...)
}


func (l *Logger) Criticalf(format string, args ...interface{}) {
	l.log(CRITICAL, &format, args...)
}


func (l *Logger) Error(args ...interface{}) {
	l.log(ERROR, nil, args...)
}


func (l *Logger) Errorf(format string, args ...interface{}) {
	l.log(ERROR, &format, args...)
}


func (l *Logger) Warning(args ...interface{}) {
	l.log(WARNING, nil, args...)
}


func (l *Logger) Warningf(format string, args ...interface{}) {
	l.log(WARNING, &format, args...)
}


func (l *Logger) Notice(args ...interface{}) {
	l.log(NOTICE, nil, args...)
}


func (l *Logger) Noticef(format string, args ...interface{}) {
	l.log(NOTICE, &format, args...)
}


func (l *Logger) Info(args ...interface{}) {
	l.log(INFO, nil, args...)
}


func (l *Logger) Infof(format string, args ...interface{}) {
	l.log(INFO, &format, args...)
}


func (l *Logger) Debug(args ...interface{}) {
	l.log(DEBUG, nil, args...)
}


func (l *Logger) Debugf(format string, args ...interface{}) {
	l.log(DEBUG, &format, args...)
}

func init() {
	Reset()
}
