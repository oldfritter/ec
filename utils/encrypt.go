package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	// "fmt"
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
	var en string
	for i, str := range strs {
		if i > 0 {
			en = en + "&"
		}
		e, err = rsa.EncryptPKCS1v15(rand.Reader, pub, []byte(str))
		if err != nil {
			return
		}
		en = en + base64.StdEncoding.EncodeToString(e)
		encrypted = []byte(en)
	}
	return
}

func PrivateKeyDecrypt(encrypted string, privateKey *rsa.PrivateKey) (message string, err error) {
	decoded, _ := base64.StdEncoding.DecodeString(encrypted)
	ms := strings.Split(string(decoded), "&")
	for _, m := range ms {
		var s []byte
		s, err = rsa.DecryptPKCS1v15(rand.Reader, privateKey, []byte(m))
		if err != nil {
			return
		}
		message += string(s)
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
