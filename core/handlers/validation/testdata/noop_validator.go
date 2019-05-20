/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	validation "github.com/mcc-github/blockchain/core/handlers/validation/api"
	"github.com/mcc-github/blockchain/protos/common"
)


type NoOpValidator struct {
}


func (*NoOpValidator) Validate(_ *common.Block, _ string, _ int, _ int, _ ...validation.ContextDatum) error {
	return nil
}


func (*NoOpValidator) Init(dependencies ...validation.Dependency) error {
	return nil
}


type NoOpValidatorFactory struct {
}


func (*NoOpValidatorFactory) New() validation.Plugin {
	return &NoOpValidator{}
}



func NewPluginFactory() validation.PluginFactory {
	return &NoOpValidatorFactory{}
}
