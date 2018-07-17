package sarama

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"reflect"
	"strconv"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
)

const (
	expectationTimeout = 500 * time.Millisecond
)

type requestHandlerFunc func(req *request) (res encoder)



type RequestNotifierFunc func(bytesRead, bytesWritten int)

























type MockBroker struct {
	brokerID     int32
	port         int32
	closing      chan none
	stopper      chan none
	expectations chan encoder
	listener     net.Listener
	t            TestReporter
	latency      time.Duration
	handler      requestHandlerFunc
	notifier     RequestNotifierFunc
	history      []RequestResponse
	lock         sync.Mutex
}


type RequestResponse struct {
	Request  protocolBody
	Response encoder
}



func (b *MockBroker) SetLatency(latency time.Duration) {
	b.latency = latency
}





func (b *MockBroker) SetHandlerByMap(handlerMap map[string]MockResponse) {
	b.setHandler(func(req *request) (res encoder) {
		reqTypeName := reflect.TypeOf(req.body).Elem().Name()
		mockResponse := handlerMap[reqTypeName]
		if mockResponse == nil {
			return nil
		}
		return mockResponse.For(req.body)
	})
}



func (b *MockBroker) SetNotifier(notifier RequestNotifierFunc) {
	b.lock.Lock()
	b.notifier = notifier
	b.lock.Unlock()
}


func (b *MockBroker) BrokerID() int32 {
	return b.brokerID
}





func (b *MockBroker) History() []RequestResponse {
	b.lock.Lock()
	history := make([]RequestResponse, len(b.history))
	copy(history, b.history)
	b.lock.Unlock()
	return history
}


func (b *MockBroker) Port() int32 {
	return b.port
}


func (b *MockBroker) Addr() string {
	return b.listener.Addr().String()
}



func (b *MockBroker) Close() {
	close(b.expectations)
	if len(b.expectations) > 0 {
		buf := bytes.NewBufferString(fmt.Sprintf("mockbroker/%d: not all expectations were satisfied! Still waiting on:\n", b.BrokerID()))
		for e := range b.expectations {
			_, _ = buf.WriteString(spew.Sdump(e))
		}
		b.t.Error(buf.String())
	}
	close(b.closing)
	<-b.stopper
}




func (b *MockBroker) setHandler(handler requestHandlerFunc) {
	b.lock.Lock()
	b.handler = handler
	b.lock.Unlock()
}

func (b *MockBroker) serverLoop() {
	defer close(b.stopper)
	var err error
	var conn net.Conn

	go func() {
		<-b.closing
		err := b.listener.Close()
		if err != nil {
			b.t.Error(err)
		}
	}()

	wg := &sync.WaitGroup{}
	i := 0
	for conn, err = b.listener.Accept(); err == nil; conn, err = b.listener.Accept() {
		wg.Add(1)
		go b.handleRequests(conn, i, wg)
		i++
	}
	wg.Wait()
	Logger.Printf("*** mockbroker/%d: listener closed, err=%v", b.BrokerID(), err)
}

func (b *MockBroker) handleRequests(conn net.Conn, idx int, wg *sync.WaitGroup) {
	defer wg.Done()
	defer func() {
		_ = conn.Close()
	}()
	Logger.Printf("*** mockbroker/%d/%d: connection opened", b.BrokerID(), idx)
	var err error

	abort := make(chan none)
	defer close(abort)
	go func() {
		select {
		case <-b.closing:
			_ = conn.Close()
		case <-abort:
		}
	}()

	resHeader := make([]byte, 8)
	for {
		req, bytesRead, err := decodeRequest(conn)
		if err != nil {
			Logger.Printf("*** mockbroker/%d/%d: invalid request: err=%+v, %+v", b.brokerID, idx, err, spew.Sdump(req))
			b.serverError(err)
			break
		}

		if b.latency > 0 {
			time.Sleep(b.latency)
		}

		b.lock.Lock()
		res := b.handler(req)
		b.history = append(b.history, RequestResponse{req.body, res})
		b.lock.Unlock()

		if res == nil {
			Logger.Printf("*** mockbroker/%d/%d: ignored %v", b.brokerID, idx, spew.Sdump(req))
			continue
		}
		Logger.Printf("*** mockbroker/%d/%d: served %v -> %v", b.brokerID, idx, req, res)

		encodedRes, err := encode(res, nil)
		if err != nil {
			b.serverError(err)
			break
		}
		if len(encodedRes) == 0 {
			b.lock.Lock()
			if b.notifier != nil {
				b.notifier(bytesRead, 0)
			}
			b.lock.Unlock()
			continue
		}

		binary.BigEndian.PutUint32(resHeader, uint32(len(encodedRes)+4))
		binary.BigEndian.PutUint32(resHeader[4:], uint32(req.correlationID))
		if _, err = conn.Write(resHeader); err != nil {
			b.serverError(err)
			break
		}
		if _, err = conn.Write(encodedRes); err != nil {
			b.serverError(err)
			break
		}

		b.lock.Lock()
		if b.notifier != nil {
			b.notifier(bytesRead, len(resHeader)+len(encodedRes))
		}
		b.lock.Unlock()
	}
	Logger.Printf("*** mockbroker/%d/%d: connection closed, err=%v", b.BrokerID(), idx, err)
}

func (b *MockBroker) defaultRequestHandler(req *request) (res encoder) {
	select {
	case res, ok := <-b.expectations:
		if !ok {
			return nil
		}
		return res
	case <-time.After(expectationTimeout):
		return nil
	}
}

func (b *MockBroker) serverError(err error) {
	isConnectionClosedError := false
	if _, ok := err.(*net.OpError); ok {
		isConnectionClosedError = true
	} else if err == io.EOF {
		isConnectionClosedError = true
	} else if err.Error() == "use of closed network connection" {
		isConnectionClosedError = true
	}

	if isConnectionClosedError {
		return
	}

	b.t.Errorf(err.Error())
}




func NewMockBroker(t TestReporter, brokerID int32) *MockBroker {
	return NewMockBrokerAddr(t, brokerID, "localhost:0")
}



func NewMockBrokerAddr(t TestReporter, brokerID int32, addr string) *MockBroker {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		t.Fatal(err)
	}
	return NewMockBrokerListener(t, brokerID, listener)
}


func NewMockBrokerListener(t TestReporter, brokerID int32, listener net.Listener) *MockBroker {
	var err error

	broker := &MockBroker{
		closing:      make(chan none),
		stopper:      make(chan none),
		t:            t,
		brokerID:     brokerID,
		expectations: make(chan encoder, 512),
		listener:     listener,
	}
	broker.handler = broker.defaultRequestHandler

	Logger.Printf("*** mockbroker/%d listening on %s\n", brokerID, broker.listener.Addr().String())
	_, portStr, err := net.SplitHostPort(broker.listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	tmp, err := strconv.ParseInt(portStr, 10, 32)
	if err != nil {
		t.Fatal(err)
	}
	broker.port = int32(tmp)

	go broker.serverLoop()

	return broker
}

func (b *MockBroker) Returns(e encoder) {
	b.expectations <- e
}
