/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/






package cscc

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/bccsp/factory"
	"github.com/mcc-github/blockchain/common/channelconfig"
	"github.com/mcc-github/blockchain/common/config"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/core/aclmgmt"
	"github.com/mcc-github/blockchain/core/aclmgmt/resources"
	"github.com/mcc-github/blockchain/core/chaincode/shim"
	"github.com/mcc-github/blockchain/core/committer/txvalidator/v20/plugindispatcher"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/ledger/util"
	"github.com/mcc-github/blockchain/core/peer"
	"github.com/mcc-github/blockchain/core/policy"
	"github.com/mcc-github/blockchain/protos/common"
	pb "github.com/mcc-github/blockchain/protos/peer"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/pkg/errors"
)



func New(
	aclProvider aclmgmt.ACLProvider,
	deployedCCInfoProvider ledger.DeployedChaincodeInfoProvider,
	lr plugindispatcher.LifecycleResources,
	nr plugindispatcher.CollectionAndLifecycleResources,
	policyChecker policy.PolicyChecker,
	p *peer.Peer,
) *PeerConfiger {
	return &PeerConfiger{
		policyChecker:          policyChecker,
		configMgr:              peer.NewConfigSupport(p),
		aclProvider:            aclProvider,
		deployedCCInfoProvider: deployedCCInfoProvider,
		legacyLifecycle:        lr,
		newLifecycle:           nr,
		peer:                   p,
	}
}

func (e *PeerConfiger) Name() string              { return "cscc" }
func (e *PeerConfiger) Chaincode() shim.Chaincode { return e }




type PeerConfiger struct {
	policyChecker          policy.PolicyChecker
	configMgr              config.Manager
	aclProvider            aclmgmt.ACLProvider
	deployedCCInfoProvider ledger.DeployedChaincodeInfoProvider
	legacyLifecycle        plugindispatcher.LifecycleResources
	newLifecycle           plugindispatcher.CollectionAndLifecycleResources
	peer                   *peer.Peer
}

var cnflogger = flogging.MustGetLogger("cscc")


const (
	JoinChain      string = "JoinChain"
	GetConfigBlock string = "GetConfigBlock"
	GetChannels    string = "GetChannels"
)


func (e *PeerConfiger) Init(stub shim.ChaincodeStubInterface) pb.Response {
	cnflogger.Info("Init CSCC")
	return shim.Success(nil)
}











func (e *PeerConfiger) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	args := stub.GetArgs()

	if len(args) < 1 {
		return shim.Error(fmt.Sprintf("Incorrect number of arguments, %d", len(args)))
	}

	fname := string(args[0])

	if fname != GetChannels && len(args) < 2 {
		return shim.Error(fmt.Sprintf("Incorrect number of arguments, %d", len(args)))
	}

	cnflogger.Debugf("Invoke function: %s", fname)

	
	
	sp, err := stub.GetSignedProposal()
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed getting signed proposal from stub: [%s]", err))
	}

	name, err := InvokedChaincodeName(sp.ProposalBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Could not identify the called chaincode: [%s]", err))
	}

	if name != e.Name() {
		return shim.Error(fmt.Sprintf("Cannot invoke CSCC from another chaincode, original invocation for '%s'", name))
	}

	return e.InvokeNoShim(args, sp)
}

func InvokedChaincodeName(proposalBytes []byte) (string, error) {
	proposal := &pb.Proposal{}
	err := proto.Unmarshal(proposalBytes, proposal)
	if err != nil {
		return "", errors.WithMessage(err, "could not unmarshal proposal")
	}

	proposalPayload := &pb.ChaincodeProposalPayload{}
	err = proto.Unmarshal(proposal.Payload, proposalPayload)
	if err != nil {
		return "", errors.WithMessage(err, "could not unmarshal chaincode proposal payload")
	}

	cis := &pb.ChaincodeInvocationSpec{}
	err = proto.Unmarshal(proposalPayload.Input, cis)
	if err != nil {
		return "", errors.WithMessage(err, "could not unmarshal chaincode invocation spec")
	}

	if cis.ChaincodeSpec == nil {
		return "", errors.Errorf("chaincode spec is nil")
	}

	if cis.ChaincodeSpec.ChaincodeId == nil {
		return "", errors.Errorf("chaincode id is nil")
	}

	return cis.ChaincodeSpec.ChaincodeId.Name, nil
}

