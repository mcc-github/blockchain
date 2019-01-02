



package syscall

import (
	"net"
	"time"

	"google.golang.org/grpc/grpclog"
)

func init() {
	grpclog.Info("CPU time info is unavailable on non-linux or appengine environment.")
}



func GetCPUTime() int64 {
	return 0
}


type Rusage struct{}


func GetRusage() (rusage *Rusage) {
	return nil
}



func CPUTimeDiff(first *Rusage, latest *Rusage) (float64, float64) {
	return 0, 0
}


func SetTCPUserTimeout(conn net.Conn, timeout time.Duration) error {
	return nil
}



func GetTCPUserTimeout(conn net.Conn) (int, error) {
	return -1, nil
}
