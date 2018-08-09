



















package observer

import "go.uber.org/zap/zapcore"



type LoggedEntry struct {
	zapcore.Entry
	Context []zapcore.Field
}


func (e LoggedEntry) ContextMap() map[string]interface{} {
	encoder := zapcore.NewMapObjectEncoder()
	for _, f := range e.Context {
		f.AddTo(encoder)
	}
	return encoder.Fields
}
