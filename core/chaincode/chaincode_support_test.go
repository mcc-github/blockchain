/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chaincode

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/common/crypto/tlsgen"
	"github.com/mcc-github/blockchain/common/flogging"
	commonledger "github.com/mcc-github/blockchain/common/ledger"
	mc "github.com/mcc-github/blockchain/common/mocks/config"
	mocklgr "github.com/mcc-github/blockchain/common/mocks/ledger"
	mockpeer "github.com/mcc-github/blockchain/common/mocks/peer"
	"github.com/mcc-github/blockchain/common/util"
	"github.com/mcc-github/blockchain/core/aclmgmt/mocks"
	"github.com/mcc-github/blockchain/core/aclmgmt/resources"
	"github.com/mcc-github/blockchain/core/chaincode/accesscontrol"
	"github.com/mcc-github/blockchain/core/chaincode/mock"
	"github.com/mcc-github/blockchain/core/chaincode/platforms"
	"github.com/mcc-github/blockchain/core/chaincode/platforms/golang"
	"github.com/mcc-github/blockchain/core/chaincode/shim"
	"github.com/mcc-github/blockchain/core/common/ccprovider"
	"github.com/mcc-github/blockchain/core/config"
	"github.com/mcc-github/blockchain/core/container"
	"github.com/mcc-github/blockchain/core/container/dockercontroller"
	"github.com/mcc-github/blockchain/core/container/inproccontroller"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/ledger/ledgermgmt"
	cmp "github.com/mcc-github/blockchain/core/mocks/peer"
	"github.com/mcc-github/blockchain/core/peer"
	"github.com/mcc-github/blockchain/core/policy"
	"github.com/mcc-github/blockchain/core/scc"
	"github.com/mcc-github/blockchain/core/scc/lscc"
	mspmgmt "github.com/mcc-github/blockchain/msp/mgmt"
	plgr "github.com/mcc-github/blockchain/protos/ledger/queryresult"
	pb "github.com/mcc-github/blockchain/protos/peer"
	putils "github.com/mcc-github/blockchain/protos/utils"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

var globalBlockNum map[string]uint64

type mockResultsIterator struct {
	current int
	kvs     []*plgr.KV
}

func (mri *mockResultsIterator) Next() (commonledger.QueryResult, error) {
	if mri.current == len(mri.kvs) {
		return nil, nil
	}
	kv := mri.kvs[mri.current]
	mri.current = mri.current + 1

	return kv, nil
}

func (mri *mockResultsIterator) Close() {
	mri.current = len(mri.kvs)
}

type mockExecQuerySimulator struct {
	txsim ledger.TxSimulator
	mocklgr.MockQueryExecutor
	resultsIter map[string]map[string]*mockResultsIterator
}

func (meqe *mockExecQuerySimulator) GetHistoryForKey(namespace, query string) (commonledger.ResultsIterator, error) {
	return meqe.commonQuery(namespace, query)
}

func (meqe *mockExecQuerySimulator) ExecuteQuery(namespace, query string) (commonledger.ResultsIterator, error) {
	return meqe.commonQuery(namespace, query)
}

func (meqe *mockExecQuerySimulator) commonQuery(namespace, query string) (commonledger.ResultsIterator, error) {
	if meqe.resultsIter == nil {
		return nil, fmt.Errorf("query executor not initialized")
	}
	nsiter := meqe.resultsIter[namespace]
	if nsiter == nil {
		return nil, fmt.Errorf("namespace %v not found for %s", namespace, query)
	}
	iter := nsiter[query]
	if iter == nil {
		fmt.Printf("iter not found for query %s\n", query)
	}
	return iter, nil
}

func (meqe *mockExecQuerySimulator) SetState(namespace string, key string, value []byte) error {
	if meqe.txsim == nil {
		return fmt.Errorf("SetState txsimulator not initialed")
	}
	return meqe.txsim.SetState(namespace, key, value)
}

func (meqe *mockExecQuerySimulator) DeleteState(namespace string, key string) error {
	if meqe.txsim == nil {
		return fmt.Errorf("SetState txsimulator not initialed")
	}
	return meqe.txsim.DeleteState(namespace, key)
}

func (meqe *mockExecQuerySimulator) SetStateMultipleKeys(namespace string, kvs map[string][]byte) error {
	if meqe.txsim == nil {
		return fmt.Errorf("SetState txsimulator not initialed")
	}
	return meqe.txsim.SetStateMultipleKeys(namespace, kvs)
}

func (meqe *mockExecQuerySimulator) ExecuteUpdate(query string) error {
	if meqe.txsim == nil {
		return fmt.Errorf("SetState txsimulator not initialed")
	}
	return meqe.txsim.ExecuteUpdate(query)
}

func (meqe *mockExecQuerySimulator) GetTxSimulationResults() ([]byte, error) {
	if meqe.txsim == nil {
		return nil, fmt.Errorf("SetState txsimulator not initialed")
	}
	simRes, err := meqe.txsim.GetTxSimulationResults()
	if err != nil {
		return nil, err
	}
	return simRes.GetPubSimulationBytes()
}

var mockAclProvider *mocks.MockACLProvider


func initMockPeer(chainIDs ...string) (*ChaincodeSupport, error) {
	msi := &cmp.MockSupportImpl{
		GetApplicationConfigRv:     &mc.MockApplication{CapabilitiesRv: &mc.MockApplicationCapabilities{}},
		GetApplicationConfigBoolRv: true,
	}

	ipRegistry := inproccontroller.NewRegistry()
	sccp := &scc.Provider{Peer: peer.Default, PeerSupport: msi, Registrar: ipRegistry}

	mockAclProvider = &mocks.MockACLProvider{}
	mockAclProvider.Reset()

	peer.MockInitialize()

	mspGetter := func(cid string) []string {
		return []string{"SampleOrg"}
	}

	peer.MockSetMSPIDGetter(mspGetter)

	ccprovider.SetChaincodesPath(ccprovider.GetChaincodeInstallPathFromViper())
	ca, _ := tlsgen.NewCA()
	certGenerator := accesscontrol.NewAuthenticator(ca)
	config := GlobalConfig()
	config.StartupTimeout = 10 * time.Second
	config.ExecuteTimeout = 1 * time.Second
	pr := platforms.NewRegistry(&golang.Platform{})
	lsccImpl := lscc.New(sccp, mockAclProvider, pr)
	chaincodeSupport := NewChaincodeSupport(
		config,
		"0.0.0.0:7052",
		true,
		ca.CertBytes(),
		certGenerator,
		&ccprovider.CCInfoFSImpl{},
		lsccImpl,
		mockAclProvider,
		container.NewVMController(
			map[string]container.VMProvider{
				dockercontroller.ContainerType: dockercontroller.NewProvider("", ""),
				inproccontroller.ContainerType: ipRegistry,
			},
		),
		sccp,
		pr,
	)
	ipRegistry.ChaincodeSupport = chaincodeSupport

	
	policy.RegisterPolicyCheckerFactory(&mockPolicyCheckerFactory{})

	ccp := &CCProviderImpl{cs: chaincodeSupport}

	sccp.RegisterSysCC(lsccImpl)

	globalBlockNum = make(map[string]uint64, len(chainIDs))
	for _, id := range chainIDs {
		if err := peer.MockCreateChain(id); err != nil {
			return nil, err
		}
		sccp.DeploySysCCs(id, ccp)
		
		if id != util.GetTestChainID() {
			mspmgmt.XXXSetMSPManager(id, mspmgmt.GetManagerForChain(util.GetTestChainID()))
		}
		globalBlockNum[id] = 1
	}

	return chaincodeSupport, nil
}

func finitMockPeer(chainIDs ...string) {
	for _, c := range chainIDs {
		if lgr := peer.GetLedger(c); lgr != nil {
			lgr.Close()
		}
	}
	ledgermgmt.CleanupTestEnv()
	ledgerPath := config.GetPath("peer.fileSystemPath")
	os.RemoveAll(ledgerPath)
	os.RemoveAll(filepath.Join(os.TempDir(), "mcc-github"))
}


var mockPeerCCSupport = mockpeer.NewMockPeerSupport()

func mockChaincodeStreamGetter(name string) (shim.PeerChaincodeStream, error) {
	return mockPeerCCSupport.GetCC(name)
}

type mockCCLauncher struct {
	execTime *time.Duration
	resp     error
	retErr   error
	notfyb   bool
}

func (ccl *mockCCLauncher) launch(ctxt context.Context, notfy chan bool) error {
	if ccl.execTime != nil {
		time.Sleep(*ccl.execTime)
	}

	
	if ccl.resp == nil {
		notfy <- ccl.notfyb
	}

	return ccl.retErr
}

func setupcc(name string) (*mockpeer.MockCCComm, *mockpeer.MockCCComm) {
	send := make(chan *pb.ChaincodeMessage)
	recv := make(chan *pb.ChaincodeMessage)
	peerSide, _ := mockPeerCCSupport.AddCC(name, recv, send)
	peerSide.SetName("peer")
	ccSide := mockPeerCCSupport.GetCCMirror(name)
	ccSide.SetPong(true)
	return peerSide, ccSide
}


