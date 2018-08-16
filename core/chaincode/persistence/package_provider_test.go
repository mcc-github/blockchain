/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package persistence_test

import (
	"github.com/mcc-github/blockchain/common/chaincode"
	"github.com/mcc-github/blockchain/core/chaincode/persistence"
	"github.com/mcc-github/blockchain/core/chaincode/persistence/mock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
)

var _ = Describe("PackageProvider", func() {
	var _ = Describe("GetChaincodeCodePackage", func() {
		var (
			mockSPP         *mock.StorePackageProvider
			mockLPP         *mock.LegacyPackageProvider
			packageProvider *persistence.PackageProvider
		)

		BeforeEach(func() {
			mockSPP = &mock.StorePackageProvider{}
			mockSPP.RetrieveHashReturns([]byte("testcchash"), nil)
			mockSPP.LoadReturns([]byte("storeCode"), "testcc", "1.0", nil)

			mockLPP = &mock.LegacyPackageProvider{}
			mockLPP.GetChaincodeCodePackageReturns([]byte("legacyCode"), nil)

			packageProvider = &persistence.PackageProvider{
				Store:    mockSPP,
				LegacyPP: mockLPP,
			}
		})

		It("gets the code package successfully", func() {
			pkgBytes, err := packageProvider.GetChaincodeCodePackage("testcc", "1.0")
			Expect(err).NotTo(HaveOccurred())
			Expect(pkgBytes).To(Equal([]byte("storeCode")))
		})

		Context("when the code package is not available in the store package provider", func() {
			BeforeEach(func() {
				mockSPP.RetrieveHashReturns(nil, &persistence.CodePackageNotFoundErr{})
			})

			It("gets the code package successfully from the legacy package provider", func() {
				pkgBytes, err := packageProvider.GetChaincodeCodePackage("testcc", "1.0")
				Expect(err).NotTo(HaveOccurred())
				Expect(pkgBytes).To(Equal([]byte("legacyCode")))
			})
		})

		Context("when retrieving the hash from the store package provider fails", func() {
			BeforeEach(func() {
				mockSPP.RetrieveHashReturns(nil, errors.New("chai"))
			})

			It("returns an error", func() {
				pkgBytes, err := packageProvider.GetChaincodeCodePackage("testcc", "1.0")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("error retrieving hash: chai"))
				Expect(pkgBytes).To(BeNil())
			})
		})

		Context("when the code package fails to load from the store package provider", func() {
			BeforeEach(func() {
				mockSPP.LoadReturns(nil, "", "", errors.New("mocha"))
			})

			It("returns an error", func() {
				pkgBytes, err := packageProvider.GetChaincodeCodePackage("testcc", "1.0")
				Expect(err).To(HaveOccurred())
				Expect(pkgBytes).To(BeNil())
			})
		})

		Context("when the code package is not available in either package provider", func() {
			BeforeEach(func() {
				mockSPP.RetrieveHashReturns(nil, &persistence.CodePackageNotFoundErr{})
				mockLPP.GetChaincodeCodePackageReturns(nil, errors.New("latte"))
			})

			It("returns an error", func() {
				pkgBytes, err := packageProvider.GetChaincodeCodePackage("testcc", "1.0")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("code package not found for chaincode with name 'testcc', version '1.0'"))
				Expect(len(pkgBytes)).To(Equal(0))
			})
		})
	})

	var _ = Describe("ListInstalledChaincodes", func() {
		var (
			mockSPP         *mock.StorePackageProvider
			mockLPP         *mock.LegacyPackageProvider
			packageProvider *persistence.PackageProvider
		)

		BeforeEach(func() {
			mockSPP = &mock.StorePackageProvider{}
			installedChaincodes := []chaincode.InstalledChaincode{
				{
					Name:    "test1",
					Version: "1.0",
					Id:      []byte("hash1"),
				},
				{
					Name:    "cc1",
					Version: "2.0",
					Id:      []byte("hash2"),
				},
			}
			mockSPP.ListInstalledChaincodesReturns(installedChaincodes, nil)

			mockLPP = &mock.LegacyPackageProvider{}
			installedChaincodesLegacy := []chaincode.InstalledChaincode{
				{
					Name:    "testLegacy",
					Version: "1.0",
					Id:      []byte("hashLegacy"),
				},
			}
			mockLPP.ListInstalledChaincodesReturns(installedChaincodesLegacy, nil)

			packageProvider = &persistence.PackageProvider{
				Store:    mockSPP,
				LegacyPP: mockLPP,
			}
		})

		It("lists the installed chaincodes successfully", func() {
			installedChaincodes, err := packageProvider.ListInstalledChaincodes()
			Expect(err).NotTo(HaveOccurred())
			Expect(len(installedChaincodes)).To(Equal(3))
		})

		Context("when listing the installed chaincodes from the persistence store fails", func() {
			BeforeEach(func() {
				mockSPP.ListInstalledChaincodesReturns(nil, errors.New("football"))
			})

			It("falls back to listing the chaincodes from the legacy package provider", func() {
				installedChaincodes, err := packageProvider.ListInstalledChaincodes()
				Expect(err).NotTo(HaveOccurred())
				Expect(len(installedChaincodes)).To(Equal(1))
			})
		})

		Context("when listing the installed chaincodes from the legacy package provider fails", func() {
			BeforeEach(func() {
				mockLPP.ListInstalledChaincodesReturns(nil, errors.New("football"))
			})

			It("lists the chaincodes from only the persistence store package provider ", func() {
				installedChaincodes, err := packageProvider.ListInstalledChaincodes()
				Expect(err).NotTo(HaveOccurred())
				Expect(len(installedChaincodes)).To(Equal(2))
			})
		})
	})
})