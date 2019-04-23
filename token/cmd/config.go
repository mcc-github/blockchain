/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package token

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/mcc-github/blockchain/bccsp/sw"
	"github.com/mcc-github/blockchain/msp"
	"github.com/mcc-github/blockchain/protos/token"
	"github.com/mcc-github/blockchain/token/client"
	"github.com/pkg/errors"
)


type ShellRecipientShare struct {
	Recipient string
	Quantity  string
}




func LoadConfig(s string) (*client.ClientConfig, error) {
	
	config := &client.ClientConfig{}
	jsonDecoder := json.NewDecoder(strings.NewReader(s))
	err := jsonDecoder.Decode(config)
	if err == nil {
		return config, nil
	}

	
	file, err := os.Open(s)
	if err != nil {
		return nil, err
	}

	
	config = &client.ClientConfig{}
	jsonDecoder = json.NewDecoder(file)
	err = jsonDecoder.Decode(config)
	if err == nil {
		return config, nil
	}

	return nil, errors.New("cannot load configuration, not a json, nor a file containing a json configuration")
}


func GetSigningIdentity(mspConfigPath, mspID, mspType string) (msp.SigningIdentity, error) {
	mspInstance, err := LoadLocalMSPAt(mspConfigPath, mspID, mspType)
	if err != nil {
		return nil, err
	}

	signingIdentity, err := mspInstance.GetDefaultSigningIdentity()
	if err != nil {
		return nil, err
	}

	return signingIdentity, nil
}





func LoadTokenOwner(s string) (*token.TokenOwner, error) {
	res, err := LoadLocalMspRecipient(s)
	if err == nil {
		return res, nil
	}

	raw, err := LoadSerialisedRecipient(s)
	if err == nil {
		return &token.TokenOwner{Raw: raw}, nil
	}

	return &token.TokenOwner{Raw: []byte(s)}, nil
}



func LoadLocalMspRecipient(s string) (*token.TokenOwner, error) {
	strs := strings.Split(s, ":")
	if len(strs) < 2 {
		return nil, errors.Errorf("invalid input '%s', expected <msp_id>:<path>", s)
	}

	mspID := strs[0]
	mspPath := strings.TrimPrefix(s, mspID+":")
	localMSP, err := LoadLocalMSPAt(mspPath, mspID, "bccsp")
	if err != nil {
		return nil, err
	}

	signer, err := localMSP.GetDefaultSigningIdentity()
	if err != nil {
		return nil, err
	}

	raw, err := signer.Serialize()
	if err != nil {
		return nil, err
	}

	return &token.TokenOwner{Raw: raw}, nil
}



func LoadLocalMSPAt(dir, id, mspType string) (msp.MSP, error) {
	if mspType != "bccsp" {
		return nil, errors.Errorf("invalid msp type, expected 'bccsp', got %s", mspType)
	}
	conf, err := msp.GetLocalMspConfig(dir, nil, id)
	if err != nil {
		return nil, err
	}
	ks, err := sw.NewFileBasedKeyStore(nil, filepath.Join(dir, "keystore"), true)
	if err != nil {
		return nil, err
	}
	thisMSP, err := msp.NewBccspMspWithKeyStore(msp.MSPv1_0, ks)
	if err != nil {
		return nil, err
	}
	err = thisMSP.Setup(conf)
	if err != nil {
		return nil, err
	}
	return thisMSP, nil
}


func LoadSerialisedRecipient(serialisedRecipientPath string) ([]byte, error) {
	fileCont, err := ioutil.ReadFile(serialisedRecipientPath)
	if err != nil {
		return nil, errors.Wrapf(err, "could not read file %s", serialisedRecipientPath)
	}

	return fileCont, nil
}




func LoadTokenIDs(s string) ([]*token.TokenId, error) {
	
	
	
	res, err := LoadTokenIDsFromJson(s)
	if err == nil {
		return res, nil
	}

	return LoadTokenIDsFromFile(s)
}


func LoadTokenIDsFromJson(s string) ([]*token.TokenId, error) {
	var tokenIDs []*token.TokenId
	err := json.Unmarshal([]byte(s), &tokenIDs)
	if err != nil {
		return nil, errors.Wrap(err, "failed unmarshalling json")
	}

	return tokenIDs, nil
}



func LoadTokenIDsFromFile(s string) ([]*token.TokenId, error) {
	fileCont, err := ioutil.ReadFile(s)
	if err != nil {
		return nil, errors.Wrapf(err, "could not read file %s", s)
	}

	return LoadTokenIDsFromJson(string(fileCont))
}




func LoadShares(s string) ([]*token.RecipientShare, error) {
	
	
	
	var err1, err2 error

	res, err1 := LoadSharesFromJson(s)
	if err1 != nil {
		res, err2 = LoadSharesFromFile(s)
	}

	if err1 == nil || err2 == nil {
		return SubstituteShareRecipient(res)
	}
	return nil, errors.Errorf("failed loading shares [%s][%s]", err1, err2)
}


func LoadSharesFromJson(s string) ([]*ShellRecipientShare, error) {
	var shares []*ShellRecipientShare
	err := json.Unmarshal([]byte(s), &shares)
	if err != nil {
		return nil, errors.Wrap(err, "failed unmarshalling json")
	}

	return shares, nil
}


func LoadSharesFromFile(s string) ([]*ShellRecipientShare, error) {
	fileCont, err := ioutil.ReadFile(s)
	if err != nil {
		return nil, errors.Wrapf(err, "could not read file %s", s)
	}

	return LoadSharesFromJson(string(fileCont))
}




func SubstituteShareRecipient(shares []*ShellRecipientShare) ([]*token.RecipientShare, error) {
	if len(shares) == 0 {
		return nil, errors.New("SubstituteShareRecipient: empty input passed")
	}

	var outputs []*token.RecipientShare
	for _, share := range shares {
		if share == nil {
			continue
		}
		if len(share.Recipient) == 0 {
			return nil, errors.New("SubstituteShareRecipient: invalid recipient share")
		}
		recipient, err := LoadTokenOwner(share.Recipient)
		if err != nil {
			return nil, errors.Wrap(err, "SubstituteShareRecipient: failed loading token owner")
		}
		outputs = append(outputs, &token.RecipientShare{Recipient: recipient, Quantity: share.Quantity})
	}

	return outputs, nil
}


type JsonLoader struct {
}

func (*JsonLoader) TokenOwner(s string) (*token.TokenOwner, error) {
	return LoadTokenOwner(s)
}

func (*JsonLoader) TokenIDs(s string) ([]*token.TokenId, error) {
	return LoadTokenIDs(s)
}

func (*JsonLoader) Shares(s string) ([]*token.RecipientShare, error) {
	return LoadShares(s)
}
