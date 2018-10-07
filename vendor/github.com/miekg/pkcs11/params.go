



package pkcs11


import "C"
import "unsafe"


type GCMParams struct {
	arena
	params  *C.CK_GCM_PARAMS
	iv      []byte
	aad     []byte
	tagSize int
}

















func NewGCMParams(iv, aad []byte, tagSize int) *GCMParams {
	return &GCMParams{
		iv:      iv,
		aad:     aad,
		tagSize: tagSize,
	}
}

func cGCMParams(p *GCMParams) []byte {
	params := C.CK_GCM_PARAMS{
		ulTagBits: C.CK_ULONG(p.tagSize),
	}
	var arena arena
	if len(p.iv) > 0 {
		iv, ivLen := arena.Allocate(p.iv)
		params.pIv = C.CK_BYTE_PTR(iv)
		params.ulIvLen = ivLen
		params.ulIvBits = ivLen * 8
	}
	if len(p.aad) > 0 {
		aad, aadLen := arena.Allocate(p.aad)
		params.pAAD = C.CK_BYTE_PTR(aad)
		params.ulAADLen = aadLen
	}
	p.Free()
	p.arena = arena
	p.params = &params
	return C.GoBytes(unsafe.Pointer(&params), C.int(unsafe.Sizeof(params)))
}





func (p *GCMParams) IV() []byte {
	if p == nil || p.params == nil {
		return nil
	}
	newIv := C.GoBytes(unsafe.Pointer(p.params.pIv), C.int(p.params.ulIvLen))
	iv := make([]byte, len(newIv))
	copy(iv, newIv)
	return iv
}





func (p *GCMParams) Free() {
	if p == nil || p.arena == nil {
		return
	}
	p.arena.Free()
	p.params = nil
	p.arena = nil
}


func NewPSSParams(hashAlg, mgf, saltLength uint) []byte {
	p := C.CK_RSA_PKCS_PSS_PARAMS{
		hashAlg: C.CK_MECHANISM_TYPE(hashAlg),
		mgf:     C.CK_RSA_PKCS_MGF_TYPE(mgf),
		sLen:    C.CK_ULONG(saltLength),
	}
	return C.GoBytes(unsafe.Pointer(&p), C.int(unsafe.Sizeof(p)))
}


type OAEPParams struct {
	HashAlg    uint
	MGF        uint
	SourceType uint
	SourceData []byte
}


func NewOAEPParams(hashAlg, mgf, sourceType uint, sourceData []byte) *OAEPParams {
	return &OAEPParams{
		HashAlg:    hashAlg,
		MGF:        mgf,
		SourceType: sourceType,
		SourceData: sourceData,
	}
}

func cOAEPParams(p *OAEPParams, arena arena) ([]byte, arena) {
	params := C.CK_RSA_PKCS_OAEP_PARAMS{
		hashAlg: C.CK_MECHANISM_TYPE(p.HashAlg),
		mgf:     C.CK_RSA_PKCS_MGF_TYPE(p.MGF),
		source:  C.CK_RSA_PKCS_OAEP_SOURCE_TYPE(p.SourceType),
	}
	if len(p.SourceData) != 0 {
		buf, len := arena.Allocate(p.SourceData)
		
		C.putOAEPParams(&params, buf, len)
	}
	return C.GoBytes(unsafe.Pointer(&params), C.int(unsafe.Sizeof(params))), arena
}
