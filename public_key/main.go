package main

import (
	"os"
	"bufio"
	"fmt"
	"strings"
)

func readPublicKey(path string) ([]byte, error) {
	s := ""
	fp, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		text := scanner.Text()
		s = s + text
	}
	return []byte(s), nil
}
func main() {
	plainText := "Bob loves Alice."
	fmt.Printf("plainText = %s\n", plainText)

	sshPublicKey := ""
	fmt.Println(strings.Split(sshPublicKey, " "))
	//publicKey, _, _, _, err := ssh.ParseAuthorizedKey(sshPublicKeyjjjj)
	//if err != nil {
	//	log.Fatalf("failed to parse openssh public key: %v", err)
	//}
	//m := ssh.MarshalAuthorizedKey(publicKey)
	//fmt.Printf("%s", string(m))

	//publicKeyBytes, err := base64.StdEncoding.DecodeString(base64EncodedPublicKey)
	//publicKeyBytes = append([]byte("---- BEGIN SSH2 PUBLIC KEY ----\n"), publicKeyBytes)
	//if err != nil {
	//	log.Fatalf("failed to decode: %v", err)
	//}
	//block, _ := pem.Decode(m)
	//if block == nil {
	//	log.Fatalln("invalid public key data")
	//}
	//if block.Type != "PUBLIC KEY" {
	//	log.Fatalf("invalid public key type : %s", block.Type)
	//}
	//fmt.Printf("%v\n", block)
	//keyInterface, err := x509.ParsePKIXPublicKey(m)
	//if err != nil {
	//	log.Fatalf("failed to parse public key string: %v", err)
	//}
	//rsaPublicKey, ok := keyInterface.(*rsa.PublicKey)
	//if !ok {
	//	log.Fatalf("failed to create rsa public key string: %v", ok)
	//}
	//
	//rsa.EncryptPKCS1v15(rand.Reader, rsaPublicKey, []byte(plainText))
}
