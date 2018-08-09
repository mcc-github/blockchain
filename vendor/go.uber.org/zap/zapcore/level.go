



















package zapcore

import (
	"bytes"
	"errors"
	"fmt"
)

var errUnmarshalNilLevel = errors.New("can't unmarshal a nil *Level")


type Level int8

const (
	
	
	DebugLevel Level = iota - 1
	
	InfoLevel
	
	
	WarnLevel
	
	
	ErrorLevel
	
	
	DPanicLevel
	
	PanicLevel
	
	FatalLevel

	_minLevel = DebugLevel
	_maxLevel = FatalLevel
)


func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	case DPanicLevel:
		return "dpanic"
	case PanicLevel:
		return "panic"
	case FatalLevel:
		return "fatal"
	default:
		return fmt.Sprintf("Level(%d)", l)
	}
}


func (l Level) CapitalString() string {
	
	
	switch l {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	case DPanicLevel:
		return "DPANIC"
	case PanicLevel:
		return "PANIC"
	case FatalLevel:
		return "FATAL"
	default:
		return fmt.Sprintf("LEVEL(%d)", l)
	}
}



func (l Level) MarshalText() ([]byte, error) {
	return []byte(l.String()), nil
}







func (l *Level) UnmarshalText(text []byte) error {
	if l == nil {
		return errUnmarshalNilLevel
	}
	if !l.unmarshalText(text) && !l.unmarshalText(bytes.ToLower(text)) {
		return fmt.Errorf("unrecognized level: %q", text)
	}
	return nil
}

func (l *Level) unmarshalText(text []byte) bool {
	switch string(text) {
	case "debug", "DEBUG":
		*l = DebugLevel
	case "info", "INFO", "": 
		*l = InfoLevel
	case "warn", "WARN":
		*l = WarnLevel
	case "error", "ERROR":
		*l = ErrorLevel
	case "dpanic", "DPANIC":
		*l = DPanicLevel
	case "panic", "PANIC":
		*l = PanicLevel
	case "fatal", "FATAL":
		*l = FatalLevel
	default:
		return false
	}
	return true
}


func (l *Level) Set(s string) error {
	return l.UnmarshalText([]byte(s))
}


func (l *Level) Get() interface{} {
	return *l
}


func (l Level) Enabled(lvl Level) bool {
	return lvl >= l
}











type LevelEnabler interface {
	Enabled(Level) bool
}
