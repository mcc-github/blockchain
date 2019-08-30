/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package privdata

import (
	"testing"

	"github.com/mcc-github/blockchain-protos-go/common"
	"github.com/mcc-github/blockchain-protos-go/peer"
	"github.com/mcc-github/blockchain/common/cauthdsl"
	lm "github.com/mcc-github/blockchain/common/mocks/ledger"
	"github.com/mcc-github/blockchain/core/common/privdata/mock"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)






type queryExecutorFactory interface{ QueryExecutorFactory }
type chaincodeInfoProvider interface{ ChaincodeInfoProvider }
type identityDeserializerFactory interface{ IdentityDeserializerFactory }

func TestNewSimpleCollectionStore(t *testing.T) {
	mockQueryExecutorFactory := &mock.QueryExecutorFactory{}
	mockCCInfoProvider := &mock.ChaincodeInfoProvider{}

	cs := NewSimpleCollectionStore(mockQueryExecutorFactory, mockCCInfoProvider)
	assert.NotNil(t, cs)
	assert.Exactly(t, mockQueryExecutorFactory, cs.qeFactory)
	assert.Exactly(t, mockCCInfoProvider, cs.ccInfoProvider)
	assert.NotNil(t, cs.idDeserializerFactory)
}

func TestCollectionStore(t *testing.T) {
	mockQueryExecutorFactory := &mock.QueryExecutorFactory{}
	mockCCInfoProvider := &mock.ChaincodeInfoProvider{}
	mockIDDeserializerFactory := &mock.IdentityDeserializerFactory{}
	mockIDDeserializerFactory.GetIdentityDeserializerReturns(&mockDeserializer{})

	cs := &SimpleCollectionStore{
		qeFactory:             mockQueryExecutorFactory,
		ccInfoProvider:        mockCCInfoProvider,
		idDeserializerFactory: mockIDDeserializerFactory,
	}

	mockQueryExecutorFactory.NewQueryExecutorReturns(nil, errors.New("new-query-executor-failed"))
	_, err := cs.RetrieveCollection(common.CollectionCriteria{})
	assert.Contains(t, err.Error(), "could not retrieve query executor for collection criteria")

	mockQueryExecutorFactory.NewQueryExecutorReturns(&lm.MockQueryExecutor{}, nil)
	_, err = cs.retrieveCollectionConfigPackage(common.CollectionCriteria{Namespace: "non-existing-chaincode"}, nil)
	assert.EqualError(t, err, "Chaincode [non-existing-chaincode] does not exist")

	_, err = cs.RetrieveCollection(common.CollectionCriteria{})
	assert.Contains(t, err.Error(), "could not be found")

	ccr := common.CollectionCriteria{Channel: "ch", Namespace: "cc", Collection: "mycollection"}
	mockCCInfoProvider.CollectionInfoReturns(nil, errors.New("collection-info-error"))
	_, err = cs.RetrieveCollection(ccr)
	assert.EqualError(t, err, "collection-info-error")

	scc := &common.StaticCollectionConfig{Name: "mycollection"}
	mockCCInfoProvider.CollectionInfoReturns(scc, nil)
	_, err = cs.RetrieveCollection(ccr)
	assert.Contains(t, err.Error(), "error setting up collection for collection criteria")

	var signers = [][]byte{[]byte("signer0"), []byte("signer1")}
	policyEnvelope := cauthdsl.Envelope(cauthdsl.Or(cauthdsl.SignedBy(0), cauthdsl.SignedBy(1)), signers)
	accessPolicy := createCollectionPolicyConfig(policyEnvelope)

	scc = &common.StaticCollectionConfig{
		Name:             "mycollection",
		MemberOrgsPolicy: accessPolicy,
		MemberOnlyRead:   false,
		MemberOnlyWrite:  false,
	}

	mockCCInfoProvider.CollectionInfoReturns(scc, nil)
	c, err := cs.RetrieveCollection(ccr)
	assert.NoError(t, err)
	assert.NotNil(t, c)

	ca, err := cs.RetrieveCollectionAccessPolicy(ccr)
	assert.NoError(t, err)
	assert.NotNil(t, ca)

	scc = &common.StaticCollectionConfig{
		Name:             "mycollection",
		MemberOrgsPolicy: accessPolicy,
		MemberOnlyRead:   true,
		MemberOnlyWrite:  true,
	}
	cc := &common.CollectionConfig{
		Payload: &common.CollectionConfig_StaticCollectionConfig{StaticCollectionConfig: scc},
	}
	ccp := &common.CollectionConfigPackage{Config: []*common.CollectionConfig{cc}}

	mockCCInfoProvider.CollectionInfoReturns(scc, nil)
	mockCCInfoProvider.ChaincodeInfoReturns(
		&ledger.DeployedChaincodeInfo{ExplicitCollectionConfigPkg: ccp},
		nil,
	)

	ccc, err := cs.RetrieveCollectionConfigPackage(ccr)
	assert.NoError(t, err)
	assert.NotNil(t, ccc)

	signedProp, _ := protoutil.MockSignedEndorserProposalOrPanic("A", &peer.ChaincodeSpec{}, []byte("signer0"), []byte("msg1"))
	readP, writeP, err := cs.RetrieveReadWritePermission(ccr, signedProp, &lm.MockQueryExecutor{})
	assert.NoError(t, err)
	assert.True(t, readP)
	assert.True(t, writeP)

	
	signedProp, _ = protoutil.MockSignedEndorserProposalOrPanic("A", &peer.ChaincodeSpec{}, []byte("signer2"), []byte("msg1"))
	readP, writeP, err = cs.RetrieveReadWritePermission(ccr, signedProp, &lm.MockQueryExecutor{})
	assert.NoError(t, err)
	assert.False(t, readP)
	assert.False(t, writeP)

	scc = &common.StaticCollectionConfig{
		Name:             "mycollection",
		MemberOrgsPolicy: accessPolicy,
		MemberOnlyRead:   false,
		MemberOnlyWrite:  false,
	}
	mockCCInfoProvider.CollectionInfoReturns(scc, nil)

	
	signedProp, _ = protoutil.MockSignedEndorserProposalOrPanic("A", &peer.ChaincodeSpec{}, []byte("signer2"), []byte("msg1"))
	readP, writeP, err = cs.RetrieveReadWritePermission(ccr, signedProp, &lm.MockQueryExecutor{})
	assert.NoError(t, err)
	assert.True(t, readP)
	assert.True(t, writeP)
}
