

```
⟩ ssh-keygen
Generating public/private rsa key pair.
Enter file in which to save the key (/Users/r_takaishi/.ssh/id_rsa): ./test_rsa
Enter passphrase (empty for no passphrase):
Enter same passphrase again:
Your identification has been saved in ./test_rsa.
Your public key has been saved in ./test_rsa.pub.
The key fingerprint is:
SHA256:O14buTjMvedxu8y4tGybrA2QGc0B8GJq4W4exmCfIVI r_takaishi@takaishiryou-no-iMac.local
The key's randomart image is:
+---[RSA 2048]----+
|      .....      |
|       . o .     |
|  E . o o o      |
| . . + . +       |
|. + =   S        |
| o B o   o .     |
|    O  oo.= o .  |
|   + . .++.XoO . |
|    .   o.=*@o=. |
+----[SHA256]-----+

~/src/github.com/takaishi/hello2018/ssh_key · (rsa_encrypt_and_decrypt±)
⟩ ls -l test_rsa*
-rw-------  1 r_takaishi  staff  1679  3 21 10:32 test_rsa
-rw-r--r--  1 r_takaishi  staff   419  3 21 10:32 test_rsa.pub
```


```
⟩ cat test_rsa
-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEA8g7SjDm5MvvCK/ML7k+x6NhAcnKidcUSBwsbf/8Iv+eW7M8u
MjAnKhctiW/qObpkIw+bqAIRF3nncm14pr3/NZV3LMtYkU7qPN1HY3+HjRxe/v6D
aTBITcznATqIHzpJkpU8QZOMJMNbI2zrH5b9ggtB7S7u3PamQgU0RoD2AQHfFTCc
Ij5nupJn47DpCOyRMTAFRfMkDD8EW3MyZvHsd5W20meinXf5b7XpVRUpDhi+/N5s
RKyXl5YCjJLX0UNwKDFVhvnCxqiYAh1362xXI6LFuB+phxoJEBT/yZGPAqs5dyBL
/mOAE+ELyk1iFkIqqll5xV9gP6wmln9HqlOlSwIDAQABAoIBAE9n02QSz5FNC36V
ZZWQ6UEEJ+gjeO3/bxGGcEgF5t3lYBphQLtQFpj1L4gFgaXcYlsqFJsByo+T+vwL
s2enrl/qn0S/lFdetvKueGvIezQsWXF3Fq7cGuwCyskZZWwxF8+RS0oL2A57U5uE
cIFVa+ZMQR1Ipy0vcIz53hM+3PSpAC24Hi6jnESaTCmZBoQWSKv8ZrM1BoDaCz12
bsBAfixWn+jIMxD+VlQJEZb+kIwIXSiYu4ZPToqYQYEHrNHpSLi6s8T9j9uTPzFk
fNa6zaxmaDRTK9ZLGeOeOXkcHCYppUc10CML0lwogMUAYxYkel1jb8z+i4t+reMF
mOeHvVkCgYEA+Y5DehLzbucRY/1HoTdUQTEc4C8WG+d5+njLnL1Mc49Qgf0hsb1M
Whc9/2BW1hxMZ+tMP4J+ffoFEpX9f6+zXH+eBZxOxLGRTCC5GeFvAevvERcZyjL/
SKqxdehPDQT0G+MNuSVZfom3f8bfR1jU2VY/drW9tgjXhBeVRYq6ALcCgYEA+E7+
O8JTNo1ivLfTpi9ABWJuHDkNnjQr6BwNhb7+YXms9CiQqZhKr3O4U4Dt0RLqAuso
meoSgxf0RNAYmk5hqmVREIxna0jjNCcsisXlCWZg58NcGtMuscnyUhnfh4BCjzOn
/DUSqG+WGMpFp01So4cdp0XOgJK02VhsVkvXRA0CgYBnbdj1jUkrW8VPZbf8T0wy
QMKw/5LwOb3KW6o36hT3iBxb46fFXKl6ZUuivjD/SHc6UsElSVZXq/nSPCv2ccGq
wpGhzaivyNBpdt6ApXg2maxZrvNXZE99tJEcRw4MXVM3A6G6bIps8XMGGEyN63k8
IoznDGf2PC/mZUfOrLJufQKBgQDVW0qoHnlRznqgnXOGv+LKvaDPL1a8MSfo8PHN
kicRqnMp+BEVKH5D87LWTVoK462fhGAGoFH3woVo+0WokODqgNP+3CWg0agoD+D9
/LyoLwflHL/vbLYaneNRGFoxG8wVL9WPqCq3/+mAs4zWDGKNkHOyXxDo+SXb+1Zb
cB8voQKBgQDPSrWPVVJ5Ht+Tjc3ZzxCE+XOTMkIV3boAxqAekqRB9qTh0Zif6/sD
SCglFlt6kyNXP17todX+/rjKCUl4ZQTxSQzeq/9zZs4YimcFNI6AU9Ld1Ga28AF5
EUJzAV4JBM/Baewq7Jo1URCPFrZ4mE0AtRg6gOJN26Vc4yBicTM5bw==
-----END RSA PRIVATE KEY-----

```

```
⟩ cat test_rsa.pub
ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDyDtKMObky+8Ir8wvuT7Ho2EBycqJ1xRIHCxt//wi/55bszy4yMCcqFy2Jb+o5umQjD5uoAhEXeedybXimvf81lXcsy1iRTuo83Udjf4eNHF7+/oNpMEhNzOcBOogfOkmSlTxBk4wkw1sjbOsflv2CC0HtLu7c9qZCBTRGgPYBAd8VMJwiPme6kmfjsOkI7JExMAVF8yQMPwRbczJm8ex3lbbSZ6Kdd/lvtelVFSkOGL783mxErJeXlgKMktfRQ3AoMVWG+cLGqJgCHXfrbFcjosW4H6mHGgkQFP/JkY8Cqzl3IEv+Y4AT4QvKTWIWQiqqWXnFX2A/rCaWf0eqU6VL r_takaishi@takaishiryou-no-iMac.local
```

```
     -m key_format
             Specify a key format for the -i (import) or -e (export) conversion options.  The supported key formats are: ``RFC4716'' (RFC 4716/SSH2 public or private key),
             ``PKCS8'' (PEM PKCS8 public key) or ``PEM'' (PEM public key).  The default conversion format is ``RFC4716''.
```

RFC4716  (RFC 4716/SSH2 public or private key),
PKCS8 (PEM PKCS8 public key)
PEM'' (PEM public key)
```


pemを指定した場合。


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


PKCS8を指定した場合。`BEGIN PUBLIC KEY`となる。`x509.ParsePKIXPublicKey`でパース可能。
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


# 参考文献

- https://ja.wikipedia.org/wiki/PKCS
https://gist.github.com/thwarted/1024558