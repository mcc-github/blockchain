/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chaincode

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/mcc-github/blockchain/core/ledger/ledgermgmt/ledgermgmttest"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/bccsp/factory"
	"github.com/mcc-github/blockchain/common/channelconfig"
	"github.com/mcc-github/blockchain/common/crypto/tlsgen"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/common/metrics/disabled"
	mockpolicies "github.com/mcc-github/blockchain/common/mocks/policies"
	"github.com/mcc-github/blockchain/common/policies"
	"github.com/mcc-github/blockchain/common/util"
	"github.com/mcc-github/blockchain/core/aclmgmt"
	"github.com/mcc-github/blockchain/core/chaincode/mock"
	cm "github.com/mcc-github/blockchain/core/chaincode/mock"
	persistence "github.com/mcc-github/blockchain/core/chaincode/persistence/intf"
	"github.com/mcc-github/blockchain/core/chaincode/platforms"
	"github.com/mcc-github/blockchain/core/chaincode/platforms/golang"
	"github.com/mcc-github/blockchain/core/chaincode/shim"
	"github.com/mcc-github/blockchain/core/comm"
	"github.com/mcc-github/blockchain/core/common/ccprovider"
	"github.com/mcc-github/blockchain/core/config"
	"github.com/mcc-github/blockchain/core/container"
	"github.com/mcc-github/blockchain/core/container/dockercontroller"
	"github.com/mcc-github/blockchain/core/container/externalbuilders"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/ledger/ledgermgmt"
	ledgermock "github.com/mcc-github/blockchain/core/ledger/mock"
	cut "github.com/mcc-github/blockchain/core/ledger/util"
	"github.com/mcc-github/blockchain/core/peer"
	"github.com/mcc-github/blockchain/core/policy"
	"github.com/mcc-github/blockchain/core/policy/mocks"
	"github.com/mcc-github/blockchain/core/scc"
	"github.com/mcc-github/blockchain/core/scc/lscc"
	"github.com/mcc-github/blockchain/internal/peer/packaging"
	"github.com/mcc-github/blockchain/msp"
	mspmgmt "github.com/mcc-github/blockchain/msp/mgmt"
	msptesttools "github.com/mcc-github/blockchain/msp/mgmt/testtools"
	"github.com/mcc-github/blockchain/protos/common"
	pb "github.com/mcc-github/blockchain/protos/peer"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)


func initPeer(chainIDs ...string) (*cm.Lifecycle, net.Listener, *ChaincodeSupport, func(), error) {
	peerInstance := &peer.Peer{}
	sccp := &scc.Provider{
		Peer:      peerInstance,
		Whitelist: scc.GlobalWhitelist(),
	}

	grpcServer := grpc.NewServer()

	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to start peer listener %s", err)
	}
	_, localPort, err := net.SplitHostPort(lis.Addr().String())
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to get port: %s", err)
	}
	localIP, err := comm.GetLocalIP()
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to get local IP: %s", err)
	}

	peerAddress := net.JoinHostPort(localIP, localPort)

	tempdir, err := ioutil.TempDir("", "chaincode")
	if err != nil {
		panic(fmt.Sprintf("failed to create temporary directory: %s", err))
	}

	lgrInitializer := ledgermgmttest.NewInitializer(filepath.Join(tempdir, "ledgersData"))
	lgrInitializer.Config.HistoryDBConfig = &ledger.HistoryDBConfig{
		Enabled: true,
	}
	peerInstance.LedgerMgr = ledgermgmt.NewLedgerMgr(lgrInitializer)
	ccprovider.SetChaincodesPath(tempdir)
	ca, _ := tlsgen.NewCA()
	pr := platforms.NewRegistry(&golang.Platform{})
	mockAclProvider := &mock.ACLProvider{}
	lsccImpl := lscc.New(sccp, mockAclProvider, peerInstance.GetMSPIDs, newPolicyChecker(peerInstance))
	ml := &cm.Lifecycle{}
	ml.ChaincodeContainerInfoStub = func(_, name string, _ ledger.SimpleQueryExecutor) (*ccprovider.ChaincodeContainerInfo, error) {
		switch name {
		case "lscc":
			return &ccprovider.ChaincodeContainerInfo{
				Name:      "lscc",
				Version:   util.GetSysCCVersion(),
				PackageID: persistence.PackageID("lscc:" + util.GetSysCCVersion()),
			}, nil
		default:
			return &ccprovider.ChaincodeContainerInfo{
				Name:      name,
				Version:   "0",
				PackageID: persistence.PackageID(name + ":0"),
			}, nil
		}
	}
	client, err := docker.NewClientFromEnv()
	if err != nil {
		return nil, nil, nil, nil, err
	}
	globalConfig := &Config{
		TLSEnabled:      false,
		Keepalive:       time.Second,
		StartupTimeout:  3 * time.Minute,
		ExecuteTimeout:  30 * time.Second,
		LogLevel:        "info",
		ShimLogLevel:    "warning",
		LogFormat:       "TEST: [%{module}] %{shortfunc} -> %{level:.4s} %{id:03x}%{color:reset} %{message}",
		TotalQueryLimit: 10000,
	}
	containerRuntime := &ContainerRuntime{
		CACert: ca.CertBytes(),
		ContainerRouter: &container.Router{
			ExternalVM: &externalbuilders.Detector{},
			DockerVM: &dockercontroller.DockerVM{
				PeerID:       "",
				NetworkID:    "",
				BuildMetrics: dockercontroller.NewBuildMetrics(&disabled.Provider{}),
				Client:       client,
				PlatformBuilder: &platforms.Builder{
					Registry: pr,
					Client:   client,
				},
			},
		},
		PeerAddress: peerAddress,
	}
	userRunsCC := false
	metricsProviders := &disabled.Provider{}
	chaincodeHandlerRegistry := NewHandlerRegistry(userRunsCC)
	chaincodeLauncher := &RuntimeLauncher{
		Metrics:         NewLaunchMetrics(metricsProviders),
		PackageProvider: &PackageProviderWrapper{FS: &ccprovider.CCInfoFSImpl{}},
		Registry:        chaincodeHandlerRegistry,
		Runtime:         containerRuntime,
		StartupTimeout:  globalConfig.StartupTimeout,
	}
	chaincodeSupport := &ChaincodeSupport{
		ACLProvider: aclmgmt.NewACLProvider(
			func(string) channelconfig.Resources { return nil },
			newPolicyChecker(peerInstance),
		),
		AppConfig:              peerInstance,
		DeployedCCInfoProvider: &ledgermock.DeployedChaincodeInfoProvider{},
		ExecuteTimeout:         globalConfig.ExecuteTimeout,
		HandlerMetrics:         NewHandlerMetrics(metricsProviders),
		HandlerRegistry:        chaincodeHandlerRegistry,
		Keepalive:              globalConfig.Keepalive,
		Launcher:               chaincodeLauncher,
		Lifecycle:              ml,
		Peer:                   peerInstance,
		Runtime:                containerRuntime,
		SystemCCProvider:       sccp,
		TotalQueryLimit:        globalConfig.TotalQueryLimit,
		UserRunsCC:             userRunsCC,
	}
	pb.RegisterChaincodeSupportServer(grpcServer, chaincodeSupport)

	sccp.RegisterSysCC(lsccImpl)

	sccp.DeploySysCCs(chaincodeSupport)

	for _, id := range chainIDs {
		if err = peer.CreateMockChannel(peerInstance, id); err != nil {
			closeListenerAndSleep(lis)
			return nil, nil, nil, nil, err
		}
		
		if id != util.GetTestChainID() {
			mspmgmt.XXXSetMSPManager(id, mspmgmt.GetManagerForChain(util.GetTestChainID()))
		}
	}

	go grpcServer.Serve(lis)

	cleanup := func() {
		finitPeer(peerInstance, lis, chainIDs...)
		os.RemoveAll(tempdir)
	}

	return ml, lis, chaincodeSupport, cleanup, nil
}

