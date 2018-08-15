/*
Copyright IBM Corp. 2017 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package plain

import (
	"encoding/asn1"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToken(t *testing.T) {
	
	name := "alice"
	credBytes := []byte(name)
	tokenType := "wild_pineapple"
	quantity := uint64(100)
	
	aliceToken := &Token{Owner: credBytes, Type: tokenType, Quantity: quantity}
	assert.NotNil(t, aliceToken)
	
	tokenBytes, err := aliceToken.MarshalBinary()
	assert.NoError(t, err)
	assert.NotNil(t, tokenBytes)
	
	parsedToken := &Token{}
	err = parsedToken.UnmarshalBinary(tokenBytes)
	assert.NoError(t, err)
	assert.Equal(t, aliceToken, parsedToken)
	
	parsedTokenBytes, err := parsedToken.MarshalBinary()
	assert.NoError(t, err)
	assert.NotNil(t, parsedTokenBytes)
	assert.Equal(t, tokenBytes, parsedTokenBytes)
}

func TestParsingEmptyToken(t *testing.T) {
	
	tokenBytes := make([]byte, 0)
	parsedToken := &Token{}
	err := parsedToken.UnmarshalBinary(tokenBytes)
	assert.Error(t, err)
}

func TestParsingBadToken(t *testing.T) {
	
	tokenBytes, err := hex.DecodeString("deadbeef")
	assert.NoError(t, err)
	parsedToken := &Token{}
	err = parsedToken.UnmarshalBinary(tokenBytes)
	assert.Error(t, err)
}

func TestParsingTokenWithIncompatibleAsn1Struct(t *testing.T) {
	
	name := "wild_pineapple"
	otherStruct := struct {
		Name string
	}{
		Name: name,
	}
	tokenBytes, err := asn1.Marshal(otherStruct)
	assert.NoError(t, err)
	assert.NotNil(t, tokenBytes)
	parsedToken := &Token{}
	err = parsedToken.UnmarshalBinary(tokenBytes)
	assert.Error(t, err)
}

func TestParsingTokenWithIncorrectNumberOfEntries(t *testing.T) {
	
	data := raw{}
	data.Entries = make([][]byte, 4)
	tokenBytes, err := asn1.Marshal(data)
	assert.NoError(t, err)
	assert.NotNil(t, tokenBytes)
	parsedToken := &Token{}
	err = parsedToken.UnmarshalBinary(tokenBytes)
	assert.Error(t, err)
}

func TestSerializeNilToken(t *testing.T) {
	var nilToken *Token
	nilToken = nil
	bytes, err := nilToken.MarshalBinary()
	assert.Error(t, err)
	assert.Nil(t, bytes)
}
