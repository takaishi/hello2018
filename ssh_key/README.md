# ssh-keygenで生成した公開鍵を扱う

`ssh-keygen` でフォーマットを指定しない場合、公開鍵として以下のようなファイルが生成される

```
⟩ cat test_rsa.pub
ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDyDtKMObky+8Ir8wvuT7Ho2EBycqJ1xRIHCxt//wi/55bszy4yMCcqFy2Jb+o5umQjD5uoAhEXeedybXimvf81lXcsy1iRTuo83Udjf4eNHF7+/oNpMEhNzOcBOogfOkmSlTxBk4wkw1sjbOsflv2CC0HtLu7c9qZCBTRGgPYBAd8VMJwiPme6kmfjsOkI7JExMAVF8yQMPwRbczJm8ex3lbbSZ6Kdd/lvtelVFSkOGL783mxErJeXlgKMktfRQ3AoMVWG+cLGqJgCHXfrbFcjosW4H6mHGgkQFP/JkY8Cqzl3IEv+Y4AT4QvKTWIWQiqqWXnFX2A/rCaWf0eqU6VL r_takaishi@takaishiryou-no-iMac.local
```

これはRFC4716の形式。`-m`オプションで指定することで、PKCS8やPEMといったフォーマットを使うこともできる。

```
     -m key_format
             Specify a key format for the -i (import) or -e (export) conversion options.  The supported key formats are: ``RFC4716'' (RFC 4716/SSH2 public or private key),
             ``PKCS8'' (PEM PKCS8 public key) or ``PEM'' (PEM public key).  The default conversion format is ``RFC4716''.
```

RFC4716  (RFC 4716/SSH2 public or private key),
PKCS8 (PEM PKCS8 public key)
PEM'' (PEM public key)


RFC4716形式からPEM形式に変換できる。

```
⟩ ssh-keygen -f ./test_rsa.pub -e -m pem > test_rsa.pem

⟩ cat test_rsa.pem
-----BEGIN RSA PUBLIC KEY-----
MIIBCgKCAQEA8g7SjDm5MvvCK/ML7k+x6NhAcnKidcUSBwsbf/8Iv+eW7M8uMjAn
KhctiW/qObpkIw+bqAIRF3nncm14pr3/NZV3LMtYkU7qPN1HY3+HjRxe/v6DaTBI
TcznATqIHzpJkpU8QZOMJMNbI2zrH5b9ggtB7S7u3PamQgU0RoD2AQHfFTCcIj5n
upJn47DpCOyRMTAFRfMkDD8EW3MyZvHsd5W20meinXf5b7XpVRUpDhi+/N5sRKyX
l5YCjJLX0UNwKDFVhvnCxqiYAh1362xXI6LFuB+phxoJEBT/yZGPAqs5dyBL/mOA
E+ELyk1iFkIqqll5xV9gP6wmln9HqlOlSwIDAQAB
-----END RSA PUBLIC KEY-----

```

PKCS8に変換もできる。この形式だと、`x509.ParsePKIXPublicKey`でパース可能。

```
⟩ ssh-keygen -f ./test_rsa.pub -e -m pkcs8 > test_rsa.pkcs8

⟩ cat test_rsa.pkcs8
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA8g7SjDm5MvvCK/ML7k+x
6NhAcnKidcUSBwsbf/8Iv+eW7M8uMjAnKhctiW/qObpkIw+bqAIRF3nncm14pr3/
NZV3LMtYkU7qPN1HY3+HjRxe/v6DaTBITcznATqIHzpJkpU8QZOMJMNbI2zrH5b9
ggtB7S7u3PamQgU0RoD2AQHfFTCcIj5nupJn47DpCOyRMTAFRfMkDD8EW3MyZvHs
d5W20meinXf5b7XpVRUpDhi+/N5sRKyXl5YCjJLX0UNwKDFVhvnCxqiYAh1362xX
I6LFuB+phxoJEBT/yZGPAqs5dyBL/mOAE+ELyk1iFkIqqll5xV9gP6wmln9HqlOl
SwIDAQAB
-----END PUBLIC KEY-----

```

## RFC4716形式の公開鍵を用いて暗号化する


```
⟩ echo "hello, world!" | go run encrypt_with_ssh_publickey.go  | go run ../rsa_encrypt_and_decrypt/decrypt/decrypt_with_rsa_privatekey.go ./test_rsa
hello, world!⏎
```



# 参考文献

- http://blog.oddbit.com/2011/05/08/converting-openssh-public-keys/
- https://ja.wikipedia.org/wiki/PKCS
- https://gist.github.com/thwarted/1024558print
- https://github.com/yosida95/golang-sshkey
- https://github.com/goken/goken/blob/master/goken09-binary/goken09-binary.md
- https://yosida95.com/2015/05/31/121709.html
- https://tools.ietf.org/html/rfc4716