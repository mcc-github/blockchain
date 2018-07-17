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
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/bccsp/factory"
	"github.com/mcc-github/blockchain/common/channelconfig"
	"github.com/mcc-github/blockchain/common/crypto/tlsgen"
	mc "github.com/mcc-github/blockchain/common/mocks/config"
	mockpolicies "github.com/mcc-github/blockchain/common/mocks/policies"
	"github.com/mcc-github/blockchain/common/policies"
	"github.com/mcc-github/blockchain/common/util"
	"github.com/mcc-github/blockchain/core/aclmgmt"
	aclmocks "github.com/mcc-github/blockchain/core/aclmgmt/mocks"
	"github.com/mcc-github/blockchain/core/chaincode/accesscontrol"
	"github.com/mcc-github/blockchain/core/chaincode/shim"
	"github.com/mcc-github/blockchain/core/common/ccprovider"
	"github.com/mcc-github/blockchain/core/config"
	"github.com/mcc-github/blockchain/core/container"
	"github.com/mcc-github/blockchain/core/container/dockercontroller"
	"github.com/mcc-github/blockchain/core/container/inproccontroller"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/ledger/ledgerconfig"
	"github.com/mcc-github/blockchain/core/ledger/ledgermgmt"
	cut "github.com/mcc-github/blockchain/core/ledger/util"
	"github.com/mcc-github/blockchain/core/ledger/util/couchdb"
	cmp "github.com/mcc-github/blockchain/core/mocks/peer"
	"github.com/mcc-github/blockchain/core/peer"
	"github.com/mcc-github/blockchain/core/policy"
	"github.com/mcc-github/blockchain/core/policy/mocks"
	"github.com/mcc-github/blockchain/core/scc"
	"github.com/mcc-github/blockchain/core/testutil"
	"github.com/mcc-github/blockchain/msp"
	mspmgmt "github.com/mcc-github/blockchain/msp/mgmt"
	"github.com/mcc-github/blockchain/msp/mgmt/testtools"
	"github.com/mcc-github/blockchain/protos/common"
	pb "github.com/mcc-github/blockchain/protos/peer"
	putils "github.com/mcc-github/blockchain/protos/utils"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var runTests bool

func testForSkip(t *testing.T) {
	
	if !runTests {
		t.SkipNow()
	}
}


