/*
Copyright IBM Corp. 2016 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package cauthdsl

import (
	"sort"

	"github.com/golang/protobuf/proto"
	cb "github.com/mcc-github/blockchain-protos-go/common"
	"github.com/mcc-github/blockchain-protos-go/msp"
	"github.com/mcc-github/blockchain/protoutil"
)


var AcceptAllPolicy *cb.SignaturePolicyEnvelope


var MarshaledAcceptAllPolicy []byte


var RejectAllPolicy *cb.SignaturePolicyEnvelope


var MarshaledRejectAllPolicy []byte

func init() {
	var err error

	AcceptAllPolicy = Envelope(NOutOf(0, []*cb.SignaturePolicy{}), [][]byte{})
	MarshaledAcceptAllPolicy, err = proto.Marshal(AcceptAllPolicy)
	if err != nil {
		panic("Error marshaling trueEnvelope")
	}

	RejectAllPolicy = Envelope(NOutOf(1, []*cb.SignaturePolicy{}), [][]byte{})
	MarshaledRejectAllPolicy, err = proto.Marshal(RejectAllPolicy)
	if err != nil {
		panic("Error marshaling falseEnvelope")
	}
}


func Envelope(policy *cb.SignaturePolicy, identities [][]byte) *cb.SignaturePolicyEnvelope {
	ids := make([]*msp.MSPPrincipal, len(identities))
	for i := range ids {
		ids[i] = &msp.MSPPrincipal{PrincipalClassification: msp.MSPPrincipal_IDENTITY, Principal: identities[i]}
	}

	return &cb.SignaturePolicyEnvelope{
		Version:    0,
		Rule:       policy,
		Identities: ids,
	}
}


func SignedBy(index int32) *cb.SignaturePolicy {
	return &cb.SignaturePolicy{
		Type: &cb.SignaturePolicy_SignedBy{
			SignedBy: index,
		},
	}
}



func SignedByMspMember(mspId string) *cb.SignaturePolicyEnvelope {
	return signedByFabricEntity(mspId, msp.MSPRole_MEMBER)
}



func SignedByMspClient(mspId string) *cb.SignaturePolicyEnvelope {
	return signedByFabricEntity(mspId, msp.MSPRole_CLIENT)
}



func SignedByMspPeer(mspId string) *cb.SignaturePolicyEnvelope {
	return signedByFabricEntity(mspId, msp.MSPRole_PEER)
}



func signedByFabricEntity(mspId string, role msp.MSPRole_MSPRoleType) *cb.SignaturePolicyEnvelope {
	
	principal := &msp.MSPPrincipal{
		PrincipalClassification: msp.MSPPrincipal_ROLE,
		Principal:               protoutil.MarshalOrPanic(&msp.MSPRole{Role: role, MspIdentifier: mspId})}

	
	p := &cb.SignaturePolicyEnvelope{
		Version:    0,
		Rule:       NOutOf(1, []*cb.SignaturePolicy{SignedBy(0)}),
		Identities: []*msp.MSPPrincipal{principal},
	}

	return p
}



func SignedByMspAdmin(mspId string) *cb.SignaturePolicyEnvelope {
	
	principal := &msp.MSPPrincipal{
		PrincipalClassification: msp.MSPPrincipal_ROLE,
		Principal:               protoutil.MarshalOrPanic(&msp.MSPRole{Role: msp.MSPRole_ADMIN, MspIdentifier: mspId})}

	
	p := &cb.SignaturePolicyEnvelope{
		Version:    0,
		Rule:       NOutOf(1, []*cb.SignaturePolicy{SignedBy(0)}),
		Identities: []*msp.MSPPrincipal{principal},
	}

	return p
}


func signedByAnyOfGivenRole(role msp.MSPRole_MSPRoleType, ids []string) *cb.SignaturePolicyEnvelope {
	return SignedByNOutOfGivenRole(1, role, ids)
}

func SignedByNOutOfGivenRole(n int32, role msp.MSPRole_MSPRoleType, ids []string) *cb.SignaturePolicyEnvelope {
	
	
	sort.Strings(ids)
	principals := make([]*msp.MSPPrincipal, len(ids))
	sigspolicy := make([]*cb.SignaturePolicy, len(ids))
	for i, id := range ids {
		principals[i] = &msp.MSPPrincipal{
			PrincipalClassification: msp.MSPPrincipal_ROLE,
			Principal:               protoutil.MarshalOrPanic(&msp.MSPRole{Role: role, MspIdentifier: id})}
		sigspolicy[i] = SignedBy(int32(i))
	}

	
	p := &cb.SignaturePolicyEnvelope{
		Version:    0,
		Rule:       NOutOf(n, sigspolicy),
		Identities: principals,
	}

	return p
}




func SignedByAnyMember(ids []string) *cb.SignaturePolicyEnvelope {
	return signedByAnyOfGivenRole(msp.MSPRole_MEMBER, ids)
}




func SignedByAnyClient(ids []string) *cb.SignaturePolicyEnvelope {
	return signedByAnyOfGivenRole(msp.MSPRole_CLIENT, ids)
}




func SignedByAnyPeer(ids []string) *cb.SignaturePolicyEnvelope {
	return signedByAnyOfGivenRole(msp.MSPRole_PEER, ids)
}




func SignedByAnyAdmin(ids []string) *cb.SignaturePolicyEnvelope {
	return signedByAnyOfGivenRole(msp.MSPRole_ADMIN, ids)
}


func And(lhs, rhs *cb.SignaturePolicy) *cb.SignaturePolicy {
	return NOutOf(2, []*cb.SignaturePolicy{lhs, rhs})
}


func Or(lhs, rhs *cb.SignaturePolicy) *cb.SignaturePolicy {
	return NOutOf(1, []*cb.SignaturePolicy{lhs, rhs})
}


func NOutOf(n int32, policies []*cb.SignaturePolicy) *cb.SignaturePolicy {
	return &cb.SignaturePolicy{
		Type: &cb.SignaturePolicy_NOutOf_{
			NOutOf: &cb.SignaturePolicy_NOutOf{
				N:     n,
				Rules: policies,
			},
		},
	}
}
