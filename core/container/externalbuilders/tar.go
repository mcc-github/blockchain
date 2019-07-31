/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package externalbuilders

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)





func Untar(buffer io.Reader, dst string) error {
	gzr, err := gzip.NewReader(buffer)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()

		if err == io.EOF {
			return nil
		}

		if err != nil {
			return errors.WithMessage(err, "could not get next tar element")
		}

		if !ValidPath(header.Name) {
			return errors.Errorf("tar contains the absolute or escaping path '%s'", header.Name)
		}

		target := filepath.Join(dst, header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0700); err != nil {
				return errors.WithMessagef(err, "could not create directory '%s'", header.Name)
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0700); err != nil {
				return errors.WithMessagef(err, "could not create directory '%s'", filepath.Dir(header.Name))
			}

			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return errors.WithMessagef(err, "could not create file '%s'", header.Name)
			}

			
			if _, err := io.Copy(f, tr); err != nil {
				return err
			}

			f.Close()
		default:
			return errors.Errorf("invalid file type '%v' contained in archive for file '%s'", header.Typeflag, header.Name)
		}
	}
}



func ValidPath(uncleanPath string) bool {
	
	
	sanitizedPath := filepath.Clean(uncleanPath)

	switch {
	case filepath.IsAbs(sanitizedPath):
		return false
	case strings.HasPrefix(sanitizedPath, ".."+string(filepath.Separator)) || sanitizedPath == "..":
		
		return false
	default:
		
		return true
	}
}