func initPeer(chainIDs ...string) (net.Listener, *ChaincodeSupport, func(), error) {
	
	finitPeer(nil, chainIDs...)

	msi := &cmp.MockSupportImpl{
		GetApplicationConfigRv:     &mc.MockApplication{CapabilitiesRv: &mc.MockApplicationCapabilities{}},
		GetApplicationConfigBoolRv: true,
	}

	ipRegistry := inproccontroller.NewRegistry()
	sccp := &scc.Provider{Peer: peer.Default, PeerSupport: msi, Registrar: ipRegistry}

	mockAclProvider = &aclmocks.MockACLProvider{}
	mockAclProvider.Reset()

	peer.MockInitialize()

	mspGetter := func(cid string) []string {
		return []string{"SampleOrg"}
	}

	peer.MockSetMSPIDGetter(mspGetter)

	
	viper.Set("peer.tls.enabled", false)

	var opts []grpc.ServerOption
	if viper.GetBool("peer.tls.enabled") {
		creds, err := credentials.NewServerTLSFromFile(config.GetPath("peer.tls.cert.file"), config.GetPath("peer.tls.key.file"))
		if err != nil {
			return nil, nil, nil, fmt.Errorf("Failed to generate credentials %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}
	grpcServer := grpc.NewServer(opts...)

	peerAddress, err := peer.GetLocalAddress()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("Error obtaining peer address: %s", err)
	}
	lis, err := net.Listen("tcp", peerAddress)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("Error starting peer listener %s", err)
	}

	ccprovider.SetChaincodesPath(ccprovider.GetCCsPath())
	ca, _ := tlsgen.NewCA()
	certGenerator := accesscontrol.NewAuthenticator(ca)
	config := GlobalConfig()
	config.StartupTimeout = 3 * time.Minute
	chaincodeSupport := NewChaincodeSupport(
		config,
		peerAddress,
		false,
		ca.CertBytes(),
		certGenerator,
		&ccprovider.CCInfoFSImpl{},
		aclmgmt.NewACLProvider(func(string) channelconfig.Resources { return nil }),
		container.NewVMController(
			map[string]container.VMProvider{
				dockercontroller.ContainerType: dockercontroller.NewProvider("", ""),
				inproccontroller.ContainerType: ipRegistry,
			},
		),
		sccp,
	)
	pb.RegisterChaincodeSupportServer(grpcServer, chaincodeSupport)

	
	policy.RegisterPolicyCheckerFactory(&mockPolicyCheckerFactory{})

	ccp := &CCProviderImpl{cs: chaincodeSupport}
	for _, cc := range scc.CreateSysCCs(ccp, sccp, mockAclProvider) {
		sccp.RegisterSysCC(cc)
	}

	for _, id := range chainIDs {
		sccp.DeDeploySysCCs(id, ccp)
		if err = peer.MockCreateChain(id); err != nil {
			closeListenerAndSleep(lis)
			return nil, nil, nil, err
		}
		sccp.DeploySysCCs(id, ccp)
		
		if id != util.GetTestChainID() {
			mspmgmt.XXXSetMSPManager(id, mspmgmt.GetManagerForChain(util.GetTestChainID()))
		}
	}

	go grpcServer.Serve(lis)

	
	return lis, chaincodeSupport, func() { finitPeer(lis, chainIDs...) }, nil
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
	ledgermgmt.CleanupTestEnv()
	ledgerPath := config.GetPath("peer.fileSystemPath")
	os.RemoveAll(ledgerPath)
	os.RemoveAll(filepath.Join(os.TempDir(), "mcc-github"))

	
	if ledgerconfig.IsCouchDBEnabled() == true {

		chainID := util.GetTestChainID()

		connectURL := viper.GetString("ledger.state.couchDBConfig.couchDBAddress")
		username := viper.GetString("ledger.state.couchDBConfig.username")
		password := viper.GetString("ledger.state.couchDBConfig.password")
		maxRetries := viper.GetInt("ledger.state.couchDBConfig.maxRetries")
		maxRetriesOnStartup := viper.GetInt("ledger.state.couchDBConfig.maxRetriesOnStartup")
		requestTimeout := viper.GetDuration("ledger.state.couchDBConfig.requestTimeout")

		couchInstance, _ := couchdb.CreateCouchInstance(connectURL, username, password, maxRetries, maxRetriesOnStartup, requestTimeout)
		db := couchdb.CouchDatabase{CouchInstance: couchInstance, DBName: chainID}
		
		db.DropDatabase()

	}
}

func startTxSimulation(ctxt context.Context, chainID string, txid string) (context.Context, ledger.TxSimulator, error) {
	lgr := peer.GetLedger(chainID)
	txsim, err := lgr.NewTxSimulator(txid)
	if err != nil {
		return nil, nil, err
	}
	historyQueryExecutor, err := lgr.NewHistoryQueryExecutor()
	if err != nil {
		return nil, nil, err
	}

	ctxt = context.WithValue(ctxt, TXSimulatorKey, txsim)
	ctxt = context.WithValue(ctxt, HistoryQueryExecutorKey, historyQueryExecutor)
	return ctxt, txsim, nil
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

	
	prop, _, err := putils.CreateDeployProposalFromCDS(chainID, cds, ss, nil, nil, nil, nil)
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

	
	prop, returnedTxid, err := putils.CreateProposalFromCISAndTxid(txid, common.HeaderType_ENDORSER_TRANSACTION, chainID, cis, ss)
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
				return nil
			}
			
			resp, err := putils.CreateProposalResponse(prop.Header, prop.Payload, &pb.Response{Status: 200},
				txSimulationBytes, nil, ccid, nil, signer)
			if err != nil {
				return err
			}

			
			env, err := putils.CreateSignedTx(prop, signer, resp)
			if err != nil {
				return err
			}

			envBytes, err := putils.GetBytesEnvelope(env)
			if err != nil {
				return err
			}

			
			block := common.NewBlock(blockNumber, []byte{})
			block.Data.Data = [][]byte{envBytes}
			txsFilter := cut.NewTxValidationFlagsSetValue(len(block.Data.Data), pb.TxValidationCode_VALID)
			block.Metadata.Metadata[common.BlockMetadataIndex_TRANSACTIONS_FILTER] = txsFilter

			

			
			_commitLock_.Lock()
			defer _commitLock_.Unlock()

			blockAndPvtData := &ledger.BlockAndPvtData{
				Block:        block,
				BlockPvtData: make(map[uint64]*ledger.TxPvtData),
			}

			
			
			
			
			

			
			seqInBlock := uint64(0)

			if txSimulationResults.PvtSimulationResults != nil {

				blockAndPvtData.BlockPvtData[seqInBlock] = &ledger.TxPvtData{
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


func getDeploymentSpec(_ context.Context, spec *pb.ChaincodeSpec) (*pb.ChaincodeDeploymentSpec, error) {
	fmt.Printf("getting deployment spec for chaincode spec: %v\n", spec)
	codePackageBytes, err := container.GetChaincodePackageBytes(spec)
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


func deploy(ctx context.Context, cccid *ccprovider.CCContext, spec *pb.ChaincodeSpec, blockNumber uint64, chaincodeSupport *ChaincodeSupport) (b []byte, err error) {
	
	cdDeploymentSpec, err := getDeploymentSpec(ctx, spec)
	if err != nil {
		return nil, err
	}
	return deploy2(ctx, cccid, cdDeploymentSpec, nil, blockNumber, chaincodeSupport)
}

func deployWithCollectionConfigs(ctx context.Context, cccid *ccprovider.CCContext, spec *pb.ChaincodeSpec,
	collectionConfigPkg *common.CollectionConfigPackage, blockNumber uint64, chaincodeSupport *ChaincodeSupport) (b []byte, err error) {
	
	cdDeploymentSpec, err := getDeploymentSpec(ctx, spec)
	if err != nil {
		return nil, err
	}
	return deploy2(ctx, cccid, cdDeploymentSpec, collectionConfigPkg, blockNumber, chaincodeSupport)
}

func deploy2(ctx context.Context, cccid *ccprovider.CCContext, chaincodeDeploymentSpec *pb.ChaincodeDeploymentSpec,
	collectionConfigPkg *common.CollectionConfigPackage, blockNumber uint64, chaincodeSupport *ChaincodeSupport) (b []byte, err error) {
	cis, err := getDeployLSCCSpec(cccid.ChainID, chaincodeDeploymentSpec, collectionConfigPkg)
	if err != nil {
		return nil, fmt.Errorf("Error creating lscc spec : %s\n", err)
	}

	uuid := util.GenerateUUID()
	cccid.TxID = uuid
	ctx, txsim, err := startTxSimulation(ctx, cccid.ChainID, cccid.TxID)
	if err != nil {
		return nil, fmt.Errorf("Failed to get handle to simulator: %s ", err)
	}

	defer func() {
		
		if err == nil {
			
			err = endTxSimulationCDS(cccid.ChainID, uuid, txsim, []byte("deployed"), true, chaincodeDeploymentSpec, blockNumber)
		} else {
			
			endTxSimulationCDS(cccid.ChainID, uuid, txsim, []byte("deployed"), false, chaincodeDeploymentSpec, blockNumber)
		}
	}()

	
	ccprovider.PutChaincodeIntoFS(chaincodeDeploymentSpec)

	sysCCVers := util.GetSysCCVersion()
	sprop, prop := putils.MockSignedEndorserProposal2OrPanic(cccid.ChainID, cis.ChaincodeSpec, signer)
	lsccid := ccprovider.NewCCContext(cccid.ChainID, cis.ChaincodeSpec.ChaincodeId.Name, sysCCVers, uuid, true, sprop, prop)

	
	if _, _, err = chaincodeSupport.Execute(ctx, lsccid, cis); err != nil {
		return nil, fmt.Errorf("Error deploying chaincode (1): %s", err)
	}

	var resp *pb.Response
	if resp, _, err = chaincodeSupport.Execute(ctx, cccid, chaincodeDeploymentSpec); err != nil {
		return nil, fmt.Errorf("Error deploying chaincode(2): %s", err)
	}

	return resp.Payload, nil
}


func invoke(ctx context.Context, chainID string, spec *pb.ChaincodeSpec, blockNumber uint64, creator []byte, chaincodeSupport *ChaincodeSupport) (ccevt *pb.ChaincodeEvent, uuid string, retval []byte, err error) {
	return invokeWithVersion(ctx, chainID, spec.GetChaincodeId().Version, spec, blockNumber, creator, chaincodeSupport)
}


func invokeWithVersion(ctx context.Context, chainID string, version string, spec *pb.ChaincodeSpec, blockNumber uint64, creator []byte, chaincodeSupport *ChaincodeSupport) (ccevt *pb.ChaincodeEvent, uuid string, retval []byte, err error) {
	cdInvocationSpec := &pb.ChaincodeInvocationSpec{ChaincodeSpec: spec}

	
	uuid = util.GenerateUUID()

	var txsim ledger.TxSimulator
	ctx, txsim, err = startTxSimulation(ctx, chainID, uuid)
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
	sprop, prop := putils.MockSignedEndorserProposalOrPanic(chainID, spec, creator, []byte("msg1"))
	cccid := ccprovider.NewCCContext(chainID, cdInvocationSpec.ChaincodeSpec.ChaincodeId.Name, version, uuid, false, sprop, prop)
	var resp *pb.Response
	resp, ccevt, err = chaincodeSupport.Execute(ctx, cccid, cdInvocationSpec)
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

func executeDeployTransaction(t *testing.T, chainID string, name string, url string) {
	_, chaincodeSupport, cleanup, err := initPeer(chainID)
	if err != nil {
		t.Fail()
		t.Logf("Error creating peer: %s", err)
	}

	defer cleanup()

	var ctxt = context.Background()

	f := "init"
	args := util.ToChaincodeArgs(f, "a", "100", "b", "200")
	spec := &pb.ChaincodeSpec{Type: 1, ChaincodeId: &pb.ChaincodeID{Name: name, Path: url, Version: "0"}, Input: &pb.ChaincodeInput{Args: args}}

	cccid := ccprovider.NewCCContext(chainID, name, "0", "", false, nil, nil)

	_, err = deploy(ctxt, cccid, spec, 0, chaincodeSupport)

	cID := spec.ChaincodeId.Name
	if err != nil {
		chaincodeSupport.Stop(ctxt, cccid, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec})
		t.Fail()
		t.Logf("Error deploying <%s>: %s", cID, err)
		return
	}

	chaincodeSupport.Stop(ctxt, cccid, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec})
}


func checkFinalState(cccid *ccprovider.CCContext, a int, b int) error {
	txid := util.GenerateUUID()
	_, txsim, err := startTxSimulation(context.Background(), cccid.ChainID, txid)
	if err != nil {
		return fmt.Errorf("Failed to get handle to simulator: %s ", err)
	}

	defer txsim.Done()

	cName := cccid.GetCanonicalName()

	
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


func invokeExample02Transaction(ctxt context.Context, cccid *ccprovider.CCContext, cID *pb.ChaincodeID, chaincodeType pb.ChaincodeSpec_Type, args []string, chaincodeSupport *ChaincodeSupport) error {
	
	var nextBlockNumber uint64 = 1
	f := "init"
	argsDeploy := util.ToChaincodeArgs(f, "a", "100", "b", "200")
	spec := &pb.ChaincodeSpec{Type: chaincodeType, ChaincodeId: cID, Input: &pb.ChaincodeInput{Args: argsDeploy}}
	_, err := deploy(ctxt, cccid, spec, nextBlockNumber, chaincodeSupport)
	nextBlockNumber++
	ccID := spec.ChaincodeId.Name
	if err != nil {
		return fmt.Errorf("Error deploying <%s>: %s", ccID, err)
	}

	time.Sleep(time.Second)

	f = "invoke"
	invokeArgs := append([]string{f}, args...)
	spec = &pb.ChaincodeSpec{ChaincodeId: cID, Input: &pb.ChaincodeInput{Args: util.ToChaincodeArgs(invokeArgs...)}}
	_, uuid, _, err := invoke(ctxt, cccid.ChainID, spec, nextBlockNumber, nil, chaincodeSupport)
	nextBlockNumber++
	if err != nil {
		return fmt.Errorf("Error invoking <%s>: %s", cccid.Name, err)
	}

	cccid.TxID = uuid
	err = checkFinalState(cccid, 90, 210)
	if err != nil {
		return fmt.Errorf("Incorrect final state after transaction for <%s>: %s", ccID, err)
	}

	
	f = "delete"
	delArgs := util.ToChaincodeArgs(f, "a")
	spec = &pb.ChaincodeSpec{ChaincodeId: cID, Input: &pb.ChaincodeInput{Args: delArgs}}
	_, _, _, err = invoke(ctxt, cccid.ChainID, spec, nextBlockNumber, nil, chaincodeSupport)
	if err != nil {
		return fmt.Errorf("Error deleting state in <%s>: %s", cccid.Name, err)
	}

	return nil
}

const (
	chaincodeExample02GolangPath   = "github.com/mcc-github/blockchain/examples/chaincode/go/example02/cmd"
	chaincodeExample04GolangPath   = "github.com/mcc-github/blockchain/examples/chaincode/go/example04/cmd"
	chaincodeEventSenderGolangPath = "github.com/mcc-github/blockchain/examples/chaincode/go/eventsender"
	chaincodePassthruGolangPath    = "github.com/mcc-github/blockchain/examples/chaincode/go/passthru"
	chaincodeExample02JavaPath     = "../../examples/chaincode/java/chaincode_example02"
	chaincodeExample04JavaPath     = "../../examples/chaincode/java/chaincode_example04"
	chaincodeExample06JavaPath     = "../../examples/chaincode/java/chaincode_example06"
	chaincodeEventSenderJavaPath   = "../../examples/chaincode/java/eventsender"
)

func runChaincodeInvokeChaincode(t *testing.T, channel1 string, channel2 string, tc tcicTc, cccid1 *ccprovider.CCContext, expectedA int, expectedB int, nextBlockNumber1, nextBlockNumber2 uint64, chaincodeSupport *ChaincodeSupport) (uint64, uint64) {
	var ctxt = context.Background()

	
	chaincode2Name := generateChaincodeName(tc.chaincodeType)
	chaincode2Version := "0"
	chaincode2Type := tc.chaincodeType
	chaincode2Path := tc.chaincodePath
	chaincode2InitArgs := util.ToChaincodeArgs("init", "e", "0")
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
		t.Fatalf("Error initializing chaincode %s(%s)", chaincode2Name, err)
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
			Args: util.ToChaincodeArgs("invoke", cccid1.Name, "e", "1"),
		},
	}
	
	_, txID, _, err := invoke(ctxt, channel1, chaincode2InvokeSpec, nextBlockNumber1, []byte("Alice"), chaincodeSupport)
	if err != nil {
		stopChaincode(ctxt, cccid1, chaincodeSupport)
		stopChaincode(ctxt, cccid2, chaincodeSupport)
		t.Fatalf("Error invoking <%s>: %s", chaincode2Name, err)
		return nextBlockNumber1, nextBlockNumber2
	}
	nextBlockNumber1++

	
	cccid1.TxID = txID

	
	err = checkFinalState(cccid1, expectedA, expectedB)
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
			Args: util.ToChaincodeArgs("invoke", cccid1.Name, "e", "1", channel1),
		},
	}

	
	_, _, _, err = invoke(ctxt, channel2, chaincode2InvokeSpec, nextBlockNumber2, []byte("Bob"), chaincodeSupport)
	if err == nil {
		
		stopChaincode(ctxt, cccid1, chaincodeSupport)
		stopChaincode(ctxt, cccid2, chaincodeSupport)
		stopChaincode(ctxt, cccid3, chaincodeSupport)
		nextBlockNumber2++
		t.Fatalf("As Bob, invoking <%s/%s> via <%s/%s> should fail, but it succeeded.", cccid1.Name, cccid1.ChainID, chaincode2Name, channel2)
		return nextBlockNumber1, nextBlockNumber2
	}
	assert.True(t, strings.Contains(err.Error(), "[Creator not recognized [Bob]]"))

	
	_, _, _, err = invoke(ctxt, channel2, chaincode2InvokeSpec, nextBlockNumber2, []byte("Alice"), chaincodeSupport)
	if err != nil {
		
		stopChaincode(ctxt, cccid1, chaincodeSupport)
		stopChaincode(ctxt, cccid2, chaincodeSupport)
		stopChaincode(ctxt, cccid3, chaincodeSupport)
		t.Fatalf("As Alice, invoking <%s/%s> via <%s/%s> should should of succeeded, but it failed: %s", cccid1.Name, cccid1.ChainID, chaincode2Name, channel2, err)
		return nextBlockNumber1, nextBlockNumber2
	}
	nextBlockNumber2++

	stopChaincode(ctxt, cccid1, chaincodeSupport)
	stopChaincode(ctxt, cccid2, chaincodeSupport)
	stopChaincode(ctxt, cccid3, chaincodeSupport)

	return nextBlockNumber1, nextBlockNumber2
}


