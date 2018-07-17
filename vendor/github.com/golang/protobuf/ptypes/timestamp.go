






























package ptypes



import (
	"errors"
	"fmt"
	"time"

	tspb "github.com/golang/protobuf/ptypes/timestamp"
)

const (
	
	
	minValidSeconds = -62135596800
	
	
	maxValidSeconds = 253402300800
)











func validateTimestamp(ts *tspb.Timestamp) error {
	if ts == nil {
		return errors.New("timestamp: nil Timestamp")
	}
	if ts.Seconds < minValidSeconds {
		return fmt.Errorf("timestamp: %v before 0001-01-01", ts)
	}
	if ts.Seconds >= maxValidSeconds {
		return fmt.Errorf("timestamp: %v after 10000-01-01", ts)
	}
	if ts.Nanos < 0 || ts.Nanos >= 1e9 {
		return fmt.Errorf("timestamp: %v: nanos not in range [0, 1e9)", ts)
	}
	return nil
}












func Timestamp(ts *tspb.Timestamp) (time.Time, error) {
	
	
	var t time.Time
	if ts == nil {
		t = time.Unix(0, 0).UTC() 
	} else {
		t = time.Unix(ts.Seconds, int64(ts.Nanos)).UTC()
	}
	return t, validateTimestamp(ts)
}



func TimestampProto(t time.Time) (*tspb.Timestamp, error) {
	seconds := t.Unix()
	nanos := int32(t.Sub(time.Unix(seconds, 0)))
	ts := &tspb.Timestamp{
		Seconds: seconds,
		Nanos:   nanos,
	}
	if err := validateTimestamp(ts); err != nil {
		return nil, err
	}
	return ts, nil
}



func TimestampString(ts *tspb.Timestamp) string {
	t, err := Timestamp(ts)
	if err != nil {
		return fmt.Sprintf("(%v)", err)
	}
	return t.Format(time.RFC3339Nano)
}
