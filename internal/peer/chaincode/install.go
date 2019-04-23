/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chaincode

import (
	"context"
	"io/ioutil"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/core/common/ccpackage"
	"github.com/mcc-github/blockchain/core/common/ccprovider"
	"github.com/mcc-github/blockchain/internal/peer/common"
	"github.com/mcc-github/blockchain/internal/pkg/identity"
	cb "github.com/mcc-github/blockchain/protos/common"
	pb "github.com/mcc-github/blockchain/protos/peer"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var chaincodeInstallCmd *cobra.Command

const (
	installCmdName = "install"
)



type Installer struct {
	Command         *cobra.Command
	EndorserClients []pb.EndorserClient
	Input           *InstallInput
	Signer          identity.SignerSerializer
}



type InstallInput struct {
	Name        string
	Version     string
	Language    string
	PackageFile string
	Path        string
}


func installCmd(cf *ChaincodeCmdFactory, i *Installer) *cobra.Command {
	chaincodeInstallCmd = &cobra.Command{
		Use:       "install",
		Short:     "Install a chaincode.",
		Long:      "Install a chaincode on a peer. This installs a chaincode deployment spec package (if provided) or packages the specified chaincode before subsequently installing it.",
		ValidArgs: []string{"1"},
		RunE: func(cmd *cobra.Command, args []string) error {
			if i == nil {
				var err error
				if cf == nil {
					cf, err = InitCmdFactory(cmd.Name(), true, false)
					if err != nil {
						return err
					}
				}
				i = &Installer{
					Command:         cmd,
					EndorserClients: cf.EndorserClients,
					Signer:          cf.Signer,
				}
			}
			return i.installChaincode(args)
		},
	}
	flagList := []string{
		"lang",
		"ctor",
		"path",
		"name",
		"version",
		"peerAddresses",
		"tlsRootCertFiles",
		"connectionProfile",
	}
	attachFlags(chaincodeInstallCmd, flagList)

	return chaincodeInstallCmd
}


func (i *Installer) installChaincode(args []string) error {
	if i.Command != nil {
		
		i.Command.SilenceUsage = true
	}

	i.setInput(args)

	
	return i.install()
}

func (i *Installer) setInput(args []string) {
	i.Input = &InstallInput{
		Name:    chaincodeName,
		Version: chaincodeVersion,
		Path:    chaincodePath,
	}

	if len(args) > 0 {
		i.Input.PackageFile = args[0]
	}
}



func (i *Installer) install() error {
	ccPkgMsg, err := i.getChaincodePackageMessage()
	if err != nil {
		return err
	}

	proposal, err := i.createInstallProposal(ccPkgMsg)
	if err != nil {
		return err
	}

	signedProposal, err := protoutil.GetSignedProposal(proposal, i.Signer)
	if err != nil {
		return errors.WithMessagef(err, "error creating signed proposal for %s", chainFuncName)
	}

	return i.submitInstallProposal(signedProposal)
}

func (i *Installer) submitInstallProposal(signedProposal *pb.SignedProposal) error {
	
	proposalResponse, err := i.EndorserClients[0].ProcessProposal(context.Background(), signedProposal)
	if err != nil {
		return errors.WithMessage(err, "error endorsing chaincode install")
	}

	if proposalResponse == nil {
		return errors.New("error during install: received nil proposal response")
	}

	if proposalResponse.Response == nil {
		return errors.New("error during install: received proposal response with nil response")
	}

	if proposalResponse.Response.Status != int32(cb.Status_SUCCESS) {
		return errors.Errorf("install failed with status: %d - %s", proposalResponse.Response.Status, proposalResponse.Response.Message)
	}
	logger.Infof("Installed remotely: %v", proposalResponse)

	return nil
}

func (i *Installer) getChaincodePackageMessage() (proto.Message, error) {
	
	if i.Input.PackageFile == "" {
		if i.Input.Path == common.UndefinedParamValue || i.Input.Version == common.UndefinedParamValue || i.Input.Name == common.UndefinedParamValue {
			return nil, errors.Errorf("must supply value for %s name, path and version parameters", chainFuncName)
		}
		
		ccPkgMsg, err := genChaincodeDeploymentSpec(i.Command, i.Input.Name, i.Input.Version)
		if err != nil {
			return nil, err
		}
		return ccPkgMsg, nil
	}

	
	
	
	ccPkgMsg, cds, err := getPackageFromFile(i.Input.PackageFile)
	if err != nil {
		return nil, err
	}

	
	cName := cds.ChaincodeSpec.ChaincodeId.Name
	cVersion := cds.ChaincodeSpec.ChaincodeId.Version

	
	if i.Input.Name != "" && i.Input.Name != cName {
		return nil, errors.Errorf("chaincode name %s does not match name %s in package", i.Input.Name, cName)
	}

	
	if i.Input.Version != "" && i.Input.Version != cVersion {
		return nil, errors.Errorf("chaincode version %s does not match version %s in packages", i.Input.Version, cVersion)
	}

	return ccPkgMsg, nil
}

func (i *Installer) createInstallProposal(msg proto.Message) (*pb.Proposal, error) {
	creator, err := i.Signer.Serialize()
	if err != nil {
		return nil, errors.WithMessage(err, "error serializing identity")
	}

	prop, _, err := protoutil.CreateInstallProposalFromCDS(msg, creator)
	if err != nil {
		return nil, errors.WithMessagef(err, "error creating proposal for %s", chainFuncName)
	}

	return prop, nil
}


func genChaincodeDeploymentSpec(cmd *cobra.Command, chaincodeName, chaincodeVersion string) (*pb.ChaincodeDeploymentSpec, error) {
	if existed, _ := ccprovider.ChaincodePackageExists(chaincodeName, chaincodeVersion); existed {
		return nil, errors.Errorf("chaincode %s:%s already exists", chaincodeName, chaincodeVersion)
	}

	spec, err := getChaincodeSpec(cmd)
	if err != nil {
		return nil, err
	}

	cds, err := getChaincodeDeploymentSpec(spec, true)
	if err != nil {
		return nil, errors.WithMessagef(err, "error getting chaincode deployment spec for %s", chaincodeName)
	}

	return cds, nil
}


func getPackageFromFile(ccPkgFile string) (proto.Message, *pb.ChaincodeDeploymentSpec, error) {
	ccPkgBytes, err := ioutil.ReadFile(ccPkgFile)
	if err != nil {
		return nil, nil, err
	}

	
	ccpack, err := ccprovider.GetCCPackage(ccPkgBytes)
	if err != nil {
		return nil, nil, err
	}

	
	o := ccpack.GetPackageObject()

	
	cds, ok := o.(*pb.ChaincodeDeploymentSpec)
	if !ok || cds == nil {
		
		env, ok := o.(*cb.Envelope)
		if !ok || env == nil {
			return nil, nil, errors.New("error extracting valid chaincode package")
		}

		
		_, sCDS, err := ccpackage.ExtractSignedCCDepSpec(env)
		if err != nil {
			return nil, nil, errors.WithMessage(err, "error extracting valid signed chaincode package")
		}

		
		cds, err = protoutil.GetChaincodeDeploymentSpec(sCDS.ChaincodeDeploymentSpec)
		if err != nil {
			return nil, nil, errors.WithMessage(err, "error extracting chaincode deployment spec")
		}

		err = platformRegistry.ValidateDeploymentSpec(cds.ChaincodeSpec.Type.String(), cds.CodePackage)
		if err != nil {
			return nil, nil, errors.WithMessage(err, "chaincode deployment spec validation failed")
		}
	}

	return o, cds, nil
}
