/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chaincode

import (
	"sync"

	"github.com/mcc-github/blockchain/core/container/ccintf"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/pkg/errors"
)


type HandlerRegistry struct {
	allowUnsolicitedRegistration bool 

	mutex     sync.Mutex                   
	handlers  map[ccintf.CCID]*Handler     
	launching map[ccintf.CCID]*LaunchState 
}

type LaunchState struct {
	mutex    sync.Mutex
	notified bool
	done     chan struct{}
	err      error
}

func NewLaunchState() *LaunchState {
	return &LaunchState{
		done: make(chan struct{}),
	}
}

func (l *LaunchState) Done() <-chan struct{} {
	return l.done
}

func (l *LaunchState) Err() error {
	l.mutex.Lock()
	err := l.err
	l.mutex.Unlock()
	return err
}

func (l *LaunchState) Notify(err error) {
	l.mutex.Lock()
	if !l.notified {
		l.notified = true
		l.err = err
		close(l.done)
	}
	l.mutex.Unlock()
}


func NewHandlerRegistry(allowUnsolicitedRegistration bool) *HandlerRegistry {
	return &HandlerRegistry{
		handlers:                     map[ccintf.CCID]*Handler{},
		launching:                    map[ccintf.CCID]*LaunchState{},
		allowUnsolicitedRegistration: allowUnsolicitedRegistration,
	}
}





func (r *HandlerRegistry) Launching(packageID ccintf.CCID) (*LaunchState, bool) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	
	if launchState, ok := r.launching[packageID]; ok {
		return launchState, true
	}

	
	if _, ok := r.handlers[packageID]; ok {
		launchState := NewLaunchState()
		launchState.Notify(nil)
		return launchState, true
	}

	
	launchState := NewLaunchState()
	r.launching[packageID] = launchState
	return launchState, false
}



func (r *HandlerRegistry) Ready(packageID ccintf.CCID) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	launchStatus := r.launching[packageID]
	if launchStatus != nil {
		launchStatus.Notify(nil)
	}
}


func (r *HandlerRegistry) Failed(packageID ccintf.CCID, err error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	launchStatus := r.launching[packageID]
	if launchStatus != nil {
		launchStatus.Notify(err)
	}
}


func (r *HandlerRegistry) Handler(packageID ccintf.CCID) *Handler {
	r.mutex.Lock()
	h := r.handlers[packageID]
	r.mutex.Unlock()
	return h
}





func (r *HandlerRegistry) Register(h *Handler) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	
	
	
	
	
	
	packageID := ccintf.CCID(h.chaincodeID.Name)

	if r.handlers[packageID] != nil {
		chaincodeLogger.Debugf("duplicate registered handler(key:%s) return error", packageID)
		return errors.Errorf("duplicate chaincodeID: %s", h.chaincodeID.Name)
	}

	
	
	if r.launching[packageID] == nil && !r.allowUnsolicitedRegistration {
		return errors.Errorf("peer will not accept external chaincode connection %v (except in dev mode)", h.chaincodeID.Name)
	}

	r.handlers[packageID] = h

	chaincodeLogger.Debugf("registered handler complete for chaincode %s", packageID)
	return nil
}




func (r *HandlerRegistry) Deregister(packageID ccintf.CCID) error {
	chaincodeLogger.Debugf("deregister handler: %s", packageID)

	r.mutex.Lock()
	handler := r.handlers[packageID]
	delete(r.handlers, packageID)
	delete(r.launching, packageID)
	r.mutex.Unlock()

	if handler == nil {
		return errors.Errorf("could not find handler: %s", packageID)
	}

	handler.Close()

	chaincodeLogger.Debugf("deregistered handler with key: %s", packageID)
	return nil
}

type TxQueryExecutorGetter struct {
	HandlerRegistry *HandlerRegistry
	PackageID       ccintf.CCID
}

func (g *TxQueryExecutorGetter) TxQueryExecutor(chainID, txID string) ledger.SimpleQueryExecutor {
	handler := g.HandlerRegistry.Handler(g.PackageID)
	if handler == nil {
		return nil
	}
	txContext := handler.TXContexts.Get(chainID, txID)
	if txContext == nil {
		return nil
	}
	return txContext.TXSimulator
}