func TestExecuteDeployTransaction(t *testing.T) {
	
	t.Skip()
	chainID := util.GetTestChainID()

	executeDeployTransaction(t, chainID, "example01", "github.com/mcc-github/blockchain/examples/chaincode/go/chaincode_example01")
}


func TestGopathExecuteDeployTransaction(t *testing.T) {
	
	t.Skip()
	chainID := util.GetTestChainID()

	
	
	os.Setenv("GOPATH", os.Getenv("GOPATH")+string(os.PathSeparator)+string(os.PathListSeparator)+"/tmp/foo"+string(os.PathListSeparator)+"/tmp/bar")
	executeDeployTransaction(t, chainID, "example01", "github.com/mcc-github/blockchain/examples/chaincode/go/chaincode_example01")
}

func TestExecuteInvokeTransaction(t *testing.T) {
	testForSkip(t)

	testCases := []struct {
		chaincodeType pb.ChaincodeSpec_Type
		chaincodePath string
	}{
		{pb.ChaincodeSpec_GOLANG, chaincodeExample02GolangPath},
		{pb.ChaincodeSpec_JAVA, chaincodeExample02JavaPath},
	}

	for _, tc := range testCases {
		t.Run(tc.chaincodeType.String(), func(t *testing.T) {

			if tc.chaincodeType == pb.ChaincodeSpec_JAVA && runtime.GOARCH != "amd64" {
				t.Skip("No Java chaincode support yet on non-amd64.")
			}

			chainID := util.GetTestChainID()

			_, chaincodeSupport, cleanup, err := initPeer(chainID)
			if err != nil {
				t.Fail()
				t.Logf("Error creating peer: %s", err)
			}

			var ctxt = context.Background()
			chaincodeName := generateChaincodeName(tc.chaincodeType)
			chaincodeVersion := "1.0.0.0"
			cccid := ccprovider.NewCCContext(chainID, chaincodeName, chaincodeVersion, "", false, nil, nil)
			ccID := &pb.ChaincodeID{Name: chaincodeName, Path: tc.chaincodePath, Version: chaincodeVersion}

			args := []string{"a", "b", "10"}
			err = invokeExample02Transaction(ctxt, cccid, ccID, tc.chaincodeType, args, chaincodeSupport)
			if err != nil {
				t.Fail()
				t.Logf("Error invoking transaction: %s", err)
			} else {
				fmt.Print("Invoke test passed\n")
				t.Log("Invoke test passed")
			}

			chaincodeSupport.Stop(ctxt, cccid, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: &pb.ChaincodeSpec{ChaincodeId: ccID}})
			cleanup()
		})
	}

}


