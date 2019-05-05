/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chaincode

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"reflect"
	"strings"

	"github.com/golang/protobuf/proto"
	cb "github.com/mcc-github/blockchain/protos/common"
	pb "github.com/mcc-github/blockchain/protos/peer"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	chaincodeListCmd          *cobra.Command
	getInstalledChaincodes    bool
	getInstantiatedChaincodes bool
)



func listCmd(cf *ChaincodeCmdFactory) *cobra.Command {
	chaincodeListCmd = &cobra.Command{
		Use:   "list",
		Short: "Get the instantiated chaincodes on a channel or installed chaincodes on a peer.",
		Long:  "Get the instantiated chaincodes in the channel if specify channel, or get installed chaincodes on the peer",
		RunE: func(cmd *cobra.Command, args []string) error {
			return getChaincodes(cmd, cf)
		},
	}

	flagList := []string{
		"channelID",
		"installed",
		"instantiated",
		"peerAddresses",
		"tlsRootCertFiles",
		"connectionProfile",
	}
	attachFlags(chaincodeListCmd, flagList)

	return chaincodeListCmd
}

func getChaincodes(cmd *cobra.Command, cf *ChaincodeCmdFactory) error {
	if getInstantiatedChaincodes && channelID == "" {
		return errors.New("The required parameter 'channelID' is empty. Rerun the command with -C flag")
	}
	
	cmd.SilenceUsage = true

	var err error
	if cf == nil {
		cf, err = InitCmdFactory(cmd.Name(), true, false)
		if err != nil {
			return err
		}
	}

	creator, err := cf.Signer.Serialize()
	if err != nil {
		return fmt.Errorf("error serializing identity: %s", err)
	}

	var proposal *pb.Proposal

	if getInstalledChaincodes == getInstantiatedChaincodes {
		return errors.New("must explicitly specify \"--installed\" or \"--instantiated\"")
	}

	if getInstalledChaincodes {
		proposal, _, err = protoutil.CreateGetInstalledChaincodesProposal(creator)
	}

	if getInstantiatedChaincodes {
		proposal, _, err = protoutil.CreateGetChaincodesProposal(channelID, creator)
	}

	if err != nil {
		return errors.WithMessage(err, "error creating proposal")
	}

	var signedProposal *pb.SignedProposal
	signedProposal, err = protoutil.GetSignedProposal(proposal, cf.Signer)
	if err != nil {
		return errors.WithMessage(err, "error creating signed proposal")
	}

	
	proposalResponse, err := cf.EndorserClients[0].ProcessProposal(context.Background(), signedProposal)
	if err != nil {
		return errors.WithMessage(err, "error endorsing proposal")
	}

	if proposalResponse.Response == nil {
		return errors.Errorf("proposal response had nil response")
	}

	if proposalResponse.Response.Status != int32(cb.Status_SUCCESS) {
		return errors.Errorf("bad response: %d - %s", proposalResponse.Response.Status, proposalResponse.Response.Message)
	}

	return printResponse(getInstalledChaincodes, getInstantiatedChaincodes, proposalResponse)
}





func printResponse(getInstalledChaincodes, getInstantiatedChaincodes bool, proposalResponse *pb.ProposalResponse) error {
	cqr := &pb.ChaincodeQueryResponse{}
	err := proto.Unmarshal(proposalResponse.Response.Payload, cqr)
	if err != nil {
		return err
	}

	if getInstalledChaincodes {
		fmt.Println("Get installed chaincodes on peer:")
	} else {
		fmt.Printf("Get instantiated chaincodes on channel %s:\n", channelID)
	}

	for _, chaincode := range cqr.Chaincodes {
		fmt.Printf("%v\n", ccInfo{chaincode}.String())
	}

	return nil
}

type ccInfo struct {
	*pb.ChaincodeInfo
}

func (cci ccInfo) String() string {
	b := bytes.Buffer{}
	md := reflect.ValueOf(*cci.ChaincodeInfo)
	md2 := reflect.Indirect(reflect.ValueOf(*cci.ChaincodeInfo)).Type()
	for i := 0; i < md.NumField(); i++ {
		f := md.Field(i)
		val := f.String()
		if isBytes(f) {
			val = hex.EncodeToString(f.Bytes())
		}
		if len(val) == 0 {
			continue
		}
		
		if strings.HasPrefix(md2.Field(i).Name, "XXX") {
			continue
		}
		b.WriteString(fmt.Sprintf("%s: %s, ", md2.Field(i).Name, val))
	}
	return b.String()[:len(b.String())-2]

}

func isBytes(v reflect.Value) bool {
	return v.Kind() == reflect.Slice && v.Type().Elem().Kind() == reflect.Uint8
}