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
	"github.com/mcc-github/blockchain/common/metrics"
	"github.com/mcc-github/blockchain/common/util"
	persistence "github.com/mcc-github/blockchain/core/chaincode/persistence/intf"
	"github.com/mcc-github/blockchain/core/chaincode/platforms"
	"github.com/mcc-github/blockchain/core/common/ccprovider"
	"github.com/mcc-github/blockchain/core/common/sysccprovider"
	"github.com/mcc-github/blockchain/core/container/ccintf"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/peer"
	pb "github.com/mcc-github/blockchain/protos/peer"
	"github.com/pkg/errors"
)

const (
	
	
	
	
	
	InitializedKeyName = "\x00" + string(utf8.MaxRune) + "initialized"
)


type Runtime interface {
	Start(ccci *ccprovider.ChaincodeContainerInfo, codePackage []byte) error
	Stop(ccci *ccprovider.ChaincodeContainerInfo) error
	Wait(ccci *ccprovider.ChaincodeContainerInfo) (int, error)
}


type Launcher interface {
	Launch(ccci *ccprovider.ChaincodeContainerInfo) error
}


type Lifecycle interface {
	
	ChaincodeDefinition(channelID, chaincodeName string, qe ledger.SimpleQueryExecutor) (ccprovider.ChaincodeDefinition, error)

	
	ChaincodeContainerInfo(channelID, chaincodeName string, qe ledger.SimpleQueryExecutor) (*ccprovider.ChaincodeContainerInfo, error)
}


type ChaincodeSupport struct {
	Keepalive              time.Duration
	ExecuteTimeout         time.Duration
	UserRunsCC             bool
	Runtime                Runtime
	ACLProvider            ACLProvider
	HandlerRegistry        *HandlerRegistry
	Launcher               Launcher
	SystemCCProvider       sysccprovider.SystemChaincodeProvider
	Lifecycle              Lifecycle
	AppConfig              ApplicationConfigRetriever
	HandlerMetrics         *HandlerMetrics
	LaunchMetrics          *LaunchMetrics
	DeployedCCInfoProvider ledger.DeployedChaincodeInfoProvider
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
	appConfig ApplicationConfigRetriever,
	metricsProvider metrics.Provider,
	deployedCCInfoProvider ledger.DeployedChaincodeInfoProvider,
) *ChaincodeSupport {
	cs := &ChaincodeSupport{
		UserRunsCC:             userRunsCC,
		Keepalive:              config.Keepalive,
		ExecuteTimeout:         config.ExecuteTimeout,
		HandlerRegistry:        NewHandlerRegistry(userRunsCC),
		ACLProvider:            aclProvider,
		SystemCCProvider:       SystemCCProvider,
		Lifecycle:              lifecycle,
		AppConfig:              appConfig,
		HandlerMetrics:         NewHandlerMetrics(metricsProvider),
		LaunchMetrics:          NewLaunchMetrics(metricsProvider),
		DeployedCCInfoProvider: deployedCCInfoProvider,
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
		Metrics:         cs.LaunchMetrics,
	}

	return cs
}




func (cs *ChaincodeSupport) LaunchInit(ccci *ccprovider.ChaincodeContainerInfo) error {
	if cs.HandlerRegistry.Handler(ccintf.New(ccci.PackageID)) != nil {
		return nil
	}

	return cs.Launcher.Launch(ccci)
}




func (cs *ChaincodeSupport) Launch(chainID string, ccci *ccprovider.ChaincodeContainerInfo) (*Handler, error) {
	ccid := ccintf.New(ccci.PackageID)

	if h := cs.HandlerRegistry.Handler(ccid); h != nil {
		return h, nil
	}

	if err := cs.Launcher.Launch(ccci); err != nil {
		return nil, errors.Wrapf(err, "[channel %s] could not launch chaincode %s", chainID, ccci.PackageID)
	}

	h := cs.HandlerRegistry.Handler(ccid)
	if h == nil {
		return nil, errors.Errorf("[channel %s] claimed to start chaincode container for %s but could not find handler", chainID, ccci.PackageID)
	}

	return h, nil
}


