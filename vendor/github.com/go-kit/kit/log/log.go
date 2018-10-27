package log

import "errors"






type Logger interface {
	Log(keyvals ...interface{}) error
}



var ErrMissingValue = errors.New("(MISSING)")







func With(logger Logger, keyvals ...interface{}) Logger {
	if len(keyvals) == 0 {
		return logger
	}
	l := newContext(logger)
	kvs := append(l.keyvals, keyvals...)
	if len(kvs)%2 != 0 {
		kvs = append(kvs, ErrMissingValue)
	}
	return &context{
		logger: l.logger,
		
		
		
		
		keyvals:   kvs[:len(kvs):len(kvs)],
		hasValuer: l.hasValuer || containsValuer(keyvals),
	}
}







func WithPrefix(logger Logger, keyvals ...interface{}) Logger {
	if len(keyvals) == 0 {
		return logger
	}
	l := newContext(logger)
	
	
	
	
	n := len(l.keyvals) + len(keyvals)
	if len(keyvals)%2 != 0 {
		n++
	}
	kvs := make([]interface{}, 0, n)
	kvs = append(kvs, keyvals...)
	if len(kvs)%2 != 0 {
		kvs = append(kvs, ErrMissingValue)
	}
	kvs = append(kvs, l.keyvals...)
	return &context{
		logger:    l.logger,
		keyvals:   kvs,
		hasValuer: l.hasValuer || containsValuer(keyvals),
	}
}




















type context struct {
	logger    Logger
	keyvals   []interface{}
	hasValuer bool
}

func newContext(logger Logger) *context {
	if c, ok := logger.(*context); ok {
		return c
	}
	return &context{logger: logger}
}




func (l *context) Log(keyvals ...interface{}) error {
	kvs := append(l.keyvals, keyvals...)
	if len(kvs)%2 != 0 {
		kvs = append(kvs, ErrMissingValue)
	}
	if l.hasValuer {
		
		
		if len(keyvals) == 0 {
			kvs = append([]interface{}{}, l.keyvals...)
		}
		bindValues(kvs[:len(l.keyvals)])
	}
	return l.logger.Log(kvs...)
}




type LoggerFunc func(...interface{}) error


func (f LoggerFunc) Log(keyvals ...interface{}) error {
	return f(keyvals...)
}
