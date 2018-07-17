/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package identity

import (
	"bytes"
	"sync"
	"sync/atomic"
	"time"

	"github.com/mcc-github/blockchain/gossip/api"
	"github.com/mcc-github/blockchain/gossip/common"
	errors "github.com/pkg/errors"
)

var (
	
	
	usageThreshold = time.Hour
)



type Mapper interface {
	
	
	Put(pkiID common.PKIidType, identity api.PeerIdentityType) error

	
	
	Get(pkiID common.PKIidType) (api.PeerIdentityType, error)

	
	
	Sign(msg []byte) ([]byte, error)

	
	Verify(vkID, signature, message []byte) error

	
	GetPKIidOfCert(api.PeerIdentityType) common.PKIidType

	
	SuspectPeers(isSuspected api.PeerSuspector)

	
	IdentityInfo() api.PeerIdentitySet

	
	Stop()
}

type purgeTrigger func(pkiID common.PKIidType, identity api.PeerIdentityType)


type identityMapperImpl struct {
	onPurge    purgeTrigger
	mcs        api.MessageCryptoService
	sa         api.SecurityAdvisor
	pkiID2Cert map[string]*storedIdentity
	sync.RWMutex
	stopChan chan struct{}
	sync.Once
	selfPKIID string
}


func NewIdentityMapper(mcs api.MessageCryptoService, selfIdentity api.PeerIdentityType, onPurge purgeTrigger, sa api.SecurityAdvisor) Mapper {
	selfPKIID := mcs.GetPKIidOfCert(selfIdentity)
	idMapper := &identityMapperImpl{
		onPurge:    onPurge,
		mcs:        mcs,
		pkiID2Cert: make(map[string]*storedIdentity),
		stopChan:   make(chan struct{}),
		selfPKIID:  string(selfPKIID),
		sa:         sa,
	}
	if err := idMapper.Put(selfPKIID, selfIdentity); err != nil {
		panic(errors.Wrap(err, "Failed putting our own identity into the identity mapper"))
	}
	go idMapper.periodicalPurgeUnusedIdentities()
	return idMapper
}

func (is *identityMapperImpl) periodicalPurgeUnusedIdentities() {
	for {
		select {
		case <-is.stopChan:
			return
		case <-time.After(usageThreshold / 10):
			is.SuspectPeers(func(_ api.PeerIdentityType) bool {
				return false
			})
		}
	}
}



func (is *identityMapperImpl) Put(pkiID common.PKIidType, identity api.PeerIdentityType) error {
	if pkiID == nil {
		return errors.New("PKIID is nil")
	}
	if identity == nil {
		return errors.New("identity is nil")
	}

	expirationDate, err := is.mcs.Expiration(identity)
	if err != nil {
		return errors.Wrap(err, "failed classifying identity")
	}

	if err := is.mcs.ValidateIdentity(identity); err != nil {
		return err
	}

	id := is.mcs.GetPKIidOfCert(identity)
	if !bytes.Equal(pkiID, id) {
		return errors.New("identity doesn't match the computed pkiID")
	}

	is.Lock()
	defer is.Unlock()
	
	
	if _, exists := is.pkiID2Cert[string(pkiID)]; exists {
		return nil
	}

	var expirationTimer *time.Timer
	if !expirationDate.IsZero() {
		if time.Now().After(expirationDate) {
			return errors.New("identity expired")
		}
		
		timeToLive := expirationDate.Add(time.Millisecond).Sub(time.Now())
		expirationTimer = time.AfterFunc(timeToLive, func() {
			is.delete(pkiID, identity)
		})
	}

	is.pkiID2Cert[string(id)] = newStoredIdentity(pkiID, identity, expirationTimer, is.sa.OrgByPeerIdentity(identity))
	return nil
}



func (is *identityMapperImpl) Get(pkiID common.PKIidType) (api.PeerIdentityType, error) {
	is.RLock()
	defer is.RUnlock()
	storedIdentity, exists := is.pkiID2Cert[string(pkiID)]
	if !exists {
		return nil, errors.New("PKIID wasn't found")
	}
	return storedIdentity.fetchIdentity(), nil
}



