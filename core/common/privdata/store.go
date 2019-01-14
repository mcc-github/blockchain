/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package privdata

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/msp"
	"github.com/mcc-github/blockchain/protos/common"
	pb "github.com/mcc-github/blockchain/protos/peer"
	"github.com/mcc-github/blockchain/protos/utils"
	"github.com/pkg/errors"
)


type Support interface {
	
	GetQueryExecutorForLedger(cid string) (ledger.QueryExecutor, error)

	
	
	GetIdentityDeserializer(chainID string) msp.IdentityDeserializer

	
	GetCollectionInfoProvider() ledger.DeployedChaincodeInfoProvider
}


type State interface {
	
	GetState(namespace string, key string) ([]byte, error)
}

type NoSuchCollectionError common.CollectionCriteria

func (f NoSuchCollectionError) Error() string {
	return fmt.Sprintf("collection %s/%s/%s could not be found", f.Channel, f.Namespace, f.Collection)
}

type simpleCollectionStore struct {
	s Support
}





func NewSimpleCollectionStore(s Support) *simpleCollectionStore {
	return &simpleCollectionStore{s}
}

func (c *simpleCollectionStore) retrieveCollectionConfigPackage(cc common.CollectionCriteria, qe ledger.QueryExecutor) (*common.CollectionConfigPackage, error) {
	var err error
	if qe == nil {
		qe, err = c.s.GetQueryExecutorForLedger(cc.Channel)
		if err != nil {
			return nil, errors.WithMessage(err, fmt.Sprintf("could not retrieve query executor for collection criteria %#v", cc))
		}
		defer qe.Done()
	}
	ccInfo, err := c.s.GetCollectionInfoProvider().ChaincodeInfo(cc.Channel, cc.Namespace, qe)
	if err != nil {
		return nil, err
	}
	if ccInfo == nil {
		return nil, errors.Errorf("Chaincode [%s] does not exist", cc.Namespace)
	}
	return ccInfo.CollectionConfigPkg, nil
}


func RetrieveCollectionConfigPackageFromState(cc common.CollectionCriteria, state State) (*common.CollectionConfigPackage, error) {

	cb, err := state.GetState("lscc", BuildCollectionKVSKey(cc.Namespace))
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("error while retrieving collection for collection criteria %#v", cc))
	}
	if cb == nil {
		return nil, NoSuchCollectionError(cc)
	}
	conf, err := ParseCollectionConfig(cb)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid configuration for collection criteria %#v", cc)
	}
	return conf, nil
}


func ParseCollectionConfig(colBytes []byte) (*common.CollectionConfigPackage, error) {
	collections := &common.CollectionConfigPackage{}
	err := proto.Unmarshal(colBytes, collections)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return collections, nil
}

func (c *simpleCollectionStore) retrieveCollectionConfig(cc common.CollectionCriteria, qe ledger.QueryExecutor) (*common.StaticCollectionConfig, error) {
	var err error
	if qe == nil {
		qe, err = c.s.GetQueryExecutorForLedger(cc.Channel)
		if err != nil {
			return nil, errors.WithMessage(err, fmt.Sprintf("could not retrieve query executor for collection criteria %#v", cc))
		}
		defer qe.Done()
	}
	collConfig, err := c.s.GetCollectionInfoProvider().CollectionInfo(cc.Channel, cc.Namespace, cc.Collection, qe)
	if err != nil {
		return nil, err
	}
	if collConfig == nil {
		return nil, NoSuchCollectionError(cc)
	}
	return collConfig, nil
}

func (c *simpleCollectionStore) retrieveSimpleCollection(cc common.CollectionCriteria, qe ledger.QueryExecutor) (*SimpleCollection, error) {
	staticCollectionConfig, err := c.retrieveCollectionConfig(cc, qe)
	if err != nil {
		return nil, err
	}
	sc := &SimpleCollection{}
	err = sc.Setup(staticCollectionConfig, c.s.GetIdentityDeserializer(cc.Channel))
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("error setting up collection for collection criteria %#v", cc))
	}
	return sc, nil
}

func (c *simpleCollectionStore) AccessFilter(channelName string, collectionPolicyConfig *common.CollectionPolicyConfig) (Filter, error) {
	sc := &SimpleCollection{}
	err := sc.setupAccessPolicy(collectionPolicyConfig, c.s.GetIdentityDeserializer(channelName))
	if err != nil {
		return nil, err
	}
	return sc.AccessFilter(), nil
}

func (c *simpleCollectionStore) RetrieveCollection(cc common.CollectionCriteria) (Collection, error) {
	return c.retrieveSimpleCollection(cc, nil)
}

func (c *simpleCollectionStore) RetrieveCollectionAccessPolicy(cc common.CollectionCriteria) (CollectionAccessPolicy, error) {
	return c.retrieveSimpleCollection(cc, nil)
}

func (c *simpleCollectionStore) RetrieveCollectionConfigPackage(cc common.CollectionCriteria) (*common.CollectionConfigPackage, error) {
	return c.retrieveCollectionConfigPackage(cc, nil)
}


func (c *simpleCollectionStore) RetrieveCollectionPersistenceConfigs(cc common.CollectionCriteria) (CollectionPersistenceConfigs, error) {
	staticCollectionConfig, err := c.retrieveCollectionConfig(cc, nil)
	if err != nil {
		return nil, err
	}
	return &SimpleCollectionPersistenceConfigs{staticCollectionConfig.BlockToLive}, nil
}

func (c *simpleCollectionStore) HasReadAccess(cc common.CollectionCriteria, signedProposal *pb.SignedProposal, qe ledger.QueryExecutor) (bool, error) {
	accessPolicy, err := c.retrieveSimpleCollection(cc, qe)
	if err != nil {
		return false, err
	}

	if !accessPolicy.IsMemberOnlyRead() {
		return true, nil
	}

	signedData, err := getSignedData(signedProposal)
	if err != nil {
		return false, err
	}

	hasReadAccess := accessPolicy.AccessFilter()
	return hasReadAccess(signedData), nil
}

func getSignedData(signedProposal *pb.SignedProposal) (common.SignedData, error) {
	proposal, err := utils.GetProposal(signedProposal.ProposalBytes)
	if err != nil {
		return common.SignedData{}, err
	}

	creator, _, err := utils.GetChaincodeProposalContext(proposal)
	if err != nil {
		return common.SignedData{}, err
	}

	return common.SignedData{
		Data:      signedProposal.ProposalBytes,
		Identity:  creator,
		Signature: signedProposal.Signature,
	}, nil
}
