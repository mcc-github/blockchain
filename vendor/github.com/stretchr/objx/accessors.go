package objx

import (
	"regexp"
	"strconv"
	"strings"
)



const arrayAccesRegexString = `^(.+)\[([0-9]+)\]$`


var arrayAccesRegex = regexp.MustCompile(arrayAccesRegexString)














func (m Map) Get(selector string) *Value {
	rawObj := access(m, selector, nil, false)
	return &Value{data: rawObj}
}











func (m Map) Set(selector string, value interface{}) Map {
	access(m, selector, value, true)
	return m
}



func access(current, selector, value interface{}, isSet bool) interface{} {
	switch selector.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		if array, ok := current.([]interface{}); ok {
			index := intFromInterface(selector)
			if index >= len(array) {
				return nil
			}
			return array[index]
		}
		return nil

	case string:
		selStr := selector.(string)
		selSegs := strings.SplitN(selStr, PathSeparator, 2)
		thisSel := selSegs[0]
		index := -1
		var err error

		if strings.Contains(thisSel, "[") {
			arrayMatches := arrayAccesRegex.FindStringSubmatch(thisSel)
			if len(arrayMatches) > 0 {
				
				thisSel = arrayMatches[1]

				
				index, err = strconv.Atoi(arrayMatches[2])

				if err != nil {
					
					
					panic("objx: Array index is not an integer.  Must use array[int].")
				}
			}
		}
		if curMap, ok := current.(Map); ok {
			current = map[string]interface{}(curMap)
		}
		
		switch current.(type) {
		case map[string]interface{}:
			curMSI := current.(map[string]interface{})
			if len(selSegs) <= 1 && isSet {
				curMSI[thisSel] = value
				return nil
			}
			current = curMSI[thisSel]
		default:
			current = nil
		}
		
		if index > -1 {
			if array, ok := current.([]interface{}); ok {
				if index < len(array) {
					current = array[index]
				} else {
					current = nil
				}
			}
		}
		if len(selSegs) > 1 {
			current = access(current, selSegs[1], value, isSet)
		}
	}
	return current
}




func intFromInterface(selector interface{}) int {
	var value int
	switch selector.(type) {
	case int:
		value = selector.(int)
	case int8:
		value = int(selector.(int8))
	case int16:
		value = int(selector.(int16))
	case int32:
		value = int(selector.(int32))
	case int64:
		value = int(selector.(int64))
	case uint:
		value = int(selector.(uint))
	case uint8:
		value = int(selector.(uint8))
	case uint16:
		value = int(selector.(uint16))
	case uint32:
		value = int(selector.(uint32))
	case uint64:
		value = int(selector.(uint64))
	default:
		return 0
	}
	return value
}
