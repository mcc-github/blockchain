/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package lscc

import "fmt"


type InvalidFunctionErr string

func (f InvalidFunctionErr) Error() string {
	return fmt.Sprintf("invalid function to lscc: %s", string(f))
}


type InvalidArgsLenErr int

func (i InvalidArgsLenErr) Error() string {
	return fmt.Sprintf("invalid number of arguments to lscc: %d", int(i))
}


type TXNotFoundErr string

func (t TXNotFoundErr) Error() string {
	return fmt.Sprintf("transaction not found: %s", string(t))
}


type InvalidDeploymentSpecErr string

func (f InvalidDeploymentSpecErr) Error() string {
	return fmt.Sprintf("invalid deployment spec: %s", string(f))
}


type ExistsErr string

func (t ExistsErr) Error() string {
	return fmt.Sprintf("chaincode with name '%s' already exists", string(t))
}


type NotFoundErr string

func (t NotFoundErr) Error() string {
	return fmt.Sprintf("could not find chaincode with name '%s'", string(t))
}


type InvalidChannelNameErr string

func (f InvalidChannelNameErr) Error() string {
	return fmt.Sprintf("invalid channel name: %s", string(f))
}


type InvalidChaincodeNameErr string

func (f InvalidChaincodeNameErr) Error() string {
	return fmt.Sprintf("invalid chaincode name '%s'. Names must start with an alphanumeric character and can only consist of alphanumerics, '_', and '-'", string(f))
}


type InvalidVersionErr string

func (f InvalidVersionErr) Error() string {
	return fmt.Sprintf("invalid chaincode version '%s'. Versions must not be empty and can only consist of alphanumerics, '_',  '-', '+', and '.'", string(f))
}


type InvalidStatedbArtifactsErr string

func (f InvalidStatedbArtifactsErr) Error() string {
	return fmt.Sprintf("invalid state database artifact: %s", string(f))
}


type ChaincodeMismatchErr string

func (f ChaincodeMismatchErr) Error() string {
	return fmt.Sprintf("chaincode name mismatch: %s", string(f))
}


type MarshallErr string

func (m MarshallErr) Error() string {
	return fmt.Sprintf("error while marshalling: %s", string(m))
}


type IdenticalVersionErr string

func (f IdenticalVersionErr) Error() string {
	return fmt.Sprintf("version already exists for chaincode with name '%s'", string(f))
}


type InvalidCCOnFSError string

func (f InvalidCCOnFSError) Error() string {
	return fmt.Sprintf("chaincode fingerprint mismatch: %s", string(f))
}


type InstantiationPolicyMissing string

func (f InstantiationPolicyMissing) Error() string {
	return "instantiation policy missing"
}


type CollectionsConfigUpgradesNotAllowed string

func (f CollectionsConfigUpgradesNotAllowed) Error() string {
	return "as V1_2 capability is not enabled, collection upgrades are not allowed"
}


type PrivateChannelDataNotAvailable string

func (f PrivateChannelDataNotAvailable) Error() string {
	return "as V1_2 or later capability is not enabled, private channel collections and data are not available"
}
