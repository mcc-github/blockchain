/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package v12

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain-chaincode-go/shim"
	"github.com/mcc-github/blockchain-chaincode-go/shimtest"
	"github.com/mcc-github/blockchain-protos-go/common"
	"github.com/mcc-github/blockchain-protos-go/ledger/rwset/kvrwset"
	mspproto "github.com/mcc-github/blockchain-protos-go/msp"
	"github.com/mcc-github/blockchain-protos-go/peer"
	"github.com/mcc-github/blockchain/bccsp/sw"
	"github.com/mcc-github/blockchain/common/capabilities"
	"github.com/mcc-github/blockchain/common/cauthdsl"
	"github.com/mcc-github/blockchain/common/channelconfig"
	mc "github.com/mcc-github/blockchain/common/mocks/config"
	lm "github.com/mcc-github/blockchain/common/mocks/ledger"
	"github.com/mcc-github/blockchain/common/mocks/scc"
	"github.com/mcc-github/blockchain/common/util"
	aclmocks "github.com/mcc-github/blockchain/core/aclmgmt/mocks"
	"github.com/mcc-github/blockchain/core/committer/txvalidator/v14"
	mocks2 "github.com/mcc-github/blockchain/core/committer/txvalidator/v14/mocks"
	"github.com/mcc-github/blockchain/core/common/ccpackage"
	"github.com/mcc-github/blockchain/core/common/ccprovider"
	"github.com/mcc-github/blockchain/core/common/privdata"
	validation "github.com/mcc-github/blockchain/core/handlers/validation/api/capabilities"
	"github.com/mcc-github/blockchain/core/handlers/validation/builtin/v12/mocks"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/rwsetutil"
	"github.com/mcc-github/blockchain/core/scc/lscc"
	"github.com/mcc-github/blockchain/msp"
	mspmgmt "github.com/mcc-github/blockchain/msp/mgmt"
	msptesttools "github.com/mcc-github/blockchain/msp/mgmt/testtools"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func createTx(endorsedByDuplicatedIdentity bool) (*common.Envelope, error) {
	ccid := &peer.ChaincodeID{Name: "foo", Version: "v1"}
	cis := &peer.ChaincodeInvocationSpec{ChaincodeSpec: &peer.ChaincodeSpec{ChaincodeId: ccid}}

	prop, _, err := protoutil.CreateProposalFromCIS(common.HeaderType_ENDORSER_TRANSACTION, util.GetTestChannelID(), cis, sid)
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

func processSignedCDS(cds *peer.ChaincodeDeploymentSpec, policy *common.SignaturePolicyEnvelope) ([]byte, error) {
	env, err := ccpackage.OwnerCreateSignedCCDepSpec(cds, policy, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create package %s", err)
	}

	b := protoutil.MarshalOrPanic(env)

	cryptoProvider, err := sw.NewDefaultSecurityLevelWithKeystore(sw.NewDummyKeyStore())
	if err != nil {
		return nil, fmt.Errorf("could not create bootBCCSP %s", err)
	}
	ccpack := &ccprovider.SignedCDSPackage{GetHasher: cryptoProvider}
	cd, err := ccpack.InitFromBuffer(b)
	if err != nil {
		return nil, fmt.Errorf("error owner creating package %s", err)
	}

	if err = ccpack.PutChaincodeToFS(); err != nil {
		return nil, fmt.Errorf("error putting package on the FS %s", err)
	}

	cd.InstantiationPolicy = protoutil.MarshalOrPanic(policy)

	return protoutil.MarshalOrPanic(cd), nil
}

func constructDeploymentSpec(name, path, version string, initArgs [][]byte, createFS bool) (*peer.ChaincodeDeploymentSpec, error) {
	spec := &peer.ChaincodeSpec{Type: 1, ChaincodeId: &peer.ChaincodeID{Name: name, Path: path, Version: version}, Input: &peer.ChaincodeInput{Args: initArgs}}

	codePackageBytes := bytes.NewBuffer(nil)
	gz := gzip.NewWriter(codePackageBytes)
	tw := tar.NewWriter(gz)

	payload := []byte(name + path + version)
	err := tw.WriteHeader(&tar.Header{
		Name: "src/garbage.go",
		Size: int64(len(payload)),
		Mode: 0100644,
	})
	if err != nil {
		return nil, err
	}

	_, err = tw.Write(payload)
	if err != nil {
		return nil, err
	}

	tw.Close()
	gz.Close()

	chaincodeDeploymentSpec := &peer.ChaincodeDeploymentSpec{ChaincodeSpec: spec, CodePackage: codePackageBytes.Bytes()}

	if createFS {
		cryptoProvider, err := sw.NewDefaultSecurityLevelWithKeystore(sw.NewDummyKeyStore())
		if err != nil {
			return nil, err
		}
		ccinfoFSImpl := &ccprovider.CCInfoFSImpl{GetHasher: cryptoProvider}
		_, err = ccinfoFSImpl.PutChaincode(chaincodeDeploymentSpec)
		if err != nil {
			return nil, err
		}
	}

	return chaincodeDeploymentSpec, nil
}

func createCCDataRWsetWithCollection(nameK, nameV, version string, policy []byte, collectionConfigPackage []byte) ([]byte, error) {
	cd := &ccprovider.ChaincodeData{
		Name:                nameV,
		Version:             version,
		InstantiationPolicy: policy,
	}

	cdbytes := protoutil.MarshalOrPanic(cd)

	rwsetBuilder := rwsetutil.NewRWSetBuilder()
	rwsetBuilder.AddToWriteSet("lscc", nameK, cdbytes)
	rwsetBuilder.AddToWriteSet("lscc", privdata.BuildCollectionKVSKey(nameK), collectionConfigPackage)
	sr, err := rwsetBuilder.GetTxSimulationResults()
	if err != nil {
		return nil, err
	}
	return sr.GetPubSimulationBytes()
}

func createCCDataRWset(nameK, nameV, version string, policy []byte) ([]byte, error) {
	cd := &ccprovider.ChaincodeData{
		Name:                nameV,
		Version:             version,
		InstantiationPolicy: policy,
	}

	cdbytes := protoutil.MarshalOrPanic(cd)

	rwsetBuilder := rwsetutil.NewRWSetBuilder()
	rwsetBuilder.AddToWriteSet("lscc", nameK, cdbytes)
	sr, err := rwsetBuilder.GetTxSimulationResults()
	if err != nil {
		return nil, err
	}
	return sr.GetPubSimulationBytes()
}

func createLSCCTxWithCollection(ccname, ccver, f string, res []byte, policy []byte, ccpBytes []byte) (*common.Envelope, error) {
	return createLSCCTxPutCdsWithCollection(ccname, ccver, f, res, nil, true, policy, ccpBytes)
}

func createLSCCTx(ccname, ccver, f string, res []byte) (*common.Envelope, error) {
	return createLSCCTxPutCds(ccname, ccver, f, res, nil, true)
}

func createLSCCTxPutCdsWithCollection(ccname, ccver, f string, res, cdsbytes []byte, putcds bool, policy []byte, ccpBytes []byte) (*common.Envelope, error) {
	cds := &peer.ChaincodeDeploymentSpec{
		ChaincodeSpec: &peer.ChaincodeSpec{
			ChaincodeId: &peer.ChaincodeID{
				Name:    ccname,
				Version: ccver,
			},
			Type: peer.ChaincodeSpec_GOLANG,
		},
	}

	cdsBytes, err := proto.Marshal(cds)
	if err != nil {
		return nil, err
	}

	var cis *peer.ChaincodeInvocationSpec
	if putcds {
		if cdsbytes != nil {
			cdsBytes = cdsbytes
		}
		cis = &peer.ChaincodeInvocationSpec{
			ChaincodeSpec: &peer.ChaincodeSpec{
				ChaincodeId: &peer.ChaincodeID{Name: "lscc"},
				Input: &peer.ChaincodeInput{
					Args: [][]byte{[]byte(f), []byte("barf"), cdsBytes, []byte("escc"), []byte("vscc"), policy, ccpBytes},
				},
				Type: peer.ChaincodeSpec_GOLANG,
			},
		}
	} else {
		cis = &peer.ChaincodeInvocationSpec{
			ChaincodeSpec: &peer.ChaincodeSpec{
				ChaincodeId: &peer.ChaincodeID{Name: "lscc"},
				Input: &peer.ChaincodeInput{
					Args: [][]byte{[]byte(f), []byte("barf")},
				},
				Type: peer.ChaincodeSpec_GOLANG,
			},
		}
	}

	prop, _, err := protoutil.CreateProposalFromCIS(common.HeaderType_ENDORSER_TRANSACTION, util.GetTestChannelID(), cis, sid)
	if err != nil {
		return nil, err
	}

	ccid := &peer.ChaincodeID{Name: ccname, Version: ccver}

	presp, err := protoutil.CreateProposalResponse(prop.Header, prop.Payload, &peer.Response{Status: 200}, res, nil, ccid, id)
	if err != nil {
		return nil, err
	}

	return protoutil.CreateSignedTx(prop, id, presp)
}

func createLSCCTxPutCds(ccname, ccver, f string, res, cdsbytes []byte, putcds bool) (*common.Envelope, error) {
	cds := &peer.ChaincodeDeploymentSpec{
		ChaincodeSpec: &peer.ChaincodeSpec{
			ChaincodeId: &peer.ChaincodeID{
				Name:    ccname,
				Version: ccver,
			},
			Type: peer.ChaincodeSpec_GOLANG,
		},
	}

	cdsBytes, err := proto.Marshal(cds)
	if err != nil {
		return nil, err
	}

	var cis *peer.ChaincodeInvocationSpec
	if putcds {
		if cdsbytes != nil {
			cdsBytes = cdsbytes
		}
		cis = &peer.ChaincodeInvocationSpec{
			ChaincodeSpec: &peer.ChaincodeSpec{
				ChaincodeId: &peer.ChaincodeID{Name: "lscc"},
				Input: &peer.ChaincodeInput{
					Args: [][]byte{[]byte(f), []byte("barf"), cdsBytes},
				},
				Type: peer.ChaincodeSpec_GOLANG,
			},
		}
	} else {
		cis = &peer.ChaincodeInvocationSpec{
			ChaincodeSpec: &peer.ChaincodeSpec{
				ChaincodeId: &peer.ChaincodeID{Name: "lscc"},
				Input: &peer.ChaincodeInput{
					Args: [][]byte{[]byte(f), []byte(ccname)},
				},
				Type: peer.ChaincodeSpec_GOLANG,
			},
		}
	}

	prop, _, err := protoutil.CreateProposalFromCIS(common.HeaderType_ENDORSER_TRANSACTION, util.GetTestChannelID(), cis, sid)
	if err != nil {
		return nil, err
	}

	ccid := &peer.ChaincodeID{Name: ccname, Version: ccver}

	presp, err := protoutil.CreateProposalResponse(prop.Header, prop.Payload, &peer.Response{Status: 200}, res, nil, ccid, id)
	if err != nil {
		return nil, err
	}

	return protoutil.CreateSignedTx(prop, id, presp)
}

func getSignedByMSPMemberPolicy(mspID string) ([]byte, error) {
	p := cauthdsl.SignedByMspMember(mspID)

	b, err := protoutil.Marshal(p)
	if err != nil {
		return nil, fmt.Errorf("Could not marshal policy, err %s", err)
	}

	return b, err
}

func getSignedByOneMemberTwicePolicy(mspID string) ([]byte, error) {
	principal := &mspproto.MSPPrincipal{
		PrincipalClassification: mspproto.MSPPrincipal_ROLE,
		Principal:               protoutil.MarshalOrPanic(&mspproto.MSPRole{Role: mspproto.MSPRole_MEMBER, MspIdentifier: mspID})}

	p := &common.SignaturePolicyEnvelope{
		Version:    0,
		Rule:       cauthdsl.NOutOf(2, []*common.SignaturePolicy{cauthdsl.SignedBy(0), cauthdsl.SignedBy(0)}),
		Identities: []*mspproto.MSPPrincipal{principal},
	}
	b, err := protoutil.Marshal(p)
	if err != nil {
		return nil, fmt.Errorf("Could not marshal policy, err %s", err)
	}

	return b, err
}

func getSignedByMSPAdminPolicy(mspID string) ([]byte, error) {
	p := cauthdsl.SignedByMspAdmin(mspID)

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
	sf := &txvalidator.StateFetcherImpl{QueryExecutorCreator: qec}
	is := &mocks.IdentityDeserializer{}
	pe := &txvalidator.PolicyEvaluator{
		IdentityDeserializer: mspmgmt.GetManagerForChain(util.GetTestChannelID()),
	}
	return New(c, sf, is, pe)
}

func TestDeduplicateIdentity(t *testing.T) {
	
	proposalResponsePayload := []byte{3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3}
	prpBuff := make([]byte, len(proposalResponsePayload), len(proposalResponsePayload)*2)
	copy(prpBuff, proposalResponsePayload)

	identity1 := protoutil.MarshalOrPanic(&mspproto.SerializedIdentity{
		IdBytes: []byte{1, 1, 1},
	})
	identity2 := protoutil.MarshalOrPanic(&mspproto.SerializedIdentity{
		IdBytes: []byte{2, 2, 2},
	})

	chaincodeActionPayload := &peer.ChaincodeActionPayload{
		Action: &peer.ChaincodeEndorsedAction{
			Endorsements: []*peer.Endorsement{
				{
					Endorser: identity1,
				},
				{
					Endorser: identity2,
				},
			},
			ProposalResponsePayload: prpBuff,
		},
	}

	signedData, err := (&Validator{}).deduplicateIdentity(chaincodeActionPayload)
	assert.NoError(t, err)
	
	assert.Equal(t, proposalResponsePayload, signedData[0].Data[:len(proposalResponsePayload)])
	assert.Equal(t, proposalResponsePayload, signedData[1].Data[:len(proposalResponsePayload)])
	
	assert.Equal(t, identity1, signedData[0].Data[len(proposalResponsePayload):])
	assert.Equal(t, identity2, signedData[1].Data[len(proposalResponsePayload):])
}

func TestInvoke(t *testing.T) {
	v := newValidationInstance(make(map[string]map[string][]byte))

	
	var err error
	b := &common.Block{Data: &common.BlockData{Data: [][]byte{[]byte("a")}}}
	err = v.Validate(b, "foo", 0, 0, []byte("a"))
	assert.Error(t, err)

	
	b = &common.Block{Data: &common.BlockData{Data: [][]byte{protoutil.MarshalOrPanic(&common.Envelope{Payload: []byte("barf")})}}}
	err = v.Validate(b, "foo", 0, 0, []byte("a"))
	assert.Error(t, err)

	
	e := protoutil.MarshalOrPanic(&common.Envelope{Payload: protoutil.MarshalOrPanic(&common.Payload{Header: &common.Header{ChannelHeader: []byte("barf")}})})
	b = &common.Block{Data: &common.BlockData{Data: [][]byte{e}}}
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

	
	b = &common.Block{Data: &common.BlockData{Data: [][]byte{envBytes}}}
	err = v.Validate(b, "foo", 0, 0, []byte("barf"))
	assert.Error(t, err)

	policy, err := getSignedByMSPMemberPolicy(mspid)
	if err != nil {
		t.Fatalf("failed getting policy, err %s", err)
	}

	
	e = protoutil.MarshalOrPanic(&common.Envelope{Payload: protoutil.MarshalOrPanic(&common.Payload{Header: &common.Header{ChannelHeader: protoutil.MarshalOrPanic(&common.ChannelHeader{Type: int32(common.HeaderType_ORDERER_TRANSACTION)})}})})
	b = &common.Block{Data: &common.BlockData{Data: [][]byte{e}}}
	err = v.Validate(b, "foo", 0, 0, policy)
	assert.Error(t, err)

	
	e = protoutil.MarshalOrPanic(&common.Envelope{Payload: protoutil.MarshalOrPanic(&common.Payload{Header: &common.Header{ChannelHeader: protoutil.MarshalOrPanic(&common.ChannelHeader{Type: int32(common.HeaderType_ORDERER_TRANSACTION)})}})})
	b = &common.Block{Data: &common.BlockData{Data: [][]byte{e}}}
	err = v.Validate(b, "foo", 0, 0, policy)
	assert.Error(t, err)

	
	b = &common.Block{Data: &common.BlockData{Data: [][]byte{envBytes}}}
	err = v.Validate(b, "foo", 0, 0, policy)
	assert.NoError(t, err)

	
	policy, err = getSignedByMSPMemberPolicy("barf")
	if err != nil {
		t.Fatalf("failed getting policy, err %s", err)
	}

	err = v.Validate(b, "foo", 0, 0, policy)
	assert.Error(t, err)

	
	policy, err = getSignedByOneMemberTwicePolicy(mspid)
	if err != nil {
		t.Fatalf("failed getting policy, err %s", err)
	}
	tx, err = createTx(true)
	if err != nil {
		t.Fatalf("createTx returned err %s", err)
	}
	envBytes, err = protoutil.GetBytesEnvelope(tx)
	if err != nil {
		t.Fatalf("GetBytesEnvelope returned err %s", err)
	}
	b = &common.Block{Data: &common.BlockData{Data: [][]byte{envBytes}}}
	err = v.Validate(b, "foo", 0, 0, policy)
	assert.Error(t, err)
}

func TestRWSetTooBig(t *testing.T) {
	state := make(map[string]map[string][]byte)
	mp := (&scc.MocksccProviderFactory{
		Qe:                    lm.NewMockQueryExecutor(state),
		ApplicationConfigBool: true,
		ApplicationConfigRv:   &mc.MockApplication{CapabilitiesRv: &mc.MockApplicationCapabilities{}},
	}).NewSystemChaincodeProvider().(*scc.MocksccProviderImpl)

	v := newValidationInstance(state)

	mockAclProvider := &aclmocks.MockACLProvider{}
	lccc := lscc.New(nil, mp, mockAclProvider, mockMSPIDGetter, &mockPolicyChecker{})
	stublccc := shimtest.NewMockStub("lscc", lccc)
	state["lscc"] = stublccc.State

	r := stublccc.MockInit("1", [][]byte{})
	if r.Status != shim.OK {
		fmt.Println("Init failed", string(r.Message))
		t.FailNow()
	}

	ccname := "mycc"
	ccver := "1"

	cd := &ccprovider.ChaincodeData{
		Name:                ccname,
		Version:             ccver,
		InstantiationPolicy: nil,
	}

	cdbytes := protoutil.MarshalOrPanic(cd)

	rwsetBuilder := rwsetutil.NewRWSetBuilder()
	rwsetBuilder.AddToWriteSet("lscc", ccname, cdbytes)
	rwsetBuilder.AddToWriteSet("lscc", "spurious", []byte("spurious"))

	sr, err := rwsetBuilder.GetTxSimulationResults()
	assert.NoError(t, err)
	srBytes, err := sr.GetPubSimulationBytes()
	assert.NoError(t, err)
	tx, err := createLSCCTx(ccname, ccver, lscc.DEPLOY, srBytes)
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

	b := &common.Block{Data: &common.BlockData{Data: [][]byte{envBytes}}}
	err = v.Validate(b, "lscc", 0, 0, policy)
	assert.Error(t, err)
}

func TestValidateDeployFail(t *testing.T) {
	state := make(map[string]map[string][]byte)
	mp := (&scc.MocksccProviderFactory{
		Qe:                    lm.NewMockQueryExecutor(state),
		ApplicationConfigBool: true,
		ApplicationConfigRv:   &mc.MockApplication{CapabilitiesRv: &mc.MockApplicationCapabilities{}},
	}).NewSystemChaincodeProvider().(*scc.MocksccProviderImpl)

	v := newValidationInstance(state)
	mockAclProvider := &aclmocks.MockACLProvider{}
	lccc := lscc.New(nil, mp, mockAclProvider, mockMSPIDGetter, &mockPolicyChecker{})
	stublccc := shimtest.NewMockStub("lscc", lccc)
	state["lscc"] = stublccc.State

	ccname := "mycc"
	ccver := "1"

	
	
	

	tx, err := createLSCCTx(ccname, ccver, lscc.DEPLOY, nil)
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

	b := &common.Block{Data: &common.BlockData{Data: [][]byte{envBytes}}}
	err = v.Validate(b, "lscc", 0, 0, policy)
	assert.Error(t, err)

	
	
	

	rwsetBuilder := rwsetutil.NewRWSetBuilder()
	rwsetBuilder.AddToWriteSet("lscc", ccname, []byte("barf"))
	sr, err := rwsetBuilder.GetTxSimulationResults()
	assert.NoError(t, err)
	resBogusBytes, err := sr.GetPubSimulationBytes()
	assert.NoError(t, err)
	tx, err = createLSCCTx(ccname, ccver, lscc.DEPLOY, resBogusBytes)
	if err != nil {
		t.Fatalf("createTx returned err %s", err)
	}

	envBytes, err = protoutil.GetBytesEnvelope(tx)
	if err != nil {
		t.Fatalf("GetBytesEnvelope returned err %s", err)
	}

	
	policy, err = getSignedByMSPMemberPolicy(mspid)
	if err != nil {
		t.Fatalf("failed getting policy, err %s", err)
	}

	b = &common.Block{Data: &common.BlockData{Data: [][]byte{envBytes}}}
	err = v.Validate(b, "lscc", 0, 0, policy)
	assert.Error(t, err)

	
	
	

	res, err := createCCDataRWset(ccname, ccname, ccver, nil)
	assert.NoError(t, err)

	tx, err = createLSCCTxPutCds(ccname, ccver, lscc.DEPLOY, res, nil, false)
	if err != nil {
		t.Fatalf("createTx returned err %s", err)
	}

	envBytes, err = protoutil.GetBytesEnvelope(tx)
	if err != nil {
		t.Fatalf("GetBytesEnvelope returned err %s", err)
	}

	
	policy, err = getSignedByMSPMemberPolicy(mspid)
	if err != nil {
		t.Fatalf("failed getting policy, err %s", err)
	}

	b = &common.Block{Data: &common.BlockData{Data: [][]byte{envBytes}}}
	err = v.Validate(b, "lscc", 0, 0, policy)
	assert.Error(t, err)

	
	
	

	res, err = createCCDataRWset(ccname, ccname, ccver, nil)
	assert.NoError(t, err)

	tx, err = createLSCCTxPutCds(ccname, ccver, lscc.DEPLOY, res, []byte("barf"), true)
	if err != nil {
		t.Fatalf("createTx returned err %s", err)
	}

	envBytes, err = protoutil.GetBytesEnvelope(tx)
	if err != nil {
		t.Fatalf("GetBytesEnvelope returned err %s", err)
	}

	
	policy, err = getSignedByMSPMemberPolicy(mspid)
	if err != nil {
		t.Fatalf("failed getting policy, err %s", err)
	}

	b = &common.Block{Data: &common.BlockData{Data: [][]byte{envBytes}}}
	err = v.Validate(b, "lscc", 0, 0, policy)
	assert.Error(t, err)

	
	
	

	res, err = createCCDataRWset(ccname, ccname, ccver+".1", nil)
	assert.NoError(t, err)

	tx, err = createLSCCTx(ccname, ccver, lscc.DEPLOY, res)
	if err != nil {
		t.Fatalf("createTx returned err %s", err)
	}

	envBytes, err = protoutil.GetBytesEnvelope(tx)
	if err != nil {
		t.Fatalf("GetBytesEnvelope returned err %s", err)
	}

	
	policy, err = getSignedByMSPMemberPolicy(mspid)
	if err != nil {
		t.Fatalf("failed getting policy, err %s", err)
	}

	b = &common.Block{Data: &common.BlockData{Data: [][]byte{envBytes}}}
	err = v.Validate(b, "lscc", 0, 0, policy)
	assert.Error(t, err)

	
	
	

	tx, err = createLSCCTx(ccname, ccver, lscc.DEPLOY, []byte("barf"))
	if err != nil {
		t.Fatalf("createTx returned err %s", err)
	}

	envBytes, err = protoutil.GetBytesEnvelope(tx)
	if err != nil {
		t.Fatalf("GetBytesEnvelope returned err %s", err)
	}

	
	policy, err = getSignedByMSPMemberPolicy(mspid)
	if err != nil {
		t.Fatalf("failed getting policy, err %s", err)
	}

	b = &common.Block{Data: &common.BlockData{Data: [][]byte{envBytes}}}
	err = v.Validate(b, "lscc", 0, 0, policy)
	assert.Error(t, err)

	
	
	

	res, err = createCCDataRWset(ccname+".badbad", ccname, ccver, nil)
	assert.NoError(t, err)

	tx, err = createLSCCTx(ccname, ccver, lscc.DEPLOY, res)
	if err != nil {
		t.Fatalf("createTx returned err %s", err)
	}

	envBytes, err = protoutil.GetBytesEnvelope(tx)
	if err != nil {
		t.Fatalf("GetBytesEnvelope returned err %s", err)
	}

	policy, err = getSignedByMSPMemberPolicy(mspid)
	if err != nil {
		t.Fatalf("failed getting policy, err %s", err)
	}

	b = &common.Block{Data: &common.BlockData{Data: [][]byte{envBytes}}}
	err = v.Validate(b, "lscc", 0, 0, policy)
	assert.Error(t, err)

	
	
	

	res, err = createCCDataRWset(ccname, ccname+".badbad", ccver, nil)
	assert.NoError(t, err)

	tx, err = createLSCCTx(ccname, ccver, lscc.DEPLOY, res)
	if err != nil {
		t.Fatalf("createTx returned err %s", err)
	}

	envBytes, err = protoutil.GetBytesEnvelope(tx)
	if err != nil {
		t.Fatalf("GetBytesEnvelope returned err %s", err)
	}

	policy, err = getSignedByMSPMemberPolicy(mspid)
	if err != nil {
		t.Fatalf("failed getting policy, err %s", err)
	}

	b = &common.Block{Data: &common.BlockData{Data: [][]byte{envBytes}}}
	err = v.Validate(b, "lscc", 0, 0, policy)
	assert.Error(t, err)

	
	
	

	cd := &ccprovider.ChaincodeData{
		Name:                ccname,
		Version:             ccver,
		InstantiationPolicy: nil,
	}

	cdbytes := protoutil.MarshalOrPanic(cd)
	rwsetBuilder = rwsetutil.NewRWSetBuilder()
	rwsetBuilder.AddToWriteSet("lscc", ccname, cdbytes)
	rwsetBuilder.AddToWriteSet("bogusbogus", "key", []byte("val"))
	sr, err = rwsetBuilder.GetTxSimulationResults()
	assert.NoError(t, err)
	srBytes, err := sr.GetPubSimulationBytes()
	assert.NoError(t, err)
	tx, err = createLSCCTx(ccname, ccver, lscc.DEPLOY, srBytes)
	if err != nil {
		t.Fatalf("createTx returned err %s", err)
	}

	envBytes, err = protoutil.GetBytesEnvelope(tx)
	if err != nil {
		t.Fatalf("GetBytesEnvelope returned err %s", err)
	}

	policy, err = getSignedByMSPMemberPolicy(mspid)
	if err != nil {
		t.Fatalf("failed getting policy, err %s", err)
	}

	b = &common.Block{Data: &common.BlockData{Data: [][]byte{envBytes}}}
	err = v.Validate(b, "lscc", 0, 0, policy)
	assert.Error(t, err)

}

func TestAlreadyDeployed(t *testing.T) {
	state := make(map[string]map[string][]byte)
	mp := (&scc.MocksccProviderFactory{
		Qe:                    lm.NewMockQueryExecutor(state),
		ApplicationConfigBool: true,
		ApplicationConfigRv:   &mc.MockApplication{CapabilitiesRv: &mc.MockApplicationCapabilities{}},
	}).NewSystemChaincodeProvider().(*scc.MocksccProviderImpl)

	v := newValidationInstance(state)
	mockAclProvider := &aclmocks.MockACLProvider{}
	lccc := lscc.New(nil, mp, mockAclProvider, mockMSPIDGetter, &mockPolicyChecker{})
	stublccc := shimtest.NewMockStub("lscc", lccc)
	state["lscc"] = stublccc.State

	ccname := "mycc"
	ccver := "alreadydeployed"
	path := "mychaincode"

	cds, err := constructDeploymentSpec(ccname, path, ccver, [][]byte{[]byte("init"), []byte("a"), []byte("100"), []byte("b"), []byte("200")}, true)
	if err != nil {
		fmt.Printf("%s\n", err)
		t.FailNow()
	}

	var b []byte
	if b, err = proto.Marshal(cds); err != nil || b == nil {
		t.FailNow()
	}

	sProp2, _ := protoutil.MockSignedEndorserProposal2OrPanic(chainId, &peer.ChaincodeSpec{}, id)
	args := [][]byte{[]byte("deploy"), []byte(ccname), b}
	if res := stublccc.MockInvokeWithSignedProposal("1", args, sProp2); res.Status != shim.OK {
		fmt.Printf("%#v\n", res)
		t.FailNow()
	}

	simresres, err := createCCDataRWset(ccname, ccname, ccver, nil)
	assert.NoError(t, err)

	tx, err := createLSCCTx(ccname, ccver, lscc.DEPLOY, simresres)
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

	bl := &common.Block{Data: &common.BlockData{Data: [][]byte{envBytes}}}
	err = v.Validate(bl, "lscc", 0, 0, policy)
	assert.Error(t, err)
}

func TestValidateDeployNoLedger(t *testing.T) {
	qec := &mocks2.QueryExecutorCreator{}
	qec.On("NewQueryExecutor").Return(nil, errors.New("failed obtaining query executor"))
	v := newCustomValidationInstance(qec, &mc.MockApplicationCapabilities{})

	ccname := "mycc"
	ccver := "1"

	defaultPolicy, err := getSignedByMSPAdminPolicy(mspid)
	assert.NoError(t, err)
	res, err := createCCDataRWset(ccname, ccname, ccver, defaultPolicy)
	assert.NoError(t, err)

	tx, err := createLSCCTx(ccname, ccver, lscc.DEPLOY, res)
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

	b := &common.Block{Data: &common.BlockData{Data: [][]byte{envBytes}}}
	err = v.Validate(b, "lscc", 0, 0, policy)
	assert.Error(t, err)
}

func TestValidateDeployNOKNilChaincodeSpec(t *testing.T) {
	state := make(map[string]map[string][]byte)
	mp := (&scc.MocksccProviderFactory{
		Qe:                    lm.NewMockQueryExecutor(state),
		ApplicationConfigBool: true,
		ApplicationConfigRv:   &mc.MockApplication{CapabilitiesRv: &mc.MockApplicationCapabilities{}},
	}).NewSystemChaincodeProvider().(*scc.MocksccProviderImpl)

	v := newValidationInstance(state)

	mockAclProvider := &aclmocks.MockACLProvider{}
	lccc := lscc.New(nil, mp, mockAclProvider, mockMSPIDGetter, &mockPolicyChecker{})
	stublccc := shimtest.NewMockStub("lscc", lccc)
	state["lscc"] = stublccc.State

	ccname := "mycc"
	ccver := "1"

	defaultPolicy, err := getSignedByMSPAdminPolicy(mspid)
	assert.NoError(t, err)
	res, err := createCCDataRWset(ccname, ccname, ccver, defaultPolicy)
	assert.NoError(t, err)

	
	cdsBytes, err := proto.Marshal(&peer.ChaincodeDeploymentSpec{})
	assert.NoError(t, err)

	
	cis := &peer.ChaincodeInvocationSpec{
		ChaincodeSpec: &peer.ChaincodeSpec{
			ChaincodeId: &peer.ChaincodeID{Name: "lscc"},
			Input: &peer.ChaincodeInput{
				Args: [][]byte{[]byte(lscc.DEPLOY), []byte("barf"), cdsBytes},
			},
			Type: peer.ChaincodeSpec_GOLANG,
		},
	}

	prop, _, err := protoutil.CreateProposalFromCIS(common.HeaderType_ENDORSER_TRANSACTION, util.GetTestChannelID(), cis, sid)
	assert.NoError(t, err)

	ccid := &peer.ChaincodeID{Name: ccname, Version: ccver}

	presp, err := protoutil.CreateProposalResponse(prop.Header, prop.Payload, &peer.Response{Status: 200}, res, nil, ccid, id)
	assert.NoError(t, err)

	env, err := protoutil.CreateSignedTx(prop, id, presp)
	assert.NoError(t, err)

	
	policy, err := getSignedByMSPMemberPolicy(mspid)
	if err != nil {
		t.Fatalf("failed getting policy, err %s", err)
	}

	b := &common.Block{Data: &common.BlockData{Data: [][]byte{protoutil.MarshalOrPanic(env)}}, Header: &common.BlockHeader{}}
	err = v.Validate(b, "lscc", 0, 0, policy)
	assert.EqualError(t, err, "VSCC error: invocation of lscc(deploy) does not have appropriate arguments")
}

func TestValidateDeployOK(t *testing.T) {
	state := make(map[string]map[string][]byte)
	mp := (&scc.MocksccProviderFactory{
		Qe:                    lm.NewMockQueryExecutor(state),
		ApplicationConfigBool: true,
		ApplicationConfigRv:   &mc.MockApplication{CapabilitiesRv: &mc.MockApplicationCapabilities{}},
	}).NewSystemChaincodeProvider().(*scc.MocksccProviderImpl)

	v := newValidationInstance(state)

	mockAclProvider := &aclmocks.MockACLProvider{}
	lccc := lscc.New(nil, mp, mockAclProvider, mockMSPIDGetter, &mockPolicyChecker{})
	stublccc := shimtest.NewMockStub("lscc", lccc)
	state["lscc"] = stublccc.State

	ccname := "mycc"
	ccver := "1"

	defaultPolicy, err := getSignedByMSPAdminPolicy(mspid)
	assert.NoError(t, err)
	res, err := createCCDataRWset(ccname, ccname, ccver, defaultPolicy)
	assert.NoError(t, err)

	tx, err := createLSCCTx(ccname, ccver, lscc.DEPLOY, res)
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

	b := &common.Block{Data: &common.BlockData{Data: [][]byte{envBytes}}}
	err = v.Validate(b, "lscc", 0, 0, policy)
	assert.NoError(t, err)
}

func TestValidateDeployNOK(t *testing.T) {
	var testCases = []struct {
		description string
		ccName      string
		ccVersion   string
		errMsg      string
	}{
		{description: "empty cc name", ccName: "", ccVersion: "1", errMsg: "invalid chaincode name ''"},
		{description: "bad first character in cc name", ccName: "_badname", ccVersion: "1.2", errMsg: "invalid chaincode name '_badname'"},
		{description: "bad character in cc name", ccName: "bad.name", ccVersion: "1-5", errMsg: "invalid chaincode name 'bad.name'"},
		{description: "empty cc version", ccName: "1good_name", ccVersion: "", errMsg: "invalid chaincode version ''"},
		{description: "bad cc version", ccName: "good-name", ccVersion: "$badversion", errMsg: "invalid chaincode version '$badversion'"},
		{description: "use system cc name", ccName: "qscc", ccVersion: "2.1", errMsg: "chaincode name 'qscc' is reserved for system chaincodes"},
	}

	
	state := make(map[string]map[string][]byte)
	mp := (&scc.MocksccProviderFactory{
		Qe:                    lm.NewMockQueryExecutor(state),
		ApplicationConfigBool: true,
		ApplicationConfigRv:   &mc.MockApplication{CapabilitiesRv: &mc.MockApplicationCapabilities{}},
	}).NewSystemChaincodeProvider().(*scc.MocksccProviderImpl)

	v := newValidationInstance(state)

	mockAclProvider := &aclmocks.MockACLProvider{}
	lccc := lscc.New(nil, mp, mockAclProvider, mockMSPIDGetter, &mockPolicyChecker{})
	stublccc := shimtest.NewMockStub("lscc", lccc)
	state["lscc"] = stublccc.State

	policy, err := getSignedByMSPAdminPolicy(mspid)
	assert.NoError(t, err)

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			testChaincodeDeployNOK(t, tc.ccName, tc.ccVersion, tc.errMsg, v, policy)
		})
	}
}

