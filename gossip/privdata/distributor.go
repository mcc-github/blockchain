/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package privdata

import (
	"bytes"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain-protos-go/common"
	protosgossip "github.com/mcc-github/blockchain-protos-go/gossip"
	"github.com/mcc-github/blockchain-protos-go/ledger/rwset"
	"github.com/mcc-github/blockchain-protos-go/transientstore"
	"github.com/mcc-github/blockchain/core/common/privdata"
	"github.com/mcc-github/blockchain/gossip/api"
	gossipCommon "github.com/mcc-github/blockchain/gossip/common"
	"github.com/mcc-github/blockchain/gossip/discovery"
	"github.com/mcc-github/blockchain/gossip/filter"
	gossipgossip "github.com/mcc-github/blockchain/gossip/gossip"
	"github.com/mcc-github/blockchain/gossip/metrics"
	"github.com/mcc-github/blockchain/gossip/protoext"
	"github.com/mcc-github/blockchain/gossip/util"
	"github.com/mcc-github/blockchain/msp"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/pkg/errors"
)


type gossipAdapter interface {
	
	SendByCriteria(message *protoext.SignedGossipMessage, criteria gossipgossip.SendCriteria) error

	
	
	PeerFilter(channel gossipCommon.ChannelID, messagePredicate api.SubChannelSelectionCriteria) (filter.RoutingFilter, error)

	
	IdentityInfo() api.PeerIdentitySet

	
	
	PeersOfChannel(gossipCommon.ChannelID) []discovery.NetworkMember
}


type PvtDataDistributor interface {
	
	Distribute(txID string, privData *transientstore.TxPvtReadWriteSetWithConfigInfo, blkHt uint64) error
}



type IdentityDeserializerFactory interface {
	
	
	GetIdentityDeserializer(chainID string) msp.IdentityDeserializer
}



type IdentityDeserializerFactoryFunc func(chainID string) msp.IdentityDeserializer

func (i IdentityDeserializerFactoryFunc) GetIdentityDeserializer(chainID string) msp.IdentityDeserializer {
	return i(chainID)
}


type distributorImpl struct {
	chainID string
	gossipAdapter
	CollectionAccessFactory
	pushAckTimeout time.Duration
	metrics        *metrics.PrivdataMetrics
}


type CollectionAccessFactory interface {
	
	AccessPolicy(config *common.CollectionConfig, chainID string) (privdata.CollectionAccessPolicy, error)
}


type policyAccessFactory struct {
	IdentityDeserializerFactory
}

func (p *policyAccessFactory) AccessPolicy(config *common.CollectionConfig, chainID string) (privdata.CollectionAccessPolicy, error) {
	colAP := &privdata.SimpleCollection{}
	switch cconf := config.Payload.(type) {
	case *common.CollectionConfig_StaticCollectionConfig:
		err := colAP.Setup(cconf.StaticCollectionConfig, p.GetIdentityDeserializer(chainID))
		if err != nil {
			return nil, errors.WithMessagef(err, "error setting up collection  %#v", cconf.StaticCollectionConfig.Name)
		}
	default:
		return nil, errors.New("unexpected collection type")
	}
	return colAP, nil
}


func NewCollectionAccessFactory(factory IdentityDeserializerFactory) CollectionAccessFactory {
	return &policyAccessFactory{
		IdentityDeserializerFactory: factory,
	}
}



func NewDistributor(chainID string, gossip gossipAdapter, factory CollectionAccessFactory,
	metrics *metrics.PrivdataMetrics, pushAckTimeout time.Duration) PvtDataDistributor {
	return &distributorImpl{
		chainID:                 chainID,
		gossipAdapter:           gossip,
		CollectionAccessFactory: factory,
		pushAckTimeout:          pushAckTimeout,
		metrics:                 metrics,
	}
}


func (d *distributorImpl) Distribute(txID string, privData *transientstore.TxPvtReadWriteSetWithConfigInfo, blkHt uint64) error {
	disseminationPlan, err := d.computeDisseminationPlan(txID, privData, blkHt)
	if err != nil {
		return errors.WithStack(err)
	}
	return d.disseminate(disseminationPlan)
}

type dissemination struct {
	msg      *protoext.SignedGossipMessage
	criteria gossipgossip.SendCriteria
}

