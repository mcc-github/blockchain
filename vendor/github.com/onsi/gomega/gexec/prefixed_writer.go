package gexec

import (
	"io"
	"sync"
)


type PrefixedWriter struct {
	prefix        []byte
	writer        io.Writer
	lock          *sync.Mutex
	atStartOfLine bool
}

func NewPrefixedWriter(prefix string, writer io.Writer) *PrefixedWriter {
	return &PrefixedWriter{
		prefix:        []byte(prefix),
		writer:        writer,
		lock:          &sync.Mutex{},
		atStartOfLine: true,
	}
}

func (w *PrefixedWriter) Write(b []byte) (int, error) {
	w.lock.Lock()
	defer w.lock.Unlock()

	toWrite := []byte{}

	for _, c := range b {
		if w.atStartOfLine {
			toWrite = append(toWrite, w.prefix...)
		}

		toWrite = append(toWrite, c)

		w.atStartOfLine = c == '\n'
	}

	_, err := w.writer.Write(toWrite)
	if err != nil {
		return 0, err
	}

	return len(b), nil
}
