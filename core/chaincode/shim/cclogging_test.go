/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package shim

import (
	"os"
	"strings"
	"testing"

	"github.com/mcc-github/blockchain/common/flogging"
	logging "github.com/op/go-logging"
	"github.com/stretchr/testify/assert"
)







func TestShimLogging(t *testing.T) {
	SetLoggingLevel(LogCritical)
	if shimLoggingLevel != LogCritical {
		t.Errorf("shimLoggingLevel is not LogCritical as expected")
	}
	if chaincodeLogger.IsEnabledFor(logging.DEBUG) {
		t.Errorf("The chaincodeLogger should not be enabled for DEBUG")
	}
	if !chaincodeLogger.IsEnabledFor(logging.CRITICAL) {
		t.Errorf("The chaincodeLogger should be enabled for CRITICAL")
	}
	var level LoggingLevel
	var err error
	level, err = LogLevel("debug")
	if err != nil {
		t.Errorf("LogLevel(debug) failed")
	}
	if level != LogDebug {
		t.Errorf("LogLevel(debug) did not return LogDebug")
	}
	level, err = LogLevel("INFO")
	if err != nil {
		t.Errorf("LogLevel(INFO) failed")
	}
	if level != LogInfo {
		t.Errorf("LogLevel(INFO) did not return LogInfo")
	}
	level, err = LogLevel("Notice")
	if err != nil {
		t.Errorf("LogLevel(Notice) failed")
	}
	if level != LogNotice {
		t.Errorf("LogLevel(Notice) did not return LogNotice")
	}
	level, err = LogLevel("WaRnInG")
	if err != nil {
		t.Errorf("LogLevel(WaRnInG) failed")
	}
	if level != LogWarning {
		t.Errorf("LogLevel(WaRnInG) did not return LogWarning")
	}
	level, err = LogLevel("ERRor")
	if err != nil {
		t.Errorf("LogLevel(ERRor) failed")
	}
	if level != LogError {
		t.Errorf("LogLevel(ERRor) did not return LogError")
	}
	level, err = LogLevel("critiCAL")
	if err != nil {
		t.Errorf("LogLevel(critiCAL) failed")
	}
	if level != LogCritical {
		t.Errorf("LogLevel(critiCAL) did not return LogCritical")
	}
	level, err = LogLevel("foo")
	if err == nil {
		t.Errorf("LogLevel(foo) did not fail")
	}
	if level != LogError {
		t.Errorf("LogLevel(foo) did not return LogError")
	}
}


func TestChaincodeLogging(t *testing.T) {

	
	format := logging.MustStringFormatter("%{time:15:04:05.000} [%{module}] %{level:.4s} : %{message}")
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	logging.SetBackend(backendFormatter).SetLevel(logging.Level(shimLoggingLevel), "shim")

	foo := NewLogger("foo")
	bar := NewLogger("bar")

	foo.Debugf("Foo is debugging: %d", 10)
	bar.Infof("Bar is informational? %s.", "Yes")
	foo.Noticef("NOTE NOTE NOTE")
	bar.Warningf("Danger, Danger %s %s", "Will", "Robinson!")
	foo.Errorf("I'm sorry Dave, I'm afraid I can't do that.")
	bar.Criticalf("PI is not equal to 3.14, we computed it as %.2f", 4.13)

	bar.Debug("Foo is debugging:", 10)
	foo.Info("Bar is informational?", "Yes.")
	bar.Notice("NOTE NOTE NOTE")
	foo.Warning("Danger, Danger", "Will", "Robinson!")
	bar.Error("I'm sorry Dave, I'm afraid I can't do that.")
	foo.Critical("PI is not equal to", 3.14, ", we computed it as", 4.13)

	foo.SetLevel(LogWarning)
	if foo.IsEnabledFor(LogDebug) {
		t.Errorf("'foo' should not be enabled for LogDebug")
	}
	if !foo.IsEnabledFor(LogCritical) {
		t.Errorf("'foo' should be enabled for LogCritical")
	}
	bar.SetLevel(LogCritical)
	if bar.IsEnabledFor(LogDebug) {
		t.Errorf("'bar' should not be enabled for LogDebug")
	}
	if !bar.IsEnabledFor(LogCritical) {
		t.Errorf("'bar' should be enabled for LogCritical")
	}
}

