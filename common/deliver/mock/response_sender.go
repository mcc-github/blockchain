
package mock

import (
	sync "sync"

	deliver "github.com/mcc-github/blockchain/common/deliver"
	common "github.com/mcc-github/blockchain/protos/common"
)

type ResponseSender struct {
	SendBlockResponseStub        func(*common.Block) error
	sendBlockResponseMutex       sync.RWMutex
	sendBlockResponseArgsForCall []struct {
		arg1 *common.Block
	}
	sendBlockResponseReturns struct {
		result1 error
	}
	sendBlockResponseReturnsOnCall map[int]struct {
		result1 error
	}
	SendStatusResponseStub        func(common.Status) error
	sendStatusResponseMutex       sync.RWMutex
	sendStatusResponseArgsForCall []struct {
		arg1 common.Status
	}
	sendStatusResponseReturns struct {
		result1 error
	}
	sendStatusResponseReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *ResponseSender) SendBlockResponse(arg1 *common.Block) error {
	fake.sendBlockResponseMutex.Lock()
	ret, specificReturn := fake.sendBlockResponseReturnsOnCall[len(fake.sendBlockResponseArgsForCall)]
	fake.sendBlockResponseArgsForCall = append(fake.sendBlockResponseArgsForCall, struct {
		arg1 *common.Block
	}{arg1})
	fake.recordInvocation("SendBlockResponse", []interface{}{arg1})
	fake.sendBlockResponseMutex.Unlock()
	if fake.SendBlockResponseStub != nil {
		return fake.SendBlockResponseStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.sendBlockResponseReturns
	return fakeReturns.result1
}

func (fake *ResponseSender) SendBlockResponseCallCount() int {
	fake.sendBlockResponseMutex.RLock()
	defer fake.sendBlockResponseMutex.RUnlock()
	return len(fake.sendBlockResponseArgsForCall)
}

func (fake *ResponseSender) SendBlockResponseCalls(stub func(*common.Block) error) {
	fake.sendBlockResponseMutex.Lock()
	defer fake.sendBlockResponseMutex.Unlock()
	fake.SendBlockResponseStub = stub
}

func (fake *ResponseSender) SendBlockResponseArgsForCall(i int) *common.Block {
	fake.sendBlockResponseMutex.RLock()
	defer fake.sendBlockResponseMutex.RUnlock()
	argsForCall := fake.sendBlockResponseArgsForCall[i]
	return argsForCall.arg1
}

func (fake *ResponseSender) SendBlockResponseReturns(result1 error) {
	fake.sendBlockResponseMutex.Lock()
	defer fake.sendBlockResponseMutex.Unlock()
	fake.SendBlockResponseStub = nil
	fake.sendBlockResponseReturns = struct {
		result1 error
	}{result1}
}

func (fake *ResponseSender) SendBlockResponseReturnsOnCall(i int, result1 error) {
	fake.sendBlockResponseMutex.Lock()
	defer fake.sendBlockResponseMutex.Unlock()
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

func (fake *ResponseSender) SendStatusResponse(arg1 common.Status) error {
	fake.sendStatusResponseMutex.Lock()
	ret, specificReturn := fake.sendStatusResponseReturnsOnCall[len(fake.sendStatusResponseArgsForCall)]
	fake.sendStatusResponseArgsForCall = append(fake.sendStatusResponseArgsForCall, struct {
		arg1 common.Status
	}{arg1})
	fake.recordInvocation("SendStatusResponse", []interface{}{arg1})
	fake.sendStatusResponseMutex.Unlock()
	if fake.SendStatusResponseStub != nil {
		return fake.SendStatusResponseStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.sendStatusResponseReturns
	return fakeReturns.result1
}

func (fake *ResponseSender) SendStatusResponseCallCount() int {
	fake.sendStatusResponseMutex.RLock()
	defer fake.sendStatusResponseMutex.RUnlock()
	return len(fake.sendStatusResponseArgsForCall)
}

func (fake *ResponseSender) SendStatusResponseCalls(stub func(common.Status) error) {
	fake.sendStatusResponseMutex.Lock()
	defer fake.sendStatusResponseMutex.Unlock()
	fake.SendStatusResponseStub = stub
}

func (fake *ResponseSender) SendStatusResponseArgsForCall(i int) common.Status {
	fake.sendStatusResponseMutex.RLock()
	defer fake.sendStatusResponseMutex.RUnlock()
	argsForCall := fake.sendStatusResponseArgsForCall[i]
	return argsForCall.arg1
}

func (fake *ResponseSender) SendStatusResponseReturns(result1 error) {
	fake.sendStatusResponseMutex.Lock()
	defer fake.sendStatusResponseMutex.Unlock()
	fake.SendStatusResponseStub = nil
	fake.sendStatusResponseReturns = struct {
		result1 error
	}{result1}
}

func (fake *ResponseSender) SendStatusResponseReturnsOnCall(i int, result1 error) {
	fake.sendStatusResponseMutex.Lock()
	defer fake.sendStatusResponseMutex.Unlock()
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

func (fake *ResponseSender) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.sendBlockResponseMutex.RLock()
	defer fake.sendBlockResponseMutex.RUnlock()
	fake.sendStatusResponseMutex.RLock()
	defer fake.sendStatusResponseMutex.RUnlock()
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
