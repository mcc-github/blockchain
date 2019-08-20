/*
Copyright IBM Corp. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package statedata

import (
	"fmt"
)


type KeyValue struct {
	Key   string
	Value []byte
}


type DataKey struct {
	Ns, Key string
}

func (d *DataKey) String() string {
	return fmt.Sprintf("Ns = %s, Key = %s", d.Ns, d.Key)
}


type PvtdataKeyHash struct {
	Ns, Coll string
	KeyHash  string
}

func (p *PvtdataKeyHash) String() string {
	return fmt.Sprintf("Ns = %s, Coll = %s, KeyHash = %x", p.Ns, p.Coll, p.KeyHash)
}




type ProposedWrites struct {
	Data        []*DataKey
	PvtdataHash []*PvtdataKeyHash
}


type ReadHint struct {
	Data        map[DataKey]*ReadHintDetails
	PvtdataHash map[PvtdataKeyHash]*ReadHintDetails
}


type ReadHintDetails struct {
	Value, Metadata bool
}



type WriteHint struct {
	Data        []*DataKey
	PvtdataHash []*PvtdataKeyHash
}
