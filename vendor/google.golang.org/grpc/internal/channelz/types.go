

package channelz

import (
	"net"
	"sync"
	"sync/atomic"
	"time"

	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
)


type entry interface {
	
	addChild(id int64, e entry)
	
	deleteChild(id int64)
	
	
	
	triggerDelete()
	
	
	deleteSelfIfReady()
	
	getParentID() int64
}


type dummyEntry struct {
	idNotFound int64
}

func (d *dummyEntry) addChild(id int64, e entry) {
	
	
	
	
	
	
	
	
	grpclog.Infof("attempt to add child of type %T with id %d to a parent (id=%d) that doesn't currently exist", e, id, d.idNotFound)
}

func (d *dummyEntry) deleteChild(id int64) {
	
	
	grpclog.Infof("attempt to delete child with id %d from a parent (id=%d) that doesn't currently exist", id, d.idNotFound)
}

func (d *dummyEntry) triggerDelete() {
	grpclog.Warningf("attempt to delete an entry (id=%d) that doesn't currently exist", d.idNotFound)
}

func (*dummyEntry) deleteSelfIfReady() {
	
}

func (*dummyEntry) getParentID() int64 {
	return 0
}




type ChannelMetric struct {
	
	ID int64
	
	RefName string
	
	
	ChannelData *ChannelInternalMetric
	
	
	NestedChans map[int64]string
	
	
	SubChans map[int64]string
	
	
	
	
	Sockets map[int64]string
	
	Trace *ChannelTrace
}




type SubChannelMetric struct {
	
	ID int64
	
	RefName string
	
	
	ChannelData *ChannelInternalMetric
	
	
	
	
	NestedChans map[int64]string
	
	
	
	
	SubChans map[int64]string
	
	
	Sockets map[int64]string
	
	Trace *ChannelTrace
}



type ChannelInternalMetric struct {
	
	State connectivity.State
	
	Target string
	
	CallsStarted int64
	
	CallsSucceeded int64
	
	CallsFailed int64
	
	LastCallStartedTimestamp time.Time
}


type ChannelTrace struct {
	
	EventNum int64
	
	CreationTime time.Time
	
	
	Events []*TraceEvent
}


type TraceEvent struct {
	
	Desc string
	
	Severity Severity
	
	Timestamp time.Time
	
	
	
	RefID int64
	
	RefName string
	
	RefType RefChannelType
}



type Channel interface {
	ChannelzMetric() *ChannelInternalMetric
}

type dummyChannel struct{}

func (d *dummyChannel) ChannelzMetric() *ChannelInternalMetric {
	return &ChannelInternalMetric{}
}

type channel struct {
	refName     string
	c           Channel
	closeCalled bool
	nestedChans map[int64]string
	subChans    map[int64]string
	id          int64
	pid         int64
	cm          *channelMap
	trace       *channelTrace
	
	
	traceRefCount int32
}

func (c *channel) addChild(id int64, e entry) {
	switch v := e.(type) {
	case *subChannel:
		c.subChans[id] = v.refName
	case *channel:
		c.nestedChans[id] = v.refName
	default:
		grpclog.Errorf("cannot add a child (id = %d) of type %T to a channel", id, e)
	}
}

func (c *channel) deleteChild(id int64) {
	delete(c.subChans, id)
	delete(c.nestedChans, id)
	c.deleteSelfIfReady()
}

func (c *channel) triggerDelete() {
	c.closeCalled = true
	c.deleteSelfIfReady()
}

func (c *channel) getParentID() int64 {
	return c.pid
}








func (c *channel) deleteSelfFromTree() (deleted bool) {
	if !c.closeCalled || len(c.subChans)+len(c.nestedChans) != 0 {
		return false
	}
	
	if c.pid != 0 {
		c.cm.findEntry(c.pid).deleteChild(c.id)
	}
	return true
}













func (c *channel) deleteSelfFromMap() (delete bool) {
	if c.getTraceRefCount() != 0 {
		c.c = &dummyChannel{}
		return false
	}
	return true
}







func (c *channel) deleteSelfIfReady() {
	if !c.deleteSelfFromTree() {
		return
	}
	if !c.deleteSelfFromMap() {
		return
	}
	c.cm.deleteEntry(c.id)
	c.trace.clear()
}

