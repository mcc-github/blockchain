
package mock

import (
	sync "sync"
	time "time"

	common "github.com/mcc-github/blockchain/protos/common"
	client "github.com/mcc-github/blockchain/token/client"
)

type FabricTxSubmitter struct {
	CreateTxEnvelopeStub        func([]byte) (*common.Envelope, string, error)
	createTxEnvelopeMutex       sync.RWMutex
	createTxEnvelopeArgsForCall []struct {
		arg1 []byte
	}
	createTxEnvelopeReturns struct {
		result1 *common.Envelope
		result2 string
		result3 error
	}
	createTxEnvelopeReturnsOnCall map[int]struct {
		result1 *common.Envelope
		result2 string
		result3 error
	}
	SubmitStub        func(*common.Envelope, time.Duration) (*common.Status, bool, error)
	submitMutex       sync.RWMutex
	submitArgsForCall []struct {
		arg1 *common.Envelope
		arg2 time.Duration
	}
	submitReturns struct {
		result1 *common.Status
		result2 bool
		result3 error
	}
	submitReturnsOnCall map[int]struct {
		result1 *common.Status
		result2 bool
		result3 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FabricTxSubmitter) CreateTxEnvelope(arg1 []byte) (*common.Envelope, string, error) {
	var arg1Copy []byte
	if arg1 != nil {
		arg1Copy = make([]byte, len(arg1))
		copy(arg1Copy, arg1)
	}
	fake.createTxEnvelopeMutex.Lock()
	ret, specificReturn := fake.createTxEnvelopeReturnsOnCall[len(fake.createTxEnvelopeArgsForCall)]
	fake.createTxEnvelopeArgsForCall = append(fake.createTxEnvelopeArgsForCall, struct {
		arg1 []byte
	}{arg1Copy})
	fake.recordInvocation("CreateTxEnvelope", []interface{}{arg1Copy})
	fake.createTxEnvelopeMutex.Unlock()
	if fake.CreateTxEnvelopeStub != nil {
		return fake.CreateTxEnvelopeStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3
	}
	fakeReturns := fake.createTxEnvelopeReturns
	return fakeReturns.result1, fakeReturns.result2, fakeReturns.result3
}

func (fake *FabricTxSubmitter) CreateTxEnvelopeCallCount() int {
	fake.createTxEnvelopeMutex.RLock()
	defer fake.createTxEnvelopeMutex.RUnlock()
	return len(fake.createTxEnvelopeArgsForCall)
}

func (fake *FabricTxSubmitter) CreateTxEnvelopeCalls(stub func([]byte) (*common.Envelope, string, error)) {
	fake.createTxEnvelopeMutex.Lock()
	defer fake.createTxEnvelopeMutex.Unlock()
	fake.CreateTxEnvelopeStub = stub
}

func (fake *FabricTxSubmitter) CreateTxEnvelopeArgsForCall(i int) []byte {
	fake.createTxEnvelopeMutex.RLock()
	defer fake.createTxEnvelopeMutex.RUnlock()
	argsForCall := fake.createTxEnvelopeArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FabricTxSubmitter) CreateTxEnvelopeReturns(result1 *common.Envelope, result2 string, result3 error) {
	fake.createTxEnvelopeMutex.Lock()
	defer fake.createTxEnvelopeMutex.Unlock()
	fake.CreateTxEnvelopeStub = nil
	fake.createTxEnvelopeReturns = struct {
		result1 *common.Envelope
		result2 string
		result3 error
	}{result1, result2, result3}
}

func (fake *FabricTxSubmitter) CreateTxEnvelopeReturnsOnCall(i int, result1 *common.Envelope, result2 string, result3 error) {
	fake.createTxEnvelopeMutex.Lock()
	defer fake.createTxEnvelopeMutex.Unlock()
	fake.CreateTxEnvelopeStub = nil
	if fake.createTxEnvelopeReturnsOnCall == nil {
		fake.createTxEnvelopeReturnsOnCall = make(map[int]struct {
			result1 *common.Envelope
			result2 string
			result3 error
		})
	}
	fake.createTxEnvelopeReturnsOnCall[i] = struct {
		result1 *common.Envelope
		result2 string
		result3 error
	}{result1, result2, result3}
}

func (fake *FabricTxSubmitter) Submit(arg1 *common.Envelope, arg2 time.Duration) (*common.Status, bool, error) {
	fake.submitMutex.Lock()
	ret, specificReturn := fake.submitReturnsOnCall[len(fake.submitArgsForCall)]
	fake.submitArgsForCall = append(fake.submitArgsForCall, struct {
		arg1 *common.Envelope
		arg2 time.Duration
	}{arg1, arg2})
	fake.recordInvocation("Submit", []interface{}{arg1, arg2})
	fake.submitMutex.Unlock()
	if fake.SubmitStub != nil {
		return fake.SubmitStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3
	}
	fakeReturns := fake.submitReturns
	return fakeReturns.result1, fakeReturns.result2, fakeReturns.result3
}

func (fake *FabricTxSubmitter) SubmitCallCount() int {
	fake.submitMutex.RLock()
	defer fake.submitMutex.RUnlock()
	return len(fake.submitArgsForCall)
}

func (fake *FabricTxSubmitter) SubmitCalls(stub func(*common.Envelope, time.Duration) (*common.Status, bool, error)) {
	fake.submitMutex.Lock()
	defer fake.submitMutex.Unlock()
	fake.SubmitStub = stub
}

func (fake *FabricTxSubmitter) SubmitArgsForCall(i int) (*common.Envelope, time.Duration) {
	fake.submitMutex.RLock()
	defer fake.submitMutex.RUnlock()
	argsForCall := fake.submitArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FabricTxSubmitter) SubmitReturns(result1 *common.Status, result2 bool, result3 error) {
	fake.submitMutex.Lock()
	defer fake.submitMutex.Unlock()
	fake.SubmitStub = nil
	fake.submitReturns = struct {
		result1 *common.Status
		result2 bool
		result3 error
	}{result1, result2, result3}
}

func (fake *FabricTxSubmitter) SubmitReturnsOnCall(i int, result1 *common.Status, result2 bool, result3 error) {
	fake.submitMutex.Lock()
	defer fake.submitMutex.Unlock()
	fake.SubmitStub = nil
	if fake.submitReturnsOnCall == nil {
		fake.submitReturnsOnCall = make(map[int]struct {
			result1 *common.Status
			result2 bool
			result3 error
		})
	}
	fake.submitReturnsOnCall[i] = struct {
		result1 *common.Status
		result2 bool
		result3 error
	}{result1, result2, result3}
}

func (fake *FabricTxSubmitter) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.createTxEnvelopeMutex.RLock()
	defer fake.createTxEnvelopeMutex.RUnlock()
	fake.submitMutex.RLock()
	defer fake.submitMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FabricTxSubmitter) recordInvocation(key string, args []interface{}) {
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

var _ client.FabricTxSubmitter = new(FabricTxSubmitter)
