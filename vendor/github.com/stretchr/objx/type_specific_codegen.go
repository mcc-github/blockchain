package objx





func (v *Value) Inter(optionalDefault ...interface{}) interface{} {
	if s, ok := v.data.(interface{}); ok {
		return s
	}
	if len(optionalDefault) == 1 {
		return optionalDefault[0]
	}
	return nil
}




func (v *Value) MustInter() interface{} {
	return v.data.(interface{})
}



func (v *Value) InterSlice(optionalDefault ...[]interface{}) []interface{} {
	if s, ok := v.data.([]interface{}); ok {
		return s
	}
	if len(optionalDefault) == 1 {
		return optionalDefault[0]
	}
	return nil
}




func (v *Value) MustInterSlice() []interface{} {
	return v.data.([]interface{})
}


func (v *Value) IsInter() bool {
	_, ok := v.data.(interface{})
	return ok
}


func (v *Value) IsInterSlice() bool {
	_, ok := v.data.([]interface{})
	return ok
}





func (v *Value) EachInter(callback func(int, interface{}) bool) *Value {
	for index, val := range v.MustInterSlice() {
		carryon := callback(index, val)
		if !carryon {
			break
		}
	}
	return v
}




func (v *Value) WhereInter(decider func(int, interface{}) bool) *Value {
	var selected []interface{}
	v.EachInter(func(index int, val interface{}) bool {
		shouldSelect := decider(index, val)
		if !shouldSelect {
			selected = append(selected, val)
		}
		return true
	})
	return &Value{data: selected}
}




func (v *Value) GroupInter(grouper func(int, interface{}) string) *Value {
	groups := make(map[string][]interface{})
	v.EachInter(func(index int, val interface{}) bool {
		group := grouper(index, val)
		if _, ok := groups[group]; !ok {
			groups[group] = make([]interface{}, 0)
		}
		groups[group] = append(groups[group], val)
		return true
	})
	return &Value{data: groups}
}




func (v *Value) ReplaceInter(replacer func(int, interface{}) interface{}) *Value {
	arr := v.MustInterSlice()
	replaced := make([]interface{}, len(arr))
	v.EachInter(func(index int, val interface{}) bool {
		replaced[index] = replacer(index, val)
		return true
	})
	return &Value{data: replaced}
}




func (v *Value) CollectInter(collector func(int, interface{}) interface{}) *Value {
	arr := v.MustInterSlice()
	collected := make([]interface{}, len(arr))
	v.EachInter(func(index int, val interface{}) bool {
		collected[index] = collector(index, val)
		return true
	})
	return &Value{data: collected}
}





func (v *Value) MSI(optionalDefault ...map[string]interface{}) map[string]interface{} {
	if s, ok := v.data.(map[string]interface{}); ok {
		return s
	}
	if len(optionalDefault) == 1 {
		return optionalDefault[0]
	}
	return nil
}




func (v *Value) MustMSI() map[string]interface{} {
	return v.data.(map[string]interface{})
}



func (v *Value) MSISlice(optionalDefault ...[]map[string]interface{}) []map[string]interface{} {
	if s, ok := v.data.([]map[string]interface{}); ok {
		return s
	}
	if len(optionalDefault) == 1 {
		return optionalDefault[0]
	}
	return nil
}




func (v *Value) MustMSISlice() []map[string]interface{} {
	return v.data.([]map[string]interface{})
}


func (v *Value) IsMSI() bool {
	_, ok := v.data.(map[string]interface{})
	return ok
}


func (v *Value) IsMSISlice() bool {
	_, ok := v.data.([]map[string]interface{})
	return ok
}





func (v *Value) EachMSI(callback func(int, map[string]interface{}) bool) *Value {
	for index, val := range v.MustMSISlice() {
		carryon := callback(index, val)
		if !carryon {
			break
		}
	}
	return v
}




func (v *Value) WhereMSI(decider func(int, map[string]interface{}) bool) *Value {
	var selected []map[string]interface{}
	v.EachMSI(func(index int, val map[string]interface{}) bool {
		shouldSelect := decider(index, val)
		if !shouldSelect {
			selected = append(selected, val)
		}
		return true
	})
	return &Value{data: selected}
}




func (v *Value) GroupMSI(grouper func(int, map[string]interface{}) string) *Value {
	groups := make(map[string][]map[string]interface{})
	v.EachMSI(func(index int, val map[string]interface{}) bool {
		group := grouper(index, val)
		if _, ok := groups[group]; !ok {
			groups[group] = make([]map[string]interface{}, 0)
		}
		groups[group] = append(groups[group], val)
		return true
	})
	return &Value{data: groups}
}




func (v *Value) ReplaceMSI(replacer func(int, map[string]interface{}) map[string]interface{}) *Value {
	arr := v.MustMSISlice()
	replaced := make([]map[string]interface{}, len(arr))
	v.EachMSI(func(index int, val map[string]interface{}) bool {
		replaced[index] = replacer(index, val)
		return true
	})
	return &Value{data: replaced}
}




func (v *Value) CollectMSI(collector func(int, map[string]interface{}) interface{}) *Value {
	arr := v.MustMSISlice()
	collected := make([]interface{}, len(arr))
	v.EachMSI(func(index int, val map[string]interface{}) bool {
		collected[index] = collector(index, val)
		return true
	})
	return &Value{data: collected}
}





func (v *Value) ObjxMap(optionalDefault ...(Map)) Map {
	if s, ok := v.data.((Map)); ok {
		return s
	}
	if len(optionalDefault) == 1 {
		return optionalDefault[0]
	}
	return New(nil)
}




func (v *Value) MustObjxMap() Map {
	return v.data.((Map))
}



func (v *Value) ObjxMapSlice(optionalDefault ...[](Map)) [](Map) {
	if s, ok := v.data.([](Map)); ok {
		return s
	}
	if len(optionalDefault) == 1 {
		return optionalDefault[0]
	}
	return nil
}




func (v *Value) MustObjxMapSlice() [](Map) {
	return v.data.([](Map))
}


func (v *Value) IsObjxMap() bool {
	_, ok := v.data.((Map))
	return ok
}


