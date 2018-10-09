/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package metadata_test

import (
	"fmt"
	"runtime"
	"testing"

	common "github.com/mcc-github/blockchain/common/metadata"
	"github.com/mcc-github/blockchain/orderer/common/metadata"
	"github.com/stretchr/testify/assert"
)

func TestGetVersionInfo(t *testing.T) {
	
	
	
	if common.Version == "" {
		common.Version = "testVersion"
	}

	expected := fmt.Sprintf(
		"%s:\n Version: %s\n Commit SHA: %s\n Go version: %s\n OS/Arch: %s\n",
		metadata.ProgramName, common.Version,
		common.CommitSHA,
		runtime.Version(),
		fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	)
	assert.Equal(t, expected, metadata.GetVersionInfo())
}
