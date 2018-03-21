package main

import (
  "io/ioutil"
  "log"
  "fmt"
  "bytes"
  "encoding/base64"
  "encoding/binary"
  "strconv"
  "math/big"
  "crypto/rsa"
  "crypto/rand"
  "bufio"
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
    //fmt.Printf("keydata[:4] = %#v\n", keydata[:4])
    //fmt.Printf("dlen = %d\n", dlen)

    data := keydata[4:dlen+4]
    keydata = keydata[4+dlen:]
    parts = append(parts, data)
    //fmt.Printf("parts = %#v\n", parts)
  }

  // parts[1]を処理してn_valを作成
  //fmt.Printf("parts[1] = %#v\n", parts[1])
  e_val_str := ""
  for _, x := range parts[1] {
    e_val_str = e_val_str + fmt.Sprintf("%02x", byte(x))

  }
  e_val, err := strconv.ParseInt(e_val_str, 16, 64)
  if err != nil {
    log.Fatal(err)
  }
  //fmt.Printf("e_val = %d\n", e_val)


  // parts[2]を処理してe_valを作成
  pubKey := &rsa.PublicKey{
    N: new(big.Int).SetBytes(parts[2]),
    E: int(e_val),
  }
  //fmt.Printf("pubKey = %#v\n", pubKey)

  stdin := bufio.NewScanner(os.Stdin)
  stdin.Scan()
  s := stdin.Text()
  e, err := rsa.EncryptPKCS1v15(rand.Reader, pubKey, []byte(s))
  if err != nil {
    log.Fatalln(err)
  }
  d := base64.StdEncoding.EncodeToString(e)
  fmt.Println(d)
    //fmt.Printf("parts[1] = %#v\n", parts[2])
  //e_val_str := ""
  //for _, x := range parts[2] {
  //  e_val_str = e_val_str + fmt.Sprintf("%02x", byte(x))
  //
  //}
  //fmt.Printf("e_val_str = %s\n", e_val_str)
  //e_val, err := strconv.ParseInt(e_val_str, 16, 64)
  //if err != nil {
  //  log.Fatal(err)
  //}
  //fmt.Printf("e_val = %d\n", e_val)


  //fmt.Printf("data = %s\n", string(data))

  //binary.Read(bytes.NewReader(keydata[:4]), binary.BigEndian, &dlen)
  //fmt.Printf("keydata[:4] = %#v\n", keydata[:4])
  //fmt.Printf("dlen = %d\n", dlen)


    //data = keydata[4:dlen+4]
    //keydata = keydata[4+dlen:]
    //fmt.Printf("%+v\n", dlen)
    //fmt.Printf("%s\n", string(data))
  //}

}