func (v *Value) IsObjxMapSlice() bool {
	_, ok := v.data.([](Map))
	return ok
}





func (v *Value) EachObjxMap(callback func(int, Map) bool) *Value {
	for index, val := range v.MustObjxMapSlice() {
		carryon := callback(index, val)
		if !carryon {
			break
		}
	}
	return v
}




func (v *Value) WhereObjxMap(decider func(int, Map) bool) *Value {
	var selected [](Map)
	v.EachObjxMap(func(index int, val Map) bool {
		shouldSelect := decider(index, val)
		if !shouldSelect {
			selected = append(selected, val)
		}
		return true
	})
	return &Value{data: selected}
}




func (v *Value) GroupObjxMap(grouper func(int, Map) string) *Value {
	groups := make(map[string][](Map))
	v.EachObjxMap(func(index int, val Map) bool {
		group := grouper(index, val)
		if _, ok := groups[group]; !ok {
			groups[group] = make([](Map), 0)
		}
		groups[group] = append(groups[group], val)
		return true
	})
	return &Value{data: groups}
}




func (v *Value) ReplaceObjxMap(replacer func(int, Map) Map) *Value {
	arr := v.MustObjxMapSlice()
	replaced := make([](Map), len(arr))
	v.EachObjxMap(func(index int, val Map) bool {
		replaced[index] = replacer(index, val)
		return true
	})
	return &Value{data: replaced}
}




func (v *Value) CollectObjxMap(collector func(int, Map) interface{}) *Value {
	arr := v.MustObjxMapSlice()
	collected := make([]interface{}, len(arr))
	v.EachObjxMap(func(index int, val Map) bool {
		collected[index] = collector(index, val)
		return true
	})
	return &Value{data: collected}
}





func (v *Value) Bool(optionalDefault ...bool) bool {
	if s, ok := v.data.(bool); ok {
		return s
	}
	if len(optionalDefault) == 1 {
		return optionalDefault[0]
	}
	return false
}




func (v *Value) MustBool() bool {
	return v.data.(bool)
}



func (v *Value) BoolSlice(optionalDefault ...[]bool) []bool {
	if s, ok := v.data.([]bool); ok {
		return s
	}
	if len(optionalDefault) == 1 {
		return optionalDefault[0]
	}
	return nil
}




func (v *Value) MustBoolSlice() []bool {
	return v.data.([]bool)
}


func (v *Value) IsBool() bool {
	_, ok := v.data.(bool)
	return ok
}


func (v *Value) IsBoolSlice() bool {
	_, ok := v.data.([]bool)
	return ok
}





func (v *Value) EachBool(callback func(int, bool) bool) *Value {
	for index, val := range v.MustBoolSlice() {
		carryon := callback(index, val)
		if !carryon {
			break
		}
	}
	return v
}




func (v *Value) WhereBool(decider func(int, bool) bool) *Value {
	var selected []bool
	v.EachBool(func(index int, val bool) bool {
		shouldSelect := decider(index, val)
		if !shouldSelect {
			selected = append(selected, val)
		}
		return true
	})
	return &Value{data: selected}
}




func (v *Value) GroupBool(grouper func(int, bool) string) *Value {
	groups := make(map[string][]bool)
	v.EachBool(func(index int, val bool) bool {
		group := grouper(index, val)
		if _, ok := groups[group]; !ok {
			groups[group] = make([]bool, 0)
		}
		groups[group] = append(groups[group], val)
		return true
	})
	return &Value{data: groups}
}




func (v *Value) ReplaceBool(replacer func(int, bool) bool) *Value {
	arr := v.MustBoolSlice()
	replaced := make([]bool, len(arr))
	v.EachBool(func(index int, val bool) bool {
		replaced[index] = replacer(index, val)
		return true
	})
	return &Value{data: replaced}
}




func (v *Value) CollectBool(collector func(int, bool) interface{}) *Value {
	arr := v.MustBoolSlice()
	collected := make([]interface{}, len(arr))
	v.EachBool(func(index int, val bool) bool {
		collected[index] = collector(index, val)
		return true
	})
	return &Value{data: collected}
}





func (v *Value) Str(optionalDefault ...string) string {
	if s, ok := v.data.(string); ok {
		return s
	}
	if len(optionalDefault) == 1 {
		return optionalDefault[0]
	}
	return ""
}




func (v *Value) MustStr() string {
	return v.data.(string)
}



func (v *Value) StrSlice(optionalDefault ...[]string) []string {
	if s, ok := v.data.([]string); ok {
		return s
	}
	if len(optionalDefault) == 1 {
		return optionalDefault[0]
	}
	return nil
}




func (v *Value) MustStrSlice() []string {
	return v.data.([]string)
}


func (v *Value) IsStr() bool {
	_, ok := v.data.(string)
	return ok
}


func (v *Value) IsStrSlice() bool {
	_, ok := v.data.([]string)
	return ok
}





func (v *Value) EachStr(callback func(int, string) bool) *Value {
	for index, val := range v.MustStrSlice() {
		carryon := callback(index, val)
		if !carryon {
			break
		}
	}
	return v
}




func (v *Value) WhereStr(decider func(int, string) bool) *Value {
	var selected []string
	v.EachStr(func(index int, val string) bool {
		shouldSelect := decider(index, val)
		if !shouldSelect {
			selected = append(selected, val)
		}
		return true
	})
	return &Value{data: selected}
}




func (v *Value) GroupStr(grouper func(int, string) string) *Value {
	groups := make(map[string][]string)
	v.EachStr(func(index int, val string) bool {
		group := grouper(index, val)
		if _, ok := groups[group]; !ok {
			groups[group] = make([]string, 0)
		}
		groups[group] = append(groups[group], val)
		return true
	})
	return &Value{data: groups}
}




func (v *Value) ReplaceStr(replacer func(int, string) string) *Value {
	arr := v.MustStrSlice()
	replaced := make([]string, len(arr))
	v.EachStr(func(index int, val string) bool {
		replaced[index] = replacer(index, val)
		return true
	})
	return &Value{data: replaced}
}




