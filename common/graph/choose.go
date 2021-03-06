/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package graph

import "math/big"

type orderedSet struct {
	elements []interface{}
}

func (s *orderedSet) add(o interface{}) {
	s.elements = append(s.elements, o)
}

type indiceSet struct {
	indices []int
}

type indiceSets []*indiceSet





func CombinationsExceed(n, k, threshold int) bool {
	if n < k {
		return false
	}
	combinations := &big.Int{}
	combinations = combinations.Binomial(int64(n), int64(k))
	t := &big.Int{}
	t.SetInt64(int64(threshold))
	return combinations.Cmp(t) > 0
}

func chooseKoutOfN(n, k int) indiceSets {
	var res indiceSets
	subGroups := &orderedSet{}
	choose(n, k, 0, nil, subGroups)
	for _, el := range subGroups.elements {
		res = append(res, el.(*indiceSet))
	}
	return res
}

func choose(n int, targetAmount int, i int, currentSubGroup []int, subGroups *orderedSet) {
	
	if len(currentSubGroup) == targetAmount {
		subGroups.add(&indiceSet{indices: currentSubGroup})
		return
	}
	
	itemsLeftToPick := n - i
	if targetAmount-len(currentSubGroup) > itemsLeftToPick {
		return
	}
	
	choose(n, targetAmount, i+1, append(currentSubGroup, i), subGroups)
	
	choose(n, targetAmount, i+1, currentSubGroup, subGroups)
}
