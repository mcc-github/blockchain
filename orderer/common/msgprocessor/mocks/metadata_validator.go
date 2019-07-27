
package mocks

import (
	"sync"
)

type FakeMetadataValidator struct {
	ValidateConsensusMetadataStub        func(oldMetadata, newMetadata []byte, newChannel bool) error
	validateConsensusMetadataMutex       sync.RWMutex
	validateConsensusMetadataArgsForCall []struct {
		oldMetadata []byte
		newMetadata []byte
		newChannel  bool
	}
	validateConsensusMetadataReturns struct {
		result1 error
	}
	validateConsensusMetadataReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeMetadataValidator) ValidateConsensusMetadata(oldMetadata []byte, newMetadata []byte, newChannel bool) error {
	var oldMetadataCopy []byte
	if oldMetadata != nil {
		oldMetadataCopy = make([]byte, len(oldMetadata))
		copy(oldMetadataCopy, oldMetadata)
	}
	var newMetadataCopy []byte
	if newMetadata != nil {
		newMetadataCopy = make([]byte, len(newMetadata))
		copy(newMetadataCopy, newMetadata)
	}
	fake.validateConsensusMetadataMutex.Lock()
	ret, specificReturn := fake.validateConsensusMetadataReturnsOnCall[len(fake.validateConsensusMetadataArgsForCall)]
	fake.validateConsensusMetadataArgsForCall = append(fake.validateConsensusMetadataArgsForCall, struct {
		oldMetadata []byte
		newMetadata []byte
		newChannel  bool
	}{oldMetadataCopy, newMetadataCopy, newChannel})
	fake.recordInvocation("ValidateConsensusMetadata", []interface{}{oldMetadataCopy, newMetadataCopy, newChannel})
	fake.validateConsensusMetadataMutex.Unlock()
	if fake.ValidateConsensusMetadataStub != nil {
		return fake.ValidateConsensusMetadataStub(oldMetadata, newMetadata, newChannel)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.validateConsensusMetadataReturns.result1
}

func (fake *FakeMetadataValidator) ValidateConsensusMetadataCallCount() int {
	fake.validateConsensusMetadataMutex.RLock()
	defer fake.validateConsensusMetadataMutex.RUnlock()
	return len(fake.validateConsensusMetadataArgsForCall)
}

func (fake *FakeMetadataValidator) ValidateConsensusMetadataArgsForCall(i int) ([]byte, []byte, bool) {
	fake.validateConsensusMetadataMutex.RLock()
	defer fake.validateConsensusMetadataMutex.RUnlock()
	return fake.validateConsensusMetadataArgsForCall[i].oldMetadata, fake.validateConsensusMetadataArgsForCall[i].newMetadata, fake.validateConsensusMetadataArgsForCall[i].newChannel
}

func (fake *FakeMetadataValidator) ValidateConsensusMetadataReturns(result1 error) {
	fake.ValidateConsensusMetadataStub = nil
	fake.validateConsensusMetadataReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeMetadataValidator) ValidateConsensusMetadataReturnsOnCall(i int, result1 error) {
	fake.ValidateConsensusMetadataStub = nil
	if fake.validateConsensusMetadataReturnsOnCall == nil {
		fake.validateConsensusMetadataReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.validateConsensusMetadataReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeMetadataValidator) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.validateConsensusMetadataMutex.RLock()
	defer fake.validateConsensusMetadataMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeMetadataValidator) recordInvocation(key string, args []interface{}) {
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
