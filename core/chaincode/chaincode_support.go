/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chaincode

import (
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/common/util"
	"github.com/mcc-github/blockchain/core/chaincode/platforms"
	"github.com/mcc-github/blockchain/core/common/ccprovider"
	"github.com/mcc-github/blockchain/core/common/sysccprovider"
	"github.com/mcc-github/blockchain/core/container/ccintf"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/peer"
	pb "github.com/mcc-github/blockchain/protos/peer"
	"github.com/pkg/errors"
)


type Runtime interface {
	Start(ccci *ccprovider.ChaincodeContainerInfo, codePackage []byte) error
	Stop(ccci *ccprovider.ChaincodeContainerInfo) error
}


type Launcher interface {
	Launch(ccci *ccprovider.ChaincodeContainerInfo) error
}


type Lifecycle interface {
	
	ChaincodeDefinition(chaincodeName string, txSim ledger.QueryExecutor) (ccprovider.ChaincodeDefinition, error)

	
	ChaincodeContainerInfo(chainID string, chaincodeID string) (*ccprovider.ChaincodeContainerInfo, error)
}


type ChaincodeSupport struct {
	Keepalive        time.Duration
	ExecuteTimeout   time.Duration
	UserRunsCC       bool
	Runtime          Runtime
	ACLProvider      ACLProvider
	HandlerRegistry  *HandlerRegistry
	Launcher         Launcher
	SystemCCProvider sysccprovider.SystemChaincodeProvider
	Lifecycle        Lifecycle
}


func NewChaincodeSupport(
	config *Config,
	peerAddress string,
	userRunsCC bool,
	caCert []byte,
	certGenerator CertGenerator,
	packageProvider PackageProvider,
	lifecycle Lifecycle,
	aclProvider ACLProvider,
	processor Processor,
	SystemCCProvider sysccprovider.SystemChaincodeProvider,
	platformRegistry *platforms.Registry,
) *ChaincodeSupport {
	cs := &ChaincodeSupport{
		UserRunsCC:       userRunsCC,
		Keepalive:        config.Keepalive,
		ExecuteTimeout:   config.ExecuteTimeout,
		HandlerRegistry:  NewHandlerRegistry(userRunsCC),
		ACLProvider:      aclProvider,
		SystemCCProvider: SystemCCProvider,
		Lifecycle:        lifecycle,
	}

	
	if !config.TLSEnabled {
		certGenerator = nil
	}

	cs.Runtime = &ContainerRuntime{
		CertGenerator:    certGenerator,
		Processor:        processor,
		CACert:           caCert,
		PeerAddress:      peerAddress,
		PlatformRegistry: platformRegistry,
		CommonEnv: []string{
			"CORE_CHAINCODE_LOGGING_LEVEL=" + config.LogLevel,
			"CORE_CHAINCODE_LOGGING_SHIM=" + config.ShimLogLevel,
			"CORE_CHAINCODE_LOGGING_FORMAT=" + config.LogFormat,
		},
	}

	cs.Launcher = &RuntimeLauncher{
		Runtime:         cs.Runtime,
		Registry:        cs.HandlerRegistry,
		PackageProvider: packageProvider,
		StartupTimeout:  config.StartupTimeout,
	}

	return cs
}




func (cs *ChaincodeSupport) LaunchInit(ccci *ccprovider.ChaincodeContainerInfo) error {
	cname := ccci.Name + ":" + ccci.Version
	if cs.HandlerRegistry.Handler(cname) != nil {
		return nil
	}

	return cs.Launcher.Launch(ccci)
}




func (cs *ChaincodeSupport) Launch(chainID, chaincodeName, chaincodeVersion string) (*Handler, error) {
	cname := chaincodeName + ":" + chaincodeVersion
	if h := cs.HandlerRegistry.Handler(cname); h != nil {
		return h, nil
	}

	ccci, err := cs.Lifecycle.ChaincodeContainerInfo(chainID, chaincodeName)
	if err != nil {
		
		if cs.UserRunsCC {
			chaincodeLogger.Error(
				"You are attempting to perform an action other than Deploy on Chaincode that is not ready and you are in developer mode. Did you forget to Deploy your chaincode?",
			)
		}

		return nil, errors.Wrapf(err, "[channel %s] failed to get chaincode container info for %s", chainID, cname)
	}

	if err := cs.Launcher.Launch(ccci); err != nil {
		return nil, errors.Wrapf(err, "[channel %s] could not launch chaincode %s", chainID, cname)
	}

	h := cs.HandlerRegistry.Handler(cname)
	if h == nil {
		return nil, errors.Wrapf(err, "[channel %s] claimed to start chaincode container for %s but could not find handler", chainID, cname)
	}

	return h, nil
}


