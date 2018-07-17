package govaluate

import (
	"bytes"
)


type expressionOutputStream struct {
	transactions []string
}

func (this *expressionOutputStream) add(transaction string) {
	this.transactions = append(this.transactions, transaction)
}

func (this *expressionOutputStream) rollback() string {

	index := len(this.transactions) - 1
	ret := this.transactions[index]

	this.transactions = this.transactions[:index]
	return ret
}

func (this *expressionOutputStream) createString(delimiter string) string {

	var retBuffer bytes.Buffer
	var transaction string

	penultimate := len(this.transactions) - 1

	for i := 0; i < penultimate; i++ {

		transaction = this.transactions[i]

		retBuffer.WriteString(transaction)
		retBuffer.WriteString(delimiter)
	}
	retBuffer.WriteString(this.transactions[penultimate])

	return retBuffer.String()
}
