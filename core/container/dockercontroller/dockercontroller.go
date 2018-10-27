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
	"encoding/hex"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"

	"github.com/fsouza/go-dockerclient"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/common/util"
	"github.com/mcc-github/blockchain/core/container"
	"github.com/mcc-github/blockchain/core/container/ccintf"
	cutil "github.com/mcc-github/blockchain/core/container/util"
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
}


type Provider struct {
	PeerID    string
	NetworkID string
}


func NewProvider(peerID, networkID string) *Provider {
	return &Provider{
		PeerID:    peerID,
		NetworkID: networkID,
	}
}


func (p *Provider) NewVM() container.VM {
	return NewDockerVM(p.PeerID, p.NetworkID)
}


func NewDockerVM(peerID, networkID string) *DockerVM {
	vm := DockerVM{
		PeerID:    peerID,
		NetworkID: networkID,
	}
	vm.getClientFnc = getDockerClient
	return &vm
}

func getDockerClient() (dockerClient, error) {
	return cutil.NewDockerClient()
}

func getDockerHostConfig() *docker.HostConfig {
	if hostConfig != nil {
		return hostConfig
	}
	dockerKey := func(key string) string {
		return "vm.docker.hostConfig." + key
	}
	getInt64 := func(key string) int64 {
		defer func() {
			if err := recover(); err != nil {
				dockerLogger.Warningf("load vm.docker.hostConfig.%s failed, error: %v", key, err)
			}
		}()
		n := viper.GetInt(dockerKey(key))
		return int64(n)
	}

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

	hostConfig = &docker.HostConfig{
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

	return hostConfig
}

func (vm *DockerVM) createContainer(client dockerClient,
	imageID string, containerID string, args []string,
	env []string, attachStdout bool) error {
	config := docker.Config{Cmd: args, Image: imageID, Env: env, AttachStdout: attachStdout, AttachStderr: attachStdout}
	copts := docker.CreateContainerOptions{Name: containerID, Config: &config, HostConfig: getDockerHostConfig()}
	dockerLogger.Debugf("Create container: %s", containerID)
	_, err := client.CreateContainer(copts)
	if err != nil {
		return err
	}
	dockerLogger.Debugf("Created container: %s", imageID)
	return nil
}

func (vm *DockerVM) deployImage(client dockerClient, ccid ccintf.CCID,
	args []string, env []string, reader io.Reader) error {
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

	if err := client.BuildImage(opts); err != nil {
		dockerLogger.Errorf("Error building images: %s", err)
		dockerLogger.Errorf("Image Output:\n********************\n%s\n********************", outputbuf.String())
		return err
	}

	dockerLogger.Debugf("Created image: %s", id)

	return nil
}


func (vm *DockerVM) Start(ccid ccintf.CCID,
	args []string, env []string, filesToUpload map[string][]byte, builder container.Builder) error {
	imageName, err := vm.GetVMNameForDocker(ccid)
	if err != nil {
		return err
	}

	client, err := vm.getClientFnc()
	if err != nil {
		dockerLogger.Debugf("start - cannot create client %s", err)
		return err
	}

	containerName := vm.GetVMName(ccid)

	attachStdout := viper.GetBool("vm.docker.attachStdout")

	
	dockerLogger.Debugf("Cleanup container %s", containerName)
	vm.stopInternal(client, containerName, 0, false, false)

	dockerLogger.Debugf("Start container %s", containerName)
	err = vm.createContainer(client, imageName, containerName, args, env, attachStdout)
	if err != nil {
		
		if err == docker.ErrNoSuchImage {
			if builder != nil {
				dockerLogger.Debugf("start-could not find image <%s> (container id <%s>), because of <%s>..."+
					"attempt to recreate image", imageName, containerName, err)

				reader, err1 := builder.Build()
				if err1 != nil {
					dockerLogger.Errorf("Error creating image builder for image <%s> (container id <%s>), "+
						"because of <%s>", imageName, containerName, err1)
				}

				if err1 = vm.deployImage(client, ccid, args, env, reader); err1 != nil {
					return err1
				}

				dockerLogger.Debug("start-recreated image successfully")
				if err1 = vm.createContainer(client, imageName, containerName, args, env, attachStdout); err1 != nil {
					dockerLogger.Errorf("start-could not recreate container post recreate image: %s", err1)
					return err1
				}
			} else {
				dockerLogger.Errorf("start-could not find image <%s>, because of %s", imageName, err)
				return err
			}
		} else {
			dockerLogger.Errorf("start-could not recreate container <%s>, because of %s", containerName, err)
			return err
		}
	}

	if attachStdout {
		
		
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
				
			case <-time.After(10 * time.Second):
				dockerLogger.Errorf("Timeout while attaching to IO channel in container %s", containerName)
				return
			}

			
			
			
			attached <- struct{}{}

			
			
			is := bufio.NewReader(r)

			
			containerLogger := flogging.MustGetLogger(containerName)
			flogging.SetModuleLevel(flogging.GetModuleLevel("peer"), containerName)

			for {
				
				
				line, err2 := is.ReadString('\n')
				if err2 != nil {
					switch err2 {
					case io.EOF:
						dockerLogger.Infof("Container %s has closed its IO channel", containerName)
					default:
						dockerLogger.Errorf("Error reading container output: %s", err2)
					}

					return
				}

				containerLogger.Info(line)
			}
		}()
	}

	
	
	if len(filesToUpload) != 0 {
		
		
		payload := bytes.NewBuffer(nil)
		gw := gzip.NewWriter(payload)
		tw := tar.NewWriter(gw)

		for path, fileToUpload := range filesToUpload {
			cutil.WriteBytesToPackage(path, fileToUpload, tw)
		}

		
		if err = tw.Close(); err != nil {
			return fmt.Errorf("Error writing files to upload to Docker instance into a temporary tar blob: %s", err)
		}

		gw.Close()

		err = client.UploadToContainer(containerName, docker.UploadToContainerOptions{
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


func (vm *DockerVM) Stop(ccid ccintf.CCID, timeout uint, dontkill bool, dontremove bool) error {
	id := vm.GetVMName(ccid)

	client, err := vm.getClientFnc()
	if err != nil {
		dockerLogger.Debugf("stop - cannot create client %s", err)
		return err
	}
	id = strings.Replace(id, ":", "_", -1)

	err = vm.stopInternal(client, id, timeout, dontkill, dontremove)

	return err
}

func (vm *DockerVM) stopInternal(client dockerClient,
	id string, timeout uint, dontkill bool, dontremove bool) error {
	err := client.StopContainer(id, timeout)
	if err != nil {
		dockerLogger.Debugf("Stop container %s(%s)", id, err)
	} else {
		dockerLogger.Debugf("Stopped container %s", id)
	}
	if !dontkill {
		err = client.KillContainer(docker.KillContainerOptions{ID: id})
		if err != nil {
			dockerLogger.Debugf("Kill container %s (%s)", id, err)
		} else {
			dockerLogger.Debugf("Killed container %s", id)
		}
	}
	if !dontremove {
		err = client.RemoveContainer(docker.RemoveContainerOptions{ID: id, Force: true})
		if err != nil {
			dockerLogger.Debugf("Remove container %s (%s)", id, err)
		} else {
			dockerLogger.Debugf("Removed container %s", id)
		}
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
		return imageName, fmt.Errorf("Error constructing Docker VM Name. '%s' breaks Docker's repository naming rules", imageName)
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