func setuperror() chan error {
	return make(chan error)
}

func processDone(t *testing.T, done chan error, expecterr bool) {
	var err error
	if done != nil {
		err = <-done
	}
	if expecterr != (err != nil) {
		if err == nil {
			t.Fatalf("Expected error but got success")
		} else {
			t.Fatalf("Expected success but got error %s", err)
		}
	}
}

func startTx(t *testing.T, chainID string, cis *pb.ChaincodeInvocationSpec, txId string) (*ccprovider.TransactionParams, ledger.TxSimulator) {
	creator := []byte([]byte("Alice"))
	sprop, prop := putils.MockSignedEndorserProposalOrPanic(chainID, cis.ChaincodeSpec, creator, []byte("msg1"))
	txsim, hqe, err := startTxSimulation(chainID, txId)
	if err != nil {
		t.Fatalf("getting txsimulator failed %s", err)
	}

	txParams := &ccprovider.TransactionParams{
		ChannelID:            chainID,
		TxID:                 txId,
		Proposal:             prop,
		SignedProp:           sprop,
		TXSimulator:          txsim,
		HistoryQueryExecutor: hqe,
	}

	return txParams, txsim
}

func endTx(t *testing.T, txParams *ccprovider.TransactionParams, txsim ledger.TxSimulator, cis *pb.ChaincodeInvocationSpec) {
	if err := endTxSimulationCIS(txParams.ChannelID, cis.ChaincodeSpec.ChaincodeId, txParams.TxID, txsim, []byte("invoke"), true, cis, globalBlockNum[txParams.ChannelID]); err != nil {
		t.Fatalf("simulation failed with error %s", err)
	}
	globalBlockNum[txParams.ChannelID] = globalBlockNum[txParams.ChannelID] + 1
}

func execCC(t *testing.T, txParams *ccprovider.TransactionParams, ccSide *mockpeer.MockCCComm, cccid *ccprovider.CCContext, waitForERROR bool, expectExecErr bool, done chan error, cis *pb.ChaincodeInvocationSpec, respSet *mockpeer.MockResponseSet, chaincodeSupport *ChaincodeSupport) error {
	ccSide.SetResponses(respSet)

	resp, _, err := chaincodeSupport.Execute(txParams, cccid, cis.ChaincodeSpec.Input)
	if err == nil && resp.Status != shim.OK {
		err = errors.New(resp.Message)
	}

	if err == nil && expectExecErr {
		t.Fatalf("expected error but succeeded")
	} else if err != nil && !expectExecErr {
		t.Fatalf("exec failed with %s", err)
	}

	
	processDone(t, done, waitForERROR)

	return nil
}


func startCC(t *testing.T, channelID string, ccname string, chaincodeSupport *ChaincodeSupport) (*mockpeer.MockCCComm, *mockpeer.MockCCComm) {
	peerSide, ccSide := setupcc(ccname)
	defer mockPeerCCSupport.RemoveCC(ccname)
	flogging.SetModuleLevel("chaincode", "debug")
	
	go func() {
		chaincodeSupport.HandleChaincodeStream(peerSide)
	}()

	done := setuperror()

	errorFunc := func(ind int, err error) {
		done <- err
	}

	ccDone := make(chan struct{})
	defer close(ccDone)

	
	go func() {
		respSet := &mockpeer.MockResponseSet{
			DoneFunc:  errorFunc,
			ErrorFunc: nil,
			Responses: []*mockpeer.MockResponse{
				{&pb.ChaincodeMessage{Type: pb.ChaincodeMessage_REGISTERED}, nil},
				{&pb.ChaincodeMessage{Type: pb.ChaincodeMessage_READY}, nil},
			},
		}
		ccSide.SetResponses(respSet)
		ccSide.Run(ccDone)
	}()

	ccSide.Send(&pb.ChaincodeMessage{Type: pb.ChaincodeMessage_REGISTER, Payload: putils.MarshalOrPanic(&pb.ChaincodeID{Name: ccname + ":0"}), Txid: "0", ChannelId: channelID})

	
	processDone(t, done, false)

	return peerSide, ccSide
}

func getTarGZ(t *testing.T, name string, contents []byte) []byte {
	startTime := time.Now()
	inputbuf := bytes.NewBuffer(nil)
	gw := gzip.NewWriter(inputbuf)
	tr := tar.NewWriter(gw)
	size := int64(len(contents))

	tr.WriteHeader(&tar.Header{Name: name, Size: size, ModTime: startTime, AccessTime: startTime, ChangeTime: startTime})
	tr.Write(contents)
	tr.Close()
	gw.Close()
	ioutil.WriteFile("/tmp/t.gz", inputbuf.Bytes(), 0644)
	return inputbuf.Bytes()
}


func deployCC(t *testing.T, txParams *ccprovider.TransactionParams, cccid *ccprovider.CCContext, spec *pb.ChaincodeSpec, chaincodeSupport *ChaincodeSupport) {
	
	code := getTarGZ(t, "src/dummy/dummy.go", []byte("code"))
	cds := &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec, CodePackage: code}

	
	ccprovider.PutChaincodeIntoFS(cds)

	b := putils.MarshalOrPanic(cds)

	sysCCVers := util.GetSysCCVersion()

	
	lsccSpec := &pb.ChaincodeInvocationSpec{ChaincodeSpec: &pb.ChaincodeSpec{Type: pb.ChaincodeSpec_GOLANG, ChaincodeId: &pb.ChaincodeID{Name: "lscc", Version: sysCCVers}, Input: &pb.ChaincodeInput{Args: [][]byte{[]byte("deploy"), []byte(txParams.ChannelID), b}}}}

	lsccid := &ccprovider.CCContext{
		Name:    lsccSpec.ChaincodeSpec.ChaincodeId.Name,
		Version: sysCCVers,
	}

	
	if _, _, err := chaincodeSupport.Execute(txParams, lsccid, lsccSpec.ChaincodeSpec.Input); err != nil {
		t.Fatalf("Error deploying chaincode %v (err: %s)", cccid, err)
	}
}

func initializeCC(t *testing.T, chainID, ccname string, ccSide *mockpeer.MockCCComm, chaincodeSupport *ChaincodeSupport) error {
	done := setuperror()

	errorFunc := func(ind int, err error) {
		done <- err
	}

	chaincodeID := &pb.ChaincodeID{Name: ccname, Version: "0"}
	ci := &pb.ChaincodeInput{Args: [][]byte{[]byte("init"), []byte("A"), []byte("100"), []byte("B"), []byte("200")}, Decorations: nil}
	cis := &pb.ChaincodeInvocationSpec{ChaincodeSpec: &pb.ChaincodeSpec{Type: pb.ChaincodeSpec_Type(pb.ChaincodeSpec_Type_value["GOLANG"]), ChaincodeId: chaincodeID, Input: ci}}

	txid := util.GenerateUUID()
	txParams, txsim := startTx(t, chainID, cis, txid)

	
	mockAclProvider.Reset()
	mockAclProvider.On("CheckACL", resources.Lscc_GetDeploymentSpec, chainID, txParams.SignedProp).Return(nil)
	mockAclProvider.On("CheckACL", resources.Lscc_GetChaincodeData, chainID, txParams.SignedProp).Return(nil)
	mockAclProvider.On("CheckACL", resources.Peer_Propose, chainID, txParams.SignedProp).Return(nil)

	
	resp := &mockpeer.MockResponse{
		RecvMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_TRANSACTION},
		RespMsg: &pb.ChaincodeMessage{
			Type:      pb.ChaincodeMessage_COMPLETED,
			Payload:   putils.MarshalOrPanic(&pb.Response{Status: shim.OK, Payload: []byte("init succeeded")}),
			Txid:      "unknowntxid",
			ChannelId: chainID,
		},
	}
	respSet := &mockpeer.MockResponseSet{
		DoneFunc:  errorFunc,
		ErrorFunc: nil,
		Responses: []*mockpeer.MockResponse{resp},
	}

	cccid := &ccprovider.CCContext{
		Name:    ccname,
		Version: "0",
	}
	execCC(t, txParams, ccSide, cccid, false, true, done, cis, respSet, chaincodeSupport)

	
	resp.RespMsg.(*pb.ChaincodeMessage).Txid = txid

	badcccid := &ccprovider.CCContext{
		Name:    ccname,
		Version: "unknownver",
	}

	
	
	execCC(t, txParams, ccSide, badcccid, false, true, nil, cis, respSet, chaincodeSupport)

	
	
	
	
	
	
	

	respSet = &mockpeer.MockResponseSet{
		DoneFunc:  errorFunc,
		ErrorFunc: nil,
		Responses: []*mockpeer.MockResponse{
			{&pb.ChaincodeMessage{Type: pb.ChaincodeMessage_TRANSACTION}, &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_PUT_STATE, Payload: putils.MarshalOrPanic(&pb.PutState{Collection: "", Key: "A", Value: []byte("100")}), Txid: txid, ChannelId: chainID}},
			{&pb.ChaincodeMessage{Type: pb.ChaincodeMessage_RESPONSE}, &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_PUT_STATE, Payload: putils.MarshalOrPanic(&pb.PutState{Collection: "", Key: "B", Value: []byte("200")}), Txid: txid, ChannelId: chainID}},
			{&pb.ChaincodeMessage{Type: pb.ChaincodeMessage_RESPONSE}, &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_COMPLETED, Payload: putils.MarshalOrPanic(&pb.Response{Status: shim.OK, Payload: []byte("OK")}), ChaincodeEvent: &pb.ChaincodeEvent{ChaincodeId: ccname}, Txid: txid, ChannelId: chainID}},
		},
	}

	execCC(t, txParams, ccSide, cccid, false, false, done, cis, respSet, chaincodeSupport)

	endTx(t, txParams, txsim, cis)

	return nil
}

