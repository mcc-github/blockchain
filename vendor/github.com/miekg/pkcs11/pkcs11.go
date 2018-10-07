




package pkcs11







import "C"
import "strings"

import "unsafe"


type Ctx struct {
	ctx *C.struct_ctx
}


func New(module string) *Ctx {
	
	
	if module == "" {
		return nil
	}
	c := new(Ctx)
	mod := C.CString(module)
	defer C.free(unsafe.Pointer(mod))
	c.ctx = C.New(mod)
	if c.ctx == nil {
		return nil
	}
	return c
}


func (c *Ctx) Destroy() {
	if c == nil || c.ctx == nil {
		return
	}
	C.Destroy(c.ctx)
	c.ctx = nil
}


func (c *Ctx) Initialize() error {
	e := C.Initialize(c.ctx)
	return toError(e)
}


func (c *Ctx) Finalize() error {
	if c.ctx == nil {
		return toError(CKR_CRYPTOKI_NOT_INITIALIZED)
	}
	e := C.Finalize(c.ctx)
	return toError(e)
}


func (c *Ctx) GetInfo() (Info, error) {
	var p C.ckInfo
	e := C.GetInfo(c.ctx, &p)
	i := Info{
		CryptokiVersion:    toVersion(p.cryptokiVersion),
		ManufacturerID:     strings.TrimRight(string(C.GoBytes(unsafe.Pointer(&p.manufacturerID[0]), 32)), " "),
		Flags:              uint(p.flags),
		LibraryDescription: strings.TrimRight(string(C.GoBytes(unsafe.Pointer(&p.libraryDescription[0]), 32)), " "),
		LibraryVersion:     toVersion(p.libraryVersion),
	}
	return i, toError(e)
}


func (c *Ctx) GetSlotList(tokenPresent bool) ([]uint, error) {
	var (
		slotList C.CK_ULONG_PTR
		ulCount  C.CK_ULONG
	)
	e := C.GetSlotList(c.ctx, cBBool(tokenPresent), &slotList, &ulCount)
	if toError(e) != nil {
		return nil, toError(e)
	}
	l := toList(slotList, ulCount)
	return l, nil
}


func (c *Ctx) GetSlotInfo(slotID uint) (SlotInfo, error) {
	var csi C.CK_SLOT_INFO
	e := C.GetSlotInfo(c.ctx, C.CK_ULONG(slotID), &csi)
	s := SlotInfo{
		SlotDescription: strings.TrimRight(string(C.GoBytes(unsafe.Pointer(&csi.slotDescription[0]), 64)), " "),
		ManufacturerID:  strings.TrimRight(string(C.GoBytes(unsafe.Pointer(&csi.manufacturerID[0]), 32)), " "),
		Flags:           uint(csi.flags),
		HardwareVersion: toVersion(csi.hardwareVersion),
		FirmwareVersion: toVersion(csi.firmwareVersion),
	}
	return s, toError(e)
}



func (c *Ctx) GetTokenInfo(slotID uint) (TokenInfo, error) {
	var cti C.CK_TOKEN_INFO
	e := C.GetTokenInfo(c.ctx, C.CK_ULONG(slotID), &cti)
	s := TokenInfo{
		Label:              strings.TrimRight(string(C.GoBytes(unsafe.Pointer(&cti.label[0]), 32)), " "),
		ManufacturerID:     strings.TrimRight(string(C.GoBytes(unsafe.Pointer(&cti.manufacturerID[0]), 32)), " "),
		Model:              strings.TrimRight(string(C.GoBytes(unsafe.Pointer(&cti.model[0]), 16)), " "),
		SerialNumber:       strings.TrimRight(string(C.GoBytes(unsafe.Pointer(&cti.serialNumber[0]), 16)), " "),
		Flags:              uint(cti.flags),
		MaxSessionCount:    uint(cti.ulMaxSessionCount),
		SessionCount:       uint(cti.ulSessionCount),
		MaxRwSessionCount:  uint(cti.ulMaxRwSessionCount),
		RwSessionCount:     uint(cti.ulRwSessionCount),
		MaxPinLen:          uint(cti.ulMaxPinLen),
		MinPinLen:          uint(cti.ulMinPinLen),
		TotalPublicMemory:  uint(cti.ulTotalPublicMemory),
		FreePublicMemory:   uint(cti.ulFreePublicMemory),
		TotalPrivateMemory: uint(cti.ulTotalPrivateMemory),
		FreePrivateMemory:  uint(cti.ulFreePrivateMemory),
		HardwareVersion:    toVersion(cti.hardwareVersion),
		FirmwareVersion:    toVersion(cti.firmwareVersion),
		UTCTime:            strings.TrimRight(string(C.GoBytes(unsafe.Pointer(&cti.utcTime[0]), 16)), " "),
	}
	return s, toError(e)
}


