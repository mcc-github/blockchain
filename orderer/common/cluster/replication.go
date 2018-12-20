/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package cluster

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"time"

	"github.com/mcc-github/blockchain/common/crypto"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/core/comm"
	"github.com/mcc-github/blockchain/orderer/common/localconfig"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/utils"
	"github.com/pkg/errors"
)

const (
	
	RetryTimeout = time.Second * 10
)




func PullerConfigFromTopLevelConfig(systemChannel string, conf *localconfig.TopLevel, tlsKey, tlsCert []byte, signer crypto.LocalSigner) PullerConfig {
	return PullerConfig{
		Channel:             systemChannel,
		MaxTotalBufferBytes: conf.General.Cluster.ReplicationBufferSize,
		Timeout:             conf.General.Cluster.RPCTimeout,
		TLSKey:              tlsKey,
		TLSCert:             tlsCert,
		Signer:              signer,
	}
}




type LedgerWriter interface {
	
	Append(block *common.Block) error

	
	Height() uint64
}




type LedgerFactory interface {
	
	
	GetOrCreate(chainID string) (LedgerWriter, error)
}




type ChannelLister interface {
	
	Channels() []string
	
	Close()
}


type Replicator struct {
	SystemChannel    string
	ChannelLister    ChannelLister
	Logger           *flogging.FabricLogger
	Puller           *BlockPuller
	BootBlock        *common.Block
	AmIPartOfChannel selfMembershipPredicate
	LedgerFactory    LedgerFactory
}



func (r *Replicator) IsReplicationNeeded() (bool, error) {
	systemChannelLedger, err := r.LedgerFactory.GetOrCreate(r.SystemChannel)
	if err != nil {
		return false, err
	}

	height := systemChannelLedger.Height()
	var lastBlockSeq uint64
	
	
	if height == 0 {
		lastBlockSeq = 0
	} else {
		lastBlockSeq = height - 1
	}

	if r.BootBlock.Header.Number > lastBlockSeq {
		return true, nil
	}
	return false, nil
}


func (r *Replicator) ReplicateChains() {
	channels := r.discoverChannels()
	channels2Pull := r.channelsToPull(channels)
	r.Logger.Info("Found myself in", len(channels2Pull), "channels:", channels2Pull)
	for _, channel := range channels2Pull {
		r.PullChannel(channel)
	}
	
	if err := r.PullChannel(r.SystemChannel); err != nil {
		r.Logger.Panicf("Failed pulling system channel: %v", err)
	}
}

func (r *Replicator) discoverChannels() []string {
	r.Logger.Debug("Entering")
	defer r.Logger.Debug("Exiting")
	channels := r.ChannelLister.Channels()
	r.Logger.Info("Discovered", len(channels), "channels:", channels)
	r.ChannelLister.Close()
	return channels
}



func (r *Replicator) PullChannel(channel string) error {
	r.Logger.Info("Pulling channel", channel)
	puller := r.Puller.Clone()
	defer puller.Close()
	puller.Channel = channel

	endpoint, latestHeight := latestHeightAndEndpoint(puller)
	if endpoint == "" {
		return errors.Errorf("failed obtaining the latest block for channel %s", channel)
	}
	r.Logger.Info("Latest block height for channel", channel, "is", latestHeight)
	
	
	
	if channel == r.SystemChannel && latestHeight-1 < r.BootBlock.Header.Number {
		return errors.Errorf("latest height found among system channel(%s) orderers is %d, but the boot block's "+
			"sequence is %d", r.SystemChannel, latestHeight, r.BootBlock.Header.Number)
	}
	return r.pullChannelBlocks(channel, puller, latestHeight)
}

func (r *Replicator) pullChannelBlocks(channel string, puller ChainPuller, latestHeight uint64) error {
	ledger, err := r.LedgerFactory.GetOrCreate(channel)
	if err != nil {
		r.Logger.Panicf("Failed to create a ledger for channel %s: %v", channel, err)
	}
	
	genesisBlock := puller.PullBlock(0)
	r.appendBlockIfNeeded(genesisBlock, ledger, channel)
	actualPrevHash := genesisBlock.Header.Hash()

	for seq := uint64(1); seq < latestHeight; seq++ {
		block := puller.PullBlock(seq)
		reportedPrevHash := block.Header.PreviousHash
		if !bytes.Equal(reportedPrevHash, actualPrevHash) {
			return errors.Errorf("block header mismatch on sequence %d, expected %x, got %x",
				block.Header.Number, actualPrevHash, reportedPrevHash)
		}
		actualPrevHash = block.Header.Hash()
		if channel == r.SystemChannel && block.Header.Number == r.BootBlock.Header.Number {
			r.compareBootBlockWithSystemChannelLastConfigBlock(block)
			r.appendBlockIfNeeded(block, ledger, channel)
			
			return nil
		}
		r.appendBlockIfNeeded(block, ledger, channel)
	}
	return nil
}

