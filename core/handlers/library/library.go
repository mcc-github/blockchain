/*
Copyright IBM Corp, SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package library

import (
	"github.com/mcc-github/blockchain/core/handlers/auth"
	"github.com/mcc-github/blockchain/core/handlers/auth/filter"
	"github.com/mcc-github/blockchain/core/handlers/decoration"
	"github.com/mcc-github/blockchain/core/handlers/decoration/decorator"
	endorsement "github.com/mcc-github/blockchain/core/handlers/endorsement/api"
	"github.com/mcc-github/blockchain/core/handlers/endorsement/builtin"
	validation "github.com/mcc-github/blockchain/core/handlers/validation/api"
	. "github.com/mcc-github/blockchain/core/handlers/validation/builtin"
)



type HandlerLibrary struct {
}






func (r *HandlerLibrary) DefaultAuth() auth.Filter {
	return filter.NewFilter()
}



func (r *HandlerLibrary) ExpirationCheck() auth.Filter {
	return filter.NewExpirationCheckFilter()
}




func (r *HandlerLibrary) DefaultDecorator() decoration.Decorator {
	return decorator.NewDecorator()
}

func (r *HandlerLibrary) DefaultEndorsement() endorsement.PluginFactory {
	return &builtin.DefaultEndorsementFactory{}
}

func (r *HandlerLibrary) DefaultValidation() validation.PluginFactory {
	return &DefaultValidationFactory{}
}
