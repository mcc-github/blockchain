/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package inproccontroller

import (
	"fmt"

	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/core/chaincode/shim"
	"github.com/mcc-github/blockchain/core/container"
	"github.com/mcc-github/blockchain/core/container/ccintf"
	pb "github.com/mcc-github/blockchain/protos/peer"

	"golang.org/x/net/context"
)



const ContainerType = "SYSTEM"

type inprocContainer struct {
	chaincode shim.Chaincode
	running   bool
	args      []string
	env       []string
	stopChan  chan struct{}
}

var (
	inprocLogger = flogging.MustGetLogger("inproccontroller")

	
	
	_shimStartInProc    = shim.StartInProc
	_inprocLoggerErrorf = inprocLogger.Errorf
)




type SysCCRegisteredErr string

func (s SysCCRegisteredErr) Error() string {
	return fmt.Sprintf("%s already registered", string(s))
}



type Registry struct {
	typeRegistry map[string]*inprocContainer
	instRegistry map[string]*inprocContainer
}


func NewRegistry() *Registry {
	return &Registry{
		typeRegistry: make(map[string]*inprocContainer),
		instRegistry: make(map[string]*inprocContainer),
	}
}


func (r *Registry) NewVM() container.VM {
	return NewInprocVM(r)
}


func (r *Registry) Register(ccid *ccintf.CCID, cc shim.Chaincode) error {
	name := ccid.GetName()
	inprocLogger.Debugf("Registering chaincode instance: %s", name)
	tmp := r.typeRegistry[name]
	if tmp != nil {
		return SysCCRegisteredErr(name)
	}

	r.typeRegistry[name] = &inprocContainer{chaincode: cc}
	return nil
}


type InprocVM struct {
	id       string
	registry *Registry
}


func NewInprocVM(r *Registry) *InprocVM {
	return &InprocVM{
		registry: r,
	}
}

func (vm *InprocVM) getInstance(ctxt context.Context, ipctemplate *inprocContainer, instName string, args []string, env []string) (*inprocContainer, error) {
	ipc := vm.registry.instRegistry[instName]
	if ipc != nil {
		inprocLogger.Warningf("chaincode instance exists for %s", instName)
		return ipc, nil
	}
	ipc = &inprocContainer{args: args, env: env, chaincode: ipctemplate.chaincode, stopChan: make(chan struct{})}
	vm.registry.instRegistry[instName] = ipc
	inprocLogger.Debugf("chaincode instance created for %s", instName)
	return ipc, nil
}

func (ipc *inprocContainer) launchInProc(ctxt context.Context, id string, args []string, env []string, ccSupport ccintf.CCSupport) error {
	peerRcvCCSend := make(chan *pb.ChaincodeMessage)
	ccRcvPeerSend := make(chan *pb.ChaincodeMessage)
	var err error
	ccchan := make(chan struct{}, 1)
	ccsupportchan := make(chan struct{}, 1)
	shimStartInProc := _shimStartInProc 
	go func() {
		defer close(ccchan)
		inprocLogger.Debugf("chaincode started for %s", id)
		if args == nil {
			args = ipc.args
		}
		if env == nil {
			env = ipc.env
		}
		err := shimStartInProc(env, args, ipc.chaincode, ccRcvPeerSend, peerRcvCCSend)
		if err != nil {
			err = fmt.Errorf("chaincode-support ended with err: %s", err)
			_inprocLoggerErrorf("%s", err)
		}
		inprocLogger.Debugf("chaincode ended for %s with err: %s", id, err)
	}()

	go func() {
		defer close(ccsupportchan)
		inprocStream := newInProcStream(peerRcvCCSend, ccRcvPeerSend)
		inprocLogger.Debugf("chaincode-support started for  %s", id)
		err := ccSupport.HandleChaincodeStream(ctxt, inprocStream)
		if err != nil {
			err = fmt.Errorf("chaincode ended with err: %s", err)
			_inprocLoggerErrorf("%s", err)
		}
		inprocLogger.Debugf("chaincode-support ended for %s with err: %s", id, err)
	}()

	select {
	case <-ccchan:
		close(peerRcvCCSend)
		inprocLogger.Debugf("chaincode %s quit", id)
	case <-ccsupportchan:
		close(ccRcvPeerSend)
		inprocLogger.Debugf("chaincode support %s quit", id)
	case <-ipc.stopChan:
		close(ccRcvPeerSend)
		close(peerRcvCCSend)
		inprocLogger.Debugf("chaincode %s stopped", id)
	}
	return err
}


func (vm *InprocVM) Start(ctxt context.Context, ccid ccintf.CCID, args []string, env []string, filesToUpload map[string][]byte, builder container.Builder) error {
	path := ccid.GetName()

	ipctemplate := vm.registry.typeRegistry[path]

	if ipctemplate == nil {
		return fmt.Errorf(fmt.Sprintf("%s not registered", path))
	}

	instName := vm.GetVMName(ccid)

	ipc, err := vm.getInstance(ctxt, ipctemplate, instName, args, env)

	if err != nil {
		return fmt.Errorf(fmt.Sprintf("could not create instance for %s", instName))
	}

	if ipc.running {
		return fmt.Errorf(fmt.Sprintf("chaincode running %s", path))
	}

	

	ccSupport, ok := ctxt.Value(ccintf.GetCCHandlerKey()).(ccintf.CCSupport)
	if !ok || ccSupport == nil {
		return fmt.Errorf("in-process communication generator not supplied")
	}

	ipc.running = true

	go func() {
		defer func() {
			if r := recover(); r != nil {
				inprocLogger.Criticalf("caught panic from chaincode  %s", instName)
			}
		}()
		ipc.launchInProc(ctxt, instName, args, env, ccSupport)
	}()

	return nil
}


func (vm *InprocVM) Stop(ctxt context.Context, ccid ccintf.CCID, timeout uint, dontkill bool, dontremove bool) error {
	path := ccid.GetName()

	ipctemplate := vm.registry.typeRegistry[path]
	if ipctemplate == nil {
		return fmt.Errorf("%s not registered", path)
	}

	instName := vm.GetVMName(ccid)

	ipc := vm.registry.instRegistry[instName]

	if ipc == nil {
		return fmt.Errorf("%s not found", instName)
	}

	if !ipc.running {
		return fmt.Errorf("%s not running", instName)
	}

	ipc.stopChan <- struct{}{}

	delete(vm.registry.instRegistry, instName)
	
	return nil
}




func (vm *InprocVM) GetVMName(ccid ccintf.CCID) string {
	return ccid.GetName()
}
