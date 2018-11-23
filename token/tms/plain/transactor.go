/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package plain

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/core/ledger/customtx"
	"github.com/mcc-github/blockchain/protos/ledger/queryresult"
	"github.com/mcc-github/blockchain/protos/token"
	"github.com/mcc-github/blockchain/token/ledger"
	"github.com/pkg/errors"
)

var namespace = "tms"

const tokenInput = "tokenInput"


type Transactor struct {
	PublicCredential []byte
	Ledger           ledger.LedgerReader
}



func (t *Transactor) RequestTransfer(request *token.TransferRequest) (*token.TokenTransaction, error) {
	var outputs []*token.PlainOutput

	inputs, tokenType, _, err := t.getInputsFromTokenIds(request.GetTokenIds())
	if err != nil {
		return nil, err
	}

	for _, ttt := range request.GetShares() {
		outputs = append(outputs, &token.PlainOutput{
			Owner:    ttt.Recipient,
			Type:     tokenType,
			Quantity: ttt.Quantity,
		})
	}

	
	transaction := &token.TokenTransaction{
		Action: &token.TokenTransaction_PlainAction{
			PlainAction: &token.PlainTokenAction{
				Data: &token.PlainTokenAction_PlainTransfer{
					PlainTransfer: &token.PlainTransfer{
						Inputs:  inputs,
						Outputs: outputs,
					},
				},
			},
		},
	}

	return transaction, nil
}


func (t *Transactor) RequestRedeem(request *token.RedeemRequest) (*token.TokenTransaction, error) {
	if len(request.GetTokenIds()) == 0 {
		return nil, errors.New("no token ids in RedeemRequest")
	}
	if request.GetQuantityToRedeem() <= 0 {
		return nil, errors.Errorf("quantity to redeem [%d] must be greater than 0", request.GetQuantityToRedeem())
	}

	inputs, tokenType, quantitySum, err := t.getInputsFromTokenIds(request.GetTokenIds())
	if err != nil {
		return nil, err
	}

	if quantitySum < request.QuantityToRedeem {
		return nil, errors.Errorf("total quantity [%d] from TokenIds is less than quantity [%d] to be redeemed", quantitySum, request.QuantityToRedeem)
	}

	
	var outputs []*token.PlainOutput
	outputs = append(outputs, &token.PlainOutput{
		Type:     tokenType,
		Quantity: request.QuantityToRedeem,
	})

	
	if quantitySum > request.QuantityToRedeem {
		outputs = append(outputs, &token.PlainOutput{
			Owner:    t.PublicCredential, 
			Type:     tokenType,
			Quantity: quantitySum - request.QuantityToRedeem,
		})
	}

	
	transaction := &token.TokenTransaction{
		Action: &token.TokenTransaction_PlainAction{
			PlainAction: &token.PlainTokenAction{
				Data: &token.PlainTokenAction_PlainRedeem{
					PlainRedeem: &token.PlainTransfer{
						Inputs:  inputs,
						Outputs: outputs,
					},
				},
			},
		},
	}

	return transaction, nil
}



