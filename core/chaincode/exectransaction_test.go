/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chaincode

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/bccsp/factory"
	"github.com/mcc-github/blockchain/common/channelconfig"
	"github.com/mcc-github/blockchain/common/crypto/tlsgen"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/common/metrics/disabled"
	mc "github.com/mcc-github/blockchain/common/mocks/config"
	mockpolicies "github.com/mcc-github/blockchain/common/mocks/policies"
	"github.com/mcc-github/blockchain/common/policies"
	"github.com/mcc-github/blockchain/common/util"
	"github.com/mcc-github/blockchain/core/aclmgmt"
	aclmocks "github.com/mcc-github/blockchain/core/aclmgmt/mocks"
	"github.com/mcc-github/blockchain/core/chaincode/accesscontrol"
	cm "github.com/mcc-github/blockchain/core/chaincode/mock"
	persistence "github.com/mcc-github/blockchain/core/chaincode/persistence/intf"
	"github.com/mcc-github/blockchain/core/chaincode/platforms"
	"github.com/mcc-github/blockchain/core/chaincode/platforms/golang"
	"github.com/mcc-github/blockchain/core/chaincode/shim"
	"github.com/mcc-github/blockchain/core/common/ccprovider"
	"github.com/mcc-github/blockchain/core/config"
	"github.com/mcc-github/blockchain/core/container"
	"github.com/mcc-github/blockchain/core/container/dockercontroller"
	"github.com/mcc-github/blockchain/core/container/inproccontroller"
	"github.com/mcc-github/blockchain/core/ledger"
	ledgermock "github.com/mcc-github/blockchain/core/ledger/mock"
	cut "github.com/mcc-github/blockchain/core/ledger/util"
	cmp "github.com/mcc-github/blockchain/core/mocks/peer"
	"github.com/mcc-github/blockchain/core/peer"
	"github.com/mcc-github/blockchain/core/policy"
	"github.com/mcc-github/blockchain/core/policy/mocks"
	"github.com/mcc-github/blockchain/core/scc"
	"github.com/mcc-github/blockchain/core/scc/lscc"
	"github.com/mcc-github/blockchain/msp"
	mspmgmt "github.com/mcc-github/blockchain/msp/mgmt"
	msptesttools "github.com/mcc-github/blockchain/msp/mgmt/testtools"
	"github.com/mcc-github/blockchain/protos/common"
	pb "github.com/mcc-github/blockchain/protos/peer"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	ma "github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)


