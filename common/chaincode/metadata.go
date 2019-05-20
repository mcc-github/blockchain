/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chaincode

import (
	"sync"

	persistence "github.com/mcc-github/blockchain/core/chaincode/persistence/intf"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/gossip"
)


type InstalledChaincode struct {
	PackageID persistence.PackageID
	Hash      []byte
	Label     string

	
	
	
	Name    string
	Version string
}


type Metadata struct {
	Name              string
	Version           string
	Policy            []byte
	Id                []byte
	CollectionsConfig *common.CollectionConfigPackage
	
	
	
	
	
	Approved  bool
	Installed bool
}


type MetadataSet []Metadata


func (ccs MetadataSet) AsChaincodes() []*gossip.Chaincode {
	var res []*gossip.Chaincode
	for _, cc := range ccs {
		res = append(res, &gossip.Chaincode{
			Name:    cc.Name,
			Version: cc.Version,
		})
	}
	return res
}


type MetadataMapping struct {
	sync.RWMutex
	mdByName map[string]Metadata
}


func NewMetadataMapping() *MetadataMapping {
	return &MetadataMapping{
		mdByName: make(map[string]Metadata),
	}
}


func (m *MetadataMapping) Lookup(cc string) (Metadata, bool) {
	m.RLock()
	defer m.RUnlock()
	md, exists := m.mdByName[cc]
	return md, exists
}


func (m *MetadataMapping) Update(ccMd Metadata) {
	m.Lock()
	defer m.Unlock()
	m.mdByName[ccMd.Name] = ccMd
}


func (m *MetadataMapping) Aggregate() MetadataSet {
	m.RLock()
	defer m.RUnlock()
	var set MetadataSet
	for _, md := range m.mdByName {
		set = append(set, md)
	}
	return set
}
