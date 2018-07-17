






























package proto



import (
	"errors"
	"fmt"
	"time"
)

const (
	
	
	minValidSeconds = -62135596800
	
	
	maxValidSeconds = 253402300800
)











func validateTimestamp(ts *timestamp) error {
	if ts == nil {
		return errors.New("timestamp: nil Timestamp")
	}
	if ts.Seconds < minValidSeconds {
		return fmt.Errorf("timestamp: %#v before 0001-01-01", ts)
	}
	if ts.Seconds >= maxValidSeconds {
		return fmt.Errorf("timestamp: %#v after 10000-01-01", ts)
	}
	if ts.Nanos < 0 || ts.Nanos >= 1e9 {
		return fmt.Errorf("timestamp: %#v: nanos not in range [0, 1e9)", ts)
	}
	return nil
}












func timestampFromProto(ts *timestamp) (time.Time, error) {
	
	
	var t time.Time
	if ts == nil {
		t = time.Unix(0, 0).UTC() 
	} else {
		t = time.Unix(ts.Seconds, int64(ts.Nanos)).UTC()
	}
	return t, validateTimestamp(ts)
}



func timestampProto(t time.Time) (*timestamp, error) {
	seconds := t.Unix()
	nanos := int32(t.Sub(time.Unix(seconds, 0)))
	ts := &timestamp{
		Seconds: seconds,
		Nanos:   nanos,
	}
	if err := validateTimestamp(ts); err != nil {
		return nil, err
	}
	return ts, nil
}
