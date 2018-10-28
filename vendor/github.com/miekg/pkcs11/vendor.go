



package pkcs11


const (
        NFCK_VENDOR_NCIPHER     = 0xde436972
        CKA_NCIPHER     = NFCK_VENDOR_NCIPHER
        CKM_NCIPHER     = NFCK_VENDOR_NCIPHER
        CKK_NCIPHER     = NFCK_VENDOR_NCIPHER
)


const (
	CKM_NC_SHA_1_HMAC_KEY_GEN       = (CKM_NCIPHER + 0x3)  
	CKM_NC_MD5_HMAC_KEY_GEN         = (CKM_NCIPHER + 0x6)  
	CKM_NC_SHA224_HMAC_KEY_GEN      = (CKM_NCIPHER + 0x24) 
	CKM_NC_SHA256_HMAC_KEY_GEN      = (CKM_NCIPHER + 0x25) 
	CKM_NC_SHA384_HMAC_KEY_GEN      = (CKM_NCIPHER + 0x26) 
	CKM_NC_SHA512_HMAC_KEY_GEN      = (CKM_NCIPHER + 0x27) 

)
