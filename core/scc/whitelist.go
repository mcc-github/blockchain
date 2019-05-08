/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package scc

import (
	"strings"

	"github.com/spf13/viper"
)


type Whitelist map[string]bool


func GlobalWhitelist() Whitelist {
	whitelist := Whitelist{}
	whitelist.load()
	return whitelist
}

func (w Whitelist) load() {
	viper.SetEnvPrefix("CORE")
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	for k, v := range viper.GetStringMapString("chaincode.system") {
		w[k] = parseBool(v)
	}
}

func parseBool(s string) bool {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "true", "t", "1", "enable", "enabled", "yes":
		return true
	default:
		return false
	}
}
