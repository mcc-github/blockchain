/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package dispatcher

import (
	"fmt"
	"reflect"

	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
)



type Dispatcher struct {
	
	Protobuf Protobuf
}





func (d *Dispatcher) Dispatch(inputBytes []byte, methodName string, receiver interface{}) ([]byte, error) {
	method := reflect.ValueOf(receiver).MethodByName(methodName)

	if method == (reflect.Value{}) {
		return nil, errors.Errorf("receiver %T.%s does not exist", receiver, methodName)
	}

	if method.Type().NumIn() != 1 {
		return nil, errors.Errorf("receiver %T.%s has %d parameters but expected 1", receiver, methodName, method.Type().NumIn())
	}

	inputType := method.Type().In(0)
	if inputType.Kind() != reflect.Ptr {
		return nil, errors.Errorf("receiver %T.%s does not accept a pointer as its argument", receiver, methodName)
	}

	if method.Type().NumOut() != 2 {
		return nil, errors.Errorf("receiver %T.%s returns %d values but expected 2", receiver, methodName, method.Type().NumOut())
	}

	if !method.Type().Out(0).Implements(reflect.TypeOf((*proto.Message)(nil)).Elem()) {
		return nil, errors.Errorf("receiver %T.%s does not return a an implementor of proto.Message as its first return value", receiver, methodName)
	}

	if !method.Type().Out(1).Implements(reflect.TypeOf((*error)(nil)).Elem()) {
		return nil, errors.Errorf("receiver %T.%s does not return an error as its second return value", receiver, methodName)
	}

	inputValue := reflect.New(inputType.Elem())
	inputMsg, ok := inputValue.Interface().(proto.Message)
	if !ok {
		return nil, errors.Errorf("receiver %T.%s does not accept a proto.Message as its argument, it is '%T'", receiver, methodName, inputValue.Interface())
	}

	err := d.Protobuf.Unmarshal(inputBytes, inputMsg)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("could not decode input arg for %T.%s", receiver, methodName))
	}

	outputVals := method.Call([]reflect.Value{inputValue})

	if !outputVals[1].IsNil() {
		return nil, outputVals[1].Interface().(error)
	}

	if outputVals[0].IsNil() {
		return nil, errors.Errorf("receiver %T.%s returned (nil, nil) which is not allowed", receiver, methodName)
	}

	outputMsg := outputVals[0].Interface().(proto.Message)

	resultBytes, err := d.Protobuf.Marshal(outputMsg)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("failed to marshal result for %T.%s", receiver, methodName))
	}

	return resultBytes, nil
}