func (r *Replicator) appendBlockIfNeeded(block *common.Block, ledger LedgerWriter, channel string) {
	currHeight := ledger.Height()
	if currHeight >= block.Header.Number+1 {
		r.Logger.Infof("Already at height %d for channel %s, skipping commit of block %d",
			currHeight, channel, block.Header.Number)
		return
	}
	if err := ledger.Append(block); err != nil {
		r.Logger.Panicf("Failed to write block %d: %v", block.Header.Number, err)
	}
	r.Logger.Infof("Committed block %d for channel %s", block.Header.Number, channel)
}

func (r *Replicator) compareBootBlockWithSystemChannelLastConfigBlock(block *common.Block) {
	
	block.Header.DataHash = block.Data.Hash()

	bootBlockHash := r.BootBlock.Header.Hash()
	retrievedBlockHash := block.Header.Hash()
	if bytes.Equal(bootBlockHash, retrievedBlockHash) {
		return
	}
	r.Logger.Panicf("Block header mismatch on last system channel block, expected %s, got %s",
		hex.EncodeToString(bootBlockHash), hex.EncodeToString(retrievedBlockHash))
}

func (r *Replicator) channelsToPull(channels []string) []string {
	r.Logger.Info("Will now pull channels:", channels)
	var channelsToPull []string
	for _, channel := range channels {
		r.Logger.Info("Pulling chain for", channel)
		puller := r.Puller.Clone()
		puller.Channel = channel
		
		
		bufferSize := puller.MaxTotalBufferBytes
		puller.MaxTotalBufferBytes = 1
		err := Participant(puller, r.AmIPartOfChannel)
		puller.Close()
		
		puller.MaxTotalBufferBytes = bufferSize
		if err == ErrNotInChannel {
			r.Logger.Info("I do not belong to channel", channel, ", skipping chain retrieval")
			continue
		}
		if err != nil {
			r.Logger.Panicf("Failed classifying whether I belong to channel %s: %v, skipping chain retrieval", channel, err)
			continue
		}
		channelsToPull = append(channelsToPull, channel)
	}
	return channelsToPull
}


type PullerConfig struct {
	TLSKey              []byte
	TLSCert             []byte
	Timeout             time.Duration
	Signer              crypto.LocalSigner
	Channel             string
	MaxTotalBufferBytes int
}


func BlockPullerFromConfigBlock(conf PullerConfig, block *common.Block) (*BlockPuller, error) {
	if block == nil {
		return nil, errors.New("nil block")
	}

	endpointconfig, err := EndpointconfigFromConfigBlock(block)
	if err != nil {
		return nil, err
	}

	dialer := &StandardDialer{
		Dialer: NewTLSPinningDialer(comm.ClientConfig{
			Timeout: conf.Timeout,
			SecOpts: &comm.SecureOptions{
				ServerRootCAs:     endpointconfig.TLSRootCAs,
				Certificate:       conf.TLSCert,
				Key:               conf.TLSKey,
				RequireClientCert: true,
				UseTLS:            true,
			},
		})}

	tlsCertAsDER, _ := pem.Decode(conf.TLSCert)
	if tlsCertAsDER == nil {
		return nil, errors.Errorf("unable to decode TLS certificate PEM: %s", base64.StdEncoding.EncodeToString(conf.TLSCert))
	}

	return &BlockPuller{
		Logger:  flogging.MustGetLogger("orderer.common.cluster.replication"),
		Dialer:  dialer,
		TLSCert: tlsCertAsDER.Bytes,
		VerifyBlockSequence: func(blocks []*common.Block) error {
			return VerifyBlocks(blocks, &NoopBlockVerifier{})
		},
		MaxTotalBufferBytes: conf.MaxTotalBufferBytes,
		Endpoints:           endpointconfig.Endpoints,
		RetryTimeout:        RetryTimeout,
		FetchTimeout:        conf.Timeout,
		Channel:             conf.Channel,
		Signer:              conf.Signer,
	}, nil
}


type NoopBlockVerifier struct{}


func (*NoopBlockVerifier) VerifyBlockSignature(sd []*common.SignedData, config *common.ConfigEnvelope) error {
	return nil
}




type ChainPuller interface {
	
	PullBlock(seq uint64) *common.Block

	
	HeightsByEndpoints() map[string]uint64

	
	Close()
}


type ChainInspector struct {
	Logger          *flogging.FabricLogger
	Puller          ChainPuller
	LastConfigBlock *common.Block
}


