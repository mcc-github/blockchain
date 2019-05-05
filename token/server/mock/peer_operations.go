
package mock

import (
	"sync"

	"github.com/mcc-github/blockchain/common/channelconfig"
	"github.com/mcc-github/blockchain/common/metrics"
	"github.com/mcc-github/blockchain/common/policies"
	"github.com/mcc-github/blockchain/core/chaincode/platforms"
	"github.com/mcc-github/blockchain/core/committer/txvalidator/plugin"
	"github.com/mcc-github/blockchain/core/committer/txvalidator/v20/plugindispatcher"
	"github.com/mcc-github/blockchain/core/common/sysccprovider"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/peer"
)

type PeerOperations struct {
	CreateChainFromBlockStub        func(*common.Block, sysccprovider.SystemChaincodeProvider, ledger.DeployedChaincodeInfoProvider, plugindispatcher.LifecycleResources, plugindispatcher.CollectionAndLifecycleResources) error
	createChainFromBlockMutex       sync.RWMutex
	createChainFromBlockArgsForCall []struct {
		arg1 *common.Block
		arg2 sysccprovider.SystemChaincodeProvider
		arg3 ledger.DeployedChaincodeInfoProvider
		arg4 plugindispatcher.LifecycleResources
		arg5 plugindispatcher.CollectionAndLifecycleResources
	}
	createChainFromBlockReturns struct {
		result1 error
	}
	createChainFromBlockReturnsOnCall map[int]struct {
		result1 error
	}
	GetChannelConfigStub        func(string) channelconfig.Resources
	getChannelConfigMutex       sync.RWMutex
	getChannelConfigArgsForCall []struct {
		arg1 string
	}
	getChannelConfigReturns struct {
		result1 channelconfig.Resources
	}
	getChannelConfigReturnsOnCall map[int]struct {
		result1 channelconfig.Resources
	}
	GetChannelsInfoStub        func() []*peer.ChannelInfo
	getChannelsInfoMutex       sync.RWMutex
	getChannelsInfoArgsForCall []struct {
	}
	getChannelsInfoReturns struct {
		result1 []*peer.ChannelInfo
	}
	getChannelsInfoReturnsOnCall map[int]struct {
		result1 []*peer.ChannelInfo
	}
	GetCurrConfigBlockStub        func(string) *common.Block
	getCurrConfigBlockMutex       sync.RWMutex
	getCurrConfigBlockArgsForCall []struct {
		arg1 string
	}
	getCurrConfigBlockReturns struct {
		result1 *common.Block
	}
	getCurrConfigBlockReturnsOnCall map[int]struct {
		result1 *common.Block
	}
	GetLedgerStub        func(string) ledger.PeerLedger
	getLedgerMutex       sync.RWMutex
	getLedgerArgsForCall []struct {
		arg1 string
	}
	getLedgerReturns struct {
		result1 ledger.PeerLedger
	}
	getLedgerReturnsOnCall map[int]struct {
		result1 ledger.PeerLedger
	}
	GetMSPIDsStub        func(string) []string
	getMSPIDsMutex       sync.RWMutex
	getMSPIDsArgsForCall []struct {
		arg1 string
	}
	getMSPIDsReturns struct {
		result1 []string
	}
	getMSPIDsReturnsOnCall map[int]struct {
		result1 []string
	}
	GetPolicyManagerStub        func(string) policies.Manager
	getPolicyManagerMutex       sync.RWMutex
	getPolicyManagerArgsForCall []struct {
		arg1 string
	}
	getPolicyManagerReturns struct {
		result1 policies.Manager
	}
	getPolicyManagerReturnsOnCall map[int]struct {
		result1 policies.Manager
	}
	GetStableChannelConfigStub        func(string) channelconfig.Resources
	getStableChannelConfigMutex       sync.RWMutex
	getStableChannelConfigArgsForCall []struct {
		arg1 string
	}
	getStableChannelConfigReturns struct {
		result1 channelconfig.Resources
	}
	getStableChannelConfigReturnsOnCall map[int]struct {
		result1 channelconfig.Resources
	}
	InitChainStub        func(string)
	initChainMutex       sync.RWMutex
	initChainArgsForCall []struct {
		arg1 string
	}
	InitializeStub        func(func(string), sysccprovider.SystemChaincodeProvider, plugin.Mapper, *platforms.Registry, ledger.DeployedChaincodeInfoProvider, ledger.MembershipInfoProvider, metrics.Provider, plugindispatcher.LifecycleResources, plugindispatcher.CollectionAndLifecycleResources, *ledger.Config, int)
	initializeMutex       sync.RWMutex
	initializeArgsForCall []struct {
		arg1  func(string)
		arg2  sysccprovider.SystemChaincodeProvider
		arg3  plugin.Mapper
		arg4  *platforms.Registry
		arg5  ledger.DeployedChaincodeInfoProvider
		arg6  ledger.MembershipInfoProvider
		arg7  metrics.Provider
		arg8  plugindispatcher.LifecycleResources
		arg9  plugindispatcher.CollectionAndLifecycleResources
		arg10 *ledger.Config
		arg11 int
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *PeerOperations) CreateChainFromBlock(arg1 *common.Block, arg2 sysccprovider.SystemChaincodeProvider, arg3 ledger.DeployedChaincodeInfoProvider, arg4 plugindispatcher.LifecycleResources, arg5 plugindispatcher.CollectionAndLifecycleResources) error {
	fake.createChainFromBlockMutex.Lock()
	ret, specificReturn := fake.createChainFromBlockReturnsOnCall[len(fake.createChainFromBlockArgsForCall)]
	fake.createChainFromBlockArgsForCall = append(fake.createChainFromBlockArgsForCall, struct {
		arg1 *common.Block
		arg2 sysccprovider.SystemChaincodeProvider
		arg3 ledger.DeployedChaincodeInfoProvider
		arg4 plugindispatcher.LifecycleResources
		arg5 plugindispatcher.CollectionAndLifecycleResources
	}{arg1, arg2, arg3, arg4, arg5})
	fake.recordInvocation("CreateChainFromBlock", []interface{}{arg1, arg2, arg3, arg4, arg5})
	fake.createChainFromBlockMutex.Unlock()
	if fake.CreateChainFromBlockStub != nil {
		return fake.CreateChainFromBlockStub(arg1, arg2, arg3, arg4, arg5)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.createChainFromBlockReturns
	return fakeReturns.result1
}

func (fake *PeerOperations) CreateChainFromBlockCallCount() int {
	fake.createChainFromBlockMutex.RLock()
	defer fake.createChainFromBlockMutex.RUnlock()
	return len(fake.createChainFromBlockArgsForCall)
}

func (fake *PeerOperations) CreateChainFromBlockCalls(stub func(*common.Block, sysccprovider.SystemChaincodeProvider, ledger.DeployedChaincodeInfoProvider, plugindispatcher.LifecycleResources, plugindispatcher.CollectionAndLifecycleResources) error) {
	fake.createChainFromBlockMutex.Lock()
	defer fake.createChainFromBlockMutex.Unlock()
	fake.CreateChainFromBlockStub = stub
}

func (fake *PeerOperations) CreateChainFromBlockArgsForCall(i int) (*common.Block, sysccprovider.SystemChaincodeProvider, ledger.DeployedChaincodeInfoProvider, plugindispatcher.LifecycleResources, plugindispatcher.CollectionAndLifecycleResources) {
	fake.createChainFromBlockMutex.RLock()
	defer fake.createChainFromBlockMutex.RUnlock()
	argsForCall := fake.createChainFromBlockArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4, argsForCall.arg5
}

func (fake *PeerOperations) CreateChainFromBlockReturns(result1 error) {
	fake.createChainFromBlockMutex.Lock()
	defer fake.createChainFromBlockMutex.Unlock()
	fake.CreateChainFromBlockStub = nil
	fake.createChainFromBlockReturns = struct {
		result1 error
	}{result1}
}

func (fake *PeerOperations) CreateChainFromBlockReturnsOnCall(i int, result1 error) {
	fake.createChainFromBlockMutex.Lock()
	defer fake.createChainFromBlockMutex.Unlock()
	fake.CreateChainFromBlockStub = nil
	if fake.createChainFromBlockReturnsOnCall == nil {
		fake.createChainFromBlockReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.createChainFromBlockReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *PeerOperations) GetChannelConfig(arg1 string) channelconfig.Resources {
	fake.getChannelConfigMutex.Lock()
	ret, specificReturn := fake.getChannelConfigReturnsOnCall[len(fake.getChannelConfigArgsForCall)]
	fake.getChannelConfigArgsForCall = append(fake.getChannelConfigArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("GetChannelConfig", []interface{}{arg1})
	fake.getChannelConfigMutex.Unlock()
	if fake.GetChannelConfigStub != nil {
		return fake.GetChannelConfigStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.getChannelConfigReturns
	return fakeReturns.result1
}

func (fake *PeerOperations) GetChannelConfigCallCount() int {
	fake.getChannelConfigMutex.RLock()
	defer fake.getChannelConfigMutex.RUnlock()
	return len(fake.getChannelConfigArgsForCall)
}

func (fake *PeerOperations) GetChannelConfigCalls(stub func(string) channelconfig.Resources) {
	fake.getChannelConfigMutex.Lock()
	defer fake.getChannelConfigMutex.Unlock()
	fake.GetChannelConfigStub = stub
}

func (fake *PeerOperations) GetChannelConfigArgsForCall(i int) string {
	fake.getChannelConfigMutex.RLock()
	defer fake.getChannelConfigMutex.RUnlock()
	argsForCall := fake.getChannelConfigArgsForCall[i]
	return argsForCall.arg1
}

func (fake *PeerOperations) GetChannelConfigReturns(result1 channelconfig.Resources) {
	fake.getChannelConfigMutex.Lock()
	defer fake.getChannelConfigMutex.Unlock()
	fake.GetChannelConfigStub = nil
	fake.getChannelConfigReturns = struct {
		result1 channelconfig.Resources
	}{result1}
}

func (fake *PeerOperations) GetChannelConfigReturnsOnCall(i int, result1 channelconfig.Resources) {
	fake.getChannelConfigMutex.Lock()
	defer fake.getChannelConfigMutex.Unlock()
	fake.GetChannelConfigStub = nil
	if fake.getChannelConfigReturnsOnCall == nil {
		fake.getChannelConfigReturnsOnCall = make(map[int]struct {
			result1 channelconfig.Resources
		})
	}
	fake.getChannelConfigReturnsOnCall[i] = struct {
		result1 channelconfig.Resources
	}{result1}
}

func (fake *PeerOperations) GetChannelsInfo() []*peer.ChannelInfo {
	fake.getChannelsInfoMutex.Lock()
	ret, specificReturn := fake.getChannelsInfoReturnsOnCall[len(fake.getChannelsInfoArgsForCall)]
	fake.getChannelsInfoArgsForCall = append(fake.getChannelsInfoArgsForCall, struct {
	}{})
	fake.recordInvocation("GetChannelsInfo", []interface{}{})
	fake.getChannelsInfoMutex.Unlock()
	if fake.GetChannelsInfoStub != nil {
		return fake.GetChannelsInfoStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.getChannelsInfoReturns
	return fakeReturns.result1
}

func (fake *PeerOperations) GetChannelsInfoCallCount() int {
	fake.getChannelsInfoMutex.RLock()
	defer fake.getChannelsInfoMutex.RUnlock()
	return len(fake.getChannelsInfoArgsForCall)
}

func (fake *PeerOperations) GetChannelsInfoCalls(stub func() []*peer.ChannelInfo) {
	fake.getChannelsInfoMutex.Lock()
	defer fake.getChannelsInfoMutex.Unlock()
	fake.GetChannelsInfoStub = stub
}

func (fake *PeerOperations) GetChannelsInfoReturns(result1 []*peer.ChannelInfo) {
	fake.getChannelsInfoMutex.Lock()
	defer fake.getChannelsInfoMutex.Unlock()
	fake.GetChannelsInfoStub = nil
	fake.getChannelsInfoReturns = struct {
		result1 []*peer.ChannelInfo
	}{result1}
}

func (fake *PeerOperations) GetChannelsInfoReturnsOnCall(i int, result1 []*peer.ChannelInfo) {
	fake.getChannelsInfoMutex.Lock()
	defer fake.getChannelsInfoMutex.Unlock()
	fake.GetChannelsInfoStub = nil
	if fake.getChannelsInfoReturnsOnCall == nil {
		fake.getChannelsInfoReturnsOnCall = make(map[int]struct {
			result1 []*peer.ChannelInfo
		})
	}
	fake.getChannelsInfoReturnsOnCall[i] = struct {
		result1 []*peer.ChannelInfo
	}{result1}
}

func (fake *PeerOperations) GetCurrConfigBlock(arg1 string) *common.Block {
	fake.getCurrConfigBlockMutex.Lock()
	ret, specificReturn := fake.getCurrConfigBlockReturnsOnCall[len(fake.getCurrConfigBlockArgsForCall)]
	fake.getCurrConfigBlockArgsForCall = append(fake.getCurrConfigBlockArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("GetCurrConfigBlock", []interface{}{arg1})
	fake.getCurrConfigBlockMutex.Unlock()
	if fake.GetCurrConfigBlockStub != nil {
		return fake.GetCurrConfigBlockStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.getCurrConfigBlockReturns
	return fakeReturns.result1
}

func (fake *PeerOperations) GetCurrConfigBlockCallCount() int {
	fake.getCurrConfigBlockMutex.RLock()
	defer fake.getCurrConfigBlockMutex.RUnlock()
	return len(fake.getCurrConfigBlockArgsForCall)
}

func (fake *PeerOperations) GetCurrConfigBlockCalls(stub func(string) *common.Block) {
	fake.getCurrConfigBlockMutex.Lock()
	defer fake.getCurrConfigBlockMutex.Unlock()
	fake.GetCurrConfigBlockStub = stub
}

func (fake *PeerOperations) GetCurrConfigBlockArgsForCall(i int) string {
	fake.getCurrConfigBlockMutex.RLock()
	defer fake.getCurrConfigBlockMutex.RUnlock()
	argsForCall := fake.getCurrConfigBlockArgsForCall[i]
	return argsForCall.arg1
}

func (fake *PeerOperations) GetCurrConfigBlockReturns(result1 *common.Block) {
	fake.getCurrConfigBlockMutex.Lock()
	defer fake.getCurrConfigBlockMutex.Unlock()
	fake.GetCurrConfigBlockStub = nil
	fake.getCurrConfigBlockReturns = struct {
		result1 *common.Block
	}{result1}
}

func (fake *PeerOperations) GetCurrConfigBlockReturnsOnCall(i int, result1 *common.Block) {
	fake.getCurrConfigBlockMutex.Lock()
	defer fake.getCurrConfigBlockMutex.Unlock()
	fake.GetCurrConfigBlockStub = nil
	if fake.getCurrConfigBlockReturnsOnCall == nil {
		fake.getCurrConfigBlockReturnsOnCall = make(map[int]struct {
			result1 *common.Block
		})
	}
	fake.getCurrConfigBlockReturnsOnCall[i] = struct {
		result1 *common.Block
	}{result1}
}

func (fake *PeerOperations) GetLedger(arg1 string) ledger.PeerLedger {
	fake.getLedgerMutex.Lock()
	ret, specificReturn := fake.getLedgerReturnsOnCall[len(fake.getLedgerArgsForCall)]
	fake.getLedgerArgsForCall = append(fake.getLedgerArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("GetLedger", []interface{}{arg1})
	fake.getLedgerMutex.Unlock()
	if fake.GetLedgerStub != nil {
		return fake.GetLedgerStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.getLedgerReturns
	return fakeReturns.result1
}

func (fake *PeerOperations) GetLedgerCallCount() int {
	fake.getLedgerMutex.RLock()
	defer fake.getLedgerMutex.RUnlock()
	return len(fake.getLedgerArgsForCall)
}

func (fake *PeerOperations) GetLedgerCalls(stub func(string) ledger.PeerLedger) {
	fake.getLedgerMutex.Lock()
	defer fake.getLedgerMutex.Unlock()
	fake.GetLedgerStub = stub
}

func (fake *PeerOperations) GetLedgerArgsForCall(i int) string {
	fake.getLedgerMutex.RLock()
	defer fake.getLedgerMutex.RUnlock()
	argsForCall := fake.getLedgerArgsForCall[i]
	return argsForCall.arg1
}

func (fake *PeerOperations) GetLedgerReturns(result1 ledger.PeerLedger) {
	fake.getLedgerMutex.Lock()
	defer fake.getLedgerMutex.Unlock()
	fake.GetLedgerStub = nil
	fake.getLedgerReturns = struct {
		result1 ledger.PeerLedger
	}{result1}
}

func (fake *PeerOperations) GetLedgerReturnsOnCall(i int, result1 ledger.PeerLedger) {
	fake.getLedgerMutex.Lock()
	defer fake.getLedgerMutex.Unlock()
	fake.GetLedgerStub = nil
	if fake.getLedgerReturnsOnCall == nil {
		fake.getLedgerReturnsOnCall = make(map[int]struct {
			result1 ledger.PeerLedger
		})
	}
	fake.getLedgerReturnsOnCall[i] = struct {
		result1 ledger.PeerLedger
	}{result1}
}

func (fake *PeerOperations) GetMSPIDs(arg1 string) []string {
	fake.getMSPIDsMutex.Lock()
	ret, specificReturn := fake.getMSPIDsReturnsOnCall[len(fake.getMSPIDsArgsForCall)]
	fake.getMSPIDsArgsForCall = append(fake.getMSPIDsArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("GetMSPIDs", []interface{}{arg1})
	fake.getMSPIDsMutex.Unlock()
	if fake.GetMSPIDsStub != nil {
		return fake.GetMSPIDsStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.getMSPIDsReturns
	return fakeReturns.result1
}

func (fake *PeerOperations) GetMSPIDsCallCount() int {
	fake.getMSPIDsMutex.RLock()
	defer fake.getMSPIDsMutex.RUnlock()
	return len(fake.getMSPIDsArgsForCall)
}

func (fake *PeerOperations) GetMSPIDsCalls(stub func(string) []string) {
	fake.getMSPIDsMutex.Lock()
	defer fake.getMSPIDsMutex.Unlock()
	fake.GetMSPIDsStub = stub
}

func (fake *PeerOperations) GetMSPIDsArgsForCall(i int) string {
	fake.getMSPIDsMutex.RLock()
	defer fake.getMSPIDsMutex.RUnlock()
	argsForCall := fake.getMSPIDsArgsForCall[i]
	return argsForCall.arg1
}

func (fake *PeerOperations) GetMSPIDsReturns(result1 []string) {
	fake.getMSPIDsMutex.Lock()
	defer fake.getMSPIDsMutex.Unlock()
	fake.GetMSPIDsStub = nil
	fake.getMSPIDsReturns = struct {
		result1 []string
	}{result1}
}

func (fake *PeerOperations) GetMSPIDsReturnsOnCall(i int, result1 []string) {
	fake.getMSPIDsMutex.Lock()
	defer fake.getMSPIDsMutex.Unlock()
	fake.GetMSPIDsStub = nil
	if fake.getMSPIDsReturnsOnCall == nil {
		fake.getMSPIDsReturnsOnCall = make(map[int]struct {
			result1 []string
		})
	}
	fake.getMSPIDsReturnsOnCall[i] = struct {
		result1 []string
	}{result1}
}

func (fake *PeerOperations) GetPolicyManager(arg1 string) policies.Manager {
	fake.getPolicyManagerMutex.Lock()
	ret, specificReturn := fake.getPolicyManagerReturnsOnCall[len(fake.getPolicyManagerArgsForCall)]
	fake.getPolicyManagerArgsForCall = append(fake.getPolicyManagerArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("GetPolicyManager", []interface{}{arg1})
	fake.getPolicyManagerMutex.Unlock()
	if fake.GetPolicyManagerStub != nil {
		return fake.GetPolicyManagerStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.getPolicyManagerReturns
	return fakeReturns.result1
}

func (fake *PeerOperations) GetPolicyManagerCallCount() int {
	fake.getPolicyManagerMutex.RLock()
	defer fake.getPolicyManagerMutex.RUnlock()
	return len(fake.getPolicyManagerArgsForCall)
}

func (fake *PeerOperations) GetPolicyManagerCalls(stub func(string) policies.Manager) {
	fake.getPolicyManagerMutex.Lock()
	defer fake.getPolicyManagerMutex.Unlock()
	fake.GetPolicyManagerStub = stub
}

func (fake *PeerOperations) GetPolicyManagerArgsForCall(i int) string {
	fake.getPolicyManagerMutex.RLock()
	defer fake.getPolicyManagerMutex.RUnlock()
	argsForCall := fake.getPolicyManagerArgsForCall[i]
	return argsForCall.arg1
}

func (fake *PeerOperations) GetPolicyManagerReturns(result1 policies.Manager) {
	fake.getPolicyManagerMutex.Lock()
	defer fake.getPolicyManagerMutex.Unlock()
	fake.GetPolicyManagerStub = nil
	fake.getPolicyManagerReturns = struct {
		result1 policies.Manager
	}{result1}
}

func (fake *PeerOperations) GetPolicyManagerReturnsOnCall(i int, result1 policies.Manager) {
	fake.getPolicyManagerMutex.Lock()
	defer fake.getPolicyManagerMutex.Unlock()
	fake.GetPolicyManagerStub = nil
	if fake.getPolicyManagerReturnsOnCall == nil {
		fake.getPolicyManagerReturnsOnCall = make(map[int]struct {
			result1 policies.Manager
		})
	}
	fake.getPolicyManagerReturnsOnCall[i] = struct {
		result1 policies.Manager
	}{result1}
}

func (fake *PeerOperations) GetStableChannelConfig(arg1 string) channelconfig.Resources {
	fake.getStableChannelConfigMutex.Lock()
	ret, specificReturn := fake.getStableChannelConfigReturnsOnCall[len(fake.getStableChannelConfigArgsForCall)]
	fake.getStableChannelConfigArgsForCall = append(fake.getStableChannelConfigArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("GetStableChannelConfig", []interface{}{arg1})
	fake.getStableChannelConfigMutex.Unlock()
	if fake.GetStableChannelConfigStub != nil {
		return fake.GetStableChannelConfigStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.getStableChannelConfigReturns
	return fakeReturns.result1
}

func (fake *PeerOperations) GetStableChannelConfigCallCount() int {
	fake.getStableChannelConfigMutex.RLock()
	defer fake.getStableChannelConfigMutex.RUnlock()
	return len(fake.getStableChannelConfigArgsForCall)
}

func (fake *PeerOperations) GetStableChannelConfigCalls(stub func(string) channelconfig.Resources) {
	fake.getStableChannelConfigMutex.Lock()
	defer fake.getStableChannelConfigMutex.Unlock()
	fake.GetStableChannelConfigStub = stub
}

func (fake *PeerOperations) GetStableChannelConfigArgsForCall(i int) string {
	fake.getStableChannelConfigMutex.RLock()
	defer fake.getStableChannelConfigMutex.RUnlock()
	argsForCall := fake.getStableChannelConfigArgsForCall[i]
	return argsForCall.arg1
}

func (fake *PeerOperations) GetStableChannelConfigReturns(result1 channelconfig.Resources) {
	fake.getStableChannelConfigMutex.Lock()
	defer fake.getStableChannelConfigMutex.Unlock()
	fake.GetStableChannelConfigStub = nil
	fake.getStableChannelConfigReturns = struct {
		result1 channelconfig.Resources
	}{result1}
}

func (fake *PeerOperations) GetStableChannelConfigReturnsOnCall(i int, result1 channelconfig.Resources) {
	fake.getStableChannelConfigMutex.Lock()
	defer fake.getStableChannelConfigMutex.Unlock()
	fake.GetStableChannelConfigStub = nil
	if fake.getStableChannelConfigReturnsOnCall == nil {
		fake.getStableChannelConfigReturnsOnCall = make(map[int]struct {
			result1 channelconfig.Resources
		})
	}
	fake.getStableChannelConfigReturnsOnCall[i] = struct {
		result1 channelconfig.Resources
	}{result1}
}

func (fake *PeerOperations) InitChain(arg1 string) {
	fake.initChainMutex.Lock()
	fake.initChainArgsForCall = append(fake.initChainArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("InitChain", []interface{}{arg1})
	fake.initChainMutex.Unlock()
	if fake.InitChainStub != nil {
		fake.InitChainStub(arg1)
	}
}

func (fake *PeerOperations) InitChainCallCount() int {
	fake.initChainMutex.RLock()
	defer fake.initChainMutex.RUnlock()
	return len(fake.initChainArgsForCall)
}

func (fake *PeerOperations) InitChainCalls(stub func(string)) {
	fake.initChainMutex.Lock()
	defer fake.initChainMutex.Unlock()
	fake.InitChainStub = stub
}

func (fake *PeerOperations) InitChainArgsForCall(i int) string {
	fake.initChainMutex.RLock()
	defer fake.initChainMutex.RUnlock()
	argsForCall := fake.initChainArgsForCall[i]
	return argsForCall.arg1
}

func (fake *PeerOperations) Initialize(arg1 func(string), arg2 sysccprovider.SystemChaincodeProvider, arg3 plugin.Mapper, arg4 *platforms.Registry, arg5 ledger.DeployedChaincodeInfoProvider, arg6 ledger.MembershipInfoProvider, arg7 metrics.Provider, arg8 plugindispatcher.LifecycleResources, arg9 plugindispatcher.CollectionAndLifecycleResources, arg10 *ledger.Config, arg11 int) {
	fake.initializeMutex.Lock()
	fake.initializeArgsForCall = append(fake.initializeArgsForCall, struct {
		arg1  func(string)
		arg2  sysccprovider.SystemChaincodeProvider
		arg3  plugin.Mapper
		arg4  *platforms.Registry
		arg5  ledger.DeployedChaincodeInfoProvider
		arg6  ledger.MembershipInfoProvider
		arg7  metrics.Provider
		arg8  plugindispatcher.LifecycleResources
		arg9  plugindispatcher.CollectionAndLifecycleResources
		arg10 *ledger.Config
		arg11 int
	}{arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8, arg9, arg10, arg11})
	fake.recordInvocation("Initialize", []interface{}{arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8, arg9, arg10, arg11})
	fake.initializeMutex.Unlock()
	if fake.InitializeStub != nil {
		fake.InitializeStub(arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8, arg9, arg10, arg11)
	}
}

func (fake *PeerOperations) InitializeCallCount() int {
	fake.initializeMutex.RLock()
	defer fake.initializeMutex.RUnlock()
	return len(fake.initializeArgsForCall)
}

func (fake *PeerOperations) InitializeCalls(stub func(func(string), sysccprovider.SystemChaincodeProvider, plugin.Mapper, *platforms.Registry, ledger.DeployedChaincodeInfoProvider, ledger.MembershipInfoProvider, metrics.Provider, plugindispatcher.LifecycleResources, plugindispatcher.CollectionAndLifecycleResources, *ledger.Config, int)) {
	fake.initializeMutex.Lock()
	defer fake.initializeMutex.Unlock()
	fake.InitializeStub = stub
}

func (fake *PeerOperations) InitializeArgsForCall(i int) (func(string), sysccprovider.SystemChaincodeProvider, plugin.Mapper, *platforms.Registry, ledger.DeployedChaincodeInfoProvider, ledger.MembershipInfoProvider, metrics.Provider, plugindispatcher.LifecycleResources, plugindispatcher.CollectionAndLifecycleResources, *ledger.Config, int) {
	fake.initializeMutex.RLock()
	defer fake.initializeMutex.RUnlock()
	argsForCall := fake.initializeArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4, argsForCall.arg5, argsForCall.arg6, argsForCall.arg7, argsForCall.arg8, argsForCall.arg9, argsForCall.arg10, argsForCall.arg11
}

func (fake *PeerOperations) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.createChainFromBlockMutex.RLock()
	defer fake.createChainFromBlockMutex.RUnlock()
	fake.getChannelConfigMutex.RLock()
	defer fake.getChannelConfigMutex.RUnlock()
	fake.getChannelsInfoMutex.RLock()
	defer fake.getChannelsInfoMutex.RUnlock()
	fake.getCurrConfigBlockMutex.RLock()
	defer fake.getCurrConfigBlockMutex.RUnlock()
	fake.getLedgerMutex.RLock()
	defer fake.getLedgerMutex.RUnlock()
	fake.getMSPIDsMutex.RLock()
	defer fake.getMSPIDsMutex.RUnlock()
	fake.getPolicyManagerMutex.RLock()
	defer fake.getPolicyManagerMutex.RUnlock()
	fake.getStableChannelConfigMutex.RLock()
	defer fake.getStableChannelConfigMutex.RUnlock()
	fake.initChainMutex.RLock()
	defer fake.initChainMutex.RUnlock()
	fake.initializeMutex.RLock()
	defer fake.initializeMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *PeerOperations) recordInvocation(key string, args []interface{}) {
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
