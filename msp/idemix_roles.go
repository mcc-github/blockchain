/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package msp

import (
	m "github.com/mcc-github/blockchain/protos/msp"
)


type Role int32


const (
	MEMBER Role = 1
	ADMIN  Role = 2
	CLIENT Role = 4
	PEER   Role = 8
	
)

func (role Role) getValue() int {
	return int(role)
}


func checkRole(bitmask int, role Role) bool {
	return (bitmask & role.getValue()) == role.getValue()
}


func getRoleMaskFromIdemixRoles(roles []Role) int {
	mask := 0
	for _, role := range roles {
		mask = mask | role.getValue()
	}
	return mask
}


func GetRoleMaskFromIdemixRole(role Role) int {
	return getRoleMaskFromIdemixRoles([]Role{role})
}


func getIdemixRoleFromMSPRole(role *m.MSPRole) int {
	return getIdemixRoleFromMSPRoleType(role.GetRole())
}


func getIdemixRoleFromMSPRoleType(rtype m.MSPRole_MSPRoleType) int {
	return getIdemixRoleFromMSPRoleValue(int(rtype))
}


func getIdemixRoleFromMSPRoleValue(role int) int {
	switch role {
	case int(m.MSPRole_ADMIN):
		return ADMIN.getValue()
	case int(m.MSPRole_CLIENT):
		return CLIENT.getValue()
	case int(m.MSPRole_MEMBER):
		return MEMBER.getValue()
	case int(m.MSPRole_PEER):
		return PEER.getValue()
	default:
		return MEMBER.getValue()
	}
}
