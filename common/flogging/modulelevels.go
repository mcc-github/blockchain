/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package flogging

import (
	"regexp"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)


type ModuleLevels struct {
	defaultLevel zapcore.Level

	mutex  sync.RWMutex
	levels map[string]zapcore.Level
}



func (m *ModuleLevels) SetDefaultLevel(l zapcore.Level) {
	m.mutex.Lock()
	m.defaultLevel = l
	m.mutex.Unlock()
}



func (m *ModuleLevels) DefaultLevel() zapcore.Level {
	m.mutex.RLock()
	l := m.defaultLevel
	m.mutex.RUnlock()
	return l
}



func (m *ModuleLevels) Reset() {
	m.mutex.Lock()
	m.levels = map[string]zapcore.Level{}
	m.defaultLevel = zapcore.InfoLevel
	m.mutex.Unlock()
}





func (m *ModuleLevels) ActivateSpec(spec string) error {
	levelAll := m.defaultLevel
	updates := map[string]zapcore.Level{}

	fields := strings.Split(spec, ":")
	for _, field := range fields {
		split := strings.Split(field, "=")
		switch len(split) {
		case 1: 
			levelAll = NameToLevel(field)

		case 2: 
			level := NameToLevel(split[1])
			if split[0] == "" {
				return errors.Errorf("invalid logging specification '%s': no module specified in segment '%s'", spec, field)
			}

			modules := strings.Split(split[0], ",")
			for _, module := range modules {
				updates[module] = level
			}

		default:
			return errors.Errorf("invalid logging specification '%s': bad segment '%s'", spec, field)
		}
	}

	m.Reset()
	m.SetDefaultLevel(NameToLevel(levelAll.String()))
	for module, level := range updates {
		m.SetLevel(module, level)
	}

	return nil
}


func (m *ModuleLevels) SetLevel(module string, l zapcore.Level) {
	m.mutex.Lock()
	if m.levels == nil {
		m.levels = map[string]zapcore.Level{}
	}
	m.levels[module] = l
	m.mutex.Unlock()
}



func (m *ModuleLevels) SetLevels(re *regexp.Regexp, l zapcore.Level) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for module := range m.levels {
		if re.MatchString(module) {
			m.levels[module] = l
		}
	}
}




func (m *ModuleLevels) Level(module string) zapcore.Level {
	m.mutex.RLock()
	l, ok := m.levels[module]
	if !ok {
		l = m.defaultLevel
	}
	m.mutex.RUnlock()
	return l
}


func (m *ModuleLevels) Levels() map[string]zapcore.Level {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	levels := make(map[string]zapcore.Level, len(m.levels))
	for k, v := range m.levels {
		levels[k] = v
	}
	return levels
}



func (m *ModuleLevels) RestoreLevels(levels map[string]zapcore.Level) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	m.levels = map[string]zapcore.Level{}
	for k, v := range levels {
		m.levels[k] = v
	}
}


func (m *ModuleLevels) LevelEnabler(module string) zapcore.LevelEnabler {
	return zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return m.Level(module).Enabled(l)
	})
}
