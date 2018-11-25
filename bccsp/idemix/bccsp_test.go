/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package idemix_test

import (
	"github.com/mcc-github/blockchain/bccsp"
	"github.com/mcc-github/blockchain/bccsp/idemix"
	"github.com/mcc-github/blockchain/bccsp/sw"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Idemix Bridge", func() {

	Describe("setting up the environment with one issuer and one user", func() {
		var (
			CSP             bccsp.BCCSP
			IssuerKey       bccsp.Key
			IssuerPublicKey bccsp.Key
			AttributeNames  []string

			UserKey      bccsp.Key
			NymKey       bccsp.Key
			NymPublicKey bccsp.Key

			credRequest []byte

			credential []byte

			RevocationKey       bccsp.Key
			RevocationPublicKey bccsp.Key
			cri                 []byte
		)

		BeforeEach(func() {
			var err error
			CSP, err = idemix.New(sw.NewDummyKeyStore())
			Expect(err).NotTo(HaveOccurred())

			
			AttributeNames = []string{"Attr1", "Attr2", "Attr3", "Attr4", "Attr5"}
			IssuerKey, err = CSP.KeyGen(&bccsp.IdemixIssuerKeyGenOpts{Temporary: true, AttributeNames: AttributeNames})
			Expect(err).NotTo(HaveOccurred())
			IssuerPublicKey, err = IssuerKey.PublicKey()
			Expect(err).NotTo(HaveOccurred())

			
			UserKey, err = CSP.KeyGen(&bccsp.IdemixUserSecretKeyGenOpts{Temporary: true})
			Expect(err).NotTo(HaveOccurred())

			
			NymKey, err = CSP.KeyDeriv(UserKey, &bccsp.IdemixNymKeyDerivationOpts{Temporary: true, IssuerPK: IssuerPublicKey})
			Expect(err).NotTo(HaveOccurred())
			NymPublicKey, err = NymKey.PublicKey()
			Expect(err).NotTo(HaveOccurred())

			
			credRequest, err = CSP.Sign(
				UserKey,
				bccsp.IdemixEmptyDigest(),
				&bccsp.IdemixCredentialRequestSignerOpts{IssuerPK: IssuerPublicKey},
			)
			Expect(err).NotTo(HaveOccurred())

			
			credential, err = CSP.Sign(
				IssuerKey,
				credRequest,
				&bccsp.IdemixCredentialSignerOpts{
					Attributes: []bccsp.IdemixAttribute{
						{Type: bccsp.IdemixBytesAttribute, Value: []byte{0}},
						{Type: bccsp.IdemixBytesAttribute, Value: []byte{0, 1}},
						{Type: bccsp.IdemixIntAttribute, Value: 1},
						{Type: bccsp.IdemixBytesAttribute, Value: []byte{0, 1, 2}},
						{Type: bccsp.IdemixBytesAttribute, Value: []byte{0, 1, 2, 3}},
					},
				},
			)
			Expect(err).NotTo(HaveOccurred())

			
			RevocationKey, err = CSP.KeyGen(&bccsp.IdemixRevocationKeyGenOpts{Temporary: true})
			Expect(err).NotTo(HaveOccurred())
			RevocationPublicKey, err = RevocationKey.PublicKey()
			Expect(err).NotTo(HaveOccurred())

			
			cri, err = CSP.Sign(
				RevocationKey,
				bccsp.IdemixEmptyDigest(),
				&bccsp.IdemixCRISignerOpts{},
			)
			Expect(err).NotTo(HaveOccurred())

		})

		It("the environment is properly set", func() {
			
			valid, err := CSP.Verify(
				IssuerPublicKey,
				credRequest,
				bccsp.IdemixEmptyDigest(),
				&bccsp.IdemixCredentialRequestSignerOpts{},
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(valid).To(BeTrue())

			
			valid, err = CSP.Verify(
				UserKey,
				credential,
				bccsp.IdemixEmptyDigest(),
				&bccsp.IdemixCredentialSignerOpts{
					IssuerPK: IssuerPublicKey,
					Attributes: []bccsp.IdemixAttribute{
						{Type: bccsp.IdemixBytesAttribute, Value: []byte{0}},
						{Type: bccsp.IdemixBytesAttribute, Value: []byte{0, 1}},
						{Type: bccsp.IdemixIntAttribute, Value: 1},
						{Type: bccsp.IdemixBytesAttribute, Value: []byte{0, 1, 2}},
						{Type: bccsp.IdemixBytesAttribute, Value: []byte{0, 1, 2, 3}},
					},
				},
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(valid).To(BeTrue())

			
			valid, err = CSP.Verify(
				RevocationPublicKey,
				cri,
				bccsp.IdemixEmptyDigest(),
				&bccsp.IdemixCRISignerOpts{},
			)
			Expect(err).NotTo(HaveOccurred())
		})

		Describe("producing an idemix signature with no disclosed attribute", func() {
			var (
				digest    []byte
				signature []byte
			)

			BeforeEach(func() {
				var err error

				digest = []byte("a digest")

				signature, err = CSP.Sign(
					UserKey,
					digest,
					&bccsp.IdemixSignerOpts{
						Credential: credential,
						Nym:        NymKey,
						IssuerPK:   IssuerPublicKey,
						Attributes: []bccsp.IdemixAttribute{
							{Type: bccsp.IdemixHiddenAttribute},
							{Type: bccsp.IdemixHiddenAttribute},
							{Type: bccsp.IdemixHiddenAttribute},
							{Type: bccsp.IdemixHiddenAttribute},
							{Type: bccsp.IdemixHiddenAttribute},
						},
						RhIndex: 4,
						Epoch:   0,
						CRI:     cri,
					},
				)
				Expect(err).NotTo(HaveOccurred())
			})

			It("the signature is valid", func() {
				valid, err := CSP.Verify(
					IssuerPublicKey,
					signature,
					digest,
					&bccsp.IdemixSignerOpts{
						RevocationPublicKey: RevocationPublicKey,
						Attributes: []bccsp.IdemixAttribute{
							{Type: bccsp.IdemixHiddenAttribute},
							{Type: bccsp.IdemixHiddenAttribute},
							{Type: bccsp.IdemixHiddenAttribute},
							{Type: bccsp.IdemixHiddenAttribute},
							{Type: bccsp.IdemixHiddenAttribute},
						},
						RhIndex: 4,
						Epoch:   0,
					},
				)
				Expect(err).NotTo(HaveOccurred())
				Expect(valid).To(BeTrue())
			})

		})

		Describe("producing an idemix signature with disclosed attributes", func() {
			var (
				digest    []byte
				signature []byte
			)

			BeforeEach(func() {
				var err error

				digest = []byte("a digest")

				signature, err = CSP.Sign(
					UserKey,
					digest,
					&bccsp.IdemixSignerOpts{
						Credential: credential,
						Nym:        NymKey,
						IssuerPK:   IssuerPublicKey,
						Attributes: []bccsp.IdemixAttribute{
							{Type: bccsp.IdemixBytesAttribute},
							{Type: bccsp.IdemixHiddenAttribute},
							{Type: bccsp.IdemixIntAttribute},
							{Type: bccsp.IdemixHiddenAttribute},
							{Type: bccsp.IdemixHiddenAttribute},
						},
						RhIndex: 4,
						Epoch:   0,
						CRI:     cri,
					},
				)
				Expect(err).NotTo(HaveOccurred())
			})

			It("the signature is valid", func() {
				valid, err := CSP.Verify(
					IssuerPublicKey,
					signature,
					digest,
					&bccsp.IdemixSignerOpts{
						RevocationPublicKey: RevocationPublicKey,
						Attributes: []bccsp.IdemixAttribute{
							{Type: bccsp.IdemixBytesAttribute, Value: []byte{0}},
							{Type: bccsp.IdemixHiddenAttribute},
							{Type: bccsp.IdemixIntAttribute, Value: 1},
							{Type: bccsp.IdemixHiddenAttribute},
							{Type: bccsp.IdemixHiddenAttribute},
						},
						RhIndex: 4,
						Epoch:   0,
					},
				)
				Expect(err).NotTo(HaveOccurred())
				Expect(valid).To(BeTrue())
			})

		})

		Describe("producing an idemix nym signature", func() {
			var (
				digest    []byte
				signature []byte
			)

			BeforeEach(func() {
				var err error

				digest = []byte("a digest")

				signature, err = CSP.Sign(
					UserKey,
					digest,
					&bccsp.IdemixNymSignerOpts{
						Nym:      NymKey,
						IssuerPK: IssuerPublicKey,
					},
				)
				Expect(err).NotTo(HaveOccurred())
			})

			It("the signature is valid", func() {
				valid, err := CSP.Verify(
					NymPublicKey,
					signature,
					digest,
					&bccsp.IdemixNymSignerOpts{
						IssuerPK: IssuerPublicKey,
					},
				)
				Expect(err).NotTo(HaveOccurred())
				Expect(valid).To(BeTrue())
			})

		})
	})
})
