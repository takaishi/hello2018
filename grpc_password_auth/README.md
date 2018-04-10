# gRPC with user/password authentication sample

## Usage

server:

```
$ go run ./main.go --username admin -p admin123 server
```

client:

```
$ go run ./main.go -u admin -p admin123 client add foo 1

$ go run ./main.go -u admin -p admin123 client list
name:"foo" age:1
```

client incorrect password:

```
$ go run ./main.go -u admin -p hoge client list
2018/04/10 09:39:16 rpc error: code = Unauthenticated desc = AccessDeniedError
exit status 1
```

## Reference

- https://qiita.com/Mamoru-Izuka/items/28724d9dd8a6b30b236d
- https://mattn.kaoriya.net/software/lang/go/20150227144125.htm
- https://github.com/grpc/grpc-go/issues/106#issuecomment-246978683
- https://github.com/buckhx/safari-zone