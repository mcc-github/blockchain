/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package inquire

import (
	"fmt"

	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/common/graph"
	"github.com/mcc-github/blockchain/common/policies"
	"github.com/mcc-github/blockchain/protos/common"
)

var logger = flogging.MustGetLogger("policies.inquire")

const (
	combinationsUpperBound = 10000
)

type inquireableSignaturePolicy struct {
	sigPol *common.SignaturePolicyEnvelope
}



func NewInquireableSignaturePolicy(sigPol *common.SignaturePolicyEnvelope) policies.InquireablePolicy {
	return &inquireableSignaturePolicy{
		sigPol: sigPol,
	}
}



func (isp *inquireableSignaturePolicy) SatisfiedBy() []policies.PrincipalSet {
	rootId := fmt.Sprintf("%d", 0)
	root := graph.NewTreeVertex(rootId, isp.sigPol.Rule)
	computePolicyTree(root)
	var res []policies.PrincipalSet
	for _, perm := range root.ToTree().Permute(combinationsUpperBound) {
		principalSet := principalsOfTree(perm, isp.sigPol.Identities)
		if len(principalSet) == 0 {
			return nil
		}
		res = append(res, principalSet)
	}
	return res
}

func principalsOfTree(tree *graph.Tree, principals policies.PrincipalSet) policies.PrincipalSet {
	var principalSet policies.PrincipalSet
	i := tree.BFS()
	for {
		v := i.Next()
		if v == nil {
			break
		}
		if !v.IsLeaf() {
			continue
		}
		pol := v.Data.(*common.SignaturePolicy)
		switch principalIndex := pol.Type.(type) {
		case *common.SignaturePolicy_SignedBy:
			if len(principals) <= int(principalIndex.SignedBy) {
				logger.Warning("Failed computing principalsOfTree, index out of bounds")
				return nil
			}
			principal := principals[principalIndex.SignedBy]
			principalSet = append(principalSet, principal)
		default:
			
			logger.Warning("Leaf vertex", v.Id, "is of type", pol.GetType())
			return nil
		}
	}
	return principalSet
}

func computePolicyTree(v *graph.TreeVertex) {
	sigPol := v.Data.(*common.SignaturePolicy)
	if p := sigPol.GetNOutOf(); p != nil {
		v.Threshold = int(p.N)
		for i, rule := range p.Rules {
			id := fmt.Sprintf("%s.%d", v.Id, i)
			u := v.AddDescendant(graph.NewTreeVertex(id, rule))
			computePolicyTree(u)
		}
	}
}
