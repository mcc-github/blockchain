
package mock

import (
	"sync"
)

type ApplicationCapabilities struct {
	ACLsStub        func() bool
	aCLsMutex       sync.RWMutex
	aCLsArgsForCall []struct {
	}
	aCLsReturns struct {
		result1 bool
	}
	aCLsReturnsOnCall map[int]struct {
		result1 bool
	}
	CollectionUpgradeStub        func() bool
	collectionUpgradeMutex       sync.RWMutex
	collectionUpgradeArgsForCall []struct {
	}
	collectionUpgradeReturns struct {
		result1 bool
	}
	collectionUpgradeReturnsOnCall map[int]struct {
		result1 bool
	}
	ForbidDuplicateTXIdInBlockStub        func() bool
	forbidDuplicateTXIdInBlockMutex       sync.RWMutex
	forbidDuplicateTXIdInBlockArgsForCall []struct {
	}
	forbidDuplicateTXIdInBlockReturns struct {
		result1 bool
	}
	forbidDuplicateTXIdInBlockReturnsOnCall map[int]struct {
		result1 bool
	}
	KeyLevelEndorsementStub        func() bool
	keyLevelEndorsementMutex       sync.RWMutex
	keyLevelEndorsementArgsForCall []struct {
	}
	keyLevelEndorsementReturns struct {
		result1 bool
	}
	keyLevelEndorsementReturnsOnCall map[int]struct {
		result1 bool
	}
	LifecycleV20Stub        func() bool
	lifecycleV20Mutex       sync.RWMutex
	lifecycleV20ArgsForCall []struct {
	}
	lifecycleV20Returns struct {
		result1 bool
	}
	lifecycleV20ReturnsOnCall map[int]struct {
		result1 bool
	}
	MetadataLifecycleStub        func() bool
	metadataLifecycleMutex       sync.RWMutex
	metadataLifecycleArgsForCall []struct {
	}
	metadataLifecycleReturns struct {
		result1 bool
	}
	metadataLifecycleReturnsOnCall map[int]struct {
		result1 bool
	}
	PrivateChannelDataStub        func() bool
	privateChannelDataMutex       sync.RWMutex
	privateChannelDataArgsForCall []struct {
	}
	privateChannelDataReturns struct {
		result1 bool
	}
	privateChannelDataReturnsOnCall map[int]struct {
		result1 bool
	}
	StorePvtDataOfInvalidTxStub        func() bool
	storePvtDataOfInvalidTxMutex       sync.RWMutex
	storePvtDataOfInvalidTxArgsForCall []struct {
	}
	storePvtDataOfInvalidTxReturns struct {
		result1 bool
	}
	storePvtDataOfInvalidTxReturnsOnCall map[int]struct {
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
	V1_1ValidationStub        func() bool
	v1_1ValidationMutex       sync.RWMutex
	v1_1ValidationArgsForCall []struct {
	}
	v1_1ValidationReturns struct {
		result1 bool
	}
	v1_1ValidationReturnsOnCall map[int]struct {
		result1 bool
	}
	V1_2ValidationStub        func() bool
	v1_2ValidationMutex       sync.RWMutex
	v1_2ValidationArgsForCall []struct {
	}
	v1_2ValidationReturns struct {
		result1 bool
	}
	v1_2ValidationReturnsOnCall map[int]struct {
		result1 bool
	}
	V1_3ValidationStub        func() bool
	v1_3ValidationMutex       sync.RWMutex
	v1_3ValidationArgsForCall []struct {
	}
	v1_3ValidationReturns struct {
		result1 bool
	}
	v1_3ValidationReturnsOnCall map[int]struct {
		result1 bool
	}
	V2_0ValidationStub        func() bool
	v2_0ValidationMutex       sync.RWMutex
	v2_0ValidationArgsForCall []struct {
	}
	v2_0ValidationReturns struct {
		result1 bool
	}
	v2_0ValidationReturnsOnCall map[int]struct {
		result1 bool
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *ApplicationCapabilities) ACLs() bool {
	fake.aCLsMutex.Lock()
	ret, specificReturn := fake.aCLsReturnsOnCall[len(fake.aCLsArgsForCall)]
	fake.aCLsArgsForCall = append(fake.aCLsArgsForCall, struct {
	}{})
	fake.recordInvocation("ACLs", []interface{}{})
	fake.aCLsMutex.Unlock()
	if fake.ACLsStub != nil {
		return fake.ACLsStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.aCLsReturns
	return fakeReturns.result1
}

func (fake *ApplicationCapabilities) ACLsCallCount() int {
	fake.aCLsMutex.RLock()
	defer fake.aCLsMutex.RUnlock()
	return len(fake.aCLsArgsForCall)
}

func (fake *ApplicationCapabilities) ACLsCalls(stub func() bool) {
	fake.aCLsMutex.Lock()
	defer fake.aCLsMutex.Unlock()
	fake.ACLsStub = stub
}

func (fake *ApplicationCapabilities) ACLsReturns(result1 bool) {
	fake.aCLsMutex.Lock()
	defer fake.aCLsMutex.Unlock()
	fake.ACLsStub = nil
	fake.aCLsReturns = struct {
		result1 bool
	}{result1}
}

func (fake *ApplicationCapabilities) ACLsReturnsOnCall(i int, result1 bool) {
	fake.aCLsMutex.Lock()
	defer fake.aCLsMutex.Unlock()
	fake.ACLsStub = nil
	if fake.aCLsReturnsOnCall == nil {
		fake.aCLsReturnsOnCall = make(map[int]struct {
			result1 bool
		})
	}
	fake.aCLsReturnsOnCall[i] = struct {
		result1 bool
	}{result1}
}

func (fake *ApplicationCapabilities) CollectionUpgrade() bool {
	fake.collectionUpgradeMutex.Lock()
	ret, specificReturn := fake.collectionUpgradeReturnsOnCall[len(fake.collectionUpgradeArgsForCall)]
	fake.collectionUpgradeArgsForCall = append(fake.collectionUpgradeArgsForCall, struct {
	}{})
	fake.recordInvocation("CollectionUpgrade", []interface{}{})
	fake.collectionUpgradeMutex.Unlock()
	if fake.CollectionUpgradeStub != nil {
		return fake.CollectionUpgradeStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.collectionUpgradeReturns
	return fakeReturns.result1
}

func (fake *ApplicationCapabilities) CollectionUpgradeCallCount() int {
	fake.collectionUpgradeMutex.RLock()
	defer fake.collectionUpgradeMutex.RUnlock()
	return len(fake.collectionUpgradeArgsForCall)
}

func (fake *ApplicationCapabilities) CollectionUpgradeCalls(stub func() bool) {
	fake.collectionUpgradeMutex.Lock()
	defer fake.collectionUpgradeMutex.Unlock()
	fake.CollectionUpgradeStub = stub
}

func (fake *ApplicationCapabilities) CollectionUpgradeReturns(result1 bool) {
	fake.collectionUpgradeMutex.Lock()
	defer fake.collectionUpgradeMutex.Unlock()
	fake.CollectionUpgradeStub = nil
	fake.collectionUpgradeReturns = struct {
		result1 bool
	}{result1}
}

func (fake *ApplicationCapabilities) CollectionUpgradeReturnsOnCall(i int, result1 bool) {
	fake.collectionUpgradeMutex.Lock()
	defer fake.collectionUpgradeMutex.Unlock()
	fake.CollectionUpgradeStub = nil
	if fake.collectionUpgradeReturnsOnCall == nil {
		fake.collectionUpgradeReturnsOnCall = make(map[int]struct {
			result1 bool
		})
	}
	fake.collectionUpgradeReturnsOnCall[i] = struct {
		result1 bool
	}{result1}
}

func (fake *ApplicationCapabilities) ForbidDuplicateTXIdInBlock() bool {
	fake.forbidDuplicateTXIdInBlockMutex.Lock()
	ret, specificReturn := fake.forbidDuplicateTXIdInBlockReturnsOnCall[len(fake.forbidDuplicateTXIdInBlockArgsForCall)]
	fake.forbidDuplicateTXIdInBlockArgsForCall = append(fake.forbidDuplicateTXIdInBlockArgsForCall, struct {
	}{})
	fake.recordInvocation("ForbidDuplicateTXIdInBlock", []interface{}{})
	fake.forbidDuplicateTXIdInBlockMutex.Unlock()
	if fake.ForbidDuplicateTXIdInBlockStub != nil {
		return fake.ForbidDuplicateTXIdInBlockStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.forbidDuplicateTXIdInBlockReturns
	return fakeReturns.result1
}

func (fake *ApplicationCapabilities) ForbidDuplicateTXIdInBlockCallCount() int {
	fake.forbidDuplicateTXIdInBlockMutex.RLock()
	defer fake.forbidDuplicateTXIdInBlockMutex.RUnlock()
	return len(fake.forbidDuplicateTXIdInBlockArgsForCall)
}

func (fake *ApplicationCapabilities) ForbidDuplicateTXIdInBlockCalls(stub func() bool) {
	fake.forbidDuplicateTXIdInBlockMutex.Lock()
	defer fake.forbidDuplicateTXIdInBlockMutex.Unlock()
	fake.ForbidDuplicateTXIdInBlockStub = stub
}

func (fake *ApplicationCapabilities) ForbidDuplicateTXIdInBlockReturns(result1 bool) {
	fake.forbidDuplicateTXIdInBlockMutex.Lock()
	defer fake.forbidDuplicateTXIdInBlockMutex.Unlock()
	fake.ForbidDuplicateTXIdInBlockStub = nil
	fake.forbidDuplicateTXIdInBlockReturns = struct {
		result1 bool
	}{result1}
}

func (fake *ApplicationCapabilities) ForbidDuplicateTXIdInBlockReturnsOnCall(i int, result1 bool) {
	fake.forbidDuplicateTXIdInBlockMutex.Lock()
	defer fake.forbidDuplicateTXIdInBlockMutex.Unlock()
	fake.ForbidDuplicateTXIdInBlockStub = nil
	if fake.forbidDuplicateTXIdInBlockReturnsOnCall == nil {
		fake.forbidDuplicateTXIdInBlockReturnsOnCall = make(map[int]struct {
			result1 bool
		})
	}
	fake.forbidDuplicateTXIdInBlockReturnsOnCall[i] = struct {
		result1 bool
	}{result1}
}

func (fake *ApplicationCapabilities) KeyLevelEndorsement() bool {
	fake.keyLevelEndorsementMutex.Lock()
	ret, specificReturn := fake.keyLevelEndorsementReturnsOnCall[len(fake.keyLevelEndorsementArgsForCall)]
	fake.keyLevelEndorsementArgsForCall = append(fake.keyLevelEndorsementArgsForCall, struct {
	}{})
	fake.recordInvocation("KeyLevelEndorsement", []interface{}{})
	fake.keyLevelEndorsementMutex.Unlock()
	if fake.KeyLevelEndorsementStub != nil {
		return fake.KeyLevelEndorsementStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.keyLevelEndorsementReturns
	return fakeReturns.result1
}

func (fake *ApplicationCapabilities) KeyLevelEndorsementCallCount() int {
	fake.keyLevelEndorsementMutex.RLock()
	defer fake.keyLevelEndorsementMutex.RUnlock()
	return len(fake.keyLevelEndorsementArgsForCall)
}

func (fake *ApplicationCapabilities) KeyLevelEndorsementCalls(stub func() bool) {
	fake.keyLevelEndorsementMutex.Lock()
	defer fake.keyLevelEndorsementMutex.Unlock()
	fake.KeyLevelEndorsementStub = stub
}

func (fake *ApplicationCapabilities) KeyLevelEndorsementReturns(result1 bool) {
	fake.keyLevelEndorsementMutex.Lock()
	defer fake.keyLevelEndorsementMutex.Unlock()
	fake.KeyLevelEndorsementStub = nil
	fake.keyLevelEndorsementReturns = struct {
		result1 bool
	}{result1}
}

func (fake *ApplicationCapabilities) KeyLevelEndorsementReturnsOnCall(i int, result1 bool) {
	fake.keyLevelEndorsementMutex.Lock()
	defer fake.keyLevelEndorsementMutex.Unlock()
	fake.KeyLevelEndorsementStub = nil
	if fake.keyLevelEndorsementReturnsOnCall == nil {
		fake.keyLevelEndorsementReturnsOnCall = make(map[int]struct {
			result1 bool
		})
	}
	fake.keyLevelEndorsementReturnsOnCall[i] = struct {
		result1 bool
	}{result1}
}

func (fake *ApplicationCapabilities) LifecycleV20() bool {
	fake.lifecycleV20Mutex.Lock()
	ret, specificReturn := fake.lifecycleV20ReturnsOnCall[len(fake.lifecycleV20ArgsForCall)]
	fake.lifecycleV20ArgsForCall = append(fake.lifecycleV20ArgsForCall, struct {
	}{})
	fake.recordInvocation("LifecycleV20", []interface{}{})
	fake.lifecycleV20Mutex.Unlock()
	if fake.LifecycleV20Stub != nil {
		return fake.LifecycleV20Stub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.lifecycleV20Returns
	return fakeReturns.result1
}

func (fake *ApplicationCapabilities) LifecycleV20CallCount() int {
	fake.lifecycleV20Mutex.RLock()
	defer fake.lifecycleV20Mutex.RUnlock()
	return len(fake.lifecycleV20ArgsForCall)
}

func (fake *ApplicationCapabilities) LifecycleV20Calls(stub func() bool) {
	fake.lifecycleV20Mutex.Lock()
	defer fake.lifecycleV20Mutex.Unlock()
	fake.LifecycleV20Stub = stub
}

func (fake *ApplicationCapabilities) LifecycleV20Returns(result1 bool) {
	fake.lifecycleV20Mutex.Lock()
	defer fake.lifecycleV20Mutex.Unlock()
	fake.LifecycleV20Stub = nil
	fake.lifecycleV20Returns = struct {
		result1 bool
	}{result1}
}

func (fake *ApplicationCapabilities) LifecycleV20ReturnsOnCall(i int, result1 bool) {
	fake.lifecycleV20Mutex.Lock()
	defer fake.lifecycleV20Mutex.Unlock()
	fake.LifecycleV20Stub = nil
	if fake.lifecycleV20ReturnsOnCall == nil {
		fake.lifecycleV20ReturnsOnCall = make(map[int]struct {
			result1 bool
		})
	}
	fake.lifecycleV20ReturnsOnCall[i] = struct {
		result1 bool
	}{result1}
}

func (fake *ApplicationCapabilities) MetadataLifecycle() bool {
	fake.metadataLifecycleMutex.Lock()
	ret, specificReturn := fake.metadataLifecycleReturnsOnCall[len(fake.metadataLifecycleArgsForCall)]
	fake.metadataLifecycleArgsForCall = append(fake.metadataLifecycleArgsForCall, struct {
	}{})
	fake.recordInvocation("MetadataLifecycle", []interface{}{})
	fake.metadataLifecycleMutex.Unlock()
	if fake.MetadataLifecycleStub != nil {
		return fake.MetadataLifecycleStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.metadataLifecycleReturns
	return fakeReturns.result1
}

func (fake *ApplicationCapabilities) MetadataLifecycleCallCount() int {
	fake.metadataLifecycleMutex.RLock()
	defer fake.metadataLifecycleMutex.RUnlock()
	return len(fake.metadataLifecycleArgsForCall)
}

func (fake *ApplicationCapabilities) MetadataLifecycleCalls(stub func() bool) {
	fake.metadataLifecycleMutex.Lock()
	defer fake.metadataLifecycleMutex.Unlock()
	fake.MetadataLifecycleStub = stub
}

func (fake *ApplicationCapabilities) MetadataLifecycleReturns(result1 bool) {
	fake.metadataLifecycleMutex.Lock()
	defer fake.metadataLifecycleMutex.Unlock()
	fake.MetadataLifecycleStub = nil
	fake.metadataLifecycleReturns = struct {
		result1 bool
	}{result1}
}

func (fake *ApplicationCapabilities) MetadataLifecycleReturnsOnCall(i int, result1 bool) {
	fake.metadataLifecycleMutex.Lock()
	defer fake.metadataLifecycleMutex.Unlock()
	fake.MetadataLifecycleStub = nil
	if fake.metadataLifecycleReturnsOnCall == nil {
		fake.metadataLifecycleReturnsOnCall = make(map[int]struct {
			result1 bool
		})
	}
	fake.metadataLifecycleReturnsOnCall[i] = struct {
		result1 bool
	}{result1}
}

func (fake *ApplicationCapabilities) PrivateChannelData() bool {
	fake.privateChannelDataMutex.Lock()
	ret, specificReturn := fake.privateChannelDataReturnsOnCall[len(fake.privateChannelDataArgsForCall)]
	fake.privateChannelDataArgsForCall = append(fake.privateChannelDataArgsForCall, struct {
	}{})
	fake.recordInvocation("PrivateChannelData", []interface{}{})
	fake.privateChannelDataMutex.Unlock()
	if fake.PrivateChannelDataStub != nil {
		return fake.PrivateChannelDataStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.privateChannelDataReturns
	return fakeReturns.result1
}

func (fake *ApplicationCapabilities) PrivateChannelDataCallCount() int {
	fake.privateChannelDataMutex.RLock()
	defer fake.privateChannelDataMutex.RUnlock()
	return len(fake.privateChannelDataArgsForCall)
}

func (fake *ApplicationCapabilities) PrivateChannelDataCalls(stub func() bool) {
	fake.privateChannelDataMutex.Lock()
	defer fake.privateChannelDataMutex.Unlock()
	fake.PrivateChannelDataStub = stub
}

func (fake *ApplicationCapabilities) PrivateChannelDataReturns(result1 bool) {
	fake.privateChannelDataMutex.Lock()
	defer fake.privateChannelDataMutex.Unlock()
	fake.PrivateChannelDataStub = nil
	fake.privateChannelDataReturns = struct {
		result1 bool
	}{result1}
}

func (fake *ApplicationCapabilities) PrivateChannelDataReturnsOnCall(i int, result1 bool) {
	fake.privateChannelDataMutex.Lock()
	defer fake.privateChannelDataMutex.Unlock()
	fake.PrivateChannelDataStub = nil
	if fake.privateChannelDataReturnsOnCall == nil {
		fake.privateChannelDataReturnsOnCall = make(map[int]struct {
			result1 bool
		})
	}
	fake.privateChannelDataReturnsOnCall[i] = struct {
		result1 bool
	}{result1}
}

func (fake *ApplicationCapabilities) StorePvtDataOfInvalidTx() bool {
	fake.storePvtDataOfInvalidTxMutex.Lock()
	ret, specificReturn := fake.storePvtDataOfInvalidTxReturnsOnCall[len(fake.storePvtDataOfInvalidTxArgsForCall)]
	fake.storePvtDataOfInvalidTxArgsForCall = append(fake.storePvtDataOfInvalidTxArgsForCall, struct {
	}{})
	fake.recordInvocation("StorePvtDataOfInvalidTx", []interface{}{})
	fake.storePvtDataOfInvalidTxMutex.Unlock()
	if fake.StorePvtDataOfInvalidTxStub != nil {
		return fake.StorePvtDataOfInvalidTxStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.storePvtDataOfInvalidTxReturns
	return fakeReturns.result1
}

func (fake *ApplicationCapabilities) StorePvtDataOfInvalidTxCallCount() int {
	fake.storePvtDataOfInvalidTxMutex.RLock()
	defer fake.storePvtDataOfInvalidTxMutex.RUnlock()
	return len(fake.storePvtDataOfInvalidTxArgsForCall)
}

func (fake *ApplicationCapabilities) StorePvtDataOfInvalidTxCalls(stub func() bool) {
	fake.storePvtDataOfInvalidTxMutex.Lock()
	defer fake.storePvtDataOfInvalidTxMutex.Unlock()
	fake.StorePvtDataOfInvalidTxStub = stub
}

func (fake *ApplicationCapabilities) StorePvtDataOfInvalidTxReturns(result1 bool) {
	fake.storePvtDataOfInvalidTxMutex.Lock()
	defer fake.storePvtDataOfInvalidTxMutex.Unlock()
	fake.StorePvtDataOfInvalidTxStub = nil
	fake.storePvtDataOfInvalidTxReturns = struct {
		result1 bool
	}{result1}
}

func (fake *ApplicationCapabilities) StorePvtDataOfInvalidTxReturnsOnCall(i int, result1 bool) {
	fake.storePvtDataOfInvalidTxMutex.Lock()
	defer fake.storePvtDataOfInvalidTxMutex.Unlock()
	fake.StorePvtDataOfInvalidTxStub = nil
	if fake.storePvtDataOfInvalidTxReturnsOnCall == nil {
		fake.storePvtDataOfInvalidTxReturnsOnCall = make(map[int]struct {
			result1 bool
		})
	}
	fake.storePvtDataOfInvalidTxReturnsOnCall[i] = struct {
		result1 bool
	}{result1}
}

func (fake *ApplicationCapabilities) Supported() error {
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

func (fake *ApplicationCapabilities) SupportedCallCount() int {
	fake.supportedMutex.RLock()
	defer fake.supportedMutex.RUnlock()
	return len(fake.supportedArgsForCall)
}

func (fake *ApplicationCapabilities) SupportedCalls(stub func() error) {
	fake.supportedMutex.Lock()
	defer fake.supportedMutex.Unlock()
	fake.SupportedStub = stub
}

func (fake *ApplicationCapabilities) SupportedReturns(result1 error) {
	fake.supportedMutex.Lock()
	defer fake.supportedMutex.Unlock()
	fake.SupportedStub = nil
	fake.supportedReturns = struct {
		result1 error
	}{result1}
}

func (fake *ApplicationCapabilities) SupportedReturnsOnCall(i int, result1 error) {
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

func (fake *ApplicationCapabilities) V1_1Validation() bool {
	fake.v1_1ValidationMutex.Lock()
	ret, specificReturn := fake.v1_1ValidationReturnsOnCall[len(fake.v1_1ValidationArgsForCall)]
	fake.v1_1ValidationArgsForCall = append(fake.v1_1ValidationArgsForCall, struct {
	}{})
	fake.recordInvocation("V1_1Validation", []interface{}{})
	fake.v1_1ValidationMutex.Unlock()
	if fake.V1_1ValidationStub != nil {
		return fake.V1_1ValidationStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.v1_1ValidationReturns
	return fakeReturns.result1
}

func (fake *ApplicationCapabilities) V1_1ValidationCallCount() int {
	fake.v1_1ValidationMutex.RLock()
	defer fake.v1_1ValidationMutex.RUnlock()
	return len(fake.v1_1ValidationArgsForCall)
}

func (fake *ApplicationCapabilities) V1_1ValidationCalls(stub func() bool) {
	fake.v1_1ValidationMutex.Lock()
	defer fake.v1_1ValidationMutex.Unlock()
	fake.V1_1ValidationStub = stub
}

func (fake *ApplicationCapabilities) V1_1ValidationReturns(result1 bool) {
	fake.v1_1ValidationMutex.Lock()
	defer fake.v1_1ValidationMutex.Unlock()
	fake.V1_1ValidationStub = nil
	fake.v1_1ValidationReturns = struct {
		result1 bool
	}{result1}
}

func (fake *ApplicationCapabilities) V1_1ValidationReturnsOnCall(i int, result1 bool) {
	fake.v1_1ValidationMutex.Lock()
	defer fake.v1_1ValidationMutex.Unlock()
	fake.V1_1ValidationStub = nil
	if fake.v1_1ValidationReturnsOnCall == nil {
		fake.v1_1ValidationReturnsOnCall = make(map[int]struct {
			result1 bool
		})
	}
	fake.v1_1ValidationReturnsOnCall[i] = struct {
		result1 bool
	}{result1}
}

func (fake *ApplicationCapabilities) V1_2Validation() bool {
	fake.v1_2ValidationMutex.Lock()
	ret, specificReturn := fake.v1_2ValidationReturnsOnCall[len(fake.v1_2ValidationArgsForCall)]
	fake.v1_2ValidationArgsForCall = append(fake.v1_2ValidationArgsForCall, struct {
	}{})
	fake.recordInvocation("V1_2Validation", []interface{}{})
	fake.v1_2ValidationMutex.Unlock()
	if fake.V1_2ValidationStub != nil {
		return fake.V1_2ValidationStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.v1_2ValidationReturns
	return fakeReturns.result1
}

func (fake *ApplicationCapabilities) V1_2ValidationCallCount() int {
	fake.v1_2ValidationMutex.RLock()
	defer fake.v1_2ValidationMutex.RUnlock()
	return len(fake.v1_2ValidationArgsForCall)
}

func (fake *ApplicationCapabilities) V1_2ValidationCalls(stub func() bool) {
	fake.v1_2ValidationMutex.Lock()
	defer fake.v1_2ValidationMutex.Unlock()
	fake.V1_2ValidationStub = stub
}

func (fake *ApplicationCapabilities) V1_2ValidationReturns(result1 bool) {
	fake.v1_2ValidationMutex.Lock()
	defer fake.v1_2ValidationMutex.Unlock()
	fake.V1_2ValidationStub = nil
	fake.v1_2ValidationReturns = struct {
		result1 bool
	}{result1}
}

func (fake *ApplicationCapabilities) V1_2ValidationReturnsOnCall(i int, result1 bool) {
	fake.v1_2ValidationMutex.Lock()
	defer fake.v1_2ValidationMutex.Unlock()
	fake.V1_2ValidationStub = nil
	if fake.v1_2ValidationReturnsOnCall == nil {
		fake.v1_2ValidationReturnsOnCall = make(map[int]struct {
			result1 bool
		})
	}
	fake.v1_2ValidationReturnsOnCall[i] = struct {
		result1 bool
	}{result1}
}

func (fake *ApplicationCapabilities) V1_3Validation() bool {
	fake.v1_3ValidationMutex.Lock()
	ret, specificReturn := fake.v1_3ValidationReturnsOnCall[len(fake.v1_3ValidationArgsForCall)]
	fake.v1_3ValidationArgsForCall = append(fake.v1_3ValidationArgsForCall, struct {
	}{})
	fake.recordInvocation("V1_3Validation", []interface{}{})
	fake.v1_3ValidationMutex.Unlock()
	if fake.V1_3ValidationStub != nil {
		return fake.V1_3ValidationStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.v1_3ValidationReturns
	return fakeReturns.result1
}

func (fake *ApplicationCapabilities) V1_3ValidationCallCount() int {
	fake.v1_3ValidationMutex.RLock()
	defer fake.v1_3ValidationMutex.RUnlock()
	return len(fake.v1_3ValidationArgsForCall)
}

func (fake *ApplicationCapabilities) V1_3ValidationCalls(stub func() bool) {
	fake.v1_3ValidationMutex.Lock()
	defer fake.v1_3ValidationMutex.Unlock()
	fake.V1_3ValidationStub = stub
}

func (fake *ApplicationCapabilities) V1_3ValidationReturns(result1 bool) {
	fake.v1_3ValidationMutex.Lock()
	defer fake.v1_3ValidationMutex.Unlock()
	fake.V1_3ValidationStub = nil
	fake.v1_3ValidationReturns = struct {
		result1 bool
	}{result1}
}

func (fake *ApplicationCapabilities) V1_3ValidationReturnsOnCall(i int, result1 bool) {
	fake.v1_3ValidationMutex.Lock()
	defer fake.v1_3ValidationMutex.Unlock()
	fake.V1_3ValidationStub = nil
	if fake.v1_3ValidationReturnsOnCall == nil {
		fake.v1_3ValidationReturnsOnCall = make(map[int]struct {
			result1 bool
		})
	}
	fake.v1_3ValidationReturnsOnCall[i] = struct {
		result1 bool
	}{result1}
}

func (fake *ApplicationCapabilities) V2_0Validation() bool {
	fake.v2_0ValidationMutex.Lock()
	ret, specificReturn := fake.v2_0ValidationReturnsOnCall[len(fake.v2_0ValidationArgsForCall)]
	fake.v2_0ValidationArgsForCall = append(fake.v2_0ValidationArgsForCall, struct {
	}{})
	fake.recordInvocation("V2_0Validation", []interface{}{})
	fake.v2_0ValidationMutex.Unlock()
	if fake.V2_0ValidationStub != nil {
		return fake.V2_0ValidationStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.v2_0ValidationReturns
	return fakeReturns.result1
}

func (fake *ApplicationCapabilities) V2_0ValidationCallCount() int {
	fake.v2_0ValidationMutex.RLock()
	defer fake.v2_0ValidationMutex.RUnlock()
	return len(fake.v2_0ValidationArgsForCall)
}

func (fake *ApplicationCapabilities) V2_0ValidationCalls(stub func() bool) {
	fake.v2_0ValidationMutex.Lock()
	defer fake.v2_0ValidationMutex.Unlock()
	fake.V2_0ValidationStub = stub
}

func (fake *ApplicationCapabilities) V2_0ValidationReturns(result1 bool) {
	fake.v2_0ValidationMutex.Lock()
	defer fake.v2_0ValidationMutex.Unlock()
	fake.V2_0ValidationStub = nil
	fake.v2_0ValidationReturns = struct {
		result1 bool
	}{result1}
}

func (fake *ApplicationCapabilities) V2_0ValidationReturnsOnCall(i int, result1 bool) {
	fake.v2_0ValidationMutex.Lock()
	defer fake.v2_0ValidationMutex.Unlock()
	fake.V2_0ValidationStub = nil
	if fake.v2_0ValidationReturnsOnCall == nil {
		fake.v2_0ValidationReturnsOnCall = make(map[int]struct {
			result1 bool
		})
	}
	fake.v2_0ValidationReturnsOnCall[i] = struct {
		result1 bool
	}{result1}
}

func (fake *ApplicationCapabilities) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.aCLsMutex.RLock()
	defer fake.aCLsMutex.RUnlock()
	fake.collectionUpgradeMutex.RLock()
	defer fake.collectionUpgradeMutex.RUnlock()
	fake.forbidDuplicateTXIdInBlockMutex.RLock()
	defer fake.forbidDuplicateTXIdInBlockMutex.RUnlock()
	fake.keyLevelEndorsementMutex.RLock()
	defer fake.keyLevelEndorsementMutex.RUnlock()
	fake.lifecycleV20Mutex.RLock()
	defer fake.lifecycleV20Mutex.RUnlock()
	fake.metadataLifecycleMutex.RLock()
	defer fake.metadataLifecycleMutex.RUnlock()
	fake.privateChannelDataMutex.RLock()
	defer fake.privateChannelDataMutex.RUnlock()
	fake.storePvtDataOfInvalidTxMutex.RLock()
	defer fake.storePvtDataOfInvalidTxMutex.RUnlock()
	fake.supportedMutex.RLock()
	defer fake.supportedMutex.RUnlock()
	fake.v1_1ValidationMutex.RLock()
	defer fake.v1_1ValidationMutex.RUnlock()
	fake.v1_2ValidationMutex.RLock()
	defer fake.v1_2ValidationMutex.RUnlock()
	fake.v1_3ValidationMutex.RLock()
	defer fake.v1_3ValidationMutex.RUnlock()
	fake.v2_0ValidationMutex.RLock()
	defer fake.v2_0ValidationMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *ApplicationCapabilities) recordInvocation(key string, args []interface{}) {
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
