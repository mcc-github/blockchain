/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package sw

import (
	"bytes"
	"crypto/aes"
	"crypto/rand"
	"io"
	"math/big"
	mrand "math/rand"
	"testing"

	"github.com/mcc-github/blockchain/bccsp"
	"github.com/mcc-github/blockchain/bccsp/mocks"
	"github.com/mcc-github/blockchain/bccsp/utils"
	"github.com/stretchr/testify/assert"
)


func TestCBCPKCS7EncryptCBCPKCS7Decrypt(t *testing.T) {
	t.Parallel()

	
	
	key := make([]byte, 32)
	rand.Reader.Read(key)

	
	var ptext = []byte("a message with arbitrary length (42 bytes)")

	encrypted, encErr := AESCBCPKCS7Encrypt(key, ptext)
	if encErr != nil {
		t.Fatalf("Error encrypting '%s': %s", ptext, encErr)
	}

	decrypted, dErr := AESCBCPKCS7Decrypt(key, encrypted)
	if dErr != nil {
		t.Fatalf("Error decrypting the encrypted '%s': %v", ptext, dErr)
	}

	if string(ptext[:]) != string(decrypted[:]) {
		t.Fatal("Decrypt( Encrypt( ptext ) ) != ptext: Ciphertext decryption with the same key must result in the original plaintext!")
	}
}


func TestPKCS7Padding(t *testing.T) {
	t.Parallel()

	
	ptext := []byte("")
	expected := []byte{16, 16, 16, 16,
		16, 16, 16, 16,
		16, 16, 16, 16,
		16, 16, 16, 16}
	result := pkcs7Padding(ptext)

	if !bytes.Equal(expected, result) {
		t.Fatal("Padding error! Expected: ", expected, "', received: '", result, "'")
	}

	
	ptext = []byte("1")
	expected = []byte{'1', 15, 15, 15,
		15, 15, 15, 15,
		15, 15, 15, 15,
		15, 15, 15, 15}
	result = pkcs7Padding(ptext)

	if !bytes.Equal(expected, result) {
		t.Fatal("Padding error! Expected: '", expected, "', received: '", result, "'")
	}

	
	ptext = []byte("12")
	expected = []byte{'1', '2', 14, 14,
		14, 14, 14, 14,
		14, 14, 14, 14,
		14, 14, 14, 14}
	result = pkcs7Padding(ptext)

	if !bytes.Equal(expected, result) {
		t.Fatal("Padding error! Expected: '", expected, "', received: '", result, "'")
	}

	
	ptext = []byte("1234567890ABCDEF")
	for i := 3; i < aes.BlockSize; i++ {
		result := pkcs7Padding(ptext[:i])

		padding := aes.BlockSize - i
		expectedPadding := bytes.Repeat([]byte{byte(padding)}, padding)
		expected = append(ptext[:i], expectedPadding...)

		if !bytes.Equal(result, expected) {
			t.Fatal("Padding error! Expected: '", expected, "', received: '", result, "'")
		}
	}

	
	ptext = bytes.Repeat([]byte{byte('x')}, aes.BlockSize)
	result = pkcs7Padding(ptext)

	expectedPadding := bytes.Repeat([]byte{byte(aes.BlockSize)}, aes.BlockSize)
	expected = append(ptext, expectedPadding...)

	if len(result) != 2*aes.BlockSize {
		t.Fatal("Padding error: expected the length of the returned slice to be 2 times aes.BlockSize")
	}

	if !bytes.Equal(expected, result) {
		t.Fatal("Padding error! Expected: '", expected, "', received: '", result, "'")
	}
}


