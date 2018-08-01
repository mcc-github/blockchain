
package mock

import (
	"sync"
)

type StorePackageProvider struct {
	LoadStub        func(hash []byte) (codePackage []byte, name, version string, err error)
	loadMutex       sync.RWMutex
	loadArgsForCall []struct {
		hash []byte
	}
	loadReturns struct {
		result1 []byte
		result2 string
		result3 string
		result4 error
	}
	loadReturnsOnCall map[int]struct {
		result1 []byte
		result2 string
		result3 string
		result4 error
	}
	RetrieveHashStub        func(name, version string) (hash []byte, err error)
	retrieveHashMutex       sync.RWMutex
	retrieveHashArgsForCall []struct {
		name    string
		version string
	}
	retrieveHashReturns struct {
		result1 []byte
		result2 error
	}
	retrieveHashReturnsOnCall map[int]struct {
		result1 []byte
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *StorePackageProvider) Load(hash []byte) (codePackage []byte, name, version string, err error) {
	var hashCopy []byte
	if hash != nil {
		hashCopy = make([]byte, len(hash))
		copy(hashCopy, hash)
	}
	fake.loadMutex.Lock()
	ret, specificReturn := fake.loadReturnsOnCall[len(fake.loadArgsForCall)]
	fake.loadArgsForCall = append(fake.loadArgsForCall, struct {
		hash []byte
	}{hashCopy})
	fake.recordInvocation("Load", []interface{}{hashCopy})
	fake.loadMutex.Unlock()
	if fake.LoadStub != nil {
		return fake.LoadStub(hash)
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3, ret.result4
	}
	return fake.loadReturns.result1, fake.loadReturns.result2, fake.loadReturns.result3, fake.loadReturns.result4
}

func (fake *StorePackageProvider) LoadCallCount() int {
	fake.loadMutex.RLock()
	defer fake.loadMutex.RUnlock()
	return len(fake.loadArgsForCall)
}

func (fake *StorePackageProvider) LoadArgsForCall(i int) []byte {
	fake.loadMutex.RLock()
	defer fake.loadMutex.RUnlock()
	return fake.loadArgsForCall[i].hash
}

func (fake *StorePackageProvider) LoadReturns(result1 []byte, result2 string, result3 string, result4 error) {
	fake.LoadStub = nil
	fake.loadReturns = struct {
		result1 []byte
		result2 string
		result3 string
		result4 error
	}{result1, result2, result3, result4}
}

func (fake *StorePackageProvider) LoadReturnsOnCall(i int, result1 []byte, result2 string, result3 string, result4 error) {
	fake.LoadStub = nil
	if fake.loadReturnsOnCall == nil {
		fake.loadReturnsOnCall = make(map[int]struct {
			result1 []byte
			result2 string
			result3 string
			result4 error
		})
	}
	fake.loadReturnsOnCall[i] = struct {
		result1 []byte
		result2 string
		result3 string
		result4 error
	}{result1, result2, result3, result4}
}

func (fake *StorePackageProvider) RetrieveHash(name string, version string) (hash []byte, err error) {
	fake.retrieveHashMutex.Lock()
	ret, specificReturn := fake.retrieveHashReturnsOnCall[len(fake.retrieveHashArgsForCall)]
	fake.retrieveHashArgsForCall = append(fake.retrieveHashArgsForCall, struct {
		name    string
		version string
	}{name, version})
	fake.recordInvocation("RetrieveHash", []interface{}{name, version})
	fake.retrieveHashMutex.Unlock()
	if fake.RetrieveHashStub != nil {
		return fake.RetrieveHashStub(name, version)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.retrieveHashReturns.result1, fake.retrieveHashReturns.result2
}

func (fake *StorePackageProvider) RetrieveHashCallCount() int {
	fake.retrieveHashMutex.RLock()
	defer fake.retrieveHashMutex.RUnlock()
	return len(fake.retrieveHashArgsForCall)
}

func (fake *StorePackageProvider) RetrieveHashArgsForCall(i int) (string, string) {
	fake.retrieveHashMutex.RLock()
	defer fake.retrieveHashMutex.RUnlock()
	return fake.retrieveHashArgsForCall[i].name, fake.retrieveHashArgsForCall[i].version
}

func (fake *StorePackageProvider) RetrieveHashReturns(result1 []byte, result2 error) {
	fake.RetrieveHashStub = nil
	fake.retrieveHashReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *StorePackageProvider) RetrieveHashReturnsOnCall(i int, result1 []byte, result2 error) {
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

func (fake *StorePackageProvider) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.loadMutex.RLock()
	defer fake.loadMutex.RUnlock()
	fake.retrieveHashMutex.RLock()
	defer fake.retrieveHashMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *StorePackageProvider) recordInvocation(key string, args []interface{}) {
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
