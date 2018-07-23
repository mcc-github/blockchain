/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package testutil

import (
	"archive/tar"
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/golang/protobuf/proto"
)


func AssertNil(t testing.TB, value interface{}) {
	if !isNil(value) {
		t.Fatalf("Value not nil. value=[%#v]\n %s", value, getCallerInfo())
	}
}


func AssertNotNil(t testing.TB, value interface{}) {
	if isNil(value) {
		t.Fatalf("Values is nil. %s", getCallerInfo())
	}
}


func AssertSame(t testing.TB, actual interface{}, expected interface{}) {
	t.Logf("%s: AssertSame [%#v] and [%#v]", getCallerInfo(), actual, expected)
	if actual != expected {
		t.Fatalf("Values actual=[%#v] and expected=[%#v] do not point to same object. %s", actual, expected, getCallerInfo())
	}
}


func AssertEquals(t testing.TB, actual interface{}, expected interface{}) {
	t.Logf("%s: AssertEquals [%#v] and [%#v]", getCallerInfo(), actual, expected)
	if expected == nil && isNil(actual) {
		return
	}

	
	actual = convertJSONToMap(actual)
	expected = convertJSONToMap(expected)

	if pa, ok := actual.(proto.Message); ok {
		if pe, ok := actual.(proto.Message); ok {
			if proto.Equal(pa, pe) {
				return
			}
		}
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Values are not equal.\n Actual=[%#v], \n Expected=[%#v]\n %s", actual, expected, getCallerInfo())
	}
}


func convertJSONToMap(value interface{}) interface{} {
	var valueMap map[string]interface{}
	valueBytes, ok := value.([]byte)
	if ok {
		if json.Unmarshal(valueBytes, &valueMap) == nil {
			return valueMap
		}
	}
	return value
}


func AssertNotEquals(t testing.TB, actual interface{}, expected interface{}) {
	if reflect.DeepEqual(actual, expected) {
		t.Fatalf("Values are not supposed to be equal. Actual=[%#v], Expected=[%#v]\n %s", actual, expected, getCallerInfo())
	}
}


func AssertError(t testing.TB, err error, message string) {
	if err == nil {
		t.Fatalf("%s\n %s", message, getCallerInfo())
	}
}


func AssertNoError(t testing.TB, err error, message string) {
	if err != nil {
		t.Fatalf("%s - Error: %s\n %s", message, err, getCallerInfo())
	}
}


func AssertContains(t testing.TB, slice interface{}, value interface{}) {
	if reflect.TypeOf(slice).Kind() != reflect.Slice && reflect.TypeOf(slice).Kind() != reflect.Array {
		t.Fatalf("Type of argument 'slice' is expected to be a slice/array, found =[%s]\n %s", reflect.TypeOf(slice), getCallerInfo())
	}

	if !Contains(slice, value) {
		t.Fatalf("Expected value [%s] not found in slice %s\n %s", value, slice, getCallerInfo())
	}
}


func AssertContainsAll(t testing.TB, sliceActual interface{}, sliceExpected interface{}) {
	if reflect.TypeOf(sliceActual).Kind() != reflect.Slice && reflect.TypeOf(sliceActual).Kind() != reflect.Array {
		t.Fatalf("Type of argument 'sliceActual' is expected to be a slice/array, found =[%s]\n %s", reflect.TypeOf(sliceActual), getCallerInfo())
	}

	if reflect.TypeOf(sliceExpected).Kind() != reflect.Slice && reflect.TypeOf(sliceExpected).Kind() != reflect.Array {
		t.Fatalf("Type of argument 'sliceExpected' is expected to be a slice/array, found =[%s]\n %s", reflect.TypeOf(sliceExpected), getCallerInfo())
	}

	array := reflect.ValueOf(sliceExpected)
	for i := 0; i < array.Len(); i++ {
		element := array.Index(i).Interface()
		if !Contains(sliceActual, element) {
			t.Fatalf("Expected value [%s] not found in slice %s\n %s", element, sliceActual, getCallerInfo())
		}
	}
}


func AssertPanic(t testing.TB, msg string) {
	x := recover()
	if x == nil {
		t.Fatal(msg)
	} else {
		t.Logf("A panic was caught successfully. Actual msg = %s", x)
	}
}


func ConstructRandomBytes(t testing.TB, size int) []byte {
	value := make([]byte, size)
	_, err := rand.Read(value)
	if err != nil {
		t.Fatalf("Error while generating random bytes: %s", err)
	}
	return value
}


func Contains(slice interface{}, value interface{}) bool {
	array := reflect.ValueOf(slice)
	for i := 0; i < array.Len(); i++ {
		element := array.Index(i).Interface()
		if value == element || reflect.DeepEqual(element, value) {
			return true
		}
	}
	return false
}

func isNil(in interface{}) bool {
	return in == nil || reflect.ValueOf(in).IsNil() || (reflect.TypeOf(in).Kind() == reflect.Slice && reflect.ValueOf(in).Len() == 0)
}

func getCallerInfo() string {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		return "Could not retrieve caller's info"
	}
	return fmt.Sprintf("CallerInfo = [%s:%d]", file, line)
}


type TarFileEntry struct {
	Name, Body string
}


func CreateTarBytesForTest(testFiles []*TarFileEntry) []byte {
	
	buffer := new(bytes.Buffer)
	tarWriter := tar.NewWriter(buffer)

	for _, file := range testFiles {
		tarHeader := &tar.Header{
			Name: file.Name,
			Mode: 0600,
			Size: int64(len(file.Body)),
		}
		err := tarWriter.WriteHeader(tarHeader)
		if err != nil {
			return nil
		}
		_, err = tarWriter.Write([]byte(file.Body))
		if err != nil {
			return nil
		}
	}
	
	tarWriter.Close()
	return buffer.Bytes()
}


func CopyDir(srcroot, destroot string) error {
	_, lastSegment := filepath.Split(srcroot)
	destroot = filepath.Join(destroot, lastSegment)

	walkFunc := func(srcpath string, info os.FileInfo, err error) error {
		srcsubpath, err := filepath.Rel(srcroot, srcpath)
		if err != nil {
			return err
		}
		destpath := filepath.Join(destroot, srcsubpath)

		if info.IsDir() { 
			if err = os.MkdirAll(destpath, info.Mode()); err != nil {
				return err
			}
			return nil
		}

		
		if err = copyFile(srcpath, destpath); err != nil {
			return err
		}
		return nil
	}

	return filepath.Walk(srcroot, walkFunc)
}

func copyFile(srcpath, destpath string) error {
	var srcFile, destFile *os.File
	var err error
	if srcFile, err = os.Open(srcpath); err != nil {
		return err
	}
	if destFile, err = os.Create(destpath); err != nil {
		return err
	}
	if _, err = io.Copy(destFile, srcFile); err != nil {
		return err
	}
	if err = srcFile.Close(); err != nil {
		return err
	}
	if err = destFile.Close(); err != nil {
		return err
	}
	return nil
}