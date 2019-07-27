/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package configtx

import (
	"fmt"
	"strings"
	"testing"

	mockpolicies "github.com/mcc-github/blockchain/common/configtx/mock"
	"github.com/mcc-github/blockchain/common/policies"
	cb "github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/stretchr/testify/assert"
)


type policyManager interface {
	policies.Manager
}


type policy interface {
	policies.Policy
}

var defaultChain = "default.chain.id"

func defaultPolicyManager() *mockpolicies.PolicyManager {
	fakePolicy := &mockpolicies.Policy{}
	fakePolicy.EvaluateReturns(nil)
	fakePolicyManager := &mockpolicies.PolicyManager{}
	fakePolicyManager.GetPolicyReturns(fakePolicy, true)
	fakePolicyManager.ManagerReturns(fakePolicyManager, true)
	return fakePolicyManager
}

type configPair struct {
	key   string
	value *cb.ConfigValue
}

func makeConfigPair(id, modificationPolicy string, lastModified uint64, data []byte) *configPair {
	return &configPair{
		key: id,
		value: &cb.ConfigValue{
			ModPolicy: modificationPolicy,
			Version:   lastModified,
			Value:     data,
		},
	}
}

func makeConfig(configPairs ...*configPair) *cb.Config {
	channelGroup := protoutil.NewConfigGroup()
	for _, pair := range configPairs {
		channelGroup.Values[pair.key] = pair.value
	}

	return &cb.Config{
		ChannelGroup: channelGroup,
	}
}

func makeConfigSet(configPairs ...*configPair) *cb.ConfigGroup {
	result := protoutil.NewConfigGroup()
	for _, pair := range configPairs {
		result.Values[pair.key] = pair.value
	}
	return result
}

func makeConfigUpdateEnvelope(chainID string, readSet, writeSet *cb.ConfigGroup) *cb.Envelope {
	return &cb.Envelope{
		Payload: protoutil.MarshalOrPanic(&cb.Payload{
			Header: &cb.Header{
				ChannelHeader: protoutil.MarshalOrPanic(&cb.ChannelHeader{
					Type: int32(cb.HeaderType_CONFIG_UPDATE),
				}),
			},
			Data: protoutil.MarshalOrPanic(&cb.ConfigUpdateEnvelope{
				ConfigUpdate: protoutil.MarshalOrPanic(&cb.ConfigUpdate{
					ChannelId: chainID,
					ReadSet:   readSet,
					WriteSet:  writeSet,
				}),
			}),
		}),
	}
}

func TestEmptyChannel(t *testing.T) {
	_, err := NewValidatorImpl("foo", &cb.Config{}, "foonamespace", defaultPolicyManager())
	assert.Error(t, err)
}


func TestDifferentChainID(t *testing.T) {
	vi, err := NewValidatorImpl(
		defaultChain,
		makeConfig(makeConfigPair("foo", "foo", 0, []byte("foo"))),
		"foonamespace",
		defaultPolicyManager())

	if err != nil {
		t.Fatalf("Error constructing config manager: %s", err)
	}

	newConfig := makeConfigUpdateEnvelope("wrongChain", makeConfigSet(), makeConfigSet(makeConfigPair("foo", "foo", 1, []byte("foo"))))

	_, err = vi.ProposeConfigUpdate(newConfig)
	if err == nil {
		t.Error("Should have errored when proposing a new config set the wrong chain ID")
	}
}


func TestOldConfigReplay(t *testing.T) {
	vi, err := NewValidatorImpl(
		defaultChain,
		makeConfig(makeConfigPair("foo", "foo", 0, []byte("foo"))),
		"foonamespace",
		defaultPolicyManager())

	if err != nil {
		t.Fatalf("Error constructing config manager: %s", err)
	}

	newConfig := makeConfigUpdateEnvelope(defaultChain, makeConfigSet(), makeConfigSet(makeConfigPair("foo", "foo", 0, []byte("foo"))))

	_, err = vi.ProposeConfigUpdate(newConfig)

	assert.EqualError(t, err, "error authorizing update: error validating DeltaSet: attempt to set key [Value]  /foonamespace/foo to version 0, but key is at version 0")
}


func TestValidConfigChange(t *testing.T) {
	vi, err := NewValidatorImpl(
		defaultChain,
		makeConfig(makeConfigPair("foo", "foo", 0, []byte("foo"))),
		"foonamespace",
		defaultPolicyManager())

	if err != nil {
		t.Fatalf("Error constructing config manager: %s", err)
	}

	newConfig := makeConfigUpdateEnvelope(defaultChain, makeConfigSet(), makeConfigSet(makeConfigPair("foo", "foo", 1, []byte("foo"))))

	configEnv, err := vi.ProposeConfigUpdate(newConfig)
	if err != nil {
		t.Errorf("Should not have errored proposing config: %s", err)
	}

	err = vi.Validate(configEnv)
	if err != nil {
		t.Errorf("Should not have errored validating config: %s", err)
	}
}