func (e *PeerConfiger) InvokeNoShim(args [][]byte, sp *pb.SignedProposal) pb.Response {
	var err error
	fname := string(args[0])

	switch fname {
	case JoinChain:
		if args[1] == nil {
			return shim.Error("Cannot join the channel <nil> configuration block provided")
		}

		block, err := protoutil.UnmarshalBlock(args[1])
		if err != nil {
			return shim.Error(fmt.Sprintf("Failed to reconstruct the genesis block, %s", err))
		}

		cid, err := protoutil.GetChainIDFromBlock(block)
		if err != nil {
			return shim.Error(fmt.Sprintf("\"JoinChain\" request failed to extract "+
				"channel id from the block due to [%s]", err))
		}

		
		if err := validateConfigBlock(block); err != nil {
			return shim.Error(fmt.Sprintf("\"JoinChain\" for chainID = %s failed because of validation "+
				"of configuration block, because of %s", cid, err))
		}

		
		if err = e.aclProvider.CheckACL(resources.Cscc_JoinChain, "", sp); err != nil {
			return shim.Error(fmt.Sprintf("access denied for [%s][%s]: [%s]", fname, cid, err))
		}

		
		
		txsFilter := util.TxValidationFlags(block.Metadata.Metadata[common.BlockMetadataIndex_TRANSACTIONS_FILTER])
		if len(txsFilter) == 0 {
			
			txsFilter = util.NewTxValidationFlagsSetValue(len(block.Data.Data), pb.TxValidationCode_VALID)
			block.Metadata.Metadata[common.BlockMetadataIndex_TRANSACTIONS_FILTER] = txsFilter
		}

		return e.joinChain(cid, block, e.deployedCCInfoProvider, e.legacyLifecycle, e.newLifecycle)
	case GetConfigBlock:
		
		if err = e.aclProvider.CheckACL(resources.Cscc_GetConfigBlock, string(args[1]), sp); err != nil {
			return shim.Error(fmt.Sprintf("access denied for [%s][%s]: %s", fname, args[1], err))
		}

		return e.getConfigBlock(args[1])
	case GetChannels:
		
		if err = e.aclProvider.CheckACL(resources.Cscc_GetChannels, "", sp); err != nil {
			return shim.Error(fmt.Sprintf("access denied for [%s]: %s", fname, err))
		}

		return e.getChannels()

	}
	return shim.Error(fmt.Sprintf("Requested function %s not found.", fname))
}


func validateConfigBlock(block *common.Block) error {
	envelopeConfig, err := protoutil.ExtractEnvelope(block, 0)
	if err != nil {
		return errors.Errorf("Failed to %s", err)
	}

	configEnv := &common.ConfigEnvelope{}
	_, err = protoutil.UnmarshalEnvelopeOfType(envelopeConfig, common.HeaderType_CONFIG, configEnv)
	if err != nil {
		return errors.Errorf("Bad configuration envelope: %s", err)
	}

	if configEnv.Config == nil {
		return errors.New("Nil config envelope Config")
	}

	if configEnv.Config.ChannelGroup == nil {
		return errors.New("Nil channel group")
	}

	if configEnv.Config.ChannelGroup.Groups == nil {
		return errors.New("No channel configuration groups are available")
	}

	_, exists := configEnv.Config.ChannelGroup.Groups[channelconfig.ApplicationGroupKey]
	if !exists {
		return errors.Errorf("Invalid configuration block, missing %s "+
			"configuration group", channelconfig.ApplicationGroupKey)
	}

	
	if err = channelconfig.ValidateCapabilities(block, factory.GetDefault()); err != nil {
		return errors.Errorf("Failed capabilities check: [%s]", err)
	}

	return nil
}




func (e *PeerConfiger) joinChain(
	chainID string,
	block *common.Block,
	deployedCCInfoProvider ledger.DeployedChaincodeInfoProvider,
	lr plugindispatcher.LifecycleResources,
	nr plugindispatcher.CollectionAndLifecycleResources,
) pb.Response {
	if err := e.peer.CreateChannel(chainID, block, deployedCCInfoProvider, lr, nr); err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}



func (e *PeerConfiger) getConfigBlock(chainID []byte) pb.Response {
	if chainID == nil {
		return shim.Error("ChainID must not be nil.")
	}

	channel := e.peer.Channel(string(chainID))
	if channel == nil {
		return shim.Error(fmt.Sprintf("Unknown chain ID, %s", string(chainID)))
	}
	block, err := peer.ConfigBlockFromLedger(channel.Ledger())
	if err != nil {
		return shim.Error(err.Error())
	}

	blockBytes, err := protoutil.Marshal(block)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(blockBytes)
}

func (e *PeerConfiger) supportByType(chainID []byte, env *common.Envelope) (config.Config, error) {
	payload := &common.Payload{}

	if err := proto.Unmarshal(env.Payload, payload); err != nil {
		return nil, errors.Errorf("failed unmarshaling payload: %v", err)
	}

	channelHdr := &common.ChannelHeader{}
	if err := proto.Unmarshal(payload.Header.ChannelHeader, channelHdr); err != nil {
		return nil, errors.Errorf("failed unmarshaling payload header: %v", err)
	}

	switch common.HeaderType(channelHdr.Type) {
	case common.HeaderType_CONFIG_UPDATE:
		return e.configMgr.GetChannelConfig(string(chainID)), nil
	}
	return nil, errors.Errorf("invalid payload header type: %d", channelHdr.Type)
}


func (e *PeerConfiger) getChannels() pb.Response {
	channelInfoArray := e.peer.GetChannelsInfo()

	
	cqr := &pb.ChannelQueryResponse{Channels: channelInfoArray}

	cqrbytes, err := proto.Marshal(cqr)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(cqrbytes)
}
