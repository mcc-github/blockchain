/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package identity

import (
	"bytes"
	"errors"
	"fmt"
	"testing"
	"time"

	"strings"

	"github.com/mcc-github/blockchain/gossip/api"
	"github.com/mcc-github/blockchain/gossip/common"
	"github.com/mcc-github/blockchain/gossip/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	msgCryptoService = &naiveCryptoService{revokedIdentities: map[string]struct{}{}}
	dummyID          = api.PeerIdentityType("dummyID")
)

type naiveCryptoService struct {
	mock.Mock
	revokedIdentities map[string]struct{}
}

var noopPurgeTrigger = func(_ common.PKIidType, _ api.PeerIdentityType) {}

func init() {
	util.SetupTestLogging()
	msgCryptoService.On("Expiration", api.PeerIdentityType(dummyID)).Return(time.Now().Add(time.Hour), nil)
	msgCryptoService.On("Expiration", api.PeerIdentityType("yacovm")).Return(time.Now().Add(time.Hour), nil)
	msgCryptoService.On("Expiration", api.PeerIdentityType("not-yacovm")).Return(time.Now().Add(time.Hour), nil)
	msgCryptoService.On("Expiration", api.PeerIdentityType("invalidIdentity")).Return(time.Now().Add(time.Hour), nil)
}

func (cs *naiveCryptoService) OrgByPeerIdentity(id api.PeerIdentityType) api.OrgIdentityType {
	found := false
	for _, call := range cs.Mock.ExpectedCalls {
		if call.Method == "OrgByPeerIdentity" {
			found = true
		}
	}
	if !found {
		return nil
	}
	return cs.Called(id).Get(0).(api.OrgIdentityType)
}

func (cs *naiveCryptoService) Expiration(peerIdentity api.PeerIdentityType) (time.Time, error) {
	args := cs.Called(peerIdentity)
	t, err := args.Get(0), args.Get(1)
	if err == nil {
		return t.(time.Time), nil
	}
	return time.Time{}, err.(error)
}

func (cs *naiveCryptoService) ValidateIdentity(peerIdentity api.PeerIdentityType) error {
	if _, isRevoked := cs.revokedIdentities[string(cs.GetPKIidOfCert(peerIdentity))]; isRevoked {
		return errors.New("revoked")
	}
	return nil
}


func (*naiveCryptoService) GetPKIidOfCert(peerIdentity api.PeerIdentityType) common.PKIidType {
	return common.PKIidType(peerIdentity)
}



func (*naiveCryptoService) VerifyBlock(chainID common.ChainID, seqNum uint64, signedBlock []byte) error {
	return nil
}



func (*naiveCryptoService) VerifyByChannel(_ common.ChainID, _ api.PeerIdentityType, _, _ []byte) error {
	return nil
}



func (*naiveCryptoService) Sign(msg []byte) ([]byte, error) {
	return msg, nil
}




func (*naiveCryptoService) Verify(peerIdentity api.PeerIdentityType, signature, message []byte) error {
	equal := bytes.Equal(signature, message)
	if !equal {
		return fmt.Errorf("Wrong certificate")
	}
	return nil
}

func TestPut(t *testing.T) {
	idStore := NewIdentityMapper(msgCryptoService, dummyID, noopPurgeTrigger, msgCryptoService)
	identity := []byte("yacovm")
	identity2 := []byte("not-yacovm")
	identity3 := []byte("invalidIdentity")
	msgCryptoService.revokedIdentities[string(identity3)] = struct{}{}
	pkiID := msgCryptoService.GetPKIidOfCert(api.PeerIdentityType(identity))
	pkiID2 := msgCryptoService.GetPKIidOfCert(api.PeerIdentityType(identity2))
	pkiID3 := msgCryptoService.GetPKIidOfCert(api.PeerIdentityType(identity3))
	assert.NoError(t, idStore.Put(pkiID, identity))
	assert.NoError(t, idStore.Put(pkiID, identity))
	assert.Error(t, idStore.Put(nil, identity))
	assert.Error(t, idStore.Put(pkiID2, nil))
	assert.Error(t, idStore.Put(pkiID2, identity))
	assert.Error(t, idStore.Put(pkiID, identity2))
	assert.Error(t, idStore.Put(pkiID3, identity3))
}