func TestPKCS7UnPadding(t *testing.T) {
	t.Parallel()

	
	expected := []byte("")
	ptext := []byte{16, 16, 16, 16,
		16, 16, 16, 16,
		16, 16, 16, 16,
		16, 16, 16, 16}

	result, _ := pkcs7UnPadding(ptext)

	if !bytes.Equal(expected, result) {
		t.Fatal("UnPadding error! Expected: '", expected, "', received: '", result, "'")
	}

	
	expected = []byte("1")
	ptext = []byte{'1', 15, 15, 15,
		15, 15, 15, 15,
		15, 15, 15, 15,
		15, 15, 15, 15}

	result, _ = pkcs7UnPadding(ptext)

	if !bytes.Equal(expected, result) {
		t.Fatal("UnPadding error! Expected: '", expected, "', received: '", result, "'")
	}

	
	expected = []byte("12")
	ptext = []byte{'1', '2', 14, 14,
		14, 14, 14, 14,
		14, 14, 14, 14,
		14, 14, 14, 14}

	result, _ = pkcs7UnPadding(ptext)

	if !bytes.Equal(expected, result) {
		t.Fatal("UnPadding error! Expected: '", expected, "', received: '", result, "'")
	}

	
	base := []byte("1234567890ABCDEF")
	for i := 3; i < aes.BlockSize; i++ {
		iPad := aes.BlockSize - i
		padding := bytes.Repeat([]byte{byte(iPad)}, iPad)
		ptext = append(base[:i], padding...)

		expected := base[:i]
		result, _ := pkcs7UnPadding(ptext)

		if !bytes.Equal(result, expected) {
			t.Fatal("UnPadding error! Expected: '", expected, "', received: '", result, "'")
		}
	}

	
	expected = bytes.Repeat([]byte{byte('x')}, aes.BlockSize)
	padding := bytes.Repeat([]byte{byte(aes.BlockSize)}, aes.BlockSize)
	ptext = append(expected, padding...)

	result, _ = pkcs7UnPadding(ptext)

	if !bytes.Equal(expected, result) {
		t.Fatal("UnPadding error! Expected: '", expected, "', received: '", result, "'")
	}
}



func TestCBCEncryptCBCPKCS7Decrypt_BlockSizeLengthPlaintext(t *testing.T) {
	t.Parallel()

	
	
	key := make([]byte, 32)
	rand.Reader.Read(key)

	
	var ptext = []byte("a 16 byte messag")

	encrypted, encErr := aesCBCEncrypt(key, ptext)
	if encErr != nil {
		t.Fatalf("Error encrypting '%s': %v", ptext, encErr)
	}

	decrypted, dErr := AESCBCPKCS7Decrypt(key, encrypted)
	if dErr == nil {
		t.Fatalf("Expected an error decrypting ptext '%s'. Decrypted to '%v'", dErr, decrypted)
	}
}



func TestCBCPKCS7EncryptCBCDecrypt_ExpectingCorruptMessage(t *testing.T) {
	t.Parallel()

	
	
	key := make([]byte, 32)
	rand.Reader.Read(key)

	
	var ptext = []byte("a 16 byte messag")

	encrypted, encErr := AESCBCPKCS7Encrypt(key, ptext)
	if encErr != nil {
		t.Fatalf("Error encrypting ptext %v", encErr)
	}

	decrypted, dErr := aesCBCDecrypt(key, encrypted)
	if dErr != nil {
		t.Fatalf("Error encrypting ptext %v, %v", dErr, decrypted)
	}

	if string(ptext[:]) != string(decrypted[:aes.BlockSize]) {
		t.Log("ptext: ", ptext)
		t.Log("decrypted: ", decrypted[:aes.BlockSize])
		t.Fatal("Encryption->Decryption with same key should result in original ptext")
	}

	if !bytes.Equal(decrypted[aes.BlockSize:], bytes.Repeat([]byte{byte(aes.BlockSize)}, aes.BlockSize)) {
		t.Fatal("Expected extra block with padding in encrypted ptext", decrypted)
	}
}


func TestCBCPKCS7Encrypt_EmptyPlaintext(t *testing.T) {
	t.Parallel()

	key := make([]byte, 32)
	rand.Reader.Read(key)

	t.Log("Generated key: ", key)

	var emptyPlaintext = []byte("")
	t.Log("Plaintext length: ", len(emptyPlaintext))

	ciphertext, encErr := AESCBCPKCS7Encrypt(key, emptyPlaintext)
	if encErr != nil {
		t.Fatalf("Error encrypting '%v'", encErr)
	}

	
	
	const expectedLength = aes.BlockSize + aes.BlockSize
	if len(ciphertext) != expectedLength {
		t.Fatalf("Wrong ciphertext length. Expected %d, received %d", expectedLength, len(ciphertext))
	}

	t.Log("Ciphertext length: ", len(ciphertext))
	t.Log("Cipher: ", ciphertext)
}


