/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package blockcutter_test

import (
	"testing"

	"github.com/mcc-github/blockchain/common/channelconfig"
	"github.com/mcc-github/blockchain/common/metrics"
	"github.com/mcc-github/blockchain/orderer/common/blockcutter"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)


type metricsHistogram interface {
	metrics.Histogram
}


type metricsProvider interface {
	metrics.Provider
}


type ordererConfigFetcher interface {
	blockcutter.OrdererConfigFetcher
}


type ordererConfig interface {
	channelconfig.Orderer
}

func TestBlockcutter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Blockcutter Suite")
}
