package humanize



import (
	"math"
	"strconv"
)

var (
	renderFloatPrecisionMultipliers = [...]float64{
		1,
		10,
		100,
		1000,
		10000,
		100000,
		1000000,
		10000000,
		100000000,
		1000000000,
	}

	renderFloatPrecisionRounders = [...]float64{
		0.5,
		0.05,
		0.005,
		0.0005,
		0.00005,
		0.000005,
		0.0000005,
		0.00000005,
		0.000000005,
		0.0000000005,
	}
)






















func FormatFloat(format string, n float64) string {
	
	
	
	
	if math.IsNaN(n) {
		return "NaN"
	}
	if n > math.MaxFloat64 {
		return "Infinity"
	}
	if n < -math.MaxFloat64 {
		return "-Infinity"
	}

	
	precision := 2
	decimalStr := "."
	thousandStr := ","
	positiveStr := ""
	negativeStr := "-"

	if len(format) > 0 {
		format := []rune(format)

		
		
		precision = 9
		thousandStr = ""

		
		formatIndx := []int{}
		for i, char := range format {
			if char != '#' && char != '0' {
				formatIndx = append(formatIndx, i)
			}
		}

		if len(formatIndx) > 0 {
			
			
			
			
			
			
			
			
			if formatIndx[0] == 0 {
				if format[formatIndx[0]] != '+' {
					panic("RenderFloat(): invalid positive sign directive")
				}
				positiveStr = "+"
				formatIndx = formatIndx[1:]
			}

			
			
			
			
			
			
			if len(formatIndx) == 2 {
				if (formatIndx[1] - formatIndx[0]) != 4 {
					panic("RenderFloat(): thousands separator directive must be followed by 3 digit-specifiers")
				}
				thousandStr = string(format[formatIndx[0]])
				formatIndx = formatIndx[1:]
			}

			
			
			
			
			
			
			if len(formatIndx) == 1 {
				decimalStr = string(format[formatIndx[0]])
				precision = len(format) - formatIndx[0] - 1
			}
		}
	}

	
	var signStr string
	if n >= 0.000000001 {
		signStr = positiveStr
	} else if n <= -0.000000001 {
		signStr = negativeStr
		n = -n
	} else {
		signStr = ""
		n = 0.0
	}

	
	intf, fracf := math.Modf(n + renderFloatPrecisionRounders[precision])

	
	intStr := strconv.FormatInt(int64(intf), 10)

	
	if len(thousandStr) > 0 {
		for i := len(intStr); i > 3; {
			i -= 3
			intStr = intStr[:i] + thousandStr + intStr[i:]
		}
	}

	
	if precision == 0 {
		return signStr + intStr
	}

	
	fracStr := strconv.Itoa(int(fracf * renderFloatPrecisionMultipliers[precision]))
	
	if len(fracStr) < precision {
		fracStr = "000000000000000"[:precision-len(fracStr)] + fracStr
	}

	return signStr + intStr + decimalStr + fracStr
}



func FormatInteger(format string, n int) string {
	return FormatFloat(format, float64(n))
}
