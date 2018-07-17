package grouper

import (
	"sync"

	"github.com/tedsuo/ifrit"
)


type DynamicClient interface {

	
	EntranceListener() <-chan EntranceEvent

	
	ExitListener() <-chan ExitEvent

	
	CloseNotifier() <-chan struct{}
	
	Inserter() chan<- Member

	
	Close()

	Get(name string) (ifrit.Process, bool)
}

type memberRequest struct {
	Name     string
	Response chan ifrit.Process
}


type dynamicClient struct {
	insertChannel       chan Member
	getMemberChannel    chan memberRequest
	completeNotifier    chan struct{}
	closeNotifier       chan struct{}
	closeOnce           *sync.Once
	entranceBroadcaster *entranceEventBroadcaster
	exitBroadcaster     *exitEventBroadcaster
}

func newClient(bufferSize int) dynamicClient {
	return dynamicClient{
		insertChannel:       make(chan Member),
		getMemberChannel:    make(chan memberRequest),
		completeNotifier:    make(chan struct{}),
		closeNotifier:       make(chan struct{}),
		closeOnce:           new(sync.Once),
		entranceBroadcaster: newEntranceEventBroadcaster(bufferSize),
		exitBroadcaster:     newExitEventBroadcaster(bufferSize),
	}
}

func (c dynamicClient) Get(name string) (ifrit.Process, bool) {
	req := memberRequest{
		Name:     name,
		Response: make(chan ifrit.Process),
	}
	select {
	case c.getMemberChannel <- req:
		p, ok := <-req.Response
		if !ok {
			return nil, false
		}
		return p, true
	case <-c.completeNotifier:
		return nil, false
	}
}

func (c dynamicClient) memberRequests() chan memberRequest {
	return c.getMemberChannel
}

func (c dynamicClient) Inserter() chan<- Member {
	return c.insertChannel
}

func (c dynamicClient) insertEventListener() <-chan Member {
	return c.insertChannel
}

func (c dynamicClient) EntranceListener() <-chan EntranceEvent {
	return c.entranceBroadcaster.Attach()
}

func (c dynamicClient) broadcastEntrance(event EntranceEvent) {
	c.entranceBroadcaster.Broadcast(event)
}

func (c dynamicClient) closeEntranceBroadcaster() {
	c.entranceBroadcaster.Close()
}

func (c dynamicClient) ExitListener() <-chan ExitEvent {
	return c.exitBroadcaster.Attach()
}

func (c dynamicClient) broadcastExit(event ExitEvent) {
	c.exitBroadcaster.Broadcast(event)
}

func (c dynamicClient) closeExitBroadcaster() {
	c.exitBroadcaster.Close()
}

func (c dynamicClient) closeBroadcasters() error {
	c.entranceBroadcaster.Close()
	c.exitBroadcaster.Close()
	close(c.completeNotifier)
	return nil
}

func (c dynamicClient) Close() {
	c.closeOnce.Do(func() {
		close(c.closeNotifier)
	})
}

func (c dynamicClient) CloseNotifier() <-chan struct{} {
	return c.closeNotifier
}
