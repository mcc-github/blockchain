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
	"github.com/pkg/errors"
)


type Support interface {
	
	GetQueryExecutorForLedger(cid string) (ledger.QueryExecutor, error)

	
	
	GetIdentityDeserializer(chainID string) msp.IdentityDeserializer
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





func NewSimpleCollectionStore(s Support) CollectionStore {
	return &simpleCollectionStore{s}
}

func (c *simpleCollectionStore) retrieveCollectionConfigPackage(cc common.CollectionCriteria) (*common.CollectionConfigPackage, error) {
	qe, err := c.s.GetQueryExecutorForLedger(cc.Channel)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("could not retrieve query executor for collection criteria %#v", cc))
	}
	defer qe.Done()
	return RetrieveCollectionConfigPackageFromState(cc, qe)
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

func (c *simpleCollectionStore) retrieveCollectionConfig(cc common.CollectionCriteria) (*common.StaticCollectionConfig, error) {
	collections, err := c.retrieveCollectionConfigPackage(cc)
	if err != nil {
		return nil, err
	}
	if collections == nil {
		return nil, nil
	}
	for _, cconf := range collections.Config {
		switch cconf := cconf.Payload.(type) {
		case *common.CollectionConfig_StaticCollectionConfig:
			if cconf.StaticCollectionConfig.Name == cc.Collection {
				return cconf.StaticCollectionConfig, nil
			}
		default:
			return nil, errors.New("unexpected collection type")
		}
	}
	return nil, NoSuchCollectionError(cc)
}

func (c *simpleCollectionStore) retrieveSimpleCollection(cc common.CollectionCriteria) (*SimpleCollection, error) {
	staticCollectionConfig, err := c.retrieveCollectionConfig(cc)
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
	return c.retrieveSimpleCollection(cc)
}

func (c *simpleCollectionStore) RetrieveCollectionAccessPolicy(cc common.CollectionCriteria) (CollectionAccessPolicy, error) {
	return c.retrieveSimpleCollection(cc)
}

func (c *simpleCollectionStore) RetrieveCollectionConfigPackage(cc common.CollectionCriteria) (*common.CollectionConfigPackage, error) {
	return c.retrieveCollectionConfigPackage(cc)
}


func (c *simpleCollectionStore) RetrieveCollectionPersistenceConfigs(cc common.CollectionCriteria) (CollectionPersistenceConfigs, error) {
	staticCollectionConfig, err := c.retrieveCollectionConfig(cc)
	if err != nil {
		return nil, err
	}
	return &SimpleCollectionPersistenceConfigs{staticCollectionConfig.BlockToLive}, nil
}
