/*
Copyright IBM Corp. 2018 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/


package mocks

import mock "github.com/stretchr/testify/mock"


type SigningIdentity struct {
	mock.Mock
}


func (_m *SigningIdentity) Serialize() ([]byte, error) {
	ret := _m.Called()

	var r0 []byte
	if rf, ok := ret.Get(0).(func() []byte); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}


func (_m *SigningIdentity) Sign(_a0 []byte) ([]byte, error) {
	ret := _m.Called(_a0)

	var r0 []byte
	if rf, ok := ret.Get(0).(func([]byte) []byte); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]byte) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
