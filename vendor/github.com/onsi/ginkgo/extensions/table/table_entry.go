package table

import (
	"reflect"

	"github.com/onsi/ginkgo"
)


type TableEntry struct {
	Description string
	Parameters  []interface{}
	Pending     bool
	Focused     bool
}

func (t TableEntry) generateIt(itBody reflect.Value) {
	if t.Pending {
		ginkgo.PIt(t.Description)
		return
	}

	values := []reflect.Value{}
	for i, param := range t.Parameters {
		var value reflect.Value

		if param == nil {
			inType := itBody.Type().In(i)
			value = reflect.Zero(inType)
		} else {
			value = reflect.ValueOf(param)
		}

		values = append(values, value)
	}

	body := func() {
		itBody.Call(values)
	}

	if t.Focused {
		ginkgo.FIt(t.Description, body)
	} else {
		ginkgo.It(t.Description, body)
	}
}


func Entry(description string, parameters ...interface{}) TableEntry {
	return TableEntry{description, parameters, false, false}
}


func FEntry(description string, parameters ...interface{}) TableEntry {
	return TableEntry{description, parameters, false, true}
}


func PEntry(description string, parameters ...interface{}) TableEntry {
	return TableEntry{description, parameters, true, false}
}


func XEntry(description string, parameters ...interface{}) TableEntry {
	return TableEntry{description, parameters, true, false}
}
