package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"strings"
)

func PublicKeyEncrypt(message, publicKey string) (encrypted []byte, err error) {
	block, _ := pem.Decode([]byte(publicKey))
	if block == nil {
		err = errors.New("public key error")
		return
	}
	pub, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return
	}
	var strs []string
	var str string
	for i, b := range []byte(message) {
		j := i % 240
		if j == 0 {
			if i > 0 {
				strs = append(strs, str)
			}
			str = ""
		}
		str = str + string(b)
	}
	strs = append(strs, str)
	var e []byte
	for i, str := range strs {
		if i > 0 {
			encrypted = append(encrypted, []byte("|")...)
		}
		e, err = rsa.EncryptPKCS1v15(rand.Reader, pub, []byte(str))
		if err != nil {
			return
		}
		encrypted = append(encrypted, e...)
	}
	return
}

func PrivateKeyDecrypt(encrypted string, privateKey *rsa.PrivateKey) (message []byte, err error) {
	ms := strings.Split(encrypted, "|")
	for _, m := range ms {
		var s []byte
		s, err = rsa.DecryptPKCS1v15(rand.Reader, privateKey, []byte(m))
		if err != nil {
			return
		}
		message = append(message, s...)
	}
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
