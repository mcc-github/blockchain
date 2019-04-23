/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package plain

import (
	"math/big"

	"github.com/pkg/errors"
)


type Quantity interface {

	
	
	Add(b Quantity) (Quantity, error)

	
	
	Sub(b Quantity) (Quantity, error)

	
	Cmp(b Quantity) (int, error)

	
	Hex() string

	
	Decimal() string
}

type BigQuantity struct {
	*big.Int
	Precision uint64
}




func ToQuantity(q string, precision uint64) (Quantity, error) {
	v, success := big.NewInt(0).SetString(q, 0)
	if !success {
		return nil, errors.New("invalid input")
	}
	if v.Cmp(big.NewInt(0)) <= 0 {
		return nil, errors.New("quantity must be larger than 0")
	}
	if precision == 0 {
		return nil, errors.New("precision be larger than 0")
	}
	if v.BitLen() > int(precision) {
		return nil, errors.Errorf("%s has precision %d > %d", q, v.BitLen(), precision)
	}

	return &BigQuantity{Int: v, Precision: precision}, nil
}



func NewZeroQuantity(precision uint64) Quantity {
	b := BigQuantity{Int: big.NewInt(0), Precision: precision}
	return &b
}

func (q *BigQuantity) Add(b Quantity) (Quantity, error) {
	bq, ok := b.(*BigQuantity)
	if !ok {
		return nil, errors.Errorf("expected uint64 quantity, got '%t", b)
	}

	
	sum := big.NewInt(0)
	sum = sum.Add(q.Int, bq.Int)

	if sum.BitLen() > int(q.Precision) {
		return nil, errors.Errorf("%d + %d = overflow", q, b)
	}

	sumq := BigQuantity{Int: sum, Precision: q.Precision}
	return &sumq, nil
}

func (q *BigQuantity) Sub(b Quantity) (Quantity, error) {
	bq, ok := b.(*BigQuantity)
	if !ok {
		return nil, errors.Errorf("expected uint64 quantity, got '%t", b)
	}

	
	if q.Int.Cmp(bq.Int) < 0 {
		return nil, errors.Errorf("%d < %d", q, b)
	}
	diff := big.NewInt(0)
	diff.Sub(q.Int, b.(*BigQuantity).Int)

	diffq := BigQuantity{Int: diff, Precision: q.Precision}
	return &diffq, nil
}

func (q *BigQuantity) Cmp(b Quantity) (int, error) {
	bq, ok := b.(*BigQuantity)
	if !ok {
		return 0, errors.Errorf("expected uint64 quantity, got '%t", b)
	}

	return q.Int.Cmp(bq.Int), nil
}

func (q *BigQuantity) Hex() string {
	return "0x" + q.Int.Text(16)
}

func (q *BigQuantity) Decimal() string {
	return q.Int.Text(10)
}
