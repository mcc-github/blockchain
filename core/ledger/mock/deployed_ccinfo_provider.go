
package mock

import (
	"sync"

	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/ledger/rwset/kvrwset"
)

type DeployedChaincodeInfoProvider struct {
	ChaincodeInfoStub        func(string, string, ledger.SimpleQueryExecutor) (*ledger.DeployedChaincodeInfo, error)
	chaincodeInfoMutex       sync.RWMutex
	chaincodeInfoArgsForCall []struct {
		arg1 string
		arg2 string
		arg3 ledger.SimpleQueryExecutor
	}
	chaincodeInfoReturns struct {
		result1 *ledger.DeployedChaincodeInfo
		result2 error
	}
	chaincodeInfoReturnsOnCall map[int]struct {
		result1 *ledger.DeployedChaincodeInfo
		result2 error
	}
	CollectionInfoStub        func(string, string, string, ledger.SimpleQueryExecutor) (*common.StaticCollectionConfig, error)
	collectionInfoMutex       sync.RWMutex
	collectionInfoArgsForCall []struct {
		arg1 string
		arg2 string
		arg3 string
		arg4 ledger.SimpleQueryExecutor
	}
	collectionInfoReturns struct {
		result1 *common.StaticCollectionConfig
		result2 error
	}
	collectionInfoReturnsOnCall map[int]struct {
		result1 *common.StaticCollectionConfig
		result2 error
	}
	ImplicitCollectionsStub        func(string, string, ledger.SimpleQueryExecutor) ([]*common.StaticCollectionConfig, error)
	implicitCollectionsMutex       sync.RWMutex
	implicitCollectionsArgsForCall []struct {
		arg1 string
		arg2 string
		arg3 ledger.SimpleQueryExecutor
	}
	implicitCollectionsReturns struct {
		result1 []*common.StaticCollectionConfig
		result2 error
	}
	implicitCollectionsReturnsOnCall map[int]struct {
		result1 []*common.StaticCollectionConfig
		result2 error
	}
	NamespacesStub        func() []string
	namespacesMutex       sync.RWMutex
	namespacesArgsForCall []struct {
	}
	namespacesReturns struct {
		result1 []string
	}
	namespacesReturnsOnCall map[int]struct {
		result1 []string
	}
	UpdatedChaincodesStub        func(map[string][]*kvrwset.KVWrite) ([]*ledger.ChaincodeLifecycleInfo, error)
	updatedChaincodesMutex       sync.RWMutex
	updatedChaincodesArgsForCall []struct {
		arg1 map[string][]*kvrwset.KVWrite
	}
	updatedChaincodesReturns struct {
		result1 []*ledger.ChaincodeLifecycleInfo
		result2 error
	}
	updatedChaincodesReturnsOnCall map[int]struct {
		result1 []*ledger.ChaincodeLifecycleInfo
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *DeployedChaincodeInfoProvider) ChaincodeInfo(arg1 string, arg2 string, arg3 ledger.SimpleQueryExecutor) (*ledger.DeployedChaincodeInfo, error) {
	fake.chaincodeInfoMutex.Lock()
	ret, specificReturn := fake.chaincodeInfoReturnsOnCall[len(fake.chaincodeInfoArgsForCall)]
	fake.chaincodeInfoArgsForCall = append(fake.chaincodeInfoArgsForCall, struct {
		arg1 string
		arg2 string
		arg3 ledger.SimpleQueryExecutor
	}{arg1, arg2, arg3})
	fake.recordInvocation("ChaincodeInfo", []interface{}{arg1, arg2, arg3})
	fake.chaincodeInfoMutex.Unlock()
	if fake.ChaincodeInfoStub != nil {
		return fake.ChaincodeInfoStub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.chaincodeInfoReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *DeployedChaincodeInfoProvider) ChaincodeInfoCallCount() int {
	fake.chaincodeInfoMutex.RLock()
	defer fake.chaincodeInfoMutex.RUnlock()
	return len(fake.chaincodeInfoArgsForCall)
}

func (fake *DeployedChaincodeInfoProvider) ChaincodeInfoCalls(stub func(string, string, ledger.SimpleQueryExecutor) (*ledger.DeployedChaincodeInfo, error)) {
	fake.chaincodeInfoMutex.Lock()
	defer fake.chaincodeInfoMutex.Unlock()
	fake.ChaincodeInfoStub = stub
}

func (fake *DeployedChaincodeInfoProvider) ChaincodeInfoArgsForCall(i int) (string, string, ledger.SimpleQueryExecutor) {
	fake.chaincodeInfoMutex.RLock()
	defer fake.chaincodeInfoMutex.RUnlock()
	argsForCall := fake.chaincodeInfoArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *DeployedChaincodeInfoProvider) ChaincodeInfoReturns(result1 *ledger.DeployedChaincodeInfo, result2 error) {
	fake.chaincodeInfoMutex.Lock()
	defer fake.chaincodeInfoMutex.Unlock()
	fake.ChaincodeInfoStub = nil
	fake.chaincodeInfoReturns = struct {
		result1 *ledger.DeployedChaincodeInfo
		result2 error
	}{result1, result2}
}

func (fake *DeployedChaincodeInfoProvider) ChaincodeInfoReturnsOnCall(i int, result1 *ledger.DeployedChaincodeInfo, result2 error) {
	fake.chaincodeInfoMutex.Lock()
	defer fake.chaincodeInfoMutex.Unlock()
	fake.ChaincodeInfoStub = nil
	if fake.chaincodeInfoReturnsOnCall == nil {
		fake.chaincodeInfoReturnsOnCall = make(map[int]struct {
			result1 *ledger.DeployedChaincodeInfo
			result2 error
		})
	}
	fake.chaincodeInfoReturnsOnCall[i] = struct {
		result1 *ledger.DeployedChaincodeInfo
		result2 error
	}{result1, result2}
}

func (fake *DeployedChaincodeInfoProvider) CollectionInfo(arg1 string, arg2 string, arg3 string, arg4 ledger.SimpleQueryExecutor) (*common.StaticCollectionConfig, error) {
	fake.collectionInfoMutex.Lock()
	ret, specificReturn := fake.collectionInfoReturnsOnCall[len(fake.collectionInfoArgsForCall)]
	fake.collectionInfoArgsForCall = append(fake.collectionInfoArgsForCall, struct {
		arg1 string
		arg2 string
		arg3 string
		arg4 ledger.SimpleQueryExecutor
	}{arg1, arg2, arg3, arg4})
	fake.recordInvocation("CollectionInfo", []interface{}{arg1, arg2, arg3, arg4})
	fake.collectionInfoMutex.Unlock()
	if fake.CollectionInfoStub != nil {
		return fake.CollectionInfoStub(arg1, arg2, arg3, arg4)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.collectionInfoReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *DeployedChaincodeInfoProvider) CollectionInfoCallCount() int {
	fake.collectionInfoMutex.RLock()
	defer fake.collectionInfoMutex.RUnlock()
	return len(fake.collectionInfoArgsForCall)
}

func (fake *DeployedChaincodeInfoProvider) CollectionInfoCalls(stub func(string, string, string, ledger.SimpleQueryExecutor) (*common.StaticCollectionConfig, error)) {
	fake.collectionInfoMutex.Lock()
	defer fake.collectionInfoMutex.Unlock()
	fake.CollectionInfoStub = stub
}

func (fake *DeployedChaincodeInfoProvider) CollectionInfoArgsForCall(i int) (string, string, string, ledger.SimpleQueryExecutor) {
	fake.collectionInfoMutex.RLock()
	defer fake.collectionInfoMutex.RUnlock()
	argsForCall := fake.collectionInfoArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4
}

func (fake *DeployedChaincodeInfoProvider) CollectionInfoReturns(result1 *common.StaticCollectionConfig, result2 error) {
	fake.collectionInfoMutex.Lock()
	defer fake.collectionInfoMutex.Unlock()
	fake.CollectionInfoStub = nil
	fake.collectionInfoReturns = struct {
		result1 *common.StaticCollectionConfig
		result2 error
	}{result1, result2}
}

func (fake *DeployedChaincodeInfoProvider) CollectionInfoReturnsOnCall(i int, result1 *common.StaticCollectionConfig, result2 error) {
	fake.collectionInfoMutex.Lock()
	defer fake.collectionInfoMutex.Unlock()
	fake.CollectionInfoStub = nil
	if fake.collectionInfoReturnsOnCall == nil {
		fake.collectionInfoReturnsOnCall = make(map[int]struct {
			result1 *common.StaticCollectionConfig
			result2 error
		})
	}
	fake.collectionInfoReturnsOnCall[i] = struct {
		result1 *common.StaticCollectionConfig
		result2 error
	}{result1, result2}
}

func (fake *DeployedChaincodeInfoProvider) ImplicitCollections(arg1 string, arg2 string, arg3 ledger.SimpleQueryExecutor) ([]*common.StaticCollectionConfig, error) {
	fake.implicitCollectionsMutex.Lock()
	ret, specificReturn := fake.implicitCollectionsReturnsOnCall[len(fake.implicitCollectionsArgsForCall)]
	fake.implicitCollectionsArgsForCall = append(fake.implicitCollectionsArgsForCall, struct {
		arg1 string
		arg2 string
		arg3 ledger.SimpleQueryExecutor
	}{arg1, arg2, arg3})
	fake.recordInvocation("ImplicitCollections", []interface{}{arg1, arg2, arg3})
	fake.implicitCollectionsMutex.Unlock()
	if fake.ImplicitCollectionsStub != nil {
		return fake.ImplicitCollectionsStub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.implicitCollectionsReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *DeployedChaincodeInfoProvider) ImplicitCollectionsCallCount() int {
	fake.implicitCollectionsMutex.RLock()
	defer fake.implicitCollectionsMutex.RUnlock()
	return len(fake.implicitCollectionsArgsForCall)
}

func (fake *DeployedChaincodeInfoProvider) ImplicitCollectionsCalls(stub func(string, string, ledger.SimpleQueryExecutor) ([]*common.StaticCollectionConfig, error)) {
	fake.implicitCollectionsMutex.Lock()
	defer fake.implicitCollectionsMutex.Unlock()
	fake.ImplicitCollectionsStub = stub
}

func (fake *DeployedChaincodeInfoProvider) ImplicitCollectionsArgsForCall(i int) (string, string, ledger.SimpleQueryExecutor) {
	fake.implicitCollectionsMutex.RLock()
	defer fake.implicitCollectionsMutex.RUnlock()
	argsForCall := fake.implicitCollectionsArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *DeployedChaincodeInfoProvider) ImplicitCollectionsReturns(result1 []*common.StaticCollectionConfig, result2 error) {
	fake.implicitCollectionsMutex.Lock()
	defer fake.implicitCollectionsMutex.Unlock()
	fake.ImplicitCollectionsStub = nil
	fake.implicitCollectionsReturns = struct {
		result1 []*common.StaticCollectionConfig
		result2 error
	}{result1, result2}
}

func (fake *DeployedChaincodeInfoProvider) ImplicitCollectionsReturnsOnCall(i int, result1 []*common.StaticCollectionConfig, result2 error) {
	fake.implicitCollectionsMutex.Lock()
	defer fake.implicitCollectionsMutex.Unlock()
	fake.ImplicitCollectionsStub = nil
	if fake.implicitCollectionsReturnsOnCall == nil {
		fake.implicitCollectionsReturnsOnCall = make(map[int]struct {
			result1 []*common.StaticCollectionConfig
			result2 error
		})
	}
	fake.implicitCollectionsReturnsOnCall[i] = struct {
		result1 []*common.StaticCollectionConfig
		result2 error
	}{result1, result2}
}

func (fake *DeployedChaincodeInfoProvider) Namespaces() []string {
	fake.namespacesMutex.Lock()
	ret, specificReturn := fake.namespacesReturnsOnCall[len(fake.namespacesArgsForCall)]
	fake.namespacesArgsForCall = append(fake.namespacesArgsForCall, struct {
	}{})
	fake.recordInvocation("Namespaces", []interface{}{})
	fake.namespacesMutex.Unlock()
	if fake.NamespacesStub != nil {
		return fake.NamespacesStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.namespacesReturns
	return fakeReturns.result1
}

func (fake *DeployedChaincodeInfoProvider) NamespacesCallCount() int {
	fake.namespacesMutex.RLock()
	defer fake.namespacesMutex.RUnlock()
	return len(fake.namespacesArgsForCall)
}

func (fake *DeployedChaincodeInfoProvider) NamespacesCalls(stub func() []string) {
	fake.namespacesMutex.Lock()
	defer fake.namespacesMutex.Unlock()
	fake.NamespacesStub = stub
}

func (fake *DeployedChaincodeInfoProvider) NamespacesReturns(result1 []string) {
	fake.namespacesMutex.Lock()
	defer fake.namespacesMutex.Unlock()
	fake.NamespacesStub = nil
	fake.namespacesReturns = struct {
		result1 []string
	}{result1}
}

func (fake *DeployedChaincodeInfoProvider) NamespacesReturnsOnCall(i int, result1 []string) {
	fake.namespacesMutex.Lock()
	defer fake.namespacesMutex.Unlock()
	fake.NamespacesStub = nil
	if fake.namespacesReturnsOnCall == nil {
		fake.namespacesReturnsOnCall = make(map[int]struct {
			result1 []string
		})
	}
	fake.namespacesReturnsOnCall[i] = struct {
		result1 []string
	}{result1}
}

func (fake *DeployedChaincodeInfoProvider) UpdatedChaincodes(arg1 map[string][]*kvrwset.KVWrite) ([]*ledger.ChaincodeLifecycleInfo, error) {
	fake.updatedChaincodesMutex.Lock()
	ret, specificReturn := fake.updatedChaincodesReturnsOnCall[len(fake.updatedChaincodesArgsForCall)]
	fake.updatedChaincodesArgsForCall = append(fake.updatedChaincodesArgsForCall, struct {
		arg1 map[string][]*kvrwset.KVWrite
	}{arg1})
	fake.recordInvocation("UpdatedChaincodes", []interface{}{arg1})
	fake.updatedChaincodesMutex.Unlock()
	if fake.UpdatedChaincodesStub != nil {
		return fake.UpdatedChaincodesStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.updatedChaincodesReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *DeployedChaincodeInfoProvider) UpdatedChaincodesCallCount() int {
	fake.updatedChaincodesMutex.RLock()
	defer fake.updatedChaincodesMutex.RUnlock()
	return len(fake.updatedChaincodesArgsForCall)
}

func (fake *DeployedChaincodeInfoProvider) UpdatedChaincodesCalls(stub func(map[string][]*kvrwset.KVWrite) ([]*ledger.ChaincodeLifecycleInfo, error)) {
	fake.updatedChaincodesMutex.Lock()
	defer fake.updatedChaincodesMutex.Unlock()
	fake.UpdatedChaincodesStub = stub
}

func (fake *DeployedChaincodeInfoProvider) UpdatedChaincodesArgsForCall(i int) map[string][]*kvrwset.KVWrite {
	fake.updatedChaincodesMutex.RLock()
	defer fake.updatedChaincodesMutex.RUnlock()
	argsForCall := fake.updatedChaincodesArgsForCall[i]
	return argsForCall.arg1
}

func (fake *DeployedChaincodeInfoProvider) UpdatedChaincodesReturns(result1 []*ledger.ChaincodeLifecycleInfo, result2 error) {
	fake.updatedChaincodesMutex.Lock()
	defer fake.updatedChaincodesMutex.Unlock()
	fake.UpdatedChaincodesStub = nil
	fake.updatedChaincodesReturns = struct {
		result1 []*ledger.ChaincodeLifecycleInfo
		result2 error
	}{result1, result2}
}

func (fake *DeployedChaincodeInfoProvider) UpdatedChaincodesReturnsOnCall(i int, result1 []*ledger.ChaincodeLifecycleInfo, result2 error) {
	fake.updatedChaincodesMutex.Lock()
	defer fake.updatedChaincodesMutex.Unlock()
	fake.UpdatedChaincodesStub = nil
	if fake.updatedChaincodesReturnsOnCall == nil {
		fake.updatedChaincodesReturnsOnCall = make(map[int]struct {
			result1 []*ledger.ChaincodeLifecycleInfo
			result2 error
		})
	}
	fake.updatedChaincodesReturnsOnCall[i] = struct {
		result1 []*ledger.ChaincodeLifecycleInfo
		result2 error
	}{result1, result2}
}

func (fake *DeployedChaincodeInfoProvider) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.chaincodeInfoMutex.RLock()
	defer fake.chaincodeInfoMutex.RUnlock()
	fake.collectionInfoMutex.RLock()
	defer fake.collectionInfoMutex.RUnlock()
	fake.implicitCollectionsMutex.RLock()
	defer fake.implicitCollectionsMutex.RUnlock()
	fake.namespacesMutex.RLock()
	defer fake.namespacesMutex.RUnlock()
	fake.updatedChaincodesMutex.RLock()
	defer fake.updatedChaincodesMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *DeployedChaincodeInfoProvider) recordInvocation(key string, args []interface{}) {
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

var _ ledger.DeployedChaincodeInfoProvider = new(DeployedChaincodeInfoProvider)
