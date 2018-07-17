package units



type Base2Bytes int64


const (
	Kibibyte Base2Bytes = 1024
	KiB                 = Kibibyte
	Mebibyte            = Kibibyte * 1024
	MiB                 = Mebibyte
	Gibibyte            = Mebibyte * 1024
	GiB                 = Gibibyte
	Tebibyte            = Gibibyte * 1024
	TiB                 = Tebibyte
	Pebibyte            = Tebibyte * 1024
	PiB                 = Pebibyte
	Exbibyte            = Pebibyte * 1024
	EiB                 = Exbibyte
)

var (
	bytesUnitMap    = MakeUnitMap("iB", "B", 1024)
	oldBytesUnitMap = MakeUnitMap("B", "B", 1024)
)



func ParseBase2Bytes(s string) (Base2Bytes, error) {
	n, err := ParseUnit(s, bytesUnitMap)
	if err != nil {
		n, err = ParseUnit(s, oldBytesUnitMap)
	}
	return Base2Bytes(n), err
}

func (b Base2Bytes) String() string {
	return ToString(int64(b), 1024, "iB", "B")
}

var (
	metricBytesUnitMap = MakeUnitMap("B", "B", 1000)
)


type MetricBytes SI


const (
	Kilobyte MetricBytes = 1000
	KB                   = Kilobyte
	Megabyte             = Kilobyte * 1000
	MB                   = Megabyte
	Gigabyte             = Megabyte * 1000
	GB                   = Gigabyte
	Terabyte             = Gigabyte * 1000
	TB                   = Terabyte
	Petabyte             = Terabyte * 1000
	PB                   = Petabyte
	Exabyte              = Petabyte * 1000
	EB                   = Exabyte
)


func ParseMetricBytes(s string) (MetricBytes, error) {
	n, err := ParseUnit(s, metricBytesUnitMap)
	return MetricBytes(n), err
}

func (m MetricBytes) String() string {
	return ToString(int64(m), 1000, "B", "B")
}



func ParseStrictBytes(s string) (int64, error) {
	n, err := ParseUnit(s, bytesUnitMap)
	if err != nil {
		n, err = ParseUnit(s, metricBytesUnitMap)
	}
	return int64(n), err
}
