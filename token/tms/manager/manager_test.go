/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package manager_test

import (
	"github.com/mcc-github/blockchain/token/identity/mock"
	"github.com/mcc-github/blockchain/token/tms/manager"
	"github.com/mcc-github/blockchain/token/tms/plain"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
)

var _ = Describe("Manager", func() {
	var (
		mgm                             *manager.Manager
		fakeIdentityDeserializerManager *mock.DeserializerManager
	)

	BeforeEach(func() {
		fakeIdentityDeserializerManager = &mock.DeserializerManager{}
		mgm = &manager.Manager{IdentityDeserializerManager: fakeIdentityDeserializerManager}
	})

	Describe("Get a TxProcessor for a non-existent channel", func() {
		BeforeEach(func() {
			fakeIdentityDeserializerManager.DeserializerReturns(nil, errors.New("GetDeserializerReturns no-way-man"))
		})
		It("returns an error", func() {
			_, err := mgm.GetTxProcessor("boguschannel")
			Expect(err.Error()).To(Equal("failed getting identity deserialiser manager for channel 'boguschannel': GetDeserializerReturns no-way-man"))
		})
	})

	Context("When a channel exists", func() {
		var (
			fakeIdentityDeserializer *mock.Deserializer
			channel                  string
		)
		BeforeEach(func() {
			channel = "ch0"
			fakeIdentityDeserializer = &mock.Deserializer{}
			fakeIdentityDeserializerManager.DeserializerReturns(fakeIdentityDeserializer, nil)
		})

		Describe("Get a TxProcessor for an existing channel", func() {
			It("returns a Verifier that implements the TxProcessor interface", func() {
				txProcessor, err := mgm.GetTxProcessor(channel)
				Expect(err).NotTo(HaveOccurred())
				Expect(txProcessor).NotTo(BeNil())
				Expect(txProcessor).To(Equal(
					&plain.Verifier{
						IssuingValidator:    &manager.AllIssuingValidator{Deserializer: fakeIdentityDeserializer},
						TokenOwnerValidator: &manager.FabricTokenOwnerValidator{Deserializer: fakeIdentityDeserializer},
					}),
				)
			})
		})
	})
})

var _ = Describe("FabricIdentityDeserializerManager", func() {
	Describe("Get an IdentityDeserializer for a non-existent channel", func() {
		var (
			blockchainIdentityDeserializerManager *manager.FabricIdentityDeserializerManager
		)
		BeforeEach(func() {
			blockchainIdentityDeserializerManager = &manager.FabricIdentityDeserializerManager{}
		})
		It("returns an error", func() {
			_, err := blockchainIdentityDeserializerManager.Deserializer("boguschannel")
			Expect(err).To(MatchError("channel not found"))
		})
	})
})
