



package http2

import "fmt"



type WriteScheduler interface {
	
	
	
	OpenStream(streamID uint32, options OpenStreamOptions)

	
	
	
	CloseStream(streamID uint32)

	
	
	
	
	AdjustStream(streamID uint32, priority PriorityParam)

	
	
	
	Push(wr FrameWriteRequest)

	
	
	
	Pop() (wr FrameWriteRequest, ok bool)
}


type OpenStreamOptions struct {
	
	
	PusherID uint32
}


type FrameWriteRequest struct {
	
	
	
	write writeFramer

	
	
	stream *stream

	
	
	
	done chan error
}



func (wr FrameWriteRequest) StreamID() uint32 {
	if wr.stream == nil {
		if se, ok := wr.write.(StreamError); ok {
			
			
			
			
			return se.StreamID
		}
		return 0
	}
	return wr.stream.id
}



func (wr FrameWriteRequest) DataSize() int {
	if wd, ok := wr.write.(*writeData); ok {
		return len(wd.p)
	}
	return 0
}











func (wr FrameWriteRequest) Consume(n int32) (FrameWriteRequest, FrameWriteRequest, int) {
	var empty FrameWriteRequest

	
	wd, ok := wr.write.(*writeData)
	if !ok || len(wd.p) == 0 {
		return wr, empty, 1
	}

	
	allowed := wr.stream.flow.available()
	if n < allowed {
		allowed = n
	}
	if wr.stream.sc.maxFrameSize < allowed {
		allowed = wr.stream.sc.maxFrameSize
	}
	if allowed <= 0 {
		return empty, empty, 0
	}
	if len(wd.p) > int(allowed) {
		wr.stream.flow.take(allowed)
		consumed := FrameWriteRequest{
			stream: wr.stream,
			write: &writeData{
				streamID: wd.streamID,
				p:        wd.p[:allowed],
				
				
				
				endStream: false,
			},
			
			
			done: nil,
		}
		rest := FrameWriteRequest{
			stream: wr.stream,
			write: &writeData{
				streamID:  wd.streamID,
				p:         wd.p[allowed:],
				endStream: wd.endStream,
			},
			done: wr.done,
		}
		return consumed, rest, 2
	}

	
	
	wr.stream.flow.take(int32(len(wd.p)))
	return wr, empty, 1
}


func (wr FrameWriteRequest) String() string {
	var des string
	if s, ok := wr.write.(fmt.Stringer); ok {
		des = s.String()
	} else {
		des = fmt.Sprintf("%T", wr.write)
	}
	return fmt.Sprintf("[FrameWriteRequest stream=%d, ch=%v, writer=%v]", wr.StreamID(), wr.done != nil, des)
}



func (wr *FrameWriteRequest) replyToWriter(err error) {
	if wr.done == nil {
		return
	}
	select {
	case wr.done <- err:
	default:
		panic(fmt.Sprintf("unbuffered done channel passed in for type %T", wr.write))
	}
	wr.write = nil 
}


type writeQueue struct {
	s []FrameWriteRequest
}

func (q *writeQueue) empty() bool { return len(q.s) == 0 }

func (q *writeQueue) push(wr FrameWriteRequest) {
	q.s = append(q.s, wr)
}

func (q *writeQueue) shift() FrameWriteRequest {
	if len(q.s) == 0 {
		panic("invalid use of queue")
	}
	wr := q.s[0]
	
	copy(q.s, q.s[1:])
	q.s[len(q.s)-1] = FrameWriteRequest{}
	q.s = q.s[:len(q.s)-1]
	return wr
}





func (q *writeQueue) consume(n int32) (FrameWriteRequest, bool) {
	if len(q.s) == 0 {
		return FrameWriteRequest{}, false
	}
	consumed, rest, numresult := q.s[0].Consume(n)
	switch numresult {
	case 0:
		return FrameWriteRequest{}, false
	case 1:
		q.shift()
	case 2:
		q.s[0] = rest
	}
	return consumed, true
}

type writeQueuePool []*writeQueue


func (p *writeQueuePool) put(q *writeQueue) {
	for i := range q.s {
		q.s[i] = FrameWriteRequest{}
	}
	q.s = q.s[:0]
	*p = append(*p, q)
}


func (p *writeQueuePool) get() *writeQueue {
	ln := len(*p)
	if ln == 0 {
		return new(writeQueue)
	}
	x := ln - 1
	q := (*p)[x]
	(*p)[x] = nil
	*p = (*p)[:x]
	return q
}