func (v *Value) CollectStr(collector func(int, string) interface{}) *Value {
	arr := v.MustStrSlice()
	collected := make([]interface{}, len(arr))
	v.EachStr(func(index int, val string) bool {
		collected[index] = collector(index, val)
		return true
	})
	return &Value{data: collected}
}





func (v *Value) Int(optionalDefault ...int) int {
	if s, ok := v.data.(int); ok {
		return s
	}
	if len(optionalDefault) == 1 {
		return optionalDefault[0]
	}
	return 0
}




func (v *Value) MustInt() int {
	return v.data.(int)
}



func (v *Value) IntSlice(optionalDefault ...[]int) []int {
	if s, ok := v.data.([]int); ok {
		return s
	}
	if len(optionalDefault) == 1 {
		return optionalDefault[0]
	}
	return nil
}




func (v *Value) MustIntSlice() []int {
	return v.data.([]int)
}


func (v *Value) IsInt() bool {
	_, ok := v.data.(int)
	return ok
}


func (v *Value) IsIntSlice() bool {
	_, ok := v.data.([]int)
	return ok
}





func (v *Value) EachInt(callback func(int, int) bool) *Value {
	for index, val := range v.MustIntSlice() {
		carryon := callback(index, val)
		if !carryon {
			break
		}
	}
	return v
}




func (v *Value) WhereInt(decider func(int, int) bool) *Value {
	var selected []int
	v.EachInt(func(index int, val int) bool {
		shouldSelect := decider(index, val)
		if !shouldSelect {
			selected = append(selected, val)
		}
		return true
	})
	return &Value{data: selected}
}




func (v *Value) GroupInt(grouper func(int, int) string) *Value {
	groups := make(map[string][]int)
	v.EachInt(func(index int, val int) bool {
		group := grouper(index, val)
		if _, ok := groups[group]; !ok {
			groups[group] = make([]int, 0)
		}
		groups[group] = append(groups[group], val)
		return true
	})
	return &Value{data: groups}
}




func (v *Value) ReplaceInt(replacer func(int, int) int) *Value {
	arr := v.MustIntSlice()
	replaced := make([]int, len(arr))
	v.EachInt(func(index int, val int) bool {
		replaced[index] = replacer(index, val)
		return true
	})
	return &Value{data: replaced}
}




func (v *Value) CollectInt(collector func(int, int) interface{}) *Value {
	arr := v.MustIntSlice()
	collected := make([]interface{}, len(arr))
	v.EachInt(func(index int, val int) bool {
		collected[index] = collector(index, val)
		return true
	})
	return &Value{data: collected}
}





func (v *Value) Int8(optionalDefault ...int8) int8 {
	if s, ok := v.data.(int8); ok {
		return s
	}
	if len(optionalDefault) == 1 {
		return optionalDefault[0]
	}
	return 0
}




func (v *Value) MustInt8() int8 {
	return v.data.(int8)
}



func (v *Value) Int8Slice(optionalDefault ...[]int8) []int8 {
	if s, ok := v.data.([]int8); ok {
		return s
	}
	if len(optionalDefault) == 1 {
		return optionalDefault[0]
	}
	return nil
}




func (v *Value) MustInt8Slice() []int8 {
	return v.data.([]int8)
}


func (v *Value) IsInt8() bool {
	_, ok := v.data.(int8)
	return ok
}


func (v *Value) IsInt8Slice() bool {
	_, ok := v.data.([]int8)
	return ok
}





func (v *Value) EachInt8(callback func(int, int8) bool) *Value {
	for index, val := range v.MustInt8Slice() {
		carryon := callback(index, val)
		if !carryon {
			break
		}
	}
	return v
}




func (v *Value) WhereInt8(decider func(int, int8) bool) *Value {
	var selected []int8
	v.EachInt8(func(index int, val int8) bool {
		shouldSelect := decider(index, val)
		if !shouldSelect {
			selected = append(selected, val)
		}
		return true
	})
	return &Value{data: selected}
}




func (v *Value) GroupInt8(grouper func(int, int8) string) *Value {
	groups := make(map[string][]int8)
	v.EachInt8(func(index int, val int8) bool {
		group := grouper(index, val)
		if _, ok := groups[group]; !ok {
			groups[group] = make([]int8, 0)
		}
		groups[group] = append(groups[group], val)
		return true
	})
	return &Value{data: groups}
}




func (v *Value) ReplaceInt8(replacer func(int, int8) int8) *Value {
	arr := v.MustInt8Slice()
	replaced := make([]int8, len(arr))
	v.EachInt8(func(index int, val int8) bool {
		replaced[index] = replacer(index, val)
		return true
	})
	return &Value{data: replaced}
}




func (v *Value) CollectInt8(collector func(int, int8) interface{}) *Value {
	arr := v.MustInt8Slice()
	collected := make([]interface{}, len(arr))
	v.EachInt8(func(index int, val int8) bool {
		collected[index] = collector(index, val)
		return true
	})
	return &Value{data: collected}
}





func (v *Value) Int16(optionalDefault ...int16) int16 {
	if s, ok := v.data.(int16); ok {
		return s
	}
	if len(optionalDefault) == 1 {
		return optionalDefault[0]
	}
	return 0
}




func (v *Value) MustInt16() int16 {
	return v.data.(int16)
}



func (v *Value) Int16Slice(optionalDefault ...[]int16) []int16 {
	if s, ok := v.data.([]int16); ok {
		return s
	}
	if len(optionalDefault) == 1 {
		return optionalDefault[0]
	}
	return nil
}




func (v *Value) MustInt16Slice() []int16 {
	return v.data.([]int16)
}


func (v *Value) IsInt16() bool {
	_, ok := v.data.(int16)
	return ok
}


func (v *Value) IsInt16Slice() bool {
	_, ok := v.data.([]int16)
	return ok
}





func (v *Value) EachInt16(callback func(int, int16) bool) *Value {
	for index, val := range v.MustInt16Slice() {
		carryon := callback(index, val)
		if !carryon {
			break
		}
	}
	return v
}




func (v *Value) WhereInt16(decider func(int, int16) bool) *Value {
	var selected []int16
	v.EachInt16(func(index int, val int16) bool {
		shouldSelect := decider(index, val)
		if !shouldSelect {
			selected = append(selected, val)
		}
		return true
	})
	return &Value{data: selected}
}




