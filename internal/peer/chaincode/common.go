/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chaincode

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"strings"
	"sync"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/common/cauthdsl"
	"github.com/mcc-github/blockchain/common/util"
	"github.com/mcc-github/blockchain/core/chaincode/shim"
	"github.com/mcc-github/blockchain/internal/peer/common"
	"github.com/mcc-github/blockchain/internal/pkg/identity"
	pcommon "github.com/mcc-github/blockchain/protos/common"
	ab "github.com/mcc-github/blockchain/protos/orderer"
	pb "github.com/mcc-github/blockchain/protos/peer"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)


func checkSpec(spec *pb.ChaincodeSpec) error {
	
	if spec == nil {
		return errors.New("expected chaincode specification, nil received")
	}
	if spec.ChaincodeId == nil {
		return errors.New("expected chaincode ID, nil received")
	}

	return platformRegistry.ValidateSpec(spec.Type.String(), spec.ChaincodeId.Path)
}


func getChaincodeDeploymentSpec(spec *pb.ChaincodeSpec, crtPkg bool) (*pb.ChaincodeDeploymentSpec, error) {
	var codePackageBytes []byte
	if crtPkg {
		var err error
		if err = checkSpec(spec); err != nil {
			return nil, err
		}

		codePackageBytes, err = platformRegistry.GetDeploymentPayload(spec.Type.String(), spec.ChaincodeId.Path)
		if err != nil {
			err = errors.WithMessage(err, "error getting chaincode package bytes")
			return nil, err
		}
	}
	chaincodeDeploymentSpec := &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec, CodePackage: codePackageBytes}
	return chaincodeDeploymentSpec, nil
}


func getChaincodeSpec(cmd *cobra.Command) (*pb.ChaincodeSpec, error) {
	spec := &pb.ChaincodeSpec{}
	if err := checkChaincodeCmdParams(cmd); err != nil {
		
		cmd.SilenceUsage = false
		return spec, err
	}

	
	input := chaincodeInput{}
	if err := json.Unmarshal([]byte(chaincodeCtorJSON), &input); err != nil {
		return spec, errors.Wrap(err, "chaincode argument error")
	}
	input.IsInit = isInit

	chaincodeLang = strings.ToUpper(chaincodeLang)
	spec = &pb.ChaincodeSpec{
		Type:        pb.ChaincodeSpec_Type(pb.ChaincodeSpec_Type_value[chaincodeLang]),
		ChaincodeId: &pb.ChaincodeID{Path: chaincodePath, Name: chaincodeName, Version: chaincodeVersion},
		Input:       &input.ChaincodeInput,
	}
	return spec, nil
}



type chaincodeInput struct {
	pb.ChaincodeInput
}



func (c *chaincodeInput) UnmarshalJSON(b []byte) error {
	sa := struct {
		Function string
		Args     []string
	}{}
	err := json.Unmarshal(b, &sa)
	if err != nil {
		return err
	}
	allArgs := sa.Args
	if sa.Function != "" {
		allArgs = append([]string{sa.Function}, sa.Args...)
	}
	c.Args = util.ToChaincodeArgs(allArgs...)
	return nil
}

