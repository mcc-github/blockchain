/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chaincode

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/internal/peer/chaincode"
	"github.com/mcc-github/blockchain/internal/peer/common"
	cb "github.com/mcc-github/blockchain/protos/common"
	pb "github.com/mcc-github/blockchain/protos/peer"
	lb "github.com/mcc-github/blockchain/protos/peer/lifecycle"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)



type Committer struct {
	Certificate     tls.Certificate
	Command         *cobra.Command
	BroadcastClient common.BroadcastClient
	EndorserClients []EndorserClient
	DeliverClients  []pb.DeliverClient
	Input           *CommitInput
	Signer          Signer
}





type CommitInput struct {
	ChannelID                string
	Name                     string
	Version                  string
	Hash                     []byte
	Sequence                 int64
	EndorsementPlugin        string
	ValidationPlugin         string
	ValidationParameterBytes []byte
	CollectionConfigPackage  *cb.CollectionConfigPackage
	InitRequired             bool
	PeerAddresses            []string
	WaitForEvent             bool
	WaitForEventTimeout      time.Duration
	TxID                     string
}


func (c *CommitInput) Validate() error {
	if c.ChannelID == "" {
		return errors.New("The required parameter 'channelID' is empty. Rerun the command with -C flag")
	}

	if c.Name == "" {
		return errors.New("The required parameter 'name' is empty. Rerun the command with -n flag")
	}

	if c.Version == "" {
		return errors.New("The required parameter 'version' is empty. Rerun the command with -v flag")
	}

	if c.Sequence == 0 {
		return errors.New("The required parameter 'sequence' is empty. Rerun the command with --sequence flag")
	}

	return nil
}


