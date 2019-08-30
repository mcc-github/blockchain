/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package protoext

import (
	"fmt"

	"github.com/mcc-github/blockchain-protos-go/gossip"
)


func IsAliveMsg(m *gossip.GossipMessage) bool {
	return m.GetAliveMsg() != nil
}


func IsDataMsg(m *gossip.GossipMessage) bool {
	return m.GetDataMsg() != nil
}


func IsStateInfoPullRequestMsg(m *gossip.GossipMessage) bool {
	return m.GetStateInfoPullReq() != nil
}


func IsStateInfoSnapshot(m *gossip.GossipMessage) bool {
	return m.GetStateSnapshot() != nil
}


func IsStateInfoMsg(m *gossip.GossipMessage) bool {
	return m.GetStateInfo() != nil
}



func IsPullMsg(m *gossip.GossipMessage) bool {
	return m.GetDataReq() != nil || m.GetDataUpdate() != nil ||
		m.GetHello() != nil || m.GetDataDig() != nil
}


func IsRemoteStateMessage(m *gossip.GossipMessage) bool {
	return m.GetStateRequest() != nil || m.GetStateResponse() != nil
}




func GetPullMsgType(m *gossip.GossipMessage) gossip.PullMsgType {
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

	return gossip.PullMsgType_UNDEFINED
}



func IsChannelRestricted(m *gossip.GossipMessage) bool {
	return m.Tag == gossip.GossipMessage_CHAN_AND_ORG || m.Tag == gossip.GossipMessage_CHAN_ONLY || m.Tag == gossip.GossipMessage_CHAN_OR_ORG
}



func IsOrgRestricted(m *gossip.GossipMessage) bool {
	return m.Tag == gossip.GossipMessage_CHAN_AND_ORG || m.Tag == gossip.GossipMessage_ORG_ONLY
}


func IsIdentityMsg(m *gossip.GossipMessage) bool {
	return m.GetPeerIdentity() != nil
}


func IsDataReq(m *gossip.GossipMessage) bool {
	return m.GetDataReq() != nil
}


func IsPrivateDataMsg(m *gossip.GossipMessage) bool {
	return m.GetPrivateReq() != nil || m.GetPrivateRes() != nil || m.GetPrivateData() != nil
}


func IsAck(m *gossip.GossipMessage) bool {
	return m.GetAck() != nil
}


func IsDataUpdate(m *gossip.GossipMessage) bool {
	return m.GetDataUpdate() != nil
}


func IsHelloMsg(m *gossip.GossipMessage) bool {
	return m.GetHello() != nil
}


func IsDigestMsg(m *gossip.GossipMessage) bool {
	return m.GetDataDig() != nil
}


func IsLeadershipMsg(m *gossip.GossipMessage) bool {
	return m.GetLeadershipMsg() != nil
}



func IsTagLegal(m *gossip.GossipMessage) error {
	if m.Tag == gossip.GossipMessage_UNDEFINED {
		return fmt.Errorf("Undefined tag")
	}
	if IsDataMsg(m) {
		if m.Tag != gossip.GossipMessage_CHAN_AND_ORG {
			return fmt.Errorf("Tag should be %s", gossip.GossipMessage_Tag_name[int32(gossip.GossipMessage_CHAN_AND_ORG)])
		}
		return nil
	}

	if IsAliveMsg(m) || m.GetMemReq() != nil || m.GetMemRes() != nil {
		if m.Tag != gossip.GossipMessage_EMPTY {
			return fmt.Errorf("Tag should be %s", gossip.GossipMessage_Tag_name[int32(gossip.GossipMessage_EMPTY)])
		}
		return nil
	}

	if IsIdentityMsg(m) {
		if m.Tag != gossip.GossipMessage_ORG_ONLY {
			return fmt.Errorf("Tag should be %s", gossip.GossipMessage_Tag_name[int32(gossip.GossipMessage_ORG_ONLY)])
		}
		return nil
	}

	if IsPullMsg(m) {
		switch GetPullMsgType(m) {
		case gossip.PullMsgType_BLOCK_MSG:
			if m.Tag != gossip.GossipMessage_CHAN_AND_ORG {
				return fmt.Errorf("Tag should be %s", gossip.GossipMessage_Tag_name[int32(gossip.GossipMessage_CHAN_AND_ORG)])
			}
			return nil
		case gossip.PullMsgType_IDENTITY_MSG:
			if m.Tag != gossip.GossipMessage_EMPTY {
				return fmt.Errorf("Tag should be %s", gossip.GossipMessage_Tag_name[int32(gossip.GossipMessage_EMPTY)])
			}
			return nil
		default:
			return fmt.Errorf("Invalid PullMsgType: %s", gossip.PullMsgType_name[int32(GetPullMsgType(m))])
		}
	}

	if IsStateInfoMsg(m) || IsStateInfoPullRequestMsg(m) || IsStateInfoSnapshot(m) || IsRemoteStateMessage(m) {
		if m.Tag != gossip.GossipMessage_CHAN_OR_ORG {
			return fmt.Errorf("Tag should be %s", gossip.GossipMessage_Tag_name[int32(gossip.GossipMessage_CHAN_OR_ORG)])
		}
		return nil
	}

	if IsLeadershipMsg(m) {
		if m.Tag != gossip.GossipMessage_CHAN_AND_ORG {
			return fmt.Errorf("Tag should be %s", gossip.GossipMessage_Tag_name[int32(gossip.GossipMessage_CHAN_AND_ORG)])
		}
		return nil
	}

	return fmt.Errorf("Unknown message type: %v", m)
}
