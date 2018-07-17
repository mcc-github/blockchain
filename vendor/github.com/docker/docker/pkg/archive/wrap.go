package archive 

import (
	"archive/tar"
	"bytes"
	"io"
)

















func Generate(input ...string) (io.Reader, error) {
	files := parseStringPairs(input...)
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	for _, file := range files {
		name, content := file[0], file[1]
		hdr := &tar.Header{
			Name: name,
			Size: int64(len(content)),
		}
		if err := tw.WriteHeader(hdr); err != nil {
			return nil, err
		}
		if _, err := tw.Write([]byte(content)); err != nil {
			return nil, err
		}
	}
	if err := tw.Close(); err != nil {
		return nil, err
	}
	return buf, nil
}

func parseStringPairs(input ...string) (output [][2]string) {
	output = make([][2]string, 0, len(input)/2+1)
	for i := 0; i < len(input); i += 2 {
		var pair [2]string
		pair[0] = input[i]
		if i+1 < len(input) {
			pair[1] = input[i+1]
		}
		output = append(output, pair)
	}
	return
}
