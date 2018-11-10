package conn

import (
	"net"
)

type conn struct {
	net.Conn
}

func (p *conn) Read(b []byte) (int, error) {
	return p.Conn.Read(b)
}

func (p *conn) Write(b []byte) (int, error) {
	return p.Conn.Write(b)
}

func NewConn(c net.Conn) (net.Conn, error) {
	conn := &conn{
		Conn: c,
	}
	return conn, nil
}