func TestCBCEncrypt_EmptyPlaintext(t *testing.T) {
	t.Parallel()

	key := make([]byte, 32)
	rand.Reader.Read(key)
	t.Log("Generated key: ", key)

	var emptyPlaintext = []byte("")
	t.Log("Message length: ", len(emptyPlaintext))

	ciphertext, encErr := aesCBCEncrypt(key, emptyPlaintext)
	assert.NoError(t, encErr)

	t.Log("Ciphertext length: ", len(ciphertext))

	
	var expectedLength = aes.BlockSize

	if len(ciphertext) != expectedLength {
		t.Fatalf("Wrong ciphertext length. Expected: '%d', received: '%d'", expectedLength, len(ciphertext))
	}
	t.Log("Ciphertext: ", ciphertext)
}


func TestCBCPKCS7Encrypt_VerifyRandomIVs(t *testing.T) {
	t.Parallel()

	key := make([]byte, aes.BlockSize)
	rand.Reader.Read(key)
	t.Log("Key 1", key)

	var ptext = []byte("a message to encrypt")

	ciphertext1, err := AESCBCPKCS7Encrypt(key, ptext)
	if err != nil {
		t.Fatalf("Error encrypting '%s': %s", ptext, err)
	}

	
	ciphertext2, err := AESCBCPKCS7Encrypt(key, ptext)
	if err != nil {
		t.Fatalf("Error encrypting '%s': %s", ptext, err)
	}

	iv1 := ciphertext1[:aes.BlockSize]
	iv2 := ciphertext2[:aes.BlockSize]

	t.Log("Ciphertext1: ", iv1)
	t.Log("Ciphertext2: ", iv2)
	t.Log("bytes.Equal: ", bytes.Equal(iv1, iv2))

	if bytes.Equal(iv1, iv2) {
		t.Fatal("Error: ciphertexts contain identical initialization vectors (IVs)")
	}
}


func TestCBCPKCS7Encrypt_CorrectCiphertextLengthCheck(t *testing.T) {
	t.Parallel()

	key := make([]byte, aes.BlockSize)
	rand.Reader.Read(key)

	
	

	var ptext = []byte("0123456789ABCDEF")

	for i := 1; i < aes.BlockSize; i++ {
		ciphertext, err := AESCBCPKCS7Encrypt(key, ptext[:i])
		if err != nil {
			t.Fatal("Error encrypting '", ptext, "'")
		}

		expectedLength := aes.BlockSize + aes.BlockSize
		if len(ciphertext) != expectedLength {
			t.Fatalf("Incorrect ciphertext incorrect: expected '%d', received '%d'", expectedLength, len(ciphertext))
		}
	}
}


func TestCBCEncryptCBCDecrypt_KeyMismatch(t *testing.T) {
	t.Parallel()

	
	key := make([]byte, aes.BlockSize)
	rand.Reader.Read(key)

	
	wrongKey := make([]byte, aes.BlockSize)
	copy(wrongKey, key[:])
	wrongKey[0] = key[0] + 1

	var ptext = []byte("1234567890ABCDEF")
	encrypted, encErr := aesCBCEncrypt(key, ptext)
	if encErr != nil {
		t.Fatalf("Error encrypting '%s': %v", ptext, encErr)
	}

	decrypted, decErr := aesCBCDecrypt(wrongKey, encrypted)
	if decErr != nil {
		t.Fatalf("Error decrypting '%s': %v", ptext, decErr)
	}

	if string(ptext[:]) == string(decrypted[:]) {
		t.Fatal("Decrypting a ciphertext with a different key than the one used for encrypting it - should not result in the original plaintext.")
	}
}


func TestCBCEncryptCBCDecrypt(t *testing.T) {
	t.Parallel()

	key := make([]byte, 32)
	rand.Reader.Read(key)

	
	var ptext = []byte("a 16 byte messag")

	encrypted, encErr := aesCBCEncrypt(key, ptext)
	if encErr != nil {
		t.Fatalf("Error encrypting '%s': %v", ptext, encErr)
	}

	decrypted, decErr := aesCBCDecrypt(key, encrypted)
	if decErr != nil {
		t.Fatalf("Error decrypting '%s': %v", ptext, decErr)
	}

	if string(ptext[:]) != string(decrypted[:]) {
		t.Fatal("Encryption->Decryption with same key should result in the original plaintext.")
	}
}


func TestCBCEncryptWithRandCBCDecrypt(t *testing.T) {
	t.Parallel()

	key := make([]byte, 32)
	rand.Reader.Read(key)

	
	var ptext = []byte("a 16 byte messag")

	encrypted, encErr := aesCBCEncryptWithRand(rand.Reader, key, ptext)
	if encErr != nil {
		t.Fatalf("Error encrypting '%s': %v", ptext, encErr)
	}

	decrypted, decErr := aesCBCDecrypt(key, encrypted)
	if decErr != nil {
		t.Fatalf("Error decrypting '%s': %v", ptext, decErr)
	}

	if string(ptext[:]) != string(decrypted[:]) {
		t.Fatal("Encryption->Decryption with same key should result in the original plaintext.")
	}
}


