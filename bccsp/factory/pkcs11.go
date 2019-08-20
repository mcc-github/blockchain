

/*
Copyright IBM Corp. 2017 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package factory

import (
	"github.com/mcc-github/blockchain/bccsp"
	"github.com/mcc-github/blockchain/bccsp/pkcs11"
	"github.com/pkg/errors"
)


type FactoryOpts struct {
	ProviderName string             `mapstructure:"default" json:"default" yaml:"Default"`
	SwOpts       *SwOpts            `mapstructure:"SW,omitempty" json:"SW,omitempty" yaml:"SwOpts"`
	PluginOpts   *PluginOpts        `mapstructure:"PLUGIN,omitempty" json:"PLUGIN,omitempty" yaml:"PluginOpts"`
	Pkcs11Opts   *pkcs11.PKCS11Opts `mapstructure:"PKCS11,omitempty" json:"PKCS11,omitempty" yaml:"PKCS11"`
}





func InitFactories(config *FactoryOpts) error {
	factoriesInitOnce.Do(func() {
		setFactories(config)
	})

	return factoriesInitError
}

func setFactories(config *FactoryOpts) error {
	
	if config == nil {
		config = GetDefaultOpts()
	}

	if config.ProviderName == "" {
		config.ProviderName = "SW"
	}

	if config.SwOpts == nil {
		config.SwOpts = GetDefaultOpts().SwOpts
	}

	
	bccspMap = make(map[string]bccsp.BCCSP)

	
	if config.SwOpts != nil {
		f := &SWFactory{}
		err := initBCCSP(f, config)
		if err != nil {
			factoriesInitError = errors.Wrap(err, "Failed initializing SW.BCCSP")
		}
	}

	
	if config.Pkcs11Opts != nil {
		f := &PKCS11Factory{}
		err := initBCCSP(f, config)
		if err != nil {
			factoriesInitError = errors.Wrapf(err, "Failed initializing PKCS11.BCCSP %s", factoriesInitError)
		}
	}

	
	if config.PluginOpts != nil {
		f := &PluginFactory{}
		err := initBCCSP(f, config)
		if err != nil {
			factoriesInitError = errors.Wrapf(err, "Failed initializing PKCS11.BCCSP %s", factoriesInitError)
		}
	}

	var ok bool
	DefaultBCCSP, ok = bccspMap[config.ProviderName]
	if !ok {
		factoriesInitError = errors.Errorf("%s\nCould not find default `%s` BCCSP", factoriesInitError, config.ProviderName)
	}

	return factoriesInitError
}


func GetBCCSPFromOpts(config *FactoryOpts) (bccsp.BCCSP, error) {
	var f BCCSPFactory
	switch config.ProviderName {
	case "SW":
		f = &SWFactory{}
	case "PKCS11":
		f = &PKCS11Factory{}
	case "PLUGIN":
		f = &PluginFactory{}
	default:
		return nil, errors.Errorf("Could not find BCCSP, no '%s' provider", config.ProviderName)
	}

	csp, err := f.Get(config)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not initialize BCCSP %s", f.Name())
	}
	return csp, nil
}
