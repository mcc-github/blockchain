/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package container

import (
	"sync"

	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/core/common/ccprovider"
	"github.com/mcc-github/blockchain/core/container/ccintf"

	"github.com/pkg/errors"
)

var vmLogger = flogging.MustGetLogger("container")




type VM interface {
	Build(ccci *ccprovider.ChaincodeContainerInfo, codePackage []byte) (Instance, error)
}






type Instance interface {
	Start(peerConnection *ccintf.PeerConnection) error
	Stop() error
	Wait() (int, error)
}

type UninitializedInstance struct{}

func (UninitializedInstance) Start(peerConnection *ccintf.PeerConnection) error {
	return errors.Errorf("instance has not yet been built, cannot be started")
}

func (UninitializedInstance) Stop() error {
	return errors.Errorf("instance has not yet been built, cannot be stopped")
}

func (UninitializedInstance) Wait() (int, error) {
	return 0, errors.Errorf("instance has not yet been built, cannot wait")
}

type Router struct {
	DockerVM   VM
	containers map[ccintf.CCID]Instance
	mutex      sync.Mutex
}

func (r *Router) getInstance(ccid ccintf.CCID) Instance {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	
	

	if r.containers == nil {
		r.containers = map[ccintf.CCID]Instance{}
	}

	vm, ok := r.containers[ccid]
	if !ok {
		return UninitializedInstance{}
	}

	return vm
}

func (r *Router) Build(ccci *ccprovider.ChaincodeContainerInfo, codePackage []byte) error {
	
	instance, err := r.DockerVM.Build(ccci, codePackage)
	if err != nil {
		return errors.WithMessage(err, "failed docker build")
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.containers == nil {
		r.containers = map[ccintf.CCID]Instance{}
	}

	r.containers[ccintf.CCID(ccci.PackageID)] = instance

	return nil
}

func (r *Router) Start(ccid ccintf.CCID, peerConnection *ccintf.PeerConnection) error {
	return r.getInstance(ccid).Start(peerConnection)
}

func (r *Router) Stop(ccid ccintf.CCID) error {
	return r.getInstance(ccid).Stop()
}

func (r *Router) Wait(ccid ccintf.CCID) (int, error) {
	return r.getInstance(ccid).Wait()
}
