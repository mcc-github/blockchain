/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package state

import (
	"crypto/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/mcc-github/blockchain/gossip/util"
	proto "github.com/mcc-github/blockchain/protos/gossip"
	"github.com/stretchr/testify/assert"
)

func init() {
	util.SetupTestLogging()
}

func randomPayloadWithSeqNum(seqNum uint64) (*proto.Payload, error) {
	data := make([]byte, 64)
	_, err := rand.Read(data)
	if err != nil {
		return nil, err
	}
	return &proto.Payload{
		SeqNum: seqNum,
		Data:   data,
	}, nil
}

func TestNewPayloadsBuffer(t *testing.T) {
	payloadsBuffer := NewPayloadsBuffer(10)
	assert.Equal(t, payloadsBuffer.Next(), uint64(10))
}

func TestPayloadsBufferImpl_Push(t *testing.T) {
	buffer := NewPayloadsBuffer(5)

	payload, err := randomPayloadWithSeqNum(4)

	if err != nil {
		t.Fatal("Wasn't able to generate random payload for test")
	}

	t.Log("Pushing new payload into buffer")
	buffer.Push(payload)

	
	
	t.Log("Getting next block sequence number")
	assert.Equal(t, buffer.Next(), uint64(5))
	t.Log("Check block buffer size")
	assert.Equal(t, buffer.Size(), 0)

	
	
	payload, err = randomPayloadWithSeqNum(5)
	if err != nil {
		t.Fatal("Wasn't able to generate random payload for test")
	}

	t.Log("Pushing new payload into buffer")
	buffer.Push(payload)
	t.Log("Getting next block sequence number")
	assert.Equal(t, buffer.Next(), uint64(5))
	t.Log("Check block buffer size")
	assert.Equal(t, buffer.Size(), 1)
}

func TestPayloadsBufferImpl_Ready(t *testing.T) {
	fin := make(chan struct{})
	buffer := NewPayloadsBuffer(1)
	assert.Equal(t, buffer.Next(), uint64(1))

	go func() {
		<-buffer.Ready()
		fin <- struct{}{}
	}()

	time.AfterFunc(100*time.Millisecond, func() {
		payload, err := randomPayloadWithSeqNum(1)

		if err != nil {
			t.Fatal("Wasn't able to generate random payload for test")
		}
		buffer.Push(payload)
	})

	select {
	case <-fin:
		payload := buffer.Pop()
		assert.Equal(t, payload.SeqNum, uint64(1))
	case <-time.After(500 * time.Millisecond):
		t.Fail()
	}
}



func TestPayloadsBufferImpl_ConcurrentPush(t *testing.T) {

	
	
	nextSeqNum := uint64(7)
	concurrency := 10

	buffer := NewPayloadsBuffer(nextSeqNum)
	assert.Equal(t, buffer.Next(), uint64(nextSeqNum))

	startWG := sync.WaitGroup{}
	startWG.Add(1)

	finishWG := sync.WaitGroup{}
	finishWG.Add(concurrency)

	payload, err := randomPayloadWithSeqNum(nextSeqNum)
	assert.NoError(t, err)

	ready := int32(0)
	readyWG := sync.WaitGroup{}
	readyWG.Add(1)
	go func() {
		
		<-buffer.Ready()
		atomic.AddInt32(&ready, 1)
		readyWG.Done()
	}()

	for i := 0; i < concurrency; i++ {
		go func() {
			buffer.Push(payload)
			startWG.Wait()
			finishWG.Done()
		}()
	}
	startWG.Done()
	finishWG.Wait()

	readyWG.Wait()
	assert.Equal(t, int32(1), atomic.LoadInt32(&ready))
	
	assert.Equal(t, 1, buffer.Size())
}


func TestPayloadsBufferImpl_Interleave(t *testing.T) {
	buffer := NewPayloadsBuffer(1)
	assert.Equal(t, buffer.Next(), uint64(1))

	
	
	
	
	
	
	
	
	
	
	

	payload, err := randomPayloadWithSeqNum(1)
	assert.NoError(t, err, "generating random payload failed")
	buffer.Push(payload)

	payload, err = randomPayloadWithSeqNum(2)
	assert.NoError(t, err, "generating random payload failed")
	buffer.Push(payload)

	select {
	case <-buffer.Ready():
	case <-time.After(500 * time.Millisecond):
		t.Error("buffer wasn't ready after 500 ms for first sequence")
	}

	
	for payload := buffer.Pop(); payload != nil; payload = buffer.Pop() {
	}

	
	select {
	case <-buffer.Ready():
		t.Error("buffer should not be ready as no new sequences have come")
	case <-time.After(500 * time.Millisecond):
	}

	
	
	
	payload, err = randomPayloadWithSeqNum(3)
	assert.NoError(t, err, "generating random payload failed")
	buffer.Push(payload)

	select {
	case <-buffer.Ready():
	case <-time.After(500 * time.Millisecond):
		t.Error("buffer wasn't ready after 500 ms for second sequence")
	}
	payload = buffer.Pop()
	assert.NotNil(t, payload, "payload should not be nil")

	

	
	payload, err = randomPayloadWithSeqNum(4)
	assert.NoError(t, err, "generating random payload failed")
	buffer.Push(payload)

	
	payload = buffer.Pop()
	assert.NotNil(t, payload, "payload should not be nil")

	
	payload, err = randomPayloadWithSeqNum(5)
	assert.NoError(t, err, "generating random payload failed")
	buffer.Push(payload)

	
	payload = buffer.Pop()
	assert.NotNil(t, payload, "payload should not be nil")

	
	
	
	select {
	case <-buffer.Ready():
		
		
		
		t.Log("buffer ready (1) -- should be error")
		t.Fail()
	case <-time.After(500 * time.Millisecond):
		t.Log("buffer not ready (1)")
	}
	payload = buffer.Pop()
	t.Logf("payload: %v", payload)
	assert.Nil(t, payload, "payload should be nil")

	select {
	case <-buffer.Ready():
		
		
		
		t.Log("buffer ready (2) -- should be error")
		t.Fail()
	case <-time.After(500 * time.Millisecond):
		t.Log("buffer not ready (2)")
	}
	payload = buffer.Pop()
	assert.Nil(t, payload, "payload should be nil")
	t.Logf("payload: %v", payload)

	select {
	case <-buffer.Ready():
		t.Error("buffer ready (3)")
	case <-time.After(500 * time.Millisecond):
		t.Log("buffer not ready (3) -- good")
	}
}
