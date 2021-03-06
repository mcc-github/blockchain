/*
Copyright IBM Corp. 2016 All Rights Reserved.

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
	"sync"

	"github.com/mcc-github/blockchain/bccsp"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/pkg/errors"
)

var (
	
	DefaultBCCSP bccsp.BCCSP

	
	
	bootBCCSP bccsp.BCCSP

	
	bccspMap map[string]bccsp.BCCSP

	
	factoriesInitOnce sync.Once
	bootBCCSPInitOnce sync.Once

	
	factoriesInitError error

	logger = flogging.MustGetLogger("bccsp")
)



type BCCSPFactory interface {

	
	Name() string

	
	Get(opts *FactoryOpts) (bccsp.BCCSP, error)
}


func GetDefault() bccsp.BCCSP {
	if DefaultBCCSP == nil {
		logger.Debug("Before using BCCSP, please call InitFactories(). Falling back to bootBCCSP.")
		bootBCCSPInitOnce.Do(func() {
			var err error
			f := &SWFactory{}
			bootBCCSP, err = f.Get(&FactoryOpts{
				ProviderName: "SW",
				SwOpts: &SwOpts{
					HashFamily: "SHA2",
					SecLevel:   256,

					Ephemeral: true,
				}})
			if err != nil {
				panic("BCCSP Internal error, failed initialization with GetDefaultOpts!")
			}
		})
		return bootBCCSP
	}
	return DefaultBCCSP
}

func initBCCSP(f BCCSPFactory, config *FactoryOpts) error {
	csp, err := f.Get(config)
	if err != nil {
		return errors.Errorf("Could not initialize BCCSP %s [%s]", f.Name(), err)
	}

	logger.Debugf("Initialize BCCSP [%s]", f.Name())
	bccspMap[f.Name()] = csp
	return nil
}
