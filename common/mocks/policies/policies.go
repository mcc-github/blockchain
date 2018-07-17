/*
Copyright IBM Corp. 2016 All Rights Reserved.

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

package policies

import (
	"github.com/mcc-github/blockchain/common/policies"
	cb "github.com/mcc-github/blockchain/protos/common"
)


type Policy struct {
	
	Err error
}


func (p *Policy) Evaluate(signatureSet []*cb.SignedData) error {
	return p.Err
}


type Manager struct {
	
	
	Policy *Policy

	
	PolicyMap map[string]policies.Policy

	
	SubManagersMap map[string]*Manager
}


func (m *Manager) Manager(path []string) (policies.Manager, bool) {
	if len(path) == 0 {
		return m, true
	}
	manager, ok := m.SubManagersMap[path[len(path)-1]]
	return manager, ok
}


func (m *Manager) GetPolicy(id string) (policies.Policy, bool) {
	if m.PolicyMap != nil {
		policy, ok := m.PolicyMap[id]
		if ok {
			return policy, true
		}
	}
	return m.Policy, m.Policy != nil
}
