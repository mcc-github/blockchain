

package table

import (
	"fmt"
	"reflect"

	"github.com/onsi/ginkgo"
)


func DescribeTable(description string, itBody interface{}, entries ...TableEntry) bool {
	describeTable(description, itBody, entries, false, false)
	return true
}


func FDescribeTable(description string, itBody interface{}, entries ...TableEntry) bool {
	describeTable(description, itBody, entries, false, true)
	return true
}


func PDescribeTable(description string, itBody interface{}, entries ...TableEntry) bool {
	describeTable(description, itBody, entries, true, false)
	return true
}


func XDescribeTable(description string, itBody interface{}, entries ...TableEntry) bool {
	describeTable(description, itBody, entries, true, false)
	return true
}

func describeTable(description string, itBody interface{}, entries []TableEntry, pending bool, focused bool) {
	itBodyValue := reflect.ValueOf(itBody)
	if itBodyValue.Kind() != reflect.Func {
		panic(fmt.Sprintf("DescribeTable expects a function, got %#v", itBody))
	}

	if pending {
		ginkgo.PDescribe(description, func() {
			for _, entry := range entries {
				entry.generateIt(itBodyValue)
			}
		})
	} else if focused {
		ginkgo.FDescribe(description, func() {
			for _, entry := range entries {
				entry.generateIt(itBodyValue)
			}
		})
	} else {
		ginkgo.Describe(description, func() {
			for _, entry := range entries {
				entry.generateIt(itBodyValue)
			}
		})
	}
}
