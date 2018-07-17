



package mux

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

type routeRegexpOptions struct {
	strictSlash    bool
	useEncodedPath bool
}

type regexpType int

const (
	regexpTypePath   regexpType = 0
	regexpTypeHost   regexpType = 1
	regexpTypePrefix regexpType = 2
	regexpTypeQuery  regexpType = 3
)











func newRouteRegexp(tpl string, typ regexpType, options routeRegexpOptions) (*routeRegexp, error) {
	
	idxs, errBraces := braceIndices(tpl)
	if errBraces != nil {
		return nil, errBraces
	}
	
	template := tpl
	
	defaultPattern := "[^/]+"
	if typ == regexpTypeQuery {
		defaultPattern = ".*"
	} else if typ == regexpTypeHost {
		defaultPattern = "[^.]+"
	}
	
	if typ != regexpTypePath {
		options.strictSlash = false
	}
	
	endSlash := false
	if options.strictSlash && strings.HasSuffix(tpl, "/") {
		tpl = tpl[:len(tpl)-1]
		endSlash = true
	}
	varsN := make([]string, len(idxs)/2)
	varsR := make([]*regexp.Regexp, len(idxs)/2)
	pattern := bytes.NewBufferString("")
	pattern.WriteByte('^')
	reverse := bytes.NewBufferString("")
	var end int
	var err error
	for i := 0; i < len(idxs); i += 2 {
		
		raw := tpl[end:idxs[i]]
		end = idxs[i+1]
		parts := strings.SplitN(tpl[idxs[i]+1:end-1], ":", 2)
		name := parts[0]
		patt := defaultPattern
		if len(parts) == 2 {
			patt = parts[1]
		}
		
		if name == "" || patt == "" {
			return nil, fmt.Errorf("mux: missing name or pattern in %q",
				tpl[idxs[i]:end])
		}
		
		fmt.Fprintf(pattern, "%s(?P<%s>%s)", regexp.QuoteMeta(raw), varGroupName(i/2), patt)

		
		fmt.Fprintf(reverse, "%s%%s", raw)

		
		varsN[i/2] = name
		varsR[i/2], err = regexp.Compile(fmt.Sprintf("^%s$", patt))
		if err != nil {
			return nil, err
		}
	}
	
	raw := tpl[end:]
	pattern.WriteString(regexp.QuoteMeta(raw))
	if options.strictSlash {
		pattern.WriteString("[/]?")
	}
	if typ == regexpTypeQuery {
		
		if queryVal := strings.SplitN(template, "=", 2)[1]; queryVal == "" {
			pattern.WriteString(defaultPattern)
		}
	}
	if typ != regexpTypePrefix {
		pattern.WriteByte('$')
	}
	reverse.WriteString(raw)
	if endSlash {
		reverse.WriteByte('/')
	}
	
	reg, errCompile := regexp.Compile(pattern.String())
	if errCompile != nil {
		return nil, errCompile
	}

	
	if reg.NumSubexp() != len(idxs)/2 {
		panic(fmt.Sprintf("route %s contains capture groups in its regexp. ", template) +
			"Only non-capturing groups are accepted: e.g. (?:pattern) instead of (pattern)")
	}

	
	return &routeRegexp{
		template:   template,
		regexpType: typ,
		options:    options,
		regexp:     reg,
		reverse:    reverse.String(),
		varsN:      varsN,
		varsR:      varsR,
	}, nil
}



type routeRegexp struct {
	
	template string
	
	regexpType regexpType
	
	options routeRegexpOptions
	
	regexp *regexp.Regexp
	
	reverse string
	
	varsN []string
	
	varsR []*regexp.Regexp
}


func (r *routeRegexp) Match(req *http.Request, match *RouteMatch) bool {
	if r.regexpType != regexpTypeHost {
		if r.regexpType == regexpTypeQuery {
			return r.matchQueryString(req)
		}
		path := req.URL.Path
		if r.options.useEncodedPath {
			path = req.URL.EscapedPath()
		}
		return r.regexp.MatchString(path)
	}

	return r.regexp.MatchString(getHost(req))
}


