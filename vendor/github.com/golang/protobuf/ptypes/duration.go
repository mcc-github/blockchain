






























package ptypes




import (
	"errors"
	"fmt"
	"time"

	durpb "github.com/golang/protobuf/ptypes/duration"
)

const (
	
	
	maxSeconds = int64(10000 * 365.25 * 24 * 60 * 60)
	minSeconds = -maxSeconds
)





func validateDuration(d *durpb.Duration) error {
	if d == nil {
		return errors.New("duration: nil Duration")
	}
	if d.Seconds < minSeconds || d.Seconds > maxSeconds {
		return fmt.Errorf("duration: %v: seconds out of range", d)
	}
	if d.Nanos <= -1e9 || d.Nanos >= 1e9 {
		return fmt.Errorf("duration: %v: nanos out of range", d)
	}
	
	if (d.Seconds < 0 && d.Nanos > 0) || (d.Seconds > 0 && d.Nanos < 0) {
		return fmt.Errorf("duration: %v: seconds and nanos have different signs", d)
	}
	return nil
}




func Duration(p *durpb.Duration) (time.Duration, error) {
	if err := validateDuration(p); err != nil {
		return 0, err
	}
	d := time.Duration(p.Seconds) * time.Second
	if int64(d/time.Second) != p.Seconds {
		return 0, fmt.Errorf("duration: %v is out of range for time.Duration", p)
	}
	if p.Nanos != 0 {
		d += time.Duration(p.Nanos)
		if (d < 0) != (p.Nanos < 0) {
			return 0, fmt.Errorf("duration: %v is out of range for time.Duration", p)
		}
	}
	return d, nil
}


func DurationProto(d time.Duration) *durpb.Duration {
	nanos := d.Nanoseconds()
	secs := nanos / 1e9
	nanos -= secs * 1e9
	return &durpb.Duration{
		Seconds: secs,
		Nanos:   int32(nanos),
	}
}
