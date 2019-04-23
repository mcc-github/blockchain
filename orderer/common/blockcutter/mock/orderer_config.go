
package mock

import (
	"sync"
	"time"

	"github.com/mcc-github/blockchain/common/channelconfig"
	"github.com/mcc-github/blockchain/protos/orderer"
)

type OrdererConfig struct {
	BatchSizeStub        func() *orderer.BatchSize
	batchSizeMutex       sync.RWMutex
	batchSizeArgsForCall []struct {
	}
	batchSizeReturns struct {
		result1 *orderer.BatchSize
	}
	batchSizeReturnsOnCall map[int]struct {
		result1 *orderer.BatchSize
	}
	BatchTimeoutStub        func() time.Duration
	batchTimeoutMutex       sync.RWMutex
	batchTimeoutArgsForCall []struct {
	}
	batchTimeoutReturns struct {
		result1 time.Duration
	}
	batchTimeoutReturnsOnCall map[int]struct {
		result1 time.Duration
	}
	CapabilitiesStub        func() channelconfig.OrdererCapabilities
	capabilitiesMutex       sync.RWMutex
	capabilitiesArgsForCall []struct {
	}
	capabilitiesReturns struct {
		result1 channelconfig.OrdererCapabilities
	}
	capabilitiesReturnsOnCall map[int]struct {
		result1 channelconfig.OrdererCapabilities
	}
	ConsensusMetadataStub        func() []byte
	consensusMetadataMutex       sync.RWMutex
	consensusMetadataArgsForCall []struct {
	}
	consensusMetadataReturns struct {
		result1 []byte
	}
	consensusMetadataReturnsOnCall map[int]struct {
		result1 []byte
	}
	ConsensusMigrationContextStub        func() uint64
	consensusMigrationContextMutex       sync.RWMutex
	consensusMigrationContextArgsForCall []struct {
	}
	consensusMigrationContextReturns struct {
		result1 uint64
	}
	consensusMigrationContextReturnsOnCall map[int]struct {
		result1 uint64
	}
	ConsensusMigrationStateStub        func() orderer.ConsensusType_MigrationState
	consensusMigrationStateMutex       sync.RWMutex
	consensusMigrationStateArgsForCall []struct {
	}
	consensusMigrationStateReturns struct {
		result1 orderer.ConsensusType_MigrationState
	}
	consensusMigrationStateReturnsOnCall map[int]struct {
		result1 orderer.ConsensusType_MigrationState
	}
	ConsensusTypeStub        func() string
	consensusTypeMutex       sync.RWMutex
	consensusTypeArgsForCall []struct {
	}
	consensusTypeReturns struct {
		result1 string
	}
	consensusTypeReturnsOnCall map[int]struct {
		result1 string
	}
	KafkaBrokersStub        func() []string
	kafkaBrokersMutex       sync.RWMutex
	kafkaBrokersArgsForCall []struct {
	}
	kafkaBrokersReturns struct {
		result1 []string
	}
	kafkaBrokersReturnsOnCall map[int]struct {
		result1 []string
	}
	MaxChannelsCountStub        func() uint64
	maxChannelsCountMutex       sync.RWMutex
	maxChannelsCountArgsForCall []struct {
	}
	maxChannelsCountReturns struct {
		result1 uint64
	}
	maxChannelsCountReturnsOnCall map[int]struct {
		result1 uint64
	}
	OrganizationsStub        func() map[string]channelconfig.OrdererOrg
	organizationsMutex       sync.RWMutex
	organizationsArgsForCall []struct {
	}
	organizationsReturns struct {
		result1 map[string]channelconfig.OrdererOrg
	}
	organizationsReturnsOnCall map[int]struct {
		result1 map[string]channelconfig.OrdererOrg
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *OrdererConfig) BatchSize() *orderer.BatchSize {
	fake.batchSizeMutex.Lock()
	ret, specificReturn := fake.batchSizeReturnsOnCall[len(fake.batchSizeArgsForCall)]
	fake.batchSizeArgsForCall = append(fake.batchSizeArgsForCall, struct {
	}{})
	fake.recordInvocation("BatchSize", []interface{}{})
	fake.batchSizeMutex.Unlock()
	if fake.BatchSizeStub != nil {
		return fake.BatchSizeStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.batchSizeReturns
	return fakeReturns.result1
}

func (fake *OrdererConfig) BatchSizeCallCount() int {
	fake.batchSizeMutex.RLock()
	defer fake.batchSizeMutex.RUnlock()
	return len(fake.batchSizeArgsForCall)
}

func (fake *OrdererConfig) BatchSizeCalls(stub func() *orderer.BatchSize) {
	fake.batchSizeMutex.Lock()
	defer fake.batchSizeMutex.Unlock()
	fake.BatchSizeStub = stub
}

func (fake *OrdererConfig) BatchSizeReturns(result1 *orderer.BatchSize) {
	fake.batchSizeMutex.Lock()
	defer fake.batchSizeMutex.Unlock()
	fake.BatchSizeStub = nil
	fake.batchSizeReturns = struct {
		result1 *orderer.BatchSize
	}{result1}
}

func (fake *OrdererConfig) BatchSizeReturnsOnCall(i int, result1 *orderer.BatchSize) {
	fake.batchSizeMutex.Lock()
	defer fake.batchSizeMutex.Unlock()
	fake.BatchSizeStub = nil
	if fake.batchSizeReturnsOnCall == nil {
		fake.batchSizeReturnsOnCall = make(map[int]struct {
			result1 *orderer.BatchSize
		})
	}
	fake.batchSizeReturnsOnCall[i] = struct {
		result1 *orderer.BatchSize
	}{result1}
}

func (fake *OrdererConfig) BatchTimeout() time.Duration {
	fake.batchTimeoutMutex.Lock()
	ret, specificReturn := fake.batchTimeoutReturnsOnCall[len(fake.batchTimeoutArgsForCall)]
	fake.batchTimeoutArgsForCall = append(fake.batchTimeoutArgsForCall, struct {
	}{})
	fake.recordInvocation("BatchTimeout", []interface{}{})
	fake.batchTimeoutMutex.Unlock()
	if fake.BatchTimeoutStub != nil {
		return fake.BatchTimeoutStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.batchTimeoutReturns
	return fakeReturns.result1
}

func (fake *OrdererConfig) BatchTimeoutCallCount() int {
	fake.batchTimeoutMutex.RLock()
	defer fake.batchTimeoutMutex.RUnlock()
	return len(fake.batchTimeoutArgsForCall)
}

func (fake *OrdererConfig) BatchTimeoutCalls(stub func() time.Duration) {
	fake.batchTimeoutMutex.Lock()
	defer fake.batchTimeoutMutex.Unlock()
	fake.BatchTimeoutStub = stub
}

func (fake *OrdererConfig) BatchTimeoutReturns(result1 time.Duration) {
	fake.batchTimeoutMutex.Lock()
	defer fake.batchTimeoutMutex.Unlock()
	fake.BatchTimeoutStub = nil
	fake.batchTimeoutReturns = struct {
		result1 time.Duration
	}{result1}
}

func (fake *OrdererConfig) BatchTimeoutReturnsOnCall(i int, result1 time.Duration) {
	fake.batchTimeoutMutex.Lock()
	defer fake.batchTimeoutMutex.Unlock()
	fake.BatchTimeoutStub = nil
	if fake.batchTimeoutReturnsOnCall == nil {
		fake.batchTimeoutReturnsOnCall = make(map[int]struct {
			result1 time.Duration
		})
	}
	fake.batchTimeoutReturnsOnCall[i] = struct {
		result1 time.Duration
	}{result1}
}

func (fake *OrdererConfig) Capabilities() channelconfig.OrdererCapabilities {
	fake.capabilitiesMutex.Lock()
	ret, specificReturn := fake.capabilitiesReturnsOnCall[len(fake.capabilitiesArgsForCall)]
	fake.capabilitiesArgsForCall = append(fake.capabilitiesArgsForCall, struct {
	}{})
	fake.recordInvocation("Capabilities", []interface{}{})
	fake.capabilitiesMutex.Unlock()
	if fake.CapabilitiesStub != nil {
		return fake.CapabilitiesStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.capabilitiesReturns
	return fakeReturns.result1
}

func (fake *OrdererConfig) CapabilitiesCallCount() int {
	fake.capabilitiesMutex.RLock()
	defer fake.capabilitiesMutex.RUnlock()
	return len(fake.capabilitiesArgsForCall)
}

func (fake *OrdererConfig) CapabilitiesCalls(stub func() channelconfig.OrdererCapabilities) {
	fake.capabilitiesMutex.Lock()
	defer fake.capabilitiesMutex.Unlock()
	fake.CapabilitiesStub = stub
}

func (fake *OrdererConfig) CapabilitiesReturns(result1 channelconfig.OrdererCapabilities) {
	fake.capabilitiesMutex.Lock()
	defer fake.capabilitiesMutex.Unlock()
	fake.CapabilitiesStub = nil
	fake.capabilitiesReturns = struct {
		result1 channelconfig.OrdererCapabilities
	}{result1}
}

func (fake *OrdererConfig) CapabilitiesReturnsOnCall(i int, result1 channelconfig.OrdererCapabilities) {
	fake.capabilitiesMutex.Lock()
	defer fake.capabilitiesMutex.Unlock()
	fake.CapabilitiesStub = nil
	if fake.capabilitiesReturnsOnCall == nil {
		fake.capabilitiesReturnsOnCall = make(map[int]struct {
			result1 channelconfig.OrdererCapabilities
		})
	}
	fake.capabilitiesReturnsOnCall[i] = struct {
		result1 channelconfig.OrdererCapabilities
	}{result1}
}

func (fake *OrdererConfig) ConsensusMetadata() []byte {
	fake.consensusMetadataMutex.Lock()
	ret, specificReturn := fake.consensusMetadataReturnsOnCall[len(fake.consensusMetadataArgsForCall)]
	fake.consensusMetadataArgsForCall = append(fake.consensusMetadataArgsForCall, struct {
	}{})
	fake.recordInvocation("ConsensusMetadata", []interface{}{})
	fake.consensusMetadataMutex.Unlock()
	if fake.ConsensusMetadataStub != nil {
		return fake.ConsensusMetadataStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.consensusMetadataReturns
	return fakeReturns.result1
}

func (fake *OrdererConfig) ConsensusMetadataCallCount() int {
	fake.consensusMetadataMutex.RLock()
	defer fake.consensusMetadataMutex.RUnlock()
	return len(fake.consensusMetadataArgsForCall)
}

func (fake *OrdererConfig) ConsensusMetadataCalls(stub func() []byte) {
	fake.consensusMetadataMutex.Lock()
	defer fake.consensusMetadataMutex.Unlock()
	fake.ConsensusMetadataStub = stub
}

func (fake *OrdererConfig) ConsensusMetadataReturns(result1 []byte) {
	fake.consensusMetadataMutex.Lock()
	defer fake.consensusMetadataMutex.Unlock()
	fake.ConsensusMetadataStub = nil
	fake.consensusMetadataReturns = struct {
		result1 []byte
	}{result1}
}

func (fake *OrdererConfig) ConsensusMetadataReturnsOnCall(i int, result1 []byte) {
	fake.consensusMetadataMutex.Lock()
	defer fake.consensusMetadataMutex.Unlock()
	fake.ConsensusMetadataStub = nil
	if fake.consensusMetadataReturnsOnCall == nil {
		fake.consensusMetadataReturnsOnCall = make(map[int]struct {
			result1 []byte
		})
	}
	fake.consensusMetadataReturnsOnCall[i] = struct {
		result1 []byte
	}{result1}
}

func (fake *OrdererConfig) ConsensusMigrationContext() uint64 {
	fake.consensusMigrationContextMutex.Lock()
	ret, specificReturn := fake.consensusMigrationContextReturnsOnCall[len(fake.consensusMigrationContextArgsForCall)]
	fake.consensusMigrationContextArgsForCall = append(fake.consensusMigrationContextArgsForCall, struct {
	}{})
	fake.recordInvocation("ConsensusMigrationContext", []interface{}{})
	fake.consensusMigrationContextMutex.Unlock()
	if fake.ConsensusMigrationContextStub != nil {
		return fake.ConsensusMigrationContextStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.consensusMigrationContextReturns
	return fakeReturns.result1
}

func (fake *OrdererConfig) ConsensusMigrationContextCallCount() int {
	fake.consensusMigrationContextMutex.RLock()
	defer fake.consensusMigrationContextMutex.RUnlock()
	return len(fake.consensusMigrationContextArgsForCall)
}

func (fake *OrdererConfig) ConsensusMigrationContextCalls(stub func() uint64) {
	fake.consensusMigrationContextMutex.Lock()
	defer fake.consensusMigrationContextMutex.Unlock()
	fake.ConsensusMigrationContextStub = stub
}

func (fake *OrdererConfig) ConsensusMigrationContextReturns(result1 uint64) {
	fake.consensusMigrationContextMutex.Lock()
	defer fake.consensusMigrationContextMutex.Unlock()
	fake.ConsensusMigrationContextStub = nil
	fake.consensusMigrationContextReturns = struct {
		result1 uint64
	}{result1}
}

func (fake *OrdererConfig) ConsensusMigrationContextReturnsOnCall(i int, result1 uint64) {
	fake.consensusMigrationContextMutex.Lock()
	defer fake.consensusMigrationContextMutex.Unlock()
	fake.ConsensusMigrationContextStub = nil
	if fake.consensusMigrationContextReturnsOnCall == nil {
		fake.consensusMigrationContextReturnsOnCall = make(map[int]struct {
			result1 uint64
		})
	}
	fake.consensusMigrationContextReturnsOnCall[i] = struct {
		result1 uint64
	}{result1}
}

func (fake *OrdererConfig) ConsensusMigrationState() orderer.ConsensusType_MigrationState {
	fake.consensusMigrationStateMutex.Lock()
	ret, specificReturn := fake.consensusMigrationStateReturnsOnCall[len(fake.consensusMigrationStateArgsForCall)]
	fake.consensusMigrationStateArgsForCall = append(fake.consensusMigrationStateArgsForCall, struct {
	}{})
	fake.recordInvocation("ConsensusMigrationState", []interface{}{})
	fake.consensusMigrationStateMutex.Unlock()
	if fake.ConsensusMigrationStateStub != nil {
		return fake.ConsensusMigrationStateStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.consensusMigrationStateReturns
	return fakeReturns.result1
}

func (fake *OrdererConfig) ConsensusMigrationStateCallCount() int {
	fake.consensusMigrationStateMutex.RLock()
	defer fake.consensusMigrationStateMutex.RUnlock()
	return len(fake.consensusMigrationStateArgsForCall)
}

func (fake *OrdererConfig) ConsensusMigrationStateCalls(stub func() orderer.ConsensusType_MigrationState) {
	fake.consensusMigrationStateMutex.Lock()
	defer fake.consensusMigrationStateMutex.Unlock()
	fake.ConsensusMigrationStateStub = stub
}

func (fake *OrdererConfig) ConsensusMigrationStateReturns(result1 orderer.ConsensusType_MigrationState) {
	fake.consensusMigrationStateMutex.Lock()
	defer fake.consensusMigrationStateMutex.Unlock()
	fake.ConsensusMigrationStateStub = nil
	fake.consensusMigrationStateReturns = struct {
		result1 orderer.ConsensusType_MigrationState
	}{result1}
}

func (fake *OrdererConfig) ConsensusMigrationStateReturnsOnCall(i int, result1 orderer.ConsensusType_MigrationState) {
	fake.consensusMigrationStateMutex.Lock()
	defer fake.consensusMigrationStateMutex.Unlock()
	fake.ConsensusMigrationStateStub = nil
	if fake.consensusMigrationStateReturnsOnCall == nil {
		fake.consensusMigrationStateReturnsOnCall = make(map[int]struct {
			result1 orderer.ConsensusType_MigrationState
		})
	}
	fake.consensusMigrationStateReturnsOnCall[i] = struct {
		result1 orderer.ConsensusType_MigrationState
	}{result1}
}

func (fake *OrdererConfig) ConsensusType() string {
	fake.consensusTypeMutex.Lock()
	ret, specificReturn := fake.consensusTypeReturnsOnCall[len(fake.consensusTypeArgsForCall)]
	fake.consensusTypeArgsForCall = append(fake.consensusTypeArgsForCall, struct {
	}{})
	fake.recordInvocation("ConsensusType", []interface{}{})
	fake.consensusTypeMutex.Unlock()
	if fake.ConsensusTypeStub != nil {
		return fake.ConsensusTypeStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.consensusTypeReturns
	return fakeReturns.result1
}

func (fake *OrdererConfig) ConsensusTypeCallCount() int {
	fake.consensusTypeMutex.RLock()
	defer fake.consensusTypeMutex.RUnlock()
	return len(fake.consensusTypeArgsForCall)
}

func (fake *OrdererConfig) ConsensusTypeCalls(stub func() string) {
	fake.consensusTypeMutex.Lock()
	defer fake.consensusTypeMutex.Unlock()
	fake.ConsensusTypeStub = stub
}

func (fake *OrdererConfig) ConsensusTypeReturns(result1 string) {
	fake.consensusTypeMutex.Lock()
	defer fake.consensusTypeMutex.Unlock()
	fake.ConsensusTypeStub = nil
	fake.consensusTypeReturns = struct {
		result1 string
	}{result1}
}

func (fake *OrdererConfig) ConsensusTypeReturnsOnCall(i int, result1 string) {
	fake.consensusTypeMutex.Lock()
	defer fake.consensusTypeMutex.Unlock()
	fake.ConsensusTypeStub = nil
	if fake.consensusTypeReturnsOnCall == nil {
		fake.consensusTypeReturnsOnCall = make(map[int]struct {
			result1 string
		})
	}
	fake.consensusTypeReturnsOnCall[i] = struct {
		result1 string
	}{result1}
}

func (fake *OrdererConfig) KafkaBrokers() []string {
	fake.kafkaBrokersMutex.Lock()
	ret, specificReturn := fake.kafkaBrokersReturnsOnCall[len(fake.kafkaBrokersArgsForCall)]
	fake.kafkaBrokersArgsForCall = append(fake.kafkaBrokersArgsForCall, struct {
	}{})
	fake.recordInvocation("KafkaBrokers", []interface{}{})
	fake.kafkaBrokersMutex.Unlock()
	if fake.KafkaBrokersStub != nil {
		return fake.KafkaBrokersStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.kafkaBrokersReturns
	return fakeReturns.result1
}

func (fake *OrdererConfig) KafkaBrokersCallCount() int {
	fake.kafkaBrokersMutex.RLock()
	defer fake.kafkaBrokersMutex.RUnlock()
	return len(fake.kafkaBrokersArgsForCall)
}

func (fake *OrdererConfig) KafkaBrokersCalls(stub func() []string) {
	fake.kafkaBrokersMutex.Lock()
	defer fake.kafkaBrokersMutex.Unlock()
	fake.KafkaBrokersStub = stub
}

func (fake *OrdererConfig) KafkaBrokersReturns(result1 []string) {
	fake.kafkaBrokersMutex.Lock()
	defer fake.kafkaBrokersMutex.Unlock()
	fake.KafkaBrokersStub = nil
	fake.kafkaBrokersReturns = struct {
		result1 []string
	}{result1}
}

func (fake *OrdererConfig) KafkaBrokersReturnsOnCall(i int, result1 []string) {
	fake.kafkaBrokersMutex.Lock()
	defer fake.kafkaBrokersMutex.Unlock()
	fake.KafkaBrokersStub = nil
	if fake.kafkaBrokersReturnsOnCall == nil {
		fake.kafkaBrokersReturnsOnCall = make(map[int]struct {
			result1 []string
		})
	}
	fake.kafkaBrokersReturnsOnCall[i] = struct {
		result1 []string
	}{result1}
}

func (fake *OrdererConfig) MaxChannelsCount() uint64 {
	fake.maxChannelsCountMutex.Lock()
	ret, specificReturn := fake.maxChannelsCountReturnsOnCall[len(fake.maxChannelsCountArgsForCall)]
	fake.maxChannelsCountArgsForCall = append(fake.maxChannelsCountArgsForCall, struct {
	}{})
	fake.recordInvocation("MaxChannelsCount", []interface{}{})
	fake.maxChannelsCountMutex.Unlock()
	if fake.MaxChannelsCountStub != nil {
		return fake.MaxChannelsCountStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.maxChannelsCountReturns
	return fakeReturns.result1
}

func (fake *OrdererConfig) MaxChannelsCountCallCount() int {
	fake.maxChannelsCountMutex.RLock()
	defer fake.maxChannelsCountMutex.RUnlock()
	return len(fake.maxChannelsCountArgsForCall)
}

func (fake *OrdererConfig) MaxChannelsCountCalls(stub func() uint64) {
	fake.maxChannelsCountMutex.Lock()
	defer fake.maxChannelsCountMutex.Unlock()
	fake.MaxChannelsCountStub = stub
}

func (fake *OrdererConfig) MaxChannelsCountReturns(result1 uint64) {
	fake.maxChannelsCountMutex.Lock()
	defer fake.maxChannelsCountMutex.Unlock()
	fake.MaxChannelsCountStub = nil
	fake.maxChannelsCountReturns = struct {
		result1 uint64
	}{result1}
}

func (fake *OrdererConfig) MaxChannelsCountReturnsOnCall(i int, result1 uint64) {
	fake.maxChannelsCountMutex.Lock()
	defer fake.maxChannelsCountMutex.Unlock()
	fake.MaxChannelsCountStub = nil
	if fake.maxChannelsCountReturnsOnCall == nil {
		fake.maxChannelsCountReturnsOnCall = make(map[int]struct {
			result1 uint64
		})
	}
	fake.maxChannelsCountReturnsOnCall[i] = struct {
		result1 uint64
	}{result1}
}

func (fake *OrdererConfig) Organizations() map[string]channelconfig.OrdererOrg {
	fake.organizationsMutex.Lock()
	ret, specificReturn := fake.organizationsReturnsOnCall[len(fake.organizationsArgsForCall)]
	fake.organizationsArgsForCall = append(fake.organizationsArgsForCall, struct {
	}{})
	fake.recordInvocation("Organizations", []interface{}{})
	fake.organizationsMutex.Unlock()
	if fake.OrganizationsStub != nil {
		return fake.OrganizationsStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.organizationsReturns
	return fakeReturns.result1
}

func (fake *OrdererConfig) OrganizationsCallCount() int {
	fake.organizationsMutex.RLock()
	defer fake.organizationsMutex.RUnlock()
	return len(fake.organizationsArgsForCall)
}

func (fake *OrdererConfig) OrganizationsCalls(stub func() map[string]channelconfig.OrdererOrg) {
	fake.organizationsMutex.Lock()
	defer fake.organizationsMutex.Unlock()
	fake.OrganizationsStub = stub
}

func (fake *OrdererConfig) OrganizationsReturns(result1 map[string]channelconfig.OrdererOrg) {
	fake.organizationsMutex.Lock()
	defer fake.organizationsMutex.Unlock()
	fake.OrganizationsStub = nil
	fake.organizationsReturns = struct {
		result1 map[string]channelconfig.OrdererOrg
	}{result1}
}

func (fake *OrdererConfig) OrganizationsReturnsOnCall(i int, result1 map[string]channelconfig.OrdererOrg) {
	fake.organizationsMutex.Lock()
	defer fake.organizationsMutex.Unlock()
	fake.OrganizationsStub = nil
	if fake.organizationsReturnsOnCall == nil {
		fake.organizationsReturnsOnCall = make(map[int]struct {
			result1 map[string]channelconfig.OrdererOrg
		})
	}
	fake.organizationsReturnsOnCall[i] = struct {
		result1 map[string]channelconfig.OrdererOrg
	}{result1}
}

func (fake *OrdererConfig) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.batchSizeMutex.RLock()
	defer fake.batchSizeMutex.RUnlock()
	fake.batchTimeoutMutex.RLock()
	defer fake.batchTimeoutMutex.RUnlock()
	fake.capabilitiesMutex.RLock()
	defer fake.capabilitiesMutex.RUnlock()
	fake.consensusMetadataMutex.RLock()
	defer fake.consensusMetadataMutex.RUnlock()
	fake.consensusMigrationContextMutex.RLock()
	defer fake.consensusMigrationContextMutex.RUnlock()
	fake.consensusMigrationStateMutex.RLock()
	defer fake.consensusMigrationStateMutex.RUnlock()
	fake.consensusTypeMutex.RLock()
	defer fake.consensusTypeMutex.RUnlock()
	fake.kafkaBrokersMutex.RLock()
	defer fake.kafkaBrokersMutex.RUnlock()
	fake.maxChannelsCountMutex.RLock()
	defer fake.maxChannelsCountMutex.RUnlock()
	fake.organizationsMutex.RLock()
	defer fake.organizationsMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *OrdererConfig) recordInvocation(key string, args []interface{}) {
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
