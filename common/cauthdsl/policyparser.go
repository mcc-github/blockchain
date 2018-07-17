/*
Copyright IBM Corp. 2017 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cauthdsl

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/Knetic/govaluate"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/msp"
	"github.com/mcc-github/blockchain/protos/utils"
)


const (
	GateAnd   = "And"
	GateOr    = "Or"
	GateOutOf = "OutOf"
)


const (
	RoleAdmin  = "admin"
	RoleMember = "member"
	RoleClient = "client"
	RolePeer   = "peer"
	
)

var (
	regex = regexp.MustCompile(
		fmt.Sprintf("^([[:alnum:].-]+)([.])(%s|%s|%s|%s)$",
			RoleAdmin, RoleMember, RoleClient, RolePeer),
	)
	regexErr = regexp.MustCompile("^No parameter '([^']+)' found[.]$")
)



func outof(args ...interface{}) (interface{}, error) {
	toret := "outof("
	if len(args) < 2 {
		return nil, fmt.Errorf("Expected at least two arguments to NOutOf. Given %d", len(args))
	}

	arg0 := args[0]
	
	if n, ok := arg0.(float64); ok {
		toret += strconv.Itoa(int(n))
	} else if n, ok := arg0.(int); ok {
		toret += strconv.Itoa(n)
	} else if n, ok := arg0.(string); ok {
		toret += n
	} else {
		return nil, fmt.Errorf("Unexpected type %s", reflect.TypeOf(arg0))
	}

	for _, arg := range args[1:] {
		toret += ", "
		switch t := arg.(type) {
		case string:
			if regex.MatchString(t) {
				toret += "'" + t + "'"
			} else {
				toret += t
			}
		default:
			return nil, fmt.Errorf("Unexpected type %s", reflect.TypeOf(arg))
		}
	}
	return toret + ")", nil
}

func and(args ...interface{}) (interface{}, error) {
	args = append([]interface{}{len(args)}, args...)
	return outof(args...)
}

func or(args ...interface{}) (interface{}, error) {
	args = append([]interface{}{1}, args...)
	return outof(args...)
}

func firstPass(args ...interface{}) (interface{}, error) {
	toret := "outof(ID"
	for _, arg := range args {
		toret += ", "
		switch t := arg.(type) {
		case string:
			if regex.MatchString(t) {
				toret += "'" + t + "'"
			} else {
				toret += t
			}
		case float32:
		case float64:
			toret += strconv.Itoa(int(t))
		default:
			return nil, fmt.Errorf("Unexpected type %s", reflect.TypeOf(arg))
		}
	}

	return toret + ")", nil
}

func secondPass(args ...interface{}) (interface{}, error) {
	
	if len(args) < 3 {
		return nil, fmt.Errorf("At least 3 arguments expected, got %d", len(args))
	}

	
	var ctx *context
	switch v := args[0].(type) {
	case *context:
		ctx = v
	default:
		return nil, fmt.Errorf("Unrecognized type, expected the context, got %s", reflect.TypeOf(args[0]))
	}

	
	var t int
	switch arg := args[1].(type) {
	case float64:
		t = int(arg)
	default:
		return nil, fmt.Errorf("Unrecognized type, expected a number, got %s", reflect.TypeOf(args[1]))
	}

	
	var n int = len(args) - 1

	
	if t > n {
		return nil, fmt.Errorf("Invalid t-out-of-n predicate, t %d, n %d", t, n)
	}

	policies := make([]*common.SignaturePolicy, 0)

	
	for _, principal := range args[2:] {
		switch t := principal.(type) {
		
		case string:
			
			subm := regex.FindAllStringSubmatch(t, -1)
			if subm == nil || len(subm) != 1 || len(subm[0]) != 4 {
				return nil, fmt.Errorf("Error parsing principal %s", t)
			}

			
			var r msp.MSPRole_MSPRoleType
			switch subm[0][3] {
			case RoleMember:
				r = msp.MSPRole_MEMBER
			case RoleAdmin:
				r = msp.MSPRole_ADMIN
			case RoleClient:
				r = msp.MSPRole_CLIENT
			case RolePeer:
				r = msp.MSPRole_PEER
			default:
				return nil, fmt.Errorf("Error parsing role %s", t)
			}

			
			p := &msp.MSPPrincipal{
				PrincipalClassification: msp.MSPPrincipal_ROLE,
				Principal:               utils.MarshalOrPanic(&msp.MSPRole{MspIdentifier: subm[0][1], Role: r})}
			ctx.principals = append(ctx.principals, p)

			
			dapolicy := SignedBy(int32(ctx.IDNum))
			policies = append(policies, dapolicy)

			
			
			ctx.IDNum++

		
		case *common.SignaturePolicy:
			policies = append(policies, t)

		default:
			return nil, fmt.Errorf("Unrecognized type, expected a principal or a policy, got %s", reflect.TypeOf(principal))
		}
	}

	return NOutOf(int32(t), policies), nil
}

type context struct {
	IDNum      int
	principals []*msp.MSPPrincipal
}

func newContext() *context {
	return &context{IDNum: 0, principals: make([]*msp.MSPPrincipal, 0)}
}



















func FromString(policy string) (*common.SignaturePolicyEnvelope, error) {
	
	intermediate, err := govaluate.NewEvaluableExpressionWithFunctions(
		policy, map[string]govaluate.ExpressionFunction{
			GateAnd:                  and,
			strings.ToLower(GateAnd): and,
			strings.ToUpper(GateAnd): and,
			GateOr:                     or,
			strings.ToLower(GateOr):    or,
			strings.ToUpper(GateOr):    or,
			GateOutOf:                  outof,
			strings.ToLower(GateOutOf): outof,
			strings.ToUpper(GateOutOf): outof,
		},
	)
	if err != nil {
		return nil, err
	}

	intermediateRes, err := intermediate.Evaluate(map[string]interface{}{})
	if err != nil {
		
		if regexErr.MatchString(err.Error()) {
			sm := regexErr.FindStringSubmatch(err.Error())
			if len(sm) == 2 {
				return nil, fmt.Errorf("unrecognized token '%s' in policy string", sm[1])
			}
		}

		return nil, err
	}

	
	
	
	
	
	
	exp, err := govaluate.NewEvaluableExpressionWithFunctions(intermediateRes.(string), map[string]govaluate.ExpressionFunction{"outof": firstPass})
	if err != nil {
		return nil, err
	}

	res, err := exp.Evaluate(map[string]interface{}{})
	if err != nil {
		
		if regexErr.MatchString(err.Error()) {
			sm := regexErr.FindStringSubmatch(err.Error())
			if len(sm) == 2 {
				return nil, fmt.Errorf("unrecognized token '%s' in policy string", sm[1])
			}
		}

		return nil, err
	}

	ctx := newContext()
	parameters := make(map[string]interface{}, 1)
	parameters["ID"] = ctx

	exp, err = govaluate.NewEvaluableExpressionWithFunctions(res.(string), map[string]govaluate.ExpressionFunction{"outof": secondPass})
	if err != nil {
		return nil, err
	}

	res, err = exp.Evaluate(parameters)
	if err != nil {
		
		if regexErr.MatchString(err.Error()) {
			sm := regexErr.FindStringSubmatch(err.Error())
			if len(sm) == 2 {
				return nil, fmt.Errorf("unrecognized token '%s' in policy string", sm[1])
			}
		}

		return nil, err
	}

	p := &common.SignaturePolicyEnvelope{
		Identities: ctx.principals,
		Version:    0,
		Rule:       res.(*common.SignaturePolicy),
	}

	return p, nil
}
