/*
Copyright IBM Corp. 2016 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package encoder

import (
	"github.com/mcc-github/blockchain/common/cauthdsl"
	"github.com/mcc-github/blockchain/common/channelconfig"
	"github.com/mcc-github/blockchain/common/crypto"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/common/genesis"
	"github.com/mcc-github/blockchain/common/policies"
	genesisconfig "github.com/mcc-github/blockchain/common/tools/configtxgen/localconfig"
	"github.com/mcc-github/blockchain/common/tools/configtxlator/update"
	"github.com/mcc-github/blockchain/common/util"
	"github.com/mcc-github/blockchain/msp"
	cb "github.com/mcc-github/blockchain/protos/common"
	pb "github.com/mcc-github/blockchain/protos/peer"
	"github.com/mcc-github/blockchain/protos/utils"

	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
)

const (
	pkgLogID                = "common/tools/configtxgen/encoder"
	ordererAdminsPolicyName = "/Channel/Orderer/Admins"

	msgVersion = int32(0)
	epoch      = 0
)

var logger = flogging.MustGetLogger(pkgLogID)

func init() {
	flogging.SetModuleLevel(pkgLogID, "info")
}

const (
	
	ConsensusTypeSolo = "solo"
	
	ConsensusTypeKafka = "kafka"

	
	BlockValidationPolicyKey = "BlockValidation"

	
	OrdererAdminsPolicy = "/Channel/Orderer/Admins"

	
	SignaturePolicyType = "Signature"

	
	ImplicitMetaPolicyType = "ImplicitMeta"
)

func addValue(cg *cb.ConfigGroup, value channelconfig.ConfigValue, modPolicy string) {
	cg.Values[value.Key()] = &cb.ConfigValue{
		Value:     utils.MarshalOrPanic(value.Value()),
		ModPolicy: modPolicy,
	}
}

func addPolicy(cg *cb.ConfigGroup, policy policies.ConfigPolicy, modPolicy string) {
	cg.Policies[policy.Key()] = &cb.ConfigPolicy{
		Policy:    policy.Value(),
		ModPolicy: modPolicy,
	}
}

func addPolicies(cg *cb.ConfigGroup, policyMap map[string]*genesisconfig.Policy, modPolicy string) error {
	for policyName, policy := range policyMap {
		switch policy.Type {
		case ImplicitMetaPolicyType:
			imp, err := policies.ImplicitMetaFromString(policy.Rule)
			if err != nil {
				return errors.Wrapf(err, "invalid implicit meta policy rule '%s'", policy.Rule)
			}
			cg.Policies[policyName] = &cb.ConfigPolicy{
				ModPolicy: modPolicy,
				Policy: &cb.Policy{
					Type:  int32(cb.Policy_IMPLICIT_META),
					Value: utils.MarshalOrPanic(imp),
				},
			}
		case SignaturePolicyType:
			sp, err := cauthdsl.FromString(policy.Rule)
			if err != nil {
				return errors.Wrapf(err, "invalid signature policy rule '%s'", policy.Rule)
			}
			cg.Policies[policyName] = &cb.ConfigPolicy{
				ModPolicy: modPolicy,
				Policy: &cb.Policy{
					Type:  int32(cb.Policy_SIGNATURE),
					Value: utils.MarshalOrPanic(sp),
				},
			}
		default:
			return errors.Errorf("unknown policy type: %s", policy.Type)
		}
	}
	return nil
}


func addImplicitMetaPolicyDefaults(cg *cb.ConfigGroup) {
	addPolicy(cg, policies.ImplicitMetaMajorityPolicy(channelconfig.AdminsPolicyKey), channelconfig.AdminsPolicyKey)
	addPolicy(cg, policies.ImplicitMetaAnyPolicy(channelconfig.ReadersPolicyKey), channelconfig.AdminsPolicyKey)
	addPolicy(cg, policies.ImplicitMetaAnyPolicy(channelconfig.WritersPolicyKey), channelconfig.AdminsPolicyKey)
}




func addSignaturePolicyDefaults(cg *cb.ConfigGroup, mspID string, devMode bool) {
	if devMode {
		logger.Warningf("Specifying AdminPrincipal is deprecated and will be removed in a future release, override the admin principal with explicit policies.")
		addPolicy(cg, policies.SignaturePolicy(channelconfig.AdminsPolicyKey, cauthdsl.SignedByMspMember(mspID)), channelconfig.AdminsPolicyKey)
	} else {
		addPolicy(cg, policies.SignaturePolicy(channelconfig.AdminsPolicyKey, cauthdsl.SignedByMspAdmin(mspID)), channelconfig.AdminsPolicyKey)
	}
	addPolicy(cg, policies.SignaturePolicy(channelconfig.ReadersPolicyKey, cauthdsl.SignedByMspMember(mspID)), channelconfig.AdminsPolicyKey)
	addPolicy(cg, policies.SignaturePolicy(channelconfig.WritersPolicyKey, cauthdsl.SignedByMspMember(mspID)), channelconfig.AdminsPolicyKey)
}






func NewChannelGroup(conf *genesisconfig.Profile) (*cb.ConfigGroup, error) {
	if conf.Orderer == nil {
		return nil, errors.New("missing orderer config section")
	}

	channelGroup := cb.NewConfigGroup()
	if len(conf.Policies) == 0 {
		logger.Warningf("Default policy emission is deprecated, please include policy specificiations for the channel group in configtx.yaml")
		addImplicitMetaPolicyDefaults(channelGroup)
	} else {
		if err := addPolicies(channelGroup, conf.Policies, channelconfig.AdminsPolicyKey); err != nil {
			return nil, errors.Wrapf(err, "error adding policies to channel group")
		}
	}

	addValue(channelGroup, channelconfig.HashingAlgorithmValue(), channelconfig.AdminsPolicyKey)
	addValue(channelGroup, channelconfig.BlockDataHashingStructureValue(), channelconfig.AdminsPolicyKey)
	addValue(channelGroup, channelconfig.OrdererAddressesValue(conf.Orderer.Addresses), ordererAdminsPolicyName)

	if conf.Consortium != "" {
		addValue(channelGroup, channelconfig.ConsortiumValue(conf.Consortium), channelconfig.AdminsPolicyKey)
	}

	if len(conf.Capabilities) > 0 {
		addValue(channelGroup, channelconfig.CapabilitiesValue(conf.Capabilities), channelconfig.AdminsPolicyKey)
	}

	var err error
	channelGroup.Groups[channelconfig.OrdererGroupKey], err = NewOrdererGroup(conf.Orderer)
	if err != nil {
		return nil, errors.Wrap(err, "could not create orderer group")
	}

	if conf.Application != nil {
		channelGroup.Groups[channelconfig.ApplicationGroupKey], err = NewApplicationGroup(conf.Application)
		if err != nil {
			return nil, errors.Wrap(err, "could not create application group")
		}
	}

	if conf.Consortiums != nil {
		channelGroup.Groups[channelconfig.ConsortiumsGroupKey], err = NewConsortiumsGroup(conf.Consortiums)
		if err != nil {
			return nil, errors.Wrap(err, "could not create consortiums group")
		}
	}

	channelGroup.ModPolicy = channelconfig.AdminsPolicyKey
	return channelGroup, nil
}




func NewOrdererGroup(conf *genesisconfig.Orderer) (*cb.ConfigGroup, error) {
	ordererGroup := cb.NewConfigGroup()
	if len(conf.Policies) == 0 {
		logger.Warningf("Default policy emission is deprecated, please include policy specificiations for the orderer group in configtx.yaml")
		addImplicitMetaPolicyDefaults(ordererGroup)
	} else {
		if err := addPolicies(ordererGroup, conf.Policies, channelconfig.AdminsPolicyKey); err != nil {
			return nil, errors.Wrapf(err, "error adding policies to orderer group")
		}
	}
	ordererGroup.Policies[BlockValidationPolicyKey] = &cb.ConfigPolicy{
		Policy:    policies.ImplicitMetaAnyPolicy(channelconfig.WritersPolicyKey).Value(),
		ModPolicy: channelconfig.AdminsPolicyKey,
	}
	addValue(ordererGroup, channelconfig.ConsensusTypeValue(conf.OrdererType), channelconfig.AdminsPolicyKey)
	addValue(ordererGroup, channelconfig.BatchSizeValue(
		conf.BatchSize.MaxMessageCount,
		conf.BatchSize.AbsoluteMaxBytes,
		conf.BatchSize.PreferredMaxBytes,
	), channelconfig.AdminsPolicyKey)
	addValue(ordererGroup, channelconfig.BatchTimeoutValue(conf.BatchTimeout.String()), channelconfig.AdminsPolicyKey)
	addValue(ordererGroup, channelconfig.ChannelRestrictionsValue(conf.MaxChannels), channelconfig.AdminsPolicyKey)

	if len(conf.Capabilities) > 0 {
		addValue(ordererGroup, channelconfig.CapabilitiesValue(conf.Capabilities), channelconfig.AdminsPolicyKey)
	}

	switch conf.OrdererType {
	case ConsensusTypeSolo:
	case ConsensusTypeKafka:
		addValue(ordererGroup, channelconfig.KafkaBrokersValue(conf.Kafka.Brokers), channelconfig.AdminsPolicyKey)
	default:
		return nil, errors.Errorf("unknown orderer type: %s", conf.OrdererType)
	}

	for _, org := range conf.Organizations {
		var err error
		ordererGroup.Groups[org.Name], err = NewOrdererOrgGroup(org)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create orderer org")
		}
	}

	ordererGroup.ModPolicy = channelconfig.AdminsPolicyKey
	return ordererGroup, nil
}



func NewOrdererOrgGroup(conf *genesisconfig.Organization) (*cb.ConfigGroup, error) {
	mspConfig, err := msp.GetVerifyingMspConfig(conf.MSPDir, conf.ID, conf.MSPType)
	if err != nil {
		return nil, errors.Wrapf(err, "1 - Error loading MSP configuration for org: %s", conf.Name)
	}

	ordererOrgGroup := cb.NewConfigGroup()
	if len(conf.Policies) == 0 {
		logger.Warningf("Default policy emission is deprecated, please include policy specificiations for the orderer org group %s in configtx.yaml", conf.Name)
		addSignaturePolicyDefaults(ordererOrgGroup, conf.ID, conf.AdminPrincipal != genesisconfig.AdminRoleAdminPrincipal)
	} else {
		if err := addPolicies(ordererOrgGroup, conf.Policies, channelconfig.AdminsPolicyKey); err != nil {
			return nil, errors.Wrapf(err, "error adding policies to orderer org group '%s'", conf.Name)
		}
	}

	addValue(ordererOrgGroup, channelconfig.MSPValue(mspConfig), channelconfig.AdminsPolicyKey)

	ordererOrgGroup.ModPolicy = channelconfig.AdminsPolicyKey
	return ordererOrgGroup, nil
}



func NewApplicationGroup(conf *genesisconfig.Application) (*cb.ConfigGroup, error) {
	applicationGroup := cb.NewConfigGroup()
	if len(conf.Policies) == 0 {
		logger.Warningf("Default policy emission is deprecated, please include policy specificiations for the application group in configtx.yaml")
		addImplicitMetaPolicyDefaults(applicationGroup)
	} else {
		if err := addPolicies(applicationGroup, conf.Policies, channelconfig.AdminsPolicyKey); err != nil {
			return nil, errors.Wrapf(err, "error adding policies to application group")
		}
	}

	if len(conf.ACLs) > 0 {
		addValue(applicationGroup, channelconfig.ACLValues(conf.ACLs), channelconfig.AdminsPolicyKey)
	}

	if len(conf.Capabilities) > 0 {
		addValue(applicationGroup, channelconfig.CapabilitiesValue(conf.Capabilities), channelconfig.AdminsPolicyKey)
	}

	for _, org := range conf.Organizations {
		var err error
		applicationGroup.Groups[org.Name], err = NewApplicationOrgGroup(org)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create application org")
		}
	}

	applicationGroup.ModPolicy = channelconfig.AdminsPolicyKey
	return applicationGroup, nil
}



func NewApplicationOrgGroup(conf *genesisconfig.Organization) (*cb.ConfigGroup, error) {
	mspConfig, err := msp.GetVerifyingMspConfig(conf.MSPDir, conf.ID, conf.MSPType)
	if err != nil {
		return nil, errors.Wrapf(err, "1 - Error loading MSP configuration for org %s", conf.Name)
	}

	applicationOrgGroup := cb.NewConfigGroup()
	if len(conf.Policies) == 0 {
		logger.Warningf("Default policy emission is deprecated, please include policy specificiations for the application org group %s in configtx.yaml", conf.Name)
		addSignaturePolicyDefaults(applicationOrgGroup, conf.ID, conf.AdminPrincipal != genesisconfig.AdminRoleAdminPrincipal)
	} else {
		if err := addPolicies(applicationOrgGroup, conf.Policies, channelconfig.AdminsPolicyKey); err != nil {
			return nil, errors.Wrapf(err, "error adding policies to application org group %s", conf.Name)
		}
	}
	addValue(applicationOrgGroup, channelconfig.MSPValue(mspConfig), channelconfig.AdminsPolicyKey)

	var anchorProtos []*pb.AnchorPeer
	for _, anchorPeer := range conf.AnchorPeers {
		anchorProtos = append(anchorProtos, &pb.AnchorPeer{
			Host: anchorPeer.Host,
			Port: int32(anchorPeer.Port),
		})
	}
	addValue(applicationOrgGroup, channelconfig.AnchorPeersValue(anchorProtos), channelconfig.AdminsPolicyKey)

	applicationOrgGroup.ModPolicy = channelconfig.AdminsPolicyKey
	return applicationOrgGroup, nil
}



func NewConsortiumsGroup(conf map[string]*genesisconfig.Consortium) (*cb.ConfigGroup, error) {
	consortiumsGroup := cb.NewConfigGroup()
	
	
	addPolicy(consortiumsGroup, policies.SignaturePolicy(channelconfig.AdminsPolicyKey, cauthdsl.AcceptAllPolicy), ordererAdminsPolicyName)

	for consortiumName, consortium := range conf {
		var err error
		consortiumsGroup.Groups[consortiumName], err = NewConsortiumGroup(consortium)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create consortium %s", consortiumName)
		}
	}

	consortiumsGroup.ModPolicy = ordererAdminsPolicyName
	return consortiumsGroup, nil
}




func NewConsortiumGroup(conf *genesisconfig.Consortium) (*cb.ConfigGroup, error) {
	consortiumGroup := cb.NewConfigGroup()

	for _, org := range conf.Organizations {
		var err error
		
		consortiumGroup.Groups[org.Name], err = NewOrdererOrgGroup(org)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create consortium org")
		}
	}

	addValue(consortiumGroup, channelconfig.ChannelCreationPolicyValue(policies.ImplicitMetaAnyPolicy(channelconfig.AdminsPolicyKey).Value()), ordererAdminsPolicyName)

	consortiumGroup.ModPolicy = ordererAdminsPolicyName
	return consortiumGroup, nil
}



func NewChannelCreateConfigUpdate(channelID string, orderingSystemChannelGroup *cb.ConfigGroup, conf *genesisconfig.Profile) (*cb.ConfigUpdate, error) {
	if conf.Application == nil {
		return nil, errors.New("cannot define a new channel with no Application section")
	}

	if conf.Consortium == "" {
		return nil, errors.New("cannot define a new channel with no Consortium value")
	}

	
	ag, err := NewApplicationGroup(conf.Application)
	if err != nil {
		return nil, errors.Wrapf(err, "could not turn channel application profile into application group")
	}

	var template, newChannelGroup *cb.ConfigGroup

	if orderingSystemChannelGroup != nil {
		
		if orderingSystemChannelGroup.Groups == nil {
			return nil, errors.New("missing all channel groups")
		}

		consortiums, ok := orderingSystemChannelGroup.Groups[channelconfig.ConsortiumsGroupKey]
		if !ok {
			return nil, errors.New("bad consortiums group")
		}

		consortium, ok := consortiums.Groups[conf.Consortium]
		if !ok {
			return nil, errors.Errorf("bad consortium: %s", conf.Consortium)
		}

		template = proto.Clone(orderingSystemChannelGroup).(*cb.ConfigGroup)
		template.Groups[channelconfig.ApplicationGroupKey] = proto.Clone(consortium).(*cb.ConfigGroup)
		
		
		
		template.Groups[channelconfig.ApplicationGroupKey].Groups["*IllegalKey*!"] = &cb.ConfigGroup{}
		delete(template.Groups, channelconfig.ConsortiumsGroupKey)

		newChannelGroup = proto.Clone(orderingSystemChannelGroup).(*cb.ConfigGroup)
		delete(newChannelGroup.Groups, channelconfig.ConsortiumsGroupKey)
		newChannelGroup.Groups[channelconfig.ApplicationGroupKey].Values = ag.Values
		newChannelGroup.Groups[channelconfig.ApplicationGroupKey].Policies = ag.Policies

		for orgName, org := range template.Groups[channelconfig.ApplicationGroupKey].Groups {
			if _, ok := ag.Groups[orgName]; ok {
				newChannelGroup.Groups[channelconfig.ApplicationGroupKey].Groups[orgName] = org
			}
		}
	} else {
		newChannelGroup = &cb.ConfigGroup{
			Groups: map[string]*cb.ConfigGroup{
				channelconfig.ApplicationGroupKey: ag,
			},
		}

		
		template = proto.Clone(newChannelGroup).(*cb.ConfigGroup)
		template.Groups[channelconfig.ApplicationGroupKey].Values = nil
		template.Groups[channelconfig.ApplicationGroupKey].Policies = nil
	}

	updt, err := update.Compute(&cb.Config{ChannelGroup: template}, &cb.Config{ChannelGroup: newChannelGroup})
	if err != nil {
		return nil, errors.Wrapf(err, "could not compute update")
	}

	
	updt.ChannelId = channelID
	updt.ReadSet.Values[channelconfig.ConsortiumKey] = &cb.ConfigValue{Version: 0}
	updt.WriteSet.Values[channelconfig.ConsortiumKey] = &cb.ConfigValue{
		Version: 0,
		Value: utils.MarshalOrPanic(&cb.Consortium{
			Name: conf.Consortium,
		}),
	}

	return updt, nil
}


func MakeChannelCreationTransaction(channelID string, signer crypto.LocalSigner, orderingSystemChannelConfigGroup *cb.ConfigGroup, conf *genesisconfig.Profile) (*cb.Envelope, error) {
	newChannelConfigUpdate, err := NewChannelCreateConfigUpdate(channelID, orderingSystemChannelConfigGroup, conf)
	if err != nil {
		return nil, errors.Wrap(err, "config update generation failure")
	}

	newConfigUpdateEnv := &cb.ConfigUpdateEnvelope{
		ConfigUpdate: utils.MarshalOrPanic(newChannelConfigUpdate),
	}

	if signer != nil {
		sigHeader, err := signer.NewSignatureHeader()
		if err != nil {
			return nil, errors.Wrap(err, "creating signature header failed")
		}

		newConfigUpdateEnv.Signatures = []*cb.ConfigSignature{{
			SignatureHeader: utils.MarshalOrPanic(sigHeader),
		}}

		newConfigUpdateEnv.Signatures[0].Signature, err = signer.Sign(util.ConcatenateBytes(newConfigUpdateEnv.Signatures[0].SignatureHeader, newConfigUpdateEnv.ConfigUpdate))
		if err != nil {
			return nil, errors.Wrap(err, "signature failure over config update")
		}

	}

	return utils.CreateSignedEnvelope(cb.HeaderType_CONFIG_UPDATE, channelID, signer, newConfigUpdateEnv, msgVersion, epoch)
}


type Bootstrapper struct {
	channelGroup *cb.ConfigGroup
}


func New(config *genesisconfig.Profile) *Bootstrapper {
	channelGroup, err := NewChannelGroup(config)
	if err != nil {
		logger.Panicf("Error creating channel group: %s", err)
	}
	return &Bootstrapper{
		channelGroup: channelGroup,
	}
}


func (bs *Bootstrapper) GenesisBlock() *cb.Block {
	block, err := genesis.NewFactoryImpl(bs.channelGroup).Block(genesisconfig.TestChainID)
	if err != nil {
		logger.Panicf("Error creating genesis block from channel group: %s", err)
	}
	return block
}


func (bs *Bootstrapper) GenesisBlockForChannel(channelID string) *cb.Block {
	block, err := genesis.NewFactoryImpl(bs.channelGroup).Block(channelID)
	if err != nil {
		logger.Panicf("Error creating genesis block from channel group: %s", err)
	}
	return block
}