func TestConfigChangeRegressedSequence(t *testing.T) {
	vi, err := NewValidatorImpl(
		defaultChain,
		makeConfig(makeConfigPair("foo", "foo", 1, []byte("foo"))),
		"foonamespace",
		defaultPolicyManager())

	if err != nil {
		t.Fatalf("Error constructing config manager: %s", err)
	}

	newConfig := makeConfigUpdateEnvelope(
		defaultChain,
		makeConfigSet(makeConfigPair("foo", "foo", 0, []byte("foo"))),
		makeConfigSet(makeConfigPair("bar", "bar", 2, []byte("bar"))),
	)

	_, err = vi.ProposeConfigUpdate(newConfig)
	assert.EqualError(t, err, "error authorizing update: error validating ReadSet: proposed update requires that key [Value]  /foonamespace/foo be at version 0, but it is currently at version 1")
}



func TestConfigChangeOldSequence(t *testing.T) {
	vi, err := NewValidatorImpl(
		defaultChain,
		makeConfig(makeConfigPair("foo", "foo", 1, []byte("foo"))),
		"foonamespace",
		defaultPolicyManager())

	if err != nil {
		t.Fatalf("Error constructing config manager: %s", err)
	}

	newConfig := makeConfigUpdateEnvelope(
		defaultChain,
		makeConfigSet(),
		makeConfigSet(
			makeConfigPair("foo", "foo", 2, []byte("foo")),
			makeConfigPair("bar", "bar", 1, []byte("bar")),
		),
	)

	_, err = vi.ProposeConfigUpdate(newConfig)

	assert.EqualError(t, err, "error authorizing update: error validating DeltaSet: attempted to set key [Value]  /foonamespace/bar to version 1, but key does not exist")
}



func TestConfigPartialUpdate(t *testing.T) {
	vi, err := NewValidatorImpl(
		defaultChain,
		makeConfig(
			makeConfigPair("foo", "foo", 0, []byte("foo")),
			makeConfigPair("bar", "bar", 0, []byte("bar")),
		),
		"foonamespace",
		defaultPolicyManager())

	if err != nil {
		t.Fatalf("Error constructing config manager: %s", err)
	}

	newConfig := makeConfigUpdateEnvelope(
		defaultChain,
		makeConfigSet(),
		makeConfigSet(makeConfigPair("bar", "bar", 1, []byte("bar"))),
	)

	_, err = vi.ProposeConfigUpdate(newConfig)
	assert.NoError(t, err, "Should have allowed partial update")
}


func TestEmptyConfigUpdate(t *testing.T) {
	vi, err := NewValidatorImpl(
		defaultChain,
		makeConfig(makeConfigPair("foo", "foo", 0, []byte("foo"))),
		"foonamespace",
		defaultPolicyManager())

	if err != nil {
		t.Fatalf("Error constructing config manager: %s", err)
	}

	newConfig := &cb.Envelope{}

	_, err = vi.ProposeConfigUpdate(newConfig)
	assert.EqualError(t, err, "error converting envelope to config update: envelope must have a Header")
}




func TestSilentConfigModification(t *testing.T) {
	vi, err := NewValidatorImpl(
		defaultChain,
		makeConfig(
			makeConfigPair("foo", "foo", 0, []byte("foo")),
			makeConfigPair("bar", "bar", 0, []byte("bar")),
		),
		"foonamespace",
		defaultPolicyManager())

	if err != nil {
		t.Fatalf("Error constructing config manager: %s", err)
	}

	newConfig := makeConfigUpdateEnvelope(
		defaultChain,
		makeConfigSet(),
		makeConfigSet(
			makeConfigPair("foo", "foo", 0, []byte("different")),
			makeConfigPair("bar", "bar", 1, []byte("bar")),
		),
	)

	_, err = vi.ProposeConfigUpdate(newConfig)
	assert.EqualError(t, err, "error authorizing update: error validating DeltaSet: attempt to set key [Value]  /foonamespace/foo to version 0, but key is at version 0")
}



func TestConfigChangeViolatesPolicy(t *testing.T) {
	pm := defaultPolicyManager()
	vi, err := NewValidatorImpl(
		defaultChain,
		makeConfig(makeConfigPair("foo", "foo", 0, []byte("foo"))),
		"foonamespace",
		pm)

	if err != nil {
		t.Fatalf("Error constructing config manager: %s", err)
	}
	
	fakePolicy := &mockpolicies.Policy{}
	fakePolicy.EvaluateReturns(fmt.Errorf("err"))
	pm.GetPolicyReturns(fakePolicy, true)

	newConfig := makeConfigUpdateEnvelope(defaultChain, makeConfigSet(), makeConfigSet(makeConfigPair("foo", "foo", 1, []byte("foo"))))

	_, err = vi.ProposeConfigUpdate(newConfig)
	assert.EqualError(t, err, "error authorizing update: error validating DeltaSet: policy for [Value]  /foonamespace/foo not satisfied: err")
}



