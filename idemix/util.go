/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package idemix

import (
	"crypto/rand"
	"crypto/sha256"

	"github.com/mcc-github/blockchain-amcl/amcl"
	"github.com/mcc-github/blockchain-amcl/amcl/FP256BN"
	"github.com/pkg/errors"
)


var GenG1 = FP256BN.NewECPbigs(
	FP256BN.NewBIGints(FP256BN.CURVE_Gx),
	FP256BN.NewBIGints(FP256BN.CURVE_Gy))


var GenG2 = FP256BN.NewECP2fp2s(
	FP256BN.NewFP2bigs(FP256BN.NewBIGints(FP256BN.CURVE_Pxa), FP256BN.NewBIGints(FP256BN.CURVE_Pxb)),
	FP256BN.NewFP2bigs(FP256BN.NewBIGints(FP256BN.CURVE_Pya), FP256BN.NewBIGints(FP256BN.CURVE_Pyb)))


var GenGT = FP256BN.Fexp(FP256BN.Ate(GenG2, GenG1))


var GroupOrder = FP256BN.NewBIGints(FP256BN.CURVE_Order)


var FieldBytes = int(FP256BN.MODBYTES)


func RandModOrder(rng *amcl.RAND) *FP256BN.BIG {
	
	q := FP256BN.NewBIGints(FP256BN.CURVE_Order)

	
	return FP256BN.Randomnum(q, rng)
}


func HashModOrder(data []byte) *FP256BN.BIG {
	digest := sha256.Sum256(data)
	digestBig := FP256BN.FromBytes(digest[:])
	digestBig.Mod(GroupOrder)
	return digestBig
}

func appendBytes(data []byte, index int, bytesToAdd []byte) int {
	copy(data[index:], bytesToAdd)
	return index + len(bytesToAdd)
}
func appendBytesG1(data []byte, index int, E *FP256BN.ECP) int {
	length := 2*FieldBytes + 1
	E.ToBytes(data[index:index+length], false)
	return index + length
}
func EcpToBytes(E *FP256BN.ECP) []byte {
	length := 2*FieldBytes + 1
	res := make([]byte, length)
	E.ToBytes(res, false)
	return res
}
func appendBytesG2(data []byte, index int, E *FP256BN.ECP2) int {
	length := 4 * FieldBytes
	E.ToBytes(data[index : index+length])
	return index + length
}
func appendBytesBig(data []byte, index int, B *FP256BN.BIG) int {
	length := FieldBytes
	B.ToBytes(data[index : index+length])
	return index + length
}
func appendBytesString(data []byte, index int, s string) int {
	bytes := []byte(s)
	copy(data[index:], bytes)
	return index + len(bytes)
}


func MakeNym(sk *FP256BN.BIG, IPk *IssuerPublicKey, rng *amcl.RAND) (*FP256BN.ECP, *FP256BN.BIG) {
	
	
	RandNym := RandModOrder(rng)
	Nym := EcpFromProto(IPk.HSk).Mul2(sk, EcpFromProto(IPk.HRand), RandNym)
	return Nym, RandNym
}


func BigToBytes(big *FP256BN.BIG) []byte {
	ret := make([]byte, FieldBytes)
	big.ToBytes(ret)
	return ret
}


func EcpToProto(p *FP256BN.ECP) *ECP {
	return &ECP{
		X: BigToBytes(p.GetX()),
		Y: BigToBytes(p.GetY())}
}


func EcpFromProto(p *ECP) *FP256BN.ECP {
	return FP256BN.NewECPbigs(FP256BN.FromBytes(p.GetX()), FP256BN.FromBytes(p.GetY()))
}


func Ecp2ToProto(p *FP256BN.ECP2) *ECP2 {
	return &ECP2{
		Xa: BigToBytes(p.GetX().GetA()),
		Xb: BigToBytes(p.GetX().GetB()),
		Ya: BigToBytes(p.GetY().GetA()),
		Yb: BigToBytes(p.GetY().GetB())}
}


func Ecp2FromProto(p *ECP2) *FP256BN.ECP2 {
	return FP256BN.NewECP2fp2s(
		FP256BN.NewFP2bigs(FP256BN.FromBytes(p.GetXa()), FP256BN.FromBytes(p.GetXb())),
		FP256BN.NewFP2bigs(FP256BN.FromBytes(p.GetYa()), FP256BN.FromBytes(p.GetYb())))
}


func GetRand() (*amcl.RAND, error) {
	seedLength := 32
	b := make([]byte, seedLength)
	_, err := rand.Read(b)
	if err != nil {
		return nil, errors.Wrap(err, "error getting randomness for seed")
	}
	rng := amcl.NewRAND()
	rng.Clean()
	rng.Seed(seedLength, b)
	return rng, nil
}


func Modadd(a, b, m *FP256BN.BIG) *FP256BN.BIG {
	c := a.Plus(b)
	c.Mod(m)
	return c
}


func Modsub(a, b, m *FP256BN.BIG) *FP256BN.BIG {
	return Modadd(a, FP256BN.Modneg(b, m), m)
}
