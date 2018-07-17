





package main




import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
)

func main() {
	fmt.Printf("
	fmt.Printf("
	fmt.Printf(`package korean 

	res, err := http.Get("http://encoding.spec.whatwg.org/index-euc-kr.txt")
	if err != nil {
		log.Fatalf("Get: %v", err)
	}
	defer res.Body.Close()

	mapping := [65536]uint16{}
	reverse := [65536]uint16{}

	scanner := bufio.NewScanner(res.Body)
	for scanner.Scan() {
		s := strings.TrimSpace(scanner.Text())
		if s == "" || s[0] == '#' {
			continue
		}
		x, y := uint16(0), uint16(0)
		if _, err := fmt.Sscanf(s, "%d 0x%x", &x, &y); err != nil {
			log.Fatalf("could not parse %q", s)
		}
		if x < 0 || 178*(0xc7-0x81)+(0xfe-0xc7)*94+(0xff-0xa1) <= x {
			log.Fatalf("EUC-KR code %d is out of range", x)
		}
		mapping[x] = y
		if reverse[y] == 0 {
			c0, c1 := uint16(0), uint16(0)
			if x < 178*(0xc7-0x81) {
				c0 = uint16(x/178) + 0x81
				c1 = uint16(x % 178)
				switch {
				case c1 < 1*26:
					c1 += 0x41
				case c1 < 2*26:
					c1 += 0x47
				default:
					c1 += 0x4d
				}
			} else {
				x -= 178 * (0xc7 - 0x81)
				c0 = uint16(x/94) + 0xc7
				c1 = uint16(x%94) + 0xa1
			}
			reverse[y] = c0<<8 | c1
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("scanner error: %v", err)
	}

	fmt.Printf("
	fmt.Printf("
	fmt.Printf("var decode = [...]uint16{\n")
	for i, v := range mapping {
		if v != 0 {
			fmt.Printf("\t%d: 0x%04X,\n", i, v)
		}
	}
	fmt.Printf("}\n\n")

	
	
	const separation = 1024

	intervals := []interval(nil)
	low, high := -1, -1
	for i, v := range reverse {
		if v == 0 {
			continue
		}
		if low < 0 {
			low = i
		} else if i-high >= separation {
			if high >= 0 {
				intervals = append(intervals, interval{low, high})
			}
			low = i
		}
		high = i + 1
	}
	if high >= 0 {
		intervals = append(intervals, interval{low, high})
	}
	sort.Sort(byDecreasingLength(intervals))

	fmt.Printf("const numEncodeTables = %d\n\n", len(intervals))
	fmt.Printf("
	fmt.Printf("
	for i, v := range intervals {
		fmt.Printf("
	}
	fmt.Printf("\n")

	for i, v := range intervals {
		fmt.Printf("const encode%dLow, encode%dHigh = %d, %d\n\n", i, i, v.low, v.high)
		fmt.Printf("var encode%d = [...]uint16{\n", i)
		for j := v.low; j < v.high; j++ {
			x := reverse[j]
			if x == 0 {
				continue
			}
			fmt.Printf("\t%d-%d: 0x%04X,\n", j, v.low, x)
		}
		fmt.Printf("}\n\n")
	}
}


type interval struct {
	low, high int
}

func (i interval) len() int { return i.high - i.low }


type byDecreasingLength []interval

func (b byDecreasingLength) Len() int           { return len(b) }
func (b byDecreasingLength) Less(i, j int) bool { return b[i].len() > b[j].len() }
func (b byDecreasingLength) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
