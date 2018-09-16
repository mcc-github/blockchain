/*
Copyright IBM Corp. 2017 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

                 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cid

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/core/chaincode/shim/ext/attrmgr"
	"github.com/mcc-github/blockchain/protos/msp"
	"github.com/pkg/errors"
)



func GetID(stub ChaincodeStubInterface) (string, error) {
	c, err := New(stub)
	if err != nil {
		return "", err
	}
	return c.GetID()
}



func GetMSPID(stub ChaincodeStubInterface) (string, error) {
	c, err := New(stub)
	if err != nil {
		return "", err
	}
	return c.GetMSPID()
}


func GetAttributeValue(stub ChaincodeStubInterface, attrName string) (value string, found bool, err error) {
	c, err := New(stub)
	if err != nil {
		return "", false, err
	}
	return c.GetAttributeValue(attrName)
}


func AssertAttributeValue(stub ChaincodeStubInterface, attrName, attrValue string) error {
	c, err := New(stub)
	if err != nil {
		return err
	}
	return c.AssertAttributeValue(attrName, attrValue)
}



func GetX509Certificate(stub ChaincodeStubInterface) (*x509.Certificate, error) {
	c, err := New(stub)
	if err != nil {
		return nil, err
	}
	return c.GetX509Certificate()
}


type clientIdentityImpl struct {
	stub  ChaincodeStubInterface
	mspID string
	cert  *x509.Certificate
	attrs *attrmgr.Attributes
}


func New(stub ChaincodeStubInterface) (ClientIdentity, error) {
	c := &clientIdentityImpl{stub: stub}
	err := c.init()
	if err != nil {
		return nil, err
	}
	return c, nil
}


func (c *clientIdentityImpl) GetID() (string, error) {
	
	
	
	id := fmt.Sprintf("x509::%s::%s", getDN(&c.cert.Subject), getDN(&c.cert.Issuer))
	return base64.StdEncoding.EncodeToString([]byte(id)), nil
}



func (c *clientIdentityImpl) GetMSPID() (string, error) {
	return c.mspID, nil
}


func (c *clientIdentityImpl) GetAttributeValue(attrName string) (value string, found bool, err error) {
	if c.attrs == nil {
		return "", false, nil
	}
	return c.attrs.Value(attrName)
}


func (c *clientIdentityImpl) AssertAttributeValue(attrName, attrValue string) error {
	val, ok, err := c.GetAttributeValue(attrName)
	if err != nil {
		return err
	}
	if !ok {
		return errors.Errorf("Attribute '%s' was not found", attrName)
	}
	if val != attrValue {
		return errors.Errorf("Attribute '%s' equals '%s', not '%s'", attrName, val, attrValue)
	}
	return nil
}



func (c *clientIdentityImpl) GetX509Certificate() (*x509.Certificate, error) {
	return c.cert, nil
}


func (c *clientIdentityImpl) init() error {
	signingID, err := c.getIdentity()
	if err != nil {
		return err
	}
	c.mspID = signingID.GetMspid()
	idbytes := signingID.GetIdBytes()
	block, _ := pem.Decode(idbytes)
	if block == nil {
		err := c.getAttributesFromIdemix()
		if err != nil {
			return errors.WithMessage(err, "identity bytes are neither X509 PEM format nor an idemix credential")
		}
		return nil
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return errors.WithMessage(err, "failed to parse certificate")
	}
	c.cert = cert
	attrs, err := attrmgr.New().GetAttributesFromCert(cert)
	if err != nil {
		return errors.WithMessage(err, "failed to get attributes from the transaction invoker's certificate")
	}
	c.attrs = attrs
	return nil
}



func (c *clientIdentityImpl) getIdentity() (*msp.SerializedIdentity, error) {
	sid := &msp.SerializedIdentity{}
	creator, err := c.stub.GetCreator()
	if err != nil || creator == nil {
		return nil, errors.WithMessage(err, "failed to get transaction invoker's identity from the chaincode stub")
	}
	err = proto.Unmarshal(creator, sid)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal transaction invoker's identity")
	}
	return sid, nil
}

func (c *clientIdentityImpl) getAttributesFromIdemix() error {
	creator, err := c.stub.GetCreator()
	attrs, err := attrmgr.New().GetAttributesFromIdemix(creator)
	if err != nil {
		return errors.WithMessage(err, "failed to get attributes from the transaction invoker's idemix credential")
	}
	c.attrs = attrs
	return nil
}





func getDN(name *pkix.Name) string {
	r := name.ToRDNSequence()
	s := ""
	for i := 0; i < len(r); i++ {
		rdn := r[len(r)-1-i]
		if i > 0 {
			s += ","
		}
		for j, tv := range rdn {
			if j > 0 {
				s += "+"
			}
			typeString := tv.Type.String()
			typeName, ok := attributeTypeNames[typeString]
			if !ok {
				derBytes, err := asn1.Marshal(tv.Value)
				if err == nil {
					s += typeString + "=#" + hex.EncodeToString(derBytes)
					continue 
				}
				typeName = typeString
			}
			valueString := fmt.Sprint(tv.Value)
			escaped := ""
			begin := 0
			for idx, c := range valueString {
				if (idx == 0 && (c == ' ' || c == '#')) ||
					(idx == len(valueString)-1 && c == ' ') {
					escaped += valueString[begin:idx]
					escaped += "\\" + string(c)
					begin = idx + 1
					continue
				}
				switch c {
				case ',', '+', '"', '\\', '<', '>', ';':
					escaped += valueString[begin:idx]
					escaped += "\\" + string(c)
					begin = idx + 1
				}
			}
			escaped += valueString[begin:]
			s += typeName + "=" + escaped
		}
	}
	return s
}

var attributeTypeNames = map[string]string{
	"2.5.4.6":  "C",
	"2.5.4.10": "O",
	"2.5.4.11": "OU",
	"2.5.4.3":  "CN",
	"2.5.4.5":  "SERIALNUMBER",
	"2.5.4.7":  "L",
	"2.5.4.8":  "ST",
	"2.5.4.9":  "STREET",
	"2.5.4.17": "POSTALCODE",
}