func invokeCC(t *testing.T, chainID, ccname string, ccSide *mockpeer.MockCCComm, chaincodeSupport *ChaincodeSupport) error {
	done := setuperror()

	errorFunc := func(ind int, err error) {
		done <- err
	}

	chaincodeID := &pb.ChaincodeID{Name: ccname, Version: "0"}
	ci := &pb.ChaincodeInput{Args: [][]byte{[]byte("invoke"), []byte("A"), []byte("B"), []byte("10")}, Decorations: nil}
	cis := &pb.ChaincodeInvocationSpec{ChaincodeSpec: &pb.ChaincodeSpec{Type: pb.ChaincodeSpec_Type(pb.ChaincodeSpec_Type_value["GOLANG"]), ChaincodeId: chaincodeID, Input: ci}}
	txid := util.GenerateUUID()
	txParams, txsim := startTx(t, chainID, cis, txid)

	
	mockAclProvider.Reset()
	mockAclProvider.On("CheckACL", resources.Lscc_GetDeploymentSpec, chainID, txParams.SignedProp).Return(nil)
	mockAclProvider.On("CheckACL", resources.Lscc_GetChaincodeData, chainID, txParams.SignedProp).Return(nil)
	mockAclProvider.On("CheckACL", resources.Peer_Propose, chainID, txParams.SignedProp).Return(nil)

	respSet := &mockpeer.MockResponseSet{
		DoneFunc:  errorFunc,
		ErrorFunc: nil,
		Responses: []*mockpeer.MockResponse{
			{&pb.ChaincodeMessage{Type: pb.ChaincodeMessage_TRANSACTION}, &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_GET_STATE, Payload: putils.MarshalOrPanic(&pb.GetState{Collection: "", Key: "A"}), Txid: txid, ChannelId: chainID}},
			{&pb.ChaincodeMessage{Type: pb.ChaincodeMessage_RESPONSE}, &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_GET_STATE, Payload: putils.MarshalOrPanic(&pb.GetState{Collection: "", Key: "B"}), Txid: txid, ChannelId: chainID}},
			{&pb.ChaincodeMessage{Type: pb.ChaincodeMessage_TRANSACTION}, &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_PUT_STATE, Payload: putils.MarshalOrPanic(&pb.PutState{Collection: "", Key: "A", Value: []byte("90")}), Txid: txid, ChannelId: chainID}},
			{&pb.ChaincodeMessage{Type: pb.ChaincodeMessage_RESPONSE}, &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_PUT_STATE, Payload: putils.MarshalOrPanic(&pb.PutState{Collection: "", Key: "B", Value: []byte("210")}), Txid: txid, ChannelId: chainID}},
			{&pb.ChaincodeMessage{Type: pb.ChaincodeMessage_RESPONSE}, &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_PUT_STATE, Payload: putils.MarshalOrPanic(&pb.PutState{Collection: "", Key: "TODEL", Value: []byte("-to-be-deleted-")}), Txid: txid, ChannelId: chainID}},
			{&pb.ChaincodeMessage{Type: pb.ChaincodeMessage_RESPONSE}, &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_COMPLETED, Payload: putils.MarshalOrPanic(&pb.Response{Status: shim.OK, Payload: []byte("OK")}), Txid: txid, ChannelId: chainID}},
		},
	}

	cccid := &ccprovider.CCContext{
		Name:    ccname,
		Version: "0",
	}
	execCC(t, txParams, ccSide, cccid, false, false, done, cis, respSet, chaincodeSupport)

	
	respSet = &mockpeer.MockResponseSet{
		DoneFunc:  errorFunc,
		ErrorFunc: nil,
		Responses: []*mockpeer.MockResponse{
			{&pb.ChaincodeMessage{Type: pb.ChaincodeMessage_TRANSACTION}, &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_GET_STATE, Payload: putils.MarshalOrPanic(&pb.GetState{Collection: "", Key: "TODEL"}), Txid: "3", ChannelId: chainID}},
			{&pb.ChaincodeMessage{Type: pb.ChaincodeMessage_RESPONSE}, &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_DEL_STATE, Payload: putils.MarshalOrPanic(&pb.DelState{Collection: "", Key: "TODEL"}), Txid: "3", ChannelId: chainID}},
			{&pb.ChaincodeMessage{Type: pb.ChaincodeMessage_RESPONSE}, &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_COMPLETED, Payload: putils.MarshalOrPanic(&pb.Response{Status: shim.OK, Payload: []byte("OK")}), Txid: "3", ChannelId: chainID}},
		},
	}

	txParams.TxID = "3"
	execCC(t, txParams, ccSide, cccid, false, false, done, cis, respSet, chaincodeSupport)

	
	respSet = &mockpeer.MockResponseSet{
		DoneFunc:  errorFunc,
		ErrorFunc: nil,
		Responses: []*mockpeer.MockResponse{
			{&pb.ChaincodeMessage{Type: pb.ChaincodeMessage_TRANSACTION}, &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_GET_STATE, Payload: putils.MarshalOrPanic(&pb.GetState{Collection: "", Key: "TODEL"}), Txid: "4", ChannelId: chainID}},
			{&pb.ChaincodeMessage{Type: pb.ChaincodeMessage_ERROR}, &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_COMPLETED, Payload: putils.MarshalOrPanic(&pb.Response{Status: shim.ERROR, Message: "variable not found"}), Txid: "4", ChannelId: chainID}},
		},
	}

	txParams.TxID = "4"
	execCC(t, txParams, ccSide, cccid, false, true, done, cis, respSet, chaincodeSupport)

	endTx(t, txParams, txsim, cis)

	return nil
}

