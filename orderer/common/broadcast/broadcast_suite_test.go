/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package broadcast_test

import (
	"testing"

	ab "github.com/mcc-github/blockchain-protos-go/orderer"
	"github.com/mcc-github/blockchain/common/metrics"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)


type abServer interface {
	ab.AtomicBroadcast_BroadcastServer
}


type metricsHistogram interface {
	metrics.Histogram
}


type metricsCounter interface {
	metrics.Counter
}


type metricsProvider interface {
	metrics.Provider
}

func TestBroadcast(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Broadcast Suite")
}
