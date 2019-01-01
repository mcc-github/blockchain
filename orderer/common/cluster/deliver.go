/*
Copyright IBM Corp. 2018 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package cluster

import (
	"context"
	"math"
	"math/rand"
	"reflect"
	"sync"
	"sync/atomic"
	"time"

	"github.com/mcc-github/blockchain/common/crypto"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/common/util"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/orderer"
	"github.com/mcc-github/blockchain/protos/utils"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)



type BlockPuller struct {
	
	MaxTotalBufferBytes int
	Signer              crypto.LocalSigner
	TLSCert             []byte
	Channel             string
	FetchTimeout        time.Duration
	RetryTimeout        time.Duration
	Logger              *flogging.FabricLogger
	Dialer              Dialer
	VerifyBlockSequence BlockSequenceVerifier
	Endpoints           []string
	
	stream       *ImpatientStream
	blockBuff    []*common.Block
	latestSeq    uint64
	endpoint     string
	conn         *grpc.ClientConn
	cancelStream func()
}



func (p *BlockPuller) Clone() *BlockPuller {
	
	copy := *p
	
	copy.stream = nil
	copy.blockBuff = nil
	copy.latestSeq = 0
	copy.endpoint = ""
	copy.conn = nil
	copy.cancelStream = nil
	return &copy
}



func (p *BlockPuller) Close() {
	if p.cancelStream != nil {
		p.cancelStream()
	}
	p.cancelStream = nil

	if p.conn != nil {
		p.conn.Close()
	}
	p.conn = nil
	p.endpoint = ""
	p.latestSeq = 0
}



func (p *BlockPuller) PullBlock(seq uint64) *common.Block {
	for {
		block := p.tryFetchBlock(seq)
		if block != nil {
			return block
		}
	}
}


func (p *BlockPuller) HeightsByEndpoints() (map[string]uint64, error) {
	endpointsInfo := p.probeEndpoints(0)
	res := make(map[string]uint64)
	for endpoint, endpointInfo := range endpointsInfo.byEndpoints() {
		endpointInfo.conn.Close()
		res[endpoint] = endpointInfo.lastBlockSeq + 1
	}
	p.Logger.Info("Returning the heights of OSNs mapped by endpoints", res)
	return res, endpointsInfo.err
}

func (p *BlockPuller) tryFetchBlock(seq uint64) *common.Block {
	var reConnected bool
	for p.isDisconnected() {
		reConnected = true
		p.connectToSomeEndpoint(seq)
		if p.isDisconnected() {
			time.Sleep(p.RetryTimeout)
		}
	}

	block := p.popBlock(seq)
	if block != nil {
		return block
	}
	
	
	if err := p.pullBlocks(seq, reConnected); err != nil {
		p.Logger.Errorf("Failed pulling blocks: %v", err)
		
		p.Close()
		
		if len(p.blockBuff) > 0 {
			return p.blockBuff[0]
		}
		return nil
	}

	if err := p.VerifyBlockSequence(p.blockBuff); err != nil {
		p.Close()
		p.Logger.Errorf("Failed verifying received blocks: %v", err)
		return nil
	}

	
	return p.popBlock(seq)
}

func (p *BlockPuller) setCancelStreamFunc(f func()) {
	p.cancelStream = f
}

func (p *BlockPuller) pullBlocks(seq uint64, reConnected bool) error {
	env, err := p.seekNextEnvelope(seq)
	if err != nil {
		p.Logger.Errorf("Failed creating seek envelope: %v", err)
		return err
	}

	stream, err := p.obtainStream(reConnected, env, seq)
	if err != nil {
		return err
	}

	var totalSize int
	p.blockBuff = nil
	nextExpectedSequence := seq
	for totalSize < p.MaxTotalBufferBytes && nextExpectedSequence <= p.latestSeq {
		resp, err := stream.Recv()
		if err != nil {
			p.Logger.Errorf("Failed receiving next block from %s: %v", p.endpoint, err)
			return err
		}

		block, err := extractBlockFromResponse(resp)
		if err != nil {
			p.Logger.Errorf("Received a bad block from %s: %v", p.endpoint, err)
			return err
		}
		seq := block.Header.Number
		if seq != nextExpectedSequence {
			p.Logger.Errorf("Expected to receive sequence %d but got %d instead", nextExpectedSequence, seq)
			return errors.Errorf("got unexpected sequence from %s - (%d) instead of (%d)", p.endpoint, seq, nextExpectedSequence)
		}
		size := blockSize(block)
		totalSize += size
		p.blockBuff = append(p.blockBuff, block)
		nextExpectedSequence++
		p.Logger.Infof("Got block %d of size %dKB from %s", seq, size/1024, p.endpoint)
	}
	return nil
}

func (p *BlockPuller) obtainStream(reConnected bool, env *common.Envelope, seq uint64) (*ImpatientStream, error) {
	var stream *ImpatientStream
	var err error
	if reConnected {
		p.Logger.Infof("Sending request for block %d to %s", seq, p.endpoint)
		stream, err = p.requestBlocks(p.endpoint, NewImpatientStream(p.conn, p.FetchTimeout), env)
		if err != nil {
			return nil, err
		}
		
		
		p.stream = stream
	} else {
		
		stream = p.stream
	}

	p.setCancelStreamFunc(stream.cancelFunc)
	return stream, nil
}




func (p *BlockPuller) popBlock(seq uint64) *common.Block {
	if len(p.blockBuff) == 0 {
		return nil
	}
	block, rest := p.blockBuff[0], p.blockBuff[1:]
	p.blockBuff = rest
	
	
	if seq != block.Header.Number {
		p.blockBuff = nil
		return nil
	}
	return block
}

func (p *BlockPuller) isDisconnected() bool {
	return p.conn == nil
}



func (p *BlockPuller) connectToSomeEndpoint(minRequestedSequence uint64) {
	
	
	endpointsInfo := p.probeEndpoints(minRequestedSequence).byEndpoints()
	if len(endpointsInfo) == 0 {
		p.Logger.Warningf("Could not connect to any endpoint of %v", p.Endpoints)
		return
	}

	
	chosenEndpoint := randomEndpoint(endpointsInfo)
	
	for endpoint, endpointInfo := range endpointsInfo {
		if endpoint == chosenEndpoint {
			continue
		}
		endpointInfo.conn.Close()
	}

	p.conn = endpointsInfo[chosenEndpoint].conn
	p.endpoint = chosenEndpoint
	p.latestSeq = endpointsInfo[chosenEndpoint].lastBlockSeq

	p.Logger.Infof("Connected to %s with last block seq of %d", p.endpoint, p.latestSeq)
}



func (p *BlockPuller) probeEndpoints(minRequestedSequence uint64) *endpointInfoBucket {
	endpointsInfo := make(chan *endpointInfo, len(p.Endpoints))

	var wg sync.WaitGroup
	wg.Add(len(p.Endpoints))

	var forbiddenErr uint32
	var unavailableErr uint32

	for _, endpoint := range p.Endpoints {
		go func(endpoint string) {
			defer wg.Done()
			ei, err := p.probeEndpoint(endpoint, minRequestedSequence)
			if err != nil {
				p.Logger.Warningf("Received error of type '%v' from %s", err, endpoint)
				if err == ErrForbidden {
					atomic.StoreUint32(&forbiddenErr, 1)
				}
				if err == ErrServiceUnavailable {
					atomic.StoreUint32(&unavailableErr, 1)
				}
				return
			}
			endpointsInfo <- ei
		}(endpoint)
	}
	wg.Wait()

	close(endpointsInfo)
	eib := &endpointInfoBucket{
		bucket: endpointsInfo,
		logger: p.Logger,
	}

	if unavailableErr == 1 && len(endpointsInfo) == 0 {
		eib.err = ErrServiceUnavailable
	}
	if forbiddenErr == 1 && len(endpointsInfo) == 0 {
		eib.err = ErrForbidden
	}
	return eib
}



func (p *BlockPuller) probeEndpoint(endpoint string, minRequestedSequence uint64) (*endpointInfo, error) {
	conn, err := p.Dialer.Dial(endpoint)
	if err != nil {
		p.Logger.Warningf("Failed connecting to %s: %v", endpoint, err)
		return nil, err
	}

	lastBlockSeq, err := p.fetchLastBlockSeq(minRequestedSequence, endpoint, conn)
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &endpointInfo{conn: conn, lastBlockSeq: lastBlockSeq, endpoint: endpoint}, nil
}


func randomEndpoint(endpointsToHeight map[string]*endpointInfo) string {
	var candidates []string
	for endpoint := range endpointsToHeight {
		candidates = append(candidates, endpoint)
	}

	rand.Seed(time.Now().UnixNano())
	return candidates[rand.Intn(len(candidates))]
}


func (p *BlockPuller) fetchLastBlockSeq(minRequestedSequence uint64, endpoint string, conn *grpc.ClientConn) (uint64, error) {
	env, err := p.seekLastEnvelope()
	if err != nil {
		p.Logger.Errorf("Failed creating seek envelope for %s: %v", endpoint, err)
		return 0, err
	}

	stream, err := p.requestBlocks(endpoint, NewImpatientStream(conn, p.FetchTimeout), env)
	if err != nil {
		return 0, err
	}
	defer stream.abort()

	resp, err := stream.Recv()
	if err != nil {
		p.Logger.Errorf("Failed receiving the latest block from %s: %v", endpoint, err)
		return 0, err
	}

	block, err := extractBlockFromResponse(resp)
	if err != nil {
		p.Logger.Errorf("Received %v from %s: %v", resp, endpoint, err)
		return 0, err
	}
	stream.CloseSend()

	seq := block.Header.Number
	if seq < minRequestedSequence {
		err := errors.Errorf("minimum requested sequence is %d but %s is at sequence %d", minRequestedSequence, endpoint, seq)
		p.Logger.Infof("Skipping pulling from %s: %v", endpoint, err)
		return 0, err
	}

	p.Logger.Infof("%s is at block sequence of %d", endpoint, seq)
	return block.Header.Number, nil
}




func (p *BlockPuller) requestBlocks(endpoint string, newStream ImpatientStreamCreator, env *common.Envelope) (*ImpatientStream, error) {
	stream, err := newStream()
	if err != nil {
		p.Logger.Warningf("Failed establishing deliver stream with %s", endpoint)
		return nil, err
	}

	if err := stream.Send(env); err != nil {
		p.Logger.Errorf("Failed sending seek envelope to %s: %v", endpoint, err)
		stream.abort()
		return nil, err
	}
	return stream, nil
}

func extractBlockFromResponse(resp *orderer.DeliverResponse) (*common.Block, error) {
	switch t := resp.Type.(type) {
	case *orderer.DeliverResponse_Block:
		block := t.Block
		if block == nil {
			return nil, errors.New("block is nil")
		}
		if block.Data == nil {
			return nil, errors.New("block data is nil")
		}
		if block.Header == nil {
			return nil, errors.New("block header is nil")
		}
		if block.Metadata == nil || len(block.Metadata.Metadata) == 0 {
			return nil, errors.New("block metadata is empty")
		}
		return block, nil
	case *orderer.DeliverResponse_Status:
		if t.Status == common.Status_FORBIDDEN {
			return nil, ErrForbidden
		}
		if t.Status == common.Status_SERVICE_UNAVAILABLE {
			return nil, ErrServiceUnavailable
		}
		return nil, errors.Errorf("faulty node, received: %v", resp)
	default:
		return nil, errors.Errorf("response is of type %v, but expected a block", reflect.TypeOf(resp.Type))
	}
}

func (p *BlockPuller) seekLastEnvelope() (*common.Envelope, error) {
	return utils.CreateSignedEnvelopeWithTLSBinding(
		common.HeaderType_DELIVER_SEEK_INFO,
		p.Channel,
		p.Signer,
		last(),
		int32(0),
		uint64(0),
		util.ComputeSHA256(p.TLSCert),
	)
}

func (p *BlockPuller) seekNextEnvelope(startSeq uint64) (*common.Envelope, error) {
	return utils.CreateSignedEnvelopeWithTLSBinding(
		common.HeaderType_DELIVER_SEEK_INFO,
		p.Channel,
		p.Signer,
		nextSeekInfo(startSeq),
		int32(0),
		uint64(0),
		util.ComputeSHA256(p.TLSCert),
	)
}

func last() *orderer.SeekInfo {
	return &orderer.SeekInfo{
		Start:    &orderer.SeekPosition{Type: &orderer.SeekPosition_Newest{Newest: &orderer.SeekNewest{}}},
		Stop:     &orderer.SeekPosition{Type: &orderer.SeekPosition_Specified{Specified: &orderer.SeekSpecified{Number: math.MaxUint64}}},
		Behavior: orderer.SeekInfo_BLOCK_UNTIL_READY,
	}
}

func nextSeekInfo(startSeq uint64) *orderer.SeekInfo {
	return &orderer.SeekInfo{
		Start:    &orderer.SeekPosition{Type: &orderer.SeekPosition_Specified{Specified: &orderer.SeekSpecified{Number: startSeq}}},
		Stop:     &orderer.SeekPosition{Type: &orderer.SeekPosition_Specified{Specified: &orderer.SeekSpecified{Number: math.MaxUint64}}},
		Behavior: orderer.SeekInfo_BLOCK_UNTIL_READY,
	}
}

func blockSize(block *common.Block) int {
	return len(utils.MarshalOrPanic(block))
}

type endpointInfo struct {
	endpoint     string
	conn         *grpc.ClientConn
	lastBlockSeq uint64
}

type endpointInfoBucket struct {
	bucket <-chan *endpointInfo
	logger *flogging.FabricLogger
	err    error
}

func (eib endpointInfoBucket) byEndpoints() map[string]*endpointInfo {
	infoByEndpoints := make(map[string]*endpointInfo)
	for endpointInfo := range eib.bucket {
		if _, exists := infoByEndpoints[endpointInfo.endpoint]; exists {
			eib.logger.Warningf("Duplicate endpoint found(%s), skipping it", endpointInfo.endpoint)
			endpointInfo.conn.Close()
			continue
		}
		infoByEndpoints[endpointInfo.endpoint] = endpointInfo
	}
	return infoByEndpoints
}


type ImpatientStreamCreator func() (*ImpatientStream, error)


type ImpatientStream struct {
	waitTimeout time.Duration
	orderer.AtomicBroadcast_DeliverClient
	cancelFunc func()
}

func (stream *ImpatientStream) abort() {
	stream.cancelFunc()
}



func (stream *ImpatientStream) Recv() (*orderer.DeliverResponse, error) {
	
	timeout := time.NewTimer(stream.waitTimeout)
	defer timeout.Stop()

	responseChan := make(chan errorAndResponse, 1)

	
	
	var receive sync.WaitGroup
	receive.Add(1)
	defer receive.Wait()

	go func() {
		defer receive.Done()
		resp, err := stream.AtomicBroadcast_DeliverClient.Recv()
		responseChan <- errorAndResponse{err: err, resp: resp}
	}()

	select {
	case <-timeout.C:
		stream.cancelFunc()
		return nil, errors.Errorf("didn't receive a response within %v", stream.waitTimeout)
	case respAndErr := <-responseChan:
		return respAndErr.resp, respAndErr.err
	}
}


func NewImpatientStream(conn *grpc.ClientConn, waitTimeout time.Duration) ImpatientStreamCreator {
	return func() (*ImpatientStream, error) {
		abc := orderer.NewAtomicBroadcastClient(conn)
		ctx, cancel := context.WithCancel(context.Background())

		stream, err := abc.Deliver(ctx)
		if err != nil {
			cancel()
			return nil, err
		}

		once := &sync.Once{}
		return &ImpatientStream{
			waitTimeout: waitTimeout,
			
			
			cancelFunc: func() {
				once.Do(cancel)
			},
			AtomicBroadcast_DeliverClient: stream,
		}, nil
	}
}

type errorAndResponse struct {
	err  error
	resp *orderer.DeliverResponse
}
