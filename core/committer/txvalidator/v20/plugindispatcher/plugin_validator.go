/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package plugindispatcher

import (
	"fmt"
	"sync"

	ledger2 "github.com/mcc-github/blockchain/common/ledger"
	"github.com/mcc-github/blockchain/common/policies"
	vp "github.com/mcc-github/blockchain/core/committer/txvalidator/plugin"
	validation "github.com/mcc-github/blockchain/core/handlers/validation/api"
	. "github.com/mcc-github/blockchain/core/handlers/validation/api/capabilities"
	. "github.com/mcc-github/blockchain/core/handlers/validation/api/identities"
	p "github.com/mcc-github/blockchain/core/handlers/validation/api/policies"
	. "github.com/mcc-github/blockchain/core/handlers/validation/api/state"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/policy"
	"github.com/mcc-github/blockchain/msp"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/pkg/errors"
)







type QueryExecutorCreator interface {
	NewQueryExecutor() (ledger.QueryExecutor, error)
}



type Context struct {
	Seq        int
	Envelope   []byte
	TxID       string
	Channel    string
	PluginName string
	Policy     []byte
	Namespace  string
	Block      *common.Block
}


func (c Context) String() string {
	return fmt.Sprintf("Tx %s, seq %d out of %d in block %d for channel %s with validation plugin %s", c.TxID, c.Seq, len(c.Block.Data.Data), c.Block.Header.Number, c.Channel, c.PluginName)
}


type PluginValidator struct {
	sync.Mutex
	pluginChannelMapping map[vp.Name]*pluginsByChannel
	vp.Mapper
	QueryExecutorCreator
	msp.IdentityDeserializer
	capabilities Capabilities
	policies.ChannelPolicyManagerGetter
	CollectionResources
}






func NewPluginValidator(pm vp.Mapper, qec QueryExecutorCreator, deserializer msp.IdentityDeserializer, capabilities Capabilities, cpmg policies.ChannelPolicyManagerGetter, cor CollectionResources) *PluginValidator {
	return &PluginValidator{
		capabilities:               capabilities,
		pluginChannelMapping:       make(map[vp.Name]*pluginsByChannel),
		Mapper:                     pm,
		QueryExecutorCreator:       qec,
		IdentityDeserializer:       deserializer,
		ChannelPolicyManagerGetter: cpmg,
		CollectionResources:        cor,
	}
}

func (pv *PluginValidator) ValidateWithPlugin(ctx *Context) error {
	plugin, err := pv.getOrCreatePlugin(ctx)
	if err != nil {
		return &validation.ExecutionFailureError{
			Reason: fmt.Sprintf("plugin with name %s couldn't be used: %v", ctx.PluginName, err),
		}
	}
	err = plugin.Validate(ctx.Block, ctx.Namespace, ctx.Seq, 0, vp.SerializedPolicy(ctx.Policy))
	validityStatus := "valid"
	if err != nil {
		validityStatus = fmt.Sprintf("invalid: %v", err)
	}
	logger.Debug("Transaction", ctx.TxID, "appears to be", validityStatus)
	return err
}

func (pv *PluginValidator) getOrCreatePlugin(ctx *Context) (validation.Plugin, error) {
	pluginFactory := pv.FactoryByName(vp.Name(ctx.PluginName))
	if pluginFactory == nil {
		return nil, errors.Errorf("plugin with name %s wasn't found", ctx.PluginName)
	}

	pluginsByChannel := pv.getOrCreatePluginChannelMapping(vp.Name(ctx.PluginName), pluginFactory)
	return pluginsByChannel.createPluginIfAbsent(ctx.Channel)

}

