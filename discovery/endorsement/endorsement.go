/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package endorsement

import (
	"fmt"

	"github.com/mcc-github/blockchain/common/chaincode"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/common/graph"
	"github.com/mcc-github/blockchain/common/policies"
	"github.com/mcc-github/blockchain/common/policies/inquire"
	"github.com/mcc-github/blockchain/gossip/api"
	"github.com/mcc-github/blockchain/gossip/common"
	. "github.com/mcc-github/blockchain/gossip/discovery"
	"github.com/mcc-github/blockchain/protos/discovery"
	"github.com/mcc-github/blockchain/protos/msp"
	"github.com/pkg/errors"
)

var (
	logger = flogging.MustGetLogger("discovery/endorsement")
)

type principalEvaluator interface {
	
	
	SatisfiesPrincipal(channel string, identity []byte, principal *msp.MSPPrincipal) error
}

type chaincodeMetadataFetcher interface {
	
	
	Metadata(channel string, cc string, loadCollections bool) *chaincode.Metadata
}

type policyFetcher interface {
	
	
	PolicyByChaincode(channel string, cc string) policies.InquireablePolicy
}

type gossipSupport interface {
	
	IdentityInfo() api.PeerIdentitySet

	
	
	PeersOfChannel(common.ChainID) Members

	
	Peers() Members
}

type endorsementAnalyzer struct {
	gossipSupport
	principalEvaluator
	policyFetcher
	chaincodeMetadataFetcher
}


func NewEndorsementAnalyzer(gs gossipSupport, pf policyFetcher, pe principalEvaluator, mf chaincodeMetadataFetcher) *endorsementAnalyzer {
	return &endorsementAnalyzer{
		gossipSupport:            gs,
		policyFetcher:            pf,
		principalEvaluator:       pe,
		chaincodeMetadataFetcher: mf,
	}
}

type peerPrincipalEvaluator func(member NetworkMember, principal *msp.MSPPrincipal) bool


func (ea *endorsementAnalyzer) PeersForEndorsement(chainID common.ChainID, interest *discovery.ChaincodeInterest) (*discovery.EndorsementDescriptor, error) {
	chanMembership, err := ea.PeersAuthorizedByCriteria(chainID, interest)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	channelMembersById := chanMembership.ByID()
	
	aliveMembership := ea.Peers().Intersect(chanMembership)
	membersById := aliveMembership.ByID()
	
	identitiesOfMembers := computeIdentitiesOfMembers(ea.IdentityInfo(), membersById)
	principalsSets, err := ea.computePrincipalSets(chainID, interest)
	if err != nil {
		logger.Warningf("Principal set computation failed: %v", err)
		return nil, errors.WithStack(err)
	}

	return ea.computeEndorsementResponse(&context{
		chaincode:           interest.Chaincodes[0].Name,
		channel:             string(chainID),
		principalsSets:      principalsSets,
		channelMembersById:  channelMembersById,
		aliveMembership:     aliveMembership,
		identitiesOfMembers: identitiesOfMembers,
	})
}

