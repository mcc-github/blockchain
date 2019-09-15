
package mock

import (
	"sync"

	"github.com/mcc-github/blockchain-protos-go/common"
	"github.com/mcc-github/blockchain/common/deliver"
	"github.com/mcc-github/blockchain/protoutil"
)

type FilteredResponseSender struct {
	DataTypeStub        func() string
	dataTypeMutex       sync.RWMutex
	dataTypeArgsForCall []struct {
	}
	dataTypeReturns struct {
		result1 string
	}
	dataTypeReturnsOnCall map[int]struct {
		result1 string
	}
	IsFilteredStub        func() bool
	isFilteredMutex       sync.RWMutex
	isFilteredArgsForCall []struct {
	}
	isFilteredReturns struct {
		result1 bool
	}
	isFilteredReturnsOnCall map[int]struct {
		result1 bool
	}
	SendBlockResponseStub        func(*common.Block, string, deliver.Chain, *protoutil.SignedData) error
	sendBlockResponseMutex       sync.RWMutex
	sendBlockResponseArgsForCall []struct {
		arg1 *common.Block
		arg2 string
		arg3 deliver.Chain
		arg4 *protoutil.SignedData
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

func (fake *FilteredResponseSender) DataType() string {
	fake.dataTypeMutex.Lock()
	ret, specificReturn := fake.dataTypeReturnsOnCall[len(fake.dataTypeArgsForCall)]
	fake.dataTypeArgsForCall = append(fake.dataTypeArgsForCall, struct {
	}{})
	fake.recordInvocation("DataType", []interface{}{})
	fake.dataTypeMutex.Unlock()
	if fake.DataTypeStub != nil {
		return fake.DataTypeStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.dataTypeReturns
	return fakeReturns.result1
}

func (fake *FilteredResponseSender) DataTypeCallCount() int {
	fake.dataTypeMutex.RLock()
	defer fake.dataTypeMutex.RUnlock()
	return len(fake.dataTypeArgsForCall)
}

func (fake *FilteredResponseSender) DataTypeCalls(stub func() string) {
	fake.dataTypeMutex.Lock()
	defer fake.dataTypeMutex.Unlock()
	fake.DataTypeStub = stub
}

func (fake *FilteredResponseSender) DataTypeReturns(result1 string) {
	fake.dataTypeMutex.Lock()
	defer fake.dataTypeMutex.Unlock()
	fake.DataTypeStub = nil
	fake.dataTypeReturns = struct {
		result1 string
	}{result1}
}

func (fake *FilteredResponseSender) DataTypeReturnsOnCall(i int, result1 string) {
	fake.dataTypeMutex.Lock()
	defer fake.dataTypeMutex.Unlock()
	fake.DataTypeStub = nil
	if fake.dataTypeReturnsOnCall == nil {
		fake.dataTypeReturnsOnCall = make(map[int]struct {
			result1 string
		})
	}
	fake.dataTypeReturnsOnCall[i] = struct {
		result1 string
	}{result1}
}

func (fake *FilteredResponseSender) IsFiltered() bool {
	fake.isFilteredMutex.Lock()
	ret, specificReturn := fake.isFilteredReturnsOnCall[len(fake.isFilteredArgsForCall)]
	fake.isFilteredArgsForCall = append(fake.isFilteredArgsForCall, struct {
	}{})
	fake.recordInvocation("IsFiltered", []interface{}{})
	fake.isFilteredMutex.Unlock()
	if fake.IsFilteredStub != nil {
		return fake.IsFilteredStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.isFilteredReturns
	return fakeReturns.result1
}

func (fake *FilteredResponseSender) IsFilteredCallCount() int {
	fake.isFilteredMutex.RLock()
	defer fake.isFilteredMutex.RUnlock()
	return len(fake.isFilteredArgsForCall)
}

func (fake *FilteredResponseSender) IsFilteredCalls(stub func() bool) {
	fake.isFilteredMutex.Lock()
	defer fake.isFilteredMutex.Unlock()
	fake.IsFilteredStub = stub
}

func (fake *FilteredResponseSender) IsFilteredReturns(result1 bool) {
	fake.isFilteredMutex.Lock()
	defer fake.isFilteredMutex.Unlock()
	fake.IsFilteredStub = nil
	fake.isFilteredReturns = struct {
		result1 bool
	}{result1}
}

func (fake *FilteredResponseSender) IsFilteredReturnsOnCall(i int, result1 bool) {
	fake.isFilteredMutex.Lock()
	defer fake.isFilteredMutex.Unlock()
	fake.IsFilteredStub = nil
	if fake.isFilteredReturnsOnCall == nil {
		fake.isFilteredReturnsOnCall = make(map[int]struct {
			result1 bool
		})
	}
	fake.isFilteredReturnsOnCall[i] = struct {
		result1 bool
	}{result1}
}

func (fake *FilteredResponseSender) SendBlockResponse(arg1 *common.Block, arg2 string, arg3 deliver.Chain, arg4 *protoutil.SignedData) error {
	fake.sendBlockResponseMutex.Lock()
	ret, specificReturn := fake.sendBlockResponseReturnsOnCall[len(fake.sendBlockResponseArgsForCall)]
	fake.sendBlockResponseArgsForCall = append(fake.sendBlockResponseArgsForCall, struct {
		arg1 *common.Block
		arg2 string
		arg3 deliver.Chain
		arg4 *protoutil.SignedData
	}{arg1, arg2, arg3, arg4})
	fake.recordInvocation("SendBlockResponse", []interface{}{arg1, arg2, arg3, arg4})
	fake.sendBlockResponseMutex.Unlock()
	if fake.SendBlockResponseStub != nil {
		return fake.SendBlockResponseStub(arg1, arg2, arg3, arg4)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.sendBlockResponseReturns
	return fakeReturns.result1
}

func (fake *FilteredResponseSender) SendBlockResponseCallCount() int {
	fake.sendBlockResponseMutex.RLock()
	defer fake.sendBlockResponseMutex.RUnlock()
	return len(fake.sendBlockResponseArgsForCall)
}

func (fake *FilteredResponseSender) SendBlockResponseCalls(stub func(*common.Block, string, deliver.Chain, *protoutil.SignedData) error) {
	fake.sendBlockResponseMutex.Lock()
	defer fake.sendBlockResponseMutex.Unlock()
	fake.SendBlockResponseStub = stub
}

func (fake *FilteredResponseSender) SendBlockResponseArgsForCall(i int) (*common.Block, string, deliver.Chain, *protoutil.SignedData) {
	fake.sendBlockResponseMutex.RLock()
	defer fake.sendBlockResponseMutex.RUnlock()
	argsForCall := fake.sendBlockResponseArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4
}

func (fake *FilteredResponseSender) SendBlockResponseReturns(result1 error) {
	fake.sendBlockResponseMutex.Lock()
	defer fake.sendBlockResponseMutex.Unlock()
	fake.SendBlockResponseStub = nil
	fake.sendBlockResponseReturns = struct {
		result1 error
	}{result1}
}

func (fake *FilteredResponseSender) SendBlockResponseReturnsOnCall(i int, result1 error) {
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

func (fake *FilteredResponseSender) SendStatusResponse(arg1 common.Status) error {
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

func (fake *FilteredResponseSender) SendStatusResponseCallCount() int {
	fake.sendStatusResponseMutex.RLock()
	defer fake.sendStatusResponseMutex.RUnlock()
	return len(fake.sendStatusResponseArgsForCall)
}

func (fake *FilteredResponseSender) SendStatusResponseCalls(stub func(common.Status) error) {
	fake.sendStatusResponseMutex.Lock()
	defer fake.sendStatusResponseMutex.Unlock()
	fake.SendStatusResponseStub = stub
}

func (fake *FilteredResponseSender) SendStatusResponseArgsForCall(i int) common.Status {
	fake.sendStatusResponseMutex.RLock()
	defer fake.sendStatusResponseMutex.RUnlock()
	argsForCall := fake.sendStatusResponseArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FilteredResponseSender) SendStatusResponseReturns(result1 error) {
	fake.sendStatusResponseMutex.Lock()
	defer fake.sendStatusResponseMutex.Unlock()
	fake.SendStatusResponseStub = nil
	fake.sendStatusResponseReturns = struct {
		result1 error
	}{result1}
}

func (fake *FilteredResponseSender) SendStatusResponseReturnsOnCall(i int, result1 error) {
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

func (fake *FilteredResponseSender) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.dataTypeMutex.RLock()
	defer fake.dataTypeMutex.RUnlock()
	fake.isFilteredMutex.RLock()
	defer fake.isFilteredMutex.RUnlock()
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

func (fake *FilteredResponseSender) recordInvocation(key string, args []interface{}) {
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
