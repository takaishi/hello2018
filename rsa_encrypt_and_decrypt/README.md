
## 準備

```
openssl genrsa 1024 > private_key.pem
openssl rsa -pubout < private_key.pem > public_key.pem
```


## 暗号化

Base64形式にエンコードし、標準出力に出力する

```
⟩ echo "hello, world!" | go run ./encrypt/encrypt_with_rsa_publickey.go ./public_key.pem
MIVSK2ahCa59p/2SEy2PQBQdop6o0ElMh0hq5xspsb+9ZsHvRtL1c/2lb1xP7bep8nJMJG5cBILiGl6zbjd4lBB3noLn4UDlO+ASMMSnihZzYNw9Dka6Le/xLEaEofu9QXybxBO5Si+VwSx/UsZGOHY1sB4ZFmWNx3iNtir0ONQ=⏎
```

## 参考文献

- [Golang で RSA 署名
](https://m0t0k1ch1st0ry.com/blog/2014/08/18/rsa-signing/)
- https://golang.org/pkg/crypto/rsa/