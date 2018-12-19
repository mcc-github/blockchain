/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package aclmgmt







type aclMgmtImpl struct {
	
	rescfgProvider ACLProvider
}




func (am *aclMgmtImpl) CheckACL(resName string, channelID string, idinfo interface{}) error {
	
	return am.rescfgProvider.CheckACL(resName, channelID, idinfo)
}




func NewACLProvider(rg ResourceGetter) ACLProvider {
	return &aclMgmtImpl{
		rescfgProvider: newResourceProvider(rg, NewDefaultACLProvider()),
	}
}
