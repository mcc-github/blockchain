/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package channelconfig

import (
	"time"

	"github.com/mcc-github/blockchain/common/configtx"
	"github.com/mcc-github/blockchain/common/policies"
	"github.com/mcc-github/blockchain/msp"
	cb "github.com/mcc-github/blockchain/protos/common"
	ab "github.com/mcc-github/blockchain/protos/orderer"
	pb "github.com/mcc-github/blockchain/protos/peer"
)


type Org interface {
	
	Name() string

	
	MSPID() string
}


type ApplicationOrg interface {
	Org

	
	AnchorPeers() []*pb.AnchorPeer
}


type OrdererOrg interface {
	Org

	
	Endpoints() []string
}


type Application interface {
	
	Organizations() map[string]ApplicationOrg

	
	APIPolicyMapper() PolicyMapper

	
	Capabilities() ApplicationCapabilities
}


type Channel interface {
	
	
	HashingAlgorithm() func(input []byte) []byte

	
	
	BlockDataHashingStructureWidth() uint32

	
	OrdererAddresses() []string

	
	Capabilities() ChannelCapabilities
}


type Consortiums interface {
	
	Consortiums() map[string]Consortium
}


type Consortium interface {
	
	ChannelCreationPolicy() *cb.Policy

	
	Organizations() map[string]Org
}


type Orderer interface {
	
	ConsensusType() string

	
	ConsensusMetadata() []byte

	
	ConsensusState() ab.ConsensusType_State

	
	BatchSize() *ab.BatchSize

	
	BatchTimeout() time.Duration

	
	MaxChannelsCount() uint64

	
	
	
	KafkaBrokers() []string

	
	Organizations() map[string]OrdererOrg

	
	Capabilities() OrdererCapabilities
}


type ChannelCapabilities interface {
	
	Supported() error

	
	
	MSPVersion() msp.MSPVersion

	
	ConsensusTypeMigration() bool

	
	OrgSpecificOrdererEndpoints() bool
}


type ApplicationCapabilities interface {
	
	Supported() error

	
	
	ForbidDuplicateTXIdInBlock() bool

	
	ACLs() bool

	
	
	
	PrivateChannelData() bool

	
	
	CollectionUpgrade() bool

	
	
	V1_1Validation() bool

	
	
	V1_2Validation() bool

	
	
	
	
	V1_3Validation() bool

	
	
	
	
	V2_0Validation() bool

	
	
	
	
	LifecycleV20() bool

	
	MetadataLifecycle() bool

	
	
	KeyLevelEndorsement() bool

	
	FabToken() bool
}


type OrdererCapabilities interface {
	
	
	PredictableChannelTemplate() bool

	
	
	Resubmission() bool

	
	Supported() error

	
	
	ExpirationCheck() bool

	
	ConsensusTypeMigration() bool

	
	
	
	UseChannelCreationPolicyAsAdmins() bool
}


type PolicyMapper interface {
	
	
	PolicyRefForAPI(apiName string) string
}




type Resources interface {
	
	ConfigtxValidator() configtx.Validator

	
	PolicyManager() policies.Manager

	
	ChannelConfig() Channel

	
	
	OrdererConfig() (Orderer, bool)

	
	
	ConsortiumsConfig() (Consortiums, bool)

	
	
	ApplicationConfig() (Application, bool)

	
	MSPManager() msp.MSPManager

	
	ValidateNew(resources Resources) error
}
