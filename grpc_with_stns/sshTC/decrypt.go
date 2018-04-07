package sshTC

import (
	"io/ioutil"
	"encoding/pem"
	"errors"
	"fmt"
	"crypto/x509"
	"encoding/base64"
	"crypto/rsa"
	"crypto/rand"
)

func (tc *sshTC) Decrypt(s string, key *rsa.PrivateKey) (string, error){
	in, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}

	decrypted, err := rsa.DecryptPKCS1v15(rand.Reader, key, in)
	if err != nil {
		return "", err
	}

	return string(decrypted), nil
}


func  (tc *sshTC) readPrivateKey(path string) (*rsa.PrivateKey, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("invalid private key data")
	}

	if got, want := block.Type, "RSA PRIVATE KEY"; got != want {
		return nil, errors.New(fmt.Sprintf("invalid key type: %s", block.Type))
	}

	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return key, nil
}