func (is *identityMapperImpl) Sign(msg []byte) ([]byte, error) {
	return is.mcs.Sign(msg)
}

func (is *identityMapperImpl) Stop() {
	is.Once.Do(func() {
		is.stopChan <- struct{}{}
	})
}


func (is *identityMapperImpl) Verify(vkID, signature, message []byte) error {
	cert, err := is.Get(vkID)
	if err != nil {
		return err
	}
	return is.mcs.Verify(cert, signature, message)
}


func (is *identityMapperImpl) GetPKIidOfCert(identity api.PeerIdentityType) common.PKIidType {
	return is.mcs.GetPKIidOfCert(identity)
}


func (is *identityMapperImpl) SuspectPeers(isSuspected api.PeerSuspector) {
	for _, identity := range is.validateIdentities(isSuspected) {
		identity.cancelExpirationTimer()
		is.delete(identity.pkiID, identity.peerIdentity)
	}
}



func (is *identityMapperImpl) validateIdentities(isSuspected api.PeerSuspector) []*storedIdentity {
	now := time.Now()
	is.RLock()
	defer is.RUnlock()
	var revokedIdentities []*storedIdentity
	for pkiID, storedIdentity := range is.pkiID2Cert {
		if pkiID != is.selfPKIID && storedIdentity.fetchLastAccessTime().Add(usageThreshold).Before(now) {
			revokedIdentities = append(revokedIdentities, storedIdentity)
			continue
		}
		if !isSuspected(storedIdentity.peerIdentity) {
			continue
		}
		if err := is.mcs.ValidateIdentity(storedIdentity.fetchIdentity()); err != nil {
			revokedIdentities = append(revokedIdentities, storedIdentity)
		}
	}
	return revokedIdentities
}


func (is *identityMapperImpl) IdentityInfo() api.PeerIdentitySet {
	var res api.PeerIdentitySet
	is.RLock()
	defer is.RUnlock()
	for _, storedIdentity := range is.pkiID2Cert {
		res = append(res, api.PeerIdentityInfo{
			Identity:     storedIdentity.peerIdentity,
			PKIId:        storedIdentity.pkiID,
			Organization: storedIdentity.orgId,
		})
	}
	return res
}

func (is *identityMapperImpl) delete(pkiID common.PKIidType, identity api.PeerIdentityType) {
	is.Lock()
	defer is.Unlock()
	is.onPurge(pkiID, identity)
	delete(is.pkiID2Cert, string(pkiID))
}

type storedIdentity struct {
	pkiID           common.PKIidType
	lastAccessTime  int64
	peerIdentity    api.PeerIdentityType
	orgId           api.OrgIdentityType
	expirationTimer *time.Timer
}

func newStoredIdentity(pkiID common.PKIidType, identity api.PeerIdentityType, expirationTimer *time.Timer, org api.OrgIdentityType) *storedIdentity {
	return &storedIdentity{
		pkiID:           pkiID,
		lastAccessTime:  time.Now().UnixNano(),
		peerIdentity:    identity,
		expirationTimer: expirationTimer,
		orgId:           org,
	}
}

func (si *storedIdentity) fetchIdentity() api.PeerIdentityType {
	atomic.StoreInt64(&si.lastAccessTime, time.Now().UnixNano())
	return si.peerIdentity
}

func (si *storedIdentity) fetchLastAccessTime() time.Time {
	return time.Unix(0, atomic.LoadInt64(&si.lastAccessTime))
}

func (si *storedIdentity) cancelExpirationTimer() {
	if si.expirationTimer == nil {
		return
	}
	si.expirationTimer.Stop()
}




func SetIdentityUsageThreshold(duration time.Duration) {
	usageThreshold = duration
}




func GetIdentityUsageThreshold() time.Duration {
	return usageThreshold
}
