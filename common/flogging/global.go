/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package flogging

import (
	"regexp"
	"strings"

	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc/grpclog"
)

const (
	defaultFormat = "%{color}%{time:2006-01-02 15:04:05.000 MST} [%{module}] %{shortfunc} -> %{level:.4s} %{id:03x}%{color:reset} %{message}"
	defaultLevel  = zapcore.InfoLevel
)

var Global *Logging
var logger *FabricLogger

func init() {
	logging, err := New(Config{})
	if err != nil {
		panic(err)
	}

	Global = logging
	logger = Global.Logger("flogging")
	grpcLogger := Global.ZapLogger("grpc")
	grpclog.SetLogger(NewGRPCLogger(grpcLogger))
}


func Init(config Config) {
	err := Global.Apply(config)
	if err != nil {
		panic(err)
	}
}




func Reset() {
	Global.ResetLevels()
	Global.Apply(Config{})
}


func GetModuleLevel(module string) string {
	return strings.ToUpper(Global.Level(module).String())
}




func SetModuleLevels(moduleRegexp, level string) error {
	re, err := regexp.Compile(moduleRegexp)
	if err != nil {
		return err
	}

	Global.SetLevels(re, NameToLevel(level))
	return nil
}


func SetModuleLevel(module string, level string) error {
	Global.SetLevel(module, NameToLevel(level))
	return nil
}



func MustGetLogger(module string) *FabricLogger {
	return Global.Logger(module)
}




func GetModuleLevels() map[string]zapcore.Level {
	return Global.Levels()
}



func RestoreLevels(levels map[string]zapcore.Level) {
	Global.RestoreLevels(levels)
}
