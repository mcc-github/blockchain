
package mock

import (
	"sync"

	"github.com/mcc-github/blockchain/core/ledger"
)

type ChaincodeLifecycleEventListener struct {
	ChaincodeDeployDoneStub        func(bool)
	chaincodeDeployDoneMutex       sync.RWMutex
	chaincodeDeployDoneArgsForCall []struct {
		arg1 bool
	}
	HandleChaincodeDeployStub        func(*ledger.ChaincodeDefinition, []byte) error
	handleChaincodeDeployMutex       sync.RWMutex
	handleChaincodeDeployArgsForCall []struct {
		arg1 *ledger.ChaincodeDefinition
		arg2 []byte
	}
	handleChaincodeDeployReturns struct {
		result1 error
	}
	handleChaincodeDeployReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *ChaincodeLifecycleEventListener) ChaincodeDeployDone(arg1 bool) {
	fake.chaincodeDeployDoneMutex.Lock()
	fake.chaincodeDeployDoneArgsForCall = append(fake.chaincodeDeployDoneArgsForCall, struct {
		arg1 bool
	}{arg1})
	fake.recordInvocation("ChaincodeDeployDone", []interface{}{arg1})
	fake.chaincodeDeployDoneMutex.Unlock()
	if fake.ChaincodeDeployDoneStub != nil {
		fake.ChaincodeDeployDoneStub(arg1)
	}
}

func (fake *ChaincodeLifecycleEventListener) ChaincodeDeployDoneCallCount() int {
	fake.chaincodeDeployDoneMutex.RLock()
	defer fake.chaincodeDeployDoneMutex.RUnlock()
	return len(fake.chaincodeDeployDoneArgsForCall)
}

func (fake *ChaincodeLifecycleEventListener) ChaincodeDeployDoneCalls(stub func(bool)) {
	fake.chaincodeDeployDoneMutex.Lock()
	defer fake.chaincodeDeployDoneMutex.Unlock()
	fake.ChaincodeDeployDoneStub = stub
}

func (fake *ChaincodeLifecycleEventListener) ChaincodeDeployDoneArgsForCall(i int) bool {
	fake.chaincodeDeployDoneMutex.RLock()
	defer fake.chaincodeDeployDoneMutex.RUnlock()
	argsForCall := fake.chaincodeDeployDoneArgsForCall[i]
	return argsForCall.arg1
}

func (fake *ChaincodeLifecycleEventListener) HandleChaincodeDeploy(arg1 *ledger.ChaincodeDefinition, arg2 []byte) error {
	var arg2Copy []byte
	if arg2 != nil {
		arg2Copy = make([]byte, len(arg2))
		copy(arg2Copy, arg2)
	}
	fake.handleChaincodeDeployMutex.Lock()
	ret, specificReturn := fake.handleChaincodeDeployReturnsOnCall[len(fake.handleChaincodeDeployArgsForCall)]
	fake.handleChaincodeDeployArgsForCall = append(fake.handleChaincodeDeployArgsForCall, struct {
		arg1 *ledger.ChaincodeDefinition
		arg2 []byte
	}{arg1, arg2Copy})
	fake.recordInvocation("HandleChaincodeDeploy", []interface{}{arg1, arg2Copy})
	fake.handleChaincodeDeployMutex.Unlock()
	if fake.HandleChaincodeDeployStub != nil {
		return fake.HandleChaincodeDeployStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.handleChaincodeDeployReturns
	return fakeReturns.result1
}

func (fake *ChaincodeLifecycleEventListener) HandleChaincodeDeployCallCount() int {
	fake.handleChaincodeDeployMutex.RLock()
	defer fake.handleChaincodeDeployMutex.RUnlock()
	return len(fake.handleChaincodeDeployArgsForCall)
}

func (fake *ChaincodeLifecycleEventListener) HandleChaincodeDeployCalls(stub func(*ledger.ChaincodeDefinition, []byte) error) {
	fake.handleChaincodeDeployMutex.Lock()
	defer fake.handleChaincodeDeployMutex.Unlock()
	fake.HandleChaincodeDeployStub = stub
}

func (fake *ChaincodeLifecycleEventListener) HandleChaincodeDeployArgsForCall(i int) (*ledger.ChaincodeDefinition, []byte) {
	fake.handleChaincodeDeployMutex.RLock()
	defer fake.handleChaincodeDeployMutex.RUnlock()
	argsForCall := fake.handleChaincodeDeployArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *ChaincodeLifecycleEventListener) HandleChaincodeDeployReturns(result1 error) {
	fake.handleChaincodeDeployMutex.Lock()
	defer fake.handleChaincodeDeployMutex.Unlock()
	fake.HandleChaincodeDeployStub = nil
	fake.handleChaincodeDeployReturns = struct {
		result1 error
	}{result1}
}

func (fake *ChaincodeLifecycleEventListener) HandleChaincodeDeployReturnsOnCall(i int, result1 error) {
	fake.handleChaincodeDeployMutex.Lock()
	defer fake.handleChaincodeDeployMutex.Unlock()
	fake.HandleChaincodeDeployStub = nil
	if fake.handleChaincodeDeployReturnsOnCall == nil {
		fake.handleChaincodeDeployReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.handleChaincodeDeployReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *ChaincodeLifecycleEventListener) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.chaincodeDeployDoneMutex.RLock()
	defer fake.chaincodeDeployDoneMutex.RUnlock()
	fake.handleChaincodeDeployMutex.RLock()
	defer fake.handleChaincodeDeployMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *ChaincodeLifecycleEventListener) recordInvocation(key string, args []interface{}) {
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

var _ ledger.ChaincodeLifecycleEventListener = new(ChaincodeLifecycleEventListener)
