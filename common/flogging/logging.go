/*
Copyright IBM Corp. 2017 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package flogging

import (
	"io"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/op/go-logging"
)

const (
	pkgLogID      = "flogging"
	defaultFormat = "%{color}%{time:2006-01-02 15:04:05.000 MST} [%{module}] %{shortfunc} -> %{level:.4s} %{id:03x}%{color:reset} %{message}"
	defaultLevel  = logging.INFO
)

var (
	logger *logging.Logger

	defaultOutput *os.File

	modules          map[string]string 
	peerStartModules map[string]string

	lock sync.RWMutex
	once sync.Once
)

func init() {
	logger = logging.MustGetLogger(pkgLogID)
	Reset()
	initgrpclogger()
}


func Reset() {
	modules = make(map[string]string)
	lock = sync.RWMutex{}

	defaultOutput = os.Stderr
	InitBackend(SetFormat(defaultFormat), defaultOutput)
	InitFromSpec("")
}


func SetFormat(formatSpec string) logging.Formatter {
	if formatSpec == "" {
		formatSpec = defaultFormat
	}
	return logging.MustStringFormatter(formatSpec)
}



func InitBackend(formatter logging.Formatter, output io.Writer) {
	backend := logging.NewLogBackend(output, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, formatter)
	logging.SetBackend(backendFormatter).SetLevel(defaultLevel, "")
}


func DefaultLevel() string {
	return defaultLevel.String()
}


func GetModuleLevel(module string) string {
	
	
	
	level := logging.GetLevel(module).String()
	return level
}




func SetModuleLevel(moduleRegExp string, level string) (string, error) {
	return setModuleLevel(moduleRegExp, level, true, false)
}

func setModuleLevel(moduleRegExp string, level string, isRegExp bool, revert bool) (string, error) {
	var re *regexp.Regexp
	logLevel, err := logging.LogLevel(level)
	if err != nil {
		logger.Warningf("Invalid logging level '%s' - ignored", level)
	} else {
		if !isRegExp || revert {
			logging.SetLevel(logLevel, moduleRegExp)
			logger.Debugf("Module '%s' logger enabled for log level '%s'", moduleRegExp, level)
		} else {
			re, err = regexp.Compile(moduleRegExp)
			if err != nil {
				logger.Warningf("Invalid regular expression: %s", moduleRegExp)
				return "", err
			}
			lock.Lock()
			defer lock.Unlock()
			for module := range modules {
				if re.MatchString(module) {
					logging.SetLevel(logging.Level(logLevel), module)
					modules[module] = logLevel.String()
					logger.Debugf("Module '%s' logger enabled for log level '%s'", module, logLevel)
				}
			}
		}
	}
	return logLevel.String(), err
}



func MustGetLogger(module string) *logging.Logger {
	l := logging.MustGetLogger(module)
	lock.Lock()
	defer lock.Unlock()
	modules[module] = GetModuleLevel(module)
	return l
}





func InitFromSpec(spec string) string {
	levelAll := defaultLevel
	var err error

	if spec != "" {
		fields := strings.Split(spec, ":")
		for _, field := range fields {
			split := strings.Split(field, "=")
			switch len(split) {
			case 1:
				if levelAll, err = logging.LogLevel(field); err != nil {
					logger.Warningf("Logging level '%s' not recognized, defaulting to '%s': %s", field, defaultLevel, err)
					levelAll = defaultLevel 
				}
			case 2:
				
				levelSingle, err := logging.LogLevel(split[1])
				if err != nil {
					logger.Warningf("Invalid logging level in '%s' ignored", field)
					continue
				}

				if split[0] == "" {
					logger.Warningf("Invalid logging override specification '%s' ignored - no module specified", field)
				} else {
					modules := strings.Split(split[0], ",")
					for _, module := range modules {
						logger.Debugf("Setting logging level for module '%s' to '%s'", module, levelSingle)
						logging.SetLevel(levelSingle, module)
					}
				}
			default:
				logger.Warningf("Invalid logging override '%s' ignored - missing ':'?", field)
			}
		}
	}

	logging.SetLevel(levelAll, "") 

	
	
	for k := range modules {
		MustGetLogger(k)
	}
	
	MustGetLogger(pkgLogID)

	return levelAll.String()
}



func SetPeerStartupModulesMap() {
	lock.Lock()
	defer lock.Unlock()

	once.Do(func() {
		peerStartModules = make(map[string]string)
		for k, v := range modules {
			peerStartModules[k] = v
		}
	})
}




func GetPeerStartupLevel(module string) string {
	if module != "" {
		if level, ok := peerStartModules[module]; ok {
			return level
		}
	}

	return ""
}



func RevertToPeerStartupLevels() error {
	lock.RLock()
	defer lock.RUnlock()
	for key := range peerStartModules {
		_, err := setModuleLevel(key, peerStartModules[key], false, true)
		if err != nil {
			return err
		}
	}
	logger.Info("Log levels reverted to the levels defined at the end of peer startup")
	return nil
}
