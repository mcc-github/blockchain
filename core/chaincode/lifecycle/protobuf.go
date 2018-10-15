/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package lifecycle

import (
	"github.com/golang/protobuf/proto"

	"github.com/pkg/errors"
)



type Protobuf interface {
	Marshal(msg proto.Message) (marshaled []byte, err error)
	Unmarshal(marshaled []byte, msg proto.Message) error
}


type ProtobufImpl struct{}


func (p ProtobufImpl) Marshal(msg proto.Message) ([]byte, error) {
	res, err := proto.Marshal(msg)
	return res, errors.WithStack(err)
}


func (p ProtobufImpl) Unmarshal(marshaled []byte, msg proto.Message) error {
	return errors.WithStack(proto.Unmarshal(marshaled, msg))
}
