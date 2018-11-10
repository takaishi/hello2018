package conn

import (
	"bytes"
	"fmt"
	"net"
	"reflect"
	"testing"
)

type testConn struct {
	net.Conn
	in  *bytes.Buffer
	out *bytes.Buffer
}

func (c *testConn) Read(b []byte) (n int, err error) {
	return c.in.Read(b)
}

func (c *testConn) Write(b []byte) (n int, err error) {
	return c.out.Write(b)
}

func (c *testConn) Close() error {
	return nil
}

func newTestRecordSecureConn(in, out *bytes.Buffer) *secureConn {
	tc := testConn{
		in:  in,
		out: out,
	}
	c, err := NewSecureConn(&tc)
	if err != nil {
		panic(fmt.Sprintf("Unexpected error creating test record connection: %v", err))
	}
	return c.(*secureConn)
}

func newConnPair() (client, server *secureConn) {
	clientBuf := new(bytes.Buffer)
	serverBuf := new(bytes.Buffer)
	clientConn := newTestRecordSecureConn(clientBuf, serverBuf)
	serverConn := newTestRecordSecureConn(serverBuf, clientBuf)
	return clientConn, serverConn
}

func TestPingPong(t *testing.T) {
	clientConn, serverConn := newConnPair()
	fmt.Printf("clientConn: %+v\n", clientConn)

	clientMsg := []byte("Client Message")
	if n, err := clientConn.Write(clientMsg); n != len(clientMsg) || err != nil {
		t.Fatalf("Client Write() = %v, %v; want %v, <nil>", n, err, len(clientMsg))
	}

	rcvClientMsg := make([]byte, len(clientMsg))
	if n, err := serverConn.Read(rcvClientMsg); n != len(rcvClientMsg) || err != nil {
		t.Fatalf("Server Read() = %v, %v; want %v, <nil>", n, err, len(rcvClientMsg))
	}

	if !reflect.DeepEqual(clientMsg, rcvClientMsg) {
		t.Fatalf("Client Write()/Server Read() = %v, want %v", rcvClientMsg, clientMsg)

	}

	serverMsg := []byte("Server Message")
	if n, err := serverConn.Write(serverMsg); n != len(serverMsg) || err != nil {
		t.Fatalf("Client Read() = %v, %v; want %v, <nil>", n, err, len(serverMsg))
	}
	rcvServerMsg := make([]byte, len(serverMsg))
	if n, err := clientConn.Read(rcvServerMsg); n != len(rcvServerMsg) || err != nil {
		t.Fatalf("Client Read() = %v, %v; want %v, <nil>", n, err, len(rcvServerMsg))
	}
	if !reflect.DeepEqual(serverMsg, rcvServerMsg) {
		t.Fatalf("Server Write()/Client Read() = %v, want %v", rcvServerMsg, serverMsg)
	}

}