func TestExecuteInvokeInvalidTransaction(t *testing.T) {
	testForSkip(t)

	chainID := util.GetTestChainID()

	_, chaincodeSupport, cleanup, err := initPeer(chainID)
	if err != nil {
		t.Fail()
		t.Logf("Error creating peer: %s", err)
	}

	defer cleanup()

	var ctxt = context.Background()

	url := "github.com/mcc-github/blockchain/examples/chaincode/go/chaincode_example02"
	ccID := &pb.ChaincodeID{Name: "example02", Path: url, Version: "0"}

	cccid := ccprovider.NewCCContext(chainID, "example02", "0", "", false, nil, nil)

	
	args := []string{"x", "-1"}
	err = invokeExample02Transaction(ctxt, cccid, ccID, pb.ChaincodeSpec_GOLANG, args, chaincodeSupport)

	
	if err != nil {
		errStr := err.Error()
		t.Logf("Got error %s\n", errStr)
		t.Log("InvalidInvoke test passed")
		chaincodeSupport.Stop(ctxt, cccid, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: &pb.ChaincodeSpec{ChaincodeId: ccID}})

		return
	}

	t.Fail()
	t.Logf("Error invoking transaction %s", err)

	chaincodeSupport.Stop(ctxt, cccid, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: &pb.ChaincodeSpec{ChaincodeId: ccID}})
}


type tcicTc struct {
	chaincodeType pb.ChaincodeSpec_Type
	chaincodePath string
}


func TestChaincodeInvokeChaincode(t *testing.T) {
	channel := util.GetTestChainID()
	channel2 := channel + "2"
	lis, chaincodeSupport, cleanup, err := initPeer(channel, channel2)
	if err != nil {
		t.Fail()
		t.Logf("Error creating peer: %s", err)
	}
	defer cleanup()

	mockAclProvider.On("CheckACL", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	testCases := []tcicTc{
		{pb.ChaincodeSpec_GOLANG, chaincodeExample04GolangPath},
	}

	ctx := context.Background()

	var nextBlockNumber1 uint64 = 1
	var nextBlockNumber2 uint64 = 1

	
	chaincode1Name := generateChaincodeName(pb.ChaincodeSpec_GOLANG)
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

	for _, tc := range testCases {
		t.Run(tc.chaincodeType.String(), func(t *testing.T) {
			expectedA = expectedA - 10
			expectedB = expectedB + 10
			nextBlockNumber1, nextBlockNumber2 = runChaincodeInvokeChaincode(
				t,
				channel,
				channel2,
				tc,
				chaincodeCtx,
				expectedA,
				expectedB,
				nextBlockNumber1,
				nextBlockNumber2,
				chaincodeSupport,
			)
		})
	}

	closeListenerAndSleep(lis)
}

func stopChaincode(ctx context.Context, chaincodeCtx *ccprovider.CCContext, chaincodeSupport *ChaincodeSupport) {
	chaincodeSupport.Stop(ctx, chaincodeCtx,
		&pb.ChaincodeDeploymentSpec{
			ChaincodeSpec: &pb.ChaincodeSpec{
				ChaincodeId: &pb.ChaincodeID{
					Name:    chaincodeCtx.Name,
					Version: chaincodeCtx.Version,
				},
			},
		})
}



func TestChaincodeInvokeChaincodeErrorCase(t *testing.T) {
	ctxt := context.Background()
	chainID := util.GetTestChainID()

	_, chaincodeSupport, cleanup, err := initPeer(chainID)
	if err != nil {
		t.Fail()
		t.Logf("Error creating peer: %s", err)
	}
	defer cleanup()

	mockAclProvider.On("CheckACL", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	
	cID1 := &pb.ChaincodeID{Name: "example02", Path: chaincodeExample02GolangPath, Version: "0"}
	f := "init"
	args := util.ToChaincodeArgs(f, "a", "100", "b", "200")

	spec1 := &pb.ChaincodeSpec{Type: 1, ChaincodeId: cID1, Input: &pb.ChaincodeInput{Args: args}}

	sProp, prop := putils.MockSignedEndorserProposalOrPanic(util.GetTestChainID(), spec1, []byte([]byte("Alice")), nil)
	cccid1 := ccprovider.NewCCContext(chainID, "example02", "0", "", false, sProp, prop)

	var nextBlockNumber uint64 = 1

	_, err = deploy(ctxt, cccid1, spec1, nextBlockNumber, chaincodeSupport)
	nextBlockNumber++
	ccID1 := spec1.ChaincodeId.Name
	if err != nil {
		t.Fail()
		t.Logf("Error initializing chaincode %s(%s)", ccID1, err)
		chaincodeSupport.Stop(ctxt, cccid1, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec1})
		return
	}

	time.Sleep(time.Second)

	
	cID2 := &pb.ChaincodeID{Name: "pthru", Path: chaincodePassthruGolangPath, Version: "0"}
	f = "init"
	args = util.ToChaincodeArgs(f)

	spec2 := &pb.ChaincodeSpec{Type: 1, ChaincodeId: cID2, Input: &pb.ChaincodeInput{Args: args}}

	cccid2 := ccprovider.NewCCContext(chainID, "pthru", "0", "", false, sProp, prop)

	_, err = deploy(ctxt, cccid2, spec2, nextBlockNumber, chaincodeSupport)
	nextBlockNumber++
	ccID2 := spec2.ChaincodeId.Name
	if err != nil {
		t.Fail()
		t.Logf("Error initializing chaincode %s(%s)", ccID2, err)
		chaincodeSupport.Stop(ctxt, cccid1, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec1})
		chaincodeSupport.Stop(ctxt, cccid2, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec2})
		return
	}

	time.Sleep(time.Second)

	
	f = ccID1
	args = util.ToChaincodeArgs(f, "invoke", "a")

	spec2 = &pb.ChaincodeSpec{Type: 1, ChaincodeId: cID2, Input: &pb.ChaincodeInput{Args: args}}
	
	_, _, _, err = invoke(ctxt, chainID, spec2, nextBlockNumber, []byte("Alice"), chaincodeSupport)

	if err == nil {
		t.Fail()
		t.Logf("Error invoking <%s>: %s", ccID2, err)
		chaincodeSupport.Stop(ctxt, cccid1, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec1})
		chaincodeSupport.Stop(ctxt, cccid2, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec2})
		return
	}

	if strings.Index(err.Error(), "Error invoking chaincode: Incorrect number of arguments. Expecting 3") < 0 {
		t.Fail()
		t.Logf("Unexpected error %s", err)
		chaincodeSupport.Stop(ctxt, cccid1, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec1})
		chaincodeSupport.Stop(ctxt, cccid2, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec2})
		return
	}

	chaincodeSupport.Stop(ctxt, cccid1, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec1})
	chaincodeSupport.Stop(ctxt, cccid2, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec2})
}


