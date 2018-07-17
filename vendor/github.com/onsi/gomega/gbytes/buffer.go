
package gbytes

import (
	"errors"
	"fmt"
	"io"
	"regexp"
	"sync"
	"time"
)


type Buffer struct {
	contents     []byte
	readCursor   uint64
	lock         *sync.Mutex
	detectCloser chan interface{}
	closed       bool
}


func NewBuffer() *Buffer {
	return &Buffer{
		lock: &sync.Mutex{},
	}
}


func BufferWithBytes(bytes []byte) *Buffer {
	return &Buffer{
		lock:     &sync.Mutex{},
		contents: bytes,
	}
}


func BufferReader(reader io.Reader) *Buffer {
	b := &Buffer{
		lock: &sync.Mutex{},
	}

	go func() {
		io.Copy(b, reader)
		b.Close()
	}()

	return b
}


func (b *Buffer) Write(p []byte) (n int, err error) {
	b.lock.Lock()
	defer b.lock.Unlock()

	if b.closed {
		return 0, errors.New("attempt to write to closed buffer")
	}

	b.contents = append(b.contents, p...)
	return len(p), nil
}


func (b *Buffer) Read(d []byte) (int, error) {
	b.lock.Lock()
	defer b.lock.Unlock()

	if b.closed {
		return 0, errors.New("attempt to read from closed buffer")
	}

	if uint64(len(b.contents)) <= b.readCursor {
		return 0, io.EOF
	}

	n := copy(d, b.contents[b.readCursor:])
	b.readCursor += uint64(n)

	return n, nil
}


func (b *Buffer) Close() error {
	b.lock.Lock()
	defer b.lock.Unlock()

	b.closed = true

	return nil
}


func (b *Buffer) Closed() bool {
	b.lock.Lock()
	defer b.lock.Unlock()

	return b.closed
}


func (b *Buffer) Contents() []byte {
	b.lock.Lock()
	defer b.lock.Unlock()

	contents := make([]byte, len(b.contents))
	copy(contents, b.contents)
	return contents
}


func (b *Buffer) Detect(desired string, args ...interface{}) chan bool {
	formattedRegexp := desired
	if len(args) > 0 {
		formattedRegexp = fmt.Sprintf(desired, args...)
	}
	re := regexp.MustCompile(formattedRegexp)

	b.lock.Lock()
	defer b.lock.Unlock()

	if b.detectCloser == nil {
		b.detectCloser = make(chan interface{})
	}

	closer := b.detectCloser
	response := make(chan bool)
	go func() {
		ticker := time.NewTicker(10 * time.Millisecond)
		defer ticker.Stop()
		defer close(response)
		for {
			select {
			case <-ticker.C:
				b.lock.Lock()
				data, cursor := b.contents[b.readCursor:], b.readCursor
				loc := re.FindIndex(data)
				b.lock.Unlock()

				if loc != nil {
					response <- true
					b.lock.Lock()
					newCursorPosition := cursor + uint64(loc[1])
					if newCursorPosition >= b.readCursor {
						b.readCursor = newCursorPosition
					}
					b.lock.Unlock()
					return
				}
			case <-closer:
				return
			}
		}
	}()

	return response
}


func (b *Buffer) CancelDetects() {
	b.lock.Lock()
	defer b.lock.Unlock()

	close(b.detectCloser)
	b.detectCloser = nil
}

func (b *Buffer) didSay(re *regexp.Regexp) (bool, []byte) {
	b.lock.Lock()
	defer b.lock.Unlock()

	unreadBytes := b.contents[b.readCursor:]
	copyOfUnreadBytes := make([]byte, len(unreadBytes))
	copy(copyOfUnreadBytes, unreadBytes)

	loc := re.FindIndex(unreadBytes)

	if loc != nil {
		b.readCursor += uint64(loc[1])
		return true, copyOfUnreadBytes
	}
	return false, copyOfUnreadBytes
}