func (r *routeRegexp) url(values map[string]string) (string, error) {
	urlValues := make([]interface{}, len(r.varsN))
	for k, v := range r.varsN {
		value, ok := values[v]
		if !ok {
			return "", fmt.Errorf("mux: missing route variable %q", v)
		}
		if r.regexpType == regexpTypeQuery {
			value = url.QueryEscape(value)
		}
		urlValues[k] = value
	}
	rv := fmt.Sprintf(r.reverse, urlValues...)
	if !r.regexp.MatchString(rv) {
		
		
		
		for k, v := range r.varsN {
			if !r.varsR[k].MatchString(values[v]) {
				return "", fmt.Errorf(
					"mux: variable %q doesn't match, expected %q", values[v],
					r.varsR[k].String())
			}
		}
	}
	return rv, nil
}




func (r *routeRegexp) getURLQuery(req *http.Request) string {
	if r.regexpType != regexpTypeQuery {
		return ""
	}
	templateKey := strings.SplitN(r.template, "=", 2)[0]
	for key, vals := range req.URL.Query() {
		if key == templateKey && len(vals) > 0 {
			return key + "=" + vals[0]
		}
	}
	return ""
}

func (r *routeRegexp) matchQueryString(req *http.Request) bool {
	return r.regexp.MatchString(r.getURLQuery(req))
}



func braceIndices(s string) ([]int, error) {
	var level, idx int
	var idxs []int
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '{':
			if level++; level == 1 {
				idx = i
			}
		case '}':
			if level--; level == 0 {
				idxs = append(idxs, idx, i+1)
			} else if level < 0 {
				return nil, fmt.Errorf("mux: unbalanced braces in %q", s)
			}
		}
	}
	if level != 0 {
		return nil, fmt.Errorf("mux: unbalanced braces in %q", s)
	}
	return idxs, nil
}


func varGroupName(idx int) string {
	return "v" + strconv.Itoa(idx)
}






type routeRegexpGroup struct {
	host    *routeRegexp
	path    *routeRegexp
	queries []*routeRegexp
}


func (v *routeRegexpGroup) setMatch(req *http.Request, m *RouteMatch, r *Route) {
	
	if v.host != nil {
		host := getHost(req)
		matches := v.host.regexp.FindStringSubmatchIndex(host)
		if len(matches) > 0 {
			extractVars(host, matches, v.host.varsN, m.Vars)
		}
	}
	path := req.URL.Path
	if r.useEncodedPath {
		path = req.URL.EscapedPath()
	}
	
	if v.path != nil {
		matches := v.path.regexp.FindStringSubmatchIndex(path)
		if len(matches) > 0 {
			extractVars(path, matches, v.path.varsN, m.Vars)
			
			if v.path.options.strictSlash {
				p1 := strings.HasSuffix(path, "/")
				p2 := strings.HasSuffix(v.path.template, "/")
				if p1 != p2 {
					u, _ := url.Parse(req.URL.String())
					if p1 {
						u.Path = u.Path[:len(u.Path)-1]
					} else {
						u.Path += "/"
					}
					m.Handler = http.RedirectHandler(u.String(), 301)
				}
			}
		}
	}
	
	for _, q := range v.queries {
		queryURL := q.getURLQuery(req)
		matches := q.regexp.FindStringSubmatchIndex(queryURL)
		if len(matches) > 0 {
			extractVars(queryURL, matches, q.varsN, m.Vars)
		}
	}
}


func getHost(r *http.Request) string {
	if r.URL.IsAbs() {
		return r.URL.Host
	}
	host := r.Host
	
	if i := strings.Index(host, ":"); i != -1 {
		host = host[:i]
	}
	return host

}

func extractVars(input string, matches []int, names []string, output map[string]string) {
	for i, name := range names {
		output[name] = input[matches[2*i+2]:matches[2*i+3]]
	}
}
