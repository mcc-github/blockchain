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

	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/core/comm"
	"github.com/mcc-github/blockchain/internal/pkg/identity"
	"github.com/mcc-github/blockchain/orderer/common/localconfig"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/pkg/errors"
)

const (
	
	RetryTimeout = time.Second * 10
)


type ChannelPredicate func(channelName string) bool


func AnyChannel(_ string) bool {
	return true
}




func PullerConfigFromTopLevelConfig(
	systemChannel string,
	conf *localconfig.TopLevel,
	tlsKey,
	tlsCert []byte,
	signer identity.SignerSerializer,
) PullerConfig {
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
	
	Channels() []ChannelGenesisBlock
	
	Close()
}


type Replicator struct {
	DoNotPanicIfClusterNotReachable bool
	Filter                          ChannelPredicate
	SystemChannel                   string
	ChannelLister                   ChannelLister
	Logger                          *flogging.FabricLogger
	Puller                          *BlockPuller
	BootBlock                       *common.Block
	AmIPartOfChannel                SelfMembershipPredicate
	LedgerFactory                   LedgerFactory
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



func (r *Replicator) ReplicateChains() []string {
	var replicatedChains []string
	channels := r.discoverChannels()
	pullHints := r.channelsToPull(channels)
	totalChannelCount := len(pullHints.channelsToPull) + len(pullHints.channelsNotToPull)
	r.Logger.Info("Found myself in", len(pullHints.channelsToPull), "channels out of", totalChannelCount, ":", pullHints)

	
	for _, channels := range [][]ChannelGenesisBlock{pullHints.channelsToPull, pullHints.channelsNotToPull} {
		for _, channel := range channels {
			ledger, err := r.LedgerFactory.GetOrCreate(channel.ChannelName)
			if err != nil {
				r.Logger.Panicf("Failed to create a ledger for channel %s: %v", channel.ChannelName, err)
			}
			gb, err := ChannelCreationBlockToGenesisBlock(channel.GenesisBlock)
			if err != nil {
				r.Logger.Panicf("Failed converting channel creation block for channel %s to genesis block: %v",
					channel.ChannelName, err)
			}
			r.appendBlock(gb, ledger, channel.ChannelName)
		}
	}

	for _, channel := range pullHints.channelsToPull {
		err := r.PullChannel(channel.ChannelName)
		if err == nil {
			replicatedChains = append(replicatedChains, channel.ChannelName)
		} else {
			r.Logger.Warningf("Failed pulling channel %s: %v", channel.ChannelName, err)
		}
	}

	
	if err := r.PullChannel(r.SystemChannel); err != nil && err != ErrSkipped {
		r.Logger.Panicf("Failed pulling system channel: %v", err)
	}
	return replicatedChains
}

func (r *Replicator) discoverChannels() []ChannelGenesisBlock {
	r.Logger.Debug("Entering")
	defer r.Logger.Debug("Exiting")
	channels := GenesisBlocks(r.ChannelLister.Channels())
	r.Logger.Info("Discovered", len(channels), "channels:", channels.Names())
	r.ChannelLister.Close()
	return channels
}



func (r *Replicator) PullChannel(channel string) error {
	if !r.Filter(channel) {
		r.Logger.Infof("Channel %s shouldn't be pulled. Skipping it", channel)
		return ErrSkipped
	}
	r.Logger.Info("Pulling channel", channel)
	puller := r.Puller.Clone()
	defer puller.Close()
	puller.Channel = channel

	ledger, err := r.LedgerFactory.GetOrCreate(channel)
	if err != nil {
		r.Logger.Panicf("Failed to create a ledger for channel %s: %v", channel, err)
	}

	endpoint, latestHeight, _ := latestHeightAndEndpoint(puller)
	if endpoint == "" {
		return errors.Errorf("failed obtaining the latest block for channel %s", channel)
	}
	r.Logger.Info("Latest block height for channel", channel, "is", latestHeight)
	
	
	
	if channel == r.SystemChannel && latestHeight-1 < r.BootBlock.Header.Number {
		return errors.Errorf("latest height found among system channel(%s) orderers is %d, but the boot block's "+
			"sequence is %d", r.SystemChannel, latestHeight, r.BootBlock.Header.Number)
	}
	return r.pullChannelBlocks(channel, puller, latestHeight, ledger)
}

func (r *Replicator) pullChannelBlocks(channel string, puller *BlockPuller, latestHeight uint64, ledger LedgerWriter) error {
	nextBlockToPull := ledger.Height()
	if nextBlockToPull == latestHeight {
		r.Logger.Infof("Latest height found (%d) is equal to our height, skipping pulling channel %s", latestHeight, channel)
		return nil
	}
	
	nextBlock := puller.PullBlock(nextBlockToPull)
	if nextBlock == nil {
		return ErrRetryCountExhausted
	}
	r.appendBlock(nextBlock, ledger, channel)
	actualPrevHash := protoutil.BlockHeaderHash(nextBlock.Header)

	for seq := uint64(nextBlockToPull + 1); seq < latestHeight; seq++ {
		block := puller.PullBlock(seq)
		if block == nil {
			return ErrRetryCountExhausted
		}
		reportedPrevHash := block.Header.PreviousHash
		if !bytes.Equal(reportedPrevHash, actualPrevHash) {
			return errors.Errorf("block header mismatch on sequence %d, expected %x, got %x",
				block.Header.Number, actualPrevHash, reportedPrevHash)
		}
		actualPrevHash = protoutil.BlockHeaderHash(block.Header)
		if channel == r.SystemChannel && block.Header.Number == r.BootBlock.Header.Number {
			r.compareBootBlockWithSystemChannelLastConfigBlock(block)
			r.appendBlock(block, ledger, channel)
			
			return nil
		}
		r.appendBlock(block, ledger, channel)
	}
	return nil
}

func (r *Replicator) appendBlock(block *common.Block, ledger LedgerWriter, channel string) {
	height := ledger.Height()
	if height > block.Header.Number {
		r.Logger.Infof("Skipping commit of block [%d] for channel %s because height is at %d", block.Header.Number, channel, height)
		return
	}
	if err := ledger.Append(block); err != nil {
		r.Logger.Panicf("Failed to write block [%d]: %v", block.Header.Number, err)
	}
	r.Logger.Infof("Committed block [%d] for channel %s", block.Header.Number, channel)
}

func (r *Replicator) compareBootBlockWithSystemChannelLastConfigBlock(block *common.Block) {
	
	block.Header.DataHash = protoutil.BlockDataHash(block.Data)

	bootBlockHash := protoutil.BlockHeaderHash(r.BootBlock.Header)
	retrievedBlockHash := protoutil.BlockHeaderHash(block.Header)
	if bytes.Equal(bootBlockHash, retrievedBlockHash) {
		return
	}
	r.Logger.Panicf("Block header mismatch on last system channel block, expected %s, got %s",
		hex.EncodeToString(bootBlockHash), hex.EncodeToString(retrievedBlockHash))
}

type channelPullHints struct {
	channelsToPull    []ChannelGenesisBlock
	channelsNotToPull []ChannelGenesisBlock
}

func (r *Replicator) channelsToPull(channels GenesisBlocks) channelPullHints {
	r.Logger.Info("Evaluating channels to pull:", channels.Names())

	
	verifier := r.Puller.VerifyBlockSequence
	
	defer func() {
		r.Puller.VerifyBlockSequence = verifier
	}()
	
	r.Puller.VerifyBlockSequence = func(blocks []*common.Block, channel string) error {
		return nil
	}

	var channelsNotToPull []ChannelGenesisBlock
	var channelsToPull []ChannelGenesisBlock
	for _, channel := range channels {
		r.Logger.Info("Probing whether I should pull channel", channel.ChannelName)
		puller := r.Puller.Clone()
		puller.Channel = channel.ChannelName
		
		
		bufferSize := puller.MaxTotalBufferBytes
		puller.MaxTotalBufferBytes = 1
		err := Participant(puller, r.AmIPartOfChannel)
		puller.Close()
		
		puller.MaxTotalBufferBytes = bufferSize
		if err == ErrNotInChannel || err == ErrForbidden {
			r.Logger.Infof("I do not belong to channel %s or am forbidden pulling it (%v), skipping chain retrieval", channel.ChannelName, err)
			channelsNotToPull = append(channelsNotToPull, channel)
			continue
		}
		if err == ErrServiceUnavailable {
			r.Logger.Infof("All orderers in the system channel are either down,"+
				"or do not service channel %s (%v), skipping chain retrieval", channel.ChannelName, err)
			channelsNotToPull = append(channelsNotToPull, channel)
			continue
		}
		if err == ErrRetryCountExhausted {
			r.Logger.Warningf("Could not obtain blocks needed for classifying whether I am in the channel,"+
				"skipping the retrieval of the chan %s", channel.ChannelName)
			channelsNotToPull = append(channelsNotToPull, channel)
			continue
		}
		if err != nil {
			if !r.DoNotPanicIfClusterNotReachable {
				r.Logger.Panicf("Failed classifying whether I belong to channel %s: %v, skipping chain retrieval", channel.ChannelName, err)
			}
			continue
		}
		r.Logger.Infof("I need to pull channel %s", channel.ChannelName)
		channelsToPull = append(channelsToPull, channel)
	}
	return channelPullHints{
		channelsToPull:    channelsToPull,
		channelsNotToPull: channelsNotToPull,
	}
}


type PullerConfig struct {
	TLSKey              []byte
	TLSCert             []byte
	Timeout             time.Duration
	Signer              identity.SignerSerializer
	Channel             string
	MaxTotalBufferBytes int
}




type VerifierRetriever interface {
	
	RetrieveVerifier(channel string) BlockVerifier
}


func BlockPullerFromConfigBlock(conf PullerConfig, block *common.Block, verifierRetriever VerifierRetriever) (*BlockPuller, error) {
	if block == nil {
		return nil, errors.New("nil block")
	}

	endpoints, err := EndpointconfigFromConfigBlock(block)
	if err != nil {
		return nil, err
	}

	clientConf := comm.ClientConfig{
		Timeout: conf.Timeout,
		SecOpts: &comm.SecureOptions{
			Certificate:       conf.TLSCert,
			Key:               conf.TLSKey,
			RequireClientCert: true,
			UseTLS:            true,
		},
	}

	dialer := &StandardDialer{
		Config: clientConf.Clone(),
	}

	tlsCertAsDER, _ := pem.Decode(conf.TLSCert)
	if tlsCertAsDER == nil {
		return nil, errors.Errorf("unable to decode TLS certificate PEM: %s", base64.StdEncoding.EncodeToString(conf.TLSCert))
	}

	return &BlockPuller{
		Logger:  flogging.MustGetLogger("orderer.common.cluster.replication"),
		Dialer:  dialer,
		TLSCert: tlsCertAsDER.Bytes,
		VerifyBlockSequence: func(blocks []*common.Block, channel string) error {
			verifier := verifierRetriever.RetrieveVerifier(channel)
			if verifier == nil {
				return errors.Errorf("couldn't acquire verifier for channel %s", channel)
			}
			return VerifyBlocks(blocks, verifier)
		},
		MaxTotalBufferBytes: conf.MaxTotalBufferBytes,
		Endpoints:           endpoints,
		RetryTimeout:        RetryTimeout,
		FetchTimeout:        conf.Timeout,
		Channel:             conf.Channel,
		Signer:              conf.Signer,
	}, nil
}


type NoopBlockVerifier struct{}


func (*NoopBlockVerifier) VerifyBlockSignature(sd []*protoutil.SignedData, config *common.ConfigEnvelope) error {
	return nil
}




type ChainPuller interface {
	
	PullBlock(seq uint64) *common.Block

	
	HeightsByEndpoints() (map[string]uint64, error)

	
	Close()
}


type ChainInspector struct {
	Logger          *flogging.FabricLogger
	Puller          ChainPuller
	LastConfigBlock *common.Block
}


var ErrSkipped = errors.New("skipped")


var ErrForbidden = errors.New("forbidden pulling the channel")


var ErrServiceUnavailable = errors.New("service unavailable")


var ErrNotInChannel = errors.New("not in the channel")

var ErrRetryCountExhausted = errors.New("retry attempts exhausted")


type SelfMembershipPredicate func(configBlock *common.Block) error










func Participant(puller ChainPuller, analyzeLastConfBlock SelfMembershipPredicate) error {
	lastConfigBlock, err := PullLastConfigBlock(puller)
	if err != nil {
		return err
	}
	return analyzeLastConfBlock(lastConfigBlock)
}


func PullLastConfigBlock(puller ChainPuller) (*common.Block, error) {
	endpoint, latestHeight, err := latestHeightAndEndpoint(puller)
	if err != nil {
		return nil, err
	}
	if endpoint == "" {
		return nil, ErrRetryCountExhausted
	}
	lastBlock := puller.PullBlock(latestHeight - 1)
	if lastBlock == nil {
		return nil, ErrRetryCountExhausted
	}
	lastConfNumber, err := lastConfigFromBlock(lastBlock)
	if err != nil {
		return nil, err
	}
	
	
	
	puller.Close()
	lastConfigBlock := puller.PullBlock(lastConfNumber)
	if lastConfigBlock == nil {
		return nil, ErrRetryCountExhausted
	}
	return lastConfigBlock, nil
}

func latestHeightAndEndpoint(puller ChainPuller) (string, uint64, error) {
	var maxHeight uint64
	var mostUpToDateEndpoint string
	heightsByEndpoints, err := puller.HeightsByEndpoints()
	if err != nil {
		return "", 0, err
	}
	for endpoint, height := range heightsByEndpoints {
		if height >= maxHeight {
			maxHeight = height
			mostUpToDateEndpoint = endpoint
		}
	}
	return mostUpToDateEndpoint, maxHeight, nil
}

func lastConfigFromBlock(block *common.Block) (uint64, error) {
	if block.Metadata == nil || len(block.Metadata.Metadata) <= int(common.BlockMetadataIndex_LAST_CONFIG) {
		return 0, errors.New("no metadata in block")
	}
	return protoutil.GetLastConfigIndexFromBlock(block)
}


func (ci *ChainInspector) Close() {
	ci.Puller.Close()
}


type ChannelGenesisBlock struct {
	ChannelName  string
	GenesisBlock *common.Block
}


type GenesisBlocks []ChannelGenesisBlock


func (gbs GenesisBlocks) Names() []string {
	var res []string
	for _, gb := range gbs {
		res = append(res, gb.ChannelName)
	}
	return res
}




func (ci *ChainInspector) Channels() []ChannelGenesisBlock {
	channels := make(map[string]ChannelGenesisBlock)
	lastConfigBlockNum := ci.LastConfigBlock.Header.Number
	var block *common.Block
	var prevHash []byte
	for seq := uint64(0); seq < lastConfigBlockNum; seq++ {
		block = ci.Puller.PullBlock(seq)
		if block == nil {
			ci.Logger.Panicf("Failed pulling block [%d] from the system channel", seq)
		}
		ci.validateHashPointer(block, prevHash)
		channel, err := IsNewChannelBlock(block)
		if err != nil {
			
			
			ci.Logger.Panicf("Failed classifying block [%d]: %s", seq, err)
			continue
		}
		
		prevHash = protoutil.BlockHeaderHash(block.Header)
		if channel == "" {
			ci.Logger.Info("Block", seq, "doesn't contain a new channel")
			continue
		}
		ci.Logger.Info("Block", seq, "contains channel", channel)
		channels[channel] = ChannelGenesisBlock{
			ChannelName:  channel,
			GenesisBlock: block,
		}
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
	ci.Logger.Panicf("Claimed previous hash of block [%d] is %x but actual previous hash is %x",
		block.Header.Number, block.Header.PreviousHash, prevHash)
}

func flattenChannelMap(m map[string]ChannelGenesisBlock) []ChannelGenesisBlock {
	var res []ChannelGenesisBlock
	for _, csb := range m {
		res = append(res, csb)
	}
	return res
}


func ChannelCreationBlockToGenesisBlock(block *common.Block) (*common.Block, error) {
	if block == nil {
		return nil, errors.New("nil block")
	}
	env, err := protoutil.ExtractEnvelope(block, 0)
	if err != nil {
		return nil, err
	}
	payload, err := protoutil.ExtractPayload(env)
	if err != nil {
		return nil, err
	}
	block.Data.Data = [][]byte{payload.Data}
	block.Header.DataHash = protoutil.BlockDataHash(block.Data)
	block.Header.Number = 0
	block.Header.PreviousHash = nil
	metadata := &common.BlockMetadata{
		Metadata: make([][]byte, 4),
	}
	block.Metadata = metadata
	metadata.Metadata[common.BlockMetadataIndex_LAST_CONFIG] = protoutil.MarshalOrPanic(&common.Metadata{
		Value: protoutil.MarshalOrPanic(&common.LastConfig{Index: 0}),
		
		
	})
	return block, nil
}



func IsNewChannelBlock(block *common.Block) (string, error) {
	if block == nil {
		return "", errors.New("nil block")
	}
	env, err := protoutil.ExtractEnvelope(block, 0)
	if err != nil {
		return "", err
	}
	payload, err := protoutil.ExtractPayload(env)
	if err != nil {
		return "", err
	}
	if payload.Header == nil {
		return "", errors.New("nil header in payload")
	}
	chdr, err := protoutil.UnmarshalChannelHeader(payload.Header.ChannelHeader)
	if err != nil {
		return "", err
	}
	
	if common.HeaderType(chdr.Type) != common.HeaderType_ORDERER_TRANSACTION {
		return "", nil
	}
	systemChannelName := chdr.ChannelId
	innerEnvelope, err := protoutil.UnmarshalEnvelope(payload.Data)
	if err != nil {
		return "", err
	}
	innerPayload, err := protoutil.UnmarshalPayload(innerEnvelope.Payload)
	if err != nil {
		return "", err
	}
	if innerPayload.Header == nil {
		return "", errors.New("inner payload's header is nil")
	}
	chdr, err = protoutil.UnmarshalChannelHeader(innerPayload.Header.ChannelHeader)
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
