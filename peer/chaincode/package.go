/*
Copyright IBM Corp. 2017 All Rights Reserved.

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

package chaincode

import (
	"fmt"
	"io/ioutil"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/common/cauthdsl"
	"github.com/mcc-github/blockchain/core/common/ccpackage"
	"github.com/mcc-github/blockchain/msp"
	mspmgmt "github.com/mcc-github/blockchain/msp/mgmt"
	pcommon "github.com/mcc-github/blockchain/protos/common"
	pb "github.com/mcc-github/blockchain/protos/peer"
	"github.com/mcc-github/blockchain/protos/utils"
	"github.com/spf13/cobra"
)

var chaincodePackageCmd *cobra.Command
var createSignedCCDepSpec bool
var signCCDepSpec bool
var instantiationPolicy string

const packageCmdName = "package"
const packageDesc = "Package the specified chaincode into a deployment spec."

type ccDepSpecFactory func(spec *pb.ChaincodeSpec) (*pb.ChaincodeDeploymentSpec, error)

func defaultCDSFactory(spec *pb.ChaincodeSpec) (*pb.ChaincodeDeploymentSpec, error) {
	return getChaincodeDeploymentSpec(spec, true)
}


func packageCmd(cf *ChaincodeCmdFactory, cdsFact ccDepSpecFactory) *cobra.Command {
	chaincodePackageCmd = &cobra.Command{
		Use:       "package",
		Short:     packageDesc,
		Long:      packageDesc,
		ValidArgs: []string{"1"},
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("output file not specified or invalid number of args (filename should be the only arg)")
			}
			
			if cdsFact == nil {
				cdsFact = defaultCDSFactory
			}
			return chaincodePackage(cmd, args, cdsFact, cf)
		},
	}
	flagList := []string{
		"lang",
		"ctor",
		"path",
		"name",
		"version",
	}
	attachFlags(chaincodePackageCmd, flagList)

	chaincodePackageCmd.Flags().BoolVarP(&createSignedCCDepSpec, "cc-package", "s", false, "create CC deployment spec for owner endorsements instead of raw CC deployment spec")
	chaincodePackageCmd.Flags().BoolVarP(&signCCDepSpec, "sign", "S", false, "if creating CC deployment spec package for owner endorsements, also sign it with local MSP")
	chaincodePackageCmd.Flags().StringVarP(&instantiationPolicy, "instantiate-policy", "i", "", "instantiation policy for the chaincode")

	return chaincodePackageCmd
}

func getInstantiationPolicy(policy string) (*pcommon.SignaturePolicyEnvelope, error) {
	p, err := cauthdsl.FromString(policy)
	if err != nil {
		return nil, fmt.Errorf("Invalid policy %s, err %s", policy, err)
	}
	return p, nil
}



func getChaincodeInstallPackage(cds *pb.ChaincodeDeploymentSpec, cf *ChaincodeCmdFactory) ([]byte, error) {
	
	var objToWrite proto.Message

	
	objToWrite = cds

	var err error

	var owner msp.SigningIdentity

	
	if createSignedCCDepSpec {
		
		
		
		if signCCDepSpec {
			if cf.Signer == nil {
				return nil, fmt.Errorf("Error getting signer")
			}
			owner = cf.Signer
		}
	}

	ip := instantiationPolicy
	if ip == "" {
		
		
		mspid, err := mspmgmt.GetLocalMSP().GetIdentifier()
		if err != nil {
			return nil, err
		}
		ip = "AND('" + mspid + ".admin')"
	}

	sp, err := getInstantiationPolicy(ip)
	if err != nil {
		return nil, err
	}

	
	objToWrite, err = ccpackage.OwnerCreateSignedCCDepSpec(cds, sp, owner)
	if err != nil {
		return nil, err
	}

	
	bytesToWrite, err := proto.Marshal(objToWrite)
	if err != nil {
		return nil, fmt.Errorf("Error marshalling chaincode package : %s", err)
	}

	return bytesToWrite, nil
}




func chaincodePackage(cmd *cobra.Command, args []string, cdsFact ccDepSpecFactory, cf *ChaincodeCmdFactory) error {
	if cdsFact == nil {
		return fmt.Errorf("Error chaincode deployment spec factory not specified")
	}
	
	cmd.SilenceUsage = true

	var err error
	if cf == nil {
		cf, err = InitCmdFactory(cmd.Name(), false, false)
		if err != nil {
			return err
		}
	}
	spec, err := getChaincodeSpec(cmd)
	if err != nil {
		return err
	}

	cds, err := cdsFact(spec)
	if err != nil {
		return fmt.Errorf("error getting chaincode code %s: %s", chaincodeName, err)
	}

	var bytesToWrite []byte
	if createSignedCCDepSpec {
		bytesToWrite, err = getChaincodeInstallPackage(cds, cf)
		if err != nil {
			return err
		}
	} else {
		bytesToWrite = utils.MarshalOrPanic(cds)
	}

	logger.Debugf("Packaged chaincode into deployment spec of size <%d>, with args = %v", len(bytesToWrite), args)
	fileToWrite := args[0]
	err = ioutil.WriteFile(fileToWrite, bytesToWrite, 0700)
	if err != nil {
		logger.Errorf("failed writing deployment spec to file [%s]: [%s]", fileToWrite, err)
		return err
	}

	return err
}
