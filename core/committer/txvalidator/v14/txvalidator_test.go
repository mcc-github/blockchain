/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package txvalidator

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/common/configtx/test"
	"github.com/mcc-github/blockchain/common/ledger/testutil"
	"github.com/mcc-github/blockchain/common/mocks/config"
	"github.com/mcc-github/blockchain/common/semaphore"
	util2 "github.com/mcc-github/blockchain/common/util"
	"github.com/mcc-github/blockchain/core/common/sysccprovider"
	ledger2 "github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/ledger/ledgermgmt"
	"github.com/mcc-github/blockchain/core/ledger/ledgermgmt/ledgermgmttest"
	"github.com/mcc-github/blockchain/core/ledger/util"
	ledgerUtil "github.com/mcc-github/blockchain/core/ledger/util"
	mocktxvalidator "github.com/mcc-github/blockchain/core/mocks/txvalidator"
	"github.com/mcc-github/blockchain/core/mocks/validator"
	"github.com/mcc-github/blockchain/msp"
	mspmgmt "github.com/mcc-github/blockchain/msp/mgmt"
	msptesttools "github.com/mcc-github/blockchain/msp/mgmt/testtools"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/peer"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/stretchr/testify/assert"
)

func testValidationWithNTXes(t *testing.T, ledger ledger2.PeerLedger, gbHash []byte, nBlocks int) {
	txid := util2.GenerateUUID()
	simulator, _ := ledger.NewTxSimulator(txid)
	simulator.SetState("ns1", "key1", []byte("value1"))
	simulator.SetState("ns1", "key2", []byte("value2"))
	simulator.SetState("ns1", "key3", []byte("value3"))
	simulator.Done()

	simRes, _ := simulator.GetTxSimulationResults()
	pubSimulationResBytes, _ := simRes.GetPubSimulationBytes()
	_, err := testutil.ConstructBytesProposalResponsePayload("v1", pubSimulationResBytes)
	if err != nil {
		t.Fatalf("Could not construct ProposalResponsePayload bytes, err: %s", err)
	}

	mockVsccValidator := &validator.MockVsccValidator{}
	tValidator := &TxValidator{
		ChainID:          "",
		Semaphore:        semaphore.New(10),
		ChannelResources: &mocktxvalidator.Support{LedgerVal: ledger, ACVal: &config.MockApplicationCapabilities{}},
		Vscc:             mockVsccValidator,
	}

	bcInfo, _ := ledger.GetBlockchainInfo()
	assert.Equal(t, &common.BlockchainInfo{
		Height: 1, CurrentBlockHash: gbHash, PreviousBlockHash: nil,
	}, bcInfo)

	sr := [][]byte{}
	for i := 0; i < nBlocks; i++ {
		sr = append(sr, pubSimulationResBytes)
	}
	block := testutil.ConstructBlock(t, 1, gbHash, sr, true)

	tValidator.Validate(block)

	txsfltr := util.TxValidationFlags(block.Metadata.Metadata[common.BlockMetadataIndex_TRANSACTIONS_FILTER])

	for i := 0; i < nBlocks; i++ {
		assert.True(t, txsfltr.IsSetTo(i, peer.TxValidationCode_VALID))
	}
}

func TestDetectTXIdDuplicates(t *testing.T) {
	txids := []string{"", "1", "2", "3", "", "2", ""}
	txsfltr := ledgerUtil.NewTxValidationFlags(len(txids))
	markTXIdDuplicates(txids, txsfltr)
	assert.True(t, txsfltr.IsSetTo(0, peer.TxValidationCode_NOT_VALIDATED))
	assert.True(t, txsfltr.IsSetTo(1, peer.TxValidationCode_NOT_VALIDATED))
	assert.True(t, txsfltr.IsSetTo(2, peer.TxValidationCode_NOT_VALIDATED))
	assert.True(t, txsfltr.IsSetTo(3, peer.TxValidationCode_NOT_VALIDATED))
	assert.True(t, txsfltr.IsSetTo(4, peer.TxValidationCode_NOT_VALIDATED))
	assert.True(t, txsfltr.IsSetTo(5, peer.TxValidationCode_DUPLICATE_TXID))
	assert.True(t, txsfltr.IsSetTo(6, peer.TxValidationCode_NOT_VALIDATED))

	txids = []string{"", "1", "2", "3", "", "21", ""}
	txsfltr = ledgerUtil.NewTxValidationFlags(len(txids))
	markTXIdDuplicates(txids, txsfltr)
	assert.True(t, txsfltr.IsSetTo(0, peer.TxValidationCode_NOT_VALIDATED))
	assert.True(t, txsfltr.IsSetTo(1, peer.TxValidationCode_NOT_VALIDATED))
	assert.True(t, txsfltr.IsSetTo(2, peer.TxValidationCode_NOT_VALIDATED))
	assert.True(t, txsfltr.IsSetTo(3, peer.TxValidationCode_NOT_VALIDATED))
	assert.True(t, txsfltr.IsSetTo(4, peer.TxValidationCode_NOT_VALIDATED))
	assert.True(t, txsfltr.IsSetTo(5, peer.TxValidationCode_NOT_VALIDATED))
	assert.True(t, txsfltr.IsSetTo(6, peer.TxValidationCode_NOT_VALIDATED))
}