func CommitCmd(c *Committer) *cobra.Command {
	chaincodeCommitCmd := &cobra.Command{
		Use:   "commit",
		Short: fmt.Sprintf("Commit the chaincode definition on the channel."),
		Long:  fmt.Sprintf("Commit the chaincode definition on the channel."),
		RunE: func(cmd *cobra.Command, args []string) error {
			if c == nil {
				
				input, err := c.createInput()
				if err != nil {
					return err
				}

				ccInput := &ClientConnectionsInput{
					CommandName:           cmd.Name(),
					EndorserRequired:      true,
					OrdererRequired:       true,
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

				endorserClients := make([]EndorserClient, len(cc.EndorserClients))
				for i, e := range cc.EndorserClients {
					endorserClients[i] = e
				}

				c = &Committer{
					Command:         cmd,
					Input:           input,
					Certificate:     cc.Certificate,
					BroadcastClient: cc.BroadcastClient,
					DeliverClients:  cc.DeliverClients,
					EndorserClients: endorserClients,
					Signer:          cc.Signer,
				}
			}
			return c.Commit()
		},
	}
	flagList := []string{
		"channelID",
		"name",
		"version",
		"sequence",
		"endorsement-plugin",
		"validation-plugin",
		"signature-policy",
		"channel-config-policy",
		"init-required",
		"collections-config",
		"peerAddresses",
		"tlsRootCertFiles",
		"connectionProfile",
		"waitForEvent",
		"waitForEventTimeout",
	}
	attachFlags(chaincodeCommitCmd, flagList)

	return chaincodeCommitCmd
}


func (c *Committer) Commit() error {
	err := c.Input.Validate()
	if err != nil {
		return err
	}

	if c.Command != nil {
		
		c.Command.SilenceUsage = true
	}

	proposal, txID, err := c.createProposal(c.Input.TxID)
	if err != nil {
		return errors.WithMessage(err, "failed to create proposal")
	}

	signedProposal, err := signProposal(proposal, c.Signer)
	if err != nil {
		return errors.WithMessage(err, "failed to create signed proposal")
	}

	var responses []*pb.ProposalResponse
	for _, endorser := range c.EndorserClients {
		proposalResponse, err := endorser.ProcessProposal(context.Background(), signedProposal)
		if err != nil {
			return errors.WithMessage(err, "failed to endorse proposal")
		}
		responses = append(responses, proposalResponse)
	}

	if len(responses) == 0 {
		
		return errors.New("no proposal responses received")
	}

	
	
	proposalResponse := responses[0]

	if proposalResponse == nil {
		return errors.New("received nil proposal response")
	}

	if proposalResponse.Response == nil {
		return errors.New("received proposal response with nil response")
	}

	if proposalResponse.Response.Status != int32(cb.Status_SUCCESS) {
		return errors.Errorf("proposal failed with status: %d - %s", proposalResponse.Response.Status, proposalResponse.Response.Message)
	}
	
	env, err := protoutil.CreateSignedTx(proposal, c.Signer, responses...)
	if err != nil {
		return errors.WithMessage(err, "failed to create signed transaction")
	}

	var dg *chaincode.DeliverGroup
	var ctx context.Context
	if c.Input.WaitForEvent {
		var cancelFunc context.CancelFunc
		ctx, cancelFunc = context.WithTimeout(context.Background(), c.Input.WaitForEventTimeout)
		defer cancelFunc()

		dg = chaincode.NewDeliverGroup(
			c.DeliverClients,
			c.Input.PeerAddresses,
			c.Signer,
			c.Certificate,
			c.Input.ChannelID,
			txID,
		)
		
		err := dg.Connect(ctx)
		if err != nil {
			return err
		}
	}

	if err = c.BroadcastClient.Send(env); err != nil {
		return errors.WithMessage(err, "failed to send transaction")
	}

	if dg != nil && ctx != nil {
		
		err = dg.Wait(ctx)
		if err != nil {
			return err
		}
	}
	return err
}


func (c *Committer) createInput() (*CommitInput, error) {
	policyBytes, err := createPolicyBytes(signaturePolicy, channelConfigPolicy)
	if err != nil {
		return nil, err
	}

	ccp, err := createCollectionConfigPackage(collectionsConfigFile)
	if err != nil {
		return nil, err
	}

	input := &CommitInput{
		ChannelID:                channelID,
		Name:                     chaincodeName,
		Version:                  chaincodeVersion,
		Sequence:                 int64(sequence),
		EndorsementPlugin:        endorsementPlugin,
		ValidationPlugin:         validationPlugin,
		ValidationParameterBytes: policyBytes,
		InitRequired:             initRequired,
		CollectionConfigPackage:  ccp,
		PeerAddresses:            peerAddresses,
		WaitForEvent:             waitForEvent,
		WaitForEventTimeout:      waitForEventTimeout,
	}

	return input, nil
}

func (c *Committer) createProposal(inputTxID string) (proposal *pb.Proposal, txID string, err error) {
	args := &lb.CommitChaincodeDefinitionArgs{
		Name:                c.Input.Name,
		Version:             c.Input.Version,
		Sequence:            c.Input.Sequence,
		EndorsementPlugin:   c.Input.EndorsementPlugin,
		ValidationPlugin:    c.Input.ValidationPlugin,
		ValidationParameter: c.Input.ValidationParameterBytes,
		InitRequired:        c.Input.InitRequired,
		Collections:         c.Input.CollectionConfigPackage,
	}

	argsBytes, err := proto.Marshal(args)
	if err != nil {
		return nil, "", err
	}
	ccInput := &pb.ChaincodeInput{Args: [][]byte{[]byte(commitFuncName), argsBytes}}

	cis := &pb.ChaincodeInvocationSpec{
		ChaincodeSpec: &pb.ChaincodeSpec{
			ChaincodeId: &pb.ChaincodeID{Name: lifecycleName},
			Input:       ccInput,
		},
	}

	creatorBytes, err := c.Signer.Serialize()
	if err != nil {
		return nil, "", errors.WithMessage(err, "failed to serialize identity")
	}

	proposal, txID, err = protoutil.CreateChaincodeProposalWithTxIDAndTransient(cb.HeaderType_ENDORSER_TRANSACTION, c.Input.ChannelID, cis, creatorBytes, inputTxID, nil)
	if err != nil {
		return nil, "", errors.WithMessage(err, "failed to create ChaincodeInvocationSpec proposal")
	}

	return proposal, txID, nil
}