func TestSetupChaincodeLogging_shim(t *testing.T) {
	var tests = []struct {
		name         string
		ccLogLevel   string
		shimLogLevel string
	}{
		{name: "ValidLevels", ccLogLevel: "debug", shimLogLevel: "warning"},
		{name: "EmptyLevels", ccLogLevel: "", shimLogLevel: ""},
		{name: "BadShimLevel", ccLogLevel: "debug", shimLogLevel: "war"},
		{name: "BadCCLevel", ccLogLevel: "deb", shimLogLevel: "notice"},
		{name: "EmptyShimLevel", ccLogLevel: "error", shimLogLevel: ""},
		{name: "EmptyCCLevel", ccLogLevel: "", shimLogLevel: "critical"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			os.Setenv("CORE_CHAINCODE_LOGGING_LEVEL", tc.ccLogLevel)
			os.Setenv("CORE_CHAINCODE_LOGGING_SHIM", tc.shimLogLevel)

			setupChaincodeLogging()

			_, ccErr := logging.LogLevel(tc.ccLogLevel)
			_, shimErr := logging.LogLevel(tc.shimLogLevel)
			if ccErr == nil {
				assert.Equal(t, strings.ToUpper(tc.ccLogLevel), logging.GetLevel("ccLogger").String())
				if shimErr == nil {
					assert.Equal(t, strings.ToUpper(tc.shimLogLevel), logging.GetLevel("shim").String())
				} else {
					assert.Equal(t, strings.ToUpper(tc.ccLogLevel), logging.GetLevel("shim").String())
				}
			} else {
				assert.Equal(t, flogging.DefaultLevel(), logging.GetLevel("ccLogger").String())
				if shimErr == nil {
					assert.Equal(t, strings.ToUpper(tc.shimLogLevel), logging.GetLevel("shim").String())
				} else {
					assert.Equal(t, flogging.DefaultLevel(), logging.GetLevel("shim").String())
				}
			}
		})
	}
}



func TestSetupChaincodeLogging_blankLevel(t *testing.T) {
	
	testLogLevelString := ""
	testLogFormat := "%{color}%{time:2006-01-02 15:04:05.000 MST} [%{module}] %{shortfunc} -> %{level:.4s} %{id:03x}%{color:reset} %{message}"

	os.Unsetenv("CORE_CHAINCODE_LOGGING_SHIM")
	os.Setenv("CORE_CHAINCODE_LOGGING_LEVEL", testLogLevelString)
	defer os.Unsetenv("CORE_CHAINCODE_LOGGING_LEVEL")
	os.Setenv("CORE_CHAINCODE_LOGGING_FORMAT", testLogFormat)
	defer os.Unsetenv("CORE_CHAINCODE_LOGGING_FORMAT")

	SetupChaincodeLogging()

	if !IsEnabledForLogLevel(flogging.DefaultLevel()) {
		t.FailNow()
	}
}



func TestSetupChaincodeLogging(t *testing.T) {
	
	testLogLevel := "debug"
	testShimLogLevel := "warning"
	testLogFormat := "%{color}%{time:2006-01-02 15:04:05.000 MST} [%{module}] %{shortfunc} -> %{level:.4s} %{id:03x}%{color:reset} %{message}"

	os.Setenv("CORE_CHAINCODE_LOGGING_LEVEL", testLogLevel)
	defer os.Unsetenv("CORE_CHAINCODE_LOGGING_LEVEL")
	os.Setenv("CORE_CHAINCODE_LOGGING_FORMAT", testLogFormat)
	defer os.Unsetenv("CORE_CHAINCODE_LOGGING_FORMAT")
	os.Setenv("CORE_CHAINCODE_LOGGING_SHIM", testShimLogLevel)
	defer os.Unsetenv("CORE_CHAINCODE_LOGGING_SHIM")

	SetupChaincodeLogging()

	if !IsEnabledForLogLevel(testShimLogLevel) {
		t.FailNow()
	}
}
