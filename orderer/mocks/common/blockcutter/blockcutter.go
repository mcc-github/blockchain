/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

                 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package blockcutter

import (
	"sync"

	"github.com/mcc-github/blockchain/common/flogging"
	cb "github.com/mcc-github/blockchain/protos/common"
)

const pkgLogID = "orderer/mocks/common/blockcutter"

var logger = flogging.MustGetLogger(pkgLogID)


type Receiver struct {
	
	IsolatedTx bool

	
	CutAncestors bool

	
	CutNext bool

	
	SkipAppendCurBatch bool

	
	mutex sync.Mutex

	
	curBatch []*cb.Envelope

	
	
	Block chan struct{}
}


func NewReceiver() *Receiver {
	return &Receiver{
		IsolatedTx:   false,
		CutAncestors: false,
		CutNext:      false,
		Block:        make(chan struct{}),
	}
}


func (mbc *Receiver) Ordered(env *cb.Envelope) ([][]*cb.Envelope, bool) {
	defer func() {
		<-mbc.Block
	}()

	mbc.mutex.Lock()
	defer mbc.mutex.Unlock()

	if mbc.IsolatedTx {
		logger.Debugf("Receiver: Returning dual batch")
		res := [][]*cb.Envelope{mbc.curBatch, {env}}
		mbc.curBatch = nil
		return res, false
	}

	if mbc.CutAncestors {
		logger.Debugf("Receiver: Returning current batch and appending newest env")
		res := [][]*cb.Envelope{mbc.curBatch}
		mbc.curBatch = []*cb.Envelope{env}
		return res, true
	}

	if !mbc.SkipAppendCurBatch {
		mbc.curBatch = append(mbc.curBatch, env)
	}

	if mbc.CutNext {
		logger.Debugf("Receiver: Returning regular batch")
		res := [][]*cb.Envelope{mbc.curBatch}
		mbc.curBatch = nil
		return res, false
	}

	logger.Debugf("Appending to batch")
	return nil, true
}


func (mbc *Receiver) Cut() []*cb.Envelope {
	mbc.mutex.Lock()
	defer mbc.mutex.Unlock()
	logger.Debugf("Cutting batch")
	res := mbc.curBatch
	mbc.curBatch = nil
	return res
}

func (mbc *Receiver) CurBatch() []*cb.Envelope {
	mbc.mutex.Lock()
	defer mbc.mutex.Unlock()
	return mbc.curBatch
}
