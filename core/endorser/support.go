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
	ACLProvider      aclmgmt.ACLProvider
	BuiltinSCCs      scc.BuiltinSCCs
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
	return s.BuiltinSCCs.IsSysCC(name)
}


func (s *SupportImpl) ExecuteLegacyInit(txParams *ccprovider.TransactionParams, name, version string, input *pb.ChaincodeInput) (*pb.Response, *pb.ChaincodeEvent, error) {
	return s.ChaincodeSupport.ExecuteLegacyInit(txParams, name, version, input)
}


func (s *SupportImpl) Execute(txParams *ccprovider.TransactionParams, name string, input *pb.ChaincodeInput) (*pb.Response, *pb.ChaincodeEvent, error) {
	
	decorators := library.InitRegistry(library.Config{}).Lookup(library.Decoration).([]decoration.Decorator)
	input.Decorations = make(map[string][]byte)
	input = decoration.Apply(txParams.Proposal, input, decorators...)
	txParams.ProposalDecorations = input.Decorations

	return s.ChaincodeSupport.Execute(txParams, name, input)
}


func (s *SupportImpl) GetChaincodeDefinition(channelID, chaincodeName string, txsim ledger.QueryExecutor) (ccprovider.ChaincodeDefinition, error) {
	return s.ChaincodeSupport.Lifecycle.ChaincodeDefinition(channelID, chaincodeName, txsim)
}



func (s *SupportImpl) CheckACL(channelID string, signedProp *pb.SignedProposal) error {
	return s.ACLProvider.CheckACL(resources.Peer_Propose, channelID, signedProp)
}



func (s *SupportImpl) GetApplicationConfig(cid string) (channelconfig.Application, bool) {
	return s.Peer.GetApplicationConfig(cid)
}


func (s *SupportImpl) GetDeployedCCInfoProvider() ledger.DeployedChaincodeInfoProvider {
	return s.ChaincodeSupport.DeployedCCInfoProvider
}