var ErrNotInChannel = errors.New("not in the channel")


type selfMembershipPredicate func(configBlock *common.Block) error






func Participant(puller ChainPuller, analyzeLastConfBlock selfMembershipPredicate) error {
	endpoint, latestHeight := latestHeightAndEndpoint(puller)
	if endpoint == "" {
		return errors.New("no available orderer")
	}
	lastBlock := puller.PullBlock(latestHeight - 1)
	lastConfNumber, err := lastConfigFromBlock(lastBlock)
	if err != nil {
		return err
	}
	
	
	
	puller.Close()
	lastConfigBlock := puller.PullBlock(lastConfNumber)
	return analyzeLastConfBlock(lastConfigBlock)
}

func latestHeightAndEndpoint(puller ChainPuller) (string, uint64) {
	var maxHeight uint64
	var mostUpToDateEndpoint string
	for endpoint, height := range puller.HeightsByEndpoints() {
		if height >= maxHeight {
			maxHeight = height
			mostUpToDateEndpoint = endpoint
		}
	}
	return mostUpToDateEndpoint, maxHeight
}

func lastConfigFromBlock(block *common.Block) (uint64, error) {
	if block.Metadata == nil || len(block.Metadata.Metadata) <= int(common.BlockMetadataIndex_LAST_CONFIG) {
		return 0, errors.New("no metadata in block")
	}
	return utils.GetLastConfigIndexFromBlock(block)
}


func (ci *ChainInspector) Close() {
	ci.Puller.Close()
}



func (ci *ChainInspector) Channels() []string {
	channels := make(map[string]struct{})
	lastConfigBlockNum := ci.LastConfigBlock.Header.Number
	var block *common.Block
	var prevHash []byte
	for seq := uint64(1); seq < lastConfigBlockNum; seq++ {
		block = ci.Puller.PullBlock(seq)
		ci.validateHashPointer(block, prevHash)
		channel, err := IsNewChannelBlock(block)
		if err != nil {
			
			
			ci.Logger.Panic("Failed classifying block", seq, ":", err)
			continue
		}
		
		prevHash = block.Header.Hash()
		if channel == "" {
			ci.Logger.Info("Block", seq, "doesn't contain a new channel")
			continue
		}
		ci.Logger.Info("Block", seq, "contains channel", channel)
		channels[channel] = struct{}{}
	}
	
	
	
	
	
	last2Blocks := []*common.Block{block, ci.LastConfigBlock}
	if err := VerifyBlockHash(1, last2Blocks); err != nil {
		ci.Logger.Panic("System channel pulled doesn't match the boot last config block:", err)
	}

	return flattenChannelMap(channels)
}

func (ci *ChainInspector) validateHashPointer(block *common.Block, prevHash []byte) {
	if prevHash == nil {
		return
	}
	if bytes.Equal(block.Header.PreviousHash, prevHash) {
		return
	}
	ci.Logger.Panicf("Claimed previous hash of block %d is %x but actual previous hash is %x",
		block.Header.Number, block.Header.PreviousHash, prevHash)
}

func flattenChannelMap(m map[string]struct{}) []string {
	var res []string
	for channel := range m {
		res = append(res, channel)
	}
	return res
}



func IsNewChannelBlock(block *common.Block) (string, error) {
	if block == nil {
		return "", errors.New("nil block")
	}
	env, err := utils.ExtractEnvelope(block, 0)
	if err != nil {
		return "", err
	}
	payload, err := utils.ExtractPayload(env)
	if err != nil {
		return "", err
	}
	if payload.Header == nil {
		return "", errors.New("nil header in payload")
	}
	chdr, err := utils.UnmarshalChannelHeader(payload.Header.ChannelHeader)
	if err != nil {
		return "", err
	}
	
	if common.HeaderType(chdr.Type) != common.HeaderType_ORDERER_TRANSACTION {
		return "", nil
	}
	systemChannelName := chdr.ChannelId
	innerEnvelope, err := utils.UnmarshalEnvelope(payload.Data)
	if err != nil {
		return "", err
	}
	innerPayload, err := utils.UnmarshalPayload(innerEnvelope.Payload)
	if err != nil {
		return "", err
	}
	if innerPayload.Header == nil {
		return "", errors.New("inner payload's header is nil")
	}
	chdr, err = utils.UnmarshalChannelHeader(innerPayload.Header.ChannelHeader)
	if err != nil {
		return "", err
	}
	
	if common.HeaderType(chdr.Type) != common.HeaderType_CONFIG {
		return "", nil
	}
	
	if chdr.ChannelId == systemChannelName {
		return "", nil
	}
	return chdr.ChannelId, nil
}