func (d *distributorImpl) computeDisseminationPlan(txID string,
	privDataWithConfig *transientstore.TxPvtReadWriteSetWithConfigInfo,
	blkHt uint64) ([]*dissemination, error) {
	privData := privDataWithConfig.PvtRwset
	var disseminationPlan []*dissemination
	for _, pvtRwset := range privData.NsPvtRwset {
		namespace := pvtRwset.Namespace
		configPackage, found := privDataWithConfig.CollectionConfigs[namespace]
		if !found {
			logger.Error("Collection config package for", namespace, "chaincode is not provided")
			return nil, errors.New(fmt.Sprint("collection config package for", namespace, "chaincode is not provided"))
		}

		for _, collection := range pvtRwset.CollectionPvtRwset {
			colCP, err := d.getCollectionConfig(configPackage, collection)
			collectionName := collection.CollectionName
			if err != nil {
				logger.Error("Could not find collection access policy for", namespace, " and collection", collectionName, "error", err)
				return nil, errors.WithMessage(err, fmt.Sprint("could not find collection access policy for", namespace, " and collection", collectionName, "error", err))
			}

			colAP, err := d.AccessPolicy(colCP, d.chainID)
			if err != nil {
				logger.Error("Could not obtain collection access policy, collection name", collectionName, "due to", err)
				return nil, errors.Wrap(err, fmt.Sprint("Could not obtain collection access policy, collection name", collectionName, "due to", err))
			}

			colFilter := colAP.AccessFilter()
			if colFilter == nil {
				logger.Error("Collection access policy for", collectionName, "has no filter")
				return nil, errors.Errorf("No collection access policy filter computed for %v", collectionName)
			}

			pvtDataMsg, err := d.createPrivateDataMessage(txID, namespace, collection, &common.CollectionConfigPackage{Config: []*common.CollectionConfig{colCP}}, blkHt)
			if err != nil {
				return nil, errors.WithStack(err)
			}

			dPlan, err := d.disseminationPlanForMsg(colAP, colFilter, pvtDataMsg)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			disseminationPlan = append(disseminationPlan, dPlan...)
		}
	}
	return disseminationPlan, nil
}

func (d *distributorImpl) getCollectionConfig(config *common.CollectionConfigPackage, collection *rwset.CollectionPvtReadWriteSet) (*common.CollectionConfig, error) {
	for _, c := range config.Config {
		if staticConfig := c.GetStaticCollectionConfig(); staticConfig != nil {
			if staticConfig.Name == collection.CollectionName {
				return c, nil
			}
		}
	}
	return nil, errors.New(fmt.Sprint("no configuration for collection", collection.CollectionName, "found"))
}

func (d *distributorImpl) disseminationPlanForMsg(colAP privdata.CollectionAccessPolicy, colFilter privdata.Filter, pvtDataMsg *protoext.SignedGossipMessage) ([]*dissemination, error) {
	var disseminationPlan []*dissemination

	routingFilter, err := d.gossipAdapter.PeerFilter(gossipCommon.ChannelID(d.chainID), func(signature api.PeerSignature) bool {
		return colFilter(protoutil.SignedData{
			Data:      signature.Message,
			Signature: signature.Signature,
			Identity:  []byte(signature.PeerIdentity),
		})
	})

	if err != nil {
		logger.Error("Failed to retrieve peer routing filter for channel", d.chainID, ":", err)
		return nil, err
	}

	eligiblePeers := d.eligiblePeersOfChannel(routingFilter)
	identitySets := d.identitiesOfEligiblePeers(eligiblePeers, colAP)

	
	maximumPeerCount := colAP.MaximumPeerCount()
	requiredPeerCount := colAP.RequiredPeerCount()

	if maximumPeerCount > 0 {
		for _, selectionPeers := range identitySets {
			required := 1
			if requiredPeerCount == 0 {
				required = 0
			}
			peer2SendPerOrg := selectionPeers[rand.Intn(len(selectionPeers))]
			sc := gossipgossip.SendCriteria{
				Timeout:  d.pushAckTimeout,
				Channel:  gossipCommon.ChannelID(d.chainID),
				MaxPeers: 1,
				MinAck:   required,
				IsEligible: func(member discovery.NetworkMember) bool {
					return bytes.Equal(member.PKIid, peer2SendPerOrg.PKIId)
				},
			}
			disseminationPlan = append(disseminationPlan, &dissemination{
				criteria: sc,
				msg: &protoext.SignedGossipMessage{
					Envelope:      proto.Clone(pvtDataMsg.Envelope).(*protosgossip.Envelope),
					GossipMessage: proto.Clone(pvtDataMsg.GossipMessage).(*protosgossip.GossipMessage),
				},
			})

			if requiredPeerCount > 0 {
				requiredPeerCount--
			}

			maximumPeerCount--
			if maximumPeerCount == 0 {
				return disseminationPlan, nil
			}
		}
	}

	
	
	sc := gossipgossip.SendCriteria{
		Timeout:  d.pushAckTimeout,
		Channel:  gossipCommon.ChannelID(d.chainID),
		MaxPeers: maximumPeerCount,
		MinAck:   requiredPeerCount,
		IsEligible: func(member discovery.NetworkMember) bool {
			return routingFilter(member)
		},
	}

	disseminationPlan = append(disseminationPlan, &dissemination{
		criteria: sc,
		msg:      pvtDataMsg,
	})

	return disseminationPlan, nil
}

