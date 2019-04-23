/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chaincode

import (
	"github.com/mcc-github/blockchain/core/container/ccintf"
	pb "github.com/mcc-github/blockchain/protos/peer"
)



func SetHandlerChaincodeID(h *Handler, chaincodeID *pb.ChaincodeID) {
	h.chaincodeID = chaincodeID
}

func SetHandlerChatStream(h *Handler, chatStream ccintf.ChaincodeStream) {
	h.chatStream = chatStream
}
