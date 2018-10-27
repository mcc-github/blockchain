package log

type nopLogger struct{}


func NewNopLogger() Logger { return nopLogger{} }

func (nopLogger) Log(...interface{}) error { return nil }
