/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chaincode

import (
	"sync"

	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/pkg/errors"
)


type HandlerRegistry struct {
	allowUnsolicitedRegistration bool 

	mutex     sync.Mutex              
	handlers  map[string]*Handler     
	launching map[string]*LaunchState 
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
		handlers:                     map[string]*Handler{},
		launching:                    map[string]*LaunchState{},
		allowUnsolicitedRegistration: allowUnsolicitedRegistration,
	}
}





func (r *HandlerRegistry) Launching(ccid string) (*LaunchState, bool) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	
	if launchState, ok := r.launching[ccid]; ok {
		return launchState, true
	}

	
	if _, ok := r.handlers[ccid]; ok {
		launchState := NewLaunchState()
		launchState.Notify(nil)
		return launchState, true
	}

	
	launchState := NewLaunchState()
	r.launching[ccid] = launchState
	return launchState, false
}



func (r *HandlerRegistry) Ready(ccid string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	launchStatus := r.launching[ccid]
	if launchStatus != nil {
		launchStatus.Notify(nil)
	}
}


func (r *HandlerRegistry) Failed(ccid string, err error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	launchStatus := r.launching[ccid]
	if launchStatus != nil {
		launchStatus.Notify(err)
	}
}


func (r *HandlerRegistry) Handler(ccid string) *Handler {
	r.mutex.Lock()
	h := r.handlers[ccid]
	r.mutex.Unlock()
	return h
}





func (r *HandlerRegistry) Register(h *Handler) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.handlers[h.chaincodeID] != nil {
		chaincodeLogger.Debugf("duplicate registered handler(key:%s) return error", h.chaincodeID)
		return errors.Errorf("duplicate chaincodeID: %s", h.chaincodeID)
	}

	
	
	if r.launching[h.chaincodeID] == nil && !r.allowUnsolicitedRegistration {
		return errors.Errorf("peer will not accept external chaincode connection %s (except in dev mode)", h.chaincodeID)
	}

	r.handlers[h.chaincodeID] = h

	chaincodeLogger.Debugf("registered handler complete for chaincode %s", h.chaincodeID)
	return nil
}




func (r *HandlerRegistry) Deregister(ccid string) error {
	chaincodeLogger.Debugf("deregister handler: %s", ccid)

	r.mutex.Lock()
	handler := r.handlers[ccid]
	delete(r.handlers, ccid)
	delete(r.launching, ccid)
	r.mutex.Unlock()

	if handler == nil {
		return errors.Errorf("could not find handler: %s", ccid)
	}

	handler.Close()

	chaincodeLogger.Debugf("deregistered handler with key: %s", ccid)
	return nil
}

type TxQueryExecutorGetter struct {
	HandlerRegistry *HandlerRegistry
	CCID            string
}

func (g *TxQueryExecutorGetter) TxQueryExecutor(chainID, txID string) ledger.SimpleQueryExecutor {
	handler := g.HandlerRegistry.Handler(g.CCID)
	if handler == nil {
		return nil
	}
	txContext := handler.TXContexts.Get(chainID, txID)
	if txContext == nil {
		return nil
	}
	return txContext.TXSimulator
}