func invokePrivateDataGetPutDelCC(t *testing.T, chainID, ccname string, ccSide *mockpeer.MockCCComm, chaincodeSupport *ChaincodeSupport) error {
	done := setuperror()

	errorFunc := func(ind int, err error) {
		done <- err
	}

	chaincodeID := &pb.ChaincodeID{Name: ccname, Version: "0"}
	ci := &pb.ChaincodeInput{Args: [][]byte{[]byte("invokePrivateData")}, Decorations: nil}
	cis := &pb.ChaincodeInvocationSpec{ChaincodeSpec: &pb.ChaincodeSpec{Type: pb.ChaincodeSpec_Type(pb.ChaincodeSpec_Type_value["GOLANG"]), ChaincodeId: chaincodeID, Input: ci}}
	txid := util.GenerateUUID()
	txParams, txsim := startTx(t, chainID, cis, txid)

	respSet := &mockpeer.MockResponseSet{
		DoneFunc:  errorFunc,
		ErrorFunc: nil,
		Responses: []*mockpeer.MockResponse{
			{&pb.ChaincodeMessage{Type: pb.ChaincodeMessage_TRANSACTION}, &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_GET_STATE, Payload: putils.MarshalOrPanic(&pb.GetState{Collection: "", Key: "C"}), Txid: txid, ChannelId: chainID}},
			{&pb.ChaincodeMessage{Type: pb.ChaincodeMessage_ERROR}, &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_COMPLETED, Payload: putils.MarshalOrPanic(&pb.Response{Status: shim.ERROR, Message: "variable not found"}), Txid: txid, ChannelId: chainID}},
		},
	}

	cccid := &ccprovider.CCContext{
		Name:    ccname,
		Version: "0",
	}
	execCC(t, txParams, ccSide, cccid, false, true, done, cis, respSet, chaincodeSupport)

	respSet = &mockpeer.MockResponseSet{
		DoneFunc:  errorFunc,
		ErrorFunc: nil,
		Responses: []*mockpeer.MockResponse{
			{&pb.ChaincodeMessage{Type: pb.ChaincodeMessage_RESPONSE}, &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_GET_STATE, Payload: putils.MarshalOrPanic(&pb.GetState{Collection: "c1", Key: "C"}), Txid: txid, ChannelId: chainID}},
			{&pb.ChaincodeMessage{Type: pb.ChaincodeMessage_RESPONSE}, &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_PUT_STATE, Payload: putils.MarshalOrPanic(&pb.PutState{Collection: "c1", Key: "C", Value: []byte("310")}), Txid: txid, ChannelId: chainID}},
			{&pb.ChaincodeMessage{Type: pb.ChaincodeMessage_RESPONSE}, &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_PUT_STATE, Payload: putils.MarshalOrPanic(&pb.PutState{Collection: "c1", Key: "A", Value: []byte("100")}), Txid: txid, ChannelId: chainID}},
			{&pb.ChaincodeMessage{Type: pb.ChaincodeMessage_RESPONSE}, &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_PUT_STATE, Payload: putils.MarshalOrPanic(&pb.PutState{Collection: "c1", Key: "B", Value: []byte("100")}), Txid: txid, ChannelId: chainID}},
			{&pb.ChaincodeMessage{Type: pb.ChaincodeMessage_RESPONSE}, &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_DEL_STATE, Payload: putils.MarshalOrPanic(&pb.DelState{Collection: "c2", Key: "C"}), Txid: txid, ChannelId: chainID}},
			{&pb.ChaincodeMessage{Type: pb.ChaincodeMessage_RESPONSE}, &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_COMPLETED, Payload: putils.MarshalOrPanic(&pb.Response{Status: shim.OK, Payload: []byte("OK")}), Txid: txid, ChannelId: chainID}},
		},
	}

	cccid = &ccprovider.CCContext{
		Name:    ccname,
		Version: "0",
	}
	execCC(t, txParams, ccSide, cccid, false, false, done, cis, respSet, chaincodeSupport)

	endTx(t, txParams, txsim, cis)

	txid = util.GenerateUUID()
	txParams, txsim = startTx(t, chainID, cis, txid)

	respSet = &mockpeer.MockResponseSet{
		DoneFunc:  errorFunc,
		ErrorFunc: nil,
		Responses: []*mockpeer.MockResponse{
			{&pb.ChaincodeMessage{Type: pb.ChaincodeMessage_TRANSACTION}, &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_GET_STATE, Payload: putils.MarshalOrPanic(&pb.GetState{Collection: "c2", Key: "C"}), Txid: txid, ChannelId: chainID}},
			{&pb.ChaincodeMessage{Type: pb.ChaincodeMessage_ERROR}, &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_COMPLETED, Payload: putils.MarshalOrPanic(&pb.Response{Status: shim.ERROR, Message: "variable not found"}), Txid: txid, ChannelId: chainID}},
		},
	}

	cccid = &ccprovider.CCContext{
		Name:    ccname,
		Version: "0",
	}
	execCC(t, txParams, ccSide, cccid, false, true, done, cis, respSet, chaincodeSupport)

	endTx(t, txParams, txsim, cis)

	return nil
}

func getQueryStateByRange(t *testing.T, collection, chainID, ccname string, ccSide *mockpeer.MockCCComm, chaincodeSupport *ChaincodeSupport) error {
	done := setuperror()

	errorFunc := func(ind int, err error) {
		done <- err
	}

	chaincodeID := &pb.ChaincodeID{Name: ccname, Version: "0"}
	ci := &pb.ChaincodeInput{Args: [][]byte{[]byte("invoke"), []byte("A"), []byte("B"), []byte("10")}, Decorations: nil}
	cis := &pb.ChaincodeInvocationSpec{ChaincodeSpec: &pb.ChaincodeSpec{Type: pb.ChaincodeSpec_Type(pb.ChaincodeSpec_Type_value["GOLANG"]), ChaincodeId: chaincodeID, Input: ci}}
	txid := util.GenerateUUID()
	txParams, txsim := startTx(t, chainID, cis, txid)

	
	mockAclProvider.Reset()
	mockAclProvider.On("CheckACL", resources.Lscc_GetDeploymentSpec, chainID, txParams.SignedProp).Return(nil)
	mockAclProvider.On("CheckACL", resources.Lscc_GetChaincodeData, chainID, txParams.SignedProp).Return(nil)
	mockAclProvider.On("CheckACL", resources.Peer_Propose, chainID, txParams.SignedProp).Return(nil)

	
	queryStateNextFunc := func(reqMsg *pb.ChaincodeMessage) *pb.ChaincodeMessage {
		qr := &pb.QueryResponse{}
		proto.Unmarshal(reqMsg.Payload, qr)
		return &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_QUERY_STATE_NEXT, Payload: putils.MarshalOrPanic(&pb.QueryStateNext{Id: qr.Id}), Txid: txid, ChannelId: chainID}
	}
	queryStateCloseFunc := func(reqMsg *pb.ChaincodeMessage) *pb.ChaincodeMessage {
		qr := &pb.QueryResponse{}
		proto.Unmarshal(reqMsg.Payload, qr)
		return &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_QUERY_STATE_CLOSE, Payload: putils.MarshalOrPanic(&pb.QueryStateClose{Id: qr.Id}), Txid: txid, ChannelId: chainID}
	}

	var mkpeer []*mockpeer.MockResponse

	mkpeer = append(mkpeer, &mockpeer.MockResponse{
		RecvMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_TRANSACTION},
		RespMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_GET_STATE_BY_RANGE, Payload: putils.MarshalOrPanic(&pb.GetStateByRange{Collection: collection, StartKey: "A", EndKey: "B"}), Txid: txid, ChannelId: chainID},
	})

	if collection == "" {
		mkpeer = append(mkpeer, &mockpeer.MockResponse{
			RecvMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_RESPONSE},
			RespMsg: queryStateNextFunc,
		})
		mkpeer = append(mkpeer, &mockpeer.MockResponse{
			RecvMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_ERROR},
			RespMsg: queryStateCloseFunc,
		})
		mkpeer = append(mkpeer, &mockpeer.MockResponse{
			RecvMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_RESPONSE},
			RespMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_COMPLETED, Payload: putils.MarshalOrPanic(&pb.Response{Status: shim.OK, Payload: []byte("OK")}), Txid: txid, ChannelId: chainID},
		})
	} else {
		
		mkpeer = append(mkpeer, &mockpeer.MockResponse{
			RecvMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_ERROR},
			RespMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_COMPLETED, Payload: putils.MarshalOrPanic(&pb.Response{Status: shim.ERROR, Message: "Not Yet Supported"}), Txid: txid, ChannelId: chainID},
		})
	}

	respSet := &mockpeer.MockResponseSet{
		DoneFunc:  errorFunc,
		ErrorFunc: nil,
		Responses: mkpeer,
	}

	cccid := &ccprovider.CCContext{
		Name:    ccname,
		Version: "0",
	}
	if collection == "" {
		execCC(t, txParams, ccSide, cccid, false, false, done, cis, respSet, chaincodeSupport)
	} else {
		execCC(t, txParams, ccSide, cccid, false, true, done, cis, respSet, chaincodeSupport)
	}

	endTx(t, txParams, txsim, cis)

	return nil
}