func (c *Ctx) GetMechanismList(slotID uint) ([]*Mechanism, error) {
	var (
		mech    C.CK_ULONG_PTR 
		mechlen C.CK_ULONG
	)
	e := C.GetMechanismList(c.ctx, C.CK_ULONG(slotID), &mech, &mechlen)
	if toError(e) != nil {
		return nil, toError(e)
	}
	
	
	m := make([]*Mechanism, int(mechlen))
	for i, typ := range toList(mech, mechlen) {
		m[i] = NewMechanism(typ, nil)
	}
	return m, nil
}



func (c *Ctx) GetMechanismInfo(slotID uint, m []*Mechanism) (MechanismInfo, error) {
	var cm C.CK_MECHANISM_INFO
	e := C.GetMechanismInfo(c.ctx, C.CK_ULONG(slotID), C.CK_MECHANISM_TYPE(m[0].Mechanism),
		C.CK_MECHANISM_INFO_PTR(&cm))
	mi := MechanismInfo{
		MinKeySize: uint(cm.ulMinKeySize),
		MaxKeySize: uint(cm.ulMaxKeySize),
		Flags:      uint(cm.flags),
	}
	return mi, toError(e)
}




func (c *Ctx) InitToken(slotID uint, pin string, label string) error {
	p := C.CString(pin)
	defer C.free(unsafe.Pointer(p))
	ll := len(label)
	for ll < 32 {
		label += " "
		ll++
	}
	l := C.CString(label[:32])
	defer C.free(unsafe.Pointer(l))
	e := C.InitToken(c.ctx, C.CK_ULONG(slotID), p, C.CK_ULONG(len(pin)), l)
	return toError(e)
}


func (c *Ctx) InitPIN(sh SessionHandle, pin string) error {
	p := C.CString(pin)
	defer C.free(unsafe.Pointer(p))
	e := C.InitPIN(c.ctx, C.CK_SESSION_HANDLE(sh), p, C.CK_ULONG(len(pin)))
	return toError(e)
}


func (c *Ctx) SetPIN(sh SessionHandle, oldpin string, newpin string) error {
	old := C.CString(oldpin)
	defer C.free(unsafe.Pointer(old))
	new := C.CString(newpin)
	defer C.free(unsafe.Pointer(new))
	e := C.SetPIN(c.ctx, C.CK_SESSION_HANDLE(sh), old, C.CK_ULONG(len(oldpin)), new, C.CK_ULONG(len(newpin)))
	return toError(e)
}


func (c *Ctx) OpenSession(slotID uint, flags uint) (SessionHandle, error) {
	var s C.CK_SESSION_HANDLE
	e := C.OpenSession(c.ctx, C.CK_ULONG(slotID), C.CK_ULONG(flags), C.CK_SESSION_HANDLE_PTR(&s))
	return SessionHandle(s), toError(e)
}


func (c *Ctx) CloseSession(sh SessionHandle) error {
	if c.ctx == nil {
		return toError(CKR_CRYPTOKI_NOT_INITIALIZED)
	}
	e := C.CloseSession(c.ctx, C.CK_SESSION_HANDLE(sh))
	return toError(e)
}


func (c *Ctx) CloseAllSessions(slotID uint) error {
	if c.ctx == nil {
		return toError(CKR_CRYPTOKI_NOT_INITIALIZED)
	}
	e := C.CloseAllSessions(c.ctx, C.CK_ULONG(slotID))
	return toError(e)
}


