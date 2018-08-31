/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package discovery

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/common/util"
	"github.com/mcc-github/blockchain/core/comm"
	common2 "github.com/mcc-github/blockchain/gossip/common"
	discovery2 "github.com/mcc-github/blockchain/gossip/discovery"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/discovery"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

var (
	logger = flogging.MustGetLogger("discovery")
)

var accessDenied = wrapError(errors.New("access denied"))



type certHashExtractor func(ctx context.Context) []byte


type dispatcher func(q *discovery.Query) *discovery.QueryResult

type service struct {
	config             Config
	channelDispatchers map[discovery.QueryType]dispatcher
	localDispatchers   map[discovery.QueryType]dispatcher
	auth               *authCache
	Support
}


type Config struct {
	TLS                          bool
	AuthCacheEnabled             bool
	AuthCacheMaxSize             int
	AuthCachePurgeRetentionRatio float64
}


func (c Config) String() string {
	if c.AuthCacheEnabled {
		return fmt.Sprintf("TLS: %t, authCacheMaxSize: %d, authCachePurgeRatio: %f", c.TLS, c.AuthCacheMaxSize, c.AuthCachePurgeRetentionRatio)
	}
	return fmt.Sprintf("TLS: %t, auth cache disabled", c.TLS)
}


type peerMapping map[string]*discovery.Peer


func NewService(config Config, sup Support) *service {
	s := &service{
		auth: newAuthCache(sup, authCacheConfig{
			enabled:             config.AuthCacheEnabled,
			maxCacheSize:        config.AuthCacheMaxSize,
			purgeRetentionRatio: config.AuthCachePurgeRetentionRatio,
		}),
		Support: sup,
	}
	s.channelDispatchers = map[discovery.QueryType]dispatcher{
		discovery.ConfigQueryType:         s.configQuery,
		discovery.ChaincodeQueryType:      s.chaincodeQuery,
		discovery.PeerMembershipQueryType: s.channelMembershipResponse,
	}
	s.localDispatchers = map[discovery.QueryType]dispatcher{
		discovery.LocalMembershipQueryType: s.localMembershipResponse,
	}
	logger.Info("Created with config", config)
	return s
}

func (s *service) Discover(ctx context.Context, request *discovery.SignedRequest) (*discovery.Response, error) {
	addr := util.ExtractRemoteAddress(ctx)
	req, err := validateStructure(ctx, request, s.config.TLS, comm.ExtractCertificateHashFromContext)
	if err != nil {
		logger.Warningf("Request from %s is malformed or invalid: %v", addr, err)
		return nil, err
	}
	logger.Debugf("Processing request from %s: %v", addr, req)
	var res []*discovery.QueryResult
	for _, q := range req.Queries {
		res = append(res, s.processQuery(q, request, req.Authentication.ClientIdentity, addr))
	}
	logger.Debugf("Returning to %s a response containing: %v", addr, res)
	return &discovery.Response{
		Results: res,
	}, nil
}

func (s *service) processQuery(query *discovery.Query, request *discovery.SignedRequest, identity []byte, addr string) *discovery.QueryResult {
	if query.Channel != "" && !s.ChannelExists(query.Channel) {
		logger.Warning("got query for channel", query.Channel, "from", addr, "but it doesn't exist")
		return accessDenied
	}
	if err := s.auth.EligibleForService(query.Channel, common.SignedData{
		Data:      request.Payload,
		Signature: request.Signature,
		Identity:  identity,
	}); err != nil {
		logger.Warning("got query for channel", query.Channel, "from", addr, "but it isn't eligible:", err)
		return accessDenied
	}
	return s.dispatch(query)
}

func (s *service) dispatch(q *discovery.Query) *discovery.QueryResult {
	dispatchers := s.channelDispatchers
	
	if q.Channel == "" {
		dispatchers = s.localDispatchers
	}
	dispatchQuery, exists := dispatchers[q.GetType()]
	if !exists {
		return wrapError(errors.New("unknown or missing request type"))
	}
	return dispatchQuery(q)
}

func (s *service) chaincodeQuery(q *discovery.Query) *discovery.QueryResult {
	if err := validateCCQuery(q.GetCcQuery()); err != nil {
		return wrapError(err)
	}
	var descriptors []*discovery.EndorsementDescriptor
	for _, interest := range q.GetCcQuery().Interests {
		desc, err := s.PeersForEndorsement(common2.ChainID(q.Channel), interest)
		if err != nil {
			logger.Errorf("Failed constructing descriptor for chaincode %s,: %v", interest, err)
			return wrapError(errors.Errorf("failed constructing descriptor for %v", interest))
		}
		descriptors = append(descriptors, desc)
	}

	return &discovery.QueryResult{
		Result: &discovery.QueryResult_CcQueryRes{
			CcQueryRes: &discovery.ChaincodeQueryResult{
				Content: descriptors,
			},
		},
	}
}

