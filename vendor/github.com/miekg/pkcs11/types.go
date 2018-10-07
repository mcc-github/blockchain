



package pkcs11


import "C"

import (
	"fmt"
	"time"
	"unsafe"
)

type arena []unsafe.Pointer

func (a *arena) Allocate(obj []byte) (C.CK_VOID_PTR, C.CK_ULONG) {
	cobj := C.calloc(C.size_t(len(obj)), 1)
	*a = append(*a, cobj)
	C.memmove(cobj, unsafe.Pointer(&obj[0]), C.size_t(len(obj)))
	return C.CK_VOID_PTR(cobj), C.CK_ULONG(len(obj))
}

func (a arena) Free() {
	for _, p := range a {
		C.free(p)
	}
}


func toList(clist C.CK_ULONG_PTR, size C.CK_ULONG) []uint {
	l := make([]uint, int(size))
	for i := 0; i < len(l); i++ {
		l[i] = uint(C.Index(clist, C.CK_ULONG(i)))
	}
	defer C.free(unsafe.Pointer(clist))
	return l
}


func cBBool(x bool) C.CK_BBOOL {
	if x {
		return C.CK_BBOOL(C.CK_TRUE)
	}
	return C.CK_BBOOL(C.CK_FALSE)
}

func uintToBytes(x uint64) []byte {
	ul := C.CK_ULONG(x)
	return C.GoBytes(unsafe.Pointer(&ul), C.int(unsafe.Sizeof(ul)))
}


type Error uint

func (e Error) Error() string {
	return fmt.Sprintf("pkcs11: 0x%X: %s", uint(e), strerror[uint(e)])
}

func toError(e C.CK_RV) error {
	if e == C.CKR_OK {
		return nil
	}
	return Error(e)
}


type SessionHandle uint


type ObjectHandle uint


type Version struct {
	Major byte
	Minor byte
}

func toVersion(version C.CK_VERSION) Version {
	return Version{byte(version.major), byte(version.minor)}
}



type SlotEvent struct {
	SlotID uint
}


type Info struct {
	CryptokiVersion    Version
	ManufacturerID     string
	Flags              uint
	LibraryDescription string
	LibraryVersion     Version
}


type SlotInfo struct {
	SlotDescription string 
	ManufacturerID  string 
	Flags           uint
	HardwareVersion Version
	FirmwareVersion Version
}


type TokenInfo struct {
	Label              string
	ManufacturerID     string
	Model              string
	SerialNumber       string
	Flags              uint
	MaxSessionCount    uint
	SessionCount       uint
	MaxRwSessionCount  uint
	RwSessionCount     uint
	MaxPinLen          uint
	MinPinLen          uint
	TotalPublicMemory  uint
	FreePublicMemory   uint
	TotalPrivateMemory uint
	FreePrivateMemory  uint
	HardwareVersion    Version
	FirmwareVersion    Version
	UTCTime            string
}


type SessionInfo struct {
	SlotID      uint
	State       uint
	Flags       uint
	DeviceError uint
}


type Attribute struct {
	Type  uint
	Value []byte
}





func NewAttribute(typ uint, x interface{}) *Attribute {
	
	
	
	a := new(Attribute)
	a.Type = typ
	if x == nil {
		return a
	}
	switch v := x.(type) {
	case bool:
		if v {
			a.Value = []byte{1}
		} else {
			a.Value = []byte{0}
		}
	case int:
		a.Value = uintToBytes(uint64(v))
	case uint:
		a.Value = uintToBytes(uint64(v))
	case string:
		a.Value = []byte(v)
	case []byte:
		a.Value = v
	case time.Time: 
		a.Value = cDate(v)
	default:
		panic("pkcs11: unhandled attribute type")
	}
	return a
}


func cAttributeList(a []*Attribute) (arena, C.CK_ATTRIBUTE_PTR, C.CK_ULONG) {
	var arena arena
	if len(a) == 0 {
		return nil, nil, 0
	}
	pa := make([]C.CK_ATTRIBUTE, len(a))
	for i, attr := range a {
		pa[i]._type = C.CK_ATTRIBUTE_TYPE(attr.Type)
		if len(attr.Value) != 0 {
			buf, len := arena.Allocate(attr.Value)
			
			C.putAttributePval(&pa[i], buf)
			pa[i].ulValueLen = len
		}
	}
	return arena, &pa[0], C.CK_ULONG(len(a))
}

func cDate(t time.Time) []byte {
	b := make([]byte, 8)
	year, month, day := t.Date()
	y := fmt.Sprintf("%4d", year)
	m := fmt.Sprintf("%02d", month)
	d1 := fmt.Sprintf("%02d", day)
	b[0], b[1], b[2], b[3] = y[0], y[1], y[2], y[3]
	b[4], b[5] = m[0], m[1]
	b[6], b[7] = d1[0], d1[1]
	return b
}


type Mechanism struct {
	Mechanism uint
	Parameter []byte
	generator interface{}
}


func NewMechanism(mech uint, x interface{}) *Mechanism {
	m := new(Mechanism)
	m.Mechanism = mech
	if x == nil {
		return m
	}

	switch p := x.(type) {
	case *GCMParams, *OAEPParams:
		
		m.generator = p
	case []byte:
		m.Parameter = p
	default:
		panic("parameter must be one of type: []byte, *GCMParams, *OAEPParams")
	}

	return m
}

func cMechanism(mechList []*Mechanism) (arena, *C.CK_MECHANISM) {
	if len(mechList) != 1 {
		panic("expected exactly one mechanism")
	}
	mech := mechList[0]
	cmech := &C.CK_MECHANISM{mechanism: C.CK_MECHANISM_TYPE(mech.Mechanism)}
	
	param := mech.Parameter
	var arena arena
	switch p := mech.generator.(type) {
	case *GCMParams:
		
		param = cGCMParams(p)
	case *OAEPParams:
		param, arena = cOAEPParams(p, arena)
	}
	if len(param) != 0 {
		buf, len := arena.Allocate(param)
		
		C.putMechanismParam(cmech, buf)
		cmech.ulParameterLen = len
	}
	return arena, cmech
}


type MechanismInfo struct {
	MinKeySize uint
	MaxKeySize uint
	Flags      uint
}


var stubData = []byte{0}


func cMessage(data []byte) (dataPtr C.CK_BYTE_PTR) {
	l := len(data)
	if l == 0 {
		
		data = stubData
	}
	return C.CK_BYTE_PTR(unsafe.Pointer(&data[0]))
}
