/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package util

import (
	"bytes"
	"fmt"
	"io"
	"runtime"
	"strings"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/common/metadata"
	"github.com/mcc-github/blockchain/core/config"
	"github.com/spf13/viper"
)

var logger = flogging.MustGetLogger("chaincode.platform.util")

type DockerBuildOptions struct {
	Image        string
	Env          []string
	Cmd          string
	InputStream  io.Reader
	OutputStream io.Writer
}

























func DockerBuild(opts DockerBuildOptions) error {
	var client *docker.Client
	var err error
	endpoint := viper.GetString("vm.endpoint")
	tlsEnabled := viper.GetBool("vm.docker.tls.enabled")
	if tlsEnabled {
		cert := config.GetPath("vm.docker.tls.cert.file")
		key := config.GetPath("vm.docker.tls.key.file")
		ca := config.GetPath("vm.docker.tls.ca.file")
		client, err = docker.NewTLSClient(endpoint, cert, key, ca)
	} else {
		client, err = docker.NewClient(endpoint)
	}

	if err != nil {
		return fmt.Errorf("Error creating docker client: %s", err)
	}
	if opts.Image == "" {
		opts.Image = GetDockerfileFromConfig("chaincode.builder")
		if opts.Image == "" {
			return fmt.Errorf("No image provided and \"chaincode.builder\" default does not exist")
		}
	}

	logger.Debugf("Attempting build with image %s", opts.Image)

	
	
	
	_, err = client.InspectImage(opts.Image)
	if err != nil {
		logger.Debugf("Image %s does not exist locally, attempt pull", opts.Image)

		err = client.PullImage(docker.PullImageOptions{Repository: opts.Image}, docker.AuthConfiguration{})
		if err != nil {
			return fmt.Errorf("Failed to pull %s: %s", opts.Image, err)
		}
	}

	
	
	
	container, err := client.CreateContainer(docker.CreateContainerOptions{
		Config: &docker.Config{
			Image:        opts.Image,
			Env:          opts.Env,
			Cmd:          []string{"/bin/sh", "-c", opts.Cmd},
			AttachStdout: true,
			AttachStderr: true,
		},
	})
	if err != nil {
		return fmt.Errorf("Error creating container: %s", err)
	}
	defer client.RemoveContainer(docker.RemoveContainerOptions{ID: container.ID})

	
	
	
	err = client.UploadToContainer(container.ID, docker.UploadToContainerOptions{
		Path:        "/chaincode/input",
		InputStream: opts.InputStream,
	})
	if err != nil {
		return fmt.Errorf("Error uploading input to container: %s", err)
	}

	
	
	
	stdout := bytes.NewBuffer(nil)
	cw, err := client.AttachToContainerNonBlocking(docker.AttachToContainerOptions{
		Container:    container.ID,
		OutputStream: stdout,
		ErrorStream:  stdout,
		Logs:         true,
		Stdout:       true,
		Stderr:       true,
		Stream:       true,
	})
	if err != nil {
		return fmt.Errorf("Error attaching to container: %s", err)
	}

	
	
	
	err = client.StartContainer(container.ID, nil)
	if err != nil {
		cw.Close()
		return fmt.Errorf("Error executing build: %s \"%s\"", err, stdout.String())
	}

	
	
	
	retval, err := client.WaitContainer(container.ID)
	if err != nil {
		cw.Close()
		return fmt.Errorf("Error waiting for container to complete: %s", err)
	}

	
	cw.Close()
	if err := cw.Wait(); err != nil {
		logger.Errorf("attach wait failed: %s", err)
	}

	if retval > 0 {
		return fmt.Errorf("Error returned from build: %d \"%s\"", retval, stdout.String())
	}

	logger.Debugf("Build output is %s", stdout.String())

	
	
	
	err = client.DownloadFromContainer(container.ID, docker.DownloadFromContainerOptions{
		Path:         "/chaincode/output/.",
		OutputStream: opts.OutputStream,
	})
	if err != nil {
		return fmt.Errorf("Error downloading output: %s", err)
	}

	return nil
}

func GetDockerfileFromConfig(path string) string {
	r := strings.NewReplacer(
		"$(ARCH)", runtime.GOARCH,
		"$(PROJECT_VERSION)", metadata.Version,
		"$(DOCKER_NS)", metadata.DockerNamespace,
		"$(BASE_DOCKER_NS)", metadata.BaseDockerNamespace)

	return r.Replace(viper.GetString(path))
}
