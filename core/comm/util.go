/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package comm

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/common/util"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
)


func AddPemToCertPool(pemCerts []byte, pool *x509.CertPool) error {
	certs, _, err := pemToX509Certs(pemCerts)
	if err != nil {
		return err
	}
	for _, cert := range certs {
		pool.AddCert(cert)
	}
	return nil
}


func pemToX509Certs(pemCerts []byte) ([]*x509.Certificate, []string, error) {

	
	certs := []*x509.Certificate{}
	subjects := []string{}
	for len(pemCerts) > 0 {
		var block *pem.Block
		block, pemCerts = pem.Decode(pemCerts)
		if block == nil {
			break
		}
		

		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, subjects, err
		} else {
			certs = append(certs, cert)
			
			subjects = append(subjects, string(cert.RawSubject))
		}
	}
	return certs, subjects, nil
}



type BindingInspector func(context.Context, proto.Message) error


type CertHashExtractor func(proto.Message) []byte




func NewBindingInspector(mutualTLS bool, extractTLSCertHash CertHashExtractor) BindingInspector {
	if extractTLSCertHash == nil {
		panic(errors.New("extractTLSCertHash parameter is nil"))
	}
	inspectMessage := mutualTLSBinding
	if !mutualTLS {
		inspectMessage = noopBinding
	}
	return func(ctx context.Context, msg proto.Message) error {
		if msg == nil {
			return errors.New("message is nil")
		}
		return inspectMessage(ctx, extractTLSCertHash(msg))
	}
}







func mutualTLSBinding(ctx context.Context, claimedTLScertHash []byte) error {
	if len(claimedTLScertHash) == 0 {
		return errors.Errorf("client didn't include its TLS cert hash")
	}
	actualTLScertHash := ExtractCertificateHashFromContext(ctx)
	if len(actualTLScertHash) == 0 {
		return errors.Errorf("client didn't send a TLS certificate")
	}
	if !bytes.Equal(actualTLScertHash, claimedTLScertHash) {
		return errors.Errorf("claimed TLS cert hash is %v but actual TLS cert hash is %v", claimedTLScertHash, actualTLScertHash)
	}
	return nil
}


func noopBinding(_ context.Context, _ []byte) error {
	return nil
}



func ExtractCertificateHashFromContext(ctx context.Context) []byte {
	rawCert := ExtractCertificateFromContext(ctx)
	if len(rawCert) == 0 {
		return nil
	}
	return util.ComputeSHA256(rawCert)
}



func ExtractCertificateFromContext(ctx context.Context) []byte {
	pr, extracted := peer.FromContext(ctx)
	if !extracted {
		return nil
	}

	authInfo := pr.AuthInfo
	if authInfo == nil {
		return nil
	}

	tlsInfo, isTLSConn := authInfo.(credentials.TLSInfo)
	if !isTLSConn {
		return nil
	}
	certs := tlsInfo.State.PeerCertificates
	if len(certs) == 0 {
		return nil
	}
	return certs[0].Raw
}