func TestQueries(t *testing.T) {
	
	

	chainID := util.GetTestChainID()

	_, chaincodeSupport, cleanup, err := initPeer(chainID)
	if err != nil {
		t.Fail()
		t.Logf("Error creating peer: %s", err)
	}

	defer cleanup()

	var ctxt = context.Background()

	url := "github.com/mcc-github/blockchain/examples/chaincode/go/map"
	cID := &pb.ChaincodeID{Name: "tmap", Path: url, Version: "0"}

	f := "init"
	args := util.ToChaincodeArgs(f)

	spec := &pb.ChaincodeSpec{Type: 1, ChaincodeId: cID, Input: &pb.ChaincodeInput{Args: args}}

	cccid := ccprovider.NewCCContext(chainID, "tmap", "0", "", false, nil, nil)

	var nextBlockNumber uint64 = 1
	_, err = deploy(ctxt, cccid, spec, nextBlockNumber, chaincodeSupport)
	nextBlockNumber++
	ccID := spec.ChaincodeId.Name
	if err != nil {
		t.Fail()
		t.Logf("Error initializing chaincode %s(%s)", ccID, err)
		chaincodeSupport.Stop(ctxt, cccid, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec})
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
		_, _, _, err = invoke(ctxt, chainID, spec, nextBlockNumber, nil, chaincodeSupport)
		nextBlockNumber++

		if err != nil {
			t.Fail()
			t.Logf("Error invoking <%s>: %s", ccID, err)
			chaincodeSupport.Stop(ctxt, cccid, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec})
			return
		}

	}

	
	f = "keys"
	args = util.ToChaincodeArgs(f, "marble001", "marble011")

	spec = &pb.ChaincodeSpec{Type: 1, ChaincodeId: cID, Input: &pb.ChaincodeInput{Args: args}}
	_, _, retval, err := invoke(ctxt, chainID, spec, nextBlockNumber, nil, chaincodeSupport)
	nextBlockNumber++
	if err != nil {
		t.Fail()
		t.Logf("Error invoking <%s>: %s", ccID, err)
		chaincodeSupport.Stop(ctxt, cccid, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec})
		return
	}

	err = json.Unmarshal(retval, &keys)
	if len(keys) != 10 {
		t.Fail()
		t.Logf("Error detected with the range query, should have returned 10 but returned %v", len(keys))
		chaincodeSupport.Stop(ctxt, cccid, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec})
		return
	}

	
	

	
	origTimeout := chaincodeSupport.ExecuteTimeout
	chaincodeSupport.ExecuteTimeout = time.Duration(1) * time.Second

	
	args = util.ToChaincodeArgs(f, "marble001", "marble002", "2000")

	spec = &pb.ChaincodeSpec{Type: 1, ChaincodeId: cID, Input: &pb.ChaincodeInput{Args: args}}
	_, _, retval, err = invoke(ctxt, chainID, spec, nextBlockNumber, nil, chaincodeSupport)
	if err == nil {
		t.Fail()
		t.Logf("expected timeout error but succeeded")
		chaincodeSupport.Stop(ctxt, cccid, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec})
		return
	}

	
	chaincodeSupport.ExecuteTimeout = origTimeout

	
	
	
	f = "keys"
	args = util.ToChaincodeArgs(f, "marble001", "marble102")

	spec = &pb.ChaincodeSpec{Type: 1, ChaincodeId: cID, Input: &pb.ChaincodeInput{Args: args}}
	_, _, retval, err = invoke(ctxt, chainID, spec, nextBlockNumber, nil, chaincodeSupport)
	nextBlockNumber++
	if err != nil {
		t.Fail()
		t.Logf("Error invoking <%s>: %s", ccID, err)
		chaincodeSupport.Stop(ctxt, cccid, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec})
		return
	}

	
	err = json.Unmarshal(retval, &keys)

	
	
	if len(keys) != 101 {
		t.Fail()
		t.Logf("Error detected with the range query, should have returned 101 but returned %v", len(keys))
		chaincodeSupport.Stop(ctxt, cccid, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec})
		return
	}

	
	
	
	f = "keys"
	args = util.ToChaincodeArgs(f, "", "")

	spec = &pb.ChaincodeSpec{Type: 1, ChaincodeId: cID, Input: &pb.ChaincodeInput{Args: args}}
	_, _, retval, err = invoke(ctxt, chainID, spec, nextBlockNumber, nil, chaincodeSupport)
	nextBlockNumber++
	if err != nil {
		t.Fail()
		t.Logf("Error invoking <%s>: %s", ccID, err)
		chaincodeSupport.Stop(ctxt, cccid, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec})
		return
	}

	
	err = json.Unmarshal(retval, &keys)

	
	
	if len(keys) != 101 {
		t.Fail()
		t.Logf("Error detected with the range query, should have returned 101 but returned %v", len(keys))
		chaincodeSupport.Stop(ctxt, cccid, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec})
		return
	}

	
	
	if ledgerconfig.IsCouchDBEnabled() == true {

		
		
		f = "query"
		args = util.ToChaincodeArgs(f, "{\"selector\":{\"color\":\"blue\"}}")

		spec = &pb.ChaincodeSpec{Type: 1, ChaincodeId: cID, Input: &pb.ChaincodeInput{Args: args}}
		_, _, _, err = invoke(ctxt, chainID, spec, nextBlockNumber, nil, chaincodeSupport)
		nextBlockNumber++

		if err != nil {
			t.Fail()
			t.Logf("Error invoking <%s>: %s", ccID, err)
			chaincodeSupport.Stop(ctxt, cccid, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec})
			return
		}

		
		err = json.Unmarshal(retval, &keys)

		
		if len(keys) != 100 {
			t.Fail()
			t.Logf("Error detected with the rich query, should have returned 100 but returned %v %s", len(keys), keys)
			chaincodeSupport.Stop(ctxt, cccid, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec})
			return
		}
		
		viper.Set("ledger.state.queryLimit", 5)

		
		f = "keys"
		args = util.ToChaincodeArgs(f, "marble001", "marble011")

		spec = &pb.ChaincodeSpec{Type: 1, ChaincodeId: cID, Input: &pb.ChaincodeInput{Args: args}}
		_, _, retval, err := invoke(ctxt, chainID, spec, nextBlockNumber, nil, chaincodeSupport)
		nextBlockNumber++
		if err != nil {
			t.Fail()
			t.Logf("Error invoking <%s>: %s", ccID, err)
			chaincodeSupport.Stop(ctxt, cccid, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec})
			return
		}

		
		err = json.Unmarshal(retval, &keys)
		
		if len(keys) != 5 {
			t.Fail()
			t.Logf("Error detected with the range query, should have returned 5 but returned %v", len(keys))
			chaincodeSupport.Stop(ctxt, cccid, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec})
			return
		}

		
		viper.Set("ledger.state.queryLimit", 10000)

		
		f = "query"
		args = util.ToChaincodeArgs(f, "{\"selector\":{\"owner\":\"jerry\"}}")

		spec = &pb.ChaincodeSpec{Type: 1, ChaincodeId: cID, Input: &pb.ChaincodeInput{Args: args}}
		_, _, retval, err = invoke(ctxt, chainID, spec, nextBlockNumber, nil, chaincodeSupport)
		nextBlockNumber++

		if err != nil {
			t.Fail()
			t.Logf("Error invoking <%s>: %s", ccID, err)
			chaincodeSupport.Stop(ctxt, cccid, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec})
			return
		}

		
		err = json.Unmarshal(retval, &keys)

		
		
		if len(keys) != 50 {
			t.Fail()
			t.Logf("Error detected with the rich query, should have returned 50 but returned %v", len(keys))
			chaincodeSupport.Stop(ctxt, cccid, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec})
			return
		}

		
		viper.Set("ledger.state.queryLimit", 5)

		
		f = "query"
		args = util.ToChaincodeArgs(f, "{\"selector\":{\"owner\":\"jerry\"}}")

		spec = &pb.ChaincodeSpec{Type: 1, ChaincodeId: cID, Input: &pb.ChaincodeInput{Args: args}}
		_, _, retval, err = invoke(ctxt, chainID, spec, nextBlockNumber, nil, chaincodeSupport)
		nextBlockNumber++
		if err != nil {
			t.Fail()
			t.Logf("Error invoking <%s>: %s", ccID, err)
			chaincodeSupport.Stop(ctxt, cccid, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec})
			return
		}

		
		err = json.Unmarshal(retval, &keys)

		
		if len(keys) != 5 {
			t.Fail()
			t.Logf("Error detected with the rich query, should have returned 5 but returned %v", len(keys))
			chaincodeSupport.Stop(ctxt, cccid, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec})
			return
		}

	}

	
	f = "put"
	args = util.ToChaincodeArgs(f, "marble012", "{\"docType\":\"marble\",\"name\":\"marble012\",\"color\":\"red\",\"size\":30,\"owner\":\"jerry\"}")
	spec = &pb.ChaincodeSpec{Type: 1, ChaincodeId: cID, Input: &pb.ChaincodeInput{Args: args}}
	_, _, _, err = invoke(ctxt, chainID, spec, nextBlockNumber, nil, chaincodeSupport)
	nextBlockNumber++
	if err != nil {
		t.Fail()
		t.Logf("Error invoking <%s>: %s", ccID, err)
		chaincodeSupport.Stop(ctxt, cccid, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec})
		return
	}

	f = "put"
	args = util.ToChaincodeArgs(f, "marble012", "{\"docType\":\"marble\",\"name\":\"marble012\",\"color\":\"red\",\"size\":30,\"owner\":\"jerry\"}")
	spec = &pb.ChaincodeSpec{Type: 1, ChaincodeId: cID, Input: &pb.ChaincodeInput{Args: args}}
	_, _, _, err = invoke(ctxt, chainID, spec, nextBlockNumber, nil, chaincodeSupport)
	nextBlockNumber++
	if err != nil {
		t.Fail()
		t.Logf("Error invoking <%s>: %s", ccID, err)
		chaincodeSupport.Stop(ctxt, cccid, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec})
		return
	}

	
	f = "history"
	args = util.ToChaincodeArgs(f, "marble012")
	spec = &pb.ChaincodeSpec{Type: 1, ChaincodeId: cID, Input: &pb.ChaincodeInput{Args: args}}
	_, _, retval, err = invoke(ctxt, chainID, spec, nextBlockNumber, nil, chaincodeSupport)
	nextBlockNumber++
	if err != nil {
		t.Fail()
		t.Logf("Error invoking <%s>: %s", ccID, err)
		chaincodeSupport.Stop(ctxt, cccid, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec})
		return
	}

	var history []interface{}
	err = json.Unmarshal(retval, &history)
	if len(history) != 3 {
		t.Fail()
		t.Logf("Error detected with the history query, should have returned 3 but returned %v", len(history))
		chaincodeSupport.Stop(ctxt, cccid, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec})
		return
	}

	chaincodeSupport.Stop(ctxt, cccid, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec})
}

