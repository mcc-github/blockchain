
package mock

import (
	"sync"

	"github.com/mcc-github/blockchain/common/channelconfig"
	"github.com/mcc-github/blockchain/common/configtx"
	"github.com/mcc-github/blockchain/common/policies"
	"github.com/mcc-github/blockchain/msp"
)

type ChannelConfig struct {
	ApplicationConfigStub        func() (channelconfig.Application, bool)
	applicationConfigMutex       sync.RWMutex
	applicationConfigArgsForCall []struct {
	}
	applicationConfigReturns struct {
		result1 channelconfig.Application
		result2 bool
	}
	applicationConfigReturnsOnCall map[int]struct {
		result1 channelconfig.Application
		result2 bool
	}
	ChannelConfigStub        func() channelconfig.Channel
	channelConfigMutex       sync.RWMutex
	channelConfigArgsForCall []struct {
	}
	channelConfigReturns struct {
		result1 channelconfig.Channel
	}
	channelConfigReturnsOnCall map[int]struct {
		result1 channelconfig.Channel
	}
	ConfigtxValidatorStub        func() configtx.Validator
	configtxValidatorMutex       sync.RWMutex
	configtxValidatorArgsForCall []struct {
	}
	configtxValidatorReturns struct {
		result1 configtx.Validator
	}
	configtxValidatorReturnsOnCall map[int]struct {
		result1 configtx.Validator
	}
	ConsortiumsConfigStub        func() (channelconfig.Consortiums, bool)
	consortiumsConfigMutex       sync.RWMutex
	consortiumsConfigArgsForCall []struct {
	}
	consortiumsConfigReturns struct {
		result1 channelconfig.Consortiums
		result2 bool
	}
	consortiumsConfigReturnsOnCall map[int]struct {
		result1 channelconfig.Consortiums
		result2 bool
	}
	MSPManagerStub        func() msp.MSPManager
	mSPManagerMutex       sync.RWMutex
	mSPManagerArgsForCall []struct {
	}
	mSPManagerReturns struct {
		result1 msp.MSPManager
	}
	mSPManagerReturnsOnCall map[int]struct {
		result1 msp.MSPManager
	}
	OrdererConfigStub        func() (channelconfig.Orderer, bool)
	ordererConfigMutex       sync.RWMutex
	ordererConfigArgsForCall []struct {
	}
	ordererConfigReturns struct {
		result1 channelconfig.Orderer
		result2 bool
	}
	ordererConfigReturnsOnCall map[int]struct {
		result1 channelconfig.Orderer
		result2 bool
	}
	PolicyManagerStub        func() policies.Manager
	policyManagerMutex       sync.RWMutex
	policyManagerArgsForCall []struct {
	}
	policyManagerReturns struct {
		result1 policies.Manager
	}
	policyManagerReturnsOnCall map[int]struct {
		result1 policies.Manager
	}
	ValidateNewStub        func(channelconfig.Resources) error
	validateNewMutex       sync.RWMutex
	validateNewArgsForCall []struct {
		arg1 channelconfig.Resources
	}
	validateNewReturns struct {
		result1 error
	}
	validateNewReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *ChannelConfig) ApplicationConfig() (channelconfig.Application, bool) {
	fake.applicationConfigMutex.Lock()
	ret, specificReturn := fake.applicationConfigReturnsOnCall[len(fake.applicationConfigArgsForCall)]
	fake.applicationConfigArgsForCall = append(fake.applicationConfigArgsForCall, struct {
	}{})
	fake.recordInvocation("ApplicationConfig", []interface{}{})
	fake.applicationConfigMutex.Unlock()
	if fake.ApplicationConfigStub != nil {
		return fake.ApplicationConfigStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.applicationConfigReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *ChannelConfig) ApplicationConfigCallCount() int {
	fake.applicationConfigMutex.RLock()
	defer fake.applicationConfigMutex.RUnlock()
	return len(fake.applicationConfigArgsForCall)
}

func (fake *ChannelConfig) ApplicationConfigCalls(stub func() (channelconfig.Application, bool)) {
	fake.applicationConfigMutex.Lock()
	defer fake.applicationConfigMutex.Unlock()
	fake.ApplicationConfigStub = stub
}

func (fake *ChannelConfig) ApplicationConfigReturns(result1 channelconfig.Application, result2 bool) {
	fake.applicationConfigMutex.Lock()
	defer fake.applicationConfigMutex.Unlock()
	fake.ApplicationConfigStub = nil
	fake.applicationConfigReturns = struct {
		result1 channelconfig.Application
		result2 bool
	}{result1, result2}
}

func (fake *ChannelConfig) ApplicationConfigReturnsOnCall(i int, result1 channelconfig.Application, result2 bool) {
	fake.applicationConfigMutex.Lock()
	defer fake.applicationConfigMutex.Unlock()
	fake.ApplicationConfigStub = nil
	if fake.applicationConfigReturnsOnCall == nil {
		fake.applicationConfigReturnsOnCall = make(map[int]struct {
			result1 channelconfig.Application
			result2 bool
		})
	}
	fake.applicationConfigReturnsOnCall[i] = struct {
		result1 channelconfig.Application
		result2 bool
	}{result1, result2}
}

func (fake *ChannelConfig) ChannelConfig() channelconfig.Channel {
	fake.channelConfigMutex.Lock()
	ret, specificReturn := fake.channelConfigReturnsOnCall[len(fake.channelConfigArgsForCall)]
	fake.channelConfigArgsForCall = append(fake.channelConfigArgsForCall, struct {
	}{})
	fake.recordInvocation("ChannelConfig", []interface{}{})
	fake.channelConfigMutex.Unlock()
	if fake.ChannelConfigStub != nil {
		return fake.ChannelConfigStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.channelConfigReturns
	return fakeReturns.result1
}

func (fake *ChannelConfig) ChannelConfigCallCount() int {
	fake.channelConfigMutex.RLock()
	defer fake.channelConfigMutex.RUnlock()
	return len(fake.channelConfigArgsForCall)
}

func (fake *ChannelConfig) ChannelConfigCalls(stub func() channelconfig.Channel) {
	fake.channelConfigMutex.Lock()
	defer fake.channelConfigMutex.Unlock()
	fake.ChannelConfigStub = stub
}

func (fake *ChannelConfig) ChannelConfigReturns(result1 channelconfig.Channel) {
	fake.channelConfigMutex.Lock()
	defer fake.channelConfigMutex.Unlock()
	fake.ChannelConfigStub = nil
	fake.channelConfigReturns = struct {
		result1 channelconfig.Channel
	}{result1}
}

func (fake *ChannelConfig) ChannelConfigReturnsOnCall(i int, result1 channelconfig.Channel) {
	fake.channelConfigMutex.Lock()
	defer fake.channelConfigMutex.Unlock()
	fake.ChannelConfigStub = nil
	if fake.channelConfigReturnsOnCall == nil {
		fake.channelConfigReturnsOnCall = make(map[int]struct {
			result1 channelconfig.Channel
		})
	}
	fake.channelConfigReturnsOnCall[i] = struct {
		result1 channelconfig.Channel
	}{result1}
}

func (fake *ChannelConfig) ConfigtxValidator() configtx.Validator {
	fake.configtxValidatorMutex.Lock()
	ret, specificReturn := fake.configtxValidatorReturnsOnCall[len(fake.configtxValidatorArgsForCall)]
	fake.configtxValidatorArgsForCall = append(fake.configtxValidatorArgsForCall, struct {
	}{})
	fake.recordInvocation("ConfigtxValidator", []interface{}{})
	fake.configtxValidatorMutex.Unlock()
	if fake.ConfigtxValidatorStub != nil {
		return fake.ConfigtxValidatorStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.configtxValidatorReturns
	return fakeReturns.result1
}

func (fake *ChannelConfig) ConfigtxValidatorCallCount() int {
	fake.configtxValidatorMutex.RLock()
	defer fake.configtxValidatorMutex.RUnlock()
	return len(fake.configtxValidatorArgsForCall)
}

func (fake *ChannelConfig) ConfigtxValidatorCalls(stub func() configtx.Validator) {
	fake.configtxValidatorMutex.Lock()
	defer fake.configtxValidatorMutex.Unlock()
	fake.ConfigtxValidatorStub = stub
}

func (fake *ChannelConfig) ConfigtxValidatorReturns(result1 configtx.Validator) {
	fake.configtxValidatorMutex.Lock()
	defer fake.configtxValidatorMutex.Unlock()
	fake.ConfigtxValidatorStub = nil
	fake.configtxValidatorReturns = struct {
		result1 configtx.Validator
	}{result1}
}

func (fake *ChannelConfig) ConfigtxValidatorReturnsOnCall(i int, result1 configtx.Validator) {
	fake.configtxValidatorMutex.Lock()
	defer fake.configtxValidatorMutex.Unlock()
	fake.ConfigtxValidatorStub = nil
	if fake.configtxValidatorReturnsOnCall == nil {
		fake.configtxValidatorReturnsOnCall = make(map[int]struct {
			result1 configtx.Validator
		})
	}
	fake.configtxValidatorReturnsOnCall[i] = struct {
		result1 configtx.Validator
	}{result1}
}

func (fake *ChannelConfig) ConsortiumsConfig() (channelconfig.Consortiums, bool) {
	fake.consortiumsConfigMutex.Lock()
	ret, specificReturn := fake.consortiumsConfigReturnsOnCall[len(fake.consortiumsConfigArgsForCall)]
	fake.consortiumsConfigArgsForCall = append(fake.consortiumsConfigArgsForCall, struct {
	}{})
	fake.recordInvocation("ConsortiumsConfig", []interface{}{})
	fake.consortiumsConfigMutex.Unlock()
	if fake.ConsortiumsConfigStub != nil {
		return fake.ConsortiumsConfigStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.consortiumsConfigReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *ChannelConfig) ConsortiumsConfigCallCount() int {
	fake.consortiumsConfigMutex.RLock()
	defer fake.consortiumsConfigMutex.RUnlock()
	return len(fake.consortiumsConfigArgsForCall)
}

func (fake *ChannelConfig) ConsortiumsConfigCalls(stub func() (channelconfig.Consortiums, bool)) {
	fake.consortiumsConfigMutex.Lock()
	defer fake.consortiumsConfigMutex.Unlock()
	fake.ConsortiumsConfigStub = stub
}

func (fake *ChannelConfig) ConsortiumsConfigReturns(result1 channelconfig.Consortiums, result2 bool) {
	fake.consortiumsConfigMutex.Lock()
	defer fake.consortiumsConfigMutex.Unlock()
	fake.ConsortiumsConfigStub = nil
	fake.consortiumsConfigReturns = struct {
		result1 channelconfig.Consortiums
		result2 bool
	}{result1, result2}
}

func (fake *ChannelConfig) ConsortiumsConfigReturnsOnCall(i int, result1 channelconfig.Consortiums, result2 bool) {
	fake.consortiumsConfigMutex.Lock()
	defer fake.consortiumsConfigMutex.Unlock()
	fake.ConsortiumsConfigStub = nil
	if fake.consortiumsConfigReturnsOnCall == nil {
		fake.consortiumsConfigReturnsOnCall = make(map[int]struct {
			result1 channelconfig.Consortiums
			result2 bool
		})
	}
	fake.consortiumsConfigReturnsOnCall[i] = struct {
		result1 channelconfig.Consortiums
		result2 bool
	}{result1, result2}
}

func (fake *ChannelConfig) MSPManager() msp.MSPManager {
	fake.mSPManagerMutex.Lock()
	ret, specificReturn := fake.mSPManagerReturnsOnCall[len(fake.mSPManagerArgsForCall)]
	fake.mSPManagerArgsForCall = append(fake.mSPManagerArgsForCall, struct {
	}{})
	fake.recordInvocation("MSPManager", []interface{}{})
	fake.mSPManagerMutex.Unlock()
	if fake.MSPManagerStub != nil {
		return fake.MSPManagerStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.mSPManagerReturns
	return fakeReturns.result1
}

func (fake *ChannelConfig) MSPManagerCallCount() int {
	fake.mSPManagerMutex.RLock()
	defer fake.mSPManagerMutex.RUnlock()
	return len(fake.mSPManagerArgsForCall)
}

func (fake *ChannelConfig) MSPManagerCalls(stub func() msp.MSPManager) {
	fake.mSPManagerMutex.Lock()
	defer fake.mSPManagerMutex.Unlock()
	fake.MSPManagerStub = stub
}

func (fake *ChannelConfig) MSPManagerReturns(result1 msp.MSPManager) {
	fake.mSPManagerMutex.Lock()
	defer fake.mSPManagerMutex.Unlock()
	fake.MSPManagerStub = nil
	fake.mSPManagerReturns = struct {
		result1 msp.MSPManager
	}{result1}
}

func (fake *ChannelConfig) MSPManagerReturnsOnCall(i int, result1 msp.MSPManager) {
	fake.mSPManagerMutex.Lock()
	defer fake.mSPManagerMutex.Unlock()
	fake.MSPManagerStub = nil
	if fake.mSPManagerReturnsOnCall == nil {
		fake.mSPManagerReturnsOnCall = make(map[int]struct {
			result1 msp.MSPManager
		})
	}
	fake.mSPManagerReturnsOnCall[i] = struct {
		result1 msp.MSPManager
	}{result1}
}

func (fake *ChannelConfig) OrdererConfig() (channelconfig.Orderer, bool) {
	fake.ordererConfigMutex.Lock()
	ret, specificReturn := fake.ordererConfigReturnsOnCall[len(fake.ordererConfigArgsForCall)]
	fake.ordererConfigArgsForCall = append(fake.ordererConfigArgsForCall, struct {
	}{})
	fake.recordInvocation("OrdererConfig", []interface{}{})
	fake.ordererConfigMutex.Unlock()
	if fake.OrdererConfigStub != nil {
		return fake.OrdererConfigStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.ordererConfigReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *ChannelConfig) OrdererConfigCallCount() int {
	fake.ordererConfigMutex.RLock()
	defer fake.ordererConfigMutex.RUnlock()
	return len(fake.ordererConfigArgsForCall)
}

func (fake *ChannelConfig) OrdererConfigCalls(stub func() (channelconfig.Orderer, bool)) {
	fake.ordererConfigMutex.Lock()
	defer fake.ordererConfigMutex.Unlock()
	fake.OrdererConfigStub = stub
}

func (fake *ChannelConfig) OrdererConfigReturns(result1 channelconfig.Orderer, result2 bool) {
	fake.ordererConfigMutex.Lock()
	defer fake.ordererConfigMutex.Unlock()
	fake.OrdererConfigStub = nil
	fake.ordererConfigReturns = struct {
		result1 channelconfig.Orderer
		result2 bool
	}{result1, result2}
}

func (fake *ChannelConfig) OrdererConfigReturnsOnCall(i int, result1 channelconfig.Orderer, result2 bool) {
	fake.ordererConfigMutex.Lock()
	defer fake.ordererConfigMutex.Unlock()
	fake.OrdererConfigStub = nil
	if fake.ordererConfigReturnsOnCall == nil {
		fake.ordererConfigReturnsOnCall = make(map[int]struct {
			result1 channelconfig.Orderer
			result2 bool
		})
	}
	fake.ordererConfigReturnsOnCall[i] = struct {
		result1 channelconfig.Orderer
		result2 bool
	}{result1, result2}
}

func (fake *ChannelConfig) PolicyManager() policies.Manager {
	fake.policyManagerMutex.Lock()
	ret, specificReturn := fake.policyManagerReturnsOnCall[len(fake.policyManagerArgsForCall)]
	fake.policyManagerArgsForCall = append(fake.policyManagerArgsForCall, struct {
	}{})
	fake.recordInvocation("PolicyManager", []interface{}{})
	fake.policyManagerMutex.Unlock()
	if fake.PolicyManagerStub != nil {
		return fake.PolicyManagerStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.policyManagerReturns
	return fakeReturns.result1
}

func (fake *ChannelConfig) PolicyManagerCallCount() int {
	fake.policyManagerMutex.RLock()
	defer fake.policyManagerMutex.RUnlock()
	return len(fake.policyManagerArgsForCall)
}

func (fake *ChannelConfig) PolicyManagerCalls(stub func() policies.Manager) {
	fake.policyManagerMutex.Lock()
	defer fake.policyManagerMutex.Unlock()
	fake.PolicyManagerStub = stub
}

func (fake *ChannelConfig) PolicyManagerReturns(result1 policies.Manager) {
	fake.policyManagerMutex.Lock()
	defer fake.policyManagerMutex.Unlock()
	fake.PolicyManagerStub = nil
	fake.policyManagerReturns = struct {
		result1 policies.Manager
	}{result1}
}

func (fake *ChannelConfig) PolicyManagerReturnsOnCall(i int, result1 policies.Manager) {
	fake.policyManagerMutex.Lock()
	defer fake.policyManagerMutex.Unlock()
	fake.PolicyManagerStub = nil
	if fake.policyManagerReturnsOnCall == nil {
		fake.policyManagerReturnsOnCall = make(map[int]struct {
			result1 policies.Manager
		})
	}
	fake.policyManagerReturnsOnCall[i] = struct {
		result1 policies.Manager
	}{result1}
}

func (fake *ChannelConfig) ValidateNew(arg1 channelconfig.Resources) error {
	fake.validateNewMutex.Lock()
	ret, specificReturn := fake.validateNewReturnsOnCall[len(fake.validateNewArgsForCall)]
	fake.validateNewArgsForCall = append(fake.validateNewArgsForCall, struct {
		arg1 channelconfig.Resources
	}{arg1})
	fake.recordInvocation("ValidateNew", []interface{}{arg1})
	fake.validateNewMutex.Unlock()
	if fake.ValidateNewStub != nil {
		return fake.ValidateNewStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.validateNewReturns
	return fakeReturns.result1
}

func (fake *ChannelConfig) ValidateNewCallCount() int {
	fake.validateNewMutex.RLock()
	defer fake.validateNewMutex.RUnlock()
	return len(fake.validateNewArgsForCall)
}

func (fake *ChannelConfig) ValidateNewCalls(stub func(channelconfig.Resources) error) {
	fake.validateNewMutex.Lock()
	defer fake.validateNewMutex.Unlock()
	fake.ValidateNewStub = stub
}

func (fake *ChannelConfig) ValidateNewArgsForCall(i int) channelconfig.Resources {
	fake.validateNewMutex.RLock()
	defer fake.validateNewMutex.RUnlock()
	argsForCall := fake.validateNewArgsForCall[i]
	return argsForCall.arg1
}

func (fake *ChannelConfig) ValidateNewReturns(result1 error) {
	fake.validateNewMutex.Lock()
	defer fake.validateNewMutex.Unlock()
	fake.ValidateNewStub = nil
	fake.validateNewReturns = struct {
		result1 error
	}{result1}
}

func (fake *ChannelConfig) ValidateNewReturnsOnCall(i int, result1 error) {
	fake.validateNewMutex.Lock()
	defer fake.validateNewMutex.Unlock()
	fake.ValidateNewStub = nil
	if fake.validateNewReturnsOnCall == nil {
		fake.validateNewReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.validateNewReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *ChannelConfig) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.applicationConfigMutex.RLock()
	defer fake.applicationConfigMutex.RUnlock()
	fake.channelConfigMutex.RLock()
	defer fake.channelConfigMutex.RUnlock()
	fake.configtxValidatorMutex.RLock()
	defer fake.configtxValidatorMutex.RUnlock()
	fake.consortiumsConfigMutex.RLock()
	defer fake.consortiumsConfigMutex.RUnlock()
	fake.mSPManagerMutex.RLock()
	defer fake.mSPManagerMutex.RUnlock()
	fake.ordererConfigMutex.RLock()
	defer fake.ordererConfigMutex.RUnlock()
	fake.policyManagerMutex.RLock()
	defer fake.policyManagerMutex.RUnlock()
	fake.validateNewMutex.RLock()
	defer fake.validateNewMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *ChannelConfig) recordInvocation(key string, args []interface{}) {
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