/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package discovery

import (
	"encoding/asn1"
	"errors"
	"sync"
	"testing"

	"github.com/mcc-github/blockchain/protoutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSignedDataToKey(t *testing.T) {
	key1, err1 := signedDataToKey(protoutil.SignedData{
		Data:      []byte{1, 2, 3, 4},
		Identity:  []byte{5, 6, 7},
		Signature: []byte{8, 9},
	})
	key2, err2 := signedDataToKey(protoutil.SignedData{
		Data:      []byte{1, 2, 3},
		Identity:  []byte{4, 5, 6},
		Signature: []byte{7, 8, 9},
	})
	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NotEqual(t, key1, key2)
}

type mockAcSupport struct {
	mock.Mock
}

func (as *mockAcSupport) EligibleForService(channel string, data protoutil.SignedData) error {
	return as.Called(channel, data).Error(0)
}

func (as *mockAcSupport) ConfigSequence(channel string) uint64 {
	return as.Called(channel).Get(0).(uint64)
}

func TestCacheDisabled(t *testing.T) {
	sd := protoutil.SignedData{
		Data:      []byte{1, 2, 3},
		Identity:  []byte("authorizedIdentity"),
		Signature: []byte{1, 2, 3},
	}

	as := &mockAcSupport{}
	as.On("ConfigSequence", "foo").Return(uint64(0))
	as.On("EligibleForService", "foo", sd).Return(nil)
	cache := newAuthCache(as, authCacheConfig{maxCacheSize: 100, purgeRetentionRatio: 0.5})

	
	cache.EligibleForService("foo", sd)
	cache.EligibleForService("foo", sd)
	as.AssertNumberOfCalls(t, "EligibleForService", 2)
}

func TestCacheUsage(t *testing.T) {
	as := &mockAcSupport{}
	as.On("ConfigSequence", "foo").Return(uint64(0))
	as.On("ConfigSequence", "bar").Return(uint64(0))
	cache := newAuthCache(as, defaultConfig())

	sd1 := protoutil.SignedData{
		Data:      []byte{1, 2, 3},
		Identity:  []byte("authorizedIdentity"),
		Signature: []byte{1, 2, 3},
	}

	sd2 := protoutil.SignedData{
		Data:      []byte{1, 2, 3},
		Identity:  []byte("authorizedIdentity"),
		Signature: []byte{1, 2, 3},
	}

	sd3 := protoutil.SignedData{
		Data:      []byte{1, 2, 3, 3},
		Identity:  []byte("unAuthorizedIdentity"),
		Signature: []byte{1, 2, 3},
	}

	testCases := []struct {
		channel     string
		expectedErr error
		sd          protoutil.SignedData
	}{
		{
			sd:      sd1,
			channel: "foo",
		},
		{
			sd:      sd2,
			channel: "bar",
		},
		{
			channel:     "bar",
			sd:          sd3,
			expectedErr: errors.New("user revoked"),
		},
	}

	for _, tst := range testCases {
		
		invoked := false
		as.On("EligibleForService", tst.channel, tst.sd).Return(tst.expectedErr).Once().Run(func(_ mock.Arguments) {
			invoked = true
		})
		t.Run("Not cached test", func(t *testing.T) {
			err := cache.EligibleForService(tst.channel, tst.sd)
			if tst.expectedErr == nil {
				assert.NoError(t, err)
			} else {
				assert.Equal(t, tst.expectedErr.Error(), err.Error())
			}
			assert.True(t, invoked)
			
			invoked = false
		})

		
		
		
		t.Run("Cached test", func(t *testing.T) {
			err := cache.EligibleForService(tst.channel, tst.sd)
			if tst.expectedErr == nil {
				assert.NoError(t, err)
			} else {
				assert.Equal(t, tst.expectedErr.Error(), err.Error())
			}
			assert.False(t, invoked)
		})
	}
}

func TestCacheMarshalFailure(t *testing.T) {
	as := &mockAcSupport{}
	cache := newAuthCache(as, defaultConfig())
	asBytes = func(_ interface{}) ([]byte, error) {
		return nil, errors.New("failed marshaling ASN1")
	}
	defer func() {
		asBytes = asn1.Marshal
	}()
	err := cache.EligibleForService("mychannel", protoutil.SignedData{})
	assert.Contains(t, err.Error(), "failed marshaling ASN1")
}