func cc2cc(t *testing.T, chainID, chainID2, ccname string, ccSide *mockpeer.MockCCComm, chaincodeSupport *ChaincodeSupport) error {
	calledCC := "calledCC"
	
	_, calledCCSide := startCC(t, chainID2, calledCC, chaincodeSupport)
	if calledCCSide == nil {
		t.Fatalf("start up failed for called CC")
	}
	defer calledCCSide.Quit()

	done := setuperror()

	errorFunc := func(ind int, err error) {
		done <- err
	}

	chaincodeID := &pb.ChaincodeID{Name: calledCC, Version: "0"}
	ci := &pb.ChaincodeInput{Args: [][]byte{[]byte("deploycc")}, Decorations: nil}
	cis := &pb.ChaincodeInvocationSpec{ChaincodeSpec: &pb.ChaincodeSpec{Type: pb.ChaincodeSpec_Type(pb.ChaincodeSpec_Type_value["GOLANG"]), ChaincodeId: chaincodeID, Input: ci}}
	txid := util.GenerateUUID()
	
	txParams, txsim := startTx(t, chainID, cis, txid)

	
	mockAclProvider.Reset()
	mockAclProvider.On("CheckACL", resources.Peer_ChaincodeToChaincode, chainID, txParams.SignedProp).Return(nil)
	mockAclProvider.On("CheckACL", resources.Lscc_GetDeploymentSpec, chainID, txParams.SignedProp).Return(nil)
	mockAclProvider.On("CheckACL", resources.Lscc_GetChaincodeData, chainID, txParams.SignedProp).Return(nil)
	mockAclProvider.On("CheckACL", resources.Peer_Propose, chainID, txParams.SignedProp).Return(nil)

	cccid := &ccprovider.CCContext{
		Name:    calledCC,
		Version: "0",
	}

	deployCC(t, txParams, cccid, cis.ChaincodeSpec, chaincodeSupport)

	
	endTx(t, txParams, txsim, cis)

	
	chaincodeID = &pb.ChaincodeID{Name: ccname, Version: "0"}
	ci = &pb.ChaincodeInput{Args: [][]byte{[]byte("invokecc")}, Decorations: nil}
	cis = &pb.ChaincodeInvocationSpec{ChaincodeSpec: &pb.ChaincodeSpec{Type: pb.ChaincodeSpec_Type(pb.ChaincodeSpec_Type_value["GOLANG"]), ChaincodeId: chaincodeID, Input: ci}}
	txid = util.GenerateUUID()
	txParams, txsim = startTx(t, chainID, cis, txid)

	
	mockAclProvider.Reset()
	mockAclProvider.On("CheckACL", resources.Peer_ChaincodeToChaincode, chainID2, txParams.SignedProp).Return(nil)
	mockAclProvider.On("CheckACL", resources.Peer_ChaincodeToChaincode, chainID, txParams.SignedProp).Return(nil)
	mockAclProvider.On("CheckACL", resources.Lscc_GetDeploymentSpec, chainID2, txParams.SignedProp).Return(nil)
	mockAclProvider.On("CheckACL", resources.Lscc_GetDeploymentSpec, chainID, txParams.SignedProp).Return(nil)
	mockAclProvider.On("CheckACL", resources.Lscc_GetChaincodeData, chainID2, txParams.SignedProp).Return(nil)
	mockAclProvider.On("CheckACL", resources.Lscc_GetChaincodeData, chainID, txParams.SignedProp).Return(nil)
	mockAclProvider.On("CheckACL", resources.Peer_Propose, chainID, txParams.SignedProp).Return(nil)

	sysCCVers := util.GetSysCCVersion()
	
	respSet := &mockpeer.MockResponseSet{
		DoneFunc:  errorFunc,
		ErrorFunc: nil,
		Responses: []*mockpeer.MockResponse{
			{&pb.ChaincodeMessage{Type: pb.ChaincodeMessage_TRANSACTION}, &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_INVOKE_CHAINCODE, Payload: putils.MarshalOrPanic(&pb.ChaincodeSpec{ChaincodeId: &pb.ChaincodeID{Name: "lscc:" + sysCCVers}, Input: &pb.ChaincodeInput{Args: [][]byte{{}}}}), Txid: txid, ChannelId: chainID}},
			{&pb.ChaincodeMessage{Type: pb.ChaincodeMessage_RESPONSE}, &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_INVOKE_CHAINCODE, Payload: putils.MarshalOrPanic(&pb.ChaincodeSpec{ChaincodeId: &pb.ChaincodeID{Name: "calledCC:0/" + chainID}, Input: &pb.ChaincodeInput{Args: [][]byte{{}}}}), Txid: txid, ChannelId: chainID}},
			{&pb.ChaincodeMessage{Type: pb.ChaincodeMessage_RESPONSE}, &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_INVOKE_CHAINCODE, Payload: putils.MarshalOrPanic(&pb.ChaincodeSpec{ChaincodeId: &pb.ChaincodeID{Name: "calledCC:0/" + chainID2}, Input: &pb.ChaincodeInput{Args: [][]byte{{}}}}), Txid: txid, ChannelId: chainID}},
			{&pb.ChaincodeMessage{Type: pb.ChaincodeMessage_RESPONSE}, &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_INVOKE_CHAINCODE, Payload: putils.MarshalOrPanic(&pb.ChaincodeSpec{ChaincodeId: &pb.ChaincodeID{Name: "vscc:" + sysCCVers}, Input: &pb.ChaincodeInput{Args: [][]byte{{}}}}), Txid: txid, ChannelId: chainID}},
		},
	}

	respSet2 := &mockpeer.MockResponseSet{
		DoneFunc:  nil,
		ErrorFunc: nil,
		Responses: []*mockpeer.MockResponse{
			{&pb.ChaincodeMessage{Type: pb.ChaincodeMessage_TRANSACTION}, &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_COMPLETED, Payload: putils.MarshalOrPanic(&pb.Response{Status: shim.OK, Payload: []byte("OK")}), Txid: txid, ChannelId: chainID}},
		},
	}
	calledCCSide.SetResponses(respSet2)

	cccid = &ccprovider.CCContext{
		Name:    ccname,
		Version: "0",
	}

	execCC(t, txParams, ccSide, cccid, false, true, done, cis, respSet, chaincodeSupport)

	endTx(t, txParams, txsim, cis)

	
	chaincodeID = &pb.ChaincodeID{Name: ccname, Version: "0"}
	ci = &pb.ChaincodeInput{Args: [][]byte{[]byte("invokecc")}, Decorations: nil}
	cis = &pb.ChaincodeInvocationSpec{ChaincodeSpec: &pb.ChaincodeSpec{Type: pb.ChaincodeSpec_Type(pb.ChaincodeSpec_Type_value["GOLANG"]), ChaincodeId: chaincodeID, Input: ci}}
	txid = util.GenerateUUID()
	txParams, txsim = startTx(t, chainID, cis, txid)

	mockAclProvider.Reset()
	mockAclProvider.On("CheckACL", resources.Peer_ChaincodeToChaincode, chainID2, txParams.SignedProp).Return(nil)
	mockAclProvider.On("CheckACL", resources.Peer_ChaincodeToChaincode, chainID, txParams.SignedProp).Return(errors.New("Bad ACL calling CC"))
	mockAclProvider.On("CheckACL", resources.Lscc_GetDeploymentSpec, chainID2, txParams.SignedProp).Return(nil)
	mockAclProvider.On("CheckACL", resources.Lscc_GetDeploymentSpec, chainID, txParams.SignedProp).Return(nil)
	mockAclProvider.On("CheckACL", resources.Lscc_GetChaincodeData, chainID2, txParams.SignedProp).Return(nil)
	mockAclProvider.On("CheckACL", resources.Lscc_GetChaincodeData, chainID, txParams.SignedProp).Return(nil)
	mockAclProvider.On("CheckACL", resources.Peer_Propose, chainID, txParams.SignedProp).Return(nil)
	
	respSet = &mockpeer.MockResponseSet{
		DoneFunc:  errorFunc,
		ErrorFunc: nil,
		Responses: []*mockpeer.MockResponse{
			{&pb.ChaincodeMessage{Type: pb.ChaincodeMessage_RESPONSE}, &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_INVOKE_CHAINCODE, Payload: putils.MarshalOrPanic(&pb.ChaincodeSpec{ChaincodeId: &pb.ChaincodeID{Name: "calledCC:0/" + chainID}, Input: &pb.ChaincodeInput{Args: [][]byte{{}}}}), Txid: txid, ChannelId: chainID}},
		},
	}

	respSet2 = &mockpeer.MockResponseSet{
		DoneFunc:  nil,
		ErrorFunc: nil,
		Responses: []*mockpeer.MockResponse{
			{&pb.ChaincodeMessage{Type: pb.ChaincodeMessage_TRANSACTION}, &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_COMPLETED, Payload: putils.MarshalOrPanic(&pb.Response{Status: shim.OK, Payload: []byte("OK")}), Txid: txid, ChannelId: chainID}},
		},
	}

	calledCCSide.SetResponses(respSet2)

	cccid = &ccprovider.CCContext{
		Name:    ccname,
		Version: "0",
	}

	execCC(t, txParams, ccSide, cccid, false, true, done, cis, respSet, chaincodeSupport)

	endTx(t, txParams, txsim, cis)

	return nil
}

