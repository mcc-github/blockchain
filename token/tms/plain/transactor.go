/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package plain

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/protos/ledger/queryresult"
	"github.com/mcc-github/blockchain/protos/token"
	"github.com/mcc-github/blockchain/token/identity"
	"github.com/mcc-github/blockchain/token/ledger"
	"github.com/pkg/errors"
)

const (
	Precision uint64 = 64
)


type Transactor struct {
	PublicCredential    []byte
	Ledger              ledger.LedgerReader
	TokenOwnerValidator identity.TokenOwnerValidator
}


func (t *Transactor) RequestTransfer(request *token.TransferRequest) (*token.TokenTransaction, error) {
	var outputs []*token.Token
	var outputSum = NewZeroQuantity(Precision)

	if len(request.GetTokenIds()) == 0 {
		return nil, errors.New("no token IDs in transfer request")
	}
	if len(request.GetShares()) == 0 {
		return nil, errors.New("no shares in transfer request")
	}

	tokenType, inputSum, _, err := t.getInputsFromTokenIds(request.GetTokenIds(), nil)
	if err != nil {
		return nil, err
	}

	for _, ttt := range request.GetShares() {
		err := t.TokenOwnerValidator.Validate(ttt.Recipient)
		if err != nil {
			return nil, errors.Errorf("invalid recipient in transfer request '%s'", err)
		}
		q, err := ToQuantity(ttt.Quantity, Precision)
		if err != nil {
			return nil, errors.Errorf("invalid quantity in transfer request '%s'", err)
		}

		outputSum, err = outputSum.Add(q)
		if err != nil {
			return nil, errors.Errorf("failed adding up output quantities, err '%s'", err)
		}

		outputs = append(outputs, &token.Token{
			Owner:    ttt.Recipient,
			Type:     tokenType,
			Quantity: q.Hex(),
		})
	}

	
	cmp, err := inputSum.Cmp(outputSum)
	if err != nil {
		return nil, errors.Errorf("cannot compare quantities '%s'", err)
	}
	if cmp < 0 {
		return nil, errors.Errorf("total quantity [%d] from TokenIds is less than total quantity [%s] for transfer", inputSum, outputSum)
	}

	
	if cmp > 0 {
		change, err := inputSum.Sub(outputSum)
		if err != nil {
			return nil, errors.Errorf("failed computing change, err '%s'", err)
		}

		outputs = append(outputs, &token.Token{
			
			Owner:    &token.TokenOwner{Type: token.TokenOwner_MSP_IDENTIFIER, Raw: t.PublicCredential},
			Type:     tokenType,
			Quantity: change.Hex(),
		})
	}

	
	transaction := &token.TokenTransaction{
		Action: &token.TokenTransaction_TokenAction{
			TokenAction: &token.TokenAction{
				Data: &token.TokenAction_Transfer{
					Transfer: &token.Transfer{
						Inputs:  request.GetTokenIds(),
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
	quantityToRedeem, err := ToQuantity(request.GetQuantity(), Precision)
	if err != nil {
		return nil, errors.Errorf("quantity to redeem [%s] is invalid, err '%s'", request.GetQuantity(), err)
	}

	tokenType, quantitySum, _, err := t.getInputsFromTokenIds(request.GetTokenIds(), nil)
	if err != nil {
		return nil, err
	}

	cmp, err := quantitySum.Cmp(quantityToRedeem)
	if err != nil {
		return nil, errors.Errorf("cannot compare quantities '%s'", err)
	}
	if cmp < 0 {
		return nil, errors.Errorf("total quantity [%d] from TokenIds is less than quantity [%s] to be redeemed", quantitySum, request.Quantity)
	}

	
	var outputs []*token.Token
	outputs = append(outputs, &token.Token{
		Type:     tokenType,
		Quantity: quantityToRedeem.Hex(),
	})

	
	if cmp > 0 {
		change, err := quantitySum.Sub(quantityToRedeem)
		if err != nil {
			return nil, errors.Errorf("failed computing change, err '%s'", err)
		}

		outputs = append(outputs, &token.Token{
			
			Owner:    &token.TokenOwner{Type: token.TokenOwner_MSP_IDENTIFIER, Raw: t.PublicCredential}, 
			Type:     tokenType,
			Quantity: change.Hex(),
		})
	}

	
	transaction := &token.TokenTransaction{
		Action: &token.TokenTransaction_TokenAction{
			TokenAction: &token.TokenAction{
				Data: &token.TokenAction_Redeem{
					Redeem: &token.Transfer{
						Inputs:  request.GetTokenIds(),
						Outputs: outputs,
					},
				},
			},
		},
	}

	return transaction, nil
}





func (t *Transactor) getInputsFromTokenIds(tokenIds []*token.TokenId, upperBound Quantity) (string, Quantity, int, error) {
	
	tokenOwner := &token.TokenOwner{Type: token.TokenOwner_MSP_IDENTIFIER, Raw: t.PublicCredential}
	ownerString, err := GetTokenOwnerString(tokenOwner)
	if err != nil {
		return "", nil, 0, err
	}

	var tokenType = ""
	var sum = NewZeroQuantity(Precision)
	counter := 0
	for _, tokenId := range tokenIds {
		
		inKey, err := createTokenKey(ownerString, tokenId.TxId, int(tokenId.Index))
		if err != nil {
			verifierLogger.Errorf("error getting creating input key: %s", err)
			return "", nil, 0, err
		}
		verifierLogger.Debugf("transferring token with ID: '%s'", inKey)

		
		verifierLogger.Debugf("getting output '%s' to spend from ledger", inKey)
		inBytes, err := t.Ledger.GetState(tokenNameSpace, inKey)
		if err != nil {
			verifierLogger.Errorf("error getting token '%s' to spend from ledger: %s", inKey, err)
			return "", nil, 0, err
		}
		if len(inBytes) == 0 {
			return "", nil, 0, errors.New(fmt.Sprintf("input TokenId (%s, %d) does not exist or not owned by the user", tokenId.TxId, tokenId.Index))
		}
		input := &token.Token{}
		err = proto.Unmarshal(inBytes, input)
		if err != nil {
			return "", nil, 0, errors.New(fmt.Sprintf("error unmarshaling input bytes: '%s'", err))
		}

		
		if !bytes.Equal(t.PublicCredential, input.Owner.Raw) {
			return "", nil, 0, errors.New(fmt.Sprintf("the requestor does not own token"))
		}

		
		if tokenType == "" {
			tokenType = input.Type
		} else if tokenType != input.Type {
			return "", nil, 0, errors.New(fmt.Sprintf("two or more token types specified in input: '%s', '%s'", tokenType, input.Type))
		}

		
		quantity, err := ToQuantity(input.Quantity, Precision)
		if err != nil {
			return "", nil, 0, errors.Errorf("quantity in input [%s] is invalid, err '%s'", input.Quantity, err)
		}

		sum, err = sum.Add(quantity)
		if err != nil {
			return "", nil, 0, errors.Errorf("failed adding up quantities, err '%s'", err)
		}
		counter++

		if upperBound != nil {
			cmp, err := sum.Cmp(upperBound)
			if err != nil {
				return "", nil, 0, errors.Errorf("failed compering with upper bound, err '%s'", err)
			}
			if cmp >= 0 {
				break
			}
		}
	}

	return tokenType, sum, counter, nil
}



func (t *Transactor) ListTokens() (*token.UnspentTokens, error) {
	
	tokenOwner := &token.TokenOwner{Type: token.TokenOwner_MSP_IDENTIFIER, Raw: t.PublicCredential}
	ownerString, err := GetTokenOwnerString(tokenOwner)
	if err != nil {
		return nil, err
	}

	startKey, err := createCompositeKey(tokenKeyPrefix, []string{ownerString})
	if err != nil {
		return nil, err
	}
	endKey := startKey + string(maxUnicodeRuneValue)

	iterator, err := t.Ledger.GetStateRangeScanIterator(tokenNameSpace, startKey, endKey)
	if err != nil {
		return nil, err
	}

	tokens := make([]*token.UnspentToken, 0)
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

			output := &token.Token{}
			err = proto.Unmarshal(result.Value, output)
			if err != nil {
				return nil, errors.New("failed to retrieve unspent tokens: casting error")
			}

			
			verifierLogger.Debugf("adding token with ID '%s' to list of unspent tokens", result.GetKey())
			id, err := getTokenIdFromKey(result.Key)
			if err != nil {
				return nil, err
			}
			
			q, err := ToQuantity(output.Quantity, Precision)
			if err != nil {
				return nil, err
			}
			tokens = append(tokens,
				&token.UnspentToken{
					Type:     output.Type,
					Quantity: q.Decimal(),
					Id:       id,
				})
		}
	}
}


func (t *Transactor) RequestTokenOperation(tokenIDs []*token.TokenId, op *token.TokenOperation) (*token.TokenTransaction, int, error) {
	if len(tokenIDs) == 0 {
		return nil, 0, errors.New("no token ids in ExpectationRequest")
	}
	if op.GetAction() == nil {
		return nil, 0, errors.New("no action in request")
	}
	if op.GetAction().GetTransfer() == nil {
		return nil, 0, errors.New("no transfer in action")
	}
	if op.GetAction().GetTransfer().GetSender() == nil {
		return nil, 0, errors.New("no sender in transfer")
	}

	
	outputs := op.GetAction().GetTransfer().GetOutputs()
	outputType, outputSum, err := parseOutputs(outputs)
	if err != nil {
		return nil, 0, err
	}

	inputType, inputSum, count, err := t.getInputsFromTokenIds(tokenIDs, outputSum)
	if err != nil {
		return nil, 0, err
	}

	if outputType != inputType {
		return nil, 0, errors.Errorf("token type mismatch in inputs and outputs for token operation (%s vs %s)", outputType, inputType)
	}
	cmp, err := outputSum.Cmp(inputSum)
	if err != nil {
		return nil, 0, errors.Errorf("cannot compare quantities '%s'", err)
	}
	if cmp > 0 {
		return nil, 0, errors.Errorf("total quantity [%d] from TokenIds is less than total quantity [%d] in token operation", inputSum, outputSum)
	}

	
	cmp, err = inputSum.Cmp(outputSum)
	if err != nil {
		return nil, 0, errors.Errorf("cannot compare quantities '%s'", err)
	}
	if cmp > 0 {
		change, err := inputSum.Sub(outputSum)
		if err != nil {
			return nil, 0, errors.Errorf("failed computing change, err '%s'", err)
		}

		outputs = append(outputs, &token.Token{
			
			Owner:    &token.TokenOwner{Type: token.TokenOwner_MSP_IDENTIFIER, Raw: t.PublicCredential},
			Type:     outputType,
			Quantity: change.Hex(),
		})
	}

	return &token.TokenTransaction{
		Action: &token.TokenTransaction_TokenAction{
			TokenAction: &token.TokenAction{
				Data: &token.TokenAction_Transfer{
					Transfer: &token.Transfer{
						Inputs:  tokenIDs[:count],
						Outputs: outputs,
					},
				},
			},
		},
	}, count, nil
}


func (t *Transactor) Done() {
	if t.Ledger != nil {
		t.Ledger.Done()
	}
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
	
	if len(components) < numComponentsInKey+1 {
		return "", nil, errors.Errorf("invalid composite key - not enough components found in key '%s'", compositeKey)
	}
	return components[0], components[1:], nil
}


func parseOutputs(outputs []*token.Token) (string, Quantity, error) {
	if len(outputs) == 0 {
		return "", nil, errors.New("no outputs in request")
	}

	outputType := ""
	outputSum := NewZeroQuantity(Precision)
	for _, output := range outputs {
		if outputType == "" {
			outputType = output.GetType()
		} else if outputType != output.GetType() {
			return "", nil, errors.Errorf("multiple token types ('%s', '%s') in outputs", outputType, output.GetType())
		}
		quantity, err := ToQuantity(output.GetQuantity(), Precision)
		if err != nil {
			return "", nil, errors.Errorf("quantity in output [%s] is invalid, err '%s'", output.GetQuantity(), err)
		}

		outputSum, err = outputSum.Add(quantity)
		if err != nil {
			return "", nil, errors.Errorf("failed adding up quantities, err '%s'", err)
		}
	}

	return outputType, outputSum, nil
}

func getTokenIdFromKey(key string) (*token.TokenId, error) {
	_, components, err := splitCompositeKey(key)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error splitting input composite key: '%s'", err))
	}

	
	if len(components) != numComponentsInKey {
		return nil, errors.New(fmt.Sprintf("not enough components in output ID composite key; expected 4, received '%s'", components))
	}

	
	txID := components[numComponentsInKey-2]
	index, err := strconv.Atoi(components[numComponentsInKey-1])
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error parsing output index '%s': '%s'", components[numComponentsInKey-1], err))
	}
	return &token.TokenId{TxId: txID, Index: uint32(index)}, nil
}
