package tc

import (
	"github.com/takaishi/hello2018/grpc_with_record_crypt/conn"
	"golang.org/x/net/context"
	"google.golang.org/grpc/credentials"
	"net"
)

type TC struct {
	info   *credentials.ProtocolInfo
	secure bool
}

func (tc *TC) ClientHandshake(ctx context.Context, addr string, rawConn net.Conn) (_ net.Conn, _ credentials.AuthInfo, err error) {
	var c net.Conn
	if tc.secure {
		c, err = conn.NewSecureConn(rawConn)
	} else {
		c, err = conn.NewConn(rawConn)
	}
	if err != nil {
		return nil, nil, err
	}
	return c, nil, err
}

func (tc *TC) ServerHandshake(rawConn net.Conn) (_ net.Conn, _ credentials.AuthInfo, err error) {
	conn, err := conn.NewConn(rawConn)
	if err != nil {
		return nil, nil, err
	}
	return conn, nil, err
}

func (tc *TC) Info() credentials.ProtocolInfo {
	return *tc.info
}

func (tc *TC) Clone() credentials.TransportCredentials {
	info := *tc.info
	return &TC{
		info: &info,
	}
}

func (tc *TC) OverrideServerName(serverNameOverride string) error {
	return nil
}

func NewServerCreds(secure bool) credentials.TransportCredentials {
	return &TC{
		info: &credentials.ProtocolInfo{
			SecurityProtocol: "test",
			SecurityVersion:  "1.0",
		},
		secure: secure,
	}
}

func NewClientCreds(secure bool) credentials.TransportCredentials {
	return &TC{
		info: &credentials.ProtocolInfo{
			SecurityProtocol: "test",
			SecurityVersion:  "1.0",
		},
		secure: secure,
	}
}
