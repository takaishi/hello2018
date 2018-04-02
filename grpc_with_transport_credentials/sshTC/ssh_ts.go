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
	mrand "math/rand"
	"encoding/pem"
	"errors"
	"crypto/x509"
	"crypto/sha256"
	"strings"
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
	//rawKey, err := ioutil.ReadFile("../invalid_rsa")
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

const rs3Letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func (tc *sshTC) randString() string {
	b := make([]byte, 10)
	for i := range b {
		b[i] = rs3Letters[int(mrand.Int63()%int64(len(rs3Letters)))]
	}
	return string(b)
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
	h := sha256.Sum256([]byte(decrypted))
	fmt.Printf("s = %s\n", decrypted)
	fmt.Printf("h = %x\n", h)

	rawConn.Write([]byte(fmt.Sprintf("%x\n", h)))
	return rawConn, nil, err
}

func (tc *sshTC) ServerHandshake(rawConn net.Conn) (_ net.Conn, _ credentials.AuthInfo, err error) {
	fmt.Printf("ServerHandshake\n")

	// 乱数を生成する
	s := tc.randString()
	h := sha256.Sum256([]byte(s))
	fmt.Printf("s = %s\n", s)
	fmt.Printf("h = %x\n", h)


	// 乱数を暗号化してクライアントに送信
	encrypted, err := tc.encrypt(s)
	if err != nil {
		fmt.Printf("Failed to encrypt: %s\n", err)
	}
	//fmt.Printf("encrypted: %s\n", encrypted)
	rawConn.Write([]byte(encrypted))

	// クライアントからハッシュ値を受け取る
	buf := make([]byte, 2014)
	n, err := rawConn.Read(buf)
	if err != nil {
		fmt.Printf("Read error: %s\n", err)
	}
	fmt.Printf("hash: %s\n", buf)

	a := make([]byte, n)
	a = buf[0:n]
	fmt.Println("===============")
	fmt.Printf("a: %#v\n", strings.TrimRight(string(a), "\n"))
	fmt.Printf("b: %#v\n", fmt.Sprintf("%x", h))
	if strings.TrimRight(string(a), "\n") == fmt.Sprintf("%x", h) {
		fmt.Println("Success!!!")
	} else {
		fmt.Println("Baaaaaaaaaaaaaaaad!!")
	}

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