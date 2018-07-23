/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package idemix

import (
	"crypto/ecdsa"

	"sort"

	"github.com/mcc-github/blockchain-amcl/amcl"
	"github.com/mcc-github/blockchain-amcl/amcl/FP256BN"
	"github.com/pkg/errors"
)


const signLabel = "sign"















func hiddenIndices(Disclosure []byte) []int {
	HiddenIndices := make([]int, 0)
	for index, disclose := range Disclosure {
		if disclose == 0 {
			HiddenIndices = append(HiddenIndices, index)
		}
	}
	return HiddenIndices
}






func NewSignature(cred *Credential, sk *FP256BN.BIG, Nym *FP256BN.ECP, RNym *FP256BN.BIG, ipk *IssuerPublicKey, Disclosure []byte, msg []byte, rhIndex int, cri *CredentialRevocationInformation, rng *amcl.RAND) (*Signature, error) {
	if cred == nil || sk == nil || Nym == nil || RNym == nil || ipk == nil || rng == nil || cri == nil {
		return nil, errors.Errorf("cannot create idemix signature: received nil input")
	}

	if rhIndex < 0 || rhIndex >= len(ipk.AttributeNames) || len(Disclosure) != len(ipk.AttributeNames) {
		return nil, errors.Errorf("cannot create idemix signature: received invalid input")
	}

	if cri.RevocationAlg != int32(ALG_NO_REVOCATION) && Disclosure[rhIndex] == 1 {
		return nil, errors.Errorf("Attribute %d is disclosed but also used as revocation handle attribute, which should remain hidden.", rhIndex)
	}

	HiddenIndices := hiddenIndices(Disclosure)

	
	r1 := RandModOrder(rng)
	r2 := RandModOrder(rng)
	r3 := FP256BN.NewBIGcopy(r1)
	r3.Invmodp(GroupOrder)

	Nonce := RandModOrder(rng)

	A := EcpFromProto(cred.A)
	B := EcpFromProto(cred.B)

	APrime := FP256BN.G1mul(A, r1) 
	ABar := FP256BN.G1mul(B, r1)
	ABar.Sub(FP256BN.G1mul(APrime, FP256BN.FromBytes(cred.E))) 

	BPrime := FP256BN.G1mul(B, r1)
	HRand := EcpFromProto(ipk.HRand)
	HSk := EcpFromProto(ipk.HSk)

	BPrime.Sub(FP256BN.G1mul(HRand, r2)) 

	S := FP256BN.FromBytes(cred.S)
	E := FP256BN.FromBytes(cred.E)
	sPrime := Modsub(S, FP256BN.Modmul(r2, r3, GroupOrder), GroupOrder)

	
	rSk := RandModOrder(rng)
	re := RandModOrder(rng)
	rR2 := RandModOrder(rng)
	rR3 := RandModOrder(rng)
	rSPrime := RandModOrder(rng)
	rRNym := RandModOrder(rng)
	rAttrs := make([]*FP256BN.BIG, len(HiddenIndices))
	for i := range HiddenIndices {
		rAttrs[i] = RandModOrder(rng)
	}

	
	prover, err := getNonRevocationProver(RevocationAlgorithm(cri.RevocationAlg))
	if err != nil {
		return nil, err
	}
	nonRevokedProofHashData, err := prover.getFSContribution(FP256BN.FromBytes(cred.Attrs[rhIndex]), rAttrs[sort.SearchInts(HiddenIndices, rhIndex)], cri, rng)
	if err != nil {
		return nil, errors.Wrap(err, "failed to compute non-revoked proof")
	}

	t1 := APrime.Mul2(re, HRand, rR2)
	t2 := FP256BN.G1mul(HRand, rSPrime)
	t2.Add(BPrime.Mul2(rR3, HSk, rSk))

	for i := 0; i < len(HiddenIndices)/2; i++ {
		t2.Add(EcpFromProto(ipk.HAttrs[HiddenIndices[2*i]]).Mul2(rAttrs[2*i], EcpFromProto(ipk.HAttrs[HiddenIndices[2*i+1]]), rAttrs[2*i+1]))
	}
	if len(HiddenIndices)%2 != 0 {
		t2.Add(FP256BN.G1mul(EcpFromProto(ipk.HAttrs[HiddenIndices[len(HiddenIndices)-1]]), rAttrs[len(HiddenIndices)-1]))
	}

	t3 := HSk.Mul2(rSk, HRand, rRNym)

	
	
	
	
	
	
	
	proofData := make([]byte, len([]byte(signLabel))+7*(2*FieldBytes+1)+FieldBytes+len(Disclosure)+len(msg)+ProofBytes[RevocationAlgorithm(cri.RevocationAlg)])
	index := 0
	index = appendBytesString(proofData, index, signLabel)
	index = appendBytesG1(proofData, index, t1)
	index = appendBytesG1(proofData, index, t2)
	index = appendBytesG1(proofData, index, t3)
	index = appendBytesG1(proofData, index, APrime)
	index = appendBytesG1(proofData, index, ABar)
	index = appendBytesG1(proofData, index, BPrime)
	index = appendBytesG1(proofData, index, Nym)
	index = appendBytes(proofData, index, nonRevokedProofHashData)
	copy(proofData[index:], ipk.Hash)
	index = index + FieldBytes
	copy(proofData[index:], Disclosure)
	index = index + len(Disclosure)
	copy(proofData[index:], msg)
	c := HashModOrder(proofData)

	
	index = 0
	proofData = proofData[:2*FieldBytes]
	index = appendBytesBig(proofData, index, c)
	index = appendBytesBig(proofData, index, Nonce)
	ProofC := HashModOrder(proofData)
	ProofSSk := Modadd(rSk, FP256BN.Modmul(ProofC, sk, GroupOrder), GroupOrder)
	ProofSE := Modsub(re, FP256BN.Modmul(ProofC, E, GroupOrder), GroupOrder)
	ProofSR2 := Modadd(rR2, FP256BN.Modmul(ProofC, r2, GroupOrder), GroupOrder)
	ProofSR3 := Modsub(rR3, FP256BN.Modmul(ProofC, r3, GroupOrder), GroupOrder)
	ProofSSPrime := Modadd(rSPrime, FP256BN.Modmul(ProofC, sPrime, GroupOrder), GroupOrder)
	ProofSRNym := Modadd(rRNym, FP256BN.Modmul(ProofC, RNym, GroupOrder), GroupOrder)

	ProofSAttrs := make([][]byte, len(HiddenIndices))
	for i, j := range HiddenIndices {
		ProofSAttrs[i] = BigToBytes(Modadd(rAttrs[i], FP256BN.Modmul(ProofC, FP256BN.FromBytes(cred.Attrs[j]), GroupOrder), GroupOrder))
	}

	nonRevokedProof, err := prover.getNonRevokedProof(ProofC)
	if err != nil {
		return nil, err
	}

	return &Signature{
			APrime:             EcpToProto(APrime),
			ABar:               EcpToProto(ABar),
			BPrime:             EcpToProto(BPrime),
			ProofC:             BigToBytes(ProofC),
			ProofSSk:           BigToBytes(ProofSSk),
			ProofSE:            BigToBytes(ProofSE),
			ProofSR2:           BigToBytes(ProofSR2),
			ProofSR3:           BigToBytes(ProofSR3),
			ProofSSPrime:       BigToBytes(ProofSSPrime),
			ProofSAttrs:        ProofSAttrs,
			Nonce:              BigToBytes(Nonce),
			Nym:                EcpToProto(Nym),
			ProofSRNym:         BigToBytes(ProofSRNym),
			RevocationEpochPk:  cri.EpochPk,
			RevocationPkSig:    cri.EpochPkSig,
			Epoch:              cri.Epoch,
			NonRevocationProof: nonRevokedProof},
		nil
}





