



















package zap

import (
	"fmt"
	"io"
	"io/ioutil"

	"go.uber.org/zap/zapcore"

	"go.uber.org/multierr"
)



















func Open(paths ...string) (zapcore.WriteSyncer, func(), error) {
	writers, close, err := open(paths)
	if err != nil {
		return nil, nil, err
	}

	writer := CombineWriteSyncers(writers...)
	return writer, close, nil
}

func open(paths []string) ([]zapcore.WriteSyncer, func(), error) {
	writers := make([]zapcore.WriteSyncer, 0, len(paths))
	closers := make([]io.Closer, 0, len(paths))
	close := func() {
		for _, c := range closers {
			c.Close()
		}
	}

	var openErr error
	for _, path := range paths {
		sink, err := newSink(path)
		if err != nil {
			openErr = multierr.Append(openErr, fmt.Errorf("couldn't open sink %q: %v", path, err))
			continue
		}
		writers = append(writers, sink)
		closers = append(closers, sink)
	}
	if openErr != nil {
		close()
		return writers, nil, openErr
	}

	return writers, close, nil
}







func CombineWriteSyncers(writers ...zapcore.WriteSyncer) zapcore.WriteSyncer {
	if len(writers) == 0 {
		return zapcore.AddSync(ioutil.Discard)
	}
	return zapcore.Lock(zapcore.NewMultiWriteSyncer(writers...))
}
