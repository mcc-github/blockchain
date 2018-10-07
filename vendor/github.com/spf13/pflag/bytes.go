package pflag

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
)


type bytesHexValue []byte


func (bytesHex bytesHexValue) String() string {
	return fmt.Sprintf("%X", []byte(bytesHex))
}


func (bytesHex *bytesHexValue) Set(value string) error {
	bin, err := hex.DecodeString(strings.TrimSpace(value))

	if err != nil {
		return err
	}

	*bytesHex = bin

	return nil
}


func (*bytesHexValue) Type() string {
	return "bytesHex"
}

func newBytesHexValue(val []byte, p *[]byte) *bytesHexValue {
	*p = val
	return (*bytesHexValue)(p)
}

func bytesHexConv(sval string) (interface{}, error) {

	bin, err := hex.DecodeString(sval)

	if err == nil {
		return bin, nil
	}

	return nil, fmt.Errorf("invalid string being converted to Bytes: %s %s", sval, err)
}


func (f *FlagSet) GetBytesHex(name string) ([]byte, error) {
	val, err := f.getFlagType(name, "bytesHex", bytesHexConv)

	if err != nil {
		return []byte{}, err
	}

	return val.([]byte), nil
}



func (f *FlagSet) BytesHexVar(p *[]byte, name string, value []byte, usage string) {
	f.VarP(newBytesHexValue(value, p), name, "", usage)
}


func (f *FlagSet) BytesHexVarP(p *[]byte, name, shorthand string, value []byte, usage string) {
	f.VarP(newBytesHexValue(value, p), name, shorthand, usage)
}



func BytesHexVar(p *[]byte, name string, value []byte, usage string) {
	CommandLine.VarP(newBytesHexValue(value, p), name, "", usage)
}


func BytesHexVarP(p *[]byte, name, shorthand string, value []byte, usage string) {
	CommandLine.VarP(newBytesHexValue(value, p), name, shorthand, usage)
}



func (f *FlagSet) BytesHex(name string, value []byte, usage string) *[]byte {
	p := new([]byte)
	f.BytesHexVarP(p, name, "", value, usage)
	return p
}


func (f *FlagSet) BytesHexP(name, shorthand string, value []byte, usage string) *[]byte {
	p := new([]byte)
	f.BytesHexVarP(p, name, shorthand, value, usage)
	return p
}



func BytesHex(name string, value []byte, usage string) *[]byte {
	return CommandLine.BytesHexP(name, "", value, usage)
}


func BytesHexP(name, shorthand string, value []byte, usage string) *[]byte {
	return CommandLine.BytesHexP(name, shorthand, value, usage)
}


type bytesBase64Value []byte


func (bytesBase64 bytesBase64Value) String() string {
	return base64.StdEncoding.EncodeToString([]byte(bytesBase64))
}


func (bytesBase64 *bytesBase64Value) Set(value string) error {
	bin, err := base64.StdEncoding.DecodeString(strings.TrimSpace(value))

	if err != nil {
		return err
	}

	*bytesBase64 = bin

	return nil
}


func (*bytesBase64Value) Type() string {
	return "bytesBase64"
}

func newBytesBase64Value(val []byte, p *[]byte) *bytesBase64Value {
	*p = val
	return (*bytesBase64Value)(p)
}

func bytesBase64ValueConv(sval string) (interface{}, error) {

	bin, err := base64.StdEncoding.DecodeString(sval)
	if err == nil {
		return bin, nil
	}

	return nil, fmt.Errorf("invalid string being converted to Bytes: %s %s", sval, err)
}


func (f *FlagSet) GetBytesBase64(name string) ([]byte, error) {
	val, err := f.getFlagType(name, "bytesBase64", bytesBase64ValueConv)

	if err != nil {
		return []byte{}, err
	}

	return val.([]byte), nil
}



func (f *FlagSet) BytesBase64Var(p *[]byte, name string, value []byte, usage string) {
	f.VarP(newBytesBase64Value(value, p), name, "", usage)
}


func (f *FlagSet) BytesBase64VarP(p *[]byte, name, shorthand string, value []byte, usage string) {
	f.VarP(newBytesBase64Value(value, p), name, shorthand, usage)
}



func BytesBase64Var(p *[]byte, name string, value []byte, usage string) {
	CommandLine.VarP(newBytesBase64Value(value, p), name, "", usage)
}


func BytesBase64VarP(p *[]byte, name, shorthand string, value []byte, usage string) {
	CommandLine.VarP(newBytesBase64Value(value, p), name, shorthand, usage)
}



func (f *FlagSet) BytesBase64(name string, value []byte, usage string) *[]byte {
	p := new([]byte)
	f.BytesBase64VarP(p, name, "", value, usage)
	return p
}


func (f *FlagSet) BytesBase64P(name, shorthand string, value []byte, usage string) *[]byte {
	p := new([]byte)
	f.BytesBase64VarP(p, name, shorthand, value, usage)
	return p
}



func BytesBase64(name string, value []byte, usage string) *[]byte {
	return CommandLine.BytesBase64P(name, "", value, usage)
}


func BytesBase64P(name, shorthand string, value []byte, usage string) *[]byte {
	return CommandLine.BytesBase64P(name, shorthand, value, usage)
}
