/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chaincode

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/proto"
	cb "github.com/mcc-github/blockchain/protos/common"
	pb "github.com/mcc-github/blockchain/protos/peer"
	lb "github.com/mcc-github/blockchain/protos/peer/lifecycle"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)



type CommittedQuerier struct {
	Command        *cobra.Command
	Input          *CommittedQueryInput
	EndorserClient EndorserClient
	Signer         Signer
}

type CommittedQueryInput struct {
	ChannelID string
	Name      string
}




func QueryCommittedCmd(c *CommittedQuerier) *cobra.Command {
	chaincodeQueryCommittedCmd := &cobra.Command{
		Use:   "querycommitted",
		Short: "Query a committed chaincode definition by channel and name on a peer.",
		Long:  "Query a committed chaincode definition by channel and name on a peer.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if c == nil {
				ccInput := &ClientConnectionsInput{
					CommandName:           cmd.Name(),
					EndorserRequired:      true,
					ChannelID:             channelID,
					PeerAddresses:         peerAddresses,
					TLSRootCertFiles:      tlsRootCertFiles,
					ConnectionProfilePath: connectionProfilePath,
					TLSEnabled:            viper.GetBool("peer.tls.enabled"),
				}

				cc, err := NewClientConnections(ccInput)
				if err != nil {
					return err
				}

				cqInput := &CommittedQueryInput{
					ChannelID: channelID,
					Name:      chaincodeName,
				}

				c = &CommittedQuerier{
					Command:        cmd,
					EndorserClient: cc.EndorserClients[0],
					Input:          cqInput,
					Signer:         cc.Signer,
				}
			}
			return c.Query()
		},
	}

	flagList := []string{
		"channelID",
		"name",
		"peerAddresses",
		"tlsRootCertFiles",
		"connectionProfile",
	}
	attachFlags(chaincodeQueryCommittedCmd, flagList)

	return chaincodeQueryCommittedCmd
}



func (c *CommittedQuerier) Query() error {
	if c.Command != nil {
		
		c.Command.SilenceUsage = true
	}

	err := c.validateInput()
	if err != nil {
		return err
	}

	proposal, err := c.createProposal()
	if err != nil {
		return errors.WithMessage(err, "failed to create proposal")
	}

	signedProposal, err := signProposal(proposal, c.Signer)
	if err != nil {
		return errors.WithMessage(err, "failed to create signed proposal")
	}

	proposalResponse, err := c.EndorserClient.ProcessProposal(context.Background(), signedProposal)
	if err != nil {
		return errors.WithMessage(err, "failed to endorse proposal")
	}

	if proposalResponse == nil {
		return errors.New("received nil proposal response")
	}

	if proposalResponse.Response == nil {
		return errors.New("received proposal response with nil response")
	}

	if proposalResponse.Response.Status != int32(cb.Status_SUCCESS) {
		return errors.Errorf("query failed with status: %d - %s", proposalResponse.Response.Status, proposalResponse.Response.Message)
	}

	return c.printResponse(proposalResponse)
}



func (c *CommittedQuerier) printResponse(proposalResponse *pb.ProposalResponse) error {
	qdcr := &lb.QueryChaincodeDefinitionResult{}
	err := proto.Unmarshal(proposalResponse.Response.Payload, qdcr)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal proposal response's response payload")
	}
	fmt.Printf("Committed chaincode definition for chaincode '%s' on channel '%s':\n", chaincodeName, channelID)
	fmt.Printf("Version: %s, Sequence: %d, Endorsement Plugin: %s, Validation Plugin: %s\n", qdcr.Version, qdcr.Sequence, qdcr.EndorsementPlugin, qdcr.ValidationPlugin)

	return nil

}

func (c *CommittedQuerier) validateInput() error {
	if c.Input.ChannelID == "" {
		return errors.New("channel name must be specified")
	}

	if c.Input.Name == "" {
		return errors.New("chaincode name must be specified")
	}

	return nil
}

func (c *CommittedQuerier) createProposal() (*pb.Proposal, error) {
	args := &lb.QueryChaincodeDefinitionArgs{
		Name: c.Input.Name,
	}

	argsBytes, err := proto.Marshal(args)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal args")
	}
	ccInput := &pb.ChaincodeInput{Args: [][]byte{[]byte("QueryChaincodeDefinition"), argsBytes}}

	cis := &pb.ChaincodeInvocationSpec{
		ChaincodeSpec: &pb.ChaincodeSpec{
			ChaincodeId: &pb.ChaincodeID{Name: lifecycleName},
			Input:       ccInput,
		},
	}

	signerSerialized, err := c.Signer.Serialize()
	if err != nil {
		return nil, errors.WithMessage(err, "failed to serialize identity")
	}

	proposal, _, err := protoutil.CreateProposalFromCIS(cb.HeaderType_ENDORSER_TRANSACTION, c.Input.ChannelID, cis, signerSerialized)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to create ChaincodeInvocationSpec proposal")
	}

	return proposal, nil
}
