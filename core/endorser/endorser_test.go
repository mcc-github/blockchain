/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package endorser_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/common/metrics/disabled"
	"github.com/mcc-github/blockchain/common/metrics/metricsfakes"
	mc "github.com/mcc-github/blockchain/common/mocks/config"
	resourceconfig "github.com/mcc-github/blockchain/common/mocks/resourcesconfig"
	"github.com/mcc-github/blockchain/common/util"
	"github.com/mcc-github/blockchain/core/common/ccprovider"
	"github.com/mcc-github/blockchain/core/endorser"
	"github.com/mcc-github/blockchain/core/endorser/mocks"
	"github.com/mcc-github/blockchain/core/handlers/endorsement/builtin"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/msp"
	"github.com/mcc-github/blockchain/msp/mgmt"
	msptesttools "github.com/mcc-github/blockchain/msp/mgmt/testtools"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/ledger/rwset"
	pb "github.com/mcc-github/blockchain/protos/peer"
	"github.com/mcc-github/blockchain/protos/transientstore"
	"github.com/mcc-github/blockchain/protoutil"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func pvtEmptyDistributor(_ string, _ string, _ *transientstore.TxPvtReadWriteSetWithConfigInfo, _ uint64) error {
	return nil
}

func getSignedPropWithCHID(ccid, ccver, chid string, t *testing.T) *pb.SignedProposal {
	ccargs := [][]byte{[]byte("args")}

	return getSignedPropWithCHIdAndArgs(chid, ccid, ccver, ccargs, t)
}

func getSignedProp(ccid, ccver string, t *testing.T) *pb.SignedProposal {
	return getSignedPropWithCHID(ccid, ccver, util.GetTestChainID(), t)
}

func getSignedPropWithCHIdAndArgs(chid, ccid, ccver string, ccargs [][]byte, t *testing.T) *pb.SignedProposal {
	spec := &pb.ChaincodeSpec{Type: 1, ChaincodeId: &pb.ChaincodeID{Name: ccid, Version: ccver}, Input: &pb.ChaincodeInput{Args: ccargs}}

	cis := &pb.ChaincodeInvocationSpec{ChaincodeSpec: spec}

	creator, err := signer.Serialize()
	assert.NoError(t, err)
	prop, _, err := protoutil.CreateChaincodeProposal(common.HeaderType_ENDORSER_TRANSACTION, chid, cis, creator)
	assert.NoError(t, err)
	propBytes, err := protoutil.GetBytesProposal(prop)
	assert.NoError(t, err)
	signature, err := signer.Sign(propBytes)
	assert.NoError(t, err)
	return &pb.SignedProposal{ProposalBytes: propBytes, Signature: signature}
}

func newMockTxSim() *mocks.MockTxSim {
	return &mocks.MockTxSim{
		GetTxSimulationResultsRv: &ledger.TxSimulationResults{
			PubSimulationResults: &rwset.TxReadWriteSet{},
		},
	}
}


type fakeEndorserMetrics struct {
	proposalDuration         *metricsfakes.Histogram
	proposalsReceived        *metricsfakes.Counter
	successfulProposals      *metricsfakes.Counter
	proposalValidationFailed *metricsfakes.Counter
	proposalACLCheckFailed   *metricsfakes.Counter
	initFailed               *metricsfakes.Counter
	endorsementsFailed       *metricsfakes.Counter
	duplicateTxsFailure      *metricsfakes.Counter
}


func initFakeMetrics(es *endorser.Endorser) *fakeEndorserMetrics {
	fakeMetrics := &fakeEndorserMetrics{
		proposalDuration:         &metricsfakes.Histogram{},
		proposalsReceived:        &metricsfakes.Counter{},
		successfulProposals:      &metricsfakes.Counter{},
		proposalValidationFailed: &metricsfakes.Counter{},
		proposalACLCheckFailed:   &metricsfakes.Counter{},
		initFailed:               &metricsfakes.Counter{},
		endorsementsFailed:       &metricsfakes.Counter{},
		duplicateTxsFailure:      &metricsfakes.Counter{},
	}

	fakeMetrics.proposalDuration.WithReturns(fakeMetrics.proposalDuration)
	fakeMetrics.proposalACLCheckFailed.WithReturns(fakeMetrics.proposalACLCheckFailed)
	fakeMetrics.initFailed.WithReturns(fakeMetrics.initFailed)
	fakeMetrics.endorsementsFailed.WithReturns(fakeMetrics.endorsementsFailed)
	fakeMetrics.duplicateTxsFailure.WithReturns(fakeMetrics.duplicateTxsFailure)

	es.Metrics.ProposalDuration = fakeMetrics.proposalDuration
	es.Metrics.ProposalsReceived = fakeMetrics.proposalsReceived
	es.Metrics.SuccessfulProposals = fakeMetrics.successfulProposals
	es.Metrics.ProposalValidationFailed = fakeMetrics.proposalValidationFailed
	es.Metrics.ProposalACLCheckFailed = fakeMetrics.proposalACLCheckFailed
	es.Metrics.InitFailed = fakeMetrics.initFailed
	es.Metrics.EndorsementsFailed = fakeMetrics.endorsementsFailed
	es.Metrics.DuplicateTxsFailure = fakeMetrics.duplicateTxsFailure

	return fakeMetrics
}

