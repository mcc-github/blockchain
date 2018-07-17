

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
	"github.com/mcc-github/blockchain/bccsp"
	"github.com/mcc-github/blockchain/bccsp/pkcs11"
	"github.com/mcc-github/blockchain/bccsp/sw"
	"github.com/pkg/errors"
)

const (
	
	PKCS11BasedFactoryName = "PKCS11"
)


type PKCS11Factory struct{}


func (f *PKCS11Factory) Name() string {
	return PKCS11BasedFactoryName
}


func (f *PKCS11Factory) Get(config *FactoryOpts) (bccsp.BCCSP, error) {
	
	if config == nil || config.Pkcs11Opts == nil {
		return nil, errors.New("Invalid config. It must not be nil.")
	}

	p11Opts := config.Pkcs11Opts

	
	var ks bccsp.KeyStore
	if p11Opts.Ephemeral == true {
		ks = sw.NewDummyKeyStore()
	} else if p11Opts.FileKeystore != nil {
		fks, err := sw.NewFileBasedKeyStore(nil, p11Opts.FileKeystore.KeyStorePath, false)
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to initialize software key store")
		}
		ks = fks
	} else {
		
		ks = sw.NewDummyKeyStore()
	}
	return pkcs11.New(*p11Opts, ks)
}
