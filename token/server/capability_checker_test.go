/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package server_test

import (
	"github.com/mcc-github/blockchain/common/channelconfig"
	"github.com/mcc-github/blockchain/core/peer"
	"github.com/mcc-github/blockchain/token/server"
	"github.com/mcc-github/blockchain/token/server/mock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)









type peerOperations interface {
	peer.Operations
}

type channelConfig interface {
	channelconfig.Resources
}

type applicationConfig interface {
	channelconfig.Application
}

type applicationCapabilities interface {
	channelconfig.ApplicationCapabilities
}

var _ = Describe("CapabilityChecker", func() {
	var (
		channelId = "mychannel"

		fakeAppCapabilities *mock.ApplicationCapabilities
		fakeAppConfig       *mock.ApplicationConfig
		fakeChannelConfig   *mock.ChannelConfig
		fakePeerOperations  *mock.PeerOperations

		capabilityChecker *server.TokenCapabilityChecker
	)

	BeforeEach(func() {
		fakeAppCapabilities = &mock.ApplicationCapabilities{}
		fakeAppCapabilities.FabTokenReturns(true)

		fakeAppConfig = &mock.ApplicationConfig{}
		fakeAppConfig.CapabilitiesReturns(fakeAppCapabilities)

		fakeChannelConfig = &mock.ChannelConfig{}
		fakeChannelConfig.ApplicationConfigReturns(fakeAppConfig, true)

		fakePeerOperations = &mock.PeerOperations{}
		fakePeerOperations.GetChannelConfigReturns(fakeChannelConfig)

		capabilityChecker = &server.TokenCapabilityChecker{PeerOps: fakePeerOperations}
	})

	It("returns FabToken true when application capabilities returns true", func() {
		fakeAppCapabilities.FabTokenReturns(true)
		result, err := capabilityChecker.FabToken(channelId)
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(true))
	})

	It("returns FabToken false when application capabilities returns false", func() {
		fakeAppCapabilities.FabTokenReturns(false)
		result, err := capabilityChecker.FabToken(channelId)
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(false))
	})

	Context("when channel config is not found", func() {
		BeforeEach(func() {
			fakePeerOperations.GetChannelConfigReturns(nil)
		})

		It("returns the error", func() {
			_, err := capabilityChecker.FabToken(channelId)
			Expect(err).To(MatchError("no channel config found for channel " + channelId))
		})
	})

	Context("when application config is not found", func() {
		BeforeEach(func() {
			fakeChannelConfig.ApplicationConfigReturns(nil, false)
		})

		It("returns the error", func() {
			_, err := capabilityChecker.FabToken(channelId)
			Expect(err).To(MatchError("no application config found for channel " + channelId))
		})
	})
})
