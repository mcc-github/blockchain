
package mocks

import (
	"sync"
)

type OrdererCapabilities struct {
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
	ExpirationCheckStub        func() bool
	expirationCheckMutex       sync.RWMutex
	expirationCheckArgsForCall []struct {
	}
	expirationCheckReturns struct {
		result1 bool
	}
	expirationCheckReturnsOnCall map[int]struct {
		result1 bool
	}
	PredictableChannelTemplateStub        func() bool
	predictableChannelTemplateMutex       sync.RWMutex
	predictableChannelTemplateArgsForCall []struct {
	}
	predictableChannelTemplateReturns struct {
		result1 bool
	}
	predictableChannelTemplateReturnsOnCall map[int]struct {
		result1 bool
	}
	ResubmissionStub        func() bool
	resubmissionMutex       sync.RWMutex
	resubmissionArgsForCall []struct {
	}
	resubmissionReturns struct {
		result1 bool
	}
	resubmissionReturnsOnCall map[int]struct {
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
	UseChannelCreationPolicyAsAdminsStub        func() bool
	useChannelCreationPolicyAsAdminsMutex       sync.RWMutex
	useChannelCreationPolicyAsAdminsArgsForCall []struct {
	}
	useChannelCreationPolicyAsAdminsReturns struct {
		result1 bool
	}
	useChannelCreationPolicyAsAdminsReturnsOnCall map[int]struct {
		result1 bool
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *OrdererCapabilities) ConsensusTypeMigration() bool {
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

func (fake *OrdererCapabilities) ConsensusTypeMigrationCallCount() int {
	fake.consensusTypeMigrationMutex.RLock()
	defer fake.consensusTypeMigrationMutex.RUnlock()
	return len(fake.consensusTypeMigrationArgsForCall)
}

func (fake *OrdererCapabilities) ConsensusTypeMigrationCalls(stub func() bool) {
	fake.consensusTypeMigrationMutex.Lock()
	defer fake.consensusTypeMigrationMutex.Unlock()
	fake.ConsensusTypeMigrationStub = stub
}

func (fake *OrdererCapabilities) ConsensusTypeMigrationReturns(result1 bool) {
	fake.consensusTypeMigrationMutex.Lock()
	defer fake.consensusTypeMigrationMutex.Unlock()
	fake.ConsensusTypeMigrationStub = nil
	fake.consensusTypeMigrationReturns = struct {
		result1 bool
	}{result1}
}

func (fake *OrdererCapabilities) ConsensusTypeMigrationReturnsOnCall(i int, result1 bool) {
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

func (fake *OrdererCapabilities) ExpirationCheck() bool {
	fake.expirationCheckMutex.Lock()
	ret, specificReturn := fake.expirationCheckReturnsOnCall[len(fake.expirationCheckArgsForCall)]
	fake.expirationCheckArgsForCall = append(fake.expirationCheckArgsForCall, struct {
	}{})
	fake.recordInvocation("ExpirationCheck", []interface{}{})
	fake.expirationCheckMutex.Unlock()
	if fake.ExpirationCheckStub != nil {
		return fake.ExpirationCheckStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.expirationCheckReturns
	return fakeReturns.result1
}

func (fake *OrdererCapabilities) ExpirationCheckCallCount() int {
	fake.expirationCheckMutex.RLock()
	defer fake.expirationCheckMutex.RUnlock()
	return len(fake.expirationCheckArgsForCall)
}

func (fake *OrdererCapabilities) ExpirationCheckCalls(stub func() bool) {
	fake.expirationCheckMutex.Lock()
	defer fake.expirationCheckMutex.Unlock()
	fake.ExpirationCheckStub = stub
}

func (fake *OrdererCapabilities) ExpirationCheckReturns(result1 bool) {
	fake.expirationCheckMutex.Lock()
	defer fake.expirationCheckMutex.Unlock()
	fake.ExpirationCheckStub = nil
	fake.expirationCheckReturns = struct {
		result1 bool
	}{result1}
}

func (fake *OrdererCapabilities) ExpirationCheckReturnsOnCall(i int, result1 bool) {
	fake.expirationCheckMutex.Lock()
	defer fake.expirationCheckMutex.Unlock()
	fake.ExpirationCheckStub = nil
	if fake.expirationCheckReturnsOnCall == nil {
		fake.expirationCheckReturnsOnCall = make(map[int]struct {
			result1 bool
		})
	}
	fake.expirationCheckReturnsOnCall[i] = struct {
		result1 bool
	}{result1}
}

func (fake *OrdererCapabilities) PredictableChannelTemplate() bool {
	fake.predictableChannelTemplateMutex.Lock()
	ret, specificReturn := fake.predictableChannelTemplateReturnsOnCall[len(fake.predictableChannelTemplateArgsForCall)]
	fake.predictableChannelTemplateArgsForCall = append(fake.predictableChannelTemplateArgsForCall, struct {
	}{})
	fake.recordInvocation("PredictableChannelTemplate", []interface{}{})
	fake.predictableChannelTemplateMutex.Unlock()
	if fake.PredictableChannelTemplateStub != nil {
		return fake.PredictableChannelTemplateStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.predictableChannelTemplateReturns
	return fakeReturns.result1
}

func (fake *OrdererCapabilities) PredictableChannelTemplateCallCount() int {
	fake.predictableChannelTemplateMutex.RLock()
	defer fake.predictableChannelTemplateMutex.RUnlock()
	return len(fake.predictableChannelTemplateArgsForCall)
}

func (fake *OrdererCapabilities) PredictableChannelTemplateCalls(stub func() bool) {
	fake.predictableChannelTemplateMutex.Lock()
	defer fake.predictableChannelTemplateMutex.Unlock()
	fake.PredictableChannelTemplateStub = stub
}

func (fake *OrdererCapabilities) PredictableChannelTemplateReturns(result1 bool) {
	fake.predictableChannelTemplateMutex.Lock()
	defer fake.predictableChannelTemplateMutex.Unlock()
	fake.PredictableChannelTemplateStub = nil
	fake.predictableChannelTemplateReturns = struct {
		result1 bool
	}{result1}
}

func (fake *OrdererCapabilities) PredictableChannelTemplateReturnsOnCall(i int, result1 bool) {
	fake.predictableChannelTemplateMutex.Lock()
	defer fake.predictableChannelTemplateMutex.Unlock()
	fake.PredictableChannelTemplateStub = nil
	if fake.predictableChannelTemplateReturnsOnCall == nil {
		fake.predictableChannelTemplateReturnsOnCall = make(map[int]struct {
			result1 bool
		})
	}
	fake.predictableChannelTemplateReturnsOnCall[i] = struct {
		result1 bool
	}{result1}
}

func (fake *OrdererCapabilities) Resubmission() bool {
	fake.resubmissionMutex.Lock()
	ret, specificReturn := fake.resubmissionReturnsOnCall[len(fake.resubmissionArgsForCall)]
	fake.resubmissionArgsForCall = append(fake.resubmissionArgsForCall, struct {
	}{})
	fake.recordInvocation("Resubmission", []interface{}{})
	fake.resubmissionMutex.Unlock()
	if fake.ResubmissionStub != nil {
		return fake.ResubmissionStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.resubmissionReturns
	return fakeReturns.result1
}

func (fake *OrdererCapabilities) ResubmissionCallCount() int {
	fake.resubmissionMutex.RLock()
	defer fake.resubmissionMutex.RUnlock()
	return len(fake.resubmissionArgsForCall)
}

func (fake *OrdererCapabilities) ResubmissionCalls(stub func() bool) {
	fake.resubmissionMutex.Lock()
	defer fake.resubmissionMutex.Unlock()
	fake.ResubmissionStub = stub
}

func (fake *OrdererCapabilities) ResubmissionReturns(result1 bool) {
	fake.resubmissionMutex.Lock()
	defer fake.resubmissionMutex.Unlock()
	fake.ResubmissionStub = nil
	fake.resubmissionReturns = struct {
		result1 bool
	}{result1}
}

func (fake *OrdererCapabilities) ResubmissionReturnsOnCall(i int, result1 bool) {
	fake.resubmissionMutex.Lock()
	defer fake.resubmissionMutex.Unlock()
	fake.ResubmissionStub = nil
	if fake.resubmissionReturnsOnCall == nil {
		fake.resubmissionReturnsOnCall = make(map[int]struct {
			result1 bool
		})
	}
	fake.resubmissionReturnsOnCall[i] = struct {
		result1 bool
	}{result1}
}

func (fake *OrdererCapabilities) Supported() error {
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

func (fake *OrdererCapabilities) SupportedCallCount() int {
	fake.supportedMutex.RLock()
	defer fake.supportedMutex.RUnlock()
	return len(fake.supportedArgsForCall)
}

func (fake *OrdererCapabilities) SupportedCalls(stub func() error) {
	fake.supportedMutex.Lock()
	defer fake.supportedMutex.Unlock()
	fake.SupportedStub = stub
}

func (fake *OrdererCapabilities) SupportedReturns(result1 error) {
	fake.supportedMutex.Lock()
	defer fake.supportedMutex.Unlock()
	fake.SupportedStub = nil
	fake.supportedReturns = struct {
		result1 error
	}{result1}
}

func (fake *OrdererCapabilities) SupportedReturnsOnCall(i int, result1 error) {
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

func (fake *OrdererCapabilities) UseChannelCreationPolicyAsAdmins() bool {
	fake.useChannelCreationPolicyAsAdminsMutex.Lock()
	ret, specificReturn := fake.useChannelCreationPolicyAsAdminsReturnsOnCall[len(fake.useChannelCreationPolicyAsAdminsArgsForCall)]
	fake.useChannelCreationPolicyAsAdminsArgsForCall = append(fake.useChannelCreationPolicyAsAdminsArgsForCall, struct {
	}{})
	fake.recordInvocation("UseChannelCreationPolicyAsAdmins", []interface{}{})
	fake.useChannelCreationPolicyAsAdminsMutex.Unlock()
	if fake.UseChannelCreationPolicyAsAdminsStub != nil {
		return fake.UseChannelCreationPolicyAsAdminsStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.useChannelCreationPolicyAsAdminsReturns
	return fakeReturns.result1
}

func (fake *OrdererCapabilities) UseChannelCreationPolicyAsAdminsCallCount() int {
	fake.useChannelCreationPolicyAsAdminsMutex.RLock()
	defer fake.useChannelCreationPolicyAsAdminsMutex.RUnlock()
	return len(fake.useChannelCreationPolicyAsAdminsArgsForCall)
}

func (fake *OrdererCapabilities) UseChannelCreationPolicyAsAdminsCalls(stub func() bool) {
	fake.useChannelCreationPolicyAsAdminsMutex.Lock()
	defer fake.useChannelCreationPolicyAsAdminsMutex.Unlock()
	fake.UseChannelCreationPolicyAsAdminsStub = stub
}

func (fake *OrdererCapabilities) UseChannelCreationPolicyAsAdminsReturns(result1 bool) {
	fake.useChannelCreationPolicyAsAdminsMutex.Lock()
	defer fake.useChannelCreationPolicyAsAdminsMutex.Unlock()
	fake.UseChannelCreationPolicyAsAdminsStub = nil
	fake.useChannelCreationPolicyAsAdminsReturns = struct {
		result1 bool
	}{result1}
}

func (fake *OrdererCapabilities) UseChannelCreationPolicyAsAdminsReturnsOnCall(i int, result1 bool) {
	fake.useChannelCreationPolicyAsAdminsMutex.Lock()
	defer fake.useChannelCreationPolicyAsAdminsMutex.Unlock()
	fake.UseChannelCreationPolicyAsAdminsStub = nil
	if fake.useChannelCreationPolicyAsAdminsReturnsOnCall == nil {
		fake.useChannelCreationPolicyAsAdminsReturnsOnCall = make(map[int]struct {
			result1 bool
		})
	}
	fake.useChannelCreationPolicyAsAdminsReturnsOnCall[i] = struct {
		result1 bool
	}{result1}
}

func (fake *OrdererCapabilities) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.consensusTypeMigrationMutex.RLock()
	defer fake.consensusTypeMigrationMutex.RUnlock()
	fake.expirationCheckMutex.RLock()
	defer fake.expirationCheckMutex.RUnlock()
	fake.predictableChannelTemplateMutex.RLock()
	defer fake.predictableChannelTemplateMutex.RUnlock()
	fake.resubmissionMutex.RLock()
	defer fake.resubmissionMutex.RUnlock()
	fake.supportedMutex.RLock()
	defer fake.supportedMutex.RUnlock()
	fake.useChannelCreationPolicyAsAdminsMutex.RLock()
	defer fake.useChannelCreationPolicyAsAdminsMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *OrdererCapabilities) recordInvocation(key string, args []interface{}) {
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
