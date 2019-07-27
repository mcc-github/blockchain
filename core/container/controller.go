/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package container

import (
	"io"
	"sync"

	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/core/container/ccintf"
)

type VMProvider interface {
	NewVM() VM
}

type Builder interface {
	Build() (io.Reader, error)
}


type VM interface {
	Start(ccid ccintf.CCID, args []string, env []string, filesToUpload map[string][]byte) error
	Stop(ccid ccintf.CCID, timeout uint, dontkill bool, dontremove bool) error
	Wait(ccid ccintf.CCID) (int, error)
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
	Args          []string
	Env           []string
	FilesToUpload map[string][]byte
}

func (si StartContainerReq) Do(v VM) error {
	return v.Start(si.CCID, si.Args, si.Env, si.FilesToUpload)
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


type BuildReq struct {
	CCID        ccintf.CCID
	Type        string
	Path        string
	Name        string
	Version     string
	CodePackage []byte
}


type ChaincodeBuilder interface {
	Build(ccid ccintf.CCID, ccType, path, name, version string, codePackage []byte) error
}

func (b BuildReq) Do(v VM) error {
	if chaincodeBuilder, ok := v.(ChaincodeBuilder); ok {
		return chaincodeBuilder.Build(b.CCID, b.Type, b.Path, b.Name, b.Version, b.CodePackage)
	}
	return nil
}

func (w BuildReq) GetCCID() ccintf.CCID {
	return w.CCID
}

func (vmc *VMController) Process(vmtype string, req VMCReq) error {
	ccid := req.GetCCID()

	v, ok := vmc.vmProviders[vmtype]
	if !ok {
		vmLogger.Panicf("Programming error: unsupported VM type: %s for ccid='%s'", vmtype, ccid.String())
	}

	vmc.lockContainer(ccid)
	defer vmc.unlockContainer(ccid)
	return req.Do(v.NewVM())
}
