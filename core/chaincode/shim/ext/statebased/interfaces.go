/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package statebased

import "fmt"


type RoleType string

const (
	
	RoleTypeMember = RoleType("MEMBER")
	
	RoleTypePeer = RoleType("PEER")
	
	RoleTypeClient = RoleType("CLIENT")
)




type RoleTypeDoesNotExistError struct {
	roleType RoleType
}

func (r *RoleTypeDoesNotExistError) Error() string {
	return fmt.Sprintf("Role type %s does not exist", r.roleType)
}





type KeyEndorsementPolicy interface {
	
	Policy() ([]byte, error)

	
	
	
	
	
	
	AddOrgs(roleType RoleType, organizations ...string) error

	
	
	DelOrgs(organizations ...string)

	
	ListOrgs() []string
}
