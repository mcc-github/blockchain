package toml







import (
	"reflect"
	"sort"
	"sync"
)


type field struct {
	name  string       
	tag   bool         
	index []int        
	typ   reflect.Type 
}




type byName []field

func (x byName) Len() int { return len(x) }

func (x byName) Swap(i, j int) { x[i], x[j] = x[j], x[i] }

func (x byName) Less(i, j int) bool {
	if x[i].name != x[j].name {
		return x[i].name < x[j].name
	}
	if len(x[i].index) != len(x[j].index) {
		return len(x[i].index) < len(x[j].index)
	}
	if x[i].tag != x[j].tag {
		return x[i].tag
	}
	return byIndex(x).Less(i, j)
}


type byIndex []field

func (x byIndex) Len() int { return len(x) }

func (x byIndex) Swap(i, j int) { x[i], x[j] = x[j], x[i] }

func (x byIndex) Less(i, j int) bool {
	for k, xik := range x[i].index {
		if k >= len(x[j].index) {
			return false
		}
		if xik != x[j].index[k] {
			return xik < x[j].index[k]
		}
	}
	return len(x[i].index) < len(x[j].index)
}




func typeFields(t reflect.Type) []field {
	
	current := []field{}
	next := []field{{typ: t}}

	
	count := map[reflect.Type]int{}
	nextCount := map[reflect.Type]int{}

	
	visited := map[reflect.Type]bool{}

	
	var fields []field

	for len(next) > 0 {
		current, next = next, current[:0]
		count, nextCount = nextCount, map[reflect.Type]int{}

		for _, f := range current {
			if visited[f.typ] {
				continue
			}
			visited[f.typ] = true

			
			for i := 0; i < f.typ.NumField(); i++ {
				sf := f.typ.Field(i)
				if sf.PkgPath != "" && !sf.Anonymous { 
					continue
				}
				opts := getOptions(sf.Tag)
				if opts.skip {
					continue
				}
				index := make([]int, len(f.index)+1)
				copy(index, f.index)
				index[len(f.index)] = i

				ft := sf.Type
				if ft.Name() == "" && ft.Kind() == reflect.Ptr {
					
					ft = ft.Elem()
				}

				
				if opts.name != "" || !sf.Anonymous || ft.Kind() != reflect.Struct {
					tagged := opts.name != ""
					name := opts.name
					if name == "" {
						name = sf.Name
					}
					fields = append(fields, field{name, tagged, index, ft})
					if count[f.typ] > 1 {
						
						
						
						
						fields = append(fields, fields[len(fields)-1])
					}
					continue
				}

				
				nextCount[ft]++
				if nextCount[ft] == 1 {
					f := field{name: ft.Name(), index: index, typ: ft}
					next = append(next, f)
				}
			}
		}
	}

	sort.Sort(byName(fields))

	
	

	
	
	
	out := fields[:0]
	for advance, i := 0, 0; i < len(fields); i += advance {
		
		
		fi := fields[i]
		name := fi.name
		for advance = 1; i+advance < len(fields); advance++ {
			fj := fields[i+advance]
			if fj.name != name {
				break
			}
		}
		if advance == 1 { 
			out = append(out, fi)
			continue
		}
		dominant, ok := dominantField(fields[i : i+advance])
		if ok {
			out = append(out, dominant)
		}
	}

	fields = out
	sort.Sort(byIndex(fields))

	return fields
}







func dominantField(fields []field) (field, bool) {
	
	
	
	length := len(fields[0].index)
	tagged := -1 
	for i, f := range fields {
		if len(f.index) > length {
			fields = fields[:i]
			break
		}
		if f.tag {
			if tagged >= 0 {
				
				
				return field{}, false
			}
			tagged = i
		}
	}
	if tagged >= 0 {
		return fields[tagged], true
	}
	
	
	
	if len(fields) > 1 {
		return field{}, false
	}
	return fields[0], true
}

var fieldCache struct {
	sync.RWMutex
	m map[reflect.Type][]field
}


func cachedTypeFields(t reflect.Type) []field {
	fieldCache.RLock()
	f := fieldCache.m[t]
	fieldCache.RUnlock()
	if f != nil {
		return f
	}

	
	
	f = typeFields(t)
	if f == nil {
		f = []field{}
	}

	fieldCache.Lock()
	if fieldCache.m == nil {
		fieldCache.m = map[reflect.Type][]field{}
	}
	fieldCache.m[t] = f
	fieldCache.Unlock()
	return f
}