func finitPeer(peerInstance *peer.Peer, lis net.Listener, chainIDs ...string) {
	if lis != nil {
		closeListenerAndSleep(lis)
	}
	for _, c := range chainIDs {
		if lgr := peerInstance.GetLedger(c); lgr != nil {
			lgr.Close()
		}
	}
	peerInstance.LedgerMgr.Close()
	ledgerPath := config.GetPath("peer.fileSystemPath")
	os.RemoveAll(ledgerPath)
	os.RemoveAll(filepath.Join(os.TempDir(), "mcc-github"))
}

func startTxSimulation(peerInstance *peer.Peer, chainID string, txid string) (ledger.TxSimulator, ledger.HistoryQueryExecutor, error) {
	lgr := peerInstance.GetLedger(chainID)
	txsim, err := lgr.NewTxSimulator(txid)
	if err != nil {
		return nil, nil, err
	}
	historyQueryExecutor, err := lgr.NewHistoryQueryExecutor()
	if err != nil {
		return nil, nil, err
	}

	return txsim, historyQueryExecutor, nil
}

func endTxSimulationCDS(peerInstance *peer.Peer, chainID string, txsim ledger.TxSimulator, payload []byte, commit bool, cds *pb.ChaincodeDeploymentSpec, blockNumber uint64) error {
	
	ss, err := signer.Serialize()
	if err != nil {
		return err
	}

	
	lsccid := &pb.ChaincodeID{
		Name:    "lscc",
		Version: util.GetSysCCVersion(),
	}

	
	prop, _, err := protoutil.CreateDeployProposalFromCDS(chainID, cds, ss, nil, nil, nil, nil)
	if err != nil {
		return err
	}

	return endTxSimulation(peerInstance, chainID, lsccid, txsim, payload, commit, prop, blockNumber)
}

func endTxSimulationCIS(peerInstance *peer.Peer, chainID string, ccid *pb.ChaincodeID, txid string, txsim ledger.TxSimulator, payload []byte, commit bool, cis *pb.ChaincodeInvocationSpec, blockNumber uint64) error {
	
	ss, err := signer.Serialize()
	if err != nil {
		return err
	}

	
	prop, returnedTxid, err := protoutil.CreateProposalFromCISAndTxid(txid, common.HeaderType_ENDORSER_TRANSACTION, chainID, cis, ss)
	if err != nil {
		return err
	}
	if returnedTxid != txid {
		return errors.New("txids are not same")
	}

	return endTxSimulation(peerInstance, chainID, ccid, txsim, payload, commit, prop, blockNumber)
}









var _commitLock_ sync.Mutex

