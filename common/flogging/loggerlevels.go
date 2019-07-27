/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package flogging

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"go.uber.org/zap/zapcore"
)


type LoggerLevels struct {
	mutex        sync.RWMutex
	levelCache   map[string]zapcore.Level
	specs        map[string]zapcore.Level
	defaultLevel zapcore.Level
	minLevel     zapcore.Level
}



func (l *LoggerLevels) DefaultLevel() zapcore.Level {
	l.mutex.RLock()
	lvl := l.defaultLevel
	l.mutex.RUnlock()
	return lvl
}





func (l *LoggerLevels) ActivateSpec(spec string) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	defaultLevel := zapcore.InfoLevel
	specs := map[string]zapcore.Level{}
	for _, field := range strings.Split(spec, ":") {
		split := strings.Split(field, "=")
		switch len(split) {
		case 1: 
			if field != "" && !IsValidLevel(field) {
				return errors.Errorf("invalid logging specification '%s': bad segment '%s'", spec, field)
			}
			defaultLevel = NameToLevel(field)

		case 2: 
			if split[0] == "" {
				return errors.Errorf("invalid logging specification '%s': no logger specified in segment '%s'", spec, field)
			}
			if field != "" && !IsValidLevel(split[1]) {
				return errors.Errorf("invalid logging specification '%s': bad segment '%s'", spec, field)
			}

			level := NameToLevel(split[1])
			loggers := strings.Split(split[0], ",")
			for _, logger := range loggers {
				
				
				
				
				if !isValidLoggerName(strings.TrimSuffix(logger, ".")) {
					return errors.Errorf("invalid logging specification '%s': bad logger name '%s'", spec, logger)
				}
				specs[logger] = level
			}

		default:
			return errors.Errorf("invalid logging specification '%s': bad segment '%s'", spec, field)
		}
	}

	minLevel := defaultLevel
	for _, lvl := range specs {
		if lvl < minLevel {
			minLevel = lvl
		}
	}

	l.minLevel = minLevel
	l.defaultLevel = defaultLevel
	l.specs = specs
	l.levelCache = map[string]zapcore.Level{}

	return nil
}


var loggerNameRegexp = regexp.MustCompile(`^[[:alnum:]_#:-]+(\.[[:alnum:]_#:-]+)*$`)





func isValidLoggerName(loggerName string) bool {
	return loggerNameRegexp.MatchString(loggerName)
}




func (l *LoggerLevels) Level(loggerName string) zapcore.Level {
	if level, ok := l.cachedLevel(loggerName); ok {
		return level
	}

	l.mutex.Lock()
	level := l.calculateLevel(loggerName)
	l.levelCache[loggerName] = level
	l.mutex.Unlock()

	return level
}



func (l *LoggerLevels) calculateLevel(loggerName string) zapcore.Level {
	candidate := loggerName + "."
	for {
		if lvl, ok := l.specs[candidate]; ok {
			return lvl
		}

		idx := strings.LastIndex(candidate, ".")
		if idx <= 0 {
			return l.defaultLevel
		}
		candidate = candidate[:idx]
	}
}



func (l *LoggerLevels) cachedLevel(loggerName string) (lvl zapcore.Level, ok bool) {
	l.mutex.RLock()
	level, ok := l.levelCache[loggerName]
	l.mutex.RUnlock()
	return level, ok
}


func (l *LoggerLevels) Spec() string {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	var fields []string
	for k, v := range l.specs {
		fields = append(fields, fmt.Sprintf("%s=%s", k, v))
	}

	sort.Strings(fields)
	fields = append(fields, l.defaultLevel.String())

	return strings.Join(fields, ":")
}



func (l *LoggerLevels) Enabled(lvl zapcore.Level) bool {
	l.mutex.RLock()
	enabled := l.minLevel.Enabled(lvl)
	l.mutex.RUnlock()
	return enabled
}
