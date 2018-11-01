/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package flogging

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"go.uber.org/zap/zapcore"
)


type ModuleLevels struct {
	defaultLevel zapcore.Level

	mutex      sync.RWMutex
	levelCache map[string]zapcore.Level
	specs      map[string]zapcore.Level
}



func (m *ModuleLevels) DefaultLevel() zapcore.Level {
	m.mutex.RLock()
	l := m.defaultLevel
	m.mutex.RUnlock()
	return l
}





func (m *ModuleLevels) ActivateSpec(spec string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

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
			for _, loggers := range loggers {
				specs[loggers] = level
			}

		default:
			return errors.Errorf("invalid logging specification '%s': bad segment '%s'", spec, field)
		}
	}

	m.defaultLevel = defaultLevel
	m.specs = specs
	m.levelCache = map[string]zapcore.Level{}

	return nil
}




func (m *ModuleLevels) Level(loggerName string) zapcore.Level {
	if level, ok := m.cachedLevel(loggerName); ok {
		return level
	}

	m.mutex.Lock()
	level := m.calculateLevel(loggerName)
	m.levelCache[loggerName] = level
	m.mutex.Unlock()

	return level
}



func (m *ModuleLevels) calculateLevel(loggerName string) zapcore.Level {
	candidate := loggerName + "."
	for {
		if lvl, ok := m.specs[candidate]; ok {
			return lvl
		}

		idx := strings.LastIndex(candidate, ".")
		if idx <= 0 {
			return m.defaultLevel
		}
		candidate = candidate[:idx]
	}
}



func (m *ModuleLevels) cachedLevel(loggerName string) (lvl zapcore.Level, ok bool) {
	m.mutex.RLock()
	level, ok := m.levelCache[loggerName]
	m.mutex.RUnlock()
	return level, ok
}


func (m *ModuleLevels) Spec() string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	var fields []string
	for k, v := range m.specs {
		fields = append(fields, fmt.Sprintf("%s=%s", k, v))
	}

	sort.Strings(fields)
	fields = append(fields, m.defaultLevel.String())

	return strings.Join(fields, ":")
}
