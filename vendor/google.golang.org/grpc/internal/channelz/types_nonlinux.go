



package channelz

import (
	"sync"

	"google.golang.org/grpc/grpclog"
)

var once sync.Once




type SocketOptionData struct {
}




func (s *SocketOptionData) Getsockopt(fd uintptr) {
	once.Do(func() {
		grpclog.Warningln("Channelz: socket options are not supported on non-linux os and appengine.")
	})
}
