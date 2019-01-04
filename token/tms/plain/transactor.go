/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package plain

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/protos/ledger/queryresult"
	"github.com/mcc-github/blockchain/protos/token"
	"github.com/mcc-github/blockchain/token/ledger"
	"github.com/pkg/errors"
)


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
			return nil, "", 0, errors.New(fmt.Sprintf("error splitting input composite key: '%s'", err))
		}
		if namespace != tokenOutput {
			return nil, "", 0, errors.New(fmt.Sprintf("namespace not '%s': '%s'", tokenOutput, namespace))
		}
		if len(components) != 2 {
			return nil, "", 0, errors.New(fmt.Sprintf("not enough components in output ID composite key; expected 2, received '%s'", components))
		}
		txID := components[0]
		index, err := strconv.Atoi(components[1])
		if err != nil {
			return nil, "", 0, errors.New(fmt.Sprintf("error parsing output index '%s': '%s'", components[1], err))
		}

		
		inBytes, err := t.Ledger.GetState(tokenNameSpace, inKey)
		if err != nil {
			return nil, "", 0, err
		}
		if len(inBytes) == 0 {
			return nil, "", 0, errors.New(fmt.Sprintf("input '%s' does not exist", inKey))
		}
		input := &token.PlainOutput{}
		err = proto.Unmarshal(inBytes, input)
		if err != nil {
			return nil, "", 0, errors.New(fmt.Sprintf("error unmarshaling input bytes: '%s'", err))
		}

		
		if !bytes.Equal(t.PublicCredential, input.Owner) {
			return nil, "", 0, errors.New(fmt.Sprintf("the requestor does not own inputs"))
		}

		
		if tokenType == "" {
			tokenType = input.Type
		} else if tokenType != input.Type {
			return nil, "", 0, errors.New(fmt.Sprintf("two or more token types specified in input: '%s', '%s'", tokenType, input.Type))
		}
		
		inputs = append(inputs, &token.InputId{TxId: txID, Index: uint32(index)})

		
		quantitySum += input.Quantity
	}

	return inputs, tokenType, quantitySum, nil
}


func (t *Transactor) ListTokens() (*token.UnspentTokens, error) {
	iterator, err := t.Ledger.GetStateRangeScanIterator(tokenNameSpace, "", "")
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
	if len(request.GetTokenIds()) == 0 {
		return nil, errors.New("no token ids in ApproveAllowanceRequest")
	}

	if len(request.AllowanceShares) == 0 {
		return nil, errors.New("no recipient shares in ApproveAllowanceRequest")
	}

	var delegatedOutputs []*token.PlainDelegatedOutput

	inputs, tokenType, sumQuantity, err := t.getInputsFromTokenIds(request.GetTokenIds())
	if err != nil {
		return nil, err
	}

	

	delegatedQuantity := uint64(0)
	for _, share := range request.GetAllowanceShares() {
		if len(share.Recipient) == 0 {
			return nil, errors.Errorf("the recipient in approve must be specified")
		}
		if share.Quantity <= 0 {
			return nil, errors.Errorf("the quantity to approve [%d] must be greater than 0", share.GetQuantity())
		}
		delegatedOutputs = append(delegatedOutputs, &token.PlainDelegatedOutput{
			Owner:      []byte(request.Credential),
			Delegatees: [][]byte{share.Recipient},
			Type:       tokenType,
			Quantity:   share.Quantity,
		})
		delegatedQuantity = delegatedQuantity + share.Quantity
	}
	if sumQuantity < delegatedQuantity {
		return nil, errors.Errorf("insufficient funds: %v < %v", sumQuantity, delegatedQuantity)

	}
	var output *token.PlainOutput
	if sumQuantity != delegatedQuantity {
		output = &token.PlainOutput{
			Owner:    request.Credential,
			Type:     tokenType,
			Quantity: sumQuantity - delegatedQuantity,
		}
	}

	transaction := &token.TokenTransaction{
		Action: &token.TokenTransaction_PlainAction{
			PlainAction: &token.PlainTokenAction{
				Data: &token.PlainTokenAction_PlainApprove{
					PlainApprove: &token.PlainApprove{
						Inputs:           inputs,
						DelegatedOutputs: delegatedOutputs,
						Output:           output,
					},
				},
			},
		},
	}

	return transaction, nil
}