func endTxSimulation(peerInstance *peer.Peer, chainID string, ccid *pb.ChaincodeID, txsim ledger.TxSimulator, _ []byte, commit bool, prop *pb.Proposal, blockNumber uint64) error {
	txsim.Done()
	if lgr := peerInstance.GetLedger(chainID); lgr != nil {
		if commit {
			var txSimulationResults *ledger.TxSimulationResults
			var txSimulationBytes []byte
			var err error

			txsim.Done()

			
			if txSimulationResults, err = txsim.GetTxSimulationResults(); err != nil {
				return err
			}
			if txSimulationBytes, err = txSimulationResults.GetPubSimulationBytes(); err != nil {
				return err
			}
			
			resp, err := protoutil.CreateProposalResponse(prop.Header, prop.Payload, &pb.Response{Status: 200},
				txSimulationBytes, nil, ccid, nil, signer)
			if err != nil {
				return err
			}

			
			env, err := protoutil.CreateSignedTx(prop, signer, resp)
			if err != nil {
				return err
			}

			envBytes, err := protoutil.GetBytesEnvelope(env)
			if err != nil {
				return err
			}

			
			bcInfo, err := lgr.GetBlockchainInfo()
			if err != nil {
				return err
			}
			block := protoutil.NewBlock(blockNumber, bcInfo.CurrentBlockHash)
			block.Data.Data = [][]byte{envBytes}
			txsFilter := cut.NewTxValidationFlagsSetValue(len(block.Data.Data), pb.TxValidationCode_VALID)
			block.Metadata.Metadata[common.BlockMetadataIndex_TRANSACTIONS_FILTER] = txsFilter

			

			
			_commitLock_.Lock()
			defer _commitLock_.Unlock()

			blockAndPvtData := &ledger.BlockAndPvtData{
				Block:   block,
				PvtData: make(ledger.TxPvtDataMap),
			}

			
			
			
			
			

			
			seqInBlock := uint64(0)

			if txSimulationResults.PvtSimulationResults != nil {
				blockAndPvtData.PvtData[seqInBlock] = &ledger.TxPvtData{
					SeqInBlock: seqInBlock,
					WriteSet:   txSimulationResults.PvtSimulationResults,
				}
			}

			if err := lgr.CommitLegacy(blockAndPvtData, &ledger.CommitOptions{}); err != nil {
				return err
			}
		}
	}

	return nil
}


func getDeploymentSpec(spec *pb.ChaincodeSpec) (*pb.ChaincodeDeploymentSpec, error) {
	fmt.Printf("getting deployment spec for chaincode spec: %v\n", spec)
	codePackageBytes, err := packaging.NewRegistry(&golang.Platform{}).GetDeploymentPayload(spec.Type.String(), spec.ChaincodeId.Path)
	if err != nil {
		return nil, err
	}

	cdDeploymentSpec := &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec, CodePackage: codePackageBytes}
	return cdDeploymentSpec, nil
}


func getDeployLSCCSpec(chainID string, cds *pb.ChaincodeDeploymentSpec, ccp *common.CollectionConfigPackage) (*pb.ChaincodeInvocationSpec, error) {
	b, err := proto.Marshal(cds)
	if err != nil {
		return nil, err
	}

	var ccpBytes []byte
	if ccp != nil {
		if ccpBytes, err = proto.Marshal(ccp); err != nil {
			return nil, err
		}
	}
	sysCCVers := util.GetSysCCVersion()

	invokeInput := &pb.ChaincodeInput{Args: [][]byte{
		[]byte("deploy"), 
		[]byte(chainID),  
		b,                
	}}

	if ccpBytes != nil {
		
		invokeInput.Args = append(invokeInput.Args, nil, nil, nil, ccpBytes)
	}

	
	lsccSpec := &pb.ChaincodeInvocationSpec{
		ChaincodeSpec: &pb.ChaincodeSpec{
			Type:        pb.ChaincodeSpec_GOLANG,
			ChaincodeId: &pb.ChaincodeID{Name: "lscc", Version: sysCCVers},
			Input:       invokeInput,
		}}

	return lsccSpec, nil
}


func deploy(chainID string, cccid *ccprovider.CCContext, spec *pb.ChaincodeSpec, blockNumber uint64, chaincodeSupport *ChaincodeSupport) (resp *pb.Response, err error) {
	
	cdDeploymentSpec, err := getDeploymentSpec(spec)
	if err != nil {
		return nil, err
	}
	return deploy2(chainID, cccid, cdDeploymentSpec, nil, blockNumber, chaincodeSupport)
}

func deployWithCollectionConfigs(chainID string, cccid *ccprovider.CCContext, spec *pb.ChaincodeSpec,
	collectionConfigPkg *common.CollectionConfigPackage, blockNumber uint64, chaincodeSupport *ChaincodeSupport) (resp *pb.Response, err error) {
	
	cdDeploymentSpec, err := getDeploymentSpec(spec)
	if err != nil {
		return nil, err
	}
	return deploy2(chainID, cccid, cdDeploymentSpec, collectionConfigPkg, blockNumber, chaincodeSupport)
}

func deploy2(chainID string, cccid *ccprovider.CCContext, chaincodeDeploymentSpec *pb.ChaincodeDeploymentSpec,
	collectionConfigPkg *common.CollectionConfigPackage, blockNumber uint64, chaincodeSupport *ChaincodeSupport) (resp *pb.Response, err error) {
	cis, err := getDeployLSCCSpec(chainID, chaincodeDeploymentSpec, collectionConfigPkg)
	if err != nil {
		return nil, fmt.Errorf("Error creating lscc spec : %s\n", err)
	}

	uuid := util.GenerateUUID()
	txsim, hqe, err := startTxSimulation(chaincodeSupport.Peer, chainID, uuid)
	sprop, prop := protoutil.MockSignedEndorserProposal2OrPanic(chainID, cis.ChaincodeSpec, signer)
	txParams := &ccprovider.TransactionParams{
		TxID:                 uuid,
		ChannelID:            chainID,
		TXSimulator:          txsim,
		HistoryQueryExecutor: hqe,
		SignedProp:           sprop,
		Proposal:             prop,
	}
	if err != nil {
		return nil, fmt.Errorf("Failed to get handle to simulator: %s ", err)
	}

	defer func() {
		
		if err == nil {
			
			err = endTxSimulationCDS(chaincodeSupport.Peer, chainID, txsim, []byte("deployed"), true, chaincodeDeploymentSpec, blockNumber)
		} else {
			
			endTxSimulationCDS(chaincodeSupport.Peer, chainID, txsim, []byte("deployed"), false, chaincodeDeploymentSpec, blockNumber)
		}
	}()

	
	ccprovider.PutChaincodeIntoFS(chaincodeDeploymentSpec)

	sysCCVers := util.GetSysCCVersion()
	lsccid := &ccprovider.CCContext{
		Name:    cis.ChaincodeSpec.ChaincodeId.Name,
		Version: sysCCVers,
	}

	
	if _, _, err = chaincodeSupport.Execute(txParams, lsccid, cis.ChaincodeSpec.Input); err != nil {
		return nil, fmt.Errorf("Error deploying chaincode (1): %s", err)
	}

	if resp, _, err = chaincodeSupport.ExecuteLegacyInit(txParams, cccid, chaincodeDeploymentSpec); err != nil {
		return nil, fmt.Errorf("Error deploying chaincode(2): %s", err)
	}

	return resp, nil
}


