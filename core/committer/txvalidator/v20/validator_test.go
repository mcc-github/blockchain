/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package txvalidator_test

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/common/cauthdsl"
	commonerrors "github.com/mcc-github/blockchain/common/errors"
	"github.com/mcc-github/blockchain/common/ledger/testutil"
	mockconfig "github.com/mcc-github/blockchain/common/mocks/config"
	"github.com/mcc-github/blockchain/common/semaphore"
	"github.com/mcc-github/blockchain/common/util"
	"github.com/mcc-github/blockchain/core/committer/txvalidator"
	vp "github.com/mcc-github/blockchain/core/committer/txvalidator/plugin"
	txvalidatorv20 "github.com/mcc-github/blockchain/core/committer/txvalidator/v20"
	mocks3 "github.com/mcc-github/blockchain/core/committer/txvalidator/v20/mocks"
	"github.com/mcc-github/blockchain/core/committer/txvalidator/v20/plugindispatcher/mocks"
	ccp "github.com/mcc-github/blockchain/core/common/ccprovider"
	validation "github.com/mcc-github/blockchain/core/handlers/validation/api"
	"github.com/mcc-github/blockchain/core/handlers/validation/builtin"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/rwsetutil"
	lutils "github.com/mcc-github/blockchain/core/ledger/util"
	mocktxvalidator "github.com/mcc-github/blockchain/core/mocks/txvalidator"
	"github.com/mcc-github/blockchain/core/scc/lscc"
	mocks2 "github.com/mcc-github/blockchain/discovery/support/mocks"
	"github.com/mcc-github/blockchain/msp"
	"github.com/mcc-github/blockchain/msp/mgmt"
	msptesttools "github.com/mcc-github/blockchain/msp/mgmt/testtools"
	"github.com/mcc-github/blockchain/protos/common"
	mb "github.com/mcc-github/blockchain/protos/msp"
	"github.com/mcc-github/blockchain/protos/peer"
	pb "github.com/mcc-github/blockchain/protos/peer"
	"github.com/mcc-github/blockchain/protos/token"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func signedByAnyMember(ids []string) []byte {
	p := cauthdsl.SignedByAnyMember(ids)
	return protoutil.MarshalOrPanic(&pb.ApplicationPolicy{Type: &pb.ApplicationPolicy_SignaturePolicy{SignaturePolicy: p}})
}

func v20Capabilities() *mockconfig.MockApplicationCapabilities {
	return &mockconfig.MockApplicationCapabilities{
		V1_2ValidationRv:      true,
		V1_3ValidationRv:      true,
		PrivateChannelDataRv:  true,
		KeyLevelEndorsementRv: true,
		V2_0ValidationRv:      true,
	}
}

func fabTokenCapabilities() *mockconfig.MockApplicationCapabilities {
	return &mockconfig.MockApplicationCapabilities{
		V1_2ValidationRv:      true,
		V1_3ValidationRv:      true,
		PrivateChannelDataRv:  true,
		KeyLevelEndorsementRv: true,
		V2_0ValidationRv:      true,
		FabTokenRv:            true,
	}
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

	proposal, _, err := protoutil.CreateProposalFromCIS(pType, util.GetTestChainID(), cis, signerSerialized)
	return proposal, err
}

const ccVersion = "1.0"

func getEnvWithType(ccID string, event []byte, res []byte, pType common.HeaderType, t *testing.T) *common.Envelope {
	
	prop, err := getProposalWithType(ccID, pType)
	assert.NoError(t, err)

	response := &peer.Response{Status: 200}

	
	presp, err := protoutil.CreateProposalResponse(prop.Header, prop.Payload, response, res, event, &peer.ChaincodeID{Name: ccID, Version: ccVersion}, nil, signer)
	assert.NoError(t, err)

	
	tx, err := protoutil.CreateSignedTx(prop, signer, presp)
	assert.NoError(t, err)

	return tx
}

func getEnv(ccID string, event []byte, res []byte, t *testing.T) *common.Envelope {
	return getEnvWithType(ccID, event, res, common.HeaderType_ENDORSER_TRANSACTION, t)
}

func getEnvWithSigner(ccID string, event []byte, res []byte, sig msp.SigningIdentity, t *testing.T) *common.Envelope {
	
	pType := common.HeaderType_ENDORSER_TRANSACTION
	cis := &peer.ChaincodeInvocationSpec{
		ChaincodeSpec: &peer.ChaincodeSpec{
			ChaincodeId: &peer.ChaincodeID{Name: ccID, Version: ccVersion},
			Input:       &peer.ChaincodeInput{Args: [][]byte{[]byte("func")}},
			Type:        peer.ChaincodeSpec_GOLANG,
		},
	}

	sID, err := sig.Serialize()
	assert.NoError(t, err)
	prop, _, err := protoutil.CreateProposalFromCIS(pType, "foochain", cis, sID)
	assert.NoError(t, err)

	response := &peer.Response{Status: 200}

	
	presp, err := protoutil.CreateProposalResponse(prop.Header, prop.Payload, response, res, event, &peer.ChaincodeID{Name: ccID, Version: ccVersion}, nil, sig)
	assert.NoError(t, err)

	
	tx, err := protoutil.CreateSignedTx(prop, sig, presp)
	assert.NoError(t, err)

	return tx
}