func testChaincodeDeployNOK(t *testing.T, ccName, ccVersion, errMsg string, v *Validator, policy []byte) {
	res, err := createCCDataRWset(ccName, ccName, ccVersion, policy)
	assert.NoError(t, err)

	tx, err := createLSCCTx(ccName, ccVersion, lscc.DEPLOY, res)
	if err != nil {
		t.Fatalf("createTx returned err %s", err)
	}

	envBytes, err := protoutil.GetBytesEnvelope(tx)
	if err != nil {
		t.Fatalf("GetBytesEnvelope returned err %s", err)
	}

	b := &common.Block{Data: &common.BlockData{Data: [][]byte{envBytes}}, Header: &common.BlockHeader{}}
	err = v.Validate(b, "lscc", 0, 0, policy)
	assert.EqualError(t, err, errMsg)
}

func TestValidateDeployWithCollection(t *testing.T) {
	state := make(map[string]map[string][]byte)
	mp := (&scc.MocksccProviderFactory{
		Qe:                    lm.NewMockQueryExecutor(state),
		ApplicationConfigBool: true,
		ApplicationConfigRv: &mc.MockApplication{CapabilitiesRv: &mc.MockApplicationCapabilities{
			PrivateChannelDataRv: true,
		}},
	}).NewSystemChaincodeProvider().(*scc.MocksccProviderImpl)

	qec := &mocks2.QueryExecutorCreator{}
	qec.On("NewQueryExecutor").Return(lm.NewMockQueryExecutor(state), nil)
	v := newCustomValidationInstance(qec, &mc.MockApplicationCapabilities{
		PrivateChannelDataRv: true,
	})

	mockAclProvider := &aclmocks.MockACLProvider{}
	lccc := lscc.New(nil, mp, mockAclProvider, mockMSPIDGetter, &mockPolicyChecker{})
	stublccc := shimtest.NewMockStub("lscc", lccc)
	state["lscc"] = stublccc.State

	r := stublccc.MockInit("1", [][]byte{})
	if r.Status != shim.OK {
		fmt.Println("Init failed", string(r.Message))
		t.FailNow()
	}

	ccname := "mycc"
	ccver := "1"

	collName1 := "mycollection1"
	collName2 := "mycollection2"
	var signers = [][]byte{[]byte("signer0"), []byte("signer1")}
	policyEnvelope := cauthdsl.Envelope(cauthdsl.Or(cauthdsl.SignedBy(0), cauthdsl.SignedBy(1)), signers)
	var requiredPeerCount, maximumPeerCount int32
	var blockToLive uint64
	requiredPeerCount = 1
	maximumPeerCount = 2
	blockToLive = 1000
	coll1 := createCollectionConfig(collName1, policyEnvelope, requiredPeerCount, maximumPeerCount, blockToLive)
	coll2 := createCollectionConfig(collName2, policyEnvelope, requiredPeerCount, maximumPeerCount, blockToLive)

	
	ccp := &common.CollectionConfigPackage{Config: []*common.CollectionConfig{coll1, coll2}}
	ccpBytes, err := proto.Marshal(ccp)
	assert.NoError(t, err)
	assert.NotNil(t, ccpBytes)

	defaultPolicy, err := getSignedByMSPAdminPolicy(mspid)
	assert.NoError(t, err)
	res, err := createCCDataRWsetWithCollection(ccname, ccname, ccver, defaultPolicy, ccpBytes)
	assert.NoError(t, err)

	tx, err := createLSCCTxWithCollection(ccname, ccver, lscc.DEPLOY, res, defaultPolicy, ccpBytes)
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

	b := &common.Block{Data: &common.BlockData{Data: [][]byte{envBytes}}}
	err = v.Validate(b, "lscc", 0, 0, policy)
	assert.NoError(t, err)

	
	
	ccp = &common.CollectionConfigPackage{Config: []*common.CollectionConfig{coll1, coll2, coll1}}
	ccpBytes, err = proto.Marshal(ccp)
	assert.NoError(t, err)
	assert.NotNil(t, ccpBytes)

	res, err = createCCDataRWsetWithCollection(ccname, ccname, ccver, defaultPolicy, ccpBytes)
	assert.NoError(t, err)

	tx, err = createLSCCTxWithCollection(ccname, ccver, lscc.DEPLOY, res, defaultPolicy, ccpBytes)
	if err != nil {
		t.Fatalf("createTx returned err %s", err)
	}

	envBytes, err = protoutil.GetBytesEnvelope(tx)
	if err != nil {
		t.Fatalf("GetBytesEnvelope returned err %s", err)
	}

	b = &common.Block{Data: &common.BlockData{Data: [][]byte{envBytes}}}
	err = v.Validate(b, "lscc", 0, 0, policy)
	assert.NoError(t, err)

	
	state = make(map[string]map[string][]byte)
	mp = (&scc.MocksccProviderFactory{
		Qe:                    lm.NewMockQueryExecutor(state),
		ApplicationConfigBool: true,
		ApplicationConfigRv: &mc.MockApplication{CapabilitiesRv: &mc.MockApplicationCapabilities{
			PrivateChannelDataRv: true,
			V1_2ValidationRv:     true,
		}},
	}).NewSystemChaincodeProvider().(*scc.MocksccProviderImpl)

	v = newValidationInstance(state)
	lccc = lscc.New(nil, mp, mockAclProvider, mockMSPIDGetter, &mockPolicyChecker{})
	stublccc = shimtest.NewMockStub("lscc", lccc)
	state["lscc"] = stublccc.State

	r = stublccc.MockInit("1", [][]byte{})
	if r.Status != shim.OK {
		fmt.Println("Init failed", string(r.Message))
		t.FailNow()
	}

	b = &common.Block{Data: &common.BlockData{Data: [][]byte{envBytes}}}
	err = v.Validate(b, "lscc", 0, 0, policy)
}

