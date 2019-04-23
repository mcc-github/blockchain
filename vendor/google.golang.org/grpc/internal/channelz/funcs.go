






package channelz

import (
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"google.golang.org/grpc/grpclog"
)

const (
	defaultMaxTraceEntry int32 = 30
)

var (
	db    dbWrapper
	idGen idGenerator
	
	EntryPerPage  = int64(50)
	curState      int32
	maxTraceEntry = defaultMaxTraceEntry
)


func TurnOn() {
	if !IsOn() {
		NewChannelzStorage()
		atomic.StoreInt32(&curState, 1)
	}
}


func IsOn() bool {
	return atomic.CompareAndSwapInt32(&curState, 1, 1)
}



func SetMaxTraceEntry(i int32) {
	atomic.StoreInt32(&maxTraceEntry, i)
}


func ResetMaxTraceEntryToDefault() {
	atomic.StoreInt32(&maxTraceEntry, defaultMaxTraceEntry)
}

func getMaxTraceEntry() int {
	i := atomic.LoadInt32(&maxTraceEntry)
	return int(i)
}



type dbWrapper struct {
	mu sync.RWMutex
	DB *channelMap
}

func (d *dbWrapper) set(db *channelMap) {
	d.mu.Lock()
	d.DB = db
	d.mu.Unlock()
}

func (d *dbWrapper) get() *channelMap {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.DB
}





func NewChannelzStorage() {
	db.set(&channelMap{
		topLevelChannels: make(map[int64]struct{}),
		channels:         make(map[int64]*channel),
		listenSockets:    make(map[int64]*listenSocket),
		normalSockets:    make(map[int64]*normalSocket),
		servers:          make(map[int64]*server),
		subChannels:      make(map[int64]*subChannel),
	})
	idGen.reset()
}







func GetTopChannels(id int64, maxResults int64) ([]*ChannelMetric, bool) {
	return db.get().GetTopChannels(id, maxResults)
}







func GetServers(id int64, maxResults int64) ([]*ServerMetric, bool) {
	return db.get().GetServers(id, maxResults)
}








func GetServerSockets(id int64, startID int64, maxResults int64) ([]*SocketMetric, bool) {
	return db.get().GetServerSockets(id, startID, maxResults)
}


func GetChannel(id int64) *ChannelMetric {
	return db.get().GetChannel(id)
}


func GetSubChannel(id int64) *SubChannelMetric {
	return db.get().GetSubChannel(id)
}


func GetSocket(id int64) *SocketMetric {
	return db.get().GetSocket(id)
}


func GetServer(id int64) *ServerMetric {
	return db.get().GetServer(id)
}





func RegisterChannel(c Channel, pid int64, ref string) int64 {
	id := idGen.genID()
	cn := &channel{
		refName:     ref,
		c:           c,
		subChans:    make(map[int64]string),
		nestedChans: make(map[int64]string),
		id:          id,
		pid:         pid,
		trace:       &channelTrace{createdTime: time.Now(), events: make([]*TraceEvent, 0, getMaxTraceEntry())},
	}
	if pid == 0 {
		db.get().addChannel(id, cn, true, pid, ref)
	} else {
		db.get().addChannel(id, cn, false, pid, ref)
	}
	return id
}




func RegisterSubChannel(c Channel, pid int64, ref string) int64 {
	if pid == 0 {
		grpclog.Error("a SubChannel's parent id cannot be 0")
		return 0
	}
	id := idGen.genID()
	sc := &subChannel{
		refName: ref,
		c:       c,
		sockets: make(map[int64]string),
		id:      id,
		pid:     pid,
		trace:   &channelTrace{createdTime: time.Now(), events: make([]*TraceEvent, 0, getMaxTraceEntry())},
	}
	db.get().addSubChannel(id, sc, pid, ref)
	return id
}



func RegisterServer(s Server, ref string) int64 {
	id := idGen.genID()
	svr := &server{
		refName:       ref,
		s:             s,
		sockets:       make(map[int64]string),
		listenSockets: make(map[int64]string),
		id:            id,
	}
	db.get().addServer(id, svr)
	return id
}





