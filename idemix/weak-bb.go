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


func WBBKeyGen(rng *amcl.RAND) (*FP256BN.BIG, *FP256BN.ECP2) {
	
	sk := RandModOrder(rng)
	
	pk := GenG2.Mul(sk)
	return sk, pk
}


func WBBSign(sk *FP256BN.BIG, m *FP256BN.BIG) *FP256BN.ECP {
	
	exp := Modadd(sk, m, GroupOrder)
	exp.Invmodp(GroupOrder)

	
	return GenG1.Mul(exp)
}


func WBBVerify(pk *FP256BN.ECP2, sig *FP256BN.ECP, m *FP256BN.BIG) error {
	if pk == nil || sig == nil || m == nil {
		return errors.Errorf("Weak-BB signature invalid: received nil input")
	}
	
	P := FP256BN.NewECP2()
	P.Copy(pk)
	P.Add(GenG2.Mul(m))
	P.Affine()
	
	if !FP256BN.Fexp(FP256BN.Ate(P, sig)).Equals(GenGT) {
		return errors.Errorf("Weak-BB signature is invalid")
	}
	return nil
}