func (d *distributorImpl) identitiesOfEligiblePeers(eligiblePeers []discovery.NetworkMember, colAP privdata.CollectionAccessPolicy) map[string]api.PeerIdentitySet {
	return d.gossipAdapter.IdentityInfo().
		Filter(func(info api.PeerIdentityInfo) bool {
			for _, orgID := range colAP.MemberOrgs() {
				if bytes.Equal(info.Organization, []byte(orgID)) {
					return true
				}
			}
			
			return false
		}).Filter(func(info api.PeerIdentityInfo) bool {
		for _, peer := range eligiblePeers {
			if bytes.Equal(info.PKIId, peer.PKIid) {
				return true
			}
		}
		
		return false
	}).ByOrg()
}

func (d *distributorImpl) eligiblePeersOfChannel(routingFilter filter.RoutingFilter) []discovery.NetworkMember {
	var eligiblePeers []discovery.NetworkMember
	for _, peer := range d.gossipAdapter.PeersOfChannel(gossipCommon.ChannelID(d.chainID)) {
		if routingFilter(peer) {
			eligiblePeers = append(eligiblePeers, peer)
		}
	}
	return eligiblePeers
}

func (d *distributorImpl) disseminate(disseminationPlan []*dissemination) error {
	var failures uint32
	var wg sync.WaitGroup
	wg.Add(len(disseminationPlan))
	start := time.Now()
	for _, dis := range disseminationPlan {
		go func(dis *dissemination) {
			defer wg.Done()
			defer d.reportSendDuration(start)
			err := d.SendByCriteria(dis.msg, dis.criteria)
			if err != nil {
				atomic.AddUint32(&failures, 1)
				m := dis.msg.GetPrivateData().Payload
				logger.Error("Failed disseminating private RWSet for TxID", m.TxId, ", namespace", m.Namespace, "collection", m.CollectionName, ":", err)
			}
		}(dis)
	}
	wg.Wait()
	failureCount := atomic.LoadUint32(&failures)
	if failureCount != 0 {
		return errors.Errorf("Failed disseminating %d out of %d private dissemination plans", failureCount, len(disseminationPlan))
	}
	return nil
}

func (d *distributorImpl) reportSendDuration(startTime time.Time) {
	d.metrics.SendDuration.With("channel", d.chainID).Observe(time.Since(startTime).Seconds())
}

func (d *distributorImpl) createPrivateDataMessage(txID, namespace string,
	collection *rwset.CollectionPvtReadWriteSet,
	ccp *common.CollectionConfigPackage,
	blkHt uint64) (*protoext.SignedGossipMessage, error) {
	msg := &protosgossip.GossipMessage{
		Channel: []byte(d.chainID),
		Nonce:   util.RandomUInt64(),
		Tag:     protosgossip.GossipMessage_CHAN_ONLY,
		Content: &protosgossip.GossipMessage_PrivateData{
			PrivateData: &protosgossip.PrivateDataMessage{
				Payload: &protosgossip.PrivatePayload{
					Namespace:         namespace,
					CollectionName:    collection.CollectionName,
					TxId:              txID,
					PrivateRwset:      collection.Rwset,
					PrivateSimHeight:  blkHt,
					CollectionConfigs: ccp,
				},
			},
		},
	}

	pvtDataMsg, err := protoext.NoopSign(msg)
	if err != nil {
		return nil, err
	}
	return pvtDataMsg, nil
}
