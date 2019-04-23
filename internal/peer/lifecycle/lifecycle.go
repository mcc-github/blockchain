/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package lifecycle

import (
	"github.com/mcc-github/blockchain/internal/peer/lifecycle/chaincode"
	"github.com/spf13/cobra"
)


func Cmd() *cobra.Command {
	lifecycleCmd := &cobra.Command{
		Use:   "lifecycle",
		Short: "Perform _lifecycle operations",
		Long:  "Perform _lifecycle operations",
	}
	lifecycleCmd.AddCommand(chaincode.Cmd())

	return lifecycleCmd
}
