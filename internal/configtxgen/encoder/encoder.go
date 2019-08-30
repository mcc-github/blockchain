/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package encoder

import (
	"github.com/gogo/protobuf/proto"
	cb "github.com/mcc-github/blockchain-protos-go/common"
	pb "github.com/mcc-github/blockchain-protos-go/peer"
	"github.com/mcc-github/blockchain/common/cauthdsl"
	"github.com/mcc-github/blockchain/common/channelconfig"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/common/genesis"
	"github.com/mcc-github/blockchain/common/policies"
	"github.com/mcc-github/blockchain/common/util"
	genesisconfig "github.com/mcc-github/blockchain/internal/configtxgen/localconfig"
	"github.com/mcc-github/blockchain/internal/configtxlator/update"
	"github.com/mcc-github/blockchain/internal/pkg/identity"
	"github.com/mcc-github/blockchain/msp"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/pkg/errors"
)

const (
	ordererAdminsPolicyName = "/Channel/Orderer/Admins"

	msgVersion = int32(0)
	epoch      = 0
)

var logger = flogging.MustGetLogger("common.tools.configtxgen.encoder")

const (
	
	ConsensusTypeSolo = "solo"
	
	ConsensusTypeKafka = "kafka"
	
	ConsensusTypeEtcdRaft = "etcdraft"

	
	BlockValidationPolicyKey = "BlockValidation"

	
	OrdererAdminsPolicy = "/Channel/Orderer/Admins"

	
	SignaturePolicyType = "Signature"

	
	ImplicitMetaPolicyType = "ImplicitMeta"
)

func addValue(cg *cb.ConfigGroup, value channelconfig.ConfigValue, modPolicy string) {
	cg.Values[value.Key()] = &cb.ConfigValue{
		Value:     protoutil.MarshalOrPanic(value.Value()),
		ModPolicy: modPolicy,
	}
}

func addPolicy(cg *cb.ConfigGroup, policy policies.ConfigPolicy, modPolicy string) {
	cg.Policies[policy.Key()] = &cb.ConfigPolicy{
		Policy:    policy.Value(),
		ModPolicy: modPolicy,
	}
}

func AddPolicies(cg *cb.ConfigGroup, policyMap map[string]*genesisconfig.Policy, modPolicy string) error {
	switch {
	case policyMap == nil:
		return errors.Errorf("no policies defined")
	case policyMap[channelconfig.AdminsPolicyKey] == nil:
		return errors.Errorf("no Admins policy defined")
	case policyMap[channelconfig.ReadersPolicyKey] == nil:
		return errors.Errorf("no Readers policy defined")
	case policyMap[channelconfig.WritersPolicyKey] == nil:
		return errors.Errorf("no Writers policy defined")
	}

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
					Value: protoutil.MarshalOrPanic(imp),
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
					Value: protoutil.MarshalOrPanic(sp),
				},
			}
		default:
			return errors.Errorf("unknown policy type: %s", policy.Type)
		}
	}
	return nil
}






