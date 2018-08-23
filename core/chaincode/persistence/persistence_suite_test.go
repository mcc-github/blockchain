/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package persistence_test

import (
	"os"
	"testing"

	"github.com/mcc-github/blockchain/core/chaincode/persistence"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)


type ioReadWriter interface {
	persistence.IOReadWriter
}


type osFileInfo interface {
	os.FileInfo
}


type storePackageProvider interface {
	persistence.StorePackageProvider
}


type legacyPackageProvider interface {
	persistence.LegacyPackageProvider
}


type packageParser interface {
	persistence.PackageParser
}

func TestPersistence(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Persistence Suite")
}
