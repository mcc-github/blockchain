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























func NewCredential(key *IssuerKey, m *CredRequest, attrs []*FP256BN.BIG, rng *amcl.RAND) (*Credential, error) {
	
	err := m.Check(key.Ipk)
	if err != nil {
		return nil, err
	}

	if len(attrs) != len(key.Ipk.AttributeNames) {
		return nil, errors.Errorf("incorrect number of attribute values passed")
	}

	
	
	E := RandModOrder(rng)
	S := RandModOrder(rng)

	B := FP256BN.NewECP()
	B.Copy(GenG1)
	Nym := EcpFromProto(m.Nym)
	B.Add(Nym)
	B.Add(EcpFromProto(key.Ipk.HRand).Mul(S))

	
	for i := 0; i < len(attrs)/2; i++ {
		B.Add(EcpFromProto(key.Ipk.HAttrs[2*i]).Mul2(attrs[2*i], EcpFromProto(key.Ipk.HAttrs[2*i+1]), attrs[2*i+1]))
	}
	if len(attrs)%2 != 0 {
		B.Add(EcpFromProto(key.Ipk.HAttrs[len(attrs)-1]).Mul(attrs[len(attrs)-1]))
	}

	Exp := Modadd(FP256BN.FromBytes(key.GetIsk()), E, GroupOrder)
	Exp.Invmodp(GroupOrder)
	A := B.Mul(Exp)

	CredAttrs := make([][]byte, len(attrs))
	for index, attribute := range attrs {
		CredAttrs[index] = BigToBytes(attribute)
	}

	return &Credential{
		A:     EcpToProto(A),
		B:     EcpToProto(B),
		E:     BigToBytes(E),
		S:     BigToBytes(S),
		Attrs: CredAttrs}, nil
}



func (cred *Credential) Ver(sk *FP256BN.BIG, ipk *IssuerPublicKey) error {

	
	A := EcpFromProto(cred.GetA())
	B := EcpFromProto(cred.GetB())
	E := FP256BN.FromBytes(cred.GetE())
	S := FP256BN.FromBytes(cred.GetS())

	
	for i := 0; i < len(cred.GetAttrs()); i++ {
		if cred.Attrs[i] == nil {
			return errors.Errorf("credential has no value for attribute %s", ipk.AttributeNames[i])
		}
	}

	
	BPrime := FP256BN.NewECP()
	BPrime.Copy(GenG1)
	BPrime.Add(EcpFromProto(ipk.HSk).Mul2(sk, EcpFromProto(ipk.HRand), S))
	for i := 0; i < len(cred.Attrs)/2; i++ {
		BPrime.Add(EcpFromProto(ipk.HAttrs[2*i]).Mul2(FP256BN.FromBytes(cred.Attrs[2*i]), EcpFromProto(ipk.HAttrs[2*i+1]), FP256BN.FromBytes(cred.Attrs[2*i+1])))
	}
	if len(cred.Attrs)%2 != 0 {
		BPrime.Add(EcpFromProto(ipk.HAttrs[len(cred.Attrs)-1]).Mul(FP256BN.FromBytes(cred.Attrs[len(cred.Attrs)-1])))
	}
	if !B.Equals(BPrime) {
		return errors.Errorf("b-value from credential does not match the attribute values")
	}

	a := GenG2.Mul(E)
	a.Add(Ecp2FromProto(ipk.W))
	a.Affine()

	if !FP256BN.Fexp(FP256BN.Ate(a, A)).Equals(FP256BN.Fexp(FP256BN.Ate(GenG2, B))) {
		return errors.Errorf("credential is not cryptographically valid")
	}
	return nil
}
