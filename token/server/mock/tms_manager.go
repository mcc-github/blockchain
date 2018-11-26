
package mock

import (
	sync "sync"

	server "github.com/mcc-github/blockchain/token/server"
)

type TMSManager struct {
	GetIssuerStub        func(string, []byte, []byte) (server.Issuer, error)
	getIssuerMutex       sync.RWMutex
	getIssuerArgsForCall []struct {
		arg1 string
		arg2 []byte
		arg3 []byte
	}
	getIssuerReturns struct {
		result1 server.Issuer
		result2 error
	}
	getIssuerReturnsOnCall map[int]struct {
		result1 server.Issuer
		result2 error
	}
	GetTransactorStub        func(string, []byte, []byte) (server.Transactor, error)
	getTransactorMutex       sync.RWMutex
	getTransactorArgsForCall []struct {
		arg1 string
		arg2 []byte
		arg3 []byte
	}
	getTransactorReturns struct {
		result1 server.Transactor
		result2 error
	}
	getTransactorReturnsOnCall map[int]struct {
		result1 server.Transactor
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *TMSManager) GetIssuer(arg1 string, arg2 []byte, arg3 []byte) (server.Issuer, error) {
	var arg2Copy []byte
	if arg2 != nil {
		arg2Copy = make([]byte, len(arg2))
		copy(arg2Copy, arg2)
	}
	var arg3Copy []byte
	if arg3 != nil {
		arg3Copy = make([]byte, len(arg3))
		copy(arg3Copy, arg3)
	}
	fake.getIssuerMutex.Lock()
	ret, specificReturn := fake.getIssuerReturnsOnCall[len(fake.getIssuerArgsForCall)]
	fake.getIssuerArgsForCall = append(fake.getIssuerArgsForCall, struct {
		arg1 string
		arg2 []byte
		arg3 []byte
	}{arg1, arg2Copy, arg3Copy})
	fake.recordInvocation("GetIssuer", []interface{}{arg1, arg2Copy, arg3Copy})
	fake.getIssuerMutex.Unlock()
	if fake.GetIssuerStub != nil {
		return fake.GetIssuerStub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getIssuerReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *TMSManager) GetIssuerCallCount() int {
	fake.getIssuerMutex.RLock()
	defer fake.getIssuerMutex.RUnlock()
	return len(fake.getIssuerArgsForCall)
}

func (fake *TMSManager) GetIssuerCalls(stub func(string, []byte, []byte) (server.Issuer, error)) {
	fake.getIssuerMutex.Lock()
	defer fake.getIssuerMutex.Unlock()
	fake.GetIssuerStub = stub
}

func (fake *TMSManager) GetIssuerArgsForCall(i int) (string, []byte, []byte) {
	fake.getIssuerMutex.RLock()
	defer fake.getIssuerMutex.RUnlock()
	argsForCall := fake.getIssuerArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *TMSManager) GetIssuerReturns(result1 server.Issuer, result2 error) {
	fake.getIssuerMutex.Lock()
	defer fake.getIssuerMutex.Unlock()
	fake.GetIssuerStub = nil
	fake.getIssuerReturns = struct {
		result1 server.Issuer
		result2 error
	}{result1, result2}
}

func (fake *TMSManager) GetIssuerReturnsOnCall(i int, result1 server.Issuer, result2 error) {
	fake.getIssuerMutex.Lock()
	defer fake.getIssuerMutex.Unlock()
	fake.GetIssuerStub = nil
	if fake.getIssuerReturnsOnCall == nil {
		fake.getIssuerReturnsOnCall = make(map[int]struct {
			result1 server.Issuer
			result2 error
		})
	}
	fake.getIssuerReturnsOnCall[i] = struct {
		result1 server.Issuer
		result2 error
	}{result1, result2}
}

func (fake *TMSManager) GetTransactor(arg1 string, arg2 []byte, arg3 []byte) (server.Transactor, error) {
	var arg2Copy []byte
	if arg2 != nil {
		arg2Copy = make([]byte, len(arg2))
		copy(arg2Copy, arg2)
	}
	var arg3Copy []byte
	if arg3 != nil {
		arg3Copy = make([]byte, len(arg3))
		copy(arg3Copy, arg3)
	}
	fake.getTransactorMutex.Lock()
	ret, specificReturn := fake.getTransactorReturnsOnCall[len(fake.getTransactorArgsForCall)]
	fake.getTransactorArgsForCall = append(fake.getTransactorArgsForCall, struct {
		arg1 string
		arg2 []byte
		arg3 []byte
	}{arg1, arg2Copy, arg3Copy})
	fake.recordInvocation("GetTransactor", []interface{}{arg1, arg2Copy, arg3Copy})
	fake.getTransactorMutex.Unlock()
	if fake.GetTransactorStub != nil {
		return fake.GetTransactorStub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getTransactorReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *TMSManager) GetTransactorCallCount() int {
	fake.getTransactorMutex.RLock()
	defer fake.getTransactorMutex.RUnlock()
	return len(fake.getTransactorArgsForCall)
}

func (fake *TMSManager) GetTransactorCalls(stub func(string, []byte, []byte) (server.Transactor, error)) {
	fake.getTransactorMutex.Lock()
	defer fake.getTransactorMutex.Unlock()
	fake.GetTransactorStub = stub
}

func (fake *TMSManager) GetTransactorArgsForCall(i int) (string, []byte, []byte) {
	fake.getTransactorMutex.RLock()
	defer fake.getTransactorMutex.RUnlock()
	argsForCall := fake.getTransactorArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *TMSManager) GetTransactorReturns(result1 server.Transactor, result2 error) {
	fake.getTransactorMutex.Lock()
	defer fake.getTransactorMutex.Unlock()
	fake.GetTransactorStub = nil
	fake.getTransactorReturns = struct {
		result1 server.Transactor
		result2 error
	}{result1, result2}
}

func (fake *TMSManager) GetTransactorReturnsOnCall(i int, result1 server.Transactor, result2 error) {
	fake.getTransactorMutex.Lock()
	defer fake.getTransactorMutex.Unlock()
	fake.GetTransactorStub = nil
	if fake.getTransactorReturnsOnCall == nil {
		fake.getTransactorReturnsOnCall = make(map[int]struct {
			result1 server.Transactor
			result2 error
		})
	}
	fake.getTransactorReturnsOnCall[i] = struct {
		result1 server.Transactor
		result2 error
	}{result1, result2}
}

func (fake *TMSManager) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getIssuerMutex.RLock()
	defer fake.getIssuerMutex.RUnlock()
	fake.getTransactorMutex.RLock()
	defer fake.getTransactorMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *TMSManager) recordInvocation(key string, args []interface{}) {
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

var _ server.TMSManager = new(TMSManager)