func (s *service) configQuery(q *discovery.Query) *discovery.QueryResult {
	conf, err := s.Config(q.Channel)
	if err != nil {
		logger.Errorf("Failed fetching config for channel %s: %v", q.Channel, err)
		return wrapError(errors.Errorf("failed fetching config for channel %s", q.Channel))
	}
	return &discovery.QueryResult{
		Result: &discovery.QueryResult_ConfigResult{
			ConfigResult: conf,
		},
	}
}

func wrapPeerResponse(peersByOrg map[string]*discovery.Peers) *discovery.QueryResult {
	return &discovery.QueryResult{
		Result: &discovery.QueryResult_Members{
			Members: &discovery.PeerMembershipResult{
				PeersByOrg: peersByOrg,
			},
		},
	}
}

func (s *service) channelMembershipResponse(q *discovery.Query) *discovery.QueryResult {
	chanPeers, err := s.PeersAuthorizedByCriteria(common2.ChainID(q.Channel), q.GetPeerQuery().Filter)
	if err != nil {
		return wrapError(err)
	}
	membersByOrgs := make(map[string]*discovery.Peers)
	chanPeerByID := discovery2.Members(chanPeers).ByID()
	for org, ids2Peers := range s.computeMembership(q) {
		membersByOrgs[org] = &discovery.Peers{}
		for id, peer := range ids2Peers {
			
			stateInfoMsg, exists := chanPeerByID[string(id)]
			
			if !exists {
				continue
			}
			peer.StateInfo = stateInfoMsg.Envelope
			membersByOrgs[org].Peers = append(membersByOrgs[org].Peers, peer)
		}
	}
	return wrapPeerResponse(membersByOrgs)
}

func (s *service) localMembershipResponse(q *discovery.Query) *discovery.QueryResult {
	membersByOrgs := make(map[string]*discovery.Peers)
	for org, ids2Peers := range s.computeMembership(q) {
		membersByOrgs[org] = &discovery.Peers{}
		for _, peer := range ids2Peers {
			membersByOrgs[org].Peers = append(membersByOrgs[org].Peers, peer)
		}
	}
	return wrapPeerResponse(membersByOrgs)
}

func (s *service) computeMembership(_ *discovery.Query) map[string]peerMapping {
	peersByOrg := make(map[string]peerMapping)
	peerAliveInfo := discovery2.Members(s.Peers()).ByID()
	for org, peerIdentities := range s.IdentityInfo().ByOrg() {
		peersForCurrentOrg := make(peerMapping)
		peersByOrg[org] = peersForCurrentOrg
		for _, id := range peerIdentities {
			
			aliveInfo, exists := peerAliveInfo[string(id.PKIId)]
			if !exists {
				continue
			}
			peersForCurrentOrg[string(id.PKIId)] = &discovery.Peer{
				Identity:       id.Identity,
				MembershipInfo: aliveInfo.Envelope,
			}
		}
	}
	return peersByOrg
}


func validateStructure(ctx context.Context, request *discovery.SignedRequest, tlsEnabled bool, certHashFromContext certHashExtractor) (*discovery.Request, error) {
	if request == nil {
		return nil, errors.New("nil request")
	}
	req, err := request.ToRequest()
	if err != nil {
		return nil, errors.Wrap(err, "failed parsing request")
	}
	if req.Authentication == nil {
		return nil, errors.New("access denied, no authentication info in request")
	}
	if len(req.Authentication.ClientIdentity) == 0 {
		return nil, errors.New("access denied, client identity wasn't supplied")
	}
	if !tlsEnabled {
		return req, nil
	}
	computedHash := certHashFromContext(ctx)
	if len(computedHash) == 0 {
		return nil, errors.New("client didn't send a TLS certificate")
	}
	if !bytes.Equal(computedHash, req.Authentication.ClientTlsCertHash) {
		claimed := hex.EncodeToString(req.Authentication.ClientTlsCertHash)
		logger.Warningf("client claimed TLS hash %s doesn't match computed TLS hash from gRPC stream %s", claimed, hex.EncodeToString(computedHash))
		return nil, errors.New("client claimed TLS hash doesn't match computed TLS hash from gRPC stream")
	}
	return req, nil
}

func validateCCQuery(ccQuery *discovery.ChaincodeQuery) error {
	if len(ccQuery.Interests) == 0 {
		return errors.New("chaincode query must have at least one chaincode interest")
	}
	for _, interest := range ccQuery.Interests {
		if interest == nil {
			return errors.New("chaincode interest is nil")
		}
		if len(interest.Chaincodes) == 0 {
			return errors.New("chaincode interest must contain at least one chaincode")
		}
		for _, cc := range interest.Chaincodes {
			if cc.Name == "" {
				return errors.New("chaincode name in interest cannot be empty")
			}
		}
	}
	return nil
}

func wrapError(err error) *discovery.QueryResult {
	return &discovery.QueryResult{
		Result: &discovery.QueryResult_Error{
			Error: &discovery.Error{
				Content: err.Error(),
			},
		},
	}
}
