
## 準備

```
openssl genrsa 1024 > private_key.pem
openssl rsa -pubout < private_key.pem > public_key.pem
```


## 暗号化

Base64形式にエンコードし、標準出力に出力する

```
⟩ echo "hello, world!" | go run ./encrypt/encrypt_with_rsa_publickey.go ./public_key.pem > encrypted.txt
```

```
⟩ cat encrypted.txt
Tf8O5oFcKDlg72EW1oP2KiEHFDhcgf9IbZ3IgiohYtDDG3M4lyEMPFTlzXBagvb80O+paqyZmGdqw/vd5QtySvn1fZTUOZaGRtCu4oPzz7Gqc86bIDXln5l7Ir50+6UZvagkE4+oRXwI2ybBrzN/5OEEf0gH1XIe/CQSgHfmkWc=⏎
```

## 複合

Base64形式の文字列を標準入力で渡すと、複合して標準出力に出力する

```
⟩ cat encrypted.txt | go run decrypt/decrypt_with_rsa_privatekey.go ./private_key.pem
hello, world!⏎
```

## 参考文献

- [Golang で RSA 署名](https://m0t0k1ch1st0ry.com/blog/2014/08/18/rsa-signing/)
- https://golang.org/pkg/crypto/rsa/