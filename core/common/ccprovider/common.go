/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package ccprovider


func LoadPackage(ccNameVersion string, path string) (CCPackage, error) {
	return (&CCInfoFSImpl{}).GetChaincodeFromPath(ccNameVersion, path)
}
