





package logging

import (
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)





func InitForTesting(level Level) *MemoryBackend {
	Reset()

	memoryBackend := NewMemoryBackend(10240)

	leveledBackend := AddModuleLevel(memoryBackend)
	leveledBackend.SetLevel(level, "")
	SetBackend(leveledBackend)

	timeNow = func() time.Time {
		return time.Unix(0, 0).UTC()
	}
	return memoryBackend
}


type node struct {
	next   *node
	Record *Record
}



func (n *node) Next() *node {
	return n.next
}



type MemoryBackend struct {
	size       int32
	maxSize    int32
	head, tail unsafe.Pointer
}


func NewMemoryBackend(size int) *MemoryBackend {
	return &MemoryBackend{maxSize: int32(size)}
}


func (b *MemoryBackend) Log(level Level, calldepth int, rec *Record) error {
	var size int32

	n := &node{Record: rec}
	np := unsafe.Pointer(n)

	
	
	
	for {
		tailp := b.tail
		swapped := atomic.CompareAndSwapPointer(
			&b.tail,
			tailp,
			np,
		)
		if swapped == true {
			if tailp == nil {
				b.head = np
			} else {
				(*node)(tailp).next = n
			}
			size = atomic.AddInt32(&b.size, 1)
			break
		}
	}

	
	
	
	if b.maxSize > 0 && size > b.maxSize {
		for {
			headp := b.head
			head := (*node)(b.head)
			if head.next == nil {
				break
			}
			swapped := atomic.CompareAndSwapPointer(
				&b.head,
				headp,
				unsafe.Pointer(head.next),
			)
			if swapped == true {
				atomic.AddInt32(&b.size, -1)
				break
			}
		}
	}
	return nil
}






func (b *MemoryBackend) Head() *node {
	return (*node)(b.head)
}

type event int

const (
	eventFlush event = iota
	eventStop
)



type ChannelMemoryBackend struct {
	maxSize    int
	size       int
	incoming   chan *Record
	events     chan event
	mu         sync.Mutex
	running    bool
	flushWg    sync.WaitGroup
	stopWg     sync.WaitGroup
	head, tail *node
}





func NewChannelMemoryBackend(size int) *ChannelMemoryBackend {
	backend := &ChannelMemoryBackend{
		maxSize:  size,
		incoming: make(chan *Record, 1024),
		events:   make(chan event),
	}
	backend.Start()
	return backend
}



func (b *ChannelMemoryBackend) Start() {
	b.mu.Lock()
	defer b.mu.Unlock()

	
	if b.running != true {
		b.running = true
		b.stopWg.Add(1)
		go b.process()
	}
}

func (b *ChannelMemoryBackend) process() {
	defer b.stopWg.Done()
	for {
		select {
		case rec := <-b.incoming:
			b.insertRecord(rec)
		case e := <-b.events:
			switch e {
			case eventStop:
				return
			case eventFlush:
				for len(b.incoming) > 0 {
					b.insertRecord(<-b.incoming)
				}
				b.flushWg.Done()
			}
		}
	}
}

func (b *ChannelMemoryBackend) insertRecord(rec *Record) {
	prev := b.tail
	b.tail = &node{Record: rec}
	if prev == nil {
		b.head = b.tail
	} else {
		prev.next = b.tail
	}

	if b.maxSize > 0 && b.size >= b.maxSize {
		b.head = b.head.next
	} else {
		b.size++
	}
}


func (b *ChannelMemoryBackend) Flush() {
	b.flushWg.Add(1)
	b.events <- eventFlush
	b.flushWg.Wait()
}


func (b *ChannelMemoryBackend) Stop() {
	b.mu.Lock()
	if b.running == true {
		b.running = false
		b.events <- eventStop
	}
	b.mu.Unlock()
	b.stopWg.Wait()
}


func (b *ChannelMemoryBackend) Log(level Level, calldepth int, rec *Record) error {
	b.incoming <- rec
	return nil
}






func (b *ChannelMemoryBackend) Head() *node {
	return b.head
}
