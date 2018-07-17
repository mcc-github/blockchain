package strslice 

import "encoding/json"



type StrSlice []string



func (e *StrSlice) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		
		
		
		return nil
	}

	p := make([]string, 0, 1)
	if err := json.Unmarshal(b, &p); err != nil {
		var s string
		if err := json.Unmarshal(b, &s); err != nil {
			return err
		}
		p = append(p, s)
	}

	*e = p
	return nil
}
