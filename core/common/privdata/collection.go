/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package privdata

import (
	"strings"

	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/protos/common"
	pb "github.com/mcc-github/blockchain/protos/peer"
	"github.com/mcc-github/blockchain/protoutil"
)


type Collection interface {
	
	
	

	
	CollectionID() string

	
	
	

	
	
	MemberOrgs() []string
}


type CollectionAccessPolicy interface {
	
	AccessFilter() Filter

	
	
	
	RequiredPeerCount() int

	
	
	MaximumPeerCount() int

	
	
	MemberOrgs() []string

	
	
	IsMemberOnlyRead() bool

	
	
	IsMemberOnlyWrite() bool
}


type CollectionPersistenceConfigs interface {
	
	
	
	BlockToLive() uint64
}







type Filter func(protoutil.SignedData) bool













type CollectionStore interface {
	
	
	
	
	
	RetrieveCollection(common.CollectionCriteria) (Collection, error)

	
	RetrieveCollectionAccessPolicy(common.CollectionCriteria) (CollectionAccessPolicy, error)

	
	
	RetrieveCollectionConfigPackage(common.CollectionCriteria) (*common.CollectionConfigPackage, error)

	
	RetrieveCollectionPersistenceConfigs(common.CollectionCriteria) (CollectionPersistenceConfigs, error)

	
	
	
	RetrieveReadWritePermission(common.CollectionCriteria, *pb.SignedProposal, ledger.QueryExecutor) (bool, bool, error)

	CollectionFilter
}

type CollectionFilter interface {
	
	AccessFilter(channelName string, collectionPolicyConfig *common.CollectionPolicyConfig) (Filter, error)
}

const (
	

	
	
	
	
	
	collectionSeparator = "~"
	
	
	collectionSuffix = "collection"
)


func BuildCollectionKVSKey(ccname string) string {
	return ccname + collectionSeparator + collectionSuffix
}


func IsCollectionConfigKey(key string) bool {
	return strings.Contains(key, collectionSeparator)
}


func GetCCNameFromCollectionConfigKey(key string) string {
	splittedKey := strings.Split(key, collectionSeparator)
	return splittedKey[0]
}
