/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package broadcast_test

import (
	"testing"

	ab "github.com/mcc-github/blockchain/protos/orderer"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)


type abServer interface {
	ab.AtomicBroadcast_BroadcastServer
}

func TestBroadcast(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Broadcast Suite")
}