func chaincodeInvokeOrQuery(cmd *cobra.Command, invoke bool, cf *ChaincodeCmdFactory) (err error) {
	spec, err := getChaincodeSpec(cmd)
	if err != nil {
		return err
	}

	
	
	txID := ""

	proposalResp, err := ChaincodeInvokeOrQuery(
		spec,
		channelID,
		txID,
		invoke,
		cf.Signer,
		cf.Certificate,
		cf.EndorserClients,
		cf.DeliverClients,
		cf.BroadcastClient,
	)

	if err != nil {
		return errors.Errorf("%s - proposal response: %v", err, proposalResp)
	}

	if invoke {
		logger.Debugf("ESCC invoke result: %v", proposalResp)
		pRespPayload, err := protoutil.GetProposalResponsePayload(proposalResp.Payload)
		if err != nil {
			return errors.WithMessage(err, "error while unmarshaling proposal response payload")
		}
		ca, err := protoutil.GetChaincodeAction(pRespPayload.Extension)
		if err != nil {
			return errors.WithMessage(err, "error while unmarshaling chaincode action")
		}
		if proposalResp.Endorsement == nil {
			return errors.Errorf("endorsement failure during invoke. response: %v", proposalResp.Response)
		}
		logger.Infof("Chaincode invoke successful. result: %v", ca.Response)
	} else {
		if proposalResp == nil {
			return errors.New("error during query: received nil proposal response")
		}
		if proposalResp.Endorsement == nil {
			return errors.Errorf("endorsement failure during query. response: %v", proposalResp.Response)
		}

		if chaincodeQueryRaw && chaincodeQueryHex {
			return fmt.Errorf("options --raw (-r) and --hex (-x) are not compatible")
		}
		if chaincodeQueryRaw {
			fmt.Println(proposalResp.Response.Payload)
			return nil
		}
		if chaincodeQueryHex {
			fmt.Printf("%x\n", proposalResp.Response.Payload)
			return nil
		}
		fmt.Println(string(proposalResp.Response.Payload))
	}
	return nil
}

type collectionConfigJson struct {
	Name            string `json:"name"`
	Policy          string `json:"policy"`
	RequiredCount   int32  `json:"requiredPeerCount"`
	MaxPeerCount    int32  `json:"maxPeerCount"`
	BlockToLive     uint64 `json:"blockToLive"`
	MemberOnlyRead  bool   `json:"memberOnlyRead"`
	MemberOnlyWrite bool   `json:"memberOnlyWrite"`
}




func GetCollectionConfigFromFile(ccFile string) (*pcommon.CollectionConfigPackage, []byte, error) {
	fileBytes, err := ioutil.ReadFile(ccFile)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "could not read file '%s'", ccFile)
	}

	return getCollectionConfigFromBytes(fileBytes)
}




func getCollectionConfigFromBytes(cconfBytes []byte) (*pcommon.CollectionConfigPackage, []byte, error) {
	cconf := &[]collectionConfigJson{}
	err := json.Unmarshal(cconfBytes, cconf)
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not parse the collection configuration")
	}

	ccarray := make([]*pcommon.CollectionConfig, 0, len(*cconf))
	for _, cconfitem := range *cconf {
		p, err := cauthdsl.FromString(cconfitem.Policy)
		if err != nil {
			return nil, nil, errors.WithMessagef(err, "invalid policy %s", cconfitem.Policy)
		}

		cpc := &pcommon.CollectionPolicyConfig{
			Payload: &pcommon.CollectionPolicyConfig_SignaturePolicy{
				SignaturePolicy: p,
			},
		}

		cc := &pcommon.CollectionConfig{
			Payload: &pcommon.CollectionConfig_StaticCollectionConfig{
				StaticCollectionConfig: &pcommon.StaticCollectionConfig{
					Name:              cconfitem.Name,
					MemberOrgsPolicy:  cpc,
					RequiredPeerCount: cconfitem.RequiredCount,
					MaximumPeerCount:  cconfitem.MaxPeerCount,
					BlockToLive:       cconfitem.BlockToLive,
					MemberOnlyRead:    cconfitem.MemberOnlyRead,
					MemberOnlyWrite:   cconfitem.MemberOnlyWrite,
				},
			},
		}

		ccarray = append(ccarray, cc)
	}

	ccp := &pcommon.CollectionConfigPackage{Config: ccarray}
	ccpBytes, err := proto.Marshal(ccp)
	return ccp, ccpBytes, err
}

