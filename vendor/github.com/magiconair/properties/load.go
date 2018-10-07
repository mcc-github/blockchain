



package properties

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)


type Encoding uint

const (
	
	
	
	
	utf8Default Encoding = iota

	
	UTF8

	
	ISO_8859_1
)

type Loader struct {
	
	
	
	Encoding Encoding

	
	
	
	
	DisableExpansion bool

	
	
	
	IgnoreMissing bool
}


func (l *Loader) LoadBytes(buf []byte) (*Properties, error) {
	return l.loadBytes(buf, l.Encoding)
}






func (l *Loader) LoadAll(names []string) (*Properties, error) {
	all := NewProperties()
	for _, name := range names {
		n, err := expandName(name)
		if err != nil {
			return nil, err
		}

		var p *Properties
		switch {
		case strings.HasPrefix(n, "http://"):
			p, err = l.LoadURL(n)
		case strings.HasPrefix(n, "https://"):
			p, err = l.LoadURL(n)
		default:
			p, err = l.LoadFile(n)
		}
		if err != nil {
			return nil, err
		}
		all.Merge(p)
	}

	all.DisableExpansion = l.DisableExpansion
	if all.DisableExpansion {
		return all, nil
	}
	return all, all.check()
}




func (l *Loader) LoadFile(filename string) (*Properties, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		if l.IgnoreMissing && os.IsNotExist(err) {
			LogPrintf("properties: %s not found. skipping", filename)
			return NewProperties(), nil
		}
		return nil, err
	}
	return l.loadBytes(data, l.Encoding)
}









func (l *Loader) LoadURL(url string) (*Properties, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("properties: error fetching %q. %s", url, err)
	}

	if resp.StatusCode == 404 && l.IgnoreMissing {
		LogPrintf("properties: %s returned %d. skipping", url, resp.StatusCode)
		return NewProperties(), nil
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("properties: %s returned %d", url, resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("properties: %s error reading response. %s", url, err)
	}
	defer resp.Body.Close()

	ct := resp.Header.Get("Content-Type")
	var enc Encoding
	switch strings.ToLower(ct) {
	case "text/plain", "text/plain; charset=iso-8859-1", "text/plain; charset=latin1":
		enc = ISO_8859_1
	case "", "text/plain; charset=utf-8":
		enc = UTF8
	default:
		return nil, fmt.Errorf("properties: invalid content type %s", ct)
	}

	return l.loadBytes(body, enc)
}

func (l *Loader) loadBytes(buf []byte, enc Encoding) (*Properties, error) {
	p, err := parse(convert(buf, enc))
	if err != nil {
		return nil, err
	}
	p.DisableExpansion = l.DisableExpansion
	if p.DisableExpansion {
		return p, nil
	}
	return p, p.check()
}


func Load(buf []byte, enc Encoding) (*Properties, error) {
	l := &Loader{Encoding: enc}
	return l.LoadBytes(buf)
}


func LoadString(s string) (*Properties, error) {
	l := &Loader{Encoding: UTF8}
	return l.LoadBytes([]byte(s))
}


func LoadMap(m map[string]string) *Properties {
	p := NewProperties()
	for k, v := range m {
		p.Set(k, v)
	}
	return p
}


func LoadFile(filename string, enc Encoding) (*Properties, error) {
	l := &Loader{Encoding: enc}
	return l.LoadAll([]string{filename})
}




func LoadFiles(filenames []string, enc Encoding, ignoreMissing bool) (*Properties, error) {
	l := &Loader{Encoding: enc, IgnoreMissing: ignoreMissing}
	return l.LoadAll(filenames)
}



func LoadURL(url string) (*Properties, error) {
	l := &Loader{Encoding: UTF8}
	return l.LoadAll([]string{url})
}





func LoadURLs(urls []string, ignoreMissing bool) (*Properties, error) {
	l := &Loader{Encoding: UTF8, IgnoreMissing: ignoreMissing}
	return l.LoadAll(urls)
}





func LoadAll(names []string, enc Encoding, ignoreMissing bool) (*Properties, error) {
	l := &Loader{Encoding: enc, IgnoreMissing: ignoreMissing}
	return l.LoadAll(names)
}



func MustLoadString(s string) *Properties {
	return must(LoadString(s))
}



func MustLoadFile(filename string, enc Encoding) *Properties {
	return must(LoadFile(filename, enc))
}




func MustLoadFiles(filenames []string, enc Encoding, ignoreMissing bool) *Properties {
	return must(LoadFiles(filenames, enc, ignoreMissing))
}



func MustLoadURL(url string) *Properties {
	return must(LoadURL(url))
}




func MustLoadURLs(urls []string, ignoreMissing bool) *Properties {
	return must(LoadURLs(urls, ignoreMissing))
}





func MustLoadAll(names []string, enc Encoding, ignoreMissing bool) *Properties {
	return must(LoadAll(names, enc, ignoreMissing))
}

func must(p *Properties, err error) *Properties {
	if err != nil {
		ErrorHandler(err)
	}
	return p
}





func expandName(name string) (string, error) {
	return expand(name, []string{}, "${", "}", make(map[string]string))
}




func convert(buf []byte, enc Encoding) string {
	switch enc {
	case utf8Default, UTF8:
		return string(buf)
	case ISO_8859_1:
		runes := make([]rune, len(buf))
		for i, b := range buf {
			runes[i] = rune(b)
		}
		return string(runes)
	default:
		ErrorHandler(fmt.Errorf("unsupported encoding %v", enc))
	}
	panic("ErrorHandler should exit")
}