func RegisterListenSocket(s Socket, pid int64, ref string) int64 {
	if pid == 0 {
		grpclog.Error("a ListenSocket's parent id cannot be 0")
		return 0
	}
	id := idGen.genID()
	ls := &listenSocket{refName: ref, s: s, id: id, pid: pid}
	db.get().addListenSocket(id, ls, pid, ref)
	return id
}





func RegisterNormalSocket(s Socket, pid int64, ref string) int64 {
	if pid == 0 {
		grpclog.Error("a NormalSocket's parent id cannot be 0")
		return 0
	}
	id := idGen.genID()
	ns := &normalSocket{refName: ref, s: s, id: id, pid: pid}
	db.get().addNormalSocket(id, ns, pid, ref)
	return id
}



func RemoveEntry(id int64) {
	db.get().removeEntry(id)
}





type TraceEventDesc struct {
	Desc     string
	Severity Severity
	Parent   *TraceEventDesc
}


func AddTraceEvent(id int64, desc *TraceEventDesc) {
	if getMaxTraceEntry() == 0 {
		return
	}
	db.get().traceEvent(id, desc)
}






type channelMap struct {
	mu               sync.RWMutex
	topLevelChannels map[int64]struct{}
	servers          map[int64]*server
	channels         map[int64]*channel
	subChannels      map[int64]*subChannel
	listenSockets    map[int64]*listenSocket
	normalSockets    map[int64]*normalSocket
}

func (c *channelMap) addServer(id int64, s *server) {
	c.mu.Lock()
	s.cm = c
	c.servers[id] = s
	c.mu.Unlock()
}

func (c *channelMap) addChannel(id int64, cn *channel, isTopChannel bool, pid int64, ref string) {
	c.mu.Lock()
	cn.cm = c
	cn.trace.cm = c
	c.channels[id] = cn
	if isTopChannel {
		c.topLevelChannels[id] = struct{}{}
	} else {
		c.findEntry(pid).addChild(id, cn)
	}
	c.mu.Unlock()
}

func (c *channelMap) addSubChannel(id int64, sc *subChannel, pid int64, ref string) {
	c.mu.Lock()
	sc.cm = c
	sc.trace.cm = c
	c.subChannels[id] = sc
	c.findEntry(pid).addChild(id, sc)
	c.mu.Unlock()
}

func (c *channelMap) addListenSocket(id int64, ls *listenSocket, pid int64, ref string) {
	c.mu.Lock()
	ls.cm = c
	c.listenSockets[id] = ls
	c.findEntry(pid).addChild(id, ls)
	c.mu.Unlock()
}

func (c *channelMap) addNormalSocket(id int64, ns *normalSocket, pid int64, ref string) {
	c.mu.Lock()
	ns.cm = c
	c.normalSockets[id] = ns
	c.findEntry(pid).addChild(id, ns)
	c.mu.Unlock()
}





func (c *channelMap) removeEntry(id int64) {
	c.mu.Lock()
	c.findEntry(id).triggerDelete()
	c.mu.Unlock()
}


func (c *channelMap) decrTraceRefCount(id int64) {
	e := c.findEntry(id)
	if v, ok := e.(tracedChannel); ok {
		v.decrTraceRefCount()
		e.deleteSelfIfReady()
	}
}


func (c *channelMap) findEntry(id int64) entry {
	var v entry
	var ok bool
	if v, ok = c.channels[id]; ok {
		return v
	}
	if v, ok = c.subChannels[id]; ok {
		return v
	}
	if v, ok = c.servers[id]; ok {
		return v
	}
	if v, ok = c.listenSockets[id]; ok {
		return v
	}
	if v, ok = c.normalSockets[id]; ok {
		return v
	}
	return &dummyEntry{idNotFound: id}
}







func (c *channelMap) deleteEntry(id int64) {
	var ok bool
	if _, ok = c.normalSockets[id]; ok {
		delete(c.normalSockets, id)
		return
	}
	if _, ok = c.subChannels[id]; ok {
		delete(c.subChannels, id)
		return
	}
	if _, ok = c.channels[id]; ok {
		delete(c.channels, id)
		delete(c.topLevelChannels, id)
		return
	}
	if _, ok = c.listenSockets[id]; ok {
		delete(c.listenSockets, id)
		return
	}
	if _, ok = c.servers[id]; ok {
		delete(c.servers, id)
		return
	}
}

