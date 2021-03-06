/*
Copyright IBM Corp. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package lockbasedtxmgr

import (
	"github.com/mcc-github/blockchain-protos-go/common"
	"github.com/mcc-github/blockchain/core/ledger"
)



type collNameValidator struct {
	ledgerID       string
	ccInfoProvider ledger.DeployedChaincodeInfoProvider
	queryExecutor  *lockBasedQueryExecutor
	cache          collConfigCache
	noop           bool
}

func newCollNameValidator(ledgerID string, ccInfoProvider ledger.DeployedChaincodeInfoProvider, qe *lockBasedQueryExecutor, noop bool) *collNameValidator {
	return &collNameValidator{ledgerID, ccInfoProvider, qe, make(collConfigCache), noop}
}

func (v *collNameValidator) validateCollName(ns, coll string) error {
	if v.noop {
		return nil
	}
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
	return nil
}

func (v *collNameValidator) retrieveCollConfigFromStateDB(ns string) (*common.CollectionConfigPackage, error) {
	logger.Debugf("retrieveCollConfigFromStateDB() begin - ns=[%s]", ns)
	ccInfo, err := v.ccInfoProvider.ChaincodeInfo(v.ledgerID, ns, v.queryExecutor)
	if err != nil {
		return nil, err
	}
	if ccInfo == nil {
		return nil, &ledger.CollConfigNotDefinedError{Ns: ns}
	}

	confPkg := ccInfo.AllCollectionsConfigPkg()
	if confPkg == nil {
		return nil, &ledger.CollConfigNotDefinedError{Ns: ns}
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
