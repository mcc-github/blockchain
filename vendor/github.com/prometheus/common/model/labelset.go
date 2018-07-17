












package model

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)






type LabelSet map[LabelName]LabelValue



func (ls LabelSet) Validate() error {
	for ln, lv := range ls {
		if !ln.IsValid() {
			return fmt.Errorf("invalid name %q", ln)
		}
		if !lv.IsValid() {
			return fmt.Errorf("invalid value %q", lv)
		}
	}
	return nil
}


func (ls LabelSet) Equal(o LabelSet) bool {
	if len(ls) != len(o) {
		return false
	}
	for ln, lv := range ls {
		olv, ok := o[ln]
		if !ok {
			return false
		}
		if olv != lv {
			return false
		}
	}
	return true
}












func (ls LabelSet) Before(o LabelSet) bool {
	if len(ls) < len(o) {
		return true
	}
	if len(ls) > len(o) {
		return false
	}

	lns := make(LabelNames, 0, len(ls)+len(o))
	for ln := range ls {
		lns = append(lns, ln)
	}
	for ln := range o {
		lns = append(lns, ln)
	}
	
	sort.Sort(lns)
	for _, ln := range lns {
		mlv, ok := ls[ln]
		if !ok {
			return true
		}
		olv, ok := o[ln]
		if !ok {
			return false
		}
		if mlv < olv {
			return true
		}
		if mlv > olv {
			return false
		}
	}
	return false
}


func (ls LabelSet) Clone() LabelSet {
	lsn := make(LabelSet, len(ls))
	for ln, lv := range ls {
		lsn[ln] = lv
	}
	return lsn
}


func (l LabelSet) Merge(other LabelSet) LabelSet {
	result := make(LabelSet, len(l))

	for k, v := range l {
		result[k] = v
	}

	for k, v := range other {
		result[k] = v
	}

	return result
}

func (l LabelSet) String() string {
	lstrs := make([]string, 0, len(l))
	for l, v := range l {
		lstrs = append(lstrs, fmt.Sprintf("%s=%q", l, v))
	}

	sort.Strings(lstrs)
	return fmt.Sprintf("{%s}", strings.Join(lstrs, ", "))
}


func (ls LabelSet) Fingerprint() Fingerprint {
	return labelSetToFingerprint(ls)
}



func (ls LabelSet) FastFingerprint() Fingerprint {
	return labelSetToFastFingerprint(ls)
}


func (l *LabelSet) UnmarshalJSON(b []byte) error {
	var m map[LabelName]LabelValue
	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}
	
	
	
	for ln := range m {
		if !ln.IsValid() {
			return fmt.Errorf("%q is not a valid label name", ln)
		}
	}
	*l = LabelSet(m)
	return nil
}