func NewChannelGroup(conf *genesisconfig.Profile) (*cb.ConfigGroup, error) {
	channelGroup := protoutil.NewConfigGroup()
	if err := AddPolicies(channelGroup, conf.Policies, channelconfig.AdminsPolicyKey); err != nil {
		return nil, errors.Wrapf(err, "error adding policies to channel group")
	}

	addValue(channelGroup, channelconfig.HashingAlgorithmValue(), channelconfig.AdminsPolicyKey)
	addValue(channelGroup, channelconfig.BlockDataHashingStructureValue(), channelconfig.AdminsPolicyKey)
	if conf.Orderer != nil && len(conf.Orderer.Addresses) > 0 {
		addValue(channelGroup, channelconfig.OrdererAddressesValue(conf.Orderer.Addresses), ordererAdminsPolicyName)
	}

	if conf.Consortium != "" {
		addValue(channelGroup, channelconfig.ConsortiumValue(conf.Consortium), channelconfig.AdminsPolicyKey)
	}

	if len(conf.Capabilities) > 0 {
		addValue(channelGroup, channelconfig.CapabilitiesValue(conf.Capabilities), channelconfig.AdminsPolicyKey)
	}

	var err error
	if conf.Orderer != nil {
		channelGroup.Groups[channelconfig.OrdererGroupKey], err = NewOrdererGroup(conf.Orderer)
		if err != nil {
			return nil, errors.Wrap(err, "could not create orderer group")
		}
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
	ordererGroup := protoutil.NewConfigGroup()
	if err := AddPolicies(ordererGroup, conf.Policies, channelconfig.AdminsPolicyKey); err != nil {
		return nil, errors.Wrapf(err, "error adding policies to orderer group")
	}
	ordererGroup.Policies[BlockValidationPolicyKey] = &cb.ConfigPolicy{
		Policy:    policies.ImplicitMetaAnyPolicy(channelconfig.WritersPolicyKey).Value(),
		ModPolicy: channelconfig.AdminsPolicyKey,
	}
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

	var consensusMetadata []byte
	var err error

	switch conf.OrdererType {
	case ConsensusTypeSolo:
	case ConsensusTypeKafka:
		addValue(ordererGroup, channelconfig.KafkaBrokersValue(conf.Kafka.Brokers), channelconfig.AdminsPolicyKey)
	case ConsensusTypeEtcdRaft:
		if consensusMetadata, err = channelconfig.MarshalEtcdRaftMetadata(conf.EtcdRaft); err != nil {
			return nil, errors.Errorf("cannot marshal metadata for orderer type %s: %s", ConsensusTypeEtcdRaft, err)
		}
	default:
		return nil, errors.Errorf("unknown orderer type: %s", conf.OrdererType)
	}

	addValue(ordererGroup, channelconfig.ConsensusTypeValue(conf.OrdererType, consensusMetadata), channelconfig.AdminsPolicyKey)

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



func NewConsortiumOrgGroup(conf *genesisconfig.Organization) (*cb.ConfigGroup, error) {
	consortiumsOrgGroup := protoutil.NewConfigGroup()
	consortiumsOrgGroup.ModPolicy = channelconfig.AdminsPolicyKey

	if conf.SkipAsForeign {
		return consortiumsOrgGroup, nil
	}

	mspConfig, err := msp.GetVerifyingMspConfig(conf.MSPDir, conf.ID, conf.MSPType)
	if err != nil {
		return nil, errors.Wrapf(err, "1 - Error loading MSP configuration for org: %s", conf.Name)
	}

	if err := AddPolicies(consortiumsOrgGroup, conf.Policies, channelconfig.AdminsPolicyKey); err != nil {
		return nil, errors.Wrapf(err, "error adding policies to consortiums org group '%s'", conf.Name)
	}

	addValue(consortiumsOrgGroup, channelconfig.MSPValue(mspConfig), channelconfig.AdminsPolicyKey)

	return consortiumsOrgGroup, nil
}



func NewOrdererOrgGroup(conf *genesisconfig.Organization) (*cb.ConfigGroup, error) {
	ordererOrgGroup := protoutil.NewConfigGroup()
	ordererOrgGroup.ModPolicy = channelconfig.AdminsPolicyKey

	if conf.SkipAsForeign {
		return ordererOrgGroup, nil
	}

	mspConfig, err := msp.GetVerifyingMspConfig(conf.MSPDir, conf.ID, conf.MSPType)
	if err != nil {
		return nil, errors.Wrapf(err, "1 - Error loading MSP configuration for org: %s", conf.Name)
	}

	if err := AddPolicies(ordererOrgGroup, conf.Policies, channelconfig.AdminsPolicyKey); err != nil {
		return nil, errors.Wrapf(err, "error adding policies to orderer org group '%s'", conf.Name)
	}

	addValue(ordererOrgGroup, channelconfig.MSPValue(mspConfig), channelconfig.AdminsPolicyKey)

	if len(conf.OrdererEndpoints) > 0 {
		addValue(ordererOrgGroup, channelconfig.EndpointsValue(conf.OrdererEndpoints), channelconfig.AdminsPolicyKey)
	}

	return ordererOrgGroup, nil
}



func NewApplicationGroup(conf *genesisconfig.Application) (*cb.ConfigGroup, error) {
	applicationGroup := protoutil.NewConfigGroup()
	if err := AddPolicies(applicationGroup, conf.Policies, channelconfig.AdminsPolicyKey); err != nil {
		return nil, errors.Wrapf(err, "error adding policies to application group")
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
	applicationOrgGroup := protoutil.NewConfigGroup()
	applicationOrgGroup.ModPolicy = channelconfig.AdminsPolicyKey

	if conf.SkipAsForeign {
		return applicationOrgGroup, nil
	}

	mspConfig, err := msp.GetVerifyingMspConfig(conf.MSPDir, conf.ID, conf.MSPType)
	if err != nil {
		return nil, errors.Wrapf(err, "1 - Error loading MSP configuration for org %s", conf.Name)
	}

	if err := AddPolicies(applicationOrgGroup, conf.Policies, channelconfig.AdminsPolicyKey); err != nil {
		return nil, errors.Wrapf(err, "error adding policies to application org group %s", conf.Name)
	}
	addValue(applicationOrgGroup, channelconfig.MSPValue(mspConfig), channelconfig.AdminsPolicyKey)

	var anchorProtos []*pb.AnchorPeer
	for _, anchorPeer := range conf.AnchorPeers {
		anchorProtos = append(anchorProtos, &pb.AnchorPeer{
			Host: anchorPeer.Host,
			Port: int32(anchorPeer.Port),
		})
	}

	
	
	
	if len(anchorProtos) > 0 {
		addValue(applicationOrgGroup, channelconfig.AnchorPeersValue(anchorProtos), channelconfig.AdminsPolicyKey)
	}

	return applicationOrgGroup, nil
}



func NewConsortiumsGroup(conf map[string]*genesisconfig.Consortium) (*cb.ConfigGroup, error) {
	consortiumsGroup := protoutil.NewConfigGroup()
	
	
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
	consortiumGroup := protoutil.NewConfigGroup()

	for _, org := range conf.Organizations {
		var err error
		consortiumGroup.Groups[org.Name], err = NewConsortiumOrgGroup(org)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create consortium org")
		}
	}

	addValue(consortiumGroup, channelconfig.ChannelCreationPolicyValue(policies.ImplicitMetaAnyPolicy(channelconfig.AdminsPolicyKey).Value()), ordererAdminsPolicyName)

	consortiumGroup.ModPolicy = ordererAdminsPolicyName
	return consortiumGroup, nil
}



func NewChannelCreateConfigUpdate(channelID string, conf *genesisconfig.Profile, templateConfig *cb.ConfigGroup) (*cb.ConfigUpdate, error) {
	if conf.Application == nil {
		return nil, errors.New("cannot define a new channel with no Application section")
	}

	if conf.Consortium == "" {
		return nil, errors.New("cannot define a new channel with no Consortium value")
	}

	newChannelGroup, err := NewChannelGroup(conf)
	if err != nil {
		return nil, errors.Wrapf(err, "could not turn parse profile into channel group")
	}

	updt, err := update.Compute(&cb.Config{ChannelGroup: templateConfig}, &cb.Config{ChannelGroup: newChannelGroup})
	if err != nil {
		return nil, errors.Wrapf(err, "could not compute update")
	}

	
	updt.ChannelId = channelID
	updt.ReadSet.Values[channelconfig.ConsortiumKey] = &cb.ConfigValue{Version: 0}
	updt.WriteSet.Values[channelconfig.ConsortiumKey] = &cb.ConfigValue{
		Version: 0,
		Value: protoutil.MarshalOrPanic(&cb.Consortium{
			Name: conf.Consortium,
		}),
	}

	return updt, nil
}




func DefaultConfigTemplate(conf *genesisconfig.Profile) (*cb.ConfigGroup, error) {
	channelGroup, err := NewChannelGroup(conf)
	if err != nil {
		return nil, errors.WithMessage(err, "error parsing configuration")
	}

	if _, ok := channelGroup.Groups[channelconfig.ApplicationGroupKey]; !ok {
		return nil, errors.New("channel template configs must contain an application section")
	}

	channelGroup.Groups[channelconfig.ApplicationGroupKey].Values = nil
	channelGroup.Groups[channelconfig.ApplicationGroupKey].Policies = nil

	return channelGroup, nil
}

func ConfigTemplateFromGroup(conf *genesisconfig.Profile, cg *cb.ConfigGroup) (*cb.ConfigGroup, error) {
	template := proto.Clone(cg).(*cb.ConfigGroup)
	if template.Groups == nil {
		return nil, errors.Errorf("supplied system channel group has no sub-groups")
	}

	template.Groups[channelconfig.ApplicationGroupKey] = &cb.ConfigGroup{
		Groups: map[string]*cb.ConfigGroup{},
		Policies: map[string]*cb.ConfigPolicy{
			channelconfig.AdminsPolicyKey: {},
		},
	}

	consortiums, ok := template.Groups[channelconfig.ConsortiumsGroupKey]
	if !ok {
		return nil, errors.Errorf("supplied system channel group does not appear to be system channel (missing consortiums group)")
	}

	if consortiums.Groups == nil {
		return nil, errors.Errorf("system channel consortiums group appears to have no consortiums defined")
	}

	consortium, ok := consortiums.Groups[conf.Consortium]
	if !ok {
		return nil, errors.Errorf("supplied system channel group is missing '%s' consortium", conf.Consortium)
	}

	if conf.Application == nil {
		return nil, errors.Errorf("supplied channel creation profile does not contain an application section")
	}

	for _, organization := range conf.Application.Organizations {
		var ok bool
		template.Groups[channelconfig.ApplicationGroupKey].Groups[organization.Name], ok = consortium.Groups[organization.Name]
		if !ok {
			return nil, errors.Errorf("consortium %s does not contain member org %s", conf.Consortium, organization.Name)
		}
	}
	delete(template.Groups, channelconfig.ConsortiumsGroupKey)

	addValue(template, channelconfig.ConsortiumValue(conf.Consortium), channelconfig.AdminsPolicyKey)

	return template, nil
}



func MakeChannelCreationTransaction(
	channelID string,
	signer identity.SignerSerializer,
	conf *genesisconfig.Profile,
) (*cb.Envelope, error) {
	template, err := DefaultConfigTemplate(conf)
	if err != nil {
		return nil, errors.WithMessage(err, "could not generate default config template")
	}
	return MakeChannelCreationTransactionFromTemplate(channelID, signer, conf, template)
}




func MakeChannelCreationTransactionWithSystemChannelContext(
	channelID string,
	signer identity.SignerSerializer,
	conf,
	systemChannelConf *genesisconfig.Profile,
) (*cb.Envelope, error) {
	cg, err := NewChannelGroup(systemChannelConf)
	if err != nil {
		return nil, errors.WithMessage(err, "could not parse system channel config")
	}

	template, err := ConfigTemplateFromGroup(conf, cg)
	if err != nil {
		return nil, errors.WithMessage(err, "could not create config template")
	}

	return MakeChannelCreationTransactionFromTemplate(channelID, signer, conf, template)
}




func MakeChannelCreationTransactionFromTemplate(
	channelID string,
	signer identity.SignerSerializer,
	conf *genesisconfig.Profile,
	template *cb.ConfigGroup,
) (*cb.Envelope, error) {
	newChannelConfigUpdate, err := NewChannelCreateConfigUpdate(channelID, conf, template)
	if err != nil {
		return nil, errors.Wrap(err, "config update generation failure")
	}

	newConfigUpdateEnv := &cb.ConfigUpdateEnvelope{
		ConfigUpdate: protoutil.MarshalOrPanic(newChannelConfigUpdate),
	}

	if signer != nil {
		sigHeader, err := protoutil.NewSignatureHeader(signer)
		if err != nil {
			return nil, errors.Wrap(err, "creating signature header failed")
		}

		newConfigUpdateEnv.Signatures = []*cb.ConfigSignature{{
			SignatureHeader: protoutil.MarshalOrPanic(sigHeader),
		}}

		newConfigUpdateEnv.Signatures[0].Signature, err = signer.Sign(util.ConcatenateBytes(newConfigUpdateEnv.Signatures[0].SignatureHeader, newConfigUpdateEnv.ConfigUpdate))
		if err != nil {
			return nil, errors.Wrap(err, "signature failure over config update")
		}

	}

	return protoutil.CreateSignedEnvelope(cb.HeaderType_CONFIG_UPDATE, channelID, signer, newConfigUpdateEnv, msgVersion, epoch)
}




func HasSkippedForeignOrgs(conf *genesisconfig.Profile) error {
	var organizations []*genesisconfig.Organization

	if conf.Orderer != nil {
		organizations = append(organizations, conf.Orderer.Organizations...)
	}

	if conf.Application != nil {
		organizations = append(organizations, conf.Application.Organizations...)
	}

	for _, consortium := range conf.Consortiums {
		organizations = append(organizations, consortium.Organizations...)
	}

	for _, org := range organizations {
		if org.SkipAsForeign {
			return errors.Errorf("organization '%s' is marked to be skipped as foreign", org.Name)
		}
	}

	return nil
}


type Bootstrapper struct {
	channelGroup *cb.ConfigGroup
}


func NewBootstrapper(config *genesisconfig.Profile) (*Bootstrapper, error) {
	if err := HasSkippedForeignOrgs(config); err != nil {
		return nil, errors.WithMessage(err, "all org definitions must be local during bootstrapping")
	}

	channelGroup, err := NewChannelGroup(config)
	if err != nil {
		return nil, errors.WithMessage(err, "could not create channel group")
	}

	return &Bootstrapper{
		channelGroup: channelGroup,
	}, nil
}


func New(config *genesisconfig.Profile) *Bootstrapper {
	bs, err := NewBootstrapper(config)
	if err != nil {
		logger.Panicf("Error creating bootsrapper: %s", err)
	}
	return bs
}


func (bs *Bootstrapper) GenesisBlock() *cb.Block {
	return genesis.NewFactoryImpl(bs.channelGroup).Block(genesisconfig.TestChainID)
}


func (bs *Bootstrapper) GenesisBlockForChannel(channelID string) *cb.Block {
	return genesis.NewFactoryImpl(bs.channelGroup).Block(channelID)
}