func TestUnchangedConfigViolatesPolicy(t *testing.T) {
	pm := defaultPolicyManager()
	vi, err := NewValidatorImpl(
		defaultChain,
		makeConfig(makeConfigPair("foo", "foo", 0, []byte("foo"))),
		"foonamespace",
		pm)

	if err != nil {
		t.Fatalf("Error constructing config manager: %s", err)
	}

	newConfig := makeConfigUpdateEnvelope(
		defaultChain,
		makeConfigSet(makeConfigPair("foo", "foo", 0, []byte("foo"))),
		makeConfigSet(makeConfigPair("bar", "bar", 0, []byte("foo"))),
	)

	configEnv, err := vi.ProposeConfigUpdate(newConfig)
	if err != nil {
		t.Errorf("Should not have errored proposing config, but got %s", err)
	}

	err = vi.Validate(configEnv)
	if err != nil {
		t.Errorf("Should not have errored validating config, but got %s", err)
	}
}



func TestInvalidProposal(t *testing.T) {
	pm := defaultPolicyManager()
	vi, err := NewValidatorImpl(
		defaultChain,
		makeConfig(makeConfigPair("foo", "foo", 0, []byte("foo"))),
		"foonamespace",
		pm)

	if err != nil {
		t.Fatalf("Error constructing config manager: %s", err)
	}

	fakePolicy := &mockpolicies.Policy{}
	fakePolicy.EvaluateReturns(fmt.Errorf("err"))
	pm.GetPolicyReturns(fakePolicy, true)

	newConfig := makeConfigUpdateEnvelope(defaultChain, makeConfigSet(), makeConfigSet(makeConfigPair("foo", "foo", 1, []byte("foo"))))

	_, err = vi.ProposeConfigUpdate(newConfig)
	assert.EqualError(t, err, "error authorizing update: error validating DeltaSet: policy for [Value]  /foonamespace/foo not satisfied: err")
}

func TestValidateErrors(t *testing.T) {
	t.Run("TestNilConfigEnv", func(t *testing.T) {
		err := (&ValidatorImpl{}).Validate(nil)
		assert.Error(t, err)
		assert.Regexp(t, "config envelope is nil", err.Error())
	})

	t.Run("TestNilConfig", func(t *testing.T) {
		err := (&ValidatorImpl{}).Validate(&cb.ConfigEnvelope{})
		assert.Error(t, err)
		assert.Regexp(t, "config envelope has nil config", err.Error())
	})

	t.Run("TestSequenceSkip", func(t *testing.T) {
		err := (&ValidatorImpl{}).Validate(&cb.ConfigEnvelope{
			Config: &cb.Config{
				Sequence: 2,
			},
		})
		assert.Error(t, err)
		assert.Regexp(t, "config currently at sequence 0", err.Error())
	})
}

func TestConstructionErrors(t *testing.T) {
	t.Run("NilConfig", func(t *testing.T) {
		v, err := NewValidatorImpl("test", nil, "foonamespace", &mockpolicies.PolicyManager{})
		assert.Nil(t, v)
		assert.Error(t, err)
		assert.Regexp(t, "nil config parameter", err.Error())
	})

	t.Run("NilChannelGroup", func(t *testing.T) {
		v, err := NewValidatorImpl("test", &cb.Config{}, "foonamespace", &mockpolicies.PolicyManager{})
		assert.Nil(t, v)
		assert.Error(t, err)
		assert.Regexp(t, "nil channel group", err.Error())
	})

	t.Run("BadChannelID", func(t *testing.T) {
		v, err := NewValidatorImpl("*&$#@*&@$#*&", &cb.Config{ChannelGroup: &cb.ConfigGroup{}}, "foonamespace", &mockpolicies.PolicyManager{})
		assert.Nil(t, v)
		assert.Error(t, err)
		assert.Regexp(t, "bad channel ID", err.Error())
		assert.EqualError(t, err, "bad channel ID: '*&$#@*&@$#*&' contains illegal characters")
	})

	t.Run("EmptyChannelID", func(t *testing.T) {
		v, err := NewValidatorImpl("", &cb.Config{ChannelGroup: &cb.ConfigGroup{}}, "foonamespace", &mockpolicies.PolicyManager{})
		assert.Nil(t, v)
		assert.Error(t, err)
		assert.Regexp(t, "bad channel ID", err.Error())
		assert.EqualError(t, err, "bad channel ID: channel ID illegal, cannot be empty")
	})

	t.Run("MaxLengthChannelID", func(t *testing.T) {
		maxChannelID := strings.Repeat("a", 250)
		v, err := NewValidatorImpl(maxChannelID, &cb.Config{ChannelGroup: &cb.ConfigGroup{}}, "foonamespace", &mockpolicies.PolicyManager{})
		assert.Nil(t, v)
		assert.Error(t, err)
		assert.EqualError(t, err, "bad channel ID: channel ID illegal, cannot be longer than 249")
	})
}
