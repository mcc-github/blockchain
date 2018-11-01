/*
Copyright IBM Corp. 2017 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package cluster

import (
	"bytes"
	"encoding/hex"
	"encoding/pem"
	"reflect"
	"sync/atomic"

	"github.com/mcc-github/blockchain/common/channelconfig"
	"github.com/mcc-github/blockchain/common/configtx"
	"github.com/mcc-github/blockchain/common/util"
	"github.com/mcc-github/blockchain/core/comm"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/utils"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)



type ConnByCertMap map[string]*grpc.ClientConn



func (cbc ConnByCertMap) Lookup(cert []byte) (*grpc.ClientConn, bool) {
	conn, ok := cbc[string(cert)]
	return conn, ok
}


func (cbc ConnByCertMap) Put(cert []byte, conn *grpc.ClientConn) {
	cbc[string(cert)] = conn
}


func (cbc ConnByCertMap) Remove(cert []byte) {
	delete(cbc, string(cert))
}


type MemberMapping map[uint64]*Stub


func (mp MemberMapping) Put(stub *Stub) {
	mp[stub.ID] = stub
}


func (mp MemberMapping) ByID(ID uint64) *Stub {
	return mp[ID]
}


func (mp MemberMapping) LookupByClientCert(cert []byte) *Stub {
	for _, stub := range mp {
		if bytes.Equal(stub.ClientTLSCert, cert) {
			return stub
		}
	}
	return nil
}



func (mp MemberMapping) ServerCertificates() StringSet {
	res := make(StringSet)
	for _, member := range mp {
		res[string(member.ServerTLSCert)] = struct{}{}
	}
	return res
}


type StringSet map[string]struct{}


func (ss StringSet) union(set StringSet) {
	for k := range set {
		ss[k] = struct{}{}
	}
}


func (ss StringSet) subtract(set StringSet) {
	for k := range set {
		delete(ss, k)
	}
}




type PredicateDialer struct {
	Config atomic.Value
}


func NewTLSPinningDialer(config comm.ClientConfig) *PredicateDialer {
	d := &PredicateDialer{}
	d.SetConfig(config)
	return d
}



func (dialer *PredicateDialer) ClientConfig() (comm.ClientConfig, error) {
	val := dialer.Config.Load()
	if val == nil {
		return comm.ClientConfig{}, errors.New("client config not initialized")
	}
	if cc, isClientConfig := val.(comm.ClientConfig); !isClientConfig {
		err := errors.Errorf("value stored is %v, not comm.ClientConfig",
			reflect.TypeOf(val))
		return comm.ClientConfig{}, err
	} else {
		if cc.SecOpts == nil {
			return comm.ClientConfig{}, errors.New("SecOpts is nil")
		}
		
		secOpts := *cc.SecOpts
		return comm.ClientConfig{
			AsyncConnect: cc.AsyncConnect,
			Timeout:      cc.Timeout,
			SecOpts:      &secOpts,
			KaOpts:       cc.KaOpts,
		}, nil
	}
}


func (dialer *PredicateDialer) SetConfig(config comm.ClientConfig) {
	configCopy := comm.ClientConfig{
		AsyncConnect: config.AsyncConnect,
		Timeout:      config.Timeout,
		SecOpts:      &comm.SecureOptions{},
		KaOpts:       &comm.KeepaliveOptions{},
	}
	
	if config.SecOpts != nil {
		*configCopy.SecOpts = *config.SecOpts
	}
	if config.KaOpts != nil {
		*configCopy.KaOpts = *config.KaOpts
	} else {
		configCopy.KaOpts = nil
	}

	dialer.Config.Store(configCopy)
}



func (dialer *PredicateDialer) Dial(address string, verifyFunc RemoteVerifier) (*grpc.ClientConn, error) {
	cfg := dialer.Config.Load().(comm.ClientConfig)
	cfg.SecOpts.VerifyCertificate = verifyFunc
	client, err := comm.NewGRPCClient(cfg)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return client.NewConnection(address, "")
}



func DERtoPEM(der []byte) string {
	return string(pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: der,
	}))
}



type StandardDialer struct {
	Dialer *PredicateDialer
}

func (bdp *StandardDialer) Dial(address string) (*grpc.ClientConn, error) {
	return bdp.Dialer.Dial(address, nil)
}




type BlockVerifier interface {
	
	
	
	
	
	
	VerifyBlockSignature(sd []*common.SignedData, config *common.ConfigEnvelope) error
}



type BlockSequenceVerifier func([]*common.Block) error


type Dialer interface {
	Dial(address string) (*grpc.ClientConn, error)
}



func VerifyBlocks(blockBuff []*common.Block, signatureVerifier BlockVerifier) error {
	if len(blockBuff) == 0 {
		return errors.New("buffer is empty")
	}
	
	
	
	for i := range blockBuff {
		if err := VerifyBlockHash(i, blockBuff); err != nil {
			return err
		}
	}

	var config *common.ConfigEnvelope
	
	
	
	for _, block := range blockBuff {
		configFromBlock, err := ConfigFromBlock(block)
		if err == notAConfig {
			continue
		}
		if err != nil {
			return err
		}
		
		if err := VerifyBlockSignature(block, signatureVerifier, config); err != nil {
			return err
		}
		config = configFromBlock
	}

	
	lastBlock := blockBuff[len(blockBuff)-1]
	return VerifyBlockSignature(lastBlock, signatureVerifier, config)
}

var notAConfig = errors.New("not a config block")



func ConfigFromBlock(block *common.Block) (*common.ConfigEnvelope, error) {
	if block == nil || block.Data == nil || len(block.Data.Data) == 0 {
		return nil, errors.New("empty block")
	}
	txn := block.Data.Data[0]
	env, err := utils.GetEnvelopeFromBlock(txn)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	payload, err := utils.GetPayload(env)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if payload.Header == nil {
		return nil, errors.New("nil header in payload")
	}
	chdr, err := utils.UnmarshalChannelHeader(payload.Header.ChannelHeader)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if common.HeaderType(chdr.Type) != common.HeaderType_CONFIG {
		return nil, notAConfig
	}
	configEnvelope, err := configtx.UnmarshalConfigEnvelope(payload.Data)
	if err != nil {
		return nil, errors.Wrap(err, "invalid config envelope")
	}
	return configEnvelope, nil
}



func VerifyBlockHash(indexInBuffer int, blockBuff []*common.Block) error {
	if len(blockBuff) <= indexInBuffer {
		return errors.Errorf("index %d out of bounds (total %d blocks)", indexInBuffer, len(blockBuff))
	}
	block := blockBuff[indexInBuffer]
	if block.Header == nil {
		return errors.New("missing block header")
	}
	seq := block.Header.Number
	dataHash := block.Data.Hash()
	
	if !bytes.Equal(dataHash, block.Header.DataHash) {
		computedHash := hex.EncodeToString(dataHash)
		claimedHash := hex.EncodeToString(block.Header.DataHash)
		return errors.Errorf("computed hash of block (%d) (%s) doesn't match claimed hash (%s)",
			seq, computedHash, claimedHash)
	}
	
	if indexInBuffer > 0 {
		prevBlock := blockBuff[indexInBuffer-1]
		currSeq := block.Header.Number
		if prevBlock.Header == nil {
			return errors.New("previous block header is nil")
		}
		prevSeq := prevBlock.Header.Number
		if prevSeq+1 != currSeq {
			return errors.Errorf("sequences %d and %d were received consecutively", prevSeq, currSeq)
		}
		if !bytes.Equal(block.Header.PreviousHash, prevBlock.Header.Hash()) {
			claimedPrevHash := hex.EncodeToString(block.Header.PreviousHash)
			actualPrevHash := hex.EncodeToString(prevBlock.Header.Hash())
			return errors.Errorf("block %d's hash (%s) mismatches %d's prev block hash (%s)",
				currSeq, actualPrevHash, prevSeq, claimedPrevHash)
		}
	}
	return nil
}


func SignatureSetFromBlock(block *common.Block) ([]*common.SignedData, error) {
	if block.Metadata == nil || len(block.Metadata.Metadata) <= int(common.BlockMetadataIndex_SIGNATURES) {
		return nil, errors.New("no metadata in block")
	}
	metadata, err := utils.GetMetadataFromBlock(block, common.BlockMetadataIndex_SIGNATURES)
	if err != nil {
		return nil, errors.Errorf("failed unmarshaling medatata for signatures: %v", err)
	}

	var signatureSet []*common.SignedData
	for _, metadataSignature := range metadata.Signatures {
		sigHdr, err := utils.GetSignatureHeader(metadataSignature.SignatureHeader)
		if err != nil {
			return nil, errors.Errorf("failed unmarshaling signature header for block with id %d: %v",
				block.Header.Number, err)
		}
		signatureSet = append(signatureSet,
			&common.SignedData{
				Identity: sigHdr.Creator,
				Data: util.ConcatenateBytes(metadata.Value,
					metadataSignature.SignatureHeader, block.Header.Bytes()),
				Signature: metadataSignature.Signature,
			},
		)
	}
	return signatureSet, nil
}


func VerifyBlockSignature(block *common.Block, verifier BlockVerifier, config *common.ConfigEnvelope) error {
	signatureSet, err := SignatureSetFromBlock(block)
	if err != nil {
		return err
	}
	return verifier.VerifyBlockSignature(signatureSet, config)
}



type EndpointConfig struct {
	TLSRootCAs [][]byte
	Endpoints  []string
}



func EndpointconfigFromConfigBlock(block *common.Block) (*EndpointConfig, error) {
	if block == nil {
		return nil, errors.New("nil block")
	}
	envelopeConfig, err := utils.ExtractEnvelope(block, 0)
	if err != nil {
		return nil, err
	}
	var tlsCACerts [][]byte
	bundle, err := channelconfig.NewBundleFromEnvelope(envelopeConfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed extracting bundle from envelope")
	}
	msps, err := bundle.MSPManager().GetMSPs()
	if err != nil {
		return nil, errors.Wrap(err, "failed obtaining MSPs from MSPManager")
	}
	ordererConfig, ok := bundle.OrdererConfig()
	if !ok {
		return nil, errors.New("failed obtaining orderer config from bundle")
	}
	for _, org := range ordererConfig.Organizations() {
		msp := msps[org.MSPID()]
		if msp == nil {
			return nil, errors.Errorf("no MSP found for MSP with ID of %s", org.MSPID())
		}
		tlsCACerts = append(tlsCACerts, msp.GetTLSRootCerts()...)
	}
	return &EndpointConfig{
		Endpoints:  bundle.ChannelConfig().OrdererAddresses(),
		TLSRootCAs: tlsCACerts,
	}, nil
}
