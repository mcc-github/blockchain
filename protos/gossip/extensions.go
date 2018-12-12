/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package gossip

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/common/util"
	"github.com/mcc-github/blockchain/gossip/api"
	"github.com/mcc-github/blockchain/gossip/common"
)


func NewGossipMessageComparator(dataBlockStorageSize int) common.MessageReplacingPolicy {
	return (&msgComparator{dataBlockStorageSize: dataBlockStorageSize}).getMsgReplacingPolicy()
}

type msgComparator struct {
	dataBlockStorageSize int
}

func (mc *msgComparator) getMsgReplacingPolicy() common.MessageReplacingPolicy {
	return func(this interface{}, that interface{}) common.InvalidationResult {
		return mc.invalidationPolicy(this, that)
	}
}

func (mc *msgComparator) invalidationPolicy(this interface{}, that interface{}) common.InvalidationResult {
	thisMsg := this.(*SignedGossipMessage)
	thatMsg := that.(*SignedGossipMessage)

	if thisMsg.IsAliveMsg() && thatMsg.IsAliveMsg() {
		return aliveInvalidationPolicy(thisMsg.GetAliveMsg(), thatMsg.GetAliveMsg())
	}

	if thisMsg.IsDataMsg() && thatMsg.IsDataMsg() {
		return mc.dataInvalidationPolicy(thisMsg.GetDataMsg(), thatMsg.GetDataMsg())
	}

	if thisMsg.IsStateInfoMsg() && thatMsg.IsStateInfoMsg() {
		return mc.stateInvalidationPolicy(thisMsg.GetStateInfo(), thatMsg.GetStateInfo())
	}

	if thisMsg.IsIdentityMsg() && thatMsg.IsIdentityMsg() {
		return mc.identityInvalidationPolicy(thisMsg.GetPeerIdentity(), thatMsg.GetPeerIdentity())
	}

	if thisMsg.IsLeadershipMsg() && thatMsg.IsLeadershipMsg() {
		return leaderInvalidationPolicy(thisMsg.GetLeadershipMsg(), thatMsg.GetLeadershipMsg())
	}

	return common.MessageNoAction
}

func (mc *msgComparator) stateInvalidationPolicy(thisStateMsg *StateInfo, thatStateMsg *StateInfo) common.InvalidationResult {
	if !bytes.Equal(thisStateMsg.PkiId, thatStateMsg.PkiId) {
		return common.MessageNoAction
	}
	return compareTimestamps(thisStateMsg.Timestamp, thatStateMsg.Timestamp)
}

func (mc *msgComparator) identityInvalidationPolicy(thisIdentityMsg *PeerIdentity, thatIdentityMsg *PeerIdentity) common.InvalidationResult {
	if bytes.Equal(thisIdentityMsg.PkiId, thatIdentityMsg.PkiId) {
		return common.MessageInvalidated
	}

	return common.MessageNoAction
}

func (mc *msgComparator) dataInvalidationPolicy(thisDataMsg *DataMessage, thatDataMsg *DataMessage) common.InvalidationResult {
	if thisDataMsg.Payload.SeqNum == thatDataMsg.Payload.SeqNum {
		return common.MessageInvalidated
	}

	diff := abs(thisDataMsg.Payload.SeqNum, thatDataMsg.Payload.SeqNum)
	if diff <= uint64(mc.dataBlockStorageSize) {
		return common.MessageNoAction
	}

	if thisDataMsg.Payload.SeqNum > thatDataMsg.Payload.SeqNum {
		return common.MessageInvalidates
	}
	return common.MessageInvalidated
}

func aliveInvalidationPolicy(thisMsg *AliveMessage, thatMsg *AliveMessage) common.InvalidationResult {
	if !bytes.Equal(thisMsg.Membership.PkiId, thatMsg.Membership.PkiId) {
		return common.MessageNoAction
	}

	return compareTimestamps(thisMsg.Timestamp, thatMsg.Timestamp)
}

func leaderInvalidationPolicy(thisMsg *LeadershipMessage, thatMsg *LeadershipMessage) common.InvalidationResult {
	if !bytes.Equal(thisMsg.PkiId, thatMsg.PkiId) {
		return common.MessageNoAction
	}

	return compareTimestamps(thisMsg.Timestamp, thatMsg.Timestamp)
}

func compareTimestamps(thisTS *PeerTime, thatTS *PeerTime) common.InvalidationResult {
	if thisTS.IncNum == thatTS.IncNum {
		if thisTS.SeqNum > thatTS.SeqNum {
			return common.MessageInvalidates
		}

		return common.MessageInvalidated
	}
	if thisTS.IncNum < thatTS.IncNum {
		return common.MessageInvalidated
	}
	return common.MessageInvalidates
}


