package main

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"
)

func main() {
	rawData, err := ioutil.ReadFile("./test_rsa.pub")
	if err != nil {
		log.Fatal(err)
	}
	key := bytes.Split(rawData, []byte(" "))[1]
	keydata, err := base64.StdEncoding.DecodeString(string(key))
	if err != nil {
		log.Fatal(err)
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

	stdin := bufio.NewScanner(os.Stdin)
	stdin.Scan()
	s := stdin.Text()
	e, err := rsa.EncryptPKCS1v15(rand.Reader, pubKey, []byte(s))
	if err != nil {
		log.Fatalln(err)
	}
	d := base64.StdEncoding.EncodeToString(e)
	fmt.Println(d)
}
