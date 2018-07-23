
package mock

import (
	"sync"

	"github.com/mcc-github/blockchain/common/config"
)

type ConfigManager struct {
	GetChannelConfigStub        func(channel string) config.Config
	getChannelConfigMutex       sync.RWMutex
	getChannelConfigArgsForCall []struct {
		channel string
	}
	getChannelConfigReturns struct {
		result1 config.Config
	}
	getChannelConfigReturnsOnCall map[int]struct {
		result1 config.Config
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *ConfigManager) GetChannelConfig(channel string) config.Config {
	fake.getChannelConfigMutex.Lock()
	ret, specificReturn := fake.getChannelConfigReturnsOnCall[len(fake.getChannelConfigArgsForCall)]
	fake.getChannelConfigArgsForCall = append(fake.getChannelConfigArgsForCall, struct {
		channel string
	}{channel})
	fake.recordInvocation("GetChannelConfig", []interface{}{channel})
	fake.getChannelConfigMutex.Unlock()
	if fake.GetChannelConfigStub != nil {
		return fake.GetChannelConfigStub(channel)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.getChannelConfigReturns.result1
}

func (fake *ConfigManager) GetChannelConfigCallCount() int {
	fake.getChannelConfigMutex.RLock()
	defer fake.getChannelConfigMutex.RUnlock()
	return len(fake.getChannelConfigArgsForCall)
}

func (fake *ConfigManager) GetChannelConfigArgsForCall(i int) string {
	fake.getChannelConfigMutex.RLock()
	defer fake.getChannelConfigMutex.RUnlock()
	return fake.getChannelConfigArgsForCall[i].channel
}

func (fake *ConfigManager) GetChannelConfigReturns(result1 config.Config) {
	fake.GetChannelConfigStub = nil
	fake.getChannelConfigReturns = struct {
		result1 config.Config
	}{result1}
}

func (fake *ConfigManager) GetChannelConfigReturnsOnCall(i int, result1 config.Config) {
	fake.GetChannelConfigStub = nil
	if fake.getChannelConfigReturnsOnCall == nil {
		fake.getChannelConfigReturnsOnCall = make(map[int]struct {
			result1 config.Config
		})
	}
	fake.getChannelConfigReturnsOnCall[i] = struct {
		result1 config.Config
	}{result1}
}

func (fake *ConfigManager) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getChannelConfigMutex.RLock()
	defer fake.getChannelConfigMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *ConfigManager) recordInvocation(key string, args []interface{}) {
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