package sshTC

import (
	"golang.org/x/net/context"
	"net"
	"google.golang.org/grpc/credentials"
	"fmt"
	mrand "math/rand"
	"crypto/sha256"
	"strings"
	"os"
)

type sshTC struct {
	info *credentials.ProtocolInfo

}

const rs3Letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func (tc *sshTC) randString() string {
	b := make([]byte, 10)
	for i := range b {
		b[i] = rs3Letters[int(mrand.Int63()%int64(len(rs3Letters)))]
	}
	return string(b)
}

func privateKeyPath() string {
	if os.Getenv("SSH_PRIVATE_KEY_PATH") != "" {
		return os.Getenv("SSH_PRIVATE_KEY_PATH")
	} else {
		return 	fmt.Sprintf("%s/.ssh/id_rsa", os.Getenv("HOME"))
	}
}

func publicKeyPath() string {
	if os.Getenv("SSH_PUBLIC_KEY_PATH") != "" {
		return os.Getenv("SSH_PUBLIC_KEY_PATH")
	} else {
		return 	fmt.Sprintf("%s/.ssh/id_rsa.pub", os.Getenv("HOME"))
	}
}

func (tc *sshTC) ClientHandshake(ctx context.Context, addr string, rawConn net.Conn) (_ net.Conn, _ credentials.AuthInfo, err error) {
	fmt.Printf("ClientHandshake\n")
	buf := make([]byte, 2014)
	n, err := rawConn.Read(buf)
	if err != nil {
		fmt.Printf("Read error: %s\n", err)
	}
	key, err := tc.readPrivateKey(privateKeyPath())

	decrypted, err := tc.Decrypt(string(buf[:n]), key)
	if err != nil {
		fmt.Printf("Failed to decrypt: %s\n", err)
	}
	h := sha256.Sum256([]byte(decrypted))
	fmt.Printf("s = %s\n", decrypted)
	fmt.Printf("h = %x\n", h)

	rawConn.Write([]byte(fmt.Sprintf("%x\n", h)))
	return rawConn, nil, err
}

func (tc *sshTC) ServerHandshake(rawConn net.Conn) (_ net.Conn, _ credentials.AuthInfo, err error) {
	fmt.Printf("ServerHandshake\n")

	// 乱数を生成する
	s := tc.randString()
	h := sha256.Sum256([]byte(s))
	fmt.Printf("s = %s\n", s)
	fmt.Printf("h = %x\n", h)


	// 乱数を暗号化してクライアントに送信
	encrypted, err := tc.Encrypt(s, publicKeyPath())
	if err != nil {
		fmt.Printf("Failed to encrypt: %s\n", err)
	}
	//fmt.Printf("encrypted: %s\n", encrypted)
	rawConn.Write([]byte(encrypted))

	// クライアントからハッシュ値を受け取る
	buf := make([]byte, 2014)
	n, err := rawConn.Read(buf)
	if err != nil {
		fmt.Printf("Read error: %s\n", err)
	}
	fmt.Printf("hash: %s\n", buf)

	buf = buf[:n]
	fmt.Println("===============")
	fmt.Printf("a: %#v\n", strings.TrimRight(string(buf), "\n"))
	fmt.Printf("b: %#v\n", fmt.Sprintf("%x", h))
	if strings.TrimRight(string(buf), "\n") == fmt.Sprintf("%x", h) {
		fmt.Println("Success!!!")
	} else {
		fmt.Println("Baaaaaaaaaaaaaaaad!!")
	}

	return rawConn, nil, err
}

func (tc *sshTC) Info() credentials.ProtocolInfo {
	return *tc.info
}

func (tc *sshTC) Clone() credentials.TransportCredentials {
	info := *tc.info
	return &sshTC{
		info: &info,
	}
}

func (tc *sshTC) OverrideServerName(serverNameOverride string) error {
	return nil
}

func NewServerCreds() credentials.TransportCredentials {
	return &sshTC{
		info: &credentials.ProtocolInfo{
			SecurityProtocol: "ssh",
			SecurityVersion: "1.0",
		},
	}
}

func NewClientCreds() credentials.TransportCredentials {
	return &sshTC{
		info: &credentials.ProtocolInfo{
			SecurityProtocol: "ssh",
			SecurityVersion:  "1.0",
		},
	}
}