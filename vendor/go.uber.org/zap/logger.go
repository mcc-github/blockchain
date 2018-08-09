



















package zap

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"time"

	"go.uber.org/zap/zapcore"
)








type Logger struct {
	core zapcore.Core

	development bool
	name        string
	errorOutput zapcore.WriteSyncer

	addCaller bool
	addStack  zapcore.LevelEnabler

	callerSkip int
}











func New(core zapcore.Core, options ...Option) *Logger {
	if core == nil {
		return NewNop()
	}
	log := &Logger{
		core:        core,
		errorOutput: zapcore.Lock(os.Stderr),
		addStack:    zapcore.FatalLevel + 1,
	}
	return log.WithOptions(options...)
}






func NewNop() *Logger {
	return &Logger{
		core:        zapcore.NewNopCore(),
		errorOutput: zapcore.AddSync(ioutil.Discard),
		addStack:    zapcore.FatalLevel + 1,
	}
}





func NewProduction(options ...Option) (*Logger, error) {
	return NewProductionConfig().Build(options...)
}





func NewDevelopment(options ...Option) (*Logger, error) {
	return NewDevelopmentConfig().Build(options...)
}





func NewExample(options ...Option) *Logger {
	encoderCfg := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		NameKey:        "logger",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
	}
	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), os.Stdout, DebugLevel)
	return New(core).WithOptions(options...)
}





func (log *Logger) Sugar() *SugaredLogger {
	core := log.clone()
	core.callerSkip += 2
	return &SugaredLogger{core}
}



func (log *Logger) Named(s string) *Logger {
	if s == "" {
		return log
	}
	l := log.clone()
	if log.name == "" {
		l.name = s
	} else {
		l.name = strings.Join([]string{l.name, s}, ".")
	}
	return l
}



func (log *Logger) WithOptions(opts ...Option) *Logger {
	c := log.clone()
	for _, opt := range opts {
		opt.apply(c)
	}
	return c
}



func (log *Logger) With(fields ...Field) *Logger {
	if len(fields) == 0 {
		return log
	}
	l := log.clone()
	l.core = l.core.With(fields)
	return l
}




func (log *Logger) Check(lvl zapcore.Level, msg string) *zapcore.CheckedEntry {
	return log.check(lvl, msg)
}



func (log *Logger) Debug(msg string, fields ...Field) {
	if ce := log.check(DebugLevel, msg); ce != nil {
		ce.Write(fields...)
	}
}



func (log *Logger) Info(msg string, fields ...Field) {
	if ce := log.check(InfoLevel, msg); ce != nil {
		ce.Write(fields...)
	}
}



func (log *Logger) Warn(msg string, fields ...Field) {
	if ce := log.check(WarnLevel, msg); ce != nil {
		ce.Write(fields...)
	}
}



func (log *Logger) Error(msg string, fields ...Field) {
	if ce := log.check(ErrorLevel, msg); ce != nil {
		ce.Write(fields...)
	}
}







func (log *Logger) DPanic(msg string, fields ...Field) {
	if ce := log.check(DPanicLevel, msg); ce != nil {
		ce.Write(fields...)
	}
}





func (log *Logger) Panic(msg string, fields ...Field) {
	if ce := log.check(PanicLevel, msg); ce != nil {
		ce.Write(fields...)
	}
}






func (log *Logger) Fatal(msg string, fields ...Field) {
	if ce := log.check(FatalLevel, msg); ce != nil {
		ce.Write(fields...)
	}
}



func (log *Logger) Sync() error {
	return log.core.Sync()
}


func (log *Logger) Core() zapcore.Core {
	return log.core
}

func (log *Logger) clone() *Logger {
	copy := *log
	return &copy
}

func (log *Logger) check(lvl zapcore.Level, msg string) *zapcore.CheckedEntry {
	
	
	const callerSkipOffset = 2

	
	
	ent := zapcore.Entry{
		LoggerName: log.name,
		Time:       time.Now(),
		Level:      lvl,
		Message:    msg,
	}
	ce := log.core.Check(ent, nil)
	willWrite := ce != nil

	
	switch ent.Level {
	case zapcore.PanicLevel:
		ce = ce.Should(ent, zapcore.WriteThenPanic)
	case zapcore.FatalLevel:
		ce = ce.Should(ent, zapcore.WriteThenFatal)
	case zapcore.DPanicLevel:
		if log.development {
			ce = ce.Should(ent, zapcore.WriteThenPanic)
		}
	}

	
	
	
	if !willWrite {
		return ce
	}

	
	ce.ErrorOutput = log.errorOutput
	if log.addCaller {
		ce.Entry.Caller = zapcore.NewEntryCaller(runtime.Caller(log.callerSkip + callerSkipOffset))
		if !ce.Entry.Caller.Defined {
			fmt.Fprintf(log.errorOutput, "%v Logger.check error: failed to get caller\n", time.Now().UTC())
			log.errorOutput.Sync()
		}
	}
	if log.addStack.Enabled(ce.Entry.Level) {
		ce.Entry.Stack = Stack("").String
	}

	return ce
}
