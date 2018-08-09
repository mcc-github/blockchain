/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package flogging

import (
	"regexp"
	"sync"

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
