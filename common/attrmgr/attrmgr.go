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



package attrmgr

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

var (
	
	
	AttrOID = asn1.ObjectIdentifier{1, 2, 3, 4, 5, 6, 7, 8, 1}
	
	AttrOIDString = "1.2.3.4.5.6.7.8.1"
)


type Attribute interface {
	
	GetName() string
	
	GetValue() string
}


type AttributeRequest interface {
	
	GetName() string
	
	IsRequired() bool
}


func New() *Mgr { return &Mgr{} }


type Mgr struct{}



func (mgr *Mgr) ProcessAttributeRequestsForCert(requests []AttributeRequest, attributes []Attribute, cert *x509.Certificate) error {
	attrs, err := mgr.ProcessAttributeRequests(requests, attributes)
	if err != nil {
		return err
	}
	return mgr.AddAttributesToCert(attrs, cert)
}



func (mgr *Mgr) ProcessAttributeRequests(requests []AttributeRequest, attributes []Attribute) (*Attributes, error) {
	attrsMap := map[string]string{}
	attrs := &Attributes{Attrs: attrsMap}
	missingRequiredAttrs := []string{}
	
	for _, req := range requests {
		
		name := req.GetName()
		attr := getAttrByName(name, attributes)
		if attr == nil {
			if req.IsRequired() {
				
				missingRequiredAttrs = append(missingRequiredAttrs, name)
			}
			
			continue
		}
		attrsMap[name] = attr.GetValue()
	}
	if len(missingRequiredAttrs) > 0 {
		return nil, errors.Errorf("The following required attributes are missing: %+v",
			missingRequiredAttrs)
	}
	return attrs, nil
}


func (mgr *Mgr) AddAttributesToCert(attrs *Attributes, cert *x509.Certificate) error {
	buf, err := json.Marshal(attrs)
	if err != nil {
		return errors.Wrap(err, "Failed to marshal attributes")
	}
	ext := pkix.Extension{
		Id:       AttrOID,
		Critical: false,
		Value:    buf,
	}
	cert.Extensions = append(cert.Extensions, ext)
	return nil
}


func (mgr *Mgr) GetAttributesFromCert(cert *x509.Certificate) (*Attributes, error) {
	
	buf, err := getAttributesFromCert(cert)
	if err != nil {
		return nil, err
	}
	
	attrs := &Attributes{}
	if buf != nil {
		err := json.Unmarshal(buf, attrs)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to unmarshal attributes from certificate")
		}
	}
	return attrs, nil
}


type Attributes struct {
	Attrs map[string]string `json:"attrs"`
}


func (a *Attributes) Names() []string {
	i := 0
	names := make([]string, len(a.Attrs))
	for name := range a.Attrs {
		names[i] = name
		i++
	}
	return names
}


func (a *Attributes) Contains(name string) bool {
	_, ok := a.Attrs[name]
	return ok
}


func (a *Attributes) Value(name string) (string, bool, error) {
	attr, ok := a.Attrs[name]
	return attr, ok, nil
}



func (a *Attributes) True(name string) error {
	val, ok, err := a.Value(name)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("Attribute '%s' was not found", name)
	}
	if val != "true" {
		return fmt.Errorf("Attribute '%s' is not true", name)
	}
	return nil
}


func getAttributesFromCert(cert *x509.Certificate) ([]byte, error) {
	for _, ext := range cert.Extensions {
		if isAttrOID(ext.Id) {
			return ext.Value, nil
		}
	}
	return nil, nil
}


func isAttrOID(oid asn1.ObjectIdentifier) bool {
	if len(oid) != len(AttrOID) {
		return false
	}
	for idx, val := range oid {
		if val != AttrOID[idx] {
			return false
		}
	}
	return true
}


func getAttrByName(name string, attrs []Attribute) Attribute {
	for _, attr := range attrs {
		if attr.GetName() == name {
			return attr
		}
	}
	return nil
}
