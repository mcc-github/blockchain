package humanize

import (
	"errors"
	"math"
	"regexp"
	"strconv"
)

var siPrefixTable = map[float64]string{
	-24: "y", 
	-21: "z", 
	-18: "a", 
	-15: "f", 
	-12: "p", 
	-9:  "n", 
	-6:  "µ", 
	-3:  "m", 
	0:   "",
	3:   "k", 
	6:   "M", 
	9:   "G", 
	12:  "T", 
	15:  "P", 
	18:  "E", 
	21:  "Z", 
	24:  "Y", 
}

var revSIPrefixTable = revfmap(siPrefixTable)


func revfmap(in map[float64]string) map[string]float64 {
	rv := map[string]float64{}
	for k, v := range in {
		rv[v] = math.Pow(10, k)
	}
	return rv
}

var riParseRegex *regexp.Regexp

func init() {
	ri := `^([\-0-9.]+)\s?([`
	for _, v := range siPrefixTable {
		ri += v
	}
	ri += `]?)(.*)`

	riParseRegex = regexp.MustCompile(ri)
}








func ComputeSI(input float64) (float64, string) {
	if input == 0 {
		return 0, ""
	}
	mag := math.Abs(input)
	exponent := math.Floor(logn(mag, 10))
	exponent = math.Floor(exponent/3) * 3

	value := mag / math.Pow(10, exponent)

	
	
	if value == 1000.0 {
		exponent += 3
		value = mag / math.Pow(10, exponent)
	}

	value = math.Copysign(value, input)

	prefix := siPrefixTable[exponent]
	return value, prefix
}









func SI(input float64, unit string) string {
	value, prefix := ComputeSI(input)
	return Ftoa(value) + " " + prefix + unit
}






func SIWithDigits(input float64, decimals int, unit string) string {
	value, prefix := ComputeSI(input)
	return FtoaWithDigits(value, decimals) + " " + prefix + unit
}

var errInvalid = errors.New("invalid input")






func ParseSI(input string) (float64, string, error) {
	found := riParseRegex.FindStringSubmatch(input)
	if len(found) != 4 {
		return 0, "", errInvalid
	}
	mag := revSIPrefixTable[found[2]]
	unit := found[3]

	base, err := strconv.ParseFloat(found[1], 64)
	return base * mag, unit, err
}