func (v *Value) GroupInt16(grouper func(int, int16) string) *Value {
	groups := make(map[string][]int16)
	v.EachInt16(func(index int, val int16) bool {
		group := grouper(index, val)
		if _, ok := groups[group]; !ok {
			groups[group] = make([]int16, 0)
		}
		groups[group] = append(groups[group], val)
		return true
	})
	return &Value{data: groups}
}




func (v *Value) ReplaceInt16(replacer func(int, int16) int16) *Value {
	arr := v.MustInt16Slice()
	replaced := make([]int16, len(arr))
	v.EachInt16(func(index int, val int16) bool {
		replaced[index] = replacer(index, val)
		return true
	})
	return &Value{data: replaced}
}




func (v *Value) CollectInt16(collector func(int, int16) interface{}) *Value {
	arr := v.MustInt16Slice()
	collected := make([]interface{}, len(arr))
	v.EachInt16(func(index int, val int16) bool {
		collected[index] = collector(index, val)
		return true
	})
	return &Value{data: collected}
}





func (v *Value) Int32(optionalDefault ...int32) int32 {
	if s, ok := v.data.(int32); ok {
		return s
	}
	if len(optionalDefault) == 1 {
		return optionalDefault[0]
	}
	return 0
}




func (v *Value) MustInt32() int32 {
	return v.data.(int32)
}



func (v *Value) Int32Slice(optionalDefault ...[]int32) []int32 {
	if s, ok := v.data.([]int32); ok {
		return s
	}
	if len(optionalDefault) == 1 {
		return optionalDefault[0]
	}
	return nil
}




func (v *Value) MustInt32Slice() []int32 {
	return v.data.([]int32)
}


func (v *Value) IsInt32() bool {
	_, ok := v.data.(int32)
	return ok
}


func (v *Value) IsInt32Slice() bool {
	_, ok := v.data.([]int32)
	return ok
}





func (v *Value) EachInt32(callback func(int, int32) bool) *Value {
	for index, val := range v.MustInt32Slice() {
		carryon := callback(index, val)
		if !carryon {
			break
		}
	}
	return v
}




func (v *Value) WhereInt32(decider func(int, int32) bool) *Value {
	var selected []int32
	v.EachInt32(func(index int, val int32) bool {
		shouldSelect := decider(index, val)
		if !shouldSelect {
			selected = append(selected, val)
		}
		return true
	})
	return &Value{data: selected}
}




func (v *Value) GroupInt32(grouper func(int, int32) string) *Value {
	groups := make(map[string][]int32)
	v.EachInt32(func(index int, val int32) bool {
		group := grouper(index, val)
		if _, ok := groups[group]; !ok {
			groups[group] = make([]int32, 0)
		}
		groups[group] = append(groups[group], val)
		return true
	})
	return &Value{data: groups}
}




func (v *Value) ReplaceInt32(replacer func(int, int32) int32) *Value {
	arr := v.MustInt32Slice()
	replaced := make([]int32, len(arr))
	v.EachInt32(func(index int, val int32) bool {
		replaced[index] = replacer(index, val)
		return true
	})
	return &Value{data: replaced}
}




func (v *Value) CollectInt32(collector func(int, int32) interface{}) *Value {
	arr := v.MustInt32Slice()
	collected := make([]interface{}, len(arr))
	v.EachInt32(func(index int, val int32) bool {
		collected[index] = collector(index, val)
		return true
	})
	return &Value{data: collected}
}





func (v *Value) Int64(optionalDefault ...int64) int64 {
	if s, ok := v.data.(int64); ok {
		return s
	}
	if len(optionalDefault) == 1 {
		return optionalDefault[0]
	}
	return 0
}




func (v *Value) MustInt64() int64 {
	return v.data.(int64)
}



func (v *Value) Int64Slice(optionalDefault ...[]int64) []int64 {
	if s, ok := v.data.([]int64); ok {
		return s
	}
	if len(optionalDefault) == 1 {
		return optionalDefault[0]
	}
	return nil
}




func (v *Value) MustInt64Slice() []int64 {
	return v.data.([]int64)
}


func (v *Value) IsInt64() bool {
	_, ok := v.data.(int64)
	return ok
}


func (v *Value) IsInt64Slice() bool {
	_, ok := v.data.([]int64)
	return ok
}





func (v *Value) EachInt64(callback func(int, int64) bool) *Value {
	for index, val := range v.MustInt64Slice() {
		carryon := callback(index, val)
		if !carryon {
			break
		}
	}
	return v
}




func (v *Value) WhereInt64(decider func(int, int64) bool) *Value {
	var selected []int64
	v.EachInt64(func(index int, val int64) bool {
		shouldSelect := decider(index, val)
		if !shouldSelect {
			selected = append(selected, val)
		}
		return true
	})
	return &Value{data: selected}
}




func (v *Value) GroupInt64(grouper func(int, int64) string) *Value {
	groups := make(map[string][]int64)
	v.EachInt64(func(index int, val int64) bool {
		group := grouper(index, val)
		if _, ok := groups[group]; !ok {
			groups[group] = make([]int64, 0)
		}
		groups[group] = append(groups[group], val)
		return true
	})
	return &Value{data: groups}
}




func (v *Value) ReplaceInt64(replacer func(int, int64) int64) *Value {
	arr := v.MustInt64Slice()
	replaced := make([]int64, len(arr))
	v.EachInt64(func(index int, val int64) bool {
		replaced[index] = replacer(index, val)
		return true
	})
	return &Value{data: replaced}
}




func (v *Value) CollectInt64(collector func(int, int64) interface{}) *Value {
	arr := v.MustInt64Slice()
	collected := make([]interface{}, len(arr))
	v.EachInt64(func(index int, val int64) bool {
		collected[index] = collector(index, val)
		return true
	})
	return &Value{data: collected}
}





func (v *Value) Uint(optionalDefault ...uint) uint {
	if s, ok := v.data.(uint); ok {
		return s
	}
	if len(optionalDefault) == 1 {
		return optionalDefault[0]
	}
	return 0
}




