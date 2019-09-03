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
	"github.com/mcc-github/blockchain-chaincode-go/shim"
	pb "github.com/mcc-github/blockchain-protos-go/peer"
	"github.com/mcc-github/blockchain-protos-go/transientstore"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/common/util"
	"github.com/mcc-github/blockchain/core/common/ccprovider"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/internal/pkg/identity"
	"github.com/mcc-github/blockchain/msp"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var endorserLogger = flogging.MustGetLogger("endorser")






type PrivateDataDistributor interface {
	DistributePrivateData(channel string, txID string, privateData *transientstore.TxPvtReadWriteSetWithConfigInfo, blkHt uint64) error
}


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

	
	EndorseWithPlugin(pluginName, channnelID string, prpBytes []byte, signedProposal *pb.SignedProposal) (*pb.Endorsement, []byte, error)

	
	GetLedgerHeight(channelID string) (uint64, error)

	
	GetDeployedCCInfoProvider() ledger.DeployedChaincodeInfoProvider
}




type ChannelFetcher interface {
	Channel(channelID string) *Channel
}

type Channel struct {
	IdentityDeserializer msp.IdentityDeserializer
}


type Endorser struct {
	ChannelFetcher         ChannelFetcher
	LocalMSP               msp.IdentityDeserializer
	PrivateDataDistributor PrivateDataDistributor
	Support                Support
	PvtRWSetAssembler      PvtRWSetAssembler
	Metrics                *Metrics
}


func (e *Endorser) callChaincode(txParams *ccprovider.TransactionParams, input *pb.ChaincodeInput, chaincodeName string) (*pb.Response, *pb.ChaincodeEvent, error) {
	endorserLogger.Infof("[%s][%s] Entry chaincode: %s", txParams.ChannelID, shorttxid(txParams.TxID), chaincodeName)
	defer func(start time.Time) {
		logger := endorserLogger.WithOptions(zap.AddCallerSkip(1))
		elapsedMilliseconds := time.Since(start).Round(time.Millisecond) / time.Millisecond
		logger.Infof("[%s][%s] Exit chaincode: %s (%dms)", txParams.ChannelID, shorttxid(txParams.TxID), chaincodeName, elapsedMilliseconds)
	}(time.Now())

	res, ccevent, err := e.Support.Execute(txParams, chaincodeName, input)
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

	
	if e.Support.IsSysCC(cds.ChaincodeSpec.ChaincodeId.Name) {
		return nil, nil, errors.Errorf("attempting to deploy a system chaincode %s/%s", cds.ChaincodeSpec.ChaincodeId.Name, txParams.ChannelID)
	}

	if len(cds.CodePackage) != 0 {
		return nil, nil, errors.Errorf("lscc upgrade/deploy should not include a code packages")
	}

	_, _, err = e.Support.ExecuteLegacyInit(txParams, cds.ChaincodeSpec.ChaincodeId.Name, cds.ChaincodeSpec.ChaincodeId.Version, cds.ChaincodeSpec.Input)
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
		pvtDataWithConfig, err := AssemblePvtRWSet(txParams.ChannelID, simResult.PvtSimulationResults, txParams.TXSimulator, e.Support.GetDeployedCCInfoProvider())
		
		
		txParams.TXSimulator.Done()

		if err != nil {
			return nil, nil, nil, errors.WithMessage(err, "failed to obtain collections config")
		}
		endorsedAt, err := e.Support.GetLedgerHeight(txParams.ChannelID)
		if err != nil {
			return nil, nil, nil, errors.WithMessage(err, fmt.Sprintf("failed to obtain ledger height for channel '%s'", txParams.ChannelID))
		}
		
		
		
		
		
		pvtDataWithConfig.EndorsedAt = endorsedAt
		if err := e.PrivateDataDistributor.DistributePrivateData(txParams.ChannelID, txParams.TxID, pvtDataWithConfig, endorsedAt); err != nil {
			return nil, nil, nil, err
		}
	}

	pubSimResBytes, err := simResult.GetPubSimulationBytes()
	if err != nil {
		return nil, nil, nil, err
	}

	return res, pubSimResBytes, ccevent, nil
}


func (e *Endorser) preProcess(up *UnpackedProposal, channel *Channel) error {
	

	err := up.Validate(channel.IdentityDeserializer)
	if err != nil {
		e.Metrics.ProposalValidationFailed.Add(1)
		return errors.WithMessage(err, "error validating proposal")
	}

	endorserLogger.Debugf("[%s][%s] processing txid: %s", up.ChannelHeader.ChannelId, shorttxid(up.ChannelHeader.TxId), up.ChannelHeader.TxId)

	if up.ChannelHeader.ChannelId == "" {
		
		
		
		
		return nil
	}

	
	meterLabels := []string{
		"channel", up.ChannelHeader.ChannelId,
		"chaincode", up.ChaincodeName,
	}

	
	
	if _, err = e.Support.GetTransactionByID(up.ChannelHeader.ChannelId, up.ChannelHeader.TxId); err == nil {
		
		
		e.Metrics.DuplicateTxsFailure.With(meterLabels...).Add(1)
		return errors.Errorf("duplicate transaction found [%s]. Creator [%x]", up.ChannelHeader.TxId, up.SignatureHeader.Creator)
	}

	
	
	if !e.Support.IsSysCC(up.ChaincodeName) {
		
		if err = e.Support.CheckACL(up.ChannelHeader.ChannelId, up.SignedProposal); err != nil {
			e.Metrics.ProposalACLCheckFailed.With(meterLabels...).Add(1)
			return err
		}
	}

	return nil
}