func TestValidateDeployWithPolicies(t *testing.T) {
	state := make(map[string]map[string][]byte)
	mp := (&scc.MocksccProviderFactory{
		Qe:                    lm.NewMockQueryExecutor(state),
		ApplicationConfigBool: true,
		ApplicationConfigRv:   &mc.MockApplication{CapabilitiesRv: &mc.MockApplicationCapabilities{}},
	}).NewSystemChaincodeProvider().(*scc.MocksccProviderImpl)

	v := newValidationInstance(state)

	mockAclProvider := &aclmocks.MockACLProvider{}
	lccc := lscc.New(nil, mp, mockAclProvider, mockMSPIDGetter, &mockPolicyChecker{})
	stublccc := shimtest.NewMockStub("lscc", lccc)
	state["lscc"] = stublccc.State

	ccname := "mycc"
	ccver := "1"

	
	
	

	res, err := createCCDataRWset(ccname, ccname, ccver, cauthdsl.MarshaledAcceptAllPolicy)
	assert.NoError(t, err)

	tx, err := createLSCCTx(ccname, ccver, lscc.DEPLOY, res)
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

	b := &common.Block{Data: &common.BlockData{Data: [][]byte{envBytes}}}
	err = v.Validate(b, "lscc", 0, 0, policy)
	assert.NoError(t, err)

	
	
	

	res, err = createCCDataRWset(ccname, ccname, ccver, cauthdsl.MarshaledRejectAllPolicy)
	assert.NoError(t, err)

	tx, err = createLSCCTx(ccname, ccver, lscc.DEPLOY, res)
	if err != nil {
		t.Fatalf("createTx returned err %s", err)
	}

	envBytes, err = protoutil.GetBytesEnvelope(tx)
	if err != nil {
		t.Fatalf("GetBytesEnvelope returned err %s", err)
	}

	
	policy, err = getSignedByMSPMemberPolicy(mspid)
	if err != nil {
		t.Fatalf("failed getting policy, err %s", err)
	}

	b = &common.Block{Data: &common.BlockData{Data: [][]byte{envBytes}}}
	err = v.Validate(b, "lscc", 0, 0, policy)
	assert.Error(t, err)
}

