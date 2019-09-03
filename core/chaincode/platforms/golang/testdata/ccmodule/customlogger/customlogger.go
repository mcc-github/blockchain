/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/


package customlogger

import "fmt"

func Logf(msg string, args ...interface{}) {
	fmt.Printf(msg, args...)
}
