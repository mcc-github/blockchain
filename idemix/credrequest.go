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


const credRequestLabel = "credRequest"


















func NewCredRequest(sk *FP256BN.BIG, IssuerNonce []byte, ipk *IssuerPublicKey, rng *amcl.RAND) *CredRequest {
	
	HSk := EcpFromProto(ipk.HSk)
	Nym := HSk.Mul(sk)

	

	
	rSk := RandModOrder(rng)

	
	t := HSk.Mul(rSk) 

	
	
	
	
	
	
	proofData := make([]byte, len([]byte(credRequestLabel))+3*(2*FieldBytes+1)+2*FieldBytes)
	index := 0
	index = appendBytesString(proofData, index, credRequestLabel)
	index = appendBytesG1(proofData, index, t)
	index = appendBytesG1(proofData, index, HSk)
	index = appendBytesG1(proofData, index, Nym)
	index = appendBytes(proofData, index, IssuerNonce)
	copy(proofData[index:], ipk.Hash)
	proofC := HashModOrder(proofData)

	
	proofS := Modadd(FP256BN.Modmul(proofC, sk, GroupOrder), rSk, GroupOrder) 

	
	return &CredRequest{
		Nym:         EcpToProto(Nym),
		IssuerNonce: IssuerNonce,
		ProofC:      BigToBytes(proofC),
		ProofS:      BigToBytes(proofS)}
}


func (m *CredRequest) Check(ipk *IssuerPublicKey) error {
	Nym := EcpFromProto(m.GetNym())
	IssuerNonce := m.GetIssuerNonce()
	ProofC := FP256BN.FromBytes(m.GetProofC())
	ProofS := FP256BN.FromBytes(m.GetProofS())

	HSk := EcpFromProto(ipk.HSk)

	if Nym == nil || IssuerNonce == nil || ProofC == nil || ProofS == nil {
		return errors.Errorf("one of the proof values is undefined")
	}

	

	
	t := HSk.Mul(ProofS)
	t.Sub(Nym.Mul(ProofC)) 

	
	proofData := make([]byte, len([]byte(credRequestLabel))+3*(2*FieldBytes+1)+2*FieldBytes)
	index := 0
	index = appendBytesString(proofData, index, credRequestLabel)
	index = appendBytesG1(proofData, index, t)
	index = appendBytesG1(proofData, index, HSk)
	index = appendBytesG1(proofData, index, Nym)
	index = appendBytes(proofData, index, IssuerNonce)
	copy(proofData[index:], ipk.Hash)

	if *ProofC != *HashModOrder(proofData) {
		return errors.Errorf("zero knowledge proof is invalid")
	}

	return nil
}
