/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package plain

import (
	"bytes"
	"fmt"
	"strconv"
	"unicode/utf8"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/core/ledger/customtx"
	"github.com/mcc-github/blockchain/protos/token"
	"github.com/mcc-github/blockchain/protos/utils"
	"github.com/mcc-github/blockchain/token/identity"
	"github.com/mcc-github/blockchain/token/ledger"
	"github.com/pkg/errors"
)

const (
	minUnicodeRuneValue   = 0            
	maxUnicodeRuneValue   = utf8.MaxRune 
	compositeKeyNamespace = "\x00"
	tokenOutput           = "tokenOutput"
	tokenTx               = "tokenTx"
)


type Verifier struct {
	IssuingValidator identity.IssuingValidator
}



func (v *Verifier) ProcessTx(txID string, creator identity.PublicInfo, ttx *token.TokenTransaction, simulator ledger.LedgerWriter) error {
	err := v.checkProcess(txID, creator, ttx, simulator)
	if err != nil {
		return err
	}

	err = v.commitProcess(txID, creator, ttx, simulator)
	if err != nil {
		return err
	}
	return nil
}

func (v *Verifier) checkProcess(txID string, creator identity.PublicInfo, ttx *token.TokenTransaction, simulator ledger.LedgerReader) error {
	action := ttx.GetPlainAction()
	if action == nil {
		return &customtx.InvalidTxError{Msg: fmt.Sprintf("check process failed for transaction '%s': missing token action", txID)}
	}

	err := v.checkAction(creator, action, txID, simulator)
	if err != nil {
		return err
	}

	err = v.checkTxDoesNotExist(txID, simulator)
	if err != nil {
		return err
	}
	return nil
}

func (v *Verifier) checkAction(creator identity.PublicInfo, plainAction *token.PlainTokenAction, txID string, simulator ledger.LedgerReader) error {
	switch action := plainAction.Data.(type) {
	case *token.PlainTokenAction_PlainImport:
		return v.checkImportAction(creator, action.PlainImport, txID, simulator)
	case *token.PlainTokenAction_PlainTransfer:
		return v.checkTransferAction(creator, action.PlainTransfer, txID, simulator)
	default:
		return &customtx.InvalidTxError{Msg: fmt.Sprintf("unknown plain token action: %T", action)}
	}
}

func (v *Verifier) checkImportAction(creator identity.PublicInfo, importAction *token.PlainImport, txID string, simulator ledger.LedgerReader) error {
	err := v.checkImportOutputs(importAction.GetOutputs(), txID, simulator)
	if err != nil {
		return err
	}
	return v.checkImportPolicy(creator, txID, importAction)
}

func (v *Verifier) checkImportOutputs(outputs []*token.PlainOutput, txID string, simulator ledger.LedgerReader) error {
	if len(outputs) == 0 {
		return &customtx.InvalidTxError{Msg: fmt.Sprintf("no outputs in transaction: %s", txID)}
	}
	for i, output := range outputs {
		err := v.checkOutputDoesNotExist(i, txID, simulator)
		if err != nil {
			return err
		}

		if output.Quantity == 0 {
			return &customtx.InvalidTxError{Msg: fmt.Sprintf("output %d quantity is 0 in transaction: %s", i, txID)}
		}
	}
	return nil
}

func (v *Verifier) checkTransferAction(creator identity.PublicInfo, transferAction *token.PlainTransfer, txID string, simulator ledger.LedgerReader) error {
	outputType, outputSum, err := v.checkTransferOutputs(transferAction.GetOutputs(), txID, simulator)
	if err != nil {
		return err
	}
	inputType, inputSum, err := v.checkTransferInputs(creator, transferAction.GetInputs(), txID, simulator)
	if err != nil {
		return err
	}
	if outputType != inputType {
		return &customtx.InvalidTxError{Msg: fmt.Sprintf("token type mismatch in inputs and outputs for transfer with ID %s (%s vs %s)", txID, outputType, inputType)}
	}
	if outputSum != inputSum {
		return &customtx.InvalidTxError{Msg: fmt.Sprintf("token sum mismatch in inputs and outputs for transfer with ID %s (%d vs %d)", txID, outputSum, inputSum)}
	}
	return nil
}

func (v *Verifier) checkOutputDoesNotExist(index int, txID string, simulator ledger.LedgerReader) error {
	outputID, err := createOutputKey(txID, index)
	if err != nil {
		return &customtx.InvalidTxError{Msg: fmt.Sprintf("error creating output ID: %s", err)}
	}

	existingOutputBytes, err := simulator.GetState(tokenNamespace, outputID)
	if err != nil {
		return err
	}

	if existingOutputBytes != nil {
		return &customtx.InvalidTxError{Msg: fmt.Sprintf("output already exists: %s", outputID)}
	}
	return nil
}

