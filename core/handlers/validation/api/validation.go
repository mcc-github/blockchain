/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package validation

import "github.com/mcc-github/blockchain-protos-go/common"


type Argument interface {
	Dependency
	
	Arg() []byte
}


type Dependency interface{}



type ContextDatum interface{}


type Plugin interface {
	
	
	Validate(block *common.Block, namespace string, txPosition int, actionPosition int, contextData ...ContextDatum) error

	
	Init(dependencies ...Dependency) error
}


type PluginFactory interface {
	New() Plugin
}




type ExecutionFailureError struct {
	Reason string
}



func (e *ExecutionFailureError) Error() string {
	return e.Reason
}