func (cs *ChaincodeSupport) Stop(ccci *ccprovider.ChaincodeContainerInfo) error {
	err := cs.Runtime.Stop(ccci)
	if err != nil {
		return err
	}

	return nil
}


func (cs *ChaincodeSupport) HandleChaincodeStream(stream ccintf.ChaincodeStream) error {
	handler := &Handler{
		Invoker:                    cs,
		DefinitionGetter:           cs.Lifecycle,
		Keepalive:                  cs.Keepalive,
		Registry:                   cs.HandlerRegistry,
		ACLProvider:                cs.ACLProvider,
		TXContexts:                 NewTransactionContexts(),
		ActiveTransactions:         NewActiveTransactions(),
		SystemCCProvider:           cs.SystemCCProvider,
		SystemCCVersion:            util.GetSysCCVersion(),
		InstantiationPolicyChecker: CheckInstantiationPolicyFunc(ccprovider.CheckInstantiationPolicy),
		QueryResponseBuilder:       &QueryResponseGenerator{MaxResultLimit: 100},
		UUIDGenerator:              UUIDGeneratorFunc(util.GenerateUUID),
		LedgerGetter:               peer.Default,
	}

	return handler.ProcessStream(stream)
}


func (cs *ChaincodeSupport) Register(stream pb.ChaincodeSupport_RegisterServer) error {
	return cs.HandleChaincodeStream(stream)
}


func createCCMessage(messageType pb.ChaincodeMessage_Type, cid string, txid string, cMsg *pb.ChaincodeInput) (*pb.ChaincodeMessage, error) {
	payload, err := proto.Marshal(cMsg)
	if err != nil {
		return nil, err
	}
	ccmsg := &pb.ChaincodeMessage{
		Type:      messageType,
		Payload:   payload,
		Txid:      txid,
		ChannelId: cid,
	}
	return ccmsg, nil
}





func (cs *ChaincodeSupport) ExecuteLegacyInit(txParams *ccprovider.TransactionParams, cccid *ccprovider.CCContext, spec *pb.ChaincodeDeploymentSpec) (*pb.Response, *pb.ChaincodeEvent, error) {
	ccci := ccprovider.DeploymentSpecToChaincodeContainerInfo(spec)
	ccci.Version = cccid.Version

	err := cs.LaunchInit(ccci)
	if err != nil {
		return nil, nil, err
	}

	cname := ccci.Name + ":" + ccci.Version
	h := cs.HandlerRegistry.Handler(cname)
	if h == nil {
		return nil, nil, errors.Wrapf(err, "[channel %s] claimed to start chaincode container for %s but could not find handler", txParams.ChannelID, cname)
	}

	resp, err := cs.execute(pb.ChaincodeMessage_INIT, txParams, cccid, spec.GetChaincodeSpec().Input, h)
	return processChaincodeExecutionResult(txParams.TxID, cccid.Name, resp, err)
}


func (cs *ChaincodeSupport) Execute(txParams *ccprovider.TransactionParams, cccid *ccprovider.CCContext, input *pb.ChaincodeInput) (*pb.Response, *pb.ChaincodeEvent, error) {
	resp, err := cs.Invoke(txParams, cccid, input)
	return processChaincodeExecutionResult(txParams.TxID, cccid.Name, resp, err)
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

func (cs *ChaincodeSupport) InvokeInit(txParams *ccprovider.TransactionParams, cccid *ccprovider.CCContext, input *pb.ChaincodeInput) (*pb.ChaincodeMessage, error) {
	h, err := cs.Launch(txParams.ChannelID, cccid.Name, cccid.Version)
	if err != nil {
		return nil, err
	}

	return cs.execute(pb.ChaincodeMessage_INIT, txParams, cccid, input, h)
}



func (cs *ChaincodeSupport) Invoke(txParams *ccprovider.TransactionParams, cccid *ccprovider.CCContext, input *pb.ChaincodeInput) (*pb.ChaincodeMessage, error) {
	h, err := cs.Launch(txParams.ChannelID, cccid.Name, cccid.Version)
	if err != nil {
		return nil, err
	}

	
	
	
	
	
	
	
	
	
	
	cctype := pb.ChaincodeMessage_TRANSACTION

	return cs.execute(cctype, txParams, cccid, input, h)
}


func (cs *ChaincodeSupport) execute(cctyp pb.ChaincodeMessage_Type, txParams *ccprovider.TransactionParams, cccid *ccprovider.CCContext, input *pb.ChaincodeInput, h *Handler) (*pb.ChaincodeMessage, error) {
	input.Decorations = txParams.ProposalDecorations
	ccMsg, err := createCCMessage(cctyp, txParams.ChannelID, txParams.TxID, input)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to create chaincode message")
	}

	ccresp, err := h.Execute(txParams, cccid, ccMsg, cs.ExecuteTimeout)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("error sending"))
	}

	return ccresp, nil
}