/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package bccsp

import "crypto"

const (
	
	IDEMIX = "IDEMIX"
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


type IdemixUserSecretKeyGenOpts struct {
	Temporary bool
}


func (*IdemixUserSecretKeyGenOpts) Algorithm() string {
	return IDEMIX
}



func (o *IdemixUserSecretKeyGenOpts) Ephemeral() bool {
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


type IdemixCredentialRequestSignerOpts struct {
	
	
	Attributes []int
	
	IssuerPK Key
	
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
	
	
	
	Disclosure []byte
	
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
