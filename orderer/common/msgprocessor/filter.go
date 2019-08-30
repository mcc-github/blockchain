/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package msgprocessor

import (
	"errors"

	ab "github.com/mcc-github/blockchain-protos-go/common"
)


var ErrEmptyMessage = errors.New("Message was empty")


type Rule interface {
	
	Apply(message *ab.Envelope) error
}


var EmptyRejectRule = Rule(emptyRejectRule{})

type emptyRejectRule struct{}

func (a emptyRejectRule) Apply(message *ab.Envelope) error {
	if message.Payload == nil {
		return ErrEmptyMessage
	}
	return nil
}


var AcceptRule = Rule(acceptRule{})

type acceptRule struct{}

func (a acceptRule) Apply(message *ab.Envelope) error {
	return nil
}


type RuleSet struct {
	rules []Rule
}


func NewRuleSet(rules []Rule) *RuleSet {
	return &RuleSet{
		rules: rules,
	}
}


func (rs *RuleSet) Apply(message *ab.Envelope) error {
	for _, rule := range rs.rules {
		err := rule.Apply(message)
		if err != nil {
			return err
		}
	}
	return nil
}
