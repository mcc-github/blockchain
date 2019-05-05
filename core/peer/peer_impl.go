





package peer

import (
	"github.com/mcc-github/blockchain/common/channelconfig"
	"github.com/mcc-github/blockchain/common/metrics"
	"github.com/mcc-github/blockchain/common/policies"
	"github.com/mcc-github/blockchain/core/chaincode/platforms"
	"github.com/mcc-github/blockchain/core/committer/txvalidator/plugin"
	"github.com/mcc-github/blockchain/core/committer/txvalidator/v20/plugindispatcher"
	"github.com/mcc-github/blockchain/core/common/sysccprovider"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/protos/common"
	pb "github.com/mcc-github/blockchain/protos/peer"
)




type Operations interface {
	CreateChainFromBlock(cb *common.Block, sccp sysccprovider.SystemChaincodeProvider, deployedCCInfoProvider ledger.DeployedChaincodeInfoProvider, lr plugindispatcher.LifecycleResources, nr plugindispatcher.CollectionAndLifecycleResources) error
	GetChannelConfig(cid string) channelconfig.Resources
	GetChannelsInfo() []*pb.ChannelInfo
	GetStableChannelConfig(cid string) channelconfig.Resources
	GetCurrConfigBlock(cid string) *common.Block
	GetLedger(cid string) ledger.PeerLedger
	GetMSPIDs(cid string) []string
	GetPolicyManager(cid string) policies.Manager
	InitChain(cid string)
	Initialize(
		init func(string),
		sccp sysccprovider.SystemChaincodeProvider,
		pm plugin.Mapper,
		pr *platforms.Registry,
		deployedCCInfoProvider ledger.DeployedChaincodeInfoProvider,
		membershipProvider ledger.MembershipInfoProvider,
		metricsProvider metrics.Provider,
		lr plugindispatcher.LifecycleResources,
		nr plugindispatcher.CollectionAndLifecycleResources,
		ledgerConfig *ledger.Config,
		nWorkers int,
	)
}

type peerImpl struct {
	createChainFromBlock   func(cb *common.Block, sccp sysccprovider.SystemChaincodeProvider, deployedCCInfoProvider ledger.DeployedChaincodeInfoProvider, lr plugindispatcher.LifecycleResources, nr plugindispatcher.CollectionAndLifecycleResources) error
	getChannelConfig       func(cid string) channelconfig.Resources
	getChannelsInfo        func() []*pb.ChannelInfo
	getStableChannelConfig func(cid string) channelconfig.Resources
	getCurrConfigBlock     func(cid string) *common.Block
	getLedger              func(cid string) ledger.PeerLedger
	getMSPIDs              func(cid string) []string
	getPolicyManager       func(cid string) policies.Manager
	initChain              func(cid string)
	initialize             func(
		init func(string),
		sccp sysccprovider.SystemChaincodeProvider,
		mapper plugin.Mapper,
		pr *platforms.Registry,
		deployedCCInfoProvider ledger.DeployedChaincodeInfoProvider,
		membershipProvider ledger.MembershipInfoProvider,
		metricsProvider metrics.Provider,
		lr plugindispatcher.LifecycleResources,
		nr plugindispatcher.CollectionAndLifecycleResources,
		ledgerConfig *ledger.Config,
		nWorkers int,
	)
}



var Default Operations = &peerImpl{
	createChainFromBlock:   CreateChainFromBlock,
	getChannelConfig:       GetChannelConfig,
	getChannelsInfo:        GetChannelsInfo,
	getStableChannelConfig: GetStableChannelConfig,
	getCurrConfigBlock:     GetCurrConfigBlock,
	getLedger:              GetLedger,
	getMSPIDs:              GetMSPIDs,
	getPolicyManager:       GetPolicyManager,
	initChain:              InitChain,
	initialize:             Initialize,
}

var DefaultSupport Support = &supportImpl{operations: Default}

func (p *peerImpl) CreateChainFromBlock(cb *common.Block, sccp sysccprovider.SystemChaincodeProvider, deployedCCInfoProvider ledger.DeployedChaincodeInfoProvider, lr plugindispatcher.LifecycleResources, nr plugindispatcher.CollectionAndLifecycleResources) error {
	return p.createChainFromBlock(cb, sccp, deployedCCInfoProvider, lr, nr)
}
func (p *peerImpl) GetChannelConfig(cid string) channelconfig.Resources {
	return p.getChannelConfig(cid)
}
func (p *peerImpl) GetChannelsInfo() []*pb.ChannelInfo { return p.getChannelsInfo() }
func (p *peerImpl) GetStableChannelConfig(cid string) channelconfig.Resources {
	return p.getStableChannelConfig(cid)
}
func (p *peerImpl) GetCurrConfigBlock(cid string) *common.Block  { return p.getCurrConfigBlock(cid) }
func (p *peerImpl) GetLedger(cid string) ledger.PeerLedger       { return p.getLedger(cid) }
func (p *peerImpl) GetMSPIDs(cid string) []string                { return p.getMSPIDs(cid) }
func (p *peerImpl) GetPolicyManager(cid string) policies.Manager { return p.getPolicyManager(cid) }
func (p *peerImpl) InitChain(cid string)                         { p.initChain(cid) }
func (p *peerImpl) Initialize(
	init func(string),
	sccp sysccprovider.SystemChaincodeProvider,
	mapper plugin.Mapper,
	pr *platforms.Registry,
	deployedCCInfoProvider ledger.DeployedChaincodeInfoProvider,
	membershipProvider ledger.MembershipInfoProvider,
	metricsProvider metrics.Provider,
	lr plugindispatcher.LifecycleResources,
	nr plugindispatcher.CollectionAndLifecycleResources,
	ledgerConfig *ledger.Config,
	nWorkers int,
) {
	p.initialize(
		init,
		sccp,
		mapper,
		pr,
		deployedCCInfoProvider,
		membershipProvider,
		metricsProvider,
		lr,
		nr,
		ledgerConfig,
		nWorkers,
	)
}
