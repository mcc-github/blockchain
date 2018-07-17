

package winio

import (
	"errors"
	"io"
	"net"
	"os"
	"syscall"
	"time"
	"unsafe"
)









const (
	cERROR_PIPE_BUSY      = syscall.Errno(231)
	cERROR_NO_DATA        = syscall.Errno(232)
	cERROR_PIPE_CONNECTED = syscall.Errno(535)
	cERROR_SEM_TIMEOUT    = syscall.Errno(121)

	cPIPE_ACCESS_DUPLEX            = 0x3
	cFILE_FLAG_FIRST_PIPE_INSTANCE = 0x80000
	cSECURITY_SQOS_PRESENT         = 0x100000
	cSECURITY_ANONYMOUS            = 0

	cPIPE_REJECT_REMOTE_CLIENTS = 0x8

	cPIPE_UNLIMITED_INSTANCES = 255

	cNMPWAIT_USE_DEFAULT_WAIT = 0
	cNMPWAIT_NOWAIT           = 1

	cPIPE_TYPE_MESSAGE = 4

	cPIPE_READMODE_MESSAGE = 2
)

var (
	
	
	ErrPipeListenerClosed = errors.New("use of closed network connection")

	errPipeWriteClosed = errors.New("pipe has been closed for write")
)

type win32Pipe struct {
	*win32File
	path string
}

type win32MessageBytePipe struct {
	win32Pipe
	writeClosed bool
	readEOF     bool
}

type pipeAddress string

func (f *win32Pipe) LocalAddr() net.Addr {
	return pipeAddress(f.path)
}

func (f *win32Pipe) RemoteAddr() net.Addr {
	return pipeAddress(f.path)
}

func (f *win32Pipe) SetDeadline(t time.Time) error {
	f.SetReadDeadline(t)
	f.SetWriteDeadline(t)
	return nil
}


func (f *win32MessageBytePipe) CloseWrite() error {
	if f.writeClosed {
		return errPipeWriteClosed
	}
	err := f.win32File.Flush()
	if err != nil {
		return err
	}
	_, err = f.win32File.Write(nil)
	if err != nil {
		return err
	}
	f.writeClosed = true
	return nil
}



func (f *win32MessageBytePipe) Write(b []byte) (int, error) {
	if f.writeClosed {
		return 0, errPipeWriteClosed
	}
	if len(b) == 0 {
		return 0, nil
	}
	return f.win32File.Write(b)
}



func (f *win32MessageBytePipe) Read(b []byte) (int, error) {
	if f.readEOF {
		return 0, io.EOF
	}
	n, err := f.win32File.Read(b)
	if err == io.EOF {
		
		
		
		
		
		f.readEOF = true
	}
	return n, err
}

func (s pipeAddress) Network() string {
	return "pipe"
}

func (s pipeAddress) String() string {
	return string(s)
}




func DialPipe(path string, timeout *time.Duration) (net.Conn, error) {
	var absTimeout time.Time
	if timeout != nil {
		absTimeout = time.Now().Add(*timeout)
	}
	var err error
	var h syscall.Handle
	for {
		h, err = createFile(path, syscall.GENERIC_READ|syscall.GENERIC_WRITE, 0, nil, syscall.OPEN_EXISTING, syscall.FILE_FLAG_OVERLAPPED|cSECURITY_SQOS_PRESENT|cSECURITY_ANONYMOUS, 0)
		if err != cERROR_PIPE_BUSY {
			break
		}
		now := time.Now()
		var ms uint32
		if absTimeout.IsZero() {
			ms = cNMPWAIT_USE_DEFAULT_WAIT
		} else if now.After(absTimeout) {
			ms = cNMPWAIT_NOWAIT
		} else {
			ms = uint32(absTimeout.Sub(now).Nanoseconds() / 1000 / 1000)
		}
		err = waitNamedPipe(path, ms)
		if err != nil {
			if err == cERROR_SEM_TIMEOUT {
				return nil, ErrTimeout
			}
			break
		}
	}
	if err != nil {
		return nil, &os.PathError{Op: "open", Path: path, Err: err}
	}

	var flags uint32
	err = getNamedPipeInfo(h, &flags, nil, nil, nil)
	if err != nil {
		return nil, err
	}

	var state uint32
	err = getNamedPipeHandleState(h, &state, nil, nil, nil, nil, 0)
	if err != nil {
		return nil, err
	}

	if state&cPIPE_READMODE_MESSAGE != 0 {
		return nil, &os.PathError{Op: "open", Path: path, Err: errors.New("message readmode pipes not supported")}
	}

	f, err := makeWin32File(h)
	if err != nil {
		syscall.Close(h)
		return nil, err
	}

	
	
	if flags&cPIPE_TYPE_MESSAGE != 0 {
		return &win32MessageBytePipe{
			win32Pipe: win32Pipe{win32File: f, path: path},
		}, nil
	}
	return &win32Pipe{win32File: f, path: path}, nil
}

type acceptResponse struct {
	f   *win32File
	err error
}

type win32PipeListener struct {
	firstHandle        syscall.Handle
	path               string
	securityDescriptor []byte
	config             PipeConfig
	acceptCh           chan (chan acceptResponse)
	closeCh            chan int
	doneCh             chan int
}