func (c *Ctx) GetSessionInfo(sh SessionHandle) (SessionInfo, error) {
	var csi C.CK_SESSION_INFO
	e := C.GetSessionInfo(c.ctx, C.CK_SESSION_HANDLE(sh), &csi)
	s := SessionInfo{SlotID: uint(csi.slotID),
		State:       uint(csi.state),
		Flags:       uint(csi.flags),
		DeviceError: uint(csi.ulDeviceError),
	}
	return s, toError(e)
}


func (c *Ctx) GetOperationState(sh SessionHandle) ([]byte, error) {
	var (
		state    C.CK_BYTE_PTR
		statelen C.CK_ULONG
	)
	e := C.GetOperationState(c.ctx, C.CK_SESSION_HANDLE(sh), &state, &statelen)
	if toError(e) != nil {
		return nil, toError(e)
	}
	b := C.GoBytes(unsafe.Pointer(state), C.int(statelen))
	C.free(unsafe.Pointer(state))
	return b, nil
}


func (c *Ctx) SetOperationState(sh SessionHandle, state []byte, encryptKey, authKey ObjectHandle) error {
	e := C.SetOperationState(c.ctx, C.CK_SESSION_HANDLE(sh), C.CK_BYTE_PTR(unsafe.Pointer(&state[0])),
		C.CK_ULONG(len(state)), C.CK_OBJECT_HANDLE(encryptKey), C.CK_OBJECT_HANDLE(authKey))
	return toError(e)
}


func (c *Ctx) Login(sh SessionHandle, userType uint, pin string) error {
	p := C.CString(pin)
	defer C.free(unsafe.Pointer(p))
	e := C.Login(c.ctx, C.CK_SESSION_HANDLE(sh), C.CK_USER_TYPE(userType), p, C.CK_ULONG(len(pin)))
	return toError(e)
}


func (c *Ctx) Logout(sh SessionHandle) error {
	if c.ctx == nil {
		return toError(CKR_CRYPTOKI_NOT_INITIALIZED)
	}
	e := C.Logout(c.ctx, C.CK_SESSION_HANDLE(sh))
	return toError(e)
}


func (c *Ctx) CreateObject(sh SessionHandle, temp []*Attribute) (ObjectHandle, error) {
	var obj C.CK_OBJECT_HANDLE
	arena, t, tcount := cAttributeList(temp)
	defer arena.Free()
	e := C.CreateObject(c.ctx, C.CK_SESSION_HANDLE(sh), t, tcount, C.CK_OBJECT_HANDLE_PTR(&obj))
	e1 := toError(e)
	if e1 == nil {
		return ObjectHandle(obj), nil
	}
	return 0, e1
}


func (c *Ctx) CopyObject(sh SessionHandle, o ObjectHandle, temp []*Attribute) (ObjectHandle, error) {
	var obj C.CK_OBJECT_HANDLE
	arena, t, tcount := cAttributeList(temp)
	defer arena.Free()

	e := C.CopyObject(c.ctx, C.CK_SESSION_HANDLE(sh), C.CK_OBJECT_HANDLE(o), t, tcount, C.CK_OBJECT_HANDLE_PTR(&obj))
	e1 := toError(e)
	if e1 == nil {
		return ObjectHandle(obj), nil
	}
	return 0, e1
}


func (c *Ctx) DestroyObject(sh SessionHandle, oh ObjectHandle) error {
	e := C.DestroyObject(c.ctx, C.CK_SESSION_HANDLE(sh), C.CK_OBJECT_HANDLE(oh))
	return toError(e)
}


func (c *Ctx) GetObjectSize(sh SessionHandle, oh ObjectHandle) (uint, error) {
	var size C.CK_ULONG
	e := C.GetObjectSize(c.ctx, C.CK_SESSION_HANDLE(sh), C.CK_OBJECT_HANDLE(oh), &size)
	return uint(size), toError(e)
}