func testEndorsementCompletedMetric(t *testing.T, fakeMetrics *fakeEndorserMetrics, callCount int32, chainID, ccnamever, succ string) {
	
	assert.EqualValues(t, callCount, fakeMetrics.proposalDuration.WithCallCount())
	labelValues := fakeMetrics.proposalDuration.WithArgsForCall(0)
	assert.EqualValues(t, labelValues, []string{"channel", chainID, "chaincode", ccnamever, "success", succ})
	assert.NotEqual(t, 0, fakeMetrics.proposalDuration.ObserveArgsForCall(0))
}

func TestEndorserNilProp(t *testing.T) {
	es := endorser.NewEndorserServer(pvtEmptyDistributor, &mocks.MockSupport{
		GetApplicationConfigBoolRv: true,
		GetApplicationConfigRv:     &mc.MockApplication{CapabilitiesRv: &mc.MockApplicationCapabilities{}},
		GetTransactionByIDErr:      errors.New(""),
		ChaincodeDefinitionRv:      &ccprovider.ChaincodeData{Escc: "ESCC"},
		ExecuteResp:                &pb.Response{Status: 200, Payload: protoutil.MarshalOrPanic(&pb.ProposalResponse{Response: &pb.Response{}})},
		GetTxSimulatorRv: &mocks.MockTxSim{
			GetTxSimulationResultsRv: &ledger.TxSimulationResults{
				PubSimulationResults: &rwset.TxReadWriteSet{},
			},
		},
	}, &disabled.Provider{})

	fakeMetrics := initFakeMetrics(es)

	pResp, err := es.ProcessProposal(context.Background(), nil)
	assert.Error(t, err)
	assert.EqualValues(t, 500, pResp.Response.Status)
	assert.Equal(t, "nil arguments", pResp.Response.Message)

	
	assert.EqualValues(t, 1, fakeMetrics.proposalsReceived.AddCallCount())
	assert.EqualValues(t, 1, fakeMetrics.proposalsReceived.AddArgsForCall(0))
	assert.EqualValues(t, 1, fakeMetrics.proposalValidationFailed.AddCallCount())
	assert.EqualValues(t, 1, fakeMetrics.proposalValidationFailed.AddArgsForCall(0))
}

func TestEndorserCCInvocationFailed(t *testing.T) {
	es := endorser.NewEndorserServer(pvtEmptyDistributor, &mocks.MockSupport{
		GetApplicationConfigBoolRv: true,
		GetApplicationConfigRv:     &mc.MockApplication{CapabilitiesRv: &mc.MockApplicationCapabilities{}},
		GetTransactionByIDErr:      errors.New(""),
		ChaincodeDefinitionRv:      &ccprovider.ChaincodeData{Escc: "ESCC"},
		ExecuteResp:                &pb.Response{Status: 1000, Payload: protoutil.MarshalOrPanic(&pb.ProposalResponse{Response: &pb.Response{}}), Message: "achoo"},
		GetTxSimulatorRv: &mocks.MockTxSim{
			GetTxSimulationResultsRv: &ledger.TxSimulationResults{
				PubSimulationResults: &rwset.TxReadWriteSet{},
			},
		},
	}, &disabled.Provider{})

	signedProp := getSignedProp("test-chaincode", "test-version", t)

	pResp, err := es.ProcessProposal(context.Background(), signedProp)
	assert.NoError(t, err)
	assert.EqualValues(t, 1000, pResp.Response.Status)
	assert.Equal(t, "achoo", pResp.Response.Message)
}

func TestEndorserNoCCDef(t *testing.T) {
	es := endorser.NewEndorserServer(pvtEmptyDistributor, &mocks.MockSupport{
		GetApplicationConfigBoolRv: true,
		GetApplicationConfigRv:     &mc.MockApplication{CapabilitiesRv: &mc.MockApplicationCapabilities{}},
		GetTransactionByIDErr:      errors.New(""),
		ChaincodeDefinitionError:   errors.New("gesundheit"),
		ExecuteResp:                &pb.Response{Status: 200, Payload: protoutil.MarshalOrPanic(&pb.ProposalResponse{Response: &pb.Response{}})},
		GetTxSimulatorRv: &mocks.MockTxSim{
			GetTxSimulationResultsRv: &ledger.TxSimulationResults{
				PubSimulationResults: &rwset.TxReadWriteSet{},
			},
		},
	}, &disabled.Provider{})

	signedProp := getSignedProp("test-chaincode", "test-version", t)

	pResp, err := es.ProcessProposal(context.Background(), signedProp)
	assert.NoError(t, err)
	assert.EqualValues(t, 500, pResp.Response.Status)
	assert.Equal(t, "make sure the chaincode test-chaincode has been successfully defined on channel testchainid and try again: gesundheit", pResp.Response.Message)
}

