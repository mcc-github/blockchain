

package mocks

import mock "github.com/stretchr/testify/mock"


type KeyLevelValidationParameterManager struct {
	mock.Mock
}


func (_m *KeyLevelValidationParameterManager) ExtractValidationParameterDependency(blockNum uint64, txNum uint64, rwset []byte) {
	_m.Called(blockNum, txNum, rwset)
}


func (_m *KeyLevelValidationParameterManager) GetValidationParameterForKey(cc string, coll string, key string, blockNum uint64, txNum uint64) ([]byte, error) {
	ret := _m.Called(cc, coll, key, blockNum, txNum)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(string, string, string, uint64, uint64) []byte); ok {
		r0 = rf(cc, coll, key, blockNum, txNum)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, string, uint64, uint64) error); ok {
		r1 = rf(cc, coll, key, blockNum, txNum)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}


func (_m *KeyLevelValidationParameterManager) SetTxValidationResult(cc string, blockNum uint64, txNum uint64, err error) {
	_m.Called(cc, blockNum, txNum, err)
}
