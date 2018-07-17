












package model

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	
	
	minimumTick = time.Millisecond
	
	second = int64(time.Second / minimumTick)
	
	nanosPerTick = int64(minimumTick / time.Nanosecond)

	
	
	Earliest = Time(math.MinInt64)
	
	
	Latest = Time(math.MaxInt64)
)



type Time int64


type Interval struct {
	Start, End Time
}


func Now() Time {
	return TimeFromUnixNano(time.Now().UnixNano())
}



func TimeFromUnix(t int64) Time {
	return Time(t * second)
}



func TimeFromUnixNano(t int64) Time {
	return Time(t / nanosPerTick)
}


func (t Time) Equal(o Time) bool {
	return t == o
}


func (t Time) Before(o Time) bool {
	return t < o
}


func (t Time) After(o Time) bool {
	return t > o
}


func (t Time) Add(d time.Duration) Time {
	return t + Time(d/minimumTick)
}


func (t Time) Sub(o Time) time.Duration {
	return time.Duration(t-o) * minimumTick
}


func (t Time) Time() time.Time {
	return time.Unix(int64(t)/second, (int64(t)%second)*nanosPerTick)
}



func (t Time) Unix() int64 {
	return int64(t) / second
}



func (t Time) UnixNano() int64 {
	return int64(t) * nanosPerTick
}


var dotPrecision = int(math.Log10(float64(second)))


func (t Time) String() string {
	return strconv.FormatFloat(float64(t)/float64(second), 'f', -1, 64)
}


func (t Time) MarshalJSON() ([]byte, error) {
	return []byte(t.String()), nil
}


func (t *Time) UnmarshalJSON(b []byte) error {
	p := strings.Split(string(b), ".")
	switch len(p) {
	case 1:
		v, err := strconv.ParseInt(string(p[0]), 10, 64)
		if err != nil {
			return err
		}
		*t = Time(v * second)

	case 2:
		v, err := strconv.ParseInt(string(p[0]), 10, 64)
		if err != nil {
			return err
		}
		v *= second

		prec := dotPrecision - len(p[1])
		if prec < 0 {
			p[1] = p[1][:dotPrecision]
		} else if prec > 0 {
			p[1] = p[1] + strings.Repeat("0", prec)
		}

		va, err := strconv.ParseInt(p[1], 10, 32)
		if err != nil {
			return err
		}

		*t = Time(v + va)

	default:
		return fmt.Errorf("invalid time %q", string(b))
	}
	return nil
}




type Duration time.Duration


func (d *Duration) Set(s string) error {
	var err error
	*d, err = ParseDuration(s)
	return err
}


func (d *Duration) Type() string {
	return "duration"
}

var durationRE = regexp.MustCompile("^([0-9]+)(y|w|d|h|m|s|ms)$")



func ParseDuration(durationStr string) (Duration, error) {
	matches := durationRE.FindStringSubmatch(durationStr)
	if len(matches) != 3 {
		return 0, fmt.Errorf("not a valid duration string: %q", durationStr)
	}
	var (
		n, _ = strconv.Atoi(matches[1])
		dur  = time.Duration(n) * time.Millisecond
	)
	switch unit := matches[2]; unit {
	case "y":
		dur *= 1000 * 60 * 60 * 24 * 365
	case "w":
		dur *= 1000 * 60 * 60 * 24 * 7
	case "d":
		dur *= 1000 * 60 * 60 * 24
	case "h":
		dur *= 1000 * 60 * 60
	case "m":
		dur *= 1000 * 60
	case "s":
		dur *= 1000
	case "ms":
		
	default:
		return 0, fmt.Errorf("invalid time unit in duration string: %q", unit)
	}
	return Duration(dur), nil
}

func (d Duration) String() string {
	var (
		ms   = int64(time.Duration(d) / time.Millisecond)
		unit = "ms"
	)
	if ms == 0 {
		return "0s"
	}
	factors := map[string]int64{
		"y":  1000 * 60 * 60 * 24 * 365,
		"w":  1000 * 60 * 60 * 24 * 7,
		"d":  1000 * 60 * 60 * 24,
		"h":  1000 * 60 * 60,
		"m":  1000 * 60,
		"s":  1000,
		"ms": 1,
	}

	switch int64(0) {
	case ms % factors["y"]:
		unit = "y"
	case ms % factors["w"]:
		unit = "w"
	case ms % factors["d"]:
		unit = "d"
	case ms % factors["h"]:
		unit = "h"
	case ms % factors["m"]:
		unit = "m"
	case ms % factors["s"]:
		unit = "s"
	}
	return fmt.Sprintf("%v%v", ms/factors[unit], unit)
}


func (d Duration) MarshalYAML() (interface{}, error) {
	return d.String(), nil
}


func (d *Duration) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}
	dur, err := ParseDuration(s)
	if err != nil {
		return err
	}
	*d = dur
	return nil
}
