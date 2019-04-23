/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package idemix

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/asn1"
	"math/big"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain-amcl/amcl"
	"github.com/mcc-github/blockchain-amcl/amcl/FP256BN"
	"github.com/pkg/errors"
)

type RevocationAlgorithm int32

const (
	ALG_NO_REVOCATION RevocationAlgorithm = iota
)

var ProofBytes = map[RevocationAlgorithm]int{
	ALG_NO_REVOCATION: 0,
}


func GenerateLongTermRevocationKey() (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
}





func CreateCRI(key *ecdsa.PrivateKey, unrevokedHandles []*FP256BN.BIG, epoch int, alg RevocationAlgorithm, rng *amcl.RAND) (*CredentialRevocationInformation, error) {
	if key == nil || rng == nil {
		return nil, errors.Errorf("CreateCRI received nil input")
	}
	cri := &CredentialRevocationInformation{}
	cri.RevocationAlg = int32(alg)
	cri.Epoch = int64(epoch)

	if alg == ALG_NO_REVOCATION {
		
		cri.EpochPk = Ecp2ToProto(GenG2)
	} else {
		
		_, epochPk := WBBKeyGen(rng)
		cri.EpochPk = Ecp2ToProto(epochPk)
	}

	
	bytesToSign, err := proto.Marshal(cri)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal CRI")
	}

	digest := sha256.Sum256(bytesToSign)

	cri.EpochPkSig, err = key.Sign(rand.Reader, digest[:], nil)
	if err != nil {
		return nil, err
	}

	if alg == ALG_NO_REVOCATION {
		return cri, nil
	} else {
		return nil, errors.Errorf("the specified revocation algorithm is not supported.")
	}
}






func VerifyEpochPK(pk *ecdsa.PublicKey, epochPK *ECP2, epochPkSig []byte, epoch int, alg RevocationAlgorithm) error {
	if pk == nil || epochPK == nil {
		return errors.Errorf("EpochPK invalid: received nil input")
	}
	cri := &CredentialRevocationInformation{}
	cri.RevocationAlg = int32(alg)
	cri.EpochPk = epochPK
	cri.Epoch = int64(epoch)
	bytesToSign, err := proto.Marshal(cri)
	if err != nil {
		return err
	}
	digest := sha256.Sum256(bytesToSign)

	var sig struct{ R, S *big.Int }
	if _, err := asn1.Unmarshal(epochPkSig, &sig); err != nil {
		return errors.Wrap(err, "failed unmashalling signature")
	}

	if !ecdsa.Verify(pk, digest[:], sig.R, sig.S) {
		return errors.Errorf("EpochPKSig invalid")
	}

	return nil
}
