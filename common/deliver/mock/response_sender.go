
package mock

import (
	"sync"

	"github.com/mcc-github/blockchain/common/deliver"
	cb "github.com/mcc-github/blockchain/protos/common"
)

type ResponseSender struct {
	SendStatusResponseStub        func(status cb.Status) error
	sendStatusResponseMutex       sync.RWMutex
	sendStatusResponseArgsForCall []struct {
		status cb.Status
	}
	sendStatusResponseReturns struct {
		result1 error
	}
	sendStatusResponseReturnsOnCall map[int]struct {
		result1 error
	}
	SendBlockResponseStub        func(block *cb.Block) error
	sendBlockResponseMutex       sync.RWMutex
	sendBlockResponseArgsForCall []struct {
		block *cb.Block
	}
	sendBlockResponseReturns struct {
		result1 error
	}
	sendBlockResponseReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *ResponseSender) SendStatusResponse(status cb.Status) error {
	fake.sendStatusResponseMutex.Lock()
	ret, specificReturn := fake.sendStatusResponseReturnsOnCall[len(fake.sendStatusResponseArgsForCall)]
	fake.sendStatusResponseArgsForCall = append(fake.sendStatusResponseArgsForCall, struct {
		status cb.Status
	}{status})
	fake.recordInvocation("SendStatusResponse", []interface{}{status})
	fake.sendStatusResponseMutex.Unlock()
	if fake.SendStatusResponseStub != nil {
		return fake.SendStatusResponseStub(status)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.sendStatusResponseReturns.result1
}

func (fake *ResponseSender) SendStatusResponseCallCount() int {
	fake.sendStatusResponseMutex.RLock()
	defer fake.sendStatusResponseMutex.RUnlock()
	return len(fake.sendStatusResponseArgsForCall)
}

func (fake *ResponseSender) SendStatusResponseArgsForCall(i int) cb.Status {
	fake.sendStatusResponseMutex.RLock()
	defer fake.sendStatusResponseMutex.RUnlock()
	return fake.sendStatusResponseArgsForCall[i].status
}

func (fake *ResponseSender) SendStatusResponseReturns(result1 error) {
	fake.SendStatusResponseStub = nil
	fake.sendStatusResponseReturns = struct {
		result1 error
	}{result1}
}

func (fake *ResponseSender) SendStatusResponseReturnsOnCall(i int, result1 error) {
	fake.SendStatusResponseStub = nil
	if fake.sendStatusResponseReturnsOnCall == nil {
		fake.sendStatusResponseReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.sendStatusResponseReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *ResponseSender) SendBlockResponse(block *cb.Block) error {
	fake.sendBlockResponseMutex.Lock()
	ret, specificReturn := fake.sendBlockResponseReturnsOnCall[len(fake.sendBlockResponseArgsForCall)]
	fake.sendBlockResponseArgsForCall = append(fake.sendBlockResponseArgsForCall, struct {
		block *cb.Block
	}{block})
	fake.recordInvocation("SendBlockResponse", []interface{}{block})
	fake.sendBlockResponseMutex.Unlock()
	if fake.SendBlockResponseStub != nil {
		return fake.SendBlockResponseStub(block)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.sendBlockResponseReturns.result1
}

func (fake *ResponseSender) SendBlockResponseCallCount() int {
	fake.sendBlockResponseMutex.RLock()
	defer fake.sendBlockResponseMutex.RUnlock()
	return len(fake.sendBlockResponseArgsForCall)
}

func (fake *ResponseSender) SendBlockResponseArgsForCall(i int) *cb.Block {
	fake.sendBlockResponseMutex.RLock()
	defer fake.sendBlockResponseMutex.RUnlock()
	return fake.sendBlockResponseArgsForCall[i].block
}

func (fake *ResponseSender) SendBlockResponseReturns(result1 error) {
	fake.SendBlockResponseStub = nil
	fake.sendBlockResponseReturns = struct {
		result1 error
	}{result1}
}

func (fake *ResponseSender) SendBlockResponseReturnsOnCall(i int, result1 error) {
	fake.SendBlockResponseStub = nil
	if fake.sendBlockResponseReturnsOnCall == nil {
		fake.sendBlockResponseReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.sendBlockResponseReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *ResponseSender) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.sendStatusResponseMutex.RLock()
	defer fake.sendStatusResponseMutex.RUnlock()
	fake.sendBlockResponseMutex.RLock()
	defer fake.sendBlockResponseMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *ResponseSender) recordInvocation(key string, args []interface{}) {
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

var _ deliver.ResponseSender = new(ResponseSender)