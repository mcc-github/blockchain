package humanize

import (
	"fmt"
	"math/big"
	"strings"
	"unicode"
)

var (
	bigIECExp = big.NewInt(1024)

	
	BigByte = big.NewInt(1)
	
	BigKiByte = (&big.Int{}).Mul(BigByte, bigIECExp)
	
	BigMiByte = (&big.Int{}).Mul(BigKiByte, bigIECExp)
	
	BigGiByte = (&big.Int{}).Mul(BigMiByte, bigIECExp)
	
	BigTiByte = (&big.Int{}).Mul(BigGiByte, bigIECExp)
	
	BigPiByte = (&big.Int{}).Mul(BigTiByte, bigIECExp)
	
	BigEiByte = (&big.Int{}).Mul(BigPiByte, bigIECExp)
	
	BigZiByte = (&big.Int{}).Mul(BigEiByte, bigIECExp)
	
	BigYiByte = (&big.Int{}).Mul(BigZiByte, bigIECExp)
)

var (
	bigSIExp = big.NewInt(1000)

	
	BigSIByte = big.NewInt(1)
	
	BigKByte = (&big.Int{}).Mul(BigSIByte, bigSIExp)
	
	BigMByte = (&big.Int{}).Mul(BigKByte, bigSIExp)
	
	BigGByte = (&big.Int{}).Mul(BigMByte, bigSIExp)
	
	BigTByte = (&big.Int{}).Mul(BigGByte, bigSIExp)
	
	BigPByte = (&big.Int{}).Mul(BigTByte, bigSIExp)
	
	BigEByte = (&big.Int{}).Mul(BigPByte, bigSIExp)
	
	BigZByte = (&big.Int{}).Mul(BigEByte, bigSIExp)
	
	BigYByte = (&big.Int{}).Mul(BigZByte, bigSIExp)
)

var bigBytesSizeTable = map[string]*big.Int{
	"b":   BigByte,
	"kib": BigKiByte,
	"kb":  BigKByte,
	"mib": BigMiByte,
	"mb":  BigMByte,
	"gib": BigGiByte,
	"gb":  BigGByte,
	"tib": BigTiByte,
	"tb":  BigTByte,
	"pib": BigPiByte,
	"pb":  BigPByte,
	"eib": BigEiByte,
	"eb":  BigEByte,
	"zib": BigZiByte,
	"zb":  BigZByte,
	"yib": BigYiByte,
	"yb":  BigYByte,
	
	"":   BigByte,
	"ki": BigKiByte,
	"k":  BigKByte,
	"mi": BigMiByte,
	"m":  BigMByte,
	"gi": BigGiByte,
	"g":  BigGByte,
	"ti": BigTiByte,
	"t":  BigTByte,
	"pi": BigPiByte,
	"p":  BigPByte,
	"ei": BigEiByte,
	"e":  BigEByte,
	"z":  BigZByte,
	"zi": BigZiByte,
	"y":  BigYByte,
	"yi": BigYiByte,
}

var ten = big.NewInt(10)

func humanateBigBytes(s, base *big.Int, sizes []string) string {
	if s.Cmp(ten) < 0 {
		return fmt.Sprintf("%d B", s)
	}
	c := (&big.Int{}).Set(s)
	val, mag := oomm(c, base, len(sizes)-1)
	suffix := sizes[mag]
	f := "%.0f %s"
	if val < 10 {
		f = "%.1f %s"
	}

	return fmt.Sprintf(f, val, suffix)

}






func BigBytes(s *big.Int) string {
	sizes := []string{"B", "kB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}
	return humanateBigBytes(s, bigSIExp, sizes)
}






func BigIBytes(s *big.Int) string {
	sizes := []string{"B", "KiB", "MiB", "GiB", "TiB", "PiB", "EiB", "ZiB", "YiB"}
	return humanateBigBytes(s, bigIECExp, sizes)
}








func ParseBigBytes(s string) (*big.Int, error) {
	lastDigit := 0
	hasComma := false
	for _, r := range s {
		if !(unicode.IsDigit(r) || r == '.' || r == ',') {
			break
		}
		if r == ',' {
			hasComma = true
		}
		lastDigit++
	}

	num := s[:lastDigit]
	if hasComma {
		num = strings.Replace(num, ",", "", -1)
	}

	val := &big.Rat{}
	_, err := fmt.Sscanf(num, "%f", val)
	if err != nil {
		return nil, err
	}

	extra := strings.ToLower(strings.TrimSpace(s[lastDigit:]))
	if m, ok := bigBytesSizeTable[extra]; ok {
		mv := (&big.Rat{}).SetInt(m)
		val.Mul(val, mv)
		rv := &big.Int{}
		rv.Div(val.Num(), val.Denom())
		return rv, nil
	}

	return nil, fmt.Errorf("unhandled size name: %v", extra)
}
