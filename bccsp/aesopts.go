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

import "io"


type AES128KeyGenOpts struct {
	Temporary bool
}


func (opts *AES128KeyGenOpts) Algorithm() string {
	return AES128
}



func (opts *AES128KeyGenOpts) Ephemeral() bool {
	return opts.Temporary
}


type AES192KeyGenOpts struct {
	Temporary bool
}


func (opts *AES192KeyGenOpts) Algorithm() string {
	return AES192
}



func (opts *AES192KeyGenOpts) Ephemeral() bool {
	return opts.Temporary
}


type AES256KeyGenOpts struct {
	Temporary bool
}


func (opts *AES256KeyGenOpts) Algorithm() string {
	return AES256
}



func (opts *AES256KeyGenOpts) Ephemeral() bool {
	return opts.Temporary
}






type AESCBCPKCS7ModeOpts struct {
	
	
	
	IV []byte
	
	
	PRNG io.Reader
}
