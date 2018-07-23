/*
Copyright IBM Corp. 2016 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package endorser

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/common/channelconfig"
	"github.com/mcc-github/blockchain/common/crypto"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/common/util"
	"github.com/mcc-github/blockchain/core/chaincode/platforms"
	"github.com/mcc-github/blockchain/core/chaincode/shim"
	"github.com/mcc-github/blockchain/core/common/ccprovider"
	"github.com/mcc-github/blockchain/core/common/validation"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/protos/common"
	pb "github.com/mcc-github/blockchain/protos/peer"
	"github.com/mcc-github/blockchain/protos/transientstore"
	putils "github.com/mcc-github/blockchain/protos/utils"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

var endorserLogger = flogging.MustGetLogger("endorser")




type privateDataDistributor func(channel string, txID string, privateData *transientstore.TxPvtReadWriteSetWithConfigInfo, blkHt uint64) error


type Support interface {
	crypto.SignerSupport
	
	
	IsSysCCAndNotInvokableExternal(name string) bool

	
	
	
	GetTxSimulator(ledgername string, txid string) (ledger.TxSimulator, error)

	
	
	GetHistoryQueryExecutor(ledgername string) (ledger.HistoryQueryExecutor, error)

	
	GetTransactionByID(chid, txID string) (*pb.ProcessedTransaction, error)

	
	
	IsSysCC(name string) bool

	
	Execute(txParams *ccprovider.TransactionParams, cid, name, version, txid string, signedProp *pb.SignedProposal, prop *pb.Proposal, input *pb.ChaincodeInput) (*pb.Response, *pb.ChaincodeEvent, error)

	
	ExecuteLegacyInit(txParams *ccprovider.TransactionParams, cid, name, version, txid string, signedProp *pb.SignedProposal, prop *pb.Proposal, spec *pb.ChaincodeDeploymentSpec) (*pb.Response, *pb.ChaincodeEvent, error)

	
	GetChaincodeDefinition(chaincodeID string, txsim ledger.QueryExecutor) (ccprovider.ChaincodeDefinition, error)

	
	
	CheckACL(signedProp *pb.SignedProposal, chdr *common.ChannelHeader, shdr *common.SignatureHeader, hdrext *pb.ChaincodeHeaderExtension) error

	
	
	IsJavaCC(buf []byte) (bool, error)

	
	
	CheckInstantiationPolicy(name, version string, cd ccprovider.ChaincodeDefinition) error

	
	GetChaincodeDeploymentSpecFS(cds *pb.ChaincodeDeploymentSpec) (*pb.ChaincodeDeploymentSpec, error)

	
	
	GetApplicationConfig(cid string) (channelconfig.Application, bool)

	
	NewQueryCreator(channel string) (QueryCreator, error)

	
	EndorseWithPlugin(ctx Context) (*pb.ProposalResponse, error)

	
	GetLedgerHeight(channelID string) (uint64, error)
}


type Endorser struct {
	distributePrivateData privateDataDistributor
	s                     Support
	PlatformRegistry      *platforms.Registry
	PvtRWSetAssembler
}


type validateResult struct {
	prop    *pb.Proposal
	hdrExt  *pb.ChaincodeHeaderExtension
	chainID string
	txid    string
	resp    *pb.ProposalResponse
}


func NewEndorserServer(privDist privateDataDistributor, s Support, pr *platforms.Registry) *Endorser {
	e := &Endorser{
		distributePrivateData: privDist,
		s:                 s,
		PlatformRegistry:  pr,
		PvtRWSetAssembler: &rwSetAssembler{},
	}
	return e
}


func (e *Endorser) callChaincode(txParams *ccprovider.TransactionParams, version string, input *pb.ChaincodeInput, cid *pb.ChaincodeID) (*pb.Response, *pb.ChaincodeEvent, error) {
	endorserLogger.Debugf("[%s][%s] Entry chaincode: %s version: %s", txParams.ChannelID, txParams.TxID, cid, version)
	defer endorserLogger.Debugf("[%s][%s] Exit", txParams.ChannelID, txParams.TxID)
	var err error
	var res *pb.Response
	var ccevent *pb.ChaincodeEvent

	
	res, ccevent, err = e.s.Execute(txParams, txParams.ChannelID, cid.Name, version, txParams.TxID, txParams.SignedProp, txParams.Proposal, input)
	if err != nil {
		return nil, nil, err
	}

	
	
	
	if res.Status >= shim.ERRORTHRESHOLD {
		return res, nil, nil
	}

	
	
	
	
	
	
	
	
	if cid.Name == "lscc" && len(input.Args) >= 3 && (string(input.Args[0]) == "deploy" || string(input.Args[0]) == "upgrade") {
		userCDS, err := putils.GetChaincodeDeploymentSpec(input.Args[2], e.PlatformRegistry)
		if err != nil {
			return nil, nil, err
		}

		var cds *pb.ChaincodeDeploymentSpec
		cds, err = e.SanitizeUserCDS(userCDS)
		if err != nil {
			return nil, nil, err
		}

		
		if e.s.IsSysCC(cds.ChaincodeSpec.ChaincodeId.Name) {
			return nil, nil, errors.Errorf("attempting to deploy a system chaincode %s/%s", cds.ChaincodeSpec.ChaincodeId.Name, txParams.ChannelID)
		}

		_, _, err = e.s.ExecuteLegacyInit(txParams, txParams.ChannelID, cds.ChaincodeSpec.ChaincodeId.Name, cds.ChaincodeSpec.ChaincodeId.Version, txParams.TxID, txParams.SignedProp, txParams.Proposal, cds)
		if err != nil {
			return nil, nil, err
		}
	}
	

	return res, ccevent, err
}

func (e *Endorser) SanitizeUserCDS(userCDS *pb.ChaincodeDeploymentSpec) (*pb.ChaincodeDeploymentSpec, error) {
	fsCDS, err := e.s.GetChaincodeDeploymentSpecFS(userCDS)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot deploy a chaincode which is not installed")
	}

	sanitizedCDS := proto.Clone(fsCDS).(*pb.ChaincodeDeploymentSpec)
	sanitizedCDS.CodePackage = nil
	sanitizedCDS.ChaincodeSpec.Input = userCDS.ChaincodeSpec.Input

	return sanitizedCDS, nil
}



func (e *Endorser) DisableJavaCCInst(cid *pb.ChaincodeID, cis *pb.ChaincodeInvocationSpec) error {
	
	if cid.Name != "lscc" {
		return nil
	}

	
	if cis.ChaincodeSpec == nil || cis.ChaincodeSpec.Input == nil {
		return nil
	}

	
	if len(cis.ChaincodeSpec.Input.Args) < 1 {
		return nil
	}

	var argNo int
	switch string(cis.ChaincodeSpec.Input.Args[0]) {
	case "install":
		argNo = 1
	case "deploy", "upgrade":
		argNo = 2
	default:
		
		return nil
	}

	if argNo >= len(cis.ChaincodeSpec.Input.Args) {
		return errors.Errorf("too few arguments passed. expected %d", argNo)
	}

	if JavaEnabled() {
		endorserLogger.Debug("java chaincode enabled")
	} else {
		endorserLogger.Debug("java chaincode disabled")
		
		isjava, err := e.s.IsJavaCC(cis.ChaincodeSpec.Input.Args[argNo])
		if err != nil {
			return err
		}
		if isjava {
			return errors.New("Java chaincode is work-in-progress and disabled")
		}
	}

	
	return nil
}


func (e *Endorser) SimulateProposal(txParams *ccprovider.TransactionParams, cid *pb.ChaincodeID) (ccprovider.ChaincodeDefinition, *pb.Response, []byte, *pb.ChaincodeEvent, error) {
	endorserLogger.Debugf("[%s][%s] Entry chaincode: %s", txParams.ChannelID, shorttxid(txParams.TxID), cid)
	defer endorserLogger.Debugf("[%s][%s] Exit", txParams.ChannelID, shorttxid(txParams.TxID))
	
	
	
	cis, err := putils.GetChaincodeInvocationSpec(txParams.Proposal)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	
	if err = e.DisableJavaCCInst(cid, cis); err != nil {
		return nil, nil, nil, nil, err
	}

	var cdLedger ccprovider.ChaincodeDefinition
	var version string

	if !e.s.IsSysCC(cid.Name) {
		cdLedger, err = e.s.GetChaincodeDefinition(cid.Name, txParams.TXSimulator)
		if err != nil {
			return nil, nil, nil, nil, errors.WithMessage(err, fmt.Sprintf("make sure the chaincode %s has been successfully instantiated and try again", cid.Name))
		}
		version = cdLedger.CCVersion()

		err = e.s.CheckInstantiationPolicy(cid.Name, version, cdLedger)
		if err != nil {
			return nil, nil, nil, nil, err
		}
	} else {
		version = util.GetSysCCVersion()
	}

	
	var simResult *ledger.TxSimulationResults
	var pubSimResBytes []byte
	var res *pb.Response
	var ccevent *pb.ChaincodeEvent
	res, ccevent, err = e.callChaincode(txParams, version, cis.ChaincodeSpec.Input, cid)
	if err != nil {
		endorserLogger.Errorf("[%s][%s] failed to invoke chaincode %s, error: %+v", txParams.ChannelID, shorttxid(txParams.TxID), cid, err)
		return nil, nil, nil, nil, err
	}

	if txParams.TXSimulator != nil {
		if simResult, err = txParams.TXSimulator.GetTxSimulationResults(); err != nil {
			txParams.TXSimulator.Done()
			return nil, nil, nil, nil, err
		}

		if simResult.PvtSimulationResults != nil {
			if cid.Name == "lscc" {
				
				txParams.TXSimulator.Done()
				return nil, nil, nil, nil, errors.New("Private data is forbidden to be used in instantiate")
			}
			pvtDataWithConfig, err := e.AssemblePvtRWSet(simResult.PvtSimulationResults, txParams.TXSimulator)
			
			
			txParams.TXSimulator.Done()

			if err != nil {
				return nil, nil, nil, nil, errors.WithMessage(err, "failed to obtain collections config")
			}
			endorsedAt, err := e.s.GetLedgerHeight(txParams.ChannelID)
			if err != nil {
				return nil, nil, nil, nil, errors.WithMessage(err, fmt.Sprint("failed to obtain ledger height for channel", txParams.ChannelID))
			}
			
			
			
			
			
			pvtDataWithConfig.EndorsedAt = endorsedAt
			if err := e.distributePrivateData(txParams.ChannelID, txParams.TxID, pvtDataWithConfig, endorsedAt); err != nil {
				return nil, nil, nil, nil, err
			}
		}

		txParams.TXSimulator.Done()
		if pubSimResBytes, err = simResult.GetPubSimulationBytes(); err != nil {
			return nil, nil, nil, nil, err
		}
	}
	return cdLedger, res, pubSimResBytes, ccevent, nil
}


func (e *Endorser) endorseProposal(_ context.Context, chainID string, txid string, signedProp *pb.SignedProposal, proposal *pb.Proposal, response *pb.Response, simRes []byte, event *pb.ChaincodeEvent, visibility []byte, ccid *pb.ChaincodeID, txsim ledger.TxSimulator, cd ccprovider.ChaincodeDefinition) (*pb.ProposalResponse, error) {
	endorserLogger.Debugf("[%s][%s] Entry chaincode: %s", chainID, shorttxid(txid), ccid)
	defer endorserLogger.Debugf("[%s][%s] Exit", chainID, shorttxid(txid))

	isSysCC := cd == nil
	
	var escc string
	
	if isSysCC {
		escc = "escc"
	} else {
		escc = cd.Endorsement()
	}

	endorserLogger.Debugf("[%s][%s] escc for chaincode %s is %s", chainID, shorttxid(txid), ccid, escc)

	
	var err error
	var eventBytes []byte
	if event != nil {
		eventBytes, err = putils.GetBytesChaincodeEvent(event)
		if err != nil {
			return nil, errors.Wrap(err, "failed to marshal event bytes")
		}
	}

	
	if isSysCC {
		
		
		ccid.Version = util.GetSysCCVersion()
	} else {
		ccid.Version = cd.CCVersion()
	}

	ctx := Context{
		PluginName:     escc,
		Channel:        chainID,
		SignedProposal: signedProp,
		ChaincodeID:    ccid,
		Event:          eventBytes,
		SimRes:         simRes,
		Response:       response,
		Visibility:     visibility,
		Proposal:       proposal,
		TxID:           txid,
	}
	return e.s.EndorseWithPlugin(ctx)
}


func (e *Endorser) preProcess(signedProp *pb.SignedProposal) (*validateResult, error) {
	vr := &validateResult{}
	
	prop, hdr, hdrExt, err := validation.ValidateProposalMessage(signedProp)

	if err != nil {
		vr.resp = &pb.ProposalResponse{Response: &pb.Response{Status: 500, Message: err.Error()}}
		return vr, err
	}

	chdr, err := putils.UnmarshalChannelHeader(hdr.ChannelHeader)
	if err != nil {
		vr.resp = &pb.ProposalResponse{Response: &pb.Response{Status: 500, Message: err.Error()}}
		return vr, err
	}

	shdr, err := putils.GetSignatureHeader(hdr.SignatureHeader)
	if err != nil {
		vr.resp = &pb.ProposalResponse{Response: &pb.Response{Status: 500, Message: err.Error()}}
		return vr, err
	}

	
	if e.s.IsSysCCAndNotInvokableExternal(hdrExt.ChaincodeId.Name) {
		endorserLogger.Errorf("Error: an attempt was made by %#v to invoke system chaincode %s", shdr.Creator, hdrExt.ChaincodeId.Name)
		err = errors.Errorf("chaincode %s cannot be invoked through a proposal", hdrExt.ChaincodeId.Name)
		vr.resp = &pb.ProposalResponse{Response: &pb.Response{Status: 500, Message: err.Error()}}
		return vr, err
	}

	chainID := chdr.ChannelId
	txid := chdr.TxId
	endorserLogger.Debugf("[%s][%s] processing txid: %s", chainID, shorttxid(txid), txid)

	if chainID != "" {
		
		
		if _, err = e.s.GetTransactionByID(chainID, txid); err == nil {
			err = errors.Errorf("duplicate transaction found [%s]. Creator [%x]", txid, shdr.Creator)
			vr.resp = &pb.ProposalResponse{Response: &pb.Response{Status: 500, Message: err.Error()}}
			return vr, err
		}

		
		
		if !e.s.IsSysCC(hdrExt.ChaincodeId.Name) {
			
			if err = e.s.CheckACL(signedProp, chdr, shdr, hdrExt); err != nil {
				vr.resp = &pb.ProposalResponse{Response: &pb.Response{Status: 500, Message: err.Error()}}
				return vr, err
			}
		}
	} else {
		
		
		
		
	}

	vr.prop, vr.hdrExt, vr.chainID, vr.txid = prop, hdrExt, chainID, txid
	return vr, nil
}


func (e *Endorser) ProcessProposal(ctx context.Context, signedProp *pb.SignedProposal) (*pb.ProposalResponse, error) {
	addr := util.ExtractRemoteAddress(ctx)
	endorserLogger.Debug("Entering: request from", addr)
	defer endorserLogger.Debug("Exit: request from", addr)

	
	vr, err := e.preProcess(signedProp)
	if err != nil {
		resp := vr.resp
		return resp, err
	}

	prop, hdrExt, chainID, txid := vr.prop, vr.hdrExt, vr.chainID, vr.txid

	
	
	
	var txsim ledger.TxSimulator
	var historyQueryExecutor ledger.HistoryQueryExecutor
	if acquireTxSimulator(chainID, vr.hdrExt.ChaincodeId) {
		if txsim, err = e.s.GetTxSimulator(chainID, txid); err != nil {
			return &pb.ProposalResponse{Response: &pb.Response{Status: 500, Message: err.Error()}}, nil
		}

		
		
		
		
		
		
		
		defer txsim.Done()

		if historyQueryExecutor, err = e.s.GetHistoryQueryExecutor(chainID); err != nil {
			return &pb.ProposalResponse{Response: &pb.Response{Status: 500, Message: err.Error()}}, nil
		}
	}

	txParams := &ccprovider.TransactionParams{
		ChannelID:            chainID,
		TxID:                 txid,
		SignedProp:           signedProp,
		Proposal:             prop,
		TXSimulator:          txsim,
		HistoryQueryExecutor: historyQueryExecutor,
	}
	

	
	
	
	

	
	cd, res, simulationResult, ccevent, err := e.SimulateProposal(txParams, hdrExt.ChaincodeId)
	if err != nil {
		return &pb.ProposalResponse{Response: &pb.Response{Status: 500, Message: err.Error()}}, nil
	}
	if res != nil {
		if res.Status >= shim.ERROR {
			endorserLogger.Errorf("[%s][%s] simulateProposal() resulted in chaincode %s response status %d for txid: %s", chainID, shorttxid(txid), hdrExt.ChaincodeId, res.Status, txid)
			var cceventBytes []byte
			if ccevent != nil {
				cceventBytes, err = putils.GetBytesChaincodeEvent(ccevent)
				if err != nil {
					return nil, errors.Wrap(err, "failed to marshal event bytes")
				}
			}
			pResp, err := putils.CreateProposalResponseFailure(prop.Header, prop.Payload, res, simulationResult, cceventBytes, hdrExt.ChaincodeId, hdrExt.PayloadVisibility)
			if err != nil {
				return &pb.ProposalResponse{Response: &pb.Response{Status: 500, Message: err.Error()}}, nil
			}

			return pResp, nil
		}
	}

	
	var pResp *pb.ProposalResponse

	
	
	if chainID == "" {
		pResp = &pb.ProposalResponse{Response: res}
	} else {
		
		pResp, err = e.endorseProposal(ctx, chainID, txid, signedProp, prop, res, simulationResult, ccevent, hdrExt.PayloadVisibility, hdrExt.ChaincodeId, txsim, cd)
		if err != nil {
			return &pb.ProposalResponse{Response: &pb.Response{Status: 500, Message: err.Error()}}, nil
		}
		if pResp.Response.Status >= shim.ERRORTHRESHOLD {
			endorserLogger.Debugf("[%s][%s] endorseProposal() resulted in chaincode %s error for txid: %s", chainID, shorttxid(txid), hdrExt.ChaincodeId, txid)
			return pResp, nil
		}
	}

	
	
	
	pResp.Response = res

	return pResp, nil
}



func acquireTxSimulator(chainID string, ccid *pb.ChaincodeID) bool {
	if chainID == "" {
		return false
	}

	
	
	
	switch ccid.Name {
	case "qscc", "cscc":
		return false
	default:
		return true
	}
}




func shorttxid(txid string) string {
	if len(txid) < 8 {
		return txid
	}
	return txid[0:8]
}