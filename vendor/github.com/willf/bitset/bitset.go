
package bitset

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
)


const wordSize = uint(64)


const log2WordSize = uint(6)


const allBits uint64 = 0xffffffffffffffff


var binaryOrder binary.ByteOrder = binary.BigEndian


var base64Encoding *base64.Encoding = base64.URLEncoding


func Base64StdEncoding() { base64Encoding = base64.StdEncoding }


func LittleEndian() { binaryOrder = binary.LittleEndian }


type BitSet struct {
	length uint
	set    []uint64
}


type Error string


func (b *BitSet) safeSet() []uint64 {
	if b.set == nil {
		b.set = make([]uint64, wordsNeeded(0))
	}
	return b.set
}


func From(buf []uint64) *BitSet {
	return &BitSet{uint(len(buf)) * 64, buf}
}


func (b *BitSet) Bytes() []uint64 {
	return b.set
}


func wordsNeeded(i uint) int {
	if i > (Cap() - wordSize + 1) {
		return int(Cap() >> log2WordSize)
	}
	return int((i + (wordSize - 1)) >> log2WordSize)
}


func New(length uint) (bset *BitSet) {
	defer func() {
		if r := recover(); r != nil {
			bset = &BitSet{
				0,
				make([]uint64, 0),
			}
		}
	}()

	bset = &BitSet{
		length,
		make([]uint64, wordsNeeded(length)),
	}

	return bset
}


func Cap() uint {
	return ^uint(0)
}


func (b *BitSet) Len() uint {
	return b.length
}


func (b *BitSet) extendSetMaybe(i uint) {
	if i >= b.length { 
		nsize := wordsNeeded(i + 1)
		if b.set == nil {
			b.set = make([]uint64, nsize)
		} else if cap(b.set) >= nsize {
			b.set = b.set[:nsize] 
		} else if len(b.set) < nsize {
			newset := make([]uint64, nsize, 2*nsize) 
			copy(newset, b.set)
			b.set = newset
		}
		b.length = i + 1
	}
}


func (b *BitSet) Test(i uint) bool {
	if i >= b.length {
		return false
	}
	return b.set[i>>log2WordSize]&(1<<(i&(wordSize-1))) != 0
}


func (b *BitSet) Set(i uint) *BitSet {
	b.extendSetMaybe(i)
	b.set[i>>log2WordSize] |= 1 << (i & (wordSize - 1))
	return b
}


func (b *BitSet) Clear(i uint) *BitSet {
	if i >= b.length {
		return b
	}
	b.set[i>>log2WordSize] &^= 1 << (i & (wordSize - 1))
	return b
}


func (b *BitSet) SetTo(i uint, value bool) *BitSet {
	if value {
		return b.Set(i)
	}
	return b.Clear(i)
}


func (b *BitSet) Flip(i uint) *BitSet {
	if i >= b.length {
		return b.Set(i)
	}
	b.set[i>>log2WordSize] ^= 1 << (i & (wordSize - 1))
	return b
}


func (b *BitSet) String() string {
	
	var buffer bytes.Buffer
	start := []byte("{")
	buffer.Write(start)
	counter := 0
	i, e := b.NextSet(0)
	for e {
		counter = counter + 1
		
		if counter > 0x40000 {
			buffer.WriteString("...")
			break
		}
		buffer.WriteString(strconv.FormatInt(int64(i), 10))
		i, e = b.NextSet(i + 1)
		if e {
			buffer.WriteString(",")
		}
	}
	buffer.WriteString("}")
	return buffer.String()
}





func (b *BitSet) NextSet(i uint) (uint, bool) {
	x := int(i >> log2WordSize)
	if x >= len(b.set) {
		return 0, false
	}
	w := b.set[x]
	w = w >> (i & (wordSize - 1))
	if w != 0 {
		return i + trailingZeroes64(w), true
	}
	x = x + 1
	for x < len(b.set) {
		if b.set[x] != 0 {
			return uint(x)*wordSize + trailingZeroes64(b.set[x]), true
		}
		x = x + 1

	}
	return 0, false
}















