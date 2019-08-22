/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package v20

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/mcc-github/blockchain/common/cauthdsl"
	commonerrors "github.com/mcc-github/blockchain/common/errors"
	mc "github.com/mcc-github/blockchain/common/mocks/config"
	lm "github.com/mcc-github/blockchain/common/mocks/ledger"
	"github.com/mcc-github/blockchain/common/util"
	"github.com/mcc-github/blockchain/core/committer/txvalidator/v14"
	mocks2 "github.com/mcc-github/blockchain/core/committer/txvalidator/v14/mocks"
	"github.com/mcc-github/blockchain/core/common/ccprovider"
	validation "github.com/mcc-github/blockchain/core/handlers/validation/api/capabilities"
	"github.com/mcc-github/blockchain/core/handlers/validation/builtin/v20/mocks"
	"github.com/mcc-github/blockchain/msp"
	mspmgmt "github.com/mcc-github/blockchain/msp/mgmt"
	msptesttools "github.com/mcc-github/blockchain/msp/mgmt/testtools"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/peer"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func createTx(endorsedByDuplicatedIdentity bool) (*common.Envelope, error) {
	ccid := &peer.ChaincodeID{Name: "foo", Version: "v1"}
	cis := &peer.ChaincodeInvocationSpec{ChaincodeSpec: &peer.ChaincodeSpec{ChaincodeId: ccid}}

	prop, _, err := protoutil.CreateProposalFromCIS(common.HeaderType_ENDORSER_TRANSACTION, util.GetTestChainID(), cis, sid)
	if err != nil {
		return nil, err
	}

	presp, err := protoutil.CreateProposalResponse(prop.Header, prop.Payload, &peer.Response{Status: 200}, []byte("res"), nil, ccid, id)
	if err != nil {
		return nil, err
	}

	var env *common.Envelope
	if endorsedByDuplicatedIdentity {
		env, err = protoutil.CreateSignedTx(prop, id, presp, presp)
	} else {
		env, err = protoutil.CreateSignedTx(prop, id, presp)
	}
	if err != nil {
		return nil, err
	}
	return env, err
}

func getSignedByMSPMemberPolicy(mspID string) ([]byte, error) {
	p := cauthdsl.SignedByMspMember(mspID)

	b, err := protoutil.Marshal(p)
	if err != nil {
		return nil, fmt.Errorf("Could not marshal policy, err %s", err)
	}

	return b, err
}

func newValidationInstance(state map[string]map[string][]byte) *Validator {
	qec := &mocks2.QueryExecutorCreator{}
	qec.On("NewQueryExecutor").Return(lm.NewMockQueryExecutor(state), nil)
	return newCustomValidationInstance(qec, &mc.MockApplicationCapabilities{})
}