func TestInvalidUpgrade(t *testing.T) {
	state := make(map[string]map[string][]byte)
	mp := (&scc.MocksccProviderFactory{
		Qe:                    lm.NewMockQueryExecutor(state),
		ApplicationConfigBool: true,
		ApplicationConfigRv:   &mc.MockApplication{CapabilitiesRv: &mc.MockApplicationCapabilities{}},
	}).NewSystemChaincodeProvider().(*scc.MocksccProviderImpl)

	v := newValidationInstance(state)

	mockAclProvider := &aclmocks.MockACLProvider{}
	lccc := lscc.New(nil, mp, mockAclProvider, mockMSPIDGetter, &mockPolicyChecker{})
	stublccc := shimtest.NewMockStub("lscc", lccc)
	state["lscc"] = stublccc.State

	ccname := "mycc"
	ccver := "2"

	simresres, err := createCCDataRWset(ccname, ccname, ccver, nil)
	assert.NoError(t, err)

	tx, err := createLSCCTx(ccname, ccver, lscc.UPGRADE, simresres)
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

	b := &common.Block{Data: &common.BlockData{Data: [][]byte{envBytes}}}
	err = v.Validate(b, "lscc", 0, 0, policy)
	assert.Error(t, err)
}

func TestValidateUpgradeOK(t *testing.T) {
	state := make(map[string]map[string][]byte)
	mp := (&scc.MocksccProviderFactory{
		Qe:                    lm.NewMockQueryExecutor(state),
		ApplicationConfigBool: true,
		ApplicationConfigRv:   &mc.MockApplication{CapabilitiesRv: &mc.MockApplicationCapabilities{}},
	}).NewSystemChaincodeProvider().(*scc.MocksccProviderImpl)

	v := newValidationInstance(state)

	mockAclProvider := &aclmocks.MockACLProvider{}
	lccc := lscc.New(nil, mp, mockAclProvider, mockMSPIDGetter, &mockPolicyChecker{})
	stublccc := shimtest.NewMockStub("lscc", lccc)
	state["lscc"] = stublccc.State

	ccname := "mycc"
	ccver := "upgradeok"
	path := "mychaincode"

	cds, err := constructDeploymentSpec(ccname, path, ccver, [][]byte{[]byte("init"), []byte("a"), []byte("100"), []byte("b"), []byte("200")}, true)
	if err != nil {
		fmt.Printf("%s\n", err)
		t.FailNow()
	}

	var b []byte
	if b, err = proto.Marshal(cds); err != nil || b == nil {
		t.FailNow()
	}

	sProp2, _ := protoutil.MockSignedEndorserProposal2OrPanic(chainId, &peer.ChaincodeSpec{}, id)
	args := [][]byte{[]byte("deploy"), []byte(ccname), b}
	if res := stublccc.MockInvokeWithSignedProposal("1", args, sProp2); res.Status != shim.OK {
		fmt.Printf("%#v\n", res)
		t.FailNow()
	}

	ccver = "2"

	simresres, err := createCCDataRWset(ccname, ccname, ccver, nil)
	assert.NoError(t, err)

	tx, err := createLSCCTx(ccname, ccver, lscc.UPGRADE, simresres)
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

	bl := &common.Block{Data: &common.BlockData{Data: [][]byte{envBytes}}}
	err = v.Validate(bl, "lscc", 0, 0, policy)
	assert.NoError(t, err)
}

