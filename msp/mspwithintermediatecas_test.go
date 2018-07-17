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

package msp

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMSPWithIntermediateCAs(t *testing.T) {
	
	
	
	
	
	thisMSP := getLocalMSP(t, "testdata/intermediate")

	

	id, err := thisMSP.GetDefaultSigningIdentity()
	assert.NoError(t, err)

	
	err = thisMSP.Validate(id.GetPublicVersion())
	assert.NoError(t, err)

	
	
	err = localMsp.Validate(id.GetPublicVersion())
	assert.Error(t, err)

	
	
	localMSPID, err := localMsp.GetDefaultSigningIdentity()
	assert.NoError(t, err)
	err = thisMSP.Validate(localMSPID.GetPublicVersion())
	assert.Error(t, err)
}

func TestMSPWithExternalIntermediateCAs(t *testing.T) {
	
	
	
	
	
	
	
	

	thisMSP := getLocalMSP(t, "testdata/external")

	

	id, err := thisMSP.GetDefaultSigningIdentity()
	assert.NoError(t, err)

	
	err = thisMSP.Validate(id.GetPublicVersion())
	assert.NoError(t, err)
}

func TestIntermediateCAIdentityValidity(t *testing.T) {
	
	
	
	
	
	thisMSP := getLocalMSP(t, "testdata/intermediate")

	id := thisMSP.(*bccspmsp).intermediateCerts[0]
	assert.Error(t, id.Validate())
}

func TestMSPWithIntermediateCAs2(t *testing.T) {
	
	
	
	
	
	
	
	thisMSP := getLocalMSP(t, filepath.Join("testdata", "intermediate2"))

	
	
	id, err := thisMSP.GetDefaultSigningIdentity()
	assert.NoError(t, err)
	err = thisMSP.Validate(id.GetPublicVersion())
	assert.NoError(t, err)

	
	pem, err := readPemFile(filepath.Join("testdata", "intermediate2", "users", "user2-cert.pem"))
	assert.NoError(t, err)
	id2, _, err := thisMSP.(*bccspmsp).getIdentityFromConf(pem)
	assert.NoError(t, err)
	err = thisMSP.Validate(id2)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid validation chain. Parent certificate should be a leaf of the certification tree ")
}
