/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package bccsp

const (
	
	
	
	
	ECDSA = "ECDSA"

	
	ECDSAP256 = "ECDSAP256"

	
	ECDSAP384 = "ECDSAP384"

	
	ECDSAReRand = "ECDSA_RERAND"

	
	
	
	RSA = "RSA"
	
	RSA1024 = "RSA1024"
	
	RSA2048 = "RSA2048"
	
	RSA3072 = "RSA3072"
	
	RSA4096 = "RSA4096"

	
	
	
	AES = "AES"
	
	AES128 = "AES128"
	
	AES192 = "AES192"
	
	AES256 = "AES256"

	
	HMAC = "HMAC"
	
	HMACTruncated256 = "HMAC_TRUNCATED_256"

	
	
	
	SHA = "SHA"

	
	SHA2 = "SHA2"
	
	SHA3 = "SHA3"

	
	SHA256 = "SHA256"
	
	SHA384 = "SHA384"
	
	SHA3_256 = "SHA3_256"
	
	SHA3_384 = "SHA3_384"

	
	X509Certificate = "X509Certificate"
)


type ECDSAKeyGenOpts struct {
	Temporary bool
}


func (opts *ECDSAKeyGenOpts) Algorithm() string {
	return ECDSA
}



func (opts *ECDSAKeyGenOpts) Ephemeral() bool {
	return opts.Temporary
}


type ECDSAPKIXPublicKeyImportOpts struct {
	Temporary bool
}


func (opts *ECDSAPKIXPublicKeyImportOpts) Algorithm() string {
	return ECDSA
}



func (opts *ECDSAPKIXPublicKeyImportOpts) Ephemeral() bool {
	return opts.Temporary
}



type ECDSAPrivateKeyImportOpts struct {
	Temporary bool
}


func (opts *ECDSAPrivateKeyImportOpts) Algorithm() string {
	return ECDSA
}



func (opts *ECDSAPrivateKeyImportOpts) Ephemeral() bool {
	return opts.Temporary
}


type ECDSAGoPublicKeyImportOpts struct {
	Temporary bool
}


func (opts *ECDSAGoPublicKeyImportOpts) Algorithm() string {
	return ECDSA
}



func (opts *ECDSAGoPublicKeyImportOpts) Ephemeral() bool {
	return opts.Temporary
}


type ECDSAReRandKeyOpts struct {
	Temporary bool
	Expansion []byte
}


func (opts *ECDSAReRandKeyOpts) Algorithm() string {
	return ECDSAReRand
}



func (opts *ECDSAReRandKeyOpts) Ephemeral() bool {
	return opts.Temporary
}


func (opts *ECDSAReRandKeyOpts) ExpansionValue() []byte {
	return opts.Expansion
}


type AESKeyGenOpts struct {
	Temporary bool
}


func (opts *AESKeyGenOpts) Algorithm() string {
	return AES
}



func (opts *AESKeyGenOpts) Ephemeral() bool {
	return opts.Temporary
}



type HMACTruncated256AESDeriveKeyOpts struct {
	Temporary bool
	Arg       []byte
}


func (opts *HMACTruncated256AESDeriveKeyOpts) Algorithm() string {
	return HMACTruncated256
}



func (opts *HMACTruncated256AESDeriveKeyOpts) Ephemeral() bool {
	return opts.Temporary
}


func (opts *HMACTruncated256AESDeriveKeyOpts) Argument() []byte {
	return opts.Arg
}


type HMACDeriveKeyOpts struct {
	Temporary bool
	Arg       []byte
}


func (opts *HMACDeriveKeyOpts) Algorithm() string {
	return HMAC
}



func (opts *HMACDeriveKeyOpts) Ephemeral() bool {
	return opts.Temporary
}


func (opts *HMACDeriveKeyOpts) Argument() []byte {
	return opts.Arg
}


type AES256ImportKeyOpts struct {
	Temporary bool
}


func (opts *AES256ImportKeyOpts) Algorithm() string {
	return AES
}



func (opts *AES256ImportKeyOpts) Ephemeral() bool {
	return opts.Temporary
}


type HMACImportKeyOpts struct {
	Temporary bool
}


func (opts *HMACImportKeyOpts) Algorithm() string {
	return HMAC
}



func (opts *HMACImportKeyOpts) Ephemeral() bool {
	return opts.Temporary
}


type SHAOpts struct {
}


func (opts *SHAOpts) Algorithm() string {
	return SHA
}


type RSAKeyGenOpts struct {
	Temporary bool
}


func (opts *RSAKeyGenOpts) Algorithm() string {
	return RSA
}



func (opts *RSAKeyGenOpts) Ephemeral() bool {
	return opts.Temporary
}


type RSAGoPublicKeyImportOpts struct {
	Temporary bool
}


func (opts *RSAGoPublicKeyImportOpts) Algorithm() string {
	return RSA
}



func (opts *RSAGoPublicKeyImportOpts) Ephemeral() bool {
	return opts.Temporary
}


type X509PublicKeyImportOpts struct {
	Temporary bool
}


func (opts *X509PublicKeyImportOpts) Algorithm() string {
	return X509Certificate
}



func (opts *X509PublicKeyImportOpts) Ephemeral() bool {
	return opts.Temporary
}
