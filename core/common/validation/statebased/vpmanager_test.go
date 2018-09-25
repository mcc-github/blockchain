/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package statebased

import (
	"fmt"
	"math/rand"
	"runtime"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/mcc-github/blockchain/common/cauthdsl"
	"github.com/mcc-github/blockchain/core/handlers/validation/api/state"
	"github.com/mcc-github/blockchain/protos/ledger/rwset"
	"github.com/mcc-github/blockchain/protos/ledger/rwset/kvrwset"
	pb "github.com/mcc-github/blockchain/protos/peer"
	"github.com/mcc-github/blockchain/protos/utils"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

type mockState struct {
	GetStateMetadataRv              map[string][]byte
	GetStateMetadataErr             error
	GetPrivateDataMetadataByHashRv  map[string][]byte
	GetPrivateDataMetadataByHashErr error
	DoneCalled                      bool
}

func (ms *mockState) GetStateMultipleKeys(namespace string, keys []string) ([][]byte, error) {
	return nil, nil
}

func (ms *mockState) GetStateRangeScanIterator(namespace string, startKey string, endKey string) (validation.ResultsIterator, error) {
	return nil, nil
}

func (ms *mockState) GetStateMetadata(namespace, key string) (map[string][]byte, error) {
	return ms.GetStateMetadataRv, ms.GetStateMetadataErr
}

func (ms *mockState) GetPrivateDataMetadataByHash(namespace, collection string, keyhash []byte) (map[string][]byte, error) {
	return ms.GetPrivateDataMetadataByHashRv, ms.GetPrivateDataMetadataByHashErr
}

func (ms *mockState) Done() {
	ms.DoneCalled = true
}

type mockStateFetcher struct {
	mutex          sync.Mutex
	returnedStates []*mockState
	FetchStateRv   *mockState
	FetchStateErr  error
}

func (ms *mockStateFetcher) DoneCalled() bool {
	for _, s := range ms.returnedStates {
		if !s.DoneCalled {
			return false
		}
	}
	return true
}

func (ms *mockStateFetcher) FetchState() (validation.State, error) {
	var rv *mockState
	if ms.FetchStateRv != nil {
		rv = &mockState{
			GetPrivateDataMetadataByHashErr: ms.FetchStateRv.GetPrivateDataMetadataByHashErr,
			GetStateMetadataErr:             ms.FetchStateRv.GetStateMetadataErr,
			GetPrivateDataMetadataByHashRv:  ms.FetchStateRv.GetPrivateDataMetadataByHashRv,
			GetStateMetadataRv:              ms.FetchStateRv.GetStateMetadataRv,
		}
		ms.mutex.Lock()
		if ms.returnedStates != nil {
			ms.returnedStates = make([]*mockState, 0, 1)
		}
		ms.returnedStates = append(ms.returnedStates, rv)
		ms.mutex.Unlock()
	}
	return rv, ms.FetchStateErr
}

func TestSimple(t *testing.T) {
	t.Parallel()

	

	vpMetadataKey := pb.MetaDataKeys_VALIDATION_PARAMETER.String()
	spe := cauthdsl.SignedByMspMember("foo")
	mr := &mockState{GetStateMetadataRv: map[string][]byte{vpMetadataKey: utils.MarshalOrPanic(spe)}, GetPrivateDataMetadataByHashRv: map[string][]byte{vpMetadataKey: utils.MarshalOrPanic(spe)}}
	ms := &mockStateFetcher{FetchStateRv: mr}
	pm := &KeyLevelValidationParameterManagerImpl{
		StateFetcher: ms,
	}

	sp, err := pm.GetValidationParameterForKey("cc", "coll", "key", 0, 0)
	assert.NoError(t, err)
	assert.Equal(t, utils.MarshalOrPanic(spe), sp)
	assert.True(t, ms.DoneCalled())
}

func rwsetUpdatingMetadataFor(cc, key string) []byte {
	vpMetadataKey := pb.MetaDataKeys_VALIDATION_PARAMETER.String()
	return utils.MarshalOrPanic(
		&rwset.TxReadWriteSet{
			NsRwset: []*rwset.NsReadWriteSet{
				{
					Namespace: cc,
					Rwset: utils.MarshalOrPanic(&kvrwset.KVRWSet{
						MetadataWrites: []*kvrwset.KVMetadataWrite{
							{
								Key: key,
								Entries: []*kvrwset.KVMetadataEntry{
									{
										Name: vpMetadataKey,
									},
								},
							},
						},
					}),
				},
			}})
}

func pvtRwsetUpdatingMetadataFor(cc, coll, key string) []byte {
	vpMetadataKey := pb.MetaDataKeys_VALIDATION_PARAMETER.String()
	return utils.MarshalOrPanic(
		&rwset.TxReadWriteSet{
			NsRwset: []*rwset.NsReadWriteSet{
				{
					Namespace: cc,
					CollectionHashedRwset: []*rwset.CollectionHashedReadWriteSet{
						{
							CollectionName: coll,
							HashedRwset: utils.MarshalOrPanic(&kvrwset.HashedRWSet{
								MetadataWrites: []*kvrwset.KVMetadataWriteHash{
									{
										KeyHash: []byte(key),
										Entries: []*kvrwset.KVMetadataEntry{
											{
												Name: vpMetadataKey,
											},
										},
									},
								},
							}),
						},
					},
				},
			}})
}

func runFunctions(t *testing.T, seed int64, funcs ...func()) {
	r := rand.New(rand.NewSource(seed))
	c := make(chan struct{})
	for _, i := range r.Perm(len(funcs)) {
		iLcl := i
		go func() {
			assert.NotPanics(t, funcs[iLcl], "assert failure occurred with seed %d", seed)
			c <- struct{}{}
		}()
	}
	for range funcs {
		<-c
	}
}

func TestDependencyNoConflict(t *testing.T) {
	t.Parallel()
	seed := time.Now().Unix()

	
	
	
	
	
	

	vpMetadataKey := pb.MetaDataKeys_VALIDATION_PARAMETER.String()
	spe := cauthdsl.SignedByMspMember("foo")
	mr := &mockState{GetStateMetadataRv: map[string][]byte{vpMetadataKey: utils.MarshalOrPanic(spe)}, GetPrivateDataMetadataByHashRv: map[string][]byte{vpMetadataKey: utils.MarshalOrPanic(spe)}}
	ms := &mockStateFetcher{FetchStateRv: mr}
	pm := &KeyLevelValidationParameterManagerImpl{StateFetcher: ms}

	cc, coll, key := "cc", "", "key"

	rwsetbytes := rwsetUpdatingMetadataFor(cc, key)

	resC := make(chan []byte, 1)
	errC := make(chan error, 1)
	runFunctions(t, seed,
		func() {
			pm.ExtractValidationParameterDependency(1, 0, rwsetbytes)
		},
		func() {
			pm.SetTxValidationResult(cc, 1, 0, errors.New(""))
		},
		func() {
			sp, err := pm.GetValidationParameterForKey(cc, coll, key, 1, 1)
			resC <- sp
			errC <- err
		})

	sp := <-resC
	err := <-errC
	assert.NoError(t, err, "assert failure occurred with seed %d", seed)
	assert.Equal(t, utils.MarshalOrPanic(spe), sp, "assert failure occurred with seed %d", seed)
	assert.True(t, ms.DoneCalled(), "assert failure occurred with seed %d", seed)
}

func TestDependencyConflict(t *testing.T) {
	t.Parallel()
	seed := time.Now().Unix()

	
	
	
	
	
	
	

	vpMetadataKey := pb.MetaDataKeys_VALIDATION_PARAMETER.String()
	spe := cauthdsl.SignedByMspMember("foo")
	mr := &mockState{GetStateMetadataRv: map[string][]byte{vpMetadataKey: utils.MarshalOrPanic(spe)}, GetPrivateDataMetadataByHashRv: map[string][]byte{vpMetadataKey: utils.MarshalOrPanic(spe)}}
	ms := &mockStateFetcher{FetchStateRv: mr}
	pm := &KeyLevelValidationParameterManagerImpl{StateFetcher: ms}

	cc, coll, key := "cc", "", "key"

	rwsetbytes := rwsetUpdatingMetadataFor(cc, key)

	resC := make(chan []byte, 1)
	errC := make(chan error, 1)
	runFunctions(t, seed,
		func() {
			pm.ExtractValidationParameterDependency(1, 0, rwsetbytes)
		},
		func() {
			pm.SetTxValidationResult(cc, 1, 0, nil)
		},
		func() {
			sp, err := pm.GetValidationParameterForKey(cc, coll, key, 1, 1)
			resC <- sp
			errC <- err
		})

	sp := <-resC
	err := <-errC
	assert.Errorf(t, err, "assert failure occurred with seed %d", seed)
	assert.IsType(t, &ValidationParameterUpdatedError{}, err, "assert failure occurred with seed %d", seed)
	assert.Nil(t, sp, "assert failure occurred with seed %d", seed)
}

func TestMultipleDependencyNoConflict(t *testing.T) {
	t.Parallel()
	seed := time.Now().Unix()

	
	
	
	
	
	

	vpMetadataKey := pb.MetaDataKeys_VALIDATION_PARAMETER.String()
	spe := cauthdsl.SignedByMspMember("foo")
	mr := &mockState{GetStateMetadataRv: map[string][]byte{vpMetadataKey: utils.MarshalOrPanic(spe)}, GetPrivateDataMetadataByHashRv: map[string][]byte{vpMetadataKey: utils.MarshalOrPanic(spe)}}
	ms := &mockStateFetcher{FetchStateRv: mr}
	pm := &KeyLevelValidationParameterManagerImpl{StateFetcher: ms}

	cc, coll, key := "cc", "", "key"

	rwsetbytes := rwsetUpdatingMetadataFor(cc, key)

	resC := make(chan []byte, 1)
	errC := make(chan error, 1)
	runFunctions(t, seed,
		func() {
			pm.ExtractValidationParameterDependency(1, 0, rwsetbytes)
		},
		func() {
			pm.SetTxValidationResult(cc, 1, 0, errors.New(""))
		},
		func() {
			pm.ExtractValidationParameterDependency(1, 1, rwsetbytes)
		},
		func() {
			pm.SetTxValidationResult(cc, 1, 1, errors.New(""))
		},
		func() {
			sp, err := pm.GetValidationParameterForKey(cc, coll, key, 1, 2)
			resC <- sp
			errC <- err
		})

	sp := <-resC
	err := <-errC
	assert.NoError(t, err, "assert failure occurred with seed %d", seed)
	assert.Equal(t, utils.MarshalOrPanic(spe), sp, "assert failure occurred with seed %d", seed)
	assert.True(t, ms.DoneCalled(), "assert failure occurred with seed %d", seed)
}

func TestMultipleDependencyConflict(t *testing.T) {
	t.Parallel()
	seed := time.Now().Unix()

	
	
	
	
	
	
	

	vpMetadataKey := pb.MetaDataKeys_VALIDATION_PARAMETER.String()
	spe := cauthdsl.SignedByMspMember("foo")
	mr := &mockState{GetStateMetadataRv: map[string][]byte{vpMetadataKey: utils.MarshalOrPanic(spe)}, GetPrivateDataMetadataByHashRv: map[string][]byte{vpMetadataKey: utils.MarshalOrPanic(spe)}}
	ms := &mockStateFetcher{FetchStateRv: mr}
	pm := &KeyLevelValidationParameterManagerImpl{StateFetcher: ms}

	cc, coll, key := "cc", "", "key"

	rwsetbytes := rwsetUpdatingMetadataFor(cc, key)

	resC := make(chan []byte, 1)
	errC := make(chan error, 1)
	runFunctions(t, seed,
		func() {
			pm.ExtractValidationParameterDependency(1, 0, rwsetbytes)
		},
		func() {
			pm.SetTxValidationResult(cc, 1, 0, errors.New(""))
		},
		func() {
			pm.ExtractValidationParameterDependency(1, 1, rwsetbytes)
		},
		func() {
			pm.SetTxValidationResult(cc, 1, 1, nil)
		},
		func() {
			sp, err := pm.GetValidationParameterForKey(cc, coll, key, 1, 2)
			resC <- sp
			errC <- err
		})

	sp := <-resC
	err := <-errC
	assert.Errorf(t, err, "assert failure occurred with seed %d", seed)
	assert.IsType(t, &ValidationParameterUpdatedError{}, err, "assert failure occurred with seed %d", seed)
	assert.Nil(t, sp, "assert failure occurred with seed %d", seed)
}

func TestPvtDependencyNoConflict(t *testing.T) {
	t.Parallel()
	seed := time.Now().Unix()

	

	vpMetadataKey := pb.MetaDataKeys_VALIDATION_PARAMETER.String()
	spe := cauthdsl.SignedByMspMember("foo")
	mr := &mockState{GetStateMetadataRv: map[string][]byte{vpMetadataKey: utils.MarshalOrPanic(spe)}, GetPrivateDataMetadataByHashRv: map[string][]byte{vpMetadataKey: utils.MarshalOrPanic(spe)}}
	ms := &mockStateFetcher{FetchStateRv: mr}
	pm := &KeyLevelValidationParameterManagerImpl{StateFetcher: ms}

	cc, coll, key := "cc", "coll", "key"

	rwsetBytes := pvtRwsetUpdatingMetadataFor(cc, coll, key)

	resC := make(chan []byte, 1)
	errC := make(chan error, 1)
	runFunctions(t, seed,
		func() {
			pm.ExtractValidationParameterDependency(1, 0, rwsetBytes)
		},
		func() {
			pm.SetTxValidationResult(cc, 1, 0, errors.New(""))
		},
		func() {
			sp, err := pm.GetValidationParameterForKey(cc, coll, key, 1, 1)
			resC <- sp
			errC <- err
		})

	sp := <-resC
	err := <-errC
	assert.NoError(t, err, "assert failure occurred with seed %d", seed)
	assert.Equal(t, utils.MarshalOrPanic(spe), sp, "assert failure occurred with seed %d", seed)
	assert.True(t, ms.DoneCalled(), "assert failure occurred with seed %d", seed)
}

func TestPvtDependencyConflict(t *testing.T) {
	t.Parallel()
	seed := time.Now().Unix()

	

	vpMetadataKey := pb.MetaDataKeys_VALIDATION_PARAMETER.String()
	spe := cauthdsl.SignedByMspMember("foo")
	mr := &mockState{GetStateMetadataRv: map[string][]byte{vpMetadataKey: utils.MarshalOrPanic(spe)}, GetPrivateDataMetadataByHashRv: map[string][]byte{vpMetadataKey: utils.MarshalOrPanic(spe)}}
	ms := &mockStateFetcher{FetchStateRv: mr}
	pm := &KeyLevelValidationParameterManagerImpl{StateFetcher: ms}

	cc, coll, key := "cc", "coll", "key"

	rwsetBytes := pvtRwsetUpdatingMetadataFor(cc, coll, key)

	resC := make(chan []byte, 1)
	errC := make(chan error, 1)
	runFunctions(t, seed,
		func() {
			pm.ExtractValidationParameterDependency(1, 0, rwsetBytes)
		},
		func() {
			pm.SetTxValidationResult(cc, 1, 0, nil)
		},
		func() {
			sp, err := pm.GetValidationParameterForKey(cc, coll, key, 1, 1)
			resC <- sp
			errC <- err
		})

	sp := <-resC
	err := <-errC
	assert.Errorf(t, err, "assert failure occurred with seed %d", seed)
	assert.IsType(t, &ValidationParameterUpdatedError{}, err, "assert failure occurred with seed %d", seed)
	assert.True(t, len(err.Error()) > 0, "assert failure occurred with seed %d", seed)
	assert.Nil(t, sp, "assert failure occurred with seed %d", seed)
}

func TestBlockValidationTerminatesBeforeNewBlock(t *testing.T) {
	t.Parallel()

	
	
	

	vpMetadataKey := pb.MetaDataKeys_VALIDATION_PARAMETER.String()
	spe := cauthdsl.SignedByMspMember("foo")
	mr := &mockState{GetStateMetadataRv: map[string][]byte{vpMetadataKey: utils.MarshalOrPanic(spe)}, GetPrivateDataMetadataByHashRv: map[string][]byte{vpMetadataKey: utils.MarshalOrPanic(spe)}}
	ms := &mockStateFetcher{FetchStateRv: mr}
	pm := &KeyLevelValidationParameterManagerImpl{StateFetcher: ms}

	cc, coll, key := "cc", "coll", "key"

	rwsetBytes := pvtRwsetUpdatingMetadataFor(cc, coll, key)

	pm.ExtractValidationParameterDependency(2, 0, rwsetBytes)
	panickingFunc := func() {
		pm.ExtractValidationParameterDependency(1, 0, rwsetBytes)
	}
	assert.Panics(t, panickingFunc)
}

func TestLedgerErrors(t *testing.T) {
	t.Parallel()
	seed := time.Now().Unix()

	
	

	mr := &mockState{
		GetStateMetadataErr:             fmt.Errorf("Ledger error"),
		GetPrivateDataMetadataByHashErr: fmt.Errorf("Ledger error"),
	}
	ms := &mockStateFetcher{FetchStateRv: mr, FetchStateErr: fmt.Errorf("Ledger error")}
	pm := &KeyLevelValidationParameterManagerImpl{StateFetcher: ms}

	cc, coll, key := "cc", "", "key"

	rwsetbytes := rwsetUpdatingMetadataFor(cc, key)

	errC := make(chan error, 1)
	runFunctions(t, seed,
		func() {
			pm.ExtractValidationParameterDependency(1, 0, rwsetbytes)
		},
		func() {
			pm.SetTxValidationResult(cc, 1, 0, errors.New(""))
		},
		func() {
			_, err := pm.GetValidationParameterForKey(cc, coll, key, 1, 1)
			errC <- err
		})

	err := <-errC
	assert.Errorf(t, err, "assert failure occurred with seed %d", seed)

	ms.FetchStateErr = nil

	runFunctions(t, seed,
		func() {
			pm.ExtractValidationParameterDependency(2, 0, rwsetbytes)
		},
		func() {
			pm.SetTxValidationResult(cc, 2, 0, errors.New(""))
		},
		func() {
			_, err := pm.GetValidationParameterForKey(cc, coll, key, 2, 1)
			errC <- err
		})

	err = <-errC
	assert.Error(t, err)

	cc, coll, key = "cc", "coll", "key"

	rwsetbytes = pvtRwsetUpdatingMetadataFor(cc, coll, key)

	runFunctions(t, seed,
		func() {
			pm.ExtractValidationParameterDependency(3, 0, rwsetbytes)
		},
		func() {
			pm.SetTxValidationResult(cc, 3, 0, errors.New(""))
		},
		func() {
			_, err = pm.GetValidationParameterForKey(cc, coll, key, 3, 1)
			errC <- err
		})

	err = <-errC
	assert.Errorf(t, err, "assert failure occurred with seed %d", seed)
	assert.True(t, ms.DoneCalled(), "assert failure occurred with seed %d", seed)
}

func TestBadRwsetIsNoDependency(t *testing.T) {
	t.Parallel()
	seed := time.Now().Unix()

	
	
	

	vpMetadataKey := pb.MetaDataKeys_VALIDATION_PARAMETER.String()
	spe := cauthdsl.SignedByMspMember("foo")
	mr := &mockState{GetStateMetadataRv: map[string][]byte{vpMetadataKey: utils.MarshalOrPanic(spe)}, GetPrivateDataMetadataByHashRv: map[string][]byte{vpMetadataKey: utils.MarshalOrPanic(spe)}}
	ms := &mockStateFetcher{FetchStateRv: mr}
	pm := &KeyLevelValidationParameterManagerImpl{StateFetcher: ms}

	cc, coll, key := "cc", "", "key"

	resC := make(chan []byte, 1)
	errC := make(chan error, 1)
	runFunctions(t, seed,
		func() {
			pm.ExtractValidationParameterDependency(1, 0, []byte("barf"))
		},
		func() {
			pm.SetTxValidationResult(cc, 1, 0, errors.New(""))
		},
		func() {
			sp, err := pm.GetValidationParameterForKey(cc, coll, key, 1, 1)
			resC <- sp
			errC <- err
		})

	sp := <-resC
	err := <-errC
	assert.NoError(t, err, "assert failure occurred with seed %d", seed)
	assert.Equal(t, utils.MarshalOrPanic(spe), sp, "assert failure occurred with seed %d", seed)
	assert.True(t, ms.DoneCalled(), "assert failure occurred with seed %d", seed)
}

func TestWritesIntoDifferentNamespaces(t *testing.T) {
	t.Parallel()
	seed := time.Now().Unix()

	
	
	

	vpMetadataKey := pb.MetaDataKeys_VALIDATION_PARAMETER.String()
	spe := cauthdsl.SignedByMspMember("foo")
	mr := &mockState{GetStateMetadataRv: map[string][]byte{vpMetadataKey: utils.MarshalOrPanic(spe)}, GetPrivateDataMetadataByHashRv: map[string][]byte{vpMetadataKey: utils.MarshalOrPanic(spe)}}
	ms := &mockStateFetcher{FetchStateRv: mr}
	pm := &KeyLevelValidationParameterManagerImpl{StateFetcher: ms}

	cc, othercc, coll, key := "cc1", "cc", "", "key"

	rwsetbytes := rwsetUpdatingMetadataFor(cc, key)

	resC := make(chan []byte, 1)
	errC := make(chan error, 1)
	runFunctions(t, seed,
		func() {
			pm.SetTxValidationResult(cc, 1, 0, nil)
		},
		func() {
			pm.ExtractValidationParameterDependency(1, 0, rwsetbytes)
			sp, err := pm.GetValidationParameterForKey(othercc, coll, key, 1, 1)
			resC <- sp
			errC <- err
		})

	sp := <-resC
	err := <-errC
	assert.NoError(t, err, "assert failure occurred with seed %d", seed)
	assert.Equal(t, utils.MarshalOrPanic(spe), sp, "assert failure occurred with seed %d", seed)
	assert.True(t, ms.DoneCalled(), "assert failure occurred with seed %d", seed)
}

func TestCombinedCalls(t *testing.T) {
	t.Parallel()
	seed := time.Now().Unix()

	
	

	vpMetadataKey := pb.MetaDataKeys_VALIDATION_PARAMETER.String()
	spe := cauthdsl.SignedByMspMember("foo")
	mr := &mockState{GetStateMetadataRv: map[string][]byte{vpMetadataKey: utils.MarshalOrPanic(spe)}, GetPrivateDataMetadataByHashRv: map[string][]byte{vpMetadataKey: utils.MarshalOrPanic(spe)}}
	ms := &mockStateFetcher{FetchStateRv: mr}
	pm := &KeyLevelValidationParameterManagerImpl{StateFetcher: ms}

	cc := "cc"
	coll := ""
	key1 := "key1"
	key2 := "key2"

	res1C := make(chan []byte, 1)
	err1C := make(chan error, 1)
	res2C := make(chan []byte, 1)
	err2C := make(chan error, 1)
	runFunctions(t, seed,
		func() {
			pm.ExtractValidationParameterDependency(1, 0, rwsetUpdatingMetadataFor(cc, key1))
		},
		func() {
			pm.SetTxValidationResult(cc, 1, 0, errors.New(""))
		},
		func() {
			pm.ExtractValidationParameterDependency(1, 1, rwsetUpdatingMetadataFor(cc, key2))
		},
		func() {
			pm.SetTxValidationResult(cc, 1, 1, nil)
		},
		func() {
			sp, err := pm.GetValidationParameterForKey(cc, coll, key1, 1, 2)
			res1C <- sp
			err1C <- err
		},
		func() {
			sp, err := pm.GetValidationParameterForKey(cc, coll, key2, 1, 2)
			res2C <- sp
			err2C <- err
		})

	sp := <-res1C
	err := <-err1C
	assert.NoError(t, err, "assert failure occurred with seed %d", seed)
	assert.Equal(t, utils.MarshalOrPanic(spe), sp, "assert failure occurred with seed %d", seed)

	sp = <-res2C
	err = <-err2C
	assert.Errorf(t, err, "assert failure occurred with seed %d", seed)
	assert.IsType(t, &ValidationParameterUpdatedError{}, err, "assert failure occurred with seed %d", seed)
	assert.Nil(t, sp, "assert failure occurred with seed %d", seed)

	assert.True(t, ms.DoneCalled(), "assert failure occurred with seed %d", seed)
}

func TestForRaces(t *testing.T) {
	seed := time.Now().Unix()

	
	
	

	vpMetadataKey := pb.MetaDataKeys_VALIDATION_PARAMETER.String()
	spe := cauthdsl.SignedByMspMember("foo")
	mr := &mockState{GetStateMetadataRv: map[string][]byte{vpMetadataKey: utils.MarshalOrPanic(spe)}, GetPrivateDataMetadataByHashRv: map[string][]byte{vpMetadataKey: utils.MarshalOrPanic(spe)}}
	ms := &mockStateFetcher{FetchStateRv: mr}
	pm := &KeyLevelValidationParameterManagerImpl{StateFetcher: ms}

	cc := "cc"
	coll := ""

	nRoutines := 1000
	funcArray := make([]func(), nRoutines)
	for i := 0; i < nRoutines; i++ {
		txnum := i
		funcArray[i] = func() {
			key := strconv.Itoa(txnum)
			pm.ExtractValidationParameterDependency(1, uint64(txnum), rwsetUpdatingMetadataFor(cc, key))

			
			runtime.Gosched()

			pm.SetTxValidationResult(cc, 1, uint64(txnum), errors.New(""))

			
			runtime.Gosched()

			sp, err := pm.GetValidationParameterForKey(cc, coll, key, 1, 2)
			assert.Equal(t, utils.MarshalOrPanic(spe), sp)
			assert.NoError(t, err)
		}
	}

	runFunctions(t, seed, funcArray...)

	assert.True(t, ms.DoneCalled(), "assert failure occurred with seed %d", seed)
}
