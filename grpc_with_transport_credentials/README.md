# gRPC with SSH Public key authentication sample

## Usage

server:

```
$ go run ./main.go server -p ./test_rsa.pub
```

client:

```
$ go run ./main.go client -i ./test_rsa add foo 1

$ go run ./main.go client -i ./test_rsa list
name:"foo" age:1
```

client with incorrect private key:

```
$ go run ./main.go client -i ./invalid_rsa list
2018/04/10 21:42:47 [ERROR] Failed to decrypt: crypto/rsa: decryption error
2018/04/10 21:42:47 rpc error: code = Unavailable desc = all SubConns are in TransientFailure, latest connection error: connection error: desc = "transport: authentication handshake failed: crypto/rsa: decryption error"
exit status 1

```

## Reference

- [grpc-goのInterceptorを使ってみる](https://qiita.com/Mamoru-Izuka/items/28724d9dd8a6b30b236d)
- https://mattn.kaoriya.net/software/lang/go/20150227144125.htm
- https://github.com/grpc/grpc-go/issues/106#issuecomment-246978683
- https://github.com/buckhx/safari-zone
- https://blog.gopheracademy.com/advent-2017/go-grpc-beyond-basics/
- [How to achieve ssh like auth based on SSL public keys?](https://github.com/grpc/grpc-go/issues/1252)
- [TransportCredentials](https://github.com/grpc/grpc-go/blob/v1.11.0/credentials/credentials.go#L82)
- [alts](https://github.com/grpc/grpc-go/blob/v1.11.0/credentials/alts/alts.go#L102)
- http://nigohiroki.hatenablog.com/entry/2013/08/18/221434
- https://qiita.com/srtkkou/items/ccbddc881d6f3549baf1