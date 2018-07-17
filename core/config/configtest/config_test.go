/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package configtest_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/mcc-github/blockchain/core/config/configtest"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestConfig_AddDevConfigPath(t *testing.T) {
	
	v := viper.New()
	err := configtest.AddDevConfigPath(v)
	assert.NoError(t, err, "Error while adding dev config path to viper")

	
	err = configtest.AddDevConfigPath(nil)
	assert.NoError(t, err, "Error while adding dev config path to default viper")

	
	gopath := os.Getenv("GOPATH")
	os.Setenv("GOPATH", "")
	defer os.Setenv("GOPATH", gopath)
	err = configtest.AddDevConfigPath(v)
	assert.Error(t, err, "GOPATH is empty, expected error from AddDevConfigPath")
}

func TestConfig_GetDevMspDir(t *testing.T) {
	
	_, err := configtest.GetDevMspDir()
	assert.NoError(t, err)

	
	gopath := os.Getenv("GOPATH")
	os.Setenv("GOPATH", "")
	defer os.Setenv("GOPATH", gopath)
	_, err = configtest.GetDevMspDir()
	assert.Error(t, err, "GOPATH is empty, expected error from GetDevMspDir")

	
	dir, err1 := ioutil.TempDir("/tmp", "devmspdir")
	assert.NoError(t, err1)
	defer os.RemoveAll(dir)
	os.Setenv("GOPATH", dir)
	_, err = configtest.GetDevMspDir()
	assert.Error(t, err, "GOPATH is set to temp dir, expected error from GetDevMspDir")
}

func TestConfig_GetDevConfigDir(t *testing.T) {
	gopath := os.Getenv("GOPATH")
	os.Setenv("GOPATH", "")
	defer os.Setenv("GOPATH", gopath)
	_, err := configtest.GetDevConfigDir()
	assert.Error(t, err, "GOPATH is empty, expected error from GetDevConfigDir")
}
