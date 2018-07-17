
package mock

import (
	"sync"

	"github.com/mcc-github/blockchain/common/deliver"
	cb "github.com/mcc-github/blockchain/protos/common"
)

type Receiver struct {
	RecvStub        func() (*cb.Envelope, error)
	recvMutex       sync.RWMutex
	recvArgsForCall []struct{}
	recvReturns     struct {
		result1 *cb.Envelope
		result2 error
	}
	recvReturnsOnCall map[int]struct {
		result1 *cb.Envelope
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *Receiver) Recv() (*cb.Envelope, error) {
	fake.recvMutex.Lock()
	ret, specificReturn := fake.recvReturnsOnCall[len(fake.recvArgsForCall)]
	fake.recvArgsForCall = append(fake.recvArgsForCall, struct{}{})
	fake.recordInvocation("Recv", []interface{}{})
	fake.recvMutex.Unlock()
	if fake.RecvStub != nil {
		return fake.RecvStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.recvReturns.result1, fake.recvReturns.result2
}

func (fake *Receiver) RecvCallCount() int {
	fake.recvMutex.RLock()
	defer fake.recvMutex.RUnlock()
	return len(fake.recvArgsForCall)
}

func (fake *Receiver) RecvReturns(result1 *cb.Envelope, result2 error) {
	fake.RecvStub = nil
	fake.recvReturns = struct {
		result1 *cb.Envelope
		result2 error
	}{result1, result2}
}

func (fake *Receiver) RecvReturnsOnCall(i int, result1 *cb.Envelope, result2 error) {
	fake.RecvStub = nil
	if fake.recvReturnsOnCall == nil {
		fake.recvReturnsOnCall = make(map[int]struct {
			result1 *cb.Envelope
			result2 error
		})
	}
	fake.recvReturnsOnCall[i] = struct {
		result1 *cb.Envelope
		result2 error
	}{result1, result2}
}

func (fake *Receiver) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.recvMutex.RLock()
	defer fake.recvMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *Receiver) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ deliver.Receiver = new(Receiver)
