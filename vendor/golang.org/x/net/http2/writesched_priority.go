



package http2

import (
	"fmt"
	"math"
	"sort"
)


const priorityDefaultWeight = 15 


type PriorityWriteSchedulerConfig struct {
	
	
	
	
	
	
	
	
	
	
	
	
	MaxClosedNodesInTree int

	
	
	
	
	
	
	
	
	
	
	MaxIdleNodesInTree int

	
	
	
	
	
	
	
	
	ThrottleOutOfOrderWrites bool
}




func NewPriorityWriteScheduler(cfg *PriorityWriteSchedulerConfig) WriteScheduler {
	if cfg == nil {
		
		
		cfg = &PriorityWriteSchedulerConfig{
			MaxClosedNodesInTree:     10,
			MaxIdleNodesInTree:       10,
			ThrottleOutOfOrderWrites: false,
		}
	}

	ws := &priorityWriteScheduler{
		nodes:                make(map[uint32]*priorityNode),
		maxClosedNodesInTree: cfg.MaxClosedNodesInTree,
		maxIdleNodesInTree:   cfg.MaxIdleNodesInTree,
		enableWriteThrottle:  cfg.ThrottleOutOfOrderWrites,
	}
	ws.nodes[0] = &ws.root
	if cfg.ThrottleOutOfOrderWrites {
		ws.writeThrottleLimit = 1024
	} else {
		ws.writeThrottleLimit = math.MaxInt32
	}
	return ws
}

type priorityNodeState int

const (
	priorityNodeOpen priorityNodeState = iota
	priorityNodeClosed
	priorityNodeIdle
)




type priorityNode struct {
	q            writeQueue        
	id           uint32            
	weight       uint8             
	state        priorityNodeState 
	bytes        int64             
	subtreeBytes int64             

	
	parent     *priorityNode
	kids       *priorityNode 
	prev, next *priorityNode 
}

func (n *priorityNode) setParent(parent *priorityNode) {
	if n == parent {
		panic("setParent to self")
	}
	if n.parent == parent {
		return
	}
	
	if parent := n.parent; parent != nil {
		if n.prev == nil {
			parent.kids = n.next
		} else {
			n.prev.next = n.next
		}
		if n.next != nil {
			n.next.prev = n.prev
		}
	}
	
	
	
	n.parent = parent
	if parent == nil {
		n.next = nil
		n.prev = nil
	} else {
		n.next = parent.kids
		n.prev = nil
		if n.next != nil {
			n.next.prev = n
		}
		parent.kids = n
	}
}

func (n *priorityNode) addBytes(b int64) {
	n.bytes += b
	for ; n != nil; n = n.parent {
		n.subtreeBytes += b
	}
}







func (n *priorityNode) walkReadyInOrder(openParent bool, tmp *[]*priorityNode, f func(*priorityNode, bool) bool) bool {
	if !n.q.empty() && f(n, openParent) {
		return true
	}
	if n.kids == nil {
		return false
	}

	
	
	if n.id != 0 {
		openParent = openParent || (n.state == priorityNodeOpen)
	}

	
	
	
	w := n.kids.weight
	needSort := false
	for k := n.kids.next; k != nil; k = k.next {
		if k.weight != w {
			needSort = true
			break
		}
	}
	if !needSort {
		for k := n.kids; k != nil; k = k.next {
			if k.walkReadyInOrder(openParent, tmp, f) {
				return true
			}
		}
		return false
	}

	
	
	*tmp = (*tmp)[:0]
	for n.kids != nil {
		*tmp = append(*tmp, n.kids)
		n.kids.setParent(nil)
	}
	sort.Sort(sortPriorityNodeSiblings(*tmp))
	for i := len(*tmp) - 1; i >= 0; i-- {
		(*tmp)[i].setParent(n) 
	}
	for k := n.kids; k != nil; k = k.next {
		if k.walkReadyInOrder(openParent, tmp, f) {
			return true
		}
	}
	return false
}

type sortPriorityNodeSiblings []*priorityNode

func (z sortPriorityNodeSiblings) Len() int      { return len(z) }
func (z sortPriorityNodeSiblings) Swap(i, k int) { z[i], z[k] = z[k], z[i] }
func (z sortPriorityNodeSiblings) Less(i, k int) bool {
	
	
	wi, bi := float64(z[i].weight+1), float64(z[i].subtreeBytes)
	wk, bk := float64(z[k].weight+1), float64(z[k].subtreeBytes)
	if bi == 0 && bk == 0 {
		return wi >= wk
	}
	if bk == 0 {
		return false
	}
	return bi/bk <= wi/wk
}

type priorityWriteScheduler struct {
	
	
	root priorityNode

	
	nodes map[uint32]*priorityNode

	
	maxID uint32

	
	
	
	closedNodes, idleNodes []*priorityNode

	
	maxClosedNodesInTree int
	maxIdleNodesInTree   int
	writeThrottleLimit   int32
	enableWriteThrottle  bool

	
	tmp []*priorityNode

	
	queuePool writeQueuePool
}

