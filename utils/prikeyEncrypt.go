package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

func PublicKeyEncrypt(message []byte, publicKey string) (encrypted []byte, err error) {
	block, _ := pem.Decode([]byte(publicKey))
	if block == nil {
		err = errors.New("public key error")
		return
	}
	pub, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return
	}
	encrypted, err = rsa.EncryptPKCS1v15(rand.Reader, pub, message)
	return
}

func PrivateKeyDecrypt(encrypted []byte, privateKey *rsa.PrivateKey) (message []byte, err error) {
	message, err = rsa.DecryptPKCS1v15(rand.Reader, privateKey, encrypted)
	return
}

// bits default use 2048
func GeneratePriKey(bits int) (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, bits)
}

func GeneratePubKey(priKey *rsa.PrivateKey) ([]byte, error) {
	block := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(&priKey.PublicKey),
	}
	return pem.EncodeToMemory(block), nil
}