func (c *channelMap) traceEvent(id int64, desc *TraceEventDesc) {
	c.mu.Lock()
	child := c.findEntry(id)
	childTC, ok := child.(tracedChannel)
	if !ok {
		c.mu.Unlock()
		return
	}
	childTC.getChannelTrace().append(&TraceEvent{Desc: desc.Desc, Severity: desc.Severity, Timestamp: time.Now()})
	if desc.Parent != nil {
		parent := c.findEntry(child.getParentID())
		var chanType RefChannelType
		switch child.(type) {
		case *channel:
			chanType = RefChannel
		case *subChannel:
			chanType = RefSubChannel
		}
		if parentTC, ok := parent.(tracedChannel); ok {
			parentTC.getChannelTrace().append(&TraceEvent{
				Desc:      desc.Parent.Desc,
				Severity:  desc.Parent.Severity,
				Timestamp: time.Now(),
				RefID:     id,
				RefName:   childTC.getRefName(),
				RefType:   chanType,
			})
			childTC.incrTraceRefCount()
		}
	}
	c.mu.Unlock()
}

type int64Slice []int64

func (s int64Slice) Len() int           { return len(s) }
func (s int64Slice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s int64Slice) Less(i, j int) bool { return s[i] < s[j] }

func copyMap(m map[int64]string) map[int64]string {
	n := make(map[int64]string)
	for k, v := range m {
		n[k] = v
	}
	return n
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func (c *channelMap) GetTopChannels(id int64, maxResults int64) ([]*ChannelMetric, bool) {
	if maxResults <= 0 {
		maxResults = EntryPerPage
	}
	c.mu.RLock()
	l := int64(len(c.topLevelChannels))
	ids := make([]int64, 0, l)
	cns := make([]*channel, 0, min(l, maxResults))

	for k := range c.topLevelChannels {
		ids = append(ids, k)
	}
	sort.Sort(int64Slice(ids))
	idx := sort.Search(len(ids), func(i int) bool { return ids[i] >= id })
	count := int64(0)
	var end bool
	var t []*ChannelMetric
	for i, v := range ids[idx:] {
		if count == maxResults {
			break
		}
		if cn, ok := c.channels[v]; ok {
			cns = append(cns, cn)
			t = append(t, &ChannelMetric{
				NestedChans: copyMap(cn.nestedChans),
				SubChans:    copyMap(cn.subChans),
			})
			count++
		}
		if i == len(ids[idx:])-1 {
			end = true
			break
		}
	}
	c.mu.RUnlock()
	if count == 0 {
		end = true
	}

	for i, cn := range cns {
		t[i].ChannelData = cn.c.ChannelzMetric()
		t[i].ID = cn.id
		t[i].RefName = cn.refName
		t[i].Trace = cn.trace.dumpData()
	}
	return t, end
}

func (c *channelMap) GetServers(id, maxResults int64) ([]*ServerMetric, bool) {
	if maxResults <= 0 {
		maxResults = EntryPerPage
	}
	c.mu.RLock()
	l := int64(len(c.servers))
	ids := make([]int64, 0, l)
	ss := make([]*server, 0, min(l, maxResults))
	for k := range c.servers {
		ids = append(ids, k)
	}
	sort.Sort(int64Slice(ids))
	idx := sort.Search(len(ids), func(i int) bool { return ids[i] >= id })
	count := int64(0)
	var end bool
	var s []*ServerMetric
	for i, v := range ids[idx:] {
		if count == maxResults {
			break
		}
		if svr, ok := c.servers[v]; ok {
			ss = append(ss, svr)
			s = append(s, &ServerMetric{
				ListenSockets: copyMap(svr.listenSockets),
			})
			count++
		}
		if i == len(ids[idx:])-1 {
			end = true
			break
		}
	}
	c.mu.RUnlock()
	if count == 0 {
		end = true
	}

	for i, svr := range ss {
		s[i].ServerData = svr.s.ChannelzMetric()
		s[i].ID = svr.id
		s[i].RefName = svr.refName
	}
	return s, end
}

func (c *channelMap) GetServerSockets(id int64, startID int64, maxResults int64) ([]*SocketMetric, bool) {
	if maxResults <= 0 {
		maxResults = EntryPerPage
	}
	var svr *server
	var ok bool
	c.mu.RLock()
	if svr, ok = c.servers[id]; !ok {
		
		c.mu.RUnlock()
		return nil, true
	}
	svrskts := svr.sockets
	l := int64(len(svrskts))
	ids := make([]int64, 0, l)
	sks := make([]*normalSocket, 0, min(l, maxResults))
	for k := range svrskts {
		ids = append(ids, k)
	}
	sort.Sort(int64Slice(ids))
	idx := sort.Search(len(ids), func(i int) bool { return ids[i] >= startID })
	count := int64(0)
	var end bool
	for i, v := range ids[idx:] {
		if count == maxResults {
			break
		}
		if ns, ok := c.normalSockets[v]; ok {
			sks = append(sks, ns)
			count++
		}
		if i == len(ids[idx:])-1 {
			end = true
			break
		}
	}
	c.mu.RUnlock()
	if count == 0 {
		end = true
	}
	var s []*SocketMetric
	for _, ns := range sks {
		sm := &SocketMetric{}
		sm.SocketData = ns.s.ChannelzMetric()
		sm.ID = ns.id
		sm.RefName = ns.refName
		s = append(s, sm)
	}
	return s, end
}

func (c *channelMap) GetChannel(id int64) *ChannelMetric {
	cm := &ChannelMetric{}
	var cn *channel
	var ok bool
	c.mu.RLock()
	if cn, ok = c.channels[id]; !ok {
		
		c.mu.RUnlock()
		return nil
	}
	cm.NestedChans = copyMap(cn.nestedChans)
	cm.SubChans = copyMap(cn.subChans)
	
	
	chanCopy := cn.c
	c.mu.RUnlock()
	cm.ChannelData = chanCopy.ChannelzMetric()
	cm.ID = cn.id
	cm.RefName = cn.refName
	cm.Trace = cn.trace.dumpData()
	return cm
}

func (c *channelMap) GetSubChannel(id int64) *SubChannelMetric {
	cm := &SubChannelMetric{}
	var sc *subChannel
	var ok bool
	c.mu.RLock()
	if sc, ok = c.subChannels[id]; !ok {
		
		c.mu.RUnlock()
		return nil
	}
	cm.Sockets = copyMap(sc.sockets)
	
	
	chanCopy := sc.c
	c.mu.RUnlock()
	cm.ChannelData = chanCopy.ChannelzMetric()
	cm.ID = sc.id
	cm.RefName = sc.refName
	cm.Trace = sc.trace.dumpData()
	return cm
}

func (c *channelMap) GetSocket(id int64) *SocketMetric {
	sm := &SocketMetric{}
	c.mu.RLock()
	if ls, ok := c.listenSockets[id]; ok {
		c.mu.RUnlock()
		sm.SocketData = ls.s.ChannelzMetric()
		sm.ID = ls.id
		sm.RefName = ls.refName
		return sm
	}
	if ns, ok := c.normalSockets[id]; ok {
		c.mu.RUnlock()
		sm.SocketData = ns.s.ChannelzMetric()
		sm.ID = ns.id
		sm.RefName = ns.refName
		return sm
	}
	c.mu.RUnlock()
	return nil
}

func (c *channelMap) GetServer(id int64) *ServerMetric {
	sm := &ServerMetric{}
	var svr *server
	var ok bool
	c.mu.RLock()
	if svr, ok = c.servers[id]; !ok {
		c.mu.RUnlock()
		return nil
	}
	sm.ListenSockets = copyMap(svr.listenSockets)
	c.mu.RUnlock()
	sm.ID = svr.id
	sm.RefName = svr.refName
	sm.ServerData = svr.s.ChannelzMetric()
	return sm
}

type idGenerator struct {
	id int64
}

func (i *idGenerator) reset() {
	atomic.StoreInt64(&i.id, 0)
}

func (i *idGenerator) genID() int64 {
	return atomic.AddInt64(&i.id, 1)
}
