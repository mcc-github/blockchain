/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package idemix

import (
	"github.com/mcc-github/blockchain-amcl/amcl"
	"github.com/mcc-github/blockchain-amcl/amcl/FP256BN"
	"github.com/pkg/errors"
)


func NewNymSignature(sk *FP256BN.BIG, Nym *FP256BN.ECP, RNym *FP256BN.BIG, ipk *IssuerPublicKey, msg []byte, rng *amcl.RAND) (*NymSignature, error) {
	if sk == nil || Nym == nil || RNym == nil || ipk == nil || rng == nil {
		return nil, errors.Errorf("cannot create NymSignature: received nil input")
	}

	Nonce := RandModOrder(rng)

	HRand := EcpFromProto(ipk.HRand)
	HSk := EcpFromProto(ipk.HSk)

	
	

	
	rSk := RandModOrder(rng)
	rRNym := RandModOrder(rng)

	
	t := HSk.Mul2(rSk, HRand, rRNym)

	
	
	
	
	
	
	
	proofData := make([]byte, len([]byte(signLabel))+2*(2*FieldBytes+1)+FieldBytes+len(msg))
	index := 0
	index = appendBytesString(proofData, index, signLabel)
	index = appendBytesG1(proofData, index, t)
	index = appendBytesG1(proofData, index, Nym)
	copy(proofData[index:], ipk.Hash)
	index = index + FieldBytes
	copy(proofData[index:], msg)
	c := HashModOrder(proofData)

	
	index = 0
	proofData = proofData[:2*FieldBytes]
	index = appendBytesBig(proofData, index, c)
	index = appendBytesBig(proofData, index, Nonce)
	ProofC := HashModOrder(proofData)

	
	ProofSSk := Modadd(rSk, FP256BN.Modmul(ProofC, sk, GroupOrder), GroupOrder)
	ProofSRNym := Modadd(rRNym, FP256BN.Modmul(ProofC, RNym, GroupOrder), GroupOrder)

	
	return &NymSignature{
		ProofC:     BigToBytes(ProofC),
		ProofSSk:   BigToBytes(ProofSSk),
		ProofSRNym: BigToBytes(ProofSRNym),
		Nonce:      BigToBytes(Nonce)}, nil
}


func (sig *NymSignature) Ver(nym *FP256BN.ECP, ipk *IssuerPublicKey, msg []byte) error {
	ProofC := FP256BN.FromBytes(sig.GetProofC())
	ProofSSk := FP256BN.FromBytes(sig.GetProofSSk())
	ProofSRNym := FP256BN.FromBytes(sig.GetProofSRNym())

	Nonce := FP256BN.FromBytes(sig.GetNonce())

	HRand := EcpFromProto(ipk.HRand)
	HSk := EcpFromProto(ipk.HSk)

	t := HSk.Mul2(ProofSSk, HRand, ProofSRNym)
	t.Sub(nym.Mul(ProofC))

	
	
	
	
	
	
	proofData := make([]byte, len([]byte(signLabel))+2*(2*FieldBytes+1)+FieldBytes+len(msg))
	index := 0
	index = appendBytesString(proofData, index, signLabel)
	index = appendBytesG1(proofData, index, t)
	index = appendBytesG1(proofData, index, nym)
	copy(proofData[index:], ipk.Hash)
	index = index + FieldBytes
	copy(proofData[index:], msg)

	c := HashModOrder(proofData)
	index = 0
	proofData = proofData[:2*FieldBytes]
	index = appendBytesBig(proofData, index, c)
	index = appendBytesBig(proofData, index, Nonce)
	if *ProofC != *HashModOrder(proofData) {
		return errors.Errorf("pseudonym signature invalid: zero-knowledge proof is invalid")
	}

	return nil
}
