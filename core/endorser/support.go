/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package endorser

import (
	"fmt"

	"github.com/mcc-github/blockchain/common/channelconfig"
	"github.com/mcc-github/blockchain/common/crypto"
	"github.com/mcc-github/blockchain/core/aclmgmt"
	"github.com/mcc-github/blockchain/core/aclmgmt/resources"
	"github.com/mcc-github/blockchain/core/chaincode"
	"github.com/mcc-github/blockchain/core/common/ccprovider"
	"github.com/mcc-github/blockchain/core/handlers/decoration"
	. "github.com/mcc-github/blockchain/core/handlers/endorsement/api/identities"
	"github.com/mcc-github/blockchain/core/handlers/library"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/peer"
	"github.com/mcc-github/blockchain/core/scc"
	"github.com/mcc-github/blockchain/protos/common"
	pb "github.com/mcc-github/blockchain/protos/peer"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)



type SupportImpl struct {
	*PluginEndorser
	crypto.SignerSupport
	Peer             peer.Operations
	PeerSupport      peer.Support
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
	return s.SignerSupport, nil
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


func (s *SupportImpl) Execute(ctxt context.Context, cid, name, version, txid string, syscc bool, signedProp *pb.SignedProposal, prop *pb.Proposal, spec ccprovider.ChaincodeSpecGetter) (*pb.Response, *pb.ChaincodeEvent, error) {
	cccid := ccprovider.NewCCContext(cid, name, version, txid, syscc, signedProp, prop)

	switch spec.(type) {
	case *pb.ChaincodeDeploymentSpec:
		return s.ChaincodeSupport.Execute(ctxt, cccid, spec)
	case *pb.ChaincodeInvocationSpec:
		cis := spec.(*pb.ChaincodeInvocationSpec)

		
		decorators := library.InitRegistry(library.Config{}).Lookup(library.Decoration).([]decoration.Decorator)
		cis.ChaincodeSpec.Input.Decorations = make(map[string][]byte)
		cis.ChaincodeSpec.Input = decoration.Apply(prop, cis.ChaincodeSpec.Input, decorators...)
		cccid.ProposalDecorations = cis.ChaincodeSpec.Input.Decorations

		return s.ChaincodeSupport.Execute(ctxt, cccid, cis)
	default:
		panic("programming error, unkwnown spec type")
	}
}


func (s *SupportImpl) GetChaincodeDefinition(ctx context.Context, chainID string, txid string, signedProp *pb.SignedProposal, prop *pb.Proposal, chaincodeID string, txsim ledger.TxSimulator) (ccprovider.ChaincodeDefinition, error) {
	ctxt := ctx
	if txsim != nil {
		ctxt = context.WithValue(ctx, chaincode.TXSimulatorKey, txsim)
	}
	lifecycle := &chaincode.Lifecycle{
		Executor: s.ChaincodeSupport,
	}
	return lifecycle.GetChaincodeDefinition(ctxt, txid, signedProp, prop, chainID, chaincodeID)
}



func (s *SupportImpl) CheckACL(signedProp *pb.SignedProposal, chdr *common.ChannelHeader, shdr *common.SignatureHeader, hdrext *pb.ChaincodeHeaderExtension) error {
	return s.ACLProvider.CheckACL(resources.Peer_Propose, chdr.ChannelId, signedProp)
}



func (s *SupportImpl) IsJavaCC(buf []byte) (bool, error) {
	
	ccpack, err := ccprovider.GetCCPackage(buf)
	if err != nil {
		return false, err
	}
	cds := ccpack.GetDepSpec()
	return (cds.ChaincodeSpec.Type == pb.ChaincodeSpec_JAVA), nil
}



func (s *SupportImpl) CheckInstantiationPolicy(name, version string, cd ccprovider.ChaincodeDefinition) error {
	return ccprovider.CheckInstantiationPolicy(name, version, cd.(*ccprovider.ChaincodeData))
}



func (s *SupportImpl) GetApplicationConfig(cid string) (channelconfig.Application, bool) {
	return s.PeerSupport.GetApplicationConfig(cid)
}
