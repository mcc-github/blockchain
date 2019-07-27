/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package dockercontroller

import (
	"archive/tar"
	"bufio"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/common/util"
	"github.com/mcc-github/blockchain/core/container"
	"github.com/mcc-github/blockchain/core/container/ccintf"
	cutil "github.com/mcc-github/blockchain/core/container/util"
	"github.com/pkg/errors"
)



const ContainerType = "DOCKER"

var (
	dockerLogger = flogging.MustGetLogger("dockercontroller")
	vmRegExp     = regexp.MustCompile("[^a-zA-Z0-9-_.]")
	imageRegExp  = regexp.MustCompile("^[a-z0-9]+(([._-][a-z0-9]+)+)?$")
)




type dockerClient interface {
	
	CreateContainer(opts docker.CreateContainerOptions) (*docker.Container, error)
	
	
	UploadToContainer(id string, opts docker.UploadToContainerOptions) error
	
	StartContainer(id string, cfg *docker.HostConfig) error
	
	
	AttachToContainer(opts docker.AttachToContainerOptions) error
	
	
	BuildImage(opts docker.BuildImageOptions) error
	
	
	StopContainer(id string, timeout uint) error
	
	
	KillContainer(opts docker.KillContainerOptions) error
	
	RemoveContainer(opts docker.RemoveContainerOptions) error
	
	
	PingWithContext(context.Context) error
	
	
	WaitContainer(containerID string) (int, error)
	
	InspectImage(imageName string) (*docker.Image, error)
}

type PlatformBuilder interface {
	GenerateDockerBuild(ccType, path, name, version string, codePackage []byte) (io.Reader, error)
}


type Provider struct {
	PeerID          string
	NetworkID       string
	BuildMetrics    *BuildMetrics
	HostConfig      *docker.HostConfig
	Client          dockerClient
	AttachStdOut    bool
	ChaincodePull   bool
	NetworkMode     string
	PlatformBuilder PlatformBuilder
}


type DockerVM struct {
	PeerID          string
	NetworkID       string
	BuildMetrics    *BuildMetrics
	HostConfig      *docker.HostConfig
	Client          dockerClient
	AttachStdOut    bool
	ChaincodePull   bool
	NetworkMode     string
	PlatformBuilder PlatformBuilder
}


func (p *Provider) NewVM() container.VM {
	return &DockerVM{
		PeerID:          p.PeerID,
		NetworkID:       p.NetworkID,
		Client:          p.Client,
		BuildMetrics:    p.BuildMetrics,
		AttachStdOut:    p.AttachStdOut,
		HostConfig:      p.HostConfig,
		ChaincodePull:   p.ChaincodePull,
		NetworkMode:     p.NetworkMode,
		PlatformBuilder: p.PlatformBuilder,
	}
}



func (p *Provider) HealthCheck(ctx context.Context) error {
	if err := p.Client.PingWithContext(ctx); err != nil {
		return errors.Wrap(err, "failed to ping to Docker daemon")
	}
	return nil
}

func (vm *DockerVM) createContainer(imageID, containerID string, args, env []string) error {
	logger := dockerLogger.With("imageID", imageID, "containerID", containerID)
	logger.Debugw("create container")
	_, err := vm.Client.CreateContainer(docker.CreateContainerOptions{
		Name: containerID,
		Config: &docker.Config{
			Cmd:          args,
			Image:        imageID,
			Env:          env,
			AttachStdout: vm.AttachStdOut,
			AttachStderr: vm.AttachStdOut,
		},
		HostConfig: vm.HostConfig,
	})
	if err != nil {
		return err
	}
	logger.Debugw("created container")
	return nil
}

func (vm *DockerVM) buildImage(ccid ccintf.CCID, reader io.Reader) error {
	id, err := vm.GetVMNameForDocker(ccid)
	if err != nil {
		return err
	}

	outputbuf := bytes.NewBuffer(nil)
	opts := docker.BuildImageOptions{
		Name:         id,
		Pull:         vm.ChaincodePull,
		NetworkMode:  vm.NetworkMode,
		InputStream:  reader,
		OutputStream: outputbuf,
	}

	startTime := time.Now()
	err = vm.Client.BuildImage(opts)

	vm.BuildMetrics.ChaincodeImageBuildDuration.With(
		"chaincode", ccid.String(),
		"success", strconv.FormatBool(err == nil),
	).Observe(time.Since(startTime).Seconds())

	if err != nil {
		dockerLogger.Errorf("Error building image: %s", err)
		dockerLogger.Errorf("Build Output:\n********************\n%s\n********************", outputbuf.String())
		return err
	}

	dockerLogger.Debugf("Created image: %s", id)
	return nil
}


func (vm *DockerVM) Build(ccid ccintf.CCID, ccType, path, name, version string, codePackage []byte) error {
	imageName, err := vm.GetVMNameForDocker(ccid)
	if err != nil {
		return err
	}

	_, err = vm.Client.InspectImage(imageName)
	if err != docker.ErrNoSuchImage {
		return errors.Wrap(err, "docker image inspection failed")
	}

	dockerfileReader, err := vm.PlatformBuilder.GenerateDockerBuild(ccType, path, name, version, codePackage)
	if err != nil {
		return errors.Wrap(err, "platform builder failed")
	}
	err = vm.buildImage(ccid, dockerfileReader)
	if err != nil {
		return errors.Wrap(err, "docker image build failed")
	}

	return nil
}