func (c *channel) getChannelTrace() *channelTrace {
	return c.trace
}

func (c *channel) incrTraceRefCount() {
	atomic.AddInt32(&c.traceRefCount, 1)
}

func (c *channel) decrTraceRefCount() {
	atomic.AddInt32(&c.traceRefCount, -1)
}

func (c *channel) getTraceRefCount() int {
	i := atomic.LoadInt32(&c.traceRefCount)
	return int(i)
}

func (c *channel) getRefName() string {
	return c.refName
}

type subChannel struct {
	refName       string
	c             Channel
	closeCalled   bool
	sockets       map[int64]string
	id            int64
	pid           int64
	cm            *channelMap
	trace         *channelTrace
	traceRefCount int32
}

func (sc *subChannel) addChild(id int64, e entry) {
	if v, ok := e.(*normalSocket); ok {
		sc.sockets[id] = v.refName
	} else {
		grpclog.Errorf("cannot add a child (id = %d) of type %T to a subChannel", id, e)
	}
}

func (sc *subChannel) deleteChild(id int64) {
	delete(sc.sockets, id)
	sc.deleteSelfIfReady()
}

func (sc *subChannel) triggerDelete() {
	sc.closeCalled = true
	sc.deleteSelfIfReady()
}

func (sc *subChannel) getParentID() int64 {
	return sc.pid
}








func (sc *subChannel) deleteSelfFromTree() (deleted bool) {
	if !sc.closeCalled || len(sc.sockets) != 0 {
		return false
	}
	sc.cm.findEntry(sc.pid).deleteChild(sc.id)
	return true
}













func (sc *subChannel) deleteSelfFromMap() (delete bool) {
	if sc.getTraceRefCount() != 0 {
		
		sc.c = &dummyChannel{}
		return false
	}
	return true
}







func (sc *subChannel) deleteSelfIfReady() {
	if !sc.deleteSelfFromTree() {
		return
	}
	if !sc.deleteSelfFromMap() {
		return
	}
	sc.cm.deleteEntry(sc.id)
	sc.trace.clear()
}

func (sc *subChannel) getChannelTrace() *channelTrace {
	return sc.trace
}

func (sc *subChannel) incrTraceRefCount() {
	atomic.AddInt32(&sc.traceRefCount, 1)
}

func (sc *subChannel) decrTraceRefCount() {
	atomic.AddInt32(&sc.traceRefCount, -1)
}

func (sc *subChannel) getTraceRefCount() int {
	i := atomic.LoadInt32(&sc.traceRefCount)
	return int(i)
}

func (sc *subChannel) getRefName() string {
	return sc.refName
}



type SocketMetric struct {
	
	ID int64
	
	RefName string
	
	
	SocketData *SocketInternalMetric
}



type SocketInternalMetric struct {
	
	StreamsStarted int64
	
	
	
	StreamsSucceeded int64
	
	
	
	StreamsFailed int64
	
	MessagesSent     int64
	MessagesReceived int64
	
	
	KeepAlivesSent int64
	
	
	LastLocalStreamCreatedTimestamp time.Time
	
	
	LastRemoteStreamCreatedTimestamp time.Time
	
	LastMessageSentTimestamp time.Time
	
	LastMessageReceivedTimestamp time.Time
	
	
	
	LocalFlowControlWindow int64
	
	
	
	RemoteFlowControlWindow int64
	
	LocalAddr net.Addr
	
	RemoteAddr net.Addr
	
	
	RemoteName    string
	SocketOptions *SocketOptionData
	Security      credentials.ChannelzSecurityValue
}



type Socket interface {
	ChannelzMetric() *SocketInternalMetric
}

type listenSocket struct {
	refName string
	s       Socket
	id      int64
	pid     int64
	cm      *channelMap
}

func (ls *listenSocket) addChild(id int64, e entry) {
	grpclog.Errorf("cannot add a child (id = %d) of type %T to a listen socket", id, e)
}

func (ls *listenSocket) deleteChild(id int64) {
	grpclog.Errorf("cannot delete a child (id = %d) from a listen socket", id)
}

func (ls *listenSocket) triggerDelete() {
	ls.cm.deleteEntry(ls.id)
	ls.cm.findEntry(ls.pid).deleteChild(ls.id)
}