func (ws *priorityWriteScheduler) OpenStream(streamID uint32, options OpenStreamOptions) {
	
	if curr := ws.nodes[streamID]; curr != nil {
		if curr.state != priorityNodeIdle {
			panic(fmt.Sprintf("stream %d already opened", streamID))
		}
		curr.state = priorityNodeOpen
		return
	}

	
	
	
	
	parent := ws.nodes[options.PusherID]
	if parent == nil {
		parent = &ws.root
	}
	n := &priorityNode{
		q:      *ws.queuePool.get(),
		id:     streamID,
		weight: priorityDefaultWeight,
		state:  priorityNodeOpen,
	}
	n.setParent(parent)
	ws.nodes[streamID] = n
	if streamID > ws.maxID {
		ws.maxID = streamID
	}
}

func (ws *priorityWriteScheduler) CloseStream(streamID uint32) {
	if streamID == 0 {
		panic("violation of WriteScheduler interface: cannot close stream 0")
	}
	if ws.nodes[streamID] == nil {
		panic(fmt.Sprintf("violation of WriteScheduler interface: unknown stream %d", streamID))
	}
	if ws.nodes[streamID].state != priorityNodeOpen {
		panic(fmt.Sprintf("violation of WriteScheduler interface: stream %d already closed", streamID))
	}

	n := ws.nodes[streamID]
	n.state = priorityNodeClosed
	n.addBytes(-n.bytes)

	q := n.q
	ws.queuePool.put(&q)
	n.q.s = nil
	if ws.maxClosedNodesInTree > 0 {
		ws.addClosedOrIdleNode(&ws.closedNodes, ws.maxClosedNodesInTree, n)
	} else {
		ws.removeNode(n)
	}
}

func (ws *priorityWriteScheduler) AdjustStream(streamID uint32, priority PriorityParam) {
	if streamID == 0 {
		panic("adjustPriority on root")
	}

	
	
	
	n := ws.nodes[streamID]
	if n == nil {
		if streamID <= ws.maxID || ws.maxIdleNodesInTree == 0 {
			return
		}
		ws.maxID = streamID
		n = &priorityNode{
			q:      *ws.queuePool.get(),
			id:     streamID,
			weight: priorityDefaultWeight,
			state:  priorityNodeIdle,
		}
		n.setParent(&ws.root)
		ws.nodes[streamID] = n
		ws.addClosedOrIdleNode(&ws.idleNodes, ws.maxIdleNodesInTree, n)
	}

	
	
	parent := ws.nodes[priority.StreamDep]
	if parent == nil {
		n.setParent(&ws.root)
		n.weight = priorityDefaultWeight
		return
	}

	
	if n == parent {
		return
	}

	
	
	
	
	
	
	
	for x := parent.parent; x != nil; x = x.parent {
		if x == n {
			parent.setParent(n.parent)
			break
		}
	}

	
	
	
	if priority.Exclusive {
		k := parent.kids
		for k != nil {
			next := k.next
			if k != n {
				k.setParent(n)
			}
			k = next
		}
	}

	n.setParent(parent)
	n.weight = priority.Weight
}

func (ws *priorityWriteScheduler) Push(wr FrameWriteRequest) {
	var n *priorityNode
	if id := wr.StreamID(); id == 0 {
		n = &ws.root
	} else {
		n = ws.nodes[id]
		if n == nil {
			
			
			
			
			
			if wr.DataSize() > 0 {
				panic("add DATA on non-open stream")
			}
			n = &ws.root
		}
	}
	n.q.push(wr)
}

func (ws *priorityWriteScheduler) Pop() (wr FrameWriteRequest, ok bool) {
	ws.root.walkReadyInOrder(false, &ws.tmp, func(n *priorityNode, openParent bool) bool {
		limit := int32(math.MaxInt32)
		if openParent {
			limit = ws.writeThrottleLimit
		}
		wr, ok = n.q.consume(limit)
		if !ok {
			return false
		}
		n.addBytes(int64(wr.DataSize()))
		
		
		
		if openParent {
			ws.writeThrottleLimit += 1024
			if ws.writeThrottleLimit < 0 {
				ws.writeThrottleLimit = math.MaxInt32
			}
		} else if ws.enableWriteThrottle {
			ws.writeThrottleLimit = 1024
		}
		return true
	})
	return wr, ok
}

func (ws *priorityWriteScheduler) addClosedOrIdleNode(list *[]*priorityNode, maxSize int, n *priorityNode) {
	if maxSize == 0 {
		return
	}
	if len(*list) == maxSize {
		
		ws.removeNode((*list)[0])
		x := (*list)[1:]
		copy(*list, x)
		*list = (*list)[:len(x)]
	}
	*list = append(*list, n)
}

func (ws *priorityWriteScheduler) removeNode(n *priorityNode) {
	for k := n.kids; k != nil; k = k.next {
		k.setParent(n.parent)
	}
	n.setParent(nil)
	delete(ws.nodes, n.id)
}
