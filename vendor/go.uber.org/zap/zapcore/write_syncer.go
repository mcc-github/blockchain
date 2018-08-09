



















package zapcore

import (
	"io"
	"sync"

	"go.uber.org/multierr"
)



type WriteSyncer interface {
	io.Writer
	Sync() error
}




func AddSync(w io.Writer) WriteSyncer {
	switch w := w.(type) {
	case WriteSyncer:
		return w
	default:
		return writerWrapper{w}
	}
}

type lockedWriteSyncer struct {
	sync.Mutex
	ws WriteSyncer
}



func Lock(ws WriteSyncer) WriteSyncer {
	if _, ok := ws.(*lockedWriteSyncer); ok {
		
		return ws
	}
	return &lockedWriteSyncer{ws: ws}
}

func (s *lockedWriteSyncer) Write(bs []byte) (int, error) {
	s.Lock()
	n, err := s.ws.Write(bs)
	s.Unlock()
	return n, err
}

func (s *lockedWriteSyncer) Sync() error {
	s.Lock()
	err := s.ws.Sync()
	s.Unlock()
	return err
}

type writerWrapper struct {
	io.Writer
}

func (w writerWrapper) Sync() error {
	return nil
}

type multiWriteSyncer []WriteSyncer



func NewMultiWriteSyncer(ws ...WriteSyncer) WriteSyncer {
	if len(ws) == 1 {
		return ws[0]
	}
	
	return multiWriteSyncer(append([]WriteSyncer(nil), ws...))
}





func (ws multiWriteSyncer) Write(p []byte) (int, error) {
	var writeErr error
	nWritten := 0
	for _, w := range ws {
		n, err := w.Write(p)
		writeErr = multierr.Append(writeErr, err)
		if nWritten == 0 && n != 0 {
			nWritten = n
		} else if n < nWritten {
			nWritten = n
		}
	}
	return nWritten, writeErr
}

func (ws multiWriteSyncer) Sync() error {
	var err error
	for _, w := range ws {
		err = multierr.Append(err, w.Sync())
	}
	return err
}
