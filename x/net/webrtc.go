package net

import (
	"net"
	"github.com/unitoftime/rtcnet"
)

func dialWebRtc(c *DialConfig) (Pipe, error) {
	conn, err := rtcnet.Dial(c.host, c.TlsConfig)
	return conn, err
}

func newWebRtcListener(c *ListenConfig) (*rtcListener, error) {
	listener, err := rtcnet.NewListener(c.host, rtcnet.ListenConfig{
		TlsConfig: c.TlsConfig,
		OriginPatterns: c.OriginPatterns,
	})
	if err != nil {
		return nil, err
	}

	sockListener := &rtcListener{
		listener: listener,
		encoder: c.Encoder,
		decoder: c.Decoder,
	}
	return sockListener, nil
}

type rtcListener struct {
	listener net.Listener
	encoder Serdes
	decoder Serdes
}

func (l *rtcListener) Accept() (Socket, error) {
	c, err := l.listener.Accept()
	if err != nil {
		return nil, err
	}

	return newAcceptedSocket(c, l.encoder, l.decoder), nil
}
func (l *rtcListener) Close() error {
	return l.listener.Close()
}
func (l *rtcListener) Addr() net.Addr {
	return l.listener.Addr()
}