func TestGet(t *testing.T) {
	idStore := NewIdentityMapper(msgCryptoService, dummyID, noopPurgeTrigger, msgCryptoService)
	identity := []byte("yacovm")
	identity2 := []byte("not-yacovm")
	pkiID := msgCryptoService.GetPKIidOfCert(api.PeerIdentityType(identity))
	pkiID2 := msgCryptoService.GetPKIidOfCert(api.PeerIdentityType(identity2))
	assert.NoError(t, idStore.Put(pkiID, identity))
	cert, err := idStore.Get(pkiID)
	assert.NoError(t, err)
	assert.Equal(t, api.PeerIdentityType(identity), cert)
	cert, err = idStore.Get(pkiID2)
	assert.Nil(t, cert)
	assert.Error(t, err)
}

func TestVerify(t *testing.T) {
	idStore := NewIdentityMapper(msgCryptoService, dummyID, noopPurgeTrigger, msgCryptoService)
	identity := []byte("yacovm")
	identity2 := []byte("not-yacovm")
	pkiID := msgCryptoService.GetPKIidOfCert(api.PeerIdentityType(identity))
	pkiID2 := msgCryptoService.GetPKIidOfCert(api.PeerIdentityType(identity2))
	idStore.Put(pkiID, api.PeerIdentityType(identity))
	signed, err := idStore.Sign([]byte("bla bla"))
	assert.NoError(t, err)
	assert.NoError(t, idStore.Verify(pkiID, signed, []byte("bla bla")))
	assert.Error(t, idStore.Verify(pkiID2, signed, []byte("bla bla")))
}

func TestListInvalidIdentities(t *testing.T) {
	deletedIdentities := make(chan string, 1)
	assertDeletedIdentity := func(expected string) {
		select {
		case <-time.After(time.Second * 10):
			t.Fatalf("Didn't detect a deleted identity, expected %s to be deleted", expected)
		case actual := <-deletedIdentities:
			assert.Equal(t, expected, actual)
		}
	}
	
	SetIdentityUsageThreshold(time.Millisecond * 500)
	assert.Equal(t, time.Millisecond*500, GetIdentityUsageThreshold())
	selfPKIID := msgCryptoService.GetPKIidOfCert(dummyID)
	idStore := NewIdentityMapper(msgCryptoService, dummyID, func(_ common.PKIidType, identity api.PeerIdentityType) {
		deletedIdentities <- string(identity)
	}, msgCryptoService)
	identity := []byte("yacovm")
	
	pkiID := msgCryptoService.GetPKIidOfCert(api.PeerIdentityType(identity))
	assert.NoError(t, idStore.Put(pkiID, api.PeerIdentityType(identity)))
	cert, err := idStore.Get(pkiID)
	assert.NoError(t, err)
	assert.NotNil(t, cert)
	
	msgCryptoService.revokedIdentities[string(pkiID)] = struct{}{}
	idStore.SuspectPeers(func(_ api.PeerIdentityType) bool {
		return true
	})
	
	cert, err = idStore.Get(pkiID)
	assert.Error(t, err)
	assert.Nil(t, cert)
	assertDeletedIdentity("yacovm")

	
	msgCryptoService.revokedIdentities = map[string]struct{}{}
	
	
	
	pkiID = msgCryptoService.GetPKIidOfCert(api.PeerIdentityType(identity))
	assert.NoError(t, idStore.Put(pkiID, api.PeerIdentityType(identity)))
	
	cert, err = idStore.Get(pkiID)
	assert.NoError(t, err)
	assert.NotNil(t, cert)
	time.Sleep(time.Second * 3)
	
	cert, err = idStore.Get(pkiID)
	assert.Error(t, err)
	assert.Nil(t, cert)
	assertDeletedIdentity("yacovm")
	
	_, err = idStore.Get(selfPKIID)
	assert.NoError(t, err)

	
	
	pkiID = msgCryptoService.GetPKIidOfCert(api.PeerIdentityType(identity))
	assert.NoError(t, idStore.Put(pkiID, api.PeerIdentityType(identity)))
	stopChan := make(chan struct{})
	go func() {
		for {
			select {
			case <-stopChan:
				return
			case <-time.After(time.Millisecond * 10):
				idStore.Get(pkiID)
			}
		}
	}()
	time.Sleep(time.Second * 3)
	
	cert, err = idStore.Get(pkiID)
	assert.NoError(t, err)
	assert.NotNil(t, cert)
	stopChan <- struct{}{}
	
	
	idStore.Stop()
	time.Sleep(time.Second * 3)
	
	cert, err = idStore.Get(pkiID)
	assert.NoError(t, err)
	assert.NotNil(t, cert)
}

