package conn

import "net"

type conn struct {
	net.Conn
	crypt HelloRecordCrypt
}

func (p *conn) Read(b []byte) (int, error) {
	return p.Conn.Read(b)
}

func (p *conn) Write(b []byte) (int, error) {
	return p.Conn.Write(b)
}

func NewConn(c net.Conn) (net.Conn, error) {
	crypt := &HelloRecordCrypt{}
	helloConn := &conn{
		Conn:  c,
		crypt: *crypt,
	}

	return helloConn, nil
}

