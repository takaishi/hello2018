package sshTC

import (
	"golang.org/x/net/context"
	"net"
	"google.golang.org/grpc/credentials"
	"fmt"
	mrand "math/rand"
	"os"
	"errors"
	"log"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"github.com/STNS/STNS/stns"
	"crypto/sha256"
	"strings"
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

func getUsername() (string, error) {
	if os.Getenv("SSH_USER") == "" {
		return "", errors.New("require SSH_USER")
	}
	return os.Getenv("SSH_USER"), nil
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
	username, err := getUsername()
	log.Printf("[DEBUG] username: %s\n", string(username))

	if err != nil {
		return nil, nil, err
	}
	rawConn.Write([]byte(username))


	buf := make([]byte, 2014)
	n, err := rawConn.Read(buf)
	if err != nil {
		fmt.Printf("Read error: %s\n", err)
		return nil, nil, err
	}
	log.Printf("[DEBUG] privateKeyPath: %s\n", privateKeyPath())
	log.Printf("[DEBUG] buf: %s\n", string(buf[:n]))
	key, err := tc.readPrivateKey(privateKeyPath())
	if err != nil {
		fmt.Printf("Failed to read private key: %s\n", err)
		return nil, nil, err
	}

	decrypted, err := tc.Decrypt(string(buf[:n]), key)
	if err != nil {
		fmt.Printf("Failed to decrypt: %s\n", err)
		return nil, nil, err
	}
	h := sha256.Sum256([]byte(decrypted))

	rawConn.Write([]byte(fmt.Sprintf("%x\n", h)))

	r := make([]byte, 64)
	n, err = rawConn.Read(r)
	if err != nil {
		fmt.Printf("Read error: %s\n", err)
		return nil, nil, err
	}
	r = r[:n]
	if string(r) != "ok" {
		fmt.Println("Failed to authenticate")
		return nil, nil, errors.New("Failed to authenticate")
	}

	return rawConn, nil, err
}

type v2Metadata struct {
	APIVersion float64 `json:"api_version"`
	Result string `json:"result"`

}
type UserResponse struct {
	Metadata v2Metadata `json:"metadata"`
	Items stns.Attributes `json:"items"`
}

func (tc *sshTC) getPubKeyFromSTNS(name string) ([]byte, error) {
	var user_resp UserResponse

	resp, err := http.Get(fmt.Sprintf("http://localhost:1104/v2/user/name/%s", name))
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &user_resp)
	return []byte(user_resp.Items[name].User.Keys[0]), nil
}

func (tc *sshTC) ServerHandshake(rawConn net.Conn) (_ net.Conn, _ credentials.AuthInfo, err error) {
	// 乱数を生成する
	s := tc.randString()
	h := sha256.Sum256([]byte(s))

	// ユーザー名を読み込み&STNSからPublicKeyを取得
	buf := make([]byte, 2014)
	n, err := rawConn.Read(buf)
	if err != nil {
		return nil, nil, errors.New(fmt.Sprintf("Read error: %s\n", err))
	}
	rawKey, err := tc.getPubKeyFromSTNS(string(buf[:n]))
	if err != nil {
		return nil, nil, err
	}
	log.Printf("[DEBUG] rawKey = %s\n", string(rawKey))

	pubKey, err := tc.ParsePublicKey(rawKey)
	if err != nil {
		return nil, nil, errors.New(fmt.Sprintf("Failed to parse: %s\n", err))
	}
	encrypted, err := tc.Encrypt(s, pubKey)
	if err != nil {
		return nil, nil, errors.New(fmt.Sprintf("Failed to encrypt: %s\n", err))
	}
	rawConn.Write([]byte(encrypted))

	buf = make([]byte, 2014)
	n, err = rawConn.Read(buf)
	if err != nil {
		return nil, nil, errors.New(fmt.Sprintf("Read error: %s\n", err))
	}
	buf = buf[:n]
	if strings.TrimRight(string(buf), "\n") == fmt.Sprintf("%x", h) {
		rawConn.Write([]byte("ok"))
		fmt.Println("Success!!!")
	} else {
		rawConn.Write([]byte("ng"))
		fmt.Println("Failed!!!")
		return nil, nil, errors.New(fmt.Sprintf("Failed to authenticate: invalid key"))
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