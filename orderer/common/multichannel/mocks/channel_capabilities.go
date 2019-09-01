
package mocks

import (
	"sync"

	"github.com/mcc-github/blockchain/msp"
)

type ChannelCapabilities struct {
	ConsensusTypeMigrationStub        func() bool
	consensusTypeMigrationMutex       sync.RWMutex
	consensusTypeMigrationArgsForCall []struct {
	}
	consensusTypeMigrationReturns struct {
		result1 bool
	}
	consensusTypeMigrationReturnsOnCall map[int]struct {
		result1 bool
	}
	MSPVersionStub        func() msp.MSPVersion
	mSPVersionMutex       sync.RWMutex
	mSPVersionArgsForCall []struct {
	}
	mSPVersionReturns struct {
		result1 msp.MSPVersion
	}
	mSPVersionReturnsOnCall map[int]struct {
		result1 msp.MSPVersion
	}
	OrgSpecificOrdererEndpointsStub        func() bool
	orgSpecificOrdererEndpointsMutex       sync.RWMutex
	orgSpecificOrdererEndpointsArgsForCall []struct {
	}
	orgSpecificOrdererEndpointsReturns struct {
		result1 bool
	}
	orgSpecificOrdererEndpointsReturnsOnCall map[int]struct {
		result1 bool
	}
	SupportedStub        func() error
	supportedMutex       sync.RWMutex
	supportedArgsForCall []struct {
	}
	supportedReturns struct {
		result1 error
	}
	supportedReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *ChannelCapabilities) ConsensusTypeMigration() bool {
	fake.consensusTypeMigrationMutex.Lock()
	ret, specificReturn := fake.consensusTypeMigrationReturnsOnCall[len(fake.consensusTypeMigrationArgsForCall)]
	fake.consensusTypeMigrationArgsForCall = append(fake.consensusTypeMigrationArgsForCall, struct {
	}{})
	fake.recordInvocation("ConsensusTypeMigration", []interface{}{})
	fake.consensusTypeMigrationMutex.Unlock()
	if fake.ConsensusTypeMigrationStub != nil {
		return fake.ConsensusTypeMigrationStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.consensusTypeMigrationReturns
	return fakeReturns.result1
}

func (fake *ChannelCapabilities) ConsensusTypeMigrationCallCount() int {
	fake.consensusTypeMigrationMutex.RLock()
	defer fake.consensusTypeMigrationMutex.RUnlock()
	return len(fake.consensusTypeMigrationArgsForCall)
}

func (fake *ChannelCapabilities) ConsensusTypeMigrationCalls(stub func() bool) {
	fake.consensusTypeMigrationMutex.Lock()
	defer fake.consensusTypeMigrationMutex.Unlock()
	fake.ConsensusTypeMigrationStub = stub
}

func (fake *ChannelCapabilities) ConsensusTypeMigrationReturns(result1 bool) {
	fake.consensusTypeMigrationMutex.Lock()
	defer fake.consensusTypeMigrationMutex.Unlock()
	fake.ConsensusTypeMigrationStub = nil
	fake.consensusTypeMigrationReturns = struct {
		result1 bool
	}{result1}
}

func (fake *ChannelCapabilities) ConsensusTypeMigrationReturnsOnCall(i int, result1 bool) {
	fake.consensusTypeMigrationMutex.Lock()
	defer fake.consensusTypeMigrationMutex.Unlock()
	fake.ConsensusTypeMigrationStub = nil
	if fake.consensusTypeMigrationReturnsOnCall == nil {
		fake.consensusTypeMigrationReturnsOnCall = make(map[int]struct {
			result1 bool
		})
	}
	fake.consensusTypeMigrationReturnsOnCall[i] = struct {
		result1 bool
	}{result1}
}

func (fake *ChannelCapabilities) MSPVersion() msp.MSPVersion {
	fake.mSPVersionMutex.Lock()
	ret, specificReturn := fake.mSPVersionReturnsOnCall[len(fake.mSPVersionArgsForCall)]
	fake.mSPVersionArgsForCall = append(fake.mSPVersionArgsForCall, struct {
	}{})
	fake.recordInvocation("MSPVersion", []interface{}{})
	fake.mSPVersionMutex.Unlock()
	if fake.MSPVersionStub != nil {
		return fake.MSPVersionStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.mSPVersionReturns
	return fakeReturns.result1
}

func (fake *ChannelCapabilities) MSPVersionCallCount() int {
	fake.mSPVersionMutex.RLock()
	defer fake.mSPVersionMutex.RUnlock()
	return len(fake.mSPVersionArgsForCall)
}

func (fake *ChannelCapabilities) MSPVersionCalls(stub func() msp.MSPVersion) {
	fake.mSPVersionMutex.Lock()
	defer fake.mSPVersionMutex.Unlock()
	fake.MSPVersionStub = stub
}

func (fake *ChannelCapabilities) MSPVersionReturns(result1 msp.MSPVersion) {
	fake.mSPVersionMutex.Lock()
	defer fake.mSPVersionMutex.Unlock()
	fake.MSPVersionStub = nil
	fake.mSPVersionReturns = struct {
		result1 msp.MSPVersion
	}{result1}
}

func (fake *ChannelCapabilities) MSPVersionReturnsOnCall(i int, result1 msp.MSPVersion) {
	fake.mSPVersionMutex.Lock()
	defer fake.mSPVersionMutex.Unlock()
	fake.MSPVersionStub = nil
	if fake.mSPVersionReturnsOnCall == nil {
		fake.mSPVersionReturnsOnCall = make(map[int]struct {
			result1 msp.MSPVersion
		})
	}
	fake.mSPVersionReturnsOnCall[i] = struct {
		result1 msp.MSPVersion
	}{result1}
}

func (fake *ChannelCapabilities) OrgSpecificOrdererEndpoints() bool {
	fake.orgSpecificOrdererEndpointsMutex.Lock()
	ret, specificReturn := fake.orgSpecificOrdererEndpointsReturnsOnCall[len(fake.orgSpecificOrdererEndpointsArgsForCall)]
	fake.orgSpecificOrdererEndpointsArgsForCall = append(fake.orgSpecificOrdererEndpointsArgsForCall, struct {
	}{})
	fake.recordInvocation("OrgSpecificOrdererEndpoints", []interface{}{})
	fake.orgSpecificOrdererEndpointsMutex.Unlock()
	if fake.OrgSpecificOrdererEndpointsStub != nil {
		return fake.OrgSpecificOrdererEndpointsStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.orgSpecificOrdererEndpointsReturns
	return fakeReturns.result1
}

func (fake *ChannelCapabilities) OrgSpecificOrdererEndpointsCallCount() int {
	fake.orgSpecificOrdererEndpointsMutex.RLock()
	defer fake.orgSpecificOrdererEndpointsMutex.RUnlock()
	return len(fake.orgSpecificOrdererEndpointsArgsForCall)
}

func (fake *ChannelCapabilities) OrgSpecificOrdererEndpointsCalls(stub func() bool) {
	fake.orgSpecificOrdererEndpointsMutex.Lock()
	defer fake.orgSpecificOrdererEndpointsMutex.Unlock()
	fake.OrgSpecificOrdererEndpointsStub = stub
}

func (fake *ChannelCapabilities) OrgSpecificOrdererEndpointsReturns(result1 bool) {
	fake.orgSpecificOrdererEndpointsMutex.Lock()
	defer fake.orgSpecificOrdererEndpointsMutex.Unlock()
	fake.OrgSpecificOrdererEndpointsStub = nil
	fake.orgSpecificOrdererEndpointsReturns = struct {
		result1 bool
	}{result1}
}

func (fake *ChannelCapabilities) OrgSpecificOrdererEndpointsReturnsOnCall(i int, result1 bool) {
	fake.orgSpecificOrdererEndpointsMutex.Lock()
	defer fake.orgSpecificOrdererEndpointsMutex.Unlock()
	fake.OrgSpecificOrdererEndpointsStub = nil
	if fake.orgSpecificOrdererEndpointsReturnsOnCall == nil {
		fake.orgSpecificOrdererEndpointsReturnsOnCall = make(map[int]struct {
			result1 bool
		})
	}
	fake.orgSpecificOrdererEndpointsReturnsOnCall[i] = struct {
		result1 bool
	}{result1}
}

func (fake *ChannelCapabilities) Supported() error {
	fake.supportedMutex.Lock()
	ret, specificReturn := fake.supportedReturnsOnCall[len(fake.supportedArgsForCall)]
	fake.supportedArgsForCall = append(fake.supportedArgsForCall, struct {
	}{})
	fake.recordInvocation("Supported", []interface{}{})
	fake.supportedMutex.Unlock()
	if fake.SupportedStub != nil {
		return fake.SupportedStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.supportedReturns
	return fakeReturns.result1
}

func (fake *ChannelCapabilities) SupportedCallCount() int {
	fake.supportedMutex.RLock()
	defer fake.supportedMutex.RUnlock()
	return len(fake.supportedArgsForCall)
}

func (fake *ChannelCapabilities) SupportedCalls(stub func() error) {
	fake.supportedMutex.Lock()
	defer fake.supportedMutex.Unlock()
	fake.SupportedStub = stub
}

func (fake *ChannelCapabilities) SupportedReturns(result1 error) {
	fake.supportedMutex.Lock()
	defer fake.supportedMutex.Unlock()
	fake.SupportedStub = nil
	fake.supportedReturns = struct {
		result1 error
	}{result1}
}

func (fake *ChannelCapabilities) SupportedReturnsOnCall(i int, result1 error) {
	fake.supportedMutex.Lock()
	defer fake.supportedMutex.Unlock()
	fake.SupportedStub = nil
	if fake.supportedReturnsOnCall == nil {
		fake.supportedReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.supportedReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *ChannelCapabilities) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.consensusTypeMigrationMutex.RLock()
	defer fake.consensusTypeMigrationMutex.RUnlock()
	fake.mSPVersionMutex.RLock()
	defer fake.mSPVersionMutex.RUnlock()
	fake.orgSpecificOrdererEndpointsMutex.RLock()
	defer fake.orgSpecificOrdererEndpointsMutex.RUnlock()
	fake.supportedMutex.RLock()
	defer fake.supportedMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *ChannelCapabilities) recordInvocation(key string, args []interface{}) {
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