func TestCacheConfigChange(t *testing.T) {
	as := &mockAcSupport{}
	sd := protoutil.SignedData{
		Data:      []byte{1, 2, 3},
		Identity:  []byte("identity"),
		Signature: []byte{1, 2, 3},
	}

	cache := newAuthCache(as, defaultConfig())

	
	as.On("EligibleForService", "mychannel", sd).Return(nil).Once()
	as.On("ConfigSequence", "mychannel").Return(uint64(0)).Times(2)
	err := cache.EligibleForService("mychannel", sd)
	assert.NoError(t, err)

	
	
	as.On("ConfigSequence", "mychannel").Return(uint64(0)).Once()
	err = cache.EligibleForService("mychannel", sd)
	assert.NoError(t, err)

	
	as.On("ConfigSequence", "mychannel").Return(uint64(1)).Times(2)
	as.On("EligibleForService", "mychannel", sd).Return(errors.New("unauthorized")).Once()
	err = cache.EligibleForService("mychannel", sd)
	assert.Contains(t, err.Error(), "unauthorized")
}

func TestCachePurgeCache(t *testing.T) {
	as := &mockAcSupport{}
	cache := newAuthCache(as, authCacheConfig{maxCacheSize: 4, purgeRetentionRatio: 0.75, enabled: true})
	as.On("ConfigSequence", "mychannel").Return(uint64(0))

	
	for _, id := range []string{"identity1", "identity2", "identity3", "identity4"} {
		sd := protoutil.SignedData{
			Data:      []byte{1, 2, 3},
			Identity:  []byte(id),
			Signature: []byte{1, 2, 3},
		}
		
		as.On("EligibleForService", "mychannel", sd).Return(nil).Once()
		err := cache.EligibleForService("mychannel", sd)
		assert.NoError(t, err)
	}

	
	var evicted int
	for _, id := range []string{"identity5", "identity1", "identity2"} {
		sd := protoutil.SignedData{
			Data:      []byte{1, 2, 3},
			Identity:  []byte(id),
			Signature: []byte{1, 2, 3},
		}
		as.On("EligibleForService", "mychannel", sd).Return(errors.New("unauthorized")).Once()
		err := cache.EligibleForService("mychannel", sd)
		if err != nil {
			evicted++
		}
	}
	assert.True(t, evicted > 0 && evicted < 4, "evicted: %d, but expected between 1 and 3 evictions", evicted)
}

func TestCacheConcurrentConfigUpdate(t *testing.T) {
	
	
	
	
	
	
	
	

	as := &mockAcSupport{}
	sd := protoutil.SignedData{
		Data:      []byte{1, 2, 3},
		Identity:  []byte{1, 2, 3},
		Signature: []byte{1, 2, 3},
	}
	var firstRequestInvoked sync.WaitGroup
	firstRequestInvoked.Add(1)
	var firstRequestFinished sync.WaitGroup
	firstRequestFinished.Add(1)
	var secondRequestFinished sync.WaitGroup
	secondRequestFinished.Add(1)
	cache := newAuthCache(as, defaultConfig())
	
	as.On("EligibleForService", "mychannel", mock.Anything).Return(nil).Once().Run(func(_ mock.Arguments) {
		firstRequestInvoked.Done()
		secondRequestFinished.Wait()
	})
	
	as.On("EligibleForService", "mychannel", mock.Anything).Return(errors.New("unauthorized")).Once()
	
	as.On("ConfigSequence", "mychannel").Return(uint64(0)).Once()
	as.On("ConfigSequence", "mychannel").Return(uint64(1)).Once()
	
	as.On("ConfigSequence", "mychannel").Return(uint64(1)).Times(2)

	
	go func() {
		defer firstRequestFinished.Done()
		firstResult := cache.EligibleForService("mychannel", sd)
		assert.NoError(t, firstResult)
	}()
	firstRequestInvoked.Wait()
	
	secondResult := cache.EligibleForService("mychannel", sd)
	
	secondRequestFinished.Done()
	
	firstRequestFinished.Wait()
	assert.Contains(t, secondResult.Error(), "unauthorized")

	
	
	as.On("ConfigSequence", "mychannel").Return(uint64(1)).Once()
	cachedResult := cache.EligibleForService("mychannel", sd)
	assert.Contains(t, cachedResult.Error(), "unauthorized")
}

func defaultConfig() authCacheConfig {
	return authCacheConfig{maxCacheSize: defaultMaxCacheSize, purgeRetentionRatio: defaultRetentionRatio, enabled: true}
}
