/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chaincode

import (
	"context"
	"path/filepath"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/core/chaincode/persistence"
	cb "github.com/mcc-github/blockchain/protos/common"
	pb "github.com/mcc-github/blockchain/protos/peer"
	lb "github.com/mcc-github/blockchain/protos/peer/lifecycle"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)



type InstalledPackageGetter struct {
	Command        *cobra.Command
	Input          *GetInstalledPackageInput
	EndorserClient EndorserClient
	Signer         Signer
	Writer         Writer
}



type GetInstalledPackageInput struct {
	PackageID       string
	OutputDirectory string
}


func (i *GetInstalledPackageInput) Validate() error {
	if i.PackageID == "" {
		return errors.New("The required parameter 'package-id' is empty. Rerun the command with --package-id flag")
	}

	return nil
}



func GetInstalledPackageCmd(i *InstalledPackageGetter) *cobra.Command {
	chaincodeGetInstalledPackageCmd := &cobra.Command{
		Use:   "getinstalledpackage [outputfile]",
		Short: "Get an installed chaincode package from a peer.",
		Long:  "Get an installed chaincode package from a peer.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if i == nil {
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

				gipInput := &GetInstalledPackageInput{
					PackageID:       packageID,
					OutputDirectory: outputDirectory,
				}

				
				
				
				i = &InstalledPackageGetter{
					Command:        cmd,
					EndorserClient: cc.EndorserClients[0],
					Input:          gipInput,
					Signer:         cc.Signer,
					Writer:         &persistence.FilesystemIO{},
				}
			}
			return i.Get()
		},
	}

	flagList := []string{
		"peerAddresses",
		"tlsRootCertFiles",
		"connectionProfile",
		"package-id",
		"output-directory",
	}
	attachFlags(chaincodeGetInstalledPackageCmd, flagList)

	return chaincodeGetInstalledPackageCmd
}


func (i *InstalledPackageGetter) Get() error {
	if i.Command != nil {
		
		i.Command.SilenceUsage = true
	}

	if err := i.Input.Validate(); err != nil {
		return err
	}

	proposal, err := i.createProposal()
	if err != nil {
		return errors.WithMessage(err, "failed to create proposal")
	}

	signedProposal, err := signProposal(proposal, i.Signer)
	if err != nil {
		return errors.WithMessage(err, "failed to create signed proposal")
	}

	proposalResponse, err := i.EndorserClient.ProcessProposal(context.Background(), signedProposal)
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
		return errors.Errorf("proposal failed with status: %d - %s", proposalResponse.Response.Status, proposalResponse.Response.Message)
	}

	return i.writePackage(proposalResponse)
}

func (i *InstalledPackageGetter) writePackage(proposalResponse *pb.ProposalResponse) error {
	result := &lb.GetInstalledChaincodePackageResult{}
	err := proto.Unmarshal(proposalResponse.Response.Payload, result)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal proposal response's response payload")
	}

	outputFile := filepath.Join(i.Input.OutputDirectory, i.Input.PackageID+".tar.gz")

	dir, name := filepath.Split(outputFile)
	
	if dir, err = filepath.Abs(dir); err != nil {
		return err
	}

	err = i.Writer.WriteFile(dir, name, result.ChaincodeInstallPackage)
	if err != nil {
		err = errors.Wrapf(err, "failed to write chaincode package to %s", outputFile)
		logger.Error(err.Error())
		return err
	}

	return nil
}

func (i *InstalledPackageGetter) createProposal() (*pb.Proposal, error) {
	args := &lb.GetInstalledChaincodePackageArgs{
		PackageId: i.Input.PackageID,
	}

	argsBytes, err := proto.Marshal(args)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal args")
	}

	ccInput := &pb.ChaincodeInput{
		Args: [][]byte{[]byte("GetInstalledChaincodePackage"), argsBytes},
	}

	cis := &pb.ChaincodeInvocationSpec{
		ChaincodeSpec: &pb.ChaincodeSpec{
			ChaincodeId: &pb.ChaincodeID{Name: lifecycleName},
			Input:       ccInput,
		},
	}

	signerSerialized, err := i.Signer.Serialize()
	if err != nil {
		return nil, errors.WithMessage(err, "failed to serialize identity")
	}

	proposal, _, err := protoutil.CreateProposalFromCIS(cb.HeaderType_ENDORSER_TRANSACTION, "", cis, signerSerialized)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to create ChaincodeInvocationSpec proposal")
	}

	return proposal, nil
}