func TestGetEvent(t *testing.T) {
	testForSkip(t)
	testCases := []struct {
		chaincodeType pb.ChaincodeSpec_Type
		chaincodePath string
	}{
		{pb.ChaincodeSpec_GOLANG, chaincodeEventSenderGolangPath},
		{pb.ChaincodeSpec_JAVA, chaincodeEventSenderJavaPath},
	}

	chainID := util.GetTestChainID()
	var nextBlockNumber uint64

	_, chaincodeSupport, cleanup, err := initPeer(chainID)
	if err != nil {
		t.Fail()
		t.Logf("Error creating peer: %s", err)
	}

	nextBlockNumber++

	defer cleanup()

	for _, tc := range testCases {
		t.Run(tc.chaincodeType.String(), func(t *testing.T) {

			if tc.chaincodeType == pb.ChaincodeSpec_JAVA && runtime.GOARCH != "amd64" {
				t.Skip("No Java chaincode support yet on non-amd64.")
			}

			var ctxt = context.Background()

			cID := &pb.ChaincodeID{Name: generateChaincodeName(tc.chaincodeType), Path: tc.chaincodePath, Version: "0"}
			f := "init"
			spec := &pb.ChaincodeSpec{Type: tc.chaincodeType, ChaincodeId: cID, Input: &pb.ChaincodeInput{Args: util.ToChaincodeArgs(f)}}

			cccid := ccprovider.NewCCContext(chainID, cID.Name, cID.Version, "", false, nil, nil)
			_, err = deploy(ctxt, cccid, spec, nextBlockNumber, chaincodeSupport)
			nextBlockNumber++
			ccID := spec.ChaincodeId.Name
			if err != nil {
				t.Fail()
				t.Logf("Error initializing chaincode %s(%s)", ccID, err)
				chaincodeSupport.Stop(ctxt, cccid, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec})
				return
			}

			time.Sleep(time.Second)

			args := util.ToChaincodeArgs("invoke", "i", "am", "satoshi")

			spec = &pb.ChaincodeSpec{Type: 1, ChaincodeId: cID, Input: &pb.ChaincodeInput{Args: args}}

			var ccevt *pb.ChaincodeEvent
			ccevt, _, _, err = invoke(ctxt, chainID, spec, nextBlockNumber, nil, chaincodeSupport)
			nextBlockNumber++

			if err != nil {
				t.Logf("Error invoking chaincode %s(%s)", ccID, err)
				t.Fail()
			}

			if ccevt == nil {
				t.Logf("Error ccevt is nil %s(%s)", ccID, err)
				t.Fail()
			}

			if ccevt.ChaincodeId != ccID {
				t.Logf("Error ccevt id(%s) != cid(%s)", ccevt.ChaincodeId, ccID)
				t.Fail()
			}

			if strings.Index(string(ccevt.Payload), "i,am,satoshi") < 0 {
				t.Logf("Error expected event not found (%s)", string(ccevt.Payload))
				t.Fail()
			}

			chaincodeSupport.Stop(ctxt, cccid, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec})
		})
	}

}



