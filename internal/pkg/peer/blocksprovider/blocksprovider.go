/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package blocksprovider

import (
	"context"
	"crypto/x509"
	"math"
	"time"

	"github.com/mcc-github/blockchain-protos-go/common"
	"github.com/mcc-github/blockchain-protos-go/gossip"
	"github.com/mcc-github/blockchain-protos-go/orderer"
	"github.com/mcc-github/blockchain/common/flogging"
	gossipcommon "github.com/mcc-github/blockchain/gossip/common"
	"github.com/mcc-github/blockchain/internal/pkg/identity"
	"github.com/mcc-github/blockchain/internal/pkg/peer/orderers"
	"github.com/mcc-github/blockchain/protoutil"

	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type sleeper struct {
	sleep func(time.Duration)
}

func (s sleeper) Sleep(d time.Duration, doneC chan struct{}) {
	if s.sleep == nil {
		timer := time.NewTimer(d)
		select {
		case <-timer.C:
		case <-doneC:
			timer.Stop()
		}
		return
	}
	s.sleep(d)
}




type LedgerInfo interface {
	
	LedgerHeight() (uint64, error)
}




type GossipServiceAdapter interface {
	
	AddPayload(chainID string, payload *gossip.Payload) error

	
	Gossip(msg *gossip.GossipMessage)
}


type BlockVerifier interface {
	VerifyBlock(channelID gossipcommon.ChannelID, blockNum uint64, block *common.Block) error
}


type OrdererConnectionSource interface {
	RandomEndpoint() (*orderers.Endpoint, error)
}


type Dialer interface {
	Dial(address string, certPool *x509.CertPool) (*grpc.ClientConn, error)
}


type DeliverStreamer interface {
	Deliver(context.Context, *grpc.ClientConn) (orderer.AtomicBroadcast_DeliverClient, error)
}


type Deliverer struct {
	ChannelID       string
	Gossip          GossipServiceAdapter
	Ledger          LedgerInfo
	BlockVerifier   BlockVerifier
	Dialer          Dialer
	Orderers        OrdererConnectionSource
	DoneC           chan struct{}
	Signer          identity.SignerSerializer
	DeliverStreamer DeliverStreamer
	Logger          *flogging.FabricLogger

	MaxRetryDelay     time.Duration
	InitialRetryDelay time.Duration
	MaxRetryDuration  time.Duration

	
	TLSCertHash []byte 

	sleeper sleeper
}

const backoffExponentBase = 1.2



func (d *Deliverer) DeliverBlocks() {
	failureCounter := 0
	totalDuration := time.Duration(0)

	
	
	
	
	maxFailures := int(math.Log(float64(d.MaxRetryDelay)/float64(d.InitialRetryDelay)) / math.Log(backoffExponentBase))
	for {
		select {
		case <-d.DoneC:
			return
		default:
		}

		if failureCounter > 0 {
			var sleepDuration time.Duration
			if failureCounter-1 > maxFailures {
				sleepDuration = d.MaxRetryDelay
			} else {
				sleepDuration = time.Duration(math.Pow(1.2, float64(failureCounter-1))*100) * time.Millisecond
			}
			totalDuration += sleepDuration
			if totalDuration > d.MaxRetryDuration {
				d.Logger.Warningf("attempted to retry block delivery for more than %v, giving up", d.MaxRetryDuration)
				return
			}
			d.sleeper.Sleep(sleepDuration, d.DoneC)
		}

		ledgerHeight, err := d.Ledger.LedgerHeight()
		if err != nil {
			d.Logger.Error("Did not return ledger height, something is critically wrong", err)
			return
		}

		seekInfoEnv, err := d.createSeekInfo(ledgerHeight)
		if err != nil {
			d.Logger.Error("Could not create a signed Deliver SeekInfo message, something is critically wrong", err)
			return
		}

		deliverClient, endpoint, cancel, err := d.connect(seekInfoEnv)
		if err != nil {
			d.Logger.Warningf("Could not connect to ordering service: %s", err)
			failureCounter++
			continue
		}

		recv := make(chan *orderer.DeliverResponse)
		go func() {
			for {
				resp, err := deliverClient.Recv()
				if err != nil {
					d.Logger.Warningf("Encountered an error reading from deliver stream: %s", err)
					close(recv)
					return
				}
				select {
				case recv <- resp:
				case <-d.DoneC:
					close(recv)
					return
				}
			}
		}()

	RecvLoop: 
		for {
			select {
			case <-endpoint.Refreshed:
				d.Logger.Infof("Ordering endpoints have been refreshed, disconnecting from deliver to reconnect using updated endpoints")
				break RecvLoop
			case response, ok := <-recv:
				if !ok {
					d.Logger.Warningf("Orderer hung up without sending status")
					failureCounter++
					break RecvLoop
				}
				err = d.processMsg(response)
				if err != nil {
					d.Logger.Warningf("Got error while attempting to receive blocks: %v", err)
					failureCounter++
					break RecvLoop
				}
				failureCounter = 0
			case <-d.DoneC:
				break RecvLoop
			}
		}

		
		cancel()
		<-recv
	}
}

func (d *Deliverer) processMsg(msg *orderer.DeliverResponse) error {
	switch t := msg.Type.(type) {
	case *orderer.DeliverResponse_Status:
		if t.Status == common.Status_SUCCESS {
			return errors.Errorf("received success for a seek that should never complete")
		}

		return errors.Errorf("received bad status %v from orderer", t.Status)
	case *orderer.DeliverResponse_Block:
		blockNum := t.Block.Header.Number
		if err := d.BlockVerifier.VerifyBlock(gossipcommon.ChannelID(d.ChannelID), blockNum, t.Block); err != nil {
			return errors.WithMessage(err, "block from orderer could not be verified")
		}

		marshaledBlock, err := proto.Marshal(t.Block)
		if err != nil {
			return errors.WithMessage(err, "block from orderer could not be re-marshaled")
		}

		
		payload := &gossip.Payload{
			Data:   marshaledBlock,
			SeqNum: blockNum,
		}

		
		gossipMsg := &gossip.GossipMessage{
			Nonce:   0,
			Tag:     gossip.GossipMessage_CHAN_AND_ORG,
			Channel: []byte(d.ChannelID),
			Content: &gossip.GossipMessage_DataMsg{
				DataMsg: &gossip.DataMessage{
					Payload: payload,
				},
			},
		}

		d.Logger.Debugf("Adding payload to local buffer, blockNum = [%d]", blockNum)
		
		if err := d.Gossip.AddPayload(d.ChannelID, payload); err != nil {
			d.Logger.Warningf("Block [%d] received from ordering service wasn't added to payload buffer: %v", blockNum, err)
			return errors.WithMessage(err, "could not add block as payload")
		}

		
		d.Logger.Debugf("Gossiping block [%d]", blockNum)
		d.Gossip.Gossip(gossipMsg)
		return nil
	default:
		d.Logger.Warningf("Received unknown: %v", t)
		return errors.Errorf("unknown message type '%T'", msg.Type)
	}
}


func (d *Deliverer) Stop() {
	close(d.DoneC)
}

func (d *Deliverer) connect(seekInfoEnv *common.Envelope) (orderer.AtomicBroadcast_DeliverClient, *orderers.Endpoint, func(), error) {
	endpoint, err := d.Orderers.RandomEndpoint()
	if err != nil {
		return nil, nil, nil, errors.WithMessage(err, "could not get orderer endpoints")
	}

	conn, err := d.Dialer.Dial(endpoint.Address, endpoint.CertPool)
	if err != nil {
		return nil, nil, nil, errors.WithMessagef(err, "could not dial endpoint '%s'", endpoint.Address)
	}

	ctx, ctxCancel := context.WithCancel(context.Background())

	deliverClient, err := d.DeliverStreamer.Deliver(ctx, conn)
	if err != nil {
		conn.Close()
		ctxCancel()
		return nil, nil, nil, errors.WithMessagef(err, "could not create deliver client to endpoints '%s'", endpoint.Address)
	}

	err = deliverClient.Send(seekInfoEnv)
	if err != nil {
		deliverClient.CloseSend()
		conn.Close()
		ctxCancel()
		return nil, nil, nil, errors.WithMessagef(err, "could not send deliver seek info handshake to '%s'", endpoint.Address)
	}

	return deliverClient, endpoint, func() {
		deliverClient.CloseSend()
		ctxCancel()
		conn.Close()
	}, nil
}

func (d *Deliverer) createSeekInfo(ledgerHeight uint64) (*common.Envelope, error) {
	return protoutil.CreateSignedEnvelopeWithTLSBinding(
		common.HeaderType_DELIVER_SEEK_INFO,
		d.ChannelID,
		d.Signer,
		&orderer.SeekInfo{
			Start: &orderer.SeekPosition{
				Type: &orderer.SeekPosition_Specified{
					Specified: &orderer.SeekSpecified{
						Number: ledgerHeight,
					},
				},
			},
			Stop: &orderer.SeekPosition{
				Type: &orderer.SeekPosition_Specified{
					Specified: &orderer.SeekSpecified{
						Number: math.MaxUint64,
					},
				},
			},
			Behavior: orderer.SeekInfo_BLOCK_UNTIL_READY,
		},
		int32(0),
		uint64(0),
		d.TLSCertHash,
	)
}