func TestEndorserSysCC(t *testing.T) {
	m := &mock.Mock{}
	m.On("Sign", mock.Anything).Return([]byte{1, 2, 3, 4, 5}, nil)
	m.On("Serialize").Return([]byte{1, 1, 1}, nil)
	m.On("GetTxSimulator", mock.Anything, mock.Anything).Return(newMockTxSim(), nil)
	support := &mocks.MockSupport{
		Mock:                       m,
		GetApplicationConfigBoolRv: true,
		GetApplicationConfigRv:     &mc.MockApplication{CapabilitiesRv: &mc.MockApplicationCapabilities{}},
		GetTransactionByIDErr:      errors.New(""),
		IsSysCCRv:                  true,
		ChaincodeDefinitionRv:      &ccprovider.ChaincodeData{Escc: "ESCC"},
		ExecuteResp:                &pb.Response{Status: 200, Payload: protoutil.MarshalOrPanic(&pb.ProposalResponse{Response: &pb.Response{}})},
	}
	attachPluginEndorser(support, nil)
	es := endorser.NewEndorserServer(pvtEmptyDistributor, support, &disabled.Provider{})

	signedProp := getSignedProp("test-chaincode", "test-version", t)

	pResp, err := es.ProcessProposal(context.Background(), signedProp)
	assert.NoError(t, err)
	assert.EqualValues(t, 200, pResp.Response.Status)
}

func TestEndorserCCInvocationError(t *testing.T) {
	es := endorser.NewEndorserServer(pvtEmptyDistributor, &mocks.MockSupport{
		GetApplicationConfigBoolRv: true,
		GetApplicationConfigRv:     &mc.MockApplication{CapabilitiesRv: &mc.MockApplicationCapabilities{}},
		GetTransactionByIDErr:      errors.New(""),
		ExecuteError:               errors.New(""),
		ChaincodeDefinitionRv:      &ccprovider.ChaincodeData{Escc: "ESCC"},
		GetTxSimulatorRv: &mocks.MockTxSim{
			GetTxSimulationResultsRv: &ledger.TxSimulationResults{
				PubSimulationResults: &rwset.TxReadWriteSet{},
			},
		},
	}, &disabled.Provider{})

	signedProp := getSignedProp("test-chaincode", "test-version", t)

	pResp, err := es.ProcessProposal(context.Background(), signedProp)
	assert.NoError(t, err)
	assert.EqualValues(t, 500, pResp.Response.Status)
}

func TestEndorserDupTXId(t *testing.T) {
	es := endorser.NewEndorserServer(pvtEmptyDistributor, &mocks.MockSupport{
		GetApplicationConfigBoolRv: true,
		GetApplicationConfigRv:     &mc.MockApplication{CapabilitiesRv: &mc.MockApplicationCapabilities{}},
		ChaincodeDefinitionRv:      &ccprovider.ChaincodeData{Escc: "ESCC"},
		ExecuteResp:                &pb.Response{Status: 200, Payload: protoutil.MarshalOrPanic(&pb.ProposalResponse{Response: &pb.Response{}})},
		GetTxSimulatorRv: &mocks.MockTxSim{
			GetTxSimulationResultsRv: &ledger.TxSimulationResults{
				PubSimulationResults: &rwset.TxReadWriteSet{},
			},
		},
	}, &disabled.Provider{})

	fakeMetrics := initFakeMetrics(es)

	signedProp := getSignedProp("test-chaincode", "test-version", t)

	pResp, err := es.ProcessProposal(context.Background(), signedProp)
	assert.Error(t, err)
	assert.EqualValues(t, 500, pResp.Response.Status)
	assert.Regexp(t, "duplicate transaction found", pResp.Response.Message)

	
	assert.EqualValues(t, 1, fakeMetrics.duplicateTxsFailure.WithCallCount())
	labelValues := fakeMetrics.duplicateTxsFailure.WithArgsForCall(0)
	assert.EqualValues(t, labelValues, []string{"channel", "testchainid", "chaincode", "test-chaincode:test-version"})
	assert.EqualValues(t, 1, fakeMetrics.duplicateTxsFailure.AddCallCount())
	assert.EqualValues(t, 1, fakeMetrics.duplicateTxsFailure.AddArgsForCall(0))
}

