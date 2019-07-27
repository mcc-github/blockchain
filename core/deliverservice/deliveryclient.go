/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package deliverservice

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/core/comm"
	"github.com/mcc-github/blockchain/core/deliverservice/blocksprovider"
	"github.com/mcc-github/blockchain/gossip/api"
	"github.com/mcc-github/blockchain/internal/pkg/identity"
	"github.com/mcc-github/blockchain/protos/orderer"
	"google.golang.org/grpc"
)

var logger = flogging.MustGetLogger("deliveryClient")



type DeliverService interface {
	
	
	
	StartDeliverForChannel(chainID string, ledgerInfo blocksprovider.LedgerInfo, finalizer func()) error

	
	
	StopDeliverForChannel(chainID string) error

	
	UpdateEndpoints(chainID string, endpoints []string) error

	
	Stop()
}




type deliverServiceImpl struct {
	conf           *Config
	blockProviders map[string]blocksprovider.BlocksProvider
	lock           sync.RWMutex
	stopping       bool
}





type Config struct {
	
	ConnFactory func(channelID string) func(endpoint string, connectionTimeout time.Duration) (*grpc.ClientConn, error)
	
	ABCFactory func(*grpc.ClientConn) orderer.AtomicBroadcastClient
	
	
	CryptoSvc api.MessageCryptoService
	
	
	Gossip blocksprovider.GossipServiceAdapter
	
	Endpoints []string
	
	Signer identity.SignerSerializer
	
	CredentialSupport *comm.CredentialSupport
	
	DeliverClientDialOpts []grpc.DialOption
	
	
	DeliverServiceConfig *DeliverServiceConfig
}


type ConnectionCriteria struct {
	
	OrdererEndpoints []string
	
	Organizations []string
	
	OrdererEndpointsByOrg map[string][]string
}

func (cc ConnectionCriteria) toEndpointCriteria() []comm.EndpointCriteria {
	var res []comm.EndpointCriteria

	
	for _, org := range cc.Organizations {
		endpoints := cc.OrdererEndpointsByOrg[org]
		if len(endpoints) == 0 {
			
			continue
		}

		for _, endpoint := range endpoints {
			res = append(res, comm.EndpointCriteria{
				Organizations: []string{org},
				Endpoint:      endpoint,
			})
		}
	}

	
	if len(res) > 0 {
		return res
	}

	for _, endpoint := range cc.OrdererEndpoints {
		res = append(res, comm.EndpointCriteria{
			Organizations: cc.Organizations,
			Endpoint:      endpoint,
		})
	}

	return res
}





func NewDeliverService(conf *Config) (DeliverService, error) {
	ds := &deliverServiceImpl{
		conf:           conf,
		blockProviders: make(map[string]blocksprovider.BlocksProvider),
	}
	if err := ds.validateConfiguration(); err != nil {
		return nil, err
	}
	return ds, nil
}

func (d *deliverServiceImpl) UpdateEndpoints(chainID string, endpoints []string) error {
	
	
	if bp, ok := d.blockProviders[chainID]; ok {
		
		bp.UpdateOrderingEndpoints(endpoints)
		return nil
	}
	return errors.New(fmt.Sprintf("Channel with %s id was not found", chainID))
}

func (d *deliverServiceImpl) validateConfiguration() error {
	conf := d.conf
	if len(conf.Endpoints) == 0 {
		return errors.New("no endpoints specified")
	}
	if conf.Gossip == nil {
		return errors.New("no gossip provider specified")
	}
	if conf.ABCFactory == nil {
		return errors.New("no AtomicBroadcast factory specified")
	}
	if conf.ConnFactory == nil {
		return errors.New("no connection factory specified")
	}
	if conf.CryptoSvc == nil {
		return errors.New("no crypto service specified")
	}
	if conf.CredentialSupport == nil {
		return errors.New("no credential support specified")
	}
	return nil
}





func (d *deliverServiceImpl) StartDeliverForChannel(chainID string, ledgerInfo blocksprovider.LedgerInfo, finalizer func()) error {
	d.lock.Lock()
	defer d.lock.Unlock()
	if d.stopping {
		errMsg := fmt.Sprintf("Delivery service is stopping cannot join a new channel %s", chainID)
		logger.Errorf(errMsg)
		return errors.New(errMsg)
	}
	if _, exist := d.blockProviders[chainID]; exist {
		errMsg := fmt.Sprintf("Delivery service - block provider already exists for %s found, can't start delivery", chainID)
		logger.Errorf(errMsg)
		return errors.New(errMsg)
	}
	client := d.newClient(chainID, ledgerInfo)
	logger.Debug("This peer will pass blocks from orderer service to other peers for channel", chainID)
	d.blockProviders[chainID] = blocksprovider.NewBlocksProvider(chainID, client, d.conf.Gossip, d.conf.CryptoSvc)
	go d.launchBlockProvider(chainID, finalizer)
	return nil
}

