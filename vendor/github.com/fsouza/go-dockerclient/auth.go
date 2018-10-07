



package docker

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
)


var ErrCannotParseDockercfg = errors.New("Failed to read authentication from dockercfg")



type AuthConfiguration struct {
	Username      string `json:"username,omitempty"`
	Password      string `json:"password,omitempty"`
	Email         string `json:"email,omitempty"`
	ServerAddress string `json:"serveraddress,omitempty"`

	
	
	
	IdentityToken string `json:"identitytoken,omitempty"`

	
	RegistryToken string `json:"registrytoken,omitempty"`
}



type AuthConfigurations struct {
	Configs map[string]AuthConfiguration `json:"configs"`
}



type AuthConfigurations119 map[string]AuthConfiguration



type dockerConfig struct {
	Auth          string `json:"auth"`
	Email         string `json:"email"`
	IdentityToken string `json:"identitytoken"`
	RegistryToken string `json:"registrytoken"`
}



func NewAuthConfigurationsFromFile(path string) (*AuthConfigurations, error) {
	r, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return NewAuthConfigurations(r)
}

func cfgPaths(dockerConfigEnv string, homeEnv string) []string {
	var paths []string
	if dockerConfigEnv != "" {
		paths = append(paths, path.Join(dockerConfigEnv, "config.json"))
	}
	if homeEnv != "" {
		paths = append(paths, path.Join(homeEnv, ".docker", "config.json"))
		paths = append(paths, path.Join(homeEnv, ".dockercfg"))
	}
	return paths
}






func NewAuthConfigurationsFromDockerCfg() (*AuthConfigurations, error) {
	err := fmt.Errorf("No docker configuration found")
	var auths *AuthConfigurations

	pathsToTry := cfgPaths(os.Getenv("DOCKER_CONFIG"), os.Getenv("HOME"))
	for _, path := range pathsToTry {
		auths, err = NewAuthConfigurationsFromFile(path)
		if err == nil {
			return auths, nil
		}
	}
	return auths, err
}



func NewAuthConfigurations(r io.Reader) (*AuthConfigurations, error) {
	var auth *AuthConfigurations
	confs, err := parseDockerConfig(r)
	if err != nil {
		return nil, err
	}
	auth, err = authConfigs(confs)
	if err != nil {
		return nil, err
	}
	return auth, nil
}

func parseDockerConfig(r io.Reader) (map[string]dockerConfig, error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r)
	byteData := buf.Bytes()

	confsWrapper := struct {
		Auths map[string]dockerConfig `json:"auths"`
	}{}
	if err := json.Unmarshal(byteData, &confsWrapper); err == nil {
		if len(confsWrapper.Auths) > 0 {
			return confsWrapper.Auths, nil
		}
	}

	var confs map[string]dockerConfig
	if err := json.Unmarshal(byteData, &confs); err != nil {
		return nil, err
	}
	return confs, nil
}


func authConfigs(confs map[string]dockerConfig) (*AuthConfigurations, error) {
	c := &AuthConfigurations{
		Configs: make(map[string]AuthConfiguration),
	}

	for reg, conf := range confs {
		if conf.Auth == "" {
			continue
		}
		data, err := base64.StdEncoding.DecodeString(conf.Auth)
		if err != nil {
			return nil, err
		}

		userpass := strings.SplitN(string(data), ":", 2)
		if len(userpass) != 2 {
			return nil, ErrCannotParseDockercfg
		}

		authConfig := AuthConfiguration{
			Email:         conf.Email,
			Username:      userpass[0],
			Password:      userpass[1],
			ServerAddress: reg,
		}

		
		if conf.IdentityToken != "" {
			authConfig.Password = ""
			authConfig.IdentityToken = conf.IdentityToken
		}

		
		if conf.RegistryToken != "" {
			authConfig.Password = ""
			authConfig.RegistryToken = conf.RegistryToken
		}
		c.Configs[reg] = authConfig
	}

	return c, nil
}


type AuthStatus struct {
	Status        string `json:"Status,omitempty" yaml:"Status,omitempty" toml:"Status,omitempty"`
	IdentityToken string `json:"IdentityToken,omitempty" yaml:"IdentityToken,omitempty" toml:"IdentityToken,omitempty"`
}






func (c *Client) AuthCheck(conf *AuthConfiguration) (AuthStatus, error) {
	var authStatus AuthStatus
	if conf == nil {
		return authStatus, errors.New("conf is nil")
	}
	resp, err := c.do("POST", "/auth", doOptions{data: conf})
	if err != nil {
		return authStatus, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return authStatus, err
	}
	if len(data) == 0 {
		return authStatus, nil
	}
	if err := json.Unmarshal(data, &authStatus); err != nil {
		return authStatus, err
	}
	return authStatus, nil
}