func invoke(chainID string, spec *pb.ChaincodeSpec, blockNumber uint64, creator []byte, chaincodeSupport *ChaincodeSupport) (ccevt *pb.ChaincodeEvent, uuid string, retval []byte, err error) {
	return invokeWithVersion(chainID, spec.GetChaincodeId().Version, spec, blockNumber, creator, chaincodeSupport)
}


func invokeWithVersion(chainID string, version string, spec *pb.ChaincodeSpec, blockNumber uint64, creator []byte, chaincodeSupport *ChaincodeSupport) (ccevt *pb.ChaincodeEvent, uuid string, retval []byte, err error) {
	cdInvocationSpec := &pb.ChaincodeInvocationSpec{ChaincodeSpec: spec}

	
	uuid = util.GenerateUUID()

	txsim, hqe, err := startTxSimulation(chaincodeSupport.Peer, chainID, uuid)
	if err != nil {
		return nil, uuid, nil, fmt.Errorf("Failed to get handle to simulator: %s ", err)
	}

	defer func() {
		
		if err == nil {
			
			err = endTxSimulationCIS(chaincodeSupport.Peer, chainID, spec.ChaincodeId, uuid, txsim, []byte("invoke"), true, cdInvocationSpec, blockNumber)
		} else {
			
			endTxSimulationCIS(chaincodeSupport.Peer, chainID, spec.ChaincodeId, uuid, txsim, []byte("invoke"), false, cdInvocationSpec, blockNumber)
		}
	}()

	if len(creator) == 0 {
		creator = []byte("Admin")
	}
	sprop, prop := protoutil.MockSignedEndorserProposalOrPanic(chainID, spec, creator, []byte("msg1"))
	cccid := &ccprovider.CCContext{
		Name:    cdInvocationSpec.ChaincodeSpec.ChaincodeId.Name,
		Version: version,
	}
	var resp *pb.Response
	txParams := &ccprovider.TransactionParams{
		TxID:                 uuid,
		ChannelID:            chainID,
		TXSimulator:          txsim,
		HistoryQueryExecutor: hqe,
		SignedProp:           sprop,
		Proposal:             prop,
	}

	resp, ccevt, err = chaincodeSupport.Execute(txParams, cccid, cdInvocationSpec.ChaincodeSpec.Input)
	if err != nil {
		return nil, uuid, nil, fmt.Errorf("Error invoking chaincode: %s", err)
	}
	if resp.Status != shim.OK {
		return nil, uuid, nil, fmt.Errorf("Error invoking chaincode: %s", resp.Message)
	}

	return ccevt, uuid, resp.Payload, err
}

func closeListenerAndSleep(l net.Listener) {
	if l != nil {
		l.Close()
		time.Sleep(2 * time.Second)
	}
}


func checkFinalState(peerInstance *peer.Peer, chainID string, cccid *ccprovider.CCContext, a int, b int) error {
	txid := util.GenerateUUID()
	txsim, _, err := startTxSimulation(peerInstance, chainID, txid)
	if err != nil {
		return fmt.Errorf("Failed to get handle to simulator: %s ", err)
	}

	defer txsim.Done()

	cName := cccid.Name + ":" + cccid.Version

	
	var Aval, Bval int
	resbytes, resErr := txsim.GetState(cccid.Name, "a")
	if resErr != nil {
		return fmt.Errorf("Error retrieving state from ledger for <%s>: %s", cName, resErr)
	}
	Aval, resErr = strconv.Atoi(string(resbytes))
	if resErr != nil {
		return fmt.Errorf("Error retrieving state from ledger for <%s>: %s", cName, resErr)
	}
	if Aval != a {
		return fmt.Errorf("Incorrect result. Aval %d != %d <%s>", Aval, a, cName)
	}

	resbytes, resErr = txsim.GetState(cccid.Name, "b")
	if resErr != nil {
		return fmt.Errorf("Error retrieving state from ledger for <%s>: %s", cName, resErr)
	}
	Bval, resErr = strconv.Atoi(string(resbytes))
	if resErr != nil {
		return fmt.Errorf("Error retrieving state from ledger for <%s>: %s", cName, resErr)
	}
	if Bval != b {
		return fmt.Errorf("Incorrect result. Bval %d != %d <%s>", Bval, b, cName)
	}

	
	fmt.Printf("Aval = %d, Bval = %d\n", Aval, Bval)
	return nil
}

const (
	chaincodeExample02GolangPath = "github.com/mcc-github/blockchain/core/chaincode/testdata/src/chaincodes/example02"
	chaincodePassthruGolangPath  = "github.com/mcc-github/blockchain/core/chaincode/testdata/src/chaincodes/passthru"
)