func (d *deliverServiceImpl) launchBlockProvider(chainID string, finalizer func()) {
	d.lock.RLock()
	pb := d.blockProviders[chainID]
	d.lock.RUnlock()
	if pb == nil {
		logger.Info("Block delivery for channel", chainID, "was stopped before block provider started")
		return
	}
	pb.DeliverBlocks()
	finalizer()
}


func (d *deliverServiceImpl) StopDeliverForChannel(chainID string) error {
	d.lock.Lock()
	defer d.lock.Unlock()
	if d.stopping {
		errMsg := fmt.Sprintf("Delivery service is stopping, cannot stop delivery for channel %s", chainID)
		logger.Errorf(errMsg)
		return errors.New(errMsg)
	}
	client, exist := d.blockProviders[chainID]
	if !exist {
		errMsg := fmt.Sprintf("Delivery service - no block provider for %s found, can't stop delivery", chainID)
		logger.Errorf(errMsg)
		return errors.New(errMsg)
	}
	client.Stop()
	delete(d.blockProviders, chainID)
	logger.Debug("This peer will stop pass blocks from orderer service to other peers")
	return nil
}


func (d *deliverServiceImpl) Stop() {
	d.lock.Lock()
	defer d.lock.Unlock()
	
	d.stopping = true

	for _, client := range d.blockProviders {
		client.Stop()
	}
}

func (d *deliverServiceImpl) newClient(chainID string, ledgerInfoProvider blocksprovider.LedgerInfo) *broadcastClient {
	reconnectBackoffThreshold := d.conf.DeliverServiceConfig.ReConnectBackoffThreshold
	reconnectTotalTimeThreshold := d.conf.DeliverServiceConfig.ReconnectTotalTimeThreshold

	requester := &blocksRequester{
		tls:         d.conf.DeliverServiceConfig.PeerTLSEnabled,
		chainID:     chainID,
		signer:      d.conf.Signer,
		credSupport: d.conf.CredentialSupport,
	}
	broadcastSetup := func(bd blocksprovider.BlocksDeliverer) error {
		return requester.RequestBlocks(ledgerInfoProvider)
	}
	backoffPolicy := func(attemptNum int, elapsedTime time.Duration) (time.Duration, bool) {
		if elapsedTime >= reconnectTotalTimeThreshold {
			return 0, false
		}
		sleepIncrement := float64(time.Millisecond * 500)
		attempt := float64(attemptNum)
		return time.Duration(math.Min(math.Pow(2, attempt)*sleepIncrement, reconnectBackoffThreshold)), true
	}
	connectionFactory := d.conf.ConnFactory(chainID)
	connProd := comm.NewConnectionProducer(
		connectionFactory,
		d.conf.Endpoints,
		d.conf.DeliverClientDialOpts,
		d.conf.DeliverServiceConfig.PeerTLSEnabled,
		d.conf.DeliverServiceConfig.ConnectionTimeout,
	)
	bClient := NewBroadcastClient(connProd, d.conf.ABCFactory, broadcastSetup, backoffPolicy)
	requester.client = bClient
	return bClient
}

type CredSupportDialerFactory struct {
	CredentialSupport *comm.CredentialSupport
	KeepaliveOptions  comm.KeepaliveOptions
	TLSEnabled        bool
}

func (c *CredSupportDialerFactory) Dialer(channelID string) func(endpoint string, connectionTimeout time.Duration) (*grpc.ClientConn, error) {
	return func(endpoint string, connectionTimeout time.Duration) (*grpc.ClientConn, error) {
		dialOpts := []grpc.DialOption{
			grpc.WithBlock(),
			grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(comm.MaxRecvMsgSize), grpc.MaxCallSendMsgSize(comm.MaxSendMsgSize)),
		}
		dialOpts = append(dialOpts, comm.ClientKeepaliveOptions(c.KeepaliveOptions)...)

		if c.TLSEnabled {
			creds, err := c.CredentialSupport.GetDeliverServiceCredentials(channelID)
			if err != nil {
				return nil, fmt.Errorf("failed obtaining credentials for channel %s: %v", channelID, err)
			}
			dialOpts = append(dialOpts, grpc.WithTransportCredentials(creds))
		} else {
			dialOpts = append(dialOpts, grpc.WithInsecure())
		}
		ctx, cancel := context.WithTimeout(context.Background(), connectionTimeout)
		defer cancel()
		return grpc.DialContext(ctx, endpoint, dialOpts...)
	}
}

func DefaultABCFactory(conn *grpc.ClientConn) orderer.AtomicBroadcastClient {
	return orderer.NewAtomicBroadcastClient(conn)
}