func TestBlockValidationDuplicateTXId(t *testing.T) {
	ledgerMgr, cleanup := constructLedgerMgrWithTestDefaults(t, "txvalidator")
	defer cleanup()

	gb, _ := test.MakeGenesisBlock("TestLedger")
	gbHash := protoutil.BlockHeaderHash(gb.Header)
	ledger, _ := ledgerMgr.CreateLedger(gb)
	defer ledger.Close()

	txid := util2.GenerateUUID()
	simulator, _ := ledger.NewTxSimulator(txid)
	simulator.SetState("ns1", "key1", []byte("value1"))
	simulator.SetState("ns1", "key2", []byte("value2"))
	simulator.SetState("ns1", "key3", []byte("value3"))
	simulator.Done()

	simRes, _ := simulator.GetTxSimulationResults()
	pubSimulationResBytes, _ := simRes.GetPubSimulationBytes()
	_, err := testutil.ConstructBytesProposalResponsePayload("v1", pubSimulationResBytes)
	if err != nil {
		t.Fatalf("Could not construct ProposalResponsePayload bytes, err: %s", err)
	}

	mockVsccValidator := &validator.MockVsccValidator{}
	acv := &config.MockApplicationCapabilities{}
	tValidator := &TxValidator{
		ChainID:          "",
		Semaphore:        semaphore.New(10),
		ChannelResources: &mocktxvalidator.Support{LedgerVal: ledger, ACVal: acv},
		Vscc:             mockVsccValidator,
	}

	bcInfo, _ := ledger.GetBlockchainInfo()
	assert.Equal(t, &common.BlockchainInfo{
		Height: 1, CurrentBlockHash: gbHash, PreviousBlockHash: nil,
	}, bcInfo)

	envs := []*common.Envelope{}
	env, _, err := testutil.ConstructTransaction(t, pubSimulationResBytes, "", true)
	envs = append(envs, env)
	envs = append(envs, env)
	block := testutil.NewBlock(envs, 1, gbHash)

	tValidator.Validate(block)

	txsfltr := util.TxValidationFlags(block.Metadata.Metadata[common.BlockMetadataIndex_TRANSACTIONS_FILTER])

	assert.True(t, txsfltr.IsSetTo(0, peer.TxValidationCode_VALID))
	assert.True(t, txsfltr.IsSetTo(1, peer.TxValidationCode_VALID))

	acv.ForbidDuplicateTXIdInBlockRv = true

	tValidator.Validate(block)

	txsfltr = util.TxValidationFlags(block.Metadata.Metadata[common.BlockMetadataIndex_TRANSACTIONS_FILTER])

	assert.True(t, txsfltr.IsSetTo(0, peer.TxValidationCode_VALID))
	assert.True(t, txsfltr.IsSetTo(1, peer.TxValidationCode_DUPLICATE_TXID))
}

func TestBlockValidation(t *testing.T) {
	ledgerMgr, cleanup := constructLedgerMgrWithTestDefaults(t, "txvalidator")
	defer cleanup()

	gb, _ := test.MakeGenesisBlock("TestLedger")
	gbHash := protoutil.BlockHeaderHash(gb.Header)
	ledger, _ := ledgerMgr.CreateLedger(gb)
	defer ledger.Close()

	
	testValidationWithNTXes(t, ledger, gbHash, 1)
}

func TestParallelBlockValidation(t *testing.T) {
	ledgerMgr, cleanup := constructLedgerMgrWithTestDefaults(t, "txvalidator")
	defer cleanup()

	gb, _ := test.MakeGenesisBlock("TestLedger")
	gbHash := protoutil.BlockHeaderHash(gb.Header)
	ledger, _ := ledgerMgr.CreateLedger(gb)
	defer ledger.Close()

	
	testValidationWithNTXes(t, ledger, gbHash, 128)
}

func TestVeryLargeParallelBlockValidation(t *testing.T) {
	ledgerMgr, cleanup := constructLedgerMgrWithTestDefaults(t, "txvalidator")
	defer cleanup()

	gb, _ := test.MakeGenesisBlock("TestLedger")
	gbHash := protoutil.BlockHeaderHash(gb.Header)
	ledger, _ := ledgerMgr.CreateLedger(gb)
	defer ledger.Close()

	
	
	
	testValidationWithNTXes(t, ledger, gbHash, 4096)
}

