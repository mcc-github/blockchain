


/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package factory

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func buildPlugin(lib string, t *testing.T) {
	t.Helper()
	
	if _, err := os.Stat(lib); err != nil {
		
		cmd := exec.Command("go", "build", "-buildmode=plugin", "github.com/mcc-github/blockchain/examples/plugins/bccsp")
		err := cmd.Run()
		if err != nil {
			t.Fatalf("Could not build plugin: [%s]", err)
		}
	}
}

func TestPluginFactoryName(t *testing.T) {
	f := &PluginFactory{}
	assert.Equal(t, f.Name(), PluginFactoryName)
}

func TestPluginFactoryInvalidConfig(t *testing.T) {
	f := &PluginFactory{}
	opts := &FactoryOpts{}

	_, err := f.Get(nil)
	assert.Error(t, err)

	_, err = f.Get(opts)
	assert.Error(t, err)

	opts.PluginOpts = &PluginOpts{}
	_, err = f.Get(opts)
	assert.Error(t, err)
}

func TestPluginFactoryValidConfig(t *testing.T) {
	
	lib := "./bccsp.so"
	defer os.Remove(lib)
	buildPlugin(lib, t)

	f := &PluginFactory{}
	opts := &FactoryOpts{
		PluginOpts: &PluginOpts{
			Library: lib,
		},
	}

	csp, err := f.Get(opts)
	assert.NoError(t, err)
	assert.NotNil(t, csp)

	_, err = csp.GetKey([]byte{123})
	assert.NoError(t, err)
}

func TestPluginFactoryFromOpts(t *testing.T) {
	
	lib := "./bccsp.so"
	defer os.Remove(lib)
	buildPlugin(lib, t)

	opts := &FactoryOpts{
		ProviderName: "PLUGIN",
		PluginOpts: &PluginOpts{
			Library: lib,
		},
	}
	csp, err := GetBCCSPFromOpts(opts)
	assert.NoError(t, err)
	assert.NotNil(t, csp)

	_, err = csp.GetKey([]byte{123})
	assert.NoError(t, err)
}