func initPeer(chainIDs ...string) (*cm.Lifecycle, net.Listener, *ChaincodeSupport, func(), error) {
	
	finitPeer(nil, chainIDs...)

	msi := &cmp.MockSupportImpl{
		GetApplicationConfigRv:     &mc.MockApplication{CapabilitiesRv: &mc.MockApplicationCapabilities{}},
		GetApplicationConfigBoolRv: true,
	}

	ipRegistry := inproccontroller.NewRegistry()
	sccp := &scc.Provider{Peer: peer.Default, PeerSupport: msi, Registrar: ipRegistry}

	mockAclProvider = &aclmocks.MockACLProvider{}
	mockAclProvider.Reset()

	ledgerCleanup, err := peer.MockInitialize()
	if err != nil {
		return nil, nil, nil, nil, err
	}

	mspGetter := func(cid string) []string {
		return []string{"SampleOrg"}
	}

	peer.MockSetMSPIDGetter(mspGetter)

	grpcServer := grpc.NewServer()

	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to start peer listener %s", err)
	}
	_, localPort, err := net.SplitHostPort(lis.Addr().String())
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to get port: %s", err)
	}
	localIP, err := peer.GetLocalIP()
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to get local IP: %s", err)
	}

	peerAddress := net.JoinHostPort(localIP, localPort)

	tempdir, err := ioutil.TempDir("", "chaincode")
	if err != nil {
		panic(fmt.Sprintf("failed to create temporary directory: %s", err))
	}

	ccprovider.SetChaincodesPath(tempdir)
	ca, _ := tlsgen.NewCA()
	certGenerator := accesscontrol.NewAuthenticator(ca)
	config := &Config{
		TLSEnabled:     false,
		Keepalive:      time.Second,
		StartupTimeout: 3 * time.Minute,
		ExecuteTimeout: 30 * time.Second,
		LogLevel:       "info",
		ShimLogLevel:   "warning",
		LogFormat:      "TEST: [%{module}] %{shortfunc} -> %{level:.4s} %{id:03x}%{color:reset} %{message}",
	}
	pr := platforms.NewRegistry(&golang.Platform{})
	lsccImpl := lscc.New(sccp, mockAclProvider, pr)
	ml := &cm.Lifecycle{}
	ml.On("ChaincodeContainerInfo", ma.Anything, "lscc", ma.Anything).Return(
		&ccprovider.ChaincodeContainerInfo{
			Name:      "lscc",
			Version:   util.GetSysCCVersion(),
			PackageID: persistence.PackageID("lscc:" + util.GetSysCCVersion()),
		}, nil)
	ml.On("ChaincodeContainerInfo", ma.Anything, "pthru", ma.Anything).Return(
		&ccprovider.ChaincodeContainerInfo{
			Name:      "pthru",
			Version:   "0",
			PackageID: persistence.PackageID("pthru:0"),
		}, nil)
	ml.On("ChaincodeContainerInfo", ma.Anything, "example02", ma.Anything).Return(
		&ccprovider.ChaincodeContainerInfo{
			Name:      "example02",
			Version:   "0",
			PackageID: persistence.PackageID("example02:0"),
		}, nil)
	ml.On("ChaincodeContainerInfo", ma.Anything, "tmap", ma.Anything).Return(
		&ccprovider.ChaincodeContainerInfo{
			Name:      "tmap",
			Version:   "0",
			PackageID: persistence.PackageID("tmap:0"),
		}, nil)
	client, err := docker.NewClientFromEnv()
	if err != nil {
		return nil, nil, nil, nil, err
	}
	provider := &dockercontroller.Provider{
		PeerID:       "",
		NetworkID:    "",
		BuildMetrics: dockercontroller.NewBuildMetrics(&disabled.Provider{}),
		Client:       client,
	}
	chaincodeSupport := NewChaincodeSupport(
		config,
		peerAddress,
		false,
		ca.CertBytes(),
		certGenerator,
		&PackageProviderWrapper{FS: &ccprovider.CCInfoFSImpl{}},
		ml,
		aclmgmt.NewACLProvider(func(string) channelconfig.Resources { return nil }),
		container.NewVMController(
			map[string]container.VMProvider{
				dockercontroller.ContainerType: provider,
				inproccontroller.ContainerType: ipRegistry,
			},
		),
		sccp,
		pr,
		peer.DefaultSupport,
		&disabled.Provider{},
		&ledgermock.DeployedChaincodeInfoProvider{},
	)
	ipRegistry.ChaincodeSupport = chaincodeSupport
	pb.RegisterChaincodeSupportServer(grpcServer, chaincodeSupport)

	
	policy.RegisterPolicyCheckerFactory(&mockPolicyCheckerFactory{})

	ccp := &CCProviderImpl{cs: chaincodeSupport}
	sccp.RegisterSysCC(lsccImpl)

	for _, id := range chainIDs {
		sccp.DeDeploySysCCs(id, ccp)
		if err = peer.MockCreateChain(id); err != nil {
			closeListenerAndSleep(lis)
			return nil, nil, nil, nil, err
		}
		sccp.DeploySysCCs(id, ccp)
		
		if id != util.GetTestChainID() {
			mspmgmt.XXXSetMSPManager(id, mspmgmt.GetManagerForChain(util.GetTestChainID()))
		}
	}

	go grpcServer.Serve(lis)

	
	return ml, lis, chaincodeSupport, func() {
		finitPeer(lis, chainIDs...)
		os.RemoveAll(tempdir)
		ledgerCleanup()
	}, nil
}

func finitPeer(lis net.Listener, chainIDs ...string) {
	if lis != nil {
		for _, c := range chainIDs {
			if lgr := peer.GetLedger(c); lgr != nil {
				lgr.Close()
			}
		}
		closeListenerAndSleep(lis)
	}
	ledgerPath := config.GetPath("peer.fileSystemPath")
	os.RemoveAll(ledgerPath)
	os.RemoveAll(filepath.Join(os.TempDir(), "mcc-github"))
}

