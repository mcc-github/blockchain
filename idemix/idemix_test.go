/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package idemix

import (
	"bytes"
	"testing"

	"github.com/mcc-github/blockchain-amcl/amcl/FP256BN"
	"github.com/stretchr/testify/assert"
)

func TestIdemix(t *testing.T) {
	
	
	rng, err := GetRand()
	assert.NoError(t, err)
	wbbsk, wbbpk := WBBKeyGen(rng)

	
	testmsg := RandModOrder(rng)

	
	wbbsig := WBBSign(wbbsk, testmsg)

	
	err = WBBVerify(wbbpk, wbbsig, testmsg)
	assert.NoError(t, err)

	
	AttributeNames := []string{"Attr1", "Attr2", "Attr3", "Attr4", "Attr5"}
	attrs := make([]*FP256BN.BIG, len(AttributeNames))
	for i := range AttributeNames {
		attrs[i] = FP256BN.NewBIGint(i)
	}

	
	if err != nil {
		t.Fatalf("Error getting rng: \"%s\"", err)
		return
	}
	
	key, err := NewIssuerKey(AttributeNames, rng)
	if err != nil {
		t.Fatalf("Issuer key generation should have succeeded but gave error \"%s\"", err)
		return
	}

	
	err = key.GetIpk().Check()
	if err != nil {
		t.Fatalf("Issuer public key should be valid")
		return
	}

	
	proofC := key.Ipk.GetProofC()
	key.Ipk.ProofC = BigToBytes(RandModOrder(rng))
	assert.Error(t, key.Ipk.Check(), "public key with broken zero-knowledge proof should be invalid")

	
	hAttrs := key.Ipk.GetHAttrs()
	key.Ipk.HAttrs = key.Ipk.HAttrs[:0]
	assert.Error(t, key.Ipk.Check(), "public key with incorrect number of HAttrs should be invalid")
	key.Ipk.HAttrs = hAttrs

	
	key.Ipk.ProofC = proofC
	h := key.Ipk.GetHash()
	assert.NoError(t, key.Ipk.Check(), "restored public key should be valid")
	assert.Zero(t, bytes.Compare(h, key.Ipk.GetHash()), "IPK hash changed on ipk Check")

	
	_, err = NewIssuerKey([]string{"Attr1", "Attr2", "Attr1"}, rng)
	assert.Error(t, err, "issuer key generation should fail with duplicate attribute names")

	
	sk := RandModOrder(rng)
	ni := RandModOrder(rng)
	m := NewCredRequest(sk, ni, key.Ipk, rng)

	cred, err := NewCredential(key, m, attrs, rng)
	assert.NoError(t, err, "Failed to issue a credential: \"%s\"", err)

	assert.NoError(t, cred.Ver(sk, key.Ipk), "credential should be valid")

	
	_, err = NewCredential(key, m, []*FP256BN.BIG{}, rng)
	assert.Error(t, err, "issuing a credential with the incorrect amount of attributes should fail")

	
	proofC = m.GetProofC()
	m.ProofC = BigToBytes(RandModOrder(rng))
	assert.Error(t, m.Check(key.Ipk), "CredRequest with broken ZK proof should not be valid")

	
	_, err = NewCredential(key, m, attrs, rng)
	assert.Error(t, err, "creating a credential from an invalid CredRequest should fail")
	m.ProofC = proofC

	
	attrsBackup := cred.GetAttrs()
	cred.Attrs = [][]byte{nil, nil, nil, nil, nil}
	assert.Error(t, cred.Ver(sk, key.Ipk), "credential with nil attribute should be invalid")
	cred.Attrs = attrsBackup

	
	revocationKey, err := GenerateLongTermRevocationKey()
	assert.NoError(t, err)

	
	epoch := 0
	cri, err := CreateCRI(revocationKey, []*FP256BN.BIG{}, epoch, ALG_NO_REVOCATION, rng)
	assert.NoError(t, err)
	err = VerifyEpochPK(&revocationKey.PublicKey, cri.EpochPk, cri.EpochPkSig, int(cri.Epoch), RevocationAlgorithm(cri.RevocationAlg))
	assert.NoError(t, err)

	
	err = VerifyEpochPK(&revocationKey.PublicKey, cri.EpochPk, cri.EpochPkSig, int(cri.Epoch)+1, RevocationAlgorithm(cri.RevocationAlg))
	assert.Error(t, err)

	
	_, err = CreateCRI(nil, []*FP256BN.BIG{}, epoch, ALG_NO_REVOCATION, rng)
	assert.Error(t, err)
	_, err = CreateCRI(revocationKey, []*FP256BN.BIG{}, epoch, ALG_NO_REVOCATION, nil)
	assert.Error(t, err)

	
	Nym, RandNym := MakeNym(sk, key.Ipk, rng)

	disclosure := []byte{0, 0, 0, 0, 0}
	msg := []byte{1, 2, 3, 4, 5}
	rhindex := 4
	sig, err := NewSignature(cred, sk, Nym, RandNym, key.Ipk, disclosure, msg, rhindex, cri, rng)
	assert.NoError(t, err)

	err = sig.Ver(disclosure, key.Ipk, msg, nil, 0, &revocationKey.PublicKey, epoch)
	if err != nil {
		t.Fatalf("Signature should be valid but verification returned error: %s", err)
		return
	}

	
	disclosure = []byte{0, 1, 1, 1, 0}
	sig, err = NewSignature(cred, sk, Nym, RandNym, key.Ipk, disclosure, msg, rhindex, cri, rng)
	assert.NoError(t, err)

	err = sig.Ver(disclosure, key.Ipk, msg, attrs, rhindex, &revocationKey.PublicKey, epoch)
	assert.NoError(t, err)

	
	nymsig, err := NewNymSignature(sk, Nym, RandNym, key.Ipk, []byte("testing"), rng)
	assert.NoError(t, err)

	err = nymsig.Ver(Nym, key.Ipk, []byte("testing"))
	if err != nil {
		t.Fatalf("NymSig should be valid but verification returned error: %s", err)
		return
	}
}
