/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/






package cscc

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/common/channelconfig"
	"github.com/mcc-github/blockchain/common/config"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/core/aclmgmt"
	"github.com/mcc-github/blockchain/core/aclmgmt/resources"
	"github.com/mcc-github/blockchain/core/chaincode/shim"
	"github.com/mcc-github/blockchain/core/common/ccprovider"
	"github.com/mcc-github/blockchain/core/common/sysccprovider"
	"github.com/mcc-github/blockchain/core/ledger/util"
	"github.com/mcc-github/blockchain/core/peer"
	"github.com/mcc-github/blockchain/core/policy"
	"github.com/mcc-github/blockchain/events/producer"
	"github.com/mcc-github/blockchain/msp/mgmt"
	"github.com/mcc-github/blockchain/protos/common"
	pb "github.com/mcc-github/blockchain/protos/peer"
	"github.com/mcc-github/blockchain/protos/utils"
	"github.com/pkg/errors"
)



func New(ccp ccprovider.ChaincodeProvider, sccp sysccprovider.SystemChaincodeProvider, aclProvider aclmgmt.ACLProvider) *PeerConfiger {
	return &PeerConfiger{
		policyChecker: policy.NewPolicyChecker(
			peer.NewChannelPolicyManagerGetter(),
			mgmt.GetLocalMSP(),
			mgmt.NewLocalMSPPrincipalGetter(),
		),
		configMgr:   peer.NewConfigSupport(),
		ccp:         ccp,
		sccp:        sccp,
		aclProvider: aclProvider,
	}
}

func (e *PeerConfiger) Name() string              { return "cscc" }
func (e *PeerConfiger) Path() string              { return "github.com/mcc-github/blockchain/core/scc/cscc" }
func (e *PeerConfiger) InitArgs() [][]byte        { return nil }
func (e *PeerConfiger) Chaincode() shim.Chaincode { return e }
func (e *PeerConfiger) InvokableExternal() bool   { return true }
func (e *PeerConfiger) InvokableCC2CC() bool      { return false }
func (e *PeerConfiger) Enabled() bool             { return true }




type PeerConfiger struct {
	policyChecker policy.PolicyChecker
	configMgr     config.Manager
	ccp           ccprovider.ChaincodeProvider
	sccp          sysccprovider.SystemChaincodeProvider
	aclProvider   aclmgmt.ACLProvider
}

var cnflogger = flogging.MustGetLogger("cscc")


const (
	JoinChain                string = "JoinChain"
	GetConfigBlock           string = "GetConfigBlock"
	GetChannels              string = "GetChannels"
	GetConfigTree            string = "GetConfigTree"
	SimulateConfigTreeUpdate string = "SimulateConfigTreeUpdate"
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

	return e.InvokeNoShim(args, sp)
}

func (e *PeerConfiger) InvokeNoShim(args [][]byte, sp *pb.SignedProposal) pb.Response {
	var err error
	fname := string(args[0])

	switch fname {
	case JoinChain:
		if args[1] == nil {
			return shim.Error("Cannot join the channel <nil> configuration block provided")
		}

		block, err := utils.GetBlockFromBlockBytes(args[1])
		if err != nil {
			return shim.Error(fmt.Sprintf("Failed to reconstruct the genesis block, %s", err))
		}

		cid, err := utils.GetChainIDFromBlock(block)
		if err != nil {
			return shim.Error(fmt.Sprintf("\"JoinChain\" request failed to extract "+
				"channel id from the block due to [%s]", err))
		}

		if err := validateConfigBlock(block); err != nil {
			return shim.Error(fmt.Sprintf("\"JoinChain\" for chainID = %s failed because of validation "+
				"of configuration block, because of %s", cid, err))
		}

		
		
		if err = e.policyChecker.CheckPolicyNoChannel(mgmt.Admins, sp); err != nil {
			return shim.Error(fmt.Sprintf("access denied for [%s][%s]: [%s]", fname, cid, err))
		}

		
		
		txsFilter := util.TxValidationFlags(block.Metadata.Metadata[common.BlockMetadataIndex_TRANSACTIONS_FILTER])
		if len(txsFilter) == 0 {
			
			txsFilter = util.NewTxValidationFlagsSetValue(len(block.Data.Data), pb.TxValidationCode_VALID)
			block.Metadata.Metadata[common.BlockMetadataIndex_TRANSACTIONS_FILTER] = txsFilter
		}

		return joinChain(cid, block, e.ccp, e.sccp)
	case GetConfigBlock:
		
		if err = e.aclProvider.CheckACL(resources.Cscc_GetConfigBlock, string(args[1]), sp); err != nil {
			return shim.Error(fmt.Sprintf("access denied for [%s][%s]: %s", fname, args[1], err))
		}

		return getConfigBlock(args[1])
	case GetConfigTree:
		
		if err = e.aclProvider.CheckACL(resources.Cscc_GetConfigTree, string(args[1]), sp); err != nil {
			return shim.Error(fmt.Sprintf("access denied for [%s][%s]: %s", fname, args[1], err))
		}

		return e.getConfigTree(args[1])
	case SimulateConfigTreeUpdate:
		
		if err = e.aclProvider.CheckACL(resources.Cscc_SimulateConfigTreeUpdate, string(args[1]), sp); err != nil {
			return shim.Error(fmt.Sprintf("access denied for [%s][%s]: %s", fname, args[1], err))
		}
		return e.simulateConfigTreeUpdate(args[1], args[2])
	case GetChannels:
		
		
		if err = e.policyChecker.CheckPolicyNoChannel(mgmt.Members, sp); err != nil {
			return shim.Error(fmt.Sprintf("access denied for [%s]: %s", fname, err))
		}

		return getChannels()

	}
	return shim.Error(fmt.Sprintf("Requested function %s not found.", fname))
}