func (b *BitSet) NextSetMany(i uint, buffer []uint) (uint, []uint) {
	myanswer := buffer
	capacity := cap(buffer)
	x := int(i >> log2WordSize)
	if x >= len(b.set) || capacity == 0 {
		return 0, myanswer[:0]
	}
	skip := i & (wordSize - 1)
	word := b.set[x] >> skip
	myanswer = myanswer[:capacity]
	size := int(0)
	for word != 0 {
		r := trailingZeroes64(word)
		t := word & ((^word) + 1)
		myanswer[size] = r + i
		size++
		if size == capacity {
			goto End
		}
		word = word ^ t
	}
	x++
	for idx, word := range b.set[x:] {
		for word != 0 {
			r := trailingZeroes64(word)
			t := word & ((^word) + 1)
			myanswer[size] = r + (uint(x+idx) << 6)
			size++
			if size == capacity {
				goto End
			}
			word = word ^ t
		}
	}
End:
	if size > 0 {
		return myanswer[size-1], myanswer[:size]
	} else {
		return 0, myanswer[:0]
	}
}




func (b *BitSet) NextClear(i uint) (uint, bool) {
	x := int(i >> log2WordSize)
	if x >= len(b.set) {
		return 0, false
	}
	w := b.set[x]
	w = w >> (i & (wordSize - 1))
	wA := allBits >> (i & (wordSize - 1))
	index := i + trailingZeroes64(^w)
	if w != wA && index < b.length {
		return index, true
	}
	x++
	for x < len(b.set) {
		index = uint(x)*wordSize + trailingZeroes64(^b.set[x])
		if b.set[x] != allBits && index < b.length {
			return index, true
		}
		x++
	}
	return 0, false
}


func (b *BitSet) ClearAll() *BitSet {
	if b != nil && b.set != nil {
		for i := range b.set {
			b.set[i] = 0
		}
	}
	return b
}


func (b *BitSet) wordCount() int {
	return len(b.set)
}


func (b *BitSet) Clone() *BitSet {
	c := New(b.length)
	if b.set != nil { 
		copy(c.set, b.set)
	}
	return c
}




func (b *BitSet) Copy(c *BitSet) (count uint) {
	if c == nil {
		return
	}
	if b.set != nil { 
		copy(c.set, b.set)
	}
	count = c.length
	if b.length < c.length {
		count = b.length
	}
	return
}


func (b *BitSet) Count() uint {
	if b != nil && b.set != nil {
		return uint(popcntSlice(b.set))
	}
	return 0
}




func (b *BitSet) Equal(c *BitSet) bool {
	if c == nil {
		return false
	}
	if b.length != c.length {
		return false
	}
	if b.length == 0 { 
		return true
	}
	

	for p, v := range b.set {
		if c.set[p] != v {
			return false
		}
	}
	return true
}

func panicIfNull(b *BitSet) {
	if b == nil {
		panic(Error("BitSet must not be null"))
	}
}



func (b *BitSet) Difference(compare *BitSet) (result *BitSet) {
	panicIfNull(b)
	panicIfNull(compare)
	result = b.Clone() 
	l := int(compare.wordCount())
	if l > int(b.wordCount()) {
		l = int(b.wordCount())
	}
	for i := 0; i < l; i++ {
		result.set[i] = b.set[i] &^ compare.set[i]
	}
	return
}


func (b *BitSet) DifferenceCardinality(compare *BitSet) uint {
	panicIfNull(b)
	panicIfNull(compare)
	l := int(compare.wordCount())
	if l > int(b.wordCount()) {
		l = int(b.wordCount())
	}
	cnt := uint64(0)
	cnt += popcntMaskSlice(b.set[:l], compare.set[:l])
	cnt += popcntSlice(b.set[l:])
	return uint(cnt)
}



