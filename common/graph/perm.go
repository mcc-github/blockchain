/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package graph

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}



type treePermutations struct {
	combinationUpperBound  int                             
	originalRoot           *TreeVertex                     
	permutations           []*TreeVertex                   
	descendantPermutations map[*TreeVertex][][]*TreeVertex 
}


func newTreePermutation(root *TreeVertex, combinationUpperBound int) *treePermutations {
	return &treePermutations{
		combinationUpperBound:  combinationUpperBound,
		descendantPermutations: make(map[*TreeVertex][][]*TreeVertex),
		originalRoot:           root,
		permutations:           []*TreeVertex{root},
	}
}



func (tp *treePermutations) permute() []*Tree {
	tp.computeDescendantPermutations()

	it := newBFSIterator(tp.originalRoot)
	for {
		v := it.Next()
		if v == nil {
			break
		}

		if len(v.Descendants) == 0 {
			continue
		}

		
		
		var permutationsWhereVexists []*TreeVertex
		var permutationsWhereVdoesntExist []*TreeVertex
		for _, perm := range tp.permutations {
			if perm.Exists(v.Id) {
				permutationsWhereVexists = append(permutationsWhereVexists, perm)
			} else {
				permutationsWhereVdoesntExist = append(permutationsWhereVdoesntExist, perm)
			}
		}

		
		tp.permutations = permutationsWhereVdoesntExist

		
		
		for _, perm := range permutationsWhereVexists {
			
			
			
			for _, permutation := range tp.descendantPermutations[v] {
				subGraph := &TreeVertex{
					Id:          v.Id,
					Data:        v.Data,
					Descendants: permutation,
				}
				newTree := perm.Clone()
				newTree.replace(v.Id, subGraph)
				
				tp.permutations = append(tp.permutations, newTree)
			}
		}
	}

	res := make([]*Tree, len(tp.permutations))
	for i, perm := range tp.permutations {
		res[i] = perm.ToTree()
	}
	return res
}



func (tp *treePermutations) computeDescendantPermutations() {
	it := newBFSIterator(tp.originalRoot)
	for {
		v := it.Next()
		if v == nil {
			return
		}

		
		if len(v.Descendants) == 0 {
			continue
		}

		
		for CombinationsExceed(len(v.Descendants), v.Threshold, tp.combinationUpperBound) {
			
			victim := rand.Intn(len(v.Descendants))
			v.Descendants = append(v.Descendants[:victim], v.Descendants[victim+1:]...)
		}

		
		for _, el := range chooseKoutOfN(len(v.Descendants), v.Threshold) {
			
			tp.descendantPermutations[v] = append(tp.descendantPermutations[v], v.selectDescendants(el.indices))
		}
	}
}


func (v *TreeVertex) selectDescendants(indices []int) []*TreeVertex {
	r := make([]*TreeVertex, len(indices))
	i := 0
	for _, index := range indices {
		r[i] = v.Descendants[index]
		i++
	}
	return r
}
