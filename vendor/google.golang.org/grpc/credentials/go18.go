



package credentials

import (
	"crypto/tls"
)

func init() {
	cipherSuiteLookup[tls.TLS_RSA_WITH_AES_128_CBC_SHA256] = "TLS_RSA_WITH_AES_128_CBC_SHA256"
	cipherSuiteLookup[tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256] = "TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256"
	cipherSuiteLookup[tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256] = "TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256"
	cipherSuiteLookup[tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305] = "TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305"
	cipherSuiteLookup[tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305] = "TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305"
}






func cloneTLSConfig(cfg *tls.Config) *tls.Config {
	if cfg == nil {
		return &tls.Config{}
	}

	return cfg.Clone()
}