func TestChaincodeInvokeChaincode(t *testing.T) {
	channel := util.GetTestChainID()
	channel2 := channel + "2"
	ml, lis, chaincodeSupport, cleanup, err := initPeer(channel, channel2)
	if err != nil {
		t.Fail()
		t.Logf("Error creating peer: %s", err)
	}
	defer cleanup()
	defer closeListenerAndSleep(lis)

	peerInstance := chaincodeSupport.Peer

	var nextBlockNumber1 uint64 = 1
	var nextBlockNumber2 uint64 = 1

	chaincode1Name := "cc_go_" + util.GenerateUUID()
	chaincode2Name := "cc_go_" + util.GenerateUUID()

	initialA, initialB := 100, 200

	
	ml.ChaincodeDefinitionReturns(&cm.ChaincodeDefinition{}, nil)
	_, cccid1, err := deployChaincode(
		chaincode1Name,
		"0",
		pb.ChaincodeSpec_GOLANG,
		chaincodeExample02GolangPath,
		util.ToChaincodeArgs("init", "a", strconv.Itoa(initialA), "b", strconv.Itoa(initialB)),
		channel,
		nextBlockNumber1,
		chaincodeSupport,
	)
	defer stopChaincode(cccid1, chaincodeSupport)
	require.NoErrorf(t, err, "error initializing chaincode %s: %s", chaincode1Name, err)
	nextBlockNumber1++
	time.Sleep(time.Second)

	
	chaincode2Version := "0"
	chaincode2Type := pb.ChaincodeSpec_GOLANG
	chaincode2Path := chaincodePassthruGolangPath

	
	_, cccid2, err := deployChaincode(
		chaincode2Name,
		chaincode2Version,
		chaincode2Type,
		chaincode2Path,
		util.ToChaincodeArgs("init"),
		channel,
		nextBlockNumber1,
		chaincodeSupport,
	)
	defer stopChaincode(cccid2, chaincodeSupport)
	require.NoErrorf(t, err, "Error initializing chaincode %s: %s", chaincode2Name, err)
	nextBlockNumber1++
	time.Sleep(time.Second)

	
	
	chaincode2InvokeSpec := &pb.ChaincodeSpec{
		Type: chaincode2Type,
		ChaincodeId: &pb.ChaincodeID{
			Name:    chaincode2Name,
			Version: chaincode2Version,
		},
		Input: &pb.ChaincodeInput{
			Args: util.ToChaincodeArgs(cccid1.Name, "invoke", "a", "b", "10", ""),
		},
	}
	_, _, _, err = invoke(channel, chaincode2InvokeSpec, nextBlockNumber1, []byte("Alice"), chaincodeSupport)
	require.NoErrorf(t, err, "error invoking %s: %s", chaincode2Name, err)
	nextBlockNumber1++

	
	err = checkFinalState(peerInstance, channel, cccid1, initialA-10, initialB+10)
	require.NoErrorf(t, err, "incorrect final state after transaction for %s: %s", chaincode1Name, err)

	
	
	
	
	pm := peerInstance.GetPolicyManager(channel)
	pm.(*mockpolicies.Manager).PolicyMap = map[string]policies.Policy{
		policies.ChannelApplicationWriters: &CreatorPolicy{Creators: [][]byte{[]byte("Alice")}},
	}

	pm = peerInstance.GetPolicyManager(channel2)
	pm.(*mockpolicies.Manager).PolicyMap = map[string]policies.Policy{
		policies.ChannelApplicationWriters: &CreatorPolicy{Creators: [][]byte{[]byte("Alice"), []byte("Bob")}},
	}

	
	_, cccid3, err := deployChaincode(
		chaincode2Name,
		chaincode2Version,
		chaincode2Type,
		chaincode2Path,
		util.ToChaincodeArgs("init"),
		channel2,
		nextBlockNumber2,
		chaincodeSupport,
	)
	defer stopChaincode(cccid3, chaincodeSupport)
	require.NoErrorf(t, err, "error initializing chaincode %s/%s: %s", chaincode2Name, channel2, err)
	nextBlockNumber2++
	time.Sleep(time.Second)

	chaincode2InvokeSpec = &pb.ChaincodeSpec{
		Type: chaincode2Type,
		ChaincodeId: &pb.ChaincodeID{
			Name:    chaincode2Name,
			Version: chaincode2Version,
		},
		Input: &pb.ChaincodeInput{
			Args: util.ToChaincodeArgs(cccid1.Name, "invoke", "a", "b", "10", channel),
		},
	}

	
	_, _, _, err = invoke(channel2, chaincode2InvokeSpec, nextBlockNumber2, []byte("Bob"), chaincodeSupport)
	require.Errorf(t, err, "as Bob, invoking <%s/%s> via <%s/%s> should fail, but it succeeded.", cccid1.Name, channel, chaincode2Name, channel2)
	assert.True(t, strings.Contains(err.Error(), "[Creator not recognized [Bob]]"))

	
	_, _, _, err = invoke(channel2, chaincode2InvokeSpec, nextBlockNumber2, []byte("Alice"), chaincodeSupport)
	require.NoError(t, err, "as Alice, invoking <%s/%s> via <%s/%s> should should of succeeded, but it failed: %s", cccid1.Name, channel, chaincode2Name, channel2, err)
	nextBlockNumber2++
}

func stopChaincode(chaincodeCtx *ccprovider.CCContext, chaincodeSupport *ChaincodeSupport) {
	chaincodeSupport.Runtime.Stop(&ccprovider.ChaincodeContainerInfo{
		PackageID: persistence.PackageID(chaincodeCtx.Name + ":" + chaincodeCtx.Version),
		Name:      chaincodeCtx.Name,
		Version:   chaincodeCtx.Version,
		Type:      "GOLANG",
	})
}



