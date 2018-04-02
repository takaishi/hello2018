package sshTC

import (
	"golang.org/x/net/context"
	"net"
	"google.golang.org/grpc/credentials"
	"fmt"
	"io/ioutil"
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"math/big"
	"crypto/rsa"
	"log"
	"crypto/rand"
	"encoding/pem"
	"errors"
	"crypto/x509"
)

type sshTC struct {
	info *credentials.ProtocolInfo

}

func (tc *sshTC) encrypt(s string) (string, error) {
	rawPubKey, err := ioutil.ReadFile("../test_rsa.pub")
	//rawPubKey, err := ioutil.ReadFile("../invalid_rsa.pub")
	if err != nil {
		return "", err
	}
	key := bytes.Split(rawPubKey, []byte(" "))[1]
	keydata, err := base64.StdEncoding.DecodeString(string(key))
	if err != nil {
		return "", err
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

	e, err := rsa.EncryptPKCS1v15(rand.Reader, pubKey, []byte(s))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(e), nil
}

func (tc *sshTC) decrypt(s string) (string, error){
	rawKey, err := ioutil.ReadFile("../test_rsa")
	if err != nil {
		return "", err
	}

	block, _ := pem.Decode(rawKey)
	if block == nil {
		log.Fatal("invalid private key data")
	}

	if got, want := block.Type, "RSA PRIVATE KEY"; got != want {
		return "", errors.New(fmt.Sprintf("invalid key type: %s", block.Type))
	}

	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}

	in, err := base64.StdEncoding.DecodeString(s)
	decrypted, err := rsa.DecryptPKCS1v15(rand.Reader, key, in)

	return string(decrypted), nil
}

func (tc *sshTC) ClientHandshake(ctx context.Context, addr string, rawConn net.Conn) (_ net.Conn, _ credentials.AuthInfo, err error) {
	fmt.Printf("ClientHandshake\n")
	buf := make([]byte, 2014)
	_, err = rawConn.Read(buf)
	if err != nil {
		fmt.Printf("Read error: %s\n", err)
	}
	decrypted, err := tc.decrypt(string(buf))
	if err != nil {
		fmt.Printf("Failed to decrypt: %s\n", err)
	}
	fmt.Printf("%s\n", decrypted)
	return rawConn, nil, err
}

func (tc *sshTC) ServerHandshake(rawConn net.Conn) (_ net.Conn, _ credentials.AuthInfo, err error) {
	fmt.Printf("ServerHandshake\n")
	encrypted, err := tc.encrypt("Hello!")
	if err != nil {
		fmt.Printf("Failed to encrypt: %s\n", err)
	}
	fmt.Printf("encrypted: %s\n", encrypted)
	rawConn.Write([]byte(encrypted))
	return rawConn, nil, err
}

func (tc *sshTC) Info() credentials.ProtocolInfo {
	return *tc.info
}

func (tc *sshTC) Clone() credentials.TransportCredentials {
	info := *tc.info
	return &sshTC{
		info: &info,
	}
}

func (tc *sshTC) OverrideServerName(serverNameOverride string) error {
	return nil
}

func NewServerCreds() credentials.TransportCredentials {
	return &sshTC{
		info: &credentials.ProtocolInfo{
			SecurityProtocol: "ssh",
			SecurityVersion: "1.0",
		},
	}
}

func NewClientCreds() credentials.TransportCredentials {
	return &sshTC{
		info: &credentials.ProtocolInfo{
			SecurityProtocol: "ssh",
			SecurityVersion:  "1.0",
		},
	}
}