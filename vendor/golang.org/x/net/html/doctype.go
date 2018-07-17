



package html

import (
	"strings"
)






func parseDoctype(s string) (n *Node, quirks bool) {
	n = &Node{Type: DoctypeNode}

	
	space := strings.IndexAny(s, whitespace)
	if space == -1 {
		space = len(s)
	}
	n.Data = s[:space]
	
	if n.Data != "html" {
		quirks = true
	}
	n.Data = strings.ToLower(n.Data)
	s = strings.TrimLeft(s[space:], whitespace)

	if len(s) < 6 {
		
		
		return n, quirks || s != ""
	}

	key := strings.ToLower(s[:6])
	s = s[6:]
	for key == "public" || key == "system" {
		s = strings.TrimLeft(s, whitespace)
		if s == "" {
			break
		}
		quote := s[0]
		if quote != '"' && quote != '\'' {
			break
		}
		s = s[1:]
		q := strings.IndexRune(s, rune(quote))
		var id string
		if q == -1 {
			id = s
			s = ""
		} else {
			id = s[:q]
			s = s[q+1:]
		}
		n.Attr = append(n.Attr, Attribute{Key: key, Val: id})
		if key == "public" {
			key = "system"
		} else {
			key = ""
		}
	}

	if key != "" || s != "" {
		quirks = true
	} else if len(n.Attr) > 0 {
		if n.Attr[0].Key == "public" {
			public := strings.ToLower(n.Attr[0].Val)
			switch public {
			case "-
				quirks = true
			default:
				for _, q := range quirkyIDs {
					if strings.HasPrefix(public, q) {
						quirks = true
						break
					}
				}
			}
			
			if len(n.Attr) == 1 && (strings.HasPrefix(public, "-
				strings.HasPrefix(public, "-
				quirks = true
			}
		}
		if lastAttr := n.Attr[len(n.Attr)-1]; lastAttr.Key == "system" &&
			strings.ToLower(lastAttr.Val) == "http://www.ibm.com/data/dtd/v11/ibmxhtml1-transitional.dtd" {
			quirks = true
		}
	}

	return n, quirks
}



var quirkyIDs = []string{
	"+
	"-
	"-
	"-
	"-
	"-
	"-
	"-
	"-
	"-
	"-
	"-
	"-
	"-
	"-
	"-
	"-
	"-
	"-
	"-
	"-
	"-
	"-
	"-
	"-
	"-
	"-
	"-
	"-
	"-
	"-
	"-
	"-
	"-
	"-
	"-
	"-
	"-
	"-
	"-
	"-
	"-
	"-
	"-
	"-
	"-
	"-
	"-
	"-
	"-
	"-
	"-
	"-
	"-
	"-
}
