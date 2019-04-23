

package humanize

import (
	"bytes"
	"math/big"
	"strings"
)



func BigCommaf(v *big.Float) string {
	buf := &bytes.Buffer{}
	if v.Sign() < 0 {
		buf.Write([]byte{'-'})
		v.Abs(v)
	}

	comma := []byte{','}

	parts := strings.Split(v.Text('f', -1), ".")
	pos := 0
	if len(parts[0])%3 != 0 {
		pos += len(parts[0]) % 3
		buf.WriteString(parts[0][:pos])
		buf.Write(comma)
	}
	for ; pos < len(parts[0]); pos += 3 {
		buf.WriteString(parts[0][pos : pos+3])
		buf.Write(comma)
	}
	buf.Truncate(buf.Len() - 1)

	if len(parts) > 1 {
		buf.Write([]byte{'.'})
		buf.WriteString(parts[1])
	}
	return buf.String()
}