func checkChaincodeCmdParams(cmd *cobra.Command) error {
	
	if chaincodeName == common.UndefinedParamValue {
		return errors.Errorf("must supply value for %s name parameter", chainFuncName)
	}

	if cmd.Name() == instantiateCmdName || cmd.Name() == installCmdName ||
		cmd.Name() == upgradeCmdName || cmd.Name() == packageCmdName {
		if chaincodeVersion == common.UndefinedParamValue {
			return errors.Errorf("chaincode version is not provided for %s", cmd.Name())
		}

		if escc != common.UndefinedParamValue {
			logger.Infof("Using escc %s", escc)
		} else {
			logger.Info("Using default escc")
			escc = "escc"
		}

		if vscc != common.UndefinedParamValue {
			logger.Infof("Using vscc %s", vscc)
		} else {
			logger.Info("Using default vscc")
			vscc = "vscc"
		}

		if policy != common.UndefinedParamValue {
			p, err := cauthdsl.FromString(policy)
			if err != nil {
				return errors.Errorf("invalid policy %s", policy)
			}
			policyMarshalled = protoutil.MarshalOrPanic(p)
		}

		if collectionsConfigFile != common.UndefinedParamValue {
			var err error
			_, collectionConfigBytes, err = GetCollectionConfigFromFile(collectionsConfigFile)
			if err != nil {
				return errors.WithMessagef(err, "invalid collection configuration in file %s", collectionsConfigFile)
			}
		}
	}

	
	
	
	
	
	if chaincodeCtorJSON != "{}" {
		var f interface{}
		err := json.Unmarshal([]byte(chaincodeCtorJSON), &f)
		if err != nil {
			return errors.Wrap(err, "chaincode argument error")
		}
		m := f.(map[string]interface{})
		sm := make(map[string]interface{})
		for k := range m {
			sm[strings.ToLower(k)] = m[k]
		}
		_, argsPresent := sm["args"]
		_, funcPresent := sm["function"]
		if !argsPresent || (len(m) == 2 && !funcPresent) || len(m) > 2 {
			return errors.New("non-empty JSON chaincode parameters must contain the following keys: 'Args' or 'Function' and 'Args'")
		}
	} else {
		if cmd == nil || (cmd != chaincodeInstallCmd && cmd != chaincodePackageCmd) {
			return errors.New("empty JSON chaincode parameters must contain the following keys: 'Args' or 'Function' and 'Args'")
		}
	}

	return nil
}

func validatePeerConnectionParameters(cmdName string) error {
	if connectionProfile != common.UndefinedParamValue {
		networkConfig, err := common.GetConfig(connectionProfile)
		if err != nil {
			return err
		}
		if len(networkConfig.Channels[channelID].Peers) != 0 {
			peerAddresses = []string{}
			tlsRootCertFiles = []string{}
			for peer, peerChannelConfig := range networkConfig.Channels[channelID].Peers {
				if peerChannelConfig.EndorsingPeer {
					peerConfig, ok := networkConfig.Peers[peer]
					if !ok {
						return errors.Errorf("peer '%s' is defined in the channel config but doesn't have associated peer config", peer)
					}
					peerAddresses = append(peerAddresses, peerConfig.URL)
					tlsRootCertFiles = append(tlsRootCertFiles, peerConfig.TLSCACerts.Path)
				}
			}
		}
	}

	
	multiplePeersAllowed := map[string]bool{
		"invoke": true,
	}
	_, ok := multiplePeersAllowed[cmdName]
	if !ok && len(peerAddresses) > 1 {
		return errors.Errorf("'%s' command can only be executed against one peer. received %d", cmdName, len(peerAddresses))
	}

	if len(tlsRootCertFiles) > len(peerAddresses) {
		logger.Warningf("received more TLS root cert files (%d) than peer addresses (%d)", len(tlsRootCertFiles), len(peerAddresses))
	}

	if viper.GetBool("peer.tls.enabled") {
		if len(tlsRootCertFiles) != len(peerAddresses) {
			return errors.Errorf("number of peer addresses (%d) does not match the number of TLS root cert files (%d)", len(peerAddresses), len(tlsRootCertFiles))
		}
	} else {
		tlsRootCertFiles = nil
	}

	return nil
}


