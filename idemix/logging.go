/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package idemix

import "log"


var logger Logger = LogFunc(log.Printf)







func SetLogger(l Logger) {
	logger = l
}




type Logger interface {
	Printf(format string, a ...interface{})
}


type LogFunc func(format string, a ...interface{})


func (l LogFunc) Printf(format string, a ...interface{}) {
	l(format, a...)
}