func TestEndorserBadACL(t *testing.T) {
	es := endorser.NewEndorserServer(pvtEmptyDistributor, &mocks.MockSupport{
		GetApplicationConfigBoolRv: true,
		GetApplicationConfigRv:     &mc.MockApplication{CapabilitiesRv: &mc.MockApplicationCapabilities{}},
		CheckACLErr:                errors.New(""),
		GetTransactionByIDErr:      errors.New(""),
		ChaincodeDefinitionRv:      &ccprovider.ChaincodeData{Escc: "ESCC"},
		ExecuteResp:                &pb.Response{Status: 200, Payload: protoutil.MarshalOrPanic(&pb.ProposalResponse{Response: &pb.Response{}})},
		GetTxSimulatorRv: &mocks.MockTxSim{
			GetTxSimulationResultsRv: &ledger.TxSimulationResults{
				PubSimulationResults: &rwset.TxReadWriteSet{},
			},
		},
	}, &disabled.Provider{})

	fakeMetrics := initFakeMetrics(es)

	signedProp := getSignedProp("test-chaincode", "test-version", t)

	pResp, err := es.ProcessProposal(context.Background(), signedProp)
	assert.Error(t, err)
	assert.EqualValues(t, 500, pResp.Response.Status)

	
	assert.EqualValues(t, 1, fakeMetrics.proposalACLCheckFailed.WithCallCount())
	labelValues := fakeMetrics.proposalACLCheckFailed.WithArgsForCall(0)
	assert.EqualValues(t, labelValues, []string{"channel", "testchainid", "chaincode", "test-chaincode:test-version"})
	assert.EqualValues(t, 1, fakeMetrics.proposalACLCheckFailed.AddCallCount())
	assert.EqualValues(t, 1, fakeMetrics.proposalACLCheckFailed.AddArgsForCall(0))
}

func TestEndorserGoodPathEmptyChannel(t *testing.T) {
	es := endorser.NewEndorserServer(pvtEmptyDistributor, &mocks.MockSupport{
		GetApplicationConfigBoolRv: true,
		GetApplicationConfigRv:     &mc.MockApplication{CapabilitiesRv: &mc.MockApplicationCapabilities{}},
		GetTransactionByIDErr:      errors.New(""),
		ChaincodeDefinitionRv:      &ccprovider.ChaincodeData{Escc: "ESCC"},
		ExecuteResp:                &pb.Response{Status: 200, Payload: protoutil.MarshalOrPanic(&pb.ProposalResponse{Response: &pb.Response{}})},
		GetTxSimulatorRv: &mocks.MockTxSim{
			GetTxSimulationResultsRv: &ledger.TxSimulationResults{
				PubSimulationResults: &rwset.TxReadWriteSet{},
			},
		},
	}, &disabled.Provider{})

	fakeMetrics := initFakeMetrics(es)

	signedProp := getSignedPropWithCHIdAndArgs("", "test-chaincode", "test-version", [][]byte{[]byte("test-args")}, t)

	pResp, err := es.ProcessProposal(context.Background(), signedProp)
	assert.NoError(t, err)
	assert.EqualValues(t, 200, pResp.Response.Status)

	
	testEndorsementCompletedMetric(t, fakeMetrics, 1, "", "test-chaincode:test-version", "true")
}

func TestEndorserLSCCInitFails(t *testing.T) {
	es := endorser.NewEndorserServer(pvtEmptyDistributor, &mocks.MockSupport{
		GetApplicationConfigBoolRv: true,
		GetApplicationConfigRv:     &mc.MockApplication{CapabilitiesRv: &mc.MockApplicationCapabilities{}},
		GetTransactionByIDErr:      errors.New(""),
		ChaincodeDefinitionRv:      &ccprovider.ChaincodeData{Escc: "ESCC"},
		ExecuteResp:                &pb.Response{Status: 200, Payload: protoutil.MarshalOrPanic(&pb.ProposalResponse{Response: &pb.Response{}})},
		GetTxSimulatorRv: &mocks.MockTxSim{
			GetTxSimulationResultsRv: &ledger.TxSimulationResults{
				PubSimulationResults: &rwset.TxReadWriteSet{},
			},
		},
		ExecuteCDSError: errors.New(""),
	}, &disabled.Provider{})

	fakeMetrics := initFakeMetrics(es)

	cds := protoutil.MarshalOrPanic(
		&pb.ChaincodeDeploymentSpec{
			ChaincodeSpec: &pb.ChaincodeSpec{
				ChaincodeId: &pb.ChaincodeID{Name: "barf", Version: "0"},
				Type:        pb.ChaincodeSpec_GOLANG,
			},
		},
	)
	signedProp := getSignedPropWithCHIdAndArgs(util.GetTestChainID(), "lscc", "0", [][]byte{[]byte("deploy"), []byte("a"), cds}, t)

	pResp, err := es.ProcessProposal(context.Background(), signedProp)
	assert.NoError(t, err)
	assert.EqualValues(t, 500, pResp.Response.Status)

	
	assert.EqualValues(t, 1, fakeMetrics.initFailed.WithCallCount())
	labelValues := fakeMetrics.initFailed.WithArgsForCall(0)
	assert.EqualValues(t, labelValues, []string{"channel", util.GetTestChainID(), "chaincode", "barf:0"})
	assert.EqualValues(t, 1, fakeMetrics.initFailed.AddCallCount())
	assert.EqualValues(t, 1, fakeMetrics.initFailed.AddArgsForCall(0))

	
	testEndorsementCompletedMetric(t, fakeMetrics, 1, util.GetTestChainID(), "lscc:0", "false")
}

