/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chaincode

import (
	"bytes"
	"time"
	"unicode/utf8"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/common/util"
	"github.com/mcc-github/blockchain/core/common/ccprovider"
	"github.com/mcc-github/blockchain/core/container/ccintf"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/peer"
	"github.com/mcc-github/blockchain/core/scc"
	pb "github.com/mcc-github/blockchain/protos/peer"
	"github.com/pkg/errors"
)

const (
	
	
	
	
	
	InitializedKeyName = "\x00" + string(utf8.MaxRune) + "initialized"
)


type Runtime interface {
	Start(ccid string) error
	Stop(ccid string) error
	Wait(ccid string) (int, error)
}


type Launcher interface {
	Launch(ccid string) error
}


type Lifecycle interface {
	
	ChaincodeDefinition(channelID, chaincodeName string, qe ledger.SimpleQueryExecutor) (ccprovider.ChaincodeDefinition, error)
}





type LegacyChaincodeDefinition interface {
	ExecuteLegacySecurityChecks() error
}


type ChaincodeSupport struct {
	ACLProvider            ACLProvider
	AppConfig              ApplicationConfigRetriever
	BuiltinSCCs            scc.BuiltinSCCs
	DeployedCCInfoProvider ledger.DeployedChaincodeInfoProvider
	ExecuteTimeout         time.Duration
	HandlerMetrics         *HandlerMetrics
	HandlerRegistry        *HandlerRegistry
	Keepalive              time.Duration
	Launcher               Launcher
	Lifecycle              Lifecycle
	Peer                   *peer.Peer
	Runtime                Runtime
	TotalQueryLimit        int
	UserRunsCC             bool
}




func (cs *ChaincodeSupport) Launch(ccid string) (*Handler, error) {
	if h := cs.HandlerRegistry.Handler(ccid); h != nil {
		return h, nil
	}

	if err := cs.Launcher.Launch(ccid); err != nil {
		return nil, errors.Wrapf(err, "could not launch chaincode %s", ccid)
	}

	h := cs.HandlerRegistry.Handler(ccid)
	if h == nil {
		return nil, errors.Errorf("claimed to start chaincode container for %s but could not find handler", ccid)
	}

	return h, nil
}


func (cs *ChaincodeSupport) LaunchInProc(ccid string) <-chan struct{} {
	launchStatus, ok := cs.HandlerRegistry.Launching(ccid)
	if ok {
		chaincodeLogger.Panicf("attempted to launch a system chaincode which has already been launched")
	}

	return launchStatus.Done()
}


func (cs *ChaincodeSupport) HandleChaincodeStream(stream ccintf.ChaincodeStream) error {
	handler := &Handler{
		Invoker:                cs,
		Keepalive:              cs.Keepalive,
		Registry:               cs.HandlerRegistry,
		ACLProvider:            cs.ACLProvider,
		TXContexts:             NewTransactionContexts(),
		ActiveTransactions:     NewActiveTransactions(),
		BuiltinSCCs:            cs.BuiltinSCCs,
		QueryResponseBuilder:   &QueryResponseGenerator{MaxResultLimit: 100},
		UUIDGenerator:          UUIDGeneratorFunc(util.GenerateUUID),
		LedgerGetter:           cs.Peer,
		DeployedCCInfoProvider: cs.DeployedCCInfoProvider,
		AppConfig:              cs.AppConfig,
		Metrics:                cs.HandlerMetrics,
		TotalQueryLimit:        cs.TotalQueryLimit,
	}

	return handler.ProcessStream(stream)
}


func (cs *ChaincodeSupport) Register(stream pb.ChaincodeSupport_RegisterServer) error {
	return cs.HandleChaincodeStream(stream)
}





func (cs *ChaincodeSupport) ExecuteLegacyInit(txParams *ccprovider.TransactionParams, ccName, ccVersion string, input *pb.ChaincodeInput) (*pb.Response, *pb.ChaincodeEvent, error) {
	
	
	
	
	ccid := ccName + ":" + ccVersion

	h, err := cs.Launch(ccid)
	if err != nil {
		return nil, nil, err
	}

	resp, err := cs.execute(pb.ChaincodeMessage_INIT, txParams, ccName, input, h)
	return processChaincodeExecutionResult(txParams.TxID, ccName, resp, err)
}


func (cs *ChaincodeSupport) Execute(txParams *ccprovider.TransactionParams, chaincodeName string, input *pb.ChaincodeInput) (*pb.Response, *pb.ChaincodeEvent, error) {
	resp, err := cs.Invoke(txParams, chaincodeName, input)
	return processChaincodeExecutionResult(txParams.TxID, chaincodeName, resp, err)
}

