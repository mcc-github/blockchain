/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package graph

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPermute(t *testing.T) {
	vR := NewTreeVertex("r", nil)
	vR.Threshold = 2

	vD := vR.AddDescendant(NewTreeVertex("D", nil))
	vD.Threshold = 2
	for _, id := range []string{"A", "B", "C"} {
		vD.AddDescendant(NewTreeVertex(id, nil))
	}

	vE := vR.AddDescendant(NewTreeVertex("E", nil))
	vE.Threshold = 2
	for _, id := range []string{"a", "b", "c"} {
		vE.AddDescendant(NewTreeVertex(id, nil))
	}

	vF := vR.AddDescendant(NewTreeVertex("F", nil))
	vF.Threshold = 2
	for _, id := range []string{"1", "2", "3"} {
		vF.AddDescendant(NewTreeVertex(id, nil))
	}

	permutations := vR.ToTree().Permute(1000)
	
	
	
	
	assert.Equal(t, 27, len(permutations))

	listCombination := func(i Iterator) []string {
		var traversal []string
		for {
			v := i.Next()
			if v == nil {
				break
			}
			traversal = append(traversal, v.Id)
		}
		return traversal
	}

	
	expectedScan := []string{"r", "D", "E", "A", "B", "a", "b"}
	assert.Equal(t, expectedScan, listCombination(permutations[0].BFS()))

	
	expectedScan = []string{"r", "E", "F", "b", "c", "2", "3"}
	assert.Equal(t, expectedScan, listCombination(permutations[26].BFS()))
}

func TestPermuteTooManyCombinations(t *testing.T) {
	root := NewTreeVertex("r", nil)
	root.Threshold = 500
	for i := 0; i < 1000; i++ {
		root.AddDescendant(NewTreeVertex(fmt.Sprintf("%d", i), nil))
	}
	permutations := root.ToTree().Permute(501)
	assert.Len(t, permutations, 501)
}
