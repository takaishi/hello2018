package main

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/rancher/os/log"
	"io/ioutil"
	"os"
)

func readPublicKey(path string) (*rsa.PublicKey, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	if block == nil {
		log.Fatal("invalid public key data")
	}

	if got, want := block.Type, "PUBLIC KEY"; got != want {
		return nil, errors.New(fmt.Sprintf("invalid key type: %s", block.Type))
	}

	keyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	key, ok := keyInterface.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not RSA public key")
	}
	return key, nil
}

func encrypt(s string, key *rsa.PublicKey) (string, error) {
	encrypted, err := rsa.EncryptPKCS1v15(rand.Reader, key, []byte(s))
	if err != nil {
		return "", err
	}

	data := base64.StdEncoding.EncodeToString(encrypted)

	return data, nil
}

func readFromStdin() string {
	stdin := bufio.NewScanner(os.Stdin)
	stdin.Scan()
	return stdin.Text()
}

func main() {
	pubKey, err := readPublicKey(os.Args[1])
	if err != nil {
		log.Fatalf("Failed to read public key: %s", err)
	}

	s := readFromStdin()
	encrypted, err := encrypt(s, pubKey)
	if err != nil {
		log.Fatalf("Failed to encrypt: %s", err)
	}
	fmt.Print(encrypted)
}