func TestCBCEncryptWithIVCBCDecrypt(t *testing.T) {
	t.Parallel()

	key := make([]byte, 32)
	rand.Reader.Read(key)

	
	var ptext = []byte("a 16 byte messag")

	iv := make([]byte, aes.BlockSize)
	_, err := io.ReadFull(rand.Reader, iv)
	assert.NoError(t, err)

	encrypted, encErr := aesCBCEncryptWithIV(iv, key, ptext)
	if encErr != nil {
		t.Fatalf("Error encrypting '%s': %v", ptext, encErr)
	}

	decrypted, decErr := aesCBCDecrypt(key, encrypted)
	if decErr != nil {
		t.Fatalf("Error decrypting '%s': %v", ptext, decErr)
	}

	if string(ptext[:]) != string(decrypted[:]) {
		t.Fatal("Encryption->Decryption with same key should result in the original plaintext.")
	}
}


func TestAESRelatedUtilFunctions(t *testing.T) {
	t.Parallel()

	key, err := GetRandomBytes(32)
	if err != nil {
		t.Fatalf("Failed generating AES key [%s]", err)
	}

	for i := 1; i < 100; i++ {
		l, err := rand.Int(rand.Reader, big.NewInt(1024))
		if err != nil {
			t.Fatalf("Failed generating AES key [%s]", err)
		}
		msg, err := GetRandomBytes(int(l.Int64()) + 1)
		if err != nil {
			t.Fatalf("Failed generating AES key [%s]", err)
		}

		ct, err := AESCBCPKCS7Encrypt(key, msg)
		if err != nil {
			t.Fatalf("Failed encrypting [%s]", err)
		}

		msg2, err := AESCBCPKCS7Decrypt(key, ct)
		if err != nil {
			t.Fatalf("Failed decrypting [%s]", err)
		}

		if 0 != bytes.Compare(msg, msg2) {
			t.Fatalf("Wrong decryption output [%x][%x]", msg, msg2)
		}
	}
}


func TestVariousAESKeyEncoding(t *testing.T) {
	t.Parallel()

	key, err := GetRandomBytes(32)
	if err != nil {
		t.Fatalf("Failed generating AES key [%s]", err)
	}

	
	pem := utils.AEStoPEM(key)
	keyFromPEM, err := utils.PEMtoAES(pem, nil)
	if err != nil {
		t.Fatalf("Failed converting PEM to AES key [%s]", err)
	}
	if 0 != bytes.Compare(key, keyFromPEM) {
		t.Fatalf("Failed converting PEM to AES key. Keys are different [%x][%x]", key, keyFromPEM)
	}

	
	pem, err = utils.AEStoEncryptedPEM(key, []byte("passwd"))
	if err != nil {
		t.Fatalf("Failed converting AES key to Encrypted PEM [%s]", err)
	}
	keyFromPEM, err = utils.PEMtoAES(pem, []byte("passwd"))
	if err != nil {
		t.Fatalf("Failed converting encrypted PEM to AES key [%s]", err)
	}
	if 0 != bytes.Compare(key, keyFromPEM) {
		t.Fatalf("Failed converting encrypted PEM to AES key. Keys are different [%x][%x]", key, keyFromPEM)
	}
}

func TestPkcs7UnPaddingInvalidInputs(t *testing.T) {
	t.Parallel()

	_, err := pkcs7UnPadding([]byte{1, 2, 3, 4, 5, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16})
	assert.Error(t, err)
	assert.Equal(t, "Invalid pkcs7 padding (pad[i] != unpadding)", err.Error())
}

func TestAESCBCEncryptInvalidInputs(t *testing.T) {
	t.Parallel()

	_, err := aesCBCEncrypt(nil, []byte{0, 1, 2, 3})
	assert.Error(t, err)
	assert.Equal(t, "Invalid plaintext. It must be a multiple of the block size", err.Error())

	_, err = aesCBCEncrypt([]byte{0}, []byte{1, 2, 3, 4, 5, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15})
	assert.Error(t, err)
}

