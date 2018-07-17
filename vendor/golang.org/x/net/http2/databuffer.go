



package http2

import (
	"errors"
	"fmt"
	"sync"
)











var (
	dataChunkSizeClasses = []int{
		1 << 10,
		2 << 10,
		4 << 10,
		8 << 10,
		16 << 10,
	}
	dataChunkPools = [...]sync.Pool{
		{New: func() interface{} { return make([]byte, 1<<10) }},
		{New: func() interface{} { return make([]byte, 2<<10) }},
		{New: func() interface{} { return make([]byte, 4<<10) }},
		{New: func() interface{} { return make([]byte, 8<<10) }},
		{New: func() interface{} { return make([]byte, 16<<10) }},
	}
)

func getDataBufferChunk(size int64) []byte {
	i := 0
	for ; i < len(dataChunkSizeClasses)-1; i++ {
		if size <= int64(dataChunkSizeClasses[i]) {
			break
		}
	}
	return dataChunkPools[i].Get().([]byte)
}

func putDataBufferChunk(p []byte) {
	for i, n := range dataChunkSizeClasses {
		if len(p) == n {
			dataChunkPools[i].Put(p)
			return
		}
	}
	panic(fmt.Sprintf("unexpected buffer len=%v", len(p)))
}






type dataBuffer struct {
	chunks   [][]byte
	r        int   
	w        int   
	size     int   
	expected int64 
}

var errReadEmpty = errors.New("read from empty dataBuffer")



func (b *dataBuffer) Read(p []byte) (int, error) {
	if b.size == 0 {
		return 0, errReadEmpty
	}
	var ntotal int
	for len(p) > 0 && b.size > 0 {
		readFrom := b.bytesFromFirstChunk()
		n := copy(p, readFrom)
		p = p[n:]
		ntotal += n
		b.r += n
		b.size -= n
		
		if b.r == len(b.chunks[0]) {
			putDataBufferChunk(b.chunks[0])
			end := len(b.chunks) - 1
			copy(b.chunks[:end], b.chunks[1:])
			b.chunks[end] = nil
			b.chunks = b.chunks[:end]
			b.r = 0
		}
	}
	return ntotal, nil
}

func (b *dataBuffer) bytesFromFirstChunk() []byte {
	if len(b.chunks) == 1 {
		return b.chunks[0][b.r:b.w]
	}
	return b.chunks[0][b.r:]
}


func (b *dataBuffer) Len() int {
	return b.size
}


func (b *dataBuffer) Write(p []byte) (int, error) {
	ntotal := len(p)
	for len(p) > 0 {
		
		
		
		want := int64(len(p))
		if b.expected > want {
			want = b.expected
		}
		chunk := b.lastChunkOrAlloc(want)
		n := copy(chunk[b.w:], p)
		p = p[n:]
		b.w += n
		b.size += n
		b.expected -= int64(n)
	}
	return ntotal, nil
}

func (b *dataBuffer) lastChunkOrAlloc(want int64) []byte {
	if len(b.chunks) != 0 {
		last := b.chunks[len(b.chunks)-1]
		if b.w < len(last) {
			return last
		}
	}
	chunk := getDataBufferChunk(want)
	b.chunks = append(b.chunks, chunk)
	b.w = 0
	return chunk
}