func TestEndorserLSCCDeploySysCC(t *testing.T) {
	SysCCMap := make(map[string]struct{})
	deployedCCName := "barf"
	SysCCMap[deployedCCName] = struct{}{}
	es := endorser.NewEndorserServer(pvtEmptyDistributor, &mocks.MockSupport{
		GetApplicationConfigBoolRv: true,
		GetApplicationConfigRv:     &mc.MockApplication{CapabilitiesRv: &mc.MockApplicationCapabilities{}},
		GetTransactionByIDErr:      errors.New(""),
		ChaincodeDefinitionRv:      &ccprovider.ChaincodeData{Escc: "ESCC"},
		ExecuteResp:                &pb.Response{Status: 200, Payload: protoutil.MarshalOrPanic(&pb.ProposalResponse{Response: &pb.Response{}})},
		GetTxSimulatorRv: &mocks.MockTxSim{
			GetTxSimulationResultsRv: &ledger.TxSimulationResults{
				PubSimulationResults: &rwset.TxReadWriteSet{},
			},
		},
		SysCCMap: SysCCMap,
	}, &disabled.Provider{})

	cds := protoutil.MarshalOrPanic(
		&pb.ChaincodeDeploymentSpec{
			ChaincodeSpec: &pb.ChaincodeSpec{
				ChaincodeId: &pb.ChaincodeID{Name: deployedCCName},
				Type:        pb.ChaincodeSpec_GOLANG,
			},
		},
	)
	signedProp := getSignedPropWithCHIdAndArgs(util.GetTestChainID(), "lscc", "0", [][]byte{[]byte("deploy"), []byte("a"), cds}, t)

	pResp, err := es.ProcessProposal(context.Background(), signedProp)
	assert.NoError(t, err)
	assert.EqualValues(t, 500, pResp.Response.Status)
	assert.Equal(t, "attempting to deploy a system chaincode barf/testchainid", pResp.Response.Message)
}

func TestEndorserGoodPathWEvents(t *testing.T) {
	m := &mock.Mock{}
	m.On("Sign", mock.Anything).Return([]byte{1, 2, 3, 4, 5}, nil)
	m.On("Serialize").Return([]byte{1, 1, 1}, nil)
	m.On("GetTxSimulator", mock.Anything, mock.Anything).Return(newMockTxSim(), nil)
	support := &mocks.MockSupport{
		Mock:                       m,
		GetApplicationConfigBoolRv: true,
		GetApplicationConfigRv:     &mc.MockApplication{CapabilitiesRv: &mc.MockApplicationCapabilities{}},
		GetTransactionByIDErr:      errors.New(""),
		ChaincodeDefinitionRv:      &ccprovider.ChaincodeData{Escc: "ESCC"},
		ExecuteResp:                &pb.Response{Status: 200, Payload: protoutil.MarshalOrPanic(&pb.ProposalResponse{Response: &pb.Response{}})},
		ExecuteEvent:               &pb.ChaincodeEvent{},
	}
	attachPluginEndorser(support, nil)
	es := endorser.NewEndorserServer(pvtEmptyDistributor, support, &disabled.Provider{})

	signedProp := getSignedProp("ccid", "0", t)

	pResp, err := es.ProcessProposal(context.Background(), signedProp)
	assert.NoError(t, err)
	assert.EqualValues(t, 200, pResp.Response.Status)
}

func TestEndorserBadChannel(t *testing.T) {
	es := endorser.NewEndorserServer(pvtEmptyDistributor, &mocks.MockSupport{
		GetApplicationConfigBoolRv: true,
		GetApplicationConfigRv:     &mc.MockApplication{CapabilitiesRv: &mc.MockApplicationCapabilities{}},
		GetTransactionByIDErr:      errors.New(""),
		ChaincodeDefinitionRv:      &ccprovider.ChaincodeData{Escc: "ESCC"},
		ExecuteResp:                &pb.Response{Status: 200, Payload: protoutil.MarshalOrPanic(&pb.ProposalResponse{Response: &pb.Response{}})},
		GetTxSimulatorRv: &mocks.MockTxSim{
			GetTxSimulationResultsRv: &ledger.TxSimulationResults{
				PubSimulationResults: &rwset.TxReadWriteSet{},
			},
		},
	}, &disabled.Provider{})

	signedProp := getSignedPropWithCHID("ccid", "0", "barfchain", t)

	pResp, err := es.ProcessProposal(context.Background(), signedProp)
	assert.Error(t, err)
	assert.EqualValues(t, 500, pResp.Response.Status)
	assert.Equal(t, "access denied: channel [barfchain] creator org [SampleOrg]", pResp.Response.Message)
}

