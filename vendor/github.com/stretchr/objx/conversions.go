package objx

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
)



func (m Map) JSON() (string, error) {
	result, err := json.Marshal(m)
	if err != nil {
		err = errors.New("objx: JSON encode failed with: " + err.Error())
	}
	return string(result), err
}



func (m Map) MustJSON() string {
	result, err := m.JSON()
	if err != nil {
		panic(err.Error())
	}
	return result
}



func (m Map) Base64() (string, error) {
	var buf bytes.Buffer

	jsonData, err := m.JSON()
	if err != nil {
		return "", err
	}

	encoder := base64.NewEncoder(base64.StdEncoding, &buf)
	_, err = encoder.Write([]byte(jsonData))
	if err != nil {
		return "", err
	}
	_ = encoder.Close()

	return buf.String(), nil
}




func (m Map) MustBase64() string {
	result, err := m.Base64()
	if err != nil {
		panic(err.Error())
	}
	return result
}




func (m Map) SignedBase64(key string) (string, error) {
	base64, err := m.Base64()
	if err != nil {
		return "", err
	}

	sig := HashWithKey(base64, key)
	return base64 + SignatureSeparator + sig, nil
}




func (m Map) MustSignedBase64(key string) string {
	result, err := m.SignedBase64(key)
	if err != nil {
		panic(err.Error())
	}
	return result
}





func (m Map) URLValues() url.Values {
	vals := make(url.Values)
	for k, v := range m {
		
		vals.Set(k, fmt.Sprintf("%v", v))
	}
	return vals
}




func (m Map) URLQuery() (string, error) {
	return m.URLValues().Encode(), nil
}