func (pv *PluginValidator) getOrCreatePluginChannelMapping(plugin vp.Name, pf validation.PluginFactory) *pluginsByChannel {
	pv.Lock()
	defer pv.Unlock()
	endorserChannelMapping, exists := pv.pluginChannelMapping[vp.Name(plugin)]
	if !exists {
		endorserChannelMapping = &pluginsByChannel{
			pluginFactory:    pf,
			channels2Plugins: make(map[string]validation.Plugin),
			pv:               pv,
		}
		pv.pluginChannelMapping[vp.Name(plugin)] = endorserChannelMapping
	}
	return endorserChannelMapping
}

type pluginsByChannel struct {
	sync.RWMutex
	pluginFactory    validation.PluginFactory
	channels2Plugins map[string]validation.Plugin
	pv               *PluginValidator
}

func (pbc *pluginsByChannel) createPluginIfAbsent(channel string) (validation.Plugin, error) {
	pbc.RLock()
	plugin, exists := pbc.channels2Plugins[channel]
	pbc.RUnlock()
	if exists {
		return plugin, nil
	}

	pbc.Lock()
	defer pbc.Unlock()
	plugin, exists = pbc.channels2Plugins[channel]
	if exists {
		return plugin, nil
	}

	pluginInstance := pbc.pluginFactory.New()
	plugin, err := pbc.initPlugin(pluginInstance, channel)
	if err != nil {
		return nil, err
	}
	pbc.channels2Plugins[channel] = plugin
	return plugin, nil
}

func (pbc *pluginsByChannel) initPlugin(plugin validation.Plugin, channel string) (validation.Plugin, error) {
	pp, err := policy.New(pbc.pv.IdentityDeserializer, channel, pbc.pv.ChannelPolicyManagerGetter)
	if err != nil {
		return nil, errors.WithMessage(err, "could not obtain a policy evaluator")
	}

	pe := &PolicyEvaluatorWrapper{IdentityDeserializer: pbc.pv.IdentityDeserializer, PolicyEvaluator: pp}
	sf := &StateFetcherImpl{QueryExecutorCreator: pbc.pv}
	if err := plugin.Init(pe, sf, pbc.pv.capabilities, pbc.pv.CollectionResources); err != nil {
		return nil, errors.Wrap(err, "failed initializing plugin")
	}
	return plugin, nil
}

type PolicyEvaluatorWrapper struct {
	msp.IdentityDeserializer
	p.PolicyEvaluator
}


func (id *PolicyEvaluatorWrapper) Evaluate(policyBytes []byte, signatureSet []*protoutil.SignedData) error {
	return id.PolicyEvaluator.Evaluate(policyBytes, signatureSet)
}


func (id *PolicyEvaluatorWrapper) DeserializeIdentity(serializedIdentity []byte) (Identity, error) {
	mspIdentity, err := id.IdentityDeserializer.DeserializeIdentity(serializedIdentity)
	if err != nil {
		return nil, err
	}
	return &identity{Identity: mspIdentity}, nil
}

type identity struct {
	msp.Identity
}

func (i *identity) GetIdentityIdentifier() *IdentityIdentifier {
	identifier := i.Identity.GetIdentifier()
	return &IdentityIdentifier{
		Id:    identifier.Id,
		Mspid: identifier.Mspid,
	}
}

type StateFetcherImpl struct {
	QueryExecutorCreator
}

func (sf *StateFetcherImpl) FetchState() (State, error) {
	qe, err := sf.NewQueryExecutor()
	if err != nil {
		return nil, err
	}
	return &StateImpl{qe}, nil
}

type StateImpl struct {
	ledger.QueryExecutor
}

func (s *StateImpl) GetStateRangeScanIterator(namespace string, startKey string, endKey string) (ResultsIterator, error) {
	it, err := s.QueryExecutor.GetStateRangeScanIterator(namespace, startKey, endKey)
	if err != nil {
		return nil, err
	}
	return &ResultsIteratorImpl{ResultsIterator: it}, nil
}

type ResultsIteratorImpl struct {
	ledger2.ResultsIterator
}

func (it *ResultsIteratorImpl) Next() (QueryResult, error) {
	return it.ResultsIterator.Next()
}
