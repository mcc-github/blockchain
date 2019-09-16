/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package blocksprovider_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/mcc-github/blockchain-protos-go/orderer"
	"github.com/mcc-github/blockchain/internal/pkg/identity"
)


type signer interface {
	identity.SignerSerializer
}


type abDeliverClient interface {
	orderer.AtomicBroadcast_DeliverClient
}

func TestBlocksprovider(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Blocksprovider Suite")
}
