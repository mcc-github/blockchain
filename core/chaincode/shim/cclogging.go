/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package shim

import (
	"os"
	"strings"
	"sync"

	logging "github.com/op/go-logging"
)



func IsEnabledForLogLevel(logLevel string) bool {
	lvl, _ := logging.LogLevel(logLevel)
	return chaincodeLogger.IsEnabledFor(lvl)
}

var loggingSetup sync.Once




func SetupChaincodeLogging() {
	loggingSetup.Do(setupChaincodeLogging)
}

func setupChaincodeLogging() {
	
	const defaultLogFormat = "%{color}%{time:2006-01-02 15:04:05.000 MST} [%{module}] %{shortfunc} -> %{level:.4s} %{id:03x}%{color:reset} %{message}"
	const defaultLevel = logging.INFO

	
	logFormat := os.Getenv("CORE_CHAINCODE_LOGGING_FORMAT")
	if logFormat == "" {
		logFormat = defaultLogFormat
	}

	formatter := logging.MustStringFormatter(logFormat)
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, formatter)
	logging.SetBackend(backendFormatter).SetLevel(defaultLevel, "")

	
	chaincodeLogLevelString := os.Getenv("CORE_CHAINCODE_LOGGING_LEVEL")
	if chaincodeLogLevelString == "" {
		chaincodeLogger.Infof("Chaincode log level not provided; defaulting to: %s", defaultLevel.String())
		chaincodeLogLevelString = defaultLevel.String()
	}

	_, err := LogLevel(chaincodeLogLevelString)
	if err != nil {
		chaincodeLogger.Warningf("Error: '%s' for chaincode log level: %s; defaulting to %s", err, chaincodeLogLevelString, defaultLevel.String())
		chaincodeLogLevelString = defaultLevel.String()
	}

	initFromSpec(chaincodeLogLevelString, defaultLevel)

	
	
	
	
	shimLogLevelString := os.Getenv("CORE_CHAINCODE_LOGGING_SHIM")
	if shimLogLevelString != "" {
		shimLogLevel, err := LogLevel(shimLogLevelString)
		if err == nil {
			SetLoggingLevel(shimLogLevel)
		} else {
			chaincodeLogger.Warningf("Error: %s for shim log level: %s", err, shimLogLevelString)
		}
	}

	
	
	buildLevel := os.Getenv("CORE_CHAINCODE_BUILDLEVEL")
	chaincodeLogger.Infof("Chaincode (build level: %s) starting up ...", buildLevel)
}


func initFromSpec(spec string, defaultLevel logging.Level) {
	levelAll := defaultLevel
	var err error

	fields := strings.Split(spec, ":")
	for _, field := range fields {
		split := strings.Split(field, "=")
		switch len(split) {
		case 1:
			if levelAll, err = logging.LogLevel(field); err != nil {
				chaincodeLogger.Warningf("Logging level '%s' not recognized, defaulting to '%s': %s", field, defaultLevel, err)
				levelAll = defaultLevel 
			}
		case 2:
			
			levelSingle, err := logging.LogLevel(split[1])
			if err != nil {
				chaincodeLogger.Warningf("Invalid logging level in '%s' ignored", field)
				continue
			}

			if split[0] == "" {
				chaincodeLogger.Warningf("Invalid logging override specification '%s' ignored - no logger specified", field)
			} else {
				loggers := strings.Split(split[0], ",")
				for _, logger := range loggers {
					chaincodeLogger.Debugf("Setting logging level for logger '%s' to '%s'", logger, levelSingle)
					logging.SetLevel(levelSingle, logger)
				}
			}
		default:
			chaincodeLogger.Warningf("Invalid logging override '%s' ignored - missing ':'?", field)
		}
	}

	logging.SetLevel(levelAll, "") 
}

































type LoggingLevel logging.Level


const (
	LogDebug    = LoggingLevel(logging.DEBUG)
	LogInfo     = LoggingLevel(logging.INFO)
	LogNotice   = LoggingLevel(logging.NOTICE)
	LogWarning  = LoggingLevel(logging.WARNING)
	LogError    = LoggingLevel(logging.ERROR)
	LogCritical = LoggingLevel(logging.CRITICAL)
)

var shimLoggingLevel = LogInfo 



func SetLoggingLevel(level LoggingLevel) {
	shimLoggingLevel = level
	logging.SetLevel(logging.Level(level), "shim")
}




func LogLevel(levelString string) (LoggingLevel, error) {
	l, err := logging.LogLevel(levelString)
	level := LoggingLevel(l)
	if err != nil {
		level = LogError
	}
	return level, err
}



type ChaincodeLogger struct {
	logger *logging.Logger
}






func NewLogger(name string) *ChaincodeLogger {
	return &ChaincodeLogger{logging.MustGetLogger(name)}
}




func (c *ChaincodeLogger) SetLevel(level LoggingLevel) {
	logging.SetLevel(logging.Level(level), c.logger.Module)
}



func (c *ChaincodeLogger) IsEnabledFor(level LoggingLevel) bool {
	return c.logger.IsEnabledFor(logging.Level(level))
}



func (c *ChaincodeLogger) Debug(args ...interface{}) {
	c.logger.Debug(args...)
}



func (c *ChaincodeLogger) Info(args ...interface{}) {
	c.logger.Info(args...)
}



func (c *ChaincodeLogger) Notice(args ...interface{}) {
	c.logger.Notice(args...)
}



func (c *ChaincodeLogger) Warning(args ...interface{}) {
	c.logger.Warning(args...)
}



func (c *ChaincodeLogger) Error(args ...interface{}) {
	c.logger.Error(args...)
}


func (c *ChaincodeLogger) Critical(args ...interface{}) {
	c.logger.Critical(args...)
}



func (c *ChaincodeLogger) Debugf(format string, args ...interface{}) {
	c.logger.Debugf(format, args...)
}



func (c *ChaincodeLogger) Infof(format string, args ...interface{}) {
	c.logger.Infof(format, args...)
}



func (c *ChaincodeLogger) Noticef(format string, args ...interface{}) {
	c.logger.Noticef(format, args...)
}



func (c *ChaincodeLogger) Warningf(format string, args ...interface{}) {
	c.logger.Warningf(format, args...)
}



func (c *ChaincodeLogger) Errorf(format string, args ...interface{}) {
	c.logger.Errorf(format, args...)
}


func (c *ChaincodeLogger) Criticalf(format string, args ...interface{}) {
	c.logger.Criticalf(format, args...)
}
