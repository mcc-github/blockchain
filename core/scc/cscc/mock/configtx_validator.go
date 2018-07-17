
package mock

import (
	"sync"

	cb "github.com/mcc-github/blockchain/protos/common"
)

type ConfigtxValidator struct {
	ValidateStub        func(configEnv *cb.ConfigEnvelope) error
	validateMutex       sync.RWMutex
	validateArgsForCall []struct {
		configEnv *cb.ConfigEnvelope
	}
	validateReturns struct {
		result1 error
	}
	validateReturnsOnCall map[int]struct {
		result1 error
	}
	ProposeConfigUpdateStub        func(configtx *cb.Envelope) (*cb.ConfigEnvelope, error)
	proposeConfigUpdateMutex       sync.RWMutex
	proposeConfigUpdateArgsForCall []struct {
		configtx *cb.Envelope
	}
	proposeConfigUpdateReturns struct {
		result1 *cb.ConfigEnvelope
		result2 error
	}
	proposeConfigUpdateReturnsOnCall map[int]struct {
		result1 *cb.ConfigEnvelope
		result2 error
	}
	ChainIDStub        func() string
	chainIDMutex       sync.RWMutex
	chainIDArgsForCall []struct{}
	chainIDReturns     struct {
		result1 string
	}
	chainIDReturnsOnCall map[int]struct {
		result1 string
	}
	ConfigProtoStub        func() *cb.Config
	configProtoMutex       sync.RWMutex
	configProtoArgsForCall []struct{}
	configProtoReturns     struct {
		result1 *cb.Config
	}
	configProtoReturnsOnCall map[int]struct {
		result1 *cb.Config
	}
	SequenceStub        func() uint64
	sequenceMutex       sync.RWMutex
	sequenceArgsForCall []struct{}
	sequenceReturns     struct {
		result1 uint64
	}
	sequenceReturnsOnCall map[int]struct {
		result1 uint64
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *ConfigtxValidator) Validate(configEnv *cb.ConfigEnvelope) error {
	fake.validateMutex.Lock()
	ret, specificReturn := fake.validateReturnsOnCall[len(fake.validateArgsForCall)]
	fake.validateArgsForCall = append(fake.validateArgsForCall, struct {
		configEnv *cb.ConfigEnvelope
	}{configEnv})
	fake.recordInvocation("Validate", []interface{}{configEnv})
	fake.validateMutex.Unlock()
	if fake.ValidateStub != nil {
		return fake.ValidateStub(configEnv)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.validateReturns.result1
}

func (fake *ConfigtxValidator) ValidateCallCount() int {
	fake.validateMutex.RLock()
	defer fake.validateMutex.RUnlock()
	return len(fake.validateArgsForCall)
}

func (fake *ConfigtxValidator) ValidateArgsForCall(i int) *cb.ConfigEnvelope {
	fake.validateMutex.RLock()
	defer fake.validateMutex.RUnlock()
	return fake.validateArgsForCall[i].configEnv
}

func (fake *ConfigtxValidator) ValidateReturns(result1 error) {
	fake.ValidateStub = nil
	fake.validateReturns = struct {
		result1 error
	}{result1}
}

func (fake *ConfigtxValidator) ValidateReturnsOnCall(i int, result1 error) {
	fake.ValidateStub = nil
	if fake.validateReturnsOnCall == nil {
		fake.validateReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.validateReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *ConfigtxValidator) ProposeConfigUpdate(configtx *cb.Envelope) (*cb.ConfigEnvelope, error) {
	fake.proposeConfigUpdateMutex.Lock()
	ret, specificReturn := fake.proposeConfigUpdateReturnsOnCall[len(fake.proposeConfigUpdateArgsForCall)]
	fake.proposeConfigUpdateArgsForCall = append(fake.proposeConfigUpdateArgsForCall, struct {
		configtx *cb.Envelope
	}{configtx})
	fake.recordInvocation("ProposeConfigUpdate", []interface{}{configtx})
	fake.proposeConfigUpdateMutex.Unlock()
	if fake.ProposeConfigUpdateStub != nil {
		return fake.ProposeConfigUpdateStub(configtx)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.proposeConfigUpdateReturns.result1, fake.proposeConfigUpdateReturns.result2
}

func (fake *ConfigtxValidator) ProposeConfigUpdateCallCount() int {
	fake.proposeConfigUpdateMutex.RLock()
	defer fake.proposeConfigUpdateMutex.RUnlock()
	return len(fake.proposeConfigUpdateArgsForCall)
}

func (fake *ConfigtxValidator) ProposeConfigUpdateArgsForCall(i int) *cb.Envelope {
	fake.proposeConfigUpdateMutex.RLock()
	defer fake.proposeConfigUpdateMutex.RUnlock()
	return fake.proposeConfigUpdateArgsForCall[i].configtx
}

func (fake *ConfigtxValidator) ProposeConfigUpdateReturns(result1 *cb.ConfigEnvelope, result2 error) {
	fake.ProposeConfigUpdateStub = nil
	fake.proposeConfigUpdateReturns = struct {
		result1 *cb.ConfigEnvelope
		result2 error
	}{result1, result2}
}

func (fake *ConfigtxValidator) ProposeConfigUpdateReturnsOnCall(i int, result1 *cb.ConfigEnvelope, result2 error) {
	fake.ProposeConfigUpdateStub = nil
	if fake.proposeConfigUpdateReturnsOnCall == nil {
		fake.proposeConfigUpdateReturnsOnCall = make(map[int]struct {
			result1 *cb.ConfigEnvelope
			result2 error
		})
	}
	fake.proposeConfigUpdateReturnsOnCall[i] = struct {
		result1 *cb.ConfigEnvelope
		result2 error
	}{result1, result2}
}

func (fake *ConfigtxValidator) ChainID() string {
	fake.chainIDMutex.Lock()
	ret, specificReturn := fake.chainIDReturnsOnCall[len(fake.chainIDArgsForCall)]
	fake.chainIDArgsForCall = append(fake.chainIDArgsForCall, struct{}{})
	fake.recordInvocation("ChainID", []interface{}{})
	fake.chainIDMutex.Unlock()
	if fake.ChainIDStub != nil {
		return fake.ChainIDStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.chainIDReturns.result1
}

func (fake *ConfigtxValidator) ChainIDCallCount() int {
	fake.chainIDMutex.RLock()
	defer fake.chainIDMutex.RUnlock()
	return len(fake.chainIDArgsForCall)
}

func (fake *ConfigtxValidator) ChainIDReturns(result1 string) {
	fake.ChainIDStub = nil
	fake.chainIDReturns = struct {
		result1 string
	}{result1}
}

func (fake *ConfigtxValidator) ChainIDReturnsOnCall(i int, result1 string) {
	fake.ChainIDStub = nil
	if fake.chainIDReturnsOnCall == nil {
		fake.chainIDReturnsOnCall = make(map[int]struct {
			result1 string
		})
	}
	fake.chainIDReturnsOnCall[i] = struct {
		result1 string
	}{result1}
}

func (fake *ConfigtxValidator) ConfigProto() *cb.Config {
	fake.configProtoMutex.Lock()
	ret, specificReturn := fake.configProtoReturnsOnCall[len(fake.configProtoArgsForCall)]
	fake.configProtoArgsForCall = append(fake.configProtoArgsForCall, struct{}{})
	fake.recordInvocation("ConfigProto", []interface{}{})
	fake.configProtoMutex.Unlock()
	if fake.ConfigProtoStub != nil {
		return fake.ConfigProtoStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.configProtoReturns.result1
}

func (fake *ConfigtxValidator) ConfigProtoCallCount() int {
	fake.configProtoMutex.RLock()
	defer fake.configProtoMutex.RUnlock()
	return len(fake.configProtoArgsForCall)
}

func (fake *ConfigtxValidator) ConfigProtoReturns(result1 *cb.Config) {
	fake.ConfigProtoStub = nil
	fake.configProtoReturns = struct {
		result1 *cb.Config
	}{result1}
}

func (fake *ConfigtxValidator) ConfigProtoReturnsOnCall(i int, result1 *cb.Config) {
	fake.ConfigProtoStub = nil
	if fake.configProtoReturnsOnCall == nil {
		fake.configProtoReturnsOnCall = make(map[int]struct {
			result1 *cb.Config
		})
	}
	fake.configProtoReturnsOnCall[i] = struct {
		result1 *cb.Config
	}{result1}
}

func (fake *ConfigtxValidator) Sequence() uint64 {
	fake.sequenceMutex.Lock()
	ret, specificReturn := fake.sequenceReturnsOnCall[len(fake.sequenceArgsForCall)]
	fake.sequenceArgsForCall = append(fake.sequenceArgsForCall, struct{}{})
	fake.recordInvocation("Sequence", []interface{}{})
	fake.sequenceMutex.Unlock()
	if fake.SequenceStub != nil {
		return fake.SequenceStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.sequenceReturns.result1
}

func (fake *ConfigtxValidator) SequenceCallCount() int {
	fake.sequenceMutex.RLock()
	defer fake.sequenceMutex.RUnlock()
	return len(fake.sequenceArgsForCall)
}

func (fake *ConfigtxValidator) SequenceReturns(result1 uint64) {
	fake.SequenceStub = nil
	fake.sequenceReturns = struct {
		result1 uint64
	}{result1}
}

func (fake *ConfigtxValidator) SequenceReturnsOnCall(i int, result1 uint64) {
	fake.SequenceStub = nil
	if fake.sequenceReturnsOnCall == nil {
		fake.sequenceReturnsOnCall = make(map[int]struct {
			result1 uint64
		})
	}
	fake.sequenceReturnsOnCall[i] = struct {
		result1 uint64
	}{result1}
}

func (fake *ConfigtxValidator) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.validateMutex.RLock()
	defer fake.validateMutex.RUnlock()
	fake.proposeConfigUpdateMutex.RLock()
	defer fake.proposeConfigUpdateMutex.RUnlock()
	fake.chainIDMutex.RLock()
	defer fake.chainIDMutex.RUnlock()
	fake.configProtoMutex.RLock()
	defer fake.configProtoMutex.RUnlock()
	fake.sequenceMutex.RLock()
	defer fake.sequenceMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *ConfigtxValidator) recordInvocation(key string, args []interface{}) {
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
