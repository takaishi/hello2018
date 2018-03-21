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
	"io/ioutil"
	"log"
	"os"
)

func readPrivateKey(path string) (*rsa.PrivateKey, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	if block == nil {
		log.Fatal("invalid private key data")
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

func decrypt(s string, key *rsa.PrivateKey) (string, error) {
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

func readFromStdin() string {
	stdin := bufio.NewScanner(os.Stdin)
	stdin.Scan()
	return stdin.Text()
}

func main() {
	priKey, err := readPrivateKey(os.Args[1])
	if err != nil {
		log.Fatalf("Failed to read private key: %s", err)
	}
	s := readFromStdin()
	decrypted, err := decrypt(s, priKey)
	if err != nil {
		log.Fatalf("Failed to decrypt: %s", err)
	}

	fmt.Print(decrypted)

}
