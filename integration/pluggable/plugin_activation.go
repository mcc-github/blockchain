/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package pluggable

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	EndorsementPluginEnvVar = "ENDORSEMENT_PLUGIN_ENV_VAR"
	ValidationPluginEnvVar  = "VALIDATION_PLUGIN_ENV_VAR"
)




func EndorsementPluginActivationFolder() string {
	return os.Getenv(EndorsementPluginEnvVar)
}




func SetEndorsementPluginActivationFolder(path string) {
	os.Setenv(EndorsementPluginEnvVar, path)
}




func ValidationPluginActivationFolder() string {
	return os.Getenv(ValidationPluginEnvVar)
}




func SetValidationPluginActivationFolder(path string) {
	os.Setenv(ValidationPluginEnvVar, path)
}

func markPluginActivation(dir string) {
	fileName := filepath.Join(dir, viper.GetString("peer.id"))
	_, err := os.Create(fileName)
	if err != nil {
		panic(fmt.Sprintf("failed to create file %s: %v", fileName, err))
	}
}



func PublishEndorsementPluginActivation() {
	markPluginActivation(EndorsementPluginActivationFolder())
}



func PublishValidationPluginActivation() {
	markPluginActivation(ValidationPluginActivationFolder())
}



func CountEndorsementPluginActivations() int {
	return listDir(EndorsementPluginActivationFolder())
}



func CountValidationPluginActivations() int {
	return listDir(ValidationPluginActivationFolder())
}

func listDir(d string) int {
	dir, err := ioutil.ReadDir(d)
	if err != nil {
		panic(fmt.Sprintf("failed listing directory %s: %v", d, err))
	}
	return len(dir)
}