func (t *Transactor) RequestTransferFrom(request *token.TransferRequest) (*token.TokenTransaction, error) {
	if len(request.GetTokenIds()) == 0 {
		return nil, errors.New("no token ids in TransferFromRequest")
	}

	if len(request.Shares) == 0 {
		return nil, errors.New("no recipient shares in TransferFromRequest")
	}

	inputs, owner, tokenType, sumQuantity, err := t.getDelegateInputsFromTokenIds(request.GetTokenIds())

	if err != nil {
		return nil, err
	}

	outputs, transferQuantity, err := getOutputsForTx(request.GetShares(), tokenType, sumQuantity)
	if err != nil {
		return nil, err
	}

	var delegatedOutput *token.PlainDelegatedOutput
	if sumQuantity != transferQuantity {
		delegatedOutput = &token.PlainDelegatedOutput{
			Owner:      owner,
			Delegatees: [][]byte{t.PublicCredential},
			Type:       tokenType,
			Quantity:   sumQuantity - transferQuantity,
		}
	}

	transaction := &token.TokenTransaction{
		Action: &token.TokenTransaction_PlainAction{
			PlainAction: &token.PlainTokenAction{
				Data: &token.PlainTokenAction_PlainTransfer_From{
					PlainTransfer_From: &token.PlainTransferFrom{
						Inputs:          inputs,
						DelegatedOutput: delegatedOutput,
						Outputs:         outputs,
					},
				},
			},
		},
	}

	return transaction, nil
}

func (t *Transactor) getDelegateInputsFromTokenIds(tokenIds [][]byte) ([]*token.InputId, []byte, string, uint64, error) {
	var inputs []*token.InputId
	tokenType := ""
	var tokenOwner []byte
	sumQuantity := uint64(0)

	input := &token.PlainDelegatedOutput{}
	for _, inKeyBytes := range tokenIds {
		
		inKey := parseCompositeKeyBytes(inKeyBytes)

		
		namespace, components, err := splitCompositeKey(inKey)
		if err != nil {
			return nil, nil, "", 0, errors.Errorf("error splitting input composite key: '%s'", err)
		}
		if namespace != tokenDelegatedOutput {
			return nil, nil, "", 0, errors.Errorf("namespace not '%s': '%s'", tokenDelegatedOutput, namespace)
		}
		if len(components) != 2 {
			return nil, nil, "", 0, errors.Errorf("the number of components in output ID composite key is not correct; expected 2, received '%d'", len(components))
		}
		txID := components[0]
		index, err := strconv.Atoi(components[1])
		if err != nil {
			return nil, nil, "", 0, errors.Errorf("error parsing output index '%s': '%s'", components[1], err)
		}

		
		inBytes, err := t.Ledger.GetState(tokenNameSpace, inKey)
		if err != nil {
			return nil, nil, "", 0, err
		}
		if len(inBytes) == 0 {
			return nil, nil, "", 0, errors.Errorf("input '%s' does not exist", inKey)
		}
		err = proto.Unmarshal(inBytes, input)
		if err != nil {
			return nil, nil, "", 0, errors.Errorf("error unmarshaling input bytes: '%s'", err)
		}
		if len(input.Delegatees) != 1 {
			panic(fmt.Sprintf("the number of delegatees of input ID is not correct; expected 1, received '%d'", len(input.Delegatees)))
		}
		
		if tokenType == "" {
			tokenType = input.Type
		} else if tokenType != input.Type {
			return nil, nil, "", 0, errors.Errorf("two or more token types specified in input: '%s', '%s'", tokenType, input.Type)
		}

		
		if tokenOwner == nil {
			tokenOwner = input.Owner
		} else if !bytes.Equal(tokenOwner, input.Owner) {
			return nil, nil, "", 0, errors.Errorf("two or more token owners specified in input: '%s', '%s'", tokenOwner, input.Owner)
		}

		
		if !bytes.Equal(input.Delegatees[0], t.PublicCredential) {
			return nil, nil, "", 0, errors.Errorf("requestor is not allowed to transfer inputs")

		}
		
		inputs = append(inputs, &token.InputId{TxId: txID, Index: uint32(index)})
		sumQuantity = sumQuantity + input.Quantity
	}
	return inputs, tokenOwner, tokenType, sumQuantity, nil
}