func (m *GossipMessage) IsAliveMsg() bool {
	return m.GetAliveMsg() != nil
}


func (m *GossipMessage) IsDataMsg() bool {
	return m.GetDataMsg() != nil
}


func (m *GossipMessage) IsStateInfoPullRequestMsg() bool {
	return m.GetStateInfoPullReq() != nil
}


func (m *GossipMessage) IsStateInfoSnapshot() bool {
	return m.GetStateSnapshot() != nil
}


func (m *GossipMessage) IsStateInfoMsg() bool {
	return m.GetStateInfo() != nil
}



func (m *GossipMessage) IsPullMsg() bool {
	return m.GetDataReq() != nil || m.GetDataUpdate() != nil ||
		m.GetHello() != nil || m.GetDataDig() != nil
}


func (m *GossipMessage) IsRemoteStateMessage() bool {
	return m.GetStateRequest() != nil || m.GetStateResponse() != nil
}




func (m *GossipMessage) GetPullMsgType() PullMsgType {
	if helloMsg := m.GetHello(); helloMsg != nil {
		return helloMsg.MsgType
	}

	if digMsg := m.GetDataDig(); digMsg != nil {
		return digMsg.MsgType
	}

	if reqMsg := m.GetDataReq(); reqMsg != nil {
		return reqMsg.MsgType
	}

	if resMsg := m.GetDataUpdate(); resMsg != nil {
		return resMsg.MsgType
	}

	return PullMsgType_UNDEFINED
}



func (m *GossipMessage) IsChannelRestricted() bool {
	return m.Tag == GossipMessage_CHAN_AND_ORG || m.Tag == GossipMessage_CHAN_ONLY || m.Tag == GossipMessage_CHAN_OR_ORG
}



func (m *GossipMessage) IsOrgRestricted() bool {
	return m.Tag == GossipMessage_CHAN_AND_ORG || m.Tag == GossipMessage_ORG_ONLY
}


func (m *GossipMessage) IsIdentityMsg() bool {
	return m.GetPeerIdentity() != nil
}


func (m *GossipMessage) IsDataReq() bool {
	return m.GetDataReq() != nil
}


func (m *GossipMessage) IsPrivateDataMsg() bool {
	return m.GetPrivateReq() != nil || m.GetPrivateRes() != nil || m.GetPrivateData() != nil
}


func (m *GossipMessage) IsAck() bool {
	return m.GetAck() != nil
}


func (m *GossipMessage) IsDataUpdate() bool {
	return m.GetDataUpdate() != nil
}


func (m *GossipMessage) IsHelloMsg() bool {
	return m.GetHello() != nil
}


func (m *GossipMessage) IsDigestMsg() bool {
	return m.GetDataDig() != nil
}


func (m *GossipMessage) IsLeadershipMsg() bool {
	return m.GetLeadershipMsg() != nil
}


type MsgConsumer func(message *SignedGossipMessage)


type IdentifierExtractor func(*SignedGossipMessage) string



func (m *GossipMessage) IsTagLegal() error {
	if m.Tag == GossipMessage_UNDEFINED {
		return fmt.Errorf("Undefined tag")
	}
	if m.IsDataMsg() {
		if m.Tag != GossipMessage_CHAN_AND_ORG {
			return fmt.Errorf("Tag should be %s", GossipMessage_Tag_name[int32(GossipMessage_CHAN_AND_ORG)])
		}
		return nil
	}

	if m.IsAliveMsg() || m.GetMemReq() != nil || m.GetMemRes() != nil {
		if m.Tag != GossipMessage_EMPTY {
			return fmt.Errorf("Tag should be %s", GossipMessage_Tag_name[int32(GossipMessage_EMPTY)])
		}
		return nil
	}

	if m.IsIdentityMsg() {
		if m.Tag != GossipMessage_ORG_ONLY {
			return fmt.Errorf("Tag should be %s", GossipMessage_Tag_name[int32(GossipMessage_ORG_ONLY)])
		}
		return nil
	}

	if m.IsPullMsg() {
		switch m.GetPullMsgType() {
		case PullMsgType_BLOCK_MSG:
			if m.Tag != GossipMessage_CHAN_AND_ORG {
				return fmt.Errorf("Tag should be %s", GossipMessage_Tag_name[int32(GossipMessage_CHAN_AND_ORG)])
			}
			return nil
		case PullMsgType_IDENTITY_MSG:
			if m.Tag != GossipMessage_EMPTY {
				return fmt.Errorf("Tag should be %s", GossipMessage_Tag_name[int32(GossipMessage_EMPTY)])
			}
			return nil
		default:
			return fmt.Errorf("Invalid PullMsgType: %s", PullMsgType_name[int32(m.GetPullMsgType())])
		}
	}

	if m.IsStateInfoMsg() || m.IsStateInfoPullRequestMsg() || m.IsStateInfoSnapshot() || m.IsRemoteStateMessage() {
		if m.Tag != GossipMessage_CHAN_OR_ORG {
			return fmt.Errorf("Tag should be %s", GossipMessage_Tag_name[int32(GossipMessage_CHAN_OR_ORG)])
		}
		return nil
	}

	if m.IsLeadershipMsg() {
		if m.Tag != GossipMessage_CHAN_AND_ORG {
			return fmt.Errorf("Tag should be %s", GossipMessage_Tag_name[int32(GossipMessage_CHAN_AND_ORG)])
		}
		return nil
	}

	return fmt.Errorf("Unknown message type: %v", m)
}




