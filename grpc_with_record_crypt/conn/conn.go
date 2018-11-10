package conn

import (
	"net"
)

const (
	MsgLenFieldSize = 4
)

type conn struct {
	net.Conn
	crypt    HelloRecordCrypt
	overhead int
}

func (p *conn) Read(b []byte) (int, error) {
	//log.Printf("[DEBUG] Read: START------------\n")
	//log.Printf("[DEBUG] len(b) = %d\n", len(b))
	//i, err := p.Conn.Read(b)
	//if err != nil {
	//	return 0, err
	//}
	//log.Printf("[DEBUG] b: %s\n", string(b))
	//log.Printf("[DEBUG] read length = %d\n", i)
	//log.Printf("[DEBUG] Read: END------------\n")
	//return i, nil

	//var msgSize uint32
	//var buf []byte
	//var decrypted []byte

	//log.Printf("[DEBUG] Read: START------------\n")
	//log.Printf("[DEBUG] len(b) = %d\n", len(b))
	//buf = make([]byte, p.overhead)
	//_, err := p.Conn.Read(buf)
	//if err != nil {
	//	log.Printf("[ERORR] failed to read %s\n", err.Error())
	//	return 0, err
	//}
	//msgSize = binary.LittleEndian.Uint32(buf)
	//log.Printf("[DEBUG] Read: msgSize = %d\n", msgSize)

	//buf = make([]byte, msgSize)
	//_, err = p.Conn.Read(buf)
	//if err != nil {
	//	log.Printf("[ERORR] failed to read msgSize%s\n", err.Error())
	//	return 0, err
	//}
	//log.Printf("[DEBUG] Read: buf = %+v\n", buf)
	//log.Printf("[DEBUG] Read: len(buf) = %+v\n", len(buf))
	//log.Printf("[DEBUG] Read: size = %d\n", s)

	//if msgSize != 0 {
	//	decrypted, err = p.crypt.Decrypt(b, buf)
	//	if err != nil {
	//		log.Printf("[ERORR] %s\n", err.Error())
	//		return 0, err
	//	}

	//log.Printf("[DEBUG] Read: decrypted = %+v\n", decrypted)
	//b = make([]byte, len(decrypted))
	//} else {
	//log.Printf("[DEBUG] msgSize == 0")
	//b = make([]byte, 0)
	//}
	//n := copy(b, decrypted)
	//log.Printf("[DEBUG] Read: b = %s\n", b)
	//log.Printf("[DEBUG] Read: b = %+v\n", b)
	//log.Printf("[DEBUG] Read: END------------\n")
	//log.Printf("[DEBUG] Read: len(b) = %+v\n", n)
	//return n, nil
	return p.Conn.Read(b)
}

func (p *conn) Write(raw []byte) (int, error) {
	//log.Printf("[DEBUG] Write: START------------\n")
	//log.Printf("[DEBUG] raw = %s\n", string(raw))
	//log.Printf("[DEBUG] len(raw) = %d\n", len(raw))
	//i, err := p.Conn.Write(raw)
	//if err != nil {
	//	return 0, err
	//}
	//log.Printf("[DEBUG] write length = %d\n", i)
	//log.Printf("[DEBUG] Write: END------------\n")
	//return i, nil
	//var tmp []byte
	//
	//log.Printf("[DEBUG] Write: START------------\n")
	//log.Printf("[DEBUG] Write: raw = %s\n", raw)
	//log.Printf("[DEBUG] Write: raw = %+v\n", raw)
	//log.Printf("[DEBUG] Write: len(raw) = %+v\n", len(raw))

	//tmp, err := p.crypt.Encrypt(tmp, raw)
	//if err != nil {
	//	log.Printf("[ERORR] failed to encrypt: %s\n", err.Error())
	//	return 0, err
	//}
	//msg := make([]byte, len(tmp)+p.overhead)
	//
	//copy(msg[4:], tmp)
	//log.Printf("[DEBUG] Write: 1: encrypted = %+v\n", msg)
	//log.Printf("[DEBUG] Write: len(encrypted) = %+v\n", len(msg))

	//msgSize := uint32(len(msg) - p.overhead)
	//log.Printf("[DEBUG] msgSize = %d\n", msgSize)
	//binary.LittleEndian.PutUint32(msg, msgSize)
	//log.Printf("[DEBUG] Write: 2: encrypted = %+v\n", msg)
	//log.Printf("[DEBUG] Write: len(encrypted) = %+v\n", len(msg))
	//_, err = p.Conn.Write(msg)
	////if err  != nil {
	//	log.Printf("[ERORR] failed to write: %s\n", err.Error())
	//	return 0, err
	//}
	//log.Printf("[DEBUG] Write: Write size: %d\n", l)
	//log.Printf("[DEBUG] Write: END------------\n")
	//return len(raw), nil
	return p.Conn.Write(raw)
}

func NewConn(c net.Conn) (net.Conn, error) {
	crypt := &HelloRecordCrypt{}
	overhead := MsgLenFieldSize
	helloConn := &conn{
		Conn:     c,
		crypt:    *crypt,
		overhead: overhead,
	}

	return helloConn, nil
}
