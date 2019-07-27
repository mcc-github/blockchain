/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package util

import (
	"fmt"
	"math/rand"
	"reflect"
	"runtime"
	"sync"
	"time"

	"github.com/spf13/viper"
)

func init() { 
	rand.Seed(time.Now().UnixNano())
}


type Equals func(a interface{}, b interface{}) bool

var viperLock sync.RWMutex


func Contains(s string, a []string) bool {
	for _, e := range a {
		if e == s {
			return true
		}
	}
	return false
}


func IndexInSlice(array interface{}, o interface{}, equals Equals) int {
	arr := reflect.ValueOf(array)
	for i := 0; i < arr.Len(); i++ {
		if equals(arr.Index(i).Interface(), o) {
			return i
		}
	}
	return -1
}



func GetRandomIndices(indiceCount, highestIndex int) []int {
	
	if highestIndex+1 < indiceCount {
		return nil
	}

	return rand.Perm(highestIndex + 1)[:indiceCount]
}



type Set struct {
	items map[interface{}]struct{}
	lock  *sync.RWMutex
}


func NewSet() *Set {
	return &Set{lock: &sync.RWMutex{}, items: make(map[interface{}]struct{})}
}


func (s *Set) Add(item interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.items[item] = struct{}{}
}


func (s *Set) Exists(item interface{}) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
	_, exists := s.items[item]
	return exists
}


func (s *Set) Size() int {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return len(s.items)
}



func (s *Set) ToArray() []interface{} {
	s.lock.RLock()
	defer s.lock.RUnlock()
	a := make([]interface{}, len(s.items))
	i := 0
	for item := range s.items {
		a[i] = item
		i++
	}
	return a
}


func (s *Set) Clear() {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.items = make(map[interface{}]struct{})
}


func (s *Set) Remove(item interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.items, item)
}



func PrintStackTrace() {
	buf := make([]byte, 1<<16)
	l := runtime.Stack(buf, true)
	fmt.Printf("%s", buf[:l])
}


func GetIntOrDefault(key string, defVal int) int {
	viperLock.RLock()
	defer viperLock.RUnlock()

	if val := viper.GetInt(key); val != 0 {
		return val
	}

	return defVal
}


func GetFloat64OrDefault(key string, defVal float64) float64 {
	viperLock.RLock()
	defer viperLock.RUnlock()

	if val := viper.GetFloat64(key); val != 0 {
		return val
	}

	return defVal
}


func GetDurationOrDefault(key string, defVal time.Duration) time.Duration {
	viperLock.RLock()
	defer viperLock.RUnlock()

	if val := viper.GetDuration(key); val != 0 {
		return val
	}

	return defVal
}


func SetVal(key string, val interface{}) {
	viperLock.Lock()
	defer viperLock.Unlock()
	viper.Set(key, val)
}



func RandomInt(n int) int {
	return rand.Intn(n)
}





func RandomUInt64() uint64 {
	return rand.Uint64()
}

func BytesToStrings(bytes [][]byte) []string {
	strings := make([]string, len(bytes))
	for i, b := range bytes {
		strings[i] = string(b)
	}
	return strings
}

func StringsToBytes(strings []string) [][]byte {
	bytes := make([][]byte, len(strings))
	for i, str := range strings {
		bytes[i] = []byte(str)
	}
	return bytes
}