func startTxSimulation(chainID string, txid string) (ledger.TxSimulator, ledger.HistoryQueryExecutor, error) {
	lgr := peer.GetLedger(chainID)
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

func endTxSimulationCDS(chainID string, txid string, txsim ledger.TxSimulator, payload []byte, commit bool, cds *pb.ChaincodeDeploymentSpec, blockNumber uint64) error {
	
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

	return endTxSimulation(chainID, lsccid, txsim, payload, commit, prop, blockNumber)
}

func endTxSimulationCIS(chainID string, ccid *pb.ChaincodeID, txid string, txsim ledger.TxSimulator, payload []byte, commit bool, cis *pb.ChaincodeInvocationSpec, blockNumber uint64) error {
	
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

	return endTxSimulation(chainID, ccid, txsim, payload, commit, prop, blockNumber)
}









var _commitLock_ sync.Mutex

func endTxSimulation(chainID string, ccid *pb.ChaincodeID, txsim ledger.TxSimulator, _ []byte, commit bool, prop *pb.Proposal, blockNumber uint64) error {
	txsim.Done()
	if lgr := peer.GetLedger(chainID); lgr != nil {
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

			if err := lgr.CommitWithPvtData(blockAndPvtData); err != nil {
				return err
			}
		}
	}

	return nil
}