func (e *Endorser) ProcessProposal(ctx context.Context, signedProp *pb.SignedProposal) (*pb.ProposalResponse, error) {
	
	startTime := time.Now()
	e.Metrics.ProposalsReceived.Add(1)

	addr := util.ExtractRemoteAddress(ctx)
	endorserLogger.Debug("Entering: request from", addr)
	defer endorserLogger.Debug("Exit: request from", addr)

	
	success := false

	up, err := UnpackProposal(signedProp)
	if err != nil {
		e.Metrics.ProposalValidationFailed.Add(1)
		return &pb.ProposalResponse{Response: &pb.Response{Status: 500, Message: err.Error()}}, err
	}

	var channel *Channel
	if up.ChannelID() != "" {
		channel = e.ChannelFetcher.Channel(up.ChannelID())
		if channel == nil {
			return &pb.ProposalResponse{Response: &pb.Response{Status: 500, Message: fmt.Sprintf("channel '%s' not found", up.ChannelHeader.ChannelId)}}, nil
		}
	} else {
		channel = &Channel{
			IdentityDeserializer: e.LocalMSP,
		}
	}

	
	err = e.preProcess(up, channel)
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

	pResp, err := e.ProcessProposalSuccessfullyOrError(up)
	if err != nil {
		return &pb.ProposalResponse{Response: &pb.Response{Status: 500, Message: err.Error()}}, nil
	}

	if pResp.Endorsement != nil || up.ChannelHeader.ChannelId == "" {
		
		
		
		success = true

		
		e.Metrics.SuccessfulProposals.Add(1)
	}
	return pResp, nil
}

func (e *Endorser) ProcessProposalSuccessfullyOrError(up *UnpackedProposal) (*pb.ProposalResponse, error) {
	txParams := &ccprovider.TransactionParams{
		ChannelID:  up.ChannelHeader.ChannelId,
		TxID:       up.ChannelHeader.TxId,
		SignedProp: up.SignedProposal,
		Proposal:   up.Proposal,
	}

	if acquireTxSimulator(up.ChannelHeader.ChannelId, up.ChaincodeName) {
		txSim, err := e.Support.GetTxSimulator(up.ChannelID(), up.TxID())
		if err != nil {
			return nil, err
		}

		
		
		
		
		
		
		
		defer txSim.Done()

		hqe, err := e.Support.GetHistoryQueryExecutor(up.ChannelID())
		if err != nil {
			return nil, err
		}

		txParams.TXSimulator = txSim
		txParams.HistoryQueryExecutor = hqe
	}

	cdLedger, err := e.Support.GetChaincodeDefinition(up.ChannelHeader.ChannelId, up.ChaincodeName, txParams.TXSimulator)
	if err != nil {
		return nil, errors.WithMessagef(err, "make sure the chaincode %s has been successfully defined on channel %s and try again", up.ChaincodeName, up.ChannelID())
	}

	
	res, simulationResult, ccevent, err := e.SimulateProposal(txParams, up.ChaincodeName, up.Input)
	if err != nil {
		return nil, errors.WithMessage(err, "error in simulation")
	}

	cceventBytes, err := CreateCCEventBytes(ccevent)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal chaincode event")
	}

	prpBytes, err := protoutil.GetBytesProposalResponsePayload(up.ProposalHash, res, simulationResult, cceventBytes, &pb.ChaincodeID{
		Name:    up.ChaincodeName,
		Version: cdLedger.CCVersion(),
	})
	if err != nil {
		endorserLogger.Warning("Failed marshaling the proposal response payload to bytes", err)
		return nil, errors.WithMessage(err, "failed to create the proposal response")
	}

	
	meterLabels := []string{
		"channel", up.ChannelID(),
		"chaincode", up.ChaincodeName,
	}

	switch {
	case res.Status >= shim.ERROR:
		return &pb.ProposalResponse{
			Response: res,
			Payload:  prpBytes,
		}, nil
	case up.ChannelID() == "":
		
		
		
		return &pb.ProposalResponse{
			Response: res,
		}, nil
	case res.Status >= shim.ERRORTHRESHOLD:
		meterLabels = append(meterLabels, "chaincodeerror", strconv.FormatBool(true))
		e.Metrics.EndorsementsFailed.With(meterLabels...).Add(1)
		endorserLogger.Debugf("[%s][%s] endorseProposal() resulted in chaincode %s error for txid: %s", up.ChannelID(), shorttxid(up.TxID()), up.ChaincodeName, up.TxID())
		return &pb.ProposalResponse{
			Response: res,
		}, nil
	}

	escc := cdLedger.Endorsement()

	endorserLogger.Debugf("[%s][%s] escc for chaincode %s is %s", up.ChannelID(), shorttxid(up.TxID()), up.ChaincodeName, escc)

	
	endorsement, mPrpBytes, err := e.Support.EndorseWithPlugin(escc, up.ChannelID(), prpBytes, up.SignedProposal)
	if err != nil {
		meterLabels = append(meterLabels, "chaincodeerror", strconv.FormatBool(false))
		e.Metrics.EndorsementsFailed.With(meterLabels...).Add(1)
		return nil, errors.WithMessage(err, "endorsing with plugin failed")
	}

	return &pb.ProposalResponse{
		Version:     1,
		Endorsement: endorsement,
		Payload:     mPrpBytes,
		Response:    res,
	}, nil
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