func TestInvalidateUpgradeBadVersion(t *testing.T) {
	state := make(map[string]map[string][]byte)
	mp := (&scc.MocksccProviderFactory{
		Qe:                    lm.NewMockQueryExecutor(state),
		ApplicationConfigBool: true,
		ApplicationConfigRv:   &mc.MockApplication{CapabilitiesRv: &mc.MockApplicationCapabilities{}},
	}).NewSystemChaincodeProvider().(*scc.MocksccProviderImpl)

	v := newValidationInstance(state)

	mockAclProvider := &aclmocks.MockACLProvider{}
	lccc := lscc.New(nil, mp, mockAclProvider, mockMSPIDGetter, &mockPolicyChecker{})
	stublccc := shimtest.NewMockStub("lscc", lccc)
	state["lscc"] = stublccc.State

	ccname := "mycc"
	ccver := "upgradebadversion"
	path := "mychaincode"

	cds, err := constructDeploymentSpec(ccname, path, ccver, [][]byte{[]byte("init"), []byte("a"), []byte("100"), []byte("b"), []byte("200")}, true)
	if err != nil {
		fmt.Printf("%s\n", err)
		t.FailNow()
	}

	var b []byte
	if b, err = proto.Marshal(cds); err != nil || b == nil {
		t.FailNow()
	}

	sProp2, _ := protoutil.MockSignedEndorserProposal2OrPanic(chainId, &peer.ChaincodeSpec{}, id)
	args := [][]byte{[]byte("deploy"), []byte(ccname), b}
	if res := stublccc.MockInvokeWithSignedProposal("1", args, sProp2); res.Status != shim.OK {
		fmt.Printf("%#v\n", res)
		t.FailNow()
	}

	simresres, err := createCCDataRWset(ccname, ccname, ccver, nil)
	assert.NoError(t, err)

	tx, err := createLSCCTx(ccname, ccver, lscc.UPGRADE, simresres)
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

	bl := &common.Block{Data: &common.BlockData{Data: [][]byte{envBytes}}}
	err = v.Validate(bl, "lscc", 0, 0, policy)
	assert.Error(t, err)
}

