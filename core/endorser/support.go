/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package endorser

import (
	"fmt"

	"github.com/mcc-github/blockchain/common/channelconfig"
	"github.com/mcc-github/blockchain/core/aclmgmt"
	"github.com/mcc-github/blockchain/core/aclmgmt/resources"
	"github.com/mcc-github/blockchain/core/chaincode"
	"github.com/mcc-github/blockchain/core/common/ccprovider"
	"github.com/mcc-github/blockchain/core/handlers/decoration"
	. "github.com/mcc-github/blockchain/core/handlers/endorsement/api/identities"
	"github.com/mcc-github/blockchain/core/handlers/library"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/scc"
	"github.com/mcc-github/blockchain/internal/pkg/identity"
	"github.com/mcc-github/blockchain/protos/common"
	pb "github.com/mcc-github/blockchain/protos/peer"
	"github.com/pkg/errors"
)



type PeerOperations interface {
	GetApplicationConfig(cid string) (channelconfig.Application, bool)
	GetLedger(cid string) ledger.PeerLedger
}



type SupportImpl struct {
	*PluginEndorser
	identity.SignerSerializer
	Peer             PeerOperations
	ChaincodeSupport *chaincode.ChaincodeSupport
	SysCCProvider    *scc.Provider
	ACLProvider      aclmgmt.ACLProvider
}

func (s *SupportImpl) NewQueryCreator(channel string) (QueryCreator, error) {
	lgr := s.Peer.GetLedger(channel)
	if lgr == nil {
		return nil, errors.Errorf("channel %s doesn't exist", channel)
	}
	return lgr, nil
}

func (s *SupportImpl) SigningIdentityForRequest(*pb.SignedProposal) (SigningIdentity, error) {
	return s.SignerSerializer, nil
}



func (s *SupportImpl) IsSysCCAndNotInvokableExternal(name string) bool {
	return s.SysCCProvider.IsSysCCAndNotInvokableExternal(name)
}




func (s *SupportImpl) GetTxSimulator(ledgername string, txid string) (ledger.TxSimulator, error) {
	lgr := s.Peer.GetLedger(ledgername)
	if lgr == nil {
		return nil, errors.Errorf("Channel does not exist: %s", ledgername)
	}
	return lgr.NewTxSimulator(txid)
}



func (s *SupportImpl) GetHistoryQueryExecutor(ledgername string) (ledger.HistoryQueryExecutor, error) {
	lgr := s.Peer.GetLedger(ledgername)
	if lgr == nil {
		return nil, errors.Errorf("Channel does not exist: %s", ledgername)
	}
	return lgr.NewHistoryQueryExecutor()
}


func (s *SupportImpl) GetTransactionByID(chid, txID string) (*pb.ProcessedTransaction, error) {
	lgr := s.Peer.GetLedger(chid)
	if lgr == nil {
		return nil, errors.Errorf("failed to look up the ledger for Channel %s", chid)
	}
	tx, err := lgr.GetTransactionByID(txID)
	if err != nil {
		return nil, errors.WithMessage(err, "GetTransactionByID failed")
	}
	return tx, nil
}


func (s *SupportImpl) GetLedgerHeight(channelID string) (uint64, error) {
	lgr := s.Peer.GetLedger(channelID)
	if lgr == nil {
		return 0, errors.Errorf("failed to look up the ledger for Channel %s", channelID)
	}

	info, err := lgr.GetBlockchainInfo()
	if err != nil {
		return 0, errors.Wrap(err, fmt.Sprintf("failed to obtain information for Channel %s", channelID))
	}

	return info.Height, nil
}



func (s *SupportImpl) IsSysCC(name string) bool {
	return s.SysCCProvider.IsSysCC(name)
}


func (s *SupportImpl) GetChaincodeDeploymentSpecFS(cds *pb.ChaincodeDeploymentSpec) (*pb.ChaincodeDeploymentSpec, error) {
	ccpack, err := ccprovider.GetChaincodeFromFS(cds.ChaincodeSpec.ChaincodeId.Name, cds.ChaincodeSpec.ChaincodeId.Version)
	if err != nil {
		return nil, errors.Wrapf(err, "could not get chaincode from fs")
	}

	return ccpack.GetDepSpec(), nil
}


func (s *SupportImpl) ExecuteLegacyInit(txParams *ccprovider.TransactionParams, cid, name, version, txid string, signedProp *pb.SignedProposal, prop *pb.Proposal, cds *pb.ChaincodeDeploymentSpec) (*pb.Response, *pb.ChaincodeEvent, error) {
	cccid := &ccprovider.CCContext{
		Name:    name,
		Version: version,
	}

	return s.ChaincodeSupport.ExecuteLegacyInit(txParams, cccid, cds)
}


func (s *SupportImpl) Execute(txParams *ccprovider.TransactionParams, cid, name, txid string, idBytes []byte, requiresInit bool, signedProp *pb.SignedProposal, prop *pb.Proposal, input *pb.ChaincodeInput) (*pb.Response, *pb.ChaincodeEvent, error) {
	cccid := &ccprovider.CCContext{
		Name:         name,
		InitRequired: requiresInit,
		ID:           idBytes,
	}

	
	decorators := library.InitRegistry(library.Config{}).Lookup(library.Decoration).([]decoration.Decorator)
	input.Decorations = make(map[string][]byte)
	input = decoration.Apply(prop, input, decorators...)
	txParams.ProposalDecorations = input.Decorations

	return s.ChaincodeSupport.Execute(txParams, cccid, input)
}


func (s *SupportImpl) GetChaincodeDefinition(channelID, chaincodeName string, txsim ledger.QueryExecutor) (ccprovider.ChaincodeDefinition, error) {
	return s.ChaincodeSupport.Lifecycle.ChaincodeDefinition(channelID, chaincodeName, txsim)
}



func (s *SupportImpl) CheckACL(signedProp *pb.SignedProposal, chdr *common.ChannelHeader, shdr *common.SignatureHeader, hdrext *pb.ChaincodeHeaderExtension) error {
	return s.ACLProvider.CheckACL(resources.Peer_Propose, chdr.ChannelId, signedProp)
}




func (s *SupportImpl) CheckInstantiationPolicy(name, version string, cd ccprovider.ChaincodeDefinition) error {
	if cData, ok := cd.(*ccprovider.ChaincodeData); ok {
		return ccprovider.CheckInstantiationPolicy(name, version, cData)
	}
	return nil
}



func (s *SupportImpl) GetApplicationConfig(cid string) (channelconfig.Application, bool) {
	return s.Peer.GetApplicationConfig(cid)
}


func (s *SupportImpl) GetDeployedCCInfoProvider() ledger.DeployedChaincodeInfoProvider {
	return s.ChaincodeSupport.DeployedCCInfoProvider
}