func (vm *DockerVM) Start(ccid ccintf.CCID, args, env []string, filesToUpload map[string][]byte) error {
	imageName, err := vm.GetVMNameForDocker(ccid)
	if err != nil {
		return err
	}

	containerName := vm.GetVMName(ccid)
	logger := dockerLogger.With("imageName", imageName, "containerName", containerName)

	vm.stopInternal(containerName, 0, false, false)

	err = vm.createContainer(imageName, containerName, args, env)
	if err != nil {
		logger.Errorf("create container failed: %s", err)
		return err
	}

	
	if vm.AttachStdOut {
		containerLogger := flogging.MustGetLogger("peer.chaincode." + containerName)
		streamOutput(dockerLogger, vm.Client, containerName, containerLogger)
	}

	
	
	if len(filesToUpload) != 0 {
		
		
		payload := bytes.NewBuffer(nil)
		gw := gzip.NewWriter(payload)
		tw := tar.NewWriter(gw)

		for path, fileToUpload := range filesToUpload {
			cutil.WriteBytesToPackage(path, fileToUpload, tw)
		}

		
		if err := tw.Close(); err != nil {
			return fmt.Errorf("Error writing files to upload to Docker instance into a temporary tar blob: %s", err)
		}

		gw.Close()

		err := vm.Client.UploadToContainer(containerName, docker.UploadToContainerOptions{
			InputStream:          bytes.NewReader(payload.Bytes()),
			Path:                 "/",
			NoOverwriteDirNonDir: false,
		})
		if err != nil {
			return fmt.Errorf("Error uploading files to the container instance %s: %s", containerName, err)
		}
	}

	
	err = vm.Client.StartContainer(containerName, nil)
	if err != nil {
		dockerLogger.Errorf("start-could not start container: %s", err)
		return err
	}

	dockerLogger.Debugf("Started container %s", containerName)
	return nil
}


func streamOutput(logger *flogging.FabricLogger, client dockerClient, containerName string, containerLogger *flogging.FabricLogger) {
	
	
	attached := make(chan struct{})
	r, w := io.Pipe()

	go func() {
		
		
		
		
		err := client.AttachToContainer(docker.AttachToContainerOptions{
			Container:    containerName,
			OutputStream: w,
			ErrorStream:  w,
			Logs:         true,
			Stdout:       true,
			Stderr:       true,
			Stream:       true,
			Success:      attached,
		})

		
		
		_ = w.CloseWithError(err)
	}()

	go func() {
		defer r.Close() 

		
		select {
		case <-attached: 
			close(attached) 

		case <-time.After(10 * time.Second):
			logger.Errorf("Timeout while attaching to IO channel in container %s", containerName)
			return
		}

		is := bufio.NewReader(r)
		for {
			
			
			line, err := is.ReadString('\n')
			switch err {
			case nil:
				containerLogger.Info(line)
			case io.EOF:
				logger.Infof("Container %s has closed its IO channel", containerName)
				return
			default:
				logger.Errorf("Error reading container output: %s", err)
				return
			}
		}
	}()
}


func (vm *DockerVM) Stop(ccid ccintf.CCID, timeout uint, dontkill bool, dontremove bool) error {
	id := vm.ccidToContainerID(ccid)
	return vm.stopInternal(id, timeout, dontkill, dontremove)
}


func (vm *DockerVM) Wait(ccid ccintf.CCID) (int, error) {
	id := vm.ccidToContainerID(ccid)
	return vm.Client.WaitContainer(id)
}

func (vm *DockerVM) ccidToContainerID(ccid ccintf.CCID) string {
	return strings.Replace(vm.GetVMName(ccid), ":", "_", -1)
}

func (vm *DockerVM) stopInternal(id string, timeout uint, dontkill, dontremove bool) error {
	logger := dockerLogger.With("id", id)

	logger.Debugw("stopping container")
	err := vm.Client.StopContainer(id, timeout)
	dockerLogger.Debugw("stop container result", "error", err)

	if !dontkill {
		logger.Debugw("killing container")
		err = vm.Client.KillContainer(docker.KillContainerOptions{ID: id})
		logger.Debugw("kill container result", "error", err)
	}

	if !dontremove {
		logger.Debugw("removing container")
		err = vm.Client.RemoveContainer(docker.RemoveContainerOptions{ID: id, Force: true})
		logger.Debugw("remove container result", "error", err)
	}

	return err
}




func (vm *DockerVM) GetVMName(ccid ccintf.CCID) string {
	
	
	return vmRegExp.ReplaceAllString(vm.preFormatImageName(ccid), "-")
}






func (vm *DockerVM) GetVMNameForDocker(ccid ccintf.CCID) (string, error) {
	name := vm.preFormatImageName(ccid)
	hash := hex.EncodeToString(util.ComputeSHA256([]byte(name)))
	saniName := vmRegExp.ReplaceAllString(name, "-")
	imageName := strings.ToLower(fmt.Sprintf("%s-%s", saniName, hash))

	
	if !imageRegExp.MatchString(imageName) {
		dockerLogger.Errorf("Error constructing Docker VM Name. '%s' breaks Docker's repository naming rules", name)
		return "", fmt.Errorf("Error constructing Docker VM Name. '%s' breaks Docker's repository naming rules", imageName)
	}

	return imageName, nil
}

func (vm *DockerVM) preFormatImageName(ccid ccintf.CCID) string {
	name := ccid.String()

	if vm.NetworkID != "" && vm.PeerID != "" {
		name = fmt.Sprintf("%s-%s-%s", vm.NetworkID, vm.PeerID, name)
	} else if vm.NetworkID != "" {
		name = fmt.Sprintf("%s-%s", vm.NetworkID, name)
	} else if vm.PeerID != "" {
		name = fmt.Sprintf("%s-%s", vm.PeerID, name)
	}

	return name
}