func (v *Value) MustUint() uint {
	return v.data.(uint)
}



func (v *Value) UintSlice(optionalDefault ...[]uint) []uint {
	if s, ok := v.data.([]uint); ok {
		return s
	}
	if len(optionalDefault) == 1 {
		return optionalDefault[0]
	}
	return nil
}




func (v *Value) MustUintSlice() []uint {
	return v.data.([]uint)
}


func (v *Value) IsUint() bool {
	_, ok := v.data.(uint)
	return ok
}


func (v *Value) IsUintSlice() bool {
	_, ok := v.data.([]uint)
	return ok
}





func (v *Value) EachUint(callback func(int, uint) bool) *Value {
	for index, val := range v.MustUintSlice() {
		carryon := callback(index, val)
		if !carryon {
			break
		}
	}
	return v
}




func (v *Value) WhereUint(decider func(int, uint) bool) *Value {
	var selected []uint
	v.EachUint(func(index int, val uint) bool {
		shouldSelect := decider(index, val)
		if !shouldSelect {
			selected = append(selected, val)
		}
		return true
	})
	return &Value{data: selected}
}




func (v *Value) GroupUint(grouper func(int, uint) string) *Value {
	groups := make(map[string][]uint)
	v.EachUint(func(index int, val uint) bool {
		group := grouper(index, val)
		if _, ok := groups[group]; !ok {
			groups[group] = make([]uint, 0)
		}
		groups[group] = append(groups[group], val)
		return true
	})
	return &Value{data: groups}
}




func (v *Value) ReplaceUint(replacer func(int, uint) uint) *Value {
	arr := v.MustUintSlice()
	replaced := make([]uint, len(arr))
	v.EachUint(func(index int, val uint) bool {
		replaced[index] = replacer(index, val)
		return true
	})
	return &Value{data: replaced}
}




func (v *Value) CollectUint(collector func(int, uint) interface{}) *Value {
	arr := v.MustUintSlice()
	collected := make([]interface{}, len(arr))
	v.EachUint(func(index int, val uint) bool {
		collected[index] = collector(index, val)
		return true
	})
	return &Value{data: collected}
}





func (v *Value) Uint8(optionalDefault ...uint8) uint8 {
	if s, ok := v.data.(uint8); ok {
		return s
	}
	if len(optionalDefault) == 1 {
		return optionalDefault[0]
	}
	return 0
}




func (v *Value) MustUint8() uint8 {
	return v.data.(uint8)
}



func (v *Value) Uint8Slice(optionalDefault ...[]uint8) []uint8 {
	if s, ok := v.data.([]uint8); ok {
		return s
	}
	if len(optionalDefault) == 1 {
		return optionalDefault[0]
	}
	return nil
}




func (v *Value) MustUint8Slice() []uint8 {
	return v.data.([]uint8)
}


func (v *Value) IsUint8() bool {
	_, ok := v.data.(uint8)
	return ok
}


func (v *Value) IsUint8Slice() bool {
	_, ok := v.data.([]uint8)
	return ok
}





func (v *Value) EachUint8(callback func(int, uint8) bool) *Value {
	for index, val := range v.MustUint8Slice() {
		carryon := callback(index, val)
		if !carryon {
			break
		}
	}
	return v
}




func (v *Value) WhereUint8(decider func(int, uint8) bool) *Value {
	var selected []uint8
	v.EachUint8(func(index int, val uint8) bool {
		shouldSelect := decider(index, val)
		if !shouldSelect {
			selected = append(selected, val)
		}
		return true
	})
	return &Value{data: selected}
}




func (v *Value) GroupUint8(grouper func(int, uint8) string) *Value {
	groups := make(map[string][]uint8)
	v.EachUint8(func(index int, val uint8) bool {
		group := grouper(index, val)
		if _, ok := groups[group]; !ok {
			groups[group] = make([]uint8, 0)
		}
		groups[group] = append(groups[group], val)
		return true
	})
	return &Value{data: groups}
}




func (v *Value) ReplaceUint8(replacer func(int, uint8) uint8) *Value {
	arr := v.MustUint8Slice()
	replaced := make([]uint8, len(arr))
	v.EachUint8(func(index int, val uint8) bool {
		replaced[index] = replacer(index, val)
		return true
	})
	return &Value{data: replaced}
}




func (v *Value) CollectUint8(collector func(int, uint8) interface{}) *Value {
	arr := v.MustUint8Slice()
	collected := make([]interface{}, len(arr))
	v.EachUint8(func(index int, val uint8) bool {
		collected[index] = collector(index, val)
		return true
	})
	return &Value{data: collected}
}





func (v *Value) Uint16(optionalDefault ...uint16) uint16 {
	if s, ok := v.data.(uint16); ok {
		return s
	}
	if len(optionalDefault) == 1 {
		return optionalDefault[0]
	}
	return 0
}




func (v *Value) MustUint16() uint16 {
	return v.data.(uint16)
}



func (v *Value) Uint16Slice(optionalDefault ...[]uint16) []uint16 {
	if s, ok := v.data.([]uint16); ok {
		return s
	}
	if len(optionalDefault) == 1 {
		return optionalDefault[0]
	}
	return nil
}




func (v *Value) MustUint16Slice() []uint16 {
	return v.data.([]uint16)
}


func (v *Value) IsUint16() bool {
	_, ok := v.data.(uint16)
	return ok
}


func (v *Value) IsUint16Slice() bool {
	_, ok := v.data.([]uint16)
	return ok
}





func (v *Value) EachUint16(callback func(int, uint16) bool) *Value {
	for index, val := range v.MustUint16Slice() {
		carryon := callback(index, val)
		if !carryon {
			break
		}
	}
	return v
}




func (v *Value) WhereUint16(decider func(int, uint16) bool) *Value {
	var selected []uint16
	v.EachUint16(func(index int, val uint16) bool {
		shouldSelect := decider(index, val)
		if !shouldSelect {
			selected = append(selected, val)
		}
		return true
	})
	return &Value{data: selected}
}




