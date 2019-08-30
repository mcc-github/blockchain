/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package cclifecycle_test

import (
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/mcc-github/blockchain-protos-go/common"
	"github.com/mcc-github/blockchain/common/chaincode"
	"github.com/mcc-github/blockchain/common/flogging/floggingtest"
	"github.com/mcc-github/blockchain/core/cclifecycle"
	"github.com/mcc-github/blockchain/core/cclifecycle/mocks"
	"github.com/mcc-github/blockchain/core/common/ccprovider"
	"github.com/mcc-github/blockchain/core/common/privdata"
	"github.com/mcc-github/blockchain/core/ledger/cceventmgmt"
	"github.com/mcc-github/blockchain/protoutil"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewQuery(t *testing.T) {
	
	var q cclifecycle.Query
	queryCreator := func() (cclifecycle.Query, error) {
		q := &mocks.Query{}
		q.On("Done")
		return q, nil
	}
	q, _ = cclifecycle.QueryCreatorFunc(queryCreator).NewQuery()
	q.Done()
}

func TestHandleMetadataUpdate(t *testing.T) {
	f := func(channel string, chaincodes chaincode.MetadataSet) {
		assert.Len(t, chaincodes, 2)
		assert.Equal(t, "mychannel", channel)
	}
	cclifecycle.HandleMetadataUpdateFunc(f).HandleMetadataUpdate("mychannel", chaincode.MetadataSet{{}, {}})
}

func TestEnumerate(t *testing.T) {
	f := func() ([]chaincode.InstalledChaincode, error) {
		return []chaincode.InstalledChaincode{{}, {}}, nil
	}
	ccs, err := cclifecycle.EnumerateFunc(f).Enumerate()
	assert.NoError(t, err)
	assert.Len(t, ccs, 2)
}

func TestLifecycleInitFailure(t *testing.T) {
	listCCs := &mocks.Enumerator{}
	listCCs.On("Enumerate").Return(nil, errors.New("failed accessing DB"))
	m, err := cclifecycle.NewMetadataManager(listCCs)
	assert.Nil(t, m)
	assert.Contains(t, err.Error(), "failed accessing DB")
}

func TestHandleChaincodeDeployGreenPath(t *testing.T) {
	recorder, restoreLogger := newLogRecorder(t)
	defer restoreLogger()

	cc1Bytes := protoutil.MarshalOrPanic(&ccprovider.ChaincodeData{
		Name:    "cc1",
		Version: "1.0",
		Id:      []byte{42},
		Policy:  []byte{1, 2, 3, 4, 5},
	})

	cc2Bytes := protoutil.MarshalOrPanic(&ccprovider.ChaincodeData{
		Name:    "cc2",
		Version: "1.0",
		Id:      []byte{42},
	})

	cc3Bytes := protoutil.MarshalOrPanic(&ccprovider.ChaincodeData{
		Name:    "cc3",
		Version: "1.0",
		Id:      []byte{42},
	})

	query := &mocks.Query{}
	query.On("GetState", "lscc", "cc1").Return(cc1Bytes, nil)
	query.On("GetState", "lscc", "cc2").Return(cc2Bytes, nil)
	query.On("GetState", "lscc", "cc3").Return(cc3Bytes, nil).Once()
	query.On("Done")
	queryCreator := &mocks.QueryCreator{}
	queryCreator.On("NewQuery").Return(query, nil)

	enum := &mocks.Enumerator{}
	enum.On("Enumerate").Return([]chaincode.InstalledChaincode{
		{
			Name:    "cc1",
			Version: "1.0",
			Hash:    []byte{42},
		},
		{
			
			Name:    "cc2",
			Version: "1.1",
			Hash:    []byte{50},
		},
		{
			
			Name:    "cc3",
			Version: "1.0",
			Hash:    []byte{50},
		},
	}, nil)

	m, err := cclifecycle.NewMetadataManager(enum)
	assert.NoError(t, err)

	lsnr := &mocks.MetadataChangeListener{}
	lsnr.On("HandleMetadataUpdate", mock.Anything, mock.Anything)
	m.AddListener(lsnr)

	sub, err := m.NewChannelSubscription("mychannel", queryCreator)
	assert.NoError(t, err)
	assert.NotNil(t, sub)

	
	assertLogged(t, recorder, "Listeners for channel mychannel invoked")
	lsnr.AssertCalled(t, "HandleMetadataUpdate", "mychannel", chaincode.MetadataSet{chaincode.Metadata{
		Name:    "cc1",
		Version: "1.0",
		Id:      []byte{42},
		Policy:  []byte{1, 2, 3, 4, 5},
	}})

	
	cc3Bytes = protoutil.MarshalOrPanic(&ccprovider.ChaincodeData{
		Name:    "cc3",
		Version: "1.0",
		Id:      []byte{50},
	})
	query.On("GetState", "lscc", "cc3").Return(cc3Bytes, nil).Once()
	sub.HandleChaincodeDeploy(&cceventmgmt.ChaincodeDefinition{Name: "cc3", Version: "1.0", Hash: []byte{50}}, nil)
	sub.ChaincodeDeployDone(true)
	
	assertLogged(t, recorder, "Listeners for channel mychannel invoked")
	assert.Len(t, lsnr.Calls, 2)
	sortedMetadata := sortedMetadataSet(lsnr.Calls[1].Arguments.Get(1).(chaincode.MetadataSet)).sort()
	assert.Equal(t, sortedMetadata, chaincode.MetadataSet{{
		Name:    "cc1",
		Version: "1.0",
		Id:      []byte{42},
		Policy:  []byte{1, 2, 3, 4, 5},
	}, {
		Name:    "cc3",
		Version: "1.0",
		Id:      []byte{50},
	}})

	
	
	cc3Bytes = protoutil.MarshalOrPanic(&ccprovider.ChaincodeData{
		Name:    "cc3",
		Version: "1.1",
		Id:      []byte{50},
	})
	query.On("GetState", "lscc", "cc3").Return(cc3Bytes, nil).Once()
	sub.HandleChaincodeDeploy(&cceventmgmt.ChaincodeDefinition{Name: "cc3", Version: "1.1", Hash: []byte{50}}, nil)
	sub.ChaincodeDeployDone(true)
	
	assertLogged(t, recorder, "Listeners for channel mychannel invoked")
	assert.Len(t, lsnr.Calls, 3)
	sortedMetadata = sortedMetadataSet(lsnr.Calls[2].Arguments.Get(1).(chaincode.MetadataSet)).sort()
	assert.Equal(t, sortedMetadata, chaincode.MetadataSet{{
		Name:    "cc1",
		Version: "1.0",
		Id:      []byte{42},
		Policy:  []byte{1, 2, 3, 4, 5},
	}, {
		Name:    "cc3",
		Version: "1.1",
		Id:      []byte{50},
	}})
}

func TestHandleChaincodeDeployFailures(t *testing.T) {
	recorder, restoreLogger := newLogRecorder(t)
	defer restoreLogger()

	cc1Bytes := protoutil.MarshalOrPanic(&ccprovider.ChaincodeData{
		Name:    "cc1",
		Version: "1.0",
		Id:      []byte{42},
		Policy:  []byte{1, 2, 3, 4, 5},
	})

	query := &mocks.Query{}
	query.On("Done")
	queryCreator := &mocks.QueryCreator{}
	enum := &mocks.Enumerator{}
	enum.On("Enumerate").Return([]chaincode.InstalledChaincode{
		{
			Name:    "cc1",
			Version: "1.0",
			Hash:    []byte{42},
		},
	}, nil)

	m, err := cclifecycle.NewMetadataManager(enum)
	assert.NoError(t, err)

	lsnr := &mocks.MetadataChangeListener{}
	lsnr.On("HandleMetadataUpdate", mock.Anything, mock.Anything)
	m.AddListener(lsnr)

	
	queryCreator.On("NewQuery").Return(nil, errors.New("failed accessing DB")).Once()
	sub, err := m.NewChannelSubscription("mychannel", queryCreator)
	assert.Nil(t, sub)
	assert.Contains(t, err.Error(), "failed accessing DB")
	lsnr.AssertNumberOfCalls(t, "HandleMetadataUpdate", 0)

	
	
	queryCreator.On("NewQuery").Return(query, nil).Once()
	queryCreator.On("NewQuery").Return(nil, errors.New("failed accessing DB")).Once()
	query.On("GetState", "lscc", "cc1").Return(cc1Bytes, nil).Once()
	sub, err = m.NewChannelSubscription("mychannel", queryCreator)
	assert.NoError(t, err)
	assert.NotNil(t, sub)
	lsnr.AssertNumberOfCalls(t, "HandleMetadataUpdate", 1)
	sub.HandleChaincodeDeploy(&cceventmgmt.ChaincodeDefinition{Name: "cc1", Version: "1.0", Hash: []byte{42}}, nil)
	sub.ChaincodeDeployDone(true)
	assertLogged(t, recorder, "Failed creating a new query for channel mychannel: failed accessing DB")
	lsnr.AssertNumberOfCalls(t, "HandleMetadataUpdate", 1)

	
	
	
	queryCreator.On("NewQuery").Return(query, nil).Once()
	query.On("GetState", "lscc", "cc1").Return(nil, errors.New("failed accessing DB")).Once()
	sub, err = m.NewChannelSubscription("mychannel", queryCreator)
	assert.NoError(t, err)
	assert.NotNil(t, sub)
	lsnr.AssertNumberOfCalls(t, "HandleMetadataUpdate", 2)
	sub.HandleChaincodeDeploy(&cceventmgmt.ChaincodeDefinition{Name: "cc1", Version: "1.0", Hash: []byte{42}}, nil)
	sub.ChaincodeDeployDone(true)
	assertLogged(t, recorder, "Query for channel mychannel for Name=cc1, Version=1.0, Hash=[]byte{0x2a} failed with error failed accessing DB")
	lsnr.AssertNumberOfCalls(t, "HandleMetadataUpdate", 2)

	
	
	
	sub, err = m.NewChannelSubscription("mychannel", queryCreator)
	lsnr.AssertNumberOfCalls(t, "HandleMetadataUpdate", 3)
	assert.NoError(t, err)
	assert.NotNil(t, sub)
	sub.HandleChaincodeDeploy(&cceventmgmt.ChaincodeDefinition{Name: "cc1", Version: "1.1", Hash: []byte{42}}, nil)
	sub.ChaincodeDeployDone(false)
	lsnr.AssertNumberOfCalls(t, "HandleMetadataUpdate", 3)
	assertLogged(t, recorder, "Chaincode deploy for updates [Name=cc1, Version=1.1, Hash=[]byte{0x2a}] failed")
}

func TestMultipleUpdates(t *testing.T) {
	recorder, restoreLogger := newLogRecorder(t)
	defer restoreLogger()

	cc1Bytes := protoutil.MarshalOrPanic(&ccprovider.ChaincodeData{
		Name:    "cc1",
		Version: "1.1",
		Id:      []byte{42},
		Policy:  []byte{1, 2, 3, 4, 5},
	})
	cc2Bytes := protoutil.MarshalOrPanic(&ccprovider.ChaincodeData{
		Name:    "cc2",
		Version: "1.0",
		Id:      []byte{50},
		Policy:  []byte{1, 2, 3, 4, 5},
	})

	query := &mocks.Query{}
	query.On("GetState", "lscc", "cc1").Return(cc1Bytes, nil)
	query.On("GetState", "lscc", "cc2").Return(cc2Bytes, nil)
	query.On("Done")
	queryCreator := &mocks.QueryCreator{}
	queryCreator.On("NewQuery").Return(query, nil)

	enum := &mocks.Enumerator{}
	enum.On("Enumerate").Return([]chaincode.InstalledChaincode{
		{
			Name:    "cc1",
			Version: "1.1",
			Hash:    []byte{42},
		},
		{
			Name:    "cc2",
			Version: "1.0",
			Hash:    []byte{50},
		},
	}, nil)

	m, err := cclifecycle.NewMetadataManager(enum)
	assert.NoError(t, err)

	var lsnrCalled sync.WaitGroup
	lsnrCalled.Add(3)
	lsnr := &mocks.MetadataChangeListener{}
	lsnr.On("HandleMetadataUpdate", mock.Anything, mock.Anything).Run(func(arguments mock.Arguments) {
		lsnrCalled.Done()
	})
	m.AddListener(lsnr)

	sub, err := m.NewChannelSubscription("mychannel", queryCreator)
	assert.NoError(t, err)

	sub.HandleChaincodeDeploy(&cceventmgmt.ChaincodeDefinition{Name: "cc1", Version: "1.1", Hash: []byte{42}}, nil)
	sub.HandleChaincodeDeploy(&cceventmgmt.ChaincodeDefinition{Name: "cc2", Version: "1.0", Hash: []byte{50}}, nil)
	sub.ChaincodeDeployDone(true)

	cc1MD := chaincode.Metadata{
		Name:    "cc1",
		Version: "1.1",
		Id:      []byte{42},
		Policy:  []byte{1, 2, 3, 4, 5},
	}
	cc2MD := chaincode.Metadata{
		Name:    "cc2",
		Version: "1.0",
		Id:      []byte{50},
		Policy:  []byte{1, 2, 3, 4, 5},
	}
	metadataSetWithBothChaincodes := chaincode.MetadataSet{cc1MD, cc2MD}

	lsnrCalled.Wait()
	
	
	expectedMetadata := sortedMetadataSet(lsnr.Calls[2].Arguments.Get(1).(chaincode.MetadataSet)).sort()
	assert.Equal(t, metadataSetWithBothChaincodes, expectedMetadata)

	
	g := NewGomegaWithT(t)
	g.Eventually(func() []string {
		return recorder.EntriesMatching("Listeners for channel mychannel invoked")
	}, time.Second*10).Should(HaveLen(3))
}

func TestMetadata(t *testing.T) {
	recorder, restoreLogger := newLogRecorder(t)
	defer restoreLogger()

	cc1Bytes := protoutil.MarshalOrPanic(&ccprovider.ChaincodeData{
		Name:    "cc1",
		Version: "1.0",
		Id:      []byte{42},
		Policy:  []byte{1, 2, 3, 4, 5},
	})

	cc2Bytes := protoutil.MarshalOrPanic(&ccprovider.ChaincodeData{
		Name:    "cc2",
		Version: "1.0",
		Id:      []byte{42},
	})

	query := &mocks.Query{}
	query.On("GetState", "lscc", "cc3").Return(cc1Bytes, nil)
	query.On("Done")
	queryCreator := &mocks.QueryCreator{}

	enum := &mocks.Enumerator{}
	enum.On("Enumerate").Return([]chaincode.InstalledChaincode{
		{
			Name:    "cc1",
			Version: "1.0",
			Hash:    []byte{42},
		},
	}, nil)

	m, err := cclifecycle.NewMetadataManager(enum)
	assert.NoError(t, err)

	
	md := m.Metadata("mychannel", "cc1", false)
	assert.Nil(t, md)
	assertLogged(t, recorder, "Requested Metadata for non-existent channel mychannel")

	
	
	query.On("GetState", "lscc", "cc1").Return(cc1Bytes, nil).Once()
	queryCreator.On("NewQuery").Return(query, nil).Once()
	sub, err := m.NewChannelSubscription("mychannel", queryCreator)
	defer sub.ChaincodeDeployDone(true)
	assert.NoError(t, err)
	assert.NotNil(t, sub)
	md = m.Metadata("mychannel", "cc1", false)
	assert.Equal(t, &chaincode.Metadata{
		Name:    "cc1",
		Version: "1.0",
		Id:      []byte{42},
		Policy:  []byte{1, 2, 3, 4, 5},
	}, md)
	assertLogged(t, recorder, "Returning metadata for channel mychannel , chaincode cc1")

	
	
	queryCreator.On("NewQuery").Return(nil, errors.New("failed obtaining query executor")).Once()
	md = m.Metadata("mychannel", "cc2", false)
	assert.Nil(t, md)
	assertLogged(t, recorder, "Failed obtaining new query for channel mychannel : failed obtaining query executor")

	
	
	queryCreator.On("NewQuery").Return(query, nil).Once()
	query.On("GetState", "lscc", "cc2").Return(nil, errors.New("GetState failed")).Once()
	md = m.Metadata("mychannel", "cc2", false)
	assert.Nil(t, md)
	assertLogged(t, recorder, "Failed querying LSCC for channel mychannel : GetState failed")

	
	
	queryCreator.On("NewQuery").Return(query, nil).Once()
	query.On("GetState", "lscc", "cc2").Return(nil, nil).Once()
	md = m.Metadata("mychannel", "cc2", false)
	assert.Nil(t, md)
	assertLogged(t, recorder, "Chaincode cc2 isn't defined in channel mychannel")

	
	
	queryCreator.On("NewQuery").Return(query, nil).Once()
	query.On("GetState", "lscc", "cc2").Return(cc2Bytes, nil).Once()
	md = m.Metadata("mychannel", "cc2", false)
	assert.Equal(t, &chaincode.Metadata{
		Name:    "cc2",
		Version: "1.0",
		Id:      []byte{42},
	}, md)

	
	
	
	queryCreator.On("NewQuery").Return(query, nil).Once()
	query.On("GetState", "lscc", "cc1").Return(cc1Bytes, nil).Once()
	query.On("GetState", "lscc", privdata.BuildCollectionKVSKey("cc1")).Return(protoutil.MarshalOrPanic(&common.CollectionConfigPackage{}), nil).Once()
	md = m.Metadata("mychannel", "cc1", true)
	assert.Equal(t, &chaincode.Metadata{
		Name:              "cc1",
		Version:           "1.0",
		Id:                []byte{42},
		Policy:            []byte{1, 2, 3, 4, 5},
		CollectionsConfig: &common.CollectionConfigPackage{},
	}, md)
	assertLogged(t, recorder, "Retrieved collection config for cc1 from cc1~collection")

	
	
	
	queryCreator.On("NewQuery").Return(query, nil).Once()
	query.On("GetState", "lscc", "cc1").Return(cc1Bytes, nil).Once()
	query.On("GetState", "lscc", privdata.BuildCollectionKVSKey("cc1")).Return(nil, errors.New("foo")).Once()
	md = m.Metadata("mychannel", "cc1", true)
	assert.Nil(t, md)
	assertLogged(t, recorder, "Failed querying lscc namespace for cc1~collection: foo")
}

func newLogRecorder(t *testing.T) (*floggingtest.Recorder, func()) {
	oldLogger := cclifecycle.Logger

	logger, recorder := floggingtest.NewTestLogger(t)
	cclifecycle.Logger = logger

	return recorder, func() { cclifecycle.Logger = oldLogger }
}

func assertLogged(t *testing.T, r *floggingtest.Recorder, msg string) {
	gt := NewGomegaWithT(t)
	gt.Eventually(r).Should(gbytes.Say(regexp.QuoteMeta(msg)))
}

type sortedMetadataSet chaincode.MetadataSet

func (mds sortedMetadataSet) Len() int {
	return len(mds)
}

func (mds sortedMetadataSet) Less(i, j int) bool {
	eI := strings.Replace(mds[i].Name, "cc", "", -1)
	eJ := strings.Replace(mds[j].Name, "cc", "", -1)
	nI, _ := strconv.ParseInt(eI, 10, 32)
	nJ, _ := strconv.ParseInt(eJ, 10, 32)
	return nI < nJ
}

func (mds sortedMetadataSet) Swap(i, j int) {
	mds[i], mds[j] = mds[j], mds[i]
}

func (mds sortedMetadataSet) sort() chaincode.MetadataSet {
	sort.Sort(mds)
	return chaincode.MetadataSet(mds)
}
