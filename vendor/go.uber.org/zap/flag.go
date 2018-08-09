



















package zap

import (
	"flag"

	"go.uber.org/zap/zapcore"
)







func LevelFlag(name string, defaultLevel zapcore.Level, usage string) *zapcore.Level {
	lvl := defaultLevel
	flag.Var(&lvl, name, usage)
	return &lvl
}