func (t *Transactor) getInputsFromTokenIds(tokenIds [][]byte) ([]*token.InputId, string, uint64, error) {
	var inputs []*token.InputId
	var tokenType string = ""
	var quantitySum uint64 = 0
	for _, inKeyBytes := range tokenIds {
		
		inKey := parseCompositeKeyBytes(inKeyBytes)

		
		namespace, components, err := splitCompositeKey(inKey)
		if err != nil {
			return nil, "", 0, &customtx.InvalidTxError{Msg: fmt.Sprintf("error splitting input composite key: '%s'", err)}
		}
		if namespace != tokenOutput {
			return nil, "", 0, &customtx.InvalidTxError{Msg: fmt.Sprintf("namespace not '%s': '%s'", tokenOutput, namespace)}
		}
		if len(components) != 2 {
			return nil, "", 0, &customtx.InvalidTxError{Msg: fmt.Sprintf("not enough components in output ID composite key; expected 2, received '%s'", components)}
		}
		txID := components[0]
		index, err := strconv.Atoi(components[1])
		if err != nil {
			return nil, "", 0, &customtx.InvalidTxError{Msg: fmt.Sprintf("error parsing output index '%s': '%s'", components[1], err)}
		}

		
		inBytes, err := t.Ledger.GetState(tokenNamespace, inKey)
		if err != nil {
			return nil, "", 0, err
		}
		if inBytes == nil {
			return nil, "", 0, &customtx.InvalidTxError{Msg: fmt.Sprintf("input '%s' does not exist", inKey)}
		}
		input := &token.PlainOutput{}
		err = proto.Unmarshal(inBytes, input)
		if err != nil {
			return nil, "", 0, &customtx.InvalidTxError{Msg: fmt.Sprintf("error unmarshaling input bytes: '%s'", err)}
		}

		
		if tokenType == "" {
			tokenType = input.Type
		} else if tokenType != input.Type {
			return nil, "", 0, &customtx.InvalidTxError{Msg: fmt.Sprintf("two or more token types specified in input: '%s', '%s'", tokenType, input.Type)}
		}

		
		inputs = append(inputs, &token.InputId{TxId: txID, Index: uint32(index)})

		
		quantitySum += input.Quantity
	}

	return inputs, tokenType, quantitySum, nil
}


func (t *Transactor) ListTokens() (*token.UnspentTokens, error) {
	iterator, err := t.Ledger.GetStateRangeScanIterator(namespace, "", "")
	if err != nil {
		return nil, err
	}

	tokens := make([]*token.TokenOutput, 0)
	prefix, err := createPrefix(tokenOutput)
	if err != nil {
		return nil, err
	}
	for {
		next, err := iterator.Next()

		switch {
		case err != nil:
			return nil, err

		case next == nil:
			
			return &token.UnspentTokens{Tokens: tokens}, nil

		default:
			result, ok := next.(*queryresult.KV)
			if !ok {
				return nil, errors.New("failed to retrieve unspent tokens: casting error")
			}
			if strings.HasPrefix(result.Key, prefix) {
				output := &token.PlainOutput{}
				err = proto.Unmarshal(result.Value, output)
				if err != nil {
					return nil, errors.New("failed to retrieve unspent tokens: casting error")
				}
				if string(output.Owner) == string(t.PublicCredential) {
					spent, err := t.isSpent(result.Key)
					if err != nil {
						return nil, err
					}
					if !spent {
						tokens = append(tokens,
							&token.TokenOutput{
								Type:     output.Type,
								Quantity: output.Quantity,
								Id:       getCompositeKeyBytes(result.Key),
							})
					}
				}
			}
		}
	}

}

func (t *Transactor) RequestApprove(request *token.ApproveRequest) (*token.TokenTransaction, error) {
	panic("implement me!")
}


func (t *Transactor) isSpent(outputID string) (bool, error) {
	key, err := createInputKey(outputID)
	if err != nil {
		return false, err
	}
	result, err := t.Ledger.GetState(namespace, key)
	if err != nil {
		return false, err
	}
	if result == nil {
		return false, nil
	}
	return true, nil
}



func createInputKey(outputID string) (string, error) {
	att := strings.Split(outputID, string(minUnicodeRuneValue))
	return createCompositeKey(tokenInput, att[1:])
}


func createPrefix(keyword string) (string, error) {
	return createCompositeKey(keyword, nil)
}


func GenerateKeyForTest(txID string, index int) (string, error) {
	return createOutputKey(txID, index)
}

func splitCompositeKey(compositeKey string) (string, []string, error) {
	componentIndex := 1
	components := []string{}
	for i := 1; i < len(compositeKey); i++ {
		if compositeKey[i] == minUnicodeRuneValue {
			components = append(components, compositeKey[componentIndex:i])
			componentIndex = i + 1
		}
	}
	if len(components) < 2 {
		return "", nil, errors.New("invalid composite key - no components found")
	}
	return components[0], components[1:], nil
}