func (sig *Signature) Ver(Disclosure []byte, ipk *IssuerPublicKey, msg []byte, attributeValues []*FP256BN.BIG, rhIndex int, revPk *ecdsa.PublicKey, epoch int) error {
	if ipk == nil || revPk == nil {
		return errors.Errorf("cannot verify idemix signature: received nil input")
	}

	if rhIndex < 0 || rhIndex >= len(ipk.AttributeNames) || len(Disclosure) != len(ipk.AttributeNames) {
		return errors.Errorf("cannot verify idemix signature: received invalid input")
	}

	if sig.NonRevocationProof.RevocationAlg != int32(ALG_NO_REVOCATION) && Disclosure[rhIndex] == 1 {
		return errors.Errorf("Attribute %d is disclosed but is also used as revocation handle, which should remain hidden.", rhIndex)
	}

	HiddenIndices := hiddenIndices(Disclosure)

	APrime := EcpFromProto(sig.GetAPrime())
	ABar := EcpFromProto(sig.GetABar())
	BPrime := EcpFromProto(sig.GetBPrime())
	Nym := EcpFromProto(sig.GetNym())
	ProofC := FP256BN.FromBytes(sig.GetProofC())
	ProofSSk := FP256BN.FromBytes(sig.GetProofSSk())
	ProofSE := FP256BN.FromBytes(sig.GetProofSE())
	ProofSR2 := FP256BN.FromBytes(sig.GetProofSR2())
	ProofSR3 := FP256BN.FromBytes(sig.GetProofSR3())
	ProofSSPrime := FP256BN.FromBytes(sig.GetProofSSPrime())
	ProofSRNym := FP256BN.FromBytes(sig.GetProofSRNym())
	ProofSAttrs := make([]*FP256BN.BIG, len(sig.GetProofSAttrs()))

	if len(sig.ProofSAttrs) != len(HiddenIndices) {
		return errors.Errorf("signature invalid: incorrect amount of s-values for AttributeProofSpec")
	}
	for i, b := range sig.ProofSAttrs {
		ProofSAttrs[i] = FP256BN.FromBytes(b)
	}

	Nonce := FP256BN.FromBytes(sig.GetNonce())

	W := Ecp2FromProto(ipk.W)
	HRand := EcpFromProto(ipk.HRand)
	HSk := EcpFromProto(ipk.HSk)

	if APrime.Is_infinity() {
		return errors.Errorf("signature invalid: APrime = 1")
	}
	temp1 := FP256BN.Ate(W, APrime)
	temp2 := FP256BN.Ate(GenG2, ABar)
	temp2.Inverse()
	temp1.Mul(temp2)
	if !FP256BN.Fexp(temp1).Isunity() {
		return errors.Errorf("signature invalid: APrime and ABar don't have the expected structure")
	}

	t1 := APrime.Mul2(ProofSE, HRand, ProofSR2)
	temp := FP256BN.NewECP()
	temp.Copy(ABar)
	temp.Sub(BPrime)
	t1.Sub(FP256BN.G1mul(temp, ProofC))

	t2 := FP256BN.G1mul(HRand, ProofSSPrime)
	t2.Add(BPrime.Mul2(ProofSR3, HSk, ProofSSk))

	for i := 0; i < len(HiddenIndices)/2; i++ {
		t2.Add(EcpFromProto(ipk.HAttrs[HiddenIndices[2*i]]).Mul2(ProofSAttrs[2*i], EcpFromProto(ipk.HAttrs[HiddenIndices[2*i+1]]), ProofSAttrs[2*i+1]))
	}
	if len(HiddenIndices)%2 != 0 {
		t2.Add(FP256BN.G1mul(EcpFromProto(ipk.HAttrs[HiddenIndices[len(HiddenIndices)-1]]), ProofSAttrs[len(HiddenIndices)-1]))
	}

	temp = FP256BN.NewECP()
	temp.Copy(GenG1)

	for index, disclose := range Disclosure {
		if disclose != 0 {
			temp.Add(FP256BN.G1mul(EcpFromProto(ipk.HAttrs[index]), attributeValues[index]))
		}
	}
	t2.Add(FP256BN.G1mul(temp, ProofC))

	t3 := HSk.Mul2(ProofSSk, HRand, ProofSRNym)
	t3.Sub(Nym.Mul(ProofC))

	nonRevokedVer, err := getNonRevocationVerifier(RevocationAlgorithm(sig.NonRevocationProof.RevocationAlg))
	if err != nil {
		return err
	}

	i := sort.SearchInts(HiddenIndices, rhIndex)
	proofSRh := ProofSAttrs[i]
	nonRevokedProofBytes, err := nonRevokedVer.recomputeFSContribution(sig.NonRevocationProof, ProofC, Ecp2FromProto(sig.RevocationEpochPk), proofSRh)
	if err != nil {
		return err
	}

	
	
	
	
	
	
	proofData := make([]byte, len([]byte(signLabel))+7*(2*FieldBytes+1)+FieldBytes+len(Disclosure)+len(msg)+ProofBytes[RevocationAlgorithm(sig.NonRevocationProof.RevocationAlg)])
	index := 0
	index = appendBytesString(proofData, index, signLabel)
	index = appendBytesG1(proofData, index, t1)
	index = appendBytesG1(proofData, index, t2)
	index = appendBytesG1(proofData, index, t3)
	index = appendBytesG1(proofData, index, APrime)
	index = appendBytesG1(proofData, index, ABar)
	index = appendBytesG1(proofData, index, BPrime)
	index = appendBytesG1(proofData, index, Nym)
	index = appendBytes(proofData, index, nonRevokedProofBytes)
	copy(proofData[index:], ipk.Hash)
	index = index + FieldBytes
	copy(proofData[index:], Disclosure)
	index = index + len(Disclosure)
	copy(proofData[index:], msg)

	c := HashModOrder(proofData)
	index = 0
	proofData = proofData[:2*FieldBytes]
	index = appendBytesBig(proofData, index, c)
	index = appendBytesBig(proofData, index, Nonce)
	if *ProofC != *HashModOrder(proofData) {
		return errors.Errorf("signature invalid: zero-knowledge proof is invalid")
	}

	return nil
}