func TestChaincodeInvokeChaincodeErrorCase(t *testing.T) {
	chainID := util.GetTestChainID()

	ml, _, chaincodeSupport, cleanup, err := initPeer(chainID)
	if err != nil {
		t.Fail()
		t.Logf("Error creating peer: %s", err)
	}
	defer cleanup()

	ml.ChaincodeDefinitionReturns(&cm.ChaincodeDefinition{}, nil)

	
	cID1 := &pb.ChaincodeID{Name: "example02", Path: chaincodeExample02GolangPath, Version: "0"}
	f := "init"
	args := util.ToChaincodeArgs(f, "a", "100", "b", "200")

	spec1 := &pb.ChaincodeSpec{Type: 1, ChaincodeId: cID1, Input: &pb.ChaincodeInput{Args: args}}

	cccid1 := &ccprovider.CCContext{
		Name:    "example02",
		Version: "0",
	}

	var nextBlockNumber uint64 = 1
	defer chaincodeSupport.Runtime.Stop(&ccprovider.ChaincodeContainerInfo{
		PackageID: persistence.PackageID(cID1.Name + ":" + cID1.Version),
		Name:      cID1.Name,
		Version:   cID1.Version,
		Path:      cID1.Path,
		Type:      "GOLANG",
	})

	_, err = deploy(chainID, cccid1, spec1, nextBlockNumber, chaincodeSupport)
	nextBlockNumber++
	ccID1 := spec1.ChaincodeId.Name
	if err != nil {
		t.Fail()
		t.Logf("Error initializing chaincode %s(%s)", ccID1, err)
		return
	}

	time.Sleep(time.Second)

	
	cID2 := &pb.ChaincodeID{Name: "pthru", Path: chaincodePassthruGolangPath, Version: "0"}
	f = "init"
	args = util.ToChaincodeArgs(f)

	spec2 := &pb.ChaincodeSpec{Type: 1, ChaincodeId: cID2, Input: &pb.ChaincodeInput{Args: args}}

	cccid2 := &ccprovider.CCContext{
		Name:    "pthru",
		Version: "0",
	}

	defer chaincodeSupport.Runtime.Stop(&ccprovider.ChaincodeContainerInfo{
		PackageID: persistence.PackageID(cID2.Name + ":" + cID2.Version),
		Name:      cID2.Name,
		Version:   cID2.Version,
		Path:      cID2.Path,
		Type:      "GOLANG",
	})
	_, err = deploy(chainID, cccid2, spec2, nextBlockNumber, chaincodeSupport)
	nextBlockNumber++
	ccID2 := spec2.ChaincodeId.Name
	if err != nil {
		t.Fail()
		t.Logf("Error initializing chaincode %s(%s)", ccID2, err)
		return
	}

	time.Sleep(time.Second)

	
	f = ccID1
	args = util.ToChaincodeArgs(f, "invoke", "a", "")

	spec2 = &pb.ChaincodeSpec{Type: 1, ChaincodeId: cID2, Input: &pb.ChaincodeInput{Args: args}}
	
	_, _, _, err = invoke(chainID, spec2, nextBlockNumber, []byte("Alice"), chaincodeSupport)

	if err == nil {
		t.Fail()
		t.Logf("Error invoking <%s>: %s", ccID2, err)
		return
	}

	if !strings.Contains(err.Error(), "Incorrect number of arguments. Expecting 3") {
		t.Fail()
		t.Logf("Unexpected error %s", err)
		return
	}
}

func TestChaincodeInit(t *testing.T) {
	chainID := util.GetTestChainID()

	_, _, chaincodeSupport, cleanup, err := initPeer(chainID)
	if err != nil {
		t.Fail()
		t.Logf("Error creating peer: %s", err)
	}

	defer cleanup()

	url := "github.com/mcc-github/blockchain/core/chaincode/testdata/src/chaincodes/init_private_data"
	cID := &pb.ChaincodeID{Name: "init_pvtdata", Path: url, Version: "0"}

	f := "init"
	args := util.ToChaincodeArgs(f)

	spec := &pb.ChaincodeSpec{Type: 1, ChaincodeId: cID, Input: &pb.ChaincodeInput{Args: args}}

	cccid := &ccprovider.CCContext{
		Name:    "tmap",
		Version: "0",
	}

	defer chaincodeSupport.Runtime.Stop(&ccprovider.ChaincodeContainerInfo{
		PackageID: persistence.PackageID(cID.Name + ":" + cID.Version),
		Name:      cID.Name,
		Version:   cID.Version,
		Path:      cID.Path,
		Type:      "GOLANG",
	})

	var nextBlockNumber uint64 = 1
	_, err = deploy(chainID, cccid, spec, nextBlockNumber, chaincodeSupport)
	assert.Contains(t, err.Error(), "private data APIs are not allowed in chaincode Init")

	url = "github.com/mcc-github/blockchain/core/chaincode/testdata/src/chaincodes/init_public_data"
	cID = &pb.ChaincodeID{Name: "init_public_data", Path: url, Version: "0"}

	f = "init"
	args = util.ToChaincodeArgs(f)

	spec = &pb.ChaincodeSpec{Type: 1, ChaincodeId: cID, Input: &pb.ChaincodeInput{Args: args}}

	cccid = &ccprovider.CCContext{
		Name:    "tmap",
		Version: "0",
	}

	resp, err := deploy(chainID, cccid, spec, nextBlockNumber, chaincodeSupport)
	assert.NoError(t, err)
	
	
	assert.Equal(t, int32(shim.OK), resp.Status)
}


