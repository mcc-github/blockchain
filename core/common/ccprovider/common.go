/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package ccprovider


func LoadPackage(ccname string, ccversion string, path string) (CCPackage, error) {
	return (&CCInfoFSImpl{}).GetChaincodeFromPath(ccname, ccversion, path)
}
