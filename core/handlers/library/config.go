

package library

import (
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)



type Config struct {
	AuthFilters []*HandlerConfig `mapstructure:"authFilters" yaml:"authFilters"`
	Decorators  []*HandlerConfig `mapstructure:"decorators" yaml:"decorators"`
	Endorsers   PluginMapping    `mapstructure:"endorsers" yaml:"endorsers"`
	Validators  PluginMapping    `mapstructure:"validators" yaml:"validators"`
}


type PluginMapping map[string]*HandlerConfig


type HandlerConfig struct {
	Name    string `mapstructure:"name" yaml:"name"`
	Library string `mapstructure:"library" yaml:"library"`
}

func LoadConfig() (Config, error) {
	var authFilters, decorators []*HandlerConfig
	if err := mapstructure.Decode(viper.Get("peer.handlers.authFilters"), &authFilters); err != nil {
		return Config{}, err
	}

	if err := mapstructure.Decode(viper.Get("peer.handlers.decorators"), &decorators); err != nil {
		return Config{}, err
	}

	endorsers, validators := make(PluginMapping), make(PluginMapping)
	e := viper.GetStringMap("peer.handlers.endorsers")
	for k := range e {
		name := viper.GetString("peer.handlers.endorsers." + k + ".name")
		library := viper.GetString("peer.handlers.endorsers." + k + ".library")
		endorsers[k] = &HandlerConfig{Name: name, Library: library}
	}

	v := viper.GetStringMap("peer.handlers.validators")
	for k := range v {
		name := viper.GetString("peer.handlers.validators." + k + ".name")
		library := viper.GetString("peer.handlers.validators." + k + ".library")
		validators[k] = &HandlerConfig{Name: name, Library: library}
	}

	return Config{
		AuthFilters: authFilters,
		Decorators:  decorators,
		Endorsers:   endorsers,
		Validators:  validators,
	}, nil
}
