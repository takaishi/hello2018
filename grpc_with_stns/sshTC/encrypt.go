package sshTC

import (
	"io/ioutil"
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"math/big"
	"crypto/rsa"
	"crypto/rand"
)

func (tc *sshTC) Encrypt(s string, pubKey *rsa.PublicKey) (string, error) {
	e, err := rsa.EncryptPKCS1v15(rand.Reader, pubKey, []byte(s))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(e), nil
}

func (tc *sshTC)ParsePublicKey(raw []byte) (*rsa.PublicKey, error) {
	key := bytes.Split(raw, []byte(" "))[1]
	keydata, err := base64.StdEncoding.DecodeString(string(key))
	if err != nil {
		return nil, err
	}

	parts := [][]byte{}
	for len(keydata) > 0 {
		var dlen uint32
		binary.Read(bytes.NewReader(keydata[:4]), binary.BigEndian, &dlen)

		data := keydata[4 : dlen+4]
		keydata = keydata[4+dlen:]
		parts = append(parts, data)
	}

	n_val := new(big.Int).SetBytes(parts[2])
	e_val := int(new(big.Int).SetBytes(parts[1]).Int64())

	pubKey := &rsa.PublicKey{
		N: n_val,
		E: e_val,
	}
	return pubKey, nil
}
func readPublicKey(path string) (*rsa.PublicKey, error) {
	rawPubKey, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	key := bytes.Split(rawPubKey, []byte(" "))[1]
	keydata, err := base64.StdEncoding.DecodeString(string(key))
	if err != nil {
		return nil, err
	}

	parts := [][]byte{}
	for len(keydata) > 0 {
		var dlen uint32
		binary.Read(bytes.NewReader(keydata[:4]), binary.BigEndian, &dlen)

		data := keydata[4 : dlen+4]
		keydata = keydata[4+dlen:]
		parts = append(parts, data)
	}

	n_val := new(big.Int).SetBytes(parts[2])
	e_val := int(new(big.Int).SetBytes(parts[1]).Int64())

	pubKey := &rsa.PublicKey{
		N: n_val,
		E: e_val,
	}
	return pubKey, nil
}
