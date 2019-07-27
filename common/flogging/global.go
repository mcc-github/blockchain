/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package flogging

import (
	"io"

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
	Global.Apply(Config{})
}



func LoggerLevel(loggerName string) string {
	return Global.Level(loggerName).String()
}



func MustGetLogger(loggerName string) *FabricLogger {
	return Global.Logger(loggerName)
}


func ActivateSpec(spec string) {
	err := Global.ActivateSpec(spec)
	if err != nil {
		panic(err)
	}
}


func DefaultLevel() string {
	return defaultLevel.String()
}



func SetWriter(w io.Writer) io.Writer {
	return Global.SetWriter(w)
}



func SetObserver(observer Observer) Observer {
	return Global.SetObserver(observer)
}
