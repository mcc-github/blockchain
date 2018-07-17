



package cldr

import (
	"fmt"
	"reflect"
	"sort"
)




type Slice struct {
	ptr reflect.Value
	typ reflect.Type
}


func (s *Slice) Value() reflect.Value {
	return s.ptr.Elem()
}





func MakeSlice(slicePtr interface{}) Slice {
	ptr := reflect.ValueOf(slicePtr)
	if ptr.Kind() != reflect.Ptr {
		panic(fmt.Sprintf("MakeSlice: argument must be pointer to slice, found %v", ptr.Type()))
	}
	sl := ptr.Elem()
	if sl.Kind() != reflect.Slice {
		panic(fmt.Sprintf("MakeSlice: argument must point to a slice, found %v", sl.Type()))
	}
	intf := reflect.TypeOf((*Elem)(nil)).Elem()
	if !sl.Type().Elem().Implements(intf) {
		panic(fmt.Sprintf("MakeSlice: element type of slice (%v) does not implement Elem", sl.Type().Elem()))
	}
	nsl := reflect.MakeSlice(sl.Type(), sl.Len(), sl.Len())
	reflect.Copy(nsl, sl)
	sl.Set(nsl)
	return Slice{
		ptr: ptr,
		typ: sl.Type().Elem().Elem(),
	}
}

func (s Slice) indexForAttr(a string) []int {
	for i := iter(reflect.Zero(s.typ)); !i.done(); i.next() {
		if n, _ := xmlName(i.field()); n == a {
			return i.index
		}
	}
	panic(fmt.Sprintf("MakeSlice: no attribute %q for type %v", a, s.typ))
}


func (s Slice) Filter(fn func(e Elem) bool) {
	k := 0
	sl := s.Value()
	for i := 0; i < sl.Len(); i++ {
		vi := sl.Index(i)
		if fn(vi.Interface().(Elem)) {
			sl.Index(k).Set(vi)
			k++
		}
	}
	sl.Set(sl.Slice(0, k))
}



func (s Slice) Group(fn func(e Elem) string) []Slice {
	m := make(map[string][]reflect.Value)
	sl := s.Value()
	for i := 0; i < sl.Len(); i++ {
		vi := sl.Index(i)
		key := fn(vi.Interface().(Elem))
		m[key] = append(m[key], vi)
	}
	keys := []string{}
	for k, _ := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	res := []Slice{}
	for _, k := range keys {
		nsl := reflect.New(sl.Type())
		nsl.Elem().Set(reflect.Append(nsl.Elem(), m[k]...))
		res = append(res, MakeSlice(nsl.Interface()))
	}
	return res
}



func (s Slice) SelectAnyOf(attr string, values ...string) {
	index := s.indexForAttr(attr)
	s.Filter(func(e Elem) bool {
		vf := reflect.ValueOf(e).Elem().FieldByIndex(index)
		return in(values, vf.String())
	})
}






func (s Slice) SelectOnePerGroup(a string, v []string) {
	index := s.indexForAttr(a)
	grouped := s.Group(func(e Elem) string { return Key(e, a) })
	sl := s.Value()
	sl.Set(sl.Slice(0, 0))
	for _, g := range grouped {
		e := reflect.Value{}
		found := len(v)
		gsl := g.Value()
		for i := 0; i < gsl.Len(); i++ {
			vi := gsl.Index(i).Elem().FieldByIndex(index)
			j := 0
			for ; j < len(v) && v[j] != vi.String(); j++ {
			}
			if j < found {
				found = j
				e = gsl.Index(i)
			}
		}
		if found < len(v) {
			sl.Set(reflect.Append(sl, e))
		}
	}
}




func (s Slice) SelectDraft(d Draft) {
	s.SelectOnePerGroup("draft", drafts[len(drafts)-2-int(d):])
}
