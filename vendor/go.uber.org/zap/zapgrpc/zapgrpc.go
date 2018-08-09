




















package zapgrpc 

import "go.uber.org/zap"


type Option interface {
	apply(*Logger)
}

type optionFunc func(*Logger)

func (f optionFunc) apply(log *Logger) {
	f(log)
}



func WithDebug() Option {
	return optionFunc(func(logger *Logger) {
		logger.print = (*zap.SugaredLogger).Debug
		logger.printf = (*zap.SugaredLogger).Debugf
	})
}




func NewLogger(l *zap.Logger, options ...Option) *Logger {
	logger := &Logger{
		log:    l.Sugar(),
		fatal:  (*zap.SugaredLogger).Fatal,
		fatalf: (*zap.SugaredLogger).Fatalf,
		print:  (*zap.SugaredLogger).Info,
		printf: (*zap.SugaredLogger).Infof,
	}
	for _, option := range options {
		option.apply(logger)
	}
	return logger
}


type Logger struct {
	log    *zap.SugaredLogger
	fatal  func(*zap.SugaredLogger, ...interface{})
	fatalf func(*zap.SugaredLogger, string, ...interface{})
	print  func(*zap.SugaredLogger, ...interface{})
	printf func(*zap.SugaredLogger, string, ...interface{})
}


func (l *Logger) Fatal(args ...interface{}) {
	l.fatal(l.log, args...)
}


func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.fatalf(l.log, format, args...)
}


func (l *Logger) Fatalln(args ...interface{}) {
	l.fatal(l.log, args...)
}


func (l *Logger) Print(args ...interface{}) {
	l.print(l.log, args...)
}


func (l *Logger) Printf(format string, args ...interface{}) {
	l.printf(l.log, format, args...)
}


func (l *Logger) Println(args ...interface{}) {
	l.print(l.log, args...)
}
