/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package container

import (
	"io"
	"sync"

	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/core/chaincode/persistence"
	"github.com/mcc-github/blockchain/core/container/ccintf"

	"github.com/pkg/errors"
)

var vmLogger = flogging.MustGetLogger("container")




type VM interface {
	Build(ccid string, metadata *persistence.ChaincodePackageMetadata, codePackageStream io.Reader) (Instance, error)
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




type PackageProvider interface {
	GetChaincodePackage(packageID string) (*persistence.ChaincodePackageMetadata, io.ReadCloser, error)
}

type Router struct {
	ExternalVM      VM
	DockerVM        VM
	containers      map[string]Instance
	PackageProvider PackageProvider
	mutex           sync.Mutex
}

func (r *Router) getInstance(ccid string) Instance {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	
	

	if r.containers == nil {
		r.containers = map[string]Instance{}
	}

	vm, ok := r.containers[ccid]
	if !ok {
		return UninitializedInstance{}
	}

	return vm
}

func (r *Router) Build(ccid string) error {
	
	
	packageID := ccid

	var instance Instance

	var externalErr error
	if r.ExternalVM != nil {
		metadata, codeStream, err := r.PackageProvider.GetChaincodePackage(packageID)
		if err != nil {
			return errors.WithMessage(err, "get chaincode package for external build failed")
		}
		instance, externalErr = r.ExternalVM.Build(ccid, metadata, codeStream)
		codeStream.Close()
	}

	var dockerErr error
	if r.ExternalVM == nil || externalErr != nil {
		metadata, codeStream, err := r.PackageProvider.GetChaincodePackage(ccid)
		if err != nil {
			return errors.WithMessage(err, "get chaincode package for docker build failed")
		}
		instance, dockerErr = r.DockerVM.Build(ccid, metadata, codeStream)
		codeStream.Close()
	}

	if dockerErr != nil {
		return errors.WithMessagef(dockerErr, "failed external (%s) and docker build", externalErr)
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.containers == nil {
		r.containers = map[string]Instance{}
	}

	r.containers[ccid] = instance

	return nil
}

func (r *Router) Start(ccid string, peerConnection *ccintf.PeerConnection) error {
	return r.getInstance(ccid).Start(peerConnection)
}

func (r *Router) Stop(ccid string) error {
	return r.getInstance(ccid).Stop()
}

func (r *Router) Wait(ccid string) (int, error) {
	return r.getInstance(ccid).Wait()
}