func TestExpiration(t *testing.T) {
	deletedIdentities := make(chan string, 1)
	SetIdentityUsageThreshold(time.Second * 500)
	idStore := NewIdentityMapper(msgCryptoService, dummyID, func(_ common.PKIidType, identity api.PeerIdentityType) {
		deletedIdentities <- string(identity)
	}, msgCryptoService)
	assertDeletedIdentity := func(expected string) {
		select {
		case <-time.After(time.Second * 10):
			t.Fatalf("Didn't detect a deleted identity, expected %s to be deleted", expected)
		case actual := <-deletedIdentities:
			assert.Equal(t, expected, actual)
		}
	}
	x509Identity := api.PeerIdentityType("x509Identity")
	expiredX509Identity := api.PeerIdentityType("expiredX509Identity")
	nonX509Identity := api.PeerIdentityType("nonX509Identity")
	notSupportedIdentity := api.PeerIdentityType("notSupportedIdentity")
	x509PkiID := idStore.GetPKIidOfCert(x509Identity)
	expiredX509PkiID := idStore.GetPKIidOfCert(expiredX509Identity)
	nonX509PkiID := idStore.GetPKIidOfCert(nonX509Identity)
	notSupportedPkiID := idStore.GetPKIidOfCert(notSupportedIdentity)
	msgCryptoService.On("Expiration", x509Identity).Return(time.Now().Add(time.Second), nil)
	msgCryptoService.On("Expiration", expiredX509Identity).Return(time.Now().Add(-time.Second), nil)
	msgCryptoService.On("Expiration", nonX509Identity).Return(time.Time{}, nil)
	msgCryptoService.On("Expiration", notSupportedIdentity).Return(time.Time{}, errors.New("no MSP supports given identity"))
	
	err := idStore.Put(x509PkiID, x509Identity)
	assert.NoError(t, err)
	err = idStore.Put(expiredX509PkiID, expiredX509Identity)
	assert.Equal(t, "identity expired", err.Error())
	err = idStore.Put(nonX509PkiID, nonX509Identity)
	assert.NoError(t, err)
	err = idStore.Put(notSupportedPkiID, notSupportedIdentity)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no MSP supports given identity")

	
	returnedIdentity, err := idStore.Get(x509PkiID)
	assert.NoError(t, err)
	assert.NotEmpty(t, returnedIdentity)

	returnedIdentity, err = idStore.Get(nonX509PkiID)
	assert.NoError(t, err)
	assert.NotEmpty(t, returnedIdentity)

	
	time.Sleep(time.Second * 3)

	
	returnedIdentity, err = idStore.Get(x509PkiID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "PKIID wasn't found")
	assert.Empty(t, returnedIdentity)
	assertDeletedIdentity("x509Identity")

	returnedIdentity, err = idStore.Get(nonX509PkiID)
	assert.NoError(t, err)
	assert.NotEmpty(t, returnedIdentity)

	
	msgCryptoService.revokedIdentities[string(nonX509PkiID)] = struct{}{}
	idStore.SuspectPeers(func(_ api.PeerIdentityType) bool {
		return true
	})
	assertDeletedIdentity("nonX509Identity")
	msgCryptoService.revokedIdentities = map[string]struct{}{}
}

func TestExpirationPanic(t *testing.T) {
	identity3 := []byte("invalidIdentity")
	msgCryptoService.revokedIdentities[string(identity3)] = struct{}{}
	assert.Panics(t, func() {
		NewIdentityMapper(msgCryptoService, identity3, noopPurgeTrigger, msgCryptoService)
	})
}

func TestIdentityInfo(t *testing.T) {
	cs := &naiveCryptoService{}
	alice := api.PeerIdentityType("alicePeer")
	bob := api.PeerIdentityType("bobPeer")
	aliceID := cs.GetPKIidOfCert(alice)
	bobId := cs.GetPKIidOfCert(bob)
	cs.On("OrgByPeerIdentity", dummyID).Return(api.OrgIdentityType("D"))
	cs.On("OrgByPeerIdentity", alice).Return(api.OrgIdentityType("A"))
	cs.On("OrgByPeerIdentity", bob).Return(api.OrgIdentityType("B"))
	cs.On("Expiration", mock.Anything).Return(time.Now().Add(time.Minute), nil)
	idStore := NewIdentityMapper(cs, dummyID, noopPurgeTrigger, cs)
	idStore.Put(aliceID, alice)
	idStore.Put(bobId, bob)
	for org, id := range idStore.IdentityInfo().ByOrg() {
		identity := string(id[0].Identity)
		pkiID := string(id[0].PKIId)
		orgId := string(id[0].Organization)
		assert.Equal(t, org, orgId)
		assert.Equal(t, strings.ToLower(org), string(identity[0]))
		assert.Equal(t, strings.ToLower(org), string(pkiID[0]))
	}
}