func getTokenTx(t *testing.T) *common.Envelope {
	transactionData := &token.TokenTransaction{
		Action: &token.TokenTransaction_TokenAction{
			TokenAction: &token.TokenAction{
				Data: &token.TokenAction_Issue{
					Issue: &token.Issue{
						Outputs: []*token.Token{
							{Owner: &token.TokenOwner{Raw: []byte("owner-1")}, Type: "TOK1", Quantity: ToHex(111)},
							{Owner: &token.TokenOwner{Raw: []byte("owner-2")}, Type: "TOK2", Quantity: ToHex(222)},
						},
					},
				},
			},
		},
	}
	tdBytes, err := proto.Marshal(transactionData)
	assert.NoError(t, err)

	signerBytes, err := signer.Serialize()
	assert.NoError(t, err)
	nonce := []byte{0, 1, 2, 3, 4}
	txID, err := protoutil.ComputeTxID(nonce, signerBytes)
	assert.NoError(t, err)

	hdr := &common.Header{
		SignatureHeader: protoutil.MarshalOrPanic(
			&common.SignatureHeader{
				Creator: signerBytes,
				Nonce:   nonce,
			},
		),
		ChannelHeader: protoutil.MarshalOrPanic(
			&common.ChannelHeader{
				Type: int32(common.HeaderType_TOKEN_TRANSACTION),
				TxId: txID,
			},
		),
	}

	
	payl := &common.Payload{Header: hdr, Data: tdBytes}
	paylBytes, err := protoutil.GetBytesPayload(payl)
	assert.NoError(t, err)

	
	sig, err := signer.Sign(paylBytes)
	assert.NoError(t, err)

	
	return &common.Envelope{Payload: paylBytes, Signature: sig}
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

func setupValidator() (*txvalidatorv20.TxValidator, *mocks3.QueryExecutor, *mocks2.Identity, *mocks3.CollectionResources) {
	mspmgr := &mocks2.MSPManager{}
	mockID := &mocks2.Identity{}
	mockID.SatisfiesPrincipalReturns(nil)
	mockID.GetIdentifierReturns(&msp.IdentityIdentifier{})
	mspmgr.DeserializeIdentityReturns(mockID, nil)

	return setupValidatorWithMspMgr(mspmgr, mockID)
}

func setupValidatorWithMspMgr(mspmgr msp.MSPManager, mockID *mocks2.Identity) (*txvalidatorv20.TxValidator, *mocks3.QueryExecutor, *mocks2.Identity, *mocks3.CollectionResources) {
	pm := &mocks.Mapper{}
	factory := &mocks.PluginFactory{}
	pm.On("FactoryByName", vp.Name("vscc")).Return(factory)
	factory.On("New").Return(&builtin.DefaultValidation{})

	mockQE := &mocks3.QueryExecutor{}
	mockQE.On("Done").Return(nil)
	mockQE.On("GetState", "lscc", "lscc").Return(nil, nil)
	mockQE.On("GetState", "lscc", "escc").Return(nil, nil)

	mockLedger := &mocks3.LedgerResources{}
	mockLedger.On("GetTransactionByID", mock.Anything).Return(nil, ledger.NotFoundInIndexErr("As idle as a painted ship upon a painted ocean"))
	mockLedger.On("NewQueryExecutor").Return(mockQE, nil)

	mockCpmg := &mocks.ChannelPolicyManagerGetter{}
	mockCpmg.On("Manager", mock.Anything).Return(nil, true)

	mockCR := &mocks3.CollectionResources{}

	v := txvalidatorv20.NewTxValidator(
		"",
		semaphore.New(10),
		&mocktxvalidator.Support{ACVal: v20Capabilities(), MSPManagerVal: mspmgr},
		mockLedger,
		&lscc.LifeCycleSysCC{},
		mockCR,
		pm,
		mockCpmg,
	)

	return v, mockQE, mockID, mockCR
}

func TestInvokeBadRWSet(t *testing.T) {
	ccID := "mycc"

	v, _, _, _ := setupValidator()

	tx := getEnv(ccID, nil, []byte("barf"), t)
	b := &common.Block{Data: &common.BlockData{Data: [][]byte{protoutil.MarshalOrPanic(tx)}}, Header: &common.BlockHeader{Number: 1}}

	err := v.Validate(b)
	assert.NoError(t, err)
	assertInvalid(b, t, peer.TxValidationCode_BAD_RWSET)
}

func TestInvokeNoPolicy(t *testing.T) {
	ccID := "mycc"

	v, mockQE, _, _ := setupValidator()

	mockQE.On("GetState", "lscc", ccID).Return(protoutil.MarshalOrPanic(&ccp.ChaincodeData{
		Name:    ccID,
		Version: ccVersion,
		Vscc:    "vscc",
		Policy:  nil,
	}), nil)

	tx := getEnv(ccID, nil, createRWset(t, ccID), t)
	b := &common.Block{Data: &common.BlockData{Data: [][]byte{protoutil.MarshalOrPanic(tx)}}, Header: &common.BlockHeader{Number: 2}}

	err := v.Validate(b)
	assert.NoError(t, err)
	assertInvalid(b, t, peer.TxValidationCode_INVALID_CHAINCODE)
}

func TestInvokeOK(t *testing.T) {
	ccID := "mycc"

	v, mockQE, _, _ := setupValidator()

	mockQE.On("GetState", "lscc", ccID).Return(protoutil.MarshalOrPanic(&ccp.ChaincodeData{
		Name:    ccID,
		Version: ccVersion,
		Vscc:    "vscc",
		Policy:  signedByAnyMember([]string{"SampleOrg"}),
	}), nil)
	mockQE.On("GetStateMetadata", ccID, "key").Return(nil, nil)

	tx := getEnv(ccID, nil, createRWset(t, ccID), t)
	b := &common.Block{Data: &common.BlockData{Data: [][]byte{protoutil.MarshalOrPanic(tx)}}, Header: &common.BlockHeader{Number: 2}}

	err := v.Validate(b)
	assert.NoError(t, err)
	assertValid(b, t)
}

func TestInvokeNoRWSet(t *testing.T) {
	ccID := "mycc"

	v, mockQE, mockID, _ := setupValidator()
	mockID.SatisfiesPrincipalReturns(errors.New("principal not satisfied"))

	mockQE.On("GetState", "lscc", ccID).Return(protoutil.MarshalOrPanic(&ccp.ChaincodeData{
		Name:    ccID,
		Version: ccVersion,
		Vscc:    "vscc",
		Policy:  signedByAnyMember([]string{"SampleOrg"}),
	}), nil)
	mockQE.On("GetStateMetadata", ccID, "key").Return(nil, nil)

	tx := getEnv(ccID, nil, createRWset(t), t)
	b := &common.Block{Data: &common.BlockData{Data: [][]byte{protoutil.MarshalOrPanic(tx)}}, Header: &common.BlockHeader{Number: 2}}

	err := v.Validate(b)
	assert.NoError(t, err)
	assertInvalid(b, t, peer.TxValidationCode_ENDORSEMENT_POLICY_FAILURE)
}


type mockSI struct {
	SerializedID []byte
	MspID        string
	SatPrinError error
}

func (msi *mockSI) ExpiresAt() time.Time {
	return time.Now()
}

func (msi *mockSI) GetIdentifier() *msp.IdentityIdentifier {
	return &msp.IdentityIdentifier{
		Mspid: msi.MspID,
		Id:    "",
	}
}

func (msi *mockSI) GetMSPIdentifier() string {
	return msi.MspID
}

func (msi *mockSI) Validate() error {
	return nil
}

func (msi *mockSI) GetOrganizationalUnits() []*msp.OUIdentifier {
	return nil
}

func (msi *mockSI) Anonymous() bool {
	return false
}

func (msi *mockSI) Verify(msg []byte, sig []byte) error {
	return nil
}

func (msi *mockSI) Serialize() ([]byte, error) {
	sid := &mb.SerializedIdentity{
		Mspid:   msi.MspID,
		IdBytes: msi.SerializedID,
	}
	sidBytes := protoutil.MarshalOrPanic(sid)
	return sidBytes, nil
}

func (msi *mockSI) SatisfiesPrincipal(principal *mb.MSPPrincipal) error {
	return msi.SatPrinError
}

func (msi *mockSI) Sign(msg []byte) ([]byte, error) {
	return msg, nil
}

func (msi *mockSI) GetPublicVersion() msp.Identity {
	return msi
}


type mockMSP struct {
	ID           msp.Identity
	SatPrinError error
	MspID        string
}

func (fake *mockMSP) DeserializeIdentity(serializedIdentity []byte) (msp.Identity, error) {
	return fake.ID, nil
}

func (fake *mockMSP) IsWellFormed(identity *mb.SerializedIdentity) error {
	return nil
}
func (fake *mockMSP) Setup(config *mb.MSPConfig) error {
	return nil
}

func (fake *mockMSP) GetVersion() msp.MSPVersion {
	return msp.MSPv1_3
}

func (fake *mockMSP) GetType() msp.ProviderType {
	return msp.FABRIC
}

func (fake *mockMSP) GetIdentifier() (string, error) {
	return fake.MspID, nil
}

func (fake *mockMSP) GetSigningIdentity(identifier *msp.IdentityIdentifier) (msp.SigningIdentity, error) {
	return nil, nil
}

func (fake *mockMSP) GetDefaultSigningIdentity() (msp.SigningIdentity, error) {
	return nil, nil
}

func (fake *mockMSP) GetTLSRootCerts() [][]byte {
	return nil
}

func (fake *mockMSP) GetTLSIntermediateCerts() [][]byte {
	return nil
}

func (fake *mockMSP) Validate(id msp.Identity) error {
	return nil
}

func (fake *mockMSP) SatisfiesPrincipal(id msp.Identity, principal *mb.MSPPrincipal) error {
	return fake.SatPrinError
}


func TestParallelValidation(t *testing.T) {
	
	txCnt := 100

	
	msp1 := &mockMSP{
		ID: &mockSI{
			MspID:        "Org1",
			SerializedID: []byte("signer0"),
			SatPrinError: nil,
		},
		SatPrinError: nil,
		MspID:        "Org1",
	}
	msp2 := &mockMSP{
		ID: &mockSI{
			MspID:        "Org2",
			SerializedID: []byte("signer1"),
			SatPrinError: errors.New("nope"),
		},
		SatPrinError: errors.New("nope"),
		MspID:        "Org2",
	}
	mgmt.GetManagerForChain("foochain")
	mgr := mgmt.GetManagerForChain("foochain")
	mgr.Setup([]msp.MSP{msp1, msp2})

	vpKey := pb.MetaDataKeys_VALIDATION_PARAMETER.String()
	ccID := "mycc"

	v, mockQE, _, mockCR := setupValidatorWithMspMgr(mgr, nil)

	mockCR.On("CollectionValidationInfo", ccID, "col1", mock.Anything).Return(nil, nil, nil)

	policy := cauthdsl.SignedByMspPeer("Org1")
	polBytes := protoutil.MarshalOrPanic(&pb.ApplicationPolicy{Type: &pb.ApplicationPolicy_SignaturePolicy{SignaturePolicy: policy}})
	mockQE.On("GetState", "lscc", ccID).Return(protoutil.MarshalOrPanic(&ccp.ChaincodeData{
		Name:    ccID,
		Version: ccVersion,
		Vscc:    "vscc",
		Policy:  polBytes,
	}), nil)
	mockQE.On("GetStateMetadata", "mycc", mock.Anything).Return(nil, nil)
	mockQE.On("GetPrivateDataMetadataByHash", "mycc", "col1", mock.Anything).Return(nil, nil)

	
	blockData := make([][]byte, 0, txCnt)
	col := "col1"
	sigID0 := &mockSI{
		SerializedID: []byte("signer0"),
		MspID:        "Org1",
	}
	sigID1 := &mockSI{
		SerializedID: []byte("signer1"),
		MspID:        "Org2",
	}
	for txNum := 0; txNum < txCnt; txNum++ {
		var sig msp.SigningIdentity
		
		key := strconv.Itoa(txNum % 10)
		rwsetBuilder := rwsetutil.NewRWSetBuilder()
		
		switch uint(txNum / 10) {
		case 0:
			
			rwsetBuilder.AddToWriteSet(ccID, key, []byte("value1"))
			sig = sigID0
		case 1:
			
			metadata := make(map[string][]byte)
			metadata[vpKey] = signedByAnyMember([]string{"SampleOrg"})
			rwsetBuilder.AddToMetadataWriteSet(ccID, key, metadata)
			sig = sigID1
		case 2:
			
			rwsetBuilder.AddToWriteSet(ccID, key, []byte("value2"))
			sig = sigID0
		case 3:
			
			metadata := make(map[string][]byte)
			metadata[vpKey] = signedByAnyMember([]string{"Org2"})
			rwsetBuilder.AddToMetadataWriteSet(ccID, key, metadata)
			sig = sigID0
		case 4:
			
			rwsetBuilder.AddToWriteSet(ccID, key, []byte("value3"))
			sig = &mockSI{
				SerializedID: []byte("signer0"),
				MspID:        "Org1",
			}
		
		case 5:
			
			rwsetBuilder.AddToPvtAndHashedWriteSet(ccID, col, key, []byte("value1"))
			sig = sigID0
		case 6:
			
			metadata := make(map[string][]byte)
			metadata[vpKey] = signedByAnyMember([]string{"SampleOrg"})
			rwsetBuilder.AddToHashedMetadataWriteSet(ccID, col, key, metadata)
			sig = sigID1
		case 7:
			
			rwsetBuilder.AddToPvtAndHashedWriteSet(ccID, col, key, []byte("value2"))
			sig = sigID0
		case 8:
			
			metadata := make(map[string][]byte)
			metadata[vpKey] = signedByAnyMember([]string{"Org2"})
			rwsetBuilder.AddToHashedMetadataWriteSet(ccID, col, key, metadata)
			sig = sigID0
		case 9:
			
			rwsetBuilder.AddToPvtAndHashedWriteSet(ccID, col, key, []byte("value3"))
			sig = sigID0
		}
		rwset, err := rwsetBuilder.GetTxSimulationResults()
		assert.NoError(t, err)
		rwsetBytes, err := rwset.GetPubSimulationBytes()
		tx := getEnvWithSigner(ccID, nil, rwsetBytes, sig, t)
		blockData = append(blockData, protoutil.MarshalOrPanic(tx))
	}

	
	b := &common.Block{Data: &common.BlockData{Data: blockData}, Header: &common.BlockHeader{Number: uint64(txCnt)}}

	
	err := v.Validate(b)
	assert.NoError(t, err)

	
	txsFilter := lutils.TxValidationFlags(b.Metadata.Metadata[common.BlockMetadataIndex_TRANSACTIONS_FILTER])
	
	for txNum := 0; txNum < txCnt; txNum += 1 {
		switch uint(txNum / 10) {
		case 1:
			fallthrough
		case 4:
			fallthrough
		case 6:
			fallthrough
		case 9:
			assert.True(t, txsFilter.IsInvalid(txNum))
		default:
			assert.False(t, txsFilter.IsInvalid(txNum))
		}
	}
}

func TestChaincodeEvent(t *testing.T) {
	ccID := "mycc"

	t.Run("MisMatchedName", func(t *testing.T) {
		v, mockQE, _, _ := setupValidator()

		mockQE.On("GetState", "lscc", ccID).Return(protoutil.MarshalOrPanic(&ccp.ChaincodeData{
			Name:    ccID,
			Version: ccVersion,
			Vscc:    "vscc",
			Policy:  signedByAnyMember([]string{"SampleOrg"}),
		}), nil)
		mockQE.On("GetStateMetadata", ccID, "key").Return(nil, nil)

		testCCEventMismatchedName(t, v, ccID)
	})

	t.Run("BadBytes", func(t *testing.T) {
		v, mockQE, _, _ := setupValidator()

		mockQE.On("GetState", "lscc", ccID).Return(protoutil.MarshalOrPanic(&ccp.ChaincodeData{
			Name:    ccID,
			Version: ccVersion,
			Vscc:    "vscc",
			Policy:  signedByAnyMember([]string{"SampleOrg"}),
		}), nil)
		mockQE.On("GetStateMetadata", ccID, "key").Return(nil, nil)

		testCCEventBadBytes(t, v, ccID)
	})

	t.Run("GoodPath", func(t *testing.T) {
		v, mockQE, _, _ := setupValidator()

		mockQE.On("GetState", "lscc", ccID).Return(protoutil.MarshalOrPanic(&ccp.ChaincodeData{
			Name:    ccID,
			Version: ccVersion,
			Vscc:    "vscc",
			Policy:  signedByAnyMember([]string{"SampleOrg"}),
		}), nil)
		mockQE.On("GetStateMetadata", ccID, "key").Return(nil, nil)

		testCCEventGoodPath(t, v, ccID)
	})
}

func testCCEventMismatchedName(t *testing.T, v txvalidator.Validator, ccID string) {
	tx := getEnv(ccID, protoutil.MarshalOrPanic(&peer.ChaincodeEvent{ChaincodeId: "wrong"}), createRWset(t), t)
	b := &common.Block{Data: &common.BlockData{Data: [][]byte{protoutil.MarshalOrPanic(tx)}}, Header: &common.BlockHeader{Number: 2}}

	err := v.Validate(b)
	assert.NoError(t, err) 
	assertInvalid(b, t, peer.TxValidationCode_INVALID_OTHER_REASON)
}

func testCCEventBadBytes(t *testing.T, v txvalidator.Validator, ccID string) {
	tx := getEnv(ccID, []byte("garbage"), createRWset(t), t)
	b := &common.Block{Data: &common.BlockData{Data: [][]byte{protoutil.MarshalOrPanic(tx)}}, Header: &common.BlockHeader{Number: 2}}

	err := v.Validate(b)
	assert.NoError(t, err) 
	assertInvalid(b, t, peer.TxValidationCode_INVALID_OTHER_REASON)
}

func testCCEventGoodPath(t *testing.T, v txvalidator.Validator, ccID string) {
	tx := getEnv(ccID, protoutil.MarshalOrPanic(&peer.ChaincodeEvent{ChaincodeId: ccID}), createRWset(t), t)
	b := &common.Block{Data: &common.BlockData{Data: [][]byte{protoutil.MarshalOrPanic(tx)}}, Header: &common.BlockHeader{Number: 2}}

	err := v.Validate(b)
	assert.NoError(t, err)
	assertValid(b, t)
}

func TestInvokeOKPvtDataOnly(t *testing.T) {
	ccID := "mycc"

	v, mockQE, mockID, mockCR := setupValidator()
	mockID.SatisfiesPrincipalReturns(errors.New("principal not satisfied"))

	mockQE.On("GetState", "lscc", ccID).Return(protoutil.MarshalOrPanic(&ccp.ChaincodeData{
		Name:    ccID,
		Version: ccVersion,
		Vscc:    "vscc",
		Policy:  signedByAnyMember([]string{"SampleOrg"}),
	}), nil)
	mockQE.On("GetPrivateDataMetadataByHash", ccID, "mycollection", mock.Anything).Return(nil, nil)

	mockCR.On("CollectionValidationInfo", ccID, "mycollection", mock.Anything).Return(nil, nil, nil)

	rwsetBuilder := rwsetutil.NewRWSetBuilder()
	rwsetBuilder.AddToPvtAndHashedWriteSet(ccID, "mycollection", "somekey", nil)
	rwset, err := rwsetBuilder.GetTxSimulationResults()
	assert.NoError(t, err)
	rwsetBytes, err := rwset.GetPubSimulationBytes()
	assert.NoError(t, err)

	tx := getEnv(ccID, nil, rwsetBytes, t)
	b := &common.Block{Data: &common.BlockData{Data: [][]byte{protoutil.MarshalOrPanic(tx)}}, Header: &common.BlockHeader{Number: 2}}

	err = v.Validate(b)
	assert.NoError(t, err)
	assertInvalid(b, t, peer.TxValidationCode_ENDORSEMENT_POLICY_FAILURE)
}

func TestInvokeOKMetaUpdateOnly(t *testing.T) {
	ccID := "mycc"

	v, mockQE, mockID, _ := setupValidator()
	mockID.SatisfiesPrincipalReturns(errors.New("principal not satisfied"))

	mockQE.On("GetState", "lscc", ccID).Return(protoutil.MarshalOrPanic(&ccp.ChaincodeData{
		Name:    ccID,
		Version: ccVersion,
		Vscc:    "vscc",
		Policy:  signedByAnyMember([]string{"SampleOrg"}),
	}), nil)
	mockQE.On("GetStateMetadata", ccID, "somekey").Return(nil, nil)

	rwsetBuilder := rwsetutil.NewRWSetBuilder()
	rwsetBuilder.AddToMetadataWriteSet(ccID, "somekey", map[string][]byte{})
	rwset, err := rwsetBuilder.GetTxSimulationResults()
	assert.NoError(t, err)
	rwsetBytes, err := rwset.GetPubSimulationBytes()
	assert.NoError(t, err)

	tx := getEnv(ccID, nil, rwsetBytes, t)
	b := &common.Block{Data: &common.BlockData{Data: [][]byte{protoutil.MarshalOrPanic(tx)}}, Header: &common.BlockHeader{Number: 2}}

	err = v.Validate(b)
	assert.NoError(t, err)
	assertInvalid(b, t, peer.TxValidationCode_ENDORSEMENT_POLICY_FAILURE)
}

func TestInvokeOKPvtMetaUpdateOnly(t *testing.T) {
	ccID := "mycc"

	v, mockQE, mockID, mockCR := setupValidator()
	mockID.SatisfiesPrincipalReturns(errors.New("principal not satisfied"))

	mockQE.On("GetState", "lscc", ccID).Return(protoutil.MarshalOrPanic(&ccp.ChaincodeData{
		Name:    ccID,
		Version: ccVersion,
		Vscc:    "vscc",
		Policy:  signedByAnyMember([]string{"SampleOrg"}),
	}), nil)
	mockQE.On("GetPrivateDataMetadataByHash", ccID, "mycollection", mock.Anything).Return(nil, nil)

	mockCR.On("CollectionValidationInfo", ccID, "mycollection", mock.Anything).Return(nil, nil, nil)

	rwsetBuilder := rwsetutil.NewRWSetBuilder()
	rwsetBuilder.AddToHashedMetadataWriteSet(ccID, "mycollection", "somekey", map[string][]byte{})
	rwset, err := rwsetBuilder.GetTxSimulationResults()
	assert.NoError(t, err)
	rwsetBytes, err := rwset.GetPubSimulationBytes()
	assert.NoError(t, err)

	tx := getEnv(ccID, nil, rwsetBytes, t)
	b := &common.Block{Data: &common.BlockData{Data: [][]byte{protoutil.MarshalOrPanic(tx)}}, Header: &common.BlockHeader{Number: 2}}

	err = v.Validate(b)
	assert.NoError(t, err)
	assertInvalid(b, t, peer.TxValidationCode_ENDORSEMENT_POLICY_FAILURE)
}

func TestInvokeNOKWritesToLSCC(t *testing.T) {
	ccID := "mycc"

	v, mockQE, _, _ := setupValidator()

	mockQE.On("GetStateMetadata", ccID, "key").Return(nil, nil)
	mockQE.On("GetState", "lscc", ccID).Return(protoutil.MarshalOrPanic(&ccp.ChaincodeData{
		Name:    ccID,
		Version: ccVersion,
		Vscc:    "vscc",
		Policy:  signedByAnyMember([]string{"SampleOrg"}),
	}), nil)

	tx := getEnv(ccID, nil, createRWset(t, ccID, "lscc"), t)
	b := &common.Block{Data: &common.BlockData{Data: [][]byte{protoutil.MarshalOrPanic(tx)}}, Header: &common.BlockHeader{Number: 2}}

	err := v.Validate(b)
	assert.NoError(t, err)
	assertInvalid(b, t, peer.TxValidationCode_INVALID_CHAINCODE)
}

func TestInvokeNOKWritesToESCC(t *testing.T) {
	ccID := "mycc"

	v, mockQE, _, _ := setupValidator()

	mockQE.On("GetStateMetadata", ccID, "key").Return(nil, nil)
	mockQE.On("GetState", "lscc", ccID).Return(protoutil.MarshalOrPanic(&ccp.ChaincodeData{
		Name:    ccID,
		Version: ccVersion,
		Vscc:    "vscc",
		Policy:  signedByAnyMember([]string{"SampleOrg"}),
	}), nil)

	tx := getEnv(ccID, nil, createRWset(t, ccID, "escc"), t)
	b := &common.Block{
		Data:   &common.BlockData{Data: [][]byte{protoutil.MarshalOrPanic(tx)}},
		Header: &common.BlockHeader{Number: 35},
	}

	err := v.Validate(b)
	assert.NoError(t, err)
	assertInvalid(b, t, peer.TxValidationCode_INVALID_CHAINCODE)
}

func TestInvokeNOKWritesToNotExt(t *testing.T) {
	ccID := "mycc"

	v, mockQE, _, _ := setupValidator()

	mockQE.On("GetStateMetadata", ccID, "key").Return(nil, nil)
	mockQE.On("GetState", "lscc", ccID).Return(protoutil.MarshalOrPanic(&ccp.ChaincodeData{
		Name:    ccID,
		Version: ccVersion,
		Vscc:    "vscc",
		Policy:  signedByAnyMember([]string{"SampleOrg"}),
	}), nil)
	mockQE.On("GetState", "lscc", "notext").Return(nil, nil)

	tx := getEnv(ccID, nil, createRWset(t, ccID, "notext"), t)
	b := &common.Block{
		Data:   &common.BlockData{Data: [][]byte{protoutil.MarshalOrPanic(tx)}},
		Header: &common.BlockHeader{Number: 35},
	}

	err := v.Validate(b)
	assert.NoError(t, err)
	assertInvalid(b, t, peer.TxValidationCode_INVALID_CHAINCODE)
}

func TestInvokeNOKInvokesNotExt(t *testing.T) {
	ccID := "notext"

	v, mockQE, _, _ := setupValidator()

	mockQE.On("GetState", "lscc", "notext").Return(nil, nil)

	tx := getEnv(ccID, nil, createRWset(t, ccID), t)
	b := &common.Block{
		Data:   &common.BlockData{Data: [][]byte{protoutil.MarshalOrPanic(tx)}},
		Header: &common.BlockHeader{},
	}

	err := v.Validate(b)
	assert.NoError(t, err)
	assertInvalid(b, t, peer.TxValidationCode_INVALID_CHAINCODE)
}

func TestInvokeNOKInvokesEmptyCCName(t *testing.T) {
	ccID := ""

	v, _, _, _ := setupValidator()

	tx := getEnv(ccID, nil, createRWset(t, ccID), t)
	b := &common.Block{
		Data:   &common.BlockData{Data: [][]byte{protoutil.MarshalOrPanic(tx)}},
		Header: &common.BlockHeader{},
	}

	err := v.Validate(b)
	assert.NoError(t, err)
	assertInvalid(b, t, peer.TxValidationCode_INVALID_CHAINCODE)
}

func TestInvokeNOKBogusActions(t *testing.T) {
	ccID := "ccid"

	v, _, _, _ := setupValidator()

	tx := getEnv(ccID, nil, []byte("barf"), t)
	b := &common.Block{
		Data:   &common.BlockData{Data: [][]byte{protoutil.MarshalOrPanic(tx)}},
		Header: &common.BlockHeader{},
	}

	err := v.Validate(b)
	assert.NoError(t, err)
	assertInvalid(b, t, peer.TxValidationCode_BAD_RWSET)
}

func TestInvokeNOKCCDoesntExist(t *testing.T) {
	ccID := "mycc"

	v, mockQE, _, _ := setupValidator()
	mockQE.On("GetState", "lscc", ccID).Return(nil, nil)

	tx := getEnv(ccID, nil, createRWset(t, ccID), t)
	b := &common.Block{
		Data:   &common.BlockData{Data: [][]byte{protoutil.MarshalOrPanic(tx)}},
		Header: &common.BlockHeader{},
	}

	err := v.Validate(b)
	assert.NoError(t, err)
	assertInvalid(b, t, peer.TxValidationCode_INVALID_CHAINCODE)
}

func TestInvokeNOKVSCCUnspecified(t *testing.T) {
	ccID := "mycc"

	v, mockQE, _, _ := setupValidator()

	mockQE.On("GetState", "lscc", ccID).Return(protoutil.MarshalOrPanic(&ccp.ChaincodeData{
		Name:    ccID,
		Version: ccVersion,
		Vscc:    "",
		Policy:  signedByAnyMember([]string{"SampleOrg"}),
	}), nil)

	tx := getEnv(ccID, nil, createRWset(t, ccID), t)
	b := &common.Block{
		Data:   &common.BlockData{Data: [][]byte{protoutil.MarshalOrPanic(tx)}},
		Header: &common.BlockHeader{},
	}

	err := v.Validate(b)
	assert.NoError(t, err)
	assertInvalid(b, t, peer.TxValidationCode_INVALID_CHAINCODE)
}

func TestInvokeNoBlock(t *testing.T) {
	v, _, _, _ := setupValidator()
	err := v.Validate(&common.Block{
		Data:   &common.BlockData{Data: [][]byte{}},
		Header: &common.BlockHeader{},
	})
	assert.NoError(t, err)
}

func TestValidateTxWithStateBasedEndorsement(t *testing.T) {
	ccID := "mycc"

	v, mockQE, _, _ := setupValidator()

	mockQE.On("GetState", "lscc", ccID).Return(protoutil.MarshalOrPanic(&ccp.ChaincodeData{
		Name:    ccID,
		Version: ccVersion,
		Vscc:    "vscc",
		Policy:  signedByAnyMember([]string{"SampleOrg"}),
	}), nil)
	mockQE.On("GetStateMetadata", ccID, "key").Return(map[string][]byte{peer.MetaDataKeys_VALIDATION_PARAMETER.String(): protoutil.MarshalOrPanic(&pb.ApplicationPolicy{Type: &pb.ApplicationPolicy_SignaturePolicy{SignaturePolicy: cauthdsl.RejectAllPolicy}})}, nil)

	tx := getEnv(ccID, nil, createRWset(t, ccID), t)
	b := &common.Block{Data: &common.BlockData{Data: [][]byte{protoutil.MarshalOrPanic(tx)}}, Header: &common.BlockHeader{Number: 3}}

	err := v.Validate(b)
	assert.NoError(t, err)
	assertInvalid(b, t, peer.TxValidationCode_ENDORSEMENT_POLICY_FAILURE)
}

func TestTokenValidTransaction(t *testing.T) {
	v, _, _, _ := setupValidator()
	v.ChannelResources.(*mocktxvalidator.Support).ACVal = fabTokenCapabilities()

	tx := getTokenTx(t)
	b := &common.Block{Data: &common.BlockData{Data: [][]byte{protoutil.MarshalOrPanic(tx)}}, Header: &common.BlockHeader{Number: 1}}

	err := v.Validate(b)
	assert.NoError(t, err)
	assertValid(b, t)
}

func TestTokenCapabilityNotEnabled(t *testing.T) {
	v, _, _, _ := setupValidator()

	tx := getTokenTx(t)
	b := &common.Block{Data: &common.BlockData{Data: [][]byte{protoutil.MarshalOrPanic(tx)}}, Header: &common.BlockHeader{Number: 1}}

	err := v.Validate(b)

	assertion := assert.New(t)
	
	assertion.NoError(err)

	
	txsfltr := lutils.TxValidationFlags(b.Metadata.Metadata[common.BlockMetadataIndex_TRANSACTIONS_FILTER])
	assertion.True(txsfltr.IsInvalid(0))
	assertion.True(txsfltr.Flag(0) == peer.TxValidationCode_UNKNOWN_TX_TYPE)
}

func TestTokenDuplicateTxId(t *testing.T) {
	v, _, _, _ := setupValidator()
	v.ChannelResources.(*mocktxvalidator.Support).ACVal = fabTokenCapabilities()

	mockLedger := &mocks3.LedgerResources{}
	v.LedgerResources = mockLedger
	mockLedger.On("GetTransactionByID", mock.Anything).Return(&peer.ProcessedTransaction{}, nil)

	tx := getTokenTx(t)

	b := testutil.NewBlock([]*common.Envelope{tx}, 0, nil)

	err := v.Validate(b)

	assertion := assert.New(t)
	
	assertion.NoError(err)

	
	txsfltr := lutils.TxValidationFlags(b.Metadata.Metadata[common.BlockMetadataIndex_TRANSACTIONS_FILTER])
	assertion.True(txsfltr.IsInvalid(0))
	assertion.True(txsfltr.Flag(0) == peer.TxValidationCode_DUPLICATE_TXID)
}

func TestDynamicCapabilitiesAndMSP(t *testing.T) {
	ccID := "mycc"

	v, mockQE, _, _ := setupValidator()

	mockQE.On("GetState", "lscc", ccID).Return(protoutil.MarshalOrPanic(&ccp.ChaincodeData{
		Name:    ccID,
		Version: ccVersion,
		Vscc:    "vscc",
		Policy:  signedByAnyMember([]string{"SampleOrg"}),
	}), nil)
	mockQE.On("GetStateMetadata", ccID, "key").Return(nil, nil)

	tx := getEnv(ccID, nil, createRWset(t, ccID), t)
	b := &common.Block{
		Data:   &common.BlockData{Data: [][]byte{protoutil.MarshalOrPanic(tx)}},
		Header: &common.BlockHeader{Number: 1},
	}

	
	err := v.Validate(b)
	assert.NoError(t, err)
	assertValid(b, t)
	
	capabilityInvokeCount := v.ChannelResources.(*mocktxvalidator.Support).CapabilitiesInvokeCount()
	mspManagerInvokeCount := v.ChannelResources.(*mocktxvalidator.Support).MSPManagerInvokeCount()

	
	err = v.Validate(b)
	assert.NoError(t, err)
	assertValid(b, t)

	
	
	assert.Equal(t, 2*capabilityInvokeCount, v.ChannelResources.(*mocktxvalidator.Support).CapabilitiesInvokeCount())
	
	
	assert.Equal(t, 2*mspManagerInvokeCount, v.ChannelResources.(*mocktxvalidator.Support).MSPManagerInvokeCount())
}








func TestLedgerIsNotAvailable(t *testing.T) {
	ccID := "mycc"

	v, mockQE, _, _ := setupValidator()
	mockQE.On("GetState", "lscc", ccID).Return(nil, errors.New("Detroit rock city"))

	tx := getEnv(ccID, nil, createRWset(t, ccID), t)

	b := &common.Block{
		Data:   &common.BlockData{Data: [][]byte{protoutil.MarshalOrPanic(tx)}},
		Header: &common.BlockHeader{},
	}

	err := v.Validate(b)

	assertion := assert.New(t)
	
	assertion.Error(err)
	
	assertion.NotNil(err.(*commonerrors.VSCCInfoLookupFailureError))
}

func TestLedgerIsNotAvailableForCheckingTxidDuplicate(t *testing.T) {
	ccID := "mycc"

	v, _, _, _ := setupValidator()

	tx := getEnv(ccID, nil, createRWset(t, ccID), t)

	mockLedger := &mocks3.LedgerResources{}
	v.LedgerResources = mockLedger
	mockLedger.On("GetTransactionByID", mock.Anything).Return(nil, errors.New("uh, oh"))

	b := &common.Block{
		Data:   &common.BlockData{Data: [][]byte{protoutil.MarshalOrPanic(tx)}},
		Header: &common.BlockHeader{Number: 1},
	}

	err := v.Validate(b)

	assertion := assert.New(t)
	
	assertion.Error(err)
}

func TestDuplicateTxId(t *testing.T) {
	ccID := "mycc"

	v, _, _, _ := setupValidator()

	mockLedger := &mocks3.LedgerResources{}
	v.LedgerResources = mockLedger
	mockLedger.On("GetTransactionByID", mock.Anything).Return(&peer.ProcessedTransaction{}, nil)

	tx := getEnv(ccID, nil, createRWset(t, ccID), t)

	b := &common.Block{
		Data:   &common.BlockData{Data: [][]byte{protoutil.MarshalOrPanic(tx)}},
		Header: &common.BlockHeader{},
	}

	err := v.Validate(b)

	assertion := assert.New(t)
	
	assertion.NoError(err)

	
	txsfltr := lutils.TxValidationFlags(b.Metadata.Metadata[common.BlockMetadataIndex_TRANSACTIONS_FILTER])
	assertion.True(txsfltr.IsInvalid(0))
	assertion.True(txsfltr.Flag(0) == peer.TxValidationCode_DUPLICATE_TXID)
}

func TestValidationInvalidEndorsing(t *testing.T) {
	ccID := "mycc"

	mspmgr := &mocks2.MSPManager{}
	mockID := &mocks2.Identity{}
	mockID.SatisfiesPrincipalReturns(nil)
	mockID.GetIdentifierReturns(&msp.IdentityIdentifier{})
	mspmgr.DeserializeIdentityReturns(mockID, nil)

	pm := &mocks.Mapper{}
	factory := &mocks.PluginFactory{}
	pm.On("FactoryByName", vp.Name("vscc")).Return(factory)
	plugin := &mocks.Plugin{}
	factory.On("New").Return(plugin)
	plugin.On("Init", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	plugin.On("Validate", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("invalid tx"))

	mockQE := &mocks3.QueryExecutor{}
	mockQE.On("Done").Return(nil)

	mockLedger := &mocks3.LedgerResources{}
	mockLedger.On("GetTransactionByID", mock.Anything).Return(nil, ledger.NotFoundInIndexErr("As idle as a painted ship upon a painted ocean"))
	mockLedger.On("NewQueryExecutor").Return(mockQE, nil)

	mockCpmg := &mocks.ChannelPolicyManagerGetter{}
	mockCpmg.On("Manager", mock.Anything).Return(nil, true)

	v := txvalidatorv20.NewTxValidator(
		"",
		semaphore.New(10),
		&mocktxvalidator.Support{ACVal: v20Capabilities(), MSPManagerVal: mspmgr},
		mockLedger,
		&lscc.LifeCycleSysCC{},
		&mocks3.CollectionResources{},
		pm,
		mockCpmg,
	)

	tx := getEnv(ccID, nil, createRWset(t, ccID), t)

	cd := &ccp.ChaincodeData{
		Name:    ccID,
		Version: ccVersion,
		Vscc:    "vscc",
		Policy:  signedByAnyMember([]string{"SampleOrg"}),
	}

	cdbytes := protoutil.MarshalOrPanic(cd)

	mockQE.On("GetState", "lscc", ccID).Return(cdbytes, nil)

	b := &common.Block{
		Data:   &common.BlockData{Data: [][]byte{protoutil.MarshalOrPanic(tx)}},
		Header: &common.BlockHeader{},
	}

	
	err := v.Validate(b)
	
	assert.NoError(t, err)
	assertInvalid(b, t, peer.TxValidationCode_ENDORSEMENT_POLICY_FAILURE)
}

func TestValidationPluginExecutionError(t *testing.T) {
	ccID := "mycc"

	mspmgr := &mocks2.MSPManager{}
	mockID := &mocks2.Identity{}
	mockID.SatisfiesPrincipalReturns(nil)
	mockID.GetIdentifierReturns(&msp.IdentityIdentifier{})
	mspmgr.DeserializeIdentityReturns(mockID, nil)

	pm := &mocks.Mapper{}
	factory := &mocks.PluginFactory{}
	pm.On("FactoryByName", vp.Name("vscc")).Return(factory)
	plugin := &mocks.Plugin{}
	factory.On("New").Return(plugin)
	plugin.On("Init", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	plugin.On("Validate", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&validation.ExecutionFailureError{
		Reason: "I/O error",
	})

	mockQE := &mocks3.QueryExecutor{}
	mockQE.On("Done").Return(nil)
	mockQE.On("GetState", "lscc", ccID).Return(protoutil.MarshalOrPanic(&ccp.ChaincodeData{
		Name:    ccID,
		Version: ccVersion,
		Vscc:    "vscc",
		Policy:  signedByAnyMember([]string{"SampleOrg"}),
	}), nil)

	mockLedger := &mocks3.LedgerResources{}
	mockLedger.On("GetTransactionByID", mock.Anything).Return(nil, ledger.NotFoundInIndexErr("As idle as a painted ship upon a painted ocean"))
	mockLedger.On("NewQueryExecutor").Return(mockQE, nil)

	mockCpmg := &mocks.ChannelPolicyManagerGetter{}
	mockCpmg.On("Manager", mock.Anything).Return(nil, true)

	v := txvalidatorv20.NewTxValidator(
		"",
		semaphore.New(10),
		&mocktxvalidator.Support{ACVal: v20Capabilities(), MSPManagerVal: mspmgr},
		mockLedger,
		&lscc.LifeCycleSysCC{},
		&mocks3.CollectionResources{},
		pm,
		mockCpmg,
	)

	tx := getEnv(ccID, nil, createRWset(t, ccID), t)
	b := &common.Block{
		Data:   &common.BlockData{Data: [][]byte{protoutil.MarshalOrPanic(tx)}},
		Header: &common.BlockHeader{},
	}

	err := v.Validate(b)
	executionErr := err.(*commonerrors.VSCCExecutionFailureError)
	assert.Contains(t, executionErr.Error(), "I/O error")
}

func TestValidationPluginNotFound(t *testing.T) {
	ccID := "mycc"

	mspmgr := &mocks2.MSPManager{}
	mockID := &mocks2.Identity{}
	mockID.SatisfiesPrincipalReturns(nil)
	mockID.GetIdentifierReturns(&msp.IdentityIdentifier{})
	mspmgr.DeserializeIdentityReturns(mockID, nil)

	pm := &mocks.Mapper{}
	pm.On("FactoryByName", vp.Name("vscc")).Return(nil)

	mockQE := &mocks3.QueryExecutor{}
	mockQE.On("Done").Return(nil)
	mockQE.On("GetState", "lscc", ccID).Return(protoutil.MarshalOrPanic(&ccp.ChaincodeData{
		Name:    ccID,
		Version: ccVersion,
		Vscc:    "vscc",
		Policy:  signedByAnyMember([]string{"SampleOrg"}),
	}), nil)

	mockLedger := &mocks3.LedgerResources{}
	mockLedger.On("GetTransactionByID", mock.Anything).Return(nil, ledger.NotFoundInIndexErr("As idle as a painted ship upon a painted ocean"))
	mockLedger.On("NewQueryExecutor").Return(mockQE, nil)

	mockCpmg := &mocks.ChannelPolicyManagerGetter{}
	mockCpmg.On("Manager", mock.Anything).Return(nil, true)

	v := txvalidatorv20.NewTxValidator(
		"",
		semaphore.New(10),
		&mocktxvalidator.Support{ACVal: v20Capabilities(), MSPManagerVal: mspmgr},
		mockLedger,
		&lscc.LifeCycleSysCC{},
		&mocks3.CollectionResources{},
		pm,
		mockCpmg,
	)

	tx := getEnv(ccID, nil, createRWset(t, ccID), t)
	b := &common.Block{
		Data:   &common.BlockData{Data: [][]byte{protoutil.MarshalOrPanic(tx)}},
		Header: &common.BlockHeader{},
	}

	err := v.Validate(b)
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

func ToHex(q uint64) string {
	return "0x" + strconv.FormatUint(q, 16)
}
