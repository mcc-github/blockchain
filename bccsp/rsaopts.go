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


type RSA1024KeyGenOpts struct {
	Temporary bool
}


func (opts *RSA1024KeyGenOpts) Algorithm() string {
	return RSA1024
}



func (opts *RSA1024KeyGenOpts) Ephemeral() bool {
	return opts.Temporary
}


type RSA2048KeyGenOpts struct {
	Temporary bool
}


func (opts *RSA2048KeyGenOpts) Algorithm() string {
	return RSA2048
}



func (opts *RSA2048KeyGenOpts) Ephemeral() bool {
	return opts.Temporary
}


type RSA3072KeyGenOpts struct {
	Temporary bool
}


func (opts *RSA3072KeyGenOpts) Algorithm() string {
	return RSA3072
}



func (opts *RSA3072KeyGenOpts) Ephemeral() bool {
	return opts.Temporary
}


type RSA4096KeyGenOpts struct {
	Temporary bool
}


func (opts *RSA4096KeyGenOpts) Algorithm() string {
	return RSA4096
}



func (opts *RSA4096KeyGenOpts) Ephemeral() bool {
	return opts.Temporary
}