func newCustomValidationInstance(qec txvalidator.QueryExecutorCreator, c validation.Capabilities) *Validator {
	sbvm := &mocks.StateBasedValidator{}
	sbvm.On("PreValidate", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	sbvm.On("PostValidate", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	sbvm.On("Validate", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	mockCR := &mocks.CollectionResources{}
	mockCR.On("CollectionValidationInfo", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil, nil)

	sf := &txvalidator.StateFetcherImpl{QueryExecutorCreator: qec}
	is := &mocks.IdentityDeserializer{}
	pe := &txvalidator.PolicyEvaluator{
		IdentityDeserializer: mspmgmt.GetManagerForChain(util.GetTestChainID()),
	}
	v := New(c, sf, is, pe, mockCR)

	v.stateBasedValidator = sbvm
	return v
}

func TestStateBasedValidationFailure(t *testing.T) {
	qec := &mocks2.QueryExecutorCreator{}
	qec.On("NewQueryExecutor").Return(lm.NewMockQueryExecutor(make(map[string]map[string][]byte)), nil)

	sbvm := &mocks.StateBasedValidator{}
	sbvm.On("PreValidate", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	sbvm.On("PostValidate", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	mockCR := &mocks.CollectionResources{}
	mockCR.On("CollectionValidationInfo", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil, nil)

	sf := &txvalidator.StateFetcherImpl{QueryExecutorCreator: qec}
	is := &mocks.IdentityDeserializer{}
	pe := &txvalidator.PolicyEvaluator{
		IdentityDeserializer: mspmgmt.GetManagerForChain(util.GetTestChainID()),
	}
	v := New(&mc.MockApplicationCapabilities{}, sf, is, pe, mockCR)
	v.stateBasedValidator = sbvm

	tx, err := createTx(false)
	if err != nil {
		t.Fatalf("createTx returned err %s", err)
	}

	envBytes, err := protoutil.GetBytesEnvelope(tx)
	if err != nil {
		t.Fatalf("GetBytesEnvelope returned err %s", err)
	}

	policy, err := getSignedByMSPMemberPolicy(mspid)
	if err != nil {
		t.Fatalf("failed getting policy, err %s", err)
	}

	b := &common.Block{Data: &common.BlockData{Data: [][]byte{envBytes}}, Header: &common.BlockHeader{}}

	
	sbvm.On("Validate", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&commonerrors.VSCCEndorsementPolicyError{Err: fmt.Errorf("some sbe validation err")}).Once()
	err = v.Validate(b, "foo", 0, 0, policy)
	assert.Error(t, err)
	assert.IsType(t, &commonerrors.VSCCEndorsementPolicyError{}, err)

	
	sbvm.On("Validate", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&commonerrors.VSCCExecutionFailureError{Err: fmt.Errorf("some sbe validation err")}).Once()
	err = v.Validate(b, "foo", 0, 0, policy)
	assert.Error(t, err)
	assert.IsType(t, &commonerrors.VSCCExecutionFailureError{}, err)

	
	sbvm.On("Validate", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
	err = v.Validate(b, "foo", 0, 0, policy)
	assert.NoError(t, err)
}

func TestInvoke(t *testing.T) {
	v := newValidationInstance(make(map[string]map[string][]byte))

	
	var err error
	b := &common.Block{Data: &common.BlockData{Data: [][]byte{[]byte("a")}}, Header: &common.BlockHeader{}}
	err = v.Validate(b, "foo", 0, 0, []byte("a"))
	assert.Error(t, err)

	
	b = &common.Block{Data: &common.BlockData{Data: [][]byte{protoutil.MarshalOrPanic(&common.Envelope{Payload: []byte("barf")})}}, Header: &common.BlockHeader{}}
	err = v.Validate(b, "foo", 0, 0, []byte("a"))
	assert.Error(t, err)

	
	e := protoutil.MarshalOrPanic(&common.Envelope{Payload: protoutil.MarshalOrPanic(&common.Payload{Header: &common.Header{ChannelHeader: []byte("barf")}})})
	b = &common.Block{Data: &common.BlockData{Data: [][]byte{e}}, Header: &common.BlockHeader{}}
	err = v.Validate(b, "foo", 0, 0, []byte("a"))
	assert.Error(t, err)

	tx, err := createTx(false)
	if err != nil {
		t.Fatalf("createTx returned err %s", err)
	}

	envBytes, err := protoutil.GetBytesEnvelope(tx)
	if err != nil {
		t.Fatalf("GetBytesEnvelope returned err %s", err)
	}

	policy, err := getSignedByMSPMemberPolicy(mspid)
	if err != nil {
		t.Fatalf("failed getting policy, err %s", err)
	}

	
	e = protoutil.MarshalOrPanic(&common.Envelope{Payload: protoutil.MarshalOrPanic(&common.Payload{Header: &common.Header{ChannelHeader: protoutil.MarshalOrPanic(&common.ChannelHeader{Type: int32(common.HeaderType_ORDERER_TRANSACTION)})}})})
	b = &common.Block{Data: &common.BlockData{Data: [][]byte{e}}, Header: &common.BlockHeader{}}
	err = v.Validate(b, "foo", 0, 0, policy)
	assert.Error(t, err)

	
	e = protoutil.MarshalOrPanic(&common.Envelope{Payload: protoutil.MarshalOrPanic(&common.Payload{Header: &common.Header{ChannelHeader: protoutil.MarshalOrPanic(&common.ChannelHeader{Type: int32(common.HeaderType_ORDERER_TRANSACTION)})}})})
	b = &common.Block{Data: &common.BlockData{Data: [][]byte{e}}, Header: &common.BlockHeader{}}
	err = v.Validate(b, "foo", 0, 0, policy)
	assert.Error(t, err)

	
	b = &common.Block{Data: &common.BlockData{Data: [][]byte{envBytes}}, Header: &common.BlockHeader{}}
	err = v.Validate(b, "foo", 0, 0, policy)
	assert.NoError(t, err)
}

func TestToApplicationPolicyTranslator_Translate(t *testing.T) {
	tr := &toApplicationPolicyTranslator{}
	res, err := tr.Translate(nil)
	assert.NoError(t, err)
	assert.Nil(t, res)

	res, err = tr.Translate([]byte("barf"))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "could not unmarshal signature policy envelope: unexpected EOF")
	assert.Nil(t, res)

	res, err = tr.Translate(protoutil.MarshalOrPanic(cauthdsl.SignedByMspMember("the right honourable member for Ipswich")))
	assert.NoError(t, err)
	assert.Equal(t, res, protoutil.MarshalOrPanic(&peer.ApplicationPolicy{
		Type: &peer.ApplicationPolicy_SignaturePolicy{
			SignaturePolicy: cauthdsl.SignedByMspMember("the right honourable member for Ipswich"),
		},
	}))
}

var id msp.SigningIdentity
var sid []byte
var mspid string
var chainId string = util.GetTestChainID()

type mockPolicyChecker struct{}

func (c *mockPolicyChecker) CheckPolicy(channelID, policyName string, signedProp *peer.SignedProposal) error {
	return nil
}

func (c *mockPolicyChecker) CheckPolicyBySignedData(channelID, policyName string, sd []*protoutil.SignedData) error {
	return nil
}

func (c *mockPolicyChecker) CheckPolicyNoChannel(policyName string, signedProp *peer.SignedProposal) error {
	return nil
}

func TestMain(m *testing.M) {
	testDir, err := ioutil.TempDir("", "v1.3-validation")
	if err != nil {
		fmt.Printf("Could not create temp dir: %s", err)
		os.Exit(-1)
	}
	defer os.RemoveAll(testDir)
	ccprovider.SetChaincodesPath(testDir)

	
	msptesttools.LoadMSPSetupForTesting()

	id, err = mspmgmt.GetLocalMSP().GetDefaultSigningIdentity()
	if err != nil {
		fmt.Printf("GetSigningIdentity failed with err %s", err)
		os.Exit(-1)
	}

	sid, err = id.Serialize()
	if err != nil {
		fmt.Printf("Serialize failed with err %s", err)
		os.Exit(-1)
	}

	
	var msp msp.MSP
	mspMgr := mspmgmt.GetManagerForChain(chainId)
	msps, err := mspMgr.GetMSPs()
	if err != nil {
		fmt.Printf("Could not retrieve the MSPs for the chain manager, err %s", err)
		os.Exit(-1)
	}
	if len(msps) == 0 {
		fmt.Printf("At least one MSP was expected")
		os.Exit(-1)
	}
	for _, m := range msps {
		msp = m
		break
	}
	mspid, err = msp.GetIdentifier()
	if err != nil {
		fmt.Printf("Failure getting the msp identifier, err %s", err)
		os.Exit(-1)
	}

	
	mspmgmt.XXXSetMSPManager("mycc", mspmgmt.GetManagerForChain(util.GetTestChainID()))

	os.Exit(m.Run())
}
