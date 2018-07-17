



package cldr

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
)


type Decoder struct {
	dirFilter     []string
	sectionFilter []string
	loader        Loader
	cldr          *CLDR
	curLocale     string
}



func (d *Decoder) SetSectionFilter(filter ...string) {
	d.sectionFilter = filter
	
}




func (d *Decoder) SetDirFilter(dir ...string) {
	d.dirFilter = dir
}


type Loader interface {
	Len() int
	Path(i int) string
	Reader(i int) (io.ReadCloser, error)
}

var fileRe = regexp.MustCompile(`.*[/\\](.*)[/\\](.*)\.xml`)


func (d *Decoder) Decode(l Loader) (cldr *CLDR, err error) {
	d.cldr = makeCLDR()
	for i := 0; i < l.Len(); i++ {
		fname := l.Path(i)
		if m := fileRe.FindStringSubmatch(fname); m != nil {
			if len(d.dirFilter) > 0 && !in(d.dirFilter, m[1]) {
				continue
			}
			var r io.Reader
			if r, err = l.Reader(i); err == nil {
				err = d.decode(m[1], m[2], r)
			}
			if err != nil {
				return nil, err
			}
		}
	}
	d.cldr.finalize(d.sectionFilter)
	return d.cldr, nil
}

func (d *Decoder) decode(dir, id string, r io.Reader) error {
	var v interface{}
	var l *LDML
	cldr := d.cldr
	switch {
	case dir == "supplemental":
		v = cldr.supp
	case dir == "transforms":
		return nil
	case dir == "bcp47":
		v = cldr.bcp47
	case dir == "validity":
		return nil
	default:
		ok := false
		if v, ok = cldr.locale[id]; !ok {
			l = &LDML{}
			v, cldr.locale[id] = l, l
		}
	}
	x := xml.NewDecoder(r)
	if err := x.Decode(v); err != nil {
		log.Printf("%s/%s: %v", dir, id, err)
		return err
	}
	if l != nil {
		if l.Identity == nil {
			return fmt.Errorf("%s/%s: missing identity element", dir, id)
		}
		
		
		
		
		
		
	}
	return nil
}

type pathLoader []string

func makePathLoader(path string) (pl pathLoader, err error) {
	err = filepath.Walk(path, func(path string, _ os.FileInfo, err error) error {
		pl = append(pl, path)
		return err
	})
	return pl, err
}

func (pl pathLoader) Len() int {
	return len(pl)
}

func (pl pathLoader) Path(i int) string {
	return pl[i]
}

func (pl pathLoader) Reader(i int) (io.ReadCloser, error) {
	return os.Open(pl[i])
}


func (d *Decoder) DecodePath(path string) (cldr *CLDR, err error) {
	loader, err := makePathLoader(path)
	if err != nil {
		return nil, err
	}
	return d.Decode(loader)
}

type zipLoader struct {
	r *zip.Reader
}

func (zl zipLoader) Len() int {
	return len(zl.r.File)
}

func (zl zipLoader) Path(i int) string {
	return zl.r.File[i].Name
}

func (zl zipLoader) Reader(i int) (io.ReadCloser, error) {
	return zl.r.File[i].Open()
}


func (d *Decoder) DecodeZip(r io.Reader) (cldr *CLDR, err error) {
	buffer, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	archive, err := zip.NewReader(bytes.NewReader(buffer), int64(len(buffer)))
	if err != nil {
		return nil, err
	}
	return d.Decode(zipLoader{archive})
}
