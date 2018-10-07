package objx

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/url"
	"strings"
)



type MSIConvertable interface {
	
	
	MSI() map[string]interface{}
}



type Map map[string]interface{}


func (m Map) Value() *Value {
	return &Value{data: m}
}


var Nil = New(nil)




func New(data interface{}) Map {
	if _, ok := data.(map[string]interface{}); !ok {
		if converter, ok := data.(MSIConvertable); ok {
			data = converter.MSI()
		} else {
			return nil
		}
	}
	return Map(data.(map[string]interface{}))
}
















func MSI(keyAndValuePairs ...interface{}) Map {
	newMap := Map{}
	keyAndValuePairsLen := len(keyAndValuePairs)
	if keyAndValuePairsLen%2 != 0 {
		return nil
	}
	for i := 0; i < keyAndValuePairsLen; i = i + 2 {
		key := keyAndValuePairs[i]
		value := keyAndValuePairs[i+1]

		
		keyString, keyStringOK := key.(string)
		if !keyStringOK {
			return nil
		}
		newMap[keyString] = value
	}
	return newMap
}







func MustFromJSON(jsonString string) Map {
	o, err := FromJSON(jsonString)
	if err != nil {
		panic("objx: MustFromJSON failed with error: " + err.Error())
	}
	return o
}





func FromJSON(jsonString string) (Map, error) {
	var data interface{}
	err := json.Unmarshal([]byte(jsonString), &data)
	if err != nil {
		return Nil, err
	}
	return New(data), nil
}





func FromBase64(base64String string) (Map, error) {
	decoder := base64.NewDecoder(base64.StdEncoding, strings.NewReader(base64String))
	decoded, err := ioutil.ReadAll(decoder)
	if err != nil {
		return nil, err
	}
	return FromJSON(string(decoded))
}





func MustFromBase64(base64String string) Map {
	result, err := FromBase64(base64String)
	if err != nil {
		panic("objx: MustFromBase64 failed with error: " + err.Error())
	}
	return result
}





func FromSignedBase64(base64String, key string) (Map, error) {
	parts := strings.Split(base64String, SignatureSeparator)
	if len(parts) != 2 {
		return nil, errors.New("objx: Signed base64 string is malformed")
	}

	sig := HashWithKey(parts[0], key)
	if parts[1] != sig {
		return nil, errors.New("objx: Signature for base64 data does not match")
	}
	return FromBase64(parts[0])
}





func MustFromSignedBase64(base64String, key string) Map {
	result, err := FromSignedBase64(base64String, key)
	if err != nil {
		panic("objx: MustFromSignedBase64 failed with error: " + err.Error())
	}
	return result
}





func FromURLQuery(query string) (Map, error) {
	vals, err := url.ParseQuery(query)
	if err != nil {
		return nil, err
	}
	m := Map{}
	for k, vals := range vals {
		m[k] = vals[0]
	}
	return m, nil
}







func MustFromURLQuery(query string) Map {
	o, err := FromURLQuery(query)
	if err != nil {
		panic("objx: MustFromURLQuery failed with error: " + err.Error())
	}
	return o
}