func TestEndorserGoodPath(t *testing.T) {
	m := &mock.Mock{}
	m.On("Sign", mock.Anything).Return([]byte{1, 2, 3, 4, 5}, nil)
	m.On("Serialize").Return([]byte{1, 1, 1}, nil)
	m.On("GetTxSimulator", mock.Anything, mock.Anything).Return(newMockTxSim(), nil)
	support := &mocks.MockSupport{
		Mock:                       m,
		GetApplicationConfigBoolRv: true,
		GetApplicationConfigRv:     &mc.MockApplication{CapabilitiesRv: &mc.MockApplicationCapabilities{}},
		GetTransactionByIDErr:      errors.New(""),
		ChaincodeDefinitionRv:      &ccprovider.ChaincodeData{Name: "ccid", Version: "0", Escc: "ESCC"},
		ExecuteResp:                &pb.Response{Status: 200, Payload: protoutil.MarshalOrPanic(&pb.ProposalResponse{Response: &pb.Response{}})},
	}
	attachPluginEndorser(support, nil)
	es := endorser.NewEndorserServer(pvtEmptyDistributor, support, &disabled.Provider{})

	fakeMetrics := initFakeMetrics(es)

	signedProp := getSignedProp("ccid", "0", t)

	pResp, err := es.ProcessProposal(context.Background(), signedProp)
	assert.NoError(t, err)
	assert.EqualValues(t, 200, pResp.Response.Status)

	
	testEndorsementCompletedMetric(t, fakeMetrics, 1, util.GetTestChainID(), "ccid:0", "true")

	
	assert.EqualValues(t, 1, fakeMetrics.successfulProposals.AddCallCount())
	assert.EqualValues(t, 1, fakeMetrics.successfulProposals.AddArgsForCall(0))
}

func TestEndorserChaincodeCallLogging(t *testing.T) {
	gt := NewGomegaWithT(t)
	m := &mock.Mock{}
	m.On("Sign", mock.Anything).Return([]byte{1, 2, 3, 4, 5}, nil)
	m.On("Serialize").Return([]byte{1, 1, 1}, nil)
	m.On("GetTxSimulator", mock.Anything, mock.Anything).Return(newMockTxSim(), nil)
	support := &mocks.MockSupport{
		Mock:                       m,
		GetApplicationConfigBoolRv: true,
		GetApplicationConfigRv:     &mc.MockApplication{CapabilitiesRv: &mc.MockApplicationCapabilities{}},
		GetTransactionByIDErr:      errors.New(""),
		ChaincodeDefinitionRv:      &ccprovider.ChaincodeData{Escc: "ESCC"},
		ExecuteResp:                &pb.Response{Status: 200, Payload: protoutil.MarshalOrPanic(&pb.ProposalResponse{Response: &pb.Response{}})},
	}
	attachPluginEndorser(support, nil)
	es := endorser.NewEndorserServer(pvtEmptyDistributor, support, &disabled.Provider{})

	buf := gbytes.NewBuffer()
	old := flogging.SetWriter(buf)
	defer flogging.SetWriter(old)

	es.ProcessProposal(context.Background(), getSignedProp("chaincode-name", "chaincode-version", t))

	t.Logf("contents:\n%s", buf.Contents())
	gt.Eventually(buf).Should(gbytes.Say(`INFO.*\[testchainid\]\[[[:xdigit:]]{8}\] Entry chaincode: chaincode-name`))
	gt.Eventually(buf).Should(gbytes.Say(`INFO.*\[testchainid\]\[[[:xdigit:]]{8}\] Exit chaincode: chaincode-name (.*ms)`))
}

func TestEndorserLSCC(t *testing.T) {
	m := &mock.Mock{}
	m.On("Sign", mock.Anything).Return([]byte{1, 2, 3, 4, 5}, nil)
	m.On("Serialize").Return([]byte{1, 1, 1}, nil)
	m.On("GetTxSimulator", mock.Anything, mock.Anything).Return(newMockTxSim(), nil)
	support := &mocks.MockSupport{
		Mock:                       m,
		GetApplicationConfigBoolRv: true,
		GetApplicationConfigRv:     &mc.MockApplication{CapabilitiesRv: &mc.MockApplicationCapabilities{}},
		GetTransactionByIDErr:      errors.New(""),
		ChaincodeDefinitionRv:      &ccprovider.ChaincodeData{Escc: "ESCC"},
		ExecuteResp:                &pb.Response{Status: 200, Payload: protoutil.MarshalOrPanic(&pb.ProposalResponse{Response: &pb.Response{}})},
	}
	attachPluginEndorser(support, nil)
	es := endorser.NewEndorserServer(pvtEmptyDistributor, support, &disabled.Provider{})

	cds := protoutil.MarshalOrPanic(
		&pb.ChaincodeDeploymentSpec{
			ChaincodeSpec: &pb.ChaincodeSpec{
				ChaincodeId: &pb.ChaincodeID{Name: "barf"},
				Type:        pb.ChaincodeSpec_GOLANG,
			},
		},
	)
	signedProp := getSignedPropWithCHIdAndArgs(util.GetTestChainID(), "lscc", "0", [][]byte{[]byte("deploy"), []byte("a"), cds}, t)

	pResp, err := es.ProcessProposal(context.Background(), signedProp)
	assert.NoError(t, err)
	assert.EqualValues(t, 200, pResp.Response.Status)
}

