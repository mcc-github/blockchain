/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package graph



type treePermutations struct {
	originalRoot           *TreeVertex                     
	permutations           []*TreeVertex                   
	descendantPermutations map[*TreeVertex][][]*TreeVertex 
}


func newTreePermutation(root *TreeVertex) *treePermutations {
	return &treePermutations{
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
