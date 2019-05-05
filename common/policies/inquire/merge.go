/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package inquire

import (
	"reflect"

	"github.com/mcc-github/blockchain/common/policies"
)


type ComparablePrincipalSets []ComparablePrincipalSet


func (cps ComparablePrincipalSets) ToPrincipalSets() policies.PrincipalSets {
	var res policies.PrincipalSets
	for _, cp := range cps {
		res = append(res, cp.ToPrincipalSet())
	}
	return res
}









func Merge(s1, s2 ComparablePrincipalSets) ComparablePrincipalSets {
	var res ComparablePrincipalSets
	setsIn1ToTheContainingSetsIn2 := computeContainedInMapping(s1, s2)
	setsIn1ThatAreIn2 := s1.OfMapping(setsIn1ToTheContainingSetsIn2, s2)
	
	
	s1 = s1.ExcludeIndices(setsIn1ToTheContainingSetsIn2)
	setsIn2ToTheContainingSetsIn1 := computeContainedInMapping(s2, s1)
	setsIn2ThatAreIn1 := s2.OfMapping(setsIn2ToTheContainingSetsIn1, s1)
	s2 = s2.ExcludeIndices(setsIn2ToTheContainingSetsIn1)

	
	
	res = append(res, setsIn1ThatAreIn2.ToMergedPrincipalSets()...)
	res = append(res, setsIn2ThatAreIn1.ToMergedPrincipalSets()...)

	
	
	
	
	s1 = s1.ExcludeIndices(setsIn2ToTheContainingSetsIn1.invert())
	s2 = s2.ExcludeIndices(setsIn1ToTheContainingSetsIn2.invert())

	
	
	if len(s1) == 0 || len(s2) == 0 {
		return res.Reduce()
	}

	
	
	
	combinedPairs := CartesianProduct(s1, s2)
	res = append(res, combinedPairs.ToMergedPrincipalSets()...)
	return res.Reduce()
}




func CartesianProduct(s1, s2 ComparablePrincipalSets) comparablePrincipalSetPairs {
	var res comparablePrincipalSetPairs
	for _, x := range s1 {
		var set comparablePrincipalSetPairs
		
		
		for _, y := range s2 {
			set = append(set, comparablePrincipalSetPair{
				contained:  x,
				containing: y,
			})
		}
		res = append(res, set...)
	}
	return res
}


type comparablePrincipalSetPair struct {
	contained  ComparablePrincipalSet
	containing ComparablePrincipalSet
}



func (pair comparablePrincipalSetPair) MergeWithPlurality() ComparablePrincipalSet {
	var principalsToAdd []*ComparablePrincipal
	used := make(map[int]struct{})
	
	for _, principal := range pair.contained {
		var covered bool
		
		for i, coveringPrincipal := range pair.containing {
			
			if _, isUsed := used[i]; isUsed {
				continue
			}
			
			if coveringPrincipal.IsA(principal) {
				used[i] = struct{}{}
				covered = true
				break
			}
		}
		
		
		if !covered {
			principalsToAdd = append(principalsToAdd, principal)
		}
	}

	res := pair.containing.Clone()
	res = append(res, principalsToAdd...)
	return res
}


type comparablePrincipalSetPairs []comparablePrincipalSetPair



func (pairs comparablePrincipalSetPairs) ToMergedPrincipalSets() ComparablePrincipalSets {
	var res ComparablePrincipalSets
	for _, pair := range pairs {
		res = append(res, pair.MergeWithPlurality())
	}
	return res
}


func (cps ComparablePrincipalSets) OfMapping(mapping map[int][]int, sets2 ComparablePrincipalSets) comparablePrincipalSetPairs {
	var res []comparablePrincipalSetPair
	for i, js := range mapping {
		for _, j := range js {
			res = append(res, comparablePrincipalSetPair{
				contained:  cps[i],
				containing: sets2[j],
			})
		}
	}
	return res
}



func (cps ComparablePrincipalSets) Reduce() ComparablePrincipalSets {
	
	
	current := cps
	for {
		currLen := len(current)
		
		reduced := current.reduce()
		newLen := len(reduced)
		if currLen == newLen {
			
			
			return reduced
		}
		
		current = reduced
	}
}

func (cps ComparablePrincipalSets) reduce() ComparablePrincipalSets {
	var res ComparablePrincipalSets
	for i, s1 := range cps {
		var isContaining bool
		for j, s2 := range cps {
			if i == j {
				continue
			}
			if s2.IsSubset(s1) {
				isContaining = true
			}
			
			
			if s1.IsSubset(s2) && i < j {
				isContaining = false
			}

		}
		if !isContaining {
			res = append(res, s1)
		}
	}
	return res
}


func (cps ComparablePrincipalSets) ExcludeIndices(mapping map[int][]int) ComparablePrincipalSets {
	var res ComparablePrincipalSets
	for i, set := range cps {
		if _, exists := mapping[i]; exists {
			continue
		}
		res = append(res, set)
	}
	return res
}





func (cps ComparablePrincipalSet) Contains(s *ComparablePrincipal) bool {
	for _, cp := range cps {
		if cp.IsA(s) {
			return true
		}
	}
	return false
}







func (cps ComparablePrincipalSet) IsContainedIn(set ComparablePrincipalSet) bool {
	for _, cp := range cps {
		if !set.Contains(cp) {
			return false
		}
	}
	return true
}




func computeContainedInMapping(s1, s2 []ComparablePrincipalSet) intMapping {
	mapping := make(map[int][]int)
	for i, ps1 := range s1 {
		for j, ps2 := range s2 {
			if !ps1.IsContainedIn(ps2) {
				continue
			}
			mapping[i] = append(mapping[i], j)
		}
	}
	return mapping
}


type intMapping map[int][]int

func (im intMapping) invert() intMapping {
	res := make(intMapping)
	for i, js := range im {
		for _, j := range js {
			res[j] = append(res[j], i)
		}
	}
	return res
}


func (cps ComparablePrincipalSet) IsSubset(sets ComparablePrincipalSet) bool {
	for _, p1 := range cps {
		var found bool
		for _, p2 := range sets {
			if reflect.DeepEqual(p1, p2) {
				found = true
			}
		}
		if !found {
			return false
		}
	}
	return true
}