func (c *Ctx) GetAttributeValue(sh SessionHandle, o ObjectHandle, a []*Attribute) ([]*Attribute, error) {
	
	
	pa := make([]C.CK_ATTRIBUTE, len(a))
	for i := 0; i < len(a); i++ {
		pa[i]._type = C.CK_ATTRIBUTE_TYPE(a[i].Type)
	}
	e := C.GetAttributeValue(c.ctx, C.CK_SESSION_HANDLE(sh), C.CK_OBJECT_HANDLE(o), &pa[0], C.CK_ULONG(len(a)))
	if err := toError(e); err != nil {
		return nil, err
	}
	a1 := make([]*Attribute, len(a))
	for i, c := range pa {
		x := new(Attribute)
		x.Type = uint(c._type)
		if int(c.ulValueLen) != -1 {
			buf := unsafe.Pointer(C.getAttributePval(&c))
			x.Value = C.GoBytes(buf, C.int(c.ulValueLen))
			C.free(buf)
		}
		a1[i] = x
	}
	return a1, nil
}


func (c *Ctx) SetAttributeValue(sh SessionHandle, o ObjectHandle, a []*Attribute) error {
	arena, pa, palen := cAttributeList(a)
	defer arena.Free()
	e := C.SetAttributeValue(c.ctx, C.CK_SESSION_HANDLE(sh), C.CK_OBJECT_HANDLE(o), pa, palen)
	return toError(e)
}



func (c *Ctx) FindObjectsInit(sh SessionHandle, temp []*Attribute) error {
	arena, t, tcount := cAttributeList(temp)
	defer arena.Free()
	e := C.FindObjectsInit(c.ctx, C.CK_SESSION_HANDLE(sh), t, tcount)
	return toError(e)
}







func (c *Ctx) FindObjects(sh SessionHandle, max int) ([]ObjectHandle, bool, error) {
	var (
		objectList C.CK_OBJECT_HANDLE_PTR
		ulCount    C.CK_ULONG
	)
	e := C.FindObjects(c.ctx, C.CK_SESSION_HANDLE(sh), &objectList, C.CK_ULONG(max), &ulCount)
	if toError(e) != nil {
		return nil, false, toError(e)
	}
	l := toList(C.CK_ULONG_PTR(unsafe.Pointer(objectList)), ulCount)
	
	
	o := make([]ObjectHandle, len(l))
	for i, v := range l {
		o[i] = ObjectHandle(v)
	}
	return o, ulCount > C.CK_ULONG(max), nil
}


func (c *Ctx) FindObjectsFinal(sh SessionHandle) error {
	e := C.FindObjectsFinal(c.ctx, C.CK_SESSION_HANDLE(sh))
	return toError(e)
}


func (c *Ctx) EncryptInit(sh SessionHandle, m []*Mechanism, o ObjectHandle) error {
	arena, mech := cMechanism(m)
	defer arena.Free()
	e := C.EncryptInit(c.ctx, C.CK_SESSION_HANDLE(sh), mech, C.CK_OBJECT_HANDLE(o))
	return toError(e)
}


func (c *Ctx) Encrypt(sh SessionHandle, message []byte) ([]byte, error) {
	var (
		enc    C.CK_BYTE_PTR
		enclen C.CK_ULONG
	)
	e := C.Encrypt(c.ctx, C.CK_SESSION_HANDLE(sh), cMessage(message), C.CK_ULONG(len(message)), &enc, &enclen)
	if toError(e) != nil {
		return nil, toError(e)
	}
	s := C.GoBytes(unsafe.Pointer(enc), C.int(enclen))
	C.free(unsafe.Pointer(enc))
	return s, nil
}


func (c *Ctx) EncryptUpdate(sh SessionHandle, plain []byte) ([]byte, error) {
	var (
		part    C.CK_BYTE_PTR
		partlen C.CK_ULONG
	)
	e := C.EncryptUpdate(c.ctx, C.CK_SESSION_HANDLE(sh), cMessage(plain), C.CK_ULONG(len(plain)), &part, &partlen)
	if toError(e) != nil {
		return nil, toError(e)
	}
	h := C.GoBytes(unsafe.Pointer(part), C.int(partlen))
	C.free(unsafe.Pointer(part))
	return h, nil
}


func (c *Ctx) EncryptFinal(sh SessionHandle) ([]byte, error) {
	var (
		enc    C.CK_BYTE_PTR
		enclen C.CK_ULONG
	)
	e := C.EncryptFinal(c.ctx, C.CK_SESSION_HANDLE(sh), &enc, &enclen)
	if toError(e) != nil {
		return nil, toError(e)
	}
	h := C.GoBytes(unsafe.Pointer(enc), C.int(enclen))
	C.free(unsafe.Pointer(enc))
	return h, nil
}