func TestChaincodeQueryChaincodeUsingInvoke(t *testing.T) {
	testForSkip(t)
	
	
	t.Skip()
	chainID := util.GetTestChainID()

	_, chaincodeSupport, cleanup, err := initPeer(chainID)
	if err != nil {
		t.Fail()
		t.Logf("Error registering user  %s", err)
		return
	}

	defer cleanup()

	var ctxt = context.Background()

	
	url1 := "github.com/mcc-github/blockchain/examples/chaincode/go/chaincode_example02"

	cID1 := &pb.ChaincodeID{Name: "example02", Path: url1, Version: "0"}
	f := "init"
	args := util.ToChaincodeArgs(f, "a", "100", "b", "200")

	spec1 := &pb.ChaincodeSpec{Type: 1, ChaincodeId: cID1, Input: &pb.ChaincodeInput{Args: args}}

	sProp, prop := putils.MockSignedEndorserProposalOrPanic(util.GetTestChainID(), spec1, []byte([]byte("Alice")), nil)
	cccid1 := ccprovider.NewCCContext(chainID, "example02", "0", "", false, sProp, prop)
	var nextBlockNumber uint64
	_, err = deploy(ctxt, cccid1, spec1, nextBlockNumber, chaincodeSupport)
	nextBlockNumber++
	ccID1 := spec1.ChaincodeId.Name
	if err != nil {
		chaincodeSupport.Stop(ctxt, cccid1, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec1})
		t.Fail()
		t.Logf("Error initializing chaincode %s(%s)", ccID1, err)
		return
	}

	time.Sleep(time.Second)

	
	url2 := "github.com/mcc-github/blockchain/examples/chaincode/go/chaincode_example05"

	cID2 := &pb.ChaincodeID{Name: "example05", Path: url2, Version: "0"}
	f = "init"
	args = util.ToChaincodeArgs(f, "sum", "0")

	spec2 := &pb.ChaincodeSpec{Type: 1, ChaincodeId: cID2, Input: &pb.ChaincodeInput{Args: args}}

	cccid2 := ccprovider.NewCCContext(chainID, "example05", "0", "", false, sProp, prop)

	_, err = deploy(ctxt, cccid2, spec2, nextBlockNumber, chaincodeSupport)
	nextBlockNumber++
	ccID2 := spec2.ChaincodeId.Name
	if err != nil {
		chaincodeSupport.Stop(ctxt, cccid1, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec1})
		chaincodeSupport.Stop(ctxt, cccid2, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec2})
		t.Fail()
		t.Logf("Error initializing chaincode %s(%s)", ccID2, err)
		return
	}

	time.Sleep(time.Second)

	
	f = "invoke"
	args = util.ToChaincodeArgs(f, ccID1, "sum")

	spec2 = &pb.ChaincodeSpec{Type: 1, ChaincodeId: cID2, Input: &pb.ChaincodeInput{Args: args}}
	
	var retVal []byte
	_, _, retVal, err = invoke(ctxt, chainID, spec2, nextBlockNumber, []byte("Alice"), chaincodeSupport)
	nextBlockNumber++
	if err != nil {
		chaincodeSupport.Stop(ctxt, cccid1, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec1})
		chaincodeSupport.Stop(ctxt, cccid2, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec2})
		t.Fail()
		t.Logf("Error invoking <%s>: %s", ccID2, err)
		return
	}

	
	result, err := strconv.Atoi(string(retVal))
	if err != nil || result != 300 {
		chaincodeSupport.Stop(ctxt, cccid1, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec1})
		chaincodeSupport.Stop(ctxt, cccid2, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec2})
		t.Fail()
		t.Logf("Incorrect final state after transaction for <%s>: %s", ccID1, err)
		return
	}

	
	f = "query"
	args = util.ToChaincodeArgs(f, ccID1, "sum")

	spec2 = &pb.ChaincodeSpec{Type: 1, ChaincodeId: cID2, Input: &pb.ChaincodeInput{Args: args}}
	
	_, _, retVal, err = invoke(ctxt, chainID, spec2, nextBlockNumber, []byte("Alice"), chaincodeSupport)

	if err != nil {
		chaincodeSupport.Stop(ctxt, cccid1, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec1})
		chaincodeSupport.Stop(ctxt, cccid2, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec2})
		t.Fail()
		t.Logf("Error querying <%s>: %s", ccID2, err)
		return
	}

	
	result, err = strconv.Atoi(string(retVal))
	if err != nil || result != 300 {
		chaincodeSupport.Stop(ctxt, cccid1, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec1})
		chaincodeSupport.Stop(ctxt, cccid2, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec2})
		t.Fail()
		t.Logf("Incorrect final value after query for <%s>: %s", ccID1, err)
		return
	}

	chaincodeSupport.Stop(ctxt, cccid1, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec1})
	chaincodeSupport.Stop(ctxt, cccid2, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec2})
}