func makeServerPipeHandle(path string, securityDescriptor []byte, c *PipeConfig, first bool) (syscall.Handle, error) {
	var flags uint32 = cPIPE_ACCESS_DUPLEX | syscall.FILE_FLAG_OVERLAPPED
	if first {
		flags |= cFILE_FLAG_FIRST_PIPE_INSTANCE
	}

	var mode uint32 = cPIPE_REJECT_REMOTE_CLIENTS
	if c.MessageMode {
		mode |= cPIPE_TYPE_MESSAGE
	}

	sa := &syscall.SecurityAttributes{}
	sa.Length = uint32(unsafe.Sizeof(*sa))
	if securityDescriptor != nil {
		len := uint32(len(securityDescriptor))
		sa.SecurityDescriptor = localAlloc(0, len)
		defer localFree(sa.SecurityDescriptor)
		copy((*[0xffff]byte)(unsafe.Pointer(sa.SecurityDescriptor))[:], securityDescriptor)
	}
	h, err := createNamedPipe(path, flags, mode, cPIPE_UNLIMITED_INSTANCES, uint32(c.OutputBufferSize), uint32(c.InputBufferSize), 0, sa)
	if err != nil {
		return 0, &os.PathError{Op: "open", Path: path, Err: err}
	}
	return h, nil
}

func (l *win32PipeListener) makeServerPipe() (*win32File, error) {
	h, err := makeServerPipeHandle(l.path, l.securityDescriptor, &l.config, false)
	if err != nil {
		return nil, err
	}
	f, err := makeWin32File(h)
	if err != nil {
		syscall.Close(h)
		return nil, err
	}
	return f, nil
}

func (l *win32PipeListener) makeConnectedServerPipe() (*win32File, error) {
	p, err := l.makeServerPipe()
	if err != nil {
		return nil, err
	}

	
	ch := make(chan error)
	go func(p *win32File) {
		ch <- connectPipe(p)
	}(p)

	select {
	case err = <-ch:
		if err != nil {
			p.Close()
			p = nil
		}
	case <-l.closeCh:
		
		p.Close()
		p = nil
		err = <-ch
		if err == nil || err == ErrFileClosed {
			err = ErrPipeListenerClosed
		}
	}
	return p, err
}

func (l *win32PipeListener) listenerRoutine() {
	closed := false
	for !closed {
		select {
		case <-l.closeCh:
			closed = true
		case responseCh := <-l.acceptCh:
			var (
				p   *win32File
				err error
			)
			for {
				p, err = l.makeConnectedServerPipe()
				
				
				if err != cERROR_NO_DATA {
					break
				}
			}
			responseCh <- acceptResponse{p, err}
			closed = err == ErrPipeListenerClosed
		}
	}
	syscall.Close(l.firstHandle)
	l.firstHandle = 0
	
	close(l.doneCh)
}


type PipeConfig struct {
	
	SecurityDescriptor string

	
	
	
	
	
	
	MessageMode bool

	
	InputBufferSize int32

	
	OutputBufferSize int32
}



func ListenPipe(path string, c *PipeConfig) (net.Listener, error) {
	var (
		sd  []byte
		err error
	)
	if c == nil {
		c = &PipeConfig{}
	}
	if c.SecurityDescriptor != "" {
		sd, err = SddlToSecurityDescriptor(c.SecurityDescriptor)
		if err != nil {
			return nil, err
		}
	}
	h, err := makeServerPipeHandle(path, sd, c, true)
	if err != nil {
		return nil, err
	}
	
	
	h2, err := createFile(path, 0, 0, nil, syscall.OPEN_EXISTING, cSECURITY_SQOS_PRESENT|cSECURITY_ANONYMOUS, 0)
	if err != nil {
		syscall.Close(h)
		return nil, err
	}
	syscall.Close(h2)
	l := &win32PipeListener{
		firstHandle:        h,
		path:               path,
		securityDescriptor: sd,
		config:             *c,
		acceptCh:           make(chan (chan acceptResponse)),
		closeCh:            make(chan int),
		doneCh:             make(chan int),
	}
	go l.listenerRoutine()
	return l, nil
}

func connectPipe(p *win32File) error {
	c, err := p.prepareIo()
	if err != nil {
		return err
	}
	defer p.wg.Done()

	err = connectNamedPipe(p.handle, &c.o)
	_, err = p.asyncIo(c, nil, 0, err)
	if err != nil && err != cERROR_PIPE_CONNECTED {
		return err
	}
	return nil
}

func (l *win32PipeListener) Accept() (net.Conn, error) {
	ch := make(chan acceptResponse)
	select {
	case l.acceptCh <- ch:
		response := <-ch
		err := response.err
		if err != nil {
			return nil, err
		}
		if l.config.MessageMode {
			return &win32MessageBytePipe{
				win32Pipe: win32Pipe{win32File: response.f, path: l.path},
			}, nil
		}
		return &win32Pipe{win32File: response.f, path: l.path}, nil
	case <-l.doneCh:
		return nil, ErrPipeListenerClosed
	}
}

func (l *win32PipeListener) Close() error {
	select {
	case l.closeCh <- 1:
		<-l.doneCh
	case <-l.doneCh:
	}
	return nil
}

func (l *win32PipeListener) Addr() net.Addr {
	return pipeAddress(l.path)
}