func getQueryResult(t *testing.T, collection, chainID, ccname string, ccSide *mockpeer.MockCCComm, chaincodeSupport *ChaincodeSupport) error {
	done := setuperror()

	errorFunc := func(ind int, err error) {
		done <- err
	}

	chaincodeID := &pb.ChaincodeID{Name: ccname, Version: "0"}
	ci := &pb.ChaincodeInput{Args: [][]byte{[]byte("invoke"), []byte("A"), []byte("B"), []byte("10")}, Decorations: nil}
	cis := &pb.ChaincodeInvocationSpec{ChaincodeSpec: &pb.ChaincodeSpec{Type: pb.ChaincodeSpec_Type(pb.ChaincodeSpec_Type_value["GOLANG"]), ChaincodeId: chaincodeID, Input: ci}}
	txid := util.GenerateUUID()
	txParams, txsim := startTx(t, chainID, cis, txid)

	mockAclProvider.Reset()
	mockAclProvider.On("CheckACL", resources.Lscc_GetDeploymentSpec, chainID, txParams.SignedProp).Return(nil)
	mockAclProvider.On("CheckACL", resources.Lscc_GetChaincodeData, chainID, txParams.SignedProp).Return(nil)
	mockAclProvider.On("CheckACL", resources.Peer_Propose, chainID, txParams.SignedProp).Return(nil)

	kvs := make([]*plgr.KV, 1000)
	for i := 0; i < 1000; i++ {
		kvs[i] = &plgr.KV{Namespace: chainID, Key: fmt.Sprintf("%d", i), Value: []byte(fmt.Sprintf("%d", i))}
	}

	queryExec := &mockExecQuerySimulator{resultsIter: make(map[string]map[string]*mockResultsIterator)}
	queryExec.resultsIter[ccname] = map[string]*mockResultsIterator{"goodquery": {kvs: kvs}}

	queryExec.txsim = txParams.TXSimulator

	
	queryStateNextFunc := func(reqMsg *pb.ChaincodeMessage) *pb.ChaincodeMessage {
		qr := &pb.QueryResponse{}
		proto.Unmarshal(reqMsg.Payload, qr)
		return &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_QUERY_STATE_NEXT, Payload: putils.MarshalOrPanic(&pb.QueryStateNext{Id: qr.Id}), Txid: txid, ChannelId: chainID}
	}
	queryStateCloseFunc := func(reqMsg *pb.ChaincodeMessage) *pb.ChaincodeMessage {
		qr := &pb.QueryResponse{}
		proto.Unmarshal(reqMsg.Payload, qr)
		return &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_QUERY_STATE_CLOSE, Payload: putils.MarshalOrPanic(&pb.QueryStateClose{Id: qr.Id}), Txid: txid, ChannelId: chainID}
	}

	var mkpeer []*mockpeer.MockResponse

	mkpeer = append(mkpeer, &mockpeer.MockResponse{
		RecvMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_TRANSACTION},
		RespMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_GET_QUERY_RESULT, Payload: putils.MarshalOrPanic(&pb.GetQueryResult{Collection: "", Query: "goodquery"}), Txid: txid, ChannelId: chainID},
	})

	if collection == "" {
		mkpeer = append(mkpeer, &mockpeer.MockResponse{
			RecvMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_RESPONSE},
			RespMsg: queryStateNextFunc,
		})
		mkpeer = append(mkpeer, &mockpeer.MockResponse{
			RecvMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_ERROR},
			RespMsg: queryStateCloseFunc,
		})
		mkpeer = append(mkpeer, &mockpeer.MockResponse{
			RecvMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_RESPONSE},
			RespMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_COMPLETED, Payload: putils.MarshalOrPanic(&pb.Response{Status: shim.OK, Payload: []byte("OK")}), Txid: txid, ChannelId: chainID},
		})
	} else {
		
		mkpeer = append(mkpeer, &mockpeer.MockResponse{
			RecvMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_ERROR},
			RespMsg: &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_COMPLETED, Payload: putils.MarshalOrPanic(&pb.Response{Status: shim.ERROR, Message: "Not Yet Supported"}), Txid: txid, ChannelId: chainID},
		})
	}

	respSet := &mockpeer.MockResponseSet{
		DoneFunc:  errorFunc,
		ErrorFunc: nil,
		Responses: mkpeer,
	}

	cccid := &ccprovider.CCContext{
		Name:    ccname,
		Version: "0",
	}
	if collection == "" {
		execCC(t, txParams, ccSide, cccid, false, false, done, cis, respSet, chaincodeSupport)
	} else {
		execCC(t, txParams, ccSide, cccid, false, true, done, cis, respSet, chaincodeSupport)
	}

	endTx(t, txParams, txsim, cis)

	return nil
}

func getHistory(t *testing.T, chainID, ccname string, ccSide *mockpeer.MockCCComm, chaincodeSupport *ChaincodeSupport) error {
	done := setuperror()

	errorFunc := func(ind int, err error) {
		done <- err
	}

	chaincodeID := &pb.ChaincodeID{Name: ccname, Version: "0"}
	ci := &pb.ChaincodeInput{Args: [][]byte{[]byte("invoke"), []byte("A"), []byte("B"), []byte("10")}, Decorations: nil}
	cis := &pb.ChaincodeInvocationSpec{ChaincodeSpec: &pb.ChaincodeSpec{Type: pb.ChaincodeSpec_Type(pb.ChaincodeSpec_Type_value["GOLANG"]), ChaincodeId: chaincodeID, Input: ci}}
	txid := util.GenerateUUID()
	txParams, txsim := startTx(t, chainID, cis, txid)

	mockAclProvider.Reset()
	mockAclProvider.On("CheckACL", resources.Lscc_GetDeploymentSpec, chainID, txParams.SignedProp).Return(nil)
	mockAclProvider.On("CheckACL", resources.Lscc_GetChaincodeData, chainID, txParams.SignedProp).Return(nil)
	mockAclProvider.On("CheckACL", resources.Peer_Propose, chainID, txParams.SignedProp).Return(nil)

	kvs := make([]*plgr.KV, 1000)
	for i := 0; i < 1000; i++ {
		kvs[i] = &plgr.KV{Namespace: chainID, Key: fmt.Sprintf("%d", i), Value: []byte(fmt.Sprintf("%d", i))}
	}

	queryExec := &mockExecQuerySimulator{resultsIter: make(map[string]map[string]*mockResultsIterator)}
	queryExec.resultsIter[ccname] = map[string]*mockResultsIterator{"goodquery": {kvs: kvs}}

	queryExec.txsim = txParams.TXSimulator

	
	queryStateNextFunc := func(reqMsg *pb.ChaincodeMessage) *pb.ChaincodeMessage {
		qr := &pb.QueryResponse{}
		proto.Unmarshal(reqMsg.Payload, qr)
		return &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_QUERY_STATE_NEXT, Payload: putils.MarshalOrPanic(&pb.QueryStateNext{Id: qr.Id}), Txid: txid}
	}
	queryStateCloseFunc := func(reqMsg *pb.ChaincodeMessage) *pb.ChaincodeMessage {
		qr := &pb.QueryResponse{}
		proto.Unmarshal(reqMsg.Payload, qr)
		return &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_QUERY_STATE_CLOSE, Payload: putils.MarshalOrPanic(&pb.QueryStateClose{Id: qr.Id}), Txid: txid}
	}

	respSet := &mockpeer.MockResponseSet{
		DoneFunc:  errorFunc,
		ErrorFunc: nil,
		Responses: []*mockpeer.MockResponse{
			{&pb.ChaincodeMessage{Type: pb.ChaincodeMessage_TRANSACTION}, &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_GET_HISTORY_FOR_KEY, Payload: putils.MarshalOrPanic(&pb.GetQueryResult{Query: "goodquery"}), Txid: txid, ChannelId: chainID}},
			{&pb.ChaincodeMessage{Type: pb.ChaincodeMessage_RESPONSE}, queryStateNextFunc},
			{&pb.ChaincodeMessage{Type: pb.ChaincodeMessage_RESPONSE}, queryStateNextFunc},
			{&pb.ChaincodeMessage{Type: pb.ChaincodeMessage_ERROR}, queryStateCloseFunc},
			{&pb.ChaincodeMessage{Type: pb.ChaincodeMessage_RESPONSE}, &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_COMPLETED, Payload: putils.MarshalOrPanic(&pb.Response{Status: shim.OK, Payload: []byte("OK")}), Txid: txid, ChannelId: chainID}},
		},
	}

	cccid := &ccprovider.CCContext{
		Name:    ccname,
		Version: "0",
	}
	execCC(t, txParams, ccSide, cccid, false, false, done, cis, respSet, chaincodeSupport)

	endTx(t, txParams, txsim, cis)

	return nil
}

func getLaunchConfigs(t *testing.T, cr *ContainerRuntime) {
	gt := NewGomegaWithT(t)
	ccContext := &ccprovider.CCContext{
		Name:    "mycc",
		Version: "v0",
	}
	lc, err := cr.LaunchConfig(ccContext.GetCanonicalName(), pb.ChaincodeSpec_GOLANG.String())
	if err != nil {
		t.Fatalf("calling getLaunchConfigs() failed with error %s", err)
	}
	args := lc.Args
	envs := lc.Envs
	filesToUpload := lc.Files

	if len(args) != 2 {
		t.Fatalf("calling getLaunchConfigs() for golang chaincode should have returned an array of 2 elements for Args, but got %v", args)
	}
	if args[0] != "chaincode" || !strings.HasPrefix(args[1], "-peer.address") {
		t.Fatalf("calling getLaunchConfigs() should have returned the start command for golang chaincode, but got %v", args)
	}

	if len(envs) != 8 {
		t.Fatalf("calling getLaunchConfigs() with TLS enabled should have returned an array of 8 elements for Envs, but got %v", envs)
	}
	gt.Expect(envs).To(ContainElement("CORE_CHAINCODE_LOGGING_LEVEL=info"))
	gt.Expect(envs).To(ContainElement("CORE_CHAINCODE_LOGGING_SHIM=warning"))
	gt.Expect(envs).To(ContainElement("CORE_CHAINCODE_ID_NAME=mycc:v0"))
	gt.Expect(envs).To(ContainElement("CORE_PEER_TLS_ENABLED=true"))
	gt.Expect(envs).To(ContainElement("CORE_TLS_CLIENT_KEY_PATH=/etc/mcc-github/blockchain/client.key"))
	gt.Expect(envs).To(ContainElement("CORE_TLS_CLIENT_CERT_PATH=/etc/mcc-github/blockchain/client.crt"))
	gt.Expect(envs).To(ContainElement("CORE_PEER_TLS_ROOTCERT_FILE=/etc/mcc-github/blockchain/peer.crt"))

	if len(filesToUpload) != 3 {
		t.Fatalf("calling getLaunchConfigs() with TLS enabled should have returned an array of 3 elements for filesToUpload, but got %v", len(filesToUpload))
	}

	cr.CertGenerator = nil 
	lc, err = cr.LaunchConfig(ccContext.GetCanonicalName(), pb.ChaincodeSpec_NODE.String())
	assert.NoError(t, err)
	args = lc.Args

	if len(args) != 3 {
		t.Fatalf("calling getLaunchConfigs() for node chaincode should have returned an array of 3 elements for Args, but got %v", args)
	}

	if args[0] != "/bin/sh" || args[1] != "-c" || !strings.HasPrefix(args[2], "cd /usr/local/src; npm start -- --peer.address") {
		t.Fatalf("calling getLaunchConfigs() should have returned the start command for node.js chaincode, but got %v", args)
	}

	lc, err = cr.LaunchConfig(ccContext.GetCanonicalName(), pb.ChaincodeSpec_GOLANG.String())
	assert.NoError(t, err)

	envs = lc.Envs
	if len(envs) != 5 {
		t.Fatalf("calling getLaunchConfigs() with TLS disabled should have returned an array of 4 elements for Envs, but got %v", envs)
	}
	gt.Expect(envs).To(ContainElement("CORE_CHAINCODE_LOGGING_LEVEL=info"))
	gt.Expect(envs).To(ContainElement("CORE_CHAINCODE_LOGGING_SHIM=warning"))
	gt.Expect(envs).To(ContainElement("CORE_CHAINCODE_ID_NAME=mycc:v0"))
	gt.Expect(envs).To(ContainElement("CORE_PEER_TLS_ENABLED=false"))
}