func (c *Ctx) DecryptInit(sh SessionHandle, m []*Mechanism, o ObjectHandle) error {
	arena, mech := cMechanism(m)
	defer arena.Free()
	e := C.DecryptInit(c.ctx, C.CK_SESSION_HANDLE(sh), mech, C.CK_OBJECT_HANDLE(o))
	return toError(e)
}


func (c *Ctx) Decrypt(sh SessionHandle, cipher []byte) ([]byte, error) {
	var (
		plain    C.CK_BYTE_PTR
		plainlen C.CK_ULONG
	)
	e := C.Decrypt(c.ctx, C.CK_SESSION_HANDLE(sh), cMessage(cipher), C.CK_ULONG(len(cipher)), &plain, &plainlen)
	if toError(e) != nil {
		return nil, toError(e)
	}
	s := C.GoBytes(unsafe.Pointer(plain), C.int(plainlen))
	C.free(unsafe.Pointer(plain))
	return s, nil
}


func (c *Ctx) DecryptUpdate(sh SessionHandle, cipher []byte) ([]byte, error) {
	var (
		part    C.CK_BYTE_PTR
		partlen C.CK_ULONG
	)
	e := C.DecryptUpdate(c.ctx, C.CK_SESSION_HANDLE(sh), cMessage(cipher), C.CK_ULONG(len(cipher)), &part, &partlen)
	if toError(e) != nil {
		return nil, toError(e)
	}
	h := C.GoBytes(unsafe.Pointer(part), C.int(partlen))
	C.free(unsafe.Pointer(part))
	return h, nil
}


func (c *Ctx) DecryptFinal(sh SessionHandle) ([]byte, error) {
	var (
		plain    C.CK_BYTE_PTR
		plainlen C.CK_ULONG
	)
	e := C.DecryptFinal(c.ctx, C.CK_SESSION_HANDLE(sh), &plain, &plainlen)
	if toError(e) != nil {
		return nil, toError(e)
	}
	h := C.GoBytes(unsafe.Pointer(plain), C.int(plainlen))
	C.free(unsafe.Pointer(plain))
	return h, nil
}


func (c *Ctx) DigestInit(sh SessionHandle, m []*Mechanism) error {
	arena, mech := cMechanism(m)
	defer arena.Free()
	e := C.DigestInit(c.ctx, C.CK_SESSION_HANDLE(sh), mech)
	return toError(e)
}


func (c *Ctx) Digest(sh SessionHandle, message []byte) ([]byte, error) {
	var (
		hash    C.CK_BYTE_PTR
		hashlen C.CK_ULONG
	)
	e := C.Digest(c.ctx, C.CK_SESSION_HANDLE(sh), cMessage(message), C.CK_ULONG(len(message)), &hash, &hashlen)
	if toError(e) != nil {
		return nil, toError(e)
	}
	h := C.GoBytes(unsafe.Pointer(hash), C.int(hashlen))
	C.free(unsafe.Pointer(hash))
	return h, nil
}


func (c *Ctx) DigestUpdate(sh SessionHandle, message []byte) error {
	e := C.DigestUpdate(c.ctx, C.CK_SESSION_HANDLE(sh), cMessage(message), C.CK_ULONG(len(message)))
	if toError(e) != nil {
		return toError(e)
	}
	return nil
}




func (c *Ctx) DigestKey(sh SessionHandle, key ObjectHandle) error {
	e := C.DigestKey(c.ctx, C.CK_SESSION_HANDLE(sh), C.CK_OBJECT_HANDLE(key))
	if toError(e) != nil {
		return toError(e)
	}
	return nil
}


func (c *Ctx) DigestFinal(sh SessionHandle) ([]byte, error) {
	var (
		hash    C.CK_BYTE_PTR
		hashlen C.CK_ULONG
	)
	e := C.DigestFinal(c.ctx, C.CK_SESSION_HANDLE(sh), &hash, &hashlen)
	if toError(e) != nil {
		return nil, toError(e)
	}
	h := C.GoBytes(unsafe.Pointer(hash), C.int(hashlen))
	C.free(unsafe.Pointer(hash))
	return h, nil
}