func getOutputsForTx(shares []*token.RecipientTransferShare, tokenType string, inputQuantity uint64) ([]*token.PlainOutput, uint64, error) {
	var outputs []*token.PlainOutput

	
	outputQuantity := uint64(0)
	for _, share := range shares {
		if len(share.Recipient) == 0 {
			return nil, 0, errors.Errorf("the recipient in transferFrom must be specified")
		}
		if share.Quantity <= 0 {
			return nil, 0, errors.Errorf("the quantity to transferFrom [%d] must be greater than 0", share.GetQuantity())
		}
		outputs = append(outputs, &token.PlainOutput{
			Owner:    share.Recipient,
			Type:     tokenType,
			Quantity: share.Quantity,
		})
		outputQuantity = outputQuantity + share.Quantity
	}
	if inputQuantity < outputQuantity {
		return nil, 0, errors.Errorf("insufficient funds: %v < %v", inputQuantity, outputQuantity)

	}
	return outputs, outputQuantity, nil
}



func (t *Transactor) RequestExpectation(request *token.ExpectationRequest) (*token.TokenTransaction, error) {
	if len(request.GetTokenIds()) == 0 {
		return nil, errors.New("no token ids in ExpectationRequest")
	}
	if request.GetExpectation() == nil {
		return nil, errors.New("no token expectation in ExpectationRequest")
	}
	if request.GetExpectation().GetPlainExpectation() == nil {
		return nil, errors.New("no plain expectation in ExpectationRequest")
	}
	if request.GetExpectation().GetPlainExpectation().GetTransferExpectation() == nil {
		return nil, errors.New("no transfer expectation in ExpectationRequest")
	}

	inputs, inputType, inputSum, err := t.getInputsFromTokenIds(request.GetTokenIds())
	if err != nil {
		return nil, err
	}

	outputs := request.GetExpectation().GetPlainExpectation().GetTransferExpectation().GetOutputs()
	outputType, outputSum, err := parseOutputs(outputs)
	if err != nil {
		return nil, err
	}
	if outputType != inputType {
		return nil, errors.Errorf("token type mismatch in inputs and outputs for expectation (%s vs %s)", outputType, inputType)
	}
	if outputSum > inputSum {
		return nil, errors.Errorf("total quantity [%d] from TokenIds is less than total quantity [%d] in expectation", inputSum, outputSum)
	}

	
	if inputSum > outputSum {
		outputs = append(outputs, &token.PlainOutput{
			Owner:    t.PublicCredential, 
			Type:     outputType,
			Quantity: inputSum - outputSum,
		})
	}

	return &token.TokenTransaction{
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
	}, nil
}


func (t *Transactor) Done() {
	if t.Ledger != nil {
		t.Ledger.Done()
	}
}


func (t *Transactor) isSpent(outputID string) (bool, error) {
	key, err := createInputKey(outputID)
	if err != nil {
		return false, err
	}
	result, err := t.Ledger.GetState(tokenNameSpace, key)
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


func parseOutputs(outputs []*token.PlainOutput) (string, uint64, error) {
	if len(outputs) == 0 {
		return "", 0, errors.New("no outputs in request")
	}

	outputType := ""
	outputSum := uint64(0)
	for _, output := range outputs {
		if outputType == "" {
			outputType = output.GetType()
		} else if outputType != output.GetType() {
			return "", 0, errors.Errorf("multiple token types ('%s', '%s') in outputs", outputType, output.GetType())
		}
		outputSum += output.GetQuantity()
	}

	return outputType, outputSum, nil
}
