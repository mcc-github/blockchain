





package syscall

import (
	"fmt"
	"net"
	"syscall"
	"time"

	"golang.org/x/sys/unix"
	"google.golang.org/grpc/grpclog"
)


func GetCPUTime() int64 {
	var ts unix.Timespec
	if err := unix.ClockGettime(unix.CLOCK_PROCESS_CPUTIME_ID, &ts); err != nil {
		grpclog.Fatal(err)
	}
	return ts.Nano()
}


type Rusage syscall.Rusage


func GetRusage() (rusage *Rusage) {
	rusage = new(Rusage)
	syscall.Getrusage(syscall.RUSAGE_SELF, (*syscall.Rusage)(rusage))
	return
}



func CPUTimeDiff(first *Rusage, latest *Rusage) (float64, float64) {
	f := (*syscall.Rusage)(first)
	l := (*syscall.Rusage)(latest)
	var (
		utimeDiffs  = l.Utime.Sec - f.Utime.Sec
		utimeDiffus = l.Utime.Usec - f.Utime.Usec
		stimeDiffs  = l.Stime.Sec - f.Stime.Sec
		stimeDiffus = l.Stime.Usec - f.Stime.Usec
	)

	uTimeElapsed := float64(utimeDiffs) + float64(utimeDiffus)*1.0e-6
	sTimeElapsed := float64(stimeDiffs) + float64(stimeDiffus)*1.0e-6

	return uTimeElapsed, sTimeElapsed
}


func SetTCPUserTimeout(conn net.Conn, timeout time.Duration) error {
	tcpconn, ok := conn.(*net.TCPConn)
	if !ok {
		
		return nil
	}
	rawConn, err := tcpconn.SyscallConn()
	if err != nil {
		return fmt.Errorf("error getting raw connection: %v", err)
	}
	err = rawConn.Control(func(fd uintptr) {
		err = syscall.SetsockoptInt(int(fd), syscall.IPPROTO_TCP, unix.TCP_USER_TIMEOUT, int(timeout/time.Millisecond))
	})
	if err != nil {
		return fmt.Errorf("error setting option on socket: %v", err)
	}

	return nil
}


func GetTCPUserTimeout(conn net.Conn) (opt int, err error) {
	tcpconn, ok := conn.(*net.TCPConn)
	if !ok {
		err = fmt.Errorf("conn is not *net.TCPConn. got %T", conn)
		return
	}
	rawConn, err := tcpconn.SyscallConn()
	if err != nil {
		err = fmt.Errorf("error getting raw connection: %v", err)
		return
	}
	err = rawConn.Control(func(fd uintptr) {
		opt, err = syscall.GetsockoptInt(int(fd), syscall.IPPROTO_TCP, unix.TCP_USER_TIMEOUT)
	})
	if err != nil {
		err = fmt.Errorf("error getting option on socket: %v", err)
		return
	}

	return
}