func (c *Ctx) SignInit(sh SessionHandle, m []*Mechanism, o ObjectHandle) error {
	arena, mech := cMechanism(m)
	defer arena.Free()
	e := C.SignInit(c.ctx, C.CK_SESSION_HANDLE(sh), mech, C.CK_OBJECT_HANDLE(o))
	return toError(e)
}



func (c *Ctx) Sign(sh SessionHandle, message []byte) ([]byte, error) {
	var (
		sig    C.CK_BYTE_PTR
		siglen C.CK_ULONG
	)
	e := C.Sign(c.ctx, C.CK_SESSION_HANDLE(sh), cMessage(message), C.CK_ULONG(len(message)), &sig, &siglen)
	if toError(e) != nil {
		return nil, toError(e)
	}
	s := C.GoBytes(unsafe.Pointer(sig), C.int(siglen))
	C.free(unsafe.Pointer(sig))
	return s, nil
}




func (c *Ctx) SignUpdate(sh SessionHandle, message []byte) error {
	e := C.SignUpdate(c.ctx, C.CK_SESSION_HANDLE(sh), cMessage(message), C.CK_ULONG(len(message)))
	return toError(e)
}


func (c *Ctx) SignFinal(sh SessionHandle) ([]byte, error) {
	var (
		sig    C.CK_BYTE_PTR
		siglen C.CK_ULONG
	)
	e := C.SignFinal(c.ctx, C.CK_SESSION_HANDLE(sh), &sig, &siglen)
	if toError(e) != nil {
		return nil, toError(e)
	}
	h := C.GoBytes(unsafe.Pointer(sig), C.int(siglen))
	C.free(unsafe.Pointer(sig))
	return h, nil
}


func (c *Ctx) SignRecoverInit(sh SessionHandle, m []*Mechanism, key ObjectHandle) error {
	arena, mech := cMechanism(m)
	defer arena.Free()
	e := C.SignRecoverInit(c.ctx, C.CK_SESSION_HANDLE(sh), mech, C.CK_OBJECT_HANDLE(key))
	return toError(e)
}


func (c *Ctx) SignRecover(sh SessionHandle, data []byte) ([]byte, error) {
	var (
		sig    C.CK_BYTE_PTR
		siglen C.CK_ULONG
	)
	e := C.SignRecover(c.ctx, C.CK_SESSION_HANDLE(sh), cMessage(data), C.CK_ULONG(len(data)), &sig, &siglen)
	if toError(e) != nil {
		return nil, toError(e)
	}
	h := C.GoBytes(unsafe.Pointer(sig), C.int(siglen))
	C.free(unsafe.Pointer(sig))
	return h, nil
}




func (c *Ctx) VerifyInit(sh SessionHandle, m []*Mechanism, key ObjectHandle) error {
	arena, mech := cMechanism(m)
	defer arena.Free()
	e := C.VerifyInit(c.ctx, C.CK_SESSION_HANDLE(sh), mech, C.CK_OBJECT_HANDLE(key))
	return toError(e)
}




func (c *Ctx) Verify(sh SessionHandle, data []byte, signature []byte) error {
	e := C.Verify(c.ctx, C.CK_SESSION_HANDLE(sh), cMessage(data), C.CK_ULONG(len(data)), cMessage(signature), C.CK_ULONG(len(signature)))
	return toError(e)
}




func (c *Ctx) VerifyUpdate(sh SessionHandle, part []byte) error {
	e := C.VerifyUpdate(c.ctx, C.CK_SESSION_HANDLE(sh), cMessage(part), C.CK_ULONG(len(part)))
	return toError(e)
}



func (c *Ctx) VerifyFinal(sh SessionHandle, signature []byte) error {
	e := C.VerifyFinal(c.ctx, C.CK_SESSION_HANDLE(sh), cMessage(signature), C.CK_ULONG(len(signature)))
	return toError(e)
}



func (c *Ctx) VerifyRecoverInit(sh SessionHandle, m []*Mechanism, key ObjectHandle) error {
	arena, mech := cMechanism(m)
	defer arena.Free()
	e := C.VerifyRecoverInit(c.ctx, C.CK_SESSION_HANDLE(sh), mech, C.CK_OBJECT_HANDLE(key))
	return toError(e)
}



