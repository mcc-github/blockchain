













package snap

import (
	"io"

	"go.etcd.io/etcd/pkg/ioutil"
	"go.etcd.io/etcd/raft/raftpb"
)









type Message struct {
	raftpb.Message
	ReadCloser io.ReadCloser
	TotalSize  int64
	closeC     chan bool
}

func NewMessage(rs raftpb.Message, rc io.ReadCloser, rcSize int64) *Message {
	return &Message{
		Message:    rs,
		ReadCloser: ioutil.NewExactReadCloser(rc, rcSize),
		TotalSize:  int64(rs.Size()) + rcSize,
		closeC:     make(chan bool, 1),
	}
}




func (m Message) CloseNotify() <-chan bool {
	return m.closeC
}

func (m Message) CloseWithError(err error) {
	if cerr := m.ReadCloser.Close(); cerr != nil {
		err = cerr
	}
	if err == nil {
		m.closeC <- true
	} else {
		m.closeC <- false
	}
}
