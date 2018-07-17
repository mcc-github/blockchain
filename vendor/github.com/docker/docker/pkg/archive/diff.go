package archive 

import (
	"archive/tar"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/docker/docker/pkg/idtools"
	"github.com/docker/docker/pkg/pools"
	"github.com/docker/docker/pkg/system"
	"github.com/sirupsen/logrus"
)




func UnpackLayer(dest string, layer io.Reader, options *TarOptions) (size int64, err error) {
	tr := tar.NewReader(layer)
	trBuf := pools.BufioReader32KPool.Get(tr)
	defer pools.BufioReader32KPool.Put(trBuf)

	var dirs []*tar.Header
	unpackedPaths := make(map[string]struct{})

	if options == nil {
		options = &TarOptions{}
	}
	if options.ExcludePatterns == nil {
		options.ExcludePatterns = []string{}
	}
	idMappings := idtools.NewIDMappingsFromMaps(options.UIDMaps, options.GIDMaps)

	aufsTempdir := ""
	aufsHardlinks := make(map[string]*tar.Header)

	
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			
			break
		}
		if err != nil {
			return 0, err
		}

		size += hdr.Size

		
		hdr.Name = filepath.Clean(hdr.Name)

		
		
		
		
		
		
		
		
		
		
		
		
		
		
		if runtime.GOOS == "windows" {
			if strings.Contains(hdr.Name, ":") {
				logrus.Warnf("Windows: Ignoring %s (is this a Linux image?)", hdr.Name)
				continue
			}
		}

		
		if !strings.HasSuffix(hdr.Name, string(os.PathSeparator)) {
			
			
			
			parent := filepath.Dir(hdr.Name)
			parentPath := filepath.Join(dest, parent)

			if _, err := os.Lstat(parentPath); err != nil && os.IsNotExist(err) {
				err = system.MkdirAll(parentPath, 0600, "")
				if err != nil {
					return 0, err
				}
			}
		}

		
		if strings.HasPrefix(hdr.Name, WhiteoutMetaPrefix) {
			
			
			
			if strings.HasPrefix(hdr.Name, WhiteoutLinkDir) && hdr.Typeflag == tar.TypeReg {
				basename := filepath.Base(hdr.Name)
				aufsHardlinks[basename] = hdr
				if aufsTempdir == "" {
					if aufsTempdir, err = ioutil.TempDir("", "dockerplnk"); err != nil {
						return 0, err
					}
					defer os.RemoveAll(aufsTempdir)
				}
				if err := createTarFile(filepath.Join(aufsTempdir, basename), dest, hdr, tr, true, nil, options.InUserNS); err != nil {
					return 0, err
				}
			}

			if hdr.Name != WhiteoutOpaqueDir {
				continue
			}
		}
		path := filepath.Join(dest, hdr.Name)
		rel, err := filepath.Rel(dest, path)
		if err != nil {
			return 0, err
		}

		
		if strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
			return 0, breakoutError(fmt.Errorf("%q is outside of %q", hdr.Name, dest))
		}
		base := filepath.Base(path)

		if strings.HasPrefix(base, WhiteoutPrefix) {
			dir := filepath.Dir(path)
			if base == WhiteoutOpaqueDir {
				_, err := os.Lstat(dir)
				if err != nil {
					return 0, err
				}
				err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
					if err != nil {
						if os.IsNotExist(err) {
							err = nil 
						}
						return err
					}
					if path == dir {
						return nil
					}
					if _, exists := unpackedPaths[path]; !exists {
						err := os.RemoveAll(path)
						return err
					}
					return nil
				})
				if err != nil {
					return 0, err
				}
			} else {
				originalBase := base[len(WhiteoutPrefix):]
				originalPath := filepath.Join(dir, originalBase)
				if err := os.RemoveAll(originalPath); err != nil {
					return 0, err
				}
			}
		} else {
			
			
			
			
			if fi, err := os.Lstat(path); err == nil {
				if !(fi.IsDir() && hdr.Typeflag == tar.TypeDir) {
					if err := os.RemoveAll(path); err != nil {
						return 0, err
					}
				}
			}

			trBuf.Reset(tr)
			srcData := io.Reader(trBuf)
			srcHdr := hdr

			
			
			if hdr.Typeflag == tar.TypeLink && strings.HasPrefix(filepath.Clean(hdr.Linkname), WhiteoutLinkDir) {
				linkBasename := filepath.Base(hdr.Linkname)
				srcHdr = aufsHardlinks[linkBasename]
				if srcHdr == nil {
					return 0, fmt.Errorf("Invalid aufs hardlink")
				}
				tmpFile, err := os.Open(filepath.Join(aufsTempdir, linkBasename))
				if err != nil {
					return 0, err
				}
				defer tmpFile.Close()
				srcData = tmpFile
			}

			if err := remapIDs(idMappings, srcHdr); err != nil {
				return 0, err
			}

			if err := createTarFile(path, dest, srcHdr, srcData, true, nil, options.InUserNS); err != nil {
				return 0, err
			}

			
			
			if hdr.Typeflag == tar.TypeDir {
				dirs = append(dirs, hdr)
			}
			unpackedPaths[path] = struct{}{}
		}
	}

	for _, hdr := range dirs {
		path := filepath.Join(dest, hdr.Name)
		if err := system.Chtimes(path, hdr.AccessTime, hdr.ModTime); err != nil {
			return 0, err
		}
	}

	return size, nil
}





func ApplyLayer(dest string, layer io.Reader) (int64, error) {
	return applyLayerHandler(dest, layer, &TarOptions{}, true)
}





func ApplyUncompressedLayer(dest string, layer io.Reader, options *TarOptions) (int64, error) {
	return applyLayerHandler(dest, layer, options, false)
}


func applyLayerHandler(dest string, layer io.Reader, options *TarOptions, decompress bool) (int64, error) {
	dest = filepath.Clean(dest)

	
	oldmask, err := system.Umask(0)
	if err != nil {
		return 0, err
	}
	defer system.Umask(oldmask) 

	if decompress {
		layer, err = DecompressStream(layer)
		if err != nil {
			return 0, err
		}
	}
	return UnpackLayer(dest, layer, options)
}
