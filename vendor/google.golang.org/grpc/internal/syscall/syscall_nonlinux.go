



package syscall

import (
	"net"
	"sync"
	"time"

	"google.golang.org/grpc/grpclog"
)

var once sync.Once

func log() {
	once.Do(func() {
		grpclog.Info("CPU time info is unavailable on non-linux or appengine environment.")
	})
}



func GetCPUTime() int64 {
	log()
	return 0
}


type Rusage struct{}


func GetRusage() (rusage *Rusage) {
	log()
	return nil
}



func CPUTimeDiff(first *Rusage, latest *Rusage) (float64, float64) {
	log()
	return 0, 0
}


func SetTCPUserTimeout(conn net.Conn, timeout time.Duration) error {
	log()
	return nil
}



func GetTCPUserTimeout(conn net.Conn) (int, error) {
	log()
	return -1, nil
}
