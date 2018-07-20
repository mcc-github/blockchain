


/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package scc


func CreatePluginSysCCs(p *Provider) []SelfDescribingSysCC {
	var sdscs []SelfDescribingSysCC
	for _, pscc := range loadSysCCs(p) {
		sdscs = append(sdscs, &SysCCWrapper{SCC: pscc})
	}
	return sdscs
}