func TestStartAndWaitSuccess(t *testing.T) {
	handlerRegistry := NewHandlerRegistry(false)
	fakeRuntime := &mock.Runtime{}
	fakeRuntime.StartStub = func(_ *ccprovider.ChaincodeContainerInfo, _ []byte) error {
		handlerRegistry.Ready("testcc:0")
		return nil
	}

	code := getTarGZ(t, "src/dummy/dummy.go", []byte("code"))
	fakePackageProvider := &mock.PackageProvider{}
	fakePackageProvider.GetChaincodeCodePackageReturns(code, nil)

	launcher := &RuntimeLauncher{
		Runtime:         fakeRuntime,
		Registry:        handlerRegistry,
		StartupTimeout:  10 * time.Second,
		PackageProvider: fakePackageProvider,
	}

	ccci := &ccprovider.ChaincodeContainerInfo{
		Type:          "GOLANG",
		Name:          "testcc",
		Version:       "0",
		ContainerType: "DOCKER",
	}

	
	err := launcher.Launch(ccci)
	if err != nil {
		t.Fatalf("expected success but failed with error %s", err)
	}
}


func TestStartAndWaitTimeout(t *testing.T) {
	fakeRuntime := &mock.Runtime{}
	fakeRuntime.StartStub = func(_ *ccprovider.ChaincodeContainerInfo, _ []byte) error {
		time.Sleep(time.Second)
		return nil
	}

	code := getTarGZ(t, "src/dummy/dummy.go", []byte("code"))
	fakePackageProvider := &mock.PackageProvider{}
	fakePackageProvider.GetChaincodeCodePackageReturns(code, nil)

	launcher := &RuntimeLauncher{
		Runtime:         fakeRuntime,
		Registry:        NewHandlerRegistry(false),
		StartupTimeout:  500 * time.Millisecond,
		PackageProvider: fakePackageProvider,
	}

	ccci := &ccprovider.ChaincodeContainerInfo{
		Type:          "GOLANG",
		Name:          "testcc",
		Version:       "0",
		ContainerType: "DOCKER",
	}

	
	err := launcher.Launch(ccci)
	if err == nil {
		t.Fatalf("expected error but succeeded")
	}
}


func TestStartAndWaitLaunchError(t *testing.T) {
	fakeRuntime := &mock.Runtime{}
	fakeRuntime.StartStub = func(_ *ccprovider.ChaincodeContainerInfo, _ []byte) error {
		return errors.New("Bad lunch; upset stomach")
	}

	code := getTarGZ(t, "src/dummy/dummy.go", []byte("code"))
	fakePackageProvider := &mock.PackageProvider{}
	fakePackageProvider.GetChaincodeCodePackageReturns(code, nil)

	launcher := &RuntimeLauncher{
		Runtime:         fakeRuntime,
		Registry:        NewHandlerRegistry(false),
		StartupTimeout:  10 * time.Second,
		PackageProvider: fakePackageProvider,
	}

	ccci := &ccprovider.ChaincodeContainerInfo{
		Type:          "GOLANG",
		Name:          "testcc",
		Version:       "0",
		ContainerType: "DOCKER",
	}

	
	err := launcher.Launch(ccci)
	if err == nil {
		t.Fatalf("expected error but succeeded")
	}
	assert.EqualError(t, err, "error starting container: Bad lunch; upset stomach")
}

func TestGetTxContextFromHandler(t *testing.T) {
	h := Handler{TXContexts: NewTransactionContexts(), SystemCCProvider: &scc.Provider{Peer: peer.Default, PeerSupport: peer.DefaultSupport, Registrar: inproccontroller.NewRegistry()}}

	chnl := "test"
	txid := "1"
	
	txContext, _ := h.getTxContextForInvoke(chnl, "1", []byte(""), "[%s]No ledger context for %s. Sending %s", 12345, "TestCC", pb.ChaincodeMessage_ERROR)
	if txContext != nil {
		t.Fatalf("expected empty txContext for empty payload")
	}

	
	peer.MockInitialize()

	err := peer.MockCreateChain(chnl)
	if err != nil {
		t.Fatalf("failed to create Peer Ledger %s", err)
	}

	pldgr := peer.GetLedger(chnl)

	
	txCtxGenerated, payload := genNewPldAndCtxFromLdgr(t, "shimTestCC", chnl, txid, pldgr, &h)

	
	txContext, ccMsg := h.getTxContextForInvoke(chnl, txid, payload, "[%s]No ledger context for %s. Sending %s", 12345, pb.ChaincodeMessage_INVOKE_CHAINCODE.String(), pb.ChaincodeMessage_ERROR)
	if txContext == nil || ccMsg != nil || txContext != txCtxGenerated {
		t.Fatalf("expected successful txContext for non empty payload and INVOKE_CHAINCODE msgType. triggerNextStateMsg: %s.", ccMsg)
	}

	
	txContext, ccMsg = h.isValidTxSim(chnl, txid,
		"[%s]No ledger context for %s. Sending %s", 12345, pb.ChaincodeMessage_PUT_STATE, pb.ChaincodeMessage_ERROR)
	if txContext == nil || ccMsg != nil || txContext != txCtxGenerated {
		t.Fatalf("expected successful txContext for non empty payload and PUT_STATE msgType. triggerNextStateMsg: %s.", ccMsg)
	}

	
	txid = "2"
	
	chnl = ""
	txCtxGenerated, payload = genNewPldAndCtxFromLdgr(t, "lscc", chnl, txid, pldgr, &h)

	
	txContext, ccMsg = h.getTxContextForInvoke(chnl, txid, payload,
		"[%s]No ledger context for %s. Sending %s", 12345, pb.ChaincodeMessage_INVOKE_CHAINCODE, pb.ChaincodeMessage_ERROR)
	if txContext == nil || ccMsg != nil || txContext != txCtxGenerated {
		t.Fatalf("expected successful txContext for non empty payload and INVOKE_CHAINCODE msgType. triggerNextStateMsg: %s.", ccMsg)
	}

	
	txid = "3"
	chnl = "TEST"
	txCtxGenerated, payload = genNewPldAndCtxFromLdgr(t, "lscc", chnl, txid, pldgr, &h)

	
	txContext, ccMsg = h.getTxContextForInvoke(chnl, txid, payload,
		"[%s]No ledger context for %s. Sending %s", 12345, pb.ChaincodeMessage_INVOKE_CHAINCODE, pb.ChaincodeMessage_ERROR)
	if txContext == nil || ccMsg != nil || txContext != txCtxGenerated {
		t.Fatalf("expected successful txContext for non empty payload and INVOKE_CHAINCODE msgType. triggerNextStateMsg: %s.", ccMsg)
	}

	
	txid = "4"
	chnl = ""
	txCtxGenerated, payload = genNewPldAndCtxFromLdgr(t, "shimTestCC", chnl, txid, pldgr, &h)
	
	txContext, ccMsg = h.getTxContextForInvoke(chnl, txid, payload,
		"[%s]No ledger context for %s. Sending %s", 12345, pb.ChaincodeMessage_INVOKE_CHAINCODE, pb.ChaincodeMessage_ERROR)
	if txContext == nil || ccMsg != nil || txContext != txCtxGenerated {
		t.Fatalf("expected successful txContext for non empty payload and INVOKE_CHAINCODE msgType. triggerNextStateMsg: %s.", ccMsg)
	}

	
	txid = "5"
	payload = genNewPld(t, "shimTestCC")
	
	txContext, ccMsg = h.getTxContextForInvoke(chnl, txid, payload,
		"[%s]No ledger context for %s. Sending %s", 12345, pb.ChaincodeMessage_INVOKE_CHAINCODE, pb.ChaincodeMessage_ERROR)
	if txContext != nil || ccMsg == nil {
		t.Fatal("expected nil txContext for non empty payload and INVOKE_CHAINCODE msgType without the ledger generating a TxContext . unexpected non nil tcContext")
	}

	
	txid = "6"
	payload = genNewPld(t, "lscc")
	
	txContext, ccMsg = h.getTxContextForInvoke(chnl, txid, payload,
		"[%s]No ledger context for %s. Sending %s", 12345, pb.ChaincodeMessage_INVOKE_CHAINCODE, pb.ChaincodeMessage_ERROR)
	if txContext != nil || ccMsg == nil {
		t.Fatal("expected nil txContext for non empty payload and INVOKE_CHAINCODE msgType without the ledger generating a TxContext . unexpected non nil tcContext")
	}
}