func TestChaincodeInvokesForbiddenSystemChaincode(t *testing.T) {
	testForSkip(t)
	chainID := util.GetTestChainID()

	_, chaincodeSupport, cleanup, err := initPeer(chainID)
	if err != nil {
		t.Fail()
		t.Logf("Error creating peer: %s", err)
	}

	defer cleanup()

	var ctxt = context.Background()

	var nextBlockNumber uint64 = 1

	
	url := "github.com/mcc-github/blockchain/examples/chaincode/go/passthru"

	cID := &pb.ChaincodeID{Name: "pthru", Path: url, Version: "0"}
	f := "init"
	args := util.ToChaincodeArgs(f)

	spec := &pb.ChaincodeSpec{Type: 1, ChaincodeId: cID, Input: &pb.ChaincodeInput{Args: args}}

	cccid := ccprovider.NewCCContext(chainID, "pthru", "0", "", false, nil, nil)

	_, err = deploy(ctxt, cccid, spec, nextBlockNumber, chaincodeSupport)
	nextBlockNumber++
	ccID := spec.ChaincodeId.Name
	if err != nil {
		t.Fail()
		t.Logf("Error initializing chaincode %s(%s)", ccID, err)
		chaincodeSupport.Stop(ctxt, cccid, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec})
		return
	}

	time.Sleep(time.Second)

	
	
	args = util.ToChaincodeArgs("escc/"+chainID, "getid", chainID, "pthru")

	spec = &pb.ChaincodeSpec{Type: 1, ChaincodeId: cID, Input: &pb.ChaincodeInput{Args: args}}
	
	_, _, _, err = invoke(ctxt, chainID, spec, nextBlockNumber, nil, chaincodeSupport)
	if err == nil {
		chaincodeSupport.Stop(ctxt, cccid, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec})
		t.Logf("invoking <%s> should have failed", ccID)
		t.Fail()
		return
	}
}



func TestChaincodeInvokesSystemChaincode(t *testing.T) {
	testForSkip(t)
	chainID := util.GetTestChainID()

	_, chaincodeSupport, cleanup, err := initPeer(chainID)
	if err != nil {
		t.Fail()
		t.Logf("Error creating peer: %s", err)
	}

	defer cleanup()

	var ctxt = context.Background()

	var nextBlockNumber uint64 = 1

	
	url := "github.com/mcc-github/blockchain/examples/chaincode/go/passthru"

	cID := &pb.ChaincodeID{Name: "pthru", Path: url, Version: "0"}
	f := "init"
	args := util.ToChaincodeArgs(f)

	spec := &pb.ChaincodeSpec{Type: 1, ChaincodeId: cID, Input: &pb.ChaincodeInput{Args: args}}

	cccid := ccprovider.NewCCContext(chainID, "pthru", "0", "", false, nil, nil)

	_, err = deploy(ctxt, cccid, spec, nextBlockNumber, chaincodeSupport)
	nextBlockNumber++
	ccID := spec.ChaincodeId.Name
	if err != nil {
		t.Fail()
		t.Logf("Error initializing chaincode %s(%s)", ccID, err)
		chaincodeSupport.Stop(ctxt, cccid, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec})
		return
	}

	time.Sleep(time.Second)

	
	
	args = util.ToChaincodeArgs("lscc/"+chainID, "getid", chainID, "pthru")

	spec = &pb.ChaincodeSpec{Type: 1, ChaincodeId: cID, Input: &pb.ChaincodeInput{Args: args}}
	
	_, _, retval, err := invoke(ctxt, chainID, spec, nextBlockNumber, nil, chaincodeSupport)

	if err != nil {
		t.Fail()
		t.Logf("Error invoking <%s>: %s", ccID, err)
		chaincodeSupport.Stop(ctxt, cccid, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec})
		return
	}

	if string(retval) != "pthru" {
		t.Fail()
		t.Logf("Expected to get back \"pthru\" from lscc but got back %s", string(retval))
		chaincodeSupport.Stop(ctxt, cccid, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec})
		return
	}

	chaincodeSupport.Stop(ctxt, cccid, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec})
}

func TestChaincodeInitializeInitError(t *testing.T) {
	testForSkip(t)
	testCases := []struct {
		name          string
		chaincodeType pb.ChaincodeSpec_Type
		chaincodePath string
		args          []string
	}{
		{"NotSuccessResponse", pb.ChaincodeSpec_GOLANG, chaincodeExample02GolangPath, []string{"init", "not", "enough", "args"}},
		{"NotSuccessResponse", pb.ChaincodeSpec_JAVA, chaincodeExample02JavaPath, []string{"init", "not", "enough", "args"}},
		{"RuntimeException", pb.ChaincodeSpec_JAVA, chaincodeExample06JavaPath, []string{"runtimeException"}},
	}

	channel := util.GetTestChainID()

	for _, tc := range testCases {
		t.Run(tc.name+"_"+tc.chaincodeType.String(), func(t *testing.T) {

			if tc.chaincodeType == pb.ChaincodeSpec_JAVA && runtime.GOARCH != "amd64" {
				t.Skip("No Java chaincode support yet on non-amd64.")
			}

			
			listener, chaincodeSupport, cleanup, err := initPeer(channel)
			if err != nil {
				t.Errorf("Error creating peer: %s", err)
			} else {
				defer finitPeer(listener, channel)
			}

			var nextBlockNumber uint64

			
			chaincodeName := generateChaincodeName(tc.chaincodeType)
			chaincodePath := tc.chaincodePath
			chaincodeVersion := "1.0.0.0"
			chaincodeType := tc.chaincodeType
			chaincodeDeployArgs := util.ArrayToChaincodeArgs(tc.args)

			
			_, chaincodeCtx, err := deployChaincode(context.Background(), chaincodeName, chaincodeVersion, chaincodeType, chaincodePath, chaincodeDeployArgs, nil, channel, nextBlockNumber, chaincodeSupport)

			
			if err == nil {
				stopChaincode(context.Background(), chaincodeCtx, chaincodeSupport)
				t.Fatal("Deployment should have failed.")
			}
			t.Log(err)
			cleanup()
		})
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

	testutil.SetupTestLogging()

	
	var numProcsDesired = viper.GetInt("peer.gomaxprocs")
	chaincodeLogger.Debugf("setting Number of procs to %d, was %d\n", numProcsDesired, runtime.GOMAXPROCS(numProcsDesired))

	
	err = factory.InitFactories(nil)
	if err != nil {
		panic(fmt.Errorf("Could not initialize BCCSP Factories [%s]", err))
	}
}

func deployChaincode(ctx context.Context, name string, version string, chaincodeType pb.ChaincodeSpec_Type, path string, args [][]byte, creator []byte, channel string, nextBlockNumber uint64, chaincodeSupport *ChaincodeSupport) ([]byte, *ccprovider.CCContext, error) {
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

	signedProposal, proposal := putils.MockSignedEndorserProposal2OrPanic(channel, chaincodeSpec, signer)

	chaincodeCtx := ccprovider.NewCCContext(channel, name, version, "", false, signedProposal, proposal)

	result, err := deploy(ctx, chaincodeCtx, chaincodeSpec, nextBlockNumber, chaincodeSupport)
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


func (c *CreatorPolicy) Evaluate(signatureSet []*common.SignedData) error {
	for _, value := range c.Creators {
		if bytes.Compare(signatureSet[0].Identity, value) == 0 {
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
