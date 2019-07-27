/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package fakes


type HistoryleveldbLogger struct {
	*HistorydbLogger
}




func (hl *HistoryleveldbLogger) Warnf(arg1 string, arg2 ...interface{}) {
	for i := 0; i < len(arg2); i++ {
		if b, ok := arg2[i].([]byte); ok {
			bCopy := make([]byte, len(b))
			copy(bCopy, b)
			arg2[i] = bCopy
		}
	}
	hl.HistorydbLogger.Warnf(arg1, arg2...)
}
