/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package txvalidator_test

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/mcc-github/blockchain/common/cauthdsl"
	ctxt "github.com/mcc-github/blockchain/common/configtx/test"
	commonerrors "github.com/mcc-github/blockchain/common/errors"
	ledger2 "github.com/mcc-github/blockchain/common/ledger"
	"github.com/mcc-github/blockchain/common/ledger/testutil"
	mockconfig "github.com/mcc-github/blockchain/common/mocks/config"
	"github.com/mcc-github/blockchain/common/mocks/scc"
	"github.com/mcc-github/blockchain/common/util"
	"github.com/mcc-github/blockchain/core/committer/txvalidator"
	"github.com/mcc-github/blockchain/core/committer/txvalidator/mocks"
	"github.com/mcc-github/blockchain/core/committer/txvalidator/testdata"
	ccp "github.com/mcc-github/blockchain/core/common/ccprovider"
	"github.com/mcc-github/blockchain/core/handlers/validation/api"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/rwsetutil"
	"github.com/mcc-github/blockchain/core/ledger/ledgermgmt"
	lutils "github.com/mcc-github/blockchain/core/ledger/util"
	mocktxvalidator "github.com/mcc-github/blockchain/core/mocks/txvalidator"
	mocks2 "github.com/mcc-github/blockchain/discovery/support/mocks"
	"github.com/mcc-github/blockchain/msp"
	"github.com/mcc-github/blockchain/msp/mgmt"
	"github.com/mcc-github/blockchain/msp/mgmt/testtools"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/peer"
	"github.com/mcc-github/blockchain/protos/utils"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/sync/semaphore"
)

func signedByAnyMember(ids []string) []byte {
	p := cauthdsl.SignedByAnyMember(ids)
	return utils.MarshalOrPanic(p)
}

