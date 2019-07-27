/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package msp

import (
	"github.com/pkg/errors"
)

type MSPVersion int

const (
	MSPv1_0 = iota
	MSPv1_1
	MSPv1_3
	MSPv1_4_2
)


type NewOpts interface {
	
	GetVersion() MSPVersion
}


type NewBaseOpts struct {
	Version MSPVersion
}

func (o *NewBaseOpts) GetVersion() MSPVersion {
	return o.Version
}


type BCCSPNewOpts struct {
	NewBaseOpts
}


type IdemixNewOpts struct {
	NewBaseOpts
}


func New(opts NewOpts) (MSP, error) {
	switch opts.(type) {
	case *BCCSPNewOpts:
		switch opts.GetVersion() {
		case MSPv1_0:
			return newBccspMsp(MSPv1_0)
		case MSPv1_1:
			return newBccspMsp(MSPv1_1)
		case MSPv1_3:
			return newBccspMsp(MSPv1_3)
		case MSPv1_4_2:
			return newBccspMsp(MSPv1_4_2)
		default:
			return nil, errors.Errorf("Invalid *BCCSPNewOpts. Version not recognized [%v]", opts.GetVersion())
		}
	case *IdemixNewOpts:
		switch opts.GetVersion() {
		case MSPv1_4_2:
			fallthrough
		case MSPv1_3:
			return newIdemixMsp(MSPv1_3)
		case MSPv1_1:
			return newIdemixMsp(MSPv1_1)
		default:
			return nil, errors.Errorf("Invalid *IdemixNewOpts. Version not recognized [%v]", opts.GetVersion())
		}
	default:
		return nil, errors.Errorf("Invalid msp.NewOpts instance. It must be either *BCCSPNewOpts or *IdemixNewOpts. It was [%v]", opts)
	}
}