func (c *Ctx) VerifyRecover(sh SessionHandle, signature []byte) ([]byte, error) {
	var (
		data    C.CK_BYTE_PTR
		datalen C.CK_ULONG
	)
	e := C.DecryptVerifyUpdate(c.ctx, C.CK_SESSION_HANDLE(sh), cMessage(signature), C.CK_ULONG(len(signature)), &data, &datalen)
	if toError(e) != nil {
		return nil, toError(e)
	}
	h := C.GoBytes(unsafe.Pointer(data), C.int(datalen))
	C.free(unsafe.Pointer(data))
	return h, nil
}


func (c *Ctx) DigestEncryptUpdate(sh SessionHandle, part []byte) ([]byte, error) {
	var (
		enc    C.CK_BYTE_PTR
		enclen C.CK_ULONG
	)
	e := C.DigestEncryptUpdate(c.ctx, C.CK_SESSION_HANDLE(sh), cMessage(part), C.CK_ULONG(len(part)), &enc, &enclen)
	if toError(e) != nil {
		return nil, toError(e)
	}
	h := C.GoBytes(unsafe.Pointer(enc), C.int(enclen))
	C.free(unsafe.Pointer(enc))
	return h, nil
}


func (c *Ctx) DecryptDigestUpdate(sh SessionHandle, cipher []byte) ([]byte, error) {
	var (
		part    C.CK_BYTE_PTR
		partlen C.CK_ULONG
	)
	e := C.DecryptDigestUpdate(c.ctx, C.CK_SESSION_HANDLE(sh), cMessage(cipher), C.CK_ULONG(len(cipher)), &part, &partlen)
	if toError(e) != nil {
		return nil, toError(e)
	}
	h := C.GoBytes(unsafe.Pointer(part), C.int(partlen))
	C.free(unsafe.Pointer(part))
	return h, nil
}


func (c *Ctx) SignEncryptUpdate(sh SessionHandle, part []byte) ([]byte, error) {
	var (
		enc    C.CK_BYTE_PTR
		enclen C.CK_ULONG
	)
	e := C.SignEncryptUpdate(c.ctx, C.CK_SESSION_HANDLE(sh), cMessage(part), C.CK_ULONG(len(part)), &enc, &enclen)
	if toError(e) != nil {
		return nil, toError(e)
	}
	h := C.GoBytes(unsafe.Pointer(enc), C.int(enclen))
	C.free(unsafe.Pointer(enc))
	return h, nil
}


func (c *Ctx) DecryptVerifyUpdate(sh SessionHandle, cipher []byte) ([]byte, error) {
	var (
		part    C.CK_BYTE_PTR
		partlen C.CK_ULONG
	)
	e := C.DecryptVerifyUpdate(c.ctx, C.CK_SESSION_HANDLE(sh), cMessage(cipher), C.CK_ULONG(len(cipher)), &part, &partlen)
	if toError(e) != nil {
		return nil, toError(e)
	}
	h := C.GoBytes(unsafe.Pointer(part), C.int(partlen))
	C.free(unsafe.Pointer(part))
	return h, nil
}


func (c *Ctx) GenerateKey(sh SessionHandle, m []*Mechanism, temp []*Attribute) (ObjectHandle, error) {
	var key C.CK_OBJECT_HANDLE
	attrarena, t, tcount := cAttributeList(temp)
	defer attrarena.Free()
	mecharena, mech := cMechanism(m)
	defer mecharena.Free()
	e := C.GenerateKey(c.ctx, C.CK_SESSION_HANDLE(sh), mech, t, tcount, C.CK_OBJECT_HANDLE_PTR(&key))
	e1 := toError(e)
	if e1 == nil {
		return ObjectHandle(key), nil
	}
	return 0, e1
}


