/*
Copyright IBM Corp. 2016 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chaincode

import (
	"fmt"
	"io/ioutil"

	"github.com/mcc-github/blockchain/core/common/ccpackage"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/spf13/cobra"
)


func signpackageCmd(cf *ChaincodeCmdFactory) *cobra.Command {
	spCmd := &cobra.Command{
		Use:       "signpackage",
		Short:     "Sign the specified chaincode package",
		Long:      "Sign the specified chaincode package",
		ValidArgs: []string{"2"},
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return fmt.Errorf("peer chaincode signpackage <inputpackage> <outputpackage>")
			}
			return signpackage(cmd, args[0], args[1], cf)
		},
	}

	return spCmd
}

func signpackage(cmd *cobra.Command, ipackageFile string, opackageFile string, cf *ChaincodeCmdFactory) error {
	
	cmd.SilenceUsage = true

	var err error
	if cf == nil {
		cf, err = InitCmdFactory(cmd.Name(), false, false)
		if err != nil {
			return err
		}
	}

	b, err := ioutil.ReadFile(ipackageFile)
	if err != nil {
		return err
	}

	env := protoutil.UnmarshalEnvelopeOrPanic(b)

	env, err = ccpackage.SignExistingPackage(env, cf.Signer)
	if err != nil {
		return err
	}

	b = protoutil.MarshalOrPanic(env)
	err = ioutil.WriteFile(opackageFile, b, 0700)
	if err != nil {
		return err
	}

	fmt.Printf("Wrote signed package to %s successfully\n", opackageFile)

	return nil
}