func (v *Value) GroupUint16(grouper func(int, uint16) string) *Value {
	groups := make(map[string][]uint16)
	v.EachUint16(func(index int, val uint16) bool {
		group := grouper(index, val)
		if _, ok := groups[group]; !ok {
			groups[group] = make([]uint16, 0)
		}
		groups[group] = append(groups[group], val)
		return true
	})
	return &Value{data: groups}
}




func (v *Value) ReplaceUint16(replacer func(int, uint16) uint16) *Value {
	arr := v.MustUint16Slice()
	replaced := make([]uint16, len(arr))
	v.EachUint16(func(index int, val uint16) bool {
		replaced[index] = replacer(index, val)
		return true
	})
	return &Value{data: replaced}
}




func (v *Value) CollectUint16(collector func(int, uint16) interface{}) *Value {
	arr := v.MustUint16Slice()
	collected := make([]interface{}, len(arr))
	v.EachUint16(func(index int, val uint16) bool {
		collected[index] = collector(index, val)
		return true
	})
	return &Value{data: collected}
}





func (v *Value) Uint32(optionalDefault ...uint32) uint32 {
	if s, ok := v.data.(uint32); ok {
		return s
	}
	if len(optionalDefault) == 1 {
		return optionalDefault[0]
	}
	return 0
}




func (v *Value) MustUint32() uint32 {
	return v.data.(uint32)
}



func (v *Value) Uint32Slice(optionalDefault ...[]uint32) []uint32 {
	if s, ok := v.data.([]uint32); ok {
		return s
	}
	if len(optionalDefault) == 1 {
		return optionalDefault[0]
	}
	return nil
}




func (v *Value) MustUint32Slice() []uint32 {
	return v.data.([]uint32)
}


func (v *Value) IsUint32() bool {
	_, ok := v.data.(uint32)
	return ok
}


func (v *Value) IsUint32Slice() bool {
	_, ok := v.data.([]uint32)
	return ok
}





func (v *Value) EachUint32(callback func(int, uint32) bool) *Value {
	for index, val := range v.MustUint32Slice() {
		carryon := callback(index, val)
		if !carryon {
			break
		}
	}
	return v
}




func (v *Value) WhereUint32(decider func(int, uint32) bool) *Value {
	var selected []uint32
	v.EachUint32(func(index int, val uint32) bool {
		shouldSelect := decider(index, val)
		if !shouldSelect {
			selected = append(selected, val)
		}
		return true
	})
	return &Value{data: selected}
}




func (v *Value) GroupUint32(grouper func(int, uint32) string) *Value {
	groups := make(map[string][]uint32)
	v.EachUint32(func(index int, val uint32) bool {
		group := grouper(index, val)
		if _, ok := groups[group]; !ok {
			groups[group] = make([]uint32, 0)
		}
		groups[group] = append(groups[group], val)
		return true
	})
	return &Value{data: groups}
}




func (v *Value) ReplaceUint32(replacer func(int, uint32) uint32) *Value {
	arr := v.MustUint32Slice()
	replaced := make([]uint32, len(arr))
	v.EachUint32(func(index int, val uint32) bool {
		replaced[index] = replacer(index, val)
		return true
	})
	return &Value{data: replaced}
}




func (v *Value) CollectUint32(collector func(int, uint32) interface{}) *Value {
	arr := v.MustUint32Slice()
	collected := make([]interface{}, len(arr))
	v.EachUint32(func(index int, val uint32) bool {
		collected[index] = collector(index, val)
		return true
	})
	return &Value{data: collected}
}





func (v *Value) Uint64(optionalDefault ...uint64) uint64 {
	if s, ok := v.data.(uint64); ok {
		return s
	}
	if len(optionalDefault) == 1 {
		return optionalDefault[0]
	}
	return 0
}




func (v *Value) MustUint64() uint64 {
	return v.data.(uint64)
}



func (v *Value) Uint64Slice(optionalDefault ...[]uint64) []uint64 {
	if s, ok := v.data.([]uint64); ok {
		return s
	}
	if len(optionalDefault) == 1 {
		return optionalDefault[0]
	}
	return nil
}




func (v *Value) MustUint64Slice() []uint64 {
	return v.data.([]uint64)
}


func (v *Value) IsUint64() bool {
	_, ok := v.data.(uint64)
	return ok
}


func (v *Value) IsUint64Slice() bool {
	_, ok := v.data.([]uint64)
	return ok
}





func (v *Value) EachUint64(callback func(int, uint64) bool) *Value {
	for index, val := range v.MustUint64Slice() {
		carryon := callback(index, val)
		if !carryon {
			break
		}
	}
	return v
}




func (v *Value) WhereUint64(decider func(int, uint64) bool) *Value {
	var selected []uint64
	v.EachUint64(func(index int, val uint64) bool {
		shouldSelect := decider(index, val)
		if !shouldSelect {
			selected = append(selected, val)
		}
		return true
	})
	return &Value{data: selected}
}




func (v *Value) GroupUint64(grouper func(int, uint64) string) *Value {
	groups := make(map[string][]uint64)
	v.EachUint64(func(index int, val uint64) bool {
		group := grouper(index, val)
		if _, ok := groups[group]; !ok {
			groups[group] = make([]uint64, 0)
		}
		groups[group] = append(groups[group], val)
		return true
	})
	return &Value{data: groups}
}




func (v *Value) ReplaceUint64(replacer func(int, uint64) uint64) *Value {
	arr := v.MustUint64Slice()
	replaced := make([]uint64, len(arr))
	v.EachUint64(func(index int, val uint64) bool {
		replaced[index] = replacer(index, val)
		return true
	})
	return &Value{data: replaced}
}




func (v *Value) CollectUint64(collector func(int, uint64) interface{}) *Value {
	arr := v.MustUint64Slice()
	collected := make([]interface{}, len(arr))
	v.EachUint64(func(index int, val uint64) bool {
		collected[index] = collector(index, val)
		return true
	})
	return &Value{data: collected}
}





func (v *Value) Uintptr(optionalDefault ...uintptr) uintptr {
	if s, ok := v.data.(uintptr); ok {
		return s
	}
	if len(optionalDefault) == 1 {
		return optionalDefault[0]
	}
	return 0
}