func (c *Ctx) GenerateKeyPair(sh SessionHandle, m []*Mechanism, public, private []*Attribute) (ObjectHandle, ObjectHandle, error) {
	var (
		pubkey  C.CK_OBJECT_HANDLE
		privkey C.CK_OBJECT_HANDLE
	)
	pubarena, pub, pubcount := cAttributeList(public)
	defer pubarena.Free()
	privarena, priv, privcount := cAttributeList(private)
	defer privarena.Free()
	mecharena, mech := cMechanism(m)
	defer mecharena.Free()
	e := C.GenerateKeyPair(c.ctx, C.CK_SESSION_HANDLE(sh), mech, pub, pubcount, priv, privcount, C.CK_OBJECT_HANDLE_PTR(&pubkey), C.CK_OBJECT_HANDLE_PTR(&privkey))
	e1 := toError(e)
	if e1 == nil {
		return ObjectHandle(pubkey), ObjectHandle(privkey), nil
	}
	return 0, 0, e1
}


func (c *Ctx) WrapKey(sh SessionHandle, m []*Mechanism, wrappingkey, key ObjectHandle) ([]byte, error) {
	var (
		wrappedkey    C.CK_BYTE_PTR
		wrappedkeylen C.CK_ULONG
	)
	arena, mech := cMechanism(m)
	defer arena.Free()
	e := C.WrapKey(c.ctx, C.CK_SESSION_HANDLE(sh), mech, C.CK_OBJECT_HANDLE(wrappingkey), C.CK_OBJECT_HANDLE(key), &wrappedkey, &wrappedkeylen)
	if toError(e) != nil {
		return nil, toError(e)
	}
	h := C.GoBytes(unsafe.Pointer(wrappedkey), C.int(wrappedkeylen))
	C.free(unsafe.Pointer(wrappedkey))
	return h, nil
}


func (c *Ctx) UnwrapKey(sh SessionHandle, m []*Mechanism, unwrappingkey ObjectHandle, wrappedkey []byte, a []*Attribute) (ObjectHandle, error) {
	var key C.CK_OBJECT_HANDLE
	attrarena, ac, aclen := cAttributeList(a)
	defer attrarena.Free()
	mecharena, mech := cMechanism(m)
	defer mecharena.Free()
	e := C.UnwrapKey(c.ctx, C.CK_SESSION_HANDLE(sh), mech, C.CK_OBJECT_HANDLE(unwrappingkey), C.CK_BYTE_PTR(unsafe.Pointer(&wrappedkey[0])), C.CK_ULONG(len(wrappedkey)), ac, aclen, &key)
	return ObjectHandle(key), toError(e)
}


func (c *Ctx) DeriveKey(sh SessionHandle, m []*Mechanism, basekey ObjectHandle, a []*Attribute) (ObjectHandle, error) {
	var key C.CK_OBJECT_HANDLE
	attrarena, ac, aclen := cAttributeList(a)
	defer attrarena.Free()
	mecharena, mech := cMechanism(m)
	defer mecharena.Free()
	e := C.DeriveKey(c.ctx, C.CK_SESSION_HANDLE(sh), mech, C.CK_OBJECT_HANDLE(basekey), ac, aclen, &key)
	return ObjectHandle(key), toError(e)
}



func (c *Ctx) SeedRandom(sh SessionHandle, seed []byte) error {
	e := C.SeedRandom(c.ctx, C.CK_SESSION_HANDLE(sh), C.CK_BYTE_PTR(unsafe.Pointer(&seed[0])), C.CK_ULONG(len(seed)))
	return toError(e)
}


func (c *Ctx) GenerateRandom(sh SessionHandle, length int) ([]byte, error) {
	var rand C.CK_BYTE_PTR
	e := C.GenerateRandom(c.ctx, C.CK_SESSION_HANDLE(sh), &rand, C.CK_ULONG(length))
	if toError(e) != nil {
		return nil, toError(e)
	}
	h := C.GoBytes(unsafe.Pointer(rand), C.int(length))
	C.free(unsafe.Pointer(rand))
	return h, nil
}



func (c *Ctx) WaitForSlotEvent(flags uint) chan SlotEvent {
	sl := make(chan SlotEvent, 1) 
	go c.waitForSlotEventHelper(flags, sl)
	return sl
}

func (c *Ctx) waitForSlotEventHelper(f uint, sl chan SlotEvent) {
	var slotID C.CK_ULONG
	C.WaitForSlotEvent(c.ctx, C.CK_FLAGS(f), &slotID)
	sl <- SlotEvent{uint(slotID)}
	close(sl) 
}