func validateUpgradeWithCollection(t *testing.T, ccver string, V1_2Validation bool) {
	state := make(map[string]map[string][]byte)
	mp := (&scc.MocksccProviderFactory{
		Qe:                    lm.NewMockQueryExecutor(state),
		ApplicationConfigBool: true,
		ApplicationConfigRv: &mc.MockApplication{CapabilitiesRv: &mc.MockApplicationCapabilities{
			PrivateChannelDataRv: true,
			V1_2ValidationRv:     V1_2Validation,
		}},
	}).NewSystemChaincodeProvider().(*scc.MocksccProviderImpl)

	qec := &mocks2.QueryExecutorCreator{}
	qec.On("NewQueryExecutor").Return(lm.NewMockQueryExecutor(state), nil)
	v := newCustomValidationInstance(qec, &mc.MockApplicationCapabilities{
		PrivateChannelDataRv: true,
		V1_2ValidationRv:     V1_2Validation,
	})

	mockAclProvider := &aclmocks.MockACLProvider{}
	lccc := lscc.New(nil, mp, mockAclProvider, mockMSPIDGetter, &mockPolicyChecker{})
	stublccc := shimtest.NewMockStub("lscc", lccc)
	state["lscc"] = stublccc.State

	r := stublccc.MockInit("1", [][]byte{})
	if r.Status != shim.OK {
		fmt.Println("Init failed", string(r.Message))
		t.FailNow()
	}

	ccname := "mycc"
	path := "mychaincode"

	cds, err := constructDeploymentSpec(ccname, path, ccver, [][]byte{[]byte("init"), []byte("a"), []byte("100"), []byte("b"), []byte("200")}, true)
	if err != nil {
		fmt.Printf("%s\n", err)
		t.FailNow()
	}

	var b []byte
	if b, err = proto.Marshal(cds); err != nil || b == nil {
		t.FailNow()
	}

	sProp2, _ := protoutil.MockSignedEndorserProposal2OrPanic(chainId, &peer.ChaincodeSpec{}, id)
	args := [][]byte{[]byte("deploy"), []byte(ccname), b}
	if res := stublccc.MockInvokeWithSignedProposal("1", args, sProp2); res.Status != shim.OK {
		fmt.Printf("%#v\n", res)
		t.FailNow()
	}

	ccver = "2"

	collName1 := "mycollection1"
	collName2 := "mycollection2"
	var signers = [][]byte{[]byte("signer0"), []byte("signer1")}
	policyEnvelope := cauthdsl.Envelope(cauthdsl.Or(cauthdsl.SignedBy(0), cauthdsl.SignedBy(1)), signers)
	var requiredPeerCount, maximumPeerCount int32
	var blockToLive uint64
	requiredPeerCount = 1
	maximumPeerCount = 2
	blockToLive = 1000
	coll1 := createCollectionConfig(collName1, policyEnvelope, requiredPeerCount, maximumPeerCount, blockToLive)
	coll2 := createCollectionConfig(collName2, policyEnvelope, requiredPeerCount, maximumPeerCount, blockToLive)

	
	
	
	
	ccp := &common.CollectionConfigPackage{Config: []*common.CollectionConfig{coll1, coll2}}
	ccpBytes, err := proto.Marshal(ccp)
	assert.NoError(t, err)
	assert.NotNil(t, ccpBytes)

	defaultPolicy, err := getSignedByMSPAdminPolicy(mspid)
	assert.NoError(t, err)
	res, err := createCCDataRWsetWithCollection(ccname, ccname, ccver, defaultPolicy, ccpBytes)
	assert.NoError(t, err)

	tx, err := createLSCCTxWithCollection(ccname, ccver, lscc.UPGRADE, res, defaultPolicy, ccpBytes)
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

	bl := &common.Block{Data: &common.BlockData{Data: [][]byte{envBytes}}}
	err = v.Validate(bl, "lscc", 0, 0, policy)
	if V1_2Validation {
		assert.NoError(t, err)
	} else {
		assert.Error(t, err, "LSCC can only issue a single putState upon deploy/upgrade")
	}

	state["lscc"][privdata.BuildCollectionKVSKey(ccname)] = ccpBytes

	if V1_2Validation {
		ccver = "3"

		collName3 := "mycollection3"
		coll3 := createCollectionConfig(collName3, policyEnvelope, requiredPeerCount, maximumPeerCount, blockToLive)

		
		
		ccp = &common.CollectionConfigPackage{Config: []*common.CollectionConfig{coll3}}
		ccpBytes, err = proto.Marshal(ccp)
		assert.NoError(t, err)
		assert.NotNil(t, ccpBytes)

		res, err = createCCDataRWsetWithCollection(ccname, ccname, ccver, defaultPolicy, ccpBytes)
		assert.NoError(t, err)

		tx, err = createLSCCTxWithCollection(ccname, ccver, lscc.UPGRADE, res, defaultPolicy, ccpBytes)
		if err != nil {
			t.Fatalf("createTx returned err %s", err)
		}

		envBytes, err = protoutil.GetBytesEnvelope(tx)
		if err != nil {
			t.Fatalf("GetBytesEnvelope returned err %s", err)
		}

		bl = &common.Block{Data: &common.BlockData{Data: [][]byte{envBytes}}}
		err = v.Validate(bl, "lscc", 0, 0, policy)
		assert.Error(t, err, "Some existing collection configurations are missing in the new collection configuration package")

		ccver = "3"

		
		
		ccp = &common.CollectionConfigPackage{Config: []*common.CollectionConfig{coll1, coll3}}
		ccpBytes, err = proto.Marshal(ccp)
		assert.NoError(t, err)
		assert.NotNil(t, ccpBytes)

		res, err = createCCDataRWsetWithCollection(ccname, ccname, ccver, defaultPolicy, ccpBytes)
		assert.NoError(t, err)

		tx, err = createLSCCTxWithCollection(ccname, ccver, lscc.UPGRADE, res, defaultPolicy, ccpBytes)
		if err != nil {
			t.Fatalf("createTx returned err %s", err)
		}

		envBytes, err = protoutil.GetBytesEnvelope(tx)
		if err != nil {
			t.Fatalf("GetBytesEnvelope returned err %s", err)
		}

		args = [][]byte{[]byte("dv"), envBytes, policy, ccpBytes}
		bl = &common.Block{Data: &common.BlockData{Data: [][]byte{envBytes}}}
		err = v.Validate(bl, "lscc", 0, 0, policy)
		assert.Error(t, err, "existing collection named mycollection2 is missing in the new collection configuration package")

		ccver = "3"

		
		ccp = &common.CollectionConfigPackage{Config: []*common.CollectionConfig{coll1, coll2, coll3}}
		ccpBytes, err = proto.Marshal(ccp)
		assert.NoError(t, err)
		assert.NotNil(t, ccpBytes)

		res, err = createCCDataRWsetWithCollection(ccname, ccname, ccver, defaultPolicy, ccpBytes)
		assert.NoError(t, err)

		tx, err = createLSCCTxWithCollection(ccname, ccver, lscc.UPGRADE, res, defaultPolicy, ccpBytes)
		if err != nil {
			t.Fatalf("createTx returned err %s", err)
		}

		envBytes, err = protoutil.GetBytesEnvelope(tx)
		if err != nil {
			t.Fatalf("GetBytesEnvelope returned err %s", err)
		}

		bl = &common.Block{Data: &common.BlockData{Data: [][]byte{envBytes}}}
		err = v.Validate(bl, "lscc", 0, 0, policy)
		assert.NoError(t, err)
	}
}

func TestValidateUpgradeWithCollection(t *testing.T) {
	
	validateUpgradeWithCollection(t, "v12-validation-enabled", true)
	
	validateUpgradeWithCollection(t, "v12-validation-disabled", false)
}

