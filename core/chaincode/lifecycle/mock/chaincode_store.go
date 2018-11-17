
package mock

import (
	sync "sync"
)

type ChaincodeStore struct {
	RetrieveHashStub        func(string, string) ([]byte, error)
	retrieveHashMutex       sync.RWMutex
	retrieveHashArgsForCall []struct {
		arg1 string
		arg2 string
	}
	retrieveHashReturns struct {
		result1 []byte
		result2 error
	}
	retrieveHashReturnsOnCall map[int]struct {
		result1 []byte
		result2 error
	}
	SaveStub        func(string, string, []byte) ([]byte, error)
	saveMutex       sync.RWMutex
	saveArgsForCall []struct {
		arg1 string
		arg2 string
		arg3 []byte
	}
	saveReturns struct {
		result1 []byte
		result2 error
	}
	saveReturnsOnCall map[int]struct {
		result1 []byte
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *ChaincodeStore) RetrieveHash(arg1 string, arg2 string) ([]byte, error) {
	fake.retrieveHashMutex.Lock()
	ret, specificReturn := fake.retrieveHashReturnsOnCall[len(fake.retrieveHashArgsForCall)]
	fake.retrieveHashArgsForCall = append(fake.retrieveHashArgsForCall, struct {
		arg1 string
		arg2 string
	}{arg1, arg2})
	fake.recordInvocation("RetrieveHash", []interface{}{arg1, arg2})
	fake.retrieveHashMutex.Unlock()
	if fake.RetrieveHashStub != nil {
		return fake.RetrieveHashStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.retrieveHashReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *ChaincodeStore) RetrieveHashCallCount() int {
	fake.retrieveHashMutex.RLock()
	defer fake.retrieveHashMutex.RUnlock()
	return len(fake.retrieveHashArgsForCall)
}

func (fake *ChaincodeStore) RetrieveHashCalls(stub func(string, string) ([]byte, error)) {
	fake.retrieveHashMutex.Lock()
	defer fake.retrieveHashMutex.Unlock()
	fake.RetrieveHashStub = stub
}

func (fake *ChaincodeStore) RetrieveHashArgsForCall(i int) (string, string) {
	fake.retrieveHashMutex.RLock()
	defer fake.retrieveHashMutex.RUnlock()
	argsForCall := fake.retrieveHashArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *ChaincodeStore) RetrieveHashReturns(result1 []byte, result2 error) {
	fake.retrieveHashMutex.Lock()
	defer fake.retrieveHashMutex.Unlock()
	fake.RetrieveHashStub = nil
	fake.retrieveHashReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *ChaincodeStore) RetrieveHashReturnsOnCall(i int, result1 []byte, result2 error) {
	fake.retrieveHashMutex.Lock()
	defer fake.retrieveHashMutex.Unlock()
	fake.RetrieveHashStub = nil
	if fake.retrieveHashReturnsOnCall == nil {
		fake.retrieveHashReturnsOnCall = make(map[int]struct {
			result1 []byte
			result2 error
		})
	}
	fake.retrieveHashReturnsOnCall[i] = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *ChaincodeStore) Save(arg1 string, arg2 string, arg3 []byte) ([]byte, error) {
	var arg3Copy []byte
	if arg3 != nil {
		arg3Copy = make([]byte, len(arg3))
		copy(arg3Copy, arg3)
	}
	fake.saveMutex.Lock()
	ret, specificReturn := fake.saveReturnsOnCall[len(fake.saveArgsForCall)]
	fake.saveArgsForCall = append(fake.saveArgsForCall, struct {
		arg1 string
		arg2 string
		arg3 []byte
	}{arg1, arg2, arg3Copy})
	fake.recordInvocation("Save", []interface{}{arg1, arg2, arg3Copy})
	fake.saveMutex.Unlock()
	if fake.SaveStub != nil {
		return fake.SaveStub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.saveReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *ChaincodeStore) SaveCallCount() int {
	fake.saveMutex.RLock()
	defer fake.saveMutex.RUnlock()
	return len(fake.saveArgsForCall)
}

func (fake *ChaincodeStore) SaveCalls(stub func(string, string, []byte) ([]byte, error)) {
	fake.saveMutex.Lock()
	defer fake.saveMutex.Unlock()
	fake.SaveStub = stub
}

func (fake *ChaincodeStore) SaveArgsForCall(i int) (string, string, []byte) {
	fake.saveMutex.RLock()
	defer fake.saveMutex.RUnlock()
	argsForCall := fake.saveArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *ChaincodeStore) SaveReturns(result1 []byte, result2 error) {
	fake.saveMutex.Lock()
	defer fake.saveMutex.Unlock()
	fake.SaveStub = nil
	fake.saveReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *ChaincodeStore) SaveReturnsOnCall(i int, result1 []byte, result2 error) {
	fake.saveMutex.Lock()
	defer fake.saveMutex.Unlock()
	fake.SaveStub = nil
	if fake.saveReturnsOnCall == nil {
		fake.saveReturnsOnCall = make(map[int]struct {
			result1 []byte
			result2 error
		})
	}
	fake.saveReturnsOnCall[i] = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *ChaincodeStore) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.retrieveHashMutex.RLock()
	defer fake.retrieveHashMutex.RUnlock()
	fake.saveMutex.RLock()
	defer fake.saveMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *ChaincodeStore) recordInvocation(key string, args []interface{}) {
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
