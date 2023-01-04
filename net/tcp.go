package net

import (
	"net"
)

type tcpPipe struct {
	conn net.Conn
}
func (t *tcpPipe) Read(b []byte) (int, error) {
	return t.conn.Read(b)
}

func (t *tcpPipe) Write(b []byte) (int, error) {
	return t.conn.Write(b)
}

func (t *tcpPipe) Close() error {
	return t.conn.Close()
}

func dialTcpSocket(c *DialConfig) (Pipe, error) {
	conn, err := net.Dial("tcp", c.host)
	if err != nil {
		return nil, err
	}

	// Create a Framed connection and set it to our connection
	framedConn := NewFrameConn(conn)
	return &tcpPipe{framedConn}, nil
}

type TcpListener struct {
	listener net.Listener
	serdes Serdes
}
func newTcpListener(c *ListenConfig) (*TcpListener, error) {
	listener, err := net.Listen(c.scheme, c.host)
	if err != nil {
		return nil, err
	}
	sockListener := &TcpListener{
		listener: listener,
		serdes: c.Serdes,
	}
	return sockListener, nil
}

func (l *TcpListener) Accept() (Socket, error) {
	c, err := l.listener.Accept()
	if err != nil {
		return nil, err
	}

	pipe := &tcpPipe{NewFrameConn(c)}
	return newAcceptedSocket(pipe, l.serdes), nil
}
func (l *TcpListener) Close() error {
	return l.listener.Close()
}
func (l *TcpListener) Addr() net.Addr {
	return l.listener.Addr()
}