func (v *Verifier) checkTransferOutputs(outputs []*token.PlainOutput, txID string, simulator ledger.LedgerReader) (string, uint64, error) {
	tokenType := ""
	tokenSum := uint64(0)
	for i, output := range outputs {
		err := v.checkOutputDoesNotExist(i, txID, simulator)
		if err != nil {
			return "", 0, err
		}
		if tokenType == "" {
			tokenType = output.GetType()
		} else if tokenType != output.GetType() {
			return "", 0, &customtx.InvalidTxError{Msg: fmt.Sprintf("multiple token types ('%s', '%s') in transfer output for txID '%s'", tokenType, output.GetType(), txID)}
		}
		tokenSum += output.GetQuantity()
	}
	return tokenType, tokenSum, nil
}

func (v *Verifier) checkTransferInputs(creator identity.PublicInfo, inputIDs []*token.InputId, txID string, simulator ledger.LedgerReader) (string, uint64, error) {
	tokenType := ""
	inputSum := uint64(0)
	processedIDs := make(map[string]bool)
	for _, id := range inputIDs {
		inputKey, err := createOutputKey(id.TxId, int(id.Index))
		if err != nil {
			return "", 0, &customtx.InvalidTxError{Msg: fmt.Sprintf("error creating output ID for transfer input: %s", err)}
		}
		input, err := v.getOutput(inputKey, simulator)
		if err != nil {
			return "", 0, err
		}
		err = v.checkInputOwner(creator, input, inputKey)
		if err != nil {
			return "", 0, err
		}
		if tokenType == "" {
			tokenType = input.GetType()
		} else if tokenType != input.GetType() {
			return "", 0, &customtx.InvalidTxError{Msg: fmt.Sprintf("multiple token types in transfer input for txID: %s (%s, %s)", txID, tokenType, input.GetType())}
		}
		if processedIDs[inputKey] {
			return "", 0, &customtx.InvalidTxError{Msg: fmt.Sprintf("token input '%s' spent more than once in single transfer with txID '%s'", inputKey, txID)}
		}
		processedIDs[inputKey] = true
		inputSum += input.GetQuantity()
		spentKey, err := createSpentKey(id.TxId, int(id.Index))
		if err != nil {
			return "", 0, err
		}
		spent, err := v.isSpent(spentKey, simulator)
		if err != nil {
			return "", 0, err
		}
		if spent {
			return "", 0, &customtx.InvalidTxError{Msg: fmt.Sprintf("input with ID %s for transfer has already been spent", inputKey)}
		}
	}
	return tokenType, inputSum, nil
}

func (v *Verifier) checkInputOwner(creator identity.PublicInfo, input *token.PlainOutput, inputID string) error {
	if !bytes.Equal(creator.Public(), input.Owner) {
		return &customtx.InvalidTxError{Msg: fmt.Sprintf("transfer input with ID %s not owned by creator", inputID)}
	}
	return nil
}

func (v *Verifier) checkTxDoesNotExist(txID string, simulator ledger.LedgerReader) error {
	txKey, err := createTxKey(txID)
	if err != nil {
		return &customtx.InvalidTxError{Msg: fmt.Sprintf("error creating txID: %s", err)}
	}

	existingTx, err := simulator.GetState(tokenNamespace, txKey)
	if err != nil {
		return err
	}

	if existingTx != nil {
		return &customtx.InvalidTxError{Msg: fmt.Sprintf("transaction already exists: %s", txID)}
	}
	return nil
}

func (v *Verifier) checkImportPolicy(creator identity.PublicInfo, txID string, importData *token.PlainImport) error {
	for _, output := range importData.Outputs {
		err := v.IssuingValidator.Validate(creator, output.Type)
		if err != nil {
			return &customtx.InvalidTxError{Msg: fmt.Sprintf("import policy check failed: %s", err)}
		}
	}
	return nil
}


const tokenNamespace = "tms"

func (v *Verifier) commitProcess(txID string, creator identity.PublicInfo, ttx *token.TokenTransaction, simulator ledger.LedgerWriter) error {
	err := v.commitAction(ttx.GetPlainAction(), txID, simulator)
	if err != nil {
		return err
	}

	err = v.addTransaction(txID, ttx, simulator)
	if err != nil {
		return err
	}

	return nil
}

func (v *Verifier) commitAction(plainAction *token.PlainTokenAction, txID string, simulator ledger.LedgerWriter) (err error) {
	switch action := plainAction.Data.(type) {
	case *token.PlainTokenAction_PlainImport:
		err = v.commitImportAction(action.PlainImport, txID, simulator)
	case *token.PlainTokenAction_PlainTransfer:
		err = v.commitTransferAction(action.PlainTransfer, txID, simulator)
	}
	return
}

