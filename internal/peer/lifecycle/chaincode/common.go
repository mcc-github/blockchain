/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chaincode

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/common/cauthdsl"
	"github.com/mcc-github/blockchain/internal/peer/chaincode"
	cb "github.com/mcc-github/blockchain/protos/common"
	pb "github.com/mcc-github/blockchain/protos/peer"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)



type EndorserClient interface {
	ProcessProposal(ctx context.Context, in *pb.SignedProposal, opts ...grpc.CallOption) (*pb.ProposalResponse, error)
}


type PeerDeliverClient interface {
	Deliver(ctx context.Context, opts ...grpc.CallOption) (pb.Deliver_DeliverClient, error)
	DeliverFiltered(ctx context.Context, opts ...grpc.CallOption) (pb.Deliver_DeliverClient, error)
}


type Signer interface {
	Sign(msg []byte) ([]byte, error)
	Serialize() ([]byte, error)
}

func signProposal(proposal *pb.Proposal, signer Signer) (*pb.SignedProposal, error) {
	
	if proposal == nil {
		return nil, errors.New("proposal cannot be nil")
	}

	if signer == nil {
		return nil, errors.New("signer cannot be nil")
	}

	proposalBytes, err := proto.Marshal(proposal)
	if err != nil {
		return nil, errors.Wrap(err, "error marshaling proposal")
	}

	signature, err := signer.Sign(proposalBytes)
	if err != nil {
		return nil, err
	}

	return &pb.SignedProposal{
		ProposalBytes: proposalBytes,
		Signature:     signature,
	}, nil
}

func createPolicyBytes(signaturePolicy, channelConfigPolicy string) ([]byte, error) {
	if signaturePolicy == "" && channelConfigPolicy == "" {
		
		return nil, nil
	}

	if signaturePolicy != "" && channelConfigPolicy != "" {
		
		return nil, errors.New("cannot specify both \"--signature-policy\" and \"--channel-config-policy\"")
	}

	var applicationPolicy *pb.ApplicationPolicy
	if signaturePolicy != "" {
		signaturePolicyEnvelope, err := cauthdsl.FromString(signaturePolicy)
		if err != nil {
			return nil, errors.Errorf("invalid signature policy: %s", signaturePolicy)
		}

		applicationPolicy = &pb.ApplicationPolicy{
			Type: &pb.ApplicationPolicy_SignaturePolicy{
				SignaturePolicy: signaturePolicyEnvelope,
			},
		}
	}

	if channelConfigPolicy != "" {
		applicationPolicy = &pb.ApplicationPolicy{
			Type: &pb.ApplicationPolicy_ChannelConfigPolicyReference{
				ChannelConfigPolicyReference: channelConfigPolicy,
			},
		}
	}

	policyBytes := protoutil.MarshalOrPanic(applicationPolicy)
	return policyBytes, nil
}

func createCollectionConfigPackage(collectionsConfigFile string) (*cb.CollectionConfigPackage, error) {
	var ccp *cb.CollectionConfigPackage
	if collectionsConfigFile != "" {
		var err error
		ccp, _, err = chaincode.GetCollectionConfigFromFile(collectionsConfigFile)
		if err != nil {
			return nil, errors.WithMessagef(err, "invalid collection configuration in file %s", collectionsConfigFile)
		}
	}
	return ccp, nil
}

func printResponseAsJSON(proposalResponse *pb.ProposalResponse, msg proto.Message, out io.Writer) error {
	err := proto.Unmarshal(proposalResponse.Response.Payload, msg)
	if err != nil {
		return errors.Wrapf(err, "failed to unmarshal proposal response's response payload as type %T", msg)
	}

	bytes, err := json.MarshalIndent(msg, "", "\t")
	if err != nil {
		return errors.Wrap(err, "failed to marshal output")
	}

	fmt.Fprintf(out, "%s\n", string(bytes))

	return nil
}
