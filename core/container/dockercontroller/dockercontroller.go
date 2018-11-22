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

	"github.com/fsouza/go-dockerclient"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/common/metrics"
	"github.com/mcc-github/blockchain/common/util"
	"github.com/mcc-github/blockchain/core/container"
	"github.com/mcc-github/blockchain/core/container/ccintf"
	cutil "github.com/mcc-github/blockchain/core/container/util"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)



const ContainerType = "DOCKER"

var (
	dockerLogger = flogging.MustGetLogger("dockercontroller")
	hostConfig   *docker.HostConfig
	vmRegExp     = regexp.MustCompile("[^a-zA-Z0-9-_.]")
	imageRegExp  = regexp.MustCompile("^[a-z0-9]+(([._-][a-z0-9]+)+)?$")
)


type getClient func() (dockerClient, error)


type DockerVM struct {
	getClientFnc getClient
	PeerID       string
	NetworkID    string
	BuildMetrics *BuildMetrics
}


type dockerClient interface {
	
	CreateContainer(opts docker.CreateContainerOptions) (*docker.Container, error)
	
	
	UploadToContainer(id string, opts docker.UploadToContainerOptions) error
	
	StartContainer(id string, cfg *docker.HostConfig) error
	
	
	AttachToContainer(opts docker.AttachToContainerOptions) error
	
	
	BuildImage(opts docker.BuildImageOptions) error
	
	
	RemoveImageExtended(id string, opts docker.RemoveImageOptions) error
	
	
	StopContainer(id string, timeout uint) error
	
	
	KillContainer(opts docker.KillContainerOptions) error
	
	RemoveContainer(opts docker.RemoveContainerOptions) error
	
	
	PingWithContext(context.Context) error
}


type Provider struct {
	PeerID       string
	NetworkID    string
	BuildMetrics *BuildMetrics
}


func NewProvider(peerID, networkID string, metricsProvider metrics.Provider) *Provider {
	return &Provider{
		PeerID:       peerID,
		NetworkID:    networkID,
		BuildMetrics: NewBuildMetrics(metricsProvider),
	}
}


func (p *Provider) NewVM() container.VM {
	return NewDockerVM(p.PeerID, p.NetworkID, p.BuildMetrics)
}


func NewDockerVM(peerID, networkID string, buildMetrics *BuildMetrics) *DockerVM {
	return &DockerVM{
		PeerID:       peerID,
		NetworkID:    networkID,
		getClientFnc: getDockerClient,
		BuildMetrics: buildMetrics,
	}
}

func getDockerClient() (dockerClient, error) {
	return cutil.NewDockerClient()
}

func getDockerHostConfig() *docker.HostConfig {
	if hostConfig != nil {
		return hostConfig
	}

	dockerKey := func(key string) string { return "vm.docker.hostConfig." + key }
	getInt64 := func(key string) int64 { return int64(viper.GetInt(dockerKey(key))) }

	var logConfig docker.LogConfig
	err := viper.UnmarshalKey(dockerKey("LogConfig"), &logConfig)
	if err != nil {
		dockerLogger.Warningf("load docker HostConfig.LogConfig failed, error: %s", err.Error())
	}
	networkMode := viper.GetString(dockerKey("NetworkMode"))
	if networkMode == "" {
		networkMode = "host"
	}
	dockerLogger.Debugf("docker container hostconfig NetworkMode: %s", networkMode)

	return &docker.HostConfig{
		CapAdd:  viper.GetStringSlice(dockerKey("CapAdd")),
		CapDrop: viper.GetStringSlice(dockerKey("CapDrop")),

		DNS:         viper.GetStringSlice(dockerKey("Dns")),
		DNSSearch:   viper.GetStringSlice(dockerKey("DnsSearch")),
		ExtraHosts:  viper.GetStringSlice(dockerKey("ExtraHosts")),
		NetworkMode: networkMode,
		IpcMode:     viper.GetString(dockerKey("IpcMode")),
		PidMode:     viper.GetString(dockerKey("PidMode")),
		UTSMode:     viper.GetString(dockerKey("UTSMode")),
		LogConfig:   logConfig,

		ReadonlyRootfs:   viper.GetBool(dockerKey("ReadonlyRootfs")),
		SecurityOpt:      viper.GetStringSlice(dockerKey("SecurityOpt")),
		CgroupParent:     viper.GetString(dockerKey("CgroupParent")),
		Memory:           getInt64("Memory"),
		MemorySwap:       getInt64("MemorySwap"),
		MemorySwappiness: getInt64("MemorySwappiness"),
		OOMKillDisable:   viper.GetBool(dockerKey("OomKillDisable")),
		CPUShares:        getInt64("CpuShares"),
		CPUSet:           viper.GetString(dockerKey("Cpuset")),
		CPUSetCPUs:       viper.GetString(dockerKey("CpusetCPUs")),
		CPUSetMEMs:       viper.GetString(dockerKey("CpusetMEMs")),
		CPUQuota:         getInt64("CpuQuota"),
		CPUPeriod:        getInt64("CpuPeriod"),
		BlkioWeight:      getInt64("BlkioWeight"),
	}
}

func (vm *DockerVM) createContainer(client dockerClient, imageID, containerID string, args, env []string, attachStdout bool) error {
	logger := dockerLogger.With("imageID", imageID, "containerID", containerID)
	logger.Debugw("create container")
	_, err := client.CreateContainer(docker.CreateContainerOptions{
		Name: containerID,
		Config: &docker.Config{
			Cmd:          args,
			Image:        imageID,
			Env:          env,
			AttachStdout: attachStdout,
			AttachStderr: attachStdout,
		},
		HostConfig: getDockerHostConfig(),
	})
	if err != nil {
		return err
	}
	logger.Debugw("created container")
	return nil
}