type ChaincodeCmdFactory struct {
	EndorserClients []pb.EndorserClient
	DeliverClients  []pb.DeliverClient
	Certificate     tls.Certificate
	Signer          identity.SignerSerializer
	BroadcastClient common.BroadcastClient
}


func InitCmdFactory(cmdName string, isEndorserRequired, isOrdererRequired bool) (*ChaincodeCmdFactory, error) {
	var err error
	var endorserClients []pb.EndorserClient
	var deliverClients []pb.DeliverClient
	if isEndorserRequired {
		if err = validatePeerConnectionParameters(cmdName); err != nil {
			return nil, errors.WithMessage(err, "error validating peer connection parameters")
		}
		for i, address := range peerAddresses {
			var tlsRootCertFile string
			if tlsRootCertFiles != nil {
				tlsRootCertFile = tlsRootCertFiles[i]
			}
			endorserClient, err := common.GetEndorserClientFnc(address, tlsRootCertFile)
			if err != nil {
				return nil, errors.WithMessagef(err, "error getting endorser client for %s", cmdName)
			}
			endorserClients = append(endorserClients, endorserClient)
			deliverClient, err := common.GetPeerDeliverClientFnc(address, tlsRootCertFile)
			if err != nil {
				return nil, errors.WithMessagef(err, "error getting deliver client for %s", cmdName)
			}
			deliverClients = append(deliverClients, deliverClient)
		}
		if len(endorserClients) == 0 {
			return nil, errors.New("no endorser clients retrieved - this might indicate a bug")
		}
	}
	certificate, err := common.GetCertificateFnc()
	if err != nil {
		return nil, errors.WithMessage(err, "error getting client cerificate")
	}

	signer, err := common.GetDefaultSignerFnc()
	if err != nil {
		return nil, errors.WithMessage(err, "error getting default signer")
	}

	var broadcastClient common.BroadcastClient
	if isOrdererRequired {
		if len(common.OrderingEndpoint) == 0 {
			if len(endorserClients) == 0 {
				return nil, errors.New("orderer is required, but no ordering endpoint or endorser client supplied")
			}
			endorserClient := endorserClients[0]

			orderingEndpoints, err := common.GetOrdererEndpointOfChainFnc(channelID, signer, endorserClient)
			if err != nil {
				return nil, errors.WithMessagef(err, "error getting channel (%s) orderer endpoint", channelID)
			}
			if len(orderingEndpoints) == 0 {
				return nil, errors.Errorf("no orderer endpoints retrieved for channel %s", channelID)
			}
			logger.Infof("Retrieved channel (%s) orderer endpoint: %s", channelID, orderingEndpoints[0])
			
			viper.Set("orderer.address", orderingEndpoints[0])
		}

		broadcastClient, err = common.GetBroadcastClientFnc()

		if err != nil {
			return nil, errors.WithMessage(err, "error getting broadcast client")
		}
	}
	return &ChaincodeCmdFactory{
		EndorserClients: endorserClients,
		DeliverClients:  deliverClients,
		Signer:          signer,
		BroadcastClient: broadcastClient,
		Certificate:     certificate,
	}, nil
}










