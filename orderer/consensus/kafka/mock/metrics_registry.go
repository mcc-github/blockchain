
package mock

import (
	sync "sync"
)

type MetricsRegistry struct {
	EachStub        func(func(string, interface{}))
	eachMutex       sync.RWMutex
	eachArgsForCall []struct {
		arg1 func(string, interface{})
	}
	GetStub        func(string) interface{}
	getMutex       sync.RWMutex
	getArgsForCall []struct {
		arg1 string
	}
	getReturns struct {
		result1 interface{}
	}
	getReturnsOnCall map[int]struct {
		result1 interface{}
	}
	GetAllStub        func() map[string]map[string]interface{}
	getAllMutex       sync.RWMutex
	getAllArgsForCall []struct {
	}
	getAllReturns struct {
		result1 map[string]map[string]interface{}
	}
	getAllReturnsOnCall map[int]struct {
		result1 map[string]map[string]interface{}
	}
	GetOrRegisterStub        func(string, interface{}) interface{}
	getOrRegisterMutex       sync.RWMutex
	getOrRegisterArgsForCall []struct {
		arg1 string
		arg2 interface{}
	}
	getOrRegisterReturns struct {
		result1 interface{}
	}
	getOrRegisterReturnsOnCall map[int]struct {
		result1 interface{}
	}
	RegisterStub        func(string, interface{}) error
	registerMutex       sync.RWMutex
	registerArgsForCall []struct {
		arg1 string
		arg2 interface{}
	}
	registerReturns struct {
		result1 error
	}
	registerReturnsOnCall map[int]struct {
		result1 error
	}
	RunHealthchecksStub        func()
	runHealthchecksMutex       sync.RWMutex
	runHealthchecksArgsForCall []struct {
	}
	UnregisterStub        func(string)
	unregisterMutex       sync.RWMutex
	unregisterArgsForCall []struct {
		arg1 string
	}
	UnregisterAllStub        func()
	unregisterAllMutex       sync.RWMutex
	unregisterAllArgsForCall []struct {
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *MetricsRegistry) Each(arg1 func(string, interface{})) {
	fake.eachMutex.Lock()
	fake.eachArgsForCall = append(fake.eachArgsForCall, struct {
		arg1 func(string, interface{})
	}{arg1})
	fake.recordInvocation("Each", []interface{}{arg1})
	fake.eachMutex.Unlock()
	if fake.EachStub != nil {
		fake.EachStub(arg1)
	}
}

func (fake *MetricsRegistry) EachCallCount() int {
	fake.eachMutex.RLock()
	defer fake.eachMutex.RUnlock()
	return len(fake.eachArgsForCall)
}

func (fake *MetricsRegistry) EachCalls(stub func(func(string, interface{}))) {
	fake.eachMutex.Lock()
	defer fake.eachMutex.Unlock()
	fake.EachStub = stub
}

func (fake *MetricsRegistry) EachArgsForCall(i int) func(string, interface{}) {
	fake.eachMutex.RLock()
	defer fake.eachMutex.RUnlock()
	argsForCall := fake.eachArgsForCall[i]
	return argsForCall.arg1
}

func (fake *MetricsRegistry) Get(arg1 string) interface{} {
	fake.getMutex.Lock()
	ret, specificReturn := fake.getReturnsOnCall[len(fake.getArgsForCall)]
	fake.getArgsForCall = append(fake.getArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("Get", []interface{}{arg1})
	fake.getMutex.Unlock()
	if fake.GetStub != nil {
		return fake.GetStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.getReturns
	return fakeReturns.result1
}

func (fake *MetricsRegistry) GetCallCount() int {
	fake.getMutex.RLock()
	defer fake.getMutex.RUnlock()
	return len(fake.getArgsForCall)
}

func (fake *MetricsRegistry) GetCalls(stub func(string) interface{}) {
	fake.getMutex.Lock()
	defer fake.getMutex.Unlock()
	fake.GetStub = stub
}

func (fake *MetricsRegistry) GetArgsForCall(i int) string {
	fake.getMutex.RLock()
	defer fake.getMutex.RUnlock()
	argsForCall := fake.getArgsForCall[i]
	return argsForCall.arg1
}

func (fake *MetricsRegistry) GetReturns(result1 interface{}) {
	fake.getMutex.Lock()
	defer fake.getMutex.Unlock()
	fake.GetStub = nil
	fake.getReturns = struct {
		result1 interface{}
	}{result1}
}

func (fake *MetricsRegistry) GetReturnsOnCall(i int, result1 interface{}) {
	fake.getMutex.Lock()
	defer fake.getMutex.Unlock()
	fake.GetStub = nil
	if fake.getReturnsOnCall == nil {
		fake.getReturnsOnCall = make(map[int]struct {
			result1 interface{}
		})
	}
	fake.getReturnsOnCall[i] = struct {
		result1 interface{}
	}{result1}
}

func (fake *MetricsRegistry) GetAll() map[string]map[string]interface{} {
	fake.getAllMutex.Lock()
	ret, specificReturn := fake.getAllReturnsOnCall[len(fake.getAllArgsForCall)]
	fake.getAllArgsForCall = append(fake.getAllArgsForCall, struct {
	}{})
	fake.recordInvocation("GetAll", []interface{}{})
	fake.getAllMutex.Unlock()
	if fake.GetAllStub != nil {
		return fake.GetAllStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.getAllReturns
	return fakeReturns.result1
}

func (fake *MetricsRegistry) GetAllCallCount() int {
	fake.getAllMutex.RLock()
	defer fake.getAllMutex.RUnlock()
	return len(fake.getAllArgsForCall)
}

func (fake *MetricsRegistry) GetAllCalls(stub func() map[string]map[string]interface{}) {
	fake.getAllMutex.Lock()
	defer fake.getAllMutex.Unlock()
	fake.GetAllStub = stub
}

func (fake *MetricsRegistry) GetAllReturns(result1 map[string]map[string]interface{}) {
	fake.getAllMutex.Lock()
	defer fake.getAllMutex.Unlock()
	fake.GetAllStub = nil
	fake.getAllReturns = struct {
		result1 map[string]map[string]interface{}
	}{result1}
}

func (fake *MetricsRegistry) GetAllReturnsOnCall(i int, result1 map[string]map[string]interface{}) {
	fake.getAllMutex.Lock()
	defer fake.getAllMutex.Unlock()
	fake.GetAllStub = nil
	if fake.getAllReturnsOnCall == nil {
		fake.getAllReturnsOnCall = make(map[int]struct {
			result1 map[string]map[string]interface{}
		})
	}
	fake.getAllReturnsOnCall[i] = struct {
		result1 map[string]map[string]interface{}
	}{result1}
}

func (fake *MetricsRegistry) GetOrRegister(arg1 string, arg2 interface{}) interface{} {
	fake.getOrRegisterMutex.Lock()
	ret, specificReturn := fake.getOrRegisterReturnsOnCall[len(fake.getOrRegisterArgsForCall)]
	fake.getOrRegisterArgsForCall = append(fake.getOrRegisterArgsForCall, struct {
		arg1 string
		arg2 interface{}
	}{arg1, arg2})
	fake.recordInvocation("GetOrRegister", []interface{}{arg1, arg2})
	fake.getOrRegisterMutex.Unlock()
	if fake.GetOrRegisterStub != nil {
		return fake.GetOrRegisterStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.getOrRegisterReturns
	return fakeReturns.result1
}

func (fake *MetricsRegistry) GetOrRegisterCallCount() int {
	fake.getOrRegisterMutex.RLock()
	defer fake.getOrRegisterMutex.RUnlock()
	return len(fake.getOrRegisterArgsForCall)
}

func (fake *MetricsRegistry) GetOrRegisterCalls(stub func(string, interface{}) interface{}) {
	fake.getOrRegisterMutex.Lock()
	defer fake.getOrRegisterMutex.Unlock()
	fake.GetOrRegisterStub = stub
}

func (fake *MetricsRegistry) GetOrRegisterArgsForCall(i int) (string, interface{}) {
	fake.getOrRegisterMutex.RLock()
	defer fake.getOrRegisterMutex.RUnlock()
	argsForCall := fake.getOrRegisterArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *MetricsRegistry) GetOrRegisterReturns(result1 interface{}) {
	fake.getOrRegisterMutex.Lock()
	defer fake.getOrRegisterMutex.Unlock()
	fake.GetOrRegisterStub = nil
	fake.getOrRegisterReturns = struct {
		result1 interface{}
	}{result1}
}

func (fake *MetricsRegistry) GetOrRegisterReturnsOnCall(i int, result1 interface{}) {
	fake.getOrRegisterMutex.Lock()
	defer fake.getOrRegisterMutex.Unlock()
	fake.GetOrRegisterStub = nil
	if fake.getOrRegisterReturnsOnCall == nil {
		fake.getOrRegisterReturnsOnCall = make(map[int]struct {
			result1 interface{}
		})
	}
	fake.getOrRegisterReturnsOnCall[i] = struct {
		result1 interface{}
	}{result1}
}

func (fake *MetricsRegistry) Register(arg1 string, arg2 interface{}) error {
	fake.registerMutex.Lock()
	ret, specificReturn := fake.registerReturnsOnCall[len(fake.registerArgsForCall)]
	fake.registerArgsForCall = append(fake.registerArgsForCall, struct {
		arg1 string
		arg2 interface{}
	}{arg1, arg2})
	fake.recordInvocation("Register", []interface{}{arg1, arg2})
	fake.registerMutex.Unlock()
	if fake.RegisterStub != nil {
		return fake.RegisterStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.registerReturns
	return fakeReturns.result1
}

func (fake *MetricsRegistry) RegisterCallCount() int {
	fake.registerMutex.RLock()
	defer fake.registerMutex.RUnlock()
	return len(fake.registerArgsForCall)
}

func (fake *MetricsRegistry) RegisterCalls(stub func(string, interface{}) error) {
	fake.registerMutex.Lock()
	defer fake.registerMutex.Unlock()
	fake.RegisterStub = stub
}

func (fake *MetricsRegistry) RegisterArgsForCall(i int) (string, interface{}) {
	fake.registerMutex.RLock()
	defer fake.registerMutex.RUnlock()
	argsForCall := fake.registerArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *MetricsRegistry) RegisterReturns(result1 error) {
	fake.registerMutex.Lock()
	defer fake.registerMutex.Unlock()
	fake.RegisterStub = nil
	fake.registerReturns = struct {
		result1 error
	}{result1}
}

func (fake *MetricsRegistry) RegisterReturnsOnCall(i int, result1 error) {
	fake.registerMutex.Lock()
	defer fake.registerMutex.Unlock()
	fake.RegisterStub = nil
	if fake.registerReturnsOnCall == nil {
		fake.registerReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.registerReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *MetricsRegistry) RunHealthchecks() {
	fake.runHealthchecksMutex.Lock()
	fake.runHealthchecksArgsForCall = append(fake.runHealthchecksArgsForCall, struct {
	}{})
	fake.recordInvocation("RunHealthchecks", []interface{}{})
	fake.runHealthchecksMutex.Unlock()
	if fake.RunHealthchecksStub != nil {
		fake.RunHealthchecksStub()
	}
}

func (fake *MetricsRegistry) RunHealthchecksCallCount() int {
	fake.runHealthchecksMutex.RLock()
	defer fake.runHealthchecksMutex.RUnlock()
	return len(fake.runHealthchecksArgsForCall)
}

func (fake *MetricsRegistry) RunHealthchecksCalls(stub func()) {
	fake.runHealthchecksMutex.Lock()
	defer fake.runHealthchecksMutex.Unlock()
	fake.RunHealthchecksStub = stub
}

func (fake *MetricsRegistry) Unregister(arg1 string) {
	fake.unregisterMutex.Lock()
	fake.unregisterArgsForCall = append(fake.unregisterArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("Unregister", []interface{}{arg1})
	fake.unregisterMutex.Unlock()
	if fake.UnregisterStub != nil {
		fake.UnregisterStub(arg1)
	}
}

func (fake *MetricsRegistry) UnregisterCallCount() int {
	fake.unregisterMutex.RLock()
	defer fake.unregisterMutex.RUnlock()
	return len(fake.unregisterArgsForCall)
}

func (fake *MetricsRegistry) UnregisterCalls(stub func(string)) {
	fake.unregisterMutex.Lock()
	defer fake.unregisterMutex.Unlock()
	fake.UnregisterStub = stub
}

func (fake *MetricsRegistry) UnregisterArgsForCall(i int) string {
	fake.unregisterMutex.RLock()
	defer fake.unregisterMutex.RUnlock()
	argsForCall := fake.unregisterArgsForCall[i]
	return argsForCall.arg1
}

func (fake *MetricsRegistry) UnregisterAll() {
	fake.unregisterAllMutex.Lock()
	fake.unregisterAllArgsForCall = append(fake.unregisterAllArgsForCall, struct {
	}{})
	fake.recordInvocation("UnregisterAll", []interface{}{})
	fake.unregisterAllMutex.Unlock()
	if fake.UnregisterAllStub != nil {
		fake.UnregisterAllStub()
	}
}

func (fake *MetricsRegistry) UnregisterAllCallCount() int {
	fake.unregisterAllMutex.RLock()
	defer fake.unregisterAllMutex.RUnlock()
	return len(fake.unregisterAllArgsForCall)
}

func (fake *MetricsRegistry) UnregisterAllCalls(stub func()) {
	fake.unregisterAllMutex.Lock()
	defer fake.unregisterAllMutex.Unlock()
	fake.UnregisterAllStub = stub
}

func (fake *MetricsRegistry) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.eachMutex.RLock()
	defer fake.eachMutex.RUnlock()
	fake.getMutex.RLock()
	defer fake.getMutex.RUnlock()
	fake.getAllMutex.RLock()
	defer fake.getAllMutex.RUnlock()
	fake.getOrRegisterMutex.RLock()
	defer fake.getOrRegisterMutex.RUnlock()
	fake.registerMutex.RLock()
	defer fake.registerMutex.RUnlock()
	fake.runHealthchecksMutex.RLock()
	defer fake.runHealthchecksMutex.RUnlock()
	fake.unregisterMutex.RLock()
	defer fake.unregisterMutex.RUnlock()
	fake.unregisterAllMutex.RLock()
	defer fake.unregisterAllMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *MetricsRegistry) recordInvocation(key string, args []interface{}) {
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
