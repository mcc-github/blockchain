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
	"time"

	"github.com/mcc-github/blockchain/common/channelconfig"
	"github.com/mcc-github/blockchain/common/configtx"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/common/policies"
	"github.com/mcc-github/blockchain/common/tools/protolator"
	"github.com/mcc-github/blockchain/common/util"
	"github.com/mcc-github/blockchain/core/comm"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protoutil"
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

func (cbc ConnByCertMap) Size() int {
	return len(cbc)
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
	cc, isClientConfig := val.(comm.ClientConfig)
	if !isClientConfig {
		err := errors.Errorf("value stored is %v, not comm.ClientConfig",
			reflect.TypeOf(val))
		return comm.ClientConfig{}, err
	}
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
	
	
	
	
	
	
	VerifyBlockSignature(sd []*protoutil.SignedData, config *common.ConfigEnvelope) error
}



type BlockSequenceVerifier func(blocks []*common.Block, channel string) error


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
		if err == errNotAConfig {
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

var errNotAConfig = errors.New("not a config block")



func ConfigFromBlock(block *common.Block) (*common.ConfigEnvelope, error) {
	if block == nil || block.Data == nil || len(block.Data.Data) == 0 {
		return nil, errors.New("empty block")
	}
	txn := block.Data.Data[0]
	env, err := protoutil.GetEnvelopeFromBlock(txn)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	payload, err := protoutil.GetPayload(env)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if block.Header.Number == 0 {
		configEnvelope, err := configtx.UnmarshalConfigEnvelope(payload.Data)
		if err != nil {
			return nil, errors.Wrap(err, "invalid config envelope")
		}
		return configEnvelope, nil
	}
	if payload.Header == nil {
		return nil, errors.New("nil header in payload")
	}
	chdr, err := protoutil.UnmarshalChannelHeader(payload.Header.ChannelHeader)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if common.HeaderType(chdr.Type) != common.HeaderType_CONFIG {
		return nil, errNotAConfig
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
	dataHash := protoutil.BlockDataHash(block.Data)
	
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
		if !bytes.Equal(block.Header.PreviousHash, protoutil.BlockHeaderHash(prevBlock.Header)) {
			claimedPrevHash := hex.EncodeToString(block.Header.PreviousHash)
			actualPrevHash := hex.EncodeToString(protoutil.BlockHeaderHash(prevBlock.Header))
			return errors.Errorf("block [%d]'s hash (%s) mismatches block [%d]'s prev block hash (%s)",
				prevSeq, actualPrevHash, currSeq, claimedPrevHash)
		}
	}
	return nil
}


func SignatureSetFromBlock(block *common.Block) ([]*protoutil.SignedData, error) {
	if block.Metadata == nil || len(block.Metadata.Metadata) <= int(common.BlockMetadataIndex_SIGNATURES) {
		return nil, errors.New("no metadata in block")
	}
	metadata, err := protoutil.GetMetadataFromBlock(block, common.BlockMetadataIndex_SIGNATURES)
	if err != nil {
		return nil, errors.Errorf("failed unmarshaling medatata for signatures: %v", err)
	}

	var signatureSet []*protoutil.SignedData
	for _, metadataSignature := range metadata.Signatures {
		sigHdr, err := protoutil.GetSignatureHeader(metadataSignature.SignatureHeader)
		if err != nil {
			return nil, errors.Errorf("failed unmarshaling signature header for block with id %d: %v",
				block.Header.Number, err)
		}
		signatureSet = append(signatureSet,
			&protoutil.SignedData{
				Identity: sigHdr.Creator,
				Data: util.ConcatenateBytes(metadata.Value,
					metadataSignature.SignatureHeader, protoutil.BlockHeaderBytes(block.Header)),
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
	envelopeConfig, err := protoutil.ExtractEnvelope(block, 0)
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




type VerifierFactory interface {
	
	VerifierFromConfig(configuration *common.ConfigEnvelope, channel string) (BlockVerifier, error)
}


type VerificationRegistry struct {
	LoadVerifier       func(chain string) BlockVerifier
	Logger             *flogging.FabricLogger
	VerifierFactory    VerifierFactory
	VerifiersByChannel map[string]BlockVerifier
}


func (vr *VerificationRegistry) RegisterVerifier(chain string) {
	if _, exists := vr.VerifiersByChannel[chain]; exists {
		vr.Logger.Debugf("No need to register verifier for chain %s", chain)
		return
	}

	v := vr.LoadVerifier(chain)
	if v == nil {
		vr.Logger.Errorf("Failed loading verifier for chain %s", chain)
		return
	}

	vr.VerifiersByChannel[chain] = v
	vr.Logger.Infof("Registered verifier for chain %s", chain)
}


func (vr *VerificationRegistry) RetrieveVerifier(channel string) BlockVerifier {
	verifier, exists := vr.VerifiersByChannel[channel]
	if exists {
		return verifier
	}
	vr.Logger.Errorf("No verifier for channel %s exists", channel)
	return nil
}



func (vr *VerificationRegistry) BlockCommitted(block *common.Block, channel string) {
	conf, err := ConfigFromBlock(block)
	
	if err == errNotAConfig {
		vr.Logger.Debugf("Committed block [%d] for channel %s that is not a config block",
			block.Header.Number, channel)
		return
	}
	
	if err != nil {
		vr.Logger.Errorf("Failed parsing block of channel %s: %v, content: %s",
			channel, err, BlockToString(block))
		return
	}

	
	verifier, err := vr.VerifierFactory.VerifierFromConfig(conf, channel)
	if err != nil {
		vr.Logger.Errorf("Failed creating a verifier from a config block for channel %s: %v, content: %s",
			channel, err, BlockToString(block))
		return
	}

	vr.VerifiersByChannel[channel] = verifier

	vr.Logger.Debugf("Committed config block [%d] for channel %s", block.Header.Number, channel)
}


func BlockToString(block *common.Block) string {
	buff := &bytes.Buffer{}
	protolator.DeepMarshalJSON(buff, block)
	return buff.String()
}


type BlockCommitFunc func(block *common.Block, channel string)


type LedgerInterceptor struct {
	Channel              string
	InterceptBlockCommit BlockCommitFunc
	LedgerWriter
}


func (interceptor *LedgerInterceptor) Append(block *common.Block) error {
	defer interceptor.InterceptBlockCommit(block, interceptor.Channel)
	return interceptor.LedgerWriter.Append(block)
}


type BlockVerifierAssembler struct {
	Logger *flogging.FabricLogger
}


func (bva *BlockVerifierAssembler) VerifierFromConfig(configuration *common.ConfigEnvelope, channel string) (BlockVerifier, error) {
	bundle, err := channelconfig.NewBundle(channel, configuration.Config)
	if err != nil {
		return nil, errors.Wrap(err, "failed extracting bundle from envelope")
	}
	policyMgr := bundle.PolicyManager()

	return &BlockValidationPolicyVerifier{
		Logger:    bva.Logger,
		PolicyMgr: policyMgr,
		Channel:   channel,
	}, nil
}


type BlockValidationPolicyVerifier struct {
	Logger    *flogging.FabricLogger
	Channel   string
	PolicyMgr policies.Manager
}


func (bv *BlockValidationPolicyVerifier) VerifyBlockSignature(sd []*protoutil.SignedData, envelope *common.ConfigEnvelope) error {
	policyMgr := bv.PolicyMgr
	
	if envelope != nil {
		bundle, err := channelconfig.NewBundle(bv.Channel, envelope.Config)
		if err != nil {
			buff := &bytes.Buffer{}
			protolator.DeepMarshalJSON(buff, envelope.Config)
			bv.Logger.Errorf("Failed creating a new bundle for channel %s, Config content is: %s", bv.Channel, buff.String())
			return err
		}
		bv.Logger.Infof("Initializing new PolicyManager for channel %s", bv.Channel)
		policyMgr = bundle.PolicyManager()
	}
	policy, exists := policyMgr.GetPolicy(policies.BlockValidation)
	if !exists {
		return errors.Errorf("policy %s wasn't found", policies.BlockValidation)
	}
	return policy.Evaluate(sd)
}




type BlockRetriever interface {
	
	
	Block(number uint64) *common.Block
}


func LastConfigBlock(block *common.Block, blockRetriever BlockRetriever) (*common.Block, error) {
	if block == nil {
		return nil, errors.New("nil block")
	}
	if blockRetriever == nil {
		return nil, errors.New("nil blockRetriever")
	}
	if block.Metadata == nil || len(block.Metadata.Metadata) <= int(common.BlockMetadataIndex_LAST_CONFIG) {
		return nil, errors.New("no metadata in block")
	}
	lastConfigBlockNum, err := protoutil.GetLastConfigIndexFromBlock(block)
	if err != nil {
		return nil, err
	}
	lastConfigBlock := blockRetriever.Block(lastConfigBlockNum)
	if lastConfigBlock == nil {
		return nil, errors.Errorf("unable to retrieve last config block [%d]", lastConfigBlockNum)
	}
	return lastConfigBlock, nil
}


type StreamCountReporter struct {
	Metrics *Metrics
	count   uint32
}

func (scr *StreamCountReporter) Increment() {
	count := atomic.AddUint32(&scr.count, 1)
	scr.Metrics.reportStreamCount(count)
}

func (scr *StreamCountReporter) Decrement() {
	count := atomic.AddUint32(&scr.count, ^uint32(0))
	scr.Metrics.reportStreamCount(count)
}

type certificateExpirationCheck struct {
	minimumExpirationWarningInterval time.Duration
	expiresAt                        time.Time
	expirationWarningThreshold       time.Duration
	lastWarning                      time.Time
	nodeName                         string
	endpoint                         string
	alert                            func(string, ...interface{})
}

func (exp *certificateExpirationCheck) checkExpiration(currentTime time.Time, channel string) {
	timeLeft := exp.expiresAt.Sub(currentTime)
	if timeLeft > exp.expirationWarningThreshold {
		return
	}

	timeSinceLastWarning := currentTime.Sub(exp.lastWarning)
	if timeSinceLastWarning < exp.minimumExpirationWarningInterval {
		return
	}

	exp.alert("Certificate of %s from %s for channel %s expires in less than %v",
		exp.nodeName, exp.endpoint, channel, timeLeft)
	exp.lastWarning = currentTime
}