func TestValidateUpgradeWithPoliciesOK(t *testing.T) {
	state := make(map[string]map[string][]byte)
	mp := (&scc.MocksccProviderFactory{
		Qe:                    lm.NewMockQueryExecutor(state),
		ApplicationConfigBool: true,
		ApplicationConfigRv:   &mc.MockApplication{CapabilitiesRv: &mc.MockApplicationCapabilities{}},
	}).NewSystemChaincodeProvider().(*scc.MocksccProviderImpl)

	v := newValidationInstance(state)

	mockAclProvider := &aclmocks.MockACLProvider{}
	lccc := lscc.New(nil, mp, mockAclProvider, mockMSPIDGetter, &mockPolicyChecker{})
	stublccc := shimtest.NewMockStub("lscc", lccc)
	state["lscc"] = stublccc.State

	ccname := "mycc"
	ccver := "upgradewithpoliciesok"
	path := "mychaincode"

	cds, err := constructDeploymentSpec(ccname, path, ccver, [][]byte{[]byte("init"), []byte("a"), []byte("100"), []byte("b"), []byte("200")}, false)
	if err != nil {
		fmt.Printf("%s\n", err)
		t.FailNow()
	}
	_, err = processSignedCDS(cds, cauthdsl.AcceptAllPolicy)
	assert.NoError(t, err)

	var b []byte
	if b, err = proto.Marshal(cds); err != nil || b == nil {
		t.FailNow()
	}

	sProp2, _ := protoutil.MockSignedEndorserProposal2OrPanic(chainId, &peer.ChaincodeSpec{}, id)
	args := [][]byte{[]byte("deploy"), []byte(ccname), b}
	if res := stublccc.MockInvokeWithSignedProposal("1", args, sProp2); res.Status != shim.OK {
		fmt.Printf("%#v\n", res)
		t.FailNow()
	}

	ccver = "2"

	simresres, err := createCCDataRWset(ccname, ccname, ccver, nil)
	assert.NoError(t, err)

	tx, err := createLSCCTx(ccname, ccver, lscc.UPGRADE, simresres)
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

	args = [][]byte{[]byte("dv"), envBytes, policy}
	bl := &common.Block{Data: &common.BlockData{Data: [][]byte{envBytes}}}
	err = v.Validate(bl, "lscc", 0, 0, policy)
	assert.NoError(t, err)
}

func TestValidateUpgradeWithNewFailAllIP(t *testing.T) {
	
	
	
	
	
	
	
	

	validateUpgradeWithNewFailAllIP(t, "v11-capabilityenabled", true, true)
	validateUpgradeWithNewFailAllIP(t, "v11-capabilitydisabled", false, false)
}

func validateUpgradeWithNewFailAllIP(t *testing.T, ccver string, v11capability, expecterr bool) {
	state := make(map[string]map[string][]byte)
	mp := (&scc.MocksccProviderFactory{
		Qe:                    lm.NewMockQueryExecutor(state),
		ApplicationConfigBool: true,
		ApplicationConfigRv:   &mc.MockApplication{CapabilitiesRv: &mc.MockApplicationCapabilities{V1_1ValidationRv: v11capability}},
	}).NewSystemChaincodeProvider().(*scc.MocksccProviderImpl)

	qec := &mocks2.QueryExecutorCreator{}
	capabilities := &mc.MockApplicationCapabilities{}
	if v11capability {
		capabilities.V1_1ValidationRv = true
	}
	qec.On("NewQueryExecutor").Return(lm.NewMockQueryExecutor(state), nil)
	v := newCustomValidationInstance(qec, capabilities)

	mockAclProvider := &aclmocks.MockACLProvider{}
	lccc := lscc.New(nil, mp, mockAclProvider, mockMSPIDGetter, &mockPolicyChecker{})
	stublccc := shimtest.NewMockStub("lscc", lccc)
	state["lscc"] = stublccc.State

	

	ccname := "mycc"
	path := "mychaincode"

	cds, err := constructDeploymentSpec(ccname, path, ccver, [][]byte{[]byte("init"), []byte("a"), []byte("100"), []byte("b"), []byte("200")}, false)
	if err != nil {
		fmt.Printf("%s\n", err)
		t.FailNow()
	}
	_, err = processSignedCDS(cds, cauthdsl.AcceptAllPolicy)
	assert.NoError(t, err)

	var b []byte
	if b, err = proto.Marshal(cds); err != nil || b == nil {
		t.FailNow()
	}

	sProp2, _ := protoutil.MockSignedEndorserProposal2OrPanic(chainId, &peer.ChaincodeSpec{}, id)
	args := [][]byte{[]byte("deploy"), []byte(ccname), b}
	if res := stublccc.MockInvokeWithSignedProposal("1", args, sProp2); res.Status != shim.OK {
		fmt.Printf("%#v\n", res)
		t.FailNow()
	}

	

	

	ccver = ccver + ".2"

	simresres, err := createCCDataRWset(ccname, ccname, ccver,
		cauthdsl.MarshaledRejectAllPolicy, 
	)
	assert.NoError(t, err)

	tx, err := createLSCCTx(ccname, ccver, lscc.UPGRADE, simresres)
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

	
	if expecterr {
		bl := &common.Block{Data: &common.BlockData{Data: [][]byte{envBytes}}}
		err = v.Validate(bl, "lscc", 0, 0, policy)
		assert.Error(t, err)
	} else {
		bl := &common.Block{Data: &common.BlockData{Data: [][]byte{envBytes}}}
		err = v.Validate(bl, "lscc", 0, 0, policy)
		assert.NoError(t, err)
	}
}

func TestValidateUpgradeWithPoliciesFail(t *testing.T) {
	state := make(map[string]map[string][]byte)
	mp := (&scc.MocksccProviderFactory{
		Qe:                    lm.NewMockQueryExecutor(state),
		ApplicationConfigBool: true,
		ApplicationConfigRv:   &mc.MockApplication{CapabilitiesRv: &mc.MockApplicationCapabilities{}},
	}).NewSystemChaincodeProvider().(*scc.MocksccProviderImpl)

	v := newValidationInstance(state)

	mockAclProvider := &aclmocks.MockACLProvider{}
	lccc := lscc.New(nil, mp, mockAclProvider, mockMSPIDGetter, &mockPolicyChecker{})
	stublccc := shimtest.NewMockStub("lscc", lccc)
	state["lscc"] = stublccc.State

	ccname := "mycc"
	ccver := "upgradewithpoliciesfail"
	path := "mychaincode"

	cds, err := constructDeploymentSpec(ccname, path, ccver, [][]byte{[]byte("init"), []byte("a"), []byte("100"), []byte("b"), []byte("200")}, false)
	if err != nil {
		fmt.Printf("%s\n", err)
		t.FailNow()
	}
	cdbytes, err := processSignedCDS(cds, cauthdsl.RejectAllPolicy)
	assert.NoError(t, err)

	var b []byte
	if b, err = proto.Marshal(cds); err != nil || b == nil {
		t.FailNow()
	}

	
	
	stublccc.MockTransactionStart("barf")
	err = stublccc.PutState(ccname, cdbytes)
	assert.NoError(t, err)
	stublccc.MockTransactionEnd("barf")

	ccver = "2"

	simresres, err := createCCDataRWset(ccname, ccname, ccver, nil)
	assert.NoError(t, err)

	tx, err := createLSCCTx(ccname, ccver, lscc.UPGRADE, simresres)
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

	bl := &common.Block{Data: &common.BlockData{Data: [][]byte{envBytes}}}
	err = v.Validate(bl, "lscc", 0, 0, policy)
	assert.Error(t, err)
}

var id msp.SigningIdentity
var sid []byte
var mspid string
var chainId string = util.GetTestChannelID()

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

func createCollectionConfig(collectionName string, signaturePolicyEnvelope *common.SignaturePolicyEnvelope,
	requiredPeerCount int32, maximumPeerCount int32, blockToLive uint64,
) *common.CollectionConfig {
	signaturePolicy := &common.CollectionPolicyConfig_SignaturePolicy{
		SignaturePolicy: signaturePolicyEnvelope,
	}
	accessPolicy := &common.CollectionPolicyConfig{
		Payload: signaturePolicy,
	}

	return &common.CollectionConfig{
		Payload: &common.CollectionConfig_StaticCollectionConfig{
			StaticCollectionConfig: &common.StaticCollectionConfig{
				Name:              collectionName,
				MemberOrgsPolicy:  accessPolicy,
				RequiredPeerCount: requiredPeerCount,
				MaximumPeerCount:  maximumPeerCount,
				BlockToLive:       blockToLive,
			},
		},
	}
}

func testValidateCollection(t *testing.T, v *Validator, collectionConfigs []*common.CollectionConfig, cdRWSet *ccprovider.ChaincodeData,
	lsccFunc string, ac channelconfig.ApplicationCapabilities, chid string,
) error {
	ccp := &common.CollectionConfigPackage{Config: collectionConfigs}
	ccpBytes, err := proto.Marshal(ccp)
	assert.NoError(t, err)
	assert.NotNil(t, ccpBytes)

	lsccargs := [][]byte{nil, nil, nil, nil, nil, ccpBytes}
	rwset := &kvrwset.KVRWSet{Writes: []*kvrwset.KVWrite{{Key: cdRWSet.Name}, {Key: privdata.BuildCollectionKVSKey(cdRWSet.Name), Value: ccpBytes}}}

	err = v.validateRWSetAndCollection(rwset, cdRWSet, lsccargs, lsccFunc, ac, chid)
	return err

}

