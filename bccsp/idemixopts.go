/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package bccsp

import (
	"crypto"
)


type RevocationAlgorithm int32

const (
	
	IDEMIX = "IDEMIX"
)

const (
	
	AlgNoRevocation RevocationAlgorithm = iota
)



type IdemixIssuerKeyGenOpts struct {
	
	Temporary bool
	
	AttributeNames []string
}


func (*IdemixIssuerKeyGenOpts) Algorithm() string {
	return IDEMIX
}



func (o *IdemixIssuerKeyGenOpts) Ephemeral() bool {
	return o.Temporary
}


type IdemixIssuerPublicKeyImportOpts struct {
	Temporary bool
	
	AttributeNames []string
}


func (*IdemixIssuerPublicKeyImportOpts) Algorithm() string {
	return IDEMIX
}



func (o *IdemixIssuerPublicKeyImportOpts) Ephemeral() bool {
	return o.Temporary
}


type IdemixUserSecretKeyGenOpts struct {
	Temporary bool
}


func (*IdemixUserSecretKeyGenOpts) Algorithm() string {
	return IDEMIX
}



func (o *IdemixUserSecretKeyGenOpts) Ephemeral() bool {
	return o.Temporary
}


type IdemixUserSecretKeyImportOpts struct {
	Temporary bool
}


func (*IdemixUserSecretKeyImportOpts) Algorithm() string {
	return IDEMIX
}



func (o *IdemixUserSecretKeyImportOpts) Ephemeral() bool {
	return o.Temporary
}



type IdemixNymKeyDerivationOpts struct {
	
	Temporary bool
	
	IssuerPK Key
}


func (*IdemixNymKeyDerivationOpts) Algorithm() string {
	return IDEMIX
}



func (o *IdemixNymKeyDerivationOpts) Ephemeral() bool {
	return o.Temporary
}



func (o *IdemixNymKeyDerivationOpts) IssuerPublicKey() Key {
	return o.IssuerPK
}


type IdemixNymPublicKeyImportOpts struct {
	
	Temporary bool
}


func (*IdemixNymPublicKeyImportOpts) Algorithm() string {
	return IDEMIX
}



func (o *IdemixNymPublicKeyImportOpts) Ephemeral() bool {
	return o.Temporary
}


type IdemixCredentialRequestSignerOpts struct {
	
	
	Attributes []int
	
	IssuerPK Key
	
	
	IssuerNonce []byte
	
	H crypto.Hash
}

func (o *IdemixCredentialRequestSignerOpts) HashFunc() crypto.Hash {
	return o.H
}



func (o *IdemixCredentialRequestSignerOpts) IssuerPublicKey() Key {
	return o.IssuerPK
}


type IdemixAttributeType int

const (
	
	IdemixHiddenAttribute IdemixAttributeType = iota
	
	IdemixBytesAttribute
	
	IdemixIntAttribute
)

type IdemixAttribute struct {
	
	Type IdemixAttributeType
	
	Value interface{}
}


type IdemixCredentialSignerOpts struct {
	
	Attributes []IdemixAttribute
	
	IssuerPK Key
	
	H crypto.Hash
}




func (o *IdemixCredentialSignerOpts) HashFunc() crypto.Hash {
	return o.H
}

func (o *IdemixCredentialSignerOpts) IssuerPublicKey() Key {
	return o.IssuerPK
}


type IdemixSignerOpts struct {
	
	Nym Key
	
	IssuerPK Key
	
	Credential []byte
	
	
	
	
	
	
	Attributes []IdemixAttribute
	
	
	RhIndex int
	
	CRI []byte
	
	Epoch int
	
	RevocationPublicKey Key
	
	H crypto.Hash
}

func (o *IdemixSignerOpts) HashFunc() crypto.Hash {
	return o.H
}


type IdemixNymSignerOpts struct {
	
	Nym Key
	
	IssuerPK Key
	
	H crypto.Hash
}




func (o *IdemixNymSignerOpts) HashFunc() crypto.Hash {
	return o.H
}


type IdemixRevocationKeyGenOpts struct {
	
	Temporary bool
}


func (*IdemixRevocationKeyGenOpts) Algorithm() string {
	return IDEMIX
}



func (o *IdemixRevocationKeyGenOpts) Ephemeral() bool {
	return o.Temporary
}


type IdemixRevocationPublicKeyImportOpts struct {
	Temporary bool
}


func (*IdemixRevocationPublicKeyImportOpts) Algorithm() string {
	return IDEMIX
}



func (o *IdemixRevocationPublicKeyImportOpts) Ephemeral() bool {
	return o.Temporary
}




type IdemixCRISignerOpts struct {
	Epoch               int
	RevocationAlgorithm RevocationAlgorithm
	UnrevokedHandles    [][]byte
	
	H crypto.Hash
}

func (o *IdemixCRISignerOpts) HashFunc() crypto.Hash {
	return o.H
}
