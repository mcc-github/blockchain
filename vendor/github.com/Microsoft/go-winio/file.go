

package winio

import (
	"errors"
	"io"
	"runtime"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)






type atomicBool int32

func (b *atomicBool) isSet() bool { return atomic.LoadInt32((*int32)(b)) != 0 }
func (b *atomicBool) setFalse()   { atomic.StoreInt32((*int32)(b), 0) }
func (b *atomicBool) setTrue()    { atomic.StoreInt32((*int32)(b), 1) }
func (b *atomicBool) swap(new bool) bool {
	var newInt int32
	if new {
		newInt = 1
	}
	return atomic.SwapInt32((*int32)(b), newInt) == 1
}

const (
	cFILE_SKIP_COMPLETION_PORT_ON_SUCCESS = 1
	cFILE_SKIP_SET_EVENT_ON_HANDLE        = 2
)

var (
	ErrFileClosed = errors.New("file has already been closed")
	ErrTimeout    = &timeoutError{}
)

type timeoutError struct{}

func (e *timeoutError) Error() string   { return "i/o timeout" }
func (e *timeoutError) Timeout() bool   { return true }
func (e *timeoutError) Temporary() bool { return true }

type timeoutChan chan struct{}

var ioInitOnce sync.Once
var ioCompletionPort syscall.Handle


type ioResult struct {
	bytes uint32
	err   error
}


type ioOperation struct {
	o  syscall.Overlapped
	ch chan ioResult
}

func initIo() {
	h, err := createIoCompletionPort(syscall.InvalidHandle, 0, 0, 0xffffffff)
	if err != nil {
		panic(err)
	}
	ioCompletionPort = h
	go ioCompletionProcessor(h)
}



type win32File struct {
	handle        syscall.Handle
	wg            sync.WaitGroup
	wgLock        sync.RWMutex
	closing       atomicBool
	readDeadline  deadlineHandler
	writeDeadline deadlineHandler
}

type deadlineHandler struct {
	setLock     sync.Mutex
	channel     timeoutChan
	channelLock sync.RWMutex
	timer       *time.Timer
	timedout    atomicBool
}


func makeWin32File(h syscall.Handle) (*win32File, error) {
	f := &win32File{handle: h}
	ioInitOnce.Do(initIo)
	_, err := createIoCompletionPort(h, ioCompletionPort, 0, 0xffffffff)
	if err != nil {
		return nil, err
	}
	err = setFileCompletionNotificationModes(h, cFILE_SKIP_COMPLETION_PORT_ON_SUCCESS|cFILE_SKIP_SET_EVENT_ON_HANDLE)
	if err != nil {
		return nil, err
	}
	f.readDeadline.channel = make(timeoutChan)
	f.writeDeadline.channel = make(timeoutChan)
	return f, nil
}

func MakeOpenFile(h syscall.Handle) (io.ReadWriteCloser, error) {
	return makeWin32File(h)
}


func (f *win32File) closeHandle() {
	f.wgLock.Lock()
	
	if !f.closing.swap(true) {
		f.wgLock.Unlock()
		
		cancelIoEx(f.handle, nil)
		f.wg.Wait()
		
		syscall.Close(f.handle)
		f.handle = 0
	} else {
		f.wgLock.Unlock()
	}
}


func (f *win32File) Close() error {
	f.closeHandle()
	return nil
}



func (f *win32File) prepareIo() (*ioOperation, error) {
	f.wgLock.RLock()
	if f.closing.isSet() {
		f.wgLock.RUnlock()
		return nil, ErrFileClosed
	}
	f.wg.Add(1)
	f.wgLock.RUnlock()
	c := &ioOperation{}
	c.ch = make(chan ioResult)
	return c, nil
}


func ioCompletionProcessor(h syscall.Handle) {
	for {
		var bytes uint32
		var key uintptr
		var op *ioOperation
		err := getQueuedCompletionStatus(h, &bytes, &key, &op, syscall.INFINITE)
		if op == nil {
			panic(err)
		}
		op.ch <- ioResult{bytes, err}
	}
}



func (f *win32File) asyncIo(c *ioOperation, d *deadlineHandler, bytes uint32, err error) (int, error) {
	if err != syscall.ERROR_IO_PENDING {
		return int(bytes), err
	}

	if f.closing.isSet() {
		cancelIoEx(f.handle, &c.o)
	}

	var timeout timeoutChan
	if d != nil {
		d.channelLock.Lock()
		timeout = d.channel
		d.channelLock.Unlock()
	}

	var r ioResult
	select {
	case r = <-c.ch:
		err = r.err
		if err == syscall.ERROR_OPERATION_ABORTED {
			if f.closing.isSet() {
				err = ErrFileClosed
			}
		}
	case <-timeout:
		cancelIoEx(f.handle, &c.o)
		r = <-c.ch
		err = r.err
		if err == syscall.ERROR_OPERATION_ABORTED {
			err = ErrTimeout
		}
	}

	
	
	
	runtime.KeepAlive(c)
	return int(r.bytes), err
}


func (f *win32File) Read(b []byte) (int, error) {
	c, err := f.prepareIo()
	if err != nil {
		return 0, err
	}
	defer f.wg.Done()

	if f.readDeadline.timedout.isSet() {
		return 0, ErrTimeout
	}

	var bytes uint32
	err = syscall.ReadFile(f.handle, b, &bytes, &c.o)
	n, err := f.asyncIo(c, &f.readDeadline, bytes, err)
	runtime.KeepAlive(b)

	
	if err == nil && n == 0 && len(b) != 0 {
		return 0, io.EOF
	} else if err == syscall.ERROR_BROKEN_PIPE {
		return 0, io.EOF
	} else {
		return n, err
	}
}


func (f *win32File) Write(b []byte) (int, error) {
	c, err := f.prepareIo()
	if err != nil {
		return 0, err
	}
	defer f.wg.Done()

	if f.writeDeadline.timedout.isSet() {
		return 0, ErrTimeout
	}

	var bytes uint32
	err = syscall.WriteFile(f.handle, b, &bytes, &c.o)
	n, err := f.asyncIo(c, &f.writeDeadline, bytes, err)
	runtime.KeepAlive(b)
	return n, err
}

func (f *win32File) SetReadDeadline(deadline time.Time) error {
	return f.readDeadline.set(deadline)
}

func (f *win32File) SetWriteDeadline(deadline time.Time) error {
	return f.writeDeadline.set(deadline)
}

func (f *win32File) Flush() error {
	return syscall.FlushFileBuffers(f.handle)
}

func (d *deadlineHandler) set(deadline time.Time) error {
	d.setLock.Lock()
	defer d.setLock.Unlock()

	if d.timer != nil {
		if !d.timer.Stop() {
			<-d.channel
		}
		d.timer = nil
	}
	d.timedout.setFalse()

	select {
	case <-d.channel:
		d.channelLock.Lock()
		d.channel = make(chan struct{})
		d.channelLock.Unlock()
	default:
	}

	if deadline.IsZero() {
		return nil
	}

	timeoutIO := func() {
		d.timedout.setTrue()
		close(d.channel)
	}

	now := time.Now()
	duration := deadline.Sub(now)
	if deadline.After(now) {
		
		d.timer = time.AfterFunc(duration, timeoutIO)
	} else {
		
		timeoutIO()
	}
	return nil
}
