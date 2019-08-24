/*
Copyright IBM Corp. 2016 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package endorser

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/common/metrics"
	"github.com/mcc-github/blockchain/common/util"
	"github.com/mcc-github/blockchain/core/chaincode/shim"
	"github.com/mcc-github/blockchain/core/common/ccprovider"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/internal/pkg/identity"
	pb "github.com/mcc-github/blockchain/protos/peer"
	"github.com/mcc-github/blockchain/protos/transientstore"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var endorserLogger = flogging.MustGetLogger("endorser")




type privateDataDistributor func(channel string, txID string, privateData *transientstore.TxPvtReadWriteSetWithConfigInfo, blkHt uint64) error


type Support interface {
	identity.SignerSerializer
	
	
	
	GetTxSimulator(ledgername string, txid string) (ledger.TxSimulator, error)

	
	
	GetHistoryQueryExecutor(ledgername string) (ledger.HistoryQueryExecutor, error)

	
	GetTransactionByID(chid, txID string) (*pb.ProcessedTransaction, error)

	
	
	IsSysCC(name string) bool

	
	Execute(txParams *ccprovider.TransactionParams, name string, input *pb.ChaincodeInput) (*pb.Response, *pb.ChaincodeEvent, error)

	
	ExecuteLegacyInit(txParams *ccprovider.TransactionParams, name, version string, spec *pb.ChaincodeInput) (*pb.Response, *pb.ChaincodeEvent, error)

	
	GetChaincodeDefinition(channelID, chaincodeID string, txsim ledger.QueryExecutor) (ccprovider.ChaincodeDefinition, error)

	
	
	CheckACL(channelID string, signedProp *pb.SignedProposal) error

	
	EndorseWithPlugin(ctx Context) (*pb.ProposalResponse, error)

	
	GetLedgerHeight(channelID string) (uint64, error)

	
	GetDeployedCCInfoProvider() ledger.DeployedChaincodeInfoProvider
}


type Endorser struct {
	distributePrivateData privateDataDistributor
	s                     Support
	PvtRWSetAssembler
	Metrics *EndorserMetrics
}


type validateResult struct {
	prop    *pb.Proposal
	hdrExt  *pb.ChaincodeHeaderExtension
	chainID string
	txid    string
	resp    *pb.ProposalResponse
}


func NewEndorserServer(privDist privateDataDistributor, s Support, metricsProv metrics.Provider) *Endorser {
	e := &Endorser{
		distributePrivateData: privDist,
		s:                     s,
		PvtRWSetAssembler:     &rwSetAssembler{},
		Metrics:               NewEndorserMetrics(metricsProv),
	}
	return e
}


func (e *Endorser) callChaincode(txParams *ccprovider.TransactionParams, input *pb.ChaincodeInput, chaincodeName string) (*pb.Response, *pb.ChaincodeEvent, error) {
	endorserLogger.Infof("[%s][%s] Entry chaincode: %s", txParams.ChannelID, shorttxid(txParams.TxID), chaincodeName)
	defer func(start time.Time) {
		logger := endorserLogger.WithOptions(zap.AddCallerSkip(1))
		elapsedMilliseconds := time.Since(start).Round(time.Millisecond) / time.Millisecond
		logger.Infof("[%s][%s] Exit chaincode: %s (%dms)", txParams.ChannelID, shorttxid(txParams.TxID), chaincodeName, elapsedMilliseconds)
	}(time.Now())

	res, ccevent, err := e.s.Execute(txParams, chaincodeName, input)
	if err != nil {
		return nil, nil, err
	}

	
	
	
	if res.Status >= shim.ERRORTHRESHOLD {
		return res, nil, nil
	}

	
	if chaincodeName != "lscc" || len(input.Args) < 3 || (string(input.Args[0]) != "deploy" && string(input.Args[0]) != "upgrade") {
		return res, ccevent, err
	}

	
	
	
	
	
	
	
	
	cds, err := protoutil.UnmarshalChaincodeDeploymentSpec(input.Args[2])
	if err != nil {
		return nil, nil, err
	}

	
	if e.s.IsSysCC(cds.ChaincodeSpec.ChaincodeId.Name) {
		return nil, nil, errors.Errorf("attempting to deploy a system chaincode %s/%s", cds.ChaincodeSpec.ChaincodeId.Name, txParams.ChannelID)
	}

	_, _, err = e.s.ExecuteLegacyInit(txParams, cds.ChaincodeSpec.ChaincodeId.Name, cds.ChaincodeSpec.ChaincodeId.Version, cds.ChaincodeSpec.Input)
	if err != nil {
		
		meterLabels := []string{
			"channel", txParams.ChannelID,
			"chaincode", cds.ChaincodeSpec.ChaincodeId.Name,
		}
		e.Metrics.InitFailed.With(meterLabels...).Add(1)
		return nil, nil, err
	}

	return res, ccevent, err

}


func (e *Endorser) SimulateProposal(txParams *ccprovider.TransactionParams, chaincodeName string, chaincodeInput *pb.ChaincodeInput) (*pb.Response, []byte, *pb.ChaincodeEvent, error) {
	endorserLogger.Debugf("[%s][%s] Entry chaincode: %s", txParams.ChannelID, shorttxid(txParams.TxID), chaincodeName)
	defer endorserLogger.Debugf("[%s][%s] Exit", txParams.ChannelID, shorttxid(txParams.TxID))

	
	res, ccevent, err := e.callChaincode(txParams, chaincodeInput, chaincodeName)
	if err != nil {
		endorserLogger.Errorf("[%s][%s] failed to invoke chaincode %s, error: %+v", txParams.ChannelID, shorttxid(txParams.TxID), chaincodeName, err)
		return nil, nil, nil, err
	}

	if txParams.TXSimulator == nil {
		return res, nil, ccevent, nil
	}

	
	
	
	defer txParams.TXSimulator.Done()

	simResult, err := txParams.TXSimulator.GetTxSimulationResults()
	if err != nil {
		return nil, nil, nil, err
	}

	if simResult.PvtSimulationResults != nil {
		if chaincodeName == "lscc" {
			
			return nil, nil, nil, errors.New("Private data is forbidden to be used in instantiate")
		}
		pvtDataWithConfig, err := e.AssemblePvtRWSet(txParams.ChannelID, simResult.PvtSimulationResults, txParams.TXSimulator, e.s.GetDeployedCCInfoProvider())
		
		
		txParams.TXSimulator.Done()

		if err != nil {
			return nil, nil, nil, errors.WithMessage(err, "failed to obtain collections config")
		}
		endorsedAt, err := e.s.GetLedgerHeight(txParams.ChannelID)
		if err != nil {
			return nil, nil, nil, errors.WithMessage(err, fmt.Sprint("failed to obtain ledger height for channel", txParams.ChannelID))
		}
		
		
		
		
		
		pvtDataWithConfig.EndorsedAt = endorsedAt
		if err := e.distributePrivateData(txParams.ChannelID, txParams.TxID, pvtDataWithConfig, endorsedAt); err != nil {
			return nil, nil, nil, err
		}
	}

	pubSimResBytes, err := simResult.GetPubSimulationBytes()
	if err != nil {
		return nil, nil, nil, err
	}

	return res, pubSimResBytes, ccevent, nil
}


func (e *Endorser) endorseProposal(up *UnpackedProposal, response *pb.Response, simRes, eventBytes []byte, cd ccprovider.ChaincodeDefinition) (*pb.ProposalResponse, error) {
	endorserLogger.Debugf("[%s][%s] Entry chaincode: %s", up.ChannelHeader.ChannelId, shorttxid(up.ChannelHeader.TxId), up.ChaincodeName)
	defer endorserLogger.Debugf("[%s][%s] Exit", up.ChannelHeader.ChannelId, shorttxid(up.ChannelHeader.TxId))

	escc := cd.Endorsement()

	endorserLogger.Debugf("[%s][%s] escc for chaincode %s is %s", up.ChannelHeader.ChannelId, shorttxid(up.ChannelHeader.TxId), up.ChaincodeName, escc)

	ctx := Context{
		PluginName:     escc,
		Channel:        up.ChannelHeader.ChannelId,
		SignedProposal: up.SignedProposal,
		ChaincodeID: &pb.ChaincodeID{
			Name:    up.ChaincodeName,
			Version: cd.CCVersion(),
		},
		Event:    eventBytes,
		SimRes:   simRes,
		Response: response,
		
		Proposal: up.Proposal,
		TxID:     up.ChannelHeader.TxId,
	}
	return e.s.EndorseWithPlugin(ctx)
}


func (e *Endorser) preProcess(signedProp *pb.SignedProposal) (*UnpackedProposal, error) {
	
	up, err := UnpackProposal(signedProp)
	if err != nil {
		e.Metrics.ProposalValidationFailed.Add(1)
		return nil, err
	}

	err = ValidateUnpackedProposal(up)
	if err != nil {
		e.Metrics.ProposalValidationFailed.Add(1)
		return nil, err
	}

	endorserLogger.Debugf("[%s][%s] processing txid: %s", up.ChannelHeader.ChannelId, shorttxid(up.ChannelHeader.TxId), up.ChannelHeader.TxId)

	if up.ChannelHeader.ChannelId == "" {
		
		
		
		
		return up, nil
	}

	
	meterLabels := []string{
		"channel", up.ChannelHeader.ChannelId,
		"chaincode", up.ChaincodeName,
	}

	
	
	if _, err = e.s.GetTransactionByID(up.ChannelHeader.ChannelId, up.ChannelHeader.TxId); err == nil {
		
		
		e.Metrics.DuplicateTxsFailure.With(meterLabels...).Add(1)
		return nil, errors.Errorf("duplicate transaction found [%s]. Creator [%x]", up.ChannelHeader.TxId, up.SignatureHeader.Creator)
	}

	
	
	if !e.s.IsSysCC(up.ChaincodeName) {
		
		if err = e.s.CheckACL(up.ChannelHeader.ChannelId, up.SignedProposal); err != nil {
			e.Metrics.ProposalACLCheckFailed.With(meterLabels...).Add(1)
			return nil, err
		}
	}

	return up, nil
}


func (e *Endorser) ProcessProposal(ctx context.Context, signedProp *pb.SignedProposal) (*pb.ProposalResponse, error) {
	
	startTime := time.Now()
	e.Metrics.ProposalsReceived.Add(1)

	addr := util.ExtractRemoteAddress(ctx)
	endorserLogger.Debug("Entering: request from", addr)
	defer endorserLogger.Debug("Exit: request from", addr)

	
	success := false

	
	up, err := e.preProcess(signedProp)
	if err != nil {
		return &pb.ProposalResponse{Response: &pb.Response{Status: 500, Message: err.Error()}}, err
	}

	defer func() {
		meterLabels := []string{
			"channel", up.ChannelHeader.ChannelId,
			"chaincode", up.ChaincodeName,
			"success", strconv.FormatBool(success),
		}
		e.Metrics.ProposalDuration.With(meterLabels...).Observe(time.Since(startTime).Seconds())

	}()

	txParams := &ccprovider.TransactionParams{
		ChannelID:  up.ChannelHeader.ChannelId,
		TxID:       up.ChannelHeader.TxId,
		SignedProp: up.SignedProposal,
		Proposal:   up.Proposal,
	}

	if acquireTxSimulator(up.ChannelHeader.ChannelId, up.ChaincodeName) {
		txSim, err := e.s.GetTxSimulator(up.ChannelHeader.ChannelId, up.ChannelHeader.TxId)
		if err != nil {
			return &pb.ProposalResponse{Response: &pb.Response{Status: 500, Message: err.Error()}}, nil
		}

		
		
		
		
		
		
		
		defer txSim.Done()

		hqe, err := e.s.GetHistoryQueryExecutor(up.ChannelHeader.ChannelId)
		if err != nil {
			return &pb.ProposalResponse{Response: &pb.Response{Status: 500, Message: err.Error()}}, nil
		}

		txParams.TXSimulator = txSim
		txParams.HistoryQueryExecutor = hqe
	}

	cdLedger, err := e.s.GetChaincodeDefinition(up.ChannelHeader.ChannelId, up.ChaincodeName, txParams.TXSimulator)
	if err != nil {
		return &pb.ProposalResponse{Response: &pb.Response{Status: 500, Message: fmt.Sprintf("make sure the chaincode %s has been successfully defined on channel %s and try again: %s", up.ChaincodeName, up.ChannelID(), err)}}, nil
	}

	
	res, simulationResult, ccevent, err := e.SimulateProposal(txParams, up.ChaincodeName, up.Input)
	if err != nil {
		return &pb.ProposalResponse{Response: &pb.Response{Status: 500, Message: err.Error()}}, nil
	}

	cceventBytes, err := CreateCCEventBytes(ccevent)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal chaincode event")
	}

	if res.Status >= shim.ERROR {
		endorserLogger.Errorf("[%s][%s] simulateProposal() resulted in chaincode %s response status %d for txid: %s", up.ChannelID(), shorttxid(up.TxID()), up.ChaincodeName, res.Status, up.TxID())
		pResp, err := protoutil.CreateProposalResponseFailure(up.Proposal.Header, up.Proposal.Payload, res, simulationResult, cceventBytes, up.ChaincodeName)
		if err != nil {
			return &pb.ProposalResponse{Response: &pb.Response{Status: 500, Message: err.Error()}}, nil
		}

		return pResp, nil
	}

	
	var pResp *pb.ProposalResponse

	
	
	if up.ChannelID() == "" {
		pResp = &pb.ProposalResponse{}
	} else {
		pResp, err = e.endorseProposal(up, res, simulationResult, cceventBytes, cdLedger)

		
		meterLabels := []string{
			"channel", up.ChannelID(),
			"chaincode", up.ChaincodeName,
		}

		if err != nil {
			meterLabels = append(meterLabels, "chaincodeerror", strconv.FormatBool(false))
			e.Metrics.EndorsementsFailed.With(meterLabels...).Add(1)
			return &pb.ProposalResponse{Response: &pb.Response{Status: 500, Message: err.Error()}}, nil
		}
		if pResp.Response.Status >= shim.ERRORTHRESHOLD {
			
			
			meterLabels = append(meterLabels, "chaincodeerror", strconv.FormatBool(true))
			e.Metrics.EndorsementsFailed.With(meterLabels...).Add(1)
			endorserLogger.Debugf("[%s][%s] endorseProposal() resulted in chaincode %s error for txid: %s", up.ChannelID(), shorttxid(up.TxID()), up.ChaincodeName, up.TxID())
			return pResp, nil
		}
	}

	
	
	
	pResp.Response = res

	
	e.Metrics.SuccessfulProposals.Add(1)
	success = true

	return pResp, nil
}



func acquireTxSimulator(chainID string, chaincodeName string) bool {
	if chainID == "" {
		return false
	}

	
	
	
	switch chaincodeName {
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

func CreateCCEventBytes(ccevent *pb.ChaincodeEvent) ([]byte, error) {
	if ccevent == nil {
		return nil, nil
	}

	return proto.Marshal(ccevent)
}
