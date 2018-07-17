package govaluate

import (
	"errors"
)


type Parameters interface {

	
	Get(name string) (interface{}, error)
}

type MapParameters map[string]interface{}

func (p MapParameters) Get(name string) (interface{}, error) {

	value, found := p[name]

	if !found {
		errorMessage := "No parameter '" + name + "' found."
		return nil, errors.New(errorMessage)
	}

	return value, nil
}
