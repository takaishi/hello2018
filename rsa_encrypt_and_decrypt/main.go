package rsa_encrypt_and_decrypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/rancher/os/log"
	"io/ioutil"
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

func main() {
	fmt.Println(">>>> Encrypt Phase")
	pubKey, err := readPublicKey("./public_key.pem")
	if err != nil {
		log.Fatalf("Failed to read public key: %s", err)
	}

	s := "hello, world!"
	encrypted, err := encrypt(s, pubKey)
	if err != nil {
		log.Fatalf("Failed to encrypt: %s", err)
	}
	//data := []byte{}
	fmt.Printf("input: %s\n", s)
	fmt.Printf("encrypted(base64 encoded): %s\n", encrypted)

	fmt.Println("")
	fmt.Println(">>>> Decrypt Phase")

	priKey, err := readPrivateKey("./private_key.pem")
	if err != nil {
		log.Fatalf("Failed to read private key: %s", err)
	}
	decrypted, err := decrypt(encrypted, priKey)
	if err != nil {
		log.Fatalf("Failed to decrypt: %s", err)
	}

	fmt.Printf("input(base64 encoded): %s\n", encrypted)
	fmt.Printf("decrypted: %s\n", string(decrypted))

}
