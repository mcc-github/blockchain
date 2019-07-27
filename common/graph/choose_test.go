/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package graph

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCombinationsExceed(t *testing.T) {
	
	assert.False(t, CombinationsExceed(20, 5, 15504))
	assert.False(t, CombinationsExceed(20, 5, 15505))
	assert.True(t, CombinationsExceed(20, 5, 15503))

	
	assert.True(t, CombinationsExceed(10000, 500, 9000))

	
	assert.False(t, CombinationsExceed(20, 30, 0))
}
