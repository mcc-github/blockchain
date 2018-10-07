package logrus

import (
	"io"
	"time"
)

var (
	
	std = New()
)

func StandardLogger() *Logger {
	return std
}


func SetOutput(out io.Writer) {
	std.SetOutput(out)
}


func SetFormatter(formatter Formatter) {
	std.SetFormatter(formatter)
}


func SetLevel(level Level) {
	std.SetLevel(level)
}


func GetLevel() Level {
	return std.GetLevel()
}


func IsLevelEnabled(level Level) bool {
	return std.IsLevelEnabled(level)
}


func AddHook(hook Hook) {
	std.AddHook(hook)
}


func WithError(err error) *Entry {
	return std.WithField(ErrorKey, err)
}






func WithField(key string, value interface{}) *Entry {
	return std.WithField(key, value)
}







func WithFields(fields Fields) *Entry {
	return std.WithFields(fields)
}






func WithTime(t time.Time) *Entry {
	return std.WithTime(t)
}


func Debug(args ...interface{}) {
	std.Debug(args...)
}


func Print(args ...interface{}) {
	std.Print(args...)
}


func Info(args ...interface{}) {
	std.Info(args...)
}


func Warn(args ...interface{}) {
	std.Warn(args...)
}


func Warning(args ...interface{}) {
	std.Warning(args...)
}


func Error(args ...interface{}) {
	std.Error(args...)
}


func Panic(args ...interface{}) {
	std.Panic(args...)
}


func Fatal(args ...interface{}) {
	std.Fatal(args...)
}


func Debugf(format string, args ...interface{}) {
	std.Debugf(format, args...)
}


func Printf(format string, args ...interface{}) {
	std.Printf(format, args...)
}


func Infof(format string, args ...interface{}) {
	std.Infof(format, args...)
}


func Warnf(format string, args ...interface{}) {
	std.Warnf(format, args...)
}


func Warningf(format string, args ...interface{}) {
	std.Warningf(format, args...)
}


func Errorf(format string, args ...interface{}) {
	std.Errorf(format, args...)
}


func Panicf(format string, args ...interface{}) {
	std.Panicf(format, args...)
}


func Fatalf(format string, args ...interface{}) {
	std.Fatalf(format, args...)
}


func Debugln(args ...interface{}) {
	std.Debugln(args...)
}


func Println(args ...interface{}) {
	std.Println(args...)
}


func Infoln(args ...interface{}) {
	std.Infoln(args...)
}


func Warnln(args ...interface{}) {
	std.Warnln(args...)
}


func Warningln(args ...interface{}) {
	std.Warningln(args...)
}


func Errorln(args ...interface{}) {
	std.Errorln(args...)
}


func Panicln(args ...interface{}) {
	std.Panicln(args...)
}


func Fatalln(args ...interface{}) {
	std.Fatalln(args...)
}