func (vm *DockerVM) deployImage(client dockerClient, ccid ccintf.CCID, reader io.Reader) error {
	id, err := vm.GetVMNameForDocker(ccid)
	if err != nil {
		return err
	}

	outputbuf := bytes.NewBuffer(nil)
	opts := docker.BuildImageOptions{
		Name:         id,
		Pull:         viper.GetBool("chaincode.pull"),
		InputStream:  reader,
		OutputStream: outputbuf,
	}

	startTime := time.Now()
	err = client.BuildImage(opts)

	vm.BuildMetrics.ChaincodeImageBuildDuration.With(
		"chaincode", ccid.Name+":"+ccid.Version,
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


func (vm *DockerVM) Start(ccid ccintf.CCID, args, env []string, filesToUpload map[string][]byte, builder container.Builder) error {
	imageName, err := vm.GetVMNameForDocker(ccid)
	if err != nil {
		return err
	}

	attachStdout := viper.GetBool("vm.docker.attachStdout")
	containerName := vm.GetVMName(ccid)
	logger := dockerLogger.With("imageName", imageName, "containerName", containerName)

	client, err := vm.getClientFnc()
	if err != nil {
		logger.Debugf("failed to get docker client", "error", err)
		return err
	}

	vm.stopInternal(client, containerName, 0, false, false)

	err = vm.createContainer(client, imageName, containerName, args, env, attachStdout)
	if err == docker.ErrNoSuchImage {
		reader, err := builder.Build()
		if err != nil {
			return errors.Wrapf(err, "failed to generate Dockerfile to build %s", containerName)
		}

		err = vm.deployImage(client, ccid, reader)
		if err != nil {
			return err
		}

		err = vm.createContainer(client, imageName, containerName, args, env, attachStdout)
		if err != nil {
			logger.Errorf("failed to create container: %s", err)
			return err
		}
	} else if err != nil {
		logger.Errorf("create container failed: %s", err)
		return err
	}

	
	if attachStdout {
		go vm.streamOutput(client, containerName)
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

		err := client.UploadToContainer(containerName, docker.UploadToContainerOptions{
			InputStream:          bytes.NewReader(payload.Bytes()),
			Path:                 "/",
			NoOverwriteDirNonDir: false,
		})
		if err != nil {
			return fmt.Errorf("Error uploading files to the container instance %s: %s", containerName, err)
		}
	}

	
	err = client.StartContainer(containerName, nil)
	if err != nil {
		dockerLogger.Errorf("start-could not start container: %s", err)
		return err
	}

	dockerLogger.Debugf("Started container %s", containerName)
	return nil
}


func (vm *DockerVM) streamOutput(client dockerClient, containerName string) {
	
	
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
		
		select {
		case <-attached: 
			close(attached) 

		case <-time.After(10 * time.Second):
			dockerLogger.Errorf("Timeout while attaching to IO channel in container %s", containerName)
			return
		}

		
		containerLogger := flogging.MustGetLogger("peer.chaincode." + containerName)

		is := bufio.NewReader(r)
		for {
			
			
			line, err := is.ReadString('\n')
			switch err {
			case nil:
				containerLogger.Info(line)
			case io.EOF:
				dockerLogger.Infof("Container %s has closed its IO channel", containerName)
			default:
				dockerLogger.Errorf("Error reading container output: %s", err)
			}
		}
	}()
}


func (vm *DockerVM) Stop(ccid ccintf.CCID, timeout uint, dontkill bool, dontremove bool) error {
	client, err := vm.getClientFnc()
	if err != nil {
		dockerLogger.Debugf("stop - cannot create client %s", err)
		return err
	}
	id := strings.Replace(vm.GetVMName(ccid), ":", "_", -1)

	return vm.stopInternal(client, id, timeout, dontkill, dontremove)
}



func (vm *DockerVM) HealthCheck(ctx context.Context) error {
	client, err := vm.getClientFnc()
	if err != nil {
		return errors.Wrap(err, "failed to connect to Docker daemon")
	}
	if err := client.PingWithContext(ctx); err != nil {
		return errors.Wrap(err, "failed to ping to Docker daemon")
	}
	return nil
}

func (vm *DockerVM) stopInternal(client dockerClient, id string, timeout uint, dontkill, dontremove bool) error {
	logger := dockerLogger.With("id", id)

	logger.Debugw("stopping container")
	err := client.StopContainer(id, timeout)
	dockerLogger.Debugw("stop container result", "error", err)

	if !dontkill {
		logger.Debugw("killing container")
		err = client.KillContainer(docker.KillContainerOptions{ID: id})
		logger.Debugw("kill container result", "error", err)
	}

	if !dontremove {
		logger.Debugw("removing container")
		err = client.RemoveContainer(docker.RemoveContainerOptions{ID: id, Force: true})
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
	name := ccid.GetName()

	if vm.NetworkID != "" && vm.PeerID != "" {
		name = fmt.Sprintf("%s-%s-%s", vm.NetworkID, vm.PeerID, name)
	} else if vm.NetworkID != "" {
		name = fmt.Sprintf("%s-%s", vm.NetworkID, name)
	} else if vm.PeerID != "" {
		name = fmt.Sprintf("%s-%s", vm.PeerID, name)
	}

	return name
}
