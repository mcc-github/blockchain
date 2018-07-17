



package statsd

import (
	"errors"
	"net"
)


type Sender interface {
	Send(data []byte) (int, error)
	Close() error
}


type SimpleSender struct {
	
	c net.PacketConn
	
	ra *net.UDPAddr
}


func (s *SimpleSender) Send(data []byte) (int, error) {
	
	
	n, err := s.c.(*net.UDPConn).WriteToUDP(data, s.ra)
	if err != nil {
		return 0, err
	}
	if n == 0 {
		return n, errors.New("Wrote no bytes")
	}
	return n, nil
}


func (s *SimpleSender) Close() error {
	err := s.c.Close()
	return err
}






func NewSimpleSender(addr string) (Sender, error) {
	c, err := net.ListenPacket("udp", ":0")
	if err != nil {
		return nil, err
	}

	ra, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		c.Close()
		return nil, err
	}

	sender := &SimpleSender{
		c:  c,
		ra: ra,
	}

	return sender, nil
}