func TestTxValidationFailure_InvalidTxid(t *testing.T) {
	ledgerMgr, cleanup := constructLedgerMgrWithTestDefaults(t, "txvalidator")
	defer cleanup()

	gb, _ := test.MakeGenesisBlock("TestLedger")
	ledger, _ := ledgerMgr.CreateLedger(gb)

	defer ledger.Close()

	tValidator := &TxValidator{
		ChainID:          "",
		Semaphore:        semaphore.New(10),
		ChannelResources: &mocktxvalidator.Support{LedgerVal: ledger, ACVal: &config.MockApplicationCapabilities{}},
		Vscc:             &validator.MockVsccValidator{},
	}

	mockSigner, err := mspmgmt.GetLocalMSP().GetDefaultSigningIdentity()
	assert.NoError(t, err)
	mockSignerSerialized, err := mockSigner.Serialize()
	assert.NoError(t, err)

	
	payload := &common.Payload{
		Header: &common.Header{
			ChannelHeader: protoutil.MarshalOrPanic(&common.ChannelHeader{
				TxId:      "INVALID TXID!!!",
				Type:      int32(common.HeaderType_ENDORSER_TRANSACTION),
				ChannelId: util2.GetTestChainID(),
			}),
			SignatureHeader: protoutil.MarshalOrPanic(&common.SignatureHeader{
				Nonce:   []byte("nonce"),
				Creator: mockSignerSerialized,
			}),
		},
		Data: []byte("test"),
	}

	payloadBytes, err := proto.Marshal(payload)

	
	assert.NoError(t, err)

	sig, err := mockSigner.Sign(payloadBytes)
	assert.NoError(t, err)

	
	envelope := &common.Envelope{
		Payload:   payloadBytes,
		Signature: sig,
	}

	envelopeBytes, err := proto.Marshal(envelope)

	
	assert.NoError(t, err)

	block := &common.Block{
		Data: &common.BlockData{
			
			Data: [][]byte{envelopeBytes},
		},
	}

	block.Header = &common.BlockHeader{
		Number:   0,
		DataHash: protoutil.BlockDataHash(block.Data),
	}

	
	protoutil.InitBlockMetadata(block)
	txsFilter := util.NewTxValidationFlagsSetValue(len(block.Data.Data), peer.TxValidationCode_VALID)
	block.Metadata.Metadata[common.BlockMetadataIndex_TRANSACTIONS_FILTER] = txsFilter

	
	ledger.CommitLegacy(&ledger2.BlockAndPvtData{Block: block}, &ledger2.CommitOptions{})

	
	
	tValidator.Validate(block)

	txsfltr := util.TxValidationFlags(block.Metadata.Metadata[common.BlockMetadataIndex_TRANSACTIONS_FILTER])
	assert.True(t, txsfltr.IsInvalid(0))

	
	assert.True(t, txsfltr.Flag(0) == peer.TxValidationCode_BAD_PROPOSAL_TXID)
}

func createCCUpgradeEnvelope(chainID, chaincodeName, chaincodeVersion string, signer msp.SigningIdentity) (*common.Envelope, error) {
	creator, err := signer.Serialize()
	if err != nil {
		return nil, err
	}

	spec := &peer.ChaincodeSpec{
		Type: peer.ChaincodeSpec_Type(peer.ChaincodeSpec_Type_value["GOLANG"]),
		ChaincodeId: &peer.ChaincodeID{
			Path:    "github.com/codePath",
			Name:    chaincodeName,
			Version: chaincodeVersion,
		},
	}

	cds := &peer.ChaincodeDeploymentSpec{ChaincodeSpec: spec, CodePackage: []byte{}}
	prop, _, err := protoutil.CreateUpgradeProposalFromCDS(chainID, cds, creator, []byte{}, []byte{}, []byte{}, nil)
	if err != nil {
		return nil, err
	}

	proposalResponse := &peer.ProposalResponse{
		Response: &peer.Response{
			Status: 200, 
		},
		Endorsement: &peer.Endorsement{},
	}

	return protoutil.CreateSignedTx(prop, signer, proposalResponse)
}

