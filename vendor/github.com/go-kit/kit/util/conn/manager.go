package conn

import (
	"errors"
	"net"
	"time"

	"github.com/go-kit/kit/log"
)



type Dialer func(network, address string) (net.Conn, error)


type AfterFunc func(time.Duration) <-chan time.Time








type Manager struct {
	dialer  Dialer
	network string
	address string
	after   AfterFunc
	logger  log.Logger

	takec chan net.Conn
	putc  chan error
}





func NewManager(d Dialer, network, address string, after AfterFunc, logger log.Logger) *Manager {
	m := &Manager{
		dialer:  d,
		network: network,
		address: address,
		after:   after,
		logger:  logger,

		takec: make(chan net.Conn),
		putc:  make(chan error),
	}
	go m.loop()
	return m
}



func NewDefaultManager(network, address string, logger log.Logger) *Manager {
	return NewManager(net.Dial, network, address, time.After, logger)
}


func (m *Manager) Take() net.Conn {
	return <-m.takec
}




func (m *Manager) Put(err error) {
	m.putc <- err
}


func (m *Manager) Write(b []byte) (int, error) {
	conn := m.Take()
	if conn == nil {
		return 0, ErrConnectionUnavailable
	}
	n, err := conn.Write(b)
	defer m.Put(err)
	return n, err
}

func (m *Manager) loop() {
	var (
		conn       = dial(m.dialer, m.network, m.address, m.logger) 
		connc      = make(chan net.Conn, 1)
		reconnectc <-chan time.Time 
		backoff    = time.Second
	)

	
	
	
	connc <- conn

	for {
		select {
		case <-reconnectc:
			reconnectc = nil 
			go func() { connc <- dial(m.dialer, m.network, m.address, m.logger) }()

		case conn = <-connc:
			if conn == nil {
				
				backoff = exponential(backoff) 
				reconnectc = m.after(backoff)  
			} else {
				
				backoff = time.Second 
				reconnectc = nil      
			}

		case m.takec <- conn:

		case err := <-m.putc:
			if err != nil && conn != nil {
				m.logger.Log("err", err)
				conn = nil                            
				reconnectc = m.after(time.Nanosecond) 
			}
		}
	}
}

func dial(d Dialer, network, address string, logger log.Logger) net.Conn {
	conn, err := d(network, address)
	if err != nil {
		logger.Log("err", err)
		conn = nil 
	}
	return conn
}

func exponential(d time.Duration) time.Duration {
	d *= 2
	if d > time.Minute {
		d = time.Minute
	}
	return d
}



var ErrConnectionUnavailable = errors.New("connection unavailable")
