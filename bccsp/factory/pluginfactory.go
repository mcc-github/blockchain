/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package factory

import (
	"errors"
	"fmt"
	"os"
	"plugin"

	"github.com/mcc-github/blockchain/bccsp"
)

const (
	
	PluginFactoryName = "PLUGIN"
)


type PluginOpts struct {
	
	Library string
	
	Config map[string]interface{}
}


type PluginFactory struct{}


func (f *PluginFactory) Name() string {
	return PluginFactoryName
}


func (f *PluginFactory) Get(config *FactoryOpts) (bccsp.BCCSP, error) {
	
	if config == nil || config.PluginOpts == nil {
		return nil, errors.New("Invalid config. It must not be nil.")
	}

	
	if config.PluginOpts.Library == "" {
		return nil, errors.New("Invalid config: missing property 'Library'")
	}

	
	if _, err := os.Stat(config.PluginOpts.Library); err != nil {
		return nil, fmt.Errorf("Could not find library '%s' [%s]", config.PluginOpts.Library, err)
	}

	
	plug, err := plugin.Open(config.PluginOpts.Library)
	if err != nil {
		return nil, fmt.Errorf("Failed to load plugin '%s' [%s]", config.PluginOpts.Library, err)
	}

	
	sym, err := plug.Lookup("New")
	if err != nil {
		return nil, fmt.Errorf("Could not find required symbol 'CryptoServiceProvider' [%s]", err)
	}

	
	new, ok := sym.(func(config map[string]interface{}) (bccsp.BCCSP, error))
	if !ok {
		return nil, fmt.Errorf("Plugin does not implement the required function signature for 'New'")
	}

	return new(config.PluginOpts.Config)
}