func TestGetTxCCInstance(t *testing.T) {
	
	err := msptesttools.LoadMSPSetupForTesting()
	if err != nil {
		t.Fatalf("Could not initialize msp, err: %s", err)
	}
	signer, err := mspmgmt.GetLocalMSP().GetDefaultSigningIdentity()
	if err != nil {
		t.Fatalf("Could not initialize signer, err: %s", err)
	}

	chainID := util2.GetTestChainID()
	upgradeCCName := "mycc"
	upgradeCCVersion := "v1"

	env, err := createCCUpgradeEnvelope(chainID, upgradeCCName, upgradeCCVersion, signer)
	assert.NoError(t, err)

	
	payload, err := protoutil.UnmarshalPayload(env.Payload)
	assert.NoError(t, err)

	expectInvokeCCIns := &sysccprovider.ChaincodeInstance{
		ChainID:          chainID,
		ChaincodeName:    "lscc",
		ChaincodeVersion: "",
	}
	expectUpgradeCCIns := &sysccprovider.ChaincodeInstance{
		ChainID:          chainID,
		ChaincodeName:    upgradeCCName,
		ChaincodeVersion: upgradeCCVersion,
	}

	tValidator := &TxValidator{}
	invokeCCIns, upgradeCCIns, err := tValidator.getTxCCInstance(payload)
	if err != nil {
		t.Fatalf("Get chaincode from tx error: %s", err)
	}
	assert.EqualValues(t, expectInvokeCCIns, invokeCCIns)
	assert.EqualValues(t, expectUpgradeCCIns, upgradeCCIns)
}

func TestInvalidTXsForUpgradeCC(t *testing.T) {
	txsChaincodeNames := map[int]*sysccprovider.ChaincodeInstance{
		0: {ChainID: "chain0", ChaincodeName: "cc0", ChaincodeVersion: "v0"}, 
		1: {ChainID: "chain1", ChaincodeName: "cc0", ChaincodeVersion: "v0"}, 
		2: {ChainID: "chain1", ChaincodeName: "lscc", ChaincodeVersion: ""},  
		3: {ChainID: "chain1", ChaincodeName: "cc0", ChaincodeVersion: "v0"}, 
		4: {ChainID: "chain1", ChaincodeName: "cc0", ChaincodeVersion: "v1"}, 
		5: {ChainID: "chain1", ChaincodeName: "cc1", ChaincodeVersion: "v0"}, 
		6: {ChainID: "chain1", ChaincodeName: "lscc", ChaincodeVersion: ""},  
		7: {ChainID: "chain1", ChaincodeName: "lscc", ChaincodeVersion: ""},  
	}
	upgradedChaincodes := map[int]*sysccprovider.ChaincodeInstance{
		2: {ChainID: "chain1", ChaincodeName: "cc0", ChaincodeVersion: "v1"},
		6: {ChainID: "chain1", ChaincodeName: "cc0", ChaincodeVersion: "v2"},
		7: {ChainID: "chain1", ChaincodeName: "cc0", ChaincodeVersion: "v3"},
	}

	txsfltr := ledgerUtil.NewTxValidationFlags(8)
	txsfltr.SetFlag(0, peer.TxValidationCode_VALID)
	txsfltr.SetFlag(1, peer.TxValidationCode_VALID)
	txsfltr.SetFlag(2, peer.TxValidationCode_VALID)
	txsfltr.SetFlag(3, peer.TxValidationCode_VALID)
	txsfltr.SetFlag(4, peer.TxValidationCode_VALID)
	txsfltr.SetFlag(5, peer.TxValidationCode_VALID)
	txsfltr.SetFlag(6, peer.TxValidationCode_VALID)
	txsfltr.SetFlag(7, peer.TxValidationCode_VALID)

	expectTxsFltr := ledgerUtil.NewTxValidationFlags(8)
	expectTxsFltr.SetFlag(0, peer.TxValidationCode_VALID)
	expectTxsFltr.SetFlag(1, peer.TxValidationCode_CHAINCODE_VERSION_CONFLICT)
	expectTxsFltr.SetFlag(2, peer.TxValidationCode_CHAINCODE_VERSION_CONFLICT)
	expectTxsFltr.SetFlag(3, peer.TxValidationCode_CHAINCODE_VERSION_CONFLICT)
	expectTxsFltr.SetFlag(4, peer.TxValidationCode_CHAINCODE_VERSION_CONFLICT)
	expectTxsFltr.SetFlag(5, peer.TxValidationCode_VALID)
	expectTxsFltr.SetFlag(6, peer.TxValidationCode_CHAINCODE_VERSION_CONFLICT)
	expectTxsFltr.SetFlag(7, peer.TxValidationCode_VALID)

	tValidator := &TxValidator{}
	tValidator.invalidTXsForUpgradeCC(txsChaincodeNames, upgradedChaincodes, txsfltr)

	assert.EqualValues(t, expectTxsFltr, txsfltr)
}

func constructLedgerMgrWithTestDefaults(t *testing.T, testDir string) (*ledgermgmt.LedgerMgr, func()) {
	testDir, err := ioutil.TempDir("", testDir)
	if err != nil {
		t.Fatalf("Failed to create ledger directory: %s", err)
	}
	ledgerMgr := ledgermgmt.NewLedgerMgr(ledgermgmttest.NewInitializer(testDir))
	cleanup := func() {
		ledgerMgr.Close()
		os.RemoveAll(testDir)
	}
	return ledgerMgr, cleanup
}