func (v *Value) MustUintptr() uintptr {
	return v.data.(uintptr)
}



func (v *Value) UintptrSlice(optionalDefault ...[]uintptr) []uintptr {
	if s, ok := v.data.([]uintptr); ok {
		return s
	}
	if len(optionalDefault) == 1 {
		return optionalDefault[0]
	}
	return nil
}




func (v *Value) MustUintptrSlice() []uintptr {
	return v.data.([]uintptr)
}


func (v *Value) IsUintptr() bool {
	_, ok := v.data.(uintptr)
	return ok
}


func (v *Value) IsUintptrSlice() bool {
	_, ok := v.data.([]uintptr)
	return ok
}





func (v *Value) EachUintptr(callback func(int, uintptr) bool) *Value {
	for index, val := range v.MustUintptrSlice() {
		carryon := callback(index, val)
		if !carryon {
			break
		}
	}
	return v
}




func (v *Value) WhereUintptr(decider func(int, uintptr) bool) *Value {
	var selected []uintptr
	v.EachUintptr(func(index int, val uintptr) bool {
		shouldSelect := decider(index, val)
		if !shouldSelect {
			selected = append(selected, val)
		}
		return true
	})
	return &Value{data: selected}
}




func (v *Value) GroupUintptr(grouper func(int, uintptr) string) *Value {
	groups := make(map[string][]uintptr)
	v.EachUintptr(func(index int, val uintptr) bool {
		group := grouper(index, val)
		if _, ok := groups[group]; !ok {
			groups[group] = make([]uintptr, 0)
		}
		groups[group] = append(groups[group], val)
		return true
	})
	return &Value{data: groups}
}




func (v *Value) ReplaceUintptr(replacer func(int, uintptr) uintptr) *Value {
	arr := v.MustUintptrSlice()
	replaced := make([]uintptr, len(arr))
	v.EachUintptr(func(index int, val uintptr) bool {
		replaced[index] = replacer(index, val)
		return true
	})
	return &Value{data: replaced}
}




func (v *Value) CollectUintptr(collector func(int, uintptr) interface{}) *Value {
	arr := v.MustUintptrSlice()
	collected := make([]interface{}, len(arr))
	v.EachUintptr(func(index int, val uintptr) bool {
		collected[index] = collector(index, val)
		return true
	})
	return &Value{data: collected}
}





func (v *Value) Float32(optionalDefault ...float32) float32 {
	if s, ok := v.data.(float32); ok {
		return s
	}
	if len(optionalDefault) == 1 {
		return optionalDefault[0]
	}
	return 0
}




func (v *Value) MustFloat32() float32 {
	return v.data.(float32)
}



func (v *Value) Float32Slice(optionalDefault ...[]float32) []float32 {
	if s, ok := v.data.([]float32); ok {
		return s
	}
	if len(optionalDefault) == 1 {
		return optionalDefault[0]
	}
	return nil
}




func (v *Value) MustFloat32Slice() []float32 {
	return v.data.([]float32)
}


func (v *Value) IsFloat32() bool {
	_, ok := v.data.(float32)
	return ok
}


func (v *Value) IsFloat32Slice() bool {
	_, ok := v.data.([]float32)
	return ok
}





func (v *Value) EachFloat32(callback func(int, float32) bool) *Value {
	for index, val := range v.MustFloat32Slice() {
		carryon := callback(index, val)
		if !carryon {
			break
		}
	}
	return v
}




func (v *Value) WhereFloat32(decider func(int, float32) bool) *Value {
	var selected []float32
	v.EachFloat32(func(index int, val float32) bool {
		shouldSelect := decider(index, val)
		if !shouldSelect {
			selected = append(selected, val)
		}
		return true
	})
	return &Value{data: selected}
}




func (v *Value) GroupFloat32(grouper func(int, float32) string) *Value {
	groups := make(map[string][]float32)
	v.EachFloat32(func(index int, val float32) bool {
		group := grouper(index, val)
		if _, ok := groups[group]; !ok {
			groups[group] = make([]float32, 0)
		}
		groups[group] = append(groups[group], val)
		return true
	})
	return &Value{data: groups}
}




func (v *Value) ReplaceFloat32(replacer func(int, float32) float32) *Value {
	arr := v.MustFloat32Slice()
	replaced := make([]float32, len(arr))
	v.EachFloat32(func(index int, val float32) bool {
		replaced[index] = replacer(index, val)
		return true
	})
	return &Value{data: replaced}
}




func (v *Value) CollectFloat32(collector func(int, float32) interface{}) *Value {
	arr := v.MustFloat32Slice()
	collected := make([]interface{}, len(arr))
	v.EachFloat32(func(index int, val float32) bool {
		collected[index] = collector(index, val)
		return true
	})
	return &Value{data: collected}
}





func (v *Value) Float64(optionalDefault ...float64) float64 {
	if s, ok := v.data.(float64); ok {
		return s
	}
	if len(optionalDefault) == 1 {
		return optionalDefault[0]
	}
	return 0
}




func (v *Value) MustFloat64() float64 {
	return v.data.(float64)
}



func (v *Value) Float64Slice(optionalDefault ...[]float64) []float64 {
	if s, ok := v.data.([]float64); ok {
		return s
	}
	if len(optionalDefault) == 1 {
		return optionalDefault[0]
	}
	return nil
}




func (v *Value) MustFloat64Slice() []float64 {
	return v.data.([]float64)
}


func (v *Value) IsFloat64() bool {
	_, ok := v.data.(float64)
	return ok
}


func (v *Value) IsFloat64Slice() bool {
	_, ok := v.data.([]float64)
	return ok
}





func (v *Value) EachFloat64(callback func(int, float64) bool) *Value {
	for index, val := range v.MustFloat64Slice() {
		carryon := callback(index, val)
		if !carryon {
			break
		}
	}
	return v
}




func (v *Value) WhereFloat64(decider func(int, float64) bool) *Value {
	var selected []float64
	v.EachFloat64(func(index int, val float64) bool {
		shouldSelect := decider(index, val)
		if !shouldSelect {
			selected = append(selected, val)
		}
		return true
	})
	return &Value{data: selected}
}