func TestValidateRWSetAndCollectionForDeploy(t *testing.T) {
	var err error
	chid := "ch"
	ccid := "mycc"
	ccver := "1.0"
	cdRWSet := &ccprovider.ChaincodeData{Name: ccid, Version: ccver}

	state := make(map[string]map[string][]byte)
	state["lscc"] = make(map[string][]byte)

	v := newValidationInstance(state)

	ac := capabilities.NewApplicationProvider(map[string]*common.Capability{
		capabilities.ApplicationV1_1: {},
	})

	lsccFunc := lscc.DEPLOY
	
	rwset := &kvrwset.KVRWSet{Writes: []*kvrwset.KVWrite{{Key: ccid}, {Key: "b"}, {Key: "c"}}}
	err = v.validateRWSetAndCollection(rwset, cdRWSet, nil, lsccFunc, ac, chid)
	assert.EqualError(t, err, "LSCC can only issue one or two putState upon deploy")

	
	rwset = &kvrwset.KVRWSet{Writes: []*kvrwset.KVWrite{{Key: ccid}, {Key: "b"}}}
	err = v.validateRWSetAndCollection(rwset, cdRWSet, nil, lsccFunc, ac, chid)
	assert.EqualError(t, err, "invalid key for the collection of chaincode mycc:1.0; expected 'mycc~collection', received 'b'")

	
	rwset = &kvrwset.KVRWSet{Writes: []*kvrwset.KVWrite{{Key: ccid}}}
	err = v.validateRWSetAndCollection(rwset, cdRWSet, nil, lsccFunc, ac, chid)
	assert.NoError(t, err)

	lsccargs := [][]byte{nil, nil, nil, nil, nil, nil}
	err = v.validateRWSetAndCollection(rwset, cdRWSet, lsccargs, lsccFunc, ac, chid)
	assert.NoError(t, err)

	
	rwset = &kvrwset.KVRWSet{Writes: []*kvrwset.KVWrite{{Key: ccid}, {Key: privdata.BuildCollectionKVSKey(ccid)}}}
	err = v.validateRWSetAndCollection(rwset, cdRWSet, lsccargs, lsccFunc, ac, chid)
	assert.NoError(t, err)

	
	lsccargs = [][]byte{nil, nil, nil, nil, nil, []byte("barf")}
	err = v.validateRWSetAndCollection(rwset, cdRWSet, lsccargs, lsccFunc, ac, chid)
	assert.EqualError(t, err, "collection configuration arguments supplied for chaincode mycc:1.0 do not match the configuration in the lscc writeset")

	
	rwset = &kvrwset.KVRWSet{Writes: []*kvrwset.KVWrite{{Key: ccid}, {Key: privdata.BuildCollectionKVSKey("mycc"), Value: []byte("barf")}}}
	err = v.validateRWSetAndCollection(rwset, cdRWSet, lsccargs, lsccFunc, ac, chid)
	assert.EqualError(t, err, "invalid collection configuration supplied for chaincode mycc:1.0")

	
	collName1 := "mycollection1"
	collName2 := "mycollection2"
	var signers = [][]byte{[]byte("signer0"), []byte("signer1")}
	policyEnvelope := cauthdsl.Envelope(cauthdsl.Or(cauthdsl.SignedBy(0), cauthdsl.SignedBy(1)), signers)
	var requiredPeerCount, maximumPeerCount int32
	var blockToLive uint64
	requiredPeerCount = 1
	maximumPeerCount = 2
	blockToLive = 10000
	coll1 := createCollectionConfig(collName1, policyEnvelope, requiredPeerCount, maximumPeerCount, blockToLive)
	coll2 := createCollectionConfig(collName2, policyEnvelope, requiredPeerCount, maximumPeerCount, blockToLive)

	err = testValidateCollection(t, v, []*common.CollectionConfig{coll1, coll2}, cdRWSet, lsccFunc, ac, chid)
	assert.NoError(t, err)

	
	err = testValidateCollection(t, v, []*common.CollectionConfig{coll1, coll2, coll1}, cdRWSet, lsccFunc, ac, chid)
	assert.NoError(t, err)

	
	collName3 := "mycollection3"
	requiredPeerCount = 2
	maximumPeerCount = 1
	blockToLive = 10000
	coll3 := createCollectionConfig(collName3, policyEnvelope, requiredPeerCount, maximumPeerCount, blockToLive)
	err = testValidateCollection(t, v, []*common.CollectionConfig{coll1, coll2, coll3}, cdRWSet, lsccFunc, ac, chid)
	assert.NoError(t, err)

	
	ac = capabilities.NewApplicationProvider(map[string]*common.Capability{
		capabilities.ApplicationV1_2: {},
	})

	
	err = testValidateCollection(t, v, []*common.CollectionConfig{coll1, coll2, coll1}, cdRWSet, lsccFunc, ac, chid)
	assert.EqualError(t, err, "collection-name: mycollection1 -- found duplicate collection configuration")

	
	requiredPeerCount = -2
	maximumPeerCount = 1
	blockToLive = 10000
	coll3 = createCollectionConfig(collName3, policyEnvelope, requiredPeerCount, maximumPeerCount, blockToLive)
	err = testValidateCollection(t, v, []*common.CollectionConfig{coll1, coll2, coll3}, cdRWSet, lsccFunc, ac, chid)
	assert.EqualError(t, err, "collection-name: mycollection3 -- requiredPeerCount (1) cannot be less than zero (-2)",
		collName3, maximumPeerCount, requiredPeerCount)

	
	requiredPeerCount = 2
	maximumPeerCount = 1
	blockToLive = 10000
	coll3 = createCollectionConfig(collName3, policyEnvelope, requiredPeerCount, maximumPeerCount, blockToLive)
	err = testValidateCollection(t, v, []*common.CollectionConfig{coll1, coll2, coll3}, cdRWSet, lsccFunc, ac, chid)
	assert.EqualError(t, err, "collection-name: mycollection3 -- maximum peer count (1) cannot be greater than the required peer count (2)")

	
	requiredPeerCount = 1
	maximumPeerCount = 2
	policyEnvelope = cauthdsl.Envelope(cauthdsl.And(cauthdsl.SignedBy(0), cauthdsl.SignedBy(1)), signers)
	coll3 = createCollectionConfig(collName3, policyEnvelope, requiredPeerCount, maximumPeerCount, blockToLive)
	err = testValidateCollection(t, v, []*common.CollectionConfig{coll3}, cdRWSet, lsccFunc, ac, chid)
	assert.EqualError(t, err, "collection-name: mycollection3 -- error in member org policy: signature policy is not an OR concatenation, NOutOf 2")

	
	ccp := &common.CollectionConfigPackage{Config: []*common.CollectionConfig{coll1}}
	ccpBytes, err := proto.Marshal(ccp)
	assert.NoError(t, err)
	state["lscc"][privdata.BuildCollectionKVSKey(ccid)] = ccpBytes
	err = testValidateCollection(t, v, []*common.CollectionConfig{coll1}, cdRWSet, lsccFunc, ac, chid)
	assert.EqualError(t, err, "collection data should not exist for chaincode mycc:1.0")
}

func TestValidateRWSetAndCollectionForUpgrade(t *testing.T) {
	chid := "ch"
	ccid := "mycc"
	ccver := "1.0"
	cdRWSet := &ccprovider.ChaincodeData{Name: ccid, Version: ccver}

	state := make(map[string]map[string][]byte)
	state["lscc"] = make(map[string][]byte)

	v := newValidationInstance(state)

	ac := capabilities.NewApplicationProvider(map[string]*common.Capability{
		capabilities.ApplicationV1_2: {},
	})

	lsccFunc := lscc.UPGRADE

	collName1 := "mycollection1"
	collName2 := "mycollection2"
	collName3 := "mycollection3"
	var signers = [][]byte{[]byte("signer0"), []byte("signer1")}
	policyEnvelope := cauthdsl.Envelope(cauthdsl.Or(cauthdsl.SignedBy(0), cauthdsl.SignedBy(1)), signers)
	var requiredPeerCount, maximumPeerCount int32
	var blockToLive uint64
	requiredPeerCount = 1
	maximumPeerCount = 2
	blockToLive = 3
	coll1 := createCollectionConfig(collName1, policyEnvelope, requiredPeerCount, maximumPeerCount, blockToLive)
	coll2 := createCollectionConfig(collName2, policyEnvelope, requiredPeerCount, maximumPeerCount, blockToLive)
	coll3 := createCollectionConfig(collName3, policyEnvelope, requiredPeerCount, maximumPeerCount, blockToLive)

	ccp := &common.CollectionConfigPackage{Config: []*common.CollectionConfig{coll1, coll2}}
	ccpBytes, err := proto.Marshal(ccp)
	assert.NoError(t, err)

	
	err = testValidateCollection(t, v, []*common.CollectionConfig{coll1}, cdRWSet, lsccFunc, ac, chid)
	assert.NoError(t, err)

	state["lscc"][privdata.BuildCollectionKVSKey(ccid)] = ccpBytes

	
	err = testValidateCollection(t, v, []*common.CollectionConfig{coll1, coll2}, cdRWSet, lsccFunc, ac, chid)
	assert.NoError(t, err)

	
	err = testValidateCollection(t, v, []*common.CollectionConfig{coll1}, cdRWSet, lsccFunc, ac, chid)
	assert.EqualError(t, err, "the following existing collections are missing in the new collection configuration package: [mycollection2]")

	
	err = testValidateCollection(t, v, []*common.CollectionConfig{coll1, coll3}, cdRWSet, lsccFunc, ac, chid)
	assert.EqualError(t, err, "the following existing collections are missing in the new collection configuration package: [mycollection2]")

	
	err = testValidateCollection(t, v, []*common.CollectionConfig{coll1, coll2, coll3}, cdRWSet, lsccFunc, ac, chid)
	assert.NoError(t, err)

	newBlockToLive := blockToLive + 1
	coll2 = createCollectionConfig(collName2, policyEnvelope, requiredPeerCount, maximumPeerCount, newBlockToLive)

	
	err = testValidateCollection(t, v, []*common.CollectionConfig{coll1, coll2, coll3}, cdRWSet, lsccFunc, ac, chid)
	assert.EqualError(t, err, "the BlockToLive in the following existing collections must not be modified: [mycollection2]")
}

var mockMSPIDGetter = func(cid string) []string {
	return []string{"SampleOrg"}
}

func TestMain(m *testing.M) {
	testDir, err := ioutil.TempDir("", "v1.2-validation")
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

	
	mspmgmt.XXXSetMSPManager("mycc", mspmgmt.GetManagerForChain(util.GetTestChannelID()))

	os.Exit(m.Run())
}

func TestInValidCollectionName(t *testing.T) {
	validNames := []string{"collection1", "collection_2"}
	inValidNames := []string{"collection.1", "collection%2", ""}

	for _, name := range validNames {
		assert.NoError(t, validateCollectionName(name), "Testing for name = "+name)
	}
	for _, name := range inValidNames {
		assert.Error(t, validateCollectionName(name), "Testing for name = "+name)
	}
}
