/*
Copyright IBM Corp, SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package library

import (
	"fmt"
	"os"
	"plugin"
	"reflect"
	"sync"

	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/core/handlers/auth"
	"github.com/mcc-github/blockchain/core/handlers/decoration"
	endorsement2 "github.com/mcc-github/blockchain/core/handlers/endorsement/api"
	validation "github.com/mcc-github/blockchain/core/handlers/validation/api"
)

var logger = flogging.MustGetLogger("core.handlers")



type Registry interface {
	
	
	Lookup(HandlerType) interface{}
}



type HandlerType int

const (
	
	Auth HandlerType = iota
	
	
	Decoration
	Endorsement
	Validation

	authPluginFactory      = "NewFilter"
	decoratorPluginFactory = "NewDecorator"
	pluginFactory          = "NewPluginFactory"
)

type registry struct {
	filters    []auth.Filter
	decorators []decoration.Decorator
	endorsers  map[string]endorsement2.PluginFactory
	validators map[string]validation.PluginFactory
}

var once sync.Once
var reg registry



func InitRegistry(c Config) Registry {
	once.Do(func() {
		reg = registry{
			endorsers:  make(map[string]endorsement2.PluginFactory),
			validators: make(map[string]validation.PluginFactory),
		}
		reg.loadHandlers(c)
	})
	return &reg
}


func (r *registry) loadHandlers(c Config) {
	for _, config := range c.AuthFilters {
		r.evaluateModeAndLoad(config, Auth)
	}
	for _, config := range c.Decorators {
		r.evaluateModeAndLoad(config, Decoration)
	}

	for chaincodeID, config := range c.Endorsers {
		r.evaluateModeAndLoad(config, Endorsement, chaincodeID)
	}

	for chaincodeID, config := range c.Validators {
		r.evaluateModeAndLoad(config, Validation, chaincodeID)
	}
}


func (r *registry) evaluateModeAndLoad(c *HandlerConfig, handlerType HandlerType, extraArgs ...string) {
	if c.Library != "" {
		r.loadPlugin(c.Library, handlerType, extraArgs...)
	} else {
		r.loadCompiled(c.Name, handlerType, extraArgs...)
	}
}


func (r *registry) loadCompiled(handlerFactory string, handlerType HandlerType, extraArgs ...string) {
	registryMD := reflect.ValueOf(&HandlerLibrary{})

	o := registryMD.MethodByName(handlerFactory)
	if !o.IsValid() {
		logger.Panicf(fmt.Sprintf("Method %s isn't a method of HandlerLibrary", handlerFactory))
	}

	inst := o.Call(nil)[0].Interface()

	if handlerType == Auth {
		r.filters = append(r.filters, inst.(auth.Filter))
	} else if handlerType == Decoration {
		r.decorators = append(r.decorators, inst.(decoration.Decorator))
	} else if handlerType == Endorsement {
		if len(extraArgs) != 1 {
			logger.Panicf("expected 1 argument in extraArgs")
		}
		r.endorsers[extraArgs[0]] = inst.(endorsement2.PluginFactory)
	} else if handlerType == Validation {
		if len(extraArgs) != 1 {
			logger.Panicf("expected 1 argument in extraArgs")
		}
		r.validators[extraArgs[0]] = inst.(validation.PluginFactory)
	}
}


func (r *registry) loadPlugin(pluginPath string, handlerType HandlerType, extraArgs ...string) {
	if _, err := os.Stat(pluginPath); err != nil {
		logger.Panicf(fmt.Sprintf("Could not find plugin at path %s: %s", pluginPath, err))
	}
	p, err := plugin.Open(pluginPath)
	if err != nil {
		logger.Panicf(fmt.Sprintf("Error opening plugin at path %s: %s", pluginPath, err))
	}

	if handlerType == Auth {
		r.initAuthPlugin(p)
	} else if handlerType == Decoration {
		r.initDecoratorPlugin(p)
	} else if handlerType == Endorsement {
		r.initEndorsementPlugin(p, extraArgs...)
	} else if handlerType == Validation {
		r.initValidationPlugin(p, extraArgs...)
	}
}


func (r *registry) initAuthPlugin(p *plugin.Plugin) {
	constructorSymbol, err := p.Lookup(authPluginFactory)
	if err != nil {
		panicWithLookupError(authPluginFactory, err)
	}
	constructor, ok := constructorSymbol.(func() auth.Filter)
	if !ok {
		panicWithDefinitionError(authPluginFactory)
	}

	filter := constructor()
	if filter != nil {
		r.filters = append(r.filters, filter)
	}
}


func (r *registry) initDecoratorPlugin(p *plugin.Plugin) {
	constructorSymbol, err := p.Lookup(decoratorPluginFactory)
	if err != nil {
		panicWithLookupError(decoratorPluginFactory, err)
	}
	constructor, ok := constructorSymbol.(func() decoration.Decorator)
	if !ok {
		panicWithDefinitionError(decoratorPluginFactory)
	}
	decorator := constructor()
	if decorator != nil {
		r.decorators = append(r.decorators, constructor())
	}
}

func (r *registry) initEndorsementPlugin(p *plugin.Plugin, extraArgs ...string) {
	if len(extraArgs) != 1 {
		logger.Panicf("expected 1 argument in extraArgs")
	}
	factorySymbol, err := p.Lookup(pluginFactory)
	if err != nil {
		panicWithLookupError(pluginFactory, err)
	}

	constructor, ok := factorySymbol.(func() endorsement2.PluginFactory)
	if !ok {
		panicWithDefinitionError(pluginFactory)
	}
	factory := constructor()
	if factory == nil {
		logger.Panicf("factory instance returned nil")
	}
	r.endorsers[extraArgs[0]] = factory
}

func (r *registry) initValidationPlugin(p *plugin.Plugin, extraArgs ...string) {
	if len(extraArgs) != 1 {
		logger.Panicf("expected 1 argument in extraArgs")
	}
	factorySymbol, err := p.Lookup(pluginFactory)
	if err != nil {
		panicWithLookupError(pluginFactory, err)
	}

	constructor, ok := factorySymbol.(func() validation.PluginFactory)
	if !ok {
		panicWithDefinitionError(pluginFactory)
	}
	factory := constructor()
	if factory == nil {
		logger.Panicf("factory instance returned nil")
	}
	r.validators[extraArgs[0]] = factory
}


func panicWithLookupError(factory string, err error) {
	logger.Panicf(fmt.Sprintf("Plugin must contain constructor with name %s. Error from lookup: %s",
		factory, err))
}



func panicWithDefinitionError(factory string) {
	logger.Panicf(fmt.Sprintf("Constructor method %s does not match expected definition",
		factory))
}



func (r *registry) Lookup(handlerType HandlerType) interface{} {
	if handlerType == Auth {
		return r.filters
	} else if handlerType == Decoration {
		return r.decorators
	} else if handlerType == Endorsement {
		return r.endorsers
	} else if handlerType == Validation {
		return r.validators
	}

	return nil
}