func attachPluginEndorser(support *mocks.MockSupport, signerReturnErr error) {
	csr := &mocks.ChannelStateRetriever{}
	queryCreator := &mocks.QueryCreator{}
	csr.On("NewQueryCreator", mock.Anything).Return(queryCreator, nil)
	sif := &mocks.SigningIdentityFetcher{}
	sif.On("SigningIdentityForRequest", mock.Anything).Return(support, signerReturnErr)
	pm := &mocks.PluginMapper{}
	pm.On("PluginFactoryByName", mock.Anything).Return(&builtin.DefaultEndorsementFactory{})
	support.PluginEndorser = endorser.NewPluginEndorser(&endorser.PluginSupport{
		ChannelStateRetriever:   csr,
		SigningIdentityFetcher:  sif,
		PluginMapper:            pm,
		TransientStoreRetriever: mockTransientStoreRetriever,
	})
}

func TestEndorseWithPlugin(t *testing.T) {
	m := &mock.Mock{}
	m.On("Sign", mock.Anything).Return([]byte{1, 2, 3, 4, 5}, nil)
	m.On("Serialize").Return([]byte{1, 1, 1}, nil)
	m.On("GetTxSimulator", mock.Anything, mock.Anything).Return(newMockTxSim(), nil)
	support := &mocks.MockSupport{
		Mock:                       m,
		GetApplicationConfigBoolRv: true,
		GetApplicationConfigRv:     &mc.MockApplication{CapabilitiesRv: &mc.MockApplicationCapabilities{}},
		GetTransactionByIDErr:      errors.New(""),
		ChaincodeDefinitionRv:      &resourceconfig.MockChaincodeDefinition{EndorsementStr: "ESCC"},
		ExecuteResp:                &pb.Response{Status: 200, Payload: []byte{1}},
	}
	attachPluginEndorser(support, nil)

	es := endorser.NewEndorserServer(pvtEmptyDistributor, support, &disabled.Provider{})

	signedProp := getSignedProp("ccid", "0", t)

	resp, err := es.ProcessProposal(context.Background(), signedProp)
	assert.NoError(t, err)
	assert.Equal(t, []byte{1, 2, 3, 4, 5}, resp.Endorsement.Signature)
	assert.Equal(t, []byte{1, 1, 1}, resp.Endorsement.Endorser)
	assert.Equal(t, 200, int(resp.Response.Status))
}

func TestEndorseEndorsementFailure(t *testing.T) {
	m := &mock.Mock{}
	m.On("Sign", mock.Anything).Return([]byte{1, 2, 3, 4, 5}, nil)
	m.On("Serialize").Return([]byte{1, 1, 1}, nil)
	m.On("GetTxSimulator", mock.Anything, mock.Anything).Return(newMockTxSim(), nil)
	support := &mocks.MockSupport{
		Mock:                       m,
		GetApplicationConfigBoolRv: true,
		GetApplicationConfigRv:     &mc.MockApplication{CapabilitiesRv: &mc.MockApplicationCapabilities{}},
		GetTransactionByIDErr:      errors.New(""),
		ChaincodeDefinitionRv:      &resourceconfig.MockChaincodeDefinition{NameRv: "ccid", VersionRv: "0", EndorsementStr: "ESCC"},
		ExecuteResp:                &pb.Response{Status: 200, Payload: []byte{1}},
	}

	
	attachPluginEndorser(support, fmt.Errorf("sign err"))

	es := endorser.NewEndorserServer(pvtEmptyDistributor, support, &disabled.Provider{})
	fakeMetrics := initFakeMetrics(es)

	signedProp := getSignedProp("ccid", "0", t)

	resp, err := es.ProcessProposal(context.Background(), signedProp)
	assert.NoError(t, err)
	assert.EqualValues(t, 500, resp.Response.Status)

	
	assert.EqualValues(t, 1, fakeMetrics.endorsementsFailed.WithCallCount())
	labelValues := fakeMetrics.endorsementsFailed.WithArgsForCall(0)
	assert.EqualValues(t, labelValues, []string{"channel", util.GetTestChainID(), "chaincode", "ccid:0", "chaincodeerror", "false"})
	assert.EqualValues(t, 1, fakeMetrics.endorsementsFailed.AddCallCount())
	assert.EqualValues(t, 1, fakeMetrics.endorsementsFailed.AddArgsForCall(0))

	
	testEndorsementCompletedMetric(t, fakeMetrics, 1, util.GetTestChainID(), "ccid:0", "false")
}