func TestAESCBCDecryptInvalidInputs(t *testing.T) {
	t.Parallel()

	_, err := aesCBCDecrypt([]byte{0}, []byte{1, 2, 3, 4, 5, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15})
	assert.Error(t, err)

	_, err = aesCBCDecrypt([]byte{1, 2, 3, 4, 5, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, []byte{0})
	assert.Error(t, err)

	_, err = aesCBCDecrypt([]byte{1, 2, 3, 4, 5, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
		[]byte{1, 2, 3, 4, 5, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16})
	assert.Error(t, err)
}



func TestAESCBCPKCS7EncryptorDecrypt(t *testing.T) {
	t.Parallel()

	raw, err := GetRandomBytes(32)
	assert.NoError(t, err)

	k := &aesPrivateKey{privKey: raw, exportable: false}

	msg := []byte("Hello World")
	encryptor := &aescbcpkcs7Encryptor{}

	_, err = encryptor.Encrypt(k, msg, nil)
	assert.Error(t, err)

	_, err = encryptor.Encrypt(k, msg, &mocks.EncrypterOpts{})
	assert.Error(t, err)

	_, err = encryptor.Encrypt(k, msg, &bccsp.AESCBCPKCS7ModeOpts{IV: []byte{1}})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Invalid IV. It must have length the block size")

	_, err = encryptor.Encrypt(k, msg, &bccsp.AESCBCPKCS7ModeOpts{IV: []byte{1}, PRNG: rand.Reader})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Invalid options. Either IV or PRNG should be different from nil, or both nil.")

	ct, err := encryptor.Encrypt(k, msg, bccsp.AESCBCPKCS7ModeOpts{})
	assert.NoError(t, err)

	ct, err = encryptor.Encrypt(k, msg, &bccsp.AESCBCPKCS7ModeOpts{})
	assert.NoError(t, err)

	decryptor := &aescbcpkcs7Decryptor{}

	_, err = decryptor.Decrypt(k, ct, nil)
	assert.Error(t, err)

	_, err = decryptor.Decrypt(k, ct, &mocks.EncrypterOpts{})
	assert.Error(t, err)

	msg2, err := decryptor.Decrypt(k, ct, &bccsp.AESCBCPKCS7ModeOpts{})
	assert.NoError(t, err)
	assert.Equal(t, msg, msg2)
}

func TestAESCBCPKCS7EncryptorWithIVSameCiphertext(t *testing.T) {
	t.Parallel()

	raw, err := GetRandomBytes(32)
	assert.NoError(t, err)

	k := &aesPrivateKey{privKey: raw, exportable: false}

	msg := []byte("Hello World")
	encryptor := &aescbcpkcs7Encryptor{}

	iv := make([]byte, aes.BlockSize)

	ct, err := encryptor.Encrypt(k, msg, &bccsp.AESCBCPKCS7ModeOpts{IV: iv})
	assert.NoError(t, err)
	assert.NotNil(t, ct)
	assert.Equal(t, iv, ct[:aes.BlockSize])

	ct2, err := encryptor.Encrypt(k, msg, &bccsp.AESCBCPKCS7ModeOpts{IV: iv})
	assert.NoError(t, err)
	assert.NotNil(t, ct2)
	assert.Equal(t, iv, ct2[:aes.BlockSize])

	assert.Equal(t, ct, ct2)
}

func TestAESCBCPKCS7EncryptorWithRandSameCiphertext(t *testing.T) {
	t.Parallel()

	raw, err := GetRandomBytes(32)
	assert.NoError(t, err)

	k := &aesPrivateKey{privKey: raw, exportable: false}

	msg := []byte("Hello World")
	encryptor := &aescbcpkcs7Encryptor{}

	r := mrand.New(mrand.NewSource(0))
	iv := make([]byte, aes.BlockSize)
	_, err = io.ReadFull(r, iv)
	assert.NoError(t, err)

	r = mrand.New(mrand.NewSource(0))
	ct, err := encryptor.Encrypt(k, msg, &bccsp.AESCBCPKCS7ModeOpts{PRNG: r})
	assert.NoError(t, err)
	assert.NotNil(t, ct)
	assert.Equal(t, iv, ct[:aes.BlockSize])

	r = mrand.New(mrand.NewSource(0))
	ct2, err := encryptor.Encrypt(k, msg, &bccsp.AESCBCPKCS7ModeOpts{PRNG: r})
	assert.NoError(t, err)
	assert.NotNil(t, ct2)
	assert.Equal(t, iv, ct2[:aes.BlockSize])

	assert.Equal(t, ct, ct2)
}
