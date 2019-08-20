/*
Copyright IBM Corp. 2017 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package qscc

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/mcc-github/blockchain/bccsp/sw"
	"github.com/mcc-github/blockchain/common/ledger/testutil"
	"github.com/mcc-github/blockchain/common/util"
	"github.com/mcc-github/blockchain/core/aclmgmt/mocks"
	"github.com/mcc-github/blockchain/core/aclmgmt/resources"
	"github.com/mcc-github/blockchain/core/chaincode/shim"
	"github.com/mcc-github/blockchain/core/chaincode/shim/shimtest"
	ledger2 "github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/ledger/ledgermgmt"
	"github.com/mcc-github/blockchain/core/ledger/ledgermgmt/ledgermgmttest"
	"github.com/mcc-github/blockchain/core/peer"
	"github.com/mcc-github/blockchain/protos/common"
	peer2 "github.com/mcc-github/blockchain/protos/peer"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestLedger(chainid string, path string) (*shimtest.MockStub, *peer.Peer, func(), error) {
	mockAclProvider.Reset()

	viper.Set("peer.fileSystemPath", path)
	testDir, err := ioutil.TempDir("", "qscc_test")
	if err != nil {
		return nil, nil, nil, err
	}
	ledgerMgr := ledgermgmt.NewLedgerMgr(ledgermgmttest.NewInitializer(testDir))

	cleanup := func() {
		ledgerMgr.Close()
		os.RemoveAll(testDir)
	}
	cryptoProvider, err := sw.NewDefaultSecurityLevelWithKeystore(sw.NewDummyKeyStore())
	if err != nil {
		return nil, nil, nil, err
	}
	peerInstance := &peer.Peer{
		LedgerMgr:      ledgerMgr,
		CryptoProvider: cryptoProvider,
	}
	peer.CreateMockChannel(peerInstance, chainid)

	lq := &LedgerQuerier{
		aclProvider: mockAclProvider,
		ledgers:     peerInstance,
	}
	stub := shimtest.NewMockStub("LedgerQuerier", lq)
	if res := stub.MockInit("1", nil); res.Status != shim.OK {
		return nil, peerInstance, cleanup, fmt.Errorf("Init failed for test ledger [%s] with message: %s", chainid, string(res.Message))
	}
	return stub, peerInstance, cleanup, nil
}


func resetProvider(res, chainid string, prop *peer2.SignedProposal, retErr error) *peer2.SignedProposal {
	mockAclProvider.Reset()
	mockAclProvider.On("CheckACL", res, chainid, prop).Return(retErr)
	return prop
}

func tempDir(t *testing.T, stem string) string {
	path, err := ioutil.TempDir("", "qscc-"+stem)
	require.NoError(t, err)
	return path
}

func TestQueryGetChainInfo(t *testing.T) {
	chainid := "mytestchainid1"
	path := tempDir(t, "test1")
	defer os.RemoveAll(path)

	stub, _, cleanup, err := setupTestLedger(chainid, path)
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer cleanup()

	args := [][]byte{[]byte(GetChainInfo), []byte(chainid)}
	prop := resetProvider(resources.Qscc_GetChainInfo, chainid, &peer2.SignedProposal{}, nil)
	res := stub.MockInvokeWithSignedProposal("1", args, prop)
	assert.Equal(t, int32(shim.OK), res.Status, "GetChainInfo failed with err: %s", res.Message)

	args = [][]byte{[]byte(GetChainInfo)}
	res = stub.MockInvoke("2", args)
	assert.Equal(t, int32(shim.ERROR), res.Status, "GetChainInfo should have failed because no channel id was provided")

	args = [][]byte{[]byte(GetChainInfo), []byte("fakechainid")}
	res = stub.MockInvoke("3", args)
	assert.Equal(t, int32(shim.ERROR), res.Status, "GetChainInfo should have failed because the channel id does not exist")
}

func TestQueryGetTransactionByID(t *testing.T) {
	chainid := "mytestchainid2"
	path := tempDir(t, "test2")
	defer os.RemoveAll(path)

	stub, _, cleanup, err := setupTestLedger(chainid, path)
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer cleanup()

	args := [][]byte{[]byte(GetTransactionByID), []byte(chainid), []byte("1")}
	prop := resetProvider(resources.Qscc_GetTransactionByID, chainid, &peer2.SignedProposal{}, nil)
	res := stub.MockInvokeWithSignedProposal("1", args, prop)
	assert.Equal(t, int32(shim.ERROR), res.Status, "GetTransactionByID should have failed with invalid txid: 1")

	args = [][]byte{[]byte(GetTransactionByID), []byte(chainid), []byte(nil)}
	res = stub.MockInvoke("2", args)
	assert.Equal(t, int32(shim.ERROR), res.Status, "GetTransactionByID should have failed with invalid txid: nil")

	
	args = [][]byte{[]byte(GetTransactionByID), []byte(chainid)}
	res = stub.MockInvoke("3", args)
	assert.Equal(t, int32(shim.ERROR), res.Status, "GetTransactionByID should have failed due to incorrect number of arguments")
}

func TestQueryGetBlockByNumber(t *testing.T) {
	chainid := "mytestchainid3"
	path := tempDir(t, "test3")
	defer os.RemoveAll(path)

	stub, _, cleanup, err := setupTestLedger(chainid, path)
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer cleanup()

	
	args := [][]byte{[]byte(GetBlockByNumber), []byte(chainid), []byte("0")}
	prop := resetProvider(resources.Qscc_GetBlockByNumber, chainid, &peer2.SignedProposal{}, nil)
	res := stub.MockInvokeWithSignedProposal("1", args, prop)
	assert.Equal(t, int32(shim.OK), res.Status, "GetBlockByNumber should have succeeded for block number: 0")

	
	args = [][]byte{[]byte(GetBlockByNumber), []byte(chainid), []byte("1")}
	res = stub.MockInvoke("2", args)
	assert.Equal(t, int32(shim.ERROR), res.Status, "GetBlockByNumber should have failed with invalid number: 1")

	
	args = [][]byte{[]byte(GetBlockByNumber), []byte(chainid), []byte(nil)}
	res = stub.MockInvoke("3", args)
	assert.Equal(t, int32(shim.ERROR), res.Status, "GetBlockByNumber should have failed with nil block number")
}

func TestQueryGetBlockByHash(t *testing.T) {
	chainid := "mytestchainid4"
	path := tempDir(t, "test4")
	defer os.RemoveAll(path)

	stub, _, cleanup, err := setupTestLedger(chainid, path)
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer cleanup()

	args := [][]byte{[]byte(GetBlockByHash), []byte(chainid), []byte("0")}
	prop := resetProvider(resources.Qscc_GetBlockByHash, chainid, &peer2.SignedProposal{}, nil)
	res := stub.MockInvokeWithSignedProposal("1", args, prop)
	assert.Equal(t, int32(shim.ERROR), res.Status, "GetBlockByHash should have failed with invalid hash: 0")

	args = [][]byte{[]byte(GetBlockByHash), []byte(chainid), []byte(nil)}
	res = stub.MockInvoke("2", args)
	assert.Equal(t, int32(shim.ERROR), res.Status, "GetBlockByHash should have failed with nil hash")
}

func TestQueryGetBlockByTxID(t *testing.T) {
	chainid := "mytestchainid5"
	path := tempDir(t, "test5")
	defer os.RemoveAll(path)

	stub, _, cleanup, err := setupTestLedger(chainid, path)
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer cleanup()

	args := [][]byte{[]byte(GetBlockByTxID), []byte(chainid), []byte("")}
	prop := resetProvider(resources.Qscc_GetBlockByTxID, chainid, &peer2.SignedProposal{}, nil)
	res := stub.MockInvokeWithSignedProposal("1", args, prop)
	assert.Equal(t, int32(shim.ERROR), res.Status, "GetBlockByTxID should have failed with blank txId.")
}

func TestFailingAccessControl(t *testing.T) {
	chainid := "mytestchainid6"
	path := tempDir(t, "test6")
	defer os.RemoveAll(path)

	_, p, cleanup, err := setupTestLedger(chainid, path)
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer cleanup()
	e := &LedgerQuerier{
		aclProvider: mockAclProvider,
		ledgers:     p,
	}
	stub := shimtest.NewMockStub("LedgerQuerier", e)

	
	args := [][]byte{[]byte(GetChainInfo), []byte(chainid)}
	sProp, _ := protoutil.MockSignedEndorserProposalOrPanic(chainid, &peer2.ChaincodeSpec{}, []byte("Alice"), []byte("msg1"))
	sProp.Signature = sProp.ProposalBytes
	
	resetProvider(resources.Qscc_GetChainInfo, chainid, sProp, errors.New("Failed access control"))
	res := stub.MockInvokeWithSignedProposal("2", args, sProp)
	assert.Equal(t, int32(shim.ERROR), res.Status, "GetChainInfo must fail: %s", res.Message)
	assert.Contains(t, res.Message, "Failed access control")
	
	mockAclProvider.AssertExpectations(t)

	
	args = [][]byte{[]byte(GetBlockByNumber), []byte(chainid), []byte("1")}
	sProp, _ = protoutil.MockSignedEndorserProposalOrPanic(chainid, &peer2.ChaincodeSpec{}, []byte("Alice"), []byte("msg1"))
	sProp.Signature = sProp.ProposalBytes
	
	resetProvider(resources.Qscc_GetBlockByNumber, chainid, sProp, errors.New("Failed access control"))
	res = stub.MockInvokeWithSignedProposal("2", args, sProp)
	assert.Equal(t, int32(shim.ERROR), res.Status, "GetBlockByNumber must fail: %s", res.Message)
	assert.Contains(t, res.Message, "Failed access control")
	
	mockAclProvider.AssertExpectations(t)

	
	args = [][]byte{[]byte(GetBlockByHash), []byte(chainid), []byte("1")}
	sProp, _ = protoutil.MockSignedEndorserProposalOrPanic(chainid, &peer2.ChaincodeSpec{}, []byte("Alice"), []byte("msg1"))
	sProp.Signature = sProp.ProposalBytes
	
	resetProvider(resources.Qscc_GetBlockByHash, chainid, sProp, errors.New("Failed access control"))
	res = stub.MockInvokeWithSignedProposal("2", args, sProp)
	assert.Equal(t, int32(shim.ERROR), res.Status, "GetBlockByHash must fail: %s", res.Message)
	assert.Contains(t, res.Message, "Failed access control")
	
	mockAclProvider.AssertExpectations(t)

	
	args = [][]byte{[]byte(GetBlockByTxID), []byte(chainid), []byte("1")}
	sProp, _ = protoutil.MockSignedEndorserProposalOrPanic(chainid, &peer2.ChaincodeSpec{}, []byte("Alice"), []byte("msg1"))
	sProp.Signature = sProp.ProposalBytes
	
	resetProvider(resources.Qscc_GetBlockByTxID, chainid, sProp, errors.New("Failed access control"))
	res = stub.MockInvokeWithSignedProposal("2", args, sProp)
	assert.Equal(t, int32(shim.ERROR), res.Status, "GetBlockByTxID must fail: %s", res.Message)
	assert.Contains(t, res.Message, "Failed access control")
	
	mockAclProvider.AssertExpectations(t)

	
	args = [][]byte{[]byte(GetTransactionByID), []byte(chainid), []byte("1")}
	sProp, _ = protoutil.MockSignedEndorserProposalOrPanic(chainid, &peer2.ChaincodeSpec{}, []byte("Alice"), []byte("msg1"))
	sProp.Signature = sProp.ProposalBytes
	
	resetProvider(resources.Qscc_GetTransactionByID, chainid, sProp, errors.New("Failed access control"))
	res = stub.MockInvokeWithSignedProposal("2", args, sProp)
	assert.Equal(t, int32(shim.ERROR), res.Status, "Qscc_GetTransactionByID must fail: %s", res.Message)
	assert.Contains(t, res.Message, "Failed access control")
	
	mockAclProvider.AssertExpectations(t)
}

func TestQueryNonexistentFunction(t *testing.T) {
	chainid := "mytestchainid7"
	path := tempDir(t, "test7")
	defer os.RemoveAll(path)

	stub, _, cleanup, err := setupTestLedger(chainid, path)
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer cleanup()

	args := [][]byte{[]byte("GetBlocks"), []byte(chainid), []byte("arg1")}
	prop := resetProvider("qscc/GetBlocks", chainid, &peer2.SignedProposal{}, nil)
	res := stub.MockInvokeWithSignedProposal("1", args, prop)
	assert.Equal(t, int32(shim.ERROR), res.Status, "GetBlocks should have failed because the function does not exist")
}



func TestQueryGeneratedBlock(t *testing.T) {
	chainid := "mytestchainid8"
	path := tempDir(t, "test8")
	defer os.RemoveAll(path)

	stub, p, cleanup, err := setupTestLedger(chainid, path)
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer cleanup()

	block1 := addBlockForTesting(t, chainid, p)

	
	args := [][]byte{[]byte(GetBlockByNumber), []byte(chainid), []byte("1")}
	prop := resetProvider(resources.Qscc_GetBlockByNumber, chainid, &peer2.SignedProposal{}, nil)
	res := stub.MockInvokeWithSignedProposal("1", args, prop)
	assert.Equal(t, int32(shim.OK), res.Status, "GetBlockByNumber should have succeeded for block number 1")

	
	args = [][]byte{[]byte(GetBlockByHash), []byte(chainid), protoutil.BlockHeaderHash(block1.Header)}
	prop = resetProvider(resources.Qscc_GetBlockByHash, chainid, &peer2.SignedProposal{}, nil)
	res = stub.MockInvokeWithSignedProposal("2", args, prop)
	assert.Equal(t, int32(shim.OK), res.Status, "GetBlockByHash should have succeeded for block 1 hash")

	
	for _, d := range block1.Data.Data {
		ebytes := d
		if ebytes != nil {
			if env, err := protoutil.GetEnvelopeFromBlock(ebytes); err != nil {
				t.Fatalf("error getting envelope from block: %s", err)
			} else if env != nil {
				payload, err := protoutil.UnmarshalPayload(env.Payload)
				if err != nil {
					t.Fatalf("error extracting payload from envelope: %s", err)
				}
				chdr, err := protoutil.UnmarshalChannelHeader(payload.Header.ChannelHeader)
				if err != nil {
					t.Fatalf(err.Error())
				}
				if common.HeaderType(chdr.Type) == common.HeaderType_ENDORSER_TRANSACTION {
					args = [][]byte{[]byte(GetBlockByTxID), []byte(chainid), []byte(chdr.TxId)}
					mockAclProvider.Reset()
					prop = resetProvider(resources.Qscc_GetBlockByTxID, chainid, &peer2.SignedProposal{}, nil)
					res = stub.MockInvokeWithSignedProposal("3", args, prop)
					assert.Equal(t, int32(shim.OK), res.Status, "GetBlockByTxId should have succeeded for txid: %s", chdr.TxId)

					args = [][]byte{[]byte(GetTransactionByID), []byte(chainid), []byte(chdr.TxId)}
					prop = resetProvider(resources.Qscc_GetTransactionByID, chainid, &peer2.SignedProposal{}, nil)
					res = stub.MockInvokeWithSignedProposal("4", args, prop)
					assert.Equal(t, int32(shim.OK), res.Status, "GetTransactionById should have succeeded for txid: %s", chdr.TxId)
				}
			}
		}
	}
}

func addBlockForTesting(t *testing.T, chainid string, p *peer.Peer) *common.Block {
	ledger := p.GetLedger(chainid)
	defer ledger.Close()

	txid1 := util.GenerateUUID()
	simulator, _ := ledger.NewTxSimulator(txid1)
	simulator.SetState("ns1", "key1", []byte("value1"))
	simulator.SetState("ns1", "key2", []byte("value2"))
	simulator.SetState("ns1", "key3", []byte("value3"))
	simulator.Done()
	simRes1, _ := simulator.GetTxSimulationResults()
	pubSimResBytes1, _ := simRes1.GetPubSimulationBytes()

	txid2 := util.GenerateUUID()
	simulator, _ = ledger.NewTxSimulator(txid2)
	simulator.SetState("ns2", "key4", []byte("value4"))
	simulator.SetState("ns2", "key5", []byte("value5"))
	simulator.SetState("ns2", "key6", []byte("value6"))
	simulator.Done()
	simRes2, _ := simulator.GetTxSimulationResults()
	pubSimResBytes2, _ := simRes2.GetPubSimulationBytes()

	bcInfo, err := ledger.GetBlockchainInfo()
	assert.NoError(t, err)
	block1 := testutil.ConstructBlock(t, 1, bcInfo.CurrentBlockHash, [][]byte{pubSimResBytes1, pubSimResBytes2}, false)
	ledger.CommitLegacy(&ledger2.BlockAndPvtData{Block: block1}, &ledger2.CommitOptions{})
	return block1
}

var mockAclProvider *mocks.MockACLProvider

func TestMain(m *testing.M) {
	mockAclProvider = &mocks.MockACLProvider{}
	mockAclProvider.Reset()

	os.Exit(m.Run())
}
