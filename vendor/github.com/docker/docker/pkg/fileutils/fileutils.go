package fileutils 

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/scanner"

	"github.com/sirupsen/logrus"
)


type PatternMatcher struct {
	patterns   []*Pattern
	exclusions bool
}



func NewPatternMatcher(patterns []string) (*PatternMatcher, error) {
	pm := &PatternMatcher{
		patterns: make([]*Pattern, 0, len(patterns)),
	}
	for _, p := range patterns {
		
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		p = filepath.Clean(p)
		newp := &Pattern{}
		if p[0] == '!' {
			if len(p) == 1 {
				return nil, errors.New("illegal exclusion pattern: \"!\"")
			}
			newp.exclusion = true
			p = p[1:]
			pm.exclusions = true
		}
		
		
		
		
		
		
		if _, err := filepath.Match(p, "."); err != nil {
			return nil, err
		}
		newp.cleanedPattern = p
		newp.dirs = strings.Split(p, string(os.PathSeparator))
		pm.patterns = append(pm.patterns, newp)
	}
	return pm, nil
}



func (pm *PatternMatcher) Matches(file string) (bool, error) {
	matched := false
	file = filepath.FromSlash(file)
	parentPath := filepath.Dir(file)
	parentPathDirs := strings.Split(parentPath, string(os.PathSeparator))

	for _, pattern := range pm.patterns {
		negative := false

		if pattern.exclusion {
			negative = true
		}

		match, err := pattern.match(file)
		if err != nil {
			return false, err
		}

		if !match && parentPath != "." {
			
			if len(pattern.dirs) <= len(parentPathDirs) {
				match, _ = pattern.match(strings.Join(parentPathDirs[:len(pattern.dirs)], string(os.PathSeparator)))
			}
		}

		if match {
			matched = !negative
		}
	}

	if matched {
		logrus.Debugf("Skipping excluded path: %s", file)
	}

	return matched, nil
}


func (pm *PatternMatcher) Exclusions() bool {
	return pm.exclusions
}


func (pm *PatternMatcher) Patterns() []*Pattern {
	return pm.patterns
}


type Pattern struct {
	cleanedPattern string
	dirs           []string
	regexp         *regexp.Regexp
	exclusion      bool
}

func (p *Pattern) String() string {
	return p.cleanedPattern
}


func (p *Pattern) Exclusion() bool {
	return p.exclusion
}

func (p *Pattern) match(path string) (bool, error) {

	if p.regexp == nil {
		if err := p.compile(); err != nil {
			return false, filepath.ErrBadPattern
		}
	}

	b := p.regexp.MatchString(path)

	return b, nil
}

func (p *Pattern) compile() error {
	regStr := "^"
	pattern := p.cleanedPattern
	
	
	var scan scanner.Scanner
	scan.Init(strings.NewReader(pattern))

	sl := string(os.PathSeparator)
	escSL := sl
	if sl == `\` {
		escSL += `\`
	}

	for scan.Peek() != scanner.EOF {
		ch := scan.Next()

		if ch == '*' {
			if scan.Peek() == '*' {
				
				scan.Next()

				
				if string(scan.Peek()) == sl {
					scan.Next()
				}

				if scan.Peek() == scanner.EOF {
					
					regStr += ".*"
				} else {
					
					
					
					regStr += "(.*" + escSL + ")?"
				}
			} else {
				
				regStr += "[^" + escSL + "]*"
			}
		} else if ch == '?' {
			
			regStr += "[^" + escSL + "]"
		} else if ch == '.' || ch == '$' {
			
			
			regStr += `\` + string(ch)
		} else if ch == '\\' {
			
			
			if sl == `\` {
				
				
				
				regStr += escSL
				continue
			}
			if scan.Peek() != scanner.EOF {
				regStr += `\` + string(scan.Next())
			} else {
				regStr += `\`
			}
		} else {
			regStr += string(ch)
		}
	}

	regStr += "$"

	re, err := regexp.Compile(regStr)
	if err != nil {
		return err
	}

	p.regexp = re
	return nil
}



func Matches(file string, patterns []string) (bool, error) {
	pm, err := NewPatternMatcher(patterns)
	if err != nil {
		return false, err
	}
	file = filepath.Clean(file)

	if file == "." {
		
		return false, nil
	}

	return pm.Matches(file)
}




func CopyFile(src, dst string) (int64, error) {
	cleanSrc := filepath.Clean(src)
	cleanDst := filepath.Clean(dst)
	if cleanSrc == cleanDst {
		return 0, nil
	}
	sf, err := os.Open(cleanSrc)
	if err != nil {
		return 0, err
	}
	defer sf.Close()
	if err := os.Remove(cleanDst); err != nil && !os.IsNotExist(err) {
		return 0, err
	}
	df, err := os.Create(cleanDst)
	if err != nil {
		return 0, err
	}
	defer df.Close()
	return io.Copy(df, sf)
}



func ReadSymlinkedDirectory(path string) (string, error) {
	var realPath string
	var err error
	if realPath, err = filepath.Abs(path); err != nil {
		return "", fmt.Errorf("unable to get absolute path for %s: %s", path, err)
	}
	if realPath, err = filepath.EvalSymlinks(realPath); err != nil {
		return "", fmt.Errorf("failed to canonicalise path for %s: %s", path, err)
	}
	realPathInfo, err := os.Stat(realPath)
	if err != nil {
		return "", fmt.Errorf("failed to stat target '%s' of '%s': %s", realPath, path, err)
	}
	if !realPathInfo.Mode().IsDir() {
		return "", fmt.Errorf("canonical path points to a file '%s'", realPath)
	}
	return realPath, nil
}


func CreateIfNotExists(path string, isDir bool) error {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			if isDir {
				return os.MkdirAll(path, 0755)
			}
			if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
				return err
			}
			f, err := os.OpenFile(path, os.O_CREATE, 0755)
			if err != nil {
				return err
			}
			f.Close()
		}
	}
	return nil
}
