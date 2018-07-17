/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package errors



type TxValidationError interface {
	error
	IsValid() bool
}



type VSCCInfoLookupFailureError struct {
	Reason string
}


func (e VSCCInfoLookupFailureError) Error() string {
	return e.Reason
}



type VSCCEndorsementPolicyError struct {
	Err error
}

func (e *VSCCEndorsementPolicyError) IsValid() bool {
	return e.Err == nil
}


func (e VSCCEndorsementPolicyError) Error() string {
	return e.Err.Error()
}




type VSCCExecutionFailureError struct {
	Err error
}


func (e VSCCExecutionFailureError) Error() string {
	return e.Err.Error()
}

func (e *VSCCExecutionFailureError) IsValid() bool {
	return e.Err == nil
}
