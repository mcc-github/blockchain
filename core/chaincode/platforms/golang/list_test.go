/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package golang

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_dependencyPackageInfo(t *testing.T) {
	t.Run("TestPeer", func(t *testing.T) {
		deps, err := dependencyPackageInfo(runtime.GOOS, runtime.GOARCH, "github.com/mcc-github/blockchain/cmd/peer")
		assert.NoError(t, err, "failed to get dependencyPackageInfo")

		var found bool
		for _, pi := range deps {
			if pi.ImportPath == "github.com/mcc-github/blockchain/cmd/peer" {
				found = true
				break
			}
		}
		assert.True(t, found, "expected to find the peer package")
	})

	t.Run("TestFromGoroot", func(t *testing.T) {
		deps, err := dependencyPackageInfo(runtime.GOOS, runtime.GOARCH, "os")
		assert.NoError(t, err)
		assert.Empty(t, deps)
	})

	t.Run("TestFailure", func(t *testing.T) {
		_, err := dependencyPackageInfo(runtime.GOOS, runtime.GOARCH, "./doesnotexist")
		assert.EqualError(t, err, "listing deps for pacakge ./doesnotexist failed: exit status 1")
	})
}
