/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chaincode

import (
	"context"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/internal/pkg/identity"
	cb "github.com/mcc-github/blockchain/protos/common"
	pb "github.com/mcc-github/blockchain/protos/peer"
	lb "github.com/mcc-github/blockchain/protos/peer/lifecycle"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)





type CommitSimulator struct {
	Command        *cobra.Command
	EndorserClient pb.EndorserClient
	Input          *CommitSimulationInput
	Signer         identity.SignerSerializer
	Writer         io.Writer
}





type CommitSimulationInput struct {
	ChannelID                string
	Name                     string
	Version                  string
	PackageID                string
	Sequence                 int64
	EndorsementPlugin        string
	ValidationPlugin         string
	ValidationParameterBytes []byte
	CollectionConfigPackage  *cb.CollectionConfigPackage
	InitRequired             bool
	PeerAddresses            []string
	TxID                     string
	OutputFormat             string
}


func (c *CommitSimulationInput) Validate() error {
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



func SimulateCommitCmd(c *CommitSimulator) *cobra.Command {
	chaincodeSimulateCommitCmd := &cobra.Command{
		Use:   "simulatecommit",
		Short: fmt.Sprintf("Simulate committing a chaincode definition."),
		Long:  fmt.Sprintf("Simulate committing a chaincode definition."),
		RunE: func(cmd *cobra.Command, args []string) error {
			if c == nil {
				
				input, err := c.createInput()
				if err != nil {
					return err
				}

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

				c = &CommitSimulator{
					Command:        cmd,
					Input:          input,
					EndorserClient: cc.EndorserClients[0],
					Signer:         cc.Signer,
					Writer:         os.Stdout,
				}
			}

			return c.Simulate()
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
		"output",
	}
	attachFlags(chaincodeSimulateCommitCmd, flagList)

	return chaincodeSimulateCommitCmd
}



func (c *CommitSimulator) Simulate() error {
	err := c.Input.Validate()
	if err != nil {
		return err
	}

	if c.Command != nil {
		
		c.Command.SilenceUsage = true
	}

	proposal, err := c.createProposal(c.Input.TxID)
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

	if strings.ToLower(c.Input.OutputFormat) == "json" {
		return printResponseAsJSON(proposalResponse, &lb.SimulateCommitChaincodeDefinitionResult{}, c.Writer)
	}
	return c.printResponse(proposalResponse)
}



func (c *CommitSimulator) printResponse(proposalResponse *pb.ProposalResponse) error {
	result := &lb.SimulateCommitChaincodeDefinitionResult{}
	err := proto.Unmarshal(proposalResponse.Response.Payload, result)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal proposal response's response payload")
	}

	orgs := []string{}
	for org := range result.Approved {
		orgs = append(orgs, org)
	}
	sort.Strings(orgs)

	fmt.Fprintf(c.Writer, "Chaincode definition for chaincode '%s', version '%s', sequence '%d' on channel '%s' approval status by org:\n", c.Input.Name, c.Input.Version, c.Input.Sequence, c.Input.ChannelID)
	for _, org := range orgs {
		fmt.Fprintf(c.Writer, "%s: %t\n", org, result.Approved[org])
	}

	return nil
}


func (c *CommitSimulator) createInput() (*CommitSimulationInput, error) {
	policyBytes, err := createPolicyBytes(signaturePolicy, channelConfigPolicy)
	if err != nil {
		return nil, err
	}

	ccp, err := createCollectionConfigPackage(collectionsConfigFile)
	if err != nil {
		return nil, err
	}

	input := &CommitSimulationInput{
		ChannelID:                channelID,
		Name:                     chaincodeName,
		Version:                  chaincodeVersion,
		PackageID:                packageID,
		Sequence:                 int64(sequence),
		EndorsementPlugin:        endorsementPlugin,
		ValidationPlugin:         validationPlugin,
		ValidationParameterBytes: policyBytes,
		InitRequired:             initRequired,
		CollectionConfigPackage:  ccp,
		PeerAddresses:            peerAddresses,
		OutputFormat:             output,
	}

	return input, nil
}

func (c *CommitSimulator) createProposal(inputTxID string) (*pb.Proposal, error) {
	args := &lb.SimulateCommitChaincodeDefinitionArgs{
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
		return nil, err
	}
	ccInput := &pb.ChaincodeInput{Args: [][]byte{[]byte(simulateCommitFuncName), argsBytes}}

	cis := &pb.ChaincodeInvocationSpec{
		ChaincodeSpec: &pb.ChaincodeSpec{
			ChaincodeId: &pb.ChaincodeID{Name: lifecycleName},
			Input:       ccInput,
		},
	}

	creatorBytes, err := c.Signer.Serialize()
	if err != nil {
		return nil, errors.WithMessage(err, "failed to serialize identity")
	}

	proposal, _, err := protoutil.CreateChaincodeProposalWithTxIDAndTransient(cb.HeaderType_ENDORSER_TRANSACTION, c.Input.ChannelID, cis, creatorBytes, inputTxID, nil)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to create ChaincodeInvocationSpec proposal")
	}

	return proposal, nil
}
