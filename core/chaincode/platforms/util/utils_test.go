/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package util

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"testing"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/mcc-github/blockchain/common/metadata"
	"github.com/mcc-github/blockchain/core/config/configtest"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestDockerPull(t *testing.T) {
	codepackage, output := io.Pipe()
	go func() {
		tw := tar.NewWriter(output)

		tw.Close()
		output.Close()
	}()

	binpackage := bytes.NewBuffer(nil)

	
	
	
	
	
	
	
	
	
	
	
	image := fmt.Sprintf("mcc-github/blockchain-ccenv:%s-1.1.0", runtime.GOARCH)
	client, err := docker.NewClientFromEnv()
	if err != nil {
		t.Errorf("failed to get docker client: %s", err)
	}

	err = DockerBuild(
		DockerBuildOptions{
			Image:        image,
			Cmd:          "/bin/true",
			InputStream:  codepackage,
			OutputStream: binpackage,
		},
		client,
	)
	if err != nil {
		t.Errorf("Error during build: %s", err)
	}
}

func TestUtil_GetDockerfileFromConfig(t *testing.T) {
	expected := "FROM " + metadata.DockerNamespace + ":" + runtime.GOARCH + "-" + metadata.Version
	path := "dt"
	viper.Set(path, "FROM $(DOCKER_NS):$(ARCH)-$(PROJECT_VERSION)")
	actual := GetDockerfileFromConfig(path)
	assert.Equal(t, expected, actual, `Error parsing Dockerfile Template. Expected "%s", got "%s"`, expected, actual)
}

func TestMain(m *testing.M) {
	viper.SetConfigName("core")
	viper.SetEnvPrefix("CORE")
	configtest.AddDevConfigPath(nil)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("could not read config %s\n", err)
		os.Exit(-1)
	}
	os.Exit(m.Run())
}
