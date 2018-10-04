

package channelz

import (
	"net"
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




type ChannelMetric struct {
	
	ID int64
	
	RefName string
	
	
	ChannelData *ChannelInternalMetric
	
	
	NestedChans map[int64]string
	
	
	SubChans map[int64]string
	
	
	
	
	Sockets map[int64]string
}




type SubChannelMetric struct {
	
	ID int64
	
	RefName string
	
	
	ChannelData *ChannelInternalMetric
	
	
	
	
	NestedChans map[int64]string
	
	
	
	
	SubChans map[int64]string
	
	
	Sockets map[int64]string
}



type ChannelInternalMetric struct {
	
	State connectivity.State
	
	Target string
	
	CallsStarted int64
	
	CallsSucceeded int64
	
	CallsFailed int64
	
	LastCallStartedTimestamp time.Time
	
}



type Channel interface {
	ChannelzMetric() *ChannelInternalMetric
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

func (c *channel) deleteSelfIfReady() {
	if !c.closeCalled || len(c.subChans)+len(c.nestedChans) != 0 {
		return
	}
	c.cm.deleteEntry(c.id)
	
	if c.pid != 0 {
		c.cm.findEntry(c.pid).deleteChild(c.id)
	}
}

type subChannel struct {
	refName     string
	c           Channel
	closeCalled bool
	sockets     map[int64]string
	id          int64
	pid         int64
	cm          *channelMap
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

func (sc *subChannel) deleteSelfIfReady() {
	if !sc.closeCalled || len(sc.sockets) != 0 {
		return
	}
	sc.cm.deleteEntry(sc.id)
	sc.cm.findEntry(sc.pid).deleteChild(sc.id)
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
