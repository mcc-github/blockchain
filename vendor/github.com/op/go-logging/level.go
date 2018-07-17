



package logging

import (
	"errors"
	"strings"
	"sync"
)


var ErrInvalidLogLevel = errors.New("logger: invalid log level")


type Level int


const (
	CRITICAL Level = iota
	ERROR
	WARNING
	NOTICE
	INFO
	DEBUG
)

var levelNames = []string{
	"CRITICAL",
	"ERROR",
	"WARNING",
	"NOTICE",
	"INFO",
	"DEBUG",
}


func (p Level) String() string {
	return levelNames[p]
}


func LogLevel(level string) (Level, error) {
	for i, name := range levelNames {
		if strings.EqualFold(name, level) {
			return Level(i), nil
		}
	}
	return ERROR, ErrInvalidLogLevel
}



type Leveled interface {
	GetLevel(string) Level
	SetLevel(Level, string)
	IsEnabledFor(Level, string) bool
}



type LeveledBackend interface {
	Backend
	Leveled
}

type moduleLeveled struct {
	levels    map[string]Level
	backend   Backend
	formatter Formatter
	once      sync.Once
}



func AddModuleLevel(backend Backend) LeveledBackend {
	var leveled LeveledBackend
	var ok bool
	if leveled, ok = backend.(LeveledBackend); !ok {
		leveled = &moduleLeveled{
			levels:  make(map[string]Level),
			backend: backend,
		}
	}
	return leveled
}


func (l *moduleLeveled) GetLevel(module string) Level {
	level, exists := l.levels[module]
	if exists == false {
		level, exists = l.levels[""]
		
		if exists == false {
			level = DEBUG
		}
	}
	return level
}


func (l *moduleLeveled) SetLevel(level Level, module string) {
	l.levels[module] = level
}


func (l *moduleLeveled) IsEnabledFor(level Level, module string) bool {
	return level <= l.GetLevel(module)
}

func (l *moduleLeveled) Log(level Level, calldepth int, rec *Record) (err error) {
	if l.IsEnabledFor(level, rec.Module) {
		
		rec.formatter = l.getFormatterAndCacheCurrent()
		err = l.backend.Log(level, calldepth+1, rec)
	}
	return
}

func (l *moduleLeveled) getFormatterAndCacheCurrent() Formatter {
	l.once.Do(func() {
		if l.formatter == nil {
			l.formatter = getFormatter()
		}
	})
	return l.formatter
}
