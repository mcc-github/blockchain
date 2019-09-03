/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package golang

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_gopathDependencyPackageInfo(t *testing.T) {
	t.Run("TestPeer", func(t *testing.T) {
		deps, err := gopathDependencyPackageInfo(runtime.GOOS, runtime.GOARCH, "github.com/mcc-github/blockchain/cmd/peer")
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

	t.Run("TestIncomplete", func(t *testing.T) {
		_, err := gopathDependencyPackageInfo(runtime.GOOS, runtime.GOARCH, "github.com/mcc-github/blockchain/core/chaincode/platforms/golang/testdata/src/chaincodes/BadImport")
		assert.EqualError(t, err, "failed to calculate dependencies: incomplete package: bogus/package")
	})

	t.Run("TestFromGoroot", func(t *testing.T) {
		deps, err := gopathDependencyPackageInfo(runtime.GOOS, runtime.GOARCH, "os")
		assert.NoError(t, err)
		assert.Empty(t, deps)
	})

	t.Run("TestFailure", func(t *testing.T) {
		_, err := gopathDependencyPackageInfo(runtime.GOOS, runtime.GOARCH, "./doesnotexist")
		assert.EqualError(t, err, "listing deps for pacakge ./doesnotexist failed: exit status 1")
	})
}

func TestPackageInfoFiles(t *testing.T) {
	packageInfo := &PackageInfo{
		GoFiles:        []string{"file1.go", "file2.go"},
		CFiles:         []string{"file1.c", "file2.c"},
		CgoFiles:       []string{"file_cgo1.go", "file_cgo2.go"},
		HFiles:         []string{"file1.h", "file2.h"},
		SFiles:         []string{"file1.s", "file2.s"},
		IgnoredGoFiles: []string{"file1_ignored.go", "file2_ignored.go"},
	}
	expected := []string{
		"file1.go", "file2.go",
		"file1.c", "file2.c",
		"file_cgo1.go", "file_cgo2.go",
		"file1.h", "file2.h",
		"file1.s", "file2.s",
		"file1_ignored.go", "file2_ignored.go",
	}
	assert.Equal(t, expected, packageInfo.Files())
}

func Test_listModuleInfo(t *testing.T) {
	cwd, err := os.Getwd()
	require.NoError(t, err, "failed to get working directory")
	defer func() {
		err := os.Chdir(cwd)
		require.NoError(t, err)
	}()

	err = os.Chdir("testdata/ccmodule")
	require.NoError(t, err, "failed to change to module directory")

	moduleDir, err := os.Getwd()
	require.NoError(t, err, "failed to get module working directory")

	mi, err := listModuleInfo()
	assert.NoError(t, err, "failed to get module info")

	expected := &ModuleInfo{
		Path:       "ccmodule",
		ImportPath: "ccmodule",
		Dir:        moduleDir,
		GoMod:      filepath.Join(moduleDir, "go.mod"),
	}
	assert.Equal(t, expected, mi)
}
