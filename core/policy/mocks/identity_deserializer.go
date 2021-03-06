

package mocks

import (
	blockchain_protos_gomsp "github.com/mcc-github/blockchain-protos-go/msp"
	mock "github.com/stretchr/testify/mock"

	msp "github.com/mcc-github/blockchain/msp"
)


type IdentityDeserializer struct {
	mock.Mock
}


func (_m *IdentityDeserializer) DeserializeIdentity(serializedIdentity []byte) (msp.Identity, error) {
	ret := _m.Called(serializedIdentity)

	var r0 msp.Identity
	if rf, ok := ret.Get(0).(func([]byte) msp.Identity); ok {
		r0 = rf(serializedIdentity)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(msp.Identity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]byte) error); ok {
		r1 = rf(serializedIdentity)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}


func (_m *IdentityDeserializer) IsWellFormed(identity *blockchain_protos_gomsp.SerializedIdentity) error {
	ret := _m.Called(identity)

	var r0 error
	if rf, ok := ret.Get(0).(func(*blockchain_protos_gomsp.SerializedIdentity) error); ok {
		r0 = rf(identity)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
