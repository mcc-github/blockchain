



















package zap

import "go.uber.org/zap/zapcore"


type Option interface {
	apply(*Logger)
}


type optionFunc func(*Logger)

func (f optionFunc) apply(log *Logger) {
	f(log)
}


func WrapCore(f func(zapcore.Core) zapcore.Core) Option {
	return optionFunc(func(log *Logger) {
		log.core = f(log.core)
	})
}








func Hooks(hooks ...func(zapcore.Entry) error) Option {
	return optionFunc(func(log *Logger) {
		log.core = zapcore.RegisterHooks(log.core, hooks...)
	})
}


func Fields(fs ...Field) Option {
	return optionFunc(func(log *Logger) {
		log.core = log.core.With(fs)
	})
}








func ErrorOutput(w zapcore.WriteSyncer) Option {
	return optionFunc(func(log *Logger) {
		log.errorOutput = w
	})
}



func Development() Option {
	return optionFunc(func(log *Logger) {
		log.development = true
	})
}



func AddCaller() Option {
	return optionFunc(func(log *Logger) {
		log.addCaller = true
	})
}





func AddCallerSkip(skip int) Option {
	return optionFunc(func(log *Logger) {
		log.callerSkip += skip
	})
}



func AddStacktrace(lvl zapcore.LevelEnabler) Option {
	return optionFunc(func(log *Logger) {
		log.addStack = lvl
	})
}
