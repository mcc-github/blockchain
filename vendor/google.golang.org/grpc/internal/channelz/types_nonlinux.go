



package channelz

import "google.golang.org/grpc/grpclog"

func init() {
	grpclog.Infof("Channelz: socket options are not supported on non-linux os and appengine.")
}




type SocketOptionData struct {
}




func (s *SocketOptionData) Getsockopt(fd uintptr) {}
