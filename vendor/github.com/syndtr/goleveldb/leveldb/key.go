





package leveldb

import (
	"encoding/binary"
	"fmt"

	"github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/storage"
)


type ErrInternalKeyCorrupted struct {
	Ikey   []byte
	Reason string
}

func (e *ErrInternalKeyCorrupted) Error() string {
	return fmt.Sprintf("leveldb: internal key %q corrupted: %s", e.Ikey, e.Reason)
}

func newErrInternalKeyCorrupted(ikey []byte, reason string) error {
	return errors.NewErrCorrupted(storage.FileDesc{}, &ErrInternalKeyCorrupted{append([]byte{}, ikey...), reason})
}

type keyType uint

func (kt keyType) String() string {
	switch kt {
	case keyTypeDel:
		return "d"
	case keyTypeVal:
		return "v"
	}
	return fmt.Sprintf("<invalid:%#x>", uint(kt))
}



const (
	keyTypeDel = keyType(0)
	keyTypeVal = keyType(1)
)






const keyTypeSeek = keyTypeVal

const (
	
	
	
	keyMaxSeq = (uint64(1) << 56) - 1
	
	keyMaxNum = (keyMaxSeq << 8) | uint64(keyTypeSeek)
)


var keyMaxNumBytes = make([]byte, 8)

func init() {
	binary.LittleEndian.PutUint64(keyMaxNumBytes, keyMaxNum)
}

type internalKey []byte

func makeInternalKey(dst, ukey []byte, seq uint64, kt keyType) internalKey {
	if seq > keyMaxSeq {
		panic("leveldb: invalid sequence number")
	} else if kt > keyTypeVal {
		panic("leveldb: invalid type")
	}

	dst = ensureBuffer(dst, len(ukey)+8)
	copy(dst, ukey)
	binary.LittleEndian.PutUint64(dst[len(ukey):], (seq<<8)|uint64(kt))
	return internalKey(dst)
}

func parseInternalKey(ik []byte) (ukey []byte, seq uint64, kt keyType, err error) {
	if len(ik) < 8 {
		return nil, 0, 0, newErrInternalKeyCorrupted(ik, "invalid length")
	}
	num := binary.LittleEndian.Uint64(ik[len(ik)-8:])
	seq, kt = uint64(num>>8), keyType(num&0xff)
	if kt > keyTypeVal {
		return nil, 0, 0, newErrInternalKeyCorrupted(ik, "invalid type")
	}
	ukey = ik[:len(ik)-8]
	return
}

func validInternalKey(ik []byte) bool {
	_, _, _, err := parseInternalKey(ik)
	return err == nil
}

func (ik internalKey) assert() {
	if ik == nil {
		panic("leveldb: nil internalKey")
	}
	if len(ik) < 8 {
		panic(fmt.Sprintf("leveldb: internal key %q, len=%d: invalid length", []byte(ik), len(ik)))
	}
}

func (ik internalKey) ukey() []byte {
	ik.assert()
	return ik[:len(ik)-8]
}

func (ik internalKey) num() uint64 {
	ik.assert()
	return binary.LittleEndian.Uint64(ik[len(ik)-8:])
}

func (ik internalKey) parseNum() (seq uint64, kt keyType) {
	num := ik.num()
	seq, kt = uint64(num>>8), keyType(num&0xff)
	if kt > keyTypeVal {
		panic(fmt.Sprintf("leveldb: internal key %q, len=%d: invalid type %#x", []byte(ik), len(ik), kt))
	}
	return
}

func (ik internalKey) String() string {
	if ik == nil {
		return "<nil>"
	}

	if ukey, seq, kt, err := parseInternalKey(ik); err == nil {
		return fmt.Sprintf("%s,%s%d", shorten(string(ukey)), kt, seq)
	}
	return fmt.Sprintf("<invalid:%#x>", []byte(ik))
}