func validateConfigBlock(block *common.Block) error {
	envelopeConfig, err := utils.ExtractEnvelope(block, 0)
	if err != nil {
		return errors.Errorf("Failed to %s", err)
	}

	configEnv := &common.ConfigEnvelope{}
	_, err = utils.UnmarshalEnvelopeOfType(envelopeConfig, common.HeaderType_CONFIG, configEnv)
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

	return nil
}




func joinChain(chainID string, block *common.Block, ccp ccprovider.ChaincodeProvider, sccp sysccprovider.SystemChaincodeProvider) pb.Response {
	if err := peer.CreateChainFromBlock(block, ccp, sccp); err != nil {
		return shim.Error(err.Error())
	}

	peer.InitChain(chainID)

	bevent, _, _, err := producer.CreateBlockEvents(block)
	if err != nil {
		cnflogger.Errorf("Error processing block events for block number [%d]: %s", block.Header.Number, err)
	} else {
		if err := producer.Send(bevent); err != nil {
			cnflogger.Errorf("Channel [%s] Error sending block event for block number [%d]: %s", chainID, block.Header.Number, err)
		}
	}

	return shim.Success(nil)
}



func getConfigBlock(chainID []byte) pb.Response {
	if chainID == nil {
		return shim.Error("ChainID must not be nil.")
	}
	block := peer.GetCurrConfigBlock(string(chainID))
	if block == nil {
		return shim.Error(fmt.Sprintf("Unknown chain ID, %s", string(chainID)))
	}
	blockBytes, err := utils.Marshal(block)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(blockBytes)
}



func (e *PeerConfiger) getConfigTree(chainID []byte) pb.Response {
	if chainID == nil {
		return shim.Error("Chain ID must not be nil")
	}
	channelCfg := e.configMgr.GetChannelConfig(string(chainID)).ConfigProto()
	if channelCfg == nil {
		return shim.Error(fmt.Sprintf("Unknown chain ID, %s", string(chainID)))
	}
	agCfg := &pb.ConfigTree{ChannelConfig: channelCfg}
	configBytes, err := utils.Marshal(agCfg)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(configBytes)
}

func (e *PeerConfiger) simulateConfigTreeUpdate(chainID []byte, envb []byte) pb.Response {
	if chainID == nil {
		return shim.Error("Chain ID must not be nil")
	}
	if envb == nil {
		return shim.Error("Config delta bytes must not be nil")
	}
	env := &common.Envelope{}
	err := proto.Unmarshal(envb, env)
	if err != nil {
		return shim.Error(err.Error())
	}
	cfg, err := supportByType(e, chainID, env)
	if err != nil {
		return shim.Error(err.Error())
	}
	_, err = cfg.ProposeConfigUpdate(env)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success([]byte("Simulation is successful"))
}

func supportByType(pc *PeerConfiger, chainID []byte, env *common.Envelope) (config.Config, error) {
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
		return pc.configMgr.GetChannelConfig(string(chainID)), nil
	}
	return nil, errors.Errorf("invalid payload header type: %d", channelHdr.Type)
}


func getChannels() pb.Response {
	channelInfoArray := peer.GetChannelsInfo()

	
	cqr := &pb.ChannelQueryResponse{Channels: channelInfoArray}

	cqrbytes, err := proto.Marshal(cqr)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(cqrbytes)
}
