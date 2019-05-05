/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package container

import (
	"context"
	"fmt"
	"io"
	"sync"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/core/chaincode/platforms"
	"github.com/mcc-github/blockchain/core/container/ccintf"
	pb "github.com/mcc-github/blockchain/protos/peer"
)

type VMProvider interface {
	NewVM() VM
}

type Builder interface {
	Build() (io.Reader, error)
}


type VM interface {
	Start(ccid ccintf.CCID, args []string, env []string, filesToUpload map[string][]byte, builder Builder) error
	Stop(ccid ccintf.CCID, timeout uint, dontkill bool, dontremove bool) error
	Wait(ccid ccintf.CCID) (int, error)
	HealthCheck(context.Context) error
}

type refCountedLock struct {
	refCount int
	lock     *sync.RWMutex
}





type VMController struct {
	sync.RWMutex
	containerLocks map[ccintf.CCID]*refCountedLock
	vmProviders    map[string]VMProvider
}

var vmLogger = flogging.MustGetLogger("container")


func NewVMController(vmProviders map[string]VMProvider) *VMController {
	return &VMController{
		containerLocks: make(map[ccintf.CCID]*refCountedLock),
		vmProviders:    vmProviders,
	}
}

func (vmc *VMController) newVM(typ string) VM {
	v, ok := vmc.vmProviders[typ]
	if !ok {
		vmLogger.Panicf("Programming error: unsupported VM type: %s", typ)
	}
	return v.NewVM()
}

func (vmc *VMController) lockContainer(id ccintf.CCID) {
	
	vmc.Lock()
	var refLck *refCountedLock
	var ok bool
	if refLck, ok = vmc.containerLocks[id]; !ok {
		refLck = &refCountedLock{refCount: 1, lock: &sync.RWMutex{}}
		vmc.containerLocks[id] = refLck
	} else {
		refLck.refCount++
		vmLogger.Debugf("refcount %d (%s)", refLck.refCount, id)
	}
	vmc.Unlock()
	vmLogger.Debugf("waiting for container(%s) lock", id)
	refLck.lock.Lock()
	vmLogger.Debugf("got container (%s) lock", id)
}

func (vmc *VMController) unlockContainer(id ccintf.CCID) {
	vmc.Lock()
	if refLck, ok := vmc.containerLocks[id]; ok {
		if refLck.refCount <= 0 {
			panic("refcnt <= 0")
		}
		refLck.lock.Unlock()
		if refLck.refCount--; refLck.refCount == 0 {
			vmLogger.Debugf("container lock deleted(%s)", id)
			delete(vmc.containerLocks, id)
		}
	} else {
		vmLogger.Debugf("no lock to unlock(%s)!!", id)
	}
	vmc.Unlock()
}





type VMCReq interface {
	Do(v VM) error
	GetCCID() ccintf.CCID
}


type StartContainerReq struct {
	ccintf.CCID
	Builder       Builder
	Args          []string
	Env           []string
	FilesToUpload map[string][]byte
}








type PlatformBuilder struct {
	Type             string
	Path             string
	Name             string
	Version          string
	CodePackage      []byte
	PlatformRegistry *platforms.Registry
	Client           *docker.Client
}


func (b *PlatformBuilder) Build() (io.Reader, error) {
	return b.PlatformRegistry.GenerateDockerBuild(
		b.Type,
		b.Path,
		b.Name,
		b.Version,
		b.CodePackage,
		b.Client,
	)
}

func (si StartContainerReq) Do(v VM) error {
	return v.Start(si.CCID, si.Args, si.Env, si.FilesToUpload, si.Builder)
}

func (si StartContainerReq) GetCCID() ccintf.CCID {
	return si.CCID
}


type StopContainerReq struct {
	ccintf.CCID
	Timeout uint
	
	Dontkill bool
	
	Dontremove bool
}

func (si StopContainerReq) Do(v VM) error {
	return v.Stop(si.CCID, si.Timeout, si.Dontkill, si.Dontremove)
}

func (si StopContainerReq) GetCCID() ccintf.CCID {
	return si.CCID
}




type ExitedFunc func(exitCode int, err error)



type WaitContainerReq struct {
	CCID   ccintf.CCID
	Exited ExitedFunc
}

func (w WaitContainerReq) Do(v VM) error {
	exited := w.Exited
	go func() {
		exitCode, err := v.Wait(w.CCID)
		exited(exitCode, err)
	}()
	return nil
}

func (w WaitContainerReq) GetCCID() ccintf.CCID {
	return w.CCID
}

func (vmc *VMController) Process(vmtype string, req VMCReq) error {
	v := vmc.newVM(vmtype)
	ccid := req.GetCCID()

	vmc.lockContainer(ccid)
	defer vmc.unlockContainer(ccid)
	return req.Do(v)
}


func GetChaincodePackageBytes(pr *platforms.Registry, spec *pb.ChaincodeSpec) ([]byte, error) {
	if spec == nil || spec.ChaincodeId == nil {
		return nil, fmt.Errorf("invalid chaincode spec")
	}

	return pr.GetDeploymentPayload(spec.Type.String(), spec.ChaincodeId.Path)
}
