



















package zap

import (
	"fmt"

	"go.uber.org/zap/zapcore"

	"go.uber.org/multierr"
)

const (
	_oddNumberErrMsg    = "Ignored key without a value."
	_nonStringKeyErrMsg = "Ignored key-value pairs with non-string keys."
)










type SugaredLogger struct {
	base *Logger
}





func (s *SugaredLogger) Desugar() *Logger {
	base := s.base.clone()
	base.callerSkip -= 2
	return base
}


func (s *SugaredLogger) Named(name string) *SugaredLogger {
	return &SugaredLogger{base: s.base.Named(name)}
}




























func (s *SugaredLogger) With(args ...interface{}) *SugaredLogger {
	return &SugaredLogger{base: s.base.With(s.sweetenFields(args)...)}
}


func (s *SugaredLogger) Debug(args ...interface{}) {
	s.log(DebugLevel, "", args, nil)
}


func (s *SugaredLogger) Info(args ...interface{}) {
	s.log(InfoLevel, "", args, nil)
}


func (s *SugaredLogger) Warn(args ...interface{}) {
	s.log(WarnLevel, "", args, nil)
}


func (s *SugaredLogger) Error(args ...interface{}) {
	s.log(ErrorLevel, "", args, nil)
}



func (s *SugaredLogger) DPanic(args ...interface{}) {
	s.log(DPanicLevel, "", args, nil)
}


func (s *SugaredLogger) Panic(args ...interface{}) {
	s.log(PanicLevel, "", args, nil)
}


func (s *SugaredLogger) Fatal(args ...interface{}) {
	s.log(FatalLevel, "", args, nil)
}


func (s *SugaredLogger) Debugf(template string, args ...interface{}) {
	s.log(DebugLevel, template, args, nil)
}


func (s *SugaredLogger) Infof(template string, args ...interface{}) {
	s.log(InfoLevel, template, args, nil)
}


func (s *SugaredLogger) Warnf(template string, args ...interface{}) {
	s.log(WarnLevel, template, args, nil)
}


func (s *SugaredLogger) Errorf(template string, args ...interface{}) {
	s.log(ErrorLevel, template, args, nil)
}



func (s *SugaredLogger) DPanicf(template string, args ...interface{}) {
	s.log(DPanicLevel, template, args, nil)
}


func (s *SugaredLogger) Panicf(template string, args ...interface{}) {
	s.log(PanicLevel, template, args, nil)
}


func (s *SugaredLogger) Fatalf(template string, args ...interface{}) {
	s.log(FatalLevel, template, args, nil)
}






func (s *SugaredLogger) Debugw(msg string, keysAndValues ...interface{}) {
	s.log(DebugLevel, msg, nil, keysAndValues)
}



func (s *SugaredLogger) Infow(msg string, keysAndValues ...interface{}) {
	s.log(InfoLevel, msg, nil, keysAndValues)
}



func (s *SugaredLogger) Warnw(msg string, keysAndValues ...interface{}) {
	s.log(WarnLevel, msg, nil, keysAndValues)
}



func (s *SugaredLogger) Errorw(msg string, keysAndValues ...interface{}) {
	s.log(ErrorLevel, msg, nil, keysAndValues)
}




func (s *SugaredLogger) DPanicw(msg string, keysAndValues ...interface{}) {
	s.log(DPanicLevel, msg, nil, keysAndValues)
}



func (s *SugaredLogger) Panicw(msg string, keysAndValues ...interface{}) {
	s.log(PanicLevel, msg, nil, keysAndValues)
}



func (s *SugaredLogger) Fatalw(msg string, keysAndValues ...interface{}) {
	s.log(FatalLevel, msg, nil, keysAndValues)
}


func (s *SugaredLogger) Sync() error {
	return s.base.Sync()
}

func (s *SugaredLogger) log(lvl zapcore.Level, template string, fmtArgs []interface{}, context []interface{}) {
	
	
	if lvl < DPanicLevel && !s.base.Core().Enabled(lvl) {
		return
	}

	
	msg := template
	if msg == "" && len(fmtArgs) > 0 {
		msg = fmt.Sprint(fmtArgs...)
	} else if msg != "" && len(fmtArgs) > 0 {
		msg = fmt.Sprintf(template, fmtArgs...)
	}

	if ce := s.base.Check(lvl, msg); ce != nil {
		ce.Write(s.sweetenFields(context)...)
	}
}

func (s *SugaredLogger) sweetenFields(args []interface{}) []Field {
	if len(args) == 0 {
		return nil
	}

	
	
	fields := make([]Field, 0, len(args))
	var invalid invalidPairs

	for i := 0; i < len(args); {
		
		if f, ok := args[i].(Field); ok {
			fields = append(fields, f)
			i++
			continue
		}

		
		if i == len(args)-1 {
			s.base.DPanic(_oddNumberErrMsg, Any("ignored", args[i]))
			break
		}

		
		
		key, val := args[i], args[i+1]
		if keyStr, ok := key.(string); !ok {
			
			if cap(invalid) == 0 {
				invalid = make(invalidPairs, 0, len(args)/2)
			}
			invalid = append(invalid, invalidPair{i, key, val})
		} else {
			fields = append(fields, Any(keyStr, val))
		}
		i += 2
	}

	
	if len(invalid) > 0 {
		s.base.DPanic(_nonStringKeyErrMsg, Array("invalid", invalid))
	}
	return fields
}

type invalidPair struct {
	position   int
	key, value interface{}
}

func (p invalidPair) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddInt64("position", int64(p.position))
	Any("key", p.key).AddTo(enc)
	Any("value", p.value).AddTo(enc)
	return nil
}

type invalidPairs []invalidPair

func (ps invalidPairs) MarshalLogArray(enc zapcore.ArrayEncoder) error {
	var err error
	for i := range ps {
		err = multierr.Append(err, enc.AppendObject(ps[i]))
	}
	return err
}