func (v *Verifier) commitImportAction(importAction *token.PlainImport, txID string, simulator ledger.LedgerWriter) error {
	for i, output := range importAction.GetOutputs() {
		outputID, err := createOutputKey(txID, i)
		if err != nil {
			return &customtx.InvalidTxError{Msg: fmt.Sprintf("error creating output ID: %s", err)}
		}

		err = v.addOutput(outputID, output, simulator)
		if err != nil {
			return err
		}
	}
	return nil
}

func (v *Verifier) commitTransferAction(transferAction *token.PlainTransfer, txID string, simulator ledger.LedgerWriter) error {
	for i, output := range transferAction.GetOutputs() {
		outputID, err := createOutputKey(txID, i)
		if err != nil {
			return &customtx.InvalidTxError{Msg: fmt.Sprintf("error creating output ID: %s", err)}
		}

		err = v.addOutput(outputID, output, simulator)
		if err != nil {
			return err
		}
	}
	return v.markInputsSpent(txID, transferAction.GetInputs(), simulator)
}

func (v *Verifier) addOutput(outputID string, output *token.PlainOutput, simulator ledger.LedgerWriter) error {
	outputBytes := utils.MarshalOrPanic(output)

	return simulator.SetState(tokenNamespace, outputID, outputBytes)
}

func (v *Verifier) addTransaction(txID string, ttx *token.TokenTransaction, simulator ledger.LedgerWriter) error {
	ttxBytes := utils.MarshalOrPanic(ttx)

	ttxID, err := createTxKey(txID)
	if err != nil {
		return &customtx.InvalidTxError{Msg: fmt.Sprintf("error creating txID: %s", err)}
	}

	return simulator.SetState(tokenNamespace, ttxID, ttxBytes)
}

var TokenInputSpentMarker = []byte{1}

func (v *Verifier) markInputsSpent(txID string, inputs []*token.InputId, simulator ledger.LedgerWriter) error {
	for _, id := range inputs {
		inputID, err := createSpentKey(id.TxId, int(id.Index))
		if err != nil {
			return &customtx.InvalidTxError{Msg: fmt.Sprintf("error creating spent key: %s", err)}
		}

		err = simulator.SetState(tokenNamespace, inputID, TokenInputSpentMarker)
		if err != nil {
			return err
		}
	}
	return nil
}

func (v *Verifier) getOutput(outputID string, simulator ledger.LedgerReader) (*token.PlainOutput, error) {
	outputBytes, err := simulator.GetState(tokenNamespace, outputID)
	if err != nil {
		return nil, err
	}
	if outputBytes == nil {
		return nil, &customtx.InvalidTxError{Msg: fmt.Sprintf("input with ID %s for transfer does not exist", outputID)}
	}
	output := &token.PlainOutput{}
	err = proto.Unmarshal(outputBytes, output)
	return output, err
}


func (v *Verifier) isSpent(spentKey string, simulator ledger.LedgerReader) (bool, error) {
	result, err := simulator.GetState(namespace, spentKey)
	return result != nil, err
}



func createOutputKey(txID string, index int) (string, error) {
	return createCompositeKey(tokenOutput, []string{txID, strconv.Itoa(index)})
}


func createTxKey(txID string) (string, error) {
	return createCompositeKey(tokenTx, []string{txID})
}



func createSpentKey(txID string, index int) (string, error) {
	return createCompositeKey("tokenInput", []string{txID, strconv.Itoa(index)})
}


func createCompositeKey(objectType string, attributes []string) (string, error) {
	if err := validateCompositeKeyAttribute(objectType); err != nil {
		return "", err
	}
	ck := compositeKeyNamespace + objectType + string(minUnicodeRuneValue)
	for _, att := range attributes {
		if err := validateCompositeKeyAttribute(att); err != nil {
			return "", err
		}
		ck += att + string(minUnicodeRuneValue)
	}
	return ck, nil
}

func validateCompositeKeyAttribute(str string) error {
	if !utf8.ValidString(str) {
		return errors.Errorf("not a valid utf8 string: [%x]", str)
	}
	for index, runeValue := range str {
		if runeValue == minUnicodeRuneValue || runeValue == maxUnicodeRuneValue {
			return errors.Errorf(`input contain unicode %#U starting at position [%d]. %#U and %#U are not allowed in the input attribute of a composite key`,
				runeValue, index, minUnicodeRuneValue, maxUnicodeRuneValue)
		}
	}
	return nil
}

func parseCompositeKeyBytes(keyBytes []byte) string {
	return string(keyBytes)
}

func getCompositeKeyBytes(compositeKey string) []byte {
	return []byte(compositeKey)
}
