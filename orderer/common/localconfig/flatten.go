/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package localconfig

import (
	"fmt"
	"reflect"
)



func Flatten(i interface{}) []string {
	var res []string
	flatten("", &res, reflect.ValueOf(i))
	return res
}



















func flatten(k string, m *[]string, v reflect.Value) {
	delimiter := "."
	if k == "" {
		delimiter = ""
	}

	switch v.Kind() {
	case reflect.Ptr:
		if v.IsNil() {
			*m = append(*m, fmt.Sprintf("%s =", k))
			return
		}
		flatten(k, m, v.Elem())
	case reflect.Struct:
		if x, ok := v.Interface().(fmt.Stringer); ok {
			*m = append(*m, fmt.Sprintf("%s = %v", k, x))
			return
		}

		for i := 0; i < v.NumField(); i++ {
			flatten(k+delimiter+v.Type().Field(i).Name, m, v.Field(i))
		}
	case reflect.String:
		
		*m = append(*m, fmt.Sprintf("%s = \"%s\"", k, v))
	default:
		*m = append(*m, fmt.Sprintf("%s = %v", k, v))
	}
}