func (b *BitSet) InPlaceDifference(compare *BitSet) {
	panicIfNull(b)
	panicIfNull(compare)
	l := int(compare.wordCount())
	if l > int(b.wordCount()) {
		l = int(b.wordCount())
	}
	for i := 0; i < l; i++ {
		b.set[i] &^= compare.set[i]
	}
}



func sortByLength(a *BitSet, b *BitSet) (ap *BitSet, bp *BitSet) {
	if a.length <= b.length {
		ap, bp = a, b
	} else {
		ap, bp = b, a
	}
	return
}



func (b *BitSet) Intersection(compare *BitSet) (result *BitSet) {
	panicIfNull(b)
	panicIfNull(compare)
	b, compare = sortByLength(b, compare)
	result = New(b.length)
	for i, word := range b.set {
		result.set[i] = word & compare.set[i]
	}
	return
}


func (b *BitSet) IntersectionCardinality(compare *BitSet) uint {
	panicIfNull(b)
	panicIfNull(compare)
	b, compare = sortByLength(b, compare)
	cnt := popcntAndSlice(b.set, compare.set)
	return uint(cnt)
}




func (b *BitSet) InPlaceIntersection(compare *BitSet) {
	panicIfNull(b)
	panicIfNull(compare)
	l := int(compare.wordCount())
	if l > int(b.wordCount()) {
		l = int(b.wordCount())
	}
	for i := 0; i < l; i++ {
		b.set[i] &= compare.set[i]
	}
	for i := l; i < len(b.set); i++ {
		b.set[i] = 0
	}
	if compare.length > 0 {
		b.extendSetMaybe(compare.length - 1)
	}
}



func (b *BitSet) Union(compare *BitSet) (result *BitSet) {
	panicIfNull(b)
	panicIfNull(compare)
	b, compare = sortByLength(b, compare)
	result = compare.Clone()
	for i, word := range b.set {
		result.set[i] = word | compare.set[i]
	}
	return
}



func (b *BitSet) UnionCardinality(compare *BitSet) uint {
	panicIfNull(b)
	panicIfNull(compare)
	b, compare = sortByLength(b, compare)
	cnt := popcntOrSlice(b.set, compare.set)
	if len(compare.set) > len(b.set) {
		cnt += popcntSlice(compare.set[len(b.set):])
	}
	return uint(cnt)
}



func (b *BitSet) InPlaceUnion(compare *BitSet) {
	panicIfNull(b)
	panicIfNull(compare)
	l := int(compare.wordCount())
	if l > int(b.wordCount()) {
		l = int(b.wordCount())
	}
	if compare.length > 0 {
		b.extendSetMaybe(compare.length - 1)
	}
	for i := 0; i < l; i++ {
		b.set[i] |= compare.set[i]
	}
	if len(compare.set) > l {
		for i := l; i < len(compare.set); i++ {
			b.set[i] = compare.set[i]
		}
	}
}



func (b *BitSet) SymmetricDifference(compare *BitSet) (result *BitSet) {
	panicIfNull(b)
	panicIfNull(compare)
	b, compare = sortByLength(b, compare)
	
	result = compare.Clone()
	for i, word := range b.set {
		result.set[i] = word ^ compare.set[i]
	}
	return
}


func (b *BitSet) SymmetricDifferenceCardinality(compare *BitSet) uint {
	panicIfNull(b)
	panicIfNull(compare)
	b, compare = sortByLength(b, compare)
	cnt := popcntXorSlice(b.set, compare.set)
	if len(compare.set) > len(b.set) {
		cnt += popcntSlice(compare.set[len(b.set):])
	}
	return uint(cnt)
}