type Verifier func(peerIdentity []byte, signature, message []byte) error



type Signer func(msg []byte) ([]byte, error)







type ReceivedMessage interface {

	
	Respond(msg *GossipMessage)

	
	GetGossipMessage() *SignedGossipMessage

	
	
	GetSourceEnvelope() *Envelope

	
	
	GetConnectionInfo() *ConnectionInfo

	
	
	
	Ack(err error)
}



type ConnectionInfo struct {
	ID       common.PKIidType
	Auth     *AuthInfo
	Identity api.PeerIdentityType
	Endpoint string
}


func (c *ConnectionInfo) String() string {
	return fmt.Sprintf("%s %v", c.Endpoint, c.ID)
}




type AuthInfo struct {
	SignedData []byte
	Signature  []byte
}




func (m *SignedGossipMessage) Sign(signer Signer) (*Envelope, error) {
	
	
	var secretEnvelope *SecretEnvelope
	if m.Envelope != nil {
		secretEnvelope = m.Envelope.SecretEnvelope
	}
	m.Envelope = nil
	payload, err := proto.Marshal(m.GossipMessage)
	if err != nil {
		return nil, err
	}
	sig, err := signer(payload)
	if err != nil {
		return nil, err
	}

	e := &Envelope{
		Payload:        payload,
		Signature:      sig,
		SecretEnvelope: secretEnvelope,
	}
	m.Envelope = e
	return e, nil
}


func (m *GossipMessage) NoopSign() (*SignedGossipMessage, error) {
	signer := func(msg []byte) ([]byte, error) {
		return nil, nil
	}
	sMsg := &SignedGossipMessage{
		GossipMessage: m,
	}
	_, err := sMsg.Sign(signer)
	return sMsg, err
}



func (m *SignedGossipMessage) Verify(peerIdentity []byte, verify Verifier) error {
	if m.Envelope == nil {
		return errors.New("Missing envelope")
	}
	if len(m.Envelope.Payload) == 0 {
		return errors.New("Empty payload")
	}
	if len(m.Envelope.Signature) == 0 {
		return errors.New("Empty signature")
	}
	payloadSigVerificationErr := verify(peerIdentity, m.Envelope.Signature, m.Envelope.Payload)
	if payloadSigVerificationErr != nil {
		return payloadSigVerificationErr
	}
	if m.Envelope.SecretEnvelope != nil {
		payload := m.Envelope.SecretEnvelope.Payload
		sig := m.Envelope.SecretEnvelope.Signature
		if len(payload) == 0 {
			return errors.New("Empty payload")
		}
		if len(sig) == 0 {
			return errors.New("Empty signature")
		}
		return verify(peerIdentity, sig, payload)
	}
	return nil
}



func (m *SignedGossipMessage) IsSigned() bool {
	return m.Envelope != nil && m.Envelope.Payload != nil && m.Envelope.Signature != nil
}




func (e *Envelope) ToGossipMessage() (*SignedGossipMessage, error) {
	if e == nil {
		return nil, errors.New("nil envelope")
	}
	msg := &GossipMessage{}
	err := proto.Unmarshal(e.Payload, msg)
	if err != nil {
		return nil, fmt.Errorf("Failed unmarshaling GossipMessage from envelope: %v", err)
	}
	return &SignedGossipMessage{
		GossipMessage: msg,
		Envelope:      e,
	}, nil
}



func (e *Envelope) SignSecret(signer Signer, secret *Secret) error {
	payload, err := proto.Marshal(secret)
	if err != nil {
		return err
	}
	sig, err := signer(payload)
	if err != nil {
		return err
	}
	e.SecretEnvelope = &SecretEnvelope{
		Payload:   payload,
		Signature: sig,
	}
	return nil
}