func processChaincodeExecutionResult(txid, ccName string, resp *pb.ChaincodeMessage, err error) (*pb.Response, *pb.ChaincodeEvent, error) {
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to execute transaction %s", txid)
	}
	if resp == nil {
		return nil, nil, errors.Errorf("nil response from transaction %s", txid)
	}

	if resp.ChaincodeEvent != nil {
		resp.ChaincodeEvent.ChaincodeId = ccName
		resp.ChaincodeEvent.TxId = txid
	}

	switch resp.Type {
	case pb.ChaincodeMessage_COMPLETED:
		res := &pb.Response{}
		err := proto.Unmarshal(resp.Payload, res)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "failed to unmarshal response for transaction %s", txid)
		}
		return res, resp.ChaincodeEvent, nil

	case pb.ChaincodeMessage_ERROR:
		return nil, resp.ChaincodeEvent, errors.Errorf("transaction returned with failure: %s", resp.Payload)

	default:
		return nil, nil, errors.Errorf("unexpected response type %d for transaction %s", resp.Type, txid)
	}
}



func (cs *ChaincodeSupport) Invoke(txParams *ccprovider.TransactionParams, chaincodeName string, input *pb.ChaincodeInput) (*pb.ChaincodeMessage, error) {
	ccid, cctype, err := cs.CheckInvocation(txParams, chaincodeName, input)
	if err != nil {
		return nil, errors.WithMessage(err, "invalid invocation")
	}

	h, err := cs.Launch(ccid)
	if err != nil {
		return nil, err
	}

	return cs.execute(cctype, txParams, chaincodeName, input, h)
}





func (cs *ChaincodeSupport) CheckInvocation(txParams *ccprovider.TransactionParams, chaincodeName string, input *pb.ChaincodeInput) (ccid string, cctype pb.ChaincodeMessage_Type, err error) {
	chaincodeLogger.Debugf("[%s] getting chaincode data for %s on channel %s", shorttxid(txParams.TxID), chaincodeName, txParams.ChannelID)
	cd, err := cs.Lifecycle.ChaincodeDefinition(txParams.ChannelID, chaincodeName, txParams.TXSimulator)
	if err != nil {
		logDevModeError(cs.UserRunsCC)
		return "", 0, errors.Wrapf(err, "[channel %s] failed to get chaincode container info for %s", txParams.ChannelID, chaincodeName)
	}

	if legacyDefinition, ok := cd.(LegacyChaincodeDefinition); ok {
		err = legacyDefinition.ExecuteLegacySecurityChecks()
		if err != nil {
			return "", 0, errors.WithMessagef(err, "[channel %s] failed the chaincode security checks for %s", txParams.ChannelID, chaincodeName)
		}
	}

	needsInitialization := false
	if cd.RequiresInit() {
		
		

		value, err := txParams.TXSimulator.GetState(chaincodeName, InitializedKeyName)
		if err != nil {
			return "", 0, errors.WithMessage(err, "could not get 'initialized' key")
		}

		needsInitialization = !bytes.Equal(value, []byte(cd.CCVersion()))
	}

	
	
	
	if input.IsInit {
		if !cd.RequiresInit() {
			return "", 0, errors.Errorf("chaincode '%s' does not require initialization but called as init", chaincodeName)
		}

		if !needsInitialization {
			return "", 0, errors.Errorf("chaincode '%s' is already initialized but called as init", chaincodeName)
		}

		err = txParams.TXSimulator.SetState(chaincodeName, InitializedKeyName, []byte(cd.CCVersion()))
		if err != nil {
			return "", 0, errors.WithMessage(err, "could not set 'initialized' key")
		}

		return cd.ChaincodeID(), pb.ChaincodeMessage_INIT, nil
	}

	if needsInitialization {
		return "", 0, errors.Errorf("chaincode '%s' has not been initialized for this version, must call as init first", chaincodeName)
	}

	return cd.ChaincodeID(), pb.ChaincodeMessage_TRANSACTION, nil
}


func (cs *ChaincodeSupport) execute(cctyp pb.ChaincodeMessage_Type, txParams *ccprovider.TransactionParams, namespace string, input *pb.ChaincodeInput, h *Handler) (*pb.ChaincodeMessage, error) {
	input.Decorations = txParams.ProposalDecorations

	payload, err := proto.Marshal(input)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to create chaincode message")
	}

	ccMsg := &pb.ChaincodeMessage{
		Type:      cctyp,
		Payload:   payload,
		Txid:      txParams.TxID,
		ChannelId: txParams.ChannelID,
	}

	ccresp, err := h.Execute(txParams, namespace, ccMsg, cs.ExecuteTimeout)
	if err != nil {
		return nil, errors.WithMessage(err, "error sending")
	}

	return ccresp, nil
}

func logDevModeError(userRunsCC bool) {
	if userRunsCC {
		chaincodeLogger.Error("You are attempting to perform an action other than Deploy on Chaincode that is not ready and you are in developer mode. Did you forget to Deploy your chaincode?")
	}
}
