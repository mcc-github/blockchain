
package mock

import mock "github.com/stretchr/testify/mock"


type MetadataProvider struct {
	mock.Mock
}


func (_m *MetadataProvider) GetDBArtifacts(codePackage []byte) ([]byte, error) {
	ret := _m.Called(codePackage)

	var r0 []byte
	if rf, ok := ret.Get(0).(func([]byte) []byte); ok {
		r0 = rf(codePackage)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]byte) error); ok {
		r1 = rf(codePackage)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