func (s *SecretEnvelope) InternalEndpoint() string {
	secret := &Secret{}
	if err := proto.Unmarshal(s.Payload, secret); err != nil {
		return ""
	}
	return secret.GetInternalEndpoint()
}



type SignedGossipMessage struct {
	*Envelope
	*GossipMessage
}

func (p *Payload) toString() string {
	return fmt.Sprintf("Block message: {Data: %d bytes, seq: %d}", len(p.Data), p.SeqNum)
}

func (du *DataUpdate) toString() string {
	mType := PullMsgType_name[int32(du.MsgType)]
	return fmt.Sprintf("Type: %s, items: %d, nonce: %d", mType, len(du.Data), du.Nonce)
}

func (mr *MembershipResponse) toString() string {
	return fmt.Sprintf("MembershipResponse with Alive: %d, Dead: %d", len(mr.Alive), len(mr.Dead))
}

func (sis *StateInfoSnapshot) toString() string {
	return fmt.Sprintf("StateInfoSnapshot with %d items", len(sis.Elements))
}



func (m *SignedGossipMessage) String() string {
	env := "No envelope"
	if m.Envelope != nil {
		var secretEnv string
		if m.SecretEnvelope != nil {
			pl := len(m.SecretEnvelope.Payload)
			sl := len(m.SecretEnvelope.Signature)
			secretEnv = fmt.Sprintf(" Secret payload: %d bytes, Secret Signature: %d bytes", pl, sl)
		}
		env = fmt.Sprintf("%d bytes, Signature: %d bytes%s", len(m.Envelope.Payload), len(m.Envelope.Signature), secretEnv)
	}
	gMsg := "No gossipMessage"
	if m.GossipMessage != nil {
		var isSimpleMsg bool
		if m.GetStateResponse() != nil {
			gMsg = fmt.Sprintf("StateResponse with %d items", len(m.GetStateResponse().Payloads))
		} else if m.IsDataMsg() && m.GetDataMsg().Payload != nil {
			gMsg = m.GetDataMsg().Payload.toString()
		} else if m.IsDataUpdate() {
			update := m.GetDataUpdate()
			gMsg = fmt.Sprintf("DataUpdate: %s", update.toString())
		} else if m.GetMemRes() != nil {
			gMsg = m.GetMemRes().toString()
		} else if m.IsStateInfoSnapshot() {
			gMsg = m.GetStateSnapshot().toString()
		} else if m.GetPrivateRes() != nil {
			gMsg = m.GetPrivateRes().ToString()
		} else {
			gMsg = m.GossipMessage.String()
			isSimpleMsg = true
		}
		if !isSimpleMsg {
			desc := fmt.Sprintf("Channel: %s, nonce: %d, tag: %s", string(m.Channel), m.Nonce, GossipMessage_Tag_name[int32(m.Tag)])
			gMsg = fmt.Sprintf("%s %s", desc, gMsg)
		}
	}
	return fmt.Sprintf("GossipMessage: %v, Envelope: %s", gMsg, env)
}

func (dd *DataRequest) FormattedDigests() []string {
	if dd.MsgType == PullMsgType_IDENTITY_MSG {
		return digestsToHex(dd.Digests)
	}

	return digestsAsStrings(dd.Digests)
}

func (dd *DataDigest) FormattedDigests() []string {
	if dd.MsgType == PullMsgType_IDENTITY_MSG {
		return digestsToHex(dd.Digests)
	}
	return digestsAsStrings(dd.Digests)
}


func (dig *PvtDataDigest) Hash() (string, error) {
	b, err := proto.Marshal(dig)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(util.ComputeSHA256(b)), nil
}


func (res *RemotePvtDataResponse) ToString() string {
	a := make([]string, len(res.Elements))
	for i, el := range res.Elements {
		a[i] = fmt.Sprintf("%s with %d elements", el.Digest.String(), len(el.Payload))
	}
	return fmt.Sprintf("%v", a)
}

func digestsAsStrings(digests [][]byte) []string {
	a := make([]string, len(digests))
	for i, dig := range digests {
		a[i] = string(dig)
	}
	return a
}

func digestsToHex(digests [][]byte) []string {
	a := make([]string, len(digests))
	for i, dig := range digests {
		a[i] = hex.EncodeToString(dig)
	}
	return a
}



func (msg *StateInfo) LedgerHeight() (uint64, error) {
	if msg.Properties != nil {
		return msg.Properties.LedgerHeight, nil
	}
	return 0, errors.New("properties undefined")
}


func abs(a, b uint64) uint64 {
	if a > b {
		return a - b
	}
	return b - a
}
