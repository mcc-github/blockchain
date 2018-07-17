



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
	
	UTF8 Encoding = 1 << iota

	
	ISO_8859_1
)


func Load(buf []byte, enc Encoding) (*Properties, error) {
	return loadBuf(buf, enc)
}


func LoadString(s string) (*Properties, error) {
	return loadBuf([]byte(s), UTF8)
}


func LoadMap(m map[string]string) *Properties {
	p := NewProperties()
	for k, v := range m {
		p.Set(k, v)
	}
	return p
}


func LoadFile(filename string, enc Encoding) (*Properties, error) {
	return loadAll([]string{filename}, enc, false)
}




func LoadFiles(filenames []string, enc Encoding, ignoreMissing bool) (*Properties, error) {
	return loadAll(filenames, enc, ignoreMissing)
}









func LoadURL(url string) (*Properties, error) {
	return loadAll([]string{url}, UTF8, false)
}





func LoadURLs(urls []string, ignoreMissing bool) (*Properties, error) {
	return loadAll(urls, UTF8, ignoreMissing)
}





func LoadAll(names []string, enc Encoding, ignoreMissing bool) (*Properties, error) {
	return loadAll(names, enc, ignoreMissing)
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

func loadBuf(buf []byte, enc Encoding) (*Properties, error) {
	p, err := parse(convert(buf, enc))
	if err != nil {
		return nil, err
	}
	return p, p.check()
}

func loadAll(names []string, enc Encoding, ignoreMissing bool) (*Properties, error) {
	result := NewProperties()
	for _, name := range names {
		n, err := expandName(name)
		if err != nil {
			return nil, err
		}
		var p *Properties
		if strings.HasPrefix(n, "http://") || strings.HasPrefix(n, "https://") {
			p, err = loadURL(n, ignoreMissing)
		} else {
			p, err = loadFile(n, enc, ignoreMissing)
		}
		if err != nil {
			return nil, err
		}
		result.Merge(p)

	}
	return result, result.check()
}

func loadFile(filename string, enc Encoding, ignoreMissing bool) (*Properties, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		if ignoreMissing && os.IsNotExist(err) {
			LogPrintf("properties: %s not found. skipping", filename)
			return NewProperties(), nil
		}
		return nil, err
	}
	p, err := parse(convert(data, enc))
	if err != nil {
		return nil, err
	}
	return p, nil
}

func loadURL(url string, ignoreMissing bool) (*Properties, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("properties: error fetching %q. %s", url, err)
	}
	if resp.StatusCode == 404 && ignoreMissing {
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
	if err = resp.Body.Close(); err != nil {
		return nil, fmt.Errorf("properties: %s error reading response. %s", url, err)
	}

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

	p, err := parse(convert(body, enc))
	if err != nil {
		return nil, err
	}
	return p, nil
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
	case UTF8:
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
