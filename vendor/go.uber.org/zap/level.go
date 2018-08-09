



















package zap

import (
	"go.uber.org/atomic"
	"go.uber.org/zap/zapcore"
)

const (
	
	
	DebugLevel = zapcore.DebugLevel
	
	InfoLevel = zapcore.InfoLevel
	
	
	WarnLevel = zapcore.WarnLevel
	
	
	ErrorLevel = zapcore.ErrorLevel
	
	
	DPanicLevel = zapcore.DPanicLevel
	
	PanicLevel = zapcore.PanicLevel
	
	FatalLevel = zapcore.FatalLevel
)







type LevelEnablerFunc func(zapcore.Level) bool


func (f LevelEnablerFunc) Enabled(lvl zapcore.Level) bool { return f(lvl) }










type AtomicLevel struct {
	l *atomic.Int32
}



func NewAtomicLevel() AtomicLevel {
	return AtomicLevel{
		l: atomic.NewInt32(int32(InfoLevel)),
	}
}



func NewAtomicLevelAt(l zapcore.Level) AtomicLevel {
	a := NewAtomicLevel()
	a.SetLevel(l)
	return a
}



func (lvl AtomicLevel) Enabled(l zapcore.Level) bool {
	return lvl.Level().Enabled(l)
}


func (lvl AtomicLevel) Level() zapcore.Level {
	return zapcore.Level(int8(lvl.l.Load()))
}


func (lvl AtomicLevel) SetLevel(l zapcore.Level) {
	lvl.l.Store(int32(l))
}


func (lvl AtomicLevel) String() string {
	return lvl.Level().String()
}




func (lvl *AtomicLevel) UnmarshalText(text []byte) error {
	if lvl.l == nil {
		lvl.l = &atomic.Int32{}
	}

	var l zapcore.Level
	if err := l.UnmarshalText(text); err != nil {
		return err
	}

	lvl.SetLevel(l)
	return nil
}




func (lvl AtomicLevel) MarshalText() (text []byte, err error) {
	return lvl.Level().MarshalText()
}
