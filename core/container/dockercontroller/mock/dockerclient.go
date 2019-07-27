
package mock

import (
	"context"
	"sync"

	docker "github.com/fsouza/go-dockerclient"
)

type DockerClient struct {
	AttachToContainerStub        func(docker.AttachToContainerOptions) error
	attachToContainerMutex       sync.RWMutex
	attachToContainerArgsForCall []struct {
		arg1 docker.AttachToContainerOptions
	}
	attachToContainerReturns struct {
		result1 error
	}
	attachToContainerReturnsOnCall map[int]struct {
		result1 error
	}
	BuildImageStub        func(docker.BuildImageOptions) error
	buildImageMutex       sync.RWMutex
	buildImageArgsForCall []struct {
		arg1 docker.BuildImageOptions
	}
	buildImageReturns struct {
		result1 error
	}
	buildImageReturnsOnCall map[int]struct {
		result1 error
	}
	CreateContainerStub        func(docker.CreateContainerOptions) (*docker.Container, error)
	createContainerMutex       sync.RWMutex
	createContainerArgsForCall []struct {
		arg1 docker.CreateContainerOptions
	}
	createContainerReturns struct {
		result1 *docker.Container
		result2 error
	}
	createContainerReturnsOnCall map[int]struct {
		result1 *docker.Container
		result2 error
	}
	InspectImageStub        func(string) (*docker.Image, error)
	inspectImageMutex       sync.RWMutex
	inspectImageArgsForCall []struct {
		arg1 string
	}
	inspectImageReturns struct {
		result1 *docker.Image
		result2 error
	}
	inspectImageReturnsOnCall map[int]struct {
		result1 *docker.Image
		result2 error
	}
	KillContainerStub        func(docker.KillContainerOptions) error
	killContainerMutex       sync.RWMutex
	killContainerArgsForCall []struct {
		arg1 docker.KillContainerOptions
	}
	killContainerReturns struct {
		result1 error
	}
	killContainerReturnsOnCall map[int]struct {
		result1 error
	}
	PingWithContextStub        func(context.Context) error
	pingWithContextMutex       sync.RWMutex
	pingWithContextArgsForCall []struct {
		arg1 context.Context
	}
	pingWithContextReturns struct {
		result1 error
	}
	pingWithContextReturnsOnCall map[int]struct {
		result1 error
	}
	RemoveContainerStub        func(docker.RemoveContainerOptions) error
	removeContainerMutex       sync.RWMutex
	removeContainerArgsForCall []struct {
		arg1 docker.RemoveContainerOptions
	}
	removeContainerReturns struct {
		result1 error
	}
	removeContainerReturnsOnCall map[int]struct {
		result1 error
	}
	StartContainerStub        func(string, *docker.HostConfig) error
	startContainerMutex       sync.RWMutex
	startContainerArgsForCall []struct {
		arg1 string
		arg2 *docker.HostConfig
	}
	startContainerReturns struct {
		result1 error
	}
	startContainerReturnsOnCall map[int]struct {
		result1 error
	}
	StopContainerStub        func(string, uint) error
	stopContainerMutex       sync.RWMutex
	stopContainerArgsForCall []struct {
		arg1 string
		arg2 uint
	}
	stopContainerReturns struct {
		result1 error
	}
	stopContainerReturnsOnCall map[int]struct {
		result1 error
	}
	UploadToContainerStub        func(string, docker.UploadToContainerOptions) error
	uploadToContainerMutex       sync.RWMutex
	uploadToContainerArgsForCall []struct {
		arg1 string
		arg2 docker.UploadToContainerOptions
	}
	uploadToContainerReturns struct {
		result1 error
	}
	uploadToContainerReturnsOnCall map[int]struct {
		result1 error
	}
	WaitContainerStub        func(string) (int, error)
	waitContainerMutex       sync.RWMutex
	waitContainerArgsForCall []struct {
		arg1 string
	}
	waitContainerReturns struct {
		result1 int
		result2 error
	}
	waitContainerReturnsOnCall map[int]struct {
		result1 int
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *DockerClient) AttachToContainer(arg1 docker.AttachToContainerOptions) error {
	fake.attachToContainerMutex.Lock()
	ret, specificReturn := fake.attachToContainerReturnsOnCall[len(fake.attachToContainerArgsForCall)]
	fake.attachToContainerArgsForCall = append(fake.attachToContainerArgsForCall, struct {
		arg1 docker.AttachToContainerOptions
	}{arg1})
	fake.recordInvocation("AttachToContainer", []interface{}{arg1})
	fake.attachToContainerMutex.Unlock()
	if fake.AttachToContainerStub != nil {
		return fake.AttachToContainerStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.attachToContainerReturns
	return fakeReturns.result1
}

func (fake *DockerClient) AttachToContainerCallCount() int {
	fake.attachToContainerMutex.RLock()
	defer fake.attachToContainerMutex.RUnlock()
	return len(fake.attachToContainerArgsForCall)
}

func (fake *DockerClient) AttachToContainerCalls(stub func(docker.AttachToContainerOptions) error) {
	fake.attachToContainerMutex.Lock()
	defer fake.attachToContainerMutex.Unlock()
	fake.AttachToContainerStub = stub
}

func (fake *DockerClient) AttachToContainerArgsForCall(i int) docker.AttachToContainerOptions {
	fake.attachToContainerMutex.RLock()
	defer fake.attachToContainerMutex.RUnlock()
	argsForCall := fake.attachToContainerArgsForCall[i]
	return argsForCall.arg1
}

func (fake *DockerClient) AttachToContainerReturns(result1 error) {
	fake.attachToContainerMutex.Lock()
	defer fake.attachToContainerMutex.Unlock()
	fake.AttachToContainerStub = nil
	fake.attachToContainerReturns = struct {
		result1 error
	}{result1}
}

func (fake *DockerClient) AttachToContainerReturnsOnCall(i int, result1 error) {
	fake.attachToContainerMutex.Lock()
	defer fake.attachToContainerMutex.Unlock()
	fake.AttachToContainerStub = nil
	if fake.attachToContainerReturnsOnCall == nil {
		fake.attachToContainerReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.attachToContainerReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *DockerClient) BuildImage(arg1 docker.BuildImageOptions) error {
	fake.buildImageMutex.Lock()
	ret, specificReturn := fake.buildImageReturnsOnCall[len(fake.buildImageArgsForCall)]
	fake.buildImageArgsForCall = append(fake.buildImageArgsForCall, struct {
		arg1 docker.BuildImageOptions
	}{arg1})
	fake.recordInvocation("BuildImage", []interface{}{arg1})
	fake.buildImageMutex.Unlock()
	if fake.BuildImageStub != nil {
		return fake.BuildImageStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.buildImageReturns
	return fakeReturns.result1
}

func (fake *DockerClient) BuildImageCallCount() int {
	fake.buildImageMutex.RLock()
	defer fake.buildImageMutex.RUnlock()
	return len(fake.buildImageArgsForCall)
}

func (fake *DockerClient) BuildImageCalls(stub func(docker.BuildImageOptions) error) {
	fake.buildImageMutex.Lock()
	defer fake.buildImageMutex.Unlock()
	fake.BuildImageStub = stub
}

func (fake *DockerClient) BuildImageArgsForCall(i int) docker.BuildImageOptions {
	fake.buildImageMutex.RLock()
	defer fake.buildImageMutex.RUnlock()
	argsForCall := fake.buildImageArgsForCall[i]
	return argsForCall.arg1
}

func (fake *DockerClient) BuildImageReturns(result1 error) {
	fake.buildImageMutex.Lock()
	defer fake.buildImageMutex.Unlock()
	fake.BuildImageStub = nil
	fake.buildImageReturns = struct {
		result1 error
	}{result1}
}

func (fake *DockerClient) BuildImageReturnsOnCall(i int, result1 error) {
	fake.buildImageMutex.Lock()
	defer fake.buildImageMutex.Unlock()
	fake.BuildImageStub = nil
	if fake.buildImageReturnsOnCall == nil {
		fake.buildImageReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.buildImageReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *DockerClient) CreateContainer(arg1 docker.CreateContainerOptions) (*docker.Container, error) {
	fake.createContainerMutex.Lock()
	ret, specificReturn := fake.createContainerReturnsOnCall[len(fake.createContainerArgsForCall)]
	fake.createContainerArgsForCall = append(fake.createContainerArgsForCall, struct {
		arg1 docker.CreateContainerOptions
	}{arg1})
	fake.recordInvocation("CreateContainer", []interface{}{arg1})
	fake.createContainerMutex.Unlock()
	if fake.CreateContainerStub != nil {
		return fake.CreateContainerStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.createContainerReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *DockerClient) CreateContainerCallCount() int {
	fake.createContainerMutex.RLock()
	defer fake.createContainerMutex.RUnlock()
	return len(fake.createContainerArgsForCall)
}

func (fake *DockerClient) CreateContainerCalls(stub func(docker.CreateContainerOptions) (*docker.Container, error)) {
	fake.createContainerMutex.Lock()
	defer fake.createContainerMutex.Unlock()
	fake.CreateContainerStub = stub
}

func (fake *DockerClient) CreateContainerArgsForCall(i int) docker.CreateContainerOptions {
	fake.createContainerMutex.RLock()
	defer fake.createContainerMutex.RUnlock()
	argsForCall := fake.createContainerArgsForCall[i]
	return argsForCall.arg1
}

func (fake *DockerClient) CreateContainerReturns(result1 *docker.Container, result2 error) {
	fake.createContainerMutex.Lock()
	defer fake.createContainerMutex.Unlock()
	fake.CreateContainerStub = nil
	fake.createContainerReturns = struct {
		result1 *docker.Container
		result2 error
	}{result1, result2}
}

func (fake *DockerClient) CreateContainerReturnsOnCall(i int, result1 *docker.Container, result2 error) {
	fake.createContainerMutex.Lock()
	defer fake.createContainerMutex.Unlock()
	fake.CreateContainerStub = nil
	if fake.createContainerReturnsOnCall == nil {
		fake.createContainerReturnsOnCall = make(map[int]struct {
			result1 *docker.Container
			result2 error
		})
	}
	fake.createContainerReturnsOnCall[i] = struct {
		result1 *docker.Container
		result2 error
	}{result1, result2}
}

func (fake *DockerClient) InspectImage(arg1 string) (*docker.Image, error) {
	fake.inspectImageMutex.Lock()
	ret, specificReturn := fake.inspectImageReturnsOnCall[len(fake.inspectImageArgsForCall)]
	fake.inspectImageArgsForCall = append(fake.inspectImageArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("InspectImage", []interface{}{arg1})
	fake.inspectImageMutex.Unlock()
	if fake.InspectImageStub != nil {
		return fake.InspectImageStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.inspectImageReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *DockerClient) InspectImageCallCount() int {
	fake.inspectImageMutex.RLock()
	defer fake.inspectImageMutex.RUnlock()
	return len(fake.inspectImageArgsForCall)
}

func (fake *DockerClient) InspectImageCalls(stub func(string) (*docker.Image, error)) {
	fake.inspectImageMutex.Lock()
	defer fake.inspectImageMutex.Unlock()
	fake.InspectImageStub = stub
}

func (fake *DockerClient) InspectImageArgsForCall(i int) string {
	fake.inspectImageMutex.RLock()
	defer fake.inspectImageMutex.RUnlock()
	argsForCall := fake.inspectImageArgsForCall[i]
	return argsForCall.arg1
}

func (fake *DockerClient) InspectImageReturns(result1 *docker.Image, result2 error) {
	fake.inspectImageMutex.Lock()
	defer fake.inspectImageMutex.Unlock()
	fake.InspectImageStub = nil
	fake.inspectImageReturns = struct {
		result1 *docker.Image
		result2 error
	}{result1, result2}
}

func (fake *DockerClient) InspectImageReturnsOnCall(i int, result1 *docker.Image, result2 error) {
	fake.inspectImageMutex.Lock()
	defer fake.inspectImageMutex.Unlock()
	fake.InspectImageStub = nil
	if fake.inspectImageReturnsOnCall == nil {
		fake.inspectImageReturnsOnCall = make(map[int]struct {
			result1 *docker.Image
			result2 error
		})
	}
	fake.inspectImageReturnsOnCall[i] = struct {
		result1 *docker.Image
		result2 error
	}{result1, result2}
}

func (fake *DockerClient) KillContainer(arg1 docker.KillContainerOptions) error {
	fake.killContainerMutex.Lock()
	ret, specificReturn := fake.killContainerReturnsOnCall[len(fake.killContainerArgsForCall)]
	fake.killContainerArgsForCall = append(fake.killContainerArgsForCall, struct {
		arg1 docker.KillContainerOptions
	}{arg1})
	fake.recordInvocation("KillContainer", []interface{}{arg1})
	fake.killContainerMutex.Unlock()
	if fake.KillContainerStub != nil {
		return fake.KillContainerStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.killContainerReturns
	return fakeReturns.result1
}

func (fake *DockerClient) KillContainerCallCount() int {
	fake.killContainerMutex.RLock()
	defer fake.killContainerMutex.RUnlock()
	return len(fake.killContainerArgsForCall)
}

func (fake *DockerClient) KillContainerCalls(stub func(docker.KillContainerOptions) error) {
	fake.killContainerMutex.Lock()
	defer fake.killContainerMutex.Unlock()
	fake.KillContainerStub = stub
}

func (fake *DockerClient) KillContainerArgsForCall(i int) docker.KillContainerOptions {
	fake.killContainerMutex.RLock()
	defer fake.killContainerMutex.RUnlock()
	argsForCall := fake.killContainerArgsForCall[i]
	return argsForCall.arg1
}

func (fake *DockerClient) KillContainerReturns(result1 error) {
	fake.killContainerMutex.Lock()
	defer fake.killContainerMutex.Unlock()
	fake.KillContainerStub = nil
	fake.killContainerReturns = struct {
		result1 error
	}{result1}
}

func (fake *DockerClient) KillContainerReturnsOnCall(i int, result1 error) {
	fake.killContainerMutex.Lock()
	defer fake.killContainerMutex.Unlock()
	fake.KillContainerStub = nil
	if fake.killContainerReturnsOnCall == nil {
		fake.killContainerReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.killContainerReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *DockerClient) PingWithContext(arg1 context.Context) error {
	fake.pingWithContextMutex.Lock()
	ret, specificReturn := fake.pingWithContextReturnsOnCall[len(fake.pingWithContextArgsForCall)]
	fake.pingWithContextArgsForCall = append(fake.pingWithContextArgsForCall, struct {
		arg1 context.Context
	}{arg1})
	fake.recordInvocation("PingWithContext", []interface{}{arg1})
	fake.pingWithContextMutex.Unlock()
	if fake.PingWithContextStub != nil {
		return fake.PingWithContextStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.pingWithContextReturns
	return fakeReturns.result1
}

func (fake *DockerClient) PingWithContextCallCount() int {
	fake.pingWithContextMutex.RLock()
	defer fake.pingWithContextMutex.RUnlock()
	return len(fake.pingWithContextArgsForCall)
}

func (fake *DockerClient) PingWithContextCalls(stub func(context.Context) error) {
	fake.pingWithContextMutex.Lock()
	defer fake.pingWithContextMutex.Unlock()
	fake.PingWithContextStub = stub
}

func (fake *DockerClient) PingWithContextArgsForCall(i int) context.Context {
	fake.pingWithContextMutex.RLock()
	defer fake.pingWithContextMutex.RUnlock()
	argsForCall := fake.pingWithContextArgsForCall[i]
	return argsForCall.arg1
}

func (fake *DockerClient) PingWithContextReturns(result1 error) {
	fake.pingWithContextMutex.Lock()
	defer fake.pingWithContextMutex.Unlock()
	fake.PingWithContextStub = nil
	fake.pingWithContextReturns = struct {
		result1 error
	}{result1}
}

func (fake *DockerClient) PingWithContextReturnsOnCall(i int, result1 error) {
	fake.pingWithContextMutex.Lock()
	defer fake.pingWithContextMutex.Unlock()
	fake.PingWithContextStub = nil
	if fake.pingWithContextReturnsOnCall == nil {
		fake.pingWithContextReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.pingWithContextReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *DockerClient) RemoveContainer(arg1 docker.RemoveContainerOptions) error {
	fake.removeContainerMutex.Lock()
	ret, specificReturn := fake.removeContainerReturnsOnCall[len(fake.removeContainerArgsForCall)]
	fake.removeContainerArgsForCall = append(fake.removeContainerArgsForCall, struct {
		arg1 docker.RemoveContainerOptions
	}{arg1})
	fake.recordInvocation("RemoveContainer", []interface{}{arg1})
	fake.removeContainerMutex.Unlock()
	if fake.RemoveContainerStub != nil {
		return fake.RemoveContainerStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.removeContainerReturns
	return fakeReturns.result1
}

func (fake *DockerClient) RemoveContainerCallCount() int {
	fake.removeContainerMutex.RLock()
	defer fake.removeContainerMutex.RUnlock()
	return len(fake.removeContainerArgsForCall)
}

func (fake *DockerClient) RemoveContainerCalls(stub func(docker.RemoveContainerOptions) error) {
	fake.removeContainerMutex.Lock()
	defer fake.removeContainerMutex.Unlock()
	fake.RemoveContainerStub = stub
}

func (fake *DockerClient) RemoveContainerArgsForCall(i int) docker.RemoveContainerOptions {
	fake.removeContainerMutex.RLock()
	defer fake.removeContainerMutex.RUnlock()
	argsForCall := fake.removeContainerArgsForCall[i]
	return argsForCall.arg1
}

func (fake *DockerClient) RemoveContainerReturns(result1 error) {
	fake.removeContainerMutex.Lock()
	defer fake.removeContainerMutex.Unlock()
	fake.RemoveContainerStub = nil
	fake.removeContainerReturns = struct {
		result1 error
	}{result1}
}

func (fake *DockerClient) RemoveContainerReturnsOnCall(i int, result1 error) {
	fake.removeContainerMutex.Lock()
	defer fake.removeContainerMutex.Unlock()
	fake.RemoveContainerStub = nil
	if fake.removeContainerReturnsOnCall == nil {
		fake.removeContainerReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.removeContainerReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *DockerClient) StartContainer(arg1 string, arg2 *docker.HostConfig) error {
	fake.startContainerMutex.Lock()
	ret, specificReturn := fake.startContainerReturnsOnCall[len(fake.startContainerArgsForCall)]
	fake.startContainerArgsForCall = append(fake.startContainerArgsForCall, struct {
		arg1 string
		arg2 *docker.HostConfig
	}{arg1, arg2})
	fake.recordInvocation("StartContainer", []interface{}{arg1, arg2})
	fake.startContainerMutex.Unlock()
	if fake.StartContainerStub != nil {
		return fake.StartContainerStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.startContainerReturns
	return fakeReturns.result1
}

func (fake *DockerClient) StartContainerCallCount() int {
	fake.startContainerMutex.RLock()
	defer fake.startContainerMutex.RUnlock()
	return len(fake.startContainerArgsForCall)
}

func (fake *DockerClient) StartContainerCalls(stub func(string, *docker.HostConfig) error) {
	fake.startContainerMutex.Lock()
	defer fake.startContainerMutex.Unlock()
	fake.StartContainerStub = stub
}

func (fake *DockerClient) StartContainerArgsForCall(i int) (string, *docker.HostConfig) {
	fake.startContainerMutex.RLock()
	defer fake.startContainerMutex.RUnlock()
	argsForCall := fake.startContainerArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *DockerClient) StartContainerReturns(result1 error) {
	fake.startContainerMutex.Lock()
	defer fake.startContainerMutex.Unlock()
	fake.StartContainerStub = nil
	fake.startContainerReturns = struct {
		result1 error
	}{result1}
}

func (fake *DockerClient) StartContainerReturnsOnCall(i int, result1 error) {
	fake.startContainerMutex.Lock()
	defer fake.startContainerMutex.Unlock()
	fake.StartContainerStub = nil
	if fake.startContainerReturnsOnCall == nil {
		fake.startContainerReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.startContainerReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *DockerClient) StopContainer(arg1 string, arg2 uint) error {
	fake.stopContainerMutex.Lock()
	ret, specificReturn := fake.stopContainerReturnsOnCall[len(fake.stopContainerArgsForCall)]
	fake.stopContainerArgsForCall = append(fake.stopContainerArgsForCall, struct {
		arg1 string
		arg2 uint
	}{arg1, arg2})
	fake.recordInvocation("StopContainer", []interface{}{arg1, arg2})
	fake.stopContainerMutex.Unlock()
	if fake.StopContainerStub != nil {
		return fake.StopContainerStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.stopContainerReturns
	return fakeReturns.result1
}

func (fake *DockerClient) StopContainerCallCount() int {
	fake.stopContainerMutex.RLock()
	defer fake.stopContainerMutex.RUnlock()
	return len(fake.stopContainerArgsForCall)
}

func (fake *DockerClient) StopContainerCalls(stub func(string, uint) error) {
	fake.stopContainerMutex.Lock()
	defer fake.stopContainerMutex.Unlock()
	fake.StopContainerStub = stub
}

func (fake *DockerClient) StopContainerArgsForCall(i int) (string, uint) {
	fake.stopContainerMutex.RLock()
	defer fake.stopContainerMutex.RUnlock()
	argsForCall := fake.stopContainerArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *DockerClient) StopContainerReturns(result1 error) {
	fake.stopContainerMutex.Lock()
	defer fake.stopContainerMutex.Unlock()
	fake.StopContainerStub = nil
	fake.stopContainerReturns = struct {
		result1 error
	}{result1}
}

func (fake *DockerClient) StopContainerReturnsOnCall(i int, result1 error) {
	fake.stopContainerMutex.Lock()
	defer fake.stopContainerMutex.Unlock()
	fake.StopContainerStub = nil
	if fake.stopContainerReturnsOnCall == nil {
		fake.stopContainerReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.stopContainerReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *DockerClient) UploadToContainer(arg1 string, arg2 docker.UploadToContainerOptions) error {
	fake.uploadToContainerMutex.Lock()
	ret, specificReturn := fake.uploadToContainerReturnsOnCall[len(fake.uploadToContainerArgsForCall)]
	fake.uploadToContainerArgsForCall = append(fake.uploadToContainerArgsForCall, struct {
		arg1 string
		arg2 docker.UploadToContainerOptions
	}{arg1, arg2})
	fake.recordInvocation("UploadToContainer", []interface{}{arg1, arg2})
	fake.uploadToContainerMutex.Unlock()
	if fake.UploadToContainerStub != nil {
		return fake.UploadToContainerStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.uploadToContainerReturns
	return fakeReturns.result1
}

func (fake *DockerClient) UploadToContainerCallCount() int {
	fake.uploadToContainerMutex.RLock()
	defer fake.uploadToContainerMutex.RUnlock()
	return len(fake.uploadToContainerArgsForCall)
}

func (fake *DockerClient) UploadToContainerCalls(stub func(string, docker.UploadToContainerOptions) error) {
	fake.uploadToContainerMutex.Lock()
	defer fake.uploadToContainerMutex.Unlock()
	fake.UploadToContainerStub = stub
}

func (fake *DockerClient) UploadToContainerArgsForCall(i int) (string, docker.UploadToContainerOptions) {
	fake.uploadToContainerMutex.RLock()
	defer fake.uploadToContainerMutex.RUnlock()
	argsForCall := fake.uploadToContainerArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *DockerClient) UploadToContainerReturns(result1 error) {
	fake.uploadToContainerMutex.Lock()
	defer fake.uploadToContainerMutex.Unlock()
	fake.UploadToContainerStub = nil
	fake.uploadToContainerReturns = struct {
		result1 error
	}{result1}
}

func (fake *DockerClient) UploadToContainerReturnsOnCall(i int, result1 error) {
	fake.uploadToContainerMutex.Lock()
	defer fake.uploadToContainerMutex.Unlock()
	fake.UploadToContainerStub = nil
	if fake.uploadToContainerReturnsOnCall == nil {
		fake.uploadToContainerReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.uploadToContainerReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *DockerClient) WaitContainer(arg1 string) (int, error) {
	fake.waitContainerMutex.Lock()
	ret, specificReturn := fake.waitContainerReturnsOnCall[len(fake.waitContainerArgsForCall)]
	fake.waitContainerArgsForCall = append(fake.waitContainerArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("WaitContainer", []interface{}{arg1})
	fake.waitContainerMutex.Unlock()
	if fake.WaitContainerStub != nil {
		return fake.WaitContainerStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.waitContainerReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *DockerClient) WaitContainerCallCount() int {
	fake.waitContainerMutex.RLock()
	defer fake.waitContainerMutex.RUnlock()
	return len(fake.waitContainerArgsForCall)
}

func (fake *DockerClient) WaitContainerCalls(stub func(string) (int, error)) {
	fake.waitContainerMutex.Lock()
	defer fake.waitContainerMutex.Unlock()
	fake.WaitContainerStub = stub
}

func (fake *DockerClient) WaitContainerArgsForCall(i int) string {
	fake.waitContainerMutex.RLock()
	defer fake.waitContainerMutex.RUnlock()
	argsForCall := fake.waitContainerArgsForCall[i]
	return argsForCall.arg1
}

func (fake *DockerClient) WaitContainerReturns(result1 int, result2 error) {
	fake.waitContainerMutex.Lock()
	defer fake.waitContainerMutex.Unlock()
	fake.WaitContainerStub = nil
	fake.waitContainerReturns = struct {
		result1 int
		result2 error
	}{result1, result2}
}

func (fake *DockerClient) WaitContainerReturnsOnCall(i int, result1 int, result2 error) {
	fake.waitContainerMutex.Lock()
	defer fake.waitContainerMutex.Unlock()
	fake.WaitContainerStub = nil
	if fake.waitContainerReturnsOnCall == nil {
		fake.waitContainerReturnsOnCall = make(map[int]struct {
			result1 int
			result2 error
		})
	}
	fake.waitContainerReturnsOnCall[i] = struct {
		result1 int
		result2 error
	}{result1, result2}
}

func (fake *DockerClient) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.attachToContainerMutex.RLock()
	defer fake.attachToContainerMutex.RUnlock()
	fake.buildImageMutex.RLock()
	defer fake.buildImageMutex.RUnlock()
	fake.createContainerMutex.RLock()
	defer fake.createContainerMutex.RUnlock()
	fake.inspectImageMutex.RLock()
	defer fake.inspectImageMutex.RUnlock()
	fake.killContainerMutex.RLock()
	defer fake.killContainerMutex.RUnlock()
	fake.pingWithContextMutex.RLock()
	defer fake.pingWithContextMutex.RUnlock()
	fake.removeContainerMutex.RLock()
	defer fake.removeContainerMutex.RUnlock()
	fake.startContainerMutex.RLock()
	defer fake.startContainerMutex.RUnlock()
	fake.stopContainerMutex.RLock()
	defer fake.stopContainerMutex.RUnlock()
	fake.uploadToContainerMutex.RLock()
	defer fake.uploadToContainerMutex.RUnlock()
	fake.waitContainerMutex.RLock()
	defer fake.waitContainerMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *DockerClient) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}