func getDeploymentSpec(spec *pb.ChaincodeSpec) (*pb.ChaincodeDeploymentSpec, error) {
	fmt.Printf("getting deployment spec for chaincode spec: %v\n", spec)
	codePackageBytes, err := container.GetChaincodePackageBytes(platforms.NewRegistry(&golang.Platform{}), spec)
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
	txsim, hqe, err := startTxSimulation(chainID, uuid)
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
			
			err = endTxSimulationCDS(chainID, uuid, txsim, []byte("deployed"), true, chaincodeDeploymentSpec, blockNumber)
		} else {
			
			endTxSimulationCDS(chainID, uuid, txsim, []byte("deployed"), false, chaincodeDeploymentSpec, blockNumber)
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

	txsim, hqe, err := startTxSimulation(chainID, uuid)
	if err != nil {
		return nil, uuid, nil, fmt.Errorf("Failed to get handle to simulator: %s ", err)
	}

	defer func() {
		
		if err == nil {
			
			err = endTxSimulationCIS(chainID, spec.ChaincodeId, uuid, txsim, []byte("invoke"), true, cdInvocationSpec, blockNumber)
		} else {
			
			endTxSimulationCIS(chainID, spec.ChaincodeId, uuid, txsim, []byte("invoke"), false, cdInvocationSpec, blockNumber)
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


func checkFinalState(chainID string, cccid *ccprovider.CCContext, a int, b int) error {
	txid := util.GenerateUUID()
	txsim, _, err := startTxSimulation(chainID, txid)
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
	fmt.Printf("Got string: %s\n", string(resbytes))
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

func runChaincodeInvokeChaincode(t *testing.T, channel1 string, channel2 string, tc tcicTc, cccid1 *ccprovider.CCContext, expectedA int, expectedB int, nextBlockNumber1, nextBlockNumber2 uint64, chaincodeSupport *ChaincodeSupport, ml *cm.Lifecycle) (uint64, uint64) {
	var ctxt = context.Background()

	
	chaincode2Name := generateChaincodeName(tc.chaincodeType)
	ml.On("ChaincodeContainerInfo", ma.Anything, chaincode2Name, ma.Anything).Return(
		&ccprovider.ChaincodeContainerInfo{
			Name:      chaincode2Name,
			Version:   "0",
			PackageID: persistence.PackageID(chaincode2Name + ":0"),
		}, nil)
	mcd := &cm.ChaincodeDefinition{}
	mcd.On("CCName").Return(chaincode2Name)
	mcd.On("CCVersion").Return("0")
	mcd.On("Hash").Return([]byte("Hulk, (sm)hash"))
	mcd.On("RequiresInit").Return(false)
	ml.On("ChaincodeDefinition", ma.Anything, chaincode2Name, ma.Anything).Return(mcd, nil)
	chaincode2Version := "0"
	chaincode2Type := tc.chaincodeType
	chaincode2Path := tc.chaincodePath
	chaincode2InitArgs := util.ToChaincodeArgs("init")
	chaincode2Creator := []byte([]byte("Alice"))

	
	_, cccid2, err := deployChaincode(
		ctxt,
		chaincode2Name,
		chaincode2Version,
		chaincode2Type,
		chaincode2Path,
		chaincode2InitArgs,
		chaincode2Creator,
		channel1,
		nextBlockNumber1,
		chaincodeSupport,
	)
	if err != nil {
		stopChaincode(ctxt, cccid1, chaincodeSupport)
		stopChaincode(ctxt, cccid2, chaincodeSupport)
		t.Fatalf("Error initializing chaincode %s(%+v)", chaincode2Name, err)
		return nextBlockNumber1, nextBlockNumber2
	}
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
	
	_, _, _, err = invoke(channel1, chaincode2InvokeSpec, nextBlockNumber1, []byte("Alice"), chaincodeSupport)
	if err != nil {
		stopChaincode(ctxt, cccid1, chaincodeSupport)
		stopChaincode(ctxt, cccid2, chaincodeSupport)
		t.Fatalf("Error invoking <%s>: %s", chaincode2Name, err)
		return nextBlockNumber1, nextBlockNumber2
	}
	nextBlockNumber1++

	
	err = checkFinalState(channel1, cccid1, expectedA, expectedB)
	if err != nil {
		stopChaincode(ctxt, cccid1, chaincodeSupport)
		stopChaincode(ctxt, cccid2, chaincodeSupport)
		t.Fatalf("Incorrect final state after transaction for <%s>: %s", cccid1.Name, err)
		return nextBlockNumber1, nextBlockNumber2
	}

	
	
	
	
	pm := peer.GetPolicyManager(channel1)
	pm.(*mockpolicies.Manager).PolicyMap = map[string]policies.Policy{
		policies.ChannelApplicationWriters: &CreatorPolicy{Creators: [][]byte{[]byte("Alice")}},
	}

	pm = peer.GetPolicyManager(channel2)
	pm.(*mockpolicies.Manager).PolicyMap = map[string]policies.Policy{
		policies.ChannelApplicationWriters: &CreatorPolicy{Creators: [][]byte{[]byte("Alice"), []byte("Bob")}},
	}

	
	_, cccid3, err := deployChaincode(
		ctxt,
		chaincode2Name,
		chaincode2Version,
		chaincode2Type,
		chaincode2Path,
		chaincode2InitArgs,
		chaincode2Creator,
		channel2,
		nextBlockNumber2,
		chaincodeSupport,
	)

	if err != nil {
		stopChaincode(ctxt, cccid1, chaincodeSupport)
		stopChaincode(ctxt, cccid2, chaincodeSupport)
		stopChaincode(ctxt, cccid3, chaincodeSupport)
		t.Fatalf("Error initializing chaincode %s/%s: %s", chaincode2Name, channel2, err)
		return nextBlockNumber1, nextBlockNumber2
	}
	nextBlockNumber2++
	time.Sleep(time.Second)

	chaincode2InvokeSpec = &pb.ChaincodeSpec{
		Type: chaincode2Type,
		ChaincodeId: &pb.ChaincodeID{
			Name:    chaincode2Name,
			Version: chaincode2Version,
		},
		Input: &pb.ChaincodeInput{
			Args: util.ToChaincodeArgs(cccid1.Name, "invoke", "a", "b", "10", channel1),
		},
	}

	
	_, _, _, err = invoke(channel2, chaincode2InvokeSpec, nextBlockNumber2, []byte("Bob"), chaincodeSupport)
	if err == nil {
		
		stopChaincode(ctxt, cccid1, chaincodeSupport)
		stopChaincode(ctxt, cccid2, chaincodeSupport)
		stopChaincode(ctxt, cccid3, chaincodeSupport)
		nextBlockNumber2++
		t.Fatalf("As Bob, invoking <%s/%s> via <%s/%s> should fail, but it succeeded.", cccid1.Name, channel1, chaincode2Name, channel2)
		return nextBlockNumber1, nextBlockNumber2
	}
	assert.True(t, strings.Contains(err.Error(), "[Creator not recognized [Bob]]"))

	
	_, _, _, err = invoke(channel2, chaincode2InvokeSpec, nextBlockNumber2, []byte("Alice"), chaincodeSupport)
	if err != nil {
		
		stopChaincode(ctxt, cccid1, chaincodeSupport)
		stopChaincode(ctxt, cccid2, chaincodeSupport)
		stopChaincode(ctxt, cccid3, chaincodeSupport)
		t.Fatalf("As Alice, invoking <%s/%s> via <%s/%s> should should of succeeded, but it failed: %s", cccid1.Name, channel1, chaincode2Name, channel2, err)
		return nextBlockNumber1, nextBlockNumber2
	}
	nextBlockNumber2++

	stopChaincode(ctxt, cccid1, chaincodeSupport)
	stopChaincode(ctxt, cccid2, chaincodeSupport)
	stopChaincode(ctxt, cccid3, chaincodeSupport)

	return nextBlockNumber1, nextBlockNumber2
}



type tcicTc struct {
	chaincodeType pb.ChaincodeSpec_Type
	chaincodePath string
}


func TestChaincodeInvokeChaincode(t *testing.T) {
	channel := util.GetTestChainID()
	channel2 := channel + "2"
	ml, lis, chaincodeSupport, cleanup, err := initPeer(channel, channel2)
	if err != nil {
		t.Fail()
		t.Logf("Error creating peer: %s", err)
	}
	defer cleanup()

	mockAclProvider.On("CheckACL", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	testCase := tcicTc{pb.ChaincodeSpec_GOLANG, chaincodePassthruGolangPath}

	ctx := context.Background()

	var nextBlockNumber1 uint64 = 1
	var nextBlockNumber2 uint64 = 1

	
	chaincode1Name := generateChaincodeName(pb.ChaincodeSpec_GOLANG)
	ml.On("ChaincodeContainerInfo", ma.Anything, chaincode1Name, ma.Anything).Return(
		&ccprovider.ChaincodeContainerInfo{
			Name:      chaincode1Name,
			Version:   "0",
			PackageID: persistence.PackageID(chaincode1Name + ":0"),
		}, nil)
	mcd := &cm.ChaincodeDefinition{}
	mcd.On("CCName").Return(chaincode1Name)
	mcd.On("CCVersion").Return("0")
	mcd.On("Hash").Return([]byte("Hulk, (sm)hash"))
	mcd.On("RequiresInit").Return(false)
	ml.On("ChaincodeDefinition", ma.Anything, chaincode1Name, ma.Anything).Return(mcd, nil)
	chaincode1Version := "0"
	chaincode1Type := pb.ChaincodeSpec_GOLANG
	chaincode1Path := chaincodeExample02GolangPath
	initialA := 100
	initialB := 200
	chaincode1InitArgs := util.ToChaincodeArgs("init", "a", strconv.Itoa(initialA), "b", strconv.Itoa(initialB))
	chaincode1Creator := []byte([]byte("Alice"))

	
	_, chaincodeCtx, err := deployChaincode(
		ctx,
		chaincode1Name,
		chaincode1Version,
		chaincode1Type,
		chaincode1Path,
		chaincode1InitArgs,
		chaincode1Creator,
		channel,
		nextBlockNumber1,
		chaincodeSupport,
	)
	if err != nil {
		stopChaincode(ctx, chaincodeCtx, chaincodeSupport)
		t.Fatalf("Error initializing chaincode %s: %s", chaincodeCtx.Name, err)
	}
	nextBlockNumber1++
	time.Sleep(time.Second)

	expectedA := initialA
	expectedB := initialB

	t.Run(testCase.chaincodeType.String(), func(t *testing.T) {
		expectedA = expectedA - 10
		expectedB = expectedB + 10
		nextBlockNumber1, nextBlockNumber2 = runChaincodeInvokeChaincode(
			t,
			channel,
			channel2,
			testCase,
			chaincodeCtx,
			expectedA,
			expectedB,
			nextBlockNumber1,
			nextBlockNumber2,
			chaincodeSupport,
			ml,
		)
	})

	closeListenerAndSleep(lis)
}

func stopChaincode(ctx context.Context, chaincodeCtx *ccprovider.CCContext, chaincodeSupport *ChaincodeSupport) {
	chaincodeSupport.Stop(&ccprovider.ChaincodeContainerInfo{
		Name:          chaincodeCtx.Name,
		Version:       chaincodeCtx.Version,
		ContainerType: "DOCKER",
		Type:          "GOLANG",
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

	mockAclProvider.On("CheckACL", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mcd := &cm.ChaincodeDefinition{}
	mcd.On("CCName").Return("example02")
	mcd.On("CCVersion").Return("0")
	mcd.On("Hash").Return([]byte("Hulk, (sm)hash"))
	mcd.On("RequiresInit").Return(false)
	ml.On("ChaincodeDefinition", ma.Anything, "example02", ma.Anything).Return(mcd, nil)

	
	cID1 := &pb.ChaincodeID{Name: "example02", Path: chaincodeExample02GolangPath, Version: "0"}
	f := "init"
	args := util.ToChaincodeArgs(f, "a", "100", "b", "200")

	spec1 := &pb.ChaincodeSpec{Type: 1, ChaincodeId: cID1, Input: &pb.ChaincodeInput{Args: args}}

	cccid1 := &ccprovider.CCContext{
		Name:    "example02",
		Version: "0",
	}

	var nextBlockNumber uint64 = 1
	defer chaincodeSupport.Stop(&ccprovider.ChaincodeContainerInfo{
		Name:          cID1.Name,
		Version:       cID1.Version,
		Path:          cID1.Path,
		ContainerType: "DOCKER",
		Type:          "GOLANG",
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

	defer chaincodeSupport.Stop(&ccprovider.ChaincodeContainerInfo{
		Name:          cID2.Name,
		Version:       cID2.Version,
		Path:          cID2.Path,
		ContainerType: "DOCKER",
		Type:          "GOLANG",
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

	defer chaincodeSupport.Stop(&ccprovider.ChaincodeContainerInfo{
		Name:          cID.Name,
		Version:       cID.Version,
		Path:          cID.Path,
		Type:          "GOLANG",
		ContainerType: "DOCKER",
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

	defer chaincodeSupport.Stop(&ccprovider.ChaincodeContainerInfo{
		Name:          cID.Name,
		Version:       cID.Version,
		Path:          cID.Path,
		Type:          "GOLANG",
		ContainerType: "DOCKER",
	})

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

	defer chaincodeSupport.Stop(&ccprovider.ChaincodeContainerInfo{
		Name:          cID.Name,
		Version:       cID.Version,
		Path:          cID.Path,
		Type:          "GOLANG",
		ContainerType: "DOCKER",
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

func deployChaincode(ctx context.Context, name string, version string, chaincodeType pb.ChaincodeSpec_Type, path string, args [][]byte, creator []byte, channel string, nextBlockNumber uint64, chaincodeSupport *ChaincodeSupport) (*pb.Response, *ccprovider.CCContext, error) {
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

var rng *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func generateChaincodeName(chaincodeType pb.ChaincodeSpec_Type) string {
	prefix := "cc_"
	switch chaincodeType {
	case pb.ChaincodeSpec_GOLANG:
		prefix = "cc_go_"
	case pb.ChaincodeSpec_JAVA:
		prefix = "cc_java_"
	case pb.ChaincodeSpec_NODE:
		prefix = "cc_js_"
	}
	return fmt.Sprintf("%s%06d", prefix, rng.Intn(999999))
}

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

type mockPolicyCheckerFactory struct{}

func (f *mockPolicyCheckerFactory) NewPolicyChecker() policy.PolicyChecker {
	return policy.NewPolicyChecker(
		peer.NewChannelPolicyManagerGetter(),
		&mocks.MockIdentityDeserializer{
			Identity: []byte("Admin"),
			Msg:      []byte("msg1"),
		},
		&mocks.MockMSPPrincipalGetter{Principal: []byte("Admin")},
	)
}