func (ls *listenSocket) deleteSelfIfReady() {
	grpclog.Errorf("cannot call deleteSelfIfReady on a listen socket")
}

func (ls *listenSocket) getParentID() int64 {
	return ls.pid
}

type normalSocket struct {
	refName string
	s       Socket
	id      int64
	pid     int64
	cm      *channelMap
}

func (ns *normalSocket) addChild(id int64, e entry) {
	grpclog.Errorf("cannot add a child (id = %d) of type %T to a normal socket", id, e)
}

func (ns *normalSocket) deleteChild(id int64) {
	grpclog.Errorf("cannot delete a child (id = %d) from a normal socket", id)
}

func (ns *normalSocket) triggerDelete() {
	ns.cm.deleteEntry(ns.id)
	ns.cm.findEntry(ns.pid).deleteChild(ns.id)
}

func (ns *normalSocket) deleteSelfIfReady() {
	grpclog.Errorf("cannot call deleteSelfIfReady on a normal socket")
}

func (ns *normalSocket) getParentID() int64 {
	return ns.pid
}




type ServerMetric struct {
	
	ID int64
	
	RefName string
	
	
	ServerData *ServerInternalMetric
	
	
	ListenSockets map[int64]string
}



type ServerInternalMetric struct {
	
	CallsStarted int64
	
	CallsSucceeded int64
	
	CallsFailed int64
	
	LastCallStartedTimestamp time.Time
}



type Server interface {
	ChannelzMetric() *ServerInternalMetric
}

type server struct {
	refName       string
	s             Server
	closeCalled   bool
	sockets       map[int64]string
	listenSockets map[int64]string
	id            int64
	cm            *channelMap
}

func (s *server) addChild(id int64, e entry) {
	switch v := e.(type) {
	case *normalSocket:
		s.sockets[id] = v.refName
	case *listenSocket:
		s.listenSockets[id] = v.refName
	default:
		grpclog.Errorf("cannot add a child (id = %d) of type %T to a server", id, e)
	}
}

func (s *server) deleteChild(id int64) {
	delete(s.sockets, id)
	delete(s.listenSockets, id)
	s.deleteSelfIfReady()
}

func (s *server) triggerDelete() {
	s.closeCalled = true
	s.deleteSelfIfReady()
}

func (s *server) deleteSelfIfReady() {
	if !s.closeCalled || len(s.sockets)+len(s.listenSockets) != 0 {
		return
	}
	s.cm.deleteEntry(s.id)
}

func (s *server) getParentID() int64 {
	return 0
}

type tracedChannel interface {
	getChannelTrace() *channelTrace
	incrTraceRefCount()
	decrTraceRefCount()
	getRefName() string
}

type channelTrace struct {
	cm          *channelMap
	createdTime time.Time
	eventCount  int64
	mu          sync.Mutex
	events      []*TraceEvent
}

func (c *channelTrace) append(e *TraceEvent) {
	c.mu.Lock()
	if len(c.events) == getMaxTraceEntry() {
		del := c.events[0]
		c.events = c.events[1:]
		if del.RefID != 0 {
			
			go func() {
				
				c.cm.mu.Lock()
				c.cm.decrTraceRefCount(del.RefID)
				c.cm.mu.Unlock()
			}()
		}
	}
	e.Timestamp = time.Now()
	c.events = append(c.events, e)
	c.eventCount++
	c.mu.Unlock()
}

func (c *channelTrace) clear() {
	c.mu.Lock()
	for _, e := range c.events {
		if e.RefID != 0 {
			
			c.cm.decrTraceRefCount(e.RefID)
		}
	}
	c.mu.Unlock()
}




type Severity int

const (
	
	CtUNKNOWN Severity = iota
	
	CtINFO
	
	CtWarning
	
	CtError
)


type RefChannelType int

const (
	
	RefChannel RefChannelType = iota
	
	RefSubChannel
)

func (c *channelTrace) dumpData() *ChannelTrace {
	c.mu.Lock()
	ct := &ChannelTrace{EventNum: c.eventCount, CreationTime: c.createdTime}
	ct.Events = c.events[:len(c.events)]
	c.mu.Unlock()
	return ct
}
