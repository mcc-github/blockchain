/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package container_test

import (
	"github.com/mcc-github/blockchain/core/container"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)


type vmProvider interface {
	container.VMProvider
}


type vm interface {
	container.VM
}


type vmcReq interface {
	container.VMCReq
}


type builder interface {
	container.Builder
}

func TestContainer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Container Suite")
}