func genNewPldAndCtxFromLdgr(t *testing.T, ccName string, chnl string, txid string, pldgr ledger.PeerLedger, h *Handler) (*TransactionContext, []byte) {
	
	txsim, err := pldgr.NewTxSimulator(txid)
	if err != nil {
		t.Fatalf("failed to create TxSimulator %s", err)
	}
	txParams := &ccprovider.TransactionParams{
		TxID:        txid,
		ChannelID:   chnl,
		TXSimulator: txsim,
	}
	newTxCtxt, err := h.TXContexts.Create(txParams)
	if err != nil {
		t.Fatalf("Error creating TxContext by the handler for cc %s and channel '%s': %s", ccName, chnl, err)
	}
	if newTxCtxt == nil {
		t.Fatalf("Error creating TxContext: newTxCtxt created by the handler is nil for cc %s and channel '%s'.", ccName, chnl)
	}
	
	payload := genNewPld(t, ccName)
	return newTxCtxt, payload
}

func genNewPld(t *testing.T, ccName string) []byte {
	
	chaincodeID := &pb.ChaincodeID{Name: ccName, Version: "0"}
	ci := &pb.ChaincodeInput{Args: [][]byte{[]byte("deploycc")}, Decorations: nil}
	cds := &pb.ChaincodeSpec{Type: pb.ChaincodeSpec_Type(pb.ChaincodeSpec_Type_value["GOLANG"]), ChaincodeId: chaincodeID, Input: ci}
	payload, err := proto.Marshal(cds)
	if err != nil {
		t.Fatalf("failed to marshal CDS %s", err)
	}
	return payload
}

func cc2SameCC(t *testing.T, chainID, chainID2, ccname string, ccSide *mockpeer.MockCCComm, chaincodeSupport *ChaincodeSupport) {
	
	chaincodeID := &pb.ChaincodeID{Name: ccname, Version: "0"}
	ci := &pb.ChaincodeInput{Args: [][]byte{[]byte("deploycc")}, Decorations: nil}
	cis := &pb.ChaincodeInvocationSpec{ChaincodeSpec: &pb.ChaincodeSpec{Type: pb.ChaincodeSpec_Type(pb.ChaincodeSpec_Type_value["GOLANG"]), ChaincodeId: chaincodeID, Input: ci}}

	txid := util.GenerateUUID()
	txParams, txsim := startTx(t, chainID2, cis, txid)

	
	mockAclProvider.Reset()
	mockAclProvider.On("CheckACL", resources.Peer_ChaincodeToChaincode, chainID2, txParams.SignedProp).Return(nil)
	mockAclProvider.On("CheckACL", resources.Lscc_GetDeploymentSpec, chainID2, txParams.SignedProp).Return(nil)
	mockAclProvider.On("CheckACL", resources.Lscc_GetChaincodeData, chainID2, txParams.SignedProp).Return(nil)
	mockAclProvider.On("CheckACL", resources.Peer_Propose, chainID2, txParams.SignedProp).Return(nil)

	cccid := &ccprovider.CCContext{
		Name:    ccname,
		Version: "0",
	}

	deployCC(t, txParams, cccid, cis.ChaincodeSpec, chaincodeSupport)

	
	endTx(t, txParams, txsim, cis)

	done := setuperror()

	errorFunc := func(ind int, err error) {
		done <- err
	}

	
	
	ci = &pb.ChaincodeInput{Args: [][]byte{[]byte("invoke"), []byte("A"), []byte("B"), []byte("10")}, Decorations: nil}
	cis = &pb.ChaincodeInvocationSpec{ChaincodeSpec: &pb.ChaincodeSpec{Type: pb.ChaincodeSpec_Type(pb.ChaincodeSpec_Type_value["GOLANG"]), ChaincodeId: chaincodeID, Input: ci}}
	txid = util.GenerateUUID()
	txParams, txsim = startTx(t, chainID, cis, txid)

	mockAclProvider.On("CheckACL", resources.Peer_ChaincodeToChaincode, chainID, txParams.SignedProp).Return(nil)
	mockAclProvider.On("CheckACL", resources.Lscc_GetDeploymentSpec, chainID, txParams.SignedProp).Return(nil)
	mockAclProvider.On("CheckACL", resources.Lscc_GetChaincodeData, chainID, txParams.SignedProp).Return(nil)
	mockAclProvider.On("CheckACL", resources.Peer_Propose, chainID, txParams.SignedProp).Return(nil)

	mockAclProvider.On("CheckACL", resources.Peer_ChaincodeToChaincode, chainID2, txParams.SignedProp).Return(nil)
	mockAclProvider.On("CheckACL", resources.Lscc_GetDeploymentSpec, chainID2, txParams.SignedProp).Return(nil)
	mockAclProvider.On("CheckACL", resources.Lscc_GetChaincodeData, chainID2, txParams.SignedProp).Return(nil)
	mockAclProvider.On("CheckACL", resources.Peer_Propose, chainID2, txParams.SignedProp).Return(nil)

	txid = "cctosamecctx"
	respSet := &mockpeer.MockResponseSet{
		DoneFunc:  errorFunc,
		ErrorFunc: nil,
		Responses: []*mockpeer.MockResponse{
			{&pb.ChaincodeMessage{Type: pb.ChaincodeMessage_TRANSACTION}, &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_INVOKE_CHAINCODE, Payload: putils.MarshalOrPanic(&pb.ChaincodeSpec{ChaincodeId: &pb.ChaincodeID{Name: ccname + ":0/" + chainID2}, Input: &pb.ChaincodeInput{Args: [][]byte{{}}}}), Txid: txid, ChannelId: chainID}},
			{&pb.ChaincodeMessage{Type: pb.ChaincodeMessage_RESPONSE}, &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_INVOKE_CHAINCODE, Payload: putils.MarshalOrPanic(&pb.ChaincodeSpec{ChaincodeId: &pb.ChaincodeID{Name: ccname + ":0/" + chainID}, Input: &pb.ChaincodeInput{Args: [][]byte{{}}}}), Txid: txid, ChannelId: chainID}},
		},
	}

	cccid = &ccprovider.CCContext{
		Name:    ccname,
		Version: "0",
	}

	execCC(t, txParams, ccSide, cccid, false, true, done, cis, respSet, chaincodeSupport)

	endTx(t, txParams, txsim, cis)
}

func TestCCFramework(t *testing.T) {
	
	chainID := "mockchainid"
	chainID2 := "secondchain"
	chaincodeSupport, err := initMockPeer(chainID, chainID2)
	if err != nil {
		t.Fatalf("%s", err)
	}
	defer finitMockPeer(chainID, chainID2)
	
	ccname := "shimTestCC"

	
	_, ccSide := startCC(t, chainID, ccname, chaincodeSupport)
	if ccSide == nil {
		t.Fatalf("start up failed")
	}
	defer ccSide.Quit()

	
	initializeCC(t, chainID, ccname, ccSide, chaincodeSupport)

	
	handler := &Handler{chaincodeID: &pb.ChaincodeID{Name: ccname + ":0"}, SystemCCProvider: chaincodeSupport.SystemCCProvider}
	if err := chaincodeSupport.HandlerRegistry.Register(handler); err == nil {
		t.Fatalf("expected re-register to fail")
	}

	
	initializeCC(t, chainID2, ccname, ccSide, chaincodeSupport)

	
	invokeCC(t, chainID, ccname, ccSide, chaincodeSupport)

	
	
	
	
	

	
	getQueryStateByRange(t, "", chainID, ccname, ccSide, chaincodeSupport)

	
	cc2SameCC(t, chainID, chainID2, ccname, ccSide, chaincodeSupport)

	
	cc2cc(t, chainID, chainID2, ccname, ccSide, chaincodeSupport)

	
	getQueryResult(t, "", chainID, ccname, ccSide, chaincodeSupport)

	
	getHistory(t, chainID, ccname, ccSide, chaincodeSupport)

	
	cr := chaincodeSupport.Runtime.(*ContainerRuntime)
	getLaunchConfigs(t, cr)

	ccSide.Quit()
}