func ChaincodeInvokeOrQuery(
	spec *pb.ChaincodeSpec,
	cID string,
	txID string,
	invoke bool,
	signer identity.SignerSerializer,
	certificate tls.Certificate,
	endorserClients []pb.EndorserClient,
	deliverClients []pb.DeliverClient,
	bc common.BroadcastClient,
) (*pb.ProposalResponse, error) {
	
	invocation := &pb.ChaincodeInvocationSpec{ChaincodeSpec: spec}

	creator, err := signer.Serialize()
	if err != nil {
		return nil, errors.WithMessage(err, "error serializing identity")
	}

	funcName := "invoke"
	if !invoke {
		funcName = "query"
	}

	
	var tMap map[string][]byte
	if transient != "" {
		if err := json.Unmarshal([]byte(transient), &tMap); err != nil {
			return nil, errors.Wrap(err, "error parsing transient string")
		}
	}

	prop, txid, err := protoutil.CreateChaincodeProposalWithTxIDAndTransient(pcommon.HeaderType_ENDORSER_TRANSACTION, cID, invocation, creator, txID, tMap)
	if err != nil {
		return nil, errors.WithMessagef(err, "error creating proposal for %s", funcName)
	}

	signedProp, err := protoutil.GetSignedProposal(prop, signer)
	if err != nil {
		return nil, errors.WithMessagef(err, "error creating signed proposal for %s", funcName)
	}
	var responses []*pb.ProposalResponse
	for _, endorser := range endorserClients {
		proposalResp, err := endorser.ProcessProposal(context.Background(), signedProp)
		if err != nil {
			return nil, errors.WithMessagef(err, "error endorsing %s", funcName)
		}
		responses = append(responses, proposalResp)
	}

	if len(responses) == 0 {
		
		return nil, errors.New("no proposal responses received - this might indicate a bug")
	}
	
	
	proposalResp := responses[0]

	if invoke {
		if proposalResp != nil {
			if proposalResp.Response.Status >= shim.ERRORTHRESHOLD {
				return proposalResp, nil
			}
			
			env, err := protoutil.CreateSignedTx(prop, signer, responses...)
			if err != nil {
				return proposalResp, errors.WithMessage(err, "could not assemble transaction")
			}
			var dg *DeliverGroup
			var ctx context.Context
			if waitForEvent {
				var cancelFunc context.CancelFunc
				ctx, cancelFunc = context.WithTimeout(context.Background(), waitForEventTimeout)
				defer cancelFunc()

				dg = NewDeliverGroup(
					deliverClients,
					peerAddresses,
					signer,
					certificate,
					channelID,
					txid,
				)
				
				err := dg.Connect(ctx)
				if err != nil {
					return nil, err
				}
			}

			
			if err = bc.Send(env); err != nil {
				return proposalResp, errors.WithMessagef(err, "error sending transaction for %s", funcName)
			}

			if dg != nil && ctx != nil {
				
				err = dg.Wait(ctx)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return proposalResp, nil
}








type DeliverGroup struct {
	Clients     []*DeliverClient
	Certificate tls.Certificate
	ChannelID   string
	TxID        string
	Signer      identity.SignerSerializer
	mutex       sync.Mutex
	Error       error
	wg          sync.WaitGroup
}



type DeliverClient struct {
	Client     pb.DeliverClient
	Connection pb.Deliver_DeliverClient
	Address    string
}

func NewDeliverGroup(
	deliverClients []pb.DeliverClient,
	peerAddresses []string,
	signer identity.SignerSerializer,
	certificate tls.Certificate,
	channelID string,
	txid string,
) *DeliverGroup {
	clients := make([]*DeliverClient, len(deliverClients))
	for i, client := range deliverClients {
		dc := &DeliverClient{
			Client:  client,
			Address: peerAddresses[i],
		}
		clients[i] = dc
	}

	dg := &DeliverGroup{
		Clients:     clients,
		Certificate: certificate,
		ChannelID:   channelID,
		TxID:        txid,
		Signer:      signer,
	}

	return dg
}





func (dg *DeliverGroup) Connect(ctx context.Context) error {
	dg.wg.Add(len(dg.Clients))
	for _, client := range dg.Clients {
		go dg.ClientConnect(ctx, client)
	}
	readyCh := make(chan struct{})
	go dg.WaitForWG(readyCh)

	select {
	case <-readyCh:
		if dg.Error != nil {
			err := errors.WithMessage(dg.Error, "failed to connect to deliver on all peers")
			return err
		}
	case <-ctx.Done():
		err := errors.New("timed out waiting for connection to deliver on all peers")
		return err
	}

	return nil
}




func (dg *DeliverGroup) ClientConnect(ctx context.Context, dc *DeliverClient) {
	defer dg.wg.Done()
	df, err := dc.Client.DeliverFiltered(ctx)
	if err != nil {
		err = errors.WithMessagef(err, "error connecting to deliver filtered at %s", dc.Address)
		dg.setError(err)
		return
	}
	defer df.CloseSend()
	dc.Connection = df

	envelope := createDeliverEnvelope(dg.ChannelID, dg.Certificate, dg.Signer)
	err = df.Send(envelope)
	if err != nil {
		err = errors.WithMessagef(err, "error sending deliver seek info envelope to %s", dc.Address)
		dg.setError(err)
		return
	}
}




func (dg *DeliverGroup) Wait(ctx context.Context) error {
	if len(dg.Clients) == 0 {
		return nil
	}

	dg.wg.Add(len(dg.Clients))
	for _, client := range dg.Clients {
		go dg.ClientWait(client)
	}
	readyCh := make(chan struct{})
	go dg.WaitForWG(readyCh)

	select {
	case <-readyCh:
		if dg.Error != nil {
			return dg.Error
		}
	case <-ctx.Done():
		err := errors.New("timed out waiting for txid on all peers")
		return err
	}

	return nil
}



func (dg *DeliverGroup) ClientWait(dc *DeliverClient) {
	defer dg.wg.Done()
	for {
		resp, err := dc.Connection.Recv()
		if err != nil {
			err = errors.WithMessagef(err, "error receiving from deliver filtered at %s", dc.Address)
			dg.setError(err)
			return
		}
		switch r := resp.Type.(type) {
		case *pb.DeliverResponse_FilteredBlock:
			filteredTransactions := r.FilteredBlock.FilteredTransactions
			for _, tx := range filteredTransactions {
				if tx.Txid == dg.TxID {
					logger.Infof("txid [%s] committed with status (%s) at %s", dg.TxID, tx.TxValidationCode, dc.Address)
					if tx.TxValidationCode != pb.TxValidationCode_VALID {
						err = errors.Errorf("transaction invalidated with status (%s)", tx.TxValidationCode)
						dg.setError(err)
					}
					return
				}
			}
		case *pb.DeliverResponse_Status:
			err = errors.Errorf("deliver completed with status (%s) before txid received", r.Status)
			dg.setError(err)
			return
		default:
			err = errors.Errorf("received unexpected response type (%T) from %s", r, dc.Address)
			dg.setError(err)
			return
		}
	}
}



func (dg *DeliverGroup) WaitForWG(readyCh chan struct{}) {
	dg.wg.Wait()
	close(readyCh)
}


func (dg *DeliverGroup) setError(err error) {
	dg.mutex.Lock()
	dg.Error = err
	dg.mutex.Unlock()
}

func createDeliverEnvelope(
	channelID string,
	certificate tls.Certificate,
	signer identity.SignerSerializer,
) *pcommon.Envelope {
	var tlsCertHash []byte
	
	if len(certificate.Certificate) > 0 {
		tlsCertHash = util.ComputeSHA256(certificate.Certificate[0])
	}

	start := &ab.SeekPosition{
		Type: &ab.SeekPosition_Newest{
			Newest: &ab.SeekNewest{},
		},
	}

	stop := &ab.SeekPosition{
		Type: &ab.SeekPosition_Specified{
			Specified: &ab.SeekSpecified{
				Number: math.MaxUint64,
			},
		},
	}

	seekInfo := &ab.SeekInfo{
		Start:    start,
		Stop:     stop,
		Behavior: ab.SeekInfo_BLOCK_UNTIL_READY,
	}

	env, err := protoutil.CreateSignedEnvelopeWithTLSBinding(
		pcommon.HeaderType_DELIVER_SEEK_INFO,
		channelID,
		signer,
		seekInfo,
		int32(0),
		uint64(0),
		tlsCertHash,
	)
	if err != nil {
		logger.Errorf("Error signing envelope: %s", err)
		return nil
	}

	return env
}
