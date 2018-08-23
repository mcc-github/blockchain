
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

func (fake *ChaincodeStore) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
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