func TestQueries(t *testing.T) {
	
	

	chainID := util.GetTestChainID()

	_, _, chaincodeSupport, cleanup, err := initPeer(chainID)
	if err != nil {
		t.Fail()
		t.Logf("Error creating peer: %s", err)
	}

	defer cleanup()

	url := "github.com/mcc-github/blockchain/core/chaincode/testdata/src/chaincodes/map"
	cID := &pb.ChaincodeID{Name: "tmap", Path: url, Version: "0"}

	f := "init"
	args := util.ToChaincodeArgs(f)

	spec := &pb.ChaincodeSpec{Type: 1, ChaincodeId: cID, Input: &pb.ChaincodeInput{Args: args}}

	cccid := &ccprovider.CCContext{
		Name:    "tmap",
		Version: "0",
	}

	defer chaincodeSupport.Runtime.Stop(&ccprovider.ChaincodeContainerInfo{
		PackageID: persistence.PackageID(cID.Name + ":" + cID.Version),
		Name:      cID.Name,
		Version:   cID.Version,
		Path:      cID.Path,
		Type:      "GOLANG",
	})

	var nextBlockNumber uint64 = 1
	_, err = deploy(chainID, cccid, spec, nextBlockNumber, chaincodeSupport)
	nextBlockNumber++
	ccID := spec.ChaincodeId.Name
	if err != nil {
		t.Fail()
		t.Logf("Error initializing chaincode %s(%s)", ccID, err)
		return
	}

	var keys []interface{}
	
	
	for i := 1; i <= 101; i++ {
		f = "put"

		
		owner := "tom"
		if i%2 == 0 {
			owner = "jerry"
		}

		
		color := "blue"
		if i == 12 {
			color = "red"
		}

		key := fmt.Sprintf("marble%03d", i)
		argsString := fmt.Sprintf("{\"docType\":\"marble\",\"name\":\"%s\",\"color\":\"%s\",\"size\":35,\"owner\":\"%s\"}", key, color, owner)
		args = util.ToChaincodeArgs(f, key, argsString)
		spec = &pb.ChaincodeSpec{Type: 1, ChaincodeId: cID, Input: &pb.ChaincodeInput{Args: args}}
		_, _, _, err = invoke(chainID, spec, nextBlockNumber, nil, chaincodeSupport)
		nextBlockNumber++

		if err != nil {
			t.Fail()
			t.Logf("Error invoking <%s>: %s", ccID, err)
			return
		}

	}

	
	f = "keys"
	args = util.ToChaincodeArgs(f, "marble001", "marble011")

	spec = &pb.ChaincodeSpec{Type: 1, ChaincodeId: cID, Input: &pb.ChaincodeInput{Args: args}}
	_, _, retval, err := invoke(chainID, spec, nextBlockNumber, nil, chaincodeSupport)
	nextBlockNumber++
	if err != nil {
		t.Fail()
		t.Logf("Error invoking <%s>: %s", ccID, err)
		return
	}

	err = json.Unmarshal(retval, &keys)
	assert.NoError(t, err)
	if len(keys) != 10 {
		t.Fail()
		t.Logf("Error detected with the range query, should have returned 10 but returned %v", len(keys))
		return
	}

	
	

	
	origTimeout := chaincodeSupport.ExecuteTimeout
	chaincodeSupport.ExecuteTimeout = time.Duration(1) * time.Second

	
	args = util.ToChaincodeArgs(f, "marble001", "marble002", "2000")

	spec = &pb.ChaincodeSpec{Type: 1, ChaincodeId: cID, Input: &pb.ChaincodeInput{Args: args}}
	_, _, _, err = invoke(chainID, spec, nextBlockNumber, nil, chaincodeSupport)
	if err == nil {
		t.Fail()
		t.Logf("expected timeout error but succeeded")
		return
	}

	
	chaincodeSupport.ExecuteTimeout = origTimeout

	
	
	
	f = "keys"
	args = util.ToChaincodeArgs(f, "marble001", "marble102")

	spec = &pb.ChaincodeSpec{Type: 1, ChaincodeId: cID, Input: &pb.ChaincodeInput{Args: args}}
	_, _, retval, err = invoke(chainID, spec, nextBlockNumber, nil, chaincodeSupport)
	nextBlockNumber++
	if err != nil {
		t.Fail()
		t.Logf("Error invoking <%s>: %s", ccID, err)
		return
	}

	
	err = json.Unmarshal(retval, &keys)
	assert.NoError(t, err)

	
	
	if len(keys) != 101 {
		t.Fail()
		t.Logf("Error detected with the range query, should have returned 101 but returned %v", len(keys))
		return
	}

	
	
	
	f = "keys"
	args = util.ToChaincodeArgs(f, "", "")

	spec = &pb.ChaincodeSpec{Type: 1, ChaincodeId: cID, Input: &pb.ChaincodeInput{Args: args}}
	_, _, retval, err = invoke(chainID, spec, nextBlockNumber, nil, chaincodeSupport)
	nextBlockNumber++
	if err != nil {
		t.Fail()
		t.Logf("Error invoking <%s>: %s", ccID, err)
		return
	}

	
	err = json.Unmarshal(retval, &keys)
	assert.NoError(t, err)

	
	
	if len(keys) != 101 {
		t.Fail()
		t.Logf("Error detected with the range query, should have returned 101 but returned %v", len(keys))
		return
	}

	type PageResponse struct {
		Bookmark string   `json:"bookmark"`
		Keys     []string `json:"keys"`
	}

	
	f = "keysByPage"
	args = util.ToChaincodeArgs(f, "marble001", "marble011", "2", "")

	spec = &pb.ChaincodeSpec{Type: 1, ChaincodeId: cID, Input: &pb.ChaincodeInput{Args: args}}
	_, _, retval, err = invoke(chainID, spec, nextBlockNumber, nil, chaincodeSupport)
	nextBlockNumber++
	if err != nil {
		t.Fail()
		t.Logf("Error invoking <%s>: %s", ccID, err)
		return
	}
	queryPage := &PageResponse{}

	json.Unmarshal(retval, &queryPage)

	expectedResult := []string{"marble001", "marble002"}

	if !reflect.DeepEqual(expectedResult, queryPage.Keys) {
		t.Fail()
		t.Logf("Error detected with the paginated range query. Returned: %v  should have returned: %v", queryPage.Keys, expectedResult)
		return
	}

	
	args = util.ToChaincodeArgs(f, "marble001", "marble011", "2", queryPage.Bookmark)
	spec = &pb.ChaincodeSpec{Type: 1, ChaincodeId: cID, Input: &pb.ChaincodeInput{Args: args}}
	_, _, retval, err = invoke(chainID, spec, nextBlockNumber, nil, chaincodeSupport)
	nextBlockNumber++
	if err != nil {
		t.Fail()
		t.Logf("Error invoking <%s>: %s", ccID, err)
		return
	}

	json.Unmarshal(retval, &queryPage)

	expectedResult = []string{"marble003", "marble004"}

	if !reflect.DeepEqual(expectedResult, queryPage.Keys) {
		t.Fail()
		t.Logf("Error detected with the paginated range query second page. Returned: %v  should have returned: %v    %v", queryPage.Keys, expectedResult, queryPage.Bookmark)
		return
	}

	
	f = "put"
	args = util.ToChaincodeArgs(f, "marble012", "{\"docType\":\"marble\",\"name\":\"marble012\",\"color\":\"red\",\"size\":30,\"owner\":\"jerry\"}")
	spec = &pb.ChaincodeSpec{Type: 1, ChaincodeId: cID, Input: &pb.ChaincodeInput{Args: args}}
	_, _, _, err = invoke(chainID, spec, nextBlockNumber, nil, chaincodeSupport)
	nextBlockNumber++
	if err != nil {
		t.Fail()
		t.Logf("Error invoking <%s>: %s", ccID, err)
		return
	}

	f = "put"
	args = util.ToChaincodeArgs(f, "marble012", "{\"docType\":\"marble\",\"name\":\"marble012\",\"color\":\"red\",\"size\":30,\"owner\":\"jerry\"}")
	spec = &pb.ChaincodeSpec{Type: 1, ChaincodeId: cID, Input: &pb.ChaincodeInput{Args: args}}
	_, _, _, err = invoke(chainID, spec, nextBlockNumber, nil, chaincodeSupport)
	nextBlockNumber++
	if err != nil {
		t.Fail()
		t.Logf("Error invoking <%s>: %s", ccID, err)
		return
	}

	
	f = "history"
	args = util.ToChaincodeArgs(f, "marble012")
	spec = &pb.ChaincodeSpec{Type: 1, ChaincodeId: cID, Input: &pb.ChaincodeInput{Args: args}}
	_, _, retval, err = invoke(chainID, spec, nextBlockNumber, nil, chaincodeSupport)
	nextBlockNumber++
	if err != nil {
		t.Fail()
		t.Logf("Error invoking <%s>: %s", ccID, err)
		return
	}

	var history []interface{}
	err = json.Unmarshal(retval, &history)
	assert.NoError(t, err)
	if len(history) != 3 {
		t.Fail()
		t.Logf("Error detected with the history query, should have returned 3 but returned %v", len(history))
		return
	}
}

