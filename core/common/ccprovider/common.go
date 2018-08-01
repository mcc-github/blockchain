/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package ccprovider

import (
	"path/filepath"

	"github.com/mcc-github/blockchain/core/config"
)


func GetChaincodeInstallPathFromViper() string {
	return filepath.Join(config.GetPath("peer.fileSystemPath"), "chaincodes")
}


func LoadPackage(ccname string, ccversion string, path string) (CCPackage, error) {
	return (&CCInfoFSImpl{}).GetChaincodeFromPath(ccname, ccversion, path)
}
