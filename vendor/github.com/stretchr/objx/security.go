package objx

import (
	"crypto/sha1"
	"encoding/hex"
)


func HashWithKey(data, key string) string {
	d := sha1.Sum([]byte(data + ":" + key))
	return hex.EncodeToString(d[:])
}
