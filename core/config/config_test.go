/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package config

import (
	"os"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestConfig_dirExists(t *testing.T) {
	tmpF := os.TempDir()
	exists := dirExists(tmpF)
	assert.True(t, exists,
		"%s directory exists but dirExists returned false", tmpF)

	tmpF = "/blah-" + time.Now().Format(time.RFC3339Nano)
	exists = dirExists(tmpF)
	assert.False(t, exists,
		"%s directory does not exist but dirExists returned true",
		tmpF)
}

func TestConfig_InitViper(t *testing.T) {
	
	v := viper.New()
	err := InitViper(v, "")
	assert.NoError(t, err, "Error returned by InitViper")

	
	err = InitViper(nil, "")
	assert.NoError(t, err, "Error returned by InitViper")
}

func TestConfig_GetPath(t *testing.T) {
	
	path := GetPath("foo")
	assert.Equal(t, "", path, "GetPath should have returned empty string for path 'foo'")

	
	viper.Set("testpath", "/test/config.yml")
	path = GetPath("testpath")
	assert.Equal(t, "/test/config.yml", path)
}

func TestConfig_TranslatePathInPlace(t *testing.T) {
	
	p := "foo"
	TranslatePathInPlace(OfficialPath, &p)
	assert.NotEqual(t, "foo", p, "TranslatePathInPlace failed to translate path %s", p)

	
	p = "/foo"
	TranslatePathInPlace(OfficialPath, &p)
	assert.Equal(t, "/foo", p, "TranslatePathInPlace failed to translate path %s", p)
}