func (v *Value) GroupFloat64(grouper func(int, float64) string) *Value {
	groups := make(map[string][]float64)
	v.EachFloat64(func(index int, val float64) bool {
		group := grouper(index, val)
		if _, ok := groups[group]; !ok {
			groups[group] = make([]float64, 0)
		}
		groups[group] = append(groups[group], val)
		return true
	})
	return &Value{data: groups}
}




func (v *Value) ReplaceFloat64(replacer func(int, float64) float64) *Value {
	arr := v.MustFloat64Slice()
	replaced := make([]float64, len(arr))
	v.EachFloat64(func(index int, val float64) bool {
		replaced[index] = replacer(index, val)
		return true
	})
	return &Value{data: replaced}
}




func (v *Value) CollectFloat64(collector func(int, float64) interface{}) *Value {
	arr := v.MustFloat64Slice()
	collected := make([]interface{}, len(arr))
	v.EachFloat64(func(index int, val float64) bool {
		collected[index] = collector(index, val)
		return true
	})
	return &Value{data: collected}
}





func (v *Value) Complex64(optionalDefault ...complex64) complex64 {
	if s, ok := v.data.(complex64); ok {
		return s
	}
	if len(optionalDefault) == 1 {
		return optionalDefault[0]
	}
	return 0
}




func (v *Value) MustComplex64() complex64 {
	return v.data.(complex64)
}



func (v *Value) Complex64Slice(optionalDefault ...[]complex64) []complex64 {
	if s, ok := v.data.([]complex64); ok {
		return s
	}
	if len(optionalDefault) == 1 {
		return optionalDefault[0]
	}
	return nil
}




func (v *Value) MustComplex64Slice() []complex64 {
	return v.data.([]complex64)
}


func (v *Value) IsComplex64() bool {
	_, ok := v.data.(complex64)
	return ok
}


func (v *Value) IsComplex64Slice() bool {
	_, ok := v.data.([]complex64)
	return ok
}





func (v *Value) EachComplex64(callback func(int, complex64) bool) *Value {
	for index, val := range v.MustComplex64Slice() {
		carryon := callback(index, val)
		if !carryon {
			break
		}
	}
	return v
}




func (v *Value) WhereComplex64(decider func(int, complex64) bool) *Value {
	var selected []complex64
	v.EachComplex64(func(index int, val complex64) bool {
		shouldSelect := decider(index, val)
		if !shouldSelect {
			selected = append(selected, val)
		}
		return true
	})
	return &Value{data: selected}
}




func (v *Value) GroupComplex64(grouper func(int, complex64) string) *Value {
	groups := make(map[string][]complex64)
	v.EachComplex64(func(index int, val complex64) bool {
		group := grouper(index, val)
		if _, ok := groups[group]; !ok {
			groups[group] = make([]complex64, 0)
		}
		groups[group] = append(groups[group], val)
		return true
	})
	return &Value{data: groups}
}




func (v *Value) ReplaceComplex64(replacer func(int, complex64) complex64) *Value {
	arr := v.MustComplex64Slice()
	replaced := make([]complex64, len(arr))
	v.EachComplex64(func(index int, val complex64) bool {
		replaced[index] = replacer(index, val)
		return true
	})
	return &Value{data: replaced}
}




func (v *Value) CollectComplex64(collector func(int, complex64) interface{}) *Value {
	arr := v.MustComplex64Slice()
	collected := make([]interface{}, len(arr))
	v.EachComplex64(func(index int, val complex64) bool {
		collected[index] = collector(index, val)
		return true
	})
	return &Value{data: collected}
}





func (v *Value) Complex128(optionalDefault ...complex128) complex128 {
	if s, ok := v.data.(complex128); ok {
		return s
	}
	if len(optionalDefault) == 1 {
		return optionalDefault[0]
	}
	return 0
}




func (v *Value) MustComplex128() complex128 {
	return v.data.(complex128)
}



func (v *Value) Complex128Slice(optionalDefault ...[]complex128) []complex128 {
	if s, ok := v.data.([]complex128); ok {
		return s
	}
	if len(optionalDefault) == 1 {
		return optionalDefault[0]
	}
	return nil
}




func (v *Value) MustComplex128Slice() []complex128 {
	return v.data.([]complex128)
}


func (v *Value) IsComplex128() bool {
	_, ok := v.data.(complex128)
	return ok
}


func (v *Value) IsComplex128Slice() bool {
	_, ok := v.data.([]complex128)
	return ok
}





func (v *Value) EachComplex128(callback func(int, complex128) bool) *Value {
	for index, val := range v.MustComplex128Slice() {
		carryon := callback(index, val)
		if !carryon {
			break
		}
	}
	return v
}




func (v *Value) WhereComplex128(decider func(int, complex128) bool) *Value {
	var selected []complex128
	v.EachComplex128(func(index int, val complex128) bool {
		shouldSelect := decider(index, val)
		if !shouldSelect {
			selected = append(selected, val)
		}
		return true
	})
	return &Value{data: selected}
}




func (v *Value) GroupComplex128(grouper func(int, complex128) string) *Value {
	groups := make(map[string][]complex128)
	v.EachComplex128(func(index int, val complex128) bool {
		group := grouper(index, val)
		if _, ok := groups[group]; !ok {
			groups[group] = make([]complex128, 0)
		}
		groups[group] = append(groups[group], val)
		return true
	})
	return &Value{data: groups}
}




func (v *Value) ReplaceComplex128(replacer func(int, complex128) complex128) *Value {
	arr := v.MustComplex128Slice()
	replaced := make([]complex128, len(arr))
	v.EachComplex128(func(index int, val complex128) bool {
		replaced[index] = replacer(index, val)
		return true
	})
	return &Value{data: replaced}
}




func (v *Value) CollectComplex128(collector func(int, complex128) interface{}) *Value {
	arr := v.MustComplex128Slice()
	collected := make([]interface{}, len(arr))
	v.EachComplex128(func(index int, val complex128) bool {
		collected[index] = collector(index, val)
		return true
	})
	return &Value{data: collected}
}