func (ea *endorsementAnalyzer) PeersAuthorizedByCriteria(chainID common.ChainID, interest *discovery.ChaincodeInterest) (Members, error) {
	peersOfChannel := ea.PeersOfChannel(chainID)
	if interest == nil || len(interest.Chaincodes) == 0 {
		return peersOfChannel, nil
	}
	identities := ea.IdentityInfo()
	identitiesByID := identities.ByID()
	metadataAndCollectionFilters, err := loadMetadataAndFilters(metadataAndFilterContext{
		identityInfoByID: identitiesByID,
		interest:         interest,
		chainID:          chainID,
		evaluator:        ea,
		fetch:            ea,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	metadata := metadataAndCollectionFilters.md
	
	chanMembership := peersOfChannel.Filter(peersWithChaincode(metadata...))
	
	return chanMembership.Filter(metadataAndCollectionFilters.isMemberAuthorized), nil
}

type context struct {
	chaincode           string
	channel             string
	aliveMembership     Members
	principalsSets      []policies.PrincipalSet
	channelMembersById  map[string]NetworkMember
	identitiesOfMembers memberIdentities
}

func (ea *endorsementAnalyzer) computeEndorsementResponse(ctx *context) (*discovery.EndorsementDescriptor, error) {
	
	
	principalGroups := mapPrincipalsToGroups(ctx.principalsSets)
	
	
	
	satGraph := principalsToPeersGraph(principalAndPeerData{
		members: ctx.aliveMembership,
		pGrps:   principalGroups,
	}, ea.satisfiesPrincipal(ctx.channel, ctx.identitiesOfMembers))

	layouts := computeLayouts(ctx.principalsSets, principalGroups, satGraph)
	if len(layouts) == 0 {
		return nil, errors.New("cannot satisfy any principal combination")
	}

	criteria := &peerMembershipCriteria{
		possibleLayouts: layouts,
		satGraph:        satGraph,
		chanMemberById:  ctx.channelMembersById,
		idOfMembers:     ctx.identitiesOfMembers,
	}

	return &discovery.EndorsementDescriptor{
		Chaincode:         ctx.chaincode,
		Layouts:           layouts,
		EndorsersByGroups: endorsersByGroup(criteria),
	}, nil
}

func (ea *endorsementAnalyzer) computePrincipalSets(chainID common.ChainID, interest *discovery.ChaincodeInterest) (policies.PrincipalSets, error) {
	var inquireablePolicies []policies.InquireablePolicy
	for _, chaincode := range interest.Chaincodes {
		pol := ea.PolicyByChaincode(string(chainID), chaincode.Name)
		if pol == nil {
			logger.Debug("Policy for chaincode '", chaincode, "'doesn't exist")
			return nil, errors.New("policy not found")
		}
		inquireablePolicies = append(inquireablePolicies, pol)
	}

	var cpss []inquire.ComparablePrincipalSets

	for _, policy := range inquireablePolicies {
		var cmpsets inquire.ComparablePrincipalSets
		for _, ps := range policy.SatisfiedBy() {
			cps := inquire.NewComparablePrincipalSet(ps)
			if cps == nil {
				return nil, errors.New("failed creating a comparable principal set")
			}
			cmpsets = append(cmpsets, cps)
		}
		if len(cmpsets) == 0 {
			return nil, errors.New("chaincode isn't installed on sufficient organizations required by the endorsement policy")
		}
		cpss = append(cpss, cmpsets)
	}

	cps, err := mergePrincipalSets(cpss)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return cps.ToPrincipalSets(), nil
}

type metadataAndFilterContext struct {
	chainID          common.ChainID
	interest         *discovery.ChaincodeInterest
	fetch            chaincodeMetadataFetcher
	identityInfoByID map[string]api.PeerIdentityInfo
	evaluator        principalEvaluator
}


type metadataAndColFilter struct {
	md                 []*chaincode.Metadata
	isMemberAuthorized memberFilter
}

func loadMetadataAndFilters(ctx metadataAndFilterContext) (*metadataAndColFilter, error) {
	var metadata []*chaincode.Metadata
	var filters []identityFilter

	for _, chaincode := range ctx.interest.Chaincodes {
		ccMD := ctx.fetch.Metadata(string(ctx.chainID), chaincode.Name, len(chaincode.CollectionNames) > 0)
		if ccMD == nil {
			return nil, errors.Errorf("No metadata was found for chaincode %s in channel %s", chaincode.Name, string(ctx.chainID))
		}
		metadata = append(metadata, ccMD)
		if len(chaincode.CollectionNames) == 0 {
			continue
		}
		principalSetByCollections, err := principalsFromCollectionConfig(ccMD.CollectionsConfig)
		if err != nil {
			logger.Warningf("Failed initializing collection filter for chaincode %s: %v", chaincode.Name, err)
			return nil, errors.WithStack(err)
		}
		filter, err := principalSetByCollections.toIdentityFilter(string(ctx.chainID), ctx.evaluator, chaincode)
		if err != nil {
			logger.Warningf("Failed computing collection principal sets for chaincode %s due to %v", chaincode.Name, err)
			return nil, errors.WithStack(err)
		}
		filters = append(filters, filter)
	}

	return computeFiltersWithMetadata(filters, metadata, ctx.identityInfoByID), nil
}

func computeFiltersWithMetadata(filters identityFilters, metadata []*chaincode.Metadata, identityInfoByID map[string]api.PeerIdentityInfo) *metadataAndColFilter {
	if len(filters) == 0 {
		return &metadataAndColFilter{
			md:                 metadata,
			isMemberAuthorized: noopMemberFilter,
		}
	}
	filter := filters.combine().toMemberFilter(identityInfoByID)
	return &metadataAndColFilter{
		isMemberAuthorized: filter,
		md:                 metadata,
	}
}


type identityFilter func(api.PeerIdentityType) bool


type identityFilters []identityFilter


type memberFilter func(member NetworkMember) bool


func noopMemberFilter(_ NetworkMember) bool {
	return true
}



func (filters identityFilters) combine() identityFilter {
	return func(identity api.PeerIdentityType) bool {
		for _, f := range filters {
			if !f(identity) {
				return false
			}
		}
		return true
	}
}



func (idf identityFilter) toMemberFilter(identityInfoByID map[string]api.PeerIdentityInfo) memberFilter {
	return func(member NetworkMember) bool {
		identity, exists := identityInfoByID[string(member.PKIid)]
		if !exists {
			return false
		}
		return idf(identity.Identity)
	}
}

func (ea *endorsementAnalyzer) satisfiesPrincipal(channel string, identitiesOfMembers memberIdentities) peerPrincipalEvaluator {
	return func(member NetworkMember, principal *msp.MSPPrincipal) bool {
		err := ea.SatisfiesPrincipal(channel, identitiesOfMembers.identityByPKIID(member.PKIid), principal)
		if err == nil {
			
			logger.Debug(member, "satisfies principal", principal)
			return true
		}
		logger.Debug(member, "doesn't satisfy principal", principal, ":", err)
		return false
	}
}

type peerMembershipCriteria struct {
	satGraph        *principalPeerGraph
	idOfMembers     memberIdentities
	chanMemberById  map[string]NetworkMember
	possibleLayouts layouts
}








func endorsersByGroup(criteria *peerMembershipCriteria) map[string]*discovery.Peers {
	satGraph := criteria.satGraph
	idOfMembers := criteria.idOfMembers
	chanMemberById := criteria.chanMemberById
	includedGroups := criteria.possibleLayouts.groupsSet()

	res := make(map[string]*discovery.Peers)
	
	
	for grp, principalVertex := range satGraph.principalVertices {
		if _, exists := includedGroups[grp]; !exists {
			
			continue
		}
		peerList := &discovery.Peers{}
		res[grp] = peerList
		for _, peerVertex := range principalVertex.Neighbors() {
			member := peerVertex.Data.(NetworkMember)
			peerList.Peers = append(peerList.Peers, &discovery.Peer{
				Identity:       idOfMembers.identityByPKIID(member.PKIid),
				StateInfo:      chanMemberById[string(member.PKIid)].Envelope,
				MembershipInfo: member.Envelope,
			})
		}
	}
	return res
}







func computeLayouts(principalsSets []policies.PrincipalSet, principalGroups principalGroupMapper, satGraph *principalPeerGraph) []*discovery.Layout {
	var layouts []*discovery.Layout
	
	
	for _, principalSet := range principalsSets {
		layout := &discovery.Layout{
			QuantitiesByGroup: make(map[string]uint32),
		}
		
		
		for principal, plurality := range principalSet.UniqueSet() {
			key := principalKey{
				cls:       int32(principal.PrincipalClassification),
				principal: string(principal.Principal),
			}
			
			layout.QuantitiesByGroup[principalGroups.group(key)] = uint32(plurality)
		}
		
		
		
		if isLayoutSatisfied(layout.QuantitiesByGroup, satGraph) {
			
			
			layouts = append(layouts, layout)
		}
	}
	return layouts
}

func isLayoutSatisfied(layout map[string]uint32, satGraph *principalPeerGraph) bool {
	for grp, plurality := range layout {
		
		if len(satGraph.principalVertices[grp].Neighbors()) < int(plurality) {
			return false
		}
	}
	return true
}

type principalPeerGraph struct {
	peerVertices      []*graph.Vertex
	principalVertices map[string]*graph.Vertex
}

type principalAndPeerData struct {
	members Members
	pGrps   principalGroupMapper
}

func principalsToPeersGraph(data principalAndPeerData, satisfiesPrincipal peerPrincipalEvaluator) *principalPeerGraph {
	
	peerVertices := make([]*graph.Vertex, len(data.members))
	for i, member := range data.members {
		peerVertices[i] = graph.NewVertex(string(member.PKIid), member)
	}

	
	principalVertices := make(map[string]*graph.Vertex)
	for pKey, grp := range data.pGrps {
		principalVertices[grp] = graph.NewVertex(grp, pKey.toPrincipal())
	}

	
	for _, principalVertex := range principalVertices {
		for _, peerVertex := range peerVertices {
			
			principal := principalVertex.Data.(*msp.MSPPrincipal)
			member := peerVertex.Data.(NetworkMember)
			if satisfiesPrincipal(member, principal) {
				peerVertex.AddNeighbor(principalVertex)
			}
		}
	}
	return &principalPeerGraph{
		peerVertices:      peerVertices,
		principalVertices: principalVertices,
	}
}

func mapPrincipalsToGroups(principalsSets []policies.PrincipalSet) principalGroupMapper {
	groupMapper := make(principalGroupMapper)
	totalPrincipals := make(map[principalKey]struct{})
	for _, principalSet := range principalsSets {
		for _, principal := range principalSet {
			totalPrincipals[principalKey{
				principal: string(principal.Principal),
				cls:       int32(principal.PrincipalClassification),
			}] = struct{}{}
		}
	}
	for principal := range totalPrincipals {
		groupMapper.group(principal)
	}
	return groupMapper
}

type memberIdentities map[string]api.PeerIdentityType

func (m memberIdentities) identityByPKIID(id common.PKIidType) api.PeerIdentityType {
	return m[string(id)]
}

func computeIdentitiesOfMembers(identitySet api.PeerIdentitySet, members map[string]NetworkMember) memberIdentities {
	identitiesByPKIID := make(map[string]api.PeerIdentityType)
	identitiesOfMembers := make(map[string]api.PeerIdentityType, len(members))
	for _, identity := range identitySet {
		identitiesByPKIID[string(identity.PKIId)] = identity.Identity
	}
	for _, member := range members {
		if identity, exists := identitiesByPKIID[string(member.PKIid)]; exists {
			identitiesOfMembers[string(member.PKIid)] = identity
		}
	}
	return identitiesOfMembers
}


type principalGroupMapper map[principalKey]string

func (mapper principalGroupMapper) group(principal principalKey) string {
	if grp, exists := mapper[principal]; exists {
		return grp
	}
	grp := fmt.Sprintf("G%d", len(mapper))
	mapper[principal] = grp
	return grp
}

type principalKey struct {
	cls       int32
	principal string
}

func (pk principalKey) toPrincipal() *msp.MSPPrincipal {
	return &msp.MSPPrincipal{
		PrincipalClassification: msp.MSPPrincipal_Classification(pk.cls),
		Principal:               []byte(pk.principal),
	}
}


type layouts []*discovery.Layout


func (l layouts) groupsSet() map[string]struct{} {
	m := make(map[string]struct{})
	for _, layout := range l {
		for grp := range layout.QuantitiesByGroup {
			m[grp] = struct{}{}
		}
	}
	return m
}

func peersWithChaincode(metadata ...*chaincode.Metadata) func(member NetworkMember) bool {
	return func(member NetworkMember) bool {
		if member.Properties == nil {
			return false
		}
		for _, ccMD := range metadata {
			var found bool
			for _, cc := range member.Properties.Chaincodes {
				if cc.Name == ccMD.Name && cc.Version == ccMD.Version {
					found = true
				}
			}
			if !found {
				return false
			}
		}
		return true
	}
}

func mergePrincipalSets(cpss []inquire.ComparablePrincipalSets) (inquire.ComparablePrincipalSets, error) {
	
	var cps inquire.ComparablePrincipalSets
	cps, cpss, err := popComparablePrincipalSets(cpss)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for _, cps2 := range cpss {
		cps = inquire.Merge(cps, cps2)
	}
	return cps, nil
}

func popComparablePrincipalSets(sets []inquire.ComparablePrincipalSets) (inquire.ComparablePrincipalSets, []inquire.ComparablePrincipalSets, error) {
	if len(sets) == 0 {
		return nil, nil, errors.New("no principal sets remained after filtering")
	}
	cps, cpss := sets[0], sets[1:]
	return cps, cpss, nil
}
