/*
Copyright IBM Corp. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package lockbasedtxmgr

import (
	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/core/common/privdata"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/protos/common"
)

const (
	lsccNamespace = "lscc"
)



type collNameValidator struct {
	queryHelper *queryHelper
	cache       collConfigCache
}

func newCollNameValidator(queryHelper *queryHelper) *collNameValidator {
	return &collNameValidator{queryHelper, make(collConfigCache)}
}

func (v *collNameValidator) validateCollName(ns, coll string) error {
	logger.Debugf("validateCollName() begin - ns=[%s], coll=[%s]", ns, coll)
	if !v.cache.isPopulatedFor(ns) {
		conf, err := v.retrieveCollConfigFromStateDB(ns)
		if err != nil {
			return err
		}
		v.cache.populate(ns, conf)
	}
	if !v.cache.containsCollName(ns, coll) {
		return &ledger.InvalidCollNameError{
			Ns:   ns,
			Coll: coll,
		}
	}
	logger.Debugf("validateCollName() validated successfully - ns=[%s], coll=[%s]", ns, coll)
	return nil
}

func (v *collNameValidator) retrieveCollConfigFromStateDB(ns string) (*common.CollectionConfigPackage, error) {
	logger.Debugf("retrieveCollConfigFromStateDB() begin - ns=[%s]", ns)
	configPkgBytes, _, err := v.queryHelper.getState(lsccNamespace, constructCollectionConfigKey(ns))
	if err != nil {
		return nil, err
	}
	if configPkgBytes == nil {
		return nil, &ledger.CollConfigNotDefinedError{Ns: ns}
	}
	confPkg := &common.CollectionConfigPackage{}
	if err := proto.Unmarshal(configPkgBytes, confPkg); err != nil {
		return nil, err
	}
	logger.Debugf("retrieveCollConfigFromStateDB() successfully retrieved - ns=[%s], confPkg=[%s]", ns, confPkg)
	return confPkg, nil
}

type collConfigCache map[collConfigkey]bool

type collConfigkey struct {
	ns, coll string
}

func (c collConfigCache) populate(ns string, pkg *common.CollectionConfigPackage) {
	
	
	c[collConfigkey{ns, ""}] = true
	for _, config := range pkg.Config {
		sConfig := config.GetStaticCollectionConfig()
		if sConfig == nil {
			continue
		}
		c[collConfigkey{ns, sConfig.Name}] = true
	}
}

func (c collConfigCache) isPopulatedFor(ns string) bool {
	return c[collConfigkey{ns, ""}]
}

func (c collConfigCache) containsCollName(ns, coll string) bool {
	return c[collConfigkey{ns, coll}]
}

func constructCollectionConfigKey(chaincodeName string) string {
	return privdata.BuildCollectionKVSKey(chaincodeName)
}