func TestMain(m *testing.M) {
	var err error

	msptesttools.LoadMSPSetupForTesting()
	signer, err = mspmgmt.GetLocalMSP().GetDefaultSigningIdentity()
	if err != nil {
		fmt.Print("Could not initialize msp/signer")
		os.Exit(-1)
		return
	}

	setupTestConfig()
	flogging.ActivateSpec("chaincode=debug")
	os.Exit(m.Run())
}

func setupTestConfig() {
	flag.Parse()

	
	viper.SetEnvPrefix("CORE")
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetConfigName("chaincodetest") 
	viper.AddConfigPath("./")            
	err := viper.ReadInConfig()          
	if err != nil {                      
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	
	err = factory.InitFactories(nil)
	if err != nil {
		panic(fmt.Errorf("Could not initialize BCCSP Factories [%s]", err))
	}
}

func deployChaincode(name string, version string, chaincodeType pb.ChaincodeSpec_Type, path string, args [][]byte, channel string, nextBlockNumber uint64, chaincodeSupport *ChaincodeSupport) (*pb.Response, *ccprovider.CCContext, error) {
	chaincodeSpec := &pb.ChaincodeSpec{
		ChaincodeId: &pb.ChaincodeID{
			Name:    name,
			Version: version,
			Path:    path,
		},
		Type: chaincodeType,
		Input: &pb.ChaincodeInput{
			Args: args,
		},
	}

	chaincodeCtx := &ccprovider.CCContext{
		Name:    name,
		Version: version,
	}

	result, err := deploy(channel, chaincodeCtx, chaincodeSpec, nextBlockNumber, chaincodeSupport)
	if err != nil {
		return nil, chaincodeCtx, fmt.Errorf("Error deploying <%s:%s>: %s", name, version, err)
	}
	return result, chaincodeCtx, nil
}

var signer msp.SigningIdentity

type CreatorPolicy struct {
	Creators [][]byte
}


func (c *CreatorPolicy) Evaluate(signatureSet []*protoutil.SignedData) error {
	for _, value := range c.Creators {
		if bytes.Equal(signatureSet[0].Identity, value) {
			return nil
		}
	}
	return fmt.Errorf("Creator not recognized [%s]", string(signatureSet[0].Identity))
}

func newPolicyChecker(peerInstance *peer.Peer) policy.PolicyChecker {
	return policy.NewPolicyChecker(
		policies.PolicyManagerGetterFunc(peerInstance.GetPolicyManager),
		&mocks.MockIdentityDeserializer{
			Identity: []byte("Admin"),
			Msg:      []byte("msg1"),
		},
		&mocks.MockMSPPrincipalGetter{Principal: []byte("Admin")},
	)
}