func (cs *ChaincodeSupport) Stop(ccci *ccprovider.ChaincodeContainerInfo) error {
	return cs.Runtime.Stop(ccci)
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
		DeployedCCInfoProvider:     cs.DeployedCCInfoProvider,
		AppConfig:                  cs.AppConfig,
		Metrics:                    cs.HandlerMetrics,
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
	
	
	
	
	ccci.PackageID = persistence.PackageID(ccci.Name + ":" + ccci.Version)

	err := cs.LaunchInit(ccci)
	if err != nil {
		return nil, nil, err
	}

	h := cs.HandlerRegistry.Handler(ccintf.New(ccci.PackageID))
	if h == nil {
		return nil, nil, errors.Wrapf(err, "[channel %s] claimed to start chaincode container for %s but could not find handler", txParams.ChannelID, ccci.PackageID)
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



func (cs *ChaincodeSupport) Invoke(txParams *ccprovider.TransactionParams, cccid *ccprovider.CCContext, input *pb.ChaincodeInput) (*pb.ChaincodeMessage, error) {
	
	var ccci *ccprovider.ChaincodeContainerInfo
	var err error

	if !cs.SystemCCProvider.IsSysCC(cccid.Name) {
		ccci, err = cs.Lifecycle.ChaincodeContainerInfo(txParams.ChannelID, cccid.Name, txParams.TXSimulator)
		if err != nil {
			
			if cs.UserRunsCC {
				chaincodeLogger.Error(
					"You are attempting to perform an action other than Deploy on Chaincode that is not ready and you are in developer mode. Did you forget to Deploy your chaincode?",
				)
			}

			return nil, errors.Wrapf(err, "[channel %s] failed to get chaincode container info for %s", txParams.ChannelID, cccid.Name)
		}
	} else {
		
		ccci = &ccprovider.ChaincodeContainerInfo{
			Version:   util.GetSysCCVersion(),
			Name:      cccid.Name,
			PackageID: persistence.PackageID(cccid.Name + ":" + util.GetSysCCVersion()),
		}
	}

	
	
	cccid.Version = ccci.Version

	h, err := cs.Launch(txParams.ChannelID, ccci)
	if err != nil {
		return nil, err
	}

	isInit, err := cs.CheckInit(txParams, cccid, input)
	if err != nil {
		return nil, err
	}

	cctype := pb.ChaincodeMessage_TRANSACTION
	if isInit {
		cctype = pb.ChaincodeMessage_INIT
	}

	return cs.execute(cctype, txParams, cccid, input, h)
}

func (cs *ChaincodeSupport) CheckInit(txParams *ccprovider.TransactionParams, cccid *ccprovider.CCContext, input *pb.ChaincodeInput) (bool, error) {
	if txParams.ChannelID == "" {
		
		return false, nil
	}

	ac, ok := cs.AppConfig.GetApplicationConfig(txParams.ChannelID)
	if !ok {
		return false, errors.Errorf("could not retrieve application config for channel '%s'", txParams.ChannelID)
	}

	if !ac.Capabilities().LifecycleV20() {
		return false, nil
	}

	if !cccid.InitRequired {
		
		
		return false, nil
	}

	

	value, err := txParams.TXSimulator.GetState(cccid.Name, InitializedKeyName)
	if err != nil {
		return false, errors.WithMessage(err, "could not get 'initialized' key")
	}

	needsInitialization := !bytes.Equal(value, []byte(cccid.Version))

	switch {
	case !input.IsInit && !needsInitialization:
		return false, nil
	case !input.IsInit && needsInitialization:
		return false, errors.Errorf("chaincode '%s' has not been initialized for this version, must call as init first", cccid.Name)
	case input.IsInit && !needsInitialization:
		return false, errors.Errorf("chaincode '%s' is already initialized but called as init", cccid.Name)
	default:
		
		err = txParams.TXSimulator.SetState(cccid.Name, InitializedKeyName, []byte(cccid.Version))
		if err != nil {
			return false, errors.WithMessage(err, "could not set 'initialized' key")
		}
		return true, nil
	}
}


func (cs *ChaincodeSupport) execute(cctyp pb.ChaincodeMessage_Type, txParams *ccprovider.TransactionParams, cccid *ccprovider.CCContext, input *pb.ChaincodeInput, h *Handler) (*pb.ChaincodeMessage, error) {
	input.Decorations = txParams.ProposalDecorations
	ccMsg, err := createCCMessage(cctyp, txParams.ChannelID, txParams.TxID, input)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to create chaincode message")
	}

	ccresp, err := h.Execute(txParams, cccid, ccMsg, cs.ExecuteTimeout)
	if err != nil {
		return nil, errors.WithMessage(err, "error sending")
	}

	return ccresp, nil
}
