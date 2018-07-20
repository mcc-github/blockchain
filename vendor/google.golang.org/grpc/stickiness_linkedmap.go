

package grpc

import (
	"container/list"
)

type linkedMapKVPair struct {
	key   string
	value *stickyStoreEntry
}







type linkedMap struct {
	m map[string]*list.Element
	l *list.List 
}


func newLinkedMap() *linkedMap {
	return &linkedMap{
		m: make(map[string]*list.Element),
		l: list.New(),
	}
}


func (m *linkedMap) put(key string, value *stickyStoreEntry) {
	if oldE, ok := m.m[key]; ok {
		
		m.l.Remove(oldE)
	}
	e := m.l.PushBack(&linkedMapKVPair{key: key, value: value})
	m.m[key] = e
}


func (m *linkedMap) get(key string) (*stickyStoreEntry, bool) {
	e, ok := m.m[key]
	if !ok {
		return nil, false
	}
	m.l.MoveToBack(e)
	return e.Value.(*linkedMapKVPair).value, true
}



func (m *linkedMap) remove(key string) (*stickyStoreEntry, bool) {
	e, ok := m.m[key]
	if !ok {
		return nil, false
	}
	delete(m.m, key)
	m.l.Remove(e)
	return e.Value.(*linkedMapKVPair).value, true
}


func (m *linkedMap) len() int {
	return len(m.m)
}


func (m *linkedMap) clear() {
	m.m = make(map[string]*list.Element)
	m.l = list.New()
}


func (m *linkedMap) removeOldest() {
	e := m.l.Front()
	m.l.Remove(e)
	delete(m.m, e.Value.(*linkedMapKVPair).key)
}