func setupLedgerAndValidator(t *testing.T) (ledger.PeerLedger, txvalidator.Validator) {
	plugin := &mocks.Plugin{}
	plugin.On("Init", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	plugin.On("Validate", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	return setupLedgerAndValidatorExplicit(t, &mockconfig.MockApplicationCapabilities{}, plugin)
}

func setupLedgerAndValidatorExplicit(t *testing.T, cpb *mockconfig.MockApplicationCapabilities, plugin *mocks.Plugin) (ledger.PeerLedger, txvalidator.Validator) {
	viper.Set("peer.fileSystemPath", "/tmp/blockchain/validatortest")
	ledgermgmt.InitializeTestEnv()
	gb, err := ctxt.MakeGenesisBlock("TestLedger")
	assert.NoError(t, err)
	theLedger, err := ledgermgmt.CreateLedger(gb)
	assert.NoError(t, err)
	vcs := struct {
		*mocktxvalidator.Support
		*semaphore.Weighted
	}{&mocktxvalidator.Support{LedgerVal: theLedger, ACVal: cpb}, semaphore.NewWeighted(10)}
	mp := (&scc.MocksccProviderFactory{}).NewSystemChaincodeProvider()
	pm := &mocks.PluginMapper{}
	factory := &mocks.PluginFactory{}
	pm.On("PluginFactoryByName", txvalidator.PluginName("vscc")).Return(factory)
	factory.On("New").Return(plugin)

	theValidator := txvalidator.NewTxValidator(vcs, mp, pm)

	return theLedger, theValidator
}

func createRWset(t *testing.T, ccnames ...string) []byte {
	rwsetBuilder := rwsetutil.NewRWSetBuilder()
	for _, ccname := range ccnames {
		rwsetBuilder.AddToWriteSet(ccname, "key", []byte("value"))
	}
	rwset, err := rwsetBuilder.GetTxSimulationResults()
	assert.NoError(t, err)
	rwsetBytes, err := rwset.GetPubSimulationBytes()
	return rwsetBytes
}

func getProposalWithType(ccID string, pType common.HeaderType) (*peer.Proposal, error) {
	cis := &peer.ChaincodeInvocationSpec{
		ChaincodeSpec: &peer.ChaincodeSpec{
			ChaincodeId: &peer.ChaincodeID{Name: ccID, Version: ccVersion},
			Input:       &peer.ChaincodeInput{Args: [][]byte{[]byte("func")}},
			Type:        peer.ChaincodeSpec_GOLANG}}

	proposal, _, err := utils.CreateProposalFromCIS(pType, util.GetTestChainID(), cis, signerSerialized)
	return proposal, err
}

const ccVersion = "1.0"

func getEnvWithType(ccID string, event []byte, res []byte, pType common.HeaderType, t *testing.T) *common.Envelope {
	
	prop, err := getProposalWithType(ccID, pType)
	assert.NoError(t, err)

	response := &peer.Response{Status: 200}

	
	presp, err := utils.CreateProposalResponse(prop.Header, prop.Payload, response, res, event, &peer.ChaincodeID{Name: ccID, Version: ccVersion}, nil, signer)
	assert.NoError(t, err)

	
	tx, err := utils.CreateSignedTx(prop, signer, presp)
	assert.NoError(t, err)

	return tx
}

func getEnv(ccID string, event []byte, res []byte, t *testing.T) *common.Envelope {
	return getEnvWithType(ccID, event, res, common.HeaderType_ENDORSER_TRANSACTION, t)
}

func putCCInfoWithVSCCAndVer(theLedger ledger.PeerLedger, ccname, vscc, ver string, policy []byte, t *testing.T) {
	cd := &ccp.ChaincodeData{
		Name:    ccname,
		Version: ver,
		Vscc:    vscc,
		Policy:  policy,
	}

	cdbytes := utils.MarshalOrPanic(cd)

	txid := util.GenerateUUID()
	simulator, err := theLedger.NewTxSimulator(txid)
	assert.NoError(t, err)
	simulator.SetState("lscc", ccname, cdbytes)
	simulator.Done()

	simRes, err := simulator.GetTxSimulationResults()
	assert.NoError(t, err)
	pubSimulationBytes, err := simRes.GetPubSimulationBytes()
	assert.NoError(t, err)
	block0 := testutil.ConstructBlock(t, 1, []byte("hash"), [][]byte{pubSimulationBytes}, true)
	err = theLedger.CommitWithPvtData(&ledger.BlockAndPvtData{
		Block: block0,
	})
	assert.NoError(t, err)
}

func putCCInfo(theLedger ledger.PeerLedger, ccname string, policy []byte, t *testing.T) {
	putCCInfoWithVSCCAndVer(theLedger, ccname, "vscc", ccVersion, policy, t)
}

func assertInvalid(block *common.Block, t *testing.T, code peer.TxValidationCode) {
	txsFilter := lutils.TxValidationFlags(block.Metadata.Metadata[common.BlockMetadataIndex_TRANSACTIONS_FILTER])
	assert.True(t, txsFilter.IsInvalid(0))
	assert.True(t, txsFilter.IsSetTo(0, code))
}

func assertValid(block *common.Block, t *testing.T) {
	txsFilter := lutils.TxValidationFlags(block.Metadata.Metadata[common.BlockMetadataIndex_TRANSACTIONS_FILTER])
	assert.False(t, txsFilter.IsInvalid(0))
}

func TestInvokeBadRWSet(t *testing.T) {
	l, v := setupLedgerAndValidator(t)
	defer ledgermgmt.CleanupTestEnv()
	defer l.Close()

	ccID := "mycc"

	tx := getEnv(ccID, nil, []byte("barf"), t)
	b := &common.Block{Data: &common.BlockData{Data: [][]byte{utils.MarshalOrPanic(tx)}}}

	err := v.Validate(b)
	assert.NoError(t, err)
	assertInvalid(b, t, peer.TxValidationCode_BAD_RWSET)
}

func TestInvokeNoPolicy(t *testing.T) {
	l, v := setupLedgerAndValidator(t)
	defer ledgermgmt.CleanupTestEnv()
	defer l.Close()

	ccID := "mycc"

	putCCInfo(l, ccID, nil, t)

	tx := getEnv(ccID, nil, createRWset(t, ccID), t)
	b := &common.Block{Data: &common.BlockData{Data: [][]byte{utils.MarshalOrPanic(tx)}}}

	err := v.Validate(b)
	assert.NoError(t, err)
	assertInvalid(b, t, peer.TxValidationCode_INVALID_OTHER_REASON)
}

func TestInvokeOK(t *testing.T) {
	l, v := setupLedgerAndValidator(t)
	defer ledgermgmt.CleanupTestEnv()
	defer l.Close()

	ccID := "mycc"

	putCCInfo(l, ccID, signedByAnyMember([]string{"SampleOrg"}), t)

	tx := getEnv(ccID, nil, createRWset(t, ccID), t)
	b := &common.Block{Data: &common.BlockData{Data: [][]byte{utils.MarshalOrPanic(tx)}}}

	err := v.Validate(b)
	assert.NoError(t, err)
	assertValid(b, t)
}

func TestInvokeNoRWSet(t *testing.T) {
	plugin := &mocks.Plugin{}
	plugin.On("Init", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	t.Run("Pre-1.2Capability", func(t *testing.T) {
		l, v := setupLedgerAndValidatorExplicit(t, &mockconfig.MockApplicationCapabilities{}, plugin)
		defer ledgermgmt.CleanupTestEnv()
		defer l.Close()

		ccID := "mycc"

		putCCInfo(l, ccID, signedByAnyMember([]string{"SampleOrg"}), t)

		tx := getEnv(ccID, nil, createRWset(t), t)
		b := &common.Block{Data: &common.BlockData{Data: [][]byte{utils.MarshalOrPanic(tx)}}}

		err := v.Validate(b)
		assert.NoError(t, err)
		assertValid(b, t)
	})

	
	
	plugin.On("Validate", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("tx is invalid")).Once()
	t.Run("Post-1.2Capability", func(t *testing.T) {
		l, v := setupLedgerAndValidatorExplicit(t, &mockconfig.MockApplicationCapabilities{V1_2ValidationRv: true}, plugin)
		defer ledgermgmt.CleanupTestEnv()
		defer l.Close()

		ccID := "mycc"

		putCCInfo(l, ccID, signedByAnyMember([]string{"SampleOrg"}), t)

		tx := getEnv(ccID, nil, createRWset(t), t)
		b := &common.Block{Data: &common.BlockData{Data: [][]byte{utils.MarshalOrPanic(tx)}}}

		err := v.Validate(b)
		assert.NoError(t, err)
		assertInvalid(b, t, peer.TxValidationCode_ENDORSEMENT_POLICY_FAILURE)
	})
}

func TestChaincodeEvent(t *testing.T) {
	plugin := &mocks.Plugin{}
	plugin.On("Init", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	plugin.On("Validate", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	t.Run("PreV1.2", func(t *testing.T) {
		t.Run("MisMatchedName", func(t *testing.T) {
			l, v := setupLedgerAndValidatorExplicit(t, &mockconfig.MockApplicationCapabilities{V1_2ValidationRv: false}, plugin)
			defer ledgermgmt.CleanupTestEnv()
			defer l.Close()

			ccID := "mycc"

			putCCInfo(l, ccID, signedByAnyMember([]string{"SampleOrg"}), t)

			tx := getEnv(ccID, utils.MarshalOrPanic(&peer.ChaincodeEvent{ChaincodeId: "wrong"}), createRWset(t), t)
			b := &common.Block{Data: &common.BlockData{Data: [][]byte{utils.MarshalOrPanic(tx)}}}

			err := v.Validate(b)
			assert.NoError(t, err)
			assertValid(b, t)
		})

		t.Run("BadBytes", func(t *testing.T) {
			l, v := setupLedgerAndValidatorExplicit(t, &mockconfig.MockApplicationCapabilities{V1_2ValidationRv: false}, plugin)
			defer ledgermgmt.CleanupTestEnv()
			defer l.Close()

			ccID := "mycc"

			putCCInfo(l, ccID, signedByAnyMember([]string{"SampleOrg"}), t)

			tx := getEnv(ccID, []byte("garbage"), createRWset(t), t)
			b := &common.Block{Data: &common.BlockData{Data: [][]byte{utils.MarshalOrPanic(tx)}}}

			err := v.Validate(b)
			assert.NoError(t, err)
			assertValid(b, t)
		})

		t.Run("GoodPath", func(t *testing.T) {
			l, v := setupLedgerAndValidatorExplicit(t, &mockconfig.MockApplicationCapabilities{V1_2ValidationRv: false}, plugin)
			defer ledgermgmt.CleanupTestEnv()
			defer l.Close()

			ccID := "mycc"

			putCCInfo(l, ccID, signedByAnyMember([]string{"SampleOrg"}), t)

			tx := getEnv(ccID, utils.MarshalOrPanic(&peer.ChaincodeEvent{ChaincodeId: ccID}), createRWset(t), t)
			b := &common.Block{Data: &common.BlockData{Data: [][]byte{utils.MarshalOrPanic(tx)}}}

			err := v.Validate(b)
			assert.NoError(t, err)
			assertValid(b, t)
		})
	})

	t.Run("PostV1.2", func(t *testing.T) {
		t.Run("MisMatchedName", func(t *testing.T) {
			l, v := setupLedgerAndValidatorExplicit(t, &mockconfig.MockApplicationCapabilities{V1_2ValidationRv: true}, plugin)
			defer ledgermgmt.CleanupTestEnv()
			defer l.Close()

			ccID := "mycc"

			putCCInfo(l, ccID, signedByAnyMember([]string{"SampleOrg"}), t)

			tx := getEnv(ccID, utils.MarshalOrPanic(&peer.ChaincodeEvent{ChaincodeId: "wrong"}), createRWset(t), t)
			b := &common.Block{Data: &common.BlockData{Data: [][]byte{utils.MarshalOrPanic(tx)}}}

			err := v.Validate(b)
			assert.NoError(t, err) 
			assertInvalid(b, t, peer.TxValidationCode_INVALID_OTHER_REASON)
		})

		t.Run("BadBytes", func(t *testing.T) {
			l, v := setupLedgerAndValidatorExplicit(t, &mockconfig.MockApplicationCapabilities{V1_2ValidationRv: true}, plugin)
			defer ledgermgmt.CleanupTestEnv()
			defer l.Close()

			ccID := "mycc"

			putCCInfo(l, ccID, signedByAnyMember([]string{"SampleOrg"}), t)

			tx := getEnv(ccID, []byte("garbage"), createRWset(t), t)
			b := &common.Block{Data: &common.BlockData{Data: [][]byte{utils.MarshalOrPanic(tx)}}}

			err := v.Validate(b)
			assert.NoError(t, err) 
			assertInvalid(b, t, peer.TxValidationCode_INVALID_OTHER_REASON)
		})

		t.Run("GoodPath", func(t *testing.T) {
			l, v := setupLedgerAndValidatorExplicit(t, &mockconfig.MockApplicationCapabilities{V1_2ValidationRv: true}, plugin)
			defer ledgermgmt.CleanupTestEnv()
			defer l.Close()

			ccID := "mycc"

			putCCInfo(l, ccID, signedByAnyMember([]string{"SampleOrg"}), t)

			tx := getEnv(ccID, utils.MarshalOrPanic(&peer.ChaincodeEvent{ChaincodeId: ccID}), createRWset(t), t)
			b := &common.Block{Data: &common.BlockData{Data: [][]byte{utils.MarshalOrPanic(tx)}}}

			err := v.Validate(b)
			assert.NoError(t, err)
			assertValid(b, t)
		})
	})
}

func TestInvokeOKPvtDataOnly(t *testing.T) {
	plugin := &mocks.Plugin{}
	plugin.On("Init", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	l, v := setupLedgerAndValidatorExplicit(t, &mockconfig.MockApplicationCapabilities{}, plugin)
	defer ledgermgmt.CleanupTestEnv()
	defer l.Close()

	v.(*txvalidator.TxValidator).Support.(struct {
		*mocktxvalidator.Support
		*semaphore.Weighted
	}).ACVal = &mockconfig.MockApplicationCapabilities{PrivateChannelDataRv: true}

	ccID := "mycc"

	putCCInfo(l, ccID, signedByAnyMember([]string{"SampleOrg"}), t)

	rwsetBuilder := rwsetutil.NewRWSetBuilder()
	rwsetBuilder.AddToPvtAndHashedWriteSet(ccID, "mycollection", "somekey", nil)
	rwset, err := rwsetBuilder.GetTxSimulationResults()
	assert.NoError(t, err)
	rwsetBytes, err := rwset.GetPubSimulationBytes()
	assert.NoError(t, err)

	tx := getEnv(ccID, nil, rwsetBytes, t)
	b := &common.Block{Data: &common.BlockData{Data: [][]byte{utils.MarshalOrPanic(tx)}}}

	plugin.On("Validate", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("tx is invalid"))

	err = v.Validate(b)
	assert.NoError(t, err)
	assertInvalid(b, t, peer.TxValidationCode_ENDORSEMENT_POLICY_FAILURE)
}

func TestInvokeOKSCC(t *testing.T) {
	l, v := setupLedgerAndValidator(t)
	defer ledgermgmt.CleanupTestEnv()
	defer l.Close()

	ccID := "lscc"

	putCCInfo(l, ccID, signedByAnyMember([]string{"SampleOrg"}), t)

	tx := getEnv(ccID, nil, createRWset(t, ccID), t)
	b := &common.Block{Data: &common.BlockData{Data: [][]byte{utils.MarshalOrPanic(tx)}}}

	err := v.Validate(b)
	assert.NoError(t, err)
	assertValid(b, t)
}

func TestInvokeNOKWritesToLSCC(t *testing.T) {
	l, v := setupLedgerAndValidator(t)
	defer ledgermgmt.CleanupTestEnv()
	defer l.Close()

	ccID := "mycc"

	putCCInfo(l, ccID, signedByAnyMember([]string{"SampleOrg"}), t)

	tx := getEnv(ccID, nil, createRWset(t, ccID, "lscc"), t)
	b := &common.Block{Data: &common.BlockData{Data: [][]byte{utils.MarshalOrPanic(tx)}}}

	err := v.Validate(b)
	assert.NoError(t, err)
	assertInvalid(b, t, peer.TxValidationCode_ILLEGAL_WRITESET)
}

func TestInvokeNOKWritesToESCC(t *testing.T) {
	l, v := setupLedgerAndValidator(t)
	defer ledgermgmt.CleanupTestEnv()
	defer l.Close()

	ccID := "mycc"

	putCCInfo(l, ccID, signedByAnyMember([]string{"SampleOrg"}), t)

	tx := getEnv(ccID, nil, createRWset(t, ccID, "escc"), t)
	b := &common.Block{Data: &common.BlockData{Data: [][]byte{utils.MarshalOrPanic(tx)}}}

	err := v.Validate(b)
	assert.NoError(t, err)
	assertInvalid(b, t, peer.TxValidationCode_ILLEGAL_WRITESET)
}

func TestInvokeNOKWritesToNotExt(t *testing.T) {
	l, v := setupLedgerAndValidator(t)
	defer ledgermgmt.CleanupTestEnv()
	defer l.Close()

	ccID := "mycc"

	putCCInfo(l, ccID, signedByAnyMember([]string{"SampleOrg"}), t)

	tx := getEnv(ccID, nil, createRWset(t, ccID, "notext"), t)
	b := &common.Block{Data: &common.BlockData{Data: [][]byte{utils.MarshalOrPanic(tx)}}}

	err := v.Validate(b)
	assert.NoError(t, err)
	assertInvalid(b, t, peer.TxValidationCode_ILLEGAL_WRITESET)
}

func TestInvokeNOKInvokesNotExt(t *testing.T) {
	l, v := setupLedgerAndValidator(t)
	defer ledgermgmt.CleanupTestEnv()
	defer l.Close()

	ccID := "notext"

	putCCInfo(l, ccID, signedByAnyMember([]string{"SampleOrg"}), t)

	tx := getEnv(ccID, nil, createRWset(t, ccID), t)
	b := &common.Block{Data: &common.BlockData{Data: [][]byte{utils.MarshalOrPanic(tx)}}}

	err := v.Validate(b)
	assert.NoError(t, err)
	assertInvalid(b, t, peer.TxValidationCode_ILLEGAL_WRITESET)
}

func TestInvokeNOKInvokesEmptyCCName(t *testing.T) {
	l, v := setupLedgerAndValidator(t)
	defer ledgermgmt.CleanupTestEnv()
	defer l.Close()

	ccID := ""

	putCCInfo(l, ccID, signedByAnyMember([]string{"SampleOrg"}), t)

	tx := getEnv(ccID, nil, createRWset(t, ccID), t)
	b := &common.Block{Data: &common.BlockData{Data: [][]byte{utils.MarshalOrPanic(tx)}}}

	err := v.Validate(b)
	assert.NoError(t, err)
	assertInvalid(b, t, peer.TxValidationCode_INVALID_OTHER_REASON)
}

func TestInvokeNOKExpiredCC(t *testing.T) {
	l, v := setupLedgerAndValidator(t)
	defer ledgermgmt.CleanupTestEnv()
	defer l.Close()

	ccID := "mycc"

	putCCInfoWithVSCCAndVer(l, ccID, "vscc", "badversion", signedByAnyMember([]string{"SampleOrg"}), t)

	tx := getEnv(ccID, nil, createRWset(t, ccID), t)
	b := &common.Block{Data: &common.BlockData{Data: [][]byte{utils.MarshalOrPanic(tx)}}}

	err := v.Validate(b)
	assert.NoError(t, err)
	assertInvalid(b, t, peer.TxValidationCode_EXPIRED_CHAINCODE)
}

func TestInvokeNOKBogusActions(t *testing.T) {
	l, v := setupLedgerAndValidator(t)
	defer ledgermgmt.CleanupTestEnv()
	defer l.Close()

	ccID := "mycc"

	putCCInfo(l, ccID, signedByAnyMember([]string{"SampleOrg"}), t)

	tx := getEnv(ccID, nil, []byte("barf"), t)
	b := &common.Block{Data: &common.BlockData{Data: [][]byte{utils.MarshalOrPanic(tx)}}}

	err := v.Validate(b)
	assert.NoError(t, err)
	assertInvalid(b, t, peer.TxValidationCode_BAD_RWSET)
}

func TestInvokeNOKCCDoesntExist(t *testing.T) {
	l, v := setupLedgerAndValidator(t)
	defer ledgermgmt.CleanupTestEnv()
	defer l.Close()

	ccID := "mycc"

	tx := getEnv(ccID, nil, createRWset(t, ccID), t)
	b := &common.Block{Data: &common.BlockData{Data: [][]byte{utils.MarshalOrPanic(tx)}}}

	err := v.Validate(b)
	assert.NoError(t, err)
	assertInvalid(b, t, peer.TxValidationCode_INVALID_OTHER_REASON)
}

func TestInvokeNOKVSCCUnspecified(t *testing.T) {
	l, v := setupLedgerAndValidator(t)
	defer ledgermgmt.CleanupTestEnv()
	defer l.Close()

	ccID := "mycc"

	putCCInfoWithVSCCAndVer(l, ccID, "", ccVersion, signedByAnyMember([]string{"SampleOrg"}), t)

	tx := getEnv(ccID, nil, createRWset(t, ccID), t)
	b := &common.Block{Data: &common.BlockData{Data: [][]byte{utils.MarshalOrPanic(tx)}}}

	err := v.Validate(b)
	assert.NoError(t, err)
	assertInvalid(b, t, peer.TxValidationCode_INVALID_OTHER_REASON)
}

func TestInvokeNoBlock(t *testing.T) {
	l, v := setupLedgerAndValidator(t)
	defer ledgermgmt.CleanupTestEnv()
	defer l.Close()

	err := v.Validate(&common.Block{Data: &common.BlockData{Data: [][]byte{}}})
	assert.NoError(t, err)
}





type mockLedger struct {
	mock.Mock
}


func (m *mockLedger) GetTransactionByID(txID string) (*peer.ProcessedTransaction, error) {
	args := m.Called(txID)
	return args.Get(0).(*peer.ProcessedTransaction), args.Error(1)
}


func (m *mockLedger) GetBlockByHash(blockHash []byte) (*common.Block, error) {
	args := m.Called(blockHash)
	return args.Get(0).(*common.Block), nil
}


func (m *mockLedger) GetBlockByTxID(txID string) (*common.Block, error) {
	args := m.Called(txID)
	return args.Get(0).(*common.Block), nil
}


func (m *mockLedger) GetTxValidationCodeByTxID(txID string) (peer.TxValidationCode, error) {
	args := m.Called(txID)
	return args.Get(0).(peer.TxValidationCode), nil
}


func (m *mockLedger) NewTxSimulator(txid string) (ledger.TxSimulator, error) {
	args := m.Called()
	return args.Get(0).(ledger.TxSimulator), nil
}


func (m *mockLedger) NewQueryExecutor() (ledger.QueryExecutor, error) {
	args := m.Called()
	return args.Get(0).(ledger.QueryExecutor), nil
}


func (m *mockLedger) NewHistoryQueryExecutor() (ledger.HistoryQueryExecutor, error) {
	args := m.Called()
	return args.Get(0).(ledger.HistoryQueryExecutor), nil
}


func (m *mockLedger) GetPvtDataAndBlockByNum(blockNum uint64, filter ledger.PvtNsCollFilter) (*ledger.BlockAndPvtData, error) {
	args := m.Called()
	return args.Get(0).(*ledger.BlockAndPvtData), nil
}


func (m *mockLedger) GetPvtDataByNum(blockNum uint64, filter ledger.PvtNsCollFilter) ([]*ledger.TxPvtData, error) {
	args := m.Called()
	return args.Get(0).([]*ledger.TxPvtData), nil
}


func (m *mockLedger) CommitWithPvtData(pvtDataAndBlock *ledger.BlockAndPvtData) error {
	return nil
}


func (m *mockLedger) PurgePrivateData(maxBlockNumToRetain uint64) error {
	return nil
}


func (m *mockLedger) PrivateDataMinBlockNum() (uint64, error) {
	return 0, nil
}


func (m *mockLedger) Prune(policy ledger2.PrunePolicy) error {
	return nil
}

func (m *mockLedger) GetBlockchainInfo() (*common.BlockchainInfo, error) {
	args := m.Called()
	return args.Get(0).(*common.BlockchainInfo), nil
}

func (m *mockLedger) GetBlockByNumber(blockNumber uint64) (*common.Block, error) {
	args := m.Called(blockNumber)
	return args.Get(0).(*common.Block), nil
}

func (m *mockLedger) GetBlocksIterator(startBlockNumber uint64) (ledger2.ResultsIterator, error) {
	args := m.Called(startBlockNumber)
	return args.Get(0).(ledger2.ResultsIterator), nil
}

func (m *mockLedger) Close() {

}

func (m *mockLedger) Commit(block *common.Block) error {
	return nil
}


func (m *mockLedger) GetConfigHistoryRetriever() (ledger.ConfigHistoryRetriever, error) {
	args := m.Called()
	return args.Get(0).(ledger.ConfigHistoryRetriever), nil
}








type mockQueryExecutor struct {
	mock.Mock
}

func (exec *mockQueryExecutor) GetState(namespace string, key string) ([]byte, error) {
	args := exec.Called(namespace, key)
	return args.Get(0).([]byte), args.Error(1)
}

func (exec *mockQueryExecutor) GetStateMultipleKeys(namespace string, keys []string) ([][]byte, error) {
	args := exec.Called(namespace, keys)
	return args.Get(0).([][]byte), args.Error(1)
}

func (exec *mockQueryExecutor) GetStateRangeScanIterator(namespace string, startKey string, endKey string) (ledger2.ResultsIterator, error) {
	args := exec.Called(namespace, startKey, endKey)
	return args.Get(0).(ledger2.ResultsIterator), args.Error(1)
}

func (exec *mockQueryExecutor) ExecuteQuery(namespace, query string) (ledger2.ResultsIterator, error) {
	args := exec.Called(namespace)
	return args.Get(0).(ledger2.ResultsIterator), args.Error(1)
}

func (exec *mockQueryExecutor) GetPrivateData(namespace, collection, key string) ([]byte, error) {
	args := exec.Called(namespace, collection, key)
	return args.Get(0).([]byte), args.Error(1)
}

func (exec *mockQueryExecutor) GetPrivateDataMultipleKeys(namespace, collection string, keys []string) ([][]byte, error) {
	args := exec.Called(namespace, collection, keys)
	return args.Get(0).([][]byte), args.Error(1)
}

func (exec *mockQueryExecutor) GetPrivateDataRangeScanIterator(namespace, collection, startKey, endKey string) (ledger2.ResultsIterator, error) {
	args := exec.Called(namespace, collection, startKey, endKey)
	return args.Get(0).(ledger2.ResultsIterator), args.Error(1)
}

func (exec *mockQueryExecutor) ExecuteQueryOnPrivateData(namespace, collection, query string) (ledger2.ResultsIterator, error) {
	args := exec.Called(namespace, collection, query)
	return args.Get(0).(ledger2.ResultsIterator), args.Error(1)
}

func (exec *mockQueryExecutor) Done() {
}

func (exec *mockQueryExecutor) GetStateMetadata(namespace, key string) (map[string][]byte, error) {
	return nil, nil
}

func (exec *mockQueryExecutor) GetPrivateDataMetadata(namespace, collection, key string) (map[string][]byte, error) {
	return nil, nil
}

func createCustomSupportAndLedger(t *testing.T) (*mocktxvalidator.Support, ledger.PeerLedger) {
	viper.Set("peer.fileSystemPath", "/tmp/blockchain/validatortest")
	ledgermgmt.InitializeTestEnv()
	gb, err := ctxt.MakeGenesisBlock("TestLedger")
	assert.NoError(t, err)
	l, err := ledgermgmt.CreateLedger(gb)
	assert.NoError(t, err)

	identity := &mocks2.Identity{}
	identity.GetIdentifierReturns(&msp.IdentityIdentifier{
		Mspid: "SampleOrg",
		Id:    "foo",
	})
	mspManager := &mocks2.MSPManager{}
	mspManager.DeserializeIdentityReturns(identity, nil)
	support := &mocktxvalidator.Support{LedgerVal: l, ACVal: &mockconfig.MockApplicationCapabilities{}, MSPManagerVal: mspManager}
	return support, l
}

func TestDynamicCapabilitiesAndMSP(t *testing.T) {
	factory := &mocks.PluginFactory{}
	factory.On("New").Return(&testdata.SampleValidationPlugin{})
	pm := &mocks.PluginMapper{}
	pm.On("PluginFactoryByName", txvalidator.PluginName("vscc")).Return(factory)

	support, l := createCustomSupportAndLedger(t)
	defer l.Close()

	vcs := struct {
		*mocktxvalidator.Support
		*semaphore.Weighted
	}{support, semaphore.NewWeighted(10)}
	mp := (&scc.MocksccProviderFactory{}).NewSystemChaincodeProvider()

	v := txvalidator.NewTxValidator(vcs, mp, pm)

	ccID := "mycc"

	putCCInfo(l, ccID, signedByAnyMember([]string{"SampleOrg"}), t)

	tx := getEnv(ccID, nil, createRWset(t, ccID), t)
	b := &common.Block{Data: &common.BlockData{Data: [][]byte{utils.MarshalOrPanic(tx)}}}

	
	err := v.Validate(b)
	assert.NoError(t, err)
	assertValid(b, t)
	
	capabilityInvokeCount := support.CapabilitiesInvokeCount
	mspManagerInvokeCount := support.MSPManagerInvokeCount

	
	err = v.Validate(b)
	assert.NoError(t, err)
	assertValid(b, t)

	
	
	assert.Equal(t, 2*capabilityInvokeCount, support.CapabilitiesInvokeCount)
	
	
	assert.Equal(t, 2*mspManagerInvokeCount, support.MSPManagerInvokeCount)
}








func TestLedgerIsNoAvailable(t *testing.T) {
	theLedger := new(mockLedger)
	vcs := struct {
		*mocktxvalidator.Support
		*semaphore.Weighted
	}{&mocktxvalidator.Support{LedgerVal: theLedger, ACVal: &mockconfig.MockApplicationCapabilities{}}, semaphore.NewWeighted(10)}
	mp := (&scc.MocksccProviderFactory{}).NewSystemChaincodeProvider()
	pm := &mocks.PluginMapper{}
	validator := txvalidator.NewTxValidator(vcs, mp, pm)

	ccID := "mycc"
	tx := getEnv(ccID, nil, createRWset(t, ccID), t)

	theLedger.On("GetTransactionByID", mock.Anything).Return(&peer.ProcessedTransaction{}, ledger.NotFoundInIndexErr(""))

	queryExecutor := new(mockQueryExecutor)
	queryExecutor.On("GetState", mock.Anything, mock.Anything).Return([]byte{}, errors.New("Unable to connect to DB"))
	theLedger.On("NewQueryExecutor", mock.Anything).Return(queryExecutor, nil)

	b := &common.Block{Data: &common.BlockData{Data: [][]byte{utils.MarshalOrPanic(tx)}}}

	err := validator.Validate(b)

	assertion := assert.New(t)
	
	assertion.Error(err)
	
	assertion.NotNil(err.(*commonerrors.VSCCInfoLookupFailureError))
}

func TestLedgerIsNotAvailableForCheckingTxidDuplicate(t *testing.T) {
	theLedger := new(mockLedger)
	vcs := struct {
		*mocktxvalidator.Support
		*semaphore.Weighted
	}{&mocktxvalidator.Support{LedgerVal: theLedger, ACVal: &mockconfig.MockApplicationCapabilities{}}, semaphore.NewWeighted(10)}
	mp := (&scc.MocksccProviderFactory{}).NewSystemChaincodeProvider()
	pm := &mocks.PluginMapper{}
	validator := txvalidator.NewTxValidator(vcs, mp, pm)

	ccID := "mycc"
	tx := getEnv(ccID, nil, createRWset(t, ccID), t)

	theLedger.On("GetTransactionByID", mock.Anything).Return(&peer.ProcessedTransaction{}, errors.New("Unable to connect to DB"))

	b := &common.Block{Data: &common.BlockData{Data: [][]byte{utils.MarshalOrPanic(tx)}}}

	err := validator.Validate(b)

	assertion := assert.New(t)
	
	assertion.Error(err)
}

func TestDuplicateTxId(t *testing.T) {
	theLedger := new(mockLedger)
	vcs := struct {
		*mocktxvalidator.Support
		*semaphore.Weighted
	}{&mocktxvalidator.Support{LedgerVal: theLedger, ACVal: &mockconfig.MockApplicationCapabilities{}}, semaphore.NewWeighted(10)}
	mp := (&scc.MocksccProviderFactory{}).NewSystemChaincodeProvider()
	pm := &mocks.PluginMapper{}
	validator := txvalidator.NewTxValidator(vcs, mp, pm)

	ccID := "mycc"
	tx := getEnv(ccID, nil, createRWset(t, ccID), t)

	theLedger.On("GetTransactionByID", mock.Anything).Return(&peer.ProcessedTransaction{}, nil)

	b := &common.Block{Data: &common.BlockData{Data: [][]byte{utils.MarshalOrPanic(tx)}}}

	err := validator.Validate(b)

	assertion := assert.New(t)
	
	assertion.NoError(err)

	
	txsfltr := lutils.TxValidationFlags(b.Metadata.Metadata[common.BlockMetadataIndex_TRANSACTIONS_FILTER])
	assertion.True(txsfltr.IsInvalid(0))
	assertion.True(txsfltr.Flag(0) == peer.TxValidationCode_DUPLICATE_TXID)
}

func TestValidationInvalidEndorsing(t *testing.T) {
	theLedger := new(mockLedger)
	vcs := struct {
		*mocktxvalidator.Support
		*semaphore.Weighted
	}{&mocktxvalidator.Support{LedgerVal: theLedger, ACVal: &mockconfig.MockApplicationCapabilities{}}, semaphore.NewWeighted(10)}
	mp := (&scc.MocksccProviderFactory{}).NewSystemChaincodeProvider()
	pm := &mocks.PluginMapper{}
	factory := &mocks.PluginFactory{}
	plugin := &mocks.Plugin{}
	factory.On("New").Return(plugin)
	plugin.On("Init", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	plugin.On("Validate", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("invalid tx"))
	pm.On("PluginFactoryByName", txvalidator.PluginName("vscc")).Return(factory)
	validator := txvalidator.NewTxValidator(vcs, mp, pm)

	ccID := "mycc"
	tx := getEnv(ccID, nil, createRWset(t, ccID), t)

	theLedger.On("GetTransactionByID", mock.Anything).Return(&peer.ProcessedTransaction{}, ledger.NotFoundInIndexErr(""))

	cd := &ccp.ChaincodeData{
		Name:    ccID,
		Version: ccVersion,
		Vscc:    "vscc",
		Policy:  signedByAnyMember([]string{"SampleOrg"}),
	}

	cdbytes := utils.MarshalOrPanic(cd)

	queryExecutor := new(mockQueryExecutor)
	queryExecutor.On("GetState", "lscc", ccID).Return(cdbytes, nil)
	theLedger.On("NewQueryExecutor", mock.Anything).Return(queryExecutor, nil)

	b := &common.Block{Data: &common.BlockData{Data: [][]byte{utils.MarshalOrPanic(tx)}}}

	
	err := validator.Validate(b)
	
	assert.NoError(t, err)
	assertInvalid(b, t, peer.TxValidationCode_ENDORSEMENT_POLICY_FAILURE)
}

func createMockLedger(t *testing.T, ccID string) *mockLedger {
	l := new(mockLedger)
	l.On("GetTransactionByID", mock.Anything).Return(&peer.ProcessedTransaction{}, ledger.NotFoundInIndexErr(""))
	cd := &ccp.ChaincodeData{
		Name:    ccID,
		Version: ccVersion,
		Vscc:    "vscc",
		Policy:  signedByAnyMember([]string{"SampleOrg"}),
	}

	cdbytes := utils.MarshalOrPanic(cd)
	queryExecutor := new(mockQueryExecutor)
	queryExecutor.On("GetState", "lscc", ccID).Return(cdbytes, nil)
	l.On("NewQueryExecutor", mock.Anything).Return(queryExecutor, nil)
	return l
}

func TestValidationPluginExecutionError(t *testing.T) {
	plugin := &mocks.Plugin{}
	plugin.On("Init", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	l, v := setupLedgerAndValidatorExplicit(t, &mockconfig.MockApplicationCapabilities{}, plugin)
	defer ledgermgmt.CleanupTestEnv()
	defer l.Close()

	ccID := "mycc"
	putCCInfo(l, ccID, signedByAnyMember([]string{"SampleOrg"}), t)

	tx := getEnv(ccID, nil, createRWset(t, ccID), t)
	b := &common.Block{Data: &common.BlockData{Data: [][]byte{utils.MarshalOrPanic(tx)}}}

	plugin.On("Validate", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&validation.ExecutionFailureError{
		Reason: "I/O error",
	})

	err := v.Validate(b)
	executionErr := err.(*commonerrors.VSCCExecutionFailureError)
	assert.Contains(t, executionErr.Error(), "I/O error")
}

func TestValidationPluginNotFound(t *testing.T) {
	ccID := "mycc"
	tx := getEnv(ccID, nil, createRWset(t, ccID), t)
	l := createMockLedger(t, ccID)
	vcs := struct {
		*mocktxvalidator.Support
		*semaphore.Weighted
	}{&mocktxvalidator.Support{LedgerVal: l, ACVal: &mockconfig.MockApplicationCapabilities{}}, semaphore.NewWeighted(10)}

	b := &common.Block{Data: &common.BlockData{Data: [][]byte{utils.MarshalOrPanic(tx)}}}

	pm := &mocks.PluginMapper{}
	pm.On("PluginFactoryByName", txvalidator.PluginName("vscc")).Return(nil)
	mp := (&scc.MocksccProviderFactory{}).NewSystemChaincodeProvider()
	validator := txvalidator.NewTxValidator(vcs, mp, pm)
	err := validator.Validate(b)
	executionErr := err.(*commonerrors.VSCCExecutionFailureError)
	assert.Contains(t, executionErr.Error(), "plugin with name vscc wasn't found")
}

var signer msp.SigningIdentity

var signerSerialized []byte

func TestMain(m *testing.M) {
	msptesttools.LoadMSPSetupForTesting()

	var err error
	signer, err = mgmt.GetLocalMSP().GetDefaultSigningIdentity()
	if err != nil {
		fmt.Println("Could not get signer")
		os.Exit(-1)
		return
	}

	signerSerialized, err = signer.Serialize()
	if err != nil {
		fmt.Println("Could not serialize identity")
		os.Exit(-1)
		return
	}

	os.Exit(m.Run())
}
