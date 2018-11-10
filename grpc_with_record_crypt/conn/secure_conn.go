package conn

import (
	"encoding/binary"
	"net"
)

const (
	MsgLenFieldSize = 4
)

type secureConn struct {
	net.Conn
	crypt    HelloRecordCrypt
	overhead int
}

func (p *secureConn) Read(buf []byte) (int, error) {
	var msgSize uint32
	var decrypted []byte
	var msgSizeBuf []byte
	var msgBuf []byte

	msgSizeBuf = make([]byte, p.overhead)
	_, err := p.Conn.Read(msgSizeBuf)
	if err != nil {
		return 0, err
	}
	msgSize = binary.LittleEndian.Uint32(msgSizeBuf)

	msgBuf = make([]byte, msgSize)
	_, err = p.Conn.Read(msgBuf)
	if err != nil {
		return 0, err
	}

	if msgSize != 0 {
		decrypted, err = p.crypt.Decrypt(decrypted, msgBuf)
		if err != nil {
			return 0, err
		}
	}
	n := copy(buf, decrypted)
	return n, nil
}

func (p *secureConn) Write(rawBuf []byte) (int, error) {
	var buf []byte

	buf, err := p.crypt.Encrypt(buf, rawBuf)
	if err != nil {
		return 0, err
	}
	msg := make([]byte, len(buf)+p.overhead)

	copy(msg[4:], buf)
	msgSize := uint32(len(msg) - p.overhead)
	binary.LittleEndian.PutUint32(msg, msgSize)
	_, err = p.Conn.Write(msg)
	if err != nil {
		return 0, err
	}
	return len(rawBuf), nil
}

func NewSecureConn(c net.Conn) (net.Conn, error) {
	crypt := &HelloRecordCrypt{}
	overhead := MsgLenFieldSize
	helloConn := &secureConn{
		Conn:     c,
		crypt:    *crypt,
		overhead: overhead,
	}

	return helloConn, nil
}
