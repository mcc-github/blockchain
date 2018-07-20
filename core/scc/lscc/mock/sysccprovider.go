
package mock

import (
	"sync"

	"github.com/mcc-github/blockchain/common/channelconfig"
	"github.com/mcc-github/blockchain/common/policies"
	"github.com/mcc-github/blockchain/core/ledger"
)

type SystemChaincodeProvider struct {
	IsSysCCStub        func(name string) bool
	isSysCCMutex       sync.RWMutex
	isSysCCArgsForCall []struct {
		name string
	}
	isSysCCReturns struct {
		result1 bool
	}
	isSysCCReturnsOnCall map[int]struct {
		result1 bool
	}
	IsSysCCAndNotInvokableCC2CCStub        func(name string) bool
	isSysCCAndNotInvokableCC2CCMutex       sync.RWMutex
	isSysCCAndNotInvokableCC2CCArgsForCall []struct {
		name string
	}
	isSysCCAndNotInvokableCC2CCReturns struct {
		result1 bool
	}
	isSysCCAndNotInvokableCC2CCReturnsOnCall map[int]struct {
		result1 bool
	}
	IsSysCCAndNotInvokableExternalStub        func(name string) bool
	isSysCCAndNotInvokableExternalMutex       sync.RWMutex
	isSysCCAndNotInvokableExternalArgsForCall []struct {
		name string
	}
	isSysCCAndNotInvokableExternalReturns struct {
		result1 bool
	}
	isSysCCAndNotInvokableExternalReturnsOnCall map[int]struct {
		result1 bool
	}
	GetQueryExecutorForLedgerStub        func(cid string) (ledger.QueryExecutor, error)
	getQueryExecutorForLedgerMutex       sync.RWMutex
	getQueryExecutorForLedgerArgsForCall []struct {
		cid string
	}
	getQueryExecutorForLedgerReturns struct {
		result1 ledger.QueryExecutor
		result2 error
	}
	getQueryExecutorForLedgerReturnsOnCall map[int]struct {
		result1 ledger.QueryExecutor
		result2 error
	}
	GetApplicationConfigStub        func(cid string) (channelconfig.Application, bool)
	getApplicationConfigMutex       sync.RWMutex
	getApplicationConfigArgsForCall []struct {
		cid string
	}
	getApplicationConfigReturns struct {
		result1 channelconfig.Application
		result2 bool
	}
	getApplicationConfigReturnsOnCall map[int]struct {
		result1 channelconfig.Application
		result2 bool
	}
	PolicyManagerStub        func(channelID string) (policies.Manager, bool)
	policyManagerMutex       sync.RWMutex
	policyManagerArgsForCall []struct {
		channelID string
	}
	policyManagerReturns struct {
		result1 policies.Manager
		result2 bool
	}
	policyManagerReturnsOnCall map[int]struct {
		result1 policies.Manager
		result2 bool
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *SystemChaincodeProvider) IsSysCC(name string) bool {
	fake.isSysCCMutex.Lock()
	ret, specificReturn := fake.isSysCCReturnsOnCall[len(fake.isSysCCArgsForCall)]
	fake.isSysCCArgsForCall = append(fake.isSysCCArgsForCall, struct {
		name string
	}{name})
	fake.recordInvocation("IsSysCC", []interface{}{name})
	fake.isSysCCMutex.Unlock()
	if fake.IsSysCCStub != nil {
		return fake.IsSysCCStub(name)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.isSysCCReturns.result1
}

func (fake *SystemChaincodeProvider) IsSysCCCallCount() int {
	fake.isSysCCMutex.RLock()
	defer fake.isSysCCMutex.RUnlock()
	return len(fake.isSysCCArgsForCall)
}

func (fake *SystemChaincodeProvider) IsSysCCArgsForCall(i int) string {
	fake.isSysCCMutex.RLock()
	defer fake.isSysCCMutex.RUnlock()
	return fake.isSysCCArgsForCall[i].name
}

func (fake *SystemChaincodeProvider) IsSysCCReturns(result1 bool) {
	fake.IsSysCCStub = nil
	fake.isSysCCReturns = struct {
		result1 bool
	}{result1}
}

func (fake *SystemChaincodeProvider) IsSysCCReturnsOnCall(i int, result1 bool) {
	fake.IsSysCCStub = nil
	if fake.isSysCCReturnsOnCall == nil {
		fake.isSysCCReturnsOnCall = make(map[int]struct {
			result1 bool
		})
	}
	fake.isSysCCReturnsOnCall[i] = struct {
		result1 bool
	}{result1}
}

func (fake *SystemChaincodeProvider) IsSysCCAndNotInvokableCC2CC(name string) bool {
	fake.isSysCCAndNotInvokableCC2CCMutex.Lock()
	ret, specificReturn := fake.isSysCCAndNotInvokableCC2CCReturnsOnCall[len(fake.isSysCCAndNotInvokableCC2CCArgsForCall)]
	fake.isSysCCAndNotInvokableCC2CCArgsForCall = append(fake.isSysCCAndNotInvokableCC2CCArgsForCall, struct {
		name string
	}{name})
	fake.recordInvocation("IsSysCCAndNotInvokableCC2CC", []interface{}{name})
	fake.isSysCCAndNotInvokableCC2CCMutex.Unlock()
	if fake.IsSysCCAndNotInvokableCC2CCStub != nil {
		return fake.IsSysCCAndNotInvokableCC2CCStub(name)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.isSysCCAndNotInvokableCC2CCReturns.result1
}

func (fake *SystemChaincodeProvider) IsSysCCAndNotInvokableCC2CCCallCount() int {
	fake.isSysCCAndNotInvokableCC2CCMutex.RLock()
	defer fake.isSysCCAndNotInvokableCC2CCMutex.RUnlock()
	return len(fake.isSysCCAndNotInvokableCC2CCArgsForCall)
}

func (fake *SystemChaincodeProvider) IsSysCCAndNotInvokableCC2CCArgsForCall(i int) string {
	fake.isSysCCAndNotInvokableCC2CCMutex.RLock()
	defer fake.isSysCCAndNotInvokableCC2CCMutex.RUnlock()
	return fake.isSysCCAndNotInvokableCC2CCArgsForCall[i].name
}

func (fake *SystemChaincodeProvider) IsSysCCAndNotInvokableCC2CCReturns(result1 bool) {
	fake.IsSysCCAndNotInvokableCC2CCStub = nil
	fake.isSysCCAndNotInvokableCC2CCReturns = struct {
		result1 bool
	}{result1}
}

func (fake *SystemChaincodeProvider) IsSysCCAndNotInvokableCC2CCReturnsOnCall(i int, result1 bool) {
	fake.IsSysCCAndNotInvokableCC2CCStub = nil
	if fake.isSysCCAndNotInvokableCC2CCReturnsOnCall == nil {
		fake.isSysCCAndNotInvokableCC2CCReturnsOnCall = make(map[int]struct {
			result1 bool
		})
	}
	fake.isSysCCAndNotInvokableCC2CCReturnsOnCall[i] = struct {
		result1 bool
	}{result1}
}

func (fake *SystemChaincodeProvider) IsSysCCAndNotInvokableExternal(name string) bool {
	fake.isSysCCAndNotInvokableExternalMutex.Lock()
	ret, specificReturn := fake.isSysCCAndNotInvokableExternalReturnsOnCall[len(fake.isSysCCAndNotInvokableExternalArgsForCall)]
	fake.isSysCCAndNotInvokableExternalArgsForCall = append(fake.isSysCCAndNotInvokableExternalArgsForCall, struct {
		name string
	}{name})
	fake.recordInvocation("IsSysCCAndNotInvokableExternal", []interface{}{name})
	fake.isSysCCAndNotInvokableExternalMutex.Unlock()
	if fake.IsSysCCAndNotInvokableExternalStub != nil {
		return fake.IsSysCCAndNotInvokableExternalStub(name)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.isSysCCAndNotInvokableExternalReturns.result1
}

func (fake *SystemChaincodeProvider) IsSysCCAndNotInvokableExternalCallCount() int {
	fake.isSysCCAndNotInvokableExternalMutex.RLock()
	defer fake.isSysCCAndNotInvokableExternalMutex.RUnlock()
	return len(fake.isSysCCAndNotInvokableExternalArgsForCall)
}

func (fake *SystemChaincodeProvider) IsSysCCAndNotInvokableExternalArgsForCall(i int) string {
	fake.isSysCCAndNotInvokableExternalMutex.RLock()
	defer fake.isSysCCAndNotInvokableExternalMutex.RUnlock()
	return fake.isSysCCAndNotInvokableExternalArgsForCall[i].name
}

func (fake *SystemChaincodeProvider) IsSysCCAndNotInvokableExternalReturns(result1 bool) {
	fake.IsSysCCAndNotInvokableExternalStub = nil
	fake.isSysCCAndNotInvokableExternalReturns = struct {
		result1 bool
	}{result1}
}

func (fake *SystemChaincodeProvider) IsSysCCAndNotInvokableExternalReturnsOnCall(i int, result1 bool) {
	fake.IsSysCCAndNotInvokableExternalStub = nil
	if fake.isSysCCAndNotInvokableExternalReturnsOnCall == nil {
		fake.isSysCCAndNotInvokableExternalReturnsOnCall = make(map[int]struct {
			result1 bool
		})
	}
	fake.isSysCCAndNotInvokableExternalReturnsOnCall[i] = struct {
		result1 bool
	}{result1}
}

func (fake *SystemChaincodeProvider) GetQueryExecutorForLedger(cid string) (ledger.QueryExecutor, error) {
	fake.getQueryExecutorForLedgerMutex.Lock()
	ret, specificReturn := fake.getQueryExecutorForLedgerReturnsOnCall[len(fake.getQueryExecutorForLedgerArgsForCall)]
	fake.getQueryExecutorForLedgerArgsForCall = append(fake.getQueryExecutorForLedgerArgsForCall, struct {
		cid string
	}{cid})
	fake.recordInvocation("GetQueryExecutorForLedger", []interface{}{cid})
	fake.getQueryExecutorForLedgerMutex.Unlock()
	if fake.GetQueryExecutorForLedgerStub != nil {
		return fake.GetQueryExecutorForLedgerStub(cid)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getQueryExecutorForLedgerReturns.result1, fake.getQueryExecutorForLedgerReturns.result2
}

func (fake *SystemChaincodeProvider) GetQueryExecutorForLedgerCallCount() int {
	fake.getQueryExecutorForLedgerMutex.RLock()
	defer fake.getQueryExecutorForLedgerMutex.RUnlock()
	return len(fake.getQueryExecutorForLedgerArgsForCall)
}

func (fake *SystemChaincodeProvider) GetQueryExecutorForLedgerArgsForCall(i int) string {
	fake.getQueryExecutorForLedgerMutex.RLock()
	defer fake.getQueryExecutorForLedgerMutex.RUnlock()
	return fake.getQueryExecutorForLedgerArgsForCall[i].cid
}

func (fake *SystemChaincodeProvider) GetQueryExecutorForLedgerReturns(result1 ledger.QueryExecutor, result2 error) {
	fake.GetQueryExecutorForLedgerStub = nil
	fake.getQueryExecutorForLedgerReturns = struct {
		result1 ledger.QueryExecutor
		result2 error
	}{result1, result2}
}

func (fake *SystemChaincodeProvider) GetQueryExecutorForLedgerReturnsOnCall(i int, result1 ledger.QueryExecutor, result2 error) {
	fake.GetQueryExecutorForLedgerStub = nil
	if fake.getQueryExecutorForLedgerReturnsOnCall == nil {
		fake.getQueryExecutorForLedgerReturnsOnCall = make(map[int]struct {
			result1 ledger.QueryExecutor
			result2 error
		})
	}
	fake.getQueryExecutorForLedgerReturnsOnCall[i] = struct {
		result1 ledger.QueryExecutor
		result2 error
	}{result1, result2}
}

func (fake *SystemChaincodeProvider) GetApplicationConfig(cid string) (channelconfig.Application, bool) {
	fake.getApplicationConfigMutex.Lock()
	ret, specificReturn := fake.getApplicationConfigReturnsOnCall[len(fake.getApplicationConfigArgsForCall)]
	fake.getApplicationConfigArgsForCall = append(fake.getApplicationConfigArgsForCall, struct {
		cid string
	}{cid})
	fake.recordInvocation("GetApplicationConfig", []interface{}{cid})
	fake.getApplicationConfigMutex.Unlock()
	if fake.GetApplicationConfigStub != nil {
		return fake.GetApplicationConfigStub(cid)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getApplicationConfigReturns.result1, fake.getApplicationConfigReturns.result2
}

func (fake *SystemChaincodeProvider) GetApplicationConfigCallCount() int {
	fake.getApplicationConfigMutex.RLock()
	defer fake.getApplicationConfigMutex.RUnlock()
	return len(fake.getApplicationConfigArgsForCall)
}

func (fake *SystemChaincodeProvider) GetApplicationConfigArgsForCall(i int) string {
	fake.getApplicationConfigMutex.RLock()
	defer fake.getApplicationConfigMutex.RUnlock()
	return fake.getApplicationConfigArgsForCall[i].cid
}

func (fake *SystemChaincodeProvider) GetApplicationConfigReturns(result1 channelconfig.Application, result2 bool) {
	fake.GetApplicationConfigStub = nil
	fake.getApplicationConfigReturns = struct {
		result1 channelconfig.Application
		result2 bool
	}{result1, result2}
}

func (fake *SystemChaincodeProvider) GetApplicationConfigReturnsOnCall(i int, result1 channelconfig.Application, result2 bool) {
	fake.GetApplicationConfigStub = nil
	if fake.getApplicationConfigReturnsOnCall == nil {
		fake.getApplicationConfigReturnsOnCall = make(map[int]struct {
			result1 channelconfig.Application
			result2 bool
		})
	}
	fake.getApplicationConfigReturnsOnCall[i] = struct {
		result1 channelconfig.Application
		result2 bool
	}{result1, result2}
}

func (fake *SystemChaincodeProvider) PolicyManager(channelID string) (policies.Manager, bool) {
	fake.policyManagerMutex.Lock()
	ret, specificReturn := fake.policyManagerReturnsOnCall[len(fake.policyManagerArgsForCall)]
	fake.policyManagerArgsForCall = append(fake.policyManagerArgsForCall, struct {
		channelID string
	}{channelID})
	fake.recordInvocation("PolicyManager", []interface{}{channelID})
	fake.policyManagerMutex.Unlock()
	if fake.PolicyManagerStub != nil {
		return fake.PolicyManagerStub(channelID)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.policyManagerReturns.result1, fake.policyManagerReturns.result2
}

func (fake *SystemChaincodeProvider) PolicyManagerCallCount() int {
	fake.policyManagerMutex.RLock()
	defer fake.policyManagerMutex.RUnlock()
	return len(fake.policyManagerArgsForCall)
}

func (fake *SystemChaincodeProvider) PolicyManagerArgsForCall(i int) string {
	fake.policyManagerMutex.RLock()
	defer fake.policyManagerMutex.RUnlock()
	return fake.policyManagerArgsForCall[i].channelID
}

func (fake *SystemChaincodeProvider) PolicyManagerReturns(result1 policies.Manager, result2 bool) {
	fake.PolicyManagerStub = nil
	fake.policyManagerReturns = struct {
		result1 policies.Manager
		result2 bool
	}{result1, result2}
}

func (fake *SystemChaincodeProvider) PolicyManagerReturnsOnCall(i int, result1 policies.Manager, result2 bool) {
	fake.PolicyManagerStub = nil
	if fake.policyManagerReturnsOnCall == nil {
		fake.policyManagerReturnsOnCall = make(map[int]struct {
			result1 policies.Manager
			result2 bool
		})
	}
	fake.policyManagerReturnsOnCall[i] = struct {
		result1 policies.Manager
		result2 bool
	}{result1, result2}
}

func (fake *SystemChaincodeProvider) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.isSysCCMutex.RLock()
	defer fake.isSysCCMutex.RUnlock()
	fake.isSysCCAndNotInvokableCC2CCMutex.RLock()
	defer fake.isSysCCAndNotInvokableCC2CCMutex.RUnlock()
	fake.isSysCCAndNotInvokableExternalMutex.RLock()
	defer fake.isSysCCAndNotInvokableExternalMutex.RUnlock()
	fake.getQueryExecutorForLedgerMutex.RLock()
	defer fake.getQueryExecutorForLedgerMutex.RUnlock()
	fake.getApplicationConfigMutex.RLock()
	defer fake.getApplicationConfigMutex.RUnlock()
	fake.policyManagerMutex.RLock()
	defer fake.policyManagerMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *SystemChaincodeProvider) recordInvocation(key string, args []interface{}) {
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