func (b *BitSet) InPlaceSymmetricDifference(compare *BitSet) {
	panicIfNull(b)
	panicIfNull(compare)
	l := int(compare.wordCount())
	if l > int(b.wordCount()) {
		l = int(b.wordCount())
	}
	if compare.length > 0 {
		b.extendSetMaybe(compare.length - 1)
	}
	for i := 0; i < l; i++ {
		b.set[i] ^= compare.set[i]
	}
	if len(compare.set) > l {
		for i := l; i < len(compare.set); i++ {
			b.set[i] = compare.set[i]
		}
	}
}


func (b *BitSet) isLenExactMultiple() bool {
	return b.length%wordSize == 0
}


func (b *BitSet) cleanLastWord() {
	if !b.isLenExactMultiple() {
		b.set[len(b.set)-1] &= allBits >> (wordSize - b.length%wordSize)
	}
}


func (b *BitSet) Complement() (result *BitSet) {
	panicIfNull(b)
	result = New(b.length)
	for i, word := range b.set {
		result.set[i] = ^word
	}
	result.cleanLastWord()
	return
}



func (b *BitSet) All() bool {
	panicIfNull(b)
	return b.Count() == b.length
}



func (b *BitSet) None() bool {
	panicIfNull(b)
	if b != nil && b.set != nil {
		for _, word := range b.set {
			if word > 0 {
				return false
			}
		}
		return true
	}
	return true
}


func (b *BitSet) Any() bool {
	panicIfNull(b)
	return !b.None()
}


func (b *BitSet) IsSuperSet(other *BitSet) bool {
	for i, e := other.NextSet(0); e; i, e = other.NextSet(i + 1) {
		if !b.Test(i) {
			return false
		}
	}
	return true
}


func (b *BitSet) IsStrictSuperSet(other *BitSet) bool {
	return b.Count() > other.Count() && b.IsSuperSet(other)
}


func (b *BitSet) DumpAsBits() string {
	if b.set == nil {
		return "."
	}
	buffer := bytes.NewBufferString("")
	i := len(b.set) - 1
	for ; i >= 0; i-- {
		fmt.Fprintf(buffer, "%064b.", b.set[i])
	}
	return string(buffer.Bytes())
}


func (b *BitSet) BinaryStorageSize() int {
	return binary.Size(uint64(0)) + binary.Size(b.set)
}


func (b *BitSet) WriteTo(stream io.Writer) (int64, error) {
	length := uint64(b.length)

	
	err := binary.Write(stream, binaryOrder, length)
	if err != nil {
		return 0, err
	}

	
	err = binary.Write(stream, binaryOrder, b.set)
	return int64(b.BinaryStorageSize()), err
}


func (b *BitSet) ReadFrom(stream io.Reader) (int64, error) {
	var length uint64

	
	err := binary.Read(stream, binaryOrder, &length)
	if err != nil {
		return 0, err
	}
	newset := New(uint(length))

	if uint64(newset.length) != length {
		return 0, errors.New("Unmarshalling error: type mismatch")
	}

	
	err = binary.Read(stream, binaryOrder, newset.set)
	if err != nil {
		return 0, err
	}

	*b = *newset
	return int64(b.BinaryStorageSize()), nil
}


func (b *BitSet) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)

	_, err := b.WriteTo(writer)
	if err != nil {
		return []byte{}, err
	}

	err = writer.Flush()

	return buf.Bytes(), err
}


func (b *BitSet) UnmarshalBinary(data []byte) error {
	buf := bytes.NewReader(data)
	reader := bufio.NewReader(buf)

	_, err := b.ReadFrom(reader)

	return err
}


func (b *BitSet) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBuffer(make([]byte, 0, b.BinaryStorageSize()))
	_, err := b.WriteTo(buffer)
	if err != nil {
		return nil, err
	}

	
	return json.Marshal(base64Encoding.EncodeToString(buffer.Bytes()))
}


func (b *BitSet) UnmarshalJSON(data []byte) error {
	
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	
	buf, err := base64Encoding.DecodeString(s)
	if err != nil {
		return err
	}

	_, err = b.ReadFrom(bytes.NewReader(buf))
	return err
}
