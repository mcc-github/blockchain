/*
Copyright IBM Corp. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package statecouchdb

import (
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/version"
)

type versions map[string]nsVersions
type revisions map[string]nsRevisions
type nsRevisions map[string]string
type nsVersions map[string]*version.Height





type versionsCache struct {
	vers versions
	revs revisions
}

func newVersionCache() *versionsCache {
	return &versionsCache{make(versions), make(revisions)}
}

func (c *versionsCache) getVersion(ns, key string) (*version.Height, bool) {
	ver, ok := c.vers[ns][key]
	if ok {
		return ver, true
	}
	return nil, false
}








func (c *versionsCache) setVerAndRev(ns, key string, ver *version.Height, rev string) {
	_, ok := c.vers[ns]
	if !ok {
		c.vers[ns] = make(nsVersions)
		c.revs[ns] = make(nsRevisions)
	}
	c.vers[ns][key] = ver
	c.revs[ns][key] = rev
}