func TestEndorseEndorsementFailureDueToCCError(t *testing.T) {
	m := &mock.Mock{}
	m.On("Sign", mock.Anything).Return([]byte{1, 2, 3, 4, 5}, nil)
	m.On("Serialize").Return([]byte{1, 1, 1}, nil)
	m.On("GetTxSimulator", mock.Anything, mock.Anything).Return(newMockTxSim(), nil)
	support := &mocks.MockSupport{
		Mock:                       m,
		GetApplicationConfigBoolRv: true,
		GetApplicationConfigRv:     &mc.MockApplication{CapabilitiesRv: &mc.MockApplicationCapabilities{}},
		GetTransactionByIDErr:      errors.New(""),
		ChaincodeDefinitionRv:      &resourceconfig.MockChaincodeDefinition{NameRv: "ccid", VersionRv: "0", EndorsementStr: "ESCC"},
		ExecuteResp:                &pb.Response{Status: 400, Message: "CC error"},
	}

	attachPluginEndorser(support, nil)

	es := endorser.NewEndorserServer(pvtEmptyDistributor, support, &disabled.Provider{})
	fakeMetrics := initFakeMetrics(es)

	signedProp := getSignedProp("ccid", "0", t)

	resp, err := es.ProcessProposal(context.Background(), signedProp)
	assert.NoError(t, err)
	assert.EqualValues(t, 400, int(resp.Response.Status))

	
	assert.EqualValues(t, 1, fakeMetrics.endorsementsFailed.WithCallCount())
	labelValues := fakeMetrics.endorsementsFailed.WithArgsForCall(0)
	assert.EqualValues(t, []string{"channel", util.GetTestChainID(), "chaincode", "ccid:0", "chaincodeerror", "true"}, labelValues)
	assert.EqualValues(t, 1, fakeMetrics.endorsementsFailed.AddCallCount())
	assert.EqualValues(t, 1, fakeMetrics.endorsementsFailed.AddArgsForCall(0))

	
	testEndorsementCompletedMetric(t, fakeMetrics, 1, util.GetTestChainID(), "ccid:0", "false")
}

func TestSimulateProposal(t *testing.T) {
	es := endorser.NewEndorserServer(pvtEmptyDistributor, &mocks.MockSupport{
		GetApplicationConfigBoolRv: true,
		GetApplicationConfigRv:     &mc.MockApplication{CapabilitiesRv: &mc.MockApplicationCapabilities{}},
		GetTransactionByIDErr:      errors.New(""),
		ChaincodeDefinitionRv:      &ccprovider.ChaincodeData{Escc: "ESCC"},
		ExecuteResp:                &pb.Response{Status: 200, Payload: protoutil.MarshalOrPanic(&pb.ProposalResponse{Response: &pb.Response{}})},
		GetTxSimulatorRv: &mocks.MockTxSim{
			GetTxSimulationResultsRv: &ledger.TxSimulationResults{
				PubSimulationResults: &rwset.TxReadWriteSet{},
			},
		},
	}, &disabled.Provider{})

	_, _, _, err := es.SimulateProposal(&ccprovider.TransactionParams{}, "")
	assert.Error(t, err)
}

func TestEndorserAcquireTxSimulator(t *testing.T) {
	tc := []struct {
		name          string
		chainID       string
		chaincodeName string
		simAcquired   bool
	}{
		{"empty channel", "", "ignored", false},
		{"query scc", util.GetTestChainID(), "qscc", false},
		{"config scc", util.GetTestChainID(), "cscc", false},
		{"mainline", util.GetTestChainID(), "chaincode", true},
	}

	expectedResponse := &pb.Response{Status: 200, Payload: protoutil.MarshalOrPanic(&pb.ProposalResponse{Response: &pb.Response{}})}
	for _, tt := range tc {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			m := &mock.Mock{}
			m.On("Sign", mock.Anything).Return([]byte{1, 2, 3, 4, 5}, nil)
			m.On("Serialize").Return([]byte{1, 1, 1}, nil)
			m.On("GetTxSimulator", mock.Anything, mock.Anything).Return(newMockTxSim(), nil)
			support := &mocks.MockSupport{
				Mock:                       m,
				GetApplicationConfigBoolRv: true,
				GetApplicationConfigRv:     &mc.MockApplication{CapabilitiesRv: &mc.MockApplicationCapabilities{}},
				GetTransactionByIDErr:      errors.New(""),
				ChaincodeDefinitionRv:      &ccprovider.ChaincodeData{Escc: "ESCC"},
				ExecuteResp:                expectedResponse,
			}
			attachPluginEndorser(support, nil)
			es := endorser.NewEndorserServer(
				pvtEmptyDistributor,
				support,
				&disabled.Provider{},
			)

			t.Parallel()
			args := [][]byte{[]byte("args")}
			signedProp := getSignedPropWithCHIdAndArgs(tt.chainID, tt.chaincodeName, "version", args, t)

			resp, err := es.ProcessProposal(context.Background(), signedProp)
			assert.NoError(t, err)
			assert.Equal(t, expectedResponse, resp.Response)

			if tt.simAcquired {
				m.AssertCalled(t, "GetTxSimulator", mock.Anything, mock.Anything)
			} else {
				m.AssertNotCalled(t, "GetTxSimulator", mock.Anything, mock.Anything)
			}
		})
	}
}

var signer msp.SigningIdentity

func TestMain(m *testing.M) {
	
	err := msptesttools.LoadMSPSetupForTesting()
	if err != nil {
		fmt.Printf("Could not initialize msp/signer, err %s", err)
		os.Exit(-1)
		return
	}
	signer, err = mgmt.GetLocalMSP().GetDefaultSigningIdentity()
	if err != nil {
		fmt.Printf("Could not initialize msp/signer")
		os.Exit(-1)
		return
	}

	retVal := m.Run()
	os.Exit(retVal)
}



type support interface {
	endorser.Support
}
