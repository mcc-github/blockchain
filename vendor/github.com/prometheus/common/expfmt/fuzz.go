















package expfmt

import "bytes"







func Fuzz(in []byte) int {
	parser := TextParser{}
	_, err := parser.TextToMetricFamilies(bytes.NewReader(in))

	if err != nil {
		return 0
	}

	return 1
}
