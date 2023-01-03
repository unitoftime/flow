package net

import (
	"fmt"
	"errors"
	// "math/rand"
	"time"
	"net"
	"net/url"
	"net/http"
	"crypto/tls"
	// "sync"
	// "sync/atomic"
	"context"

	"nhooyr.io/websocket"
)

// tomorrow: look at adding just a reconnect goroutine which runs until the socket is closed. I can't really remember why I didn't like having a reconnect loop but I dont think its too bad
// func newWebsocket(c *DialConfig) Socket {
// 	// TODO - centralize
// 	return &TransportSocket{
// 		encoder: encoder,
// 		recvBuf: make([]byte, MaxRecvMsgSize),
// 	}
// }

// Returns a connected socket or fails with an error
func dialWebsocket(c *DialConfig) (Transport, error) {
	ctx := context.Background()
	wsConn, err := dialWs(ctx, c.Url, c.TlsConfig)
	if err != nil {
		return nil, err
	}

	// Note: This connection is automagically framed by websockets
	conn := websocket.NetConn(ctx, wsConn, websocket.MessageBinary)

	return conn, nil
	// return newTransportSocket(conn, c.Serdes), nil
}

// --------------------------------------------------------------------------------
// - Listener
// --------------------------------------------------------------------------------

// TODO - (When I migrate to TCP) TCP will send 0 byte messages to indicate closes, websockets sends them without closing
type WebsocketListener struct {
	httpServer http.Server
	originPatterns []string
	addr net.Addr
	serdes Serdes
	pendingAccepts chan Socket // TODO - should this get buffered?
	pendingAcceptErrors chan error // TODO - should this get buffered?
}
func newWebsocketListener(c *ListenConfig) (*WebsocketListener, error) {
	u, err := url.Parse(c.Url)
	if err != nil {
		return nil, err
	}

	// TODO - is TCP always correct?
	listener, err := tls.Listen("tcp", u.Host, c.TlsConfig)
	if err != nil {
		panic(err)
	}

	wsl := &WebsocketListener{
		serdes: c.Serdes,
		addr: listener.Addr(),
		pendingAccepts: make(chan Socket),
		pendingAcceptErrors: make(chan error),
		originPatterns: c.OriginPatterns,
	}

	httpServer := c.HttpServer
	httpServer.Handler = wsl

	go func() {
		for {
			err := httpServer.Serve(listener)
			// TODO - what happens if this continually fails, how do we notify back?
			// ErrServerClosed is returned when shutdown or close is called
			fmt.Println("Serving Error:", err)

			if errors.Is(err, http.ErrServerClosed) {
				return // Just close if the server is shutdown or closed
			}

			time.Sleep(1 * time.Second)
		}
	}()

	return wsl, nil
}

func (l *WebsocketListener) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: l.originPatterns,
	})
	if err != nil {
		// Return as an accept error
		l.pendingAcceptErrors <- err
		return
	}

	// Build the socket and push to channel
	ctx := context.Background()
	conn := websocket.NetConn(ctx, c, websocket.MessageBinary)
	sock := newAcceptedSocket(conn, l.serdes)
	l.pendingAccepts <- sock
}

func (l *WebsocketListener) Accept() (Socket, error) {
	select{
	case sock := <-l.pendingAccepts:
		return sock, nil
	case err := <-l.pendingAcceptErrors:
		return nil, err
	}
}
func (l *WebsocketListener) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()
	return l.httpServer.Shutdown(ctx)
}
func (l *WebsocketListener) Addr() net.Addr {
	return l.addr
}
