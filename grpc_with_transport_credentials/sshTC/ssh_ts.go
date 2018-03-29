package sshTC

import (
	"golang.org/x/net/context"
	"net"
	"google.golang.org/grpc/credentials"
)

type sshTC struct {
	info *credentials.ProtocolInfo

}

func (tc *sshTC) ClientHandshake(ctx context.Context, addr string, rawConn net.Conn) (_ net.Conn, _ credentials.AuthInfo, err error) {
	return nil, nil, err
}

func (tc *sshTC) ServerHandshake(rawConn net.Conn) (_ net.Conn, _ credentials.AuthInfo, err error) {
	return nil, nil, err
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
