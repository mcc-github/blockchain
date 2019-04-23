/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package persistence


type PackageID string


func (p PackageID) String() string {
	return string(p)
}
