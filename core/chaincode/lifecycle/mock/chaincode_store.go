
package mock

import (
	"sync"
)

type ChaincodeStore struct {
	SaveStub        func(name, version string, ccInstallPkg []byte) (hash []byte, err error)
	saveMutex       sync.RWMutex
	saveArgsForCall []struct {
		name         string
		version      string
		ccInstallPkg []byte
	}
	saveReturns struct {
		result1 []byte
		result2 error
	}
	saveReturnsOnCall map[int]struct {
		result1 []byte
		result2 error
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

func (fake *ChaincodeStore) Save(name string, version string, ccInstallPkg []byte) (hash []byte, err error) {
	var ccInstallPkgCopy []byte
	if ccInstallPkg != nil {
		ccInstallPkgCopy = make([]byte, len(ccInstallPkg))
		copy(ccInstallPkgCopy, ccInstallPkg)
	}
	fake.saveMutex.Lock()
	ret, specificReturn := fake.saveReturnsOnCall[len(fake.saveArgsForCall)]
	fake.saveArgsForCall = append(fake.saveArgsForCall, struct {
		name         string
		version      string
		ccInstallPkg []byte
	}{name, version, ccInstallPkgCopy})
	fake.recordInvocation("Save", []interface{}{name, version, ccInstallPkgCopy})
	fake.saveMutex.Unlock()
	if fake.SaveStub != nil {
		return fake.SaveStub(name, version, ccInstallPkg)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.saveReturns.result1, fake.saveReturns.result2
}

func (fake *ChaincodeStore) SaveCallCount() int {
	fake.saveMutex.RLock()
	defer fake.saveMutex.RUnlock()
	return len(fake.saveArgsForCall)
}

func (fake *ChaincodeStore) SaveArgsForCall(i int) (string, string, []byte) {
	fake.saveMutex.RLock()
	defer fake.saveMutex.RUnlock()
	return fake.saveArgsForCall[i].name, fake.saveArgsForCall[i].version, fake.saveArgsForCall[i].ccInstallPkg
}

func (fake *ChaincodeStore) SaveReturns(result1 []byte, result2 error) {
	fake.SaveStub = nil
	fake.saveReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *ChaincodeStore) SaveReturnsOnCall(i int, result1 []byte, result2 error) {
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

func (fake *ChaincodeStore) RetrieveHash(name string, version string) (hash []byte, err error) {
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

func (fake *ChaincodeStore) RetrieveHashCallCount() int {
	fake.retrieveHashMutex.RLock()
	defer fake.retrieveHashMutex.RUnlock()
	return len(fake.retrieveHashArgsForCall)
}

func (fake *ChaincodeStore) RetrieveHashArgsForCall(i int) (string, string) {
	fake.retrieveHashMutex.RLock()
	defer fake.retrieveHashMutex.RUnlock()
	return fake.retrieveHashArgsForCall[i].name, fake.retrieveHashArgsForCall[i].version
}

func (fake *ChaincodeStore) RetrieveHashReturns(result1 []byte, result2 error) {
	fake.RetrieveHashStub = nil
	fake.retrieveHashReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *ChaincodeStore) RetrieveHashReturnsOnCall(i int, result1 []byte, result2 error) {
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

func (fake *ChaincodeStore) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.saveMutex.RLock()
	defer fake.saveMutex.RUnlock()
	fake.retrieveHashMutex.RLock()
	defer fake.retrieveHashMutex.RUnlock()
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
