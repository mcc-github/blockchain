package nat

import (
	"sort"
	"strings"
)

type portSorter struct {
	ports []Port
	by    func(i, j Port) bool
}

func (s *portSorter) Len() int {
	return len(s.ports)
}

func (s *portSorter) Swap(i, j int) {
	s.ports[i], s.ports[j] = s.ports[j], s.ports[i]
}

func (s *portSorter) Less(i, j int) bool {
	ip := s.ports[i]
	jp := s.ports[j]

	return s.by(ip, jp)
}




func Sort(ports []Port, predicate func(i, j Port) bool) {
	s := &portSorter{ports, predicate}
	sort.Sort(s)
}

type portMapEntry struct {
	port    Port
	binding PortBinding
}

type portMapSorter []portMapEntry

func (s portMapSorter) Len() int      { return len(s) }
func (s portMapSorter) Swap(i, j int) { s[i], s[j] = s[j], s[i] }





func (s portMapSorter) Less(i, j int) bool {
	pi, pj := s[i].port, s[j].port
	hpi, hpj := toInt(s[i].binding.HostPort), toInt(s[j].binding.HostPort)
	return hpi > hpj || pi.Int() > pj.Int() || (pi.Int() == pj.Int() && strings.ToLower(pi.Proto()) == "tcp")
}



func SortPortMap(ports []Port, bindings PortMap) {
	s := portMapSorter{}
	for _, p := range ports {
		if binding, ok := bindings[p]; ok {
			for _, b := range binding {
				s = append(s, portMapEntry{port: p, binding: b})
			}
			bindings[p] = []PortBinding{}
		} else {
			s = append(s, portMapEntry{port: p})
		}
	}

	sort.Sort(s)
	var (
		i  int
		pm = make(map[Port]struct{})
	)
	
	for _, entry := range s {
		if _, ok := pm[entry.port]; !ok {
			ports[i] = entry.port
			pm[entry.port] = struct{}{}
			i++
		}
		
		if _, ok := bindings[entry.port]; ok {
			bindings[entry.port] = append(bindings[entry.port], entry.binding)
		}
	}
}

func toInt(s string) uint64 {
	i, _, err := ParsePortRange(s)
	if err != nil {
		i = 0
	}
	return i
}
