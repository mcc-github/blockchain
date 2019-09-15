

package mocks

import (
	common "github.com/mcc-github/blockchain/cmd/common"

	kingpin "gopkg.in/alecthomas/kingpin.v2"

	mock "github.com/stretchr/testify/mock"
)


type CommandRegistrar struct {
	mock.Mock
}


func (_m *CommandRegistrar) Command(name string, help string, onCommand common.CLICommand) *kingpin.CmdClause {
	ret := _m.Called(name, help, onCommand)

	var r0 *kingpin.CmdClause
	if rf, ok := ret.Get(0).(func(string, string, common.CLICommand) *kingpin.CmdClause); ok {
		r0 = rf(name, help, onCommand)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*kingpin.CmdClause)
		}
	}

	return r0
